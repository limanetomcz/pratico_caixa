package models

import (
	"time"
)

// CategoriaMovimentacao representa uma categoria de movimentação
type CategoriaMovimentacao struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	EmpresaID   int       `json:"empresa_id" gorm:"not null;index"` // Obrigatório: categoria por empresa
	Nome        string    `json:"nome" gorm:"type:varchar(100);not null"`
	Descricao   string    `json:"descricao" gorm:"type:text"`
	Tipo        TipoMovimentacao `json:"tipo" gorm:"type:enum('entrada','saida');not null"` // Tipo padrão da categoria
	Cor         string    `json:"cor" gorm:"type:varchar(7);default:'#3B82F6'"` // Cor em hex para UI
	Icone       string    `json:"icone" gorm:"type:varchar(50)"` // Nome do ícone FontAwesome
	Ativo       bool      `json:"ativo" gorm:"default:true"`
	Ordem       int       `json:"ordem" gorm:"default:0"` // Ordem de exibição
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relacionamentos
	Movimentacoes []Movimentacao `json:"movimentacoes,omitempty" gorm:"foreignKey:CategoriaID"`
}

// TableName define o nome da tabela
func (CategoriaMovimentacao) TableName() string {
	return "categorias_movimentacao"
}

