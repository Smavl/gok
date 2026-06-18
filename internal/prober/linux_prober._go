package prober

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

func NewLinuxProber(sess SessionInterface, probeOpts ProberOptions) *LinuxProber {
	return &LinuxProber{
		Session:    sess,
		Binaries:   make([]string, 0),
		CmdTimeout: probeOpts.CmdTimeout,
	}
}

type LinuxProber struct {
	Session    SessionInterface
	Binaries   []string
	CmdTimeout time.Duration
}


func (prober *LinuxProber) GetBinaries() []string {
	return prober.Binaries
}

// WIP FUNCTION
func (prober *LinuxProber) binariesPresentFast(binaries []string) (map[string]bool, error) {
	session := prober.Session
	session.ClearProbingBuffer()

	delimiter := "¤"
	// whichCmd := fmt.Sprintf("echo \"$(which %s >/dev/null 2>&1; echo $?)\";echo '%s'\n", binary, delimiter)
	var whichCompoundCmds strings.Builder
	whichCompoundCmds .WriteString("echo \"$(")
	for _, binary := range binaries {
		whichCompoundCmds.WriteString("which " + binary + " >/dev/null 2>&1; echo $?; ")
	}
	fmt.Fprintf(&whichCompoundCmds, ")\";echo '%s'\n", delimiter)
	fmt.Printf("Executing compound which command: %s", whichCompoundCmds.String())
	ExecuteCmd(session, []byte(whichCompoundCmds.String()))
	waitForOutputUsingDelimeter(session, delimiter, prober.CmdTimeout)

	output := session.GetProbingLines()
	// Skip first line (always command echo)
	if len(output) > 0 {
		output = output[1:]
	}

	results := make(map[string]bool)
	for _, binary := range binaries {
		fmt.Printf("Checking binary: %s\n", binary)
		whichExitCode, err := getExitCode(output, delimiter)
		if err != nil {
			return nil, err
		}
		if whichExitCode == Success {
			prober.Binaries = append(prober.Binaries, binary)
			results[binary] = true
		} else {
			results[binary] = false
		}
	}
	return results, nil
}

func (prober *LinuxProber) binaryPresentFast(binary string) (bool, error) {
	session := prober.Session
	session.ClearProbingBuffer()

	delimiter := "¤"
	whichCmd := fmt.Sprintf("echo \"$(which %s >/dev/null 2>&1; echo $?)\";echo '%s'\n", binary, delimiter)
	ExecuteCmd(session, []byte(whichCmd))
	waitForOutputUsingDelimeter(session, delimiter, prober.CmdTimeout)

	output := session.GetProbingLines()
	// Skip first line (always command echo)
	if len(output) > 0 {
		output = output[1:]
	}

	whichExitCode, err := getExitCode(output, delimiter)
	if err != nil {
		return false, err
	}
	if whichExitCode == Success {
		prober.Binaries = append(prober.Binaries, binary)
		return true, nil
	}
	return false, nil
}

func (prober *LinuxProber) binaryPresent(binary string) (bool, error) {
	session := prober.Session

	// Should look like:
	//> $ echo "$(which which > /dev/null 2>&1; echo $?)"
	//> 0
	whichCmd := fmt.Sprintf("echo \"$(which %s>/dev/null 2>&1; echo $?)\"\n", binary)
	ExecuteCmd(session, []byte(whichCmd))
	waitForOutput(session.GetProbingDataChannel(), prober.CmdTimeout)

	output := session.GetProbingLines()
	whichExitCode, err := getExitCode(output, "")
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
	session := prober.Session
	session.ClearProbingBuffer()

	delimiter := "¤"
	whichCmd := fmt.Sprintf("echo \"$(which which>/dev/null 2>&1; echo $?)\";echo '%s'\n", delimiter)
	ExecuteCmd(session, []byte(whichCmd))
	waitForOutputUsingDelimeter(session, delimiter, prober.CmdTimeout)

	lines := session.GetProbingLines()

	whichExitCode, err := getExitCode(lines, delimiter)
	if err != nil {
		return
	}

	if whichExitCode == Success {
		prober.Binaries = append(prober.Binaries, "which")
	}
}

func (prober *LinuxProber) EnumerateBinaries() {
	prober.handleWhichEnumeration()
	gotWhich := slices.Contains(prober.Binaries, "which")
	if !gotWhich {
		return
	}

	interestingBinaries := []string{
		// programing languages / interpreters
		"python", "python3", "python2", "perl",
		// Utilities
		"base64", "find", "grep",
		// network tools
		"nc", "netcat", "socat",
		// HTTP tools
		"wget", "curl",
	}
	for _, binary := range interestingBinaries {
		_, err := prober.binaryPresentFast(binary)
		if err != nil {
			continue
		}
	}
	// prober.binariesPresentFast(interestingBinaries)
}
