import 'package:flutter/material.dart';
import 'package:photo_view/photo_view.dart';

class FloorPlanEditor extends StatefulWidget {
  final String imagePath;

  const FloorPlanEditor({super.key, required this.imagePath});

  @override
  State<FloorPlanEditor> createState() => _FloorPlanEditorState();
}

class _FloorPlanEditorState extends State<FloorPlanEditor> {
  List<Offset> wallPoints = [];
  bool isCalibrating = false;
  double scale = 1.0; // pixels per meter

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Разметка стен на плане'),
        actions: [
          IconButton(
            icon: const Icon(Icons.square_foot),
            onPressed: () => setState(() => isCalibrating = !isCalibrating),
            color: isCalibrating ? Colors.orange : null,
          ),
          IconButton(
            icon: const Icon(Icons.save),
            onPressed: () {
              // Сохранение геометрии и отправка на сервер
            },
          ),
        ],
      ),
      body: GestureDetector(
        onTapUp: (details) {
          setState(() {
            wallPoints.add(details.localPosition);
          });
        },
        child: Stack(
          children: [
            PhotoView.customChild(
              child: Image.network(widget.imagePath),
            ),
            CustomPaint(
              painter: WallPainter(wallPoints),
              size: Size.infinite,
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => setState(() => wallPoints.clear()),
        child: const Icon(Icons.delete),
      ),
    );
  }
}

class WallPainter extends CustomPainter {
  final List<Offset> points;
  WallPainter(this.points);

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = Colors.red
      ..strokeWidth = 5.0
      ..strokeCap = StrokeCap.round;

    for (int i = 0; i < points.length - 1; i += 2) {
      if (i + 1 < points.length) {
        canvas.drawLine(points[i], points[i + 1], paint);
      }
    }

    for (var point in points) {
      canvas.drawCircle(point, 6.0, paint..color = Colors.blue);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => true;
}
