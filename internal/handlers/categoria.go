package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"caixa/internal/models"
	"gorm.io/gorm"
)

// ListarCategorias lista todas as categorias da empresa do usuário
func ListarCategorias(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		empresaID, _ := c.Get("empresa_id")
		tipo := c.Query("tipo") // Filtro opcional por tipo (entrada/saida)

		if empresaID == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Empresa não identificada",
			})
			return
		}

		empresaIDInt := empresaID.(*int)
		if empresaIDInt == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Empresa não identificada",
			})
			return
		}

		query := db.Where("ativo = ? AND empresa_id = ?", true, *empresaIDInt)

		// Filtro por tipo se fornecido
		if tipo == "entrada" || tipo == "saida" {
			query = query.Where("tipo = ?", tipo)
		}

		var categorias []models.CategoriaMovimentacao
		if err := query.Order("ordem ASC, nome ASC").Find(&categorias).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao listar categorias",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    categorias,
		})
	}
}

// CriarCategoria cria uma nova categoria de movimentação (apenas admin)
func CriarCategoria(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("user_role")
		empresaID, _ := c.Get("empresa_id")

		// Verificar se é admin
		role, ok := userRole.(string)
		if !ok || (role != "admin" && role != "super_admin") {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Apenas administradores podem criar categorias",
			})
			return
		}

		if empresaID == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Empresa não identificada",
			})
			return
		}

		empresaIDInt := empresaID.(*int)
		if empresaIDInt == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Empresa não identificada",
			})
			return
		}

		var req struct {
			Nome      string `json:"nome" binding:"required"`
			Descricao string `json:"descricao"`
			Tipo      string `json:"tipo" binding:"required,oneof=entrada saida"`
			Cor       string `json:"cor"`
			Icone     string `json:"icone"`
			Ordem     int    `json:"ordem"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Dados inválidos",
				"error":   err.Error(),
			})
			return
		}

		// Definir valores padrão
		if req.Cor == "" {
			req.Cor = "#3B82F6"
		}
		if req.Icone == "" {
			if req.Tipo == "entrada" {
				req.Icone = "fa-arrow-down"
			} else {
				req.Icone = "fa-arrow-up"
			}
		}

		categoria := models.CategoriaMovimentacao{
			EmpresaID: *empresaIDInt,
			Nome:      req.Nome,
			Descricao: req.Descricao,
			Tipo:      models.TipoMovimentacao(req.Tipo),
			Cor:       req.Cor,
			Icone:     req.Icone,
			Ordem:     req.Ordem,
			Ativo:     true,
		}

		if err := db.Create(&categoria).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao criar categoria",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Categoria criada com sucesso",
			"data":    categoria,
		})
	}
}

// AtualizarCategoria atualiza uma categoria (apenas admin)
func AtualizarCategoria(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("user_role")
		empresaID, _ := c.Get("empresa_id")
		categoriaID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

		// Verificar se é admin
		role, ok := userRole.(string)
		if !ok || (role != "admin" && role != "super_admin") {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Apenas administradores podem atualizar categorias",
			})
			return
		}

		if empresaID == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Empresa não identificada",
			})
			return
		}

		empresaIDInt := empresaID.(*int)
		if empresaIDInt == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Empresa não identificada",
			})
			return
		}

		// Verificar se a categoria existe e pertence à empresa
		var categoria models.CategoriaMovimentacao
		if err := db.Where("id = ? AND empresa_id = ?", categoriaID, *empresaIDInt).First(&categoria).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "Categoria não encontrada",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao buscar categoria",
				"error":   err.Error(),
			})
			return
		}

		var req struct {
			Nome      string `json:"nome"`
			Descricao string `json:"descricao"`
			Tipo      string `json:"tipo"`
			Cor       string `json:"cor"`
			Icone     string `json:"icone"`
			Ordem     *int   `json:"ordem"`
			Ativo     *bool  `json:"ativo"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Dados inválidos",
				"error":   err.Error(),
			})
			return
		}

		// Atualizar campos fornecidos
		if req.Nome != "" {
			categoria.Nome = req.Nome
		}
		// Descricao pode ser vazia, então sempre atualiza se fornecida
		categoria.Descricao = req.Descricao
		if req.Tipo == "entrada" || req.Tipo == "saida" {
			categoria.Tipo = models.TipoMovimentacao(req.Tipo)
		}
		if req.Cor != "" {
			categoria.Cor = req.Cor
		}
		// Icone pode ser vazio, então sempre atualiza se fornecido
		categoria.Icone = req.Icone
		if req.Ordem != nil {
			categoria.Ordem = *req.Ordem
		}
		if req.Ativo != nil {
			categoria.Ativo = *req.Ativo
		}

		if err := db.Save(&categoria).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao atualizar categoria",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Categoria atualizada com sucesso",
			"data":    categoria,
		})
	}
}

// DeletarCategoria deleta uma categoria (apenas admin)
func DeletarCategoria(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("user_role")
		empresaID, _ := c.Get("empresa_id")
		categoriaID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

		// Verificar se é admin
		role, ok := userRole.(string)
		if !ok || (role != "admin" && role != "super_admin") {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Apenas administradores podem deletar categorias",
			})
			return
		}

		if empresaID == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Empresa não identificada",
			})
			return
		}

		empresaIDInt := empresaID.(*int)
		if empresaIDInt == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Empresa não identificada",
			})
			return
		}

		// Verificar se a categoria existe e pertence à empresa
		var categoria models.CategoriaMovimentacao
		if err := db.Where("id = ? AND empresa_id = ?", categoriaID, *empresaIDInt).First(&categoria).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "Categoria não encontrada",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao buscar categoria",
				"error":   err.Error(),
			})
			return
		}

		// Verificar se há movimentações usando esta categoria
		var count int64
		db.Model(&models.Movimentacao{}).Where("categoria_id = ?", categoriaID).Count(&count)
		if count > 0 {
			// Em vez de deletar, desativar
			categoria.Ativo = false
			if err := db.Save(&categoria).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Erro ao desativar categoria",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Categoria desativada (há movimentações usando esta categoria)",
				"data":    categoria,
			})
			return
		}

		// Se não há movimentações, pode deletar
		if err := db.Delete(&categoria).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao deletar categoria",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Categoria deletada com sucesso",
		})
	}
}

