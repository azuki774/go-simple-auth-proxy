package auth

import (
	"fmt"
	"net/http"
)

const (
	// TODO
	basicuser = "user"
	basicpass = "pass"
)

type Store interface {
	CheckCookieValue(value string) bool
}

type CookieManager struct {
	authStore Store
}

// 新しい cookie を生成
func GenerateCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:  "token",
		Value: "example_token_value", // TODO
	}
	return cookie
}

// cookie が正当かどうか確認
func (c *CookieManager) IsValidCookie(r *http.Request) (ok bool, err error) {
	token, err := r.Cookie("token")
	if err != nil {
		return false, fmt.Errorf("unknown error: %w", err)
	}
	v := token.Value
	ok = c.authStore.CheckCookieValue(v)

	return ok, nil
}

func CheckBasicAuth(r *http.Request) bool {
	// 認証情報取得
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		// 存在しなければ false
		return false
	}
	return clientID == basicuser && clientSecret == basicpass
}
