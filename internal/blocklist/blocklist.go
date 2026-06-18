// Package blocklist, Türkiye'de DPI'dan etkilendiği bilinen servisler için
// küratörlü, kategorize edilmiş varsayılan alan adı listelerini sağlar.
package blocklist

// Category, ilişkili bir servis grubudur (CLI'de seçilebilir birim).
type Category struct {
	// Key, makine-okunur kimlik (örn. "discord").
	Key string
	// Title, kullanıcıya gösterilen ad (örn. "Discord").
	Title string
	// Domains, bu kategori için bypass edilecek alan adları.
	Domains []string
}

// Categories, varsayılan küratörlü kategorileri döndürür.
// Bu liste topluluk katkılarıyla zamanla güncellenmelidir.
func Categories() []Category {
	return []Category{
		{
			Key:   "discord",
			Title: "Discord",
			// Düz domain yazmak yeterli: PAC her domaini + tüm alt alanlarını
			// (cdn.discordapp.com, gateway.discord.gg vb.) otomatik kapsar.
			Domains: []string{
				"discord.com",
				"discordapp.com",
				"discord.gg",
				"discordapp.net",
				"discord.media",
				"discordcdn.com",
				"discord.dev",
				"discordstatus.com",
			},
		},
		// İleride topluluk katkısıyla genişletilebilecek diğer kategoriler.
		// Yer tutucu örnekler kasıtlı olarak EKLENMEDİ; yalnızca doğrulanmış,
		// gerçekten etkilenen servisler eklenmeli (yanlış pozitif kullanıcı
		// trafiğini gereksiz yere proxy'ye sokar).
	}
}

// Get, anahtara göre kategoriyi döndürür; bulunamazsa ok=false.
func Get(key string) (Category, bool) {
	for _, c := range Categories() {
		if c.Key == key {
			return c, true
		}
	}
	return Category{}, false
}

// DefaultDomains, ilk kurulum için önerilen varsayılan domain kümesini döndürür
// (şimdilik yalnızca Discord — en bilinen vaka).
func DefaultDomains() []string {
	d, _ := Get("discord")
	return d.Domains
}
