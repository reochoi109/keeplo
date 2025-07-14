package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Mode       string
	Port       string
	UseTLS     bool
	Debug      bool
	LogLevel   string
	HMACSecret string

	DB         DBConfig
	Recaptcha  RecaptchaConfig
	CORSOrigin []string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RecaptchaConfig struct {
	SiteKey   string
	SecretKey string
}

var AppConfig Config

func Init() {
	mode := parseMode()

	_ = godotenv.Load(".env")                           // 공통 기본값
	_ = godotenv.Overload(fmt.Sprintf(".env.%s", mode)) // 선택한 환경 오버라이드

	AppConfig = Config{
		Mode:       mode,
		Port:       get("PORT", ":8080"),
		UseTLS:     get("USE_TLS", "false") == "true",
		Debug:      get("DEBUG", "false") == "true",
		LogLevel:   get("LOG_LEVEL", "info"),
		HMACSecret: get("HMAC_SECRET", ""),

		DB: DBConfig{
			Host:     get("PG_DB_HOST", "localhost"),
			Port:     get("PG_DB_PORT", "5432"),
			User:     get("PG_DB_USER", "postgres"),
			Password: get("PG_DB_PASSWORD", ""),
			Name:     get("PG_DB_NAME", "keeplo"),
		},

		Recaptcha: RecaptchaConfig{
			SiteKey:   get("RECAPTCHA_SITE_KEY", ""),
			SecretKey: get("RECAPTCHA_SECRET_KEY", ""),
		},

		CORSOrigin: strings.Split(get("WHITE_LIST", ""), ","),
	}

	log.Printf("[Config] Loaded: mode=%s, port=%s", AppConfig.Mode, AppConfig.Port)
}

func get(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

// Data Source Name
func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}
