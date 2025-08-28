package config

import (
	"os"

	tools "github.com/Wuli-Giao-Giao/tools/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func LoadConfig() {
	_ = godotenv.Load()

	cfg := &Config{}

	// 建立 logger (先用預設 info，讀檔後再覆蓋)
	Logger = tools.NewLogrusLogger("info", os.Stdout)

	// 讀取 YAML + 環境變數
	err := cleanenv.ReadConfig("./config.yaml", cfg)
	if err != nil {
		Logger.Errorf("Failed to read config: %v", err)
	}

	// 根據設定檔更新 logger level
	Logger = tools.NewLogrusLogger(cfg.Logger.Level, os.Stdout)

	AppConfig = cfg
}
