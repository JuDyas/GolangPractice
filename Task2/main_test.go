package main

import "testing"

func Test_cezar(t *testing.T) {
	type args struct {
		text  string
		shift int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Shift for 3 letters",
			args: args{
				text:  "I`ll be back",
				shift: 3,
			},
			want: "L`oo eh edfn",
		},
		{
			name: "Shift for 0 letters",
			args: args{
				text:  "I`ll be back",
				shift: 0,
			},
			want: "I`ll be back",
		},
		{
			name: "Shift for 27 (1) letters",
			args: args{
				text:  "I`ll be back",
				shift: 27,
			},
			want: "J`mm cf cbdl",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cezar(tt.args.text, tt.args.shift); got != tt.want {
				t.Errorf("cezar() = %v, want %v", got, tt.want)
			}
		})
	}
}
