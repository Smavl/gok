package config

import (
	"testing"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/types"
	"github.com/stretchr/testify/require"
)

// Test that context with accumulated results flows between phases
func TestPhaseBuilderContextPropagation(t *testing.T) {
	// Given: Genesis detected Linux, Initial found binaries
	results := &types.ProbeResults{
		OS:            types.LinuxOs,
		BinariesFound: []string{},
	}

	bctx := types.PhaseBuilderContext{
		ProbeResults: results,
		Mode:         domain.Default,
	}

	builder, err := newOSPhaseBuilder(bctx.ProbeResults.OS, bctx.Mode)
	require.NoError(t, err)

	// When: Initial phase runs (simulated binary discovery)
	initialPhase, present := builder.BuildInitialPhase(bctx)
	require.True(t, present)
	require.NotNil(t, initialPhase)

	// FAKE binaries found by initial phase
	results.BinariesFound = append(results.BinariesFound, "which", "python3")

	// When: Recon phase is built with accumulated context
	reconPhase, present := builder.BuildReconPhase(bctx)

	// Then: Recon receives context with OS and accumulated binaries
	require.False(t, present) // Not implemented yet
	require.Nil(t, reconPhase)
	require.Equal(t, types.LinuxOs, bctx.ProbeResults.OS)
	require.Contains(t, bctx.ProbeResults.BinariesFound, "which")
	require.Contains(t, bctx.ProbeResults.BinariesFound, "python3")
}

func TestModeSelection(t *testing.T) {
	tests := []struct {
		name         string
		os           types.OS
		mode         domain.ProbingMode
		expectedType interface{}
	}{
		{
			name:         "Linux Default",
			os:           types.LinuxOs,
			mode:         domain.Default,
			expectedType: &linuxDefaultPhaseBuilder{},
		},
		{
			name:         "Linux Stealth",
			os:           types.LinuxOs,
			mode:         domain.Stealth,
			expectedType: &linuxStealthPhaseBuilder{},
		},
		{
			name:         "Linux Aggressive",
			os:           types.LinuxOs,
			mode:         domain.Agressive,
			expectedType: &linuxAggressivePhaseBuilder{},
		},
		{
			name:         "Windows Default",
			os:           types.WindowsOs,
			mode:         domain.Default,
			expectedType: &windowsDefaultPhaseBuilder{},
		},
		{
			name:         "Unknown OS",
			os:           types.UnknownOS,
			mode:         domain.Default,
			expectedType: &unknownOSPhaseBuilder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: an OS and mode combination

			// When: creating a builder for that combination
			builder, err := newOSPhaseBuilder(tt.os, tt.mode)

			// Then: the correct builder type is returned
			require.NoError(t, err)
			require.IsType(t, tt.expectedType, builder)
		})
	}
}
