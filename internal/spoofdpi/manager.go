// Package spoofdpi, resmî xvzc/spoofdpi binary'sini indirir, doğrular ve yönetir.
package spoofdpi

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	owner = "xvzc"
	repo  = "spoofdpi"
)

// ghAsset, GitHub releases API asset alt kümesi.
type ghAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Digest      string `json:"digest"` // "sha256:HEX"
}

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

// assetArch, Go runtime.GOARCH'ı release asset arch adına çevirir.
func assetArch() (string, error) {
	switch runtime.GOARCH {
	case "arm64":
		return "arm64", nil
	case "amd64":
		return "x86_64", nil
	default:
		return "", fmt.Errorf("desteklenmeyen mimari: %s", runtime.GOARCH)
	}
}

// BinPath, yönetilen binary'nin tam yolunu döndürür: ~/.spoofdpi-tr/bin/spoofdpi
func BinPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".spoofdpi-tr", "bin", "spoofdpi"), nil
}

var httpClient = &http.Client{Timeout: 60 * time.Second}

// fetchRelease, verilen tag için release'i getirir. tag boşsa "latest" kullanılır.
func fetchRelease(tag string) (*ghRelease, error) {
	var url string
	if tag == "" {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	} else {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("release sorgulanamadı: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("release API durumu %d", resp.StatusCode)
	}
	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	return &rel, nil
}

// pickAsset, darwin/<arch> tar.gz asset'ini ve sha256 digest'ini bulur.
func pickAsset(rel *ghRelease) (ghAsset, error) {
	arch, err := assetArch()
	if err != nil {
		return ghAsset{}, err
	}
	want := fmt.Sprintf("darwin_%s.tar.gz", arch)
	for _, a := range rel.Assets {
		// .sbom.json gibi yan dosyaları ele (yalnızca tam tar.gz adı).
		if strings.HasSuffix(a.Name, want) && strings.HasSuffix(a.Name, ".tar.gz") {
			return a, nil
		}
	}
	return ghAsset{}, fmt.Errorf("uygun asset bulunamadı (%s)", want)
}

// verifyChecksum, data'nın sha256'sını expected (hex string) ile karşılaştırır.
func verifyChecksum(data []byte, expected string) error {
	sum := sha256.Sum256(data)
	got := hex.EncodeToString(sum[:])
	if !strings.EqualFold(got, expected) {
		return fmt.Errorf("checksum uyuşmazlığı: beklenen %s, hesaplanan %s", expected, got)
	}
	return nil
}

// Install, verilen sürümü (tag boşsa latest) indirir, sha256 doğrular,
// tarball'dan spoofdpi binary'sini çıkarır ve kurar. Kurulan sürümü döndürür.
func Install(tag string) (installedVersion string, err error) {
	rel, err := fetchRelease(tag)
	if err != nil {
		return "", err
	}
	asset, err := pickAsset(rel)
	if err != nil {
		return "", err
	}
	expected, ok := strings.CutPrefix(asset.Digest, "sha256:")
	if !ok || expected == "" {
		return "", fmt.Errorf("asset için sha256 digest yok; güvenlik nedeniyle iptal")
	}

	// İndir.
	resp, err := httpClient.Get(asset.DownloadURL)
	if err != nil {
		return "", fmt.Errorf("indirme hatası: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("indirme durumu %d", resp.StatusCode)
	}

	// Tüm tarball'ı belleğe al (birkaç MB) ve checksum doğrula.
	const maxDownload = 64 << 20 // 64 MB
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxDownload+1))
	if err != nil {
		return "", err
	}
	if len(data) > maxDownload {
		return "", fmt.Errorf("indirme çok büyük (>64MB), iptal edildi")
	}
	if err := verifyChecksum(data, expected); err != nil {
		return "", err
	}

	// spoofdpi binary'sini tarball'dan çıkar.
	bin, err := extractBinary(data)
	if err != nil {
		return "", err
	}

	// Hedefe yaz (önce .tmp, sonra atomik rename).
	dest, err := BinPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", err
	}
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, bin, 0o755); err != nil {
		return "", err
	}
	// Rename başarılı olursa Remove no-op olur; başarısız olursa .tmp dosyasını temizler.
	defer os.Remove(tmp)
	if err := os.Rename(tmp, dest); err != nil {
		return "", err
	}
	return strings.TrimPrefix(rel.TagName, "v"), nil
}

// extractBinary, gzip+tar içinden "spoofdpi" adlı dosyayı bayt olarak döndürür.
func extractBinary(targz []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(targz))
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// Tarball kökünde "spoofdpi" (alt dizin olabilir → base ad kontrol).
		if hdr.Typeflag == tar.TypeReg && filepath.Base(hdr.Name) == "spoofdpi" {
			return io.ReadAll(tr)
		}
	}
	return nil, fmt.Errorf("tarball içinde spoofdpi binary bulunamadı")
}

// IsInstalled, yönetilen binary diskte var mı kontrol eder.
func IsInstalled() bool {
	p, err := BinPath()
	if err != nil {
		return false
	}
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

// LatestVersion, upstream'deki en son sürümü (v önekiz) döndürür.
func LatestVersion() (string, error) {
	rel, err := fetchRelease("")
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(rel.TagName, "v"), nil
}
