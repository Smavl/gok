package upgrader

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/prober/executor"
	"github.com/smavl/gok/internal/prober/types"
	"golang.org/x/term"
)


type Upgrader struct {
	ctx context.Context
	Session domain.CommandSession
	ProbeResults *types.ProbeResults
	Executor executor.Executor
}

func NewUpgrader(ctx context.Context, session domain.CommandSession, results *types.ProbeResults, executor executor.Executor ) *Upgrader {

	return &Upgrader{
		ctx: ctx,
		Session: session,
		ProbeResults: results,
		Executor: executor,
	}
}

// func (u *Upgrader) Upgrade(session domain.Session, results *types.ProbeResults) error {
func (u *Upgrader) Upgrade() error {
	// upgrade shell (PTY)
	err := u.UpgradePTY()
	if err != nil {
		return fmt.Errorf("failed to upgrade PTY: %w", err)
	}

	// Export envs
	err = u.exportENVs()
	if err != nil {
		return fmt.Errorf("failed to export envs: %w", err)
	}
	// set tty dims 
	rows, cols, err := GetTTYSize()
	if err != nil {
		return fmt.Errorf("failed to get TTY size: %w", err)
	}

	err = u.SetTTYSize(rows, cols)
	if err != nil {
		return fmt.Errorf("failed to set TTY size: %w", err)
	}

	return nil
}

// TODO: Implement returning fallbacks too
func determineBestPtyUpgrader(results *types.ProbeResults) (PtySpawner, error){
	// WIP
	switch {
	case results.HasBinary("python3"):
		return newPython3(results),nil
	case results.HasBinary("socat"):
		return &Socat{},nil
	case results.HasBinary("script"):
		return &Script{}, nil
	// TODO: Move up when implemented
	case results.HasBinary("python2"):
		return &Python2{}, nil
	default:
		return nil, os.ErrNotExist
	}
}

func (u *Upgrader) Execute(cmd string) error {
	// timeout context
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(u.ctx, timeout)
	defer cancel()

	_, err := u.Executor.Execute(ctx, u.Session, cmd)
	if err != nil {
		return fmt.Errorf("failed to execute '%s': %w", cmd, err)
	}
	return nil
}

// run pty spawner
func (u *Upgrader) UpgradePTY() error {
	ptySpawner, err := determineBestPtyUpgrader(u.ProbeResults)
	if err != nil {
		return err
	}
	ptyUpgradePayload := ptySpawner.GetPayload()

	// NOTE:
	// We cant use the Executor for this (e.g. the delimiter based one)
	// as we cant detect when its done, since the payload will never return
	bytes := []byte(ptyUpgradePayload)
	bytes = append(bytes, '\n')
	_, err = u.Session.Write(bytes)
	if err != nil {
		return err
	}

	// Wait for PTY to spawn
	time.Sleep(500 * time.Millisecond)

	// Clear any output from PTY spawn (shell prompt, etc.)
	u.Session.ClearProbingBuffer()

	return nil
}


// export env's 

func (u *Upgrader) exportENVs() error {
	binary, err := u.ProbeResults.GetBinary("bash")
	if err != nil { 
		return err
	}

	bashPath := binary.Path

	cmds := []string{
		"export SHELL=" + bashPath,
		"export TERM=xterm-256color",
	}

	for _, cmd := range cmds {
		err := u.Execute(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

// get cols and rows

func GetTTYSize() (int, int, error) {
	cols,rows,err := term.GetSize(int(os.Stdin.Fd()))

	return rows, cols, err
}

// set cols and rows
func (u *Upgrader) SetTTYSize(rows, cols int) error {

	rs := strconv.Itoa(rows)
	cs := strconv.Itoa(cols)
	cmd := "stty rows " + rs + " columns " + cs
	err := u.Execute(cmd)
	if err != nil {
		return err
	}

	return nil
}

// TODO: listener on window change
// requires developing a custom protocol for the best outcome
