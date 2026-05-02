package profile

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type BirthProfile struct {
	UserID    int64
	BirthDate time.Time
	BirthTime *CivilTime
	City      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CivilTime struct {
	Hour   int
	Minute int
}

func (t CivilTime) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
}

type Store interface {
	Save(ctx context.Context, birthProfile BirthProfile) error
	Get(ctx context.Context, userID int64) (BirthProfile, bool, error)
}

type MemoryStore struct {
	mu       sync.RWMutex
	profiles map[int64]BirthProfile
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		profiles: make(map[int64]BirthProfile),
	}
}

func (s *MemoryStore) Save(_ context.Context, birthProfile BirthProfile) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.profiles[birthProfile.UserID] = birthProfile
	return nil
}

func (s *MemoryStore) Get(_ context.Context, userID int64) (BirthProfile, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	birthProfile, ok := s.profiles[userID]
	return birthProfile, ok, nil
}
