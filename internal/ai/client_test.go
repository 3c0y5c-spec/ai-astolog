package ai

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

func TestOpenAIClientReturnsNotConfiguredWithoutAPIKey(t *testing.T) {
	client := NewOpenAIClient(Config{})

	_, err := client.AnswerQuestion(context.Background(), QuestionRequest{Question: "Что делать?"})
	if !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("AnswerQuestion() error = %v, want ErrNotConfigured", err)
	}
}

func TestOpenAIClientSendsChatCompletionRequest(t *testing.T) {
	var gotAuthorization string
	var gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthorization = r.Header.Get("Authorization")
		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		gotBody = string(body)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"Ответ AI"}}]}`))
	}))
	defer server.Close()

	client := NewOpenAIClient(Config{
		APIKey:  "secret",
		BaseURL: server.URL,
		Model:   "test-model",
		Timeout: time.Second,
	})

	got, err := client.AnswerQuestion(context.Background(), QuestionRequest{
		ProfileContext: "Солнечный знак: Овен ♈",
		Question:       "Что важно сегодня?",
	})
	if err != nil {
		t.Fatalf("AnswerQuestion() error = %v", err)
	}
	if got != "Ответ AI" {
		t.Fatalf("AnswerQuestion() = %q, want %q", got, "Ответ AI")
	}
	if gotAuthorization != "Bearer secret" {
		t.Fatalf("Authorization = %q, want bearer token", gotAuthorization)
	}
	for _, want := range []string{"test-model", "Солнечный знак: Овен ♈", "Что важно сегодня?"} {
		if !strings.Contains(gotBody, want) {
			t.Fatalf("request body = %q, want substring %q", gotBody, want)
		}
	}
}

func TestBuildFallbackAnswerUsesSunSign(t *testing.T) {
	got := BuildFallbackAnswer(domainprofile.BirthProfile{
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		City:      "Москва",
	}, "Что важно сегодня?")

	for _, want := range []string{"AI-провайдер пока не настроен", "Вопрос: Что важно сегодня?", "Овен ♈", "инициатива, смелость и быстрый старт"} {
		if !strings.Contains(got, want) {
			t.Fatalf("BuildFallbackAnswer() = %q, want substring %q", got, want)
		}
	}
}

func TestBuildProfileContextIncludesProfileAndSign(t *testing.T) {
	birthTime := domainprofile.CivilTime{Hour: 8, Minute: 30}
	got := BuildProfileContext(domainprofile.BirthProfile{
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		BirthTime: &birthTime,
		City:      "Москва",
	})

	for _, want := range []string{"Дата рождения: 24.03.1992", "Время рождения: 08:30", "Город рождения: Москва", "Солнечный знак: Овен ♈"} {
		if !strings.Contains(got, want) {
			t.Fatalf("BuildProfileContext() = %q, want substring %q", got, want)
		}
	}
}
