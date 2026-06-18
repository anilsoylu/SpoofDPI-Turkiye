package config

import (
	"os"
	"testing"
)

func TestDefault(t *testing.T) {
	c := Default()
	if c.Port != 8080 {
		t.Errorf("varsayılan port 8080 olmalı, bulundu %d", c.Port)
	}
	if !c.EnableDoH {
		t.Error("varsayılan DoH açık olmalı")
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	// HOME'u geçici dizine yönlendir, böylece gerçek kullanıcı config'i ezilmez.
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	want := Default()
	want.Port = 9090
	want.Domains = []string{"discord.com", "discord.gg"}
	if err := want.Save(); err != nil {
		t.Fatalf("Save hatası: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load hatası: %v", err)
	}
	if got.Port != 9090 {
		t.Errorf("port 9090 bekleniyordu, bulundu %d", got.Port)
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
	if c.Port != 8080 {
		t.Errorf("eksik config Default() döndürmeli, port %d", c.Port)
	}
	_ = os.Getenv // (lint susturucu; gerçek kullanım yukarıda)
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
