package database

import (
	"fmt"
	"log"

	"github.com/your-org/oj-platform/internal/models"
	"github.com/your-org/oj-platform/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.DatabaseConfig) error {
	var err error

	gormConfig := &gorm.Config{}
	if config.AppConfig.Server.Mode == "debug" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	// 优先使用SQLite进行开发测试
	if cfg.Host == "" || cfg.Host == "sqlite" {
		DB, err = gorm.Open(sqlite.Open("oj_platform.db"), gormConfig)
	} else {
		DB, err = gorm.Open(postgres.Open(cfg.DSN()), gormConfig)
	}

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移
	if err = autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Problem{},
		&models.TestCase{},
		&models.Submission{},
	)
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
