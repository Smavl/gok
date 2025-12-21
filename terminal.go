package main

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/term"
)

type Terminal struct {
	mu         sync.Mutex
	savedState *term.State
	display    Display
	rawMode    bool
}

func NewTerminal(display Display) (*Terminal, error) {
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
