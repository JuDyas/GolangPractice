package main

import "testing"

func Test_palindrome(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "One word palindrome",
			args: args{text: "madam"},
			want: true,
		},
		{
			name: "Not palindrome",
			args: args{text: "hello"},
			want: false,
		},
		{
			name: "Palindrome with many spaces",
			args: args{text: "A man a plan a canal Panama"},
			want: true,
		},
		{
			name: "Empty text",
			args: args{text: ""},
			want: true,
		},
		{
			name: "Mixed case",
			args: args{text: "Ra.]ce:23Car"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := palindrome(tt.args.text); got != tt.want {
				t.Errorf("palindrome() = %v, want %v", got, tt.want)
			}
		})
	}
}
