package database

import (
	"log"
	"caixa/internal/models"
	"gorm.io/gorm"
)

// SeedCategoriasPadrao cria categorias padrão para cada empresa
// NOTA: Este seeder não cria categorias automaticamente
// As categorias devem ser criadas pelos admins de cada empresa
func SeedCategoriasPadrao(db *gorm.DB) error {
	log.Println("ℹ️ Categorias devem ser criadas pelos administradores de cada empresa")
	return nil
}

// SeedCategoriasParaEmpresa cria categorias padrão para uma empresa específica
func SeedCategoriasParaEmpresa(db *gorm.DB, empresaID int) error {
	log.Printf("🌱 Criando categorias padrão para empresa ID: %d", empresaID)

	categorias := []models.CategoriaMovimentacao{
		// Categorias de Entrada
		{
			EmpresaID: empresaID,
			Nome:      "Vendas",
			Descricao: "Recebimento de vendas",
			Tipo:      models.TipoEntrada,
			Cor:       "#10B981", // Verde
			Icone:     "fa-shopping-cart",
			Ativo:     true,
			Ordem:     1,
		},
		{
			EmpresaID: empresaID,
			Nome:      "Recebimentos",
			Descricao: "Recebimentos diversos",
			Tipo:      models.TipoEntrada,
			Cor:       "#3B82F6", // Azul
			Icone:     "fa-hand-holding-usd",
			Ativo:     true,
			Ordem:     2,
		},
		{
			EmpresaID: empresaID,
			Nome:      "Depósito",
			Descricao: "Depósitos em caixa",
			Tipo:      models.TipoEntrada,
			Cor:       "#8B5CF6", // Roxo
			Icone:     "fa-piggy-bank",
			Ativo:     true,
			Ordem:     3,
		},
		{
			EmpresaID: empresaID,
			Nome:      "Outros Recebimentos",
			Descricao: "Outros tipos de recebimentos",
			Tipo:      models.TipoEntrada,
			Cor:       "#06B6D4", // Ciano
			Icone:     "fa-money-bill-wave",
			Ativo:     true,
			Ordem:     4,
		},
		// Categorias de Saída
		{
			EmpresaID: empresaID,
			Nome:      "Despesas",
			Descricao: "Despesas operacionais",
			Tipo:      models.TipoSaida,
			Cor:       "#EF4444", // Vermelho
			Icone:     "fa-file-invoice-dollar",
			Ativo:     true,
			Ordem:     1,
		},
		{
			EmpresaID: empresaID,
			Nome:      "Pagamentos",
			Descricao: "Pagamentos diversos",
			Tipo:      models.TipoSaida,
			Cor:       "#F59E0B", // Laranja
			Icone:     "fa-credit-card",
			Ativo:     true,
			Ordem:     2,
		},
		{
			EmpresaID: empresaID,
			Nome:      "Retirada",
			Descricao: "Retiradas do caixa",
			Tipo:      models.TipoSaida,
			Cor:       "#EC4899", // Rosa
			Icone:     "fa-wallet",
			Ativo:     true,
			Ordem:     3,
		},
		{
			EmpresaID: empresaID,
			Nome:      "Outros Pagamentos",
			Descricao: "Outros tipos de pagamentos",
			Tipo:      models.TipoSaida,
			Cor:       "#DC2626", // Vermelho escuro
			Icone:     "fa-money-check-alt",
			Ativo:     true,
			Ordem:     4,
		},
	}

	for _, categoria := range categorias {
		// Verificar se já existe (por nome, tipo e empresa)
		var existente models.CategoriaMovimentacao
		err := db.Where("nome = ? AND tipo = ? AND empresa_id = ?", categoria.Nome, categoria.Tipo, empresaID).First(&existente).Error
		
		if err == gorm.ErrRecordNotFound {
			// Criar categoria se não existir
			if err := db.Create(&categoria).Error; err != nil {
				log.Printf("⚠️ Erro ao criar categoria %s: %v", categoria.Nome, err)
			} else {
				log.Printf("✅ Categoria criada: %s (%s)", categoria.Nome, categoria.Tipo)
			}
		} else if err != nil {
			log.Printf("⚠️ Erro ao verificar categoria %s: %v", categoria.Nome, err)
		}
		// Se já existe, não faz nada
	}

	log.Printf("✅ Categorias padrão verificadas/criadas para empresa ID: %d", empresaID)
	return nil
}

