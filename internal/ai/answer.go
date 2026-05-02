package ai

import (
	"fmt"
	"strings"

	domainastrology "github.com/3c0y5c-spec/ai-astolog/internal/domain/astrology"
	domainprofile "github.com/3c0y5c-spec/ai-astolog/internal/domain/profile"
)

func BuildProfileContext(birthProfile domainprofile.BirthProfile) string {
	birthTime := "не указано"
	if birthProfile.BirthTime != nil {
		birthTime = birthProfile.BirthTime.String()
	}

	sign := domainastrology.SunSignForDate(birthProfile.BirthDate)
	return fmt.Sprintf(
		"Дата рождения: %s\nВремя рождения: %s\nГород рождения: %s\nСолнечный знак: %s %s\nСтихия: %s\nКлючевая тема: %s",
		birthProfile.BirthDate.Format("02.01.2006"),
		birthTime,
		birthProfile.City,
		sign.Name,
		sign.Symbol,
		sign.Element,
		sign.Theme,
	)
}

func BuildFallbackAnswer(birthProfile domainprofile.BirthProfile, question string) string {
	sign := domainastrology.SunSignForDate(birthProfile.BirthDate)
	normalizedQuestion := strings.TrimSpace(question)
	if normalizedQuestion == "" {
		normalizedQuestion = "твой вопрос"
	}

	return fmt.Sprintf(
		"AI-провайдер пока не настроен, поэтому отвечаю в MVP-режиме по солнечному знаку.\n\nВопрос: %s\n\n%s %s подсказывает смотреть на ситуацию через тему: %s. Начни с одного честного шага: сформулируй, чего ты хочешь на самом деле, и выбери действие, которое можно сделать сегодня без давления.\n\nЭто развлекательная астрологическая интерпретация для саморефлексии, а не инструкция к действию.",
		normalizedQuestion,
		sign.Name,
		sign.Symbol,
		sign.Theme,
	)
}
