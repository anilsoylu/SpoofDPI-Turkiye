#!/bin/bash
set -euo pipefail

OWNER="anilsoylu"
REPO="SpoofDPI-Turkiye"
BIN_NAME="spoofdpi-tr"

# tpws motoru (bol-van/zapret) kaynaktan derlenir — indirilmiş binary YOK.
ZAPRET_REPO="https://github.com/bol-van/zapret"
SPOOF_DIR="$HOME/.spoofdpi-tr"
TPWS_DIR="$SPOOF_DIR/bin"
TPWS_BIN="$TPWS_DIR/tpws"

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

# ---------------------------------------------------------------------------
# 1) Yönetici (manager) binary'sini edin: önce release, yoksa kaynaktan derle.
# ---------------------------------------------------------------------------
install_manager_from_release() {
  TAG="$(curl -fsSL "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" \
    | grep '"tag_name"' | head -1 | cut -d'"' -f4 || true)"
  if [ -z "${TAG:-}" ]; then
    return 1
  fi

  ASSET="${BIN_NAME}_${TAG#v}_darwin_${ARCH}.tar.gz"
  URL="https://github.com/${OWNER}/${REPO}/releases/download/${TAG}/${ASSET}"

  echo "İndiriliyor: ${ASSET}"
  local tmp
  tmp="$(mktemp -d)"
  if ! curl -fsSL "$URL" -o "$tmp/pkg.tar.gz"; then
    rm -rf "$tmp"
    return 1
  fi
  tar -xzf "$tmp/pkg.tar.gz" -C "$tmp"
  install -m 0755 "$tmp/${BIN_NAME}" "$INSTALL_DIR/${BIN_NAME}"
  rm -rf "$tmp"
  echo "✓ ${BIN_NAME} kuruldu (release): ${INSTALL_DIR}/${BIN_NAME}"
}

install_manager_from_source() {
  if ! command -v go &>/dev/null; then
    echo "Go bulunamadı. Önce kurun: brew install go" >&2
    echo "Sonra bu betiği tekrar çalıştırın." >&2
    exit 1
  fi
  echo "Kaynaktan derleniyor: ${BIN_NAME} (go install)..."
  GOBIN="$INSTALL_DIR" go install "github.com/${OWNER}/${REPO}/cmd/${BIN_NAME}@latest"
  echo "✓ ${BIN_NAME} kuruldu (kaynak): ${INSTALL_DIR}/${BIN_NAME}"
}

if ! install_manager_from_release; then
  echo "Release bulunamadı — kaynaktan derlemeye geçiliyor."
  install_manager_from_source
fi

# PATH uyarısı.
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) echo "Not: $INSTALL_DIR PATH'te değil. ~/.zshrc içine ekleyin: export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
esac

# ---------------------------------------------------------------------------
# 2) tpws motorunu KAYNAKTAN derle (bol-van/zapret -> make mac).
#    İndirilmiş binary kullanılmaz (güvenlik: denetlenebilir, kaynaktan derleme).
# ---------------------------------------------------------------------------
build_tpws_from_source() {
  if [ -x "$TPWS_BIN" ]; then
    echo "✓ tpws zaten kurulu: $TPWS_BIN"
    return 0
  fi

  # Derleyici (cc/clang) gerekir — Xcode Command Line Tools.
  if ! command -v cc &>/dev/null && ! command -v clang &>/dev/null; then
    echo "C derleyicisi (cc/clang) bulunamadı." >&2
    echo "Xcode Command Line Tools kurun: xcode-select --install" >&2
    exit 1
  fi
  if ! command -v git &>/dev/null; then
    echo "git bulunamadı. Xcode Command Line Tools kurun: xcode-select --install" >&2
    exit 1
  fi
  if ! command -v make &>/dev/null; then
    echo "make bulunamadı. Xcode Command Line Tools kurun: xcode-select --install" >&2
    exit 1
  fi

  echo "tpws motoru kaynaktan derleniyor (bol-van/zapret)..."
  local tmp
  tmp="$(mktemp -d)"
  git clone --depth 1 "$ZAPRET_REPO" "$tmp/zapret"
  ( cd "$tmp/zapret" && make mac )

  local built="$tmp/zapret/binaries/my/tpws"
  if [ ! -f "$built" ]; then
    echo "tpws derlenemedi: $built bulunamadı" >&2
    rm -rf "$tmp"
    exit 1
  fi

  mkdir -p "$TPWS_DIR"
  install -m 0755 "$built" "$TPWS_BIN"
  rm -rf "$tmp"
  echo "✓ tpws kuruldu: $TPWS_BIN"
}

build_tpws_from_source

# ---------------------------------------------------------------------------
# 3) İnteraktif kurulumu başlat (tek seferlik admin onayı içerir).
# ---------------------------------------------------------------------------
if [ -t 0 ]; then
  echo
  "$INSTALL_DIR/${BIN_NAME}" install
else
  echo "Kurulumu tamamlamak için çalıştırın: ${BIN_NAME} install"
fi
