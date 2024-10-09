package main

import (
	"reflect"
	"testing"
)

func Test_bubbleSort(t *testing.T) {
	type args struct {
		sl []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "No sort",
			args: args{sl: []int{1, 2, 3, 4, 5}},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "One number",
			args: args{sl: []int{0, 0, 0, 0, 0, 0, 0}},
			want: []int{0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "Simple sort",
			args: args{sl: []int{234, 4, 56, 2, 548, 0, 43, 12, 53, 23, 89}},
			want: []int{0, 2, 4, 12, 23, 43, 53, 56, 89, 234, 548},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bubbleSort(tt.args.sl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bubbleSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
