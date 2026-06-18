import SwiftUI
import AppKit

@main
struct SpoofDPITurkiyeApp: App {
    @StateObject private var state = AppState()

    var body: some Scene {
        // Ana pencere
        Window("SpoofDPI Türkiye", id: "main") {
            MainView()
                .environmentObject(state)
        }
        .windowStyle(.hiddenTitleBar)
        .defaultSize(width: 1040, height: 720)
        .windowResizability(.contentSize)

        // Menu bar ikonu
        MenuBarExtra("SpoofDPI Türkiye", systemImage: "shield.lefthalf.filled") {
            MenuBarMenuView()
                .environmentObject(state)
        }
    }
}

// MARK: - Menu Bar küçük menüsü

struct MenuBarMenuView: View {
    @EnvironmentObject private var state: AppState
    @Environment(\.openWindow) private var openWindow

    var body: some View {
        Button(state.running ? state.t("btn.stop") : state.t("btn.start")) {
            state.toggle()
        }

        Divider()

        Button(state.t("menu.open")) {
            openWindow(id: "main")
            NSApplication.shared.activate(ignoringOtherApps: true)
        }

        Divider()

        Button(state.t("menu.quit")) {
            NSApplication.shared.terminate(nil)
        }
    }
}
