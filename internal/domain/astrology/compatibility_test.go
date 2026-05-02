package astrology

import (
	"strings"
	"testing"
	"time"

	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

func TestBuildCompatibilityReport(t *testing.T) {
	report := BuildCompatibilityReport(domainprofile.BirthProfile{
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
	}, time.Date(1990, time.October, 2, 0, 0, 0, 0, time.UTC))

	for _, want := range []string{
		"Совместимость (MVP):",
		"Твой знак: Овен ♈",
		"Знак партнёра: Весы ♎",
		"Стихии: Огонь + Воздух",
		"Огонь и Воздух усиливают движение",
		"Качества: кардинальное + кардинальное",
		"Фокус отношений:",
		"развлекательная астрологическая интерпретация",
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("BuildCompatibilityReport() = %q, want substring %q", report, want)
		}
	}
}
