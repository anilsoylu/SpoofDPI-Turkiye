// Package macos, macOS sistem entegrasyonunu (LaunchAgent + ağ proxy) yönetir.
package macos

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const launchLabel = "com.spoofdpi-tr"

// plistPath, LaunchAgent plist yolunu döndürür.
func plistPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "LaunchAgents", launchLabel+".plist"), nil
}

// buildPlist, verilen argümanlarla spoofdpi'yi çalıştıran plist XML'ini üretir.
// Saf fonksiyon — test edilebilir.
func buildPlist(binPath string, args []string) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">` + "\n")
	sb.WriteString(`<plist version="1.0">` + "\n<dict>\n")
	sb.WriteString("  <key>Label</key>\n  <string>" + launchLabel + "</string>\n")
	sb.WriteString("  <key>ProgramArguments</key>\n  <array>\n")
	sb.WriteString("    <string>" + xmlEscape(binPath) + "</string>\n")
	for _, a := range args {
		sb.WriteString("    <string>" + xmlEscape(a) + "</string>\n")
	}
	sb.WriteString("  </array>\n")
	sb.WriteString("  <key>RunAtLoad</key>\n  <true/>\n")
	sb.WriteString("  <key>KeepAlive</key>\n  <true/>\n")
	sb.WriteString("  <key>StandardOutPath</key>\n  <string>/tmp/spoofdpi-tr.log</string>\n")
	sb.WriteString("  <key>StandardErrorPath</key>\n  <string>/tmp/spoofdpi-tr.error.log</string>\n")
	sb.WriteString("</dict>\n</plist>\n")
	return sb.String()
}

func xmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return r.Replace(s)
}

// writeAndLoad, plist'i yazar ve launchctl ile yükler. Zaten yüklüyse önce çıkarır.
func writeAndLoad(binPath string, args []string) error {
	p, err := plistPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	// Zaten yüklüyse temiz başlamak için çıkar (hata yok say).
	_ = exec.Command("launchctl", "unload", "-w", p).Run()
	if err := os.WriteFile(p, []byte(buildPlist(binPath, args)), 0o644); err != nil {
		return err
	}
	out, err := exec.Command("launchctl", "load", "-w", p).CombinedOutput()
	if err != nil {
		return fmt.Errorf("launchctl load başarısız: %v: %s", err, out)
	}
	return nil
}

// unload, servisi durdurur ve plist'i kaldırır (idempotent).
func unload() error {
	p, err := plistPath()
	if err != nil {
		return err
	}
	_ = exec.Command("launchctl", "unload", "-w", p).Run()
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// isLoaded, servis launchctl'de yüklü mü kontrol eder.
func isLoaded() bool {
	out, err := exec.Command("launchctl", "list").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), launchLabel)
}
