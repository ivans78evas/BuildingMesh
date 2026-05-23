import 'package:flutter/material.dart';

class LayerManager extends StatelessWidget {
  final Function(String, bool) onLayerToggle;
  final Map<String, bool> activeLayers;

  const LayerManager({
    super.key,
    required this.onLayerToggle,
    required this.activeLayers,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.black.withValues(alpha: 0.8),
        borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Text('Управление слоями', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
          const Divider(),
          _buildLayerTile('BIM: Трубы', 'pipes'),
          _buildLayerTile('BIM: Электрика', 'electric'),
          _buildLayerTile('3DGS: Скан реальности', 'splat'),
          _buildLayerTile('RuView: WiFi Сенсоры', 'ruview'),
        ],
      ),
    );
  }

  Widget _buildLayerTile(String title, String key) {
    return CheckboxListTile(
      title: Text(title),
      value: activeLayers[key] ?? true,
      onChanged: (val) => onLayerToggle(key, val ?? false),
    );
  }
}
