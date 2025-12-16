package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Mode interface {
	OnNewSession(sess *Session)
	HandleInput(input string, c *Core)
}

type MenuMode struct {
	Display Display
}

type ShellMode struct{
	Display Display
}

// MenuMode
func (m *MenuMode) OnNewSession(sess *Session) {
	msg := fmt.Sprintf("\n[+] New session #%d from %s\n", sess.ID, sess.Addr)
	m.Display.Message(msg)
	m.Display.Prompt()
}

func (m *MenuMode) HandleInput(input string, c *Core) {
	// split on all whitespace
	args := strings.Fields(input)

	defer c.Prompt()

	if len(args) == 0 {
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

	case "sessions", "sesh", "sess", "s":
		if c.SessionManager.GetAmount() == 0 {

			c.Message("\n[!] No active sessions\n")
		} else {
			c.Message("\nActive Sessions:\n")
			for _, sess := range c.SessionManager.GetSessions() {
				c.Message("\t[%d] %s\n", sess.ID, sess.Addr)
			}
		}

	case "interact", "int", "i":
		// check if session arg is supplied
		if len(args) == 2 {
			id, err := strconv.Atoi(args[1])
			if err != nil {
				c.Message("Id: %v is not a number\n", id)
			}

			sessionExists := c.SessionManager.Exists(id)

			if sessionExists {
				c.activeShellID = id
				c.Message("[*] Session #%v: Dropping into shell..\n", id)
				c.EnableShellMode()
			} else {
				c.Message("Session #%d not found\n", id)
			}

		} else {
			c.Message("[!] No session chosen, or invalid argument\n")
		}

	case "exit", "quit":
		// TODO: Need to do more?
		os.Exit(0)

	default:
		// TODO: Add help suggestion
		c.Message("[!] Unknown command: %s\n", subCmd)
	}
}

// ShellMode
func (m *ShellMode) OnNewSession(sess *Session) {
	// NOTE: ATM we dont announce anything while in a shell!
	// TODO: Add flag to announce
	m.Display.Message("kaj\n")
}

func (m *ShellMode) HandleInput(input string, c *Core) {
	// ShellMode input is handled by the specialized runShellReader loop,
	// not this line-based handler.
	// TODO: Fix into here?
}
