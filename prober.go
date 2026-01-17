package main

import (
	"crypto/rand"
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
	waitForOutput(session.probinDataArrived, cmdTimeout)

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

type LinuxProber struct {
	binaries []string
}
