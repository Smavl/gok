package config

import (
	"time"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/operations"
	"github.com/smavl/gok/internal/prober/strategy"
	"github.com/smavl/gok/internal/prober/strategy/binary"
	"github.com/smavl/gok/internal/prober/types"
)

func newLinuxBuilder(mode domain.ProbingMode) OSPhaseBuilder {
	switch mode {
	case domain.Stealth:
		return &linuxStealthPhaseBuilder{}
	case domain.Agressive:
		return &linuxAggressivePhaseBuilder{}
	case domain.Default:
			return &linuxDefaultPhaseBuilder{}
	default:
		panic("unsupported mode for Linux OS")
	}
}

// ===== Linux Default =====

type linuxDefaultPhaseBuilder struct{}

func (b *linuxDefaultPhaseBuilder) BuildInitialPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	basicBinaries := []string{
		"which", "base64",
		"python", "python3",
		"bash",
	}

	initalPhase := types.PhaseConfig{
		Operations: []types.ProbeOperation{
			operations.EnumerateBinaries(
				basicBinaries,
				binary.NewWhichStrategy(),
				[]strategy.BinaryCheckStrategy{},
			),
		},
		TimeoutPerOp: 500 * time.Millisecond,
	}
	
	return &initalPhase, true
}

func (b *linuxDefaultPhaseBuilder) BuildReconPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	// TODO: Implement recon phase
	// panic("not implemented")
	return nil, false
}

func (b *linuxDefaultPhaseBuilder) BuildDeepScanPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	// TODO: Implement deep scan phase
	// panic("not implemented")
	return nil, false
}

// ===== Linux Stealth =====

type linuxStealthPhaseBuilder struct{}

func (b *linuxStealthPhaseBuilder) BuildInitialPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

func (b *linuxStealthPhaseBuilder) BuildReconPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

func (b *linuxStealthPhaseBuilder) BuildDeepScanPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

// ===== Linux Aggressive =====

type linuxAggressivePhaseBuilder struct{}

func (b *linuxAggressivePhaseBuilder) BuildInitialPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

func (b *linuxAggressivePhaseBuilder) BuildReconPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

func (b *linuxAggressivePhaseBuilder) BuildDeepScanPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}
