import SwiftUI

// MARK: - Sabit kategori listesi

private struct Category: Identifiable {
    let id: String  // key
    let title: String
    let domains: [String]
}

private let builtinCategories: [Category] = [
    Category(
        id: "discord",
        title: "Discord",
        domains: [
            "discord.com",
            "discordapp.com",
            "discord.gg",
            "discordapp.net",
            "discord.media",
            "discordcdn.com"
        ]
    )
]

// MARK: - ContentView

struct ContentView: View {
    @EnvironmentObject private var state: AppState

    @State private var portInput: String = ""
    @State private var newDomain: String = ""
    @State private var portError: Bool = false

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 0) {
                headerSection
                Divider()

                if !state.cliInstalled {
                    notInstalledSection
                } else {
                    mainSection
                }
            }
            .padding(.vertical, 8)
        }
        .frame(width: 320)
        .onAppear {
            portInput = "\(state.port)"
        }
        .onChange(of: state.port) { _, newVal in
            portInput = "\(newVal)"
        }
    }

    // MARK: - Başlık

    private var headerSection: some View {
        HStack {
            Text("SpoofDPI Türkiye")
                .font(.headline)
            Spacer()
            Circle()
                .fill(state.running ? Color.green : Color.secondary)
                .frame(width: 10, height: 10)
            Text(state.running ? "Çalışıyor" : "Durdu")
                .font(.caption)
                .foregroundStyle(.secondary)
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 10)
    }

    // MARK: - CLI kurulu değil

    private var notInstalledSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Label("CLI kurulu değil", systemImage: "exclamationmark.triangle.fill")
                .foregroundStyle(.orange)
                .font(.subheadline.bold())

            Text("Önce CLI'ı kur:")
                .font(.caption)
                .foregroundStyle(.secondary)

            Text("curl -fsSL https://raw.githubusercontent.com/anilsoylu/SpoofDPI-Turkiye/main/install.sh | bash")
                .font(.system(size: 10, design: .monospaced))
                .textSelection(.enabled)
                .padding(8)
                .background(Color.secondary.opacity(0.12))
                .clipShape(RoundedRectangle(cornerRadius: 6))

            Button("Yeniden Kontrol Et") {
                state.refresh()
            }
            .buttonStyle(.borderedProminent)
        }
        .padding(16)
    }

    // MARK: - Ana içerik

    private var mainSection: some View {
        VStack(alignment: .leading, spacing: 0) {
            toggleSection
            Divider()
            portSection
            Divider()
            categorySection
            Divider()
            domainListSection
            Divider()
            addDomainSection
            Divider()
            footerSection
        }
    }

    // MARK: - Başlat / Durdur

    private var toggleSection: some View {
        HStack {
            if state.busy {
                ProgressView()
                    .scaleEffect(0.8)
                    .padding(.trailing, 4)
                Text("İşleniyor...")
                    .foregroundStyle(.secondary)
            } else {
                Button(state.running ? "Durdur" : "Başlat") {
                    state.toggle()
                }
                .buttonStyle(.borderedProminent)
                .tint(state.running ? .red : .green)
            }
            Spacer()
        }
        .padding(16)
    }

    // MARK: - Port

    private var portSection: some View {
        VStack(alignment: .leading, spacing: 6) {
            Text("Proxy Portu")
                .font(.subheadline.bold())

            HStack {
                TextField("Port", text: $portInput)
                    .textFieldStyle(.roundedBorder)
                    .frame(width: 80)
                    .onChange(of: portInput) { _, val in
                        portError = false
                        _ = val  // suppress warning
                    }

                if portError {
                    Text("1-65535")
                        .font(.caption)
                        .foregroundStyle(.red)
                }

                Spacer()

                Button("Uygula") {
                    applyPort()
                }
                .disabled(state.busy)
            }
        }
        .padding(16)
    }

    private func applyPort() {
        guard let n = Int(portInput), n >= 1, n <= 65535 else {
            portError = true
            return
        }
        portError = false
        state.setPort(n)
    }

    // MARK: - Kategoriler

    private var categorySection: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text("Bypass Kategorileri")
                .font(.subheadline.bold())

            ForEach(builtinCategories) { cat in
                categoryToggle(cat)
            }
        }
        .padding(16)
    }

    private func categoryToggle(_ cat: Category) -> some View {
        let isOn = cat.domains.allSatisfy { state.domains.contains($0) }
        return Toggle(cat.title, isOn: Binding(
            get: { isOn },
            set: { newVal in
                if newVal {
                    for domain in cat.domains {
                        state.addDomain(domain)
                    }
                } else {
                    for domain in cat.domains {
                        state.removeDomain(domain)
                    }
                }
            }
        ))
        .disabled(state.busy)
    }

    // MARK: - Domain listesi

    private var domainListSection: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text("Bypass Edilen Domainler (\(state.domains.count))")
                .font(.subheadline.bold())

            if state.domains.isEmpty {
                Text("(domain yok)")
                    .font(.caption)
                    .foregroundStyle(.secondary)
            } else {
                ForEach(state.domains, id: \.self) { domain in
                    HStack {
                        Text(domain)
                            .font(.system(size: 12, design: .monospaced))
                        Spacer()
                        Button {
                            state.removeDomain(domain)
                        } label: {
                            Image(systemName: "xmark.circle.fill")
                                .foregroundStyle(.secondary)
                        }
                        .buttonStyle(.plain)
                        .disabled(state.busy)
                    }
                }
            }
        }
        .padding(16)
    }

    // MARK: - Özel domain ekleme

    private var addDomainSection: some View {
        VStack(alignment: .leading, spacing: 6) {
            Text("Özel Domain Ekle")
                .font(.subheadline.bold())

            HStack {
                TextField("örn. netflix.com", text: $newDomain)
                    .textFieldStyle(.roundedBorder)
                    .onSubmit { addNewDomain() }

                Button("Ekle") { addNewDomain() }
                    .disabled(newDomain.trimmingCharacters(in: .whitespaces).isEmpty || state.busy)
            }
        }
        .padding(16)
    }

    private func addNewDomain() {
        let trimmed = newDomain.trimmingCharacters(in: .whitespaces)
        guard !trimmed.isEmpty else { return }
        state.addDomain(trimmed)
        newDomain = ""
    }

    // MARK: - Alt bar

    private var footerSection: some View {
        HStack {
            Button("Güncelle") {
                state.update()
            }
            .disabled(state.busy)

            Spacer()

            if !state.version.isEmpty {
                Text("v\(state.version)")
                    .font(.caption2)
                    .foregroundStyle(.tertiary)
            }

            Spacer()

            Button("Çıkış") {
                NSApplication.shared.terminate(nil)
            }
        }
        .padding(16)
    }
}
