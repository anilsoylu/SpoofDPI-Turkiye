import Foundation

/// spoofdpi-tr binary'sini bulup Process ile çalıştıran katman.
enum CLI {

    /// Bilinen konumlarda spoofdpi-tr binary'sini arar; bulursa tam yolunu döndürür.
    static func findBinary() -> String? {
        // 0. Uygulama paketine gömülü CLI (kendi kendine yeten kurulum) — önce buna bak.
        if let res = Bundle.main.resourceURL?.appendingPathComponent("spoofdpi-tr"),
           FileManager.default.isExecutableFile(atPath: res.path) {
            return res.path
        }
        // 1. PATH env değişkeninden ara.
        if let pathEnv = ProcessInfo.processInfo.environment["PATH"] {
            for dir in pathEnv.split(separator: ":") {
                let candidate = "\(dir)/spoofdpi-tr"
                if FileManager.default.isExecutableFile(atPath: candidate) {
                    return candidate
                }
            }
        }
        // 2. Sabit konumlar.
        let fixed = [
            "/usr/local/bin/spoofdpi-tr",
            "/opt/homebrew/bin/spoofdpi-tr",
            (NSHomeDirectory() as NSString).appendingPathComponent(".local/bin/spoofdpi-tr")
        ]
        for path in fixed {
            if FileManager.default.isExecutableFile(atPath: path) {
                return path
            }
        }
        return nil
    }

    /// Binary'yi verilen argümanlarla çalıştırır; stdout+stderr toplar.
    /// - Returns: (out: birleşik çıktı, ok: exit code 0 ise true)
    @discardableResult
    static func run(_ args: [String]) -> (out: String, ok: Bool) {
        guard let bin = findBinary() else {
            return ("spoofdpi-tr bulunamadı", false)
        }
        let process = Process()
        process.executableURL = URL(fileURLWithPath: bin)
        process.arguments = args

        let pipe = Pipe()
        process.standardOutput = pipe
        process.standardError = pipe

        do {
            try process.run()
        } catch {
            return ("Process başlatılamadı: \(error.localizedDescription)", false)
        }
        let data = pipe.fileHandleForReading.readDataToEndOfFile()
        process.waitUntilExit()

        let out = String(data: data, encoding: .utf8) ?? ""
        let ok = process.terminationStatus == 0
        return (out, ok)
    }

    // MARK: - Yardımcı metodlar

    @discardableResult
    static func on() -> (out: String, ok: Bool) { run(["on"]) }

    @discardableResult
    static func off() -> (out: String, ok: Bool) { run(["off"]) }

    @discardableResult
    static func add(_ domain: String) -> (out: String, ok: Bool) { run(["add", domain]) }

    @discardableResult
    static func remove(_ domain: String) -> (out: String, ok: Bool) { run(["remove", domain]) }

    @discardableResult
    static func setPort(_ p: Int) -> (out: String, ok: Bool) { run(["port", "\(p)"]) }

    /// `spoofdpi-tr status` çıktısını döndürür (durum tespiti için).
    @discardableResult
    static func status() -> (out: String, ok: Bool) { run(["status"]) }

    /// Domain listesini tamamen verilen liste ile değiştirir.
    @discardableResult
    static func set(_ domains: [String]) -> (out: String, ok: Bool) { run(["set"] + domains) }

    /// Tüm yapılandırmayı kaldırır (-y ile onay sorulmadan).
    @discardableResult
    static func uninstall() -> (out: String, ok: Bool) { run(["uninstall", "-y"]) }
}
