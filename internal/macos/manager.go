package macos

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/engine"
)

// CurrentUser, sudoers kuralı ve plist için gerçek (login) kullanıcı adını döndürür.
func CurrentUser() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

// IsRunning, helper'a status sorarak tpws daemon'unun yüklü olup olmadığını döndürür.
// Parolasızdır (sudoers). Hata durumunda false döner.
func IsRunning() bool {
	out, err := exec.Command("sudo", "-n", HelperPath, "status").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "tpws: calisiyor")
}

// On, helper'ı parolasız sudo ile çağırarak servisi başlatır (PF + tpws daemon).
// Önce domainleri tpws'in okuyacağı hostlist dosyasına yazar, sonra helper'ı
// çağırır. Daemon, plist'teki --hostlist yolundan domainleri DOSYADAN okur;
// helper start her zaman bootout+bootstrap yaptığından tpws taze başlar ve
// güncel hostlist'i yeniden okur (add/remove sonrası yeni domain bypass görür).
func On(tpwsPort int, domains []string) error {
	if err := engine.WriteHostlist(domains); err != nil {
		return fmt.Errorf("hostlist yazılamadı: %w", err)
	}
	cmd := exec.Command("sudo", "-n", HelperPath, "start", strconv.Itoa(tpwsPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Off, helper ile servisi durdurur (tpws daemon bootout + anchor flush).
func Off() error {
	cmd := exec.Command("sudo", "-n", HelperPath, "stop")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Status, helper'dan ham durum metnini döndürür (parolasız).
func Status() (string, error) {
	out, err := exec.Command("sudo", "-n", HelperPath, "status").CombinedOutput()
	return string(out), err
}

// InstallScript, install sırasında TEK bir osascript admin diyaloğuyla
// çalıştırılacak kök kurulum bash betiğini üretir. Bu betik:
//   - helper'ı /usr/local/libexec'e yazar (root:wheel 0755)
//   - sudoers kuralını yazar (0440, visudo -cf ile doğrular)
//   - pf.conf'u yamar (anchor satırları)
//
// PLIST TEK KAYNAK (#2): LaunchDaemon plist'i bu betik YAZMAZ. Plist'in TEK
// üreticisi helper'ın write_plist'idir; install sonrası çağrılan
// `helper start <port>` plist'i güncel portla yazar. Böylece plist iki yerde
// (Go + bash) üretilmez ve sapma riski ortadan kalkar.
//
// Betik, gömülü heredoc'lar üzerinden dosya içeriklerini taşır; çağıran taraf
// bunu `osascript -e 'do shell script "..." with administrator privileges'` ile
// çalıştırır. Bu fonksiyon SAFtır (yalnızca metin üretir).
func InstallScript(patchedPFConf string, tpwsBin string, sudoUser string) string {
	helper := HelperScript(tpwsBin, engine.HostlistPath())
	sudoers := SudoersRule(sudoUser)

	var b strings.Builder
	b.WriteString("set -e\n")
	b.WriteString("mkdir -p /usr/local/libexec\n")

	// helper
	writeHeredoc(&b, HelperPath, helper)
	fmt.Fprintf(&b, "chown root:wheel %q\n", HelperPath)
	fmt.Fprintf(&b, "chmod 755 %q\n", HelperPath)

	// sudoers (önce geçici, visudo -cf ile doğrula, sonra taşı)
	tmpSudoers := SudoersPath + ".tmp"
	writeHeredoc(&b, tmpSudoers, sudoers)
	fmt.Fprintf(&b, "chmod 440 %q\n", tmpSudoers)
	fmt.Fprintf(&b, "visudo -cf %q\n", tmpSudoers)
	fmt.Fprintf(&b, "mv %q %q\n", tmpSudoers, SudoersPath)

	// pf.conf yaması
	writeHeredoc(&b, PFConfPath, patchedPFConf)

	// LaunchDaemon dizininin var olduğundan emin ol (plist'i helper yazacak).
	fmt.Fprintf(&b, "mkdir -p %q\n", filepath.Dir(LaunchDaemonPath))

	return b.String()
}

// UninstallScript, uninstall için TEK osascript admin bloğunda çalışacak betiği
// üretir. SAFtır (yalnızca metin üretir).
//
// BOOT GÜVENLİĞİ (#1): pf.conf'ta `load anchor ... from <ANCHOR_PATH>` satırı
// varken anchor dosyasını silmek, boot'ta pfctl -f /etc/pf.conf'un BAŞARISIZ
// olmasına yol açar (dangling load anchor → kural yüklenemez → internet bozulur).
// Bu yüzden SIRA katıdır ve pfctl hatası YUTULMAZ:
//
//  1. pf.conf'tan anchor satırlarını çıkar (yeni içerik yazılır).
//  2. pfctl -f /etc/pf.conf çalıştır ve BAŞARISINI doğrula.
//  3. SADECE (2) başarılıysa anchor dosyasını sil. pfctl -f başarısızsa anchor
//     dosyası KORUNUR (load anchor hâlâ geçerli kalır) ve hata raporlanır.
//
// LaunchDaemon/helper/sudoers silinmesi pf'den bağımsızdır ve her durumda yapılır.
func UninstallScript(unpatchedPFConf string) string {
	var b strings.Builder
	// Hata olsa bile betiğin geri kalanı kontrollü ilerlesin; pfctl başarısını
	// AÇIKÇA bayrakla izleriz (|| true ile yutmadan).
	b.WriteString("set +e\n")

	// Daemon'u durdur (best-effort; pf güvenliğini etkilemez).
	fmt.Fprintf(&b, "%q stop 2>/dev/null\n", HelperPath)

	// 1. pf.conf'u anchor satırları olmadan yeniden yaz.
	writeHeredoc(&b, PFConfPath, unpatchedPFConf)

	// 2. pf.conf'u yeniden yükle ve sonucu izle. BAŞARISIZ olursa anchor dosyasını
	//    SİLME (dangling load anchor önlenir) ve net bir hata raporla.
	fmt.Fprintf(&b, "if pfctl -f %q; then\n", PFConfPath)
	//    3. SADECE başarılıysa anchor dosyasını sil.
	fmt.Fprintf(&b, "  rm -f %q\n", AnchorPath)
	b.WriteString("else\n")
	fmt.Fprintf(&b, "  echo \"HATA: pfctl -f %s basarisiz — anchor dosyasi (%s) KORUNDU; pf.conf butunlugu icin manuel kontrol gerekli.\" >&2\n", PFConfPath, AnchorPath)
	b.WriteString("fi\n")

	// LaunchDaemon, helper ve sudoers'ı sil (pf'den bağımsız; her durumda).
	for _, p := range []string{LaunchDaemonPath, HelperPath, SudoersPath} {
		fmt.Fprintf(&b, "rm -f %q\n", p)
	}
	return b.String()
}

// writeHeredoc, içeriği tek tırnaklı (genişletmesiz) heredoc ile dosyaya yazar.
// Sınırlayıcı tırnaklı olduğundan içerikteki $ ve “ kaçışsız korunur.
func writeHeredoc(b *strings.Builder, path, content string) {
	const delim = "SPOOFDPI_TR_EOF"
	fmt.Fprintf(b, "cat > %q <<'%s'\n", path, delim)
	b.WriteString(content)
	if !strings.HasSuffix(content, "\n") {
		b.WriteString("\n")
	}
	fmt.Fprintf(b, "%s\n", delim)
}

// RunAdmin, verilen bash betiğini tek bir macOS yönetici (admin) diyaloğuyla
// çalıştırır: osascript ... with administrator privileges. Betik geçici bir
// dosyaya yazılır ve `bash <dosya>` admin olarak koşturulur (tırnak kaçışı
// dertlerini önler). YAN ETKİLİDİR — yalnızca install/uninstall'da çağrılır.
func RunAdmin(bashScript, prompt string) error {
	tmp, err := os.CreateTemp("", "spoofdpi-tr-*.sh")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(bashScript); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	// osascript: bash <tmp> komutunu admin yetkisiyle çalıştır.
	inner := fmt.Sprintf("/bin/bash %s", shellQuote(tmp.Name()))
	osa := fmt.Sprintf("do shell script %s with prompt %s with administrator privileges",
		appleScriptQuote(inner), appleScriptQuote(prompt))

	cmd := exec.Command("osascript", "-e", osa)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// shellQuote, tek bir yolu güvenli biçimde tek tırnağa alır.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// appleScriptQuote, AppleScript string literali üretir (çift tırnak + kaçış).
func appleScriptQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}
