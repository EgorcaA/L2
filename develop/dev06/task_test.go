package main

import (
	"reflect"
	"testing"
)

func TestGrepLines(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		lines    []string
		options  Options
		expected []string
	}{
		{
			name:    "basic match",
			pattern: "Go",
			lines: []string{
				"Go is a great language.",
				"I love programming in Go.",
				"Gophers are amazing.",
			},
			options: Options{},
			expected: []string{
				"Go is a great language.",
				"I love programming in Go.",
			},
		},
		{
			name:    "case insensitive match",
			pattern: "go",
			lines: []string{
				"Go is a great language.",
				"I love programming in Go.",
				"Gophers are amazing.",
			},
			options: Options{ignoreCase: true},
			expected: []string{
				"Go is a great language.",
				"I love programming in Go.",
			},
		},
		{
			name:    "inverted match",
			pattern: "Go",
			lines: []string{
				"Go is a great language.",
				"I love programming in Go.",
				"Gophers are amazing.",
			},
			options: Options{invert: true},
			expected: []string{
				"Gophers are amazing.",
			},
		},
		{
			name:    "line numbers",
			pattern: "Go",
			lines: []string{
				"Go is a great language.",
				"I love programming in Go.",
				"Gophers are amazing.",
			},
			options: Options{lineNumbers: true},
			expected: []string{
				"1: Go is a great language.",
				"2: I love programming in Go.",
			},
		},
		{
			name:    "count matches",
			pattern: "Go",
			lines: []string{
				"Go is a great language.",
				"I love programming in Go.",
				"Gophers are amazing.",
			},
			options:  Options{count: true},
			expected: []string{"2"},
		},
		{
			name:    "context match",
			pattern: "Go",
			lines: []string{
				"Go is a great language.",
				"I love programming in Go.",
				"Gophers are amazing.",
				"Go rocks!",
			},
			options: Options{context: 1},
			expected: []string{
				"Go is a great language.",
				"I love programming in Go.",
				"Gophers are amazing.",
				"Go rocks!",
			},
		},
		{
			name:    "exact match",
			pattern: "Go",
			lines: []string{
				"Go is a great language.",
				"Go",
				"I love programming in Go.",
			},
			options: Options{fixed: true},
			expected: []string{
				"Go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := grepLines(tt.pattern, tt.lines, tt.options)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
