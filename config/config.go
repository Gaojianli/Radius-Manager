package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBPort            int
	DBUser            string
	DBPassword        string
	DBName            string
	JWTSecret         string
	ServerPort        string
	DefaultAdminUser  string
	DefaultAdminPass  string
	DefaultAdminEmail string
}

var AppConfig *Config

func LoadConfig() error {
	godotenv.Load()

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "3306"))
	if err != nil {
		dbPort = 3306
	}

	AppConfig = &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            dbPort,
		DBUser:            getEnv("DB_USER", "root"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", "radius_mgnt"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
		ServerPort:        getEnv("SERVER_PORT", ":8080"),
		DefaultAdminUser:  getEnv("DEFAULT_ADMIN_USER", "admin"),
		DefaultAdminPass:  getEnv("DEFAULT_ADMIN_PASSWORD", "admin123"),
		DefaultAdminEmail: getEnv("DEFAULT_ADMIN_EMAIL", "admin@example.com"),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
