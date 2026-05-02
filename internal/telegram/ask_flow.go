package telegram

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/3c0y5c-spec/ai-astolog/internal/ai"
	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

type askManager struct {
	mu       sync.Mutex
	sessions map[int64]struct{}
	client   ai.Client
}

func newAskManager(client ai.Client) *askManager {
	return &askManager{
		sessions: make(map[int64]struct{}),
		client:   client,
	}
}

func (m *askManager) start(userID int64) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[userID] = struct{}{}
	return "Задай вопрос AI-астрологу одним сообщением. Например: «На что обратить внимание в отношениях?»\n\nЧтобы отменить вопрос, отправь /cancel."
}

func (m *askManager) cancel(userID int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[userID]; !ok {
		return false
	}

	delete(m.sessions, userID)
	return true
}

func (m *askManager) active(userID int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.sessions[userID]
	return ok
}

func (m *askManager) handle(ctx context.Context, userID int64, birthProfile domainprofile.BirthProfile, text string) (string, bool) {
	m.mu.Lock()
	_, ok := m.sessions[userID]
	if !ok {
		m.mu.Unlock()
		return "", false
	}

	question := strings.TrimSpace(text)
	if question == "" || strings.HasPrefix(question, "/") {
		m.mu.Unlock()
		return "Задай вопрос текстом, например: «Как лучше распределить силы сегодня?»", true
	}

	delete(m.sessions, userID)
	m.mu.Unlock()

	if m.client == nil {
		return ai.BuildFallbackAnswer(birthProfile, question), true
	}

	answer, err := m.client.AnswerQuestion(ctx, ai.QuestionRequest{
		ProfileContext: ai.BuildProfileContext(birthProfile),
		Question:       question,
	})
	if err != nil {
		if errors.Is(err, ai.ErrNotConfigured) {
			return ai.BuildFallbackAnswer(birthProfile, question), true
		}
		return "Не смог получить ответ AI-астролога. Попробуй /ask ещё раз чуть позже.", true
	}

	return answer, true
}
