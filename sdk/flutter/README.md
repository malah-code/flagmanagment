# FlagManagment — Flutter & Dart SDK

Official Flutter and Dart client SDK for **FlagManagment**, providing offline caching, contextual user targeting, and sub-millisecond local evaluations for cross-platform iOS & Android mobile apps.

---

## Installation

Add to your `pubspec.yaml`:
```yaml
dependencies:
  flagmanagment_flutter: ^1.0.0
```

---

## Quickstart

```dart
import 'package:flagmanagment_flutter/flagmanagment_flutter.dart';

void main() async {
  final client = FlagClient(
    apiKey: 'flg_client_your_mobile_key',
    apiUrl: 'https://flags.example.com',
    initialContext: {
      'targetingKey': 'mobile-user-88',
      'tier': 'pro',
    },
  );

  await client.connect();

  // Evaluate flags synchronously without network lag
  final isDarkMode = client.getBooleanValue('dark-mode', false);
  final promoBanner = client.getStringValue('promotional-banner', 'Welcome!');

  print('Dark mode active: $isDarkMode');
}
```
