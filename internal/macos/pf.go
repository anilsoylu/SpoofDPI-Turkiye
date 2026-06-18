// Package macos, macOS Packet Filter (PF) ile tpws redirect kurallarını üreten
// ve /etc/pf.conf'u yamayan saf fonksiyonlar sağlar. Bu dosyadaki fonksiyonlar
// yan etkisizdir (metin üretir); gerçek pfctl çağrıları helper script'tedir.
package macos

import (
	"fmt"
	"strings"
)

// AnchorName, /etc/pf.anchors altındaki anchor dosyasının ve pf.conf'taki
// anchor referanslarının ortak adıdır.
const AnchorName = "spoofdpi-tr"

// Yol sabitleri — helper script ve manager bu yolları paylaşır.
const (
	AnchorPath       = "/etc/pf.anchors/" + AnchorName
	PFConfPath       = "/etc/pf.conf"
	HelperPath       = "/usr/local/libexec/spoofdpi-tr-helper"
	SudoersPath      = "/etc/sudoers.d/spoofdpi-tr"
	LaunchDaemonPath = "/Library/LaunchDaemons/com.spoofdpi-tr.plist"
	LaunchDaemonID   = "com.spoofdpi-tr"
)

// AnchorRules, /etc/pf.anchors/spoofdpi-tr içeriğini üretir.
//
// rdr kuralı: lo0 üzerinden 443'e giden (localhost kaynaklı olmayan) trafiği
// tpws portuna yönlendirir. pass out route-to: trafiği lo0 üzerinden geri sokar.
// `user { >root }` döngüyü önler (ŞART): root olarak çalışan tpws'in kendi
// çıkışı yeniden yönlendirilmez.
//
// LAN dışlama: Expo/LAN geliştirmesi (cihaz<->Mac) bozulmasın diye özel ağ
// blokları kurallardan dışlanır.
func AnchorRules(tpwsPort int) string {
	// 127.0.0.0/8 + özel LAN blokları dışlanır. Çoklu negatif adresler PF'de
	// liste `{ }` içine alınmalıdır (yoksa pfctl "syntax error" verir).
	exclude := "{ !127.0.0.0/8 !192.168.0.0/16 !10.0.0.0/8 !172.16.0.0/12 }"
	var b strings.Builder
	fmt.Fprintf(&b, "rdr on lo0 inet proto tcp from %s to any port { 443 } -> 127.0.0.1 port %d\n", exclude, tpwsPort)
	fmt.Fprintf(&b, "pass out route-to (lo0 127.0.0.1) inet proto tcp from %s to any port { 443 } user { >root }\n", exclude)
	return b.String()
}

// anchorLines, pf.conf'a eklenecek üç satırdır. `load anchor` satırı ŞART
// (yoksa anchor boş kalır — PoC'de bulundu).
func anchorLines() []string {
	return []string{
		fmt.Sprintf("rdr-anchor %q", AnchorName),
		fmt.Sprintf("anchor %q", AnchorName),
		fmt.Sprintf("load anchor %q from %q", AnchorName, AnchorPath),
	}
}

// PFConfWithAnchor, verilen pf.conf metnine üç anchor satırını ekler.
//
// PF, kuralların TÜR SIRASINA uymasını zorunlu kılar:
// options → normalization(scrub) → queueing → translation(nat/rdr) → filtering(anchor).
// Bu yüzden satırlar tek bir blok olarak eklenemez; her satır kendi türünün
// doğru bölümüne yerleştirilmelidir:
//   - rdr-anchor (translation): son rdr-anchor/nat-anchor satırından HEMEN SONRA.
//   - anchor + load anchor (filtering): dosyanın SONUNA (tüm filtering kurallarından sonra).
//
// Idempotent: satırlar zaten varsa metin değişmeden döner.
func PFConfWithAnchor(original string) string {
	lines := anchorLines()

	// Idempotency: her satır zaten varsa dokunma.
	already := true
	for _, l := range lines {
		if !strings.Contains(original, l) {
			already = false
			break
		}
	}
	if already {
		return original
	}

	// Önce eski (kısmi) anchor satırlarını temizle.
	cleaned := PFConfWithoutAnchor(original)

	rdrLine := fmt.Sprintf("rdr-anchor %q", AnchorName)
	filterBlock := fmt.Sprintf("anchor %q\nload anchor %q from %q",
		AnchorName, AnchorName, AnchorPath)

	src := strings.Split(cleaned, "\n")

	// rdr-anchor satırını, var olan son translation (rdr-anchor/nat-anchor)
	// satırından sonra yerleştir. Yoksa filtering bloğuyla birlikte sona düşer.
	rdrAfter := -1
	for i, l := range src {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "rdr-anchor ") || strings.HasPrefix(t, "nat-anchor ") {
			rdrAfter = i
		}
	}

	var b strings.Builder
	for i, l := range src {
		b.WriteString(l)
		if i < len(src)-1 {
			b.WriteString("\n")
		}
		if i == rdrAfter {
			b.WriteString(rdrLine + "\n")
		}
	}

	out := b.String()
	if out != "" && !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	// translation (rdr) satırı için uygun bir konum bulunamadıysa, filtering
	// bloğundan ÖNCE ekle (PF tür sırasını korumak için).
	if rdrAfter < 0 {
		out += rdrLine + "\n"
	}
	// anchor + load anchor (filtering) en sona — tüm filtering kurallarından sonra.
	out += filterBlock + "\n"
	return out
}

// PFConfWithoutAnchor, eklenen üç anchor satırını metinden çıkarır (uninstall).
// Anchor'a ait herhangi bir satır içeren tüm satırlar atılır.
func PFConfWithoutAnchor(patched string) string {
	src := strings.Split(patched, "\n")
	var kept []string
	for _, l := range src {
		t := strings.TrimSpace(l)
		if isAnchorLine(t) {
			continue
		}
		kept = append(kept, l)
	}
	return strings.Join(kept, "\n")
}

// isAnchorLine, bir satırın bizim anchor referanslarımızdan biri olup olmadığını
// belirler. Hem rdr-anchor/anchor/load anchor varyantlarını yakalar.
func isAnchorLine(t string) bool {
	q := `"` + AnchorName + `"`
	if !strings.Contains(t, q) {
		return false
	}
	return strings.HasPrefix(t, "rdr-anchor ") ||
		strings.HasPrefix(t, "anchor ") ||
		strings.HasPrefix(t, "load anchor ")
}
