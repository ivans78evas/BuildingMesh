#!/bin/bash

# Скрипт развертывания системы мониторинга стройки для Debian 13
# Оптимизировано под Intel N100

set -e

echo "--- Настройка окружения Debian 13 ---"
sudo apt update
sudo apt install -y curl git build-essential sqlite3 lsof

# Установка Go
if ! command -v go &> /dev/null
then
    echo "Установка Go..."
    # Для Debian 13 (Trixie) можно использовать системный пакет или скачать свежий
    sudo apt install -y golang-go
fi

echo "--- Сборка Backend (N100 Optimized) ---"
cd backend
go mod tidy
go build -o construction_server ./cmd/server/main.go

# Создание директорий для данных
mkdir -p uploads/plans uploads/splats

# Очистка порта если занят
sudo kill $(lsof -t -i :8080) 2>/dev/null || true

echo "--- Запуск сервера в фоновом режиме ---"
nohup ./construction_server > server.log 2>&1 &

echo "--- Проверка запуска ---"
sleep 3
if curl -s http://localhost:8080/api/v1/projects > /dev/null; then
    echo "===================================================="
    echo "Backend успешно запущен на порту 8080!"
    echo "IP сервера для подключения мобильного приложения:"
    hostname -I | awk '{print $1}'
    echo "===================================================="
else
    echo "Ошибка запуска. Проверьте backend/server.log"
    cat server.log
fi
