package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	RedisURL      string
	FrontendURL   string
	LogLevel      string
	Environment   string

	DBMaxOpenConns int
	DBMaxIdleConns int
	DBMaxIdleTime  string

	// Конфигурация JWT
	JWTSecret            string
	JWTAccessExpiryMin   int
	JWTRefreshExpiryDays int
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		ServerAddress:        getEnv("SERVER_ADDRESS", ":8080"),
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		RedisURL:             getEnv("REDIS_URL", "redis://localhost:6379/0"),
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:8081"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		Environment:          getEnv("ENVIRONMENT", "development"),
		DBMaxOpenConns:       getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:       getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBMaxIdleTime:        getEnv("DB_MAX_IDLE_TIME", "15m"),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTAccessExpiryMin:   getEnvInt("JWT_ACCESS_EXPIRY_MIN", 15),
		JWTRefreshExpiryDays: getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
