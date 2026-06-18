import SwiftUI
import AppKit

// MARK: - Menü çubuğu popover paneli (~360pt genişlik)

struct MenuBarPanel: View {
    @EnvironmentObject private var state: AppState
    @Environment(\.openWindow) private var openWindow
    @State private var showUninstallAlert = false

    private var domainCount: Int { state.domains.count }

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {

                // 1. Üst satır: wordmark + dil toggle + durum noktası
                headerRow
                    .padding(.horizontal, 14)
                    .padding(.top, 14)
                    .padding(.bottom, 12)

                // 2. Kompakt ShieldHero
                ShieldHero(running: state.running, size: 96, animated: false)
                    .frame(maxWidth: .infinity)
                    .padding(.bottom, 8)

                // 3. Durum metni
                statusBlock
                    .frame(maxWidth: .infinity)
                    .padding(.bottom, 14)

                // 4. Birincil buton
                primaryActionButton
                    .padding(.horizontal, 14)
                    .padding(.bottom, 10)

                // 5. İkincil butonlar
                secondaryActionButtons
                    .padding(.horizontal, 14)
                    .padding(.bottom, 12)

                // 6. İnce ayraç
                Rectangle()
                    .fill(Color.dsBorder)
                    .frame(height: 1)
                    .padding(.horizontal, 14)
                    .padding(.bottom, 10)

                // 7. Domain özeti satırı
                domainSummaryRow
                    .padding(.horizontal, 14)
                    .padding(.bottom, 10)

                // 8. Detaylar butonu
                detailsButton
                    .padding(.horizontal, 14)
                    .padding(.bottom, 12)

                // 9. İnce ayraç + alt satır
                Rectangle()
                    .fill(Color.dsBorder)
                    .frame(height: 1)
                    .padding(.horizontal, 14)
                    .padding(.bottom, 10)

                footerRow
                    .padding(.horizontal, 14)
                    .padding(.bottom, 12)

                // 10. Son mesaj (varsa)
                if !state.lastMessage.isEmpty {
                    Text(state.lastMessage)
                        .font(.system(size: 10))
                        .foregroundStyle(Color.dsSecondaryText.opacity(0.7))
                        .lineLimit(2)
                        .multilineTextAlignment(.leading)
                        .frame(maxWidth: .infinity, alignment: .leading)
                        .padding(.horizontal, 14)
                        .padding(.bottom, 10)
                }
            }
        .frame(width: 320)
        .background(Color.dsPanel)
        .alert(state.t("uninstall.confirm"), isPresented: $showUninstallAlert) {
            Button(state.t("uninstall.confirm.btn"), role: .destructive) {
                state.uninstall()
            }
            Button(state.t("cancel"), role: .cancel) {}
        }
    }

    // MARK: - Üst satır

    private var headerRow: some View {
        HStack(alignment: .center, spacing: 8) {
            // Wordmark
            VStack(alignment: .leading, spacing: 2) {
                Text("SPOOFDPI")
                    .font(.system(size: 15, weight: .bold))
                    .foregroundStyle(Color.dsPrimaryText)
                Text("TÜRKİYE · v\(state.version.isEmpty ? "—" : state.version)")
                    .font(.system(size: 9, weight: .semibold))
                    .tracking(1.2)
                    .foregroundStyle(Color.dsSecondaryText)
            }

            Spacer()

            // TR/EN pill toggle
            MenuBarLangToggle()

            // Durum noktası
            Circle()
                .fill(state.running ? Color.dsTeal : Color.dsSecondaryText.opacity(0.4))
                .frame(width: 8, height: 8)
                .shadow(color: state.running ? Color.dsTeal.opacity(0.6) : .clear, radius: 4)
        }
    }

    // MARK: - Durum metni

    private var statusBlock: some View {
        VStack(spacing: 4) {
            Text(state.running ? state.t("status.protected") : state.t("status.off"))
                .font(.system(size: 18, weight: .semibold))
                .foregroundStyle(state.running ? Color.dsTeal : Color.dsPrimaryText)

            Text(state.running
                    ? "\(state.t("status.sub.on")) · \(domainCount) \(state.t("status.sub.domains"))"
                    : state.t("status.sub.off"))
                .font(.system(size: 12))
                .foregroundStyle(Color.dsSecondaryText)
        }
        .multilineTextAlignment(.center)
    }

    // MARK: - Birincil buton

    private var primaryActionButton: some View {
        Group {
            if state.busy {
                HStack(spacing: 10) {
                    ProgressView()
                        .scaleEffect(0.75)
                        .tint(.black.opacity(0.7))
                    Text(state.t("info.processing"))
                        .font(.system(size: 14, weight: .semibold))
                        .foregroundStyle(.black.opacity(0.75))
                }
                .frame(maxWidth: .infinity)
                .frame(height: 44)
                .background(
                    state.running
                        ? LinearGradient(colors: [.dsCoral, .dsCoralDark], startPoint: .topLeading, endPoint: .bottomTrailing)
                        : LinearGradient(colors: [.dsTeal, .dsTealDark], startPoint: .topLeading, endPoint: .bottomTrailing)
                )
                .clipShape(RoundedRectangle(cornerRadius: 12, style: .continuous))
            } else {
                Button {
                    state.toggle()
                } label: {
                    HStack(spacing: 7) {
                        Image(systemName: state.running ? "stop.fill" : "play.fill")
                            .font(.system(size: 13, weight: .bold))
                        Text(state.running
                                ? state.t("btn.stop")
                                : state.t("btn.start"))
                    }
                }
                .buttonStyle(CompactPrimaryButtonStyle(variant: state.running ? .coral : .teal))
                .disabled(state.busy)
            }
        }
    }

    // MARK: - İkincil butonlar

    private var secondaryActionButtons: some View {
        HStack(spacing: 10) {
            Button {
                state.restart()
            } label: {
                Label(state.t("btn.restart"), systemImage: "arrow.clockwise")
                    .frame(maxWidth: .infinity)
            }
            .buttonStyle(OutlinePillButtonStyle())
            .disabled(state.busy)

            Button {
                state.refresh()
            } label: {
                Label(state.t("btn.refresh"), systemImage: "arrow.triangle.2.circlepath")
                    .frame(maxWidth: .infinity)
            }
            .buttonStyle(OutlinePillButtonStyle())
            .disabled(state.busy)
        }
    }

    // MARK: - Domain özeti

    private var domainSummaryRow: some View {
        HStack {
            Label("\(domainCount) domain filtreleniyor", systemImage: "globe")
                .font(.system(size: 12))
                .foregroundStyle(Color.dsSecondaryText)

            Spacer()

            Button {
                openMainWindow()
            } label: {
                Text("Düzenle →")
                    .font(.system(size: 12, weight: .medium))
                    .foregroundStyle(Color.dsTeal)
            }
            .buttonStyle(.plain)
        }
    }

    // MARK: - Detaylar butonu

    private var detailsButton: some View {
        Button {
            openMainWindow()
        } label: {
            Text("Bağlantı Testi & Detaylar")
                .font(.system(size: 13, weight: .medium))
                .foregroundStyle(Color.dsSecondaryText)
                .frame(maxWidth: .infinity)
                .frame(height: 36)
                .background(Color.dsPanel)
                .clipShape(RoundedRectangle(cornerRadius: 10, style: .continuous))
                .overlay(
                    RoundedRectangle(cornerRadius: 10, style: .continuous)
                        .stroke(Color.dsBorder, lineWidth: 1)
                )
        }
        .buttonStyle(.plain)
    }

    // MARK: - Alt satır

    private var footerRow: some View {
        VStack(spacing: 8) {
            HStack(spacing: 6) {
                // Açılışta otomatik bilgisi
                Image(systemName: "checkmark.circle.fill")
                    .font(.system(size: 11))
                    .foregroundStyle(Color.dsTeal)
                Text(state.t("info.autostart"))
                    .font(.system(size: 11))
                    .foregroundStyle(Color.dsSecondaryText)

                Spacer()

                // Çıkış
                Button {
                    NSApplication.shared.terminate(nil)
                } label: {
                    Text(state.t("menu.quit"))
                        .font(.system(size: 11, weight: .medium))
                        .foregroundStyle(Color.dsSecondaryText)
                }
                .buttonStyle(.plain)
            }

            HStack {
                Spacer()
                // Kaldır
                Button {
                    showUninstallAlert = true
                } label: {
                    Text(state.t("btn.uninstall"))
                        .font(.system(size: 11, weight: .medium))
                        .foregroundStyle(Color.dsCoral)
                }
                .buttonStyle(.plain)
                .disabled(state.busy)
            }
        }
    }

    // MARK: - Yardımcı

    private func openMainWindow() {
        openWindow(id: "main")
        NSApp.activate(ignoringOtherApps: true)
    }
}

