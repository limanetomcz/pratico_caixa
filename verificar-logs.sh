#!/bin/bash

echo "🔍 Verificando status do container caixa_app..."
docker ps | grep caixa_app || echo "⚠️ Container não está rodando"

echo ""
echo "📋 Últimos logs do container (últimas 100 linhas):"
docker logs caixa_app --tail 100 2>&1

echo ""
echo "🧪 Testando health check:"
curl -v http://localhost:9286/health 2>&1 | head -20

echo ""
echo "🧪 Testando acesso ao dashboard (sem token - deve redirecionar):"
curl -v http://localhost:9286/dashboard 2>&1 | head -20

echo ""
echo "🧪 Testando acesso à página de login:"
curl -v http://localhost:9286/login 2>&1 | head -20

echo ""
echo "📊 Verificando se a porta está aberta:"
netstat -tuln | grep 9286 || ss -tuln | grep 9286 || echo "Comando netstat/ss não disponível"

