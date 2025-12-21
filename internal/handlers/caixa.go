package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"caixa/internal/models"
	"gorm.io/gorm"
)

// AbrirCaixa abre um novo caixa
func AbrirCaixa(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		empresaID, _ := c.Get("empresa_id")

		// Verificar se já existe caixa aberto para este usuário
		var caixaExistente models.Caixa
		if err := db.Where("user_id = ? AND status = ?", userID, models.StatusAberto).First(&caixaExistente).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Já existe um caixa aberto para este usuário",
				"data":    caixaExistente,
			})
			return
		}

		var req struct {
			ValorAbertura float64 `json:"valor_abertura" binding:"required"`
			Observacoes   string  `json:"observacoes"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Dados inválidos",
				"error":   err.Error(),
			})
			return
		}

		caixa := models.Caixa{
			UserID:       userID.(int),
			EmpresaID:    getEmpresaID(empresaID),
			DataAbertura: time.Now(),
			ValorAbertura: req.ValorAbertura,
			Status:       models.StatusAberto,
			Observacoes:  req.Observacoes,
		}

		if err := db.Create(&caixa).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao abrir caixa",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Caixa aberto com sucesso",
			"data":    caixa,
		})
	}
}

// FecharCaixa fecha um caixa
func FecharCaixa(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		caixaID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

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

		if !caixa.IsAberto() {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Caixa já está fechado",
			})
			return
		}

		// Calcular saldo atual
		valorFechamento := caixa.GetSaldoAtual(db)
		agora := time.Now()

		log.Printf("🔒 Fechando caixa ID: %d, Saldo calculado: %.2f", caixa.ID, valorFechamento)

		caixa.Status = models.StatusFechado
		caixa.DataFechamento = &agora
		caixa.ValorFechamento = &valorFechamento

		if err := db.Save(&caixa).Error; err != nil {
			log.Printf("❌ Erro ao salvar caixa fechado: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao fechar caixa",
				"error":   err.Error(),
			})
			return
		}

		log.Printf("✅ Caixa fechado com sucesso. ID: %d, Valor Fechamento: %.2f", caixa.ID, valorFechamento)

		// Carregar movimentações para o relatório
		var movimentacoes []models.Movimentacao
		db.Where("caixa_id = ?", caixa.ID).Order("created_at ASC").Find(&movimentacoes)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Caixa fechado com sucesso",
			"data": gin.H{
				"caixa":         caixa,
				"movimentacoes": movimentacoes,
			},
		})
	}
}

// GetCaixaAtual retorna o caixa aberto do usuário
func GetCaixaAtual(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		var caixa models.Caixa
		if err := db.Where("user_id = ? AND status = ?", userID, models.StatusAberto).First(&caixa).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"data":    nil,
					"message": "Nenhum caixa aberto",
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

		// Calcular saldo atual
		saldoAtual := caixa.GetSaldoAtual(db)

		// Carregar movimentações com categorias
		var movimentacoes []models.Movimentacao
		db.Preload("Categoria").Where("caixa_id = ?", caixa.ID).Order("created_at DESC").Find(&movimentacoes)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"caixa":         caixa,
				"saldo_atual":   saldoAtual,
				"movimentacoes": movimentacoes,
			},
		})
	}
}

// ListarCaixas lista todos os caixas do usuário com filtro opcional por data
func ListarCaixas(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		query := db.Where("user_id = ?", userID)

		// Filtro por data de início (data de abertura)
		dataInicio := c.Query("data_inicio")
		if dataInicio != "" {
			// Usar DATE() para comparar apenas a data, ignorando a hora
			query = query.Where("DATE(data_abertura) >= ?", dataInicio)
		}

		// Filtro por data de fim (data de abertura)
		dataFim := c.Query("data_fim")
		if dataFim != "" {
			// Usar DATE() para comparar apenas a data, ignorando a hora
			query = query.Where("DATE(data_abertura) <= ?", dataFim)
		}

		var caixas []models.Caixa
		if err := query.Order("created_at DESC").Find(&caixas).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao listar caixas",
				"error":   err.Error(),
			})
			return
		}

		// Calcular saldo para cada caixa
		result := make([]gin.H, len(caixas))
		for i, caixa := range caixas {
			saldo := caixa.GetSaldoAtual(db)
			result[i] = gin.H{
				"caixa":       caixa,
				"saldo_atual": saldo,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	}
}

// GetCaixa retorna um caixa específico com movimentações
func GetCaixa(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		caixaID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

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

		// Calcular saldo atual
		saldoAtual := caixa.GetSaldoAtual(db)

		// Carregar movimentações com categorias
		var movimentacoes []models.Movimentacao
		db.Preload("Categoria").Where("caixa_id = ?", caixa.ID).Order("created_at ASC").Find(&movimentacoes)

		// Calcular totais
		var totalEntradas, totalSaidas float64
		for _, mov := range movimentacoes {
			if mov.Tipo == models.TipoEntrada {
				totalEntradas += mov.Valor
			} else {
				totalSaidas += mov.Valor
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"caixa":          caixa,
				"saldo_atual":    saldoAtual,
				"movimentacoes":  movimentacoes,
				"total_entradas": totalEntradas,
				"total_saidas":   totalSaidas,
			},
		})
	}
}

// getEmpresaID converte empresa_id do contexto para *int
func getEmpresaID(empresaID interface{}) *int {
	if empresaID == nil {
		return nil
	}
	if id, ok := empresaID.(int); ok {
		return &id
	}
	return nil
}

