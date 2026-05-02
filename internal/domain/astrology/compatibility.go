package astrology

import (
	"fmt"
	"time"

	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

func BuildCompatibilityReport(birthProfile domainprofile.BirthProfile, partnerBirthDate time.Time) string {
	userSign := SunSignForDate(birthProfile.BirthDate)
	partnerSign := SunSignForDate(partnerBirthDate)
	elementInsight := elementCompatibility(userSign.Element, partnerSign.Element)
	modalityInsight := modalityCompatibility(userSign.Modality, partnerSign.Modality)

	return fmt.Sprintf(
		"Совместимость (MVP):\nТвой знак: %s %s\nЗнак партнёра: %s %s\n\nСтихии: %s + %s\n%s\n\nКачества: %s + %s\n%s\n\nФокус отношений: соединить твою тему «%s» с темой партнёра «%s».\n\nЭто развлекательная астрологическая интерпретация для саморефлексии, а не прогноз судьбы отношений.",
		userSign.Name,
		userSign.Symbol,
		partnerSign.Name,
		partnerSign.Symbol,
		userSign.Element,
		partnerSign.Element,
		elementInsight,
		userSign.Modality,
		partnerSign.Modality,
		modalityInsight,
		userSign.Theme,
		partnerSign.Theme,
	)
}

func elementCompatibility(first, second string) string {
	if first == second {
		return "Похожие стихии дают естественное чувство ритма: легче понимать реакции, желания и базовый темперамент друг друга."
	}

	switch {
	case sameElementPair(first, second, "Огонь", "Воздух"):
		return "Огонь и Воздух усиливают движение: один добавляет импульс, другой идеи и пространство для разговора."
	case sameElementPair(first, second, "Земля", "Вода"):
		return "Земля и Вода создают мягкую опору: практичность помогает чувствам становиться надёжнее, а эмпатия смягчает контроль."
	case sameElementPair(first, second, "Огонь", "Вода"):
		return "Огонь и Вода могут цеплять глубиной и страстью, но важно беречь эмоциональные границы и темп реакции."
	case sameElementPair(first, second, "Огонь", "Земля"):
		return "Огонь хочет быстро действовать, Земля — проверять почву. Союзу помогает договорённость о темпе и конкретных шагах."
	case sameElementPair(first, second, "Воздух", "Вода"):
		return "Воздух объясняет словами, Вода считывает чувства. Здесь важно не спорить с эмоциями и не молчать о мыслях."
	default:
		return "Воздух и Земля соединяют идеи и практику: притяжение растёт, когда планы превращаются в понятные действия."
	}
}

func modalityCompatibility(first, second string) string {
	if first == second {
		return "Одинаковое качество помогает стартовать в похожем стиле, но иногда усиливает упрямство или повторяющиеся сценарии."
	}

	switch {
	case sameElementPair(first, second, "кардинальное", "фиксированное"):
		return "Кардинальность задаёт направление, фиксированность удерживает курс. Важно заранее делить лидерство и зоны ответственности."
	case sameElementPair(first, second, "кардинальное", "мутабельное"):
		return "Кардинальность запускает процессы, мутабельность адаптирует маршрут. Союзу помогает гибкий план без давления."
	default:
		return "Фиксированность даёт устойчивость, мутабельность — гибкость. Баланс появляется, когда стабильность не превращается в застой."
	}
}

func sameElementPair(first, second, expectedFirst, expectedSecond string) bool {
	return (first == expectedFirst && second == expectedSecond) || (first == expectedSecond && second == expectedFirst)
}
