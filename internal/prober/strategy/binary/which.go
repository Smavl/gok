package binary

import (
	"context"
	"fmt"
	"strings"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/executor"
	"github.com/smavl/gok/internal/prober/types"
)

// WhichStrategy is a implementation of BinaryCheckStrategy
type WhichStrategy struct {
	executor executor.Executor
}

func NewWhichStrategy() *WhichStrategy {
	return &WhichStrategy{
		executor: executor.NewDefaultExecutor(),
	}
}

func PathWasReturned(s string) bool {
	// check if starts with `/`
	hasPrefix := strings.HasPrefix(s, "/")
	// check that it does not start with `which: no`
	hasNegWhichPattern := strings.HasPrefix(s, "which: no")

	return hasPrefix && !hasNegWhichPattern
}

func (s *WhichStrategy) CheckExists(ctx context.Context, sess domain.CommandSession, binary string) (types.BinaryResult, error) {
	// Build the which command - redirect output to suppress noise
	cmd := fmt.Sprintf("which %s 2>&1", binary)

	// Execute and get exit code
	res, err := s.executor.Execute(ctx, sess, cmd)
	if err != nil {
		// something went wrong
		return types.BinaryResult{
			Name:  binary,
			Found: false,
		}, err
	}

	// 
	found := false
	path := ""
	for _, line := range res {
		path = strings.TrimSpace(line)
		if PathWasReturned(path) { found = true; break }

	}

	result := types.BinaryResult{
		Name:  binary,
		Path:  path,
		Found: found,
	}

	// Exit code 0 => binary was found
	return result, nil
}


// exit code base approach (does not return path)
// Uses `which` to determin the existence of a binary
// func (s *WhichStrategy) CheckExists(ctx context.Context, sess domain.SessionInterface, binary string) (bool, error) {
// 	// Build the which command - redirect output to suppress noise
// 	cmd := fmt.Sprintf("which %s >/dev/null 2>&1", binary)
//
// 	// Execute and get exit code
// 	exitCode, err := s.executor.ExecuteWithExitCode(ctx, sess, cmd)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to execute which command: %w", err)
// 	}
//
// 	// Exit code 0 => binary was found
// 	return exitCode == 0, nil
// }
