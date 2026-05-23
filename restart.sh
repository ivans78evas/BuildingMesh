#!/bin/bash
# Скрипт быстрой перезагрузки сервера

echo "Перезапуск сервера Construction AR..."

# Находим и убиваем процесс на порту 8080
PID=$(lsof -t -i :8080)
if [ -z "$PID" ]; then
    echo "Сервер не был запущен."
else
    echo "Остановка сервера (PID: $PID)..."
    kill $PID
    sleep 1
fi

# Переходим в папку и запускаем
cd backend
if [ -f "./construction_server" ]; then
    nohup ./construction_server > server.log 2>&1 &
    echo "Сервер запущен в фоновом режиме (запись в backend/server.log)"
else
    echo "Бинарный файл не найден. Собираю..."
    go build -o construction_server ./cmd/server/main.go
    nohup ./construction_server > server.log 2>&1 &
    echo "Сервер собран и запущен."
fi

echo "Готово! API доступно на порту 8080."
