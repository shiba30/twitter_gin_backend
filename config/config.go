package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerAddress     string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	SmtpHost          string
	SmtpPort          string
	From              string
	SecretKey         string
	RedisAddr         string
	RedisPassword     string
	RedisDB           int
	UploadedImagesDir string `yaml:"uploaded_images_dir"`
	DefaultPageSize   string `yaml:"defaultPageSize"`
}

const config_dir = "config/config.yaml"

// 設定値読込
func LoadConfig() (Config, error) {
	var cfg Config

	// config.yml
	data, err := ioutil.ReadFile(config_dir)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal yaml data: %v", err)
	}

	// 環境変数
	cfg.ServerAddress = os.Getenv("SERVER_ADDRESS")
	if cfg.ServerAddress == "" {
		return cfg, errors.New("SERVER_ADDRESS is not set")
	}
	cfg.DBHost = os.Getenv("DB_HOST")
	if cfg.DBHost == "" {
		return cfg, errors.New("DB_HOST is not set")
	}
	cfg.DBPort = os.Getenv("DB_PORT")
	if cfg.DBPort == "" {
		return cfg, errors.New("DB_PORT is not set")
	}
	cfg.DBUser = os.Getenv("DB_USER")
	if cfg.DBUser == "" {
		return cfg, errors.New("DB_USER is not set")
	}
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	if cfg.DBPassword == "" {
		return cfg, errors.New("DB_PASSWORD is not set")
	}
	cfg.DBName = os.Getenv("DB_NAME")
	if cfg.DBName == "" {
		return cfg, errors.New("DB_NAME is not set")
	}
	cfg.SmtpHost = os.Getenv("SMTP_HOST")
	if cfg.SmtpHost == "" {
		return cfg, errors.New("SMTP_HOST is not set")
	}
	cfg.SmtpPort = os.Getenv("SMTP_PORT")
	if cfg.SmtpPort == "" {
		return cfg, errors.New("SMTP_PORT is not set")
	}
	cfg.From = os.Getenv("FROM_EMAIL")
	if cfg.From == "" {
		return cfg, errors.New("FROM_EMAIL is not set")
	}
	cfg.SecretKey = os.Getenv("SECRET_KEY")
	if cfg.SecretKey == "" {
		return cfg, errors.New("SECRET_KEY is not set")
	}
	cfg.RedisAddr = os.Getenv("REDIS_ADDRESS")
	if cfg.RedisAddr == "" {
		return cfg, errors.New("REDIS_ADDRESS is not set")
	}
	cfg.RedisPassword = os.Getenv("REDIS_PASSWORD")
	cfg.RedisDB, err = strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return cfg, fmt.Errorf("failed to convert REDIS_DB to int: %v", err)
	}

	return cfg, nil
}
