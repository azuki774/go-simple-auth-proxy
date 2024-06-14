package repository

import (
	"log/slog"
	"slices"
	"sync"
)

type StoreInMemory struct {
	CookieStore      []string
	BasicAuthStore   map[string]string
	Mu               *sync.Mutex
	MaxAuthStoreSize int // これ以上保存すると古いものから削除する
}

// CookieStore に値が存在するかを確認する
func (s *StoreInMemory) CheckCookieValue(value string) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	return slices.Contains(s.CookieStore, value)
}

func (s *StoreInMemory) InsertCookieValue(value string) (err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.CookieStore = append(s.CookieStore, value)

	if len(s.CookieStore) > s.MaxAuthStoreSize { // 多すぎたら削除する
		s.CookieStore = s.CookieStore[1:]
	}

	slog.Info("add cookie store", "count", len(s.CookieStore))
	return nil
}

func (s *StoreInMemory) GetBasicAuthPassword(user string) string {
	// if 'user' is not found, return ""
	return s.BasicAuthStore[user]
}
