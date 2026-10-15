# Linux–Android Continuity: Brick-by-Brick Build Plan

As researched: 2026-10-10

## 1. Product and constraints

- [ ] 1. Define the first supported pair as one Android phone and one Linux desktop.
- [ ] 2. Freeze v1 features: secure pairing, file transfer, clipboard actions, notifications, and phone-screen mirroring.
- [ ] 3. Move media control, link handoff, photo sync, webcam/microphone, presence, and multi-device support to later milestones.
- [ ] 4. Exclude root access, hidden capture, arbitrary account/token scraping, and bypasses of Android consent.
- [ ] 5. Set the architecture: shared Rust networking core, Kotlin/Jetpack Compose Android app, Qt Quick/QML Linux UI, Rust Linux daemon.
- [ ] 6. Set Android compileSdk and targetSdk to API 37; choose minSdk after checking the actual target phone and support range.
- [ ] 7. Support Arch/Omarchy first; keep desktop protocol and UI layers separable for later Linux distributions.
- [ ] 8. Make local-network operation the default; make cloud accounts and a relay unnecessary.
- [ ] 9. Create a feature-to-permission matrix before requesting any Android permission.
- [ ] 10. Write user journeys for pairing, reconnecting, sending, receiving, revoking, and troubleshooting.

## 2. Workspace and tools

- [ ] 11. Create a monorepo for Rust, Android, Linux UI, protocol schemas, packaging, and docs.
- [ ] 12. Create Rust workspace crates for protocol, identity, transport, feature services, and platform adapters.
- [ ] 13. Create the Android app in Kotlin with the current stable Android Gradle Plugin and Gradle wrapper.
- [ ] 14. Create the Linux UI with Qt 6.12, Qt Quick Controls, QML, and CMake.
- [ ] 15. Pin Rust stable, Android SDK/build tools, Qt modules, and all dependency versions.
- [ ] 16. Use Rust 2024 edition where every selected dependency and tool supports it.
- [ ] 17. Select the newest Protobuf Edition supported by the selected Rust code generator and runtime.
- [ ] 18. Set up rustfmt, Clippy, Kotlin formatting/lint, qmllint, and C++ formatting.
- [ ] 19. Define module naming, error types, logging rules, and review conventions.
- [ ] 20. Add CI jobs for formatting, linting, Linux builds, Android builds, and dependency checks.
- [ ] 21. Generate a software bill of materials for each release build.
- [ ] 22. Add Rust dependency/license auditing and Android dependency/license reporting.
- [ ] 23. Choose the project license and record licenses for every bundled dependency.

## 3. Architecture and security foundation

- [ ] 24. Draw process boundaries for Android UI, Android service, shared Rust core, Linux daemon, and Linux UI.
- [ ] 25. Define typed interfaces between the shared core and each platform adapter.
- [ ] 26. Keep Android lifecycle and permissions in Kotlin; keep QUIC, protocol, and transfer state in Rust.
- [ ] 27. Keep the Linux daemon independent of the Qt UI and control it over the user D-Bus session.
- [ ] 28. Divide code into pairing, connection, files, clipboard, notifications, media, and screen-stream feature modules.
- [ ] 29. Model each feature as explicit states, commands, events, and recoverable errors.
- [ ] 30. Define one source of truth for paired devices, transfer queues, and user preferences.
- [ ] 31. Define database schemas, migrations, deletion, backup, and reset behavior.
- [ ] 32. Keep secrets, private keys, file contents, notification text, and clipboard data out of logs.
- [ ] 33. Store Android identity material using Android Keystore protection.
- [ ] 34. Store Linux identity material through the Freedesktop Secret Service or an encrypted user-owned fallback.
- [ ] 35. Generate device keys with vetted cryptographic libraries; never implement cryptographic primitives.
- [ ] 36. Write a threat model covering hostile LAN peers, QR interception, replay, compromised devices, and malicious files.
- [ ] 37. Define per-device and per-feature authorization, including revocation and key rotation.
- [ ] 38. Map mobile security work to the current OWASP MASVS controls.

## 4. Wire protocol and transport

