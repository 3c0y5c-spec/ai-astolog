package telegram

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/3c0y5c-spec/ai-astolog/internal/ai"
	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

type stubAIClient struct {
	answer string
	err    error
	got    ai.QuestionRequest
}

func (c *stubAIClient) AnswerQuestion(_ context.Context, request ai.QuestionRequest) (string, error) {
	c.got = request
	return c.answer, c.err
}

type failingProfileStore struct{}

func (s failingProfileStore) Save(context.Context, domainprofile.BirthProfile) error {
	return nil
}

func (s failingProfileStore) Get(context.Context, int64) (domainprofile.BirthProfile, bool, error) {
	return domainprofile.BirthProfile{}, false, errors.New("store unavailable")
}

type disappearingProfileStore struct {
	profile domainprofile.BirthProfile
	calls   int
}

func (s *disappearingProfileStore) Save(_ context.Context, birthProfile domainprofile.BirthProfile) error {
	s.profile = birthProfile
	return nil
}

func (s *disappearingProfileStore) Get(context.Context, int64) (domainprofile.BirthProfile, bool, error) {
	s.calls++
	if s.calls == 1 {
		return s.profile, true, nil
	}
	return domainprofile.BirthProfile{}, false, nil
}

func TestServiceRoutesCommandsThroughActiveProfileFlow(t *testing.T) {
	service := &Service{
		profiles: newProfileManager(domainprofile.NewMemoryStore()),
	}

	ctx := context.Background()
	userID := int64(42)

	service.replyForText(ctx, userID, "/profile")
	service.replyForText(ctx, userID, "24.03.1992")
	service.replyForText(ctx, userID, "08:30")

	got := service.replyForText(ctx, userID, "/help")
	if !strings.Contains(got, "Введи город рождения текстом") {
		t.Fatalf("replyForText(/help during city step) = %q, want city validation", got)
	}

	got = service.replyForText(ctx, userID, "Москва")
	if !strings.Contains(got, "Профиль сохранён") {
		t.Fatalf("replyForText(city) = %q, want saved profile", got)
	}
}

func TestServiceCancelStopsActiveProfileFlow(t *testing.T) {
	service := &Service{
		profiles: newProfileManager(domainprofile.NewMemoryStore()),
	}

	ctx := context.Background()
	userID := int64(42)

	service.replyForText(ctx, userID, "/profile")
	got := service.replyForText(ctx, userID, "/cancel")
	if !strings.Contains(got, "Анкета отменена") {
		t.Fatalf("replyForText(/cancel) = %q, want cancellation", got)
	}

	got = service.replyForText(ctx, userID, "24.03.1992")
	if got != HelpText {
		t.Fatalf("replyForText(date after cancel) = %q, want HelpText", got)
	}
}

func TestServiceChartRequiresProfile(t *testing.T) {
	service := &Service{
		profiles: newProfileManager(domainprofile.NewMemoryStore()),
	}

	got := service.replyForText(context.Background(), 42, "/chart")
	if !strings.Contains(got, "Сначала заполни анкету рождения через /profile") {
		t.Fatalf("replyForText(/chart) = %q, want missing profile prompt", got)
	}
}

func TestServiceChartUsesSavedProfile(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	service := &Service{
		profiles: newProfileManager(store),
	}
	birthTime := domainprofile.CivilTime{Hour: 8, Minute: 30}
	err := store.Save(context.Background(), domainprofile.BirthProfile{
		UserID:    42,
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		BirthTime: &birthTime,
		City:      "Москва",
	})
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got := service.replyForText(context.Background(), 42, "/chart")
	for _, want := range []string{"Натальная карта (MVP):", "Солнечный знак: Овен ♈", "Дата рождения: 24.03.1992", "Город рождения: Москва"} {
		if !strings.Contains(got, want) {
			t.Fatalf("replyForText(/chart) = %q, want substring %q", got, want)
		}
	}
}

func TestServiceDailyRequiresProfile(t *testing.T) {
	service := &Service{
		profiles: newProfileManager(domainprofile.NewMemoryStore()),
	}

	got := service.replyForText(context.Background(), 42, "/daily")
	if !strings.Contains(got, "Сначала заполни анкету рождения через /profile") {
		t.Fatalf("replyForText(/daily) = %q, want missing profile prompt", got)
	}
}

func TestServiceDailyUsesSavedProfile(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	service := &Service{
		profiles: newProfileManager(store),
	}
	service.profiles.now = func() time.Time {
		return time.Date(2026, time.May, 2, 12, 0, 0, 0, time.UTC)
	}
	err := store.Save(context.Background(), domainprofile.BirthProfile{
		UserID:    42,
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		City:      "Москва",
	})
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got := service.replyForText(context.Background(), 42, "/daily")
	for _, want := range []string{"Ежедневный прогноз на 02.05.2026:", "Солнечный знак: Овен ♈", "Фокус дня:", "Совет:", "Вопрос для себя:"} {
		if !strings.Contains(got, want) {
			t.Fatalf("replyForText(/daily) = %q, want substring %q", got, want)
		}
	}
}

func TestServiceAskRequiresProfile(t *testing.T) {
	service := &Service{
		profiles: newProfileManager(domainprofile.NewMemoryStore()),
		asks:     newAskManager(nil),
	}

	got := service.replyForText(context.Background(), 42, "/ask")
	if !strings.Contains(got, "Сначала заполни анкету рождения через /profile") {
		t.Fatalf("replyForText(/ask) = %q, want missing profile prompt", got)
	}
}

