package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Core struct {
	mu     sync.RWMutex
	Config Config

	// Managers
	SessionManager       *SessionManager
	ShellListenerManager *ShellListenerManager
	InputMan             *InputManagerImpl
	commander            CommandHandler
	shellMode            *RawShellMode
	terminal             *Terminal

	eventChan chan Event
}

func NewCore(cfg Config) *Core {
	terminalDisplay := NewTerminalDisplay()
	terminal, err := NewTerminal(terminalDisplay)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Fatal: Could not setup terminal: %v\n", err)
		os.Exit(1)
	}

	eventChan := make(chan Event)
	inputMan := NewInputManager(NewLineReader(), eventChan)
	sm := NewSessionManager()
	slm := NewShellListenerManager(sm, terminal, eventChan)
	shellMode := NewRawShellMode(sm, inputMan, terminal)

	core := &Core{
		Config:               cfg,
		SessionManager:       sm,
		ShellListenerManager: slm,
		terminal:             terminal,
		eventChan:            eventChan,
		InputMan:             inputMan,
		shellMode:            shellMode,
		commander:            *NewCommandHandler(sm, slm, terminal, shellMode),
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		core.terminal.Restore()
		os.Exit(0)
	}()

	core.ShellListenerManager.Init(context.Background(), cfg)
	return core
}

func (c *Core) Prompt() {
	c.terminal.Prompt()
}

func (c *Core) Message(format string, a ...any) {
	c.terminal.Message(format, a...)
}

func (c *Core) Start() {
	// Show initial prompt
	c.Prompt()

	// Start inputReader
	c.InputMan.Start(context.Background())

	// event loop
	for {
		select {
		case event := <-c.eventChan:
			event.Handle(c)
			// does the above handle NewSessionEvent?
		}
	}
}
