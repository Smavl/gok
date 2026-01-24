package terminal

import (
	"sync"

	"github.com/smavl/gok/internal/event"
	"github.com/smavl/gok/internal/session"
)


type SessionInterface interface {
	Foreground()
	Background()
	Write([]byte) (int, error)
}

type RawShellMode struct {
	mu            sync.Mutex
	activeShellID int

	currentSession *session.Session

	im       *InputManagerImpl
	terminal TerminalController
	shellCh  chan<- event.ShellByteEvent
	menuCh   chan<- event.MenuCmdEvent
}

func NewRawShellMode(inputMan *InputManagerImpl, terminal TerminalController, shellCh chan<- event.ShellByteEvent, menuCh chan<- event.MenuCmdEvent) *RawShellMode {
	return &RawShellMode{
		im:       inputMan,
		terminal: terminal,
		shellCh:  shellCh,
		menuCh:   menuCh,
	}
}

func (m *RawShellMode) Enter(session *session.Session) {
	ID := session.GetID()

	m.terminal.Message("[*] Session #%v: Dropping into shell..\n", ID)
	m.activeShellID = ID
	m.currentSession = session

	m.terminal.SetRaw()
	session.Foreground()
	m.im.SwapReader(NewByteReader(m.shellCh))
	session.Write([]byte{'\n'})
}

func (m *RawShellMode) GetActiveSessionId() int {
	return m.activeShellID
}

func (m *RawShellMode) Exit() {
	// Get current session
	session := m.currentSession

	// background
	session.Background()

	// restore the user terminal
	m.terminal.Restore()
	m.terminal.Message("\n[*] Returning to main menu\n")

	// Swap back to line reader
	m.im.SwapReader(NewLineReader(m.menuCh))
	defer m.terminal.Prompt()
}
