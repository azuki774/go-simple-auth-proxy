package repository

import (
	"sync"
	"testing"
)

func TestStore_CheckCookieValue(t *testing.T) {
	type fields struct {
		CookieStore []string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "exist: ZmFkNThkODAtMjdjNi00Y2ExLWE0OTAtOTM3ZjNmODE3YWVl", // fad58d80-27c6-4ca1-a490-937f3f817aee
			fields: fields{
				CookieStore: []string{"a", "b", "ZmFkNThkODAtMjdjNi00Y2ExLWE0OTAtOTM3ZjNmODE3YWVl"},
			},
			args: args{
				value: "ZmFkNThkODAtMjdjNi00Y2ExLWE0OTAtOTM3ZjNmODE3YWVl",
			},
			want: true,
		},
		{
			name: "not exist: YmE5MzA4NGQtOWExYy00NjhjLThlZmItZmVlZjU2YmNhODlm", // ba93084d-9a1c-468c-8efb-feef56bca89f
			fields: fields{
				CookieStore: []string{"a", "b", "ZmFkNThkODAtMjdjNi00Y2ExLWE0OTAtOTM3ZjNmODE3YWVl"},
			},
			args: args{
				value: "YmE5MzA4NGQtOWExYy00NjhjLThlZmItZmVlZjU2YmNhODlm",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StoreInMemory{
				CookieStore: tt.fields.CookieStore,
				Mu:          &sync.Mutex{}, // テストのため新規作成
			}
			if got := s.CheckCookieValue(tt.args.value); got != tt.want {
				t.Errorf("StoreInMemory.CheckCookieValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_InsertCookieValue(t *testing.T) {
	type fields struct {
		CookieStore []string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "0 -> 1",
			fields: fields{
				CookieStore: []string{},
			},
			args: args{
				value: "SUPER_SUGOI_COOKIE",
			},
			wantErr: false,
		},
		{
			name: "1 -> 2",
			fields: fields{
				CookieStore: []string{"EXTERME_SUPER_SUGOI_COOKIE"},
			},
			args: args{
				value: "SUPER_SUGOI_COOKIE",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StoreInMemory{
				CookieStore: tt.fields.CookieStore,
				Mu:          &sync.Mutex{},
			}
			if err := s.InsertCookieValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("StoreInMemory.InsertCookieValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
