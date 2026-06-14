package config

import (
	// "time"

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
	}



	defaultConfig := types.ProbeConfig{
		Genesis: genesisPhase,
	}

	return defaultConfig
}

// func DefaultConfig() types.ProbeConfig {
// 	basicBinaries:= []string{
// 		"which", "base64", 
// 		"python", "python3", 
// 	}
// 	commonBinaries := []string{
// 		"bash", "sh",
// 		"curl", "wget", "nc", 
// 		"find", "grep",
// 	}
//
// 	initialPhase := types.PhaseConfig{
// 		Operations: []types.ProbeOperation{
// 			// TODO: Detect os
// 			// Enumerate Essential Binaries
// 			operations.EnumerateBinaries(
// 				basicBinaries,
// 				binary.NewWhichStrategy(),
// 				[]strategy.BinaryCheckStrategy{},
// 			),
// 		},
// 		// TODO: FAKE-IT
// 		TimeoutPerOp: 500 * time.Millisecond,
// 	}
//
// 	reconPhase := types.PhaseConfig{
// 		Operations: []types.ProbeOperation{
// 			// Enumerate binaries using which command
// 			operations.EnumerateBinaries(
// 				commonBinaries,
// 				binary.NewWhichStrategy(),
// 				[]strategy.BinaryCheckStrategy{},
// 			),
// 		},
// 		// TODO: FAKE-IT
// 		TimeoutPerOp: 500 * time.Millisecond,
// 	}
//
// 	defaultConfig := map[types.ProbePhase]types.PhaseConfig{
// 		types.PhaseInitial: initialPhase,
// 		types.PhaseRecon:   reconPhase,
// 	}
//
// 	return types.ProbeConfig{Phases: defaultConfig}
// }
