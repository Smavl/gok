package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/term"
)

type Core struct {
	mu               sync.RWMutex
	Config           Config
	menuReaderCancel context.CancelFunc

	SessionManager       *SessionManager
	ShellListenerManager *ShellListenerManager

	Display Display
	// Mode Mode
	// for restoring
	terminalState *term.State	

	// shell
	activeShellID int

	//input management 
	inputReader InputReader
	eventChan chan Event
	stopChan chan struct{}

}

func NewCore(cfg Config) *Core {
	initialState, _ := term.GetState(int(os.Stdin.Fd()))
	terminalDisplay := NewTerminalDisplay()
	// NOTE: Do we need a proper constructor?
	// menuMode := MenuMode{Display: terminalDisplay}
	sm := NewSessionManager(terminalDisplay)
	eventChan := make(chan Event)
	slm := NewShellListenerManager(sm, terminalDisplay, eventChan)
	core := Core{
		Config: cfg,
		// managers
		SessionManager:       sm,
		ShellListenerManager: slm,

		Display: terminalDisplay,
		// Mode:    &menuMode,

		// channels
		eventChan: eventChan,
		inputReader: NewLineReader(),
		terminalState: initialState,	
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		term.Restore(int(os.Stdin.Fd()), core.terminalState)
		os.Exit(0)
	}()


	core.ShellListenerManager.Init(cfg)
	return &core
}

func (c *Core) swapReader(newReader InputReader) {
	// Stop the old goroutine
	c.inputReader.Stop()
	close(c.stopChan)

	// Swap to new reader
	c.inputReader = newReader

	// Start new goroutine
	c.startInputReader()
}

func (c *Core) EnterShell(id int) {
	c.activeShellID = id
	session, _ := c.SessionManager.Get(id)

	// raw mode
	term.MakeRaw(int(os.Stdin.Fd()))
	if td, ok := c.Display.(*TerminalDisplay); ok {
		td.SetRawMode(true)
	}

	session.Foreground()

	// Swap to byte reader (shell reader)
	c.swapReader(NewByteReader())
}

func (c *Core) ExitShell() {
	session, _ := c.SessionManager.Get(c.activeShellID)
	session.Background()

	term.Restore(int(os.Stdin.Fd()), c.terminalState)
	if td, ok := c.Display.(*TerminalDisplay); ok {
		td.SetRawMode(false)
	}

	c.Message("\r\n")

	// Swap back to line reader (non-blocking!)
	c.swapReader(NewLineReader())
}


func (c *Core) Prompt() {
	c.Display.Prompt()
}

func (c *Core) Message(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	c.Display.Message(msg)
}



func (c *Core) startInputReader() {
	c.stopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-c.stopChan:
				return  // Exit goroutine cleanly
			default:
				event := c.inputReader.Read()
				if event != nil {
					select {
					case c.eventChan <- event:
					case <-c.stopChan:
						return
					}
				}
			}
		}
	}()
}

func (c *Core) Start() {
	// Show initial prompt
	c.Prompt()

	// Start inputReader
	c.startInputReader()

	// event loop
	for {
		select {
		case event := <- c.eventChan:
			event.Handle(c)
		// does the above handle NewSessionEvent?
		}
	}
}
