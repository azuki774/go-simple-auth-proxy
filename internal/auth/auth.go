package auth

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

type Store interface {
	CheckCookieValue(value string) bool
	InsertCookieValue(value string) (err error)
	GetBasicAuthPassword(user string) string
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
		return false, nil
	}
	v := token.Value
	ok = a.AuthStore.CheckCookieValue(v)

	return ok, nil
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
