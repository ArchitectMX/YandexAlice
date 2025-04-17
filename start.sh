#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -eu

# Проверка, что Go установлен
if ! command -v go &> /dev/null; then
    echo "Go не установлен. Пожалуйста, установи Go перед запуском."
    exit 1
fi

# Опционально: установка зависимостей
go mod tidy

# Запуск основного файла
echo "🚀 Запуск сервера..."
go run alice.go
