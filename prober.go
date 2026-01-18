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
	Linux OS = iota
	Unknown
)

const defaultTimeout = 200*time.Millisecond
const TestingTimeout = 2*time.Millisecond

var cmdTimeout = defaultTimeout

type ExitCode int
const (
	// 0 => Success
	Success ExitCode = iota
	// 127 => "Command not found: The command is not recognized or available in the environment’s PATH."
	CommandNotFound = 127
	ExitStatusOutOfRange = 255
)

// type Binary int
//
// const (
// 	which binary = iota
// )

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

func (o OS) GetProber(session *Session) (Prober, error) {
	switch o {
	case Linux: 
		return NewLinuxProber(session), nil
	default: return nil, NoProberForOs 
	}
}

type DetermineOSStrategy interface {
	Determine() OS
}

type RandomCommandStrategy struct {
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

func (rs *RandomCommandStrategy) DetermineOS(session *Session) (OS, error) {
	session.ClearProbingBuffer()

	// five random letters
	rCmd := rand.Text()[:5] + "\n"

	// execute random command
	ExectuteCmd(session, []byte(rCmd))

	// wait for the buffer to populate
	waitForOutput(session.probingDataArrived, cmdTimeout)

	// capture outout and determine os
	output := session.GetProbingLines()

	// DEBUG
	//    session.display.Write([]byte("Detected Output:\n"))
	//    for i, line := range output {
	// session.display.Write([]byte(fmt.Sprint(i) +" : "+ line))
	//    }

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

type Prober interface {
	EnumerateBinaries()
	// EnumerateUser()
	// EnumerateUsers()
}

func NewLinuxProber(session *Session) *LinuxProber {
	return &LinuxProber{
		Session: session,
		binaries: make([]string, 1),
	}
}

type LinuxProber struct {
	Session *Session
	binaries []string
}

func getExitCode(session *Session) (ExitCode, error) {
	lastLine := session.GetProbingLines()[len(session.GetProbingLines())-1]
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

func (prober *LinuxProber) binaryPresent(binary string) (bool,error) {
	session := prober.Session

	whichCmd := fmt.Sprintf("which %s;echo $?\n", binary)
	ExectuteCmd(session, []byte(whichCmd))
	waitForOutput(session.probingDataArrived, cmdTimeout)

	whichExitCode,err := getExitCode(session)
	if err != nil {
		return false, err
	}
	if whichExitCode == Success {
		prober.binaries = append(prober.binaries, binary)
		return true, nil
	}
	return false, nil
}

func (prober *LinuxProber) handleWhichEnumeration() {
	// Check for which
	session := prober.Session
	whichCmd := "which which;echo $?\n"
	ExectuteCmd(session, []byte(whichCmd))
	waitForOutput(session.probingDataArrived, cmdTimeout)
	// output := session.GetProbingLines()
	// DEBUG:
	// session.display.Write([]byte("Detected Output:\n"))
	// for i, line := range output {
	// 	session.display.Write([]byte(fmt.Sprint(i) +" : "+ line))
	// }
	whichExitCode,err := getExitCode(session)
	if err != nil {
		// session.display.Write([]byte("Error getting exit code for which command\n"))
		return
	}

	if whichExitCode == Success {
		prober.binaries = append(prober.binaries, "which")
	}
}



func (prober *LinuxProber) EnumerateBinaries() {
	session := prober.Session
	session.ClearProbingBuffer()

	// Check for which
	prober.handleWhichEnumeration()
	gotWhich := slices.Contains(prober.binaries,"which")
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
			session.display.Write(fmt.Appendf(nil, "Binary found: %s\n", binary))
		}
		
	}

}
