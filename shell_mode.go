package main

import (
	"errors"
	"sync"
)

// type ShellMode interface {
// 	Enter(ID int)
// 	Exit()
// }

type RawShellMode struct {
	mu sync.Mutex
	activeShellID int

	sm *SessionManager
	im *InputManagerImpl
	terminal *Terminal
}

func NewRawShellMode(sessionMan *SessionManager, inputMan *InputManagerImpl, terminal *Terminal) *RawShellMode {
	return &RawShellMode{
		sm: sessionMan,
		im: inputMan,
		terminal: terminal,
	}
}

func (m *RawShellMode) Enter(ID int) {
	session, err := m.sm.Get(ID)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			m.terminal.ShowError(err)
		}
		m.terminal.Message("[!] Error: Session #%d not found\n", ID)
		m.terminal.Prompt()
		return
	}

	m.terminal.Message("[*] Session #%v: Dropping into shell..\n", ID)
	m.activeShellID = ID

	m.terminal.SetRaw()
	session.Foreground()
	m.im.SwapReader(NewByteReader())
}
func (m *RawShellMode) GetActiveSessionId() int {
	return m.activeShellID
}

func (m *RawShellMode) Exit() {
	session, err := m.sm.Get(m.GetActiveSessionId())
	if err != nil {
		// Unexpcted but do restore terminal anyway
		m.terminal.Message("\r\n[!] Warning: Session no longer exists\r\n")
	} else {
		session.Background()
	}

	m.terminal.Restore()

	m.terminal.Message("\r\n")
	m.terminal.Message("[*] Returning to main menu\n")

	// Swap back to line reader
	m.im.SwapReader(NewLineReader())
	defer m.terminal.Prompt()
}
