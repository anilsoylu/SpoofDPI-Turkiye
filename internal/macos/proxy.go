package macos

import (
	"fmt"
	"os/exec"
	"strings"
)

// activeServices, etkin ağ servislerinin adlarını döndürür (devre dışı olanları atlar).
func activeServices() ([]string, error) {
	out, err := exec.Command("networksetup", "-listallnetworkservices").Output()
	if err != nil {
		return nil, fmt.Errorf("ağ servisleri listelenemedi: %w", err)
	}
	var services []string
	for i, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if i == 0 || line == "" {
			continue // ilk satır başlık
		}
		if strings.HasPrefix(line, "*") {
			continue // * = devre dışı servis
		}
		services = append(services, line)
	}
	return services, nil
}

// enablePAC, tüm etkin servislerde PAC URL'ini ayarlar ve açar.
func enablePAC(pacURL string) error {
	services, err := activeServices()
	if err != nil {
		return err
	}
	for _, s := range services {
		if out, err := exec.Command("networksetup", "-setautoproxyurl", s, pacURL).CombinedOutput(); err != nil {
			return fmt.Errorf("%s için PAC URL ayarlanamadı: %v: %s", s, err, out)
		}
		if out, err := exec.Command("networksetup", "-setautoproxystate", s, "on").CombinedOutput(); err != nil {
			return fmt.Errorf("%s için PAC durumu açılamadı: %v: %s", s, err, out)
		}
	}
	return nil
}

// disablePAC, tüm etkin servislerde PAC proxy'yi kapatır.
func disablePAC() error {
	services, err := activeServices()
	if err != nil {
		return err
	}
	for _, s := range services {
		// Hataları topla ama devam et; mümkün olduğunca çok servisi temizle.
		_ = exec.Command("networksetup", "-setautoproxystate", s, "off").Run()
	}
	return nil
}
