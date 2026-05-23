#!/bin/bash

# Скрипт полной настройки окружения и запуска (Debian 13)
# Оптимизировано под Intel N100

set -e

# Цвета для вывода
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Шаг 1: Проверка и установка системных зависимостей ===${NC}"
sudo apt update
sudo apt install -y git curl build-essential sqlite3 lsof clang cmake ninja-build pkg-config libgtk-3-dev liblzma-dev

# 1. Проверка Go
if ! command -v go &> /dev/null
then
    echo "Установка Go..."
    sudo apt install -y golang-go
else
    echo -e "${GREEN}Go уже установлен: $(go version)${NC}"
fi

# 2. Проверка Flutter
if ! command -v flutter &> /dev/null
then
    echo "Flutter не найден. Пожалуйста, установите Flutter вручную: https://docs.flutter.dev/get-started/install/linux"
else
    echo -e "${GREEN}Flutter уже установлен${NC}"
fi

echo -e "${BLUE}=== Шаг 1.5: Подготовка статики Web-панели ===${NC}"
mkdir -p backend/static/web
# В данной версии статика уже создана в backend/static/web/index.html

echo -e "${BLUE}=== Шаг 2: Настройка Backend ===${NC}"
cd backend
echo "Проверка зависимостей Go..."
go mod tidy
echo "Сборка сервера..."
go build -o construction_server ./cmd/server/main.go
cd ..

echo -e "${BLUE}=== Шаг 3: Настройка Frontend App ===${NC}"
cd frontend/construction_ar_app
echo "Загрузка пакетов Flutter..."
flutter pub get
cd ../..

echo -e "${BLUE}=== Шаг 4: Запуск системы ===${NC}"
chmod +x deploy_server.sh
./deploy_server.sh

echo -e "${GREEN}Настройка завершена! Web-панель доступна по адресу сервера.${NC}"
