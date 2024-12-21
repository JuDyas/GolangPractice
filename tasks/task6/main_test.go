package main

import (
	"io"
	"strings"
	"testing"
	"unicode/utf8"
)

func Test_countCharacters(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "simple string",
			args:    args{reader: strings.NewReader("Hello, world!")},
			want:    utf8.RuneCountInString("Hello, world!"),
			wantErr: false,
		},
		{
			name:    "string with spaces",
			args:    args{reader: strings.NewReader("  Go is great  ")},
			want:    utf8.RuneCountInString("Go is great"),
			wantErr: false,
		},
		{
			name:    "empty string",
			args:    args{reader: strings.NewReader("")},
			want:    0,
			wantErr: false,
		},
		{
			name:    "unicode string",
			args:    args{reader: strings.NewReader("Привет, мир!")},
			want:    utf8.RuneCountInString("Привет, мир!"),
			wantErr: false,
		},
		{
			name:    "newline characters",
			args:    args{reader: strings.NewReader("\n\n\n")},
			want:    0,
			wantErr: false,
		},
		{
			name:    "string with special characters",
			args:    args{reader: strings.NewReader("Hello\tworld\n")},
			want:    utf8.RuneCountInString("Hello\tworld"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countCharacters(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("countCharacters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("countCharacters() got = %v, want %v", got, tt.want)
			}
		})
	}
}
