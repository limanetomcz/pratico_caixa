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

// createDatabaseIfNotExists cria o banco de dados se ele não existir
func createDatabaseIfNotExists() error {
	dbName := os.Getenv("DB_DATABASE")
	if dbName == "" {
		return fmt.Errorf("DB_DATABASE não configurado")
	}

	// Conectar ao MySQL sem especificar o banco de dados
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	tempDB, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("erro ao conectar ao MySQL: %w", err)
	}
	defer func() {
		sqlDB, _ := tempDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	// Verificar se o banco existe contando os resultados
	var count int64
	err = tempDB.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", dbName).Scan(&count).Error
	if err != nil {
		// Se der erro na verificação, tentar criar mesmo assim (mais seguro)
		log.Printf("⚠️  Não foi possível verificar se o banco existe, tentando criar...")
	}

	if count == 0 {
		log.Printf("📦 Criando banco de dados '%s'...", dbName)
		// Criar o banco de dados (IF NOT EXISTS garante que não dará erro se já existir)
		err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)).Error
		if err != nil {
			return fmt.Errorf("erro ao criar banco de dados: %w", err)
		}
		log.Printf("✅ Banco de dados '%s' criado com sucesso", dbName)
	} else {
		log.Printf("✅ Banco de dados '%s' já existe", dbName)
	}

	return nil
}

func Connect() (*gorm.DB, error) {
	// Criar banco de dados se não existir
	if err := createDatabaseIfNotExists(); err != nil {
		log.Printf("⚠️  Aviso ao criar banco de dados: %v", err)
		log.Println("⚠️  Tentando conectar mesmo assim...")
	}

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

