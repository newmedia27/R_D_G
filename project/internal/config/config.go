package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	LogLevel          string
	MongoURI          string
	MongoDatabase     string
	MongoUser         string
	MongoUserPassword string
	JWTSecret         string
	JWTExpiration     time.Duration
	RefreshExpiration time.Duration
	AllowedOrigins    []string
	AllowCredentials  bool
	IsProduction      bool
}

type parseEnvTypes interface {
	string | bool | time.Duration
}

func getEnv[T parseEnvTypes](key string, fallback T) T {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	var res any
	switch any(fallback).(type) {
	case string:
		res = value
	case bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		res = b
	case time.Duration:
		t, err := time.ParseDuration(value)
		if err != nil {
			return fallback
		}
		res = t
	}
	return res.(T)
}

func splitEnv(key, fallback string) []string {
	v := os.Getenv(key)
	if v == "" {
		return []string{fallback}
	}
	return strings.Split(v, ",")
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		Port:              getEnv("PORT", "8000"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		MongoURI:          getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:     getEnv("MONGO_DATABASE", "app"),
		MongoUser:         getEnv("MONGO_USER", "admin"),
		MongoUserPassword: getEnv("MONGO_USER_PASSWORD", "admin"),
		JWTSecret:         getEnv("JWT_SECRET", "secret"),
		JWTExpiration:     getEnv("JWT_EXPIRATION", time.Minute*5),
		RefreshExpiration: getEnv("REFRESH_EXPIRATION", time.Hour*24*7),
		AllowedOrigins:    splitEnv("ALLOWED_ORIGINS", "*"),
		AllowCredentials:  getEnv("ALLOWED_CREDENTIALS", false),
		IsProduction:      getEnv("IS_PRODUCTION", false),
	}
	return cfg, nil
}
