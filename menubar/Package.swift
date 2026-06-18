// swift-tools-version: 6.0
import PackageDescription

let package = Package(
    name: "SpoofDPITurkiye",
    platforms: [.macOS(.v14)],
    targets: [
        .executableTarget(name: "SpoofDPITurkiye", path: "Sources/SpoofDPITurkiye")
    ]
)
