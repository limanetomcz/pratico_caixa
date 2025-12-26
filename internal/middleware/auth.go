package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"caixa/internal/services"
)

// min retorna o menor valor entre dois inteiros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// JWTAuth valida o token JWT nas requisições
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar se é requisição web (não API)
		// Verifica Accept header, Content-Type, ou se o path não começa com /api
		accept := c.Request.Header.Get("Accept")
		path := c.Request.URL.Path
		contentType := c.Request.Header.Get("Content-Type")
		
		// É requisição web se:
		// 1. Accept contém text/html OU
		// 2. Accept está vazio (navegador padrão) OU
		// 3. Path não começa com /api (rotas web não são APIs) OU
		// 4. Content-Type é text/html
		// MAS não é se for uma requisição XHR/AJAX (Accept contém application/json)
		isWebRequest := !strings.Contains(accept, "application/json") && (
			strings.Contains(accept, "text/html") || 
			accept == "" ||
			!strings.HasPrefix(path, "/api") ||
			contentType == "text/html")
		
		// Log para debug (remover em produção)
		if os.Getenv("APP_ENV") != "production" {
			log.Printf("🔍 Auth Middleware - Path: %s, Accept: %s, IsWeb: %v", path, accept, isWebRequest)
		}

		// Tentar obter token do header Authorization
		authHeader := c.GetHeader("Authorization")
		var token string
		
		if authHeader == "" {
			// Tentar obter do cookie
			var err error
			token, err = c.Cookie("auth_token")
			if err != nil || token == "" {
				// Tentar obter da URL (vindo do módulo de empresas)
				tokenFromURL := c.Query("token")
				log.Printf("🔎 Verificando token na URL: encontrado=%v, tamanho=%d", tokenFromURL != "", len(tokenFromURL))
				if tokenFromURL != "" {
					// Decodificar token da URL (pode estar codificado)
					token = tokenFromURL
					// Salvar token em cookie para próximas requisições
					c.SetCookie("auth_token", token, 86400, "/", "", false, false)
					log.Printf("✅ Token recebido via URL e salvo em cookie (tamanho: %d)", len(token))
				} else {
					log.Printf("⚠️ Token não encontrado - nem no cookie nem na URL")
				}
			} else {
				if os.Getenv("APP_ENV") != "production" {
					log.Printf("✅ Token encontrado no cookie (tamanho: %d)", len(token))
				}
			}
			
			if token == "" {
				// Se for requisição web, redirecionar para login
				if isWebRequest {
					log.Printf("🔄 Redirecionando para /login (sem token)")
					c.Redirect(http.StatusFound, "/login")
					c.Abort()
					return
				}
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "Token não fornecido",
				})
				c.Abort()
				return
			}
			authHeader = "Bearer " + token
		} else {
			// Extrair token do header Authorization
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// Verificar se temos um token válido
		if token == "" {
			// Tentar extrair do header Authorization se ainda não tiver
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
		
		if token == "" {
			if isWebRequest {
				log.Printf("🔄 Redirecionando para /login (formato de token inválido)")
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Formato de token inválido",
			})
			c.Abort()
			return
		}

		// Log para debug
		log.Printf("🔑 Validando token (tamanho: %d, primeiros 20 chars): %s...", len(token), token[:min(20, len(token))])

		// Validar token com o serviço de autenticação
		user, err := services.ValidateJWT(token)
		if err != nil {
			log.Printf("❌ Erro ao validar token: %v", err)
			if isWebRequest {
				// Verificar se o erro é de token expirado
				errMsg := strings.ToLower(err.Error())
				if strings.Contains(errMsg, "expired") || strings.Contains(errMsg, "expirado") {
					c.Redirect(http.StatusFound, "/login?expired=1")
				} else {
					c.Redirect(http.StatusFound, "/login")
				}
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token inválido ou expirado",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// Log para debug
		log.Printf("✅ Token válido - User ID: %d, Email: %s, Role: %s", user.ID, user.Email, user.Role)

		// Adicionar informações do usuário ao contexto
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user_name", user.Name)
		c.Set("user_role", user.Role)
		c.Set("empresa_id", user.EmpresaID)

		c.Next()
	}
}

