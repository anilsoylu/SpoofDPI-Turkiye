#!/bin/bash
set -euo pipefail

OWNER="anilsoylu"
REPO="SpoofDPI-Turkiye"
BIN_NAME="spoofdpi-tr"

# Platform kontrolü — yalnızca macOS.
if [ "$(uname -s)" != "Darwin" ]; then
  echo "Bu araç şu an yalnızca macOS destekler." >&2
  exit 1
fi

# Mimari.
case "$(uname -m)" in
  arm64) ARCH="arm64" ;;
  x86_64) ARCH="x86_64" ;;
  *) echo "Desteklenmeyen mimari: $(uname -m)" >&2; exit 1 ;;
esac

# Kurulum dizini (yazılabilir olanı seç).
if [ -w "/usr/local/bin" ]; then
  INSTALL_DIR="/usr/local/bin"
else
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
fi

# En son sürümü bul.
TAG="$(curl -fsSL "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" \
  | grep '"tag_name"' | head -1 | cut -d'"' -f4)"
if [ -z "${TAG:-}" ]; then
  echo "Release bulunamadı. Kaynaktan kurmak için: go install github.com/${OWNER}/${REPO}/cmd/spoofdpi-tr@latest" >&2
  exit 1
fi

ASSET="${BIN_NAME}_${TAG#v}_darwin_${ARCH}.tar.gz"
URL="https://github.com/${OWNER}/${REPO}/releases/download/${TAG}/${ASSET}"

echo "İndiriliyor: ${ASSET}"
TMP="$(mktemp -d)"
curl -fsSL "$URL" -o "$TMP/pkg.tar.gz"
tar -xzf "$TMP/pkg.tar.gz" -C "$TMP"
install -m 0755 "$TMP/${BIN_NAME}" "$INSTALL_DIR/${BIN_NAME}"
rm -rf "$TMP"

echo "✓ ${BIN_NAME} kuruldu: ${INSTALL_DIR}/${BIN_NAME}"

# PATH uyarısı.
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) echo "Not: $INSTALL_DIR PATH'te değil. ~/.zshrc içine ekleyin: export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
esac

# İnteraktif kurulumu başlat (TTY varsa).
if [ -t 0 ]; then
  echo
  "$INSTALL_DIR/${BIN_NAME}" install
else
  echo "Kurulumu tamamlamak için çalıştırın: ${BIN_NAME} install"
fi
