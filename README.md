# Интерактивная система AR-мониторинга стройки (SaaS)

Комплексное решение для визуализации скрытых коммуникаций и слоев стен.

## Технологии
- **AR & 3D**: Google Geospatial API, 3D Gaussian Splatting (PlayCanvas), ARCore.
- **Sensing**: RuView (WiFi Sensing AI) через ESP32-S3 (OTG/Bluetooth).
- **Backend**: Go + SQLite (оптимизировано под Intel N100).
- **Frontend**: Flutter (вычисления на GPU мобильного устройства).
- **Intelligence**: Google Cloud Vision & Gemini Flash для анализа чертежей.

## Роли доступа
- **Владелец**: Полный аудит проекта.
- **Управляющий**: Настройка объекта и прав.
- **Архитектор**: Работа с BIM и планами.
- **Прораб**: Съемка 360-панорам и Splats.
- **Инженер/Рабочий**: AR X-Ray ("прозрачные стены").

## Установка и запуск
1. **Быстрый деплой на сервер (Debian 13)**:
   ```bash
   chmod +x deploy_server.sh
   ./deploy_server.sh
   ```
2. **Backend (Вручную)**:
   \`\`\`bash
   cd backend
   go run cmd/server/main.go
   \`\`\`
2. **Frontend**:
   \`\`\`bash
   cd frontend/construction_ar_app
   flutter run
   \`\`\`

## Подключение ESP32-S3 (RuView)
Подключите устройство через USB OTG. В приложении выберите порт в меню "Сенсоры". Данные WiFi Sensing будут автоматически накладываться в AR.
