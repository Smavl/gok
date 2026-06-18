package os

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/smavl/gok/internal/prober/executor"
	"github.com/smavl/gok/internal/prober/types"
)

// OSError is an implementation of OSDetectionStrategy
// Which relies on error messages to determine the OS
type OSErrorDetectionStrategy struct {
	executor *executor.CommandExecutor
}

func NewOSErrorDetectionStrategy() *OSErrorDetectionStrategy {
	return &OSErrorDetectionStrategy{
		executor: executor.NewCommandExecutor(),
	}
}

// supports usage like:
// hasErrorPattern(output, "not recognized", "is not recognized"):
func hasErrorPattern(lines []string, patterns ...string) bool {
	// join strings and normalize to lower
	input := strings.ToLower(strings.Join(lines, "\n"))
	
	for _, pattern := range patterns {
		if strings.Contains(input, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func inferOsByError(output []string) (types.OS, error) {
	var resultOS types.OS
	var resultErr error

	switch {
	// "bash: rCmd: command not found"
	case hasErrorPattern(output, "command not found", "not found:"):
		resultOS, resultErr = types.LinuxOs, nil
	case hasErrorPattern(output, "is not recognized as an internal or external command", "not recognized"):
		resultOS, resultErr = types.WindowsOs, nil
	default:
		resultOS, resultErr = types.UnknownOS, fmt.Errorf("unable to determine OS from error message: %v", output)
	}

	return resultOS, resultErr
}

func (s *OSErrorDetectionStrategy) DetermineOS(ctx context.Context, sess types.SessionInterface) (types.OS, error) {
	// random shell command / builtin / binary
	rcmd := rand.Text()[:8]
	cmd := fmt.Sprintf("%s 2>&1", rcmd)

	output, err := s.executor.Execute(ctx, sess, cmd)
	if err != nil {
		return types.UnknownOS, fmt.Errorf("failed to execute OS command: %w", err)
	}


	OS, err := inferOsByError(output)
	if err != nil {
		return types.UnknownOS, err
	}

	return OS, nil
}
