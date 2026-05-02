package astrology

import (
	"strings"
	"testing"
	"time"

	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

func TestSunSignForDate(t *testing.T) {
	tests := map[string]string{
		"1992-03-20": "Рыбы",
		"1992-03-21": "Овен",
		"1992-04-19": "Овен",
		"1992-04-20": "Телец",
		"1992-12-21": "Стрелец",
		"1992-12-22": "Козерог",
		"1992-01-19": "Козерог",
		"1992-01-20": "Водолей",
	}

	for rawDate, want := range tests {
		t.Run(rawDate, func(t *testing.T) {
			date, err := time.Parse("2006-01-02", rawDate)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			got := SunSignForDate(date)
			if got.Name != want {
				t.Fatalf("SunSignForDate(%s) = %q, want %q", rawDate, got.Name, want)
			}
		})
	}
}

func TestBuildNatalSummary(t *testing.T) {
	birthTime := domainprofile.CivilTime{Hour: 8, Minute: 30}
	summary := BuildNatalSummary(domainprofile.BirthProfile{
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		BirthTime: &birthTime,
		City:      "Москва",
	})

	for _, want := range []string{
		"Натальная карта (MVP):",
		"Дата рождения: 24.03.1992",
		"Время рождения: 08:30",
		"Город рождения: Москва",
		"Солнечный знак: Овен ♈",
		"Луна, Асцендент, дома и аспекты появятся в следующем этапе.",
	} {
		if !strings.Contains(summary, want) {
			t.Fatalf("BuildNatalSummary() = %q, want substring %q", summary, want)
		}
	}
}

func TestBuildNatalSummaryWithoutBirthTime(t *testing.T) {
	summary := BuildNatalSummary(domainprofile.BirthProfile{
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		City:      "Москва",
	})

	if !strings.Contains(summary, "Время рождения: не указано") {
		t.Fatalf("BuildNatalSummary() = %q, want unknown time", summary)
	}
}
