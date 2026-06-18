package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/blocklist"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/cli"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/config"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/engine"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/macos"
)

// version, derleme sırasında ldflags ile gömülür; varsayılan "dev".
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	var err error
	switch cmd {
	case "version", "-v", "--version":
		fmt.Printf("spoofdpi-tr %s\n", version)
	case "install":
		err = runInstall(args)
	case "on", "start":
		err = runOn()
	case "off", "stop":
		err = runOff()
	case "status":
		err = runStatus()
	case "add":
		err = runAdd(args)
	case "remove":
		err = runRemove(args)
	case "list":
		err = runList()
	case "set":
		err = runSet(args)
	case "port":
		err = runPort(args)
	case "uninstall":
		err = runUninstall(args)
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "bilinmeyen komut: %s\n\n", cmd)
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "hata: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Print(`spoofdpi-tr — Türkiye için tpws+PF transparan DPI bypass yöneticisi

Kullanım:
  spoofdpi-tr <komut> [argümanlar]

Komutlar:
  install         İnteraktif kurulum (port, servisler, domainler) — bir kez admin onayı
  on              Bypass'ı başlat (parolasız)
  off             Bypass'ı durdur (parolasız)
  status          Servis ve yapılandırma durumunu göster
  add <domain>    Bypass listesine domain ekle
  remove <dom>    Bypass listesinden domain çıkar
  set <domain...> Bypass listesini verilen domainlerle tamamen değiştir
  list            Bypass edilen domainleri göster
  port <numara>   tpws redirect portunu değiştir (1-65535)
  uninstall       Tüm yapılandırma ve sistem dosyalarını kaldır [-y]
  version         Sürümü göster
