package config

import "os"

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port        string
	Env         string
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "8090"),
			Env:         getEnv("ENV", "development"),
			DatabaseURL: getEnv("DATABASE_URL", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
