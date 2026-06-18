package pac

import (
	"strings"
	"testing"
)

func TestGenerateContainsPortAndDomain(t *testing.T) {
	out := Generate(9090, []string{"discord.com"})
	if !strings.Contains(out, "PROXY 127.0.0.1:9090") {
		t.Error("PAC port içermeli")
	}
	if !strings.Contains(out, `"discord.com"`) {
		t.Error("PAC domain içermeli")
	}
	if !strings.Contains(out, "function FindProxyForURL") {
		t.Error("PAC fonksiyon imzası içermeli")
	}
}

func TestGenerateEmptyReturnsDirectOnly(t *testing.T) {
	out := Generate(8080, nil)
	if !strings.Contains(out, `return "DIRECT"`) {
		t.Error("boş listede DIRECT dönmeli")
	}
	if strings.Contains(out, "domains = [") {
		t.Error("boş listede domain dizisi olmamalı")
	}
}

func TestNormalizeDedupAndSort(t *testing.T) {
	got := Normalize([]string{"B.com", "a.com", "b.com", "  a.com  "})
	if len(got) != 2 || got[0] != "a.com" || got[1] != "b.com" {
		t.Errorf("beklenen [a.com b.com], bulundu %v", got)
	}
}

func TestNormalizeRejectsInvalid(t *testing.T) {
	got := Normalize([]string{"good.com", "no-dot", "bad;inject.com", "", "x.com\");evil"})
	for _, d := range got {
		if d != "good.com" {
			t.Errorf("yalnızca good.com geçmeli, geçersiz girdi: %q", d)
		}
	}
	if len(got) != 1 {
		t.Errorf("yalnızca 1 geçerli domain bekleniyordu, bulundu %d: %v", len(got), got)
	}
}

func TestValidDomainPunctuation(t *testing.T) {
	cases := map[string]bool{
		"discord.com": true, "a.b.c.com": true, "x-y.com": true,
		"nodot": false, ".lead.com": false, "trail.com.": false,
		"sp ace.com": false, "semi;.com": false,
	}
	for in, want := range cases {
		if got := validDomain(in); got != want {
			t.Errorf("validDomain(%q)=%v, beklenen %v", in, got, want)
		}
	}
}
