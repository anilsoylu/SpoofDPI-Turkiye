[Türkçe](#türkçe) | [English](#english)

---

## Türkçe

### SpoofDPI Türkiye

Türkiye'de DPI (Derin Paket İnceleme) engellerini aşmak için macOS'a özel, tek komutla kurulabilen bir yönetici CLI.

#### Ne yapar?

macOS'ta DPI tabanlı engelleri aşar. Yalnızca tarayıcıda değil, **masaüstü uygulamalarında da (Discord masaüstü dahil)** çalışır. Bunu macOS PF (Packet Filter) ile **transparan yönlendirme** sayesinde yapar: 443 (HTTPS) trafiği çekirdek seviyesinde yerel motora yönlendirilir, böylece sistem proxy'sini dinlemeyen uygulamalar bile kapsanır.

Yalnızca seçtiğiniz domainler etkilenir (örn. Discord); geri kalan tüm trafik dokunulmadan, doğrudan akar.

#### Nasıl çalışır?

- **Motor: tpws** — açık kaynak ([bol-van/zapret](https://github.com/bol-van/zapret)) bir DPI-bypass motoru. Kurulum sırasında **kaynaktan derlenir** (`git clone --depth 1 zapret` → `make mac`) ve `~/.spoofdpi-tr/bin/tpws` altına yerleştirilir. İndirilmiş hazır binary kullanılmaz.
- **macOS PF transparan redirect** — PF, 443 trafiğini tpws'in dinlediği porta (varsayılan **988**) yönlendirir. TLS şifresi çözülmez; TLS `ClientHello` bölme noktalarından parçalanır ve out-of-band bayt ile gönderilir (`--split-pos=1,midsld --oob=tls`) — böylece DPI, SNI'yi yeniden birleştirip göremez.
- **Seçici** — yalnızca hostlist'teki domainler desync edilir; gerisi pass-through (etkilenmez).
- **Tek seferlik admin** — kurulumda bir kez yönetici onayı (parolasız helper + sudoers) verilir; sonrasında `on` / `off` parolasız çalışır.
- **Port çakışmaz** — varsayılan port 988, Expo ve benzeri geliştirici araçlarıyla çakışmaz; istediğiniz portu seçebilirsiniz.

#### Kurulum (macOS)

Terminale yapıştırın:

```bash
curl -fsSL https://raw.githubusercontent.com/anilsoylu/SpoofDPI-Turkiye/master/install.sh | bash
```

Kurulum betiği şunları yapar:
1. Yönetici (manager) binary'sini indirir (release varsa) ya da kaynaktan derler (Go gerektirir — yoksa `brew install go` uyarısı verir).
2. **tpws motorunu kaynaktan derler** (Xcode Command Line Tools / `cc` gerektirir — yoksa `xcode-select --install` uyarısı verir).
3. İnteraktif kurulumu başlatır: port ve hangi servislerin (Discord vb.) bypass edileceğini sorar, tek seferlik yönetici onayı alır.

#### Kullanım

```bash
spoofdpi-tr on               # bypass'ı başlat (parolasız)
spoofdpi-tr off              # durdur — 443 trafiği doğrudan akar (parolasız)
spoofdpi-tr status           # servis ve yapılandırma durumunu gör
spoofdpi-tr add discord.com  # bypass listesine domain ekle
spoofdpi-tr remove discord.com # listeden domain çıkar
spoofdpi-tr set a.com b.com  # listeyi verilen domainlerle tamamen değiştir
spoofdpi-tr list             # bypass edilen domainleri listele
spoofdpi-tr port 9090        # tpws redirect portunu değiştir (1-65535)
spoofdpi-tr uninstall        # tamamen kaldır
```

#### Menü Çubuğu Uygulaması

Proje, terminal sevmeyenler için native bir **macOS menü çubuğu uygulaması** (SwiftUI) içerir.

- **Ne yapar:** Menü çubuğunda (üst barda) yaşar, **dock ikonu yoktur**. Tıklayınca: koruma durumu (Korunuyorsunuz/Kapalı), **aç/kapa anahtarı**, açılır-kapanır alan adı düzenleyici (text area + Discord profili + Kaydet ve Uygula), bağlantı testi (Discord/OpenAI/Anthropic/GitHub), ayarlar (port, dil TR/EN, kaldır). **Light/Dark moda otomatik uyum**, native macOS görünümü.
- **Nasıl derlenir/çalıştırılır:**
  ```bash
  cd menubar
  ./build-app.sh
  open "SpoofDPI-Türkiye.app"
  ```
  (Go + Xcode Command Line Tools gerekir. CLI bundle'a gömülüdür; ama sistem kurulumu için yine de `spoofdpi-tr install` veya `curl|bash` gerekir — GUI on/off/durum için CLI'ı kullanır.)
- NOT: GUI, kurulu CLI'a (`spoofdpi-tr`) shell-out eder; önce CLI kurulu olmalı (install.sh).

#### Kaldırma (macOS)

Tek satır (binary kurulu olmasa da çalışır):

```bash
curl -fsSL https://raw.githubusercontent.com/anilsoylu/SpoofDPI-Turkiye/master/uninstall.sh | bash
```

Binary kuruluysa alternatif:

```bash
spoofdpi-tr uninstall
```

Kaldırma; helper'ı durdurur, `pf.conf`'u eski haline getirir ve sudoers / helper / anchor / LaunchDaemon dosyalarıyla `~/.spoofdpi-tr` dizinini (tpws binary dahil) siler.

#### Güvenlik notları

- **tpws kaynaktan derlenir** — hazır binary indirilmez; derlenen kaynak denetlenebilirdir (açık kaynak, [bol-van/zapret](https://github.com/bol-van/zapret)).
- **TLS şifresi çözülmez** — araç trafiğinizin içeriğini göremez; yalnızca TLS el sıkışmasındaki SNI kaydını yeniden düzenler.
- **Root neden gerekli?** — 443 trafiğini çekirdek seviyesinde yönlendirebilmek için PF/transparan redirect, root ayrıcalığı ister. Bu nedenle kurulumda bir kez yönetici onayı alınır; sonrası parolasız helper ile yürür.

#### Önemli notlar

- Yalnızca **macOS** desteklenmektedir (Apple Silicon + Intel).
- Teşekkürler: [bol-van/zapret](https://github.com/bol-van/zapret) — tpws DPI-bypass motoru.

---

## English

### SpoofDPI Turkiye

A macOS-only manager CLI for bypassing DPI (Deep Packet Inspection) blocks in Turkey. Installs with a single command.

#### What does it do?

It bypasses DPI-based blocks on macOS. It works not only in the **browser** but also in **desktop applications (including the Discord desktop app)**. It achieves this via **transparent redirection** with macOS PF (Packet Filter): port 443 (HTTPS) traffic is redirected to the local engine at the kernel level, so even apps that do not honor the system proxy are covered.

Only the domains you choose are affected (e.g. Discord); all other traffic flows directly, untouched.

#### How it works

- **Engine: tpws** — an open-source DPI-bypass engine ([bol-van/zapret](https://github.com/bol-van/zapret)). It is **built from source** during installation (`git clone --depth 1 zapret` → `make mac`) and placed at `~/.spoofdpi-tr/bin/tpws`. No prebuilt binary is downloaded.
- **macOS PF transparent redirect** — PF redirects port 443 traffic to the port tpws listens on (default **988**). TLS is never decrypted; the TLS `ClientHello` is fragmented at split points and sent with an out-of-band byte (`--split-pos=1,midsld --oob=tls`) so the DPI cannot reassemble the SNI.
- **Selective** — only domains in the hostlist are desynced; everything else passes through unaffected.
- **One-time admin** — installation asks for admin approval once (passwordless helper + sudoers); afterwards `on` / `off` run without a password.
- **No port conflicts** — the default port 988 does not collide with Expo and similar developer tools; you can pick any port.

#### Installation (macOS)

Paste into your terminal:

```bash
curl -fsSL https://raw.githubusercontent.com/anilsoylu/SpoofDPI-Turkiye/master/install.sh | bash
```

The installer:
1. Downloads the manager binary (if a release exists) or builds it from source (requires Go — otherwise it tells you to `brew install go`).
2. **Builds the tpws engine from source** (requires Xcode Command Line Tools / `cc` — otherwise it tells you to run `xcode-select --install`).
3. Launches the interactive setup: asks for your port and which services (Discord etc.) to bypass, and takes a one-time admin approval.

#### Usage

```bash
spoofdpi-tr on               # start bypass (passwordless)
spoofdpi-tr off              # stop — 443 traffic flows directly (passwordless)
spoofdpi-tr status           # show service and config status
spoofdpi-tr add discord.com  # add a domain to the bypass list
spoofdpi-tr remove discord.com # remove a domain from the list
spoofdpi-tr set a.com b.com  # replace the list entirely with given domains
spoofdpi-tr list             # list bypassed domains
spoofdpi-tr port 9090        # change the tpws redirect port (1-65535)
spoofdpi-tr uninstall        # remove everything
```

#### Menu Bar App

The project includes a native **macOS menu bar app** (SwiftUI) for those who do not like the terminal.

- **What it does:** It lives in the menu bar (top bar) and has **no dock icon**. When clicked: protection status (Protected/Off), an **on/off toggle**, a collapsible domain editor (text area + Discord profile + Save and Apply), connection test (Discord/OpenAI/Anthropic/GitHub), settings (port, language TR/EN, uninstall). **Automatic Light/Dark mode adaptation**, native macOS look.
- **How to build/run:**
  ```bash
  cd menubar
  ./build-app.sh
  open "SpoofDPI-Türkiye.app"
  ```
  (Requires Go + Xcode Command Line Tools. The CLI is embedded in the bundle; but system setup still needs `spoofdpi-tr install` or `curl|bash` — the GUI uses the CLI for on/off/status.)
- NOTE: The GUI shells out to the installed CLI (`spoofdpi-tr`); the CLI must be installed first (install.sh).

#### Uninstallation (macOS)

One-liner (works even if the binary is not in PATH):

```bash
curl -fsSL https://raw.githubusercontent.com/anilsoylu/SpoofDPI-Turkiye/master/uninstall.sh | bash
```

If the binary is already installed, you can also run:

```bash
spoofdpi-tr uninstall
```

Uninstall stops the helper, restores `pf.conf`, and removes the sudoers / helper / anchor / LaunchDaemon files along with the `~/.spoofdpi-tr` directory (including the tpws binary).

#### Security notes

- **tpws is built from source** — no prebuilt binary is downloaded; the compiled source is auditable (open source, [bol-van/zapret](https://github.com/bol-van/zapret)).
- **TLS is not decrypted** — the tool cannot see the contents of your traffic; it only rewrites the SNI record in the TLS handshake.
- **Why root is required** — redirecting port 443 traffic at the kernel level (PF / transparent redirect) requires root privileges. That is why installation takes a one-time admin approval; everything after runs via a passwordless helper.

#### Notes

- **macOS only** (Apple Silicon + Intel).
- Credit and thanks: [bol-van/zapret](https://github.com/bol-van/zapret) — the tpws DPI-bypass engine.

---

MIT License — Copyright (c) 2026 Anıl Soylu