- [ ] 39. Write a versioned protocol specification before implementing feature messages.
- [ ] 40. Define Protobuf schemas for identity metadata, capabilities, requests, responses, errors, and events.
- [ ] 41. Add protocol-version negotiation and explicit capability negotiation.
- [ ] 42. Add request IDs, cancellation, idempotency rules, and structured error codes.
- [ ] 43. Set strict limits for message sizes, nesting, file counts, names, and concurrent operations.
- [ ] 44. Separate small control messages from bulk file streams.
- [ ] 45. Add protocol fixtures shared by Linux and Android builds.
- [ ] 46. Implement IETF QUIC v1 with Tokio, Quinn, and rustls.
- [ ] 47. Require TLS 1.3 and mutually authenticate paired device identities.
- [ ] 48. Pin the verified peer identity after pairing; reject changed identities until re-paired.
- [ ] 49. Disable 0-RTT for pairing, file writes, clipboard writes, and remote actions.
- [ ] 50. Use independent QUIC streams for concurrent transfers and control operations.
- [ ] 51. Use QUIC datagrams only for latency-sensitive screen/audio frames that can be dropped.
- [ ] 52. Apply backpressure, bounded queues, transfer timeouts, and cancellation across every stream.
- [ ] 53. Add reconnect backoff, keepalive, network-change handling, and stale-connection cleanup.
- [ ] 54. Support IPv4 and IPv6 and refresh endpoints after DHCP or interface changes.
- [ ] 55. Record transport metrics without logging private payloads or stable network identifiers.

## 5. Linux daemon and desktop UI

- [ ] 56. Build the Rust daemon as a per-user service with no root requirement.
- [ ] 57. Add a systemd user unit with restart policy, startup ordering, and clean shutdown.
- [ ] 58. Apply systemd sandboxing, restrictive filesystem access, and a dedicated runtime directory.
- [ ] 59. Bind network listeners only to intended active LAN interfaces.
- [ ] 60. Document the minimum inbound UDP firewall rules for QUIC and discovery.
- [ ] 61. Expose a small, versioned D-Bus API using zbus.
- [ ] 62. Define D-Bus methods for pairing, device status, transfers, and feature settings.
- [ ] 63. Define D-Bus signals for connection, progress, permission, and error changes.
- [ ] 64. Build the Qt Quick UI as a separate desktop process.
- [ ] 65. Create a custom design system for spacing, typography, color, motion, icons, and states.
- [ ] 66. Design the device dashboard, pairing window, transfer queue, phone files, clipboard, and settings screens.
- [ ] 67. Make all screens work with keyboard navigation, scaling, screen readers, and reduced motion.
- [ ] 68. Use xdg-desktop-portal for user-selected files and sandbox-compatible access.
- [ ] 69. Add desktop notifications, launcher integration, and an optional tray/status item.

## 6. Android app and platform adapters

- [ ] 70. Build the Android UI with Kotlin, Jetpack Compose, Material 3, and adaptive layouts.
- [ ] 71. Use layered architecture, repositories, ViewModels, StateFlow, and unidirectional data flow.
- [ ] 72. Build the Kotlin-to-Rust bridge with JNI and keep all blocking/native work off the main thread.
- [ ] 73. Cross-compile the Rust core for the required Android ABIs and package native libraries in the APK.
- [ ] 74. Implement Android service lifecycle adapters separately from the shared Rust connection state machine.
- [ ] 75. Declare only permissions required by an enabled feature.
- [ ] 76. Declare and request ACCESS_LOCAL_NETWORK for Android 17/API 37 LAN discovery and connections.
- [ ] 77. Request local-network permission at the point the user enables pairing or LAN features.
- [ ] 78. Declare the correct foreground-service types and type-specific permissions for each long-running feature.
- [ ] 79. Use connectedDevice service behavior only while maintaining an active companion connection.
- [ ] 80. Use CompanionDeviceManager for supported pairing and presence flows.
- [ ] 81. Use WorkManager for deferrable work instead of a permanent background service.
- [ ] 82. Make the connection service visibly user-controlled, stoppable, and reflected by a persistent notification.
- [ ] 83. Request notification-listener, camera, nearby-device, and screen-capture access only when the matching feature is enabled.
- [ ] 84. Build permission-denied, permission-revoked, battery-restricted, and OEM-killed recovery states.
- [ ] 85. Use Android Keystore-protected keys and handle devices without StrongBox.
- [ ] 86. Use SAF for user-selected folders and MediaStore/photo picker for media access.
- [ ] 87. Open CameraX only during pairing and scan the short-lived QR invitation.
- [ ] 88. Verify install, upgrade, reboot, battery restriction, and permission recovery on the target Xiaomi phone.

