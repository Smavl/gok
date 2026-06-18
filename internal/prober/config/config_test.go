package config

import (
	"testing"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/types"
	"github.com/stretchr/testify/require"
)

func TestPhaseBuilderContextPropagation(t *testing.T) {
	// Given: results accumulated from Genesis and Initial phase
	// OS: Found during Genesis
	// BinariesFound: Found during Initial phase
	results := &types.ProbeResults{
		OS:            types.LinuxOs,
		BinariesFound: []string{"which", "base64", "python3"},
	}

	bctx := types.PhaseBuilderContext{
		ProbeResults: results,
		Mode:         domain.Default,
	}

	// When: getting the builder for the next phase
	builder, err := newOSPhaseBuilder(bctx.ProbeResults.OS, bctx.Mode)

	// Then: builder receives the accumulated results
	require.NoError(t, err)
	require.NotNil(t, builder)
	_, ok := builder.(*linuxDefaultPhaseBuilder)
	require.True(t, ok, "Expected linuxDefaultPhaseBuilder for Linux OS with Default mode")

	// When: building the recon phase
	reconPhase, present := builder.BuildReconPhase(bctx)

	// Then: builder received context with accumulated results
	require.False(t, present)
	require.Nil(t, reconPhase)
	require.Equal(t, types.LinuxOs, bctx.ProbeResults.OS)
	require.Contains(t, bctx.ProbeResults.BinariesFound, "which")
	require.Contains(t, bctx.ProbeResults.BinariesFound, "base64")
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

func TestOSDetectionToBuilderSelection(t *testing.T) {
	// Given: Genesis phase detected Linux
	genesisResults := &types.ProbeResults{
		OS: types.LinuxOs,
	}

	bctx := types.PhaseBuilderContext{
		ProbeResults: genesisResults,
		Mode:         domain.Default,
	}

	// When: building the initial phase
	phaseConfig, present := buildInitialPhase(bctx)

	// Then: Linux builder is selected and returns a valid phase
	require.True(t, present)
	require.NotNil(t, phaseConfig)
	require.NotEmpty(t, phaseConfig.Operations)
	require.Greater(t, len(phaseConfig.Operations), 0)
}

func TestBuilderReceivesContext(t *testing.T) {
	// Given: Genesis phase detected Linux OS
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

	// When: initial phase runs and finds binaries
	initialPhase, present := builder.BuildInitialPhase(bctx)
	require.True(t, present)
	require.NotNil(t, initialPhase)

	// Simulate that initial phase found some binaries
	results.BinariesFound = append(results.BinariesFound, "which", "python3")

	// When: recon phase is built with updated context
	reconPhase, present := builder.BuildReconPhase(bctx)

	// Then: recon phase receives the accumulated results
	require.False(t, present)
	require.Nil(t, reconPhase)
	require.Contains(t, bctx.ProbeResults.BinariesFound, "which")
	require.Contains(t, bctx.ProbeResults.BinariesFound, "python3")
}
