package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Gaojianli/raduis_mgnt/config"
	"github.com/Gaojianli/raduis_mgnt/dao"
	"github.com/Gaojianli/raduis_mgnt/models"
)

var (
	DB  *gorm.DB
	DAO *dao.DAOManager
)

func Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DAO = dao.NewDAOManager(DB)

	err = createDefaultAdmin()
	if err != nil {
		log.Printf("Warning: failed to create default admin: %v", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func createDefaultAdmin() error {
	ctx := context.Background()

	count, err := DAO.User.CountAdmins(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		admin := models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "admin123",
			IsAdmin:  true,
			Status:   true,
		}

		if err := DAO.User.Create(ctx, &admin); err != nil {
			return err
		}

		log.Println("Default admin created: username=admin, password=admin123")
	}

	return nil
}
