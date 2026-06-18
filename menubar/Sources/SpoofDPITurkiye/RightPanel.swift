import SwiftUI

// MARK: - Sağ panel

struct RightPanel: View {
    @EnvironmentObject private var state: AppState

    var body: some View {
        ScrollView(.vertical, showsIndicators: false) {
            VStack(spacing: 16) {
                DomainsCard()
                    .environmentObject(state)

                ConnectionTestCard()
                    .environmentObject(state)

                Spacer(minLength: 20)
            }
            .padding(20)
        }
        .background(Color.dsBackground)
    }
}
