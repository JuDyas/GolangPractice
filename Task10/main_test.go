package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Test_directoryTree(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}
	fmt.Println("Current working directory:", wd)
	type args struct {
		path  string
		level int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Directory tree",
			args: args{
				path: "./",
			},
			want: []string{
				"Test Dir 1",
				".    Test.py",
				".    Test.txt",
				"Test Dir 2",
				".    Test Dir 3",
				".    .    Test Dir 4",
				".    .    .    trst.js",
				".    .    test.css",
				".    test.html",
				"main.go",
				"main_test.go",
			},
			wantErr: false,
		},
		{
			name: "Dont exist directory",
			args: args{
				path: "./Koko",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := directoryTree(tt.args.path, tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("directoryTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("directoryTree() got = %v, want %v", got, tt.want)
			}
		})
	}
}
