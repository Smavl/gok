package terminal

import (
	"fmt"
	"os"
	"sync"

	"github.com/smavl/gok/internal/domain"
	"golang.org/x/term"
)


type Terminal struct {
	mu         sync.Mutex
	savedState *term.State
	display    Display
	rawMode    bool
}
type HeadlessTerminal struct {
	mu         sync.Mutex
	display    Display
	// rawMode    bool
	// savedState *term.State
}

func NewTerminal(display Display) (domain.TerminalController, error) {
	state, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	return &Terminal{
		savedState: state,
		display:    display,
		rawMode:    false,
	}, nil
}

func (t *Terminal) SetRaw() {
	t.mu.Lock()
	defer t.mu.Unlock()
	term.MakeRaw(int(os.Stdin.Fd()))
	t.rawMode = true
}

func (t *Terminal) Restore() {
	t.mu.Lock()
	defer t.mu.Unlock()
	term.Restore(int(os.Stdin.Fd()), t.savedState)
	t.rawMode = false
}

func (t *Terminal) Write(b []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// In raw mode, convert LF to CRLF for proper display
	if t.rawMode {
		converted := make([]byte, 0, len(b)*2)
		for _, ch := range b {
			if ch == '\n' {
				converted = append(converted, '\r', '\n')
			} else {
				converted = append(converted, ch)
			}
		}
		return t.display.Write(converted)
	}

	return t.display.Write(b)
}

func (t *Terminal) Message(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	t.Write([]byte(msg))
}

func (t *Terminal) Prompt() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Don't show prompt in raw mode (shell mode)
	if t.rawMode {
		return
	}

	t.display.Write([]byte("GOK > "))
}

func (t *Terminal) ShowError(err error) {
	t.Write([]byte(err.Error()))
}


func NewHeadlessTerminal(display Display) (domain.TerminalController, error) {
	return &HeadlessTerminal{
		display:    display,
		// savedState: nil,
		// rawMode:    false,
	}, nil
}

func (t *HeadlessTerminal) Message(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	t.display.Write([]byte(msg))
}

func (t *HeadlessTerminal) Prompt() {
	// No prompt in headless mode
}

func (t *HeadlessTerminal) SetRaw() {
	// No terminal state to set in headless mode
}

func (t *HeadlessTerminal) Restore() {
	// No terminal state to restore in headless mode
}

func (t *HeadlessTerminal) Write(b []byte) (int, error) {
	return t.display.Write(b)
}

func (t *HeadlessTerminal) ShowError(err error) {
	t.Write([]byte(err.Error()))
}
