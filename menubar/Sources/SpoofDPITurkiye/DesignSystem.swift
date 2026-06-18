import SwiftUI

// MARK: - Renk paleti

extension Color {
    static let dsBackground      = Color(hex: "#0E0F13")
    static let dsPanel           = Color(hex: "#16181C")
    static let dsBorder          = Color(hex: "#262A31")
    static let dsPrimaryText     = Color.white
    static let dsSecondaryText   = Color(hex: "#8A8F98")
    static let dsTeal            = Color(hex: "#34D3A6")
    static let dsTealDark        = Color(hex: "#2BB894")
    static let dsIndigo          = Color(hex: "#6C7CF0")
    static let dsIndigoDark      = Color(hex: "#8B7DF0")
    static let dsCoral           = Color(hex: "#FF6B5E")
    static let dsCoralDark       = Color(hex: "#F0584A")
    static let dsRowBg           = Color(hex: "#1B1E24")

    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 6:  (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8:  (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default: (a, r, g, b) = (255, 0, 0, 0)
        }
        self.init(
            .sRGB,
            red:   Double(r) / 255,
            green: Double(g) / 255,
            blue:  Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}

// MARK: - Kart ViewModifier

struct CardModifier: ViewModifier {
    func body(content: Content) -> some View {
        content
            .background(Color.dsPanel)
            .clipShape(RoundedRectangle(cornerRadius: 18, style: .continuous))
            .overlay(
                RoundedRectangle(cornerRadius: 18, style: .continuous)
                    .stroke(Color.dsBorder, lineWidth: 1)
            )
            .shadow(color: .black.opacity(0.25), radius: 8, x: 0, y: 4)
    }
}

extension View {
    func dsCard() -> some View {
        modifier(CardModifier())
    }
}

// MARK: - Küçük başlık metni

struct SectionLabel: View {
    let text: String
    var icon: String? = nil

    var body: some View {
        HStack(spacing: 6) {
            if let icon {
                Image(systemName: icon)
                    .font(.system(size: 11, weight: .semibold))
                    .foregroundStyle(Color.dsSecondaryText)
            }
            Text(text)
                .font(.system(size: 11, weight: .semibold))
                .tracking(1.2)
                .foregroundStyle(Color.dsSecondaryText)
        }
    }
}

// MARK: - Birincil buton stili (gradient + glow)

struct PrimaryButtonStyle: ButtonStyle {
    enum Variant { case teal, coral, indigo }
    let variant: Variant

    private var gradient: LinearGradient {
        switch variant {
        case .teal:
            return LinearGradient(
                colors: [.dsTeal, .dsTealDark],
                startPoint: .topLeading, endPoint: .bottomTrailing
            )
        case .coral:
            return LinearGradient(
                colors: [.dsCoral, .dsCoralDark],
                startPoint: .topLeading, endPoint: .bottomTrailing
            )
        case .indigo:
            return LinearGradient(
                colors: [.dsIndigo, .dsIndigoDark],
                startPoint: .topLeading, endPoint: .bottomTrailing
            )
        }
    }

    private var glowColor: Color {
        switch variant {
        case .teal:   return .dsTeal.opacity(0.45)
        case .coral:  return .dsCoral.opacity(0.45)
        case .indigo: return .dsIndigo.opacity(0.45)
        }
    }

    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .font(.system(size: 15, weight: .semibold))
            .foregroundStyle(.black.opacity(0.85))
            .frame(maxWidth: .infinity)
            .frame(height: 50)
            .background(gradient)
            .clipShape(RoundedRectangle(cornerRadius: 14, style: .continuous))
            .shadow(color: glowColor, radius: configuration.isPressed ? 4 : 12, x: 0, y: 0)
            .scaleEffect(configuration.isPressed ? 0.97 : 1)
            .animation(.easeOut(duration: 0.12), value: configuration.isPressed)
    }
}

// MARK: - İkincil (outline pill) buton stili

struct OutlinePillButtonStyle: ButtonStyle {
    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .font(.system(size: 13, weight: .medium))
            .foregroundStyle(Color.dsSecondaryText)
            .padding(.horizontal, 14)
            .padding(.vertical, 9)
            .background(Color.dsPanel)
            .clipShape(Capsule())
            .overlay(
                Capsule().stroke(Color.dsBorder, lineWidth: 1)
            )
            .scaleEffect(configuration.isPressed ? 0.96 : 1)
            .animation(.easeOut(duration: 0.1), value: configuration.isPressed)
    }
}