## 7. Discovery and pairing

- [ ] 89. Publish a DNS-SD service over mDNS with only the minimum discovery metadata.
- [ ] 90. Discover peers over mDNS and treat every advertisement as untrusted input.
- [ ] 91. Keep QR/manual endpoint pairing available when multicast discovery is blocked.
- [ ] 92. Add optional BLE presence discovery only after the LAN connection flow works.
- [ ] 93. Generate a short-lived, single-use pairing invitation on the Linux desktop.
- [ ] 94. Put endpoint, protocol version, public identity, expiry, and a random nonce in the invitation.
- [ ] 95. Keep private keys, account data, file paths, and reusable credentials out of the QR payload.
- [ ] 96. Display a matching authentication code on both devices before trust is committed.
- [ ] 97. Require explicit user confirmation on both devices before pinning identities.
- [ ] 98. Rate-limit failed pairing attempts and expire unused invitations.
- [ ] 99. Implement pair, reconnect, rename, revoke, forget, and re-pair flows.
- [ ] 100. Add a manual address entry path with certificate fingerprint confirmation.
- [ ] 101. Show actionable connection diagnostics for discovery, routing, firewall, and peer isolation failures.

## 8. Files and network resilience

- [ ] 102. Build one-shot Android-to-Linux file sending through the Android share sheet.
- [ ] 103. Build Linux-to-Android sending through the desktop file chooser and Android receive confirmation.
- [ ] 104. Add remote folder browsing only within user-granted SAF tree URIs.
- [ ] 105. Preserve and revalidate Android SAF URI grants; handle revoked grants cleanly.
- [ ] 106. Implement chunked streaming, checksums, progress, cancellation, and resumable transfers.
- [ ] 107. Write incoming files to a temporary name and atomically finalize after validation.
- [ ] 108. Prevent path traversal, symlink escape, unsafe overwrite, and filename spoofing.
- [ ] 109. Add conflict handling, free-space checks, MIME detection, and user-selected destinations.
- [ ] 110. Add per-transfer authorization and a visible receive confirmation setting.
- [ ] 111. Test large files, interrupted transfers, repeated transfers, low disk space, and storage permission loss.
- [ ] 112. Add a local-only hotspot fallback for routers that isolate Wi-Fi clients.
- [ ] 113. Add user-guided USB tethering as an alternate local transport.

## 9. Clipboard, notifications, and media

- [ ] 114. Implement explicit PC-to-phone clipboard send and phone-to-PC clipboard send actions.
- [ ] 115. Respect Android background clipboard restrictions; never promise silent universal clipboard reads.
- [ ] 116. Require foreground or user-triggered clipboard access on Android.
- [ ] 117. Add sensitive-content filtering, per-device clipboard toggles, size limits, and short retention.
- [ ] 118. Mark sensitive clipboard items so Android can redact previews.
- [ ] 119. Add a user-controlled NotificationListenerService with clear onboarding and per-app filters.
- [ ] 120. Mirror only notification fields selected by the user and avoid storing notification bodies.
- [ ] 121. Forward notification actions through their original user-approved PendingIntents.
- [ ] 122. Add dismissal/reply only for notifications that expose supported actions.
- [ ] 123. Read media sessions through the supported Android media/notification APIs.
- [ ] 124. Implement play, pause, next, previous, and now-playing metadata with per-device controls.
- [ ] 125. Add optional battery/status reporting without collecting unrelated phone or account data.

## 10. Phone-screen mirroring

