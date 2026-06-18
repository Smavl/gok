package config

import (
	"time"

	"github.com/smavl/gok/internal/prober/operations"
	"github.com/smavl/gok/internal/prober/strategy"
	// "github.com/smavl/gok/internal/prober/strategy/binary"
	"github.com/smavl/gok/internal/prober/strategy/os"
	"github.com/smavl/gok/internal/prober/types"
)

func DefaultConfig() types.ProbeConfig {

	// Genesis phase (e.g. OS Detection)
	genesisPhase := types.PhaseConfig{
		Operations: []types.ProbeOperation{
			operations.DetectOS(
				os.NewOSErrorDetectionStrategy(),
				[]strategy.OSDetectionStrategy{},
				),
		},
		TimeoutPerOp: 500 * time.Millisecond,

	}



	defaultConfig := types.ProbeConfig{
		Genesis: genesisPhase,
		Phases: map[types.ProbePhase]types.PhaseBuilder{
			types.PhaseInitial:		buildInitialPhase,
			types.PhaseRecon:		buildReconPhase,
			types.PhaseDeepScan:	buildDeepScanPhase,
		},
	}

	return defaultConfig
}


func buildInitialPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	builder, err := newOSPhaseBuilder(bctx.ProbeResults.OS, bctx.Mode)
	if err != nil {
		// ERROR: failed to get os phase builder
		return nil, false
	}

	return builder.BuildInitialPhase(bctx)
}

func buildReconPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	builder, err := newOSPhaseBuilder(bctx.ProbeResults.OS, bctx.Mode)
	if err != nil {
		// ERROR: failed to get os phase builder
		return nil, false
	}

	return builder.BuildReconPhase(bctx)
}

func buildDeepScanPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	builder, err := newOSPhaseBuilder(bctx.ProbeResults.OS, bctx.Mode)
	if err != nil {
		// ERROR: failed to get os phase builder
		return nil, false
	}

	return builder.BuildDeepScanPhase(bctx)
}

