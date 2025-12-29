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

	// Handlers
	commander            *CommandHandler
	shellMode            *RawShellMode
	terminal             *Terminal

	// event
	eventChan chan Event
	ctx context.Context
	cancel context.CancelFunc
	wg sync.WaitGroup
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
		commander:            NewCommandHandler(sm, slm, terminal, shellMode),
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		core.Shutdown()
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
	c.ctx, c.cancel = context.WithCancel(context.Background())
	// Show initial prompt
	c.Prompt()

	c.InputMan.Start(c.ctx)

	// event loop
	for {
		select {
		case event := <-c.eventChan:
			event.Handle(c)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Core) Shutdown() {
	c.terminal.Message("\n[*] Shutting down gracefully...\n")

	// Cancel all contexts
	if c.cancel != nil {
		c.cancel()
	}

	c.InputMan.Stop()

	for _, sess := range c.SessionManager.GetSessions() {
		sess.Stop()
	}

	c.terminal.Restore()
	os.Exit(0)
}
