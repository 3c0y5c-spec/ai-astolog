package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultEnvironment = "development"
	defaultHTTPAddr    = ":8080"
	defaultBotMode     = "polling"
	defaultShutdown    = 10 * time.Second
	defaultAITimeout   = 20 * time.Second
)

type Config struct {
	Environment      string
	HTTPAddr         string
	TelegramBotToken string
	BotMode          string
	PublicWebhookURL string
	WebhookSecret    string
	ShutdownTimeout  time.Duration
	AIAPIKey         string
	AIBaseURL        string
	AIModel          string
	AITimeout        time.Duration
}

func Load() (Config, error) {
	shutdownTimeout, err := durationFromEnv("SHUTDOWN_TIMEOUT_SECONDS", defaultShutdown)
	if err != nil {
		return Config{}, err
	}
	aiTimeout, err := durationFromEnv("AI_TIMEOUT_SECONDS", defaultAITimeout)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Environment:      stringFromEnv("APP_ENV", defaultEnvironment),
		HTTPAddr:         stringFromEnv("HTTP_ADDR", defaultHTTPAddr),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		BotMode:          stringFromEnv("BOT_MODE", defaultBotMode),
		PublicWebhookURL: os.Getenv("PUBLIC_WEBHOOK_URL"),
		WebhookSecret:    os.Getenv("WEBHOOK_SECRET"),
		ShutdownTimeout:  shutdownTimeout,
		AIAPIKey:         os.Getenv("AI_API_KEY"),
		AIBaseURL:        os.Getenv("AI_BASE_URL"),
		AIModel:          os.Getenv("AI_MODEL"),
		AITimeout:        aiTimeout,
	}

	if cfg.TelegramBotToken == "" {
		return Config{}, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.BotMode != "polling" && cfg.BotMode != "webhook" {
		return Config{}, fmt.Errorf("BOT_MODE must be polling or webhook")
	}
	if cfg.BotMode == "webhook" && cfg.PublicWebhookURL == "" {
		return Config{}, fmt.Errorf("PUBLIC_WEBHOOK_URL is required when BOT_MODE=webhook")
	}

	return cfg, nil
}

func stringFromEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func durationFromEnv(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}

	return time.Duration(seconds) * time.Second, nil
}
