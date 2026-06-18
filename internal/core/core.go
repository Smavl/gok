package core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/smavl/gok/internal/cli"
	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/event"
	"github.com/smavl/gok/internal/session"
	"github.com/smavl/gok/internal/terminal"
)

// TODO: Move elsewhere
const (
	CtrlD = 0x04 // Exit shell
	CtrlC = 0x03 // Interrupt
	CtrlL = 0x0C // Clear screen
)


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
	terminal  domain.TerminalController

	// event
	eventBus *event.EventBus

	// Context things
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewCore(cfg cli.Config) *Core {
	terminalDisplay := terminal.NewTerminalDisplay()

	var term domain.TerminalController
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

	probingOptions := domain.ProbingOptions{
		ProbingMode:       cfg.ProbingMode,
		DisableProber:     cfg.DisableProber,
	}

	eventBus := event.NewEventBus()
	inputMan := terminal.NewInputManager(terminal.NewLineReader(eventBus.Menu))
	sm := session.NewSessionManager(probingOptions)
	slm := session.NewShellListenerManager(sm, term, eventBus.Session)
	shellMode := terminal.NewRawShellMode(inputMan, term, eventBus.Shell, eventBus.Menu)

	core := &Core{
		Config:               cfg,
		SessionManager:       sm,
		ShellListenerManager: slm,
		terminal:             term,
		eventBus:             eventBus,
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
		case e := <-c.eventBus.Session:
			c.handleNewSessionEvent(e)
		case e := <-c.eventBus.Menu:
			c.handleMenuCmdEvent(e)
		case e := <-c.eventBus.Shell:
			c.handleShellByteEvent(e)
		case <-c.ctx.Done(): return
		}
	}
}


// Triggered when a shell lands
func (c *Core ) handleNewSessionEvent(e event.NewSessionEvent) {
	ID := e.SessionID
	addr := e.SessionAddr
	systemOS := e.SystemOS
	c.terminal.Message("\n[+] %s => New session #%d | %s \n", addr, ID, systemOS)
	c.terminal.Prompt()
}

func (c *Core ) handleMenuCmdEvent(e event.MenuCmdEvent) {
	input := e.Input
	c.commander.Execute(input)
}


func (c *Core ) handleShellByteEvent(e event.ShellByteEvent) {
	// WARN: Should be buffer instead?
	if e.Byte == CtrlD {
		// h.ShellMode().Exit()
		c.shellMode.Exit()
		return
	}

	// Convert CR to LF for Unix shells
	// FIX: HMMM without this the overal shell doesnt wrap properly (even when exiting???)
	b := e.Byte
	if b == '\r' {
		b = '\n'
	}

	// Forward to remote shell
	sess, err := c.SessionManager.Get(c.shellMode.GetActiveSessionId())
	if err != nil {
		// Session died while we were in it
		c.shellMode.Exit()
		return
	}
	sess.Write([]byte{b})
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
