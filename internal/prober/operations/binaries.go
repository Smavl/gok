package operations

import (
	"context"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/strategy"
	"github.com/smavl/gok/internal/prober/types"
)

// EnumerateBinaries creates a ProbeOperation that checks for the existence of specified binaries
// It uses a primary strategy and falls back to alternative strategies on failure
func EnumerateBinaries(
	binaries []string,
	primaryStrategy strategy.BinaryCheckStrategy,
	fallbackStrategies []strategy.BinaryCheckStrategy,
) types.ProbeOperation {
	return func(ctx context.Context, sess domain.CommandSession) (types.ProbeResult, error) {
		found := types.BinaryResults{}

		for _, binary := range binaries {
			// Try primary strategy
			result, err := primaryStrategy.CheckExists(ctx, sess, binary)

			// If primary strategy failed, try fallbacks
			if err != nil {
				for _, fallback := range fallbackStrategies {
					result, err = fallback.CheckExists(ctx, sess, binary)
					if err == nil {
						// Fallback succeeded, stop trying
						break
					}
				}
			}

			// If we still have an error after all strategies, skip this binary
			// (could also return error here depending on desired behavior)
			if err != nil {
				continue
			}

			// Binary exists, add to found list
			if result.Found {
				found.Binaries = append(found.Binaries, result)
			}
		}

		return found, nil
	}
}
