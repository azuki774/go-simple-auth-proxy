package repository

import (
	"sync"
)

type StoreInMemory struct {
	BasicAuthStore map[string]string
	Mu             *sync.Mutex
}

func (s *StoreInMemory) GetBasicAuthPassword(user string) string {
	// if 'user' is not found, return ""
	return s.BasicAuthStore[user]
}
