import Foundation

// MARK: - Dil

enum Lang: String, CaseIterable {
    case tr = "TR"
    case en = "EN"
}

// MARK: - Çeviri sözlüğü

private let translations: [String: [Lang: String]] = [
    // Durum
    "status.protected":         [.tr: "Korunuyor",             .en: "Protected"],
    "status.off":               [.tr: "Kapalı",                .en: "Off"],
    "status.sub.on":            [.tr: "spoofdpi etkin",        .en: "spoofdpi active"],
    "status.sub.domains":       [.tr: "domain filtreleniyor",  .en: "domains filtered"],
    "status.sub.off":           [.tr: "koruma kapalı",         .en: "protection off"],
    // Butonlar
    "btn.stop":                 [.tr: "Korumayı Durdur",       .en: "Stop Protection"],
    "btn.start":                [.tr: "Korumayı Başlat",       .en: "Start Protection"],
    "btn.restart":              [.tr: "Yeniden Başlat",        .en: "Restart"],
    "btn.refresh":              [.tr: "Yenile",                .en: "Refresh"],
    "btn.uninstall":            [.tr: "SpoofDPI'ı Kaldır",     .en: "Uninstall SpoofDPI"],
    "btn.saveapply":            [.tr: "Kaydet ve Uygula",      .en: "Save & Apply"],
    "btn.discord":              [.tr: "Discord profili",       .en: "Discord profile"],
    "btn.runtest":              [.tr: "Test Et",               .en: "Run Test"],
    // Bilgi satırları
    "info.autostart":           [.tr: "Açılışta otomatik başlar", .en: "Starts at login"],
    "info.processing":          [.tr: "İşleniyor...",          .en: "Processing..."],
    // Kart başlıkları
    "card.domains.title":       [.tr: "HEDEF ALAN ADLARI",    .en: "TARGET DOMAINS"],
    "card.domains.hint":        [.tr: "Her satıra bir kök alan adı yazın. Alt alan adları otomatik kapsanır.",
                                 .en: "Enter one root domain per line. Subdomains are automatically included."],
    "card.test.title":          [.tr: "BAĞLANTI TESTİ",       .en: "CONNECTION TEST"],
    // Menu bar
    "menu.open":                [.tr: "Pencereyi Aç",          .en: "Open Window"],
    "menu.quit":                [.tr: "Çıkış",                 .en: "Quit"],
    // Wordmark alt yazı
    "wordmark.sub":             [.tr: "TÜRKİYE",              .en: "TURKEY"],
    // Onay
    "uninstall.confirm":        [.tr: "SpoofDPI Türkiye'yi tamamen kaldırmak istediğinizden emin misiniz?",
                                 .en: "Are you sure you want to completely uninstall SpoofDPI Türkiye?"],
    "uninstall.confirm.btn":    [.tr: "Kaldır",               .en: "Uninstall"],
    "cancel":                   [.tr: "İptal",                 .en: "Cancel"],
]

// MARK: - Çeviri fonksiyonu

func t(_ key: String, lang: Lang) -> String {
    translations[key]?[lang] ?? key
}

// MARK: - Namespace (AppState.Localization.t erişimi için)

enum Localization {
    static func t(_ key: String, lang: Lang) -> String {
        translations[key]?[lang] ?? key
    }
}
