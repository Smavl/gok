package core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/smavl/gok/internal/cli"
	"github.com/smavl/gok/internal/prober"
	"github.com/smavl/gok/internal/session"
	"github.com/smavl/gok/internal/terminal"
)

// sessionManagerAdapter adapts session.SessionManager to terminal.SessionManager interface
type sessionManagerAdapter struct {
	sm *session.SessionManager
}

func (sma *sessionManagerAdapter) Get(id int) (any, error) {
	return sma.sm.Get(id)
}

type Core struct {
	mu     sync.RWMutex
	Config cli.Config

	// Managers
	SessionManager       *session.SessionManager
	ShellListenerManager *session.ShellListenerManager
	InputMan             *terminal.InputManagerImpl

	// Handlers
	commander *CommandHandler
	shellMode *terminal.RawShellMode
	terminal  terminal.TerminalController

	// event
	eventChan chan any
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewCore(cfg cli.Config) *Core {
	terminalDisplay := terminal.NewTerminalDisplay()

	var term terminal.TerminalController
	var err error

	// Hacky headless switch
	if cfg.HeadlessMode {
		term, err = terminal.NewHeadlessTerminal(terminalDisplay)
	} else {
		term, err = terminal.NewTerminal(terminalDisplay)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Fatal: Could not setup terminal: %v\n", err)
		os.Exit(1)
	}

	eventChan := make(chan any)
	inputMan := terminal.NewInputManager(terminal.NewLineReader(), eventChan)
	sm := session.NewSessionManager(prober.ProberOptions{
		CmdTimeout: cfg.ProbingCmdTimeout,
	})
	slm := session.NewShellListenerManager(sm, term, eventChan)
	// Create adapter for terminal.SessionManager interface
	smAdapter := &sessionManagerAdapter{sm: sm}
	shellMode := terminal.NewRawShellMode(smAdapter, inputMan, term)

	core := &Core{
		Config:               cfg,
		SessionManager:       sm,
		ShellListenerManager: slm,
		terminal:             term,
		eventChan:            eventChan,
		InputMan:             inputMan,
		shellMode:            shellMode,
		commander:            NewCommandHandler(sm, slm, term, shellMode),
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
		case eventAny := <-c.eventChan:
			switch e := eventAny.(type) {
			case session.SessionConnectedEvent:
				handleNewSessionEvent(e, c)
			case terminal.ShellByteEvent:
				handleShellByteEvent(e, c)
			case terminal.MenuCmdEvent:
				handleMenuCmdEvent(e, c)
			}
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
