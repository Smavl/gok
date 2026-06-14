package prober

import (
	"context"
	"fmt"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/config"
	"github.com/smavl/gok/internal/prober/types"
)

type Prober struct {
	sess    types.SessionInterface
	config  types.ProbeConfig
	results *types.ProbeResults
	done bool
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
		done: false,
	}, nil
}

func (p *Prober) setDone() {
	p.done = true
}

func (p *Prober) IsDone() bool {
	return p.done
}

func newPhaseBuilderContext(p *Prober) types.PhaseBuilderContext {
	return types.PhaseBuilderContext{
		ProbeResults: p.results,
	}
}

func (p *Prober) Run(ctx context.Context) (error) {
	bctx := newPhaseBuilderContext(p)
	phases := []types.ProbePhase{types.PhaseInitial, types.PhaseRecon, types.PhaseDeepScan}
	// TODO: Should a failed run also just be "done"?
	defer p.setDone()

	// Run Genesis phase for os detection
	err := p.runPhase(ctx,p.config.Genesis)
	if err != nil {
		// ERROR: Genesis phase failed
		return err
	}

	// Run each dynamic phase
	for _, phase := range phases {
		phaseBuilder, exists := p.config.Phases[phase]

		if !exists {
			continue
		}

		// build phase
		builtPhase, present := phaseBuilder(bctx)
		if !present {
			// NOTE: error handling???
			// phase was not built 
			continue
		}
		if err := p.runPhase(ctx, *builtPhase); err != nil {
			return err
		}
	}

	return nil
}

func (p *Prober) runPhase(ctx context.Context, cfg types.PhaseConfig) error {
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

		// Apply whatever result respective to the operation
		result.Apply(p.results)

	}

	return nil
}

func (p* Prober) GetProbingResultsIfDone() (*types.ProbeResults, error) {
	if !p.done {
		return nil, fmt.Errorf("probing not completed yet")
	}
	return p.results, nil
}

// GetBinaries returns the list of found binaries 
func (p *Prober) GetBinaries() []string {
	return p.results.BinariesFound
}

