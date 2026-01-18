//go:build integration

package main

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestContainer struct {
  Container testcontainers.Container
  ctx      context.Context
}

// Helper function to start a test container
func StartContainer(t *testing.T, image string) *TestContainer {
  t.Helper()
  ctx := context.Background()

  container, err := testcontainers.Run(
    ctx,
    image,
    testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
	    // hacky command that keeps the container running but can exit fast!
            Cmd: []string{"sh", "-c", "trap 'exit 0' TERM INT; echo ready; sleep infinity & wait"},
        },
    }),
    testcontainers.WithWaitStrategy(
      wait.ForLog("ready"),
      ),
    )
  if err != nil {
    t.Fatalf("Failed to start container: %v", err)
  }

  testcontainers.CleanupContainer(t, container)  // Auto cleanup

  return &TestContainer{
    Container: container,
    ctx:       ctx,
  }
}


func (tc *TestContainer) Exec(t *testing.T, cmd []string) (string, error) {
  t.Helper()

  exitCode, reader, err := tc.Container.Exec(tc.ctx, cmd)
  require.NoError(t,err)

  buf := new(strings.Builder)
  _, err = io.Copy(buf, reader)
  require.NoError(t, err)

  output := buf.String()

  require.Equal(t, 0, exitCode, "Non-zero exit code: %d, output: %s", exitCode, output)

  return output, nil
}
