package main

import (
	"fmt"
)

type Event interface {
	Handle(c *Core)
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

func (e NewSessionEvent) Handle(c *Core) {
	c.Message("%s", fmt.Sprintf("\n[+] New session #%d from %s\n", e.Session.ID, e.Session.Addr))
	c.Prompt()
}

func (e ShellByteEvent) Handle(c *Core) {
	// WARN: Should be buffer instead?
	if e.Byte == CtrlD {
		c.shellMode.Exit()
		return
	}

	// Convert CR to LF for Unix shells
	// FIX: HMMM without this the overal shell doesnt wrap properly (even when exiting???)
	if e.Byte == '\r' {
		e.Byte = '\n'
	}

	// Forward to remote shell
	session, err := c.SessionManager.Get(c.shellMode.GetActiveSessionId())
	if err != nil {
		// Session died while we were in it
		c.shellMode.Exit()
		return
	}
	session.Write([]byte{e.Byte})

}

func (e MenuCmdEvent) Handle(c *Core) {
	c.commander.Execute(e.Input)
}
