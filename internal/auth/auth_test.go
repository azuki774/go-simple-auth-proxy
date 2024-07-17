package auth

import (
	"azuki774/go-simple-auth-proxy/internal/timeutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

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
				ExpirationTime: 999, // now: 1721142000 -> 1721142999
				HmacSecret:     "super_sugoi_secret",
			},
			// {"alg":"HS256","typ":"JWT"} -> eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
			// {"exp":1721142999,"iss":"testprogram"} -> eyJpc3MiOiJ0ZXN0cHJvZ3JhbSIsImV4cCI6MTcyMTE0Mjk5OX0
			// sign (super_sugoi_secret): eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0ZXN0cHJvZ3JhbSIsImV4cCI6MTcyMTE0Mjk5OX0 => x9E6MEisgT3eTTZXMCaK0BGVdeVPuuN1ZPsP7w89eiE
			wantCookieValue: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0ZXN0cHJvZ3JhbSIsImV4cCI6MTcyMTE0Mjk5OX0.MJd9moHsqxrUs3ujOUcwR6AEQNZzbqj8yOudrHfCBpg",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeutil.NowFunc = func() time.Time { return time.Unix(1721142000, 0) }
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
