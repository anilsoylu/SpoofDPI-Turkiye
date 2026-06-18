package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Config, spoofdpi-tr'nin kalıcı kullanıcı yapılandırmasıdır.
// ~/.spoofdpi-tr/config.json içinde saklanır.
type Config struct {
	// Port, tpws motorunun dinleyeceği yerel TPROXY/redirect portu (varsayılan 988).
	// PF, 443 trafiğini bu porta yönlendirir.
	Port int `json:"port"`
	// Domains, tpws hostlist'ine yazılan, DPI-bypass uygulanacak alan adları
	// (örn. "discord.com"). tpws "subdomain auto apply" yapar; düz domain yeterli.
	Domains []string `json:"domains"`
	// EnableDoH / DNSAddr, eski PAC+spoofdpi mimarisinden kalmıştır. tpws DNS
	// çözümlemesi YAPMAZ; bu alanlar artık kullanılmaz ama geriye dönük JSON
	// uyumluluğu için saklanır.
	EnableDoH bool   `json:"enable_doh,omitempty"`
	DNSAddr   string `json:"dns_addr,omitempty"`
}

// Default, ilk kurulumda kullanılan makul Türkiye varsayılanlarını döndürür.
// Domains kasıtlı boştur; blocklist kategorileri sağlayacak.
func Default() *Config {
	return &Config{
		// 988 tpws redirect portu — yaygın geliştirici portlarıyla (8080/8081/19000
		// Expo/Metro) çakışmaz. Kullanıcı 'port' komutuyla değiştirebilir.
		Port:    988,
		Domains: []string{},
	}
}

// NormalizeDomains, ham domain listesini standart biçime getirir:
// trim, lower, "*." / "." önek strip, boşları at, tekilleştir, sırala.
func NormalizeDomains(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, d := range in {
		d = strings.ToLower(strings.TrimSpace(d))
		d = strings.TrimPrefix(d, "*.")
		d = strings.TrimPrefix(d, ".")
		if d == "" || seen[d] {
			continue
		}
		seen[d] = true
		out = append(out, d)
	}
	sort.Strings(out)
	return out
}

// Dir, yapılandırma dizinini döndürür: ~/.spoofdpi-tr
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".spoofdpi-tr"), nil
}

// Path, config.json'un tam yolunu döndürür.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load, config.json'u okur. Dosya yoksa Default() döndürür (hata değil).
func Load() (*Config, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return Default(), nil
	}
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// ValidatePort, portun geçerli aralıkta olduğunu denetler.
func ValidatePort(p int) error {
	if p < 1 || p > 65535 {
		return fmt.Errorf("port 1-65535 aralığında olmalı, verilen: %d", p)
	}
	return nil
}

// Save, yapılandırmayı diske yazar; dizini gerekirse oluşturur (0700).
func (c *Config) Save() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	p, err := Path()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}
