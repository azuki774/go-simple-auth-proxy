package repository

import (
	"log/slog"
	"slices"
	"sync"
)

type Store struct {
	CookieStore []string
	Mu          *sync.Mutex
}

// CookieStore に値が存在するかを確認する
func (s *Store) CheckCookieValue(value string) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	return slices.Contains(s.CookieStore, value)
}

func (s *Store) InsertCookieValue(value string) (err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.CookieStore = append(s.CookieStore, value)
	slog.Info("add cookie store", "count", len(s.CookieStore))
	return nil
}
