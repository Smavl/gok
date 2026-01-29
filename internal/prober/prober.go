package prober

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/misc"
)

// SessionInterface defines the minimal interface of functions that prober needs from a session
type SessionInterface interface {
	Write([]byte) (int, error)
	GetProbingLines() []string
	ClearProbingBuffer()
	GetProbingDataChannel() <-chan struct{}
}

type OS int

const (
	Unknown OS = iota
	Linux
)

type BashExitCode int
const (
	// 0 => Success
	Success BashExitCode = 0
	// 127 => "Command not found: The command is not recognized or available in the environment’s PATH."
	CommandNotFound = 127
	// 255 => "Exit status out of range: Typically, this happens when a script or command exits with a number > 255"
	ExitStatusOutOfRange = 255
)

func (o OS) String() string {
	switch o {
	case Linux:
		return "Linux"
	case Unknown:
		return "Unknown OS"
	default:
		return "Invalid"
	}
}

func (e BashExitCode) String() string {
	switch e {
	case Success:
		return "Success"
	case CommandNotFound:
		return "Command Not Found"
	case ExitStatusOutOfRange:
		return "Exit Status Out Of Range"
	default:
		return "Unknown Exit Code"
	}
}

// TODO: Rename to GetNewProber
func (o OS) GetNewProber(sess SessionInterface, probeOpts ProberOptions) (Prober, error) {
	switch o {
	case Linux:
		return NewLinuxProber(sess, probeOpts), nil
	default: return nil, misc.NoProberForOs
	}
}

type DetermineOSStrategy interface {
	Determine() OS
}

type RandomCommandStrategy struct {
	CmdTimeout time.Duration
}

func ExecuteCmd(sess SessionInterface, cmd []byte) {
	sess.Write([]byte(cmd))
}

func waitForOutput(dataArrived <-chan struct{}, timeout time.Duration) {
	for {
		select {
		case <-dataArrived:
		case <-time.After(timeout):
			return
		}
	}
}
func waitForOutputUsingDelimeter(sess SessionInterface, delimiter string, timeout time.Duration) {
	// Make cancelable
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-sess.GetProbingDataChannel():
			// check for delimiter to minimize wait time
			if endDelimiterFound(sess, delimiter) {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func endDelimiterFound(sess SessionInterface, delimiter string) bool {
	lines := sess.GetProbingLines()
	if len(lines) == 0 {
		return false
	}
	joined := strings.Join(lines, "")
	return strings.Contains(joined, delimiter)
}

// supports usage like:
// hasErrorPattern(output, "not recognized", "is not recognized"):
func hasErrorPattern(lines []string, patterns ...string) bool {
	for _, line := range lines {
		for _, pattern := range patterns {
			if strings.Contains(line, strings.ToLower(pattern)) {
				return true
			}
		}
	}
	return false
}

func inferOsByError(output []string) (OS, error) {
	var resultOS OS
	var resultErr error

	switch {
	// "bash: rCmd: command not found"
	case hasErrorPattern(output, "command not found"):
		resultOS, resultErr = Linux, nil
	default:
		resultOS, resultErr = Unknown, misc.CouldNotDetermineOSError
	}

	return resultOS, resultErr
}

func (rs *RandomCommandStrategy) DetermineOS(sess SessionInterface) (OS, error) {
	sess.ClearProbingBuffer()

	// five random letters
	rCmd := rand.Text()[:5] + "\n"

	// execute random command
	ExecuteCmd(sess, []byte(rCmd))
	// wait for the response (buffer to populate)
	waitForOutput(sess.GetProbingDataChannel(), rs.CmdTimeout)

	// capture outout and determine os
	output := sess.GetProbingLines()

	return inferOsByError(output)
}

type Prober interface {
	EnumerateBinaries()
	GetBinaries() []string
	// EnumerateUser()
	// EnumerateUsers()
}

type ProberOptions struct {
	CmdTimeout time.Duration
	ProbingMode domain.ProbingMode
}

func getExitCode(output []string, delimiter string) (BashExitCode, error) {

	if len(output) == 0 {
		return 0, fmt.Errorf("no output received")
	}
	// getIndex of line with delimiter
	var idxDelim int
	delimiterFound := false
	for i := len(output) -1 ; i >=0 ; i-- {
		if strings.Contains(output[i], delimiter) {
			idxDelim = i
			delimiterFound = true
			break
		}
	}

	// Check if delimiter was found and there's a line before it
	if !delimiterFound {
		return 0, fmt.Errorf("delimiter not found in output!")
	}
	if idxDelim == 0 {
		return 0, fmt.Errorf("no exit code line before delimiter")
	}

	exitCodeLine := output[idxDelim-1]
	// convert to int - Atoi
	s := strings.TrimSpace(exitCodeLine)
	// replace delimiter
	s = strings.ReplaceAll(s, delimiter, "")
	codeInt, err := strconv.Atoi(s)
	if err != nil {
		// Could not convert to int
		return 0, err
	}
	// cast integer to ExitCode
	return BashExitCode(codeInt), nil
}

