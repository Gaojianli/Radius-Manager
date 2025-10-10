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

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin123"
	DefaultAdminEmail    = "admin@example.com"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 连接数据库
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	ctx := context.Background()

	switch command {
	case "reset-password":
		resetAdminPassword(ctx)
	case "create-admin":
		createAdminUser(ctx)
	case "init":
		initAdminUser(ctx)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showUsage()
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Println("Usage: admin_tool <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  reset-password  Reset admin password to default (admin123)")
	fmt.Println("  create-admin    Create default admin user if not exists")
	fmt.Println("  init            Initialize admin user (create if not exists, reset password)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/admin_tool.go reset-password")
	fmt.Println("  go run cmd/admin_tool.go create-admin")
	fmt.Println("  go run cmd/admin_tool.go init")
}

func resetAdminPassword(ctx context.Context) {
	// 查找 admin 用户
	adminUser, err := database.DAO.User.GetByUsername(ctx, DefaultAdminUsername)
	if err != nil {
		log.Fatal("Failed to find admin user. Run 'create-admin' command first:", err)
	}

	// 重置密码为默认密码
	if err := adminUser.HashPassword(DefaultAdminPassword); err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// 更新数据库中的密码
	if err := database.DAO.User.Update(ctx, adminUser); err != nil {
		log.Fatal("Failed to update admin password:", err)
	}

	fmt.Printf("✅ Admin user '%s' password has been reset to: %s\n", adminUser.Username, DefaultAdminPassword)
	fmt.Println("⚠️  Please change the default password after login for security reasons.")
}

func createAdminUser(ctx context.Context) {
	// 检查 admin 用户是否已存在
	_, err := database.DAO.User.GetByUsername(ctx, DefaultAdminUsername)
	if err == nil {
		fmt.Printf("⚠️  Admin user '%s' already exists. Use 'reset-password' to reset password.\n", DefaultAdminUsername)
		return
	}

	// 创建新的 admin 用户
	adminUser := &models.User{
		Username: DefaultAdminUsername,
		Email:    DefaultAdminEmail,
		Password: DefaultAdminPassword, // 这将在 BeforeCreate 钩子中被哈希
		IsAdmin:  true,
		Banned:   false,
	}

	if err := database.DAO.User.Create(ctx, adminUser); err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	fmt.Printf("✅ Admin user created successfully:\n")
	fmt.Printf("   Username: %s\n", DefaultAdminUsername)
	fmt.Printf("   Password: %s\n", DefaultAdminPassword)
	fmt.Printf("   Email: %s\n", DefaultAdminEmail)
	fmt.Println("⚠️  Please change the default password after login for security reasons.")
}

func initAdminUser(ctx context.Context) {
	// 检查 admin 用户是否已存在
	adminUser, err := database.DAO.User.GetByUsername(ctx, DefaultAdminUsername)
	if err != nil {
		// 用户不存在，创建一个新的
		fmt.Println("Admin user not found, creating...")
		createAdminUser(ctx)
	} else {
		// 用户存在，重置密码
		fmt.Printf("Admin user '%s' found, resetting password...\n", adminUser.Username)
		resetAdminPassword(ctx)
	}
}
