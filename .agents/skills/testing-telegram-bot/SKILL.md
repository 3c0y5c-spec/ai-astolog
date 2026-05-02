---
name: testing-telegram-bot
description: Test the ai-astolog Telegram bot runtime flows end-to-end. Use when verifying Telegram command handling, onboarding questionnaires, chart summaries, daily forecasts, ask flows, or bot runtime changes.
---

# Testing ai-astolog Telegram Bot

## Devin Secrets Needed

- `TELEGRAM_BOT_TOKEN`: BotFather token for the test Telegram bot. Reference it as `${TELEGRAM_BOT_TOKEN}`; do not print or write the value to files.
- `AI_API_KEY`: Optional OpenAI-compatible key for provider-backed `/ask` testing. If it is absent, `/ask` should use local MVP fallback mode.

## Setup

1. Use the managed Go toolchain from `/home/ubuntu/.local/go/bin`:
   ```bash
   PATH=/home/ubuntu/.local/go/bin:$PATH go test ./...
   PATH=/home/ubuntu/.local/go/bin:$PATH go vet ./...
   ```
2. Build the Docker image from the repo root. Use a task-specific tag if helpful:
   ```bash
   docker build -t ai-astolog:test .
   ```
3. Stop any existing polling containers before starting a new one. Telegram allows only one active `getUpdates` poller per bot token:
   ```bash
   docker rm -f ai-astolog-ask-test ai-astolog-daily-test ai-astolog-chart-test ai-astolog-profile-test ai-astolog-profile-e2e 2>/dev/null || true
   ```
4. Start the bot in polling mode without putting the secret value in the command line:
   ```bash
   docker run --rm --name ai-astolog-ask-test \
     -p 18083:8080 \
     --env TELEGRAM_BOT_TOKEN \
     -e BOT_MODE=polling \
     -e APP_ENV=test \
     -e AI_API_KEY= \
     ai-astolog:test
   ```
5. Verify runtime basics from another shell:
   ```bash
   curl -fsS http://localhost:18083/healthz
   curl -fsS "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/getMe" | python3 -m json.tool
   ```

Expected basics:
- `/healthz` returns `{"status":"ok"}`.
- `getMe` returns `ok: true`, `is_bot: true`, and the expected test bot username.
- Runtime logs include `http server started` and `telegram bot polling started`.

## `/profile` E2E Flow

Use a real Telegram client to send user-originated messages to the bot. The core adversarial sequence is:

1. `/profile`
2. `31.02.1992`
3. `24.03.1992`
4. `25:99`
5. `08:30`
6. `/help`
7. `Москва`
8. `/profile`
9. `/cancel`
10. `24.03.1992`

Expected assertions:
- `/profile` asks for birth date in `ДД.ММ.ГГГГ`.
- `31.02.1992` replies with `Не смог распознать дату` and stays on date.
- `24.03.1992` asks for time in `ЧЧ:ММ`.
- `25:99` replies with `Не смог распознать время` and stays on time.
- `08:30` asks for city.
- `/help` during city input replies with `Введи город рождения текстом`, not the generic help text.
- `Москва` saves the profile and shows exact values: `Дата рождения: 24.03.1992`, `Время рождения: 08:30`, `Город рождения: Москва`.
- `/cancel` replies `Анкета отменена. Отправь /profile, чтобы начать заново.`
- A date sent after `/cancel` should receive the generic help/disclaimer response, not continue the questionnaire.

## `/chart` E2E Flow

Use a fresh bot process/in-memory store when checking missing-profile behavior. The core sequence is:

1. `/chart`
2. `/profile`
3. `24.03.1992`
4. `08:30`
5. `Москва`
6. `/chart`

Expected assertions:
- First `/chart` replies with `Сначала заполни анкету рождения через /profile`.
- `/profile` collects and saves exact values `24.03.1992`, `08:30`, and `Москва`.
- Final `/chart` contains `Натальная карта (MVP):`.
- Final `/chart` includes exact saved fields: `Дата рождения: 24.03.1992`, `Время рождения: 08:30`, `Город рождения: Москва`.
- Final `/chart` maps the date to `Солнечный знак: Овен ♈` and includes `Стихия: Огонь`, `Качество: кардинальное`, plus the Moon/Ascendant disclaimer.

## `/daily` E2E Flow

Use a fresh bot process/in-memory store when checking missing-profile behavior. The core sequence is:

1. `/daily`
2. `/profile`
3. `24.03.1992`
4. `08:30`
5. `Москва`
6. `/daily`

Expected assertions:
- First `/daily` replies with `Сначала заполни анкету рождения через /profile`.
- `/profile` collects and saves exact values `24.03.1992`, `08:30`, and `Москва`.
- Final `/daily` contains `Ежедневный прогноз на <current date in DD.MM.YYYY>:`.
- Final `/daily` maps the date to `Солнечный знак: Овен ♈`.
- Final `/daily` includes `Фокус дня:`, `Совет:`, `Вопрос для себя:`, and the entertainment disclaimer.

## `/ask` E2E Flow

Use a fresh bot process/in-memory store when checking missing-profile behavior. If `AI_API_KEY` is absent or intentionally empty, test fallback mode. The core sequence is:

1. `/ask`
2. `/profile`
3. `24.03.1992`
4. `08:30`
5. `Москва`
6. `/ask`
7. `/help`
8. `Что важно сегодня?`

Expected assertions:
- First `/ask` replies with `Сначала заполни анкету рождения через /profile`.
- `/profile` collects and saves exact values `24.03.1992`, `08:30`, and `Москва`.
- Second `/ask` prompts with `Задай вопрос AI-астрологу одним сообщением` and mentions `/cancel`.
- `/help` while waiting for the question replies with `Задай вопрос текстом`, not generic help.
- The final fallback answer includes `AI-провайдер пока не настроен`.
- The final fallback answer includes `Вопрос: Что важно сегодня?`.
- The final fallback answer maps the date to `Овен ♈` and includes `инициатива, смелость и быстрый старт`.
- The final fallback answer includes the self-reflection disclaimer.

## Notes

- If logs show `terminated by other getUpdates request`, another bot instance is polling with the same token. Stop the older container/process before retesting.
- For Telegram E2E evidence, ask the user for a screenshot of the real chat unless you have an authenticated Telegram client available in the VM.
