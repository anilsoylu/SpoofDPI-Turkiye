package blocklist

import "testing"

func TestCategoriesNotEmpty(t *testing.T) {
	cats := Categories()
	if len(cats) == 0 {
		t.Fatal("en az bir kategori olmalı")
	}
}

func TestGetDiscord(t *testing.T) {
	c, ok := Get("discord")
	if !ok {
		t.Fatal("discord kategorisi bulunmalı")
	}
	if len(c.Domains) == 0 {
		t.Error("discord domainleri boş olmamalı")
	}
}

func TestGetUnknown(t *testing.T) {
	if _, ok := Get("yokboyle"); ok {
		t.Error("bilinmeyen anahtar ok=false dönmeli")
	}
}

func TestDefaultDomainsIncludesDiscordCom(t *testing.T) {
	found := false
	for _, d := range DefaultDomains() {
		if d == "discord.com" {
			found = true
		}
	}
	if !found {
		t.Error("varsayılan domainler discord.com içermeli")
	}
}
