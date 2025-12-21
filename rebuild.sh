#!/bin/bash

echo "🔨 Reconstruindo o container do módulo Caixa..."

cd "$(dirname "$0")"

echo "📦 Parando container..."
docker compose down

echo "🔨 Reconstruindo imagem..."
docker compose build --no-cache

echo "🚀 Iniciando container..."
docker compose up -d

echo "⏳ Aguardando container iniciar..."
sleep 3

echo "📋 Verificando logs..."
docker logs caixa_app --tail 20

echo ""
echo "✅ Rebuild concluído!"
echo "🌐 Acesse: http://localhost:9286/login"

