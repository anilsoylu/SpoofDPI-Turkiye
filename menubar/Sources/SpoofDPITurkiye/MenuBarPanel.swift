import SwiftUI
import AppKit

// MARK: - Menü çubuğu paneli (native macOS / Control Center stili)

struct MenuBarPanel: View {
    @EnvironmentObject private var state: AppState
    @State private var domainsText: String = ""

    var body: some View {
        VStack(alignment: .leading, spacing: 14) {
            header
            protectionCard
            domainsSection
            Divider()
            footer
        }
        .padding(14)
        .frame(width: 312)
        // Materyal/popover varsayılan zemini — sistem otomatik adapte eder.
        .onAppear { domainsText = state.domains.joined(separator: "\n") }
        .onChange(of: state.domains) { domainsText = state.domains.joined(separator: "\n") }
        .sheet(isPresented: $state.showTestSheet) {
            ConnectionTestSheet().environmentObject(state)
        }
        .sheet(isPresented: $state.showSettingsSheet) {
            SettingsSheet().environmentObject(state)
        }
    }

    // MARK: - Header

    private var header: some View {
        HStack(spacing: 8) {
            Image(systemName: "shield.fill")
                .font(.system(size: 15))
                .foregroundStyle(.green)
            Text(state.t("app.title"))
                .font(.headline)
            Spacer()
            HStack(spacing: 5) {
                Circle()
                    .fill(state.running ? Color.green : Color.secondary)
                    .frame(width: 7, height: 7)
                Text(state.running ? state.t("status.on") : state.t("status.off.short"))
                    .font(.caption)
                    .foregroundStyle(.secondary)
            }
        }
    }

    // MARK: - Koruma kartı

    private var protectionCard: some View {
        HStack(alignment: .center) {
            VStack(alignment: .leading, spacing: 2) {
                Text(state.t("row.protection"))
                    .font(.body)
                Text("\(state.t("row.port")) \(state.port) · \(state.domains.count) \(state.t("row.domains.count"))")
                    .font(.caption)
                    .foregroundStyle(.secondary)
            }
            Spacer()
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
        .padding(12)
        .background(.quaternary, in: RoundedRectangle(cornerRadius: 10, style: .continuous))
    }

    // MARK: - Alan adları

    private var domainsSection: some View {
        VStack(alignment: .leading, spacing: 6) {
            Text(state.t("section.domains"))
                .font(.caption)
                .fontWeight(.semibold)
                .foregroundStyle(.secondary)

            Text(state.t("domains.hint"))
                .font(.caption)
                .foregroundStyle(.secondary)

            TextEditor(text: $domainsText)
                .font(.body.monospaced())
                .frame(height: 110)
                .scrollContentBackground(.hidden)
                .padding(6)
                .background(.quaternary, in: RoundedRectangle(cornerRadius: 8))

            Button {
                state.applyDomains(domainsText)
            } label: {
                Text(state.t("btn.saveapply"))
                    .frame(maxWidth: .infinity)
            }
            .buttonStyle(.borderedProminent)
            .tint(.green)
            .disabled(state.busy)
        }
    }

    // MARK: - Footer

    private var footer: some View {
        HStack {
            Button {
                state.showTestSheet = true
            } label: {
                Label(state.t("footer.test"), systemImage: "clock.arrow.circlepath")
            }
            Spacer()
            Button {
                state.showSettingsSheet = true
            } label: {
                Label(state.t("footer.settings"), systemImage: "gearshape")
            }
        }
        .buttonStyle(.plain)
        .font(.subheadline)
        .foregroundStyle(.primary)
    }
}
