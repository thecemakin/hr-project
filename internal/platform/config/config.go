package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv            string
	HTTPPort          string
	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPassword        string
	DBSSLMode         string
	JWTSecret         string
	JWTAccessTokenTTL string
}

func Load() (*Config, error) {
	// Try loading .env file, ignore error if it doesn't exist (e.g., in production)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		HTTPPort:          getEnv("HTTP_PORT", "8080"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBName:            getEnv("DB_NAME", "hr_project"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		JWTSecret:         getEnv("JWT_SECRET", "super-secret-key-change-me"),
		JWTAccessTokenTTL: getEnv("JWT_ACCESS_TOKEN_TTL", "15m"),
	}

	// Security warning for production
	if cfg.AppEnv == "production" && cfg.JWTSecret == "super-secret-key-change-me" {
		log.Println("[WARNING] JWT_SECRET is set to the default value in production. This is HIGHLY INSECURE.")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	
	// Only log warnings for non-critical fallbacks or when specific keys are missing
	if fallback != "" {
		log.Printf("[INFO] Configuration: %s not set, using default: %s", key, fallback)
	}
	return fallback
}
