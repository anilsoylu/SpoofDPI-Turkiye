package macos

import (
	"fmt"
	"os"
	"strings"

	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/engine"
)

// ExpectedHelper, diske kurulması gereken güncel helper metnini döndürür.
// Yollar (tpws bin + hostlist) kurulum kullanıcısına göre sabit gömülür;
// karşılaştırma ve yenileme aynı şablonu kullanır.
func ExpectedHelper() string {
	return HelperScript(engine.BinPath(), engine.HostlistPath(), engine.ExcludePath())
}

// helperStaleContent, diskteki helper birebir güncel şablonla eşleşmiyorsa
// true döner. Saf karşılaştırma (test edilebilir); dosya okuma HelperStale'de.
func helperStaleContent(onDisk, expected string) bool {
	return onDisk != expected
}

// HelperStale, kurulu helper'ın güncel olup olmadığını döndürür. Helper
// root:sahipli ama dünyaca okunabilir (0755) olduğundan karşılaştırma
// ayrıcalıksız yapılır. Dosya yoksa/okunamazsa (true, err) döner.
func HelperStale() (bool, error) {
	data, err := os.ReadFile(HelperPath)
	if err != nil {
		return true, err
	}
	return helperStaleContent(string(data), ExpectedHelper()), nil
}

// refreshHelperScript, helper dosyasını güncel şablonla yenileyen kök betiği
// üretir. SAFtır (yalnızca metin üretir); yürütme RefreshHelper'dadır.
// Yollar/ID'ler değişmediğinden yalnızca helper dosyası yenilenir
// (sudoers/pf.conf/plist'e dokunulmaz).
func refreshHelperScript() string {
	var b strings.Builder
	writeHeredoc(&b, HelperPath, ExpectedHelper())
	fmt.Fprintf(&b, "chown root:wheel %q\n", HelperPath)
	fmt.Fprintf(&b, "chmod 755 %q\n", HelperPath)
	return b.String()
}

// RefreshHelper, eski helper'ı güncel şablonla değiştirir. TEK yönetici
// diyaloğu ister (helper root:sahipli). YAN ETKİLİDİR — yalnızca
// upgrade-helper komutundan çağrılır.
func RefreshHelper() error {
	return RunAdmin(refreshHelperScript(), "SpoofDPI helper güncelleme izni (tek seferlik)")
}
