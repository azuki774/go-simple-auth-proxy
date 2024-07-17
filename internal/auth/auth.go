package auth

import (
	"azuki774/go-simple-auth-proxy/internal/timeutil"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt"
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
	fmt.Println(tokenString) // FOR TEST
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return a.HmacSecret, nil
	})
	if err != nil {
		// TODO: 今のままだとここに落ちてしまう
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["exp"], claims["iss"])
		// TODO: Issuer が正しいか
		// TODO: 有効期限が切れていないか
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
