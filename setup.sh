#!/bin/bash

echo "💰 Configurando Módulo Caixa (Golang)..."

# Verificar se .env existe
if [ ! -f .env ]; then
    echo "📝 Criando arquivo .env..."
    cp env.example .env
    echo "✅ Arquivo .env criado. Edite conforme necessário."
fi

# Criar banco de dados no MySQL compartilhado
echo "🗄️  Criando banco de dados..."
docker exec pratico_mysql mysql -uroot -proot -e "CREATE DATABASE IF NOT EXISTS caixa_db;" 2>/dev/null || {
    echo "⚠️  Não foi possível criar o banco automaticamente."
    echo "   Execute manualmente: CREATE DATABASE IF NOT EXISTS caixa_db;"
}

# Build e iniciar containers
echo "🔨 Construindo e iniciando containers..."
docker compose build
docker compose up -d

echo "✅ Setup concluído!"
echo ""
echo "📋 Próximos passos:"
echo "   1. Acesse: http://localhost:9286"
echo "   2. Faça login com suas credenciais do módulo Auth"
echo "   3. Abra um caixa e comece a trabalhar!"
echo ""
echo "📚 Documentação: veja README.md"

