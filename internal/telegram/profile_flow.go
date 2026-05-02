package telegram

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

const (
	profileStepBirthDate = iota + 1
	profileStepBirthTime
	profileStepCity
)

type profileDraft struct {
	step      int
	birthDate time.Time
	birthTime *domainprofile.CivilTime
}

type profileManager struct {
	mu       sync.Mutex
	store    domainprofile.Store
	sessions map[int64]profileDraft
	now      func() time.Time
}

func newProfileManager(store domainprofile.Store) *profileManager {
	return &profileManager{
		store:    store,
		sessions: make(map[int64]profileDraft),
		now:      time.Now,
	}
}

func (m *profileManager) start(userID int64) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[userID] = profileDraft{step: profileStepBirthDate}
	return ProfileStartText
}

func (m *profileManager) cancel(userID int64) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[userID]; !ok {
		return "Сейчас нет активной анкеты. Отправь /profile, чтобы начать."
	}

	delete(m.sessions, userID)
	return "Анкета отменена. Отправь /profile, чтобы начать заново."
}

func (m *profileManager) get(ctx context.Context, userID int64) (domainprofile.BirthProfile, bool, error) {
	return m.store.Get(ctx, userID)
}

func (m *profileManager) handle(ctx context.Context, userID int64, text string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	draft, ok := m.sessions[userID]
	if !ok {
		return "", false
	}

	text = strings.TrimSpace(text)
	switch draft.step {
	case profileStepBirthDate:
		birthDate, ok := parseBirthDate(text)
		if !ok {
			return "Не смог распознать дату. Введи дату рождения в формате ДД.ММ.ГГГГ, например 24.03.1992.", true
		}
		draft.birthDate = birthDate
		draft.step = profileStepBirthTime
		m.sessions[userID] = draft
		return "Введи время рождения в формате ЧЧ:ММ, например 08:30. Если точного времени нет, напиши «нет».", true
	case profileStepBirthTime:
		birthTime, ok := parseBirthTime(text)
		if !ok {
			return "Не смог распознать время. Введи в формате ЧЧ:ММ или напиши «нет».", true
		}
		draft.birthTime = birthTime
		draft.step = profileStepCity
		m.sessions[userID] = draft
		return "Введи город рождения, например Москва.", true
	case profileStepCity:
		city := strings.TrimSpace(text)
		if len([]rune(city)) < 2 || strings.HasPrefix(city, "/") {
			return "Введи город рождения текстом, например Москва.", true
		}

		now := m.now().UTC()
		birthProfile := domainprofile.BirthProfile{
			UserID:    userID,
			BirthDate: draft.birthDate,
			BirthTime: draft.birthTime,
			City:      city,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := m.store.Save(ctx, birthProfile); err != nil {
			return "Не смог сохранить профиль. Попробуй /profile ещё раз.", true
		}

		delete(m.sessions, userID)
		return formatProfileSaved(birthProfile), true
	default:
		delete(m.sessions, userID)
		return "Анкета сброшена. Отправь /profile, чтобы начать заново.", true
	}
}

func parseBirthDate(text string) (time.Time, bool) {
	for _, layout := range []string{"02.01.2006", "2006-01-02"} {
		value, err := time.ParseInLocation(layout, strings.TrimSpace(text), time.UTC)
		if err == nil {
			return value, true
		}
	}
	return time.Time{}, false
}

func parseBirthTime(text string) (*domainprofile.CivilTime, bool) {
	normalized := strings.ToLower(strings.TrimSpace(text))
	switch normalized {
	case "нет", "не знаю", "unknown", "skip", "-":
		return nil, true
	}

	value, err := time.Parse("15:04", normalized)
	if err != nil {
		return nil, false
	}

	return &domainprofile.CivilTime{Hour: value.Hour(), Minute: value.Minute()}, true
}

func formatProfileSaved(birthProfile domainprofile.BirthProfile) string {
	birthTime := "не указано"
	if birthProfile.BirthTime != nil {
		birthTime = birthProfile.BirthTime.String()
	}

	return fmt.Sprintf(
		"Профиль сохранён:\nДата рождения: %s\nВремя рождения: %s\nГород рождения: %s\n\nДальше можно открыть /chart (скоро).",
		birthProfile.BirthDate.Format("02.01.2006"),
		birthTime,
		birthProfile.City,
	)
}
