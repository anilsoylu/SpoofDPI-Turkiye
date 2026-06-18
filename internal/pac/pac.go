// Package pac, bir alan adı listesinden Proxy Auto-Config (PAC) dosyası üretir.
// Yalnızca eşleşen domainler proxy'ye, geri kalan her şey DIRECT'e yönlendirilir.
package pac

import (
	"fmt"
	"sort"
	"strings"
)

// validDomain, basit bir host adı doğrulaması yapar (PAC enjeksiyonunu önler).
// Yalnızca harf, rakam, nokta ve tire kabul edilir; en az bir nokta gerekir.
func validDomain(d string) bool {
	d = strings.TrimSpace(d)
	if d == "" || !strings.Contains(d, ".") {
		return false
	}
	for _, r := range d {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '.' || r == '-':
		default:
			return false
		}
	}
	// Baş/son nokta veya tire kabul edilmez.
	if strings.HasPrefix(d, ".") || strings.HasPrefix(d, "-") ||
		strings.HasSuffix(d, ".") || strings.HasSuffix(d, "-") {
		return false
	}
	return true
}

// Normalize, listeyi temizler: trim, küçük harf, geçersizleri at, tekilleştir, sırala.
func Normalize(domains []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, d := range domains {
		d = strings.ToLower(strings.TrimSpace(d))
		if !validDomain(d) || seen[d] {
			continue
		}
		seen[d] = true
		out = append(out, d)
	}
	sort.Strings(out)
	return out
}

// Generate, verilen port ve domainler için PAC dosyası metnini döndürür.
// Domainler Normalize edilmiş kabul edilir; etmemişse içeride normalize edilir.
func Generate(port int, domains []string) string {
	clean := Normalize(domains)

	var b strings.Builder
	b.WriteString("// spoofdpi-tr tarafından üretildi — elle düzenlemeyin\n")
	b.WriteString("function FindProxyForURL(url, host) {\n")
	b.WriteString(fmt.Sprintf("  var proxy = \"PROXY 127.0.0.1:%d\";\n", port))
	b.WriteString("  host = host.toLowerCase();\n")
	if len(clean) == 0 {
		b.WriteString("  return \"DIRECT\";\n")
		b.WriteString("}\n")
		return b.String()
	}
	b.WriteString("  var domains = [\n")
	for _, d := range clean {
		b.WriteString(fmt.Sprintf("    %q,\n", d))
	}
	b.WriteString("  ];\n")
	b.WriteString("  for (var i = 0; i < domains.length; i++) {\n")
	b.WriteString("    var d = domains[i];\n")
	b.WriteString("    if (host === d || dnsDomainIs(host, \".\" + d)) {\n")
	b.WriteString("      return proxy;\n")
	b.WriteString("    }\n")
	b.WriteString("  }\n")
	b.WriteString("  return \"DIRECT\";\n")
	b.WriteString("}\n")
	return b.String()
}
