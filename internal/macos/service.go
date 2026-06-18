package macos

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/config"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/pac"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/spoofdpi"
)

// pacPath, üretilen PAC dosyasının yolunu döndürür.
func pacPath() (string, error) {
	dir, err := config.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "proxy.pac"), nil
}

// spoofdpiArgs, config'ten spoofdpi komut argümanlarını kurar.
func spoofdpiArgs(c *config.Config) []string {
	args := []string{
		"-port", strconv.Itoa(c.Port),
		"-system-proxy=false",
		"-silent",
	}
	if c.EnableDoH {
		args = append(args, "-enable-doh")
	}
	if c.DNSAddr != "" {
		args = append(args, "-dns-addr", c.DNSAddr)
	}
	return args
}

// On, PAC'i yazar, LaunchAgent'ı yükler ve sistem proxy'sini PAC'e yönlendirir.
func On(c *config.Config) error {
	if !spoofdpi.IsInstalled() {
		return fmt.Errorf("spoofdpi binary kurulu değil; önce 'spoofdpi-tr update' çalıştırın")
	}
	// 1. PAC yaz.
	pp, err := pacPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(pp), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(pp, []byte(pac.Generate(c.Port, c.Domains)), 0o644); err != nil {
		return err
	}
	// 2. LaunchAgent yükle.
	bin, err := spoofdpi.BinPath()
	if err != nil {
		return err
	}
	if err := writeAndLoad(bin, spoofdpiArgs(c)); err != nil {
		return err
	}
	// 3. Sistem proxy'yi PAC'e yönlendir.
	return enablePAC("file://" + pp)
}

// Off, sistem proxy'yi kapatır ve LaunchAgent'ı çıkarır.
func Off() error {
	if err := disablePAC(); err != nil {
		return err
	}
	return unload()
}

// Status, mevcut durumu temsil eder.
type Status struct {
	ServiceLoaded   bool
	BinaryInstalled bool
}

// CurrentStatus, servis ve binary durumunu döndürür.
func CurrentStatus() Status {
	return Status{
		ServiceLoaded:   isLoaded(),
		BinaryInstalled: spoofdpi.IsInstalled(),
	}
}

// RefreshPAC, config değiştiğinde (add/remove) PAC'i yeniden yazar.
// Servis çalışıyorsa proxy zaten PAC'i otomatik okur (dosya değişir).
func RefreshPAC(c *config.Config) error {
	pp, err := pacPath()
	if err != nil {
		return err
	}
	return os.WriteFile(pp, []byte(pac.Generate(c.Port, c.Domains)), 0o644)
}
