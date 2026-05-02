package telegram

import (
	"context"
	"log/slog"

	domainastrology "github.com/3c0y5c-spec/ai-astolog/internal/domain/astrology"
	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service struct {
	bot      *bot.Bot
	logger   *slog.Logger
	profiles *profileManager
}

func New(token string, logger *slog.Logger) (*Service, error) {
	return NewWithProfileStore(token, logger, domainprofile.NewMemoryStore())
}

func NewWithProfileStore(token string, logger *slog.Logger, store domainprofile.Store) (*Service, error) {
	service := &Service{
		logger:   logger,
		profiles: newProfileManager(store),
	}

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
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/cancel", bot.MatchTypeExact, s.handleCommand)
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/chart", bot.MatchTypeExact, s.handleCommand)
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/daily", bot.MatchTypeExact, s.handleCommand)
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/ask", bot.MatchTypeExact, s.handleCommand)
}

func (s *Service) handleCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	message := update.Message
	s.sendText(ctx, b, message.Chat.ID, s.replyForText(ctx, userIDFromMessage(message), message.Text))
}

func (s *Service) handleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	message := update.Message
	s.sendText(ctx, b, message.Chat.ID, s.replyForText(ctx, userIDFromMessage(message), message.Text))
}

func (s *Service) replyForText(ctx context.Context, userID int64, text string) string {
	switch text {
	case "/profile":
		return s.profiles.start(userID)
	case "/cancel":
		return s.profiles.cancel(userID)
	case "/chart":
		return s.chartReply(ctx, userID)
	}

	if reply, handled := s.profiles.handle(ctx, userID, text); handled {
		return reply
	}

	return ReplyForCommand(text)
}

func (s *Service) chartReply(ctx context.Context, userID int64) string {
	birthProfile, ok, err := s.profiles.get(ctx, userID)
	if err != nil {
		return "Не смог загрузить профиль. Попробуй /chart ещё раз."
	}
	if !ok {
		return "Сначала заполни анкету рождения через /profile, чтобы построить натальную карту."
	}

	return domainastrology.BuildNatalSummary(birthProfile)
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

func userIDFromMessage(message *models.Message) int64 {
	if message.From != nil {
		return message.From.ID
	}
	return message.Chat.ID
}
