package telegram

const StartText = "Привет! Я AI-астролог. Помогу собрать данные рождения, построить базовую натальную карту, прогноз, совместимость и ответить на вопросы в астрологическом стиле.\n\nДля MVP доступны команды:\n/help — что умеет бот\n/profile — заполнить анкету рождения\n/cancel — отменить анкету, вопрос или совместимость\n/chart — базовая натальная карта\n/daily — ежедневный прогноз\n/ask — вопрос AI-астрологу\n/compatibility — совместимость по датам рождения"

const HelpText = "AI-астролог — развлекательный эзотерический сервис, а не медицинская, финансовая или юридическая консультация.\n\nДоступно сейчас:\n• анкета рождения через /profile\n• базовая натальная карта через /chart\n• ежедневный прогноз через /daily\n• свободный вопрос AI-астрологу через /ask\n• совместимость через /compatibility"

const ProfileStartText = "Начнём анкету рождения.\n\nВведи дату рождения в формате ДД.ММ.ГГГГ, например 24.03.1992.\n\nЧтобы отменить анкету, отправь /cancel."

const ComingSoonText = "Эта функция появится в следующем этапе MVP."

func ReplyForCommand(command string) string {
	switch command {
	case "/start":
		return StartText
	case "/help":
		return HelpText
	case "/profile":
		return ProfileStartText
	case "/ask":
		return ComingSoonText
	default:
		return HelpText
	}
}
