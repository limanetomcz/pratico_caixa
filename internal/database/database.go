package database

import (
	"fmt"
	"log"
	"os"

	"caixa/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	log.Println("✅ Conectado ao banco de dados MySQL")

	return DB, nil
}

// Migrate executa as migrations
func Migrate(db *gorm.DB) error {
	log.Println("🔄 Executando migrations...")
	
	err := db.AutoMigrate(
		&models.Caixa{},
		&models.Movimentacao{},
		&models.CategoriaMovimentacao{},
	)
	
	if err != nil {
		return fmt.Errorf("erro ao executar migrations: %w", err)
	}
	
	log.Println("✅ Migrations executadas com sucesso")
	return nil
}

