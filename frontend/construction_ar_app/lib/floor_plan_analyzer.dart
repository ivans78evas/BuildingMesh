import 'dart:math';
import 'package:flutter/material.dart';

class CVWallDetector {
  // Прототип алгоритма детектирования стен.
  // В продакшене будет использовать OpenCV через FFI или TensorFlow Lite.
  static Future<List<Rect>> detectWalls(String imagePath) async {
    // Имитация задержки обработки на GPU
    await Future.delayed(const Duration(milliseconds: 500));

    // Возвращаем примеры найденных стен
    return [
      const Rect.fromLTWH(50, 100, 200, 10),
      const Rect.fromLTWH(50, 100, 10, 300),
    ];
  }
}

class FloorPlanAnalyzer extends StatefulWidget {
  final String imagePath;
  const FloorPlanAnalyzer({super.key, required this.imagePath});

  @override
  State<FloorPlanAnalyzer> createState() => _FloorPlanAnalyzerState();
}

class _FloorPlanAnalyzerState extends State<FloorPlanAnalyzer> {
  List<Rect> detectedWalls = [];
  bool isLoading = false;

  void _runDetection() async {
    setState(() => isLoading = true);
    final walls = await CVWallDetector.detectWalls(widget.imagePath);
    setState(() {
      detectedWalls = walls;
      isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('AI Анализ плана')),
      body: Stack(
        children: [
          Image.network(widget.imagePath),
          if (isLoading) const Center(child: CircularProgressIndicator()),
          CustomPaint(
            painter: DetectionPainter(detectedWalls),
            size: Size.infinite,
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _runDetection,
        child: const Icon(Icons.analytics),
      ),
    );
  }
}

class DetectionPainter extends CustomPainter {
  final List<Rect> rects;
  DetectionPainter(this.rects);

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = Colors.green.withOpacity(0.5)
      ..style = PaintingStyle.fill;

    for (var rect in rects) {
      canvas.drawRect(rect, paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => true;
}
