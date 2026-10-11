package engine

import (
	"os"
	"path/filepath"
	"strings"
)

// DesyncExcludes, desync uygulanmadan düz relay'lenen host'lardır (tpws
// --hostlist-exclude; subdomain auto apply). Bunlar DPI tarafından engelli
// DEĞİLDIR ama split+oob desync'i tolere edemezler:
//
//	dl2.discordapp.net — Discord masaüstü güncelleme CDN'i (Google ön yüz).
//	  Desync'li ClientHello'ya yanıt vermez (handshake takılır, 15sn timeout);
//	  düz relay ve direkt bağlantı ~0.2sn'de 200 döner. Liste dar tutulur:
//	  desync'ten muafiyet, engelli bir host'a yanlışlıkla uygulanmamalı.
var DesyncExcludes = []string{
	"dl2.discordapp.net",
}

// ExcludePath, tpws'in okuyacağı exclude dosyasının tam yolunu döndürür:
// ~/.spoofdpi-tr/excludes.txt (mutlak; root tarafından okunabilir).
func ExcludePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".spoofdpi-tr", "excludes.txt")
	}
	return filepath.Join(home, ".spoofdpi-tr", "excludes.txt")
}

// WriteExcludes, DesyncExcludes listesini ExcludePath()'e yazar (her satır
// bir host). Dizini gerekirse oluşturur.
func WriteExcludes() error {
	path := ExcludePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strings.Join(DesyncExcludes, "\n")+"\n"), 0o644)
}
