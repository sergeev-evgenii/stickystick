package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	JWTSecret    string
	Port         string
	Env          string
	UploadDir    string
	BaseURL      string
	VKToken          string
	VKGroupID        string
	TelegramBotToken string
	TelegramChatID   string
	MaxBotToken      string
	MaxChatID        string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/stickystick?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Port:        getEnv("PORT", "4000"),
		Env:         getEnv("ENV", "development"),
		UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
		BaseURL:     getEnv("BASE_URL", "0.0.0.0:4000"),
		VKToken:          getEnv("VK_TOKEN", "vk1.a.F7c2Dp34TzTrYERFcdcwGu66C-V770jWmHElrafFxS7LEaiQ4ym3vTAwqD0k3bwar0qFa2C8flnD55JSXJoJ7EpoDfZoIkr5zyPQ4I3CuC03CCINrO-afDtxMRfeGqb0zzY4wqnXFuNMZPAFdt5ZjTYqPh3FsWk2G1-G17cIaQ8I3qACGqZA9dHGPzqtICmltppV9LfreqPqJBBAHrgnnw"),
		VKGroupID:         getEnv("VK_GROUP_ID", "236352692"),
		TelegramBotToken:  getEnv("TELEGRAM_BOT_TOKEN", "8723537638:AAGDNfwlTGQ8Ue_3RTxAiMIt50KHRNP_8ak"),
		TelegramChatID:    getEnv("TELEGRAM_CHAT_ID", "-1002030852315"),
		MaxBotToken:       getEnv("MAX_BOT_TOKEN", ""),
		MaxChatID:         getEnv("MAX_CHAT_ID", ""),
	}, nil
}


func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
