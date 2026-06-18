package executor

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/smavl/gok/internal/prober/types"
)

// CommandExecutor: Execution of commands on a session
type CommandExecutor struct{}

// TODO: Add strategies for the 
func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

// ExecuteWithExitCode executes a command and returns its exit code
func (e *CommandExecutor) ExecuteWithExitCode(ctx context.Context, sess types.SessionInterface, cmd string) (int, error) {
	sess.ClearProbingBuffer()

	delimiter := generateDelimiter()

	// Wrap command to capture exit code
	wrappedCmd := fmt.Sprintf("%s; echo $?; echo '%s'\n", cmd, delimiter)

	_, err := sess.Write([]byte(wrappedCmd))
	if err != nil {
		return -1, fmt.Errorf("failed to write command: %w", err)
	}

	// Wait for completion
	if !e.waitForDelimiter(ctx, sess, delimiter) {
		return -1, fmt.Errorf("timeout waiting for command completion")
	}

	// Parse exit code
	lines := sess.GetProbingLines()
	exitCode, err := parseExitCode(lines, delimiter)
	if err != nil {
		return -1, fmt.Errorf("failed to parse exit code: %w", err)
	}

	return exitCode, nil
}

// Execute runs a command and returns the output lines
func (e *CommandExecutor) Execute(ctx context.Context, sess types.SessionInterface, cmd string) ([]string, error) {
	sess.ClearProbingBuffer()

	delimiter := generateDelimiter()

	// Wrap command with delimiter
	wrappedCmd := fmt.Sprintf("%s; echo '%s'\n", cmd, delimiter)

	_, err := sess.Write([]byte(wrappedCmd))
	if err != nil {
		return nil, fmt.Errorf("failed to write command: %w", err)
	}

	// Wait for completion
	if !e.waitForDelimiter(ctx, sess, delimiter) {
		return nil, fmt.Errorf("timeout waiting for command completion")
	}

	// Get output lines (excluding delimiter line)
	lines := sess.GetProbingLines()
	return filterDelimiterLine(lines, delimiter), nil
}

// waitForDelimiter waits for the delimiter to appear in output
func (e *CommandExecutor) waitForDelimiter(ctx context.Context, sess types.SessionInterface, delimiter string) bool {
	dataChannel := sess.GetProbingDataChannel()

	for {
		select {
		case <-dataChannel:
			lines := sess.GetProbingLines()
			if containsDelimiter(lines, delimiter) {
				return true
			}
		case <-ctx.Done():
			return false
		}
	}
}

// generateDelimiter creates a unique delimiter for command execution
func generateDelimiter() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return "¤" + hex.EncodeToString(bytes) + "¤"
}

// containsDelimiter checks if delimiter is present in any line
func containsDelimiter(lines []string, delimiter string) bool {
	for _, line := range lines {
		if strings.Contains(line, delimiter) {
			return true
		}
	}
	return false
}

// parseExitCode extracts exit code from command output
func parseExitCode(lines []string, delimiter string) (int, error) {
	if len(lines) < 2 {
		return -1, fmt.Errorf("insufficient output lines")
	}

	// Find delimiter line
	delimiterIdx := -1
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], delimiter) {
			delimiterIdx = i
			break
		}
	}

	if delimiterIdx == -1 {
		return -1, fmt.Errorf("delimiter not found")
	}

	if delimiterIdx < 1 {
		return -1, fmt.Errorf("no exit code line")
	}

	// Exit code is on the line before delimiter
	exitCodeLine := strings.TrimSpace(lines[delimiterIdx-1])

	var exitCode int
	_, err := fmt.Sscanf(exitCodeLine, "%d", &exitCode)
	if err != nil {
		return -1, fmt.Errorf("failed to parse exit code '%s': %w", exitCodeLine, err)
	}

	return exitCode, nil
}

// filterDelimiterLine removes the delimiter line from output
func filterDelimiterLine(lines []string, delimiter string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if !strings.Contains(line, delimiter) {
			result = append(result, line)
		}
	}
	return result
}
