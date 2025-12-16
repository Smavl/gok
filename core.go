package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"golang.org/x/term"
)

type Core struct {
	mu               sync.RWMutex
	Config           Config
	menuReaderCancel context.CancelFunc

	SessionManager       *SessionManager
	ShellListenerManager *ShellListenerManager

	Display Display
	Mode Mode

	// shell
	activeShellID int

	// Event Channels
	eventChan chan Event
}

func NewCore(cfg Config) *Core {
	terminalDisplay := NewTerminalDisplay()
	// NOTE: Do we need a proper constructor?
	menuMode := MenuMode{Display: terminalDisplay}
	sm := NewSessionManager(terminalDisplay)
	eventChan := make(chan Event)
	slm := NewShellListenerManager(sm, terminalDisplay, eventChan)
	core := Core{
		Config: cfg,
		// managers
		SessionManager:       sm,
		ShellListenerManager: slm,

		Display: terminalDisplay,
		Mode:    &menuMode,

		// channels
		eventChan: eventChan,
	}
	core.ShellListenerManager.Init(cfg)
	return &core
}

func (c *Core) EnableShellMode() {
	session, _ := c.SessionManager.Get(c.activeShellID)
	// TODO: error handling

	// Stop the menu reader goroutine so it releases stdin
	c.mu.Lock()
	if c.menuReaderCancel != nil {
		c.menuReaderCancel()
		c.menuReaderCancel = nil
	}
	c.mu.Unlock()

	// Bring session to foreground (redirect buffer output to "stdout")
	session.Foreground()
	oldMenuMode := c.Mode
	c.Mode = &ShellMode{Display: c.Display}

	// drop into shell (blocking)
	c.runShellReader()

	// When session is escaped: background the session
	session.Background()
	c.Mode = oldMenuMode

	// Restart the menu reader goroutine now that we're back in menu mode
	c.startMenuReader()
}



func (c *Core) Prompt() {
	c.Display.Prompt()
}

func (c *Core) Message(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	c.Display.Message(msg)
}


// read handlers

// Menu reader, reads byte-by-byte but is line-buffered
func (c *Core) runMenuReader(ctx context.Context) {
	// This reader is cancellable. 
	var lineBuffer []byte
	buf := make([]byte, 1)

	c.Prompt()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				// NOTE: This will cause the reader to exit if stdin is closed.
				return
			}
			
			// Simple line assembly
			if buf[0] == '\n' {
				c.eventChan <- UserInputEvent{io: string(lineBuffer)}
				lineBuffer = []byte{}
			} else if buf[0] != '\r' { // Ignore carriage returns
				lineBuffer = append(lineBuffer, buf[0])
			}
		}
	}
}

// Shell reader (raw terminal):
func (c *Core) runShellReader() {
	session, _ := c.SessionManager.Get(c.activeShellID)
	// TODO: error handling

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		c.Message("Failed to enter raw mode %v\n", err)
		return
	}
	defer term.Restore(fd, oldState)

	// Enable raw mode output conversion
	if td, ok := c.Display.(*TerminalDisplay); ok {
		td.SetRawMode(true)
		defer td.SetRawMode(false)
	}

	// workaround for resuming session
	session.Write([]byte("\n"))

	// Read byte-by-byte in raw mode
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return
		}
		if n > 0 {
			// Check for escape (Ctrl+D = 0x04)
			if buf[0] == 0x04 {
				c.Message("\r\n")
				return
			}

			// Convert CR to LF for Unix shells
			// FIX: HMMM without this the overal shell doesnt wrap properly (even when exiting???)
			if buf[0] == '\r' {
				buf[0] = '\n'
			}

			// Forward to remote shell
			session.Write(buf[:n])
		}
	}
}

func (c *Core) handleEvent(e Event) {
	switch evt := e.(type) {
	case NewSessionEvent:
		c.Mode.OnNewSession(evt.Session)
	case UserInputEvent:
		c.Mode.HandleInput(evt.io, c)
	}
}

func (c *Core) startMenuReader() {
	c.mu.Lock()
	// Create a new context and cancel function for the reader
	ctx, cancel := context.WithCancel(context.Background())
	c.menuReaderCancel = cancel
	go c.runMenuReader(ctx)
	c.mu.Unlock()
}

func (c *Core) Start() {
	// Start the first instance of the menu reader
	c.startMenuReader()

	// Process events from the single channel
	for e := range c.eventChan {
		c.handleEvent(e)
	}
}
