#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "Swift release build başlatılıyor..."
swift build -c release

BINARY_PATH="$(swift build -c release --show-bin-path 2>/dev/null)/SpoofDPITurkiye"

APP_NAME="SpoofDPI-Türkiye.app"
APP_DIR="$SCRIPT_DIR/$APP_NAME"
CONTENTS="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS/MacOS"

# Önceki bundle'ı temizle
rm -rf "$APP_DIR"
mkdir -p "$MACOS_DIR"

# Binary'yi kopyala
cp "$BINARY_PATH" "$MACOS_DIR/SpoofDPITurkiye"

# spoofdpi-tr CLI'ını derle ve bundle'a göm → app ayrı kuruluma gerek duymaz
RESOURCES_DIR="$CONTENTS/Resources"
mkdir -p "$RESOURCES_DIR"
echo "spoofdpi-tr CLI derleniyor ve bundle'a gömülüyor..."
( cd "$SCRIPT_DIR/.." && go build -ldflags "-X main.version=0.1.0" -o "$RESOURCES_DIR/spoofdpi-tr" ./cmd/spoofdpi-tr )

# Uygulama ikonunu bundle'a ekle
if [ -f "$SCRIPT_DIR/AppIcon.icns" ]; then
  cp "$SCRIPT_DIR/AppIcon.icns" "$RESOURCES_DIR/AppIcon.icns"
fi

# Info.plist yaz
cat > "$CONTENTS/Info.plist" << 'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>SpoofDPI Türkiye</string>
    <key>CFBundleIdentifier</key>
    <string>com.anilsoylu.spoofdpi-tr.menubar</string>
    <key>CFBundleExecutable</key>
    <string>SpoofDPITurkiye</string>
    <key>CFBundleIconFile</key>
    <string>AppIcon</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSMinimumSystemVersion</key>
    <string>14.0</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSAppTransportSecurity</key>
    <dict>
        <key>NSAllowsArbitraryLoads</key>
        <true/>
    </dict>
    <key>LSUIElement</key>
    <true/>
</dict>
</plist>
PLIST

echo ""
echo "Bundle oluşturuldu: $APP_DIR"
echo "Çalıştırmak için: open \"$APP_DIR\""
