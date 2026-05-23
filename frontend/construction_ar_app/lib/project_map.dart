import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';

class ProjectMapView extends StatefulWidget {
  const ProjectMapView({super.key});

  @override
  State<ProjectMapView> createState() => _ProjectMapViewState();
}

class _ProjectMapViewState extends State<ProjectMapView> {
  static const CameraPosition _initialPosition = CameraPosition(
    target: LatLng(55.7558, 37.6173), // Пример: Москва
    zoom: 18,
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Карта объектов Google 3D')),
      body: GoogleMap(
        initialCameraPosition: _initialPosition,
        mapType: MapType.hybrid, // Для 3D вида со спутника
        markers: {
          const Marker(
            markerId: MarkerId('proj1'),
            position: LatLng(55.7558, 37.6173),
            infoWindow: InfoWindow(title: 'Объект: ЖК Центр'),
          ),
        },
      ),
    );
  }
}
