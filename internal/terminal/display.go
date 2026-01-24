package terminal

import (
	"io"
	"os"
	"sync"
)

type Display interface {
	io.Writer
}

type TerminalDisplay struct {
	mu sync.Mutex
}

func NewTerminalDisplay() *TerminalDisplay {
	return &TerminalDisplay{}
}

func (d *TerminalDisplay) Write(b []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return os.Stdout.Write(b)
}
