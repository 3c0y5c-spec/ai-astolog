package astrology

import (
	"strings"
	"testing"
	"time"

	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

func TestBuildDailyForecast(t *testing.T) {
	forecast := BuildDailyForecast(domainprofile.BirthProfile{
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		City:      "Москва",
	}, time.Date(2026, time.May, 2, 12, 0, 0, 0, time.UTC))

	for _, want := range []string{
		"Ежедневный прогноз на 02.05.2026:",
		"Солнечный знак: Овен ♈",
		"Фокус дня:",
		"Совет:",
		"Вопрос для себя:",
		"Это развлекательная астрологическая интерпретация",
	} {
		if !strings.Contains(forecast, want) {
			t.Fatalf("BuildDailyForecast() = %q, want substring %q", forecast, want)
		}
	}
}

func TestDailyForecastForProfileIsDeterministic(t *testing.T) {
	birthProfile := domainprofile.BirthProfile{
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
	}
	date := time.Date(2026, time.May, 2, 12, 0, 0, 0, time.UTC)

	first := DailyForecastForProfile(birthProfile, date)
	second := DailyForecastForProfile(birthProfile, date)

	if first != second {
		t.Fatalf("DailyForecastForProfile() = %+v and %+v, want deterministic result", first, second)
	}
	if first.Sign.Name != "Овен" {
		t.Fatalf("DailyForecastForProfile().Sign = %q, want Овен", first.Sign.Name)
	}
}
