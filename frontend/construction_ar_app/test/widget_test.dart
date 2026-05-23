import 'package:flutter_test/flutter_test.dart';
import 'package:construction_ar_app/main.dart';

void main() {
  testWidgets('Dashboard loads test', (WidgetTester tester) async {
    await tester.pumpWidget(const ConstructionApp());
    expect(find.text('Мои Объекты'), findsOneWidget);
  });
}
