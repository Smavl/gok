package prober

import (
	"context"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/config"
	"github.com/smavl/gok/internal/prober/types"
)

type Prober struct {
	sess    types.SessionInterface
	config  types.ProbeConfig
	results *types.ProbeResults
}

func NewProber(sess types.SessionInterface, opts domain.ProbingOptions) (*Prober, error) {
	cfg, err := config.ConfigForMode(opts.ProbingMode)
	if err != nil {
		return nil, err
	}
	return &Prober{
		sess:    sess,
		config:  cfg,
		results: &types.ProbeResults{},
	}, nil
}

func (p *Prober) Run(ctx context.Context) (*types.ProbeResults, error) {
	phases := []types.ProbePhase{types.PhaseInitial, types.PhaseRecon, types.PhaseDeepScan}

	// Run each phase
	for _, phase := range phases {
		phaseConfig, exists := p.config.Phases[phase]
		if !exists {
			continue
		}

		if err := p.runPhase(ctx, phase, phaseConfig); err != nil {
			return p.results, err
		}
	}

	return p.results, nil
}

func (p *Prober) runPhase(ctx context.Context, phase types.ProbePhase, cfg types.PhaseConfig) error {
	for _, op := range cfg.Operations {
		// Create timeout context for this operation
		opCtx, cancel := context.WithTimeout(ctx, cfg.TimeoutPerOp)

		// Execute operation
		result, err := op(opCtx, p.sess)
		cancel()

		if err != nil {
			// TODO: decide on error handling strategy (continue vs fail fast)
			// For now, continue on error
			continue
		}

		p.results.Add(result)
	}

	return nil
}

// GetBinaries returns the list of found binaries 
func (p *Prober) GetBinaries() []string {
	if p.results == nil {
		return []string{}
	}

	binaries := []string{}
	for _, result := range p.results.GetResults() {
		if result.Name == "binaries" {
			if binaryResult, ok := result.Data.([]string); ok {
				binaries = append(binaries, binaryResult...)
			}
		}
	}

	return binaries
}

