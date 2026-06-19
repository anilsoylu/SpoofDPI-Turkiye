import SwiftUI
import AppKit

// MARK: - Menü çubuğu paneli (native macOS / Control Center stili)

// Popover içi ekranlar — sheet yerine tek pencere içinde gezinme (BUG4 fix).
enum PanelScreen {
    case main
    case settings
    case test
}

struct MenuBarPanel: View {
    @EnvironmentObject private var state: AppState
    @State private var domainsText: String = ""
    // Aktif ekran. Sheet kullanmıyoruz; geri butonuyla main'e döneriz.
    @State private var screen: PanelScreen = .main
    // TextEditor'ın odakta olup olmadığı. domainsText'i state.domains ile yalnızca
    // editor ODAKTA DEĞİLKEN senkronlarız: kullanıcı yazarken ezilmez (BUG3),
    // harici/kaydet sonrası değişimler yansır (#5: stale TextEditor düzeltildi).
    @FocusState private var domainsFocused: Bool
    // ALAN ADLARI bölümü açılır-kapanır. Varsayılan KAPALI — panel daha ferah.
    @State private var domainsExpanded = false

    // Discord hızlı profili — her satır bir kök alan adı.
    private let discordDomains = [
        "discord.com",
        "discordapp.com",
        "discord.gg",
        "discordapp.net",
        "discord.media",
        "discordcdn.com",
    ]

    var body: some View {
        Group {
            switch screen {
            case .main:
                mainScreen
            case .settings:
                SettingsScreen(onBack: { screen = .main })
                    .environmentObject(state)
            case .test:
                ConnectionTestScreen(onBack: { screen = .main })
                    .environmentObject(state)
            }
        }
        .frame(width: 332)
        // İlk açılışta domains metnini state ile doldur (editor henüz odakta değil).
        .onAppear { syncDomainsTextIfUnfocused() }
        // state.domains harici/kaydet sonrası değişirse metni senkronla — ama
        // YALNIZCA editor odakta değilken (kullanıcı yazarken ezme; BUG3 korunur,
        // #5: stale TextEditor çözülür).
        .onChange(of: state.domains) { _, _ in syncDomainsTextIfUnfocused() }
    }

    // MARK: - Ana ekran

    private var mainScreen: some View {
        VStack(alignment: .leading, spacing: 16) {
            header

            if state.cliInstalled {
                heroSection
                domainsSection
            } else {
                installCard
            }

            Divider()
            footer
        }
        .padding(16)
    }

    // domainsText'i state.domains ile senkronlar — YALNIZCA editor odakta
    // değilken. Kullanıcı yazarken çağrılsa bile odak nedeniyle metni ezmez
    // (BUG3); harici/kaydet sonrası değişimleri yansıtır (#5).
    private func syncDomainsTextIfUnfocused() {
        guard !domainsFocused else { return }
        let joined = state.domains.joined(separator: "\n")
        if domainsText != joined {
            domainsText = joined
        }
    }

    // MARK: - Header (minik marka, durum noktası kahraman karta taşındı)

    private var header: some View {
        HStack(spacing: 7) {
            Image(systemName: "shield.fill")
                .font(.system(size: 13))
                .foregroundStyle(.green)
            Text(state.t("app.title"))
                .font(.headline)
            Spacer()
        }
    }

    // MARK: - Kahraman: koruma durumu kartı + aksiyon geri bildirimi

    private var heroSection: some View {
        VStack(alignment: .leading, spacing: 8) {
            protectionCard
            if !state.lastMessage.isEmpty {
                feedbackRow
            }
        }
        .animation(.easeInOut(duration: 0.2), value: state.lastMessage)
    }

    private var protectionCard: some View {
        HStack(alignment: .center, spacing: 12) {
            Image(systemName: state.running ? "checkmark.shield.fill" : "shield.slash")
                .font(.system(size: 28))
                .foregroundStyle(state.running ? AnyShapeStyle(.green) : AnyShapeStyle(.secondary))
                .symbolRenderingMode(.hierarchical)

            VStack(alignment: .leading, spacing: 2) {
                Text(state.running ? state.t("status.protected.title") : state.t("status.off.title"))
                    .font(.headline)
                Text(subtitle)
                    .font(.caption)
                    .foregroundStyle(.secondary)
                    .fixedSize(horizontal: false, vertical: true)
            }

            Spacer(minLength: 8)

            if state.busy {
                ProgressView().controlSize(.small)
            } else {
                Toggle("", isOn: Binding(
                    get: { state.running },
                    set: { _ in state.toggle() }
                ))
                .labelsHidden()
                .toggleStyle(.switch)
                .tint(.green)
            }
        }
        .padding(14)
        .background(heroBackground)
        .clipShape(RoundedRectangle(cornerRadius: 12, style: .continuous))
        .animation(.easeInOut(duration: 0.2), value: state.running)
    }

