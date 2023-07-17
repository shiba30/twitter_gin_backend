package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
}

// 設定ファイルから値を取得する
func LoadConfig(filename string) (Config, error) {
	var cfg Config

	cfgFile, err := os.Open(filename)
	if err != nil {
		return cfg, err
	}
	defer cfgFile.Close()

	jsonParser := json.NewDecoder(cfgFile)
	err = jsonParser.Decode(&cfg)
	return cfg, err
}
