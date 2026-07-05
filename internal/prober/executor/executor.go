package executor

import (
	"context"

	"github.com/smavl/gok/internal/domain"
)

// TODO: Move Executor out of prober dir/package
// Interface for the executor
type Executor interface {
	Execute(ctx context.Context,sess domain.ProbingSession,	cmd string) ([]string, error)
	ExecuteWithExitCode(ctx context.Context, sess domain.ProbingSession, cmd string) (int, error)
}

// CommandExecutor: Execution of commands on a session
func NewDefaultExecutor() Executor {
	return NewRandomDelimExecutor()
}

