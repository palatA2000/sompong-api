package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
	DBSSL           string
	CORSOrigins     string
	QuizAPIKey      string
	QuizChoiceCount int
	GeminiAPIKey    string
	GeminiModel     string
	Gemini3Model    string
	GeminiBaseURL   string
}

var App *Config

func Load() {
	_ = godotenv.Load()

	App = &Config{
		AppPort:         getEnv("APP_PORT", "3000"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPass:          getEnv("DB_PASSWORD", "password"),
		DBName:          getEnv("DB_NAME", "sompong_db"),
		DBSSL:           getEnv("DB_SSLMODE", "disable"),
		CORSOrigins:     getEnv("CORS_ORIGINS", "*"),
		QuizAPIKey:      getEnv("QUIZ_API_KEY", ""),
		QuizChoiceCount: getEnvInt("QUIZ_CHOICE_COUNT", 4),
		GeminiAPIKey:    getEnv("GEMINI_API_KEY", ""),
		GeminiModel:     getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
		Gemini3Model:    getEnv("GEMINI_3_MODEL", "gemini-3.1-flash-lite-preview"),
		GeminiBaseURL:   getEnv("GEMINI_BASE_URL", "https://generativelanguage.googleapis.com/v1beta"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSL,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}
