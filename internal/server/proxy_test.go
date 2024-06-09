package server

import (
	"net/http"
	"testing"
)

func TestServer_proxyMain(t *testing.T) {
	type fields struct {
		ListenPort    string
		ProxyAddr     string
		Client        Client
		Authenticater Authenticater
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantResultCode ProxyResultCode
	}{
		{
			name: "cookie OK",
			fields: fields{
				Client:        &mockClient{},
				Authenticater: &mockAuthenticater{},
			},
			args: args{
				w: &mockResponseWriter{},
				r: &http.Request{},
			},
			wantResultCode: ProxyResultCookieOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				ListenPort:    tt.fields.ListenPort,
				ProxyAddr:     tt.fields.ProxyAddr,
				Client:        tt.fields.Client,
				Authenticater: tt.fields.Authenticater,
			}
			if gotResultCode := s.proxyMain(tt.args.w, tt.args.r); gotResultCode != tt.wantResultCode {
				t.Errorf("Server.proxyMain() = %v, want %v", gotResultCode, tt.wantResultCode)
			}
		})
	}
}
