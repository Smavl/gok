package main

import "fmt"

type Mode interface {
	OnNewSession(sess *Session)
}


type MenuMode struct {
	Display Display
}
type ShellMode struct {

}

// MenuMode

func (m *MenuMode) OnNewSession(sess *Session) {
	msg := fmt.Sprintf("\n[+] New session #%d from %s\n", sess.ID, sess.Addr)
	m.Display.Message(msg)
	m.Display.Prompt()
}

// ShellMode

func (m *ShellMode) OnNewSession(sess *Session) {
	// NOTE: ATM we dont announce anything while in a shell!
}

