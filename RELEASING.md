# Sürüm Çıkarma / Releasing

Bu proje sürümleri **GoReleaser** + **GitHub Actions** ile otomatik üretir.
`v*` biçiminde bir etiket push edildiğinde `.github/workflows/release.yml`
tetiklenir, darwin `arm64` ve `x86_64` için `.tar.gz` binary'leri üretir ve
GitHub Release'e yükler.

## Yeni sürüm çıkarmak

```bash
# 1. master güncel ve testler geçiyor olsun
go test ./...

# 2. Etiketle ve push et (semver, 'v' önekiyle)
git tag v0.1.0
git push origin v0.1.0
```

Bu kadar. Actions iş akışı geri kalanı halleder. Birkaç dakika içinde
Releases sayfasında şu asset'ler oluşur:

```
spoofdpi-tr_0.1.0_darwin_arm64.tar.gz
spoofdpi-tr_0.1.0_darwin_x86_64.tar.gz
checksums.txt
```

Bu isimlendirme `install.sh` ve `uninstall.sh`'ın beklediği desenle birebir
eşleşir — yani release çıktığı an `curl -fsSL .../install.sh | bash` çalışır.

## Yerel test (opsiyonel)

GoReleaser kuruluysa yayınlamadan dene:

```bash
brew install goreleaser
goreleaser release --snapshot --clean   # yayınlamaz, sadece dist/ üretir
goreleaser check                         # config doğrulaması
```

## Notlar

- İlk release çıkana kadar `install.sh` 404 döner ve `go install` fallback'ini
  önerir — bu beklenen davranıştır.
- Sürüm `main.version`'a `-ldflags "-X main.version={{.Version}}"` ile gömülür;
  `spoofdpi-tr version` bunu gösterir.
