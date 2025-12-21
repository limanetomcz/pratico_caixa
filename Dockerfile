# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependências do sistema
RUN apk add --no-cache git

# Copiar go mod files primeiro
COPY go.mod ./

# Baixar dependências (sem go.sum ainda)
RUN go mod download

# Copiar código fonte
COPY . .

# Gerar go.sum e fazer build (binário estático para Alpine)
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o main . && ls -lh main

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copiar binário do builder
COPY --from=builder /app/main /app/main
COPY --from=builder /app/web /app/web

# Verificar se o binário existe e torná-lo executável
RUN ls -la /app/ && chmod +x /app/main && ls -la /app/main

# Expor porta
EXPOSE 9286

# Comando para executar (usar caminho absoluto)
CMD ["/app/main"]

