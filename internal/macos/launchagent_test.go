package macos

import (
	"strings"
	"testing"

	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/config"
)

// proxy.go — autoProxyURLMatches testleri.
func TestAutoProxyURLMatches(t *testing.T) {
	pac := "/Users/x/.spoofdpi-tr/proxy.pac"
	matching := "URL: file:///Users/x/.spoofdpi-tr/proxy.pac\nEnabled: Yes\n"
	other := "URL: file:///Users/x/other.pac\nEnabled: Yes\n"
	empty := ""

	if !autoProxyURLMatches(matching, pac) {
		t.Error("eşleşen çıktı true döndürmeli")
	}
	if autoProxyURLMatches(other, pac) {
		t.Error("farklı URL false döndürmeli")
	}
	if autoProxyURLMatches(empty, pac) {
		t.Error("boş çıktı false döndürmeli")
	}
	// URL satırı yoksa false.
	noURL := "Enabled: Yes\n"
	if autoProxyURLMatches(noURL, pac) {
		t.Error("URL satırı olmayan çıktı false döndürmeli")
	}
}

// service.go — fileURL testleri.
func TestFileURL(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"/Users/anil/.spoofdpi-tr/proxy.pac", "file:///Users/anil/.spoofdpi-tr/proxy.pac"},
		{"/Users/John Doe/x.pac", "file:///Users/John%20Doe/x.pac"},
		{"/simple/path", "file:///simple/path"},
	}
	for _, tc := range cases {
		got := fileURL(tc.path)
		if got != tc.want {
			t.Errorf("fileURL(%q) = %q, beklenen %q", tc.path, got, tc.want)
		}
	}
}

func TestBuildPlistContainsArgs(t *testing.T) {
	out := buildPlist("/bin/spoofdpi", []string{"-port", "9090", "-silent"})
	for _, want := range []string{
		"<string>/bin/spoofdpi</string>",
		"<string>-port</string>",
		"<string>9090</string>",
		"com.spoofdpi-tr",
		"RunAtLoad",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("plist %q içermeli", want)
		}
	}
}

func TestXMLEscape(t *testing.T) {
	if got := xmlEscape("a&b<c>"); got != "a&amp;b&lt;c&gt;" {
		t.Errorf("xmlEscape yanlış: %q", got)
	}
}

func TestSpoofdpiArgs(t *testing.T) {
	c := &config.Config{Port: 9090, EnableDoH: true, DNSAddr: "1.1.1.1"}
	args := strings.Join(spoofdpiArgs(c), " ")
	for _, want := range []string{"-port 9090", "-system-proxy=false", "-silent", "-enable-doh", "-dns-addr 1.1.1.1"} {
		if !strings.Contains(args, want) {
			t.Errorf("args %q içermeli, bulundu: %s", want, args)
		}
	}
}

func TestSpoofdpiArgsNoDoH(t *testing.T) {
	c := &config.Config{Port: 8080}
	args := strings.Join(spoofdpiArgs(c), " ")
	if strings.Contains(args, "-enable-doh") {
		t.Error("DoH kapalıyken -enable-doh olmamalı")
	}
}
