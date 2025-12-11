package main

import (
	"fmt"
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
	mu sync.Mutex
}

func NewTerminalDisplay() *TerminalDisplay {
	return &TerminalDisplay{}
}

func (d *TerminalDisplay) Write(b []byte) (n int, err error){
	d.mu.Lock()
	defer d.mu.Unlock()
	return os.Stdout.Write(b)
}

func (d *TerminalDisplay ) Prompt() {
	d.mu.Lock()
	defer d.mu.Unlock()
	fmt.Printf("GOK > ")
}
func (d *TerminalDisplay) Message(msg string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	fmt.Print(msg)
}
// Error(msg string) 
//
// NewSession(sess *Session)
//
// SessionList(sessions []*Session)
// ListenerList(listeners map[string]*Listener)
//
// EnteringShell(sessionID int)

