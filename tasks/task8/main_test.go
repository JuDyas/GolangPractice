package main

import (
	"errors"
	"testing"
	"time"
)

type mockSuccessRequest struct{}

func (r *mockSuccessRequest) Request() error {
	return nil
}

type mockBadRequest struct{}

func (r *mockBadRequest) Request() error {
	return errors.New("bad request")
}

type mockServerStartingUp struct{}

func (r *mockServerStartingUp) Request() error {
	return errors.New("server starting up")
}

type mockServerStartingUpSuccess struct {
	callCount   int
	successCall int
}

func (m *mockServerStartingUpSuccess) Request() error {
	m.callCount++
	if m.callCount >= m.successCall {
		return nil
	}
	return errors.New("server starting up")
}

func Test_retry(t *testing.T) {
	type args struct {
		attempts int
		sleep    time.Duration
		r        Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success on first attempt",
			args: args{
				attempts: 3,
				sleep:    time.Millisecond * 10,
				r:        &mockSuccessRequest{},
			},
			wantErr: false,
		},
		{
			name: "Fail due to bad request",
			args: args{
				attempts: 3,
				sleep:    time.Millisecond * 25,
				r:        &mockBadRequest{},
			},
			wantErr: true,
		},
		{
			name: "Fail due to server starting up",
			args: args{
				attempts: 3,
				sleep:    time.Millisecond * 25,
				r:        &mockServerStartingUp{},
			},
			wantErr: true,
		},
		{
			name: "Server starts up after 2 attempts",
			args: args{
				attempts: 5,
				sleep:    time.Millisecond * 25,
				r:        &mockServerStartingUpSuccess{successCall: 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := retry(tt.args.attempts, tt.args.sleep, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("retry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
