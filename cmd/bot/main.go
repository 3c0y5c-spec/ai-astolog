package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/3c0y5c-spec/ai-astolog/internal/ai"
	"github.com/3c0y5c-spec/ai-astolog/internal/config"
	"github.com/3c0y5c-spec/ai-astolog/internal/httpserver"
	telegrambot "github.com/3c0y5c-spec/ai-astolog/internal/telegram"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer := httpserver.New(cfg.HTTPAddr, logger)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Error("http server failed", "error", err)
			stop()
		}
	}()

	if cfg.BotMode == "webhook" {
		logger.Info("webhook mode is configured for future implementation", "public_url", cfg.PublicWebhookURL)
		<-ctx.Done()
		shutdownHTTPServer(logger, httpServer, cfg.ShutdownTimeout)
		return
	}

	aiClient := ai.NewOpenAIClient(ai.Config{
		APIKey:  cfg.AIAPIKey,
		BaseURL: cfg.AIBaseURL,
		Model:   cfg.AIModel,
		Timeout: cfg.AITimeout,
	})

	botService, err := telegrambot.NewWithDependencies(cfg.TelegramBotToken, logger, nil, aiClient)
	if err != nil {
		logger.Error("telegram bot init failed", "error", err)
		os.Exit(1)
	}

	botService.Start(ctx)
	shutdownHTTPServer(logger, httpServer, cfg.ShutdownTimeout)
}

func shutdownHTTPServer(logger *slog.Logger, httpServer *httpserver.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("http server shutdown failed", "error", err)
	}
}
