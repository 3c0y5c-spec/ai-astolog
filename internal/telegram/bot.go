package telegram

import (
	"context"
	"log/slog"

	"github.com/3c0y5c-spec/ai-astolog/internal/ai"
	domainastrology "github.com/3c0y5c-spec/ai-astolog/internal/domain/astrology"
	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service struct {
	bot      *bot.Bot
	logger   *slog.Logger
	profiles *profileManager
	asks     *askManager
	compat   *compatibilityManager
}

func New(token string, logger *slog.Logger) (*Service, error) {
	return NewWithDependencies(token, logger, domainprofile.NewMemoryStore(), nil)
}

func NewWithProfileStore(token string, logger *slog.Logger, store domainprofile.Store) (*Service, error) {
	return NewWithDependencies(token, logger, store, nil)
}

func NewWithDependencies(token string, logger *slog.Logger, store domainprofile.Store, aiClient ai.Client) (*Service, error) {
	if store == nil {
		store = domainprofile.NewMemoryStore()
	}

	service := &Service{
		logger:   logger,
		profiles: newProfileManager(store),
		asks:     newAskManager(aiClient),
		compat:   newCompatibilityManager(),
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
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/compatibility", bot.MatchTypeExact, s.handleCommand)
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
	if text == "/cancel" {
		return s.cancelReply(userID)
	}

	if reply, handled := s.profiles.handle(ctx, userID, text); handled {
		return reply
	}
	if reply, handled := s.askReply(ctx, userID, text); handled {
		return reply
	}
	if reply, handled := s.compatibilityReply(ctx, userID, text); handled {
		return reply
	}

	switch text {
	case "/profile":
		return s.profiles.start(userID)
	case "/chart":
		return s.chartReply(ctx, userID)
	case "/daily":
		return s.dailyReply(ctx, userID)
	case "/ask":
		return s.askStartReply(ctx, userID)
	case "/compatibility":
		return s.compatibilityStartReply(ctx, userID)
	}

	return ReplyForCommand(text)
}

func (s *Service) cancelReply(userID int64) string {
	if s.asks != nil && s.asks.cancel(userID) {
		return "Вопрос отменён. Отправь /ask, чтобы задать новый вопрос."
	}
	if s.compat != nil && s.compat.cancel(userID) {
		return "Расчёт совместимости отменён. Отправь /compatibility, чтобы начать заново."
	}

	return s.profiles.cancel(userID)
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

func (s *Service) dailyReply(ctx context.Context, userID int64) string {
	birthProfile, ok, err := s.profiles.get(ctx, userID)
	if err != nil {
		return "Не смог загрузить профиль. Попробуй /daily ещё раз."
	}
	if !ok {
		return "Сначала заполни анкету рождения через /profile, чтобы получить ежедневный прогноз."
	}

	return domainastrology.BuildDailyForecast(birthProfile, s.profiles.currentDate())
}

func (s *Service) askStartReply(ctx context.Context, userID int64) string {
	birthProfile, ok, err := s.profiles.get(ctx, userID)
	if err != nil {
		return "Не смог загрузить профиль. Попробуй /ask ещё раз."
	}
	if !ok {
		return "Сначала заполни анкету рождения через /profile, чтобы задать вопрос AI-астрологу."
	}
	if s.asks == nil {
		return ai.BuildFallbackAnswer(birthProfile, "общий вопрос")
	}

	return s.asks.start(userID)
}

func (s *Service) askReply(ctx context.Context, userID int64, text string) (string, bool) {
	if s.asks == nil || !s.asks.active(userID) {
		return "", false
	}

	birthProfile, ok, err := s.profiles.get(ctx, userID)
	if err != nil {
		s.asks.clear(userID)
		return "Не смог загрузить профиль. Попробуй /ask ещё раз.", true
	}
	if !ok {
		s.asks.clear(userID)
		return "Сначала заполни анкету рождения через /profile, чтобы задать вопрос AI-астрологу.", true
	}

	return s.asks.handle(ctx, userID, birthProfile, text)
}

func (s *Service) compatibilityStartReply(ctx context.Context, userID int64) string {
	_, ok, err := s.profiles.get(ctx, userID)
	if err != nil {
		return "Не смог загрузить профиль. Попробуй /compatibility ещё раз."
	}
	if !ok {
		return "Сначала заполни анкету рождения через /profile, чтобы рассчитать совместимость."
	}
	if s.compat == nil {
		return "Не смог начать расчёт совместимости. Попробуй позже."
	}

	return s.compat.start(userID)
}

func (s *Service) compatibilityReply(ctx context.Context, userID int64, text string) (string, bool) {
	if s.compat == nil || !s.compat.active(userID) {
		return "", false
	}

	birthProfile, ok, err := s.profiles.get(ctx, userID)
	if err != nil {
		s.compat.clear(userID)
		return "Не смог загрузить профиль. Попробуй /compatibility ещё раз.", true
	}
	if !ok {
		s.compat.clear(userID)
		return "Сначала заполни анкету рождения через /profile, чтобы рассчитать совместимость.", true
	}

	return s.compat.handle(ctx, userID, birthProfile, text)
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
