package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ZapretRepo, tpws motorunun derlendiği kaynaktır. install.sh içindeki
// ZAPRET_REPO ile aynı değeri taşır; biri değişirse diğeri de güncellenmeli.
const ZapretRepo = "https://github.com/bol-van/zapret"

// BuildFromSource, zapret'i sığ klonlayıp `make mac` ile derler ve çıkan tpws
// binary'sini BinPath()'e kurar (0755). Mevcut binary, kurulum BAŞARILI
// olana kadar değiştirilmez; derleme başarısızsa eski binary yerinde kalır.
// Yönetici yetkisi gerektirmez (kurulum ~ altına yapılır).
func BuildFromSource(repoURL string) error {
	for _, bin := range []string{"git", "make"} {
		if _, err := exec.LookPath(bin); err != nil {
			return fmt.Errorf("%s bulunamadı (Xcode Command Line Tools): %w", bin, err)
		}
	}
	if _, err := exec.LookPath("cc"); err != nil {
		if _, err := exec.LookPath("clang"); err != nil {
			return fmt.Errorf("C derleyicisi (cc/clang) bulunamadı (Xcode Command Line Tools): %w", err)
		}
	}

	tmp, err := os.MkdirTemp("", "zapret-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	src := filepath.Join(tmp, "zapret")
	clone := exec.Command("git", "clone", "--depth", "1", repoURL, src)
	if out, err := clone.CombinedOutput(); err != nil {
		return fmt.Errorf("zapret klonlanamadı: %w\n%s", err, out)
	}

	make := exec.Command("make", "mac")
	make.Dir = src
	if out, err := make.CombinedOutput(); err != nil {
		return fmt.Errorf("tpws derlenemedi: %w\n%s", err, out)
	}

	built := filepath.Join(src, "binaries", "my", "tpws")
	if fi, err := os.Stat(built); err != nil || fi.IsDir() {
		return fmt.Errorf("tpws derlenemedi: %s bulunamadı", built)
	}

	if err := os.MkdirAll(filepath.Dir(BinPath()), 0o755); err != nil {
		return err
	}
	data, err := os.ReadFile(built)
	if err != nil {
		return err
	}
	return os.WriteFile(BinPath(), data, 0o755)
}
