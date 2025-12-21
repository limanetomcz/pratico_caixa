#!/bin/bash

echo "🧪 Testando recebimento de token no módulo Caixa..."
echo ""

echo "1️⃣ Testando acesso sem token (deve redirecionar para /login):"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -v http://localhost:9286/dashboard 2>&1 | grep -E "(HTTP|Location|Set-Cookie)" | head -5
echo ""

echo "2️⃣ Testando acesso com token na URL (simulando acesso do módulo de empresas):"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "⚠️  Para testar, você precisa de um token válido. Substitua TOKEN_AQUI pelo token real."
echo "curl -v 'http://localhost:9286/dashboard?token=TOKEN_AQUI' 2>&1 | grep -E '(HTTP|Location|Set-Cookie)' | head -10"
echo ""

echo "3️⃣ Verificando logs do container:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker logs caixa_app --tail 30 2>&1 | grep -E "(🔍|🔑|✅|❌|🔄|Auth|token|Token)" || echo "Nenhum log relevante encontrado"
echo ""

echo "4️⃣ Verificando se o serviço de autenticação está acessível:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker exec caixa_app wget -qO- http://auth_nginx:80/api/v1/health 2>&1 | head -5 || echo "Erro ao acessar serviço de autenticação"
echo ""

