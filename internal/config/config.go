package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret       string
	ExpireHours  int
	ExpiresHours int
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

func (c CloudinaryConfig) IsConfigured() bool {
	return c.CloudName != "" && c.APIKey != "" && c.APISecret != ""
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

type Config struct {
	Server     ServerConfig
	JWT        JWTConfig
	DB         DBConfig
	Cloudinary CloudinaryConfig
	SMTP       SMTPConfig
}

func (d DBConfig) DSN() string {

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.Username, d.Password, d.Database, d.SSLMode)

}


// loadDotEnv подгружает .env: сначала из родительской папки (корень репо при go run из cmd/api),
// затем из текущей директории — второй файл переопределяет ключи первого.
// Ошибки «файл не найден» игнорируются, чтобы работали только переменные ОС.
func loadDotEnv() {
	paths := []string{
		filepath.Join("..", ".env"),
		".env",
	}
	for _, p := range paths {
		_ = godotenv.Load(p)
	}
}

func trimBOM(s string) string {
	return strings.TrimPrefix(strings.TrimSpace(s), "\ufeff")
}

func Load() (*Config, error) {
	loadDotEnv()

	expiresHours, err := strconv.Atoi(getEnv("JWT_EXPIRES_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRES_HOURS: %w", err)
	}

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "1234"),
			Database: getEnv("DB_NAME", "db_pizza_delivery"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:       getEnv("JWT_SECRET", "secret"),
			ExpireHours:  expiresHours,
			ExpiresHours: expiresHours,
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8085"),
		},
		Cloudinary: CloudinaryConfig{
			CloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
			APIKey:    getEnv("CLOUDINARY_API_KEY", ""),
			APISecret: getEnv("CLOUDINARY_API_SECRET", ""),
		},
		SMTP: SMTPConfig{
			Host:     trimBOM(getEnv("SMTP_HOST", "smtp.gmail.com")),
			Port:     trimBOM(getEnv("SMTP_PORT", "587")),
			User:     trimBOM(getEnv("SMTP_USER", "")),
			Password: trimBOM(getEnv("SMTP_PASSWORD", "")),
			From:     trimBOM(getEnv("SMTP_FROM", "")),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
