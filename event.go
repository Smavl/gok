package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Event interface {
	Handle(c *Core)
}

//# Events

// Triggered when a shell lands
type NewSessionEvent struct {
	Session *Session
}

type ShellByteEvent struct {
	Byte byte
}

type MenuCmdEvent struct {
	input string 
}

func (e NewSessionEvent) Handle(c *Core) {
	c.Message("%s", fmt.Sprintf("\n[+] New session #%d from %s\n", e.Session.ID, e.Session.Addr))
	c.Prompt()
}

func (e ShellByteEvent) Handle(c *Core) {
	// WARN: Should be buffer instead?

	// ctrl+d
	if e.Byte == 0x04 {
		c.ExitShell()
		return
	}

	// Convert CR to LF for Unix shells
	// FIX: HMMM without this the overal shell doesnt wrap properly (even when exiting???)
	if e.Byte == '\r' {
		e.Byte = '\n'
	}

	// Forward to remote shell
	session, _ := c.SessionManager.Get(c.activeShellID)
	session.Write([]byte{e.Byte})

}

func (e MenuCmdEvent) Handle(c *Core) {
	// c.Prompt()
	// split on all whitespace
	args := strings.Fields(e.input)

	if len(args) == 0 {
		c.Prompt()
		return
	}

	subCmd := args[0]
	switch subCmd {
	// Management
	case "listeners", "lis", "l":
		c.mu.RLock()
		addresses := c.ShellListenerManager.GetAddresses()
		if len(addresses) == 0 {
			c.Message("[!] No active listeners\n")
		} else {
			c.Message("\nListeners:\n")
			for _,addr := range addresses {
				c.Message("- %v\n", addr)
			}
		}
		c.mu.RUnlock()
		defer c.Prompt()

	case "sessions", "sesh", "sess", "s":
		if c.SessionManager.GetAmount() == 0 {

			c.Message("\n[!] No active sessions\n")
		} else {
			c.Message("\nActive Sessions:\n")
			for _, sess := range c.SessionManager.GetSessions() {
				c.Message("\t[%d] %s\n", sess.ID, sess.Addr)
			}
		}
		defer c.Prompt()

	case "interact", "int", "i":
		// check if session arg is supplied
		if len(args) == 2 {
			id, err := strconv.Atoi(args[1])
			if err != nil {
				c.Message("Id: %v is not a number\n", id)
				c.Prompt()
				return
			}

			sessionExists := c.SessionManager.Exists(id)

			if sessionExists {
				c.activeShellID = id
				c.Message("[*] Session #%v: Dropping into shell..\n", id)
				c.EnterShell(id)
				// No prompt - we're in shell mode now!
			} else {
				c.Message("Session #%d not found\n", id)
				c.Prompt()
			}

		} else {
			c.Message("[!] No session chosen, or invalid argument\n")
			c.Prompt()
		}

	case "exit", "quit":
		// TODO: Need to do more?
		os.Exit(0)

	default:
		// TODO: Add help suggestion
		c.Message("[!] Unknown command: %s\n", subCmd)
		defer c.Prompt()
	}
}
