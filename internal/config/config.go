package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config contém toda a configuração da aplicação
type Config struct {
	Port        string
	Environment string
	Database    DatabaseConfig
	JWT         JWTConfig
}

// JWTConfig contém configurações de autenticação JWT
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// DatabaseConfig contém configurações do banco
type DatabaseConfig struct {
	URL            string
	MaxConnections int
	MaxIdleConns   int
	MaxLifetime    time.Duration
}

// Load carrega as configurações do ambiente
func Load() *Config {
	// Carregar .env apenas em desenvolvimento
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	return &Config{
		Port:        getEnvOrDefault("PORT", "8080"),
		Environment: getEnvOrDefault("ENV", "development"),
		Database: DatabaseConfig{
			URL:            getEnvOrDefault("DATABASE_URL", ""),
			MaxConnections: getEnvAsInt("DB_MAX_CONNECTIONS", 25),
			MaxIdleConns:   getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
			MaxLifetime:    time.Duration(getEnvAsInt("DB_MAX_LIFETIME_MINUTES", 5)) * time.Minute,
		},
		JWT: JWTConfig{
			Secret: getEnvOrDefault("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiry: time.Duration(getEnvAsInt("JWT_EXPIRY_HOURS", 24)) * time.Hour,
		},
	}
}

// getEnvOrDefault retorna o valor da env var ou o padrão
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retorna o valor da env var como int ou o padrão
func getEnvAsInt(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}

	return defaultValue
}
