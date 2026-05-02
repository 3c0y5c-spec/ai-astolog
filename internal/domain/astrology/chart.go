package astrology

import (
	"fmt"
	"time"

	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

type SunSign struct {
	Name        string
	Symbol      string
	Element     string
	Modality    string
	Theme       string
	Description string
}

func BuildNatalSummary(birthProfile domainprofile.BirthProfile) string {
	sign := SunSignForDate(birthProfile.BirthDate)

	birthTime := "не указано"
	if birthProfile.BirthTime != nil {
		birthTime = birthProfile.BirthTime.String()
	}

	return fmt.Sprintf(
		"Натальная карта (MVP):\nДата рождения: %s\nВремя рождения: %s\nГород рождения: %s\n\nСолнечный знак: %s %s\nСтихия: %s\nКачество: %s\nКлючевая тема: %s\n\n%s\n\nЛуна, Асцендент, дома и аспекты появятся в следующем этапе.",
		birthProfile.BirthDate.Format("02.01.2006"),
		birthTime,
		birthProfile.City,
		sign.Name,
		sign.Symbol,
		sign.Element,
		sign.Modality,
		sign.Theme,
		sign.Description,
	)
}

func SunSignForDate(date time.Time) SunSign {
	month := int(date.Month())
	day := date.Day()

	switch {
	case afterOrEqual(month, day, 3, 21) && beforeOrEqual(month, day, 4, 19):
		return signs[0]
	case afterOrEqual(month, day, 4, 20) && beforeOrEqual(month, day, 5, 20):
		return signs[1]
	case afterOrEqual(month, day, 5, 21) && beforeOrEqual(month, day, 6, 20):
		return signs[2]
	case afterOrEqual(month, day, 6, 21) && beforeOrEqual(month, day, 7, 22):
		return signs[3]
	case afterOrEqual(month, day, 7, 23) && beforeOrEqual(month, day, 8, 22):
		return signs[4]
	case afterOrEqual(month, day, 8, 23) && beforeOrEqual(month, day, 9, 22):
		return signs[5]
	case afterOrEqual(month, day, 9, 23) && beforeOrEqual(month, day, 10, 22):
		return signs[6]
	case afterOrEqual(month, day, 10, 23) && beforeOrEqual(month, day, 11, 21):
		return signs[7]
	case afterOrEqual(month, day, 11, 22) && beforeOrEqual(month, day, 12, 21):
		return signs[8]
	case afterOrEqual(month, day, 12, 22) || beforeOrEqual(month, day, 1, 19):
		return signs[9]
	case afterOrEqual(month, day, 1, 20) && beforeOrEqual(month, day, 2, 18):
		return signs[10]
	default:
		return signs[11]
	}
}

func afterOrEqual(month, day, startMonth, startDay int) bool {
	return month > startMonth || (month == startMonth && day >= startDay)
}

func beforeOrEqual(month, day, endMonth, endDay int) bool {
	return month < endMonth || (month == endMonth && day <= endDay)
}

var signs = []SunSign{
	{
		Name:        "Овен",
		Symbol:      "♈",
		Element:     "Огонь",
		Modality:    "кардинальное",
		Theme:       "инициатива, смелость и быстрый старт",
		Description: "Солнце в Овне подчёркивает прямоту, самостоятельность и желание действовать первым.",
	},
	{
		Name:        "Телец",
		Symbol:      "♉",
		Element:     "Земля",
		Modality:    "фиксированное",
		Theme:       "устойчивость, ценности и телесный комфорт",
		Description: "Солнце в Тельце подчёркивает терпение, практичность и умение создавать надёжную опору.",
	},
	{
		Name:        "Близнецы",
		Symbol:      "♊",
		Element:     "Воздух",
		Modality:    "мутабельное",
		Theme:       "общение, обучение и гибкость мышления",
		Description: "Солнце в Близнецах подчёркивает любознательность, лёгкость контакта и интерес к разным идеям.",
	},
	{
		Name:        "Рак",
		Symbol:      "♋",
		Element:     "Вода",
		Modality:    "кардинальное",
		Theme:       "эмоциональная безопасность, дом и близость",
		Description: "Солнце в Раке подчёркивает чувствительность, заботу и связь с личной историей.",
	},
	{
		Name:        "Лев",
		Symbol:      "♌",
		Element:     "Огонь",
		Modality:    "фиксированное",
		Theme:       "самовыражение, творчество и признание",
		Description: "Солнце во Льве подчёркивает щедрость, харизму и потребность проявляться ярко.",
	},
	{
		Name:        "Дева",
		Symbol:      "♍",
		Element:     "Земля",
		Modality:    "мутабельное",
		Theme:       "порядок, польза и внимательность к деталям",
		Description: "Солнце в Деве подчёркивает наблюдательность, практическую помощь и стремление улучшать процессы.",
	},
	{
		Name:        "Весы",
		Symbol:      "♎",
		Element:     "Воздух",
		Modality:    "кардинальное",
		Theme:       "баланс, партнёрство и эстетика",
		Description: "Солнце в Весах подчёркивает дипломатичность, чувство меры и талант видеть разные стороны ситуации.",
	},
	{
		Name:        "Скорпион",
		Symbol:      "♏",
		Element:     "Вода",
		Modality:    "фиксированное",
		Theme:       "глубина, трансформация и внутренняя сила",
		Description: "Солнце в Скорпионе подчёркивает интенсивность, проницательность и готовность идти в суть.",
	},
	{
		Name:        "Стрелец",
		Symbol:      "♐",
		Element:     "Огонь",
		Modality:    "мутабельное",
		Theme:       "смысл, свобода и расширение горизонтов",
		Description: "Солнце в Стрельце подчёркивает оптимизм, тягу к знаниям и поиск большого направления.",
	},
	{
		Name:        "Козерог",
		Symbol:      "♑",
		Element:     "Земля",
		Modality:    "кардинальное",
		Theme:       "структура, ответственность и долгий результат",
		Description: "Солнце в Козероге подчёркивает дисциплину, стратегичность и уважение к реальным достижениям.",
	},
	{
		Name:        "Водолей",
		Symbol:      "♒",
		Element:     "Воздух",
		Modality:    "фиксированное",
		Theme:       "идеи, сообщества и независимое мышление",
		Description: "Солнце в Водолее подчёркивает оригинальность, дружелюбие и интерес к будущему.",
	},
	{
		Name:        "Рыбы",
		Symbol:      "♓",
		Element:     "Вода",
		Modality:    "мутабельное",
		Theme:       "интуиция, воображение и эмпатия",
		Description: "Солнце в Рыбах подчёркивает мягкость, восприимчивость и богатый внутренний мир.",
	},
}
