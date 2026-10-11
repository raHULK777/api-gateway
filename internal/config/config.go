package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port         string `yaml:"port"`
	DatabaseURL  string `yaml:"database_url"`
	RedisURL     string `yaml:"redis_url"`
	JWTSecret    string `yaml:"jwt_secret"`
	RateLimitRPS float64 `yaml:"rate_limit_rps"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/apigateway?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		RateLimitRPS: parseFloat(getEnv("RATE_LIMIT_RPS", "100")),
	}

	// Try to load from config.yaml if exists
	if data, err := os.ReadFile("config.yaml"); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func parseFloat(s string) float64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return 0
}
