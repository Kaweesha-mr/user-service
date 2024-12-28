package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"user-service/pkg/model"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatalf("DB_DSN is not set in Env")
		return nil, fmt.Errorf("DB_DSN is not set in Env")
	}

	var err error

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automatically migrate the schema, creating the tables if they do not exist
	err = db.AutoMigrate(&model.User{}) // You can add more models as needed
	if err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}
	log.Println("Connected to PostgreSQL successfully!")
	return db, nil
}
