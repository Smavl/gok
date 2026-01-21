package core

import (
	"github.com/smavl/gok/internal/session"
	"github.com/smavl/gok/internal/terminal"
)

type EventHandler interface {
	Sessions() *session.SessionManager
	ShellMode() *terminal.RawShellMode
	Terminal() terminal.TerminalController
	Commander() *CommandHandler
}

func (c *Core) Sessions() *session.SessionManager         { return c.SessionManager }
func (c *Core) ShellMode() *terminal.RawShellMode         { return c.shellMode }
func (c *Core) Terminal() terminal.TerminalController     { return c.terminal }
func (c *Core) Commander() *CommandHandler                { return c.commander }

type Event interface {
	Handle(h EventHandler)
}

//# Control Characters

const (
	CtrlD = 0x04 // Exit shell
	CtrlC = 0x03 // Interrupt
	CtrlL = 0x0C // Clear screen
)

// Triggered when a shell lands
func handleNewSessionEvent(e session.SessionConnectedEvent, h EventHandler) {
	h.Terminal().Message("\n[+] %s => New session #%d | %s \n", e.Session.Addr, e.Session.ID, e.Session.SystemInfo.OS.String())
	h.Terminal().Prompt()
}

func handleShellByteEvent(e terminal.ShellByteEvent, h EventHandler) {
	// WARN: Should be buffer instead?
	if e.Byte == CtrlD {
		h.ShellMode().Exit()
		return
	}

	// Convert CR to LF for Unix shells
	// FIX: HMMM without this the overal shell doesnt wrap properly (even when exiting???)
	b := e.Byte
	if b == '\r' {
		b = '\n'
	}

	// Forward to remote shell
	sess, err := h.Sessions().Get(h.ShellMode().GetActiveSessionId())
	if err != nil {
		// Session died while we were in it
		h.ShellMode().Exit()
		return
	}
	sess.Write([]byte{b})
}

func handleMenuCmdEvent(e terminal.MenuCmdEvent, h EventHandler) {
	h.Commander().Execute(e.Input)
}
