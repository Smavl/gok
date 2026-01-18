//go:build !integration
package main

import "testing"


// test error based os detection
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

// Infer OS based on error messages
func TestInferErrorBasedOSDetection(t *testing.T) {
	tests := []struct {
		name    string
		output  []string
		wantOS  OS
		wantErr error
	}{
		// Positive testing
		{
			name:    "Linux detection: bash",
			output:  []string{"bash: xyzwv: command not found"},
			wantOS:  Linux,
			wantErr: nil,
		},
		{
			name:    "Linux detection: sh",
			output:  []string{"sh: ashfa: command not found"},
			wantOS:  Linux,
			wantErr: nil,
		},
		// Negative testing
		{
			name:    "Garbage string",
			output:  []string{"sdhgbjskdb"},
			wantOS:  Unknown,
			wantErr: CouldNotDetermineOSError,
		},
		{
			name:    "Empty string",
			output:  []string{""},
			wantOS:  Unknown,
			wantErr: CouldNotDetermineOSError,
		},
}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOS, gotErr := inferOsByError(tt.output)
			if gotOS != tt.wantOS {
				t.Errorf("got OS %v, want %v", gotOS, tt.wantOS)
			}
			if gotErr != tt.wantErr {
				t.Errorf("got error %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

