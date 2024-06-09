package auth

import (
	"fmt"
	"net/http"
)

type AuthStore interface {
	CheckCookieValue(value string) bool
}

type CookieManager struct {
	authStore AuthStore
}

// 新しい cookie を生成
func GenerateCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:  "token",
		Value: "example_token_value",
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
