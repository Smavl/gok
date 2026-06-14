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

		OS, err := primaryStrategy.DetermineOS(ctx, sess)

		return types.OSResult{DetectedOS: OS}, err
	}
}
