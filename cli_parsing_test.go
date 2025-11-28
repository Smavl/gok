package main

import (
	"testing"
	"github.com/alecthomas/kong"
)

// Add your CLI parsing tests here

func TestPortRangeSimple(t *testing.T) {
	// Given
	args := []string{"-p", "9001-9005"}

	// When
	parser := kong.Must(&Flags)
	_, err := parser.Parse(args)


	// Then
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}

	got := Flags.PortRange.Ports
	want := []int{9001,9002,9003,9004,9005}

	if len(want) != len(got) {
		t.Fatalf("Length mismatch: got %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Index %d: got %v, want %v", i, got[i], want[i])
		}
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
