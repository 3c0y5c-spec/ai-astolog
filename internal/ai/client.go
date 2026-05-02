package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
	defaultModel   = "gpt-4o-mini"
	defaultTimeout = 20 * time.Second
)

var ErrNotConfigured = errors.New("ai provider is not configured")

type Config struct {
	APIKey  string
	BaseURL string
	Model   string
	Timeout time.Duration
}

type QuestionRequest struct {
	ProfileContext string
	Question       string
}

type Client interface {
	AnswerQuestion(ctx context.Context, request QuestionRequest) (string, error)
}

type OpenAIClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewOpenAIClient(config Config) *OpenAIClient {
	timeout := config.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return &OpenAIClient{
		apiKey:  strings.TrimSpace(config.APIKey),
		baseURL: strings.TrimRight(stringWithDefault(config.BaseURL, defaultBaseURL), "/"),
		model:   stringWithDefault(config.Model, defaultModel),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *OpenAIClient) AnswerQuestion(ctx context.Context, request QuestionRequest) (string, error) {
	if c.apiKey == "" {
		return "", ErrNotConfigured
	}

	body, err := json.Marshal(chatCompletionRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: fmt.Sprintf("Профиль пользователя:\n%s\n\nВопрос:\n%s", request.ProfileContext, request.Question)},
		},
		Temperature: 0.7,
	})
	if err != nil {
		return "", err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpRequest.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("ai provider returned %s", response.Status)
	}

	var completion chatCompletionResponse
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return "", err
	}
	if len(completion.Choices) == 0 {
		return "", errors.New("ai provider returned no choices")
	}

	answer := strings.TrimSpace(completion.Choices[0].Message.Content)
	if answer == "" {
		return "", errors.New("ai provider returned empty answer")
	}

	return answer, nil
}

func stringWithDefault(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

const systemPrompt = "Ты AI-астролог в Telegram-боте. Отвечай на русском языке, тепло и конкретно, в развлекательном эзотерическом стиле. Используй только предоставленный профиль и общий астрологический контекст. Не давай медицинских, финансовых, юридических или психологических инструкций. Заверши коротким дисклеймером, что это интерпретация для саморефлексии."

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []chatCompletionChoice `json:"choices"`
}

type chatCompletionChoice struct {
	Message chatMessage `json:"message"`
}
