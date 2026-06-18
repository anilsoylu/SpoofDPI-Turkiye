import SwiftUI

// MARK: - Parlayan kalkan hero bileşeni

struct ShieldHero: View {
    let running: Bool
    var size: CGFloat = 240
    var animated: Bool = true

    @State private var pulseScale: CGFloat = 1.0
    @State private var pulseOpacity: Double = 0.6

    private var glowColor: Color { running ? .dsTeal.opacity(0.4) : .clear }

    // Orantılı boyutlar (tüm sabitler size'a göre ölçeklendi)
    private var iconSize:   CGFloat { size * 0.40 }   // 240 → 96
    private var halo0:      CGFloat { size * 0.583 }  // 240 → 140
    private var haloStep:   CGFloat { size * 0.183 }  // 240 → 44
    private var gradRadius: CGFloat { size * 0.50 }   // 240 → 120
    private var gradStart:  CGFloat { size * 0.125 }  // 240 → 30

    var body: some View {
        ZStack {
            // Radyal halo katmanları
            if running {
                ForEach(0..<3, id: \.self) { i in
                    let d = halo0 + CGFloat(i) * haloStep
                    Circle()
                        .stroke(
                            Color.dsTeal.opacity(0.08 - Double(i) * 0.025),
                            lineWidth: 1.5
                        )
                        .frame(width: d, height: d)
                        .scaleEffect(pulseScale)
                        .opacity(pulseOpacity)
                }
                // Arka radyal gradient
                RadialGradient(
                    colors: [Color.dsTeal.opacity(0.18), .clear],
                    center: .center,
                    startRadius: gradStart,
                    endRadius: gradRadius
                )
                .frame(width: size, height: size)
                .blendMode(.screen)
            }

            // Kalkan ikonu
            Image(systemName: running ? "shield.fill" : "shield")
                .font(.system(size: iconSize, weight: .bold))
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
        .frame(width: size, height: size)
        .onAppear { startPulse() }
        .onChange(of: running) { startPulse() }
    }

    private func startPulse() {
        // Sık açılan menü panelinde sürekli hareket istemiyoruz:
        // animated=false ise glow statik kalır, pulse döngüsü yok.
        guard running && animated else {
            withAnimation(.easeOut(duration: 0.3)) {
                pulseScale = 1.0
                pulseOpacity = running ? 0.5 : 0.6
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
