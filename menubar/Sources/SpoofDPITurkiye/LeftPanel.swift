import SwiftUI

// MARK: - Sol panel

struct LeftPanel: View {
    @EnvironmentObject private var state: AppState
    @State private var showUninstallAlert = false

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            // 1. Wordmark + TR/EN toggle
            wordmarkRow
                .padding(.horizontal, 24)
                .padding(.top, 28)
                .padding(.bottom, 20)

            // 2. Kalkan hero
            ShieldHero(running: state.running)
                .frame(maxWidth: .infinity)
                .padding(.bottom, 16)

            // 3. Durum metni
            statusText
                .frame(maxWidth: .infinity)
                .padding(.bottom, 24)

            // 4. Ana buton
            primaryButton
                .padding(.horizontal, 24)
                .padding(.bottom, 12)

            // 5. İkincil butonlar
            secondaryButtons
                .padding(.horizontal, 24)
                .padding(.bottom, 20)

            // 6. Ayraç + bilgi satırı
            Divider()
                .background(Color.dsBorder)
                .padding(.horizontal, 24)
                .padding(.bottom, 16)

            infoRow
                .padding(.horizontal, 24)
                .padding(.bottom, 12)

            // 7. Kaldır butonu
            uninstallButton
                .padding(.horizontal, 24)
                .padding(.bottom, 16)

            Spacer(minLength: 0)

            // 8. Son mesaj
            if !state.lastMessage.isEmpty {
                lastMessageBar
                    .padding(.horizontal, 24)
                    .padding(.bottom, 20)
            }
        }
        .background(Color.dsPanel)
        .alert(state.t("uninstall.confirm"), isPresented: $showUninstallAlert) {
            Button(state.t("uninstall.confirm.btn"), role: .destructive) {
                state.uninstall()
            }
            Button(state.t("cancel"), role: .cancel) {}
        }
    }

    // MARK: - Wordmark + dil toggle

    private var wordmarkRow: some View {
        HStack(alignment: .center) {
            VStack(alignment: .leading, spacing: 2) {
                Text("SPOOFDPI")
                    .font(.system(size: 18, weight: .bold))
                    .foregroundStyle(Color.dsPrimaryText)
                HStack(spacing: 4) {
                    Text(state.t("wordmark.sub"))
                        .font(.system(size: 10, weight: .semibold))
                        .tracking(1.5)
                        .foregroundStyle(Color.dsSecondaryText)
                    if !state.version.isEmpty {
                        Text("· v\(state.version)")
                            .font(.system(size: 10, weight: .regular))
                            .tracking(0.5)
                            .foregroundStyle(Color.dsSecondaryText.opacity(0.7))
                    }
                }
            }

            Spacer()

            LangToggle()
        }
    }

    // MARK: - Durum metni

    private var statusText: some View {
        VStack(spacing: 6) {
            Text(state.running ? state.t("status.protected") : state.t("status.off"))
                .font(.system(size: 26, weight: .bold))
                .foregroundStyle(state.running ? Color.dsTeal : Color.dsPrimaryText)

            if state.running {
                Text("\(state.t("status.sub.on")) · \(state.domains.count) \(state.t("status.sub.domains"))")
                    .font(.system(size: 13))
                    .foregroundStyle(Color.dsSecondaryText)
            } else {
                Text(state.t("status.sub.off"))
                    .font(.system(size: 13))
                    .foregroundStyle(Color.dsSecondaryText)
            }
        }
    }

    // MARK: - Birincil buton

    private var primaryButton: some View {
        Group {
            if state.busy {
                HStack(spacing: 10) {
                    ProgressView()
                        .scaleEffect(0.8)
                        .tint(.black.opacity(0.7))
                    Text(state.t("info.processing"))
                        .font(.system(size: 15, weight: .semibold))
                        .foregroundStyle(.black.opacity(0.75))
                }
                .frame(maxWidth: .infinity)
                .frame(height: 50)
                .background(
                    state.running
                        ? LinearGradient(colors: [.dsCoral, .dsCoralDark], startPoint: .topLeading, endPoint: .bottomTrailing)
                        : LinearGradient(colors: [.dsTeal, .dsTealDark], startPoint: .topLeading, endPoint: .bottomTrailing)
                )
                .clipShape(RoundedRectangle(cornerRadius: 14, style: .continuous))
            } else {
                Button {
                    state.toggle()
                } label: {
                    HStack(spacing: 8) {
                        Image(systemName: state.running ? "stop.fill" : "play.fill")
                            .font(.system(size: 14, weight: .bold))
                        Text(state.running ? state.t("btn.stop") : state.t("btn.start"))
                    }
                }
                .buttonStyle(PrimaryButtonStyle(variant: state.running ? .coral : .teal))
                .disabled(state.busy)
            }
        }
    }

    // MARK: - İkincil butonlar

    private var secondaryButtons: some View {
        HStack(spacing: 10) {
            Button {
                state.restart()
            } label: {
                Label(state.t("btn.restart"), systemImage: "arrow.clockwise")
            }
            .buttonStyle(OutlinePillButtonStyle())
            .disabled(state.busy)

            Button {
                state.refresh()
            } label: {
                Label(state.t("btn.refresh"), systemImage: "arrow.triangle.2.circlepath")
            }
            .buttonStyle(OutlinePillButtonStyle())
            .disabled(state.busy)
        }
    }

    // MARK: - Bilgi satırı

    private var infoRow: some View {
        HStack(spacing: 6) {
            Image(systemName: "checkmark.circle.fill")
                .font(.system(size: 13))
                .foregroundStyle(Color.dsTeal)
            Text(state.t("info.autostart"))
                .font(.system(size: 12))
                .foregroundStyle(Color.dsSecondaryText)
        }
    }

    // MARK: - Kaldır butonu

    private var uninstallButton: some View {
        Button {
            showUninstallAlert = true
        } label: {
            HStack(spacing: 6) {
                Image(systemName: "trash")
                    .font(.system(size: 13))
                Text(state.t("btn.uninstall"))
                    .font(.system(size: 13, weight: .medium))
            }
            .foregroundStyle(Color.dsCoral)
        }
        .buttonStyle(.plain)
        .disabled(state.busy)
    }

    // MARK: - Son mesaj

    private var lastMessageBar: some View {
        Text(state.lastMessage)
            .font(.system(size: 11))
            .foregroundStyle(Color.dsSecondaryText)
            .lineLimit(2)
            .multilineTextAlignment(.leading)
            .frame(maxWidth: .infinity, alignment: .leading)
    }
}

// MARK: - Dil toggle

private struct LangToggle: View {
    @EnvironmentObject private var state: AppState

    var body: some View {
        HStack(spacing: 0) {
            ForEach(Lang.allCases, id: \.self) { lang in
                Button {
                    state.lang = lang
                } label: {
                    Text(lang.rawValue)
                        .font(.system(size: 12, weight: .semibold))
                        .foregroundStyle(state.lang == lang ? Color.black.opacity(0.8) : Color.dsSecondaryText)
                        .frame(width: 36, height: 26)
                        .background(state.lang == lang ? Color.dsTeal : Color.clear)
                        .clipShape(RoundedRectangle(cornerRadius: 8, style: .continuous))
                }
                .buttonStyle(.plain)
            }
        }
        .background(Color.dsBackground)
        .clipShape(RoundedRectangle(cornerRadius: 10, style: .continuous))
        .overlay(
            RoundedRectangle(cornerRadius: 10, style: .continuous)
                .stroke(Color.dsBorder, lineWidth: 1)
        )
    }
}
