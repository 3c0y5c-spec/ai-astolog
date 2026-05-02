package telegram

import (
	"context"
	"strings"
	"sync"

	domainastrology "github.com/3c0y5c-spec/ai-astolog/internal/domain/astrology"
	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

type compatibilityManager struct {
	mu       sync.Mutex
	sessions map[int64]struct{}
}

func newCompatibilityManager() *compatibilityManager {
	return &compatibilityManager{
		sessions: make(map[int64]struct{}),
	}
}

func (m *compatibilityManager) start(userID int64) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[userID] = struct{}{}
	return "Введи дату рождения партнёра в формате ДД.ММ.ГГГГ, например 02.10.1990.\n\nЧтобы отменить расчёт совместимости, отправь /cancel."
}

func (m *compatibilityManager) cancel(userID int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[userID]; !ok {
		return false
	}

	delete(m.sessions, userID)
	return true
}

func (m *compatibilityManager) active(userID int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.sessions[userID]
	return ok
}

func (m *compatibilityManager) clear(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, userID)
}

func (m *compatibilityManager) handle(_ context.Context, userID int64, birthProfile domainprofile.BirthProfile, text string) (string, bool) {
	m.mu.Lock()
	_, ok := m.sessions[userID]
	if !ok {
		m.mu.Unlock()
		return "", false
	}

	partnerBirthDate, valid := parseBirthDate(strings.TrimSpace(text))
	if !valid {
		m.mu.Unlock()
		return "Не смог распознать дату партнёра. Введи дату в формате ДД.ММ.ГГГГ, например 02.10.1990.", true
	}

	delete(m.sessions, userID)
	m.mu.Unlock()

	return domainastrology.BuildCompatibilityReport(birthProfile, partnerBirthDate), true
}
