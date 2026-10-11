package engine

import (
	"os"
	"testing"
)

func TestWriteExcludes(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	if err := WriteExcludes(); err != nil {
		t.Fatalf("WriteExcludes hata: %v", err)
	}
	data, err := os.ReadFile(ExcludePath())
	if err != nil {
		t.Fatalf("excludes okunamadı: %v", err)
	}
	got := string(data)
	want := "dl2.discordapp.net\n"
	if got != want {
		t.Errorf("excludes içeriği=%q, beklenen %q", got, want)
	}
}
