package binary

import (
	"context"
	"fmt"

	"github.com/smavl/gok/internal/prober/executor"
	"github.com/smavl/gok/internal/prober/types"
)

// WhichStrategy is a implementation of BinaryCheckStrategy
type WhichStrategy struct {
	executor *executor.CommandExecutor
}

func NewWhichStrategy() *WhichStrategy {
	return &WhichStrategy{
		executor: executor.NewCommandExecutor(),
	}
}

// Uses `which` to determin the existence of a binary
func (s *WhichStrategy) CheckExists(ctx context.Context, sess types.SessionInterface, binary string) (bool, error) {
	// Build the which command - redirect output to suppress noise
	cmd := fmt.Sprintf("which %s >/dev/null 2>&1", binary)

	// Execute and get exit code
	exitCode, err := s.executor.ExecuteWithExitCode(ctx, sess, cmd)
	if err != nil {
		return false, fmt.Errorf("failed to execute which command: %w", err)
	}

	// Exit code 0 => binary was found
	return exitCode == 0, nil
}
