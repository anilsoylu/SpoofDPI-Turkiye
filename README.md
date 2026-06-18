[Türkçe](#türkçe) | [English](#english)

---

## Türkçe

### SpoofDPI Türkiye

Türkiye'de DPI (Derin Paket İnceleme) engellerini aşmak için macOS'a özel, tek komutla kurulabilen bir yönetici CLI.

#### Ne yapar?

[`xvzc/spoofdpi`](https://github.com/xvzc/spoofdpi) adlı resmî açık kaynak motorunu **sarmalayan** (fork değil) bir araçtır. Motor fork'lanmadığı için upstream güvenlik güncellemeleri ve iyileştirmeleri otomatik olarak kullanılabilir hale gelir.

#### Neden?

Eski [`renardozt/SpoofDPI-Turkiye`](https://github.com/renardozt/SpoofDPI-Turkiye) fork'unun iki temel sorunu vardı:

1. **"Tüm internet yavaşlıyordu"** — eski araç sistem proxy'sini genel olarak açıyordu; tüm trafik proxy'den geçtiği için her şey yavaşlıyor ya da kesiliyordu.
2. **"Port seçilemiyordu, Expo ve benzeri araçlarla çakışıyordu"** — port sabitti, değiştirilemiyordu.

Bu araç her iki sorunu da çözer:

- **PAC tabanlı seçici yönlendirme**: Yalnızca seçtiğiniz domainler (Discord vb.) yerel proxy'ye yönlendirilir; geri kalan tüm internet trafiği `DIRECT` gider. Proxy çökse bile yalnızca hedef domainler etkilenir.
- **Yapılandırılabilir port**: Kurulum sırasında istediğiniz portu seçebilirsiniz.
- **Tek komutla kurulum**: Non-developer kullanıcılar için interaktif kurulum sihirbazı.

#### Kurulum (macOS)

Terminale yapıştırın:

```bash
curl -fsSL https://raw.githubusercontent.com/anilsoylu/SpoofDPI-Turkiye/master/install.sh | bash
```

Kurulum size port ve hangi servislerin (Discord vb.) bypass edileceğini sorar. Geri kalan tüm internet trafiğiniz doğrudan (`DIRECT`) gider — yalnızca seçtiğiniz siteler proxy'den geçer.

#### Kullanım

```bash
spoofdpi-tr on              # bypass'ı başlat
spoofdpi-tr off             # durdur (proxy → DIRECT)
spoofdpi-tr status          # durumu gör
spoofdpi-tr add twitch.tv   # domain ekle
spoofdpi-tr remove twitch.tv # domain çıkar
spoofdpi-tr list            # bypass edilen domainleri listele
spoofdpi-tr update          # spoofdpi motorunu güncelle
spoofdpi-tr uninstall       # tamamen kaldır
```

#### Nasıl çalışır?

1. `spoofdpi-tr on` çalıştırıldığında araç bir PAC (Proxy Auto-Config) dosyası oluşturur.
2. PAC dosyası, yalnızca seçilen domainler için `PROXY 127.0.0.1:<port>` döndürür; diğer tüm istekler `DIRECT` ile doğrudan bağlanır.
3. macOS'un ağ ayarlarına bu PAC dosyası `file://` URL'i olarak tanıtılır.
4. Resmî `spoofdpi` binary'si `-system-proxy=false` ile başlatılır (kendi sistem proxy'sini kurmasın; PAC bu işi üstlenir).
5. LaunchAgent ile servis arka planda çalışır ve sistem yeniden başladığında otomatik olarak başlar.

#### Önemli notlar

- Yalnızca **macOS** desteklenmektedir.
- `xvzc/spoofdpi` projesini fork'lamıyoruz; sarmalıyoruz. Tüm DPI bypass mantığı onlara aittir. Bu projenin katkısı PAC seçiciliği, yapılandırma yönetimi ve kullanıcı deneyimidir.
- Upstream'e teşekkürler: [xvzc/SpoofDPI](https://github.com/xvzc/SpoofDPI) — Apache 2.0 lisanslı.

---

## English

### SpoofDPI Turkiye

A macOS-only manager CLI for bypassing DPI (Deep Packet Inspection) blocks in Turkey. Installs with a single command.

#### What does it do?

It **wraps** (not forks) the official [`xvzc/spoofdpi`](https://github.com/xvzc/spoofdpi) engine. Because we don't fork the engine, upstream security patches and improvements are available automatically.

#### Why?

The older [`renardozt/SpoofDPI-Turkiye`](https://github.com/renardozt/SpoofDPI-Turkiye) fork had two core problems:

1. **"The entire internet slowed down"** — the old tool set a global system proxy, routing all traffic through the proxy, causing slowdowns and disconnections.
2. **"The port couldn't be changed, conflicting with Expo and similar tools"** — the port was hardcoded.

This tool solves both:

- **PAC-based selective routing**: Only your chosen domains (Discord etc.) are routed through the local proxy; all other internet traffic goes `DIRECT`. If the proxy crashes, only the target domains are affected.
- **Configurable port**: Choose any port during setup.
- **Single-command installation**: An interactive setup wizard for non-developers.

#### Installation (macOS)

Paste into your terminal:

```bash
curl -fsSL https://raw.githubusercontent.com/anilsoylu/SpoofDPI-Turkiye/master/install.sh | bash
```

The installer asks for your preferred port and which services (Discord etc.) to bypass. All other internet traffic goes directly (`DIRECT`) — only your selected sites go through the proxy.

#### Usage

```bash
spoofdpi-tr on              # start bypass
spoofdpi-tr off             # stop (proxy → DIRECT)
spoofdpi-tr status          # show status
spoofdpi-tr add twitch.tv   # add a domain
spoofdpi-tr remove twitch.tv # remove a domain
spoofdpi-tr list            # list bypassed domains
spoofdpi-tr update          # update the spoofdpi engine
spoofdpi-tr uninstall       # remove everything
```

#### How it works

1. When `spoofdpi-tr on` is run, it generates a PAC (Proxy Auto-Config) file.
2. The PAC file returns `PROXY 127.0.0.1:<port>` only for selected domains; all other requests use `DIRECT`.
3. macOS network settings are pointed at this PAC file via a `file://` URL.
4. The official `spoofdpi` binary is launched with `-system-proxy=false` (so it does not set its own system proxy; the PAC handles routing).
5. A LaunchAgent runs the service in the background and restarts it on login.

#### Notes

- **macOS only.**
- We wrap, not fork, `xvzc/spoofdpi`. All DPI bypass logic is theirs. This project's contribution is PAC-based selectivity, configuration management, and user experience.
- Credit and thanks: [xvzc/SpoofDPI](https://github.com/xvzc/SpoofDPI) — Apache 2.0 licensed.

---

MIT License — Copyright (c) 2026 Anıl Soylu