func TestServiceAskUsesFallbackWithoutAIKey(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	service := &Service{
		profiles: newProfileManager(store),
		asks:     newAskManager(nil),
	}
	err := store.Save(context.Background(), domainprofile.BirthProfile{
		UserID:    42,
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		City:      "Москва",
	})
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got := service.replyForText(context.Background(), 42, "/ask")
	if !strings.Contains(got, "Задай вопрос AI-астрологу") {
		t.Fatalf("replyForText(/ask) = %q, want question prompt", got)
	}

	got = service.replyForText(context.Background(), 42, "Что важно сегодня?")
	for _, want := range []string{"AI-провайдер пока не настроен", "Вопрос: Что важно сегодня?", "Овен ♈"} {
		if !strings.Contains(got, want) {
			t.Fatalf("replyForText(question) = %q, want substring %q", got, want)
		}
	}
}

func TestServiceAskUsesAIClient(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	aiClient := &stubAIClient{answer: "AI ответ"}
	service := &Service{
		profiles: newProfileManager(store),
		asks:     newAskManager(aiClient),
	}
	err := store.Save(context.Background(), domainprofile.BirthProfile{
		UserID:    42,
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		City:      "Москва",
	})
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	service.replyForText(context.Background(), 42, "/ask")
	got := service.replyForText(context.Background(), 42, "Что важно сегодня?")
	if got != "AI ответ" {
		t.Fatalf("replyForText(question) = %q, want AI answer", got)
	}
	if aiClient.got.Question != "Что важно сегодня?" {
		t.Fatalf("AI question = %q, want submitted question", aiClient.got.Question)
	}
	if !strings.Contains(aiClient.got.ProfileContext, "Солнечный знак: Овен ♈") {
		t.Fatalf("AI profile context = %q, want Aries context", aiClient.got.ProfileContext)
	}
}

func TestServiceAskFallsBackWhenAINotConfigured(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	service := &Service{
		profiles: newProfileManager(store),
		asks:     newAskManager(&stubAIClient{err: ai.ErrNotConfigured}),
	}
	err := store.Save(context.Background(), domainprofile.BirthProfile{
		UserID:    42,
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		City:      "Москва",
	})
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	service.replyForText(context.Background(), 42, "/ask")
	got := service.replyForText(context.Background(), 42, "Что важно сегодня?")
	if !strings.Contains(got, "AI-провайдер пока не настроен") {
		t.Fatalf("replyForText(question) = %q, want fallback", got)
	}
}

func TestServiceAskHandlesAIError(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	service := &Service{
		profiles: newProfileManager(store),
		asks:     newAskManager(&stubAIClient{err: errors.New("provider down")}),
	}
	err := store.Save(context.Background(), domainprofile.BirthProfile{
		UserID:    42,
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		City:      "Москва",
	})
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	service.replyForText(context.Background(), 42, "/ask")
	got := service.replyForText(context.Background(), 42, "Что важно сегодня?")
	if !strings.Contains(got, "Не смог получить ответ AI-астролога") {
		t.Fatalf("replyForText(question) = %q, want AI error", got)
	}
}

func TestServiceCancelStopsActiveAskFlow(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	service := &Service{
		profiles: newProfileManager(store),
		asks:     newAskManager(nil),
	}
	err := store.Save(context.Background(), domainprofile.BirthProfile{
		UserID:    42,
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		City:      "Москва",
	})
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	service.replyForText(context.Background(), 42, "/ask")
	got := service.replyForText(context.Background(), 42, "/cancel")
	if !strings.Contains(got, "Вопрос отменён") {
		t.Fatalf("replyForText(/cancel) = %q, want ask cancellation", got)
	}

	got = service.replyForText(context.Background(), 42, "Что важно сегодня?")
	if got != HelpText {
		t.Fatalf("replyForText(question after cancel) = %q, want HelpText", got)
	}
}

func TestServiceAskClearsSessionWhenProfileGetFails(t *testing.T) {
	service := &Service{
		profiles: newProfileManager(failingProfileStore{}),
		asks:     newAskManager(nil),
	}
	userID := int64(42)
	service.asks.start(userID)

	got := service.replyForText(context.Background(), userID, "Что важно сегодня?")
	if !strings.Contains(got, "Не смог загрузить профиль") {
		t.Fatalf("replyForText(question) = %q, want profile load error", got)
	}

	got = service.replyForText(context.Background(), userID, "Ещё вопрос")
	if got != HelpText {
		t.Fatalf("replyForText(after failure) = %q, want HelpText after cleared session", got)
	}
}

func TestServiceAskClearsSessionWhenProfileDisappears(t *testing.T) {
	store := &disappearingProfileStore{
		profile: domainprofile.BirthProfile{
			UserID:    42,
			BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
			City:      "Москва",
		},
	}
	service := &Service{
		profiles: newProfileManager(store),
		asks:     newAskManager(nil),
	}
	userID := int64(42)

	got := service.replyForText(context.Background(), userID, "/ask")
	if !strings.Contains(got, "Задай вопрос AI-астрологу") {
		t.Fatalf("replyForText(/ask) = %q, want question prompt", got)
	}

	got = service.replyForText(context.Background(), userID, "Что важно сегодня?")
	if !strings.Contains(got, "Сначала заполни анкету рождения через /profile") {
		t.Fatalf("replyForText(question) = %q, want missing profile prompt", got)
	}

	got = service.replyForText(context.Background(), userID, "Ещё вопрос")
	if got != HelpText {
		t.Fatalf("replyForText(after missing profile) = %q, want HelpText after cleared session", got)
	}
}
