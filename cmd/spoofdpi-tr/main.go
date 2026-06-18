package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/blocklist"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/cli"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/config"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/macos"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/pac"
	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/spoofdpi"
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
	case "update":
		err = runUpdate(args)
	case "uninstall":
		err = runUninstall(args)
	case "port":
		err = runPort(args)
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
	fmt.Print(`spoofdpi-tr — Türkiye için DPI bypass yöneticisi

Kullanım:
  spoofdpi-tr <komut> [argümanlar]

Komutlar:
  install        İnteraktif kurulum (port, servisler, otomatik başlatma)
  on             Bypass'ı başlat (PAC'i devreye al)
  off            Bypass'ı durdur (proxy → DIRECT)
  status         Servis ve yapılandırma durumunu göster
  add <domain>   Bypass listesine domain ekle
  remove <dom>   Bypass listesinden domain çıkar
  list           Bypass edilen domainleri göster
  update         Resmî spoofdpi binary'sini güncelle
  uninstall      Tüm yapılandırma ve servisi kaldır [-y]
  port <numara>  Proxy portunu değiştir (1-65535)
  version        Sürümü göster
`)
}

func runUpdate(args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	fmt.Println("En son spoofdpi sürümü kontrol ediliyor...")

	// Madde 9: Zaten güncel mi kontrol et.
	latest, err := spoofdpi.LatestVersion()
	if err != nil {
		return err
	}
	if spoofdpi.IsInstalled() && latest == cfg.SpoofDPIVersion {
		fmt.Printf("✓ Zaten güncel (v%s)\n", latest)
		return nil
	}

	ver, err := spoofdpi.Install("") // latest
	if err != nil {
		return err
	}
	cfg.SpoofDPIVersion = ver
	if err := cfg.Save(); err != nil {
		return err
	}
	bp, _ := spoofdpi.BinPath()
	fmt.Printf("✓ spoofdpi %s kuruldu: %s\n", ver, bp)
	return nil
}

func runOn() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := macos.On(cfg); err != nil {
		return err
	}
	// Madde 10: 0 domain uyarısı.
	if len(cfg.Domains) == 0 {
		fmt.Printf("⚠ Hiç domain seçili değil — bypass hiçbir trafiği etkilemez. " +
			"'spoofdpi-tr add <domain>' veya 'spoofdpi-tr install' ile ekleyin.\n")
	} else {
		fmt.Printf("✓ Açık — port %d, %d domain proxy'ye yönlendiriliyor (kalan trafik DIRECT)\n",
			cfg.Port, len(cfg.Domains))
	}
	return nil
}

func runOff() error {
	if err := macos.Off(); err != nil {
		return err
	}
	fmt.Println("✓ Kapalı — sistem proxy DIRECT'e döndürüldü")
	return nil
}

func runStatus() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	st := macos.CurrentStatus()
	state := "durdu"
	if st.ServiceLoaded {
		state = "çalışıyor"
	}
	fmt.Printf("Servis : %s\n", state)
	fmt.Printf("Binary : %v\n", st.BinaryInstalled)
	fmt.Printf("Port   : %d\n", cfg.Port)
	fmt.Printf("Domain : %d\n", len(cfg.Domains))
	// Madde 9: SpoofDPIVersion göster.
	ver := cfg.SpoofDPIVersion
	if ver == "" {
		ver = "kurulu değil"
	}
	fmt.Printf("Sürüm  : %s\n", ver)
	return nil
}

