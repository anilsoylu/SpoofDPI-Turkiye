package macos

import (
	"strings"
	"testing"
)

func TestHelperScriptCommands(t *testing.T) {
	s := HelperScript("/Users/x/.spoofdpi-tr/bin/tpws", "/Users/x/.spoofdpi-tr/hostlist.txt")
	for _, want := range []string{"start)", "stop)", "status)", "pfctl", "launchctl"} {
		if !strings.Contains(s, want) {
			t.Errorf("HelperScript %q içermeli", want)
		}
	}
	if !strings.HasPrefix(s, "#!/bin/bash") {
		t.Error("HelperScript shebang ile başlamalı")
	}
	// BUG1: helper plist'i kendisi GÜNCEL portla üretmeli (gömülü tpws bin +
	// hostlist). write_plist fonksiyonu ve gömülü yollar mevcut olmalı.
	for _, want := range []string{
		"write_plist",
		"/Users/x/.spoofdpi-tr/bin/tpws",
		"/Users/x/.spoofdpi-tr/hostlist.txt",
		"--split-pos=1,midsld",
		"--oob=tls",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("HelperScript plist üretimi %q içermeli", want)
		}
	}
	// BUG2: tpws'in portu gerçekten dinlediğini doğrulayan kontrol olmalı.
	if !strings.Contains(s, "verify_listen") || !strings.Contains(s, "LISTEN") {
		t.Error("HelperScript tpws listen doğrulaması içermeli")
	}
}

// TestHelperScriptValidatesPort, helper'ın ayrıcalık sınırında port'u
// doğruladığını garanti eder. sudoers 'helper *' olduğundan kullanıcı keyfi
// argümanla çağırabilir; helper sayısal-olmayan port'u reddetmelidir.
func TestHelperScriptValidatesPort(t *testing.T) {
	s := HelperScript("/Users/x/.spoofdpi-tr/bin/tpws", "/Users/x/.spoofdpi-tr/hostlist.txt")
	// Rakam-dışı reddi (case kalıbı) ve aralık kontrolü mevcut olmalı.
	for _, want := range []string{"*[!0-9]*", "65535", "-lt 1"} {
		if !strings.Contains(s, want) {
			t.Errorf("HelperScript port doğrulaması %q içermeli", want)
		}
	}
	// set -euo pipefail hâlâ korunmalı.
	if !strings.Contains(s, "set -euo pipefail") {
		t.Error("HelperScript 'set -euo pipefail' içermeli")
	}
}

// TestHelperPlistSingleSource, plist'in TEK üreticisinin helper write_plist
// olduğunu (#2) ve doğru anahtarları/argümanları ürettiğini doğrular. Go tarafı
// artık plist üretmediğinden plist içeriği yalnızca helper script'te yaşar.
func TestHelperPlistSingleSource(t *testing.T) {
	s := HelperScript("/Users/x/.spoofdpi-tr/bin/tpws", "/Users/x/.spoofdpi-tr/hostlist.txt")
	for _, want := range []string{
		"write_plist",
		"<key>Label</key>",
		"<key>RunAtLoad</key>",
		"<key>KeepAlive</key>",
		"<key>SuccessfulExit</key>",
		"<key>ThrottleInterval</key>",
		"<key>ProcessType</key>",
		"<string>Adaptive</string>",
		"<key>StandardOutPath</key>",
		"<key>StandardErrorPath</key>",
		"--user=root",
		"--bind-addr=127.0.0.1",
		"--hostlist=${HOSTLIST}",
		"--split-pos=1,midsld",
		"--oob=tls",
		"--port=${port}",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("helper write_plist %q içermeli", want)
		}
	}
}

// TestUninstallScriptSafeOrdering, #1 boot güvenliğini garanti eder: pf.conf'tan
// anchor satırları çıkarılır, pfctl -f çalıştırılır ve anchor dosyası YALNIZCA
// pfctl başarılıysa silinir. pfctl hatası YUTULMAMALIdır (|| true yok).
func TestUninstallScriptSafeOrdering(t *testing.T) {
	s := UninstallScript("# pf.conf without anchor\n")

	// pfctl çağrısı bir koşulun içinde olmalı (başarı doğrulanmalı).
	if !strings.Contains(s, "if pfctl -f") {
		t.Error("UninstallScript pfctl başarısını koşulla doğrulamalı (if pfctl -f)")
	}
	// pfctl hatası yutulmamalı.
	if strings.Contains(s, "pfctl -f \""+PFConfPath+"\" 2>/dev/null || true") ||
		strings.Contains(s, "pfctl -f "+PFConfPath+" || true") {
		t.Error("UninstallScript pfctl hatasını '|| true' ile YUTMAMALI")
	}
	// SIRA: anchor dosyası rm'i, pfctl -f satırından SONRA gelmeli.
	pfIdx := strings.Index(s, "if pfctl -f")
	anchorRmIdx := strings.Index(s, "rm -f \""+AnchorPath+"\"")
	if pfIdx < 0 || anchorRmIdx < 0 {
		t.Fatalf("beklenen satırlar yok: pfIdx=%d anchorRmIdx=%d", pfIdx, anchorRmIdx)
	}
	if anchorRmIdx < pfIdx {
		t.Error("anchor dosyası silme pfctl -f'den ÖNCE gelmemeli (#1)")
	}
	// pf.conf yeniden yazımı (heredoc) anchor rm'den ÖNCE olmalı.
	pfConfWriteIdx := strings.Index(s, "cat > \""+PFConfPath+"\"")
	if pfConfWriteIdx < 0 || pfConfWriteIdx > anchorRmIdx {
		t.Error("pf.conf yeniden yazımı anchor silmeden ÖNCE olmalı (#1)")
	}
}

func TestSudoersRule(t *testing.T) {
	r := SudoersRule("anil")
	want := "anil ALL=(root) NOPASSWD: /usr/local/libexec/spoofdpi-tr-helper *\n"
	if r != want {
		t.Errorf("SudoersRule=%q, beklenen %q", r, want)
	}
}
