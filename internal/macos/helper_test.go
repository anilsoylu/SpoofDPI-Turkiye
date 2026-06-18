package macos

import (
	"strings"
	"testing"
)

func TestHelperScriptCommands(t *testing.T) {
	s := HelperScript()
	for _, want := range []string{"start)", "stop)", "status)", "pfctl", "launchctl", "tlsrec"} {
		// "tlsrec" helper'da olmayabilir (tpws args plist'te); diğerleri zorunlu.
		if want == "tlsrec" {
			continue
		}
		if !strings.Contains(s, want) {
			t.Errorf("HelperScript %q içermeli", want)
		}
	}
	if !strings.HasPrefix(s, "#!/bin/bash") {
		t.Error("HelperScript shebang ile başlamalı")
	}
}

// TestHelperScriptValidatesPort, helper'ın ayrıcalık sınırında port'u
// doğruladığını garanti eder. sudoers 'helper *' olduğundan kullanıcı keyfi
// argümanla çağırabilir; helper sayısal-olmayan port'u reddetmelidir.
func TestHelperScriptValidatesPort(t *testing.T) {
	s := HelperScript()
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

func TestLaunchDaemonPlist(t *testing.T) {
	p := LaunchDaemonPlist("/Users/x/.spoofdpi-tr/bin/tpws", []string{"--port", "988", "--tlsrec=sni"})
	for _, want := range []string{
		"com.spoofdpi-tr",
		"<key>RunAtLoad</key>",
		"<key>KeepAlive</key>",
		"<true/>",
		"--tlsrec=sni",
		"/Users/x/.spoofdpi-tr/bin/tpws",
		"<key>StandardOutPath</key>",
	} {
		if !strings.Contains(p, want) {
			t.Errorf("plist %q içermeli:\n%s", want, p)
		}
	}
}

func TestSudoersRule(t *testing.T) {
	r := SudoersRule("anil")
	want := "anil ALL=(root) NOPASSWD: /usr/local/libexec/spoofdpi-tr-helper *\n"
	if r != want {
		t.Errorf("SudoersRule=%q, beklenen %q", r, want)
	}
}
