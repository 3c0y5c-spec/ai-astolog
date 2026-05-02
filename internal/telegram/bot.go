package telegram

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service struct {
	bot    *bot.Bot
	logger *slog.Logger
}

func New(token string, logger *slog.Logger) (*Service, error) {
	service := &Service{logger: logger}

	b, err := bot.New(token, bot.WithDefaultHandler(service.handleMessage))
	if err != nil {
		return nil, err
	}

	service.bot = b
	service.registerHandlers()

	return service, nil
}

func (s *Service) Start(ctx context.Context) {
	s.logger.Info("telegram bot polling started")
	s.bot.Start(ctx)
}

func (s *Service) registerHandlers() {
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, s.handleCommand)
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, s.handleCommand)
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/profile", bot.MatchTypeExact, s.handleCommand)
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/chart", bot.MatchTypeExact, s.handleCommand)
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/daily", bot.MatchTypeExact, s.handleCommand)
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/ask", bot.MatchTypeExact, s.handleCommand)
}

func (s *Service) handleCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	s.sendText(ctx, b, update.Message.Chat.ID, ReplyForCommand(update.Message.Text))
}

func (s *Service) handleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	s.sendText(ctx, b, update.Message.Chat.ID, HelpText)
}

func (s *Service) sendText(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		s.logger.Error("telegram send message failed", "error", err)
	}
}
