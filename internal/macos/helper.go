package macos

import (
	"fmt"
	"strings"
)

// HelperScript, root yetkisiyle çalışan bash helper script metnini üretir.
// /usr/local/libexec/spoofdpi-tr-helper olarak yazılır (root:wheel 0755).
//
// Sudoers kuralı sayesinde bu script parolasız `sudo` ile çağrılabilir; tüm
// ayrıcalıklı işlemler (pfctl, launchctl, anchor + plist dosyası yazımı) burada
// yapılır.
//
// tpwsBin ve hostlist install anında gömülür (kurulum kullanıcısına göre sabit
// yollar). Böylece helper, plist'i HER `start <port>` çağrısında GÜNCEL portla
// kendisi yeniden üretebilir (BUG1 fix): port değişince anchor da plist de aynı
// porta hizalanır ve daemon taze başlar.
//
// Komutlar:
//
//	start <tpwsPort>               : anchor + plist dosyalarını GÜNCEL portla yaz,
//	                                 pf'yi yükle+aktif et, anchor kurallarını
//	                                 yükle, tpws daemon'u (yeniden) başlat ve
//	                                 portu gerçekten dinleyene kadar doğrula.
//	                                 Domainler tpws tarafından --hostlist
//	                                 dosyasından okunur; daemon taze başladığında
//	                                 güncel hostlist'i yeniden okur.
//	stop                           : tpws daemon'u durdur, anchor kurallarını boşalt.
//	status                         : tpws çalışıyor mu + anchor dolu mu.
func HelperScript(tpwsBin, hostlist string) string {
	// %[n]$s ile alan adlarını sabitlere bağlıyoruz; bash içindeki $1/$2 gibi
	// değişkenler Go format dizisinde %% kaçışıyla korunur.
	return fmt.Sprintf(`#!/bin/bash
# spoofdpi-tr root helper — sudoers ile parolasız çağrılır.
# Tüm ayrıcalıklı PF/launchctl işlemleri burada yapılır.
set -euo pipefail

ANCHOR_NAME=%[1]q
ANCHOR_PATH=%[2]q
PF_CONF=%[3]q
PLIST=%[4]q
DAEMON_ID=%[5]q
TPWS_BIN=%[6]q
HOSTLIST=%[7]q

write_anchor() {
  local port="$1"
  cat > "$ANCHOR_PATH" <<EOF
rdr on lo0 inet proto tcp from { !127.0.0.0/8 !192.168.0.0/16 !10.0.0.0/8 !172.16.0.0/12 } to any port { 443 } -> 127.0.0.1 port ${port}
pass out route-to (lo0 127.0.0.1) inet proto tcp from { !127.0.0.0/8 !192.168.0.0/16 !10.0.0.0/8 !172.16.0.0/12 } to any port { 443 } user { >root }
EOF
  chmod 644 "$ANCHOR_PATH"
}

# write_plist, LaunchDaemon plist'ini GÜNCEL portla yeniden üretir (BUG1 fix).
# tpws bin + hostlist yolu install anında gömülüdür; yalnızca port değişkendir.
write_plist() {
  local port="$1"
  cat > "$PLIST" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>${DAEMON_ID}</string>
	<key>ProgramArguments</key>
	<array>
		<string>${TPWS_BIN}</string>
		<string>--user=root</string>
		<string>--port=${port}</string>
		<string>--bind-addr=127.0.0.1</string>
		<string>--hostlist=${HOSTLIST}</string>
		<string>--tlsrec=sni</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<dict>
		<key>SuccessfulExit</key>
		<false/>
	</dict>
	<key>ThrottleInterval</key>
	<integer>5</integer>
	<key>StandardOutPath</key>
	<string>/tmp/spoofdpi-tr-tpws.out.log</string>
	<key>StandardErrorPath</key>
	<string>/tmp/spoofdpi-tr-tpws.err.log</string>
</dict>
</plist>
EOF
  chown root:wheel "$PLIST" 2>/dev/null || true
  chmod 644 "$PLIST"
}

# port_listener_pids, 127.0.0.1:<port>'u DİNLEYEN PID'leri (varsa) yazdırır.
port_listener_pids() {
  lsof -nP -iTCP:"$1" -sTCP:LISTEN -t 2>/dev/null || true
}

# daemon_start, eski örneği TAMAMEN temizleyip (port boşalana kadar bekleyerek)
# plist'i taze bootstrap eder ve YENİ daemon'ın gerçekten dinlediğini doğrular.
#
# #4: bootout asenkron tamamlanabilir ve eski tpws portu hâlâ tutuyor olabilir;
# bu durumda hemen bootstrap edilirse verify_listen ESKİ process'i "başladı"
# sanır (yanlış pozitif). Bu yüzden bootstrap'tan ÖNCE port BOŞALANA kadar
# bekleriz; takılı kalan stale tpws varsa öldürürüz.
# #7: sabit süre yerine port-boşaldı sinyaliyle beklenir (yük-bağımsız, ~max 10sn).
daemon_start() {
  local port="$1" i pids
  if launchctl print "system/${DAEMON_ID}" >/dev/null 2>&1; then
    launchctl bootout "system/${DAEMON_ID}" >/dev/null 2>&1 || true
  fi
  # bootout'un oturmasını VE portun boşalmasını bekle (yük-bağımsız, ~max 10sn).
  for i in $(seq 1 20); do
    if ! launchctl print "system/${DAEMON_ID}" >/dev/null 2>&1 \
       && [ -z "$(port_listener_pids "$port")" ]; then
      break
    fi
    sleep 0.5
  done
  # Hâlâ portu tutan stale dinleyici kaldıysa zorla öldür (eski tpws), sonra
  # portun gerçekten boşalmasını kısa süre bekle.
  pids="$(port_listener_pids "$port")"
  if [ -n "$pids" ]; then
    kill -TERM $pids 2>/dev/null || true
    for i in 1 2 3 4 5 6; do
      [ -z "$(port_listener_pids "$port")" ] && break
      sleep 0.5
    done
    pids="$(port_listener_pids "$port")"
    if [ -n "$pids" ]; then
      kill -KILL $pids 2>/dev/null || true
      sleep 0.5
    fi
  fi
  launchctl enable "system/${DAEMON_ID}" >/dev/null 2>&1 || true
  launchctl bootstrap system "$PLIST"
}

daemon_stop() {
  launchctl bootout "system/${DAEMON_ID}" >/dev/null 2>&1 || \
    launchctl unload -w "$PLIST" >/dev/null 2>&1 || true
}

# verify_listen, YENİ tpws daemon'ının verilen portu gerçekten DİNLEDİĞİNİ
# doğrular (#4: yanlış pozitif önleme). Yalnızca "biri dinliyor" yetmez; dinleyen
# PID, launchctl'in bizim DAEMON_ID için bildirdiği PID ile EŞLEŞMELİDİR. Böylece
# bootout'tan sağ kalan eski bir process "yeni daemon başladı" sanılmaz.
# KeepAlive ile daemon yüklü görünse bile süreç çökebilir; ~max 5sn dene.
verify_listen() {
  local port="$1" i lpids dpid
  for i in $(seq 1 10); do
    lpids="$(port_listener_pids "$port")"
    if [ -n "$lpids" ]; then
      # launchctl'in bildirdiği PID dinleyiciler arasında mı?
      dpid="$(launchctl print "system/${DAEMON_ID}" 2>/dev/null \
        | sed -n 's/.*[[:space:]]pid = \([0-9][0-9]*\).*/\1/p' | head -n1 || true)"
      if [ -n "$dpid" ]; then
        for p in $lpids; do
          [ "$p" = "$dpid" ] && return 0
        done
      else
        # launchctl PID raporlamadıysa (bazı sürümler), en azından port dinleniyor.
        return 0
      fi
    fi
    sleep 0.5
  done
  return 1
}

cmd="${1:-}"
case "$cmd" in
  start)
    port="${2:-}"
    # Domainler tpws tarafından --hostlist dosyasından okunur; helper yalnızca
    # PF + daemon yaşam döngüsünü yönetir. daemon_start her zaman bootout+
    # bootstrap yaptığından tpws taze başlar ve güncel hostlist'i yeniden okur.
    if [ -z "$port" ]; then echo "kullanim: start <port>" >&2; exit 2; fi
    # GÜVENLİK: sudoers kuralı 'helper *' olduğundan kullanıcı bu root script'i
    # doğrudan keyfi argümanla çağırabilir; Go tarafının int doğrulamasına
    # GÜVENME. port'u burada, ayrıcalık sınırında doğrula: yalnızca 1-65535
    # aralığında ondalık sayı. Aksi halde reddet (enjeksiyon/bozuk PF savunması).
    case "$port" in
      ''|*[!0-9]*) echo "gecersiz port: $port" >&2; exit 2 ;;
    esac
    if [ "$port" -lt 1 ] || [ "$port" -gt 65535 ]; then
      echo "port 1-65535 araliginda olmali: $port" >&2; exit 2
    fi
    write_anchor "$port"
    write_plist "$port"
    pfctl -f "$PF_CONF"
    pfctl -e 2>/dev/null || true
    pfctl -a "$ANCHOR_NAME" -f "$ANCHOR_PATH"
    daemon_start "$port"
    if ! verify_listen "$port"; then
      echo "HATA: tpws ${port} portunu dinlemiyor (daemon baslatilamadi)." >&2
      echo "Teshis: launchctl print system/${DAEMON_ID} ve /tmp/spoofdpi-tr-tpws.err.log" >&2
      exit 1
    fi
    echo "spoofdpi-tr basladi (port ${port})"
    ;;
  stop)
    daemon_stop
    pfctl -a "$ANCHOR_NAME" -F all 2>/dev/null || true
    echo "spoofdpi-tr durdu"
    ;;
  status)
    # "calisiyor" yalnızca daemon yüklü VE port dinleniyorsa raporlanır; böylece
    # çökmüş ama yüklü görünen daemon "calisiyor" demez (BUG2 doğrulaması).
    listening=0
    if lsof -nP -iTCP -sTCP:LISTEN 2>/dev/null | grep -q tpws; then
      listening=1
    fi
    if launchctl print "system/${DAEMON_ID}" >/dev/null 2>&1 && [ "$listening" -eq 1 ]; then
      echo "tpws: calisiyor"
    else
      echo "tpws: durdu"
    fi
    if pfctl -a "$ANCHOR_NAME" -s nat 2>/dev/null | grep -q rdr; then
      echo "pf-anchor: dolu"
    else
      echo "pf-anchor: bos"
    fi
    ;;
  *)
    echo "kullanim: spoofdpi-tr-helper start|stop|status" >&2
    exit 2
    ;;
esac
`, AnchorName, AnchorPath, PFConfPath, LaunchDaemonPath, LaunchDaemonID, tpwsBin, hostlist)
}

// Plist üretimi TEK KAYNAK olarak helper'ın write_plist fonksiyonundadır (#2);
// Go tarafı plist üretmez. Bkz. HelperScript içindeki write_plist.

// SudoersRule, helper'ın parolasız çağrılmasını sağlayan sudoers satırını üretir.
// /etc/sudoers.d/spoofdpi-tr olarak yazılır (chmod 440, visudo -cf ile doğrula).
func SudoersRule(user string) string {
	return fmt.Sprintf("%s ALL=(root) NOPASSWD: %s *\n", user, HelperPath)
}

// xmlEscape, plist string'leri için minimal XML kaçışı yapar.
func xmlEscape(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
	return r.Replace(s)
}
