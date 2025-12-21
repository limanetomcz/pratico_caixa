#!/bin/bash

echo "🧪 Testando acesso ao dashboard..."
echo ""

echo "1️⃣ Testando sem autenticação (deve redirecionar para /login):"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -v -L http://localhost:9286/dashboard 2>&1 | grep -E "(HTTP|Location|Content-Type|<!DOCTYPE)" | head -10
echo ""

echo "2️⃣ Testando com header Accept: text/html:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -v -H "Accept: text/html" http://localhost:9286/dashboard 2>&1 | grep -E "(HTTP|Location|Content-Type|<!DOCTYPE)" | head -10
echo ""

echo "3️⃣ Verificando logs do container:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker logs caixa_app --tail 20 2>&1 | grep -E "(🔍|🔄|📊|Auth|dashboard)" || echo "Nenhum log relevante encontrado"
echo ""

echo "4️⃣ Verificando se os templates existem no container:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker exec caixa_app ls -la /app/web/templates/ 2>&1 || echo "Erro ao listar templates"
echo ""

