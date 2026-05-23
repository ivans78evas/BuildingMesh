import 'package:flutter/material.dart';
import 'package:panorama/panorama.dart';

class PanoramaViewer extends StatelessWidget {
  final String imagePath;

  const PanoramaViewer({super.key, required this.imagePath});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Просмотр 360°'),
      ),
      body: Center(
        child: Panorama(
          animSpeed: 1.0,
          sensorControl: SensorControl.Orientation,
          child: Image.network(imagePath),
        ),
      ),
    );
  }
}
