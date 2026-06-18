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

- Release yalnızca **yönetici (manager) binary'sini** (`spoofdpi-tr`) içerir.
  **tpws motoru release'e dahil DEĞİLDİR.** `install.sh`, tpws'i kullanıcı
  makinesinde kaynaktan derler (`git clone --depth 1 https://github.com/bol-van/zapret`
  → `make mac` → `~/.spoofdpi-tr/bin/tpws`). Güvenlik gereği indirilmiş binary
  dağıtılmaz; bu yüzden GoReleaser yapılandırmasında tpws yoktur.
- İlk release çıkana kadar `install.sh`, release bulamayınca `go install`
  ile manager'ı kaynaktan derler — bu beklenen davranıştır.
- Sürüm `main.version`'a `-ldflags "-X main.version={{.Version}}"` ile gömülür;
  `spoofdpi-tr version` bunu gösterir.
