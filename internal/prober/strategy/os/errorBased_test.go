package os

import (
	"testing"

	"github.com/smavl/gok/internal/prober/types"
)

func TestInferOsByError(t *testing.T) {
	tests := []struct {
		name     string
		output   []string
		expected types.OS
		wantErr  bool
	}{
		{"Linux bash", []string{"bash: xyz: command not found"}, types.LinuxOs, false},
		{"Linux sh", []string{"sh: xyz: command not found:"}, types.LinuxOs, false},
		{"Windows", []string{"'xyzqwerty' is not recognized as an internal or external command"}, types.WindowsOs, false},
		{"Windows variant", []string{"something not recognized"}, types.WindowsOs, false},
		{"Unknown", []string{"random error"}, types.UnknownOS, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inferOsByError(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("inferOsByError() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.expected {
				t.Errorf("inferOsByError() = %v, want %v", got, tt.expected)
			}
		})
	}
}
