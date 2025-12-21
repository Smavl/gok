package main

import (
	"fmt"
)

type EventHandler interface {
	Sessions()	*SessionManager
	ShellMode() *RawShellMode
	Terminal()	*Terminal
	Commander() *CommandHandler
}

func (c *Core) Sessions() *SessionManager { return c.SessionManager }
func (c *Core) ShellMode() *RawShellMode { return c.shellMode }
func (c *Core) Terminal() *Terminal { return c.terminal }
func (c *Core) Commander() *CommandHandler { return c.commander}

type Event interface {
	Handle(h EventHandler)
}

//# Control Characters

const (
	CtrlD = 0x04  // Exit shell
	CtrlC = 0x03  // Interrupt
	CtrlL = 0x0C  // Clear screen
)

// Triggered when a shell lands
type NewSessionEvent struct {
	Session *Session
}

type ShellByteEvent struct {
	Byte byte
}

type MenuCmdEvent struct {
	Input string 
}

func (e NewSessionEvent) Handle(h EventHandler) {
	h.Terminal().Message("%s", fmt.Sprintf("\n[+] New session #%d from %s\n", e.Session.ID, e.Session.Addr))
	h.Terminal().Prompt()
}

func (e ShellByteEvent) Handle(h EventHandler) {
	// WARN: Should be buffer instead?
	if e.Byte == CtrlD {
		h.ShellMode().Exit()
		return
	}

	// Convert CR to LF for Unix shells
	// FIX: HMMM without this the overal shell doesnt wrap properly (even when exiting???)
	if e.Byte == '\r' {
		e.Byte = '\n'
	}

	// Forward to remote shell
	session, err := h.Sessions().Get(h.ShellMode().GetActiveSessionId())
	if err != nil {
		// Session died while we were in it
		h.ShellMode().Exit()
		return
	}
	session.Write([]byte{e.Byte})

}

func (e MenuCmdEvent) Handle(h EventHandler) {
	h.Commander().Execute(e.Input)
}
