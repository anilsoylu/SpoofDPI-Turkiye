package macos

import (
	"strings"
	"testing"
)

func TestHelperStaleContent(t *testing.T) {
	want := ExpectedHelper()
	if !helperStaleContent("eski helper içeriği", want) {
		t.Error("farklı içerik stale sayılmalı")
	}
	if helperStaleContent(want, want) {
		t.Error("birebir eşleşen içerik güncel sayılmalı")
	}
	if !helperStaleContent(want+"\n", want) {
		t.Error("tek bayt farkı bile stale sayılmalı (sessiz sapma yok)")
	}
}

func TestRefreshHelperScript(t *testing.T) {
	s := refreshHelperScript()
	for _, want := range []string{
		"cat > \"" + HelperPath + "\"",
		"chown root:wheel \"" + HelperPath + "\"",
		"chmod 755 \"" + HelperPath + "\"",
		"SPOOFDPI_TR_EOF",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("yenileme betiği %q içermeli", want)
		}
	}
	// Yenileme SADECE helper dosyasına dokunur; sudoers/pf.conf/plist yazmamalı.
	for _, dont := range []string{SudoersPath, PFConfPath, LaunchDaemonPath} {
		if strings.Contains(s, "cat > \""+dont+"\"") {
			t.Errorf("yenileme betiği %q dosyasına YAZMAMALI", dont)
		}
	}
}
