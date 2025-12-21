package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"caixa/internal/database"
	"caixa/internal/middleware"
	"caixa/internal/routes"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Configurar banco de dados
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	// Executar migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Erro ao executar migrations:", err)
	}

	// Categorias devem ser criadas pelos administradores de cada empresa
	// Não criamos categorias automaticamente

	// Configurar Gin
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware global
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Servir arquivos estáticos
	router.Static("/static", "./web/static")
	
	// Carregar templates HTML
	templatePath := "web/templates/*"
	if _, err := os.Stat("web/templates"); os.IsNotExist(err) {
		log.Printf("⚠️ Diretório de templates não encontrado em: web/templates")
		log.Printf("📂 Diretório atual: %s", os.Getenv("PWD"))
		if wd, err := os.Getwd(); err == nil {
			log.Printf("📂 Working directory: %s", wd)
		}
	} else {
		log.Printf("✅ Templates encontrados em: %s", templatePath)
	}
	router.LoadHTMLGlob(templatePath)

	// Rotas
	routes.SetupRoutes(router, db)

	// Porta do servidor
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "9286"
	}

	log.Printf("🚀 Módulo Caixa iniciado na porta %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}

