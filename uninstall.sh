#!/bin/bash
set -euo pipefail

BIN_NAME="spoofdpi-tr"
PAC_NAME="proxy.pac"

# Platform kontrolü — yalnızca macOS.
if [ "$(uname -s)" != "Darwin" ]; then
  echo "Bu araç şu an yalnızca macOS destekler." >&2
  exit 1
fi

echo "SpoofDPI Türkiye kaldırılıyor..."

# Yöntem 1: binary PATH'teyse kendi kendine temizlesin (-y ile onay atlama).
if command -v "$BIN_NAME" &>/dev/null; then
  "$BIN_NAME" uninstall -y
  echo "✓ Kaldırıldı"
  exit 0
fi

# Yöntem 2: Fallback — binary bulunamazsa manuel temizlik.

# LaunchAgent'ı durdur ve kaldır.
PLIST="$HOME/Library/LaunchAgents/com.spoofdpi-tr.plist"
if [ -f "$PLIST" ]; then
  launchctl unload -w "$PLIST" 2>/dev/null || true
  rm -f "$PLIST"
fi

# Tüm ağ servislerinde bizim PAC'imize işaret eden autoproxy'yi kapat.
while IFS= read -r line; do
  # İlk satır başlık; devre dışı servisler '*' ile başlar.
  [[ "$line" == An* ]] && continue
  [[ "$line" == \** ]] && continue
  [ -z "$line" ] && continue
  svc="$line"
  url_output="$(networksetup -getautoproxyurl "$svc" 2>/dev/null || true)"
  if echo "$url_output" | grep -q "$PAC_NAME"; then
    networksetup -setautoproxystate "$svc" off 2>/dev/null || true
  fi
done < <(networksetup -listallnetworkservices 2>/dev/null || true)

# Yapılandırma dizinini sil.
if [ -d "$HOME/.spoofdpi-tr" ]; then
  rm -rf "$HOME/.spoofdpi-tr"
fi

# Binary'leri sil (iki yaygın konum).
rm -f "/usr/local/bin/$BIN_NAME" "$HOME/.local/bin/$BIN_NAME"

echo "✓ Kaldırıldı"
