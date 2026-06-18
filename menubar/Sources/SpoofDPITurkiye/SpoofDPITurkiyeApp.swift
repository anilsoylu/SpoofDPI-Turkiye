import SwiftUI
import AppKit

@main
struct SpoofDPITurkiyeApp: App {
    @StateObject private var state = AppState()

    var body: some Scene {
        MenuBarExtra("SpoofDPI Türkiye", systemImage: "shield.lefthalf.filled") {
            ContentView()
                .environmentObject(state)
        }
        .menuBarExtraStyle(.window)
    }
}
