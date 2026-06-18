# SpoofDPI Türkiye — Menu Bar Uygulaması

macOS menu bar'dan spoofdpi-tr CLI'ını yöneten SwiftUI uygulaması.

## Gereksinimler

- macOS 14+
- `spoofdpi-tr` CLI kurulu olmalı (PATH'te veya `/usr/local/bin`, `/opt/homebrew/bin`, `~/.local/bin`)

## Build

```bash
cd menubar
./build-app.sh
```

Çıktı: `menubar/SpoofDPI-Türkiye.app`

Debug build için:

```bash
cd menubar
swift build
```

## Not

Bu uygulama tamamen `spoofdpi-tr` CLI'ına bağımlıdır. CLI kurulu değilse uygulama kurulum talimatlarını gösterir.
