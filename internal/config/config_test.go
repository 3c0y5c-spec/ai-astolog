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
}

func TestLoadRequiresTelegramToken(t *testing.T) {
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
