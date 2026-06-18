import SwiftUI
import AppKit

// MARK: - Menü çubuğu paneli (native macOS / Control Center stili)

struct MenuBarPanel: View {
    @EnvironmentObject private var state: AppState

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

            if state.domains.isEmpty {
                Text(state.t("domains.empty"))
                    .font(.caption)
                    .foregroundStyle(.secondary)
                    .padding(.vertical, 4)
            } else {
                ScrollView {
                    VStack(alignment: .leading, spacing: 0) {
                        ForEach(state.domains, id: \.self) { domain in
                            DomainRow(domain: domain) {
                                state.removeDomain(domain)
                            }
                        }
                    }
                }
                .frame(maxHeight: 160)
            }

            AddDomainRow(
                addLabel: state.t("domains.add"),
                placeholder: state.t("domains.add.placeholder")
            ) { newDomain in
                state.addDomain(newDomain)
            }
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

// MARK: - Domain satırı

private struct DomainRow: View {
    let domain: String
    let onRemove: () -> Void
    @State private var hovering = false

    var body: some View {
        HStack {
            Text(domain)
                .font(.body)
            Spacer()
            Button(action: onRemove) {
                Image(systemName: "xmark.circle.fill")
                    .foregroundStyle(.secondary)
            }
            .buttonStyle(.plain)
        }
        .padding(.vertical, 5)
        .padding(.horizontal, 6)
        .background(hovering ? AnyShapeStyle(.quinary) : AnyShapeStyle(.clear),
                    in: RoundedRectangle(cornerRadius: 6, style: .continuous))
        .onHover { hovering = $0 }
    }
}

// MARK: - Alan adı ekleme satırı

private struct AddDomainRow: View {
    let addLabel: String
    let placeholder: String
    let onAdd: (String) -> Void
    @State private var editing = false
    @State private var text = ""
    @FocusState private var focused: Bool

    var body: some View {
        if editing {
            HStack(spacing: 6) {
                Image(systemName: "globe")
                    .foregroundStyle(.secondary)
                TextField(placeholder, text: $text)
                    .textFieldStyle(.plain)
                    .font(.body)
                    .focused($focused)
                    .onSubmit(commit)
                Button(action: commit) {
                    Image(systemName: "return")
                        .foregroundStyle(.secondary)
                }
                .buttonStyle(.plain)
                .disabled(text.trimmingCharacters(in: .whitespaces).isEmpty)
            }
            .padding(.vertical, 5)
            .padding(.horizontal, 6)
            .background(.quaternary, in: RoundedRectangle(cornerRadius: 6, style: .continuous))
            .onAppear { focused = true }
        } else {
            Button {
                editing = true
            } label: {
                HStack(spacing: 6) {
                    Image(systemName: "plus.circle.fill")
                        .foregroundStyle(.green)
                    Text(addLabel)
                        .font(.body)
                    Spacer()
                }
                .padding(.vertical, 5)
                .padding(.horizontal, 6)
                .contentShape(Rectangle())
            }
            .buttonStyle(.plain)
        }
    }

    private func commit() {
        let trimmed = text.trimmingCharacters(in: .whitespaces)
        guard !trimmed.isEmpty else { return }
        onAdd(trimmed)
        text = ""
        editing = false
    }
}
