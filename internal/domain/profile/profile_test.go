package profile

import (
	"context"
	"testing"
	"time"
)

func TestMemoryStoreSavesAndLoadsProfile(t *testing.T) {
	store := NewMemoryStore()
	birthTime := CivilTime{Hour: 8, Minute: 30}
	want := BirthProfile{
		UserID:    42,
		BirthDate: time.Date(1992, time.March, 24, 0, 0, 0, 0, time.UTC),
		BirthTime: &birthTime,
		City:      "Москва",
		CreatedAt: time.Date(2026, time.May, 2, 5, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, time.May, 2, 5, 0, 0, 0, time.UTC),
	}

	if err := store.Save(context.Background(), want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, ok, err := store.Get(context.Background(), 42)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !ok {
		t.Fatal("Get() ok = false, want true")
	}
	if got.UserID != want.UserID || got.City != want.City || !got.BirthDate.Equal(want.BirthDate) || got.BirthTime.String() != want.BirthTime.String() {
		t.Fatalf("Get() = %+v, want %+v", got, want)
	}
}

func TestMemoryStoreMissingProfile(t *testing.T) {
	store := NewMemoryStore()

	_, ok, err := store.Get(context.Background(), 42)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if ok {
		t.Fatal("Get() ok = true, want false")
	}
}
