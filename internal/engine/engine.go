// Package engine, tpws DPI-bypass motoru için saf yardımcılar sağlar.
// tpws binary'si ~/.spoofdpi-tr/bin/tpws altında derli/hazır bulunur (PoC'de
// doğrulanmıştır). Bu paket yalnızca yol/varlık ve hostlist helpers yapar;
// kök (root) işlemleri helper script tarafından yürütülür.
package engine

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/config"
)

// BinPath, tpws binary'sinin tam yolunu döndürür: ~/.spoofdpi-tr/bin/tpws
func BinPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Home çözümlenemezse göreli olmayan en makul yedek; pratikte oluşmaz.
		return filepath.Join(".spoofdpi-tr", "bin", "tpws")
	}
	return filepath.Join(home, ".spoofdpi-tr", "bin", "tpws")
}

// IsInstalled, tpws binary'sinin mevcut olup olmadığını döndürür.
func IsInstalled() bool {
	info, err := os.Stat(BinPath())
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// HostlistPath, tpws'in okuyacağı hostlist dosyasının tam yolunu döndürür:
// ~/.spoofdpi-tr/hostlist.txt (mutlak; root tarafından okunabilir).
func HostlistPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Home çözümlenemezse göreli olmayan en makul yedek; pratikte oluşmaz.
		return filepath.Join(".spoofdpi-tr", "hostlist.txt")
	}
	return filepath.Join(home, ".spoofdpi-tr", "hostlist.txt")
}

// WriteHostlist, verilen domainleri (NormalizeDomains uygulanmış, her satır bir
// domain) HostlistPath()'e yazar. Dizini gerekirse oluşturur. tpws "subdomain
// auto apply" yaptığından düz domain yeterlidir.
func WriteHostlist(domains []string) error {
	path := HostlistPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	norm := config.NormalizeDomains(domains)
	var sb strings.Builder
	for _, d := range norm {
		sb.WriteString(d)
		sb.WriteString("\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
