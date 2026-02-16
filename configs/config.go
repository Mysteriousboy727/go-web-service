package configs

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisAddr  string
	JWTSecret  string
	RateLimit  int
}

// Load reads config from environment variables with sensible defaults.
// In production: values come from Kubernetes Secrets or AWS Parameter Store.
// In development: values come from a .env file or docker-compose environment block.
func Load() *Config {
	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT", "100"))
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "userdb"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:  getEnv("JWT_SECRET", "change-me-in-production"),
		RateLimit:  rateLimit,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
