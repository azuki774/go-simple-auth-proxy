package auth

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
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
func (a *Authenticater) GenerateCookie() (*http.Cookie, error) {
	// base64 ( uuid v4 : dt.Unix() )
	rowv := fmt.Sprintf("%s:%d", uuid.New().String(), time.Now().Unix())
	v := base64.StdEncoding.EncodeToString([]byte(rowv))
	cookie := &http.Cookie{
		Name:  "token",
		Value: v,
	}
	slog.Info("generate cookie", "value", v)
	err := a.AuthStore.InsertCookieValue(v)
	if err != nil {
		slog.Error("failed to store new cookie", "err", err)
		return nil, err
	}
	return cookie, nil
}

// cookie が正当かどうか確認
func (a *Authenticater) IsValidCookie(r *http.Request) (ok bool, err error) {
	token, err := r.Cookie("token")
	if err != nil {
		// unknown error: http: named cookie not present
		// token の key がない場合もここに落ちるので、この場合は ok = false とする
		slog.Info("undetected cookie: token")
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
