package terminal

import (
	"sync"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/event"
	"github.com/smavl/gok/internal/misc"
)

type RawShellMode struct {
	mu            sync.Mutex
	activeShellID int

	currentSession domain.InteractiveSession

	im       *InputManagerImpl
	terminal domain.TerminalController
	shellCh  chan<- event.ShellByteEvent
	menuCh   chan<- event.MenuCmdEvent
	IsEntered bool
}

func NewRawShellMode(inputMan *InputManagerImpl, terminal domain.TerminalController, shellCh chan<- event.ShellByteEvent, menuCh chan<- event.MenuCmdEvent) *RawShellMode {
	return &RawShellMode{
		im:       inputMan,
		terminal: terminal,
		shellCh:  shellCh,
		menuCh:   menuCh,
		IsEntered: false,
	}
}

func (m *RawShellMode) isActiveDirty() bool {
	return m.IsEntered
}

func (m *RawShellMode) setActiveDirty(active bool) {
	m.IsEntered = active
}


func (m *RawShellMode) IsActive() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.IsEntered
}

func (m *RawShellMode) SetActive(active bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.IsEntered = active
}

func (m *RawShellMode) Enter(s domain.InteractiveSession) error {
	m.mu.Lock()
	if m.isActiveDirty() {
		m.mu.Unlock()
		m.terminal.Message("[!] Already in shell mode. How did you do this??\n")
		return misc.ErrAlreadyInShellMode
	}
	ID := s.GetID()

	m.activeShellID = ID
	m.currentSession = s
	m.setActiveDirty(true)
	m.mu.Unlock()

	m.terminal.Message("[*] Session #%v: Dropping into shell..\n", ID)
	m.terminal.SetRaw()
	// TODO: Refactor this. Move Foreground() and maybe embed SetActive() ?
	s.Foreground()
	// TODO: Refactor this hacky call at some point
	s.Write([]byte{'\n'})
	m.im.SwapReader(NewByteReader(m.shellCh))

	return nil
}

func (m *RawShellMode) GetActiveSessionId() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeShellID
}

func (m *RawShellMode) Exit() {
	m.mu.Lock()
	// Get current session
	session := m.currentSession

	// background
	m.setActiveDirty(false)
	m.mu.Unlock()

	// restore the user terminal
	session.Background()
	m.terminal.Restore()
	m.terminal.Message("\n[*] Returning to main menu\n")

	// Swap back to line reader
	m.im.SwapReader(NewLineReader(m.menuCh))
	defer m.terminal.Prompt()
}
