package types

import (
	"context"
	"time"

	"github.com/smavl/gok/internal/domain"
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

// func ()

type BinaryResults struct {
	Binaries []BinaryResult
}

type BinaryResult struct {
	Name string
	Path string
	Found bool
}

func (r BinaryResults) Apply(pr *ProbeResults){
	pr.BinariesFound.Binaries = append(pr.BinariesFound.Binaries, r.Binaries...)

	// TODO: maybe update capabilties / derived capabilties here
}

type ProbeResults struct {
	// Results []*ProbeResult
	OS OS 
	BinariesFound BinaryResults

	// future results:
	// Users []User
	// Files []File
	// NetworkInfo *NetworkInfo
	// Capabilities Capabilities
}

func (pr *ProbeResults) HasBinary(binaryName string) bool {
	for _, b := range pr.BinariesFound.Binaries {
		if b.Name == binaryName && b.Found {
			return true
		}
	}
	return false
}

// type Capabilities struct {
// 	HasWhich
// }

type BinaryCapability struct {
	Name string
	Path string
	// Version string
}

type ProbeOperation func(ctx context.Context, sess SessionInterface) (ProbeResult, error)

type PhaseConfig struct {
	Operations      []ProbeOperation
	TimeoutPerOp	time.Duration
}

type PhaseBuilderContext struct {
	ProbeResults	*ProbeResults
	Mode			domain.ProbingMode
	// Capabilities Capabilities
}


type PhaseBuilder func(PhaseBuilderContext) (*PhaseConfig, bool)

type ProbeConfig struct {
	// First static phase (e.g. OS detection)
	Genesis PhaseConfig

	// Rest of the phases are dynamically build based on intel from prior probes
	Phases map[ProbePhase]PhaseBuilder

}
