package config

import (
	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/types"
)

func newWindowsBuilder(mode domain.ProbingMode) OSPhaseBuilder {
	switch mode {
	case domain.Stealth:
		return &windowsStealthPhaseBuilder{}
	case domain.Agressive:
		return &windowsAggressivePhaseBuilder{}
	default:
		return &windowsDefaultPhaseBuilder{}
	}
}

// ===== Windows Default =====

type windowsDefaultPhaseBuilder struct{}

func (b *windowsDefaultPhaseBuilder) BuildInitialPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

func (b *windowsDefaultPhaseBuilder) BuildReconPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

func (b *windowsDefaultPhaseBuilder) BuildDeepScanPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

// ===== Windows Stealth =====

type windowsStealthPhaseBuilder struct{}

func (b *windowsStealthPhaseBuilder) BuildInitialPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

func (b *windowsStealthPhaseBuilder) BuildReconPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

func (b *windowsStealthPhaseBuilder) BuildDeepScanPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

// ===== Windows Aggressive =====

type windowsAggressivePhaseBuilder struct{}

func (b *windowsAggressivePhaseBuilder) BuildInitialPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

func (b *windowsAggressivePhaseBuilder) BuildReconPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}

func (b *windowsAggressivePhaseBuilder) BuildDeepScanPhase(bctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	panic("not implemented")
}