    private var subtitle: String {
        if state.running {
            // "\(sayı) site korunuyor · port \(port)"
            return "\(state.domains.count) \(state.t("status.protected.sub")) · \(state.t("row.port").lowercased()) \(state.port)"
        } else {
            return state.t("status.off.sub")
        }
    }

    @ViewBuilder
    private var heroBackground: some View {
        if state.running {
            Color.green.opacity(0.12)
        } else {
            Color(nsColor: .quaternaryLabelColor).opacity(0.6)
        }
    }

    private var feedbackRow: some View {
        HStack(spacing: 6) {
            Image(systemName: "checkmark.circle.fill")
                .font(.caption)
                .foregroundStyle(.green)
            Text(state.lastMessage)
                .font(.caption)
                .foregroundStyle(.secondary)
                .lineLimit(2)
        }
        .padding(.horizontal, 2)
        .transition(.opacity)
    }

    // MARK: - Kurulu değil — net uyarı kartı

    private var installCommand = "curl -fsSL https://raw.githubusercontent.com/anilsoylu/SpoofDPI-Turkiye/master/install.sh | bash"

    private var installCard: some View {
        VStack(alignment: .leading, spacing: 10) {
            HStack(alignment: .center, spacing: 12) {
                Image(systemName: "exclamationmark.triangle.fill")
                    .font(.system(size: 26))
                    .foregroundStyle(.orange)
                    .symbolRenderingMode(.hierarchical)
                VStack(alignment: .leading, spacing: 2) {
                    Text(state.t("install.incomplete"))
                        .font(.headline)
                    Text(state.t("install.run"))
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }
                Spacer(minLength: 0)
            }

            HStack(spacing: 8) {
                Text(installCommand)
                    .font(.caption.monospaced())
                    .textSelection(.enabled)
                    .lineLimit(3)
                    .frame(maxWidth: .infinity, alignment: .leading)
                Button {
                    let pb = NSPasteboard.general
                    pb.clearContents()
                    pb.setString(installCommand, forType: .string)
                } label: {
                    Image(systemName: "doc.on.doc")
                }
                .buttonStyle(.borderless)
                .help("Kopyala")
            }
            .padding(10)
            .background(.quaternary, in: RoundedRectangle(cornerRadius: 8, style: .continuous))
        }
        .padding(14)
        .background(Color.orange.opacity(0.10))
        .clipShape(RoundedRectangle(cornerRadius: 12, style: .continuous))
    }

    // MARK: - Alan adları (sürtünmeyi azalt)

    private var domainsSection: some View {
        VStack(alignment: .leading, spacing: 8) {
            // Açılır-kapanır başlık satırı — tüm satır tıklanabilir.
            Button {
                withAnimation(.easeInOut(duration: 0.2)) { domainsExpanded.toggle() }
            } label: {
                HStack {
                    Text(state.t("section.domains"))
                        .font(.caption)
                        .fontWeight(.semibold)
                        .foregroundStyle(.secondary)
                    Text("\(state.domains.count)")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                        .monospacedDigit()
                    Spacer()
                    Image(systemName: domainsExpanded ? "chevron.down" : "chevron.right")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }
            }
            .buttonStyle(.plain)
            .contentShape(Rectangle())

            if domainsExpanded {
                Text(state.t("domains.hint"))
                    .font(.caption)
                    .foregroundStyle(.secondary)

                TextEditor(text: $domainsText)
                    .font(.body.monospaced())
                    .frame(height: 110)
                    .scrollContentBackground(.hidden)
                    .padding(6)
                    .background(.quaternary, in: RoundedRectangle(cornerRadius: 8))
                    .focused($domainsFocused)

                HStack(spacing: 8) {
                    Button {
                        domainsText = discordDomains.joined(separator: "\n")
                    } label: {
                        Label(state.t("btn.discordprofile"), systemImage: "bolt.fill")
                            .font(.callout)
                    }
                    .buttonStyle(.bordered)
                    .controlSize(.small)
                    .disabled(state.busy)

                    Spacer()
                }

                Button {
                    // Kaydet → CLI set → persist. Tamamlanınca normalize edilmiş
                    // listeyle metni yeniden senkronla (BUG3: kaydet gerçekten
                    // persist eder, panel kapanıp açılınca domainler durur).
                    state.applyDomains(domainsText) { normalized in
                        domainsText = normalized.joined(separator: "\n")
                    }
                } label: {
                    Text(state.t("btn.saveapply"))
                        .frame(maxWidth: .infinity)
                }
                .buttonStyle(.borderedProminent)
                .tint(.green)
                .disabled(state.busy)
            }
        }
    }

    // MARK: - Footer

    private var footer: some View {
        HStack {
            Button {
                screen = .test
            } label: {
                Label(state.t("footer.test"), systemImage: "clock.arrow.circlepath")
            }
            Spacer()
            Button {
                screen = .settings
            } label: {
                Label(state.t("footer.settings"), systemImage: "gearshape")
            }
        }
        .buttonStyle(.plain)
        .font(.subheadline)
        .foregroundStyle(.primary)
    }
}