func runInstall(args []string) error {
	p := cli.New(os.Stdin, os.Stdout)

	// Madde 1: Mevcut config'i koru; yoksa Default() döndürür.
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("SpoofDPI Türkiye kurulumu")
	fmt.Println()

	// 1. Resmî binary'yi kur.
	if !spoofdpi.IsInstalled() {
		fmt.Println("spoofdpi indiriliyor ve doğrulanıyor...")
		ver, err := spoofdpi.Install("")
		if err != nil {
			return err
		}
		cfg.SpoofDPIVersion = ver
		fmt.Printf("✓ spoofdpi %s kuruldu\n\n", ver)
	}

	// 2. Port sor (varsayılan mevcut config'ten gelir — yeniden kurulumda korunur).
	// Madde 4: Geçersiz port girilirse tekrar sor.
	for {
		port, err := p.AskInt("Dinleme portu", cfg.Port)
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
	// Madde 4: <1024 ayrıcalıklı port uyarısı (engelleme değil).
	if cfg.Port < 1024 {
		fmt.Fprintf(os.Stderr, "Uyarı: %d ayrıcalıklı port (<1024); macOS'ta root yetkisi gerekebilir.\n", cfg.Port)
	}

	// 3. Kategori seç. Her kategori için evet/hayır.
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

	// 4. Özel domain ekleme.
	if add, _ := p.AskYesNo("Özel domain eklemek ister misiniz?", false); add {
		for {
			line, _ := p.AskLine("  Domain (boş bırak=bitir)")
			if line == "" {
				break
			}
			domains = append(domains, line)
		}
	}
	cfg.Domains = pac.Normalize(domains)

	// 5. Kaydet.
	if err := cfg.Save(); err != nil {
		return err
	}

	// 6. Otomatik başlat.
	start, err := p.AskYesNo("\nŞimdi başlatılsın ve açılışta otomatik çalışsın mı?", true)
	if err != nil {
		return err
	}
	if start {
		if err := macos.On(cfg); err != nil {
			return err
		}
		fmt.Printf("\n✓ Kuruldu ve başlatıldı — port %d, %d domain proxy'den, kalan internet DIRECT\n",
			cfg.Port, len(cfg.Domains))
	} else {
		fmt.Println("\n✓ Kuruldu. Başlatmak için: spoofdpi-tr on")
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
	cfg.Domains = pac.Normalize(append(cfg.Domains, args...))
	if err := cfg.Save(); err != nil {
		return err
	}
	// Servis çalışıyorsa PAC'i tazele.
	if macos.CurrentStatus().ServiceLoaded {
		if err := macos.RefreshPAC(cfg); err != nil {
			return err
		}
	}
	fmt.Printf("✓ Eklendi. Toplam %d domain.\n", len(cfg.Domains))
	return nil
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
	for _, a := range args {
		rm[strings.ToLower(strings.TrimSpace(a))] = true
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
	if macos.CurrentStatus().ServiceLoaded {
		if err := macos.RefreshPAC(cfg); err != nil {
			return err
		}
	}
	fmt.Printf("✓ Çıkarıldı. Toplam %d domain.\n", len(cfg.Domains))
	return nil
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
	// Port LaunchAgent plist argümanında geçtiğinden servis yeniden başlatılmalı.
	if macos.CurrentStatus().ServiceLoaded {
		if err := macos.Off(); err != nil {
			return fmt.Errorf("servis durdurulurken hata: %w", err)
		}
		if err := macos.On(cfg); err != nil {
			return fmt.Errorf("servis yeniden başlatılırken hata: %w", err)
		}
		fmt.Printf("✓ Port %d olarak ayarlandı — servis yeni portla yeniden başlatıldı\n", n)
	} else {
		fmt.Printf("✓ Port %d olarak ayarlandı\n", n)
	}
	return nil
}

// hasFlag, args içinde -y veya --yes bayrağı var mı kontrol eder.
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

func runUninstall(args []string) error {
	// Madde 11: TTY kontrolü ve onay sorusu.
	fi, _ := os.Stdin.Stat()
	isTTY := (fi.Mode() & os.ModeCharDevice) != 0
	forceYes := hasFlag(args, "-y", "--yes")

	if isTTY && !forceYes {
		p := cli.New(os.Stdin, os.Stdout)
		ok, err := p.AskYesNo("Tüm yapılandırma ve servis silinsin mi?", false)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("İptal edildi.")
			return nil
		}
	}

	// 1. Servisi durdur ve proxy'yi geri al (kullanıcı internetsiz kalmasın).
	if err := macos.Off(); err != nil {
		fmt.Fprintf(os.Stderr, "uyarı: servis durdurulurken hata: %v\n", err)
	}
	// 2. Yapılandırma dizinini sil.
	dir, err := config.Dir()
	if err != nil {
		return err
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	// 3. CLI binary'sinin kendisini sil (best-effort).
	if exe, err := os.Executable(); err == nil {
		if err := os.Remove(exe); err != nil {
			fmt.Fprintf(os.Stderr, "uyarı: binary silinemedi (%s): %v\n", exe, err)
		}
	}
	fmt.Println("✓ Tamamen kaldırıldı — servis durduruldu, proxy DIRECT'e döndü, yapılandırma ve uygulama silindi.")
	return nil
}
