package main

import (
	"reflect"
	"testing"
)

func Test_matrix(t *testing.T) {
	type args struct {
		A [][]int
		B [][]int
	}
	tests := []struct {
		name    string
		args    args
		want    [][]int
		wantErr bool
	}{
		{
			name: "Valid multiplication 2x3 and 3x2",
			args: args{A: [][]int{
				{1, 2, 3},
				{4, 5, 6},
			},
				B: [][]int{
					{7, 8},
					{9, 10},
					{11, 12},
				}},
			want: [][]int{
				{58, 64},
				{139, 154},
			},
			wantErr: false,
		},
		{
			name: "Invalid multiplication",
			args: args{A: [][]int{
				{1, 2},
				{3, 4},
				{5, 6},
				{7, 7},
				{8, 9},
			},
				B: [][]int{
					{5},
					{8},
					{9},
				}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matrix(tt.args.A, tt.args.B)
			if (err != nil) != tt.wantErr {
				t.Errorf("matrix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("matrix() got = %v, want %v", got, tt.want)
			}
		})
	}
}
