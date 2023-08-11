package config

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SmtpHost   string
	SmtpPort   string
	From       string
	SecretKey  string
}

// 環境変数から設定値読込
func LoadConfig() Config {
	var cfg Config

	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBPort = os.Getenv("DB_PORT")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.SmtpHost = os.Getenv("SMTP_HOST")
	cfg.SmtpPort = os.Getenv("SMTP_PORT")
	cfg.From = os.Getenv("FROM_EMAIL")
	cfg.SecretKey = os.Getenv("SECRET_KEY")

	return cfg
}
