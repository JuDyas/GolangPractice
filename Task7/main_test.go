package main

import (
	"reflect"
	"testing"
)

func Test_hashJson(t *testing.T) {
	jsonData := map[string]interface{}{
		"name":     "Denys",
		"age":      "123",
		"password": "UpbQ2X*&FQ$L",
		"info": map[string]interface{}{
			"weight": "60",
			"ip":     "192.168.1.1",
		},
	}
	sha1Type := "sha1"
	md5Type := "md5"
	sha256Type := "sha256"

	type args struct {
		data     map[string]interface{}
		whatHash []string
		hashType *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test SHA1 hashing",
			args: args{
				data: map[string]interface{}{
					"name":     "Denys",
					"age":      "123",
					"password": "UpbQ2X*&FQ$L",
					"info": map[string]interface{}{
						"weight": "60",
						"ip":     "192.168.1.1",
					},
				},
				whatHash: []string{"password", "ip"},
				hashType: &sha1Type,
			},
			want: `{
  "age": "123",
  "info": {
    "ip": "90aa44756bd2f4fc2390f903a6f25f43216b0790",
    "weight": "60"
  },
  "name": "Denys",
  "password": "394225babd1c019e83bc52332b227b56c8827d45"
}`,
			wantErr: false,
		},

		{
			name: "Test MD5 hashing",
			args: args{
				data: map[string]interface{}{
					"name":     "Denys",
					"age":      "123",
					"password": "UpbQ2X*&FQ$L",
					"info": map[string]interface{}{
						"weight": "60",
						"ip":     "192.168.1.1",
					},
				},
				whatHash: []string{"password", "ip"},
				hashType: &md5Type,
			},
			want: `{
  "age": "123",
  "info": {
    "ip": "66efff4c945d3c3b87fc271b47d456db",
    "weight": "60"
  },
  "name": "Denys",
  "password": "ca245dbdf8422d9640c036842545b7e5"
}`,
			wantErr: false,
		},

		{
			name: "Test unsupported hash type",
			args: args{
				data:     jsonData,
				whatHash: []string{"password"},
				hashType: &sha256Type,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hashJson(tt.args.data, tt.args.whatHash, tt.args.hashType)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hashJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processFlags(t *testing.T) {
	jsonData := `{"name": "Denys", "age": "123", "password": "UpbQ2X*&FQ$L", "info": {"weight": "60", "ip": "192.168.1.1"}}`
	errJsonData := `{"name": "Denys"`
	hashElements := "password,ip"

	type args struct {
		jsonData     *string
		hashElements *string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		want1   []string
		wantErr bool
	}{
		{
			name: "Successful processing input data",
			args: args{
				jsonData:     &jsonData,
				hashElements: &hashElements,
			},
			want: map[string]interface{}{
				"name":     "Denys",
				"age":      "123",
				"password": "UpbQ2X*&FQ$L",
				"info": map[string]interface{}{
					"weight": "60",
					"ip":     "192.168.1.1",
				},
			},
			want1:   []string{"password", "ip"},
			wantErr: false,
		},
		{
			name: "Incorrect JSON data",
			args: args{
				jsonData:     &errJsonData,
				hashElements: &hashElements,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := processFlags(tt.args.jsonData, tt.args.hashElements)
			if (err != nil) != tt.wantErr {
				t.Errorf("processFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processFlags() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("processFlags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_md5Hash(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Successful md5 hash test",
			args: args{data: "jdjsjskfdfjslfj"},
			want: "8b1d905f8988ac97032749b73d2ae662",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := md5Hash(tt.args.data); got != tt.want {
				t.Errorf("md5Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sha1Hash(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Successful md5 hash test",
			args: args{data: "jdjsjskfdfjslfj"},
			want: "c59041b3ac235125ae61090eea87af31b0d96022",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sha1Hash(tt.args.data); got != tt.want {
				t.Errorf("sha1Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
