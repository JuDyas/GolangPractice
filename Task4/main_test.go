package main

import (
	"reflect"
	"testing"
)

func Test_calculate(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple addition",
			args: args{text: "2 + 2"},
			want: "4.000000",
		},
		{
			name: "Simple subtraction",
			args: args{text: "5 - 3"},
			want: "2.000000",
		},
		{
			name: "Mixed operations",
			args: args{text: "2 + 3 * 4"},
			want: "14.000000",
		},
		{
			name: "Mixed with division",
			args: args{"10 / 2 + 3"},
			want: "8.000000",
		},
		{
			name: "Multiple operations",
			args: args{"2 + 2 * 2"},
			want: "6.000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculate(tt.args.text); got != tt.want {
				t.Errorf("calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addsubst(t *testing.T) {
	type args struct {
		expression []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Simple addition", args{[]string{"2", "+", "3"}}, "5.000000"},
		{"Simple subtraction", args{[]string{"5", "-", "2"}}, "3.000000"},
		{"Mixed addition and subtraction", args{[]string{"10", "-", "3", "+", "5"}}, "12.000000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addsubst(tt.args.expression); got != tt.want {
				t.Errorf("addsubst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_multyply(t *testing.T) {
	type args struct {
		expression []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"Simple multiplication", args{[]string{"2", "*", "3"}}, []string{"6.000000"}},
		{"Multiplication and division", args{[]string{"6", "/", "2"}}, []string{"3.000000"}},
		{"Mixed operations", args{[]string{"4", "*", "2", "/", "4"}}, []string{"2.000000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := multyply(tt.args.expression); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("multyply() = %v, want %v", got, tt.want)
			}
		})
	}
}
