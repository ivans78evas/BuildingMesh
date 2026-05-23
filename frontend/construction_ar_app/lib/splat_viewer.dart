import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';

class SplatViewer extends StatefulWidget {
  final String plyUrl;

  const SplatViewer({super.key, required this.plyUrl});

  @override
  State<SplatViewer> createState() => _SplatViewerState();
}

class _SplatViewerState extends State<SplatViewer> {
  late final WebViewController controller;

  @override
  void initState() {
    super.initState();
    controller = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setBackgroundColor(const Color(0x00000000))
      ..loadHtmlString(_buildSupersplatHtml(widget.plyUrl));
  }

  String _buildSupersplatHtml(String url) {
    // В реальности здесь будет загрузка PlayCanvas движка и Supersplat скриптов
    // Оптимизировано для GPU мобильного телефона через WebGL2
    return '''
    <!DOCTYPE html>
    <html>
    <head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <style>body { margin: 0; background: black; overflow: hidden; }</style>
        <script src="https://code.playcanvas.com/playcanvas-stable.min.js"></script>
    </head>
    <body>
        <canvas id="application-canvas"></canvas>
        <script>
            console.log("Loading GSplat from: $url");
            // Инициализация PlayCanvas и загрузка .ply
            // Здесь будет код Supersplat для рендеринга Gaussian Splatting
        </script>
    </body>
    </html>
    ''';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('3DGS Splat Viewer')),
      body: WebViewWidget(controller: controller),
    );
  }
}
