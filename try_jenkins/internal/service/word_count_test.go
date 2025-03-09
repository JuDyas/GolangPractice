package service

import (
	"testing"
)

func TestWordCountService_CountWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "single word",
			input:    "hello",
			expected: 1,
		},
		{
			name:     "multiple words",
			input:    "hello world",
			expected: 2,
		},
		{
			name:     "words with symbols",
			input:    "hello, world!",
			expected: 2,
		},
		{
			name:     "words with numbers",
			input:    "g2g ready",
			expected: 2,
		},
		{
			name:     "text with newline characters",
			input:    "hello\nworld",
			expected: 2,
		},
		{
			name:     "mixed case text",
			input:    "Hello WoRLd",
			expected: 2,
		},
		{
			name:     "text with consecutive spaces",
			input:    "hello      world",
			expected: 2,
		},
		{
			name:     "text with no valid words",
			input:    "12345 6789",
			expected: 0,
		},
		{
			name:     "text with special characters only",
			input:    "!@#$%^&*()",
			expected: 0,
		},
	}

	svc := NewWordCountService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.CountWords(tt.input)
			if got != tt.expected {
				t.Errorf("CountWords() = %d, want %d", got, tt.expected)
			}
		})
	}
}
