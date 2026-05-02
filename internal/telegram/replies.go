package telegram

const StartText = "Привет! Я AI-астролог. Помогу собрать данные рождения, построить базовую натальную карту и ответить на вопросы в астрологическом стиле.\n\nДля MVP доступны команды:\n/help — что умеет бот\n/profile — анкета рождения (скоро)\n/chart — натальная карта (скоро)\n/daily — ежедневный прогноз (скоро)\n/ask — вопрос AI-астрологу (скоро)"

const HelpText = "AI-астролог — развлекательный эзотерический сервис, а не медицинская, финансовая или юридическая консультация.\n\nПланируемые функции:\n• анкета рождения\n• натальная карта\n• ежедневный прогноз\n• совместимость\n• свободный вопрос AI-астрологу"

const ComingSoonText = "Эта функция появится в следующем этапе MVP."

func ReplyForCommand(command string) string {
	switch command {
	case "/start":
		return StartText
	case "/help":
		return HelpText
	case "/profile", "/chart", "/daily", "/ask":
		return ComingSoonText
	default:
		return HelpText
	}
}
