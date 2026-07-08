package strategy

import (
	"context"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/types"
)

// In this file we have the strategy interfaces for the prober

// BinaryCheckStrategy: Techniques to check the existence of a binary
type BinaryCheckStrategy interface {
	CheckExists(ctx context.Context, sess domain.CommandSession, binary string) (types.BinaryResult, error)
}

type OSDetectionStrategy interface {
	DetermineOS(ctx context.Context, sess domain.CommandSession) (types.OS, error)	
}
