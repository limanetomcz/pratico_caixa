# 💰 Módulo Caixa - Golang

Módulo de gerenciamento de caixa desenvolvido em **Golang** com **Gin Framework** e interface moderna com **Tailwind CSS**.

## 📋 Funcionalidades

- ✅ **Abertura e Fechamento de Caixa** - Controle completo do ciclo de vida do caixa
- ✅ **Entradas e Saídas** - Registro de todas as movimentações financeiras
- ✅ **Relatório de Fechamento** - Relatório completo com todas as movimentações
- ✅ **Histórico de Caixas** - Consulta de caixas antigos
- ✅ **Isolamento por Usuário** - Cada usuário gerencia apenas seus próprios caixas
- ✅ **API para Integração** - Endpoints para integração com outros módulos (ex: contas a receber)

## 🏗️ Estrutura

```
caixa/
├── internal/
│   ├── config/        # Configurações
│   ├── database/      # Conexão com banco e migrations
│   ├── handlers/      # Handlers HTTP (web e API)
│   ├── middleware/    # Middlewares (Auth, CORS, Logger)
│   ├── models/        # Models (Caixa, Movimentacao)
│   ├── routes/        # Definição de rotas
│   └── services/      # Serviços (Auth)
├── web/
│   ├── static/        # Arquivos estáticos
│   └── templates/     # Templates HTML
├── main.go            # Entry point
├── go.mod             # Dependências Go
├── Dockerfile         # Imagem Docker
└── docker-compose.yml # Orquestração
```

## 🚀 Como Usar

### 1. Pré-requisitos

- Docker e Docker Compose
- MySQL compartilhado (`pratico_mysql`) rodando

### 2. Configuração

Copie o arquivo `env.example` para `.env`:

```bash
cp env.example .env
```

Edite o `.env` conforme necessário.

### 3. Criar Banco de Dados

Execute no MySQL compartilhado:

```sql
CREATE DATABASE IF NOT EXISTS caixa_db;
```

Ou use o script de setup:

```bash
./setup.sh
```

### 4. Executar com Docker

```bash
docker compose up -d
```

### 5. Executar Localmente (Desenvolvimento)

```bash
# Instalar dependências
go mod download

# Executar
go run main.go
```

## 🔗 Integração com Sistema

### Autenticação

O módulo se integra com o módulo **Auth** existente:

- **Login**: Redireciona para o módulo Auth ou faz login via API
- **JWT**: Valida tokens JWT emitidos pelo módulo Auth
- **Middleware**: Protege rotas com autenticação JWT

### Banco de Dados

Usa o MySQL compartilhado (`pratico_mysql`):

- **Host**: `pratico_mysql` (nome do container)
- **Porta**: `3306`
- **Database**: `caixa_db`

### Porta

- **Porta**: `9286`
- **URL Local**: `http://localhost:9286`
- **URL Produção**: Configurada via `APP_URL` no `.env`

## 📝 Rotas

### Públicas

- `GET /` - Página inicial (redireciona para login)
- `GET /login` - Página de login
- `GET /health` - Health check

### Protegidas (Requerem JWT)

#### Web
- `GET /dashboard` - Dashboard principal (gerenciamento de caixa atual)
- `GET /caixas` - Listagem de todos os caixas
- `GET /caixas/:id` - Detalhes de um caixa específico
- `GET /caixas/:id/relatorio` - Relatório de fechamento

#### API - Caixas
- `GET /api/v1/caixa/atual` - Obter caixa aberto atual do usuário
- `POST /api/v1/caixa/abrir` - Abrir um novo caixa
- `POST /api/v1/caixa/:id/fechar` - Fechar um caixa
- `GET /api/v1/caixas` - Listar todos os caixas do usuário
- `GET /api/v1/caixas/:id` - Obter detalhes de um caixa
- `GET /api/v1/caixas/abertos` - Listar caixas abertos (para integração)

#### API - Movimentações
- `POST /api/v1/movimentacoes` - Criar nova movimentação
- `GET /api/v1/caixas/:caixa_id/movimentacoes` - Listar movimentações de um caixa
- `DELETE /api/v1/movimentacoes/:id` - Deletar movimentação

## 🔌 API para Integração (Futuro)

### Endpoint: `GET /api/v1/caixas/abertos`

Retorna os caixas abertos do usuário para integração com outros módulos (ex: contas a receber).

**Resposta:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "data_abertura": "2025-12-12T10:00:00Z",
      "valor_abertura": 100.00,
      "saldo_atual": 150.00,
      "empresa_id": 1
    }
  ]
}
```

## 🎨 Layout

O módulo usa **Tailwind CSS** via CDN para um design moderno e responsivo:

- Design limpo e profissional
- Componentes reutilizáveis
- Responsivo (mobile-first)
- Ícones Font Awesome
- Relatório imprimível

## 🔧 Desenvolvimento

### Adicionar Nova Funcionalidade

1. Criar model em `internal/models/` (se necessário)
2. Criar handler em `internal/handlers/`
3. Adicionar rota em `internal/routes/routes.go`
4. Criar template em `web/templates/` (se for página web)

### Modificar Layout

Edite os templates em `web/templates/` usando classes Tailwind CSS.

## 📦 Dependências Principais

- `github.com/gin-gonic/gin` - Framework web
- `gorm.io/gorm` - ORM
- `gorm.io/driver/mysql` - Driver MySQL
- `github.com/golang-jwt/jwt/v5` - JWT
- `github.com/joho/godotenv` - Variáveis de ambiente

## 🐛 Troubleshooting

### Erro de conexão com banco

Verifique se o container `pratico_mysql` está rodando:

```bash
docker ps | grep pratico_mysql
```

### Erro de autenticação

Verifique se o módulo Auth está rodando e acessível:

```bash
curl http://auth_nginx:80/api/v1/health
```

### Porta já em uso

Altere a porta no `.env`:

```env
APP_PORT=9287
```

## 📚 Próximos Passos

- [x] Estrutura básica
- [x] Abertura e fechamento de caixa
- [x] Entradas e saídas
- [x] Relatório de fechamento
- [x] Histórico de caixas
- [x] Isolamento por usuário
- [ ] Integração com módulo de contas a receber
- [ ] API para receber valores de outros módulos
- [ ] API para pagar valores de outros módulos

## 🤝 Contribuindo

Este módulo está pronto para uso! Adapte conforme suas necessidades.

