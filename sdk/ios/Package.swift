// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "FlagManagment",
    platforms: [
        .iOS(.v15),
        .macOS(.v12)
    ],
    products: [
        .library(
            name: "FlagManagment",
            targets: ["FlagManagment"]),
    ],
    dependencies: [
        .package(url: "https://github.com/open-feature/swift-sdk.git", from: "0.1.0"),
    ],
    targets: [
        .target(
            name: "FlagManagment",
            dependencies: [
                .product(name: "OpenFeature", package: "swift-sdk")
            ]),
        .testTarget(
            name: "FlagManagmentTests",
            dependencies: ["FlagManagment"]),
    ]
)
