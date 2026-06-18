import SwiftUI

// MARK: - Ana pencere görünümü

struct MainView: View {
    @EnvironmentObject private var state: AppState

    var body: some View {
        HStack(spacing: 0) {
            // Sol panel ~38%
            LeftPanel()
                .frame(width: 395)
                .frame(maxHeight: .infinity)

            // Dikey ayraç
            Rectangle()
                .fill(Color.dsBorder)
                .frame(width: 1)
                .frame(maxHeight: .infinity)

            // Sağ panel ~62%
            RightPanel()
                .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
        .frame(width: 1040, height: 720)
        .background(Color.dsBackground)
        .preferredColorScheme(.dark)
    }
}
