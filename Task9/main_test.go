package main

import (
	"reflect"
	"testing"
)

func Test_countLetters(t *testing.T) {
	type args struct {
		letters      []string
		countWorkers int
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "Simple count",
			args: args{
				letters:      []string{"a", "а", "ы", "ї"},
				countWorkers: 2,
			},
			want:    []int{1, 1, 1, 1},
			wantErr: false,
		},
		{
			name: "Empty input",
			args: args{
				letters:      []string{},
				countWorkers: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countLetters(tt.args.letters, tt.args.countWorkers)
			if (err != nil) != tt.wantErr {
				t.Errorf("countLetters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countLetters() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readInput(t *testing.T) {
	type args struct {
		inputString string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Simple input string",
			args: args{
				inputString: "lorem ipsumпривет я тут це що",
			},
			want: []string{"l", "o", "r", "e", "m", "i", "p", "s", "u", "m", "п", "р", "и", "в", "е", "т", "я", "т", "у", "т", "ц", "е", "щ", "о"},
		},
		{
			name: "Empty input string",
			args: args{
				inputString: "",
			},
			want: nil,
		},
		{
			name: "String with non char",
			args: args{
				inputString: "lor242.,/4em ips^+um.,прив]е353т я 67ту!т 123ц$е що",
			},
			want: []string{"l", "o", "r", "e", "m", "i", "p", "s", "u", "m", "п", "р", "и", "в", "е", "т", "я", "т", "у", "т", "ц", "е", "щ", "о"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readInput(tt.args.inputString)
			if (err != nil) != tt.wantErr {
				t.Errorf("readInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readInput() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isEn(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "is it english", args: args{r: 'a'}, want: true},
		{name: "is not english", args: args{r: 'ы'}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEn(tt.args.r); got != tt.want {
				t.Errorf("isEn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isRu(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "is it russian", args: args{r: 'ы'}, want: true},
		{name: "is not russian", args: args{r: 'a'}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRu(tt.args.r); got != tt.want {
				t.Errorf("isRu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isRuUa(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "is it Cyrillic", args: args{r: 'р'}, want: true},
		{name: "is not Cyrillic", args: args{r: 'ы'}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRuUa(tt.args.r); got != tt.want {
				t.Errorf("isRuUa() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isUa(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "is it ukrainian", args: args{r: 'ї'}, want: true},
		{name: "is not ukrainian", args: args{r: 'п'}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUa(tt.args.r); got != tt.want {
				t.Errorf("isUa() = %v, want %v", got, tt.want)
			}
		})
	}
}
