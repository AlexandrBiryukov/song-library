// database/database.go
package database

import (
	"fmt"
	"song-library/config"
	"song-library/logger"
	"song-library/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Миграции
	err = db.AutoMigrate(&models.Song{})
	if err != nil {
		return err
	}

	DB = db
	logger.Log.Info("Database connected and migrated")
	return nil
}
