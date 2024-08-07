package auth

import (
	"azuki774/go-simple-auth-proxy/internal/timeutil"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const CookieJWTName = "jwt"

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
		Name:  CookieJWTName,
		Value: tokenString,
	}

	return cookie, nil
}

func (a *Authenticater) IsValidCookie(r *http.Request) (ok bool, err error) {
	tokenCookie, err := r.Cookie(CookieJWTName)
	if err != nil {
		// unknown error: http: named cookie not present
		// token の key がない場合もここに落ちるので、この場合は ok = false とする
		return false, nil
	}

	tokenString := tokenCookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(a.HmacSecret), nil
	})
	if err != nil {
		// token expired も含む
		if errors.Is(err, jwt.ErrTokenExpired) {
			slog.Info("token expired", "jwt", maskedJwt(tokenString))
			return false, nil
		}
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["iss"] != a.Issuer {
			slog.Info("issuer mismatched", "jwt", maskedJwt(tokenString))
			return false, nil
		}
	} else {
		return false, err
	}

	return true, nil
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

func maskedJwt(tokenString string) string {
	splitsToken := strings.Fields(tokenString) // 'AAA.BBB.CCC' -> ['AAA','BBB','CCC']
	if len(splitsToken) != 3 {
		return tokenString
	}
	splitsToken[2] = "***"
	return fmt.Sprintf("%s.%s.%s", splitsToken[0], splitsToken[1], splitsToken[2])
}
