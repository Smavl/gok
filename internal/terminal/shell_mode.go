package terminal

import (
	"sync"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/event"
)

type RawShellMode struct {
	mu            sync.Mutex
	activeShellID int

	currentSession domain.Session

	im       *InputManagerImpl
	terminal domain.TerminalController
	shellCh  chan<- event.ShellByteEvent
	menuCh   chan<- event.MenuCmdEvent
}

func NewRawShellMode(inputMan *InputManagerImpl, terminal domain.TerminalController, shellCh chan<- event.ShellByteEvent, menuCh chan<- event.MenuCmdEvent) *RawShellMode {
	return &RawShellMode{
		im:       inputMan,
		terminal: terminal,
		shellCh:  shellCh,
		menuCh:   menuCh,
	}
}

func (m *RawShellMode) Enter(s domain.Session) {
	ID := s.GetID()

	m.terminal.Message("[*] Session #%v: Dropping into shell..\n", ID)
	m.activeShellID = ID
	m.currentSession = s

	m.terminal.SetRaw()
	s.Foreground()
	s.Write([]byte{'\n'})
	m.im.SwapReader(NewByteReader(m.shellCh))
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
