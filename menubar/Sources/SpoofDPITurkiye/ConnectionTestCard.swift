import SwiftUI

// MARK: - Bağlantı testi kartı

struct ConnectionTestCard: View {
    @EnvironmentObject private var state: AppState

    private let targets: [(name: String, url: String)] = [
        ("Discord",            "https://discord.com"),
        ("OpenAI / Codex",     "https://api.openai.com"),
        ("Anthropic / Claude", "https://api.anthropic.com"),
        ("GitHub",             "https://github.com"),
    ]

    var body: some View {
        VStack(alignment: .leading, spacing: 14) {
            // Başlık satırı
            HStack {
                SectionLabel(text: state.t("card.test.title"), icon: "antenna.radiowaves.left.and.right")
                Spacer()
                Button {
                    state.runConnectionTests()
                } label: {
                    Text(state.t("btn.runtest"))
                }
                .buttonStyle(OutlinePillButtonStyle())
            }

            // Test satırları
            VStack(spacing: 8) {
                ForEach(targets, id: \.name) { target in
                    TestRow(name: target.name, result: state.testResults[target.name])
                }
            }
        }
        .padding(20)
        .dsCard()
        .onAppear {
            if state.testResults.isEmpty {
                state.runConnectionTests()
            }
        }
    }
}

// MARK: - Tek test satırı

private struct TestRow: View {
    let name: String
    let result: TestResult?

    private var dotColor: Color {
        guard let r = result else { return Color.dsSecondaryText.opacity(0.4) }
        if r.testing { return Color.dsSecondaryText.opacity(0.5) }
        return r.reachable ? .dsTeal : .dsCoral
    }

    private var statusText: String {
        guard let r = result else { return "..." }
        return r.status
    }

    var body: some View {
        HStack(spacing: 10) {
            Circle()
                .fill(dotColor)
                .frame(width: 8, height: 8)

            Text(name)
                .font(.system(size: 13, weight: .medium))
                .foregroundStyle(Color.dsPrimaryText)

            Spacer()

            Text(statusText)
                .font(.system(size: 12, design: .monospaced))
                .foregroundStyle(Color.dsSecondaryText)
        }
        .padding(.horizontal, 14)
        .padding(.vertical, 10)
        .background(Color.dsRowBg)
        .clipShape(RoundedRectangle(cornerRadius: 10, style: .continuous))
    }
}
