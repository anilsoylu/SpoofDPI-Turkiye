package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config, spoofdpi-tr'nin kalıcı kullanıcı yapılandırmasıdır.
// ~/.spoofdpi-tr/config.json içinde saklanır.
type Config struct {
	// SpoofDPIVersion, indirilen resmî spoofdpi binary sürümü (örn. "1.5.3").
	SpoofDPIVersion string `json:"spoofdpi_version"`
	// Port, spoofdpi'nin dinleyeceği yerel proxy portu (varsayılan 8080).
	Port int `json:"port"`
	// Domains, PAC tarafından proxy'ye yönlendirilecek alan adları (örn. "discord.com").
	Domains []string `json:"domains"`
	// DNS ayarları (Türkiye'de DNS zehirlemesi yaygın → DoH varsayılan açık).
	EnableDoH bool   `json:"enable_doh"`
	DNSAddr   string `json:"dns_addr"`
}

// Default, ilk kurulumda kullanılan makul Türkiye varsayılanlarını döndürür.
// Domains kasıtlı boştur; Plan 003 küratörlü blocklist'i sağlayacak.
func Default() *Config {
	return &Config{
		Port:      8080,
		Domains:   []string{},
		EnableDoH: true,
		DNSAddr:   "1.1.1.1",
	}
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
