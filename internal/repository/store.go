package repository

import (
	"log/slog"
	"slices"
	"sync"
)

type Store struct {
	CookieStore []string
	mu          *sync.Mutex
}

// CookieStore に値が存在するかを確認する
func (s *Store) CheckCookieValue(value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Contains(s.CookieStore, value)
}

func (s *Store) InsertCookieValue(value string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CookieStore = append(s.CookieStore, value)
	slog.Info("add cookie store", "count", len(s.CookieStore))
	return nil
}
