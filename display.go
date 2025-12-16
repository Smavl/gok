package main

import (
	// "fmt"
	"io"
	"os"
	"sync"
)

type Display interface {
	io.Writer 

	Prompt() // for "GOK>"
	Message(msg string) // Event: messages
	// Error(msg string) // Event: errors
	//
	// // Event notifications
	// NewSession(sess *Session)
	//
	// // Menu things
	// // Lists
	// SessionList(sessions []*Session)
	// ListenerList(listeners map[string]*Listener)
	//
	// // Shell interaction
	// EnteringShell(sessionID int)
}



// Handles all terminal output (Not to be confused with mode)
type TerminalDisplay struct {
	mu      sync.Mutex
	rawMode bool
}

func NewTerminalDisplay() *TerminalDisplay {
	return &TerminalDisplay{}
}

func (d *TerminalDisplay) SetRawMode(enabled bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.rawMode = enabled
}

func (d *TerminalDisplay) Write(b []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// In raw mode, convert LF to CRLF for proper display
	if d.rawMode {
		converted := make([]byte, 0, len(b)*2)
		for _, ch := range b {
			if ch == '\n' {
				converted = append(converted, '\r', '\n')
			} else {
				converted = append(converted, ch)
			}
		}
		return os.Stdout.Write(converted)
	}

	return os.Stdout.Write(b)
}

func (d *TerminalDisplay ) Prompt() {
	d.Write([]byte("GOK > "))
}
func (d *TerminalDisplay) Message(msg string) {
	d.Write([]byte(msg))
}
// Error(msg string) 
//
// NewSession(sess *Session)
//
// SessionList(sessions []*Session)
// ListenerList(listeners map[string]*Listener)
//
// EnteringShell(sessionID int)

