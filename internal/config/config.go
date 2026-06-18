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
//
// GÜVENLİK: Yalnızca GEÇERLİ domainler (bkz. ValidDomain) korunur. Geçersiz
// veya enjeksiyon karakterli girdiler (`evil.com;rm`, `$(...)`, boşluk, yeni
// satır vb.) SESSİZCE ATILIR. Bu, root yetkili helper/PF/hostlist akışlarına
// kontrolsüz kullanıcı metninin sızmasını önler (defense-in-depth) ve
// hostlist dosyasının satır bütünlüğünü korur.
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
		if !ValidDomain(d) {
			continue
		}
		seen[d] = true
		out = append(out, d)
	}
	sort.Strings(out)
	return out
}

// ValidDomain, normalize edilmiş bir domain'in güvenli ve sözdizimsel olarak
// geçerli bir host adı olup olmadığını döndürür. Kabul kriterleri:
//   - yalnızca [a-z0-9.-] karakterleri (küçük harf; çağıran NormalizeDomains
//     zaten ToLower uygular),
//   - en az bir nokta (TLD ayıracı),
//   - başında/sonunda nokta veya tire yok, ardışık nokta yok,
//   - her etiket 1-63, toplam uzunluk <= 253.
//
// Bu, kabuk metakarakterleri (`;`, `|`, `$`, “ ` “, `(`, `)`, boşluk, yeni
// satır) ve diğer enjeksiyon yüklerini reddeder.
func ValidDomain(d string) bool {
	if len(d) == 0 || len(d) > 253 {
		return false
	}
	if !strings.Contains(d, ".") {
		return false
	}
	if d[0] == '.' || d[0] == '-' || d[len(d)-1] == '.' || d[len(d)-1] == '-' {
		return false
	}
	labelLen := 0
	prevDot := false
	for i := 0; i < len(d); i++ {
		c := d[i]
		switch {
		case c == '.':
			if prevDot || labelLen == 0 {
				return false // ardışık veya boş etiket
			}
			prevDot = true
			labelLen = 0
			continue
		case c >= 'a' && c <= 'z':
		case c >= '0' && c <= '9':
		case c == '-':
			// Etiket başında/sonunda tire olamaz; nokta sınırlarını kontrol et.
			if labelLen == 0 || (i+1 < len(d) && d[i+1] == '.') {
				return false
			}
		default:
			return false // izin verilmeyen karakter (enjeksiyon savunması)
		}
		prevDot = false
		labelLen++
		if labelLen > 63 {
			return false
		}
	}
	return true
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
