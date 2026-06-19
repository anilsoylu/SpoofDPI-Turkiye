import SwiftUI

// MARK: - Bağlantı testi ekranı (popover içi inline — sheet YOK, BUG4 fix)

struct ConnectionTestScreen: View {
    @EnvironmentObject private var state: AppState
    var onBack: () -> Void

    private let targets = ["Discord", "OpenAI / Codex", "Anthropic / Claude", "GitHub"]

    private var anyTesting: Bool {
        state.testResults.values.contains { $0.testing }
    }

    var body: some View {
        VStack(spacing: 0) {
            // Başlık çubuğu — geri butonu güvenilir biçimde main'e döner.
            HStack {
                Button {
                    onBack()
                } label: {
                    Label(state.t("done"), systemImage: "chevron.left")
                        .labelStyle(.titleAndIcon)
                }
                .buttonStyle(.plain)
                Spacer()
                Text(state.t("test.title"))
                    .font(.headline)
                Spacer()
                Label(state.t("done"), systemImage: "chevron.left")
                    .labelStyle(.titleAndIcon)
                    .opacity(0)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)

            Divider()

            List {
                ForEach(targets, id: \.self) { name in
                    HStack {
                        Text(name)
                            .font(.body)
                        Spacer()
                        statusView(for: state.testResults[name])
                    }
                }
            }
            .listStyle(.inset)
            .frame(height: 220)

            Divider()

            HStack {
                Spacer()
                Button {
                    state.runConnectionTests()
                } label: {
                    Label(state.t("test.run"), systemImage: "play.fill")
                }
                .buttonStyle(.borderedProminent)
                .tint(.green)
                .disabled(anyTesting)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
        }
    }

    @ViewBuilder
    private func statusView(for result: TestResult?) -> some View {
        if let result {
            if result.testing {
                ProgressView().controlSize(.small)
            } else {
                HStack(spacing: 5) {
                    Image(systemName: result.reachable ? "checkmark.circle.fill" : "xmark.circle.fill")
                        .foregroundStyle(result.reachable ? .green : .secondary)
                    Text(result.status)
                        .font(.caption)
                        .foregroundStyle(.secondary)
                        .monospacedDigit()
                }
            }
        } else {
            Text(state.t("test.idle"))
                .font(.caption)
                .foregroundStyle(.secondary)
        }
    }
}
