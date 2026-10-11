package engine

import (
	"os"
	"path/filepath"
	"testing"
)

// TestBuildFromSourceInvalidRepo, klon başarısızsa hata dönüldüğünü ve
// mevcut binary'ye dokunulmadığını garanti eder (ağ gerektirmez).
func TestBuildFromSourceInvalidRepo(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	if err := BuildFromSource("file:///var/yok/olmayan-repo"); err == nil {
		t.Fatal("geçersiz repo ile BuildFromSource hata dönmeli")
	}
	if _, err := os.Stat(filepath.Join(tmp, ".spoofdpi-tr", "bin", "tpws")); !os.IsNotExist(err) {
		t.Error("başarısız derleme binary yazmamalı")
	}
}
