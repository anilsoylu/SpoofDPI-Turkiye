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

// MARK: - AppState

@MainActor
final class AppState: ObservableObject {
    @Published var running: Bool = false
    @Published var port: Int = 8080
    @Published var domains: [String] = []
    @Published var version: String = ""
    @Published var cliInstalled: Bool = false
    @Published var busy: Bool = false

    init() {
        refresh()
    }

    func refresh() {
        cliInstalled = CLI.findBinary() != nil

        // config.json oku
        let configPath = (NSHomeDirectory() as NSString)
            .appendingPathComponent(".spoofdpi-tr/config.json")
        if let data = try? Data(contentsOf: URL(fileURLWithPath: configPath)),
           let cfg = try? JSONDecoder().decode(ConfigJSON.self, from: data) {
            port = cfg.port
            domains = cfg.domains
            version = cfg.spoofDPIVersion
        } else {
            port = 8080
            domains = []
            version = ""
        }

        // Servis durumu: launchctl list içinde com.spoofdpi-tr ara
        running = isServiceRunning()
    }

    // MARK: - Aksiyonlar

    func toggle() {
        busy = true
        Task {
            if running {
                CLI.off()
            } else {
                CLI.on()
            }
            refresh()
            busy = false
        }
    }

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

    func update() {
        busy = true
        Task {
            CLI.update()
            refresh()
            busy = false
        }
    }

    // MARK: - Yardımcı

    private func isServiceRunning() -> Bool {
        let process = Process()
        process.executableURL = URL(fileURLWithPath: "/bin/launchctl")
        process.arguments = ["list"]

        let pipe = Pipe()
        process.standardOutput = pipe
        process.standardError = Pipe() // sessizce yut

        do {
            try process.run()
        } catch {
            return false
        }
        let data = pipe.fileHandleForReading.readDataToEndOfFile()
        process.waitUntilExit()

        let output = String(data: data, encoding: .utf8) ?? ""
        return output.contains("com.spoofdpi-tr")
    }
}
