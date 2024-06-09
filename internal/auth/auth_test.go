package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
)

func TestCookieManager_IsValidCookie(t *testing.T) {
	type fields struct {
		authStore Store
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				authStore: &mockStore{},
			},
			args: args{
				r: &http.Request{},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "ng(not same value)",
			fields: fields{
				authStore: &mockStore{
					CheckCookieValueErr: true,
				},
			},
			args: args{
				r: &http.Request{},
			},
			wantOk:  false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Authenticater{
				AuthStore: tt.fields.authStore,
			}

			// cookie いれる
			tt.args.r.Header = map[string][]string{"Cookie": []string{"token"}}

			gotOk, err := a.IsValidCookie(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("CookieManager.IsValidCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("CookieManager.IsValidCookie() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

// 乱数生成なので乱数で作って base64 -d するだけ
func TestAuthenticater_GenerateCookie(t *testing.T) {
	var Authenticater Authenticater

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			c := Authenticater.GenerateCookie()
			dec, err := base64.StdEncoding.DecodeString(c.Value)
			if err != nil {
				t.Errorf("TestAuthenticater_GenerateCookie = %v", err)
			}
			fmt.Printf("%s <-> %s\n", c.Value, dec)
		})
	}
}
