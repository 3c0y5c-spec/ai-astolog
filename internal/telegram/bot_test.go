package telegram

import (
	"context"
	"strings"
	"testing"
	"time"

	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

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
