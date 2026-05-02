package astrology

import (
	"fmt"
	"time"

	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

type DailyForecast struct {
	Date       time.Time
	Sign       SunSign
	Focus      string
	Advice     string
	Reflection string
}

func BuildDailyForecast(birthProfile domainprofile.BirthProfile, date time.Time) string {
	forecast := DailyForecastForProfile(birthProfile, date)

	return fmt.Sprintf(
		"Ежедневный прогноз на %s:\nСолнечный знак: %s %s\nФокус дня: %s\nСовет: %s\nВопрос для себя: %s\n\nЭто развлекательная астрологическая интерпретация, а не инструкция к действию.",
		forecast.Date.Format("02.01.2006"),
		forecast.Sign.Name,
		forecast.Sign.Symbol,
		forecast.Focus,
		forecast.Advice,
		forecast.Reflection,
	)
}

func DailyForecastForProfile(birthProfile domainprofile.BirthProfile, date time.Time) DailyForecast {
	sign := SunSignForDate(birthProfile.BirthDate)
	index := (date.YearDay() + int(birthProfile.BirthDate.Month()) + birthProfile.BirthDate.Day()) % len(dailyThemes)
	theme := dailyThemes[index]

	return DailyForecast{
		Date:       date,
		Sign:       sign,
		Focus:      fmt.Sprintf("%s через тему «%s»", sign.Theme, theme.Focus),
		Advice:     theme.Advice,
		Reflection: theme.Reflection,
	}
}

type dailyTheme struct {
	Focus      string
	Advice     string
	Reflection string
}

var dailyThemes = []dailyTheme{
	{
		Focus:      "внимание к телу и темпу",
		Advice:     "выбери одно важное дело и доведи его до спокойного результата.",
		Reflection: "где сегодня стоит замедлиться, чтобы действовать точнее?",
	},
	{
		Focus:      "разговоры и ясные договорённости",
		Advice:     "формулируй просьбы прямо и проверяй, что тебя поняли одинаково.",
		Reflection: "какой разговор лучше не откладывать?",
	},
	{
		Focus:      "дом, опора и эмоциональная честность",
		Advice:     "оставь место для восстановления и не обещай больше, чем можешь дать.",
		Reflection: "что сегодня делает тебя устойчивее?",
	},
	{
		Focus:      "смелый шаг и личная инициатива",
		Advice:     "начни с маленького действия, которое давно просится наружу.",
		Reflection: "где пора выбрать себя без лишней борьбы?",
	},
	{
		Focus:      "порядок, детали и полезные привычки",
		Advice:     "разбери один узкий участок хаоса вместо попытки исправить всё сразу.",
		Reflection: "какая простая настройка сэкономит тебе силы завтра?",
	},
	{
		Focus:      "творчество и живой интерес",
		Advice:     "дай себе право сделать задачу чуть красивее, легче или интереснее.",
		Reflection: "что сегодня вернёт ощущение игры?",
	},
}
