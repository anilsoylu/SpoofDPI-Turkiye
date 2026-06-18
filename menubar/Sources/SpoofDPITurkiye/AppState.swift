import Foundation
import SwiftUI

// MARK: - Config JSON modeli

private struct ConfigJSON: Codable {
    var spoofDPIVersion: String
    var port: Int
    var domains: [String]
    var enableDoH: Bool
    var dnsAddr: String

    enum CodingKeys: String, CodingKey {
        case spoofDPIVersion = "spoofdpi_version"
        case port
        case domains
        case enableDoH = "enable_doh"
        case dnsAddr = "dns_addr"
    }

    static var empty: ConfigJSON {
        ConfigJSON(spoofDPIVersion: "", port: 8080, domains: [], enableDoH: true, dnsAddr: "1.1.1.1")
    }
}

// MARK: - Bağlantı testi sonucu

struct TestResult {
    var status: String   // "HTTP 200", "—", "..."
    var reachable: Bool
    var testing: Bool
}

// MARK: - AppState

@MainActor
final class AppState: ObservableObject {
    @Published var running: Bool = false
    @Published var port: Int = 8080
    @Published var domains: [String] = []
    @Published var version: String = (Bundle.main.infoDictionary?["CFBundleShortVersionString"] as? String) ?? "1.0.0"
    @Published var cliInstalled: Bool = false
    @Published var busy: Bool = false

    @Published var lang: Lang = .tr
    @Published var lastMessage: String = ""
    @Published var testResults: [String: TestResult] = [:]

    // Sheet sunumu
    @Published var showTestSheet: Bool = false
    @Published var showSettingsSheet: Bool = false

    // Çeviri kısayolu
    func t(_ key: String) -> String {
        Localization.t(key, lang: lang)
    }

    init() {
        refresh()
    }

    func refresh() {
        cliInstalled = CLI.findBinary() != nil

        let configPath = (NSHomeDirectory() as NSString)
            .appendingPathComponent(".spoofdpi-tr/config.json")
        if let data = try? Data(contentsOf: URL(fileURLWithPath: configPath)),
           let cfg = try? JSONDecoder().decode(ConfigJSON.self, from: data) {
            port = cfg.port
            domains = cfg.domains
        } else {
            port = 8080
            domains = []
        }

        running = isServiceRunning()
    }

    // MARK: - Servis durumu

    private func isServiceRunning() -> Bool {
        // Birincil: yeni CLI'ın `status` çıktısını ayrıştır. Helper "tpws: calisiyor"
        // satırını üretir (status komutu bunu yansıtır).
        if cliInstalled {
            let r = CLI.status()
            let out = r.out.lowercased()
            if out.contains("calisiyor") || out.contains("çalışıyor") {
                return true
            }
            if out.contains("durdu") || out.contains("bos") || out.contains("boş") {
                return false
            }
        }

        // Yedek: launchctl üzerinden LaunchDaemon yüklü mü kontrol et.
        let process = Process()
        process.executableURL = URL(fileURLWithPath: "/bin/launchctl")
        process.arguments = ["list"]

        let pipe = Pipe()
        process.standardOutput = pipe
        process.standardError = Pipe()

        do { try process.run() } catch { return false }
        let data = pipe.fileHandleForReading.readDataToEndOfFile()
        process.waitUntilExit()

        let output = String(data: data, encoding: .utf8) ?? ""
        return output.contains("com.spoofdpi-tr")
    }

    // MARK: - Aksiyonlar

    func toggle() {
        busy = true
        Task {
            if running {
                let r = CLI.off()
                lastMessage = r.out.trimmingCharacters(in: .whitespacesAndNewlines)
            } else {
                let r = CLI.on()
                lastMessage = r.out.trimmingCharacters(in: .whitespacesAndNewlines)
            }
            refresh()
            busy = false
        }
    }

    func restart() {
        busy = true
        Task {
            CLI.off()
            let r = CLI.on()
            lastMessage = r.out.trimmingCharacters(in: .whitespacesAndNewlines)
            refresh()
            busy = false
        }
    }

    func applyDomains(_ text: String) {
        busy = true
        Task {
            let lines = text
                .components(separatedBy: .newlines)
                .map { $0.trimmingCharacters(in: .whitespaces) }
                .filter { !$0.isEmpty }
            let r = CLI.set(lines)
            lastMessage = r.out.trimmingCharacters(in: .whitespacesAndNewlines)
            if lastMessage.isEmpty {
                lastMessage = lang == .tr
                    ? "Domainler kaydedildi ve uygulandı."
                    : "Domains saved and applied."
            }
            refresh()
            busy = false
        }
    }

    func uninstall() {
        busy = true
        Task {
            CLI.uninstall()
            lastMessage = lang == .tr ? "Kaldırıldı." : "Uninstalled."
            refresh()
            busy = false
        }
    }

    // MARK: - Eski yardımcılar (geriye uyumluluk)

    func addDomain(_ domain: String) {
        let trimmed = domain.trimmingCharacters(in: .whitespaces)
        guard !trimmed.isEmpty else { return }
        busy = true
        Task {
            CLI.add(trimmed)
            refresh()
            busy = false
        }
    }

    func removeDomain(_ domain: String) {
        busy = true
        Task {
            CLI.remove(domain)
            refresh()
            busy = false
        }
    }

    func setPort(_ p: Int) {
        busy = true
        Task {
            CLI.setPort(p)
            refresh()
            busy = false
        }
    }

    // MARK: - Bağlantı testi

    func runConnectionTests() {
        let targets: [(String, String)] = [
            ("Discord",            "https://discord.com"),
            ("OpenAI / Codex",     "https://api.openai.com"),
            ("Anthropic / Claude", "https://api.anthropic.com"),
            ("GitHub",             "https://github.com"),
        ]

        for (name, _) in targets {
            testResults[name] = TestResult(status: "...", reachable: false, testing: true)
        }

        Task {
            await withTaskGroup(of: (String, TestResult).self) { group in
                for (name, urlStr) in targets {
                    group.addTask {
                        guard let url = URL(string: urlStr) else {
                            return (name, TestResult(status: "—", reachable: false, testing: false))
                        }
                        var request = URLRequest(url: url, timeoutInterval: 6)
                        request.httpMethod = "HEAD"
                        do {
                            let (_, response) = try await URLSession.shared.data(for: request)
                            let code = (response as? HTTPURLResponse)?.statusCode ?? 0
                            return (name, TestResult(status: "HTTP \(code)", reachable: true, testing: false))
                        } catch {
                            return (name, TestResult(status: "—", reachable: false, testing: false))
                        }
                    }
                }
                for await (name, result) in group {
                    testResults[name] = result
                }
            }
        }
    }
}
