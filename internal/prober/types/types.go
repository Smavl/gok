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


type ProbeResult interface {
	Apply(results *ProbeResults)
}

type OSResult struct {
	DetectedOS OS 
}

func (r OSResult) Apply(pr *ProbeResults){
	pr.OS = r.DetectedOS
}

type BinariesResult struct {
	Binaries []string
}

func (r BinariesResult) Apply(pr *ProbeResults){
	pr.BinariesFound = append(pr.BinariesFound, r.Binaries...)

	// TODO: maybe update capabilties / derived capabilties here
}

type ProbeResults struct {
	// Results []*ProbeResult
	OS OS 
	BinariesFound []string

	// future results:
	// Users []User
	// Files []File
	// NetworkInfo *NetworkInfo
	// Capabilities Capabilities
}

// type Capabilities struct {
// 	HasWhich
// }

type ProbeOperation func(ctx context.Context, sess SessionInterface) (ProbeResult, error)

type PhaseConfig struct {
	Operations      []ProbeOperation
	TimeoutPerOp	time.Duration
}

type PhaseBuilder func(accResults *ProbeResults) *PhaseConfig

type ProbeConfig struct {
	// Phases map[ProbePhase]PhaseConfig

	// First static phase (e.g. OS detection)
	Genesis PhaseConfig

	// Rest of the phases are dynamically build on intel from prior probes
	Phases map[ProbePhase]PhaseBuilder

}
