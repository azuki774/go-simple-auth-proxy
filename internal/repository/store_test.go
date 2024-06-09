package repository

import (
	"sync"
	"testing"
)

func TestStore_CheckCookieValue(t *testing.T) {
	type fields struct {
		CookieStore []string
		mu          sync.Mutex
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
				mu:          sync.Mutex{},
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
				mu:          sync.Mutex{},
			},
			args: args{
				value: "YmE5MzA4NGQtOWExYy00NjhjLThlZmItZmVlZjU2YmNhODlm",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				CookieStore: tt.fields.CookieStore,
				mu:          tt.fields.mu,
			}
			if got := s.CheckCookieValue(tt.args.value); got != tt.want {
				t.Errorf("Store.CheckCookieValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_InsertCookieValue(t *testing.T) {
	type fields struct {
		CookieStore []string
		mu          sync.Mutex
	}
	tests := []struct {
		name      string
		fields    fields
		wantValue string
		wantErr   bool
	}{
		{
			name: "empty -> 1",
			fields: fields{
				CookieStore: []string{},
				mu:          sync.Mutex{},
			},
			wantErr: false,
		},
		{
			name: "1 -> 2",
			fields: fields{
				CookieStore: []string{"abcde"},
				mu:          sync.Mutex{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				CookieStore: tt.fields.CookieStore,
				mu:          tt.fields.mu,
			}
			_, err := s.InsertCookieValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.InsertCookieValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// ランダム生成なのでパス
			// if gotValue != tt.wantValue {
			// 	t.Errorf("Store.InsertCookieValue() = %v, want %v", gotValue, tt.wantValue)
			// }
		})
	}
}
