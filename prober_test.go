package main

import "testing"


func TestHasErrorPattern(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		patterns []string
		want     bool
	}{
		{
			name:     "Linux bash error",
			lines:    []string{"bash: xyz: command not found"},
			patterns: []string{"command not found"},
			want:     true,
		},
		{
			name:     "Linux sh error",
			lines:    []string{"sh: ahfja: command not found"},
			patterns: []string{"command not found"},
			want:     true,
		},
		{
			name:     "No match",
			lines:    []string{"some output"},
			patterns: []string{"error"},
			want:     false,
		},
		// TODO: Add windows when implemented
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasErrorPattern(tt.lines, tt.patterns...)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
