package auth

import (
	"azuki774/go-simple-auth-proxy/internal/timeutil"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	GetBasicAuthPassword(user string) string
}

type Authenticater struct {
	AuthStore      Store
	Issuer         string // use JWT Payload
	ExpirationTime int64
	HmacSecret     string
}

func (a *Authenticater) GenerateCookie() (*http.Cookie, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": timeutil.NowFunc().Unix() + a.ExpirationTime,
		"iss": a.Issuer,
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(a.HmacSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT access token: %w", err)
	}
	cookie := &http.Cookie{
		Name:  "jwt",
		Value: tokenString,
	}

	return cookie, nil
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
