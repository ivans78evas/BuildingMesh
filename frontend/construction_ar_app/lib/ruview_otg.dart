import 'dart:async';
import 'package:flutter_libserialport/flutter_libserialport.dart';

class RuViewOTGService {
  SerialPort? _port;
  StreamSubscription? _subscription;

  Future<void> connect(String portName) async {
    _port = SerialPort(portName);
    if (_port!.openReadWrite()) {
      print("Connected to ESP32-S3 via OTG");
      // Слушаем поток данных от RuView
      _subscription = _port!.config.stream.listen((data) {
        _parseRuViewData(data);
      });
    }
  }

  void _parseRuViewData(dynamic data) {
    // Логика обработки WiFi Sensing данных (RuView)
    // Определение поз, расстояний и типов объектов за стеной
  }

  void dispose() {
    _subscription?.cancel();
    _port?.close();
  }
}
