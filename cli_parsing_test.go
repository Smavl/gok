package main

import (
	"github.com/alecthomas/kong"
	"testing"
)

// table driven tests for port range parsing
func TestPortRangeTableDriven(t *testing.T) {
	tests := []struct {
		name      string
		arg       string
		wantPorts []int
		wantErr   bool
	}{
		// Positive test cases
		{
			name:      "Single port",
			arg:       "9001",
			wantPorts: []int{9001},
			wantErr:   false,
		},
		{
			name:      "Simple range",
			arg:       "9000-9002",
			wantPorts: []int{9000, 9001, 9002},
			wantErr:   false,
		},
		{
			name:      "Bigger range",
			arg:       "9000-9069",
			wantPorts: []int{
				9000, 9001, 9002, 9003, 9004, 9005, 9006, 9007, 9008, 9009, 9010,
				9011, 9012, 9013, 9014, 9015, 9016, 9017, 9018, 9019, 9020,
				9021, 9022, 9023, 9024, 9025, 9026, 9027, 9028, 9029, 9030,
				9031, 9032, 9033, 9034, 9035, 9036, 9037, 9038, 9039, 9040,
				9041, 9042, 9043, 9044, 9045, 9046, 9047, 9048, 9049, 9050,
				9051, 9052, 9053, 9054, 9055, 9056, 9057, 9058, 9059, 9060,
				9061, 9062, 9063, 9064, 9065, 9066, 9067, 9068, 9069},
			wantErr: false,
		},
		{
			name:    "Max valid port: 65535",
			arg:     "65535",
			wantPorts: []int{65535},
			wantErr: false,
		},
		// Negative test cases
		{
			name:    "Wrong Delimiter: _",
			arg:     "9000_9005",
			wantErr: true,
		},
		{
			name:    "Wrong Delimiter: ...",
			arg:     "9000...9005",
			wantErr: true,
		},
		{
			name:    "Iverse order",
			arg:     "9005-9000",
			wantErr: true,
		},
		{
			name:    "Non-numeric start",
			arg:     "abc-9005",
			wantErr: true,
		},
		{
			name:    "Non-numeric end",
			arg:     "9000-xyz",
			wantErr: true,
		},
		{
			name:    "No starting port in port range",
			arg:     "-123",
			wantErr: true,
		},
		{
			name:    "No end port in port range",
			arg:     "123-",
			wantErr: true,
		},
		{
			name:    "Invalid port: 0",
			arg:     "0",
			wantErr: true,
		},
		{
			name:    "Invalid port: 0",
			arg:     "0-2",
			wantErr: true,
		},
		{
			name:    "Invalid port:	65536",
			arg:     "65536",
			wantErr: true,
		},
		{
			name:    "Invalid PortRange: 65535-65536",
			arg:     "65535-65536",
			wantErr: true,
		},
		{
			name:    "Empty string",
			arg:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			args := []string{"-p", tt.arg}

			// When
			parser := kong.Must(&Flags)
			_, err := parser.Parse(args)

			// Then
			// TEST: check for expected error
			if (err != nil) != tt.wantErr {
				t.Fatalf("Got error: %v, wantErr: %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			got := Flags.PortRange.Ports
			want := tt.wantPorts

			// TEST: check for expected ports
			if len(want) != len(got) {
				t.Fatalf("Length mismatch: got %d, want %d", len(got), len(want))
			}
			// TEST: check each port
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("Index %d: got %v, want %v", i, got[i], want[i])
				}
			}
		})
	}
}

// func TestXX(t *testing.T) {
// 	// Given
//
// 	// When
//
// 	// Then
//
// }
