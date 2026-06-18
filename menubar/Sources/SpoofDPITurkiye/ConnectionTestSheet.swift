import SwiftUI

// MARK: - Bağlantı testi sheet'i (native Form/List)

struct ConnectionTestSheet: View {
    @EnvironmentObject private var state: AppState
    @Environment(\.dismiss) private var dismiss

    private let targets = ["Discord", "OpenAI / Codex", "Anthropic / Claude", "GitHub"]

    private var anyTesting: Bool {
        state.testResults.values.contains { $0.testing }
    }

    var body: some View {
        VStack(spacing: 0) {
            // Başlık çubuğu
            HStack {
                Text(state.t("test.title"))
                    .font(.headline)
                Spacer()
                Button(state.t("done")) { dismiss() }
                    .keyboardShortcut(.defaultAction)
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
        .frame(width: 360, height: 320)
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
