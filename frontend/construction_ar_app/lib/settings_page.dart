import 'package:flutter/material.dart';

class SettingsPage extends StatefulWidget {
  const SettingsPage({super.key});

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  final _apiKeyController = TextEditingController();
  String _selectedOrg = 'Organization A';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Настройки SaaS')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          ListTile(
            title: const Text('Организация'),
            trailing: DropdownButton<String>(
              value: _selectedOrg,
              onChanged: (val) => setState(() => _selectedOrg = val!),
              items: ['Organization A', 'Organization B'].map((e) => DropdownMenuItem(value: e, child: Text(e))).toList(),
            ),
          ),
          const Divider(),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 8.0),
            child: TextField(
              controller: _apiKeyController,
              decoration: const InputDecoration(
                labelText: 'Google Cloud API Key',
                border: OutlineInputBorder(),
                helperText: 'Используется для Maps, Vision и Gemini',
              ),
              obscureText: true,
            ),
          ),
          const SizedBox(height: 20),
          ElevatedButton(
            onPressed: () {
              ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Настройки сохранены')));
            },
            child: const Text('Сохранить'),
          ),
        ],
      ),
    );
  }
}
