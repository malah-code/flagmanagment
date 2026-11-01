import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'evaluator.dart';

class FlagClient {
  final String apiKey;
  final String apiUrl;
  Map<String, dynamic> _context = {};
  Map<String, dynamic> _flags = {};

  final _changeController = StreamController<String>.broadcast();
  Stream<String> get onFlagChange => _changeController.stream;

  FlagClient({
    required this.apiKey,
    required this.apiUrl,
    Map<String, dynamic>? initialContext,
  }) {
    if (initialContext != null) {
      _context = initialContext;
    }
  }

  Future<void> connect() async {
    await fetchEvaluations();
  }

  void setContext(Map<String, dynamic> newContext) {
    _context = {..._context, ...newContext};
    fetchEvaluations();
  }

  Future<void> fetchEvaluations() async {
    try {
      final uri = Uri.parse('$apiUrl/api/v1/sdk/evaluations');
      final res = await http.post(
        uri,
        headers: {
          'Content-Type': 'application/json',
          'X-Client-Key': apiKey,
        },
        body: jsonEncode({'context': _context}),
      );

      if (res.statusCode == 200) {
        final data = jsonDecode(res.body) as Map<String, dynamic>;
        _flags = data['flags'] ?? data;
        _changeController.add('ALL');
      }
    } catch (e) {
      // offline fallback
    }
  }

  bool getBooleanValue(String flagKey, bool defaultValue) {
    final flag = _flags[flagKey];
    if (flag is Map<String, dynamic>) {
      final res = FlagEvaluator.evaluate<bool>(
        flag: flag,
        defaultValue: defaultValue,
        targetingKey: _context['targetingKey'] as String?,
      );
      return res.value;
    }
    return defaultValue;
  }

  String getStringValue(String flagKey, String defaultValue) {
    final flag = _flags[flagKey];
    if (flag is Map<String, dynamic>) {
      final res = FlagEvaluator.evaluate<String>(
        flag: flag,
        defaultValue: defaultValue,
        targetingKey: _context['targetingKey'] as String?,
      );
      return res.value;
    }
    return defaultValue;
  }

  void dispose() {
    _changeController.close();
  }
}
