package main

import (
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

	SessionManager       *SessionManager
	ShellListenerManager *ShellListenerManager

	Display Display

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
	initialState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Fatal: Could not get terminal state: %v\n", err)
		os.Exit(1)
	}
	terminalDisplay := NewTerminalDisplay()
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
	session, err := c.SessionManager.Get(id)
	if err != nil {
		c.Message("[!] Error: Session #%d not found\n", id)
		c.Prompt()
		return
	}

	c.SetActiveShell(id)

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
	session, err := c.SessionManager.Get(c.GetActiveShell())
	if err != nil {
		// Session might have died - still restore terminal
		c.Message("\r\n[!] Warning: Session no longer exists\r\n")
	} else {
		session.Background()
	}

	term.Restore(int(os.Stdin.Fd()), c.terminalState)
	if td, ok := c.Display.(*TerminalDisplay); ok {
		td.SetRawMode(false)
	}

	c.Message("\r\n")

	// Swap back to line reader (non-blocking!)
	c.swapReader(NewLineReader())
	defer c.Prompt()
}


func (c *Core) Prompt() {
	c.Display.Prompt()
}

func (c *Core) Message(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	c.Display.Message(msg)
}

// Thread-safe accessors for activeShellID
func (c *Core) SetActiveShell(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.activeShellID = id
}

func (c *Core) GetActiveShell() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.activeShellID
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
