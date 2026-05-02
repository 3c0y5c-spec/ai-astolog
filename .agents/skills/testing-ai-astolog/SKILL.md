---
name: testing-ai-astolog
description: Test the ai-astolog Go Telegram bot scaffold end-to-end. Use when verifying bot startup, healthcheck, Docker runtime, or Telegram command replies.
---

# Testing ai-astolog

## Devin Secrets Needed

- `TELEGRAM_BOT_TOKEN`: Telegram BotFather token for the bot under test. Request it through Devin secure secrets; do not paste the value into shell commands or reports.

## Environment

- Use Go from `/home/ubuntu/.local/go/bin` when available:
  ```bash
  PATH=/home/ubuntu/.local/go/bin:$PATH go test ./...
  PATH=/home/ubuntu/.local/go/bin:$PATH go vet ./...
  ```
- The app can be tested through Docker without needing a deployed environment.

## Standard checks

```bash
PATH=/home/ubuntu/.local/go/bin:$PATH go fmt ./...
PATH=/home/ubuntu/.local/go/bin:$PATH go test ./...
PATH=/home/ubuntu/.local/go/bin:$PATH go vet ./...
docker build -t ai-astolog:e2e .
```

## Runtime healthcheck E2E

Run webhook mode to avoid initializing Telegram polling while verifying the built container serves HTTP:

```bash
docker rm -f ai-astolog-e2e >/dev/null 2>&1 || true
docker run -d --name ai-astolog-e2e \
  -e TELEGRAM_BOT_TOKEN="${TELEGRAM_BOT_TOKEN}" \
  -e BOT_MODE=webhook \
  -e PUBLIC_WEBHOOK_URL=https://example.com \
  -p 18080:8080 \
  ai-astolog:e2e
```

Assert `GET http://127.0.0.1:18080/healthz` returns:

- HTTP status `200`
- `Content-Type` containing `application/json`
- exact body `{"status":"ok"}\n`

## Telegram token verification

Use Telegram Bot API `getMe` with `${TELEGRAM_BOT_TOKEN}` and only log redacted/safe fields. Assert:

- `ok: true`
- numeric bot id is present
- `is_bot: true`
- username is present

## Polling and command reply verification

Run polling mode with the secure token and verify logs include `telegram bot polling started` and do not include `telegram bot init failed`.

A bot token alone cannot fabricate a user-originated `/start` or `/help` message. To verify command replies, ask a human to send `/start` and `/help` to the bot in Telegram and provide a screenshot or text transcript. Expected replies are defined in `internal/telegram/replies.go`.

## Cleanup

Always stop the test container after testing:

```bash
docker rm -f ai-astolog-e2e >/dev/null 2>&1 || true
```
