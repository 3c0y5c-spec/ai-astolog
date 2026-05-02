# ai-astolog

AI Astolog is a Telegram bot scaffold for an AI-powered astrology assistant.

## Planned MVP

- Telegram onboarding for birth date, birth time, and city.
- Basic natal chart summary.
- Daily personalized forecast.
- Compatibility flow.
- Free-form AI astrologer questions.

The service is intended for entertainment and esoteric self-reflection. It is not medical, legal, financial, or psychological advice.

## Tech stack

- Go 1.26.2
- [`github.com/go-telegram/bot`](https://github.com/go-telegram/bot) for Telegram Bot API integration
- Standard `net/http` health endpoint
- PostgreSQL, migrations, and AI provider integration are planned for later PRs

## Project layout

```text
cmd/bot/                    application entrypoint
internal/config/            environment configuration
internal/domain/astrology/  MVP astrology summaries
internal/domain/profile/    birth profile model and in-memory storage
internal/httpserver/        healthcheck HTTP server
internal/telegram/          Telegram handlers and reply text
```

## Configuration

Copy `.env.example` to `.env` for local development and fill in real values.

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `TELEGRAM_BOT_TOKEN` | yes | - | Bot token from BotFather |
| `APP_ENV` | no | `development` | Runtime environment label |
| `HTTP_ADDR` | no | `:8080` | Healthcheck HTTP listen address |
| `BOT_MODE` | no | `polling` | `polling` or `webhook` |
| `PUBLIC_WEBHOOK_URL` | webhook only | - | Public webhook base URL |
| `WEBHOOK_SECRET` | no | - | Secret for Telegram webhook validation |
| `SHUTDOWN_TIMEOUT_SECONDS` | no | `10` | Graceful shutdown timeout |

## Local development

```bash
cp .env.example .env
export TELEGRAM_BOT_TOKEN="<token-from-botfather>"
go run ./cmd/bot
```

Healthcheck:

```bash
curl http://localhost:8080/healthz
```

## Bot commands

- `/start` — greeting and command list
- `/help` — disclaimer and planned features
- `/profile` — birth profile onboarding:
  1. birth date in `ДД.ММ.ГГГГ` or `YYYY-MM-DD`
  2. birth time in `ЧЧ:ММ`, or `нет` if unknown
  3. birth city text
- `/cancel` — cancel the active profile onboarding flow
- `/chart` — basic natal chart MVP summary from the saved birth profile
- `/daily` — daily MVP forecast from the saved birth profile
- `/ask` — planned next feature

## Checks

```bash
gofmt -w ./cmd ./internal
go test ./...
go vet ./...
```

## Docker

```bash
docker build -t ai-astolog .
docker run --rm -p 8080:8080 --env-file .env ai-astolog
```