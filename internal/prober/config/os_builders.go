package config

import (

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/types"
)


type OSPhaseBuilder interface {
	BuildInitialPhase(ctx types.PhaseBuilderContext)	(*types.PhaseConfig, bool)
	BuildReconPhase(ctx types.PhaseBuilderContext)		(*types.PhaseConfig, bool)
	BuildDeepScanPhase(ctx types.PhaseBuilderContext)	(*types.PhaseConfig, bool)
}

func newOSPhaseBuilder(os types.OS, mode domain.ProbingMode) (OSPhaseBuilder, error) {
	switch os {
	case types.LinuxOs:
		return newLinuxBuilder(mode), nil
	case types.WindowsOs:
		return newWindowsBuilder(mode), nil
	case types.UnknownOS:
		return &unknownOSPhaseBuilder{}, nil
	default:
		return &unknownOSPhaseBuilder{}, nil
	}
}

// Unknown / none: Does nothing for now

type unknownOSPhaseBuilder struct {}

func (b *unknownOSPhaseBuilder) BuildInitialPhase(ctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	return nil, false
}

func (b *unknownOSPhaseBuilder) BuildReconPhase(ctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	return nil, false
}

func (b *unknownOSPhaseBuilder) BuildDeepScanPhase(ctx types.PhaseBuilderContext) (*types.PhaseConfig, bool) {
	return nil, false
}
