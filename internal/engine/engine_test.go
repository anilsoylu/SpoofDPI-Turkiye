package engine

import (
	"os"
	"strings"
	"testing"
)

func TestArgsContainsTLSRecAndPort(t *testing.T) {
	args := Args(988)
	joined := strings.Join(args, " ")

	if !strings.Contains(joined, "--tlsrec=sni") {
		t.Errorf("Args --tlsrec=sni içermeli: %v", args)
	}
	if !strings.Contains(joined, "--bind-addr=127.0.0.1") {
		t.Errorf("Args --bind-addr=127.0.0.1 içermeli: %v", args)
	}
	if !strings.Contains(joined, "--user=root") {
		t.Errorf("Args --user=root içermeli: %v", args)
	}

	// Port, --port'tan hemen sonra ayrı argüman olarak gelmeli.
	foundPort := false
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "--port" && args[i+1] == "988" {
			foundPort = true
		}
	}
	if !foundPort {
		t.Errorf("Args port 988 içermeli: %v", args)
	}

	// Domainler artık bir DOSYADAN okunur: --hostlist <path> içermeli,
	// eski csv argümanı --hostlist-domains İÇERMEMELİ.
	if strings.Contains(joined, "--hostlist-domains") {
		t.Errorf("Args artık --hostlist-domains İÇERMEMELİ: %v", args)
	}
	foundHost := false
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "--hostlist" && args[i+1] == HostlistPath() {
			foundHost = true
		}
	}
	if !foundHost {
		t.Errorf("Args --hostlist %q içermeli: %v", HostlistPath(), args)
	}
}

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
