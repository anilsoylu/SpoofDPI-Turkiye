package macos

import (
	"strings"
	"testing"

	"github.com/anilsoylu/SpoofDPI-Turkiye/internal/config"
)

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
