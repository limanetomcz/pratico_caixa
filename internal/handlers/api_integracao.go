package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"caixa/internal/models"
	"gorm.io/gorm"
)

// GetCaixasAbertos retorna os caixas abertos do usuário (para integração com outros módulos)
func GetCaixasAbertos(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		var caixas []models.Caixa
		if err := db.Where("user_id = ? AND status = ?", userID, models.StatusAberto).
			Order("created_at DESC").
			Find(&caixas).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao buscar caixas abertos",
				"error":   err.Error(),
			})
			return
		}

		// Formatar resposta para integração
		result := make([]gin.H, len(caixas))
		for i, caixa := range caixas {
			saldo := caixa.GetSaldoAtual(db)
			result[i] = gin.H{
				"id":            caixa.ID,
				"data_abertura":  caixa.DataAbertura,
				"valor_abertura": caixa.ValorAbertura,
				"saldo_atual":   saldo,
				"empresa_id":    caixa.EmpresaID,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	}
}

