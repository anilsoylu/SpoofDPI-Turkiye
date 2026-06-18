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
        // ANA arayüz: menü çubuğu popover paneli
        MenuBarExtra {
            MenuBarPanel()
                .environmentObject(state)
        } label: {
            // Duruma göre ikon: açık=dolu kalkan, kapalı=içi boş kalkan
            Image(systemName: state.running ? "shield.lefthalf.filled" : "shield")
        }
        .menuBarExtraStyle(.window)

        // OPSİYONEL: tam pencere (otomatik açılmaz; popover'dan "Detaylar" ile açılır)
        Window("SpoofDPI Türkiye", id: "main") {
            MainView()
                .environmentObject(state)
        }
        .defaultSize(width: 1040, height: 720)
        .windowResizability(.contentSize)
    }
}
