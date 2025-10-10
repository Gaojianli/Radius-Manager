package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Gaojianli/raduis_mgnt/config"
	"github.com/Gaojianli/raduis_mgnt/database"
	"github.com/Gaojianli/raduis_mgnt/models"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 连接数据库
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	ctx := context.Background()

	// 查找 admin 用户
	adminUser, err := database.DAO.User.GetByUsername(ctx, "admin")
	if err != nil {
		log.Fatal("Failed to find admin user:", err)
	}

	// 重置密码为默认密码 "admin123"
	defaultPassword := "admin123"

	// 使用 HashPassword 方法来正确哈希密码
	if err := adminUser.HashPassword(defaultPassword); err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// 更新数据库中的密码
	if err := database.DAO.User.Update(ctx, adminUser); err != nil {
		log.Fatal("Failed to update admin password:", err)
	}

	fmt.Printf("✅ Admin user '%s' password has been reset to default password: %s\n", adminUser.Username, defaultPassword)
	fmt.Println("⚠️  Please change the default password after login for security reasons.")
}
