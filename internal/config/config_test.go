package config

import (
	"os"
	"testing"
)

func TestDefault(t *testing.T) {
	c := Default()
	if c.Port != 988 {
		t.Errorf("varsayılan port 988 olmalı, bulundu %d", c.Port)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	// HOME'u geçici dizine yönlendir, böylece gerçek kullanıcı config'i ezilmez.
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	want := Default()
	want.Port = 988
	want.Domains = []string{"discord.com", "discord.gg"}
	if err := want.Save(); err != nil {
		t.Fatalf("Save hatası: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load hatası: %v", err)
	}
	if got.Port != 988 {
		t.Errorf("port 988 bekleniyordu, bulundu %d", got.Port)
	}
	if len(got.Domains) != 2 {
		t.Errorf("2 domain bekleniyordu, bulundu %d", len(got.Domains))
	}
}

func TestLoadMissingReturnsDefault(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	c, err := Load()
	if err != nil {
		t.Fatalf("eksik config hata vermemeli: %v", err)
	}
	if c.Port != 988 {
		t.Errorf("eksik config Default() döndürmeli, port %d", c.Port)
	}
	_ = os.Getenv // (lint susturucu; gerçek kullanım yukarıda)
}

func TestNormalizeDomains(t *testing.T) {
	in := []string{
		"  Discord.com ",
		"*.discord.gg",
		".discordapp.com",
		"discord.com", // tekrar
		"",
		"  ",
		"DISCORD.GG", // *.discord.gg ile aynı normalize sonucu
	}
	got := NormalizeDomains(in)
	want := []string{"discord.com", "discord.gg", "discordapp.com"}
	if len(got) != len(want) {
		t.Fatalf("normalize: %d öğe bekleniyordu, %d bulundu: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("normalize[%d]=%q, beklenen %q (tümü: %v)", i, got[i], want[i], got)
		}
	}
}

func TestValidDomain(t *testing.T) {
	valid := []string{
		"discord.com",
		"a.b.c.example.co.uk",
		"xn--nxasmq6b.example.com", // punycode
		"my-host.example.org",
		"1.2.3.4.sslip.io",
	}
	for _, d := range valid {
		if !ValidDomain(d) {
			t.Errorf("ValidDomain(%q)=false, geçerli olmalı", d)
		}
	}

	invalid := []string{
		"",                  // boş
		"nodot",             // nokta yok
		".leading.com",      // baş nokta
		"trailing.com.",     // son nokta
		"-lead.com",         // baş tire
		"trail-.com",        // etiket sonu tire
		"a..b.com",          // ardışık nokta
		"evil.com;rm -rf /", // shell metakarakter
		"$(reboot).com",     // command substitution
		"a`whoami`.com",     // backtick
		"a b.com",           // boşluk
		"a\n.com",           // yeni satır
		"under_score.com",   // alt çizgi
		"UPPER.com",         // büyük harf (normalize öncesi reddedilmeli)
		"a|b.com",           // pipe
		"a/b.com",           // slash
	}
	for _, d := range invalid {
		if ValidDomain(d) {
			t.Errorf("ValidDomain(%q)=true, reddedilmeli (enjeksiyon/geçersiz)", d)
		}
	}
}

func TestNormalizeDomainsRejectsInjection(t *testing.T) {
	in := []string{
		"discord.com",
		"evil.com;rm -rf /",
		"$(reboot).com",
		"a b.com",
		"good.example.org",
		"a`id`.com",
		"nodot",
	}
	got := NormalizeDomains(in)
	want := map[string]bool{"discord.com": true, "good.example.org": true}
	if len(got) != len(want) {
		t.Fatalf("normalize enjeksiyonu süzmeli: bulundu %v", got)
	}
	for _, d := range got {
		if !want[d] {
			t.Errorf("beklenmeyen/zararlı domain geçti: %q", d)
		}
	}
}

func TestValidatePort(t *testing.T) {
	// Geçersiz portlar hata döndürmeli.
	for _, p := range []int{0, -1, 65536, -1000} {
		if err := ValidatePort(p); err == nil {
			t.Errorf("ValidatePort(%d) hata vermeli, vermedi", p)
		}
	}
	// Geçerli portlar nil döndürmeli.
	for _, p := range []int{1, 80, 443, 8080, 65535} {
		if err := ValidatePort(p); err != nil {
			t.Errorf("ValidatePort(%d) hata vermemeli, verdi: %v", p, err)
		}
	}
}
