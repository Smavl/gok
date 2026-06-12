package types

import (
	"context"
	"time"
)

// SessionInterface defines the minimal interface that prober needs from a session
type SessionInterface interface {
	Write([]byte) (int, error)
	GetProbingLines() []string
	ClearProbingBuffer()
	GetProbingDataChannel() <-chan struct{}
}

// OS represents the detected operating system
type OS int

const (
	UnknownOS OS = iota
	LinuxOs
	WindowsOs
)

func (o OS) String() string {
	switch o {
	case LinuxOs:
		return "Linux"
	case WindowsOs:
		return "Windows"
	case UnknownOS:
		return "Unknown OS"
	default:
		return "Invalid"
	}
}

type ProbePhase int
const (
	PhaseInitial ProbePhase = iota
	PhaseRecon
	PhaseDeepScan
)


type ProbeResult struct {
	Name string
	Data interface{}
}

type ProbeResults struct {
	Results []*ProbeResult
}

func (pr *ProbeResults) Add(result *ProbeResult) {
	if result != nil {
		pr.Results = append(pr.Results, result)
	}
}

func (pr *ProbeResults) GetResults() []*ProbeResult {
	return pr.Results
}

type ProbeOperation func(ctx context.Context, sess SessionInterface) (*ProbeResult, error)

type PhaseConfig struct {
	Operations      []ProbeOperation
	TimeoutPerOp	time.Duration
}

type ProbeConfig struct {
	Phases map[ProbePhase]PhaseConfig
}
