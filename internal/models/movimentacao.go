package models

import (
	"time"
)

// TipoMovimentacao representa o tipo de movimentação
type TipoMovimentacao string

const (
	TipoEntrada TipoMovimentacao = "entrada"
	TipoSaida   TipoMovimentacao = "saida"
)

// Movimentacao representa uma movimentação de caixa
type Movimentacao struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	CaixaID     uint            `json:"caixa_id" gorm:"not null;index"`
	UserID      int             `json:"user_id" gorm:"not null;index"`
	CategoriaID *uint           `json:"categoria_id" gorm:"index"` // Opcional: pode ser null
	Tipo        TipoMovimentacao `json:"tipo" gorm:"type:enum('entrada','saida');not null"`
	Valor       float64         `json:"valor" gorm:"type:decimal(15,2);not null"`
	Descricao   string          `json:"descricao" gorm:"type:varchar(255);not null"`
	Observacoes string          `json:"observacoes" gorm:"type:text"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	
	// Relacionamentos
	Caixa     Caixa                `json:"caixa,omitempty" gorm:"foreignKey:CaixaID"`
	Categoria *CategoriaMovimentacao `json:"categoria,omitempty" gorm:"foreignKey:CategoriaID"`
}

// TableName define o nome da tabela
func (Movimentacao) TableName() string {
	return "movimentacoes"
}

