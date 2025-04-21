#!/bin/bash

# Устанавливаем все зависимости
echo "📦 Установка зависимостей..."
pip install -r requirements.txt

# Запуск Flask сервера
echo "🚀 Запуск сервера..."
flask run --host=0.0.0.0 --port=8080
