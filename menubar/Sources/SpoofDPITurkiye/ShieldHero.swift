import SwiftUI

// MARK: - Parlayan kalkan hero bileşeni

struct ShieldHero: View {
    let running: Bool

    @State private var pulseScale: CGFloat = 1.0
    @State private var pulseOpacity: Double = 0.6

    private var shieldColor: Color { running ? .dsTeal : Color.dsSecondaryText }
    private var glowColor: Color   { running ? .dsTeal.opacity(0.4) : .clear }

    var body: some View {
        ZStack {
            // Radyal halo katmanları
            if running {
                ForEach(0..<3, id: \.self) { i in
                    Circle()
                        .stroke(
                            Color.dsTeal.opacity(0.08 - Double(i) * 0.025),
                            lineWidth: 1.5
                        )
                        .frame(width: CGFloat(140 + i * 44), height: CGFloat(140 + i * 44))
                        .scaleEffect(pulseScale)
                        .opacity(pulseOpacity)
                }
                // Arka radyal gradient
                RadialGradient(
                    colors: [Color.dsTeal.opacity(0.18), .clear],
                    center: .center,
                    startRadius: 30,
                    endRadius: 120
                )
                .frame(width: 240, height: 240)
                .blendMode(.screen)
            }

            // Kalkan ikonu
            Image(systemName: running ? "shield.fill" : "shield")
                .font(.system(size: 96, weight: .bold))
                .foregroundStyle(
                    running
                        ? LinearGradient(
                            colors: [.dsTeal, .dsTealDark],
                            startPoint: .top, endPoint: .bottom
                          )
                        : LinearGradient(
                            colors: [Color.dsSecondaryText.opacity(0.7), Color.dsSecondaryText.opacity(0.4)],
                            startPoint: .top, endPoint: .bottom
                          )
                )
                .shadow(color: running ? glowColor : .clear, radius: 20, x: 0, y: 0)
                .scaleEffect(running ? pulseScale : 1.0)
        }
        .frame(width: 240, height: 240)
        .onAppear { startPulse() }
        .onChange(of: running) { startPulse() }
    }

    private func startPulse() {
        guard running else {
            withAnimation(.easeOut(duration: 0.3)) {
                pulseScale = 1.0
                pulseOpacity = 0.6
            }
            return
        }
        withAnimation(
            .easeInOut(duration: 1.8).repeatForever(autoreverses: true)
        ) {
            pulseScale = 1.06
            pulseOpacity = 0.3
        }
    }
}
