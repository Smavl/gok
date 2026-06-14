package config

import (
	"time"

	"github.com/smavl/gok/internal/prober/operations"
	"github.com/smavl/gok/internal/prober/strategy"
	"github.com/smavl/gok/internal/prober/strategy/binary"
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
			types.PhaseInitial: buildInitialPhase,
			// types.PhaseRecon:   reconPhase,
			// types.PhaseDeepScan: deepScanPhase,
		},
	}

	return defaultConfig
}

func buildInitialPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {

	switch bctx.ProbeResults.OS {
	case types.LinuxOs:
		return buildLinuxInitialPhase(bctx)
	case types.WindowsOs:
		return buildWindowsInitialPhase(bctx)
	default:
		// OS not detected or unsupported, skip this phase
		return nil, false
	}
}

func buildLinuxInitialPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {

	basicBinaries:= []string{
		"which", "base64", 
		"python", "python3", 
	}

	return &types.PhaseConfig{
		Operations: []types.ProbeOperation{
			operations.EnumerateBinaries(
				basicBinaries,
				binary.NewWhichStrategy(),
				[]strategy.BinaryCheckStrategy{},
				),
		},
		TimeoutPerOp: 500 * time.Millisecond,
	}, true
}

func buildWindowsInitialPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("windows initial phase not implemented yet")
}

