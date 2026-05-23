import 'package:flutter/material.dart';
import 'package:arcore_flutter_plugin/arcore_flutter_plugin.dart';
import 'package:vector_math/vector_math_64.dart' as vector;

class ARWallOverlay extends StatefulWidget {
  const ARWallOverlay({super.key});

  @override
  State<ARWallOverlay> createState() => _ARWallOverlayState();
}

class _ARWallOverlayState extends State<ARWallOverlay> {
  late ArCoreController arCoreController;

  void _onArCoreViewCreated(ArCoreController controller) {
    arCoreController = controller;
    _addAnchorMarker(controller);
  }

  void _addAnchorMarker(ArCoreController controller) {
    // В реальности здесь будет поиск QR-кода на щитке
    final node = ArCoreNode(
      shape: ArCoreCube(
        materials: [ArCoreMaterial(color: Colors.blue)],
        size: vector.Vector3(0.1, 0.1, 0.1),
      ),
      position: vector.Vector3(0, 0, -1),
      name: "Anchor_Electric_Panel",
    );
    controller.addArCoreNode(node);
  }


  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('AR X-Ray стен')),
      body: ArCoreView(
        onArCoreViewCreated: _onArCoreViewCreated,
        enableUpdateListener: true,
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          // Ручное выравнивание: смещение модели
        },
        child: const Icon(Icons.sync),
      ),
    );
  }

  @override
  void dispose() {
    arCoreController.dispose();
    super.dispose();
  }
}
