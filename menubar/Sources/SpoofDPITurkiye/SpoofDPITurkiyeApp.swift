import SwiftUI
import AppKit

// MARK: - AppDelegate (dock ikonu kapatır)

final class AppDelegate: NSObject, NSApplicationDelegate {
    func applicationDidFinishLaunching(_ notification: Notification) {
        // Dock ikonu ve uygulama menüsünü gizle; sadece menü çubuğu simgesi
        NSApp.setActivationPolicy(.accessory)
    }
}

// MARK: - Uygulama

@main
struct SpoofDPITurkiyeApp: App {
    @NSApplicationDelegateAdaptor(AppDelegate.self) private var appDelegate
    @StateObject private var state = AppState()

    var body: some Scene {
        // Tek arayüz: menü çubuğu popover paneli (ayrı pencere yok)
        MenuBarExtra {
            MenuBarPanel()
                .environmentObject(state)
        } label: {
            Image(systemName: state.running ? "shield.fill" : "shield")
        }
        .menuBarExtraStyle(.window)
    }
}
