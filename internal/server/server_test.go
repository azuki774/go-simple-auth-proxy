package server

import (
	"errors"
	"testing"
)

func TestServer_CheckReadiness(t *testing.T) {
	type fields struct {
		ListenPort    string
		Client        Client
		Authenticater Authenticater
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Client: &mockClient{},
			},
			wantErr: false,
		},
		{
			name: "ng",
			fields: fields{
				Client: &mockClient{
					err: errors.New("i/o timeout"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				ListenPort:    tt.fields.ListenPort,
				Client:        tt.fields.Client,
				Authenticater: tt.fields.Authenticater,
			}
			if err := s.CheckReadiness(); (err != nil) != tt.wantErr {
				t.Errorf("Server.CheckReadiness() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
