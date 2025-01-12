package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
}

type AppConfig struct {
	Port string
}

type PostgresConfig struct {
	User   string
	Pass   string
	Dbname string
	Port   string
	Host   string
}

func InitConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error reading it; relying on environment variables.")
	}

	cfg := &Config{
		App: AppConfig{
			getEnv("PORT", "8080"),
		},
		Postgres: PostgresConfig{
			User:   os.Getenv("DB_USERNAME"),
			Pass:   os.Getenv("DB_PASSWORD"),
			Dbname: os.Getenv("DB_DATABASE"),
			Port:   os.Getenv("DB_PORT"),
			Host:   os.Getenv("DB_HOST"),
		},
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	return v
}
