package auth

import (
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	GetBasicAuthPassword(user string) string
}

type Authenticater struct {
	AuthStore Store
}

func (a *Authenticater) GenerateCookie() (*http.Cookie, error) {
	// TODO
	return nil, nil
}

func (a *Authenticater) IsValidCookie(r *http.Request) (ok bool, err error) {
	// TODO
	return false, nil
}

func (a *Authenticater) CheckBasicAuth(r *http.Request) bool {
	// 認証情報取得
	reqUser, reqPass, ok := r.BasicAuth()
	if !ok {
		return false
	}
	hashPass := a.AuthStore.GetBasicAuthPassword(reqUser) // 正しいパスワードのハッシュを取得
	if err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(reqPass)); err != nil {
		slog.Info("basic auth mismatched", "user", reqUser)
		return false
	}

	return true
}
