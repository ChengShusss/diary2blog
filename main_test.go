package main

import (
	"reflect"
	"testing"
	"time"
)

func TestNormalizeTime(t *testing.T) {

	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				s: "20230101",
			},
			want:    time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		}, {
			name: "test",
			args: args{
				s: "20230228",
			},
			want:    time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		}, {
			name: "test",
			args: args{
				s: "20230531",
			},
			want:    time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeDate(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalizeTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunc(t *testing.T) {
	got, err := normalizeDate("20240614")
	if err != nil {
		t.Errorf("normalizeTime() error = %v", err)
		return
	}
	t.Logf("normalizeTime() got = %v", got)
}
