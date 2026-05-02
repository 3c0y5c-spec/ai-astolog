package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Environment != defaultEnvironment {
		t.Fatalf("Environment = %q, want %q", cfg.Environment, defaultEnvironment)
	}
	if cfg.HTTPAddr != defaultHTTPAddr {
		t.Fatalf("HTTPAddr = %q, want %q", cfg.HTTPAddr, defaultHTTPAddr)
	}
	if cfg.BotMode != defaultBotMode {
		t.Fatalf("BotMode = %q, want %q", cfg.BotMode, defaultBotMode)
	}
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Fatalf("ShutdownTimeout = %s, want 10s", cfg.ShutdownTimeout)
	}
	if cfg.AITimeout != 20*time.Second {
		t.Fatalf("AITimeout = %s, want 20s", cfg.AITimeout)
	}
}

func TestLoadRequiresTelegramToken(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want token error")
	}
}

func TestLoadValidatesWebhookURL(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("BOT_MODE", "webhook")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want webhook URL error")
	}
}

func TestLoadReadsShutdownTimeout(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "3")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.ShutdownTimeout != 3*time.Second {
		t.Fatalf("ShutdownTimeout = %s, want 3s", cfg.ShutdownTimeout)
	}
}

func TestLoadReadsAIConfig(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("AI_API_KEY", "ai-key")
	t.Setenv("AI_BASE_URL", "https://example.test/v1")
	t.Setenv("AI_MODEL", "custom-model")
	t.Setenv("AI_TIMEOUT_SECONDS", "7")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AIAPIKey != "ai-key" {
		t.Fatalf("AIAPIKey = %q, want configured key", cfg.AIAPIKey)
	}
	if cfg.AIBaseURL != "https://example.test/v1" {
		t.Fatalf("AIBaseURL = %q, want configured base URL", cfg.AIBaseURL)
	}
	if cfg.AIModel != "custom-model" {
		t.Fatalf("AIModel = %q, want configured model", cfg.AIModel)
	}
	if cfg.AITimeout != 7*time.Second {
		t.Fatalf("AITimeout = %s, want 7s", cfg.AITimeout)
	}
}
