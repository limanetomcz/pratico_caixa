package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index renderiza a página inicial
func Index(c *gin.Context) {
	// Verificar se já tem token no cookie (pode vir de outro módulo se estiver no mesmo domínio)
	token, _ := c.Cookie("auth_token")
	
	// Se não tem no cookie, verificar se há token na URL (vindo do módulo de empresas)
	if token == "" {
		tokenFromURL := c.Query("token")
		if tokenFromURL != "" {
			// Salvar token em cookie e redirecionar
			// httpOnly=false para permitir que JavaScript leia o cookie
			c.SetCookie("auth_token", tokenFromURL, 86400, "/", "", false, false)
			c.Redirect(http.StatusFound, "/dashboard")
			return
		}
		// Se não tem nem no cookie nem na URL, redirecionar para login
		c.Redirect(http.StatusFound, "/login")
		return
	}
	
	// Se já tem token no cookie, ir direto para o dashboard
	c.Redirect(http.StatusFound, "/dashboard")
}

// LoginPage renderiza a página de login
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login - Caixa",
	})
}

// Dashboard renderiza o dashboard
// NOTA: O middleware já validou o token e adicionou os dados do usuário ao contexto
// Não precisamos verificar o token novamente aqui
func Dashboard(c *gin.Context) {
	// Garantir que o cookie está salvo (caso o token tenha vindo da URL)
	tokenFromURL := c.Query("token")
	if tokenFromURL != "" {
		// Salvar token em cookie para que o JavaScript possa ler
		// Parâmetros: nome, valor, maxAge (segundos), path, domain, secure, httpOnly
		c.SetCookie("auth_token", tokenFromURL, 86400, "/", "", false, false)
		log.Printf("✅ Token da URL salvo em cookie no handler Dashboard (tamanho: %d)", len(tokenFromURL))
	} else {
		// Verificar se já tem cookie
		token, _ := c.Cookie("auth_token")
		if token != "" {
			log.Printf("✅ Token já existe no cookie (tamanho: %d)", len(token))
		} else {
			log.Printf("⚠️ Token não encontrado nem na URL nem no cookie")
		}
	}
	
	log.Printf("📊 Renderizando dashboard...")
	
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")
	userEmail, _ := c.Get("user_email")
	userRole, _ := c.Get("user_role")

	// Valores padrão se não houver dados do usuário
	if userName == nil || userName == "" {
		userName = "Usuário"
	}
	if userEmail == nil {
		userEmail = ""
	}
	if userRole == nil {
		userRole = "usuario"
	}

	log.Printf("📊 Dados do usuário - ID: %v, Nome: %v, Email: %v, Role: %v", userID, userName, userEmail, userRole)

	// Obter token para passar ao JavaScript
	token, _ := c.Cookie("auth_token")
	if token == "" {
		token = c.Query("token")
	}
	
	// Codificar token como JSON para passar ao JavaScript de forma segura
	var tokenJSON string
	if token != "" {
		tokenBytes, err := json.Marshal(token)
		if err == nil {
			tokenJSON = string(tokenBytes)
		} else {
			log.Printf("⚠️ Erro ao codificar token como JSON: %v", err)
		}
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":    "Dashboard - Caixa",
		"userID":   userID,
		"userName": userName,
		"userEmail": userEmail,
		"userRole": userRole,
		"tokenJSON": tokenJSON, // Token codificado como JSON (pode ser vazio)
	})
}

// CaixasPage renderiza a página de listagem de caixas
func CaixasPage(c *gin.Context) {
	c.HTML(http.StatusOK, "caixas.html", gin.H{
		"title": "Caixas - Sistema",
	})
}

// CaixaDetalhesPage renderiza a página de detalhes do caixa
func CaixaDetalhesPage(c *gin.Context) {
	caixaID := c.Param("id")
	c.HTML(http.StatusOK, "caixa_detalhes.html", gin.H{
		"title":   "Detalhes do Caixa",
		"caixaID": caixaID,
	})
}

// RelatorioFechamentoPage renderiza o relatório de fechamento
func RelatorioFechamentoPage(c *gin.Context) {
	caixaID := c.Param("id")
	c.HTML(http.StatusOK, "relatorio_fechamento.html", gin.H{
		"title":   "Relatório de Fechamento",
		"caixaID": caixaID,
	})
}

// RelatorioExtratoPage renderiza o extrato simplificado do caixa
func RelatorioExtratoPage(c *gin.Context) {
	caixaID := c.Param("id")
	c.HTML(http.StatusOK, "relatorio_extrato.html", gin.H{
		"title":   "Extrato de Caixa",
		"caixaID": caixaID,
	})
}

// CategoriasPage renderiza a página de gerenciamento de categorias (apenas admin)
func CategoriasPage(c *gin.Context) {
	userRole, _ := c.Get("user_role")
	
	role, ok := userRole.(string)
	if !ok || (role != "admin" && role != "super_admin") {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"title":   "Acesso Negado",
			"message": "Apenas administradores podem acessar esta página",
		})
		return
	}
	
	c.HTML(http.StatusOK, "categorias.html", gin.H{
		"title": "Gerenciar Categorias",
	})
}

// Health retorna o status do serviço
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "Caixa",
		"version": "1.0.0",
	})
}
