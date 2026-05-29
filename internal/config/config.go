package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type RetryConfig struct {
	Max   int           `mapstructure:"max" validate:"required"`
	Delay time.Duration `mapstructure:"delay" validate:"required"`
}

type Config struct {
	AppPort             string `validate:"required,numeric"`
	AppEnv              string `validate:"required"`
	AppHost             string `validate:"required"`
	DBHost              string `validate:"required"`
	DBPort              string `validate:"required"`
	DBUsername          string `validate:"required"`
	DBPassword          string `validate:"required"`
	DBDatabase          string `validate:"required"`
	DBTimezone          string `validate:"required"`
	LogLevel            string `validate:"required"`
	JwtSecretKey        string `validate:"required"`
	SupabaseURL         string `validate:"required"`
	SupabaseKey         string `validate:"required"`
	GoogleClientID      string
	GoogleSecret        string
	GoogleCallbackUrl   string
	GithubClientID      string
	GithubSecret        string
	GithubCallbackUrl   string
	FrontendCallbackURL string `validate:"required,url"`
	AdminPass           string `validate:"required"`
	AdminEmail          string `validate:"required"`
	XenditSecretKey     string
	MailerAPIKey        string
	MailerInboxID       string
	MailerDriver        string
	RecommendationURL    string
	RecommendationAPIKey string
	RetryConfig
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file is missing")
	}

	delay, _ := time.ParseDuration(os.Getenv("RETRY_DELAY"))
	if delay == 0 {
		delay = 5 * time.Second
	}

	maxRetry, _ := strconv.Atoi(os.Getenv("RETRY_MAX"))

	return &Config{
		AppPort:             os.Getenv("APP_PORT"),
		AppEnv:              os.Getenv("APP_ENV"),
		AppHost:             os.Getenv("APP_HOST"),
		DBHost:              os.Getenv("DB_HOST"),
		DBPort:              os.Getenv("DB_PORT"),
		DBUsername:          os.Getenv("DB_USERNAME"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBDatabase:          os.Getenv("DB_NAME"),
		DBTimezone:          os.Getenv("DB_TIMEZONE"),
		LogLevel:            os.Getenv("LOG_LEVEL"),
		JwtSecretKey:        os.Getenv("JWT_SECRET_KEY"),
		SupabaseURL:         os.Getenv("SUPABASE_URL"),
		SupabaseKey:         os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleSecret:        os.Getenv("GOOGLE_SECRET"),
		GoogleCallbackUrl:   os.Getenv("GOOGLE_CALLBACK_URL"),
		GithubClientID:      os.Getenv("GITHUB_CLIENT_ID"),
		GithubSecret:        os.Getenv("GITHUB_SECRET"),
		GithubCallbackUrl:   os.Getenv("GITHUB_CALLBACK_URL"),
		FrontendCallbackURL: os.Getenv("FRONTEND_CALLBACK_URL"),
		AdminPass:           os.Getenv("ADMIN_PASS"),
		AdminEmail:          os.Getenv("ADMIN_EMAIL"),
		XenditSecretKey:     os.Getenv("XENDIT_SECRET_KEY"),
		MailerAPIKey:        os.Getenv("MAILER_API_KEY"),
		MailerInboxID:       os.Getenv("MAILER_INBOX_ID"),
		MailerDriver:        os.Getenv("MAILER_DRIVER"),
		RecommendationURL:    os.Getenv("RECOMMENDATION_URL"),
		RecommendationAPIKey: os.Getenv("RECOMMENDATION_API_KEY"),
		RetryConfig: RetryConfig{
			Max:   maxRetry,
			Delay: delay,
		},
	}

}

func (c *Config) Validate(v *validator.Validate) {
	if err := v.Struct(c); err != nil {
		log.Fatalf("invalid data config: %v", err)
	}
}