// MARK: - Kompakt birincil buton stili (popover için daha küçük)

private struct CompactPrimaryButtonStyle: ButtonStyle {
    enum Variant { case teal, coral }
    let variant: Variant

    private var gradient: LinearGradient {
        switch variant {
        case .teal:
            return LinearGradient(colors: [.dsTeal, .dsTealDark], startPoint: .topLeading, endPoint: .bottomTrailing)
        case .coral:
            return LinearGradient(colors: [.dsCoral, .dsCoralDark], startPoint: .topLeading, endPoint: .bottomTrailing)
        }
    }

    private var glowColor: Color {
        switch variant {
        case .teal:  return Color.dsTeal.opacity(0.5)
        case .coral: return Color.dsCoral.opacity(0.5)
        }
    }

    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .font(.system(size: 14, weight: .semibold))
            .foregroundStyle(.black.opacity(0.85))
            .frame(maxWidth: .infinity)
            .frame(height: 44)
            .background(gradient)
            .clipShape(RoundedRectangle(cornerRadius: 12, style: .continuous))
            .shadow(color: glowColor, radius: configuration.isPressed ? 4 : 10, x: 0, y: 0)
            .scaleEffect(configuration.isPressed ? 0.97 : 1)
            .animation(.easeOut(duration: 0.12), value: configuration.isPressed)
    }
}

// MARK: - Menü çubuğu dil toggle (panel içi, internal erişim)

struct MenuBarLangToggle: View {
    @EnvironmentObject private var state: AppState

    var body: some View {
        HStack(spacing: 0) {
            ForEach(Lang.allCases, id: \.self) { lang in
                Button {
                    state.lang = lang
                } label: {
                    Text(lang.rawValue)
                        .font(.system(size: 11, weight: .semibold))
                        .foregroundStyle(state.lang == lang ? Color.black.opacity(0.8) : Color.dsSecondaryText)
                        .frame(width: 30, height: 22)
                        .background(state.lang == lang ? Color.dsTeal : Color.clear)
                        .clipShape(RoundedRectangle(cornerRadius: 7, style: .continuous))
                }
                .buttonStyle(.plain)
            }
        }
        .background(Color.dsBackground)
        .clipShape(RoundedRectangle(cornerRadius: 8, style: .continuous))
        .overlay(
            RoundedRectangle(cornerRadius: 8, style: .continuous)
                .stroke(Color.dsBorder, lineWidth: 1)
        )
    }
}
