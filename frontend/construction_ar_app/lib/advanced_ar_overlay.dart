import 'package:flutter/material.dart';
import 'package:arcore_flutter_plugin/arcore_flutter_plugin.dart';
import 'package:vector_math/vector_math_64.dart' as vector;

class AdvancedAROverlay extends StatefulWidget {
  const AdvancedAROverlay({super.key});

  @override
  State<AdvancedAROverlay> createState() => _AdvancedAROverlayState();
}

class _AdvancedAROverlayState extends State<AdvancedAROverlay> {
  late ArCoreController arCoreController;

  void _onArCoreViewCreated(ArCoreController controller) {
    arCoreController = controller;
    _loadBIMLayers();
    _loadRuViewGhosts();
  }

  void _loadBIMLayers() {
    // Отображение труб и кабелей из BIM
    final pipe = ArCoreNode(
      shape: ArCoreCylinder(
        materials: [ArCoreMaterial(color: Colors.blue)],
        radius: 0.05,
        height: 2.0,
      ),
      position: vector.Vector3(0, 1, -2),
    );
    arCoreController.addArCoreNode(pipe);
  }

  void _loadRuViewGhosts() {
    // Визуализация объектов обнаруженных RuView (через ESP32-S3)
    final ghost = ArCoreNode(
      shape: ArCoreSphere(
        materials: [ArCoreMaterial(color: Colors.red.withOpacity(0.5))],
        radius: 0.3,
      ),
      position: vector.Vector3(1, 1, -3), // Объект за стеной
      name: "RuView_Person_Ghost",
    );
    arCoreController.addArCoreNode(ghost);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Multi-Layer AR (BIM+Splat+RuView)')),
      body: ArCoreView(onArCoreViewCreated: _onArCoreViewCreated),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          // Использование Google Geospatial API для перекалибровки
        },
        child: const Icon(Icons.public),
      ),
    );
  }
}