`)
}

// runInstall, interaktif kurulumu yürütür. Tek bir osascript admin diyaloğuyla
// tüm sistem dosyalarını yerleştirir; sonrası parolasızdır.
func runInstall(args []string) error {
	_ = args
	p := cli.New(os.Stdin, os.Stdout)

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("SpoofDPI Türkiye kurulumu (tpws + PF)")
	fmt.Println()

	// tpws motoru hazır mı?
	if !engine.IsInstalled() {
		return fmt.Errorf("tpws motoru bulunamadı: %s — önce tpws'i bu yola yerleştirin", engine.BinPath())
	}

	// 1. Port sor (geçersizse tekrar sor).
	for {
		port, err := p.AskInt("tpws redirect portu", cfg.Port)
		if err != nil {
			return err
		}
		if err := config.ValidatePort(port); err != nil {
			fmt.Fprintf(os.Stderr, "Hata: %v\n", err)
			continue
		}
		cfg.Port = port
		break
	}

	// 2. Kategori seç.
	fmt.Println("\nHangi servisler bypass edilsin?")
	var domains []string
	for _, c := range blocklist.Categories() {
		yes, err := p.AskYesNo("  "+c.Title, true)
		if err != nil {
			return err
		}
		if yes {
			domains = append(domains, c.Domains...)
		}
	}

	// 3. Özel domain.
	if add, _ := p.AskYesNo("Özel domain eklemek ister misiniz?", false); add {
		for {
			line, _ := p.AskLine("  Domain (boş bırak=bitir)")
			if line == "" {
				break
			}
			domains = append(domains, line)
		}
	}
	cfg.Domains = config.NormalizeDomains(domains)

	// 4. Kaydet.
	if err := cfg.Save(); err != nil {
		return err
	}

	// 5. ROOT KURULUM — tek admin diyaloğu.
	if err := installSystemFiles(cfg); err != nil {
		return err
	}

	// 6. Başlat (parolasız).
	if err := macos.On(cfg.Port, cfg.Domains); err != nil {
		return fmt.Errorf("başlatma hatası: %w", err)
	}

	if len(cfg.Domains) == 0 {
		fmt.Printf("\n✓ Kuruldu ve başlatıldı. ⚠ Hiç domain seçili değil — 'spoofdpi-tr add <domain>' ekleyin.\n")
	} else {
		fmt.Printf("\n✓ Kuruldu ve başlatıldı — port %d, %d domain için DPI bypass aktif.\n", cfg.Port, len(cfg.Domains))
	}
	fmt.Println("Bundan sonra 'on'/'off'/'status' parolasız çalışır.")
	return nil
}

// installSystemFiles, yamalanmış pf.conf, plist, helper ve sudoers'ı tek admin
// diyaloğuyla yerleştirir.
func installSystemFiles(cfg *config.Config) error {
	// Mevcut pf.conf'u oku ve yama hazırla.
	orig, err := os.ReadFile(macos.PFConfPath)
	if err != nil {
		return fmt.Errorf("%s okunamadı: %w", macos.PFConfPath, err)
	}
	patched := macos.PFConfWithAnchor(string(orig))

	user, err := macos.CurrentUser()
	if err != nil {
		return err
	}

	tpwsBin := engine.BinPath()
	tpwsArgs := engine.Args(cfg.Port)

	script := macos.InstallScript(patched, tpwsBin, tpwsArgs, user)
	fmt.Println("\nSistem dosyaları yükleniyor — yönetici parolası istenebilir (tek seferlik)...")
	return macos.RunAdmin(script, "SpoofDPI Türkiye kurulum izni")
}

func runOn() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if !installedSystemFiles() {
		return fmt.Errorf("kurulu değil — önce 'spoofdpi-tr install' çalıştırın")
	}
	if err := macos.On(cfg.Port, cfg.Domains); err != nil {
		return err
	}
	if len(cfg.Domains) == 0 {
		fmt.Printf("⚠ Hiç domain seçili değil — bypass hiçbir trafiği etkilemez.\n")
	} else {
		fmt.Printf("✓ Açık — port %d, %d domain için DPI bypass aktif.\n", cfg.Port, len(cfg.Domains))
	}
	return nil
}

func runOff() error {
	if !installedSystemFiles() {
		return fmt.Errorf("kurulu değil — önce 'spoofdpi-tr install' çalıştırın")
	}
	if err := macos.Off(); err != nil {
		return err
	}
	fmt.Println("✓ Kapalı — 443 trafiği doğrudan akıyor.")
	return nil
}

func runStatus() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	fmt.Printf("Port   : %d\n", cfg.Port)
	fmt.Printf("Domain : %d\n", len(cfg.Domains))
	fmt.Printf("tpws   : %v\n", engine.IsInstalled())
	if installedSystemFiles() {
		out, _ := macos.Status()
		fmt.Print(strings.TrimRight(out, "\n") + "\n")
	} else {
		fmt.Println("Kurulum: yok (install gerekli)")
	}
	return nil
}

func runAdd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("kullanım: spoofdpi-tr add <domain>")
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	cfg.Domains = config.NormalizeDomains(append(cfg.Domains, args...))
	if err := cfg.Save(); err != nil {
		return err
	}
	fmt.Printf("✓ Eklendi. Toplam %d domain.\n", len(cfg.Domains))
	return refreshIfRunning(cfg)
}

func runRemove(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("kullanım: spoofdpi-tr remove <domain>")
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	rm := map[string]bool{}
	for _, a := range config.NormalizeDomains(args) {
		rm[a] = true
	}
	var kept []string
	for _, d := range cfg.Domains {
		if !rm[d] {
			kept = append(kept, d)
		}
	}
	cfg.Domains = kept
	if err := cfg.Save(); err != nil {
		return err
	}
	fmt.Printf("✓ Çıkarıldı. Toplam %d domain.\n", len(cfg.Domains))
	return refreshIfRunning(cfg)
}

func runSet(args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	cfg.Domains = config.NormalizeDomains(args)
	if err := cfg.Save(); err != nil {
		return err
	}
	fmt.Printf("✓ %d domain ayarlandı.\n", len(cfg.Domains))
	return refreshIfRunning(cfg)
}

func runList() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if len(cfg.Domains) == 0 {
		fmt.Println("(domain yok)")
		return nil
	}
	for _, d := range cfg.Domains {
		fmt.Println(d)
	}
	return nil
}

func runPort(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("kullanım: spoofdpi-tr port <1-65535>")
	}
	n := 0
	if _, err := fmt.Sscanf(args[0], "%d", &n); err != nil || args[0] == "" {
		return fmt.Errorf("kullanım: spoofdpi-tr port <1-65535>")
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := config.ValidatePort(n); err != nil {
		return err
	}
	cfg.Port = n
	if err := cfg.Save(); err != nil {
		return err
	}
	fmt.Printf("✓ Port %d olarak ayarlandı.\n", n)
	return refreshIfRunning(cfg)
}

func runUninstall(args []string) error {
	fi, _ := os.Stdin.Stat()
	isTTY := (fi.Mode() & os.ModeCharDevice) != 0
	forceYes := hasFlag(args, "-y", "--yes")

	if isTTY && !forceYes {
		p := cli.New(os.Stdin, os.Stdout)
		ok, err := p.AskYesNo("Tüm yapılandırma ve sistem dosyaları silinsin mi?", false)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("İptal edildi.")
			return nil
		}
	}

	// pf.conf'tan anchor'ı çıkar (mevcut içerikten).
	if data, err := os.ReadFile(macos.PFConfPath); err == nil {
		unpatched := macos.PFConfWithoutAnchor(string(data))
		script := macos.UninstallScript(unpatched)
		fmt.Println("Sistem dosyaları kaldırılıyor — yönetici parolası istenebilir...")
		if err := macos.RunAdmin(script, "SpoofDPI Türkiye kaldırma izni"); err != nil {
			fmt.Fprintf(os.Stderr, "uyarı: sistem temizliği hatası: %v\n", err)
		}
	}

	// Kullanıcı config dizinini sil (tpws binary'si dahil).
	if dir, err := config.Dir(); err == nil {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Fprintf(os.Stderr, "uyarı: config silinemedi: %v\n", err)
		}
	}

	fmt.Println("✓ Tamamen kaldırıldı.")
	return nil
}

// refreshIfRunning, servis çalışıyorsa yeni hostlist/port ile yeniden başlatır.
func refreshIfRunning(cfg *config.Config) error {
	if !installedSystemFiles() || !macos.IsRunning() {
		return nil
	}
	if err := macos.On(cfg.Port, cfg.Domains); err != nil {
		return fmt.Errorf("servis tazelenirken hata: %w", err)
	}
	fmt.Println("✓ Servis yeni ayarlarla yeniden yüklendi.")
	return nil
}

// installedSystemFiles, helper'ın yerinde olup olmadığını döndürür (kurulum işareti).
func installedSystemFiles() bool {
	info, err := os.Stat(macos.HelperPath)
	return err == nil && !info.IsDir()
}

// hasFlag, args içinde verilen bayraklardan biri var mı kontrol eder.
func hasFlag(args []string, flags ...string) bool {
	for _, a := range args {
		for _, f := range flags {
			if a == f {
				return true
			}
		}
	}
	return false
}
