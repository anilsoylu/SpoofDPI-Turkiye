#!/bin/bash
set -euo pipefail

BIN_NAME="spoofdpi-tr"

# Sistem dosyası yolları (internal/macos ile birebir aynı).
HELPER_PATH="/usr/local/libexec/spoofdpi-tr-helper"
SUDOERS_PATH="/etc/sudoers.d/spoofdpi-tr"
ANCHOR_PATH="/etc/pf.anchors/spoofdpi-tr"
LAUNCHD_PATH="/Library/LaunchDaemons/com.spoofdpi-tr.plist"
LAUNCHD_ID="com.spoofdpi-tr"
PF_CONF="/etc/pf.conf"
SPOOF_DIR="$HOME/.spoofdpi-tr"

# Platform kontrolü — yalnızca macOS.
if [ "$(uname -s)" != "Darwin" ]; then
  echo "Bu araç şu an yalnızca macOS destekler." >&2
  exit 1
fi

echo "SpoofDPI Türkiye kaldırılıyor..."

# Yöntem 1: binary PATH'teyse kendi kendine temizlesin (-y ile onay atlama).
# Bu, pf.conf'tan anchor'ı çıkarır, helper/sudoers/anchor/LaunchDaemon'u siler
# ve ~/.spoofdpi-tr dizinini (tpws binary dahil) kaldırır.
if command -v "$BIN_NAME" &>/dev/null; then
  "$BIN_NAME" uninstall -y
  echo "✓ Kaldırıldı"
  exit 0
fi

# Yöntem 2: Fallback — binary bulunamazsa manuel temizlik (root gerekir).
echo "Binary bulunamadı — manuel temizlik (yönetici parolası istenebilir)..."

# LaunchDaemon'u durdur (tpws root daemon).
sudo launchctl bootout system "$LAUNCHD_PATH" 2>/dev/null \
  || sudo launchctl unload -w "$LAUNCHD_PATH" 2>/dev/null || true

# pf.conf'tan anchor satırlarımızı çıkar ve pf'yi yeniden yükle.
if [ -f "$PF_CONF" ]; then
  sudo sed -i '' \
    -e '/^[[:space:]]*rdr-anchor[[:space:]]*"spoofdpi-tr"/d' \
    -e '/^[[:space:]]*anchor[[:space:]]*"spoofdpi-tr"/d' \
    -e '/^[[:space:]]*load anchor[[:space:]]*"spoofdpi-tr"/d' \
    "$PF_CONF" 2>/dev/null || true
  sudo pfctl -f "$PF_CONF" 2>/dev/null || true
fi

# Sistem dosyalarını sil.
sudo rm -f "$HELPER_PATH" "$SUDOERS_PATH" "$ANCHOR_PATH" "$LAUNCHD_PATH" 2>/dev/null || true

# Yapılandırma + tpws binary dizinini sil.
rm -rf "$SPOOF_DIR"

# Manager binary'lerini sil (iki yaygın konum).
rm -f "/usr/local/bin/$BIN_NAME" "$HOME/.local/bin/$BIN_NAME"

echo "✓ Kaldırıldı"
