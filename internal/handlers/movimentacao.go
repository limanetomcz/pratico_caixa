package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"caixa/internal/models"
	"gorm.io/gorm"
)

// CriarMovimentacao cria uma nova movimentação
func CriarMovimentacao(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		var req struct {
			CaixaID     uint    `json:"caixa_id" binding:"required"`
			CategoriaID *uint   `json:"categoria_id"` // Opcional
			Tipo        string  `json:"tipo" binding:"required,oneof=entrada saida"`
			Valor       float64 `json:"valor" binding:"required,gt=0"`
			Descricao   string  `json:"descricao" binding:"required"`
			Observacoes string  `json:"observacoes"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Dados inválidos",
				"error":   err.Error(),
			})
			return
		}

		// Verificar se o caixa existe e pertence ao usuário
		var caixa models.Caixa
		if err := db.Where("id = ? AND user_id = ?", req.CaixaID, userID).First(&caixa).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "Caixa não encontrado",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao buscar caixa",
				"error":   err.Error(),
			})
			return
		}

		// Verificar se o caixa está aberto
		if !caixa.IsAberto() {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Não é possível adicionar movimentação em caixa fechado",
			})
			return
		}

		// Validar categoria se fornecida
		if req.CategoriaID != nil && *req.CategoriaID > 0 {
			var categoria models.CategoriaMovimentacao
			if err := db.First(&categoria, *req.CategoriaID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Categoria não encontrada",
				})
				return
			}
			// Verificar se o tipo da categoria corresponde ao tipo da movimentação
			if categoria.Tipo != models.TipoMovimentacao(req.Tipo) {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "A categoria selecionada não corresponde ao tipo de movimentação",
				})
				return
			}
		}

		movimentacao := models.Movimentacao{
			CaixaID:     req.CaixaID,
			CategoriaID: req.CategoriaID,
			UserID:      userID.(int),
			Tipo:        models.TipoMovimentacao(req.Tipo),
			Valor:       req.Valor,
			Descricao:   req.Descricao,
			Observacoes: req.Observacoes,
		}

		if err := db.Create(&movimentacao).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao criar movimentação",
				"error":   err.Error(),
			})
			return
		}

		// Carregar categoria se houver
		if movimentacao.CategoriaID != nil {
			db.Preload("Categoria").First(&movimentacao, movimentacao.ID)
		}

		// Recalcular saldo
		saldoAtual := caixa.GetSaldoAtual(db)

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Movimentação criada com sucesso",
			"data": gin.H{
				"movimentacao": movimentacao,
				"saldo_atual":  saldoAtual,
			},
		})
	}
}

// ListarMovimentacoes lista movimentações de um caixa
func ListarMovimentacoes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		caixaID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

		// Verificar se o caixa pertence ao usuário
		var caixa models.Caixa
		if err := db.Where("id = ? AND user_id = ?", caixaID, userID).First(&caixa).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "Caixa não encontrado",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao buscar caixa",
				"error":   err.Error(),
			})
			return
		}

		var movimentacoes []models.Movimentacao
		if err := db.Preload("Categoria").Where("caixa_id = ?", caixaID).Order("created_at DESC").Find(&movimentacoes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao listar movimentações",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    movimentacoes,
		})
	}
}

// DeletarMovimentacao deleta uma movimentação
func DeletarMovimentacao(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		movID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

		// Verificar se a movimentação existe e pertence ao usuário
		var movimentacao models.Movimentacao
		if err := db.Where("id = ? AND user_id = ?", movID, userID).First(&movimentacao).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "Movimentação não encontrada",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao buscar movimentação",
				"error":   err.Error(),
			})
			return
		}

		// Verificar se o caixa está aberto
		var caixa models.Caixa
		if err := db.First(&caixa, movimentacao.CaixaID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao buscar caixa",
				"error":   err.Error(),
			})
			return
		}

		if !caixa.IsAberto() {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Não é possível deletar movimentação de caixa fechado",
			})
			return
		}

		if err := db.Delete(&movimentacao).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao deletar movimentação",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Movimentação deletada com sucesso",
		})
	}
}

