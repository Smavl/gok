package main

import (
	"crypto/rand"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

type OS int

const (
	Unknown OS = iota
	Linux 
)

// var cmdTimeout = defaultTimeout

type ExitCode int
const (
	// 0 => Success
	Success ExitCode = 0
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

func (e ExitCode) String() string {
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
func (o OS) GetProber(session *Session, probeOpts ProberOptions) (Prober, error) {
	switch o {
	case Linux: 
		return NewLinuxProber(session, probeOpts), nil
	default: return nil, NoProberForOs 
	}
}

type DetermineOSStrategy interface {
	Determine() OS
}

type RandomCommandStrategy struct {
	cmdTimeout time.Duration
}

func ExectuteCmd(session *Session, cmd []byte) {
	session.Write([]byte(cmd))
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

// supports usage like:
// hasErrorPattern(output, "not recognized", "is not recognized"):
func hasErrorPattern(lines []string, patterns ...string) bool {
	for _, line := range lines {
		// lowerLine := strings.ToLower(line)
		for _, pattern := range patterns {
			// if strings.Contains(lowerLine, strings.ToLower(pattern)) {
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
		resultOS, resultErr = Unknown, CouldNotDetermineOSError
	}

	return resultOS, resultErr
}

func (rs *RandomCommandStrategy) DetermineOS(session *Session) (OS, error) {
	// DEBUG:
	fmt.Printf("DEBUG: DetermineOS - cmdTimeout is %v\n", rs.cmdTimeout)
	session.ClearProbingBuffer()

	// five random letters
	rCmd := rand.Text()[:5] + "\n"

	// execute random command
	ExectuteCmd(session, []byte(rCmd))
	// wait for the response (buffer to populate)
	waitForOutput(session.probingDataArrived, rs.cmdTimeout)

	// capture outout and determine os
	output := session.GetProbingLines()
	fmt.Printf("DEBUG: Probing Output: %q\n", output) // DEBUG

	return inferOsByError(output)
}

type Prober interface {
	EnumerateBinaries()
	getBinaries() []string
	// EnumerateUser()
	// EnumerateUsers()
}

type ProberOptions struct {
	cmdTimeout time.Duration
}

func NewLinuxProber(session *Session, probeOpts ProberOptions) *LinuxProber {
	return &LinuxProber{
		Session: session,
		Binaries: make([]string, 1),
		cmdTimeout: probeOpts.cmdTimeout,
	}
}

type LinuxProber struct {
	Session *Session
	Binaries []string
	cmdTimeout time.Duration
}

func getExitCode(output []string) (ExitCode, error) {
	if len(output) == 0 {
		return 0, fmt.Errorf("no output received")
	}
	lastLine := output[len(output)-1]
	// convert to int - Atoi
	s := strings.TrimSpace(lastLine)
	codeInt, err := strconv.Atoi(s)
	if err != nil {
		// Could not convert to int
		return 0, err
	}
	// cast integer to ExitCode 
	return ExitCode(codeInt), nil
}

func (prober *LinuxProber) getBinaries() []string {
	return prober.Binaries
}

func (prober *LinuxProber) binaryPresent(binary string) (bool,error) {
	session := prober.Session

	whichCmd := fmt.Sprintf("which %s;echo $?\n", binary)
	ExectuteCmd(session, []byte(whichCmd))
	waitForOutput(session.probingDataArrived, prober.cmdTimeout)

	output := session.GetProbingLines()
	whichExitCode,err := getExitCode(output)
	if err != nil {
		return false, err
	}
	if whichExitCode == Success {
		prober.Binaries = append(prober.Binaries, binary)
		return true, nil
	}
	return false, nil
}

func (prober *LinuxProber) handleWhichEnumeration() {
	// Check for which
	session := prober.Session
	fmt.Printf("DEBUG: cmdTimeout is %v\n", prober.cmdTimeout) // DEBUG
	whichCmd := "which which;echo $?\n"
	ExectuteCmd(session, []byte(whichCmd))
	waitForOutput(session.probingDataArrived, prober.cmdTimeout)

	lines := session.GetProbingLines()
	fmt.Printf("DEBUG: Probing Output: %q\n", lines) // DEBUG

	whichExitCode,err := getExitCode(lines)
	if err != nil {
		fmt.Printf("DEBUG: getExitCode error: %v\n", err) // DEBUG
		// TODO: Return error?
		return
	}

	if whichExitCode == Success {
		prober.Binaries = append(prober.Binaries, "which")
	}
}



func (prober *LinuxProber) EnumerateBinaries() {
	session := prober.Session
	session.ClearProbingBuffer()

	// Check for which
	prober.handleWhichEnumeration()
	gotWhich := slices.Contains(prober.Binaries,"which")
	if !gotWhich {
		return
	}

	interestingBinaries := []string{
		"python","python3","python2",
		"perl",
		"wget","curl",
		"nc","netcat","socat",
	}
	for _, binary := range interestingBinaries {
		binaryPresent, err := prober.binaryPresent(binary)
		if err != nil {
			// session.display.Write([]byte(fmt.Sprintf("Error checking for binary %s: %v\n", binary, err)))
			continue
		}
		if binaryPresent {
			// debug:
			session.display.Write(fmt.Appendf(nil, "Binary found: %s\n", binary))
		}
	}

}
