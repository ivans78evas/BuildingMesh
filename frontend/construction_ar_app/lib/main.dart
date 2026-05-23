import 'package:flutter/material.dart';
import 'package:construction_ar_app/advanced_ar_overlay.dart';
import 'package:construction_ar_app/project_map.dart';

void main() {
  runApp(const ConstructionApp());
}

class ConstructionApp extends StatelessWidget {
  const ConstructionApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'AR Build SaaS',
      theme: ThemeData(
        brightness: Brightness.dark,
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      home: const DashboardPage(),
    );
  }
}

class DashboardPage extends StatelessWidget {
  const DashboardPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Мои Объекты'),
        actions: [
          IconButton(icon: const Icon(Icons.person), onPressed: () {}),
        ],
      ),
      body: GridView.count(
        crossAxisCount: 2,
        padding: const EdgeInsets.all(16),
        children: [
          _buildProjectCard(
            context,
            'ЖК Центр',
            'ул. Ленина, 10',
            Icons.apartment,
            () => Navigator.push(context, MaterialPageRoute(builder: (c) => const AdvancedAROverlay())),
          ),
          _buildProjectCard(
            context,
            'БЦ Квартал',
            'пр. Мира, 45',
            Icons.business,
            () => Navigator.push(context, MaterialPageRoute(builder: (c) => const ProjectMapView())),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        label: const Text('Новый объект'),
        icon: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildProjectCard(BuildContext context, String title, String subtitle, IconData icon, VoidCallback onTap) {
    return Card(
      child: InkWell(
        onTap: onTap,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, size: 48, color: Colors.blue),
            const SizedBox(height: 8),
            Text(title, style: Theme.of(context).textTheme.titleMedium),
            Text(subtitle, style: Theme.of(context).textTheme.bodySmall),
          ],
        ),
      ),
    );
  }
}
