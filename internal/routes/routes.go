package routes

import (
	"github.com/gin-gonic/gin"
	"caixa/internal/handlers"
	"caixa/internal/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Rotas públicas
	public := router.Group("/")
	{
		public.GET("/", handlers.Index)
		public.GET("/health", handlers.Health)
		public.GET("/login", handlers.LoginPage)
	}

	// Rotas protegidas (requerem autenticação)
	protected := router.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		// Páginas web
		protected.GET("/dashboard", handlers.Dashboard)
		protected.GET("/caixas", handlers.CaixasPage)
		protected.GET("/categorias", handlers.CategoriasPage)
		// Rotas específicas devem vir ANTES da rota genérica /caixas/:id
		protected.GET("/caixas/:id/relatorio", handlers.RelatorioFechamentoPage)
		protected.GET("/caixas/:id/extrato", handlers.RelatorioExtratoPage)
		protected.GET("/caixas/:id", handlers.CaixaDetalhesPage)

		// API - Caixas
		protected.GET("/api/v1/caixa/atual", handlers.GetCaixaAtual(db))
		protected.POST("/api/v1/caixa/abrir", handlers.AbrirCaixa(db))
		protected.GET("/api/v1/caixas", handlers.ListarCaixas(db))
		
		// API - Caixas abertos (para integração com outros módulos) - DEVE VIR ANTES DE /caixas/:id
		protected.GET("/api/v1/caixas/abertos", handlers.GetCaixasAbertos(db))
		
		// API - Categorias
		protected.GET("/api/v1/categorias", handlers.ListarCategorias(db))
		protected.POST("/api/v1/categorias", handlers.CriarCategoria(db))
		protected.PUT("/api/v1/categorias/:id", handlers.AtualizarCategoria(db))
		protected.DELETE("/api/v1/categorias/:id", handlers.DeletarCategoria(db))
		
		// API - Movimentações (DEVE VIR ANTES DE /caixas/:id para evitar conflito)
		protected.POST("/api/v1/movimentacoes", handlers.CriarMovimentacao(db))
		protected.GET("/api/v1/caixas/:id/movimentacoes", handlers.ListarMovimentacoes(db))
		protected.DELETE("/api/v1/movimentacoes/:id", handlers.DeletarMovimentacao(db))
		
		// API - Caixa específico (DEVE VIR DEPOIS das rotas mais específicas)
		protected.GET("/api/v1/caixas/:id", handlers.GetCaixa(db))
		protected.POST("/api/v1/caixas/:id/fechar", handlers.FecharCaixa(db))
	}
}

