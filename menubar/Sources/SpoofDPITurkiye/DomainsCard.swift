import SwiftUI

// MARK: - Domain kartı

struct DomainsCard: View {
    @EnvironmentObject private var state: AppState
    @State private var domainsText: String = ""

    private let discordDomains = [
        "discord.com",
        "discordapp.com",
        "discord.gg",
        "discordapp.net",
        "discord.media",
        "discordcdn.com"
    ]

    var body: some View {
        VStack(alignment: .leading, spacing: 14) {
            // Başlık
            SectionLabel(text: state.t("card.domains.title"), icon: "globe")

            // Yardımcı metin
            Text(state.t("card.domains.hint"))
                .font(.system(size: 12))
                .foregroundStyle(Color.dsSecondaryText)
                .fixedSize(horizontal: false, vertical: true)

            // TextEditor
            TextEditor(text: $domainsText)
                .font(.system(size: 13, design: .monospaced))
                .foregroundStyle(Color.dsPrimaryText)
                .scrollContentBackground(.hidden)
                .background(Color.dsBackground)
                .clipShape(RoundedRectangle(cornerRadius: 12, style: .continuous))
                .overlay(
                    RoundedRectangle(cornerRadius: 12, style: .continuous)
                        .stroke(Color.dsBorder, lineWidth: 1)
                )
                .frame(minHeight: 180, maxHeight: 260)

            // Alt buton satırı
            HStack(spacing: 10) {
                Button {
                    domainsText = discordDomains.joined(separator: "\n")
                } label: {
                    Label(state.t("btn.discord"), systemImage: "bolt.fill")
                }
                .buttonStyle(OutlinePillButtonStyle())

                Spacer()

                Button {
                    state.applyDomains(domainsText)
                } label: {
                    Text(state.t("btn.saveapply"))
                }
                .buttonStyle(PrimaryButtonStyle(variant: .indigo))
                .frame(width: 160)
                .disabled(state.busy)
            }
        }
        .padding(20)
        .dsCard()
        .onAppear {
            domainsText = state.domains.joined(separator: "\n")
        }
        .onChange(of: state.domains) {
            let newDomains = state.domains
            // Sadece dışarıdan refresh gelirse güncelle
            let current = domainsText
                .components(separatedBy: .newlines)
                .map { $0.trimmingCharacters(in: .whitespaces) }
                .filter { !$0.isEmpty }
                .sorted()
            if current != newDomains.sorted() {
                domainsText = newDomains.joined(separator: "\n")
            }
        }
    }
}
