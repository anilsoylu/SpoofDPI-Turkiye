package spoofdpi

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestAssetArch(t *testing.T) {
	// assetArch dönen değer "arm64" veya "x86_64" olmalı (bu makinede arm64).
	arch, err := assetArch()
	if err != nil {
		t.Fatalf("assetArch hata verdi: %v", err)
	}
	if arch != "arm64" && arch != "x86_64" {
		t.Errorf("beklenen arm64 veya x86_64, bulundu %q", arch)
	}
}

func TestPickAsset(t *testing.T) {
	arch, _ := assetArch()
	rel := &ghRelease{
		TagName: "v1.5.3",
		Assets: []ghAsset{
			{Name: "spoofdpi_1.5.3_darwin_" + arch + ".tar.gz", DownloadURL: "http://x", Digest: "sha256:abc"},
			{Name: "spoofdpi_1.5.3_linux_amd64.tar.gz", DownloadURL: "http://y", Digest: "sha256:def"},
			{Name: "spoofdpi_1.5.3_darwin_" + arch + ".tar.gz.sbom.json", DownloadURL: "http://z", Digest: ""},
		},
	}
	got, err := pickAsset(rel)
	if err != nil {
		t.Fatalf("pickAsset hata: %v", err)
	}
	if got.Digest != "sha256:abc" {
		t.Errorf("yanlış asset seçildi: %+v", got)
	}
}

func makeTarGz(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	hdr := &tar.Header{Name: name, Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}
	tw.Close()
	gz.Close()
	return buf.Bytes()
}

func TestExtractBinary(t *testing.T) {
	content := []byte("#!/bin/sh\necho hi")
	tgz := makeTarGz(t, "spoofdpi", content)
	got, err := extractBinary(tgz)
	if err != nil {
		t.Fatalf("extractBinary hata: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("içerik uyuşmuyor: %q", got)
	}
}

func TestExtractBinaryMissing(t *testing.T) {
	tgz := makeTarGz(t, "baskabir", []byte("data"))
	_, err := extractBinary(tgz)
	if err == nil {
		t.Error("spoofdpi yokken hata bekleniyor")
	}
}

func TestVerifyChecksumMatch(t *testing.T) {
	data := []byte("test verisi")
	sum := sha256.Sum256(data)
	hex := hex.EncodeToString(sum[:])
	if err := verifyChecksum(data, hex); err != nil {
		t.Errorf("eşleşen checksum hata vermemeli: %v", err)
	}
}

func TestVerifyChecksumMismatch(t *testing.T) {
	data := []byte("test verisi")
	if err := verifyChecksum(data, "yanlis_hex_degeri_000000000000000000000000000000000000000000000000"); err == nil {
		t.Error("uyuşmayan checksum hata vermeli")
	}
}
