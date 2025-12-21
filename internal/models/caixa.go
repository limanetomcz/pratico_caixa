package models

import (
	"time"
	"gorm.io/gorm"
)

// StatusCaixa representa o status do caixa
type StatusCaixa string

const (
	StatusAberto  StatusCaixa = "aberto"
	StatusFechado StatusCaixa = "fechado"
)

// Caixa representa um caixa
type Caixa struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	UserID        int         `json:"user_id" gorm:"not null;index"`
	EmpresaID     *int        `json:"empresa_id" gorm:"index"`
	DataAbertura  time.Time   `json:"data_abertura" gorm:"not null"`
	DataFechamento *time.Time `json:"data_fechamento"`
	ValorAbertura float64    `json:"valor_abertura" gorm:"type:decimal(15,2);default:0"`
	ValorFechamento *float64 `json:"valor_fechamento" gorm:"type:decimal(15,2)"`
	Status        StatusCaixa `json:"status" gorm:"type:enum('aberto','fechado');default:'aberto'"`
	Observacoes   string      `json:"observacoes" gorm:"type:text"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	
	// Relacionamentos
	Movimentacoes []Movimentacao `json:"movimentacoes,omitempty" gorm:"foreignKey:CaixaID"`
}

// TableName define o nome da tabela
func (Caixa) TableName() string {
	return "caixas"
}

// GetSaldoAtual calcula o saldo atual do caixa
func (c *Caixa) GetSaldoAtual(db *gorm.DB) float64 {
	var totalEntradas float64
	var totalSaidas float64

	db.Model(&Movimentacao{}).
		Where("caixa_id = ? AND tipo = ?", c.ID, TipoEntrada).
		Select("COALESCE(SUM(valor), 0)").
		Scan(&totalEntradas)

	db.Model(&Movimentacao{}).
		Where("caixa_id = ? AND tipo = ?", c.ID, TipoSaida).
		Select("COALESCE(SUM(valor), 0)").
		Scan(&totalSaidas)

	return c.ValorAbertura + totalEntradas - totalSaidas
}

// IsAberto verifica se o caixa está aberto
func (c *Caixa) IsAberto() bool {
	return c.Status == StatusAberto
}

