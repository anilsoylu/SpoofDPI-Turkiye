import SwiftUI
import AppKit

// MARK: - Ayarlar sheet'i (native Form)

struct SettingsSheet: View {
    @EnvironmentObject private var state: AppState
    @Environment(\.dismiss) private var dismiss

    @State private var portValue: Int = 8080
    @State private var showUninstallAlert = false

    var body: some View {
        VStack(spacing: 0) {
            // Başlık çubuğu
            HStack {
                Text(state.t("settings.title"))
                    .font(.headline)
                Spacer()
                Button(state.t("done")) { dismiss() }
                    .keyboardShortcut(.defaultAction)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)

            Divider()

            Form {
                Section {
                    // Port
                    HStack {
                        Text(state.t("settings.port"))
                        Spacer()
                        TextField("", value: $portValue, format: .number.grouping(.never))
                            .textFieldStyle(.roundedBorder)
                            .frame(width: 70)
                            .multilineTextAlignment(.trailing)
                        Stepper("", value: $portValue, in: 1...65535)
                            .labelsHidden()
                    }
                    .onChange(of: portValue) { _, newValue in
                        if newValue != state.port {
                            state.setPort(newValue)
                        }
                    }

                    // Dil
                    Picker(state.t("settings.language"), selection: $state.lang) {
                        ForEach(Lang.allCases, id: \.self) { lang in
                            Text(lang.rawValue).tag(lang)
                        }
                    }
                    .pickerStyle(.segmented)
                }

                Section {
                    LabeledContent(state.t("settings.version"),
                                   value: state.version.isEmpty ? "—" : state.version)
                }

                Section {
                    Button(role: .destructive) {
                        showUninstallAlert = true
                    } label: {
                        Label(state.t("settings.uninstall"), systemImage: "trash")
                    }
                    .disabled(state.busy)

                    Button {
                        NSApplication.shared.terminate(nil)
                    } label: {
                        Label(state.t("settings.quit"), systemImage: "power")
                    }
                }
            }
            .formStyle(.grouped)
        }
        .frame(width: 360, height: 360)
        .onAppear { portValue = state.port }
        .alert(state.t("uninstall.confirm"), isPresented: $showUninstallAlert) {
            Button(state.t("uninstall.confirm.btn"), role: .destructive) {
                state.uninstall()
                dismiss()
            }
            Button(state.t("cancel"), role: .cancel) {}
        }
    }
}
