package telegram

import (
	"context"
	"strings"
	"testing"
	"time"

	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

func TestProfileManagerCompletesProfileWithBirthTime(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	manager := newProfileManager(store)
	manager.now = func() time.Time {
		return time.Date(2026, time.May, 2, 5, 0, 0, 0, time.UTC)
	}

	if got := manager.start(42); got != ProfileStartText {
		t.Fatalf("start() = %q, want ProfileStartText", got)
	}

	got, handled := manager.handle(context.Background(), 42, "24.03.1992")
	if !handled || !strings.Contains(got, "Введи время рождения") {
		t.Fatalf("date handle = (%q, %v), want time prompt", got, handled)
	}

	got, handled = manager.handle(context.Background(), 42, "08:30")
	if !handled || !strings.Contains(got, "Введи город рождения") {
		t.Fatalf("time handle = (%q, %v), want city prompt", got, handled)
	}

	got, handled = manager.handle(context.Background(), 42, "Москва")
	if !handled {
		t.Fatal("city handle handled = false, want true")
	}
	for _, want := range []string{"Профиль сохранён", "Дата рождения: 24.03.1992", "Время рождения: 08:30", "Город рождения: Москва"} {
		if !strings.Contains(got, want) {
			t.Fatalf("saved text = %q, want substring %q", got, want)
		}
	}

	profile, ok, err := store.Get(context.Background(), 42)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !ok {
		t.Fatal("Get() ok = false, want true")
	}
	if profile.City != "Москва" || profile.BirthTime.String() != "08:30" {
		t.Fatalf("profile = %+v, want Moscow 08:30", profile)
	}
}

func TestProfileManagerAllowsUnknownBirthTime(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	manager := newProfileManager(store)

	manager.start(42)
	manager.handle(context.Background(), 42, "1992-03-24")
	manager.handle(context.Background(), 42, "нет")
	got, handled := manager.handle(context.Background(), 42, "Санкт-Петербург")
	if !handled {
		t.Fatal("city handle handled = false, want true")
	}
	if !strings.Contains(got, "Время рождения: не указано") {
		t.Fatalf("saved text = %q, want unknown time", got)
	}
}

func TestProfileManagerRejectsInvalidInputs(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	manager := newProfileManager(store)

	manager.start(42)
	got, handled := manager.handle(context.Background(), 42, "31.02.1992")
	if !handled || !strings.Contains(got, "Не смог распознать дату") {
		t.Fatalf("invalid date handle = (%q, %v), want date error", got, handled)
	}

	got, handled = manager.handle(context.Background(), 42, "24.03.1992")
	if !handled || !strings.Contains(got, "Введи время рождения") {
		t.Fatalf("date handle = (%q, %v), want time prompt", got, handled)
	}

	got, handled = manager.handle(context.Background(), 42, "25:99")
	if !handled || !strings.Contains(got, "Не смог распознать время") {
		t.Fatalf("invalid time handle = (%q, %v), want time error", got, handled)
	}

	got, handled = manager.handle(context.Background(), 42, "08:30")
	if !handled || !strings.Contains(got, "Введи город рождения") {
		t.Fatalf("time handle = (%q, %v), want city prompt", got, handled)
	}

	got, handled = manager.handle(context.Background(), 42, "/help")
	if !handled || !strings.Contains(got, "Введи город рождения текстом") {
		t.Fatalf("invalid city handle = (%q, %v), want city error", got, handled)
	}
}

func TestProfileManagerCancel(t *testing.T) {
	store := domainprofile.NewMemoryStore()
	manager := newProfileManager(store)

	manager.start(42)
	got := manager.cancel(42)
	if !strings.Contains(got, "Анкета отменена") {
		t.Fatalf("cancel() = %q, want cancellation", got)
	}

	_, handled := manager.handle(context.Background(), 42, "24.03.1992")
	if handled {
		t.Fatal("handle() handled = true after cancel, want false")
	}
}
