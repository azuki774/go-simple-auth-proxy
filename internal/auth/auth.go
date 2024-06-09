package auth

import (
	"log/slog"
	"net/http"
)

const (
	// TODO
	basicuser = "user"
	basicpass = "pass"
)

type Store interface {
	CheckCookieValue(value string) bool
	InsertCookieValue(value string) (err error)
}

type Authenticater struct {
	AuthStore Store
}

// 新しい cookie を生成
func (a *Authenticater) GenerateCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:  "token",
		Value: "example_token_value", // TODO
	}
	return cookie
}

// cookie が正当かどうか確認
func (a *Authenticater) IsValidCookie(r *http.Request) (ok bool, err error) {
	token, err := r.Cookie("token")
	if err != nil {
		// unknown error: http: named cookie not present
		// token の key がない場合もここに落ちるので、この場合は ok = false とする
		slog.Debug("undetected cookie: token")
		return false, nil
	}
	v := token.Value
	ok = a.AuthStore.CheckCookieValue(v)

	return ok, nil
}

func (a *Authenticater) CheckBasicAuth(r *http.Request) bool {
	// 認証情報取得
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		// 存在しなければ false
		return false
	}
	return clientID == basicuser && clientSecret == basicpass
}
