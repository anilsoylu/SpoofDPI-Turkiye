import SwiftUI
import AppKit

// MARK: - Ayarlar ekranı (popover içi inline — sheet YOK, BUG4 fix)

struct SettingsScreen: View {
    @EnvironmentObject private var state: AppState
    var onBack: () -> Void

    @State private var portValue: Int = 8080
    @State private var showUninstallAlert = false

    // Kaydedilmemiş port değişikliği var mı? Sadece farklıysa "Uygula" aktif.
    private var portDirty: Bool { portValue != state.port }

    // Port geçerli mi? (#6) UI, CLI'a göndermeden ÖNCE 1-65535 aralığını doğrular;
    // geçersizse "Uygula" devre dışı kalır ve inline hata gösterilir (sessiz
    // fail / desync yerine). Go tarafı da ValidatePort yapar (defense-in-depth).
    private var portValid: Bool { portValue >= 1 && portValue <= 65535 }

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
                Text(state.t("settings.title"))
                    .font(.headline)
                Spacer()
                // Simetri için görünmez yer tutucu.
                Label(state.t("done"), systemImage: "chevron.left")
                    .labelStyle(.titleAndIcon)
                    .opacity(0)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)

            Divider()

            Form {
                Section {
                    // Port — değer değiştir + "Uygula" ile async setPort
                    // (her tuş vuruşunda servis yeniden başlatılmaz).
                    HStack {
                        Text(state.t("settings.port"))
                        Spacer()
                        TextField("", value: $portValue, format: .number.grouping(.never))
                            .textFieldStyle(.roundedBorder)
                            .frame(width: 70)
                            .multilineTextAlignment(.trailing)
                            .disabled(state.busy)
                        Stepper("", value: $portValue, in: 1...65535)
                            .labelsHidden()
                            .disabled(state.busy)
                    }

                    // Geçersiz port (aralık dışı) inline uyarısı (#6).
                    if !portValid {
                        Text(state.t("settings.port.invalid"))
                            .font(.caption)
                            .foregroundStyle(.red)
                    }

                    HStack {
                        Spacer()
                        if state.busy {
                            ProgressView().controlSize(.small)
                        }
                        Button(state.t("btn.apply")) {
                            // UI'da doğrulanmadan CLI'a gönderilmez (#6); buton
                            // zaten geçersizken devre dışı, yine de savunma için
                            // guard ile kontrol et.
                            guard portValid else { return }
                            state.setPort(portValue)
                        }
                        .disabled(state.busy || !portDirty || !portValid)
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
        .frame(height: 380)
        .onAppear { portValue = state.port }
        // state.port harici değişirse (ör. refresh) alanı senkronla — ama
        // kullanıcı düzenlerken değil: yalnızca busy değilken ve dirty değilken.
        .onChange(of: state.port) { _, newValue in
            if !state.busy { portValue = newValue }
        }
        .alert(state.t("uninstall.confirm"), isPresented: $showUninstallAlert) {
            Button(state.t("uninstall.confirm.btn"), role: .destructive) {
                state.uninstall()
                onBack()
            }
            Button(state.t("cancel"), role: .cancel) {}
        }
    }
}
