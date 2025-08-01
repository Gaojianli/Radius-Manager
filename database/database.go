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

	// 重试连接数据库，最多重试 30 次，每次间隔 2 秒
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			// 测试连接
			sqlDB, testErr := DB.DB()
			if testErr == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					log.Printf("Database connected successfully on attempt %d/%d", i+1, maxRetries)
					break
				}
			}
		}

		if i == maxRetries-1 {
			return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
		}

		log.Printf("Database connection attempt %d/%d failed: %v, retrying in 2 seconds...", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if DB == nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = DB.AutoMigrate(&models.User{}, &models.AuthLog{})
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
			Username: config.AppConfig.DefaultAdminUser,
			Email:    config.AppConfig.DefaultAdminEmail,
			Password: config.AppConfig.DefaultAdminPass,
			IsAdmin:  true,
			Banned:   false,
		}

		if err := DAO.User.Create(ctx, &admin); err != nil {
			return err
		}

		log.Printf("Default admin created: username=%s, email=%s",
			config.AppConfig.DefaultAdminUser,
			config.AppConfig.DefaultAdminEmail)
	}

	return nil
}
