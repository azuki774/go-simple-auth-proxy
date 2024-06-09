package server

import (
	"errors"
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
				Authenticater: &mockAuthenticater{cookieok: true},
			},
			args: args{
				w: &mockResponseWriter{},
				r: &http.Request{},
			},
			wantResultCode: ProxyResultCookieOK,
		},
		{
			name: "cookie NG -> BasicAuth OK",
			fields: fields{
				Client:        &mockClient{},
				Authenticater: &mockAuthenticater{basicok: true},
			},
			args: args{
				w: &mockResponseWriter{},
				r: &http.Request{},
			},
			wantResultCode: ProxyResultBasicAuthOK,
		},
		{
			name: "cookie NG -> BasicAuth NG",
			fields: fields{
				Client:        &mockClient{},
				Authenticater: &mockAuthenticater{},
			},
			args: args{
				w: &mockResponseWriter{},
				r: &http.Request{},
			},
			wantResultCode: ProxyResultBasicAuthNG,
		},
		{
			name: "Fetch NG",
			fields: fields{
				Client:        &mockClient{err: errors.New("nanka error")},
				Authenticater: &mockAuthenticater{basicok: true},
			},
			args: args{
				w: &mockResponseWriter{},
				r: &http.Request{},
			},
			wantResultCode: ProxyResultFetchNG,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				ListenPort:    tt.fields.ListenPort,
				Client:        tt.fields.Client,
				Authenticater: tt.fields.Authenticater,
			}
			if gotResultCode := s.proxyMain(tt.args.w, tt.args.r); gotResultCode != tt.wantResultCode {
				t.Errorf("Server.proxyMain() = %v, want %v", gotResultCode, tt.wantResultCode)
			}
		})
	}
}