- [ ] 126. Start each screen-capture session from a visible user action and request fresh MediaProjection consent.
- [ ] 127. Declare the mediaProjection foreground-service type and stop it when capture ends.
- [ ] 128. Capture the user-selected screen or app window and honor protected/blank content.
- [ ] 129. Encode with Android MediaCodec hardware support and negotiate a supported codec and resolution.
- [ ] 130. Stream time-sensitive frames over QUIC datagrams with bounded latency and frame dropping.
- [ ] 131. Adapt frame rate and resolution to bandwidth, device temperature, and decode performance.
- [ ] 132. Decode with GStreamer on Linux and render frames in the Qt Quick UI.
- [ ] 133. Stop capture immediately on user stop, screen lock, consent revocation, disconnect, or projection callback.

## 11. Expanded continuity features

- [ ] 134. Add link and text handoff through Android share intents and Linux URI handlers.
- [ ] 135. Add photo sync as an opt-in, user-selected MediaStore/SAF workflow.
- [ ] 136. Add phone webcam/microphone streaming only after screen streaming and device lifecycle are stable.
- [ ] 137. Defer remote touch injection unless a reviewed, explicit Android accessibility use case justifies AccessibilityService.
- [ ] 138. Map account-related requests to official provider APIs, browser OAuth, or user-initiated share flows.
- [ ] 139. Add multiple paired devices with independent identities, permissions, and revoke controls.
- [ ] 140. Add optional proximity lock only after false-positive handling and safe recovery are designed.

## 12. Verification and release

- [ ] 141. Test across Android API levels, phone vendors, Wi-Fi bands, IPv4/IPv6, VPNs, and changing networks.
- [ ] 142. Test AP/client isolation, blocked multicast, DHCP changes, hotspot fallback, and USB tethering.
- [ ] 143. Fuzz protocol parsing, file metadata, QR payloads, and transfer path handling.
- [ ] 144. Review mobile storage, cryptography, network, platform, code, and privacy controls against OWASP MASVS.
- [ ] 145. Create signed release APKs, protect and back up signing keys, and register the package before broader distribution.
- [ ] 146. Build Arch packages, document clean install/update/uninstall, and publish a troubleshooting guide.

## Primary references

- [Android 17 and API 37](https://developer.android.com/about/versions/17/)
- [Android local-network permission](https://developer.android.com/privacy-and-security/local-network-permission)
- [Android foreground-service types](https://developer.android.com/develop/background-work/services/fgs/service-types)
- [Android CompanionDeviceManager](https://developer.android.com/reference/android/companion/CompanionDeviceManager)
- [Android architecture recommendations](https://developer.android.com/topic/architecture/recommendations)
- [Jetpack Compose](https://developer.android.com/develop/ui/compose)
- [Android file access and SAF](https://developer.android.com/training/data-storage/shared/documents-files)
- [Android clipboard](https://developer.android.com/develop/ui/views/touch-and-input/copy-paste)
- [Android notification listener](https://developer.android.com/reference/android/service/notification/NotificationListenerService)
- [Android MediaProjection](https://developer.android.com/media/grow/media-projection)
- [Android Keystore](https://developer.android.com/privacy-and-security/keystore)
- [Android local-only Wi-Fi hotspot](https://developer.android.com/develop/connectivity/wifi/localonlyhotspot)
- [Android developer verification](https://developer.android.com/developer-verification)
- [QUIC RFC 9000](https://www.rfc-editor.org/rfc/rfc9000)
- [TLS 1.3 RFC 8446](https://www.rfc-editor.org/rfc/rfc8446)
- [mDNS RFC 6762](https://www.rfc-editor.org/rfc/rfc6762)
- [DNS-SD RFC 6763](https://www.rfc-editor.org/rfc/rfc6763)
- [Protocol Buffers Editions](https://protobuf.dev/programming-guides/editions/)
- [Tokio](https://tokio.rs/tokio/tutorial)
- [Quinn](https://quinn-rs.github.io/quinn/)
- [Qt Quick](https://doc.qt.io/qt-6/qtquick-index.html)
- [D-Bus specification](https://dbus.freedesktop.org/doc/dbus-specification.html)
- [XDG Desktop Portal API](https://flatpak.github.io/xdg-desktop-portal/docs/api-reference.html)
- [XDG ScreenCast portal](https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.ScreenCast.html)
- [OWASP MASVS](https://mas.owasp.org/MASVS/)
