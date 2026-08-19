package engine

import (
	"os"
	"strings"
	"testing"
)

func TestHostlistPathAbsolute(t *testing.T) {
	p := HostlistPath()
	if !strings.HasSuffix(p, "hostlist.txt") {
		t.Errorf("HostlistPath hostlist.txt ile bitmeli: %q", p)
	}
}

func TestWriteHostlist(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	if err := WriteHostlist([]string{"Discord.GG", "*.discord.com", "discord.gg"}); err != nil {
		t.Fatalf("WriteHostlist hata: %v", err)
	}
	data, err := os.ReadFile(HostlistPath())
	if err != nil {
		t.Fatalf("hostlist okunamadı: %v", err)
	}
	got := string(data)
	// NormalizeDomains: lower, "*." strip, tekilleştir, sırala; her satır bir domain.
	want := "discord.com\ndiscord.gg\n"
	if got != want {
		t.Errorf("hostlist içeriği=%q, beklenen %q", got, want)
	}
}
