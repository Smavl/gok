package strategy

import (
	"context"

	"github.com/smavl/gok/internal/prober/types"
)

// In this file we have the strategy interfaces for the prober

// BinaryCheckStrategy: Techniques to check the existence of a binary 
type BinaryCheckStrategy interface {
	CheckExists(ctx context.Context, sess types.SessionInterface, binary string) (bool, error)
}

type OSDetectionStrategy interface {
	DetermineOS(ctx context.Context, sess types.SessionInterface) (types.OS, error)	
}
