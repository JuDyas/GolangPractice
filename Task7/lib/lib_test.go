package lib

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

type MockHasher struct {
	ReturnError bool
	HashType    string
}

func (m *MockHasher) Hash(data string) (string, error) {
	if m.ReturnError {
		return "", errors.New("hash error")
	} else if m.HashType != "" {
		if m.HashType == "md5" {
			hash := md5.New()
			hash.Write([]byte(data))
			return hex.EncodeToString(hash.Sum(nil)), nil
		} else if m.HashType == "sha1" {
			hash := sha1.New()
			hash.Write([]byte(data))
			return hex.EncodeToString(hash.Sum(nil)), nil
		}
	}
	return "", errors.New("hash error")
}

func TestHashJson(t *testing.T) {
	type args struct {
		data     map[string]interface{}
		whatHash []string
		hasher   Hasher
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Simple hashing md5",
			args: args{
				data: map[string]interface{}{
					"password": "101010110122",
					"email":    "opana@golang.com",
				},
				whatHash: []string{"password", "email"},
				hasher:   &MockHasher{ReturnError: false, HashType: "md5"},
			},
			want: map[string]interface{}{
				"password": "c5e325ca8d26c4c2163b714d69745792",
				"email":    "ef7171a01d8ce1b8b5dc8b5c937bd75e",
			},
			wantErr: false,
		},
		{
			name: "Simple hashing sha1",
			args: args{
				data: map[string]interface{}{
					"password": "101010110122",
					"email":    "opana@golang.com",
				},
				whatHash: []string{"password", "email"},
				hasher:   &MockHasher{ReturnError: false, HashType: "sha1"},
			},
			want: map[string]interface{}{
				"password": "8efa617034be6a6e3d8ee9ee71842805fe4196e3",
				"email":    "f94bf8a9ec4e2f551a5a62a8ccc4688c6556d90b",
			},
			wantErr: false,
		},
		{
			name: "Hashing nested object",
			args: args{
				data: map[string]interface{}{
					"password": "101010110122",
					"info": map[string]interface{}{
						"email": "test@example.com",
					},
				},
				whatHash: []string{"password", "email"},
				hasher:   &MockHasher{ReturnError: false, HashType: "sha1"},
			},
			want: map[string]interface{}{
				"password": "8efa617034be6a6e3d8ee9ee71842805fe4196e3",
				"info": map[string]interface{}{
					"email": "567159d622ffbb50b11b0efd307be358624a26ee",
				},
			},
			wantErr: false,
		},
		{
			name: "Error hashing",
			args: args{
				data: map[string]interface{}{
					"password": "101010110122",
					"email":    "opana@golang.com",
				},
				whatHash: []string{"password", "email"},
				hasher:   &MockHasher{ReturnError: true, HashType: "sha1"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashJson(tt.args.data, tt.args.whatHash, tt.args.hasher)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != "" {
				var gotMap map[string]interface{}
				if err := json.Unmarshal([]byte(got), &gotMap); err != nil {
					t.Errorf("Failed to unmarshal got JSON: %v", err)
				}
				if !reflect.DeepEqual(gotMap, tt.want) {
					t.Errorf("HashJson() got = %v, want %v", gotMap, tt.want)
				}
			}
		})
	}
}
