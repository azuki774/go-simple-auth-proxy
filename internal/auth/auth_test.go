package auth

import (
	"net/http"
	"testing"
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
