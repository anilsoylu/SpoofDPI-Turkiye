package macos

import (
	"strings"
	"testing"
)

func TestAnchorRules(t *testing.T) {
	r := AnchorRules(988)
	if !strings.Contains(r, "rdr on lo0") {
		t.Errorf("AnchorRules rdr kuralı içermeli:\n%s", r)
	}
	if !strings.Contains(r, "port 988") {
		t.Errorf("AnchorRules tpws portunu içermeli:\n%s", r)
	}
	if !strings.Contains(r, "user { >root }") {
		t.Errorf("AnchorRules döngü önleyici 'user { >root }' içermeli:\n%s", r)
	}
	if !strings.Contains(r, "!127.0.0.0/8") {
		t.Errorf("AnchorRules localhost dışlamalı:\n%s", r)
	}
	if !strings.Contains(r, "!192.168.0.0/16") {
		t.Errorf("AnchorRules LAN dışlamalı (Expo dev):\n%s", r)
	}
}

func TestPFConfWithAnchorPreservesPFOrdering(t *testing.T) {
	// Gerçek macOS varsayılan pf.conf'taki com.apple anchor bloğu.
	orig := "scrub-anchor \"com.apple/*\"\nnat-anchor \"com.apple/*\"\nrdr-anchor \"com.apple/*\"\n" +
		"dummynet-anchor \"com.apple/*\"\nanchor \"com.apple/*\"\nload anchor \"com.apple\" from \"/etc/pf.anchors/com.apple\"\n"
	out := PFConfWithAnchor(orig)

	for _, want := range anchorLines() {
		if !strings.Contains(out, want) {
			t.Errorf("yamada eksik satır: %q\nçıktı:\n%s", want, out)
		}
	}

	// PF tür sırası: translation (rdr) filtering (anchor) ÖNCESİNDE olmalı.
	idxOurRdr := strings.Index(out, `rdr-anchor "spoofdpi-tr"`)
	idxOurFilter := strings.Index(out, `anchor "spoofdpi-tr"`)
	idxOurLoad := strings.Index(out, `load anchor "spoofdpi-tr"`)
	if idxOurRdr < 0 || idxOurFilter < 0 || idxOurLoad < 0 {
		t.Fatalf("anchor satırları eksik:\n%s", out)
	}
	if idxOurRdr > idxOurFilter {
		t.Errorf("rdr-anchor (translation) filtering anchor'dan ÖNCE gelmeli:\n%s", out)
	}

	// Bizim rdr-anchor'ımız com.apple translation satırlarından SONRA gelmeli
	// (com.apple filtering anchor'ından önce kalmalı ki tür sırası bozulmasın).
	idxAppleRdr := strings.Index(out, `rdr-anchor "com.apple/*"`)
	idxAppleFilter := strings.Index(out, `anchor "com.apple/*"`)
	if idxOurRdr < idxAppleRdr {
		t.Errorf("bizim rdr-anchor com.apple rdr-anchor'ından sonra gelmeli:\n%s", out)
	}
	// Bizim filtering bloğumuz com.apple filtering anchor'ından sonra (en sonda).
	if idxOurFilter < idxAppleFilter {
		t.Errorf("bizim anchor (filtering) com.apple anchor'ından sonra gelmeli:\n%s", out)
	}
}

func TestPFConfWithAnchorIdempotent(t *testing.T) {
	orig := "rdr-anchor \"com.apple/*\"\n"
	once := PFConfWithAnchor(orig)
	twice := PFConfWithAnchor(once)
	if once != twice {
		t.Errorf("PFConfWithAnchor idempotent değil:\nbir:\n%s\niki:\n%s", once, twice)
	}
}

func TestPFConfWithAnchorNoApple(t *testing.T) {
	orig := "# bos pf.conf\nset skip on lo0\n"
	out := PFConfWithAnchor(orig)
	if !strings.Contains(out, `load anchor "spoofdpi-tr"`) {
		t.Errorf("com.apple olmadan da anchor eklenmeli:\n%s", out)
	}
}

func TestPFConfWithoutAnchorRemoves(t *testing.T) {
	orig := "rdr-anchor \"com.apple/*\"\n"
	patched := PFConfWithAnchor(orig)
	restored := PFConfWithoutAnchor(patched)

	for _, l := range anchorLines() {
		if strings.Contains(restored, l) {
			t.Errorf("anchor satırı silinmedi: %q\n%s", l, restored)
		}
	}
	if !strings.Contains(restored, "com.apple") {
		t.Errorf("orijinal com.apple satırı korunmalı:\n%s", restored)
	}
}
