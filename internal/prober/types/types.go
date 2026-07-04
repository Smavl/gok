package types

import (
	"context"
	"fmt"
	"time"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/misc"
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

type EnvResult struct {
	BashPath string
	// Variables map[string]string
}

func (br BinaryResult) GetPath() (string, error)  {
	if br.Found {
		return br.Path, nil
	}
	return "", fmt.Errorf("binary %s not found", br.Name)
}

func (r BinaryResults) Apply(pr *ProbeResults){
	pr.BinariesResults.Binaries = append(pr.BinariesResults.Binaries, r.Binaries...)

	// TODO: maybe update capabilties / derived capabilties here
}

type ProbeResults struct {
	// Results []*ProbeResult
	OS OS 
	BinariesResults BinaryResults
	EnvResults EnvResult

	// future results:
	// Users []User
	// Files []File
	// NetworkInfo *NetworkInfo
	// Capabilities Capabilities
}

// TODO: move
func (pr *ProbeResults) HasBinary(binaryName string) bool {
	for _, b := range pr.BinariesResults.Binaries {
		if b.Name == binaryName && b.Found {
			return true
		}
	}
	return false
}

func (pr *ProbeResults) GetBinary(binaryName string) (BinaryResult, error) {
	for _, b := range pr.BinariesResults.Binaries {
		if b.Name == binaryName && b.Found {
			return b, nil
		}
	}
	return BinaryResult{}, misc.ErrBinaryNotFound
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
