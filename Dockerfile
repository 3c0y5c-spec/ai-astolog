FROM golang:1.26.2-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/ai-astolog ./cmd/bot

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=build /out/ai-astolog /app/ai-astolog

USER app
EXPOSE 8080

ENTRYPOINT ["/app/ai-astolog"]
