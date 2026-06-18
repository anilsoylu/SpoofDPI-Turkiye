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

// autoProxyURLMatches, networksetup -getautoproxyurl çıktısının pacPath'i
// içerip içermediğini denetler. Saf fonksiyon — test edilebilir.
func autoProxyURLMatches(getautoproxyOutput, pacPath string) bool {
	for _, line := range strings.Split(getautoproxyOutput, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "URL:") {
			return strings.Contains(line, pacPath)
		}
	}
	return false
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

// disablePAC, yalnızca BİZİM PAC'imizi kullanan servislerde PAC proxy'yi kapatır.
// Diğer servislerin mevcut PAC yapılandırmasına dokunulmaz.
func disablePAC(pacPath string) error {
	services, err := activeServices()
	if err != nil {
		return err
	}
	for _, s := range services {
		// Önce bu servisin hangi PAC URL'ini kullandığını sorgula.
		out, err := exec.Command("networksetup", "-getautoproxyurl", s).Output()
		if err != nil {
			// Sorgu başarısız olursa bu servise dokunma.
			continue
		}
		// Yalnızca bizim PAC dosyamıza işaret ediyorsa kapat.
		if autoProxyURLMatches(string(out), pacPath) {
			_ = exec.Command("networksetup", "-setautoproxystate", s, "off").Run()
		}
	}
	return nil
}
