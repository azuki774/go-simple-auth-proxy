package auth

import (
	"azuki774/go-simple-auth-proxy/internal/timeutil"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const testBaseTime = 1721142000

func TestAuthenticater_CheckBasicAuth(t *testing.T) {
	type fields struct {
		AuthStore Store
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		user   string // basic user
		pass   string // basic pass
		want   bool
	}{
		{
			name:   "no basicauth in request",
			fields: fields{AuthStore: &mockStore{}},
			args:   args{r: &http.Request{}},
			want:   false,
		},
		{
			name:   "basicauth OK",
			fields: fields{AuthStore: &mockStore{ReturnGetBasicAuthPassword: "$2a$10$etIpH1oxl4Ky5koV2AzyYe42caqi/tvtme/UTwxA7lHlB2loLDOte"}}, // pass
			args:   args{r: &http.Request{Header: http.Header{}}},
			user:   "user",
			pass:   "pass",
			want:   true,
		},
		{
			name:   "basicauth NG 1",
			fields: fields{AuthStore: &mockStore{ReturnGetBasicAuthPassword: ""}},
			args:   args{r: &http.Request{Header: http.Header{}}},
			user:   "root",
			pass:   "pass",
			want:   false,
		},
		{
			name:   "basicauth NG 2",
			fields: fields{AuthStore: &mockStore{ReturnGetBasicAuthPassword: "$2a$10$etIpH1oxl4Ky5koV2AzyYe42caqi/tvtme/UTwxA7lHlB2loLDOte"}}, // pass
			args:   args{r: &http.Request{Header: http.Header{}}},
			user:   "user",
			pass:   "passWORD",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Authenticater{
				AuthStore: tt.fields.AuthStore,
			}
			if tt.user != "" {
				tt.args.r.SetBasicAuth(tt.user, tt.pass)
			}
			if got := a.CheckBasicAuth(tt.args.r); got != tt.want {
				t.Errorf("Authenticater.CheckBasicAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticater_GenerateCookie(t *testing.T) {
	type fields struct {
		AuthStore      Store
		Issuer         string
		ExpirationTime int64
		HmacSecret     string
	}
	tests := []struct {
		name            string
		fields          fields
		wantCookieValue string // from *http.Cookie
		wantErr         bool
		Nowtime         time.Time
	}{
		{
			name: "ok",
			fields: fields{
				AuthStore:      &mockStore{},
				Issuer:         "testprogram",
				ExpirationTime: 999, // now: testBaseTime = 1721142000 -> 1721142999
				HmacSecret:     "super_sugoi_secret",
			},
			// {"alg":"HS256","typ":"JWT"} -> eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
			// {"exp":1721142999,"iss":"testprogram"} -> eyJleHAiOjE3MjExNDI5OTksImlzcyI6InRlc3Rwcm9ncmFtIn0
			// sign (super_sugoi_secret): eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjExNDI5OTksImlzcyI6InRlc3Rwcm9ncmFtIn0 => MJd9moHsqxrUs3ujOUcwR6AEQNZzbqj8yOudrHfCBpg
			wantCookieValue: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjExNDI5OTksImlzcyI6InRlc3Rwcm9ncmFtIn0.MJd9moHsqxrUs3ujOUcwR6AEQNZzbqj8yOudrHfCBpg",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeutil.NowFunc = func() time.Time { return time.Unix(testBaseTime, 0) }

			a := &Authenticater{
				AuthStore:      tt.fields.AuthStore,
				Issuer:         tt.fields.Issuer,
				ExpirationTime: tt.fields.ExpirationTime,
				HmacSecret:     tt.fields.HmacSecret,
			}
			got, err := a.GenerateCookie()
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticater.GenerateCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// token の中身を比較
			gottoken := got.Value

			if !reflect.DeepEqual(gottoken, tt.wantCookieValue) {
				t.Errorf("Authenticater.GenerateCookie() = %v, want %v", gottoken, tt.wantCookieValue)

			}
		})
	}
}

func TestAuthenticater_IsValidCookie(t *testing.T) {
	type fields struct {
		AuthStore      Store
		Issuer         string
		ExpirationTime int64
		HmacSecret     string
	}
	type args struct {
		r           *http.Request
		tokenString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args // cookie データは後で入れる
		wantOk  bool
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				AuthStore:  &mockStore{},
				Issuer:     "testprogram",
				HmacSecret: "super_sugoi_secret",
			},
			args: args{
				r:           &http.Request{},
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTksImlzcyI6InRlc3Rwcm9ncmFtIn0.JQddrOcvLCTzKfPG3oCqwSe0LLcI-xcoIbrZ-DKbbJ4",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "expired",
			fields: fields{
				AuthStore:  &mockStore{},
				Issuer:     "testprogram",
				HmacSecret: "super_sugoi_secret",
			},
			args: args{
				r:           &http.Request{},
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5LCJpc3MiOiJ0ZXN0cHJvZ3JhbSJ9.5Xx7MYFjl60ASmTChS_SROGt9Y9-4Al6ZjcWHlQGGp8",
			},
			wantOk:  false,
			wantErr: false,
		},
		{
			name: "invalid sign",
			fields: fields{
				AuthStore:  &mockStore{},
				Issuer:     "testprogram",
				HmacSecret: "super_sugoi_secret",
			},
			args: args{
				r:           &http.Request{},
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5LCJpc3MiOiJ0ZXN0cHJvZ3JhbSJ9.AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			},
			wantOk:  false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Authenticater{
				AuthStore:      tt.fields.AuthStore,
				Issuer:         tt.fields.Issuer,
				ExpirationTime: tt.fields.ExpirationTime,
				HmacSecret:     tt.fields.HmacSecret,
			}

			// cookie いれる
			tt.args.r.Header = map[string][]string{
				"Cookie": {fmt.Sprintf("%s=%s", CookieJWTName, tt.args.tokenString)},
			}

			gotOk, err := a.IsValidCookie(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticater.IsValidCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("Authenticater.IsValidCookie() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
