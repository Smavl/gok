package operations

import (
	"context"

	"github.com/smavl/gok/internal/prober/strategy"
	"github.com/smavl/gok/internal/prober/types"
)


// Operation for detecting os
// This is run prior to the main prober phases and will determine them based 
// on what os is detected as it should be used to determine what strategies 
// that is relevant in other probing phases
func DetectOS(
	primaryStrategy strategy.OSDetectionStrategy,
	fallbackStrategies []strategy.OSDetectionStrategy,
) types.ProbeOperation {
	return func(ctx context.Context, sess types.SessionInterface) (types.ProbeResult, error) {

		// Run primary strategy 
		OSres, err := primaryStrategy.DetermineOS(ctx, sess)
		// Invert error check to return early on sucess
		if err == nil {
			return types.OSResult{DetectedOS: OSres}, nil
		}

		// Run fallback strategies if primary failed
		for _, fallback := range fallbackStrategies {
			OSres, err = fallback.DetermineOS(ctx, sess)
			if err == nil {
				return types.OSResult{DetectedOS: OSres}, nil
			}
		}

		// If all strategies failed return: res + error from primaryStrategy
		return types.OSResult{DetectedOS: OSres}, err
	}
}
