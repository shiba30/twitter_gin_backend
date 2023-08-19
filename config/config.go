package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	SmtpHost      string
	SmtpPort      string
	From          string
	SecretKey     string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

// 環境変数から設定値読込
func LoadConfig() Config {
	var cfg Config
	var err error

	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBPort = os.Getenv("DB_PORT")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.SmtpHost = os.Getenv("SMTP_HOST")
	cfg.SmtpPort = os.Getenv("SMTP_PORT")
	cfg.From = os.Getenv("FROM_EMAIL")
	cfg.SecretKey = os.Getenv("SECRET_KEY")
	cfg.RedisAddr = os.Getenv("REDIS_ADDRESS")
	cfg.RedisPassword = os.Getenv("REDIS_PASSWORD")
	cfg.RedisDB, err = strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Failed to convert REDIS_DB to int: %v", err)
	}

	return cfg
}
