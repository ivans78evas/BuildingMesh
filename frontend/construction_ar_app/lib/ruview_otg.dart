import 'dart:async';
import 'dart:typed_data';
import 'package:flutter_libserialport/flutter_libserialport.dart';

class RuViewOTGService {
  SerialPort? _port;
  StreamSubscription? _subscription;
  final _controller = StreamController<Uint8List>();

  Stream<Uint8List> get dataStream => _controller.stream;

  Future<void> connect(String portName) async {
    _port = SerialPort(portName);
    if (_port!.openReadWrite()) {
      // Использование SerialPortReader для чтения потока данных
      final reader = SerialPortReader(_port!);
      _subscription = reader.stream.listen((data) {
        _controller.add(data);
        _parseRuViewData(data);
      });
    }
  }

  void _parseRuViewData(Uint8List data) {
    // Логика обработки WiFi Sensing данных (RuView)
  }

  void dispose() {
    _subscription?.cancel();
    _port?.close();
    _controller.close();
  }
}
