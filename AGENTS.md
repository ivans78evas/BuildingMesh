# ИНСТРУКЦИИ ДЛЯ АГЕНТА (BUILDINGMESH)

Этот файл содержит фундаментальные правила разработки и архитектурные контракты проекта. Агент обязан следовать им беспрекословно.

## 1. АРХИТЕКТУРНЫЙ СТЭК
- **Backend:** Go 1.23, `modernc.org/sqlite` (CGO-free), Redis, RabbitMQ.
- **Frontend:** Flutter (Mobile AR), React (Web Engineering Console).
- **Core Engine:** Изолированное ядро обработки LiDAR/CSI на Go.
- **Enterprise Head:** Twenty CRM (Metadata & Workflow) через Proxy-паттерн.

## 2. IP PROTECTION & LICENSING (ПРАВИЛА ТОРГОВОЙ ЧИСТОТЫ)
- **Лицензионная изоляция:** Запрещено импортировать или наследовать код Twenty CRM (AGPL-3.0) в проприетарные сервисы Go.
- **Proxy Only:** Все взаимодействие с Twenty происходит строго через API (GraphQL/REST).
- **Multi-tenancy:** Каждая SQL-структура и API-запрос ДОЛЖНЫ содержать `organization_id`. Смешивание данных тенантов является критической ошибкой.

## 3. ВЫСОКАЯ НАГРУЗКА И ХРАНЕНИЕ (SCALABILITY)
- **Бинарные данные:** Категорически ЗАПРЕЩЕНО хранить тяжелые файлы (Point Clouds, CSI-матрицы, панорамы 360) в SQLite.
- **S3-First:** Все файлы > 1MB сбрасываются в `S3Service` (или локальный Volume). В БД хранятся только метаданные и ссылки.
- **RabbitMQ Damping:** Все операции загрузки и обработки 3D-данных должны быть асинхронными. API отвечает `202 Accepted` сразу после сохранения сырого файла в хранилище.

## 4. ПРАВИЛА UI/UX И ФИЛОСОФИЯ REACT (DEFINITION OF DONE)
- **Стиль:** Dark Mode, шрифт 'JetBrains Mono' для технических данных.
- **Time-Machine:** Любая работа со слоями стен должна поддерживать хронологическую навигацию (`TemporalScanLayer`).
- **Data Integrity:** После подтверждения (Health Score > 95%) объект `InspectionCommit` должен переходить в статус `Read-Only`.
- **React Philosophy:**
  - **Component-Based:** UI строится через атомарные компоненты (ThreeDViewer, StatusBadge, AuthGate).
  - **Single Source of Truth:** Состояние проекта и текущего слоя управляется централизованно.
  - **Hooks Only:** Использование функциональных компонентов и React Hooks для управления жизненным циклом и эффектами (useEffect для Three.js).
  - **Composition:** Сложные интерфейсы собираются через композицию компонентов, а не наследование.

## 5. ПРОВЕРКА ПЕРЕД КОММИТОМ
1. `cd backend && go build ./...` — проверка компиляции.
2. `grep "organization_id" backend/internal/models/models.go` — проверка multi-tenant контракта.
3. Проверка отсутствия бинарных файлов в репозитории.
