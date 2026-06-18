package terminal

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/smavl/gok/internal/event"
)

func writeInput(t *testing.T, w io.Writer, input string) {
	t.Helper()

	if _, err := w.Write([]byte(input)); err != nil {
		t.Fatalf("write input: %v", err)
	}
}

func readMenu(t *testing.T, menu <-chan event.MenuCmdEvent) string {
	t.Helper()

	select {
	case e := <-menu:
		return e.Input
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for menu event")
		return ""
	}
}

func readShellByte(t *testing.T, shell <-chan event.ShellByteEvent) byte {
	t.Helper()

	select {
	case e := <-shell:
		return e.Byte
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for shell byte")
		return 0
	}
}

// This is a regressive test for testing the swapping of readers
// It test the handling of:
// i 0
// <swap reader to raw shell mode>
// id
// <swap reader back to line reader>
// exit
//
// all these inputs should be recieved!
func TestInputManagerModeSwitchDoesNotDropFirstMenuByte(t *testing.T) {
	// test pipes for I/O
	ioR, ioW := io.Pipe()
	bus := event.NewEventBus()
	ctx, cancel := context.WithCancel(context.Background())

	// create input manager with with test pipes
	im := NewInputManagerWithInput(NewLineReader(bus.Menu), ioR)
	im.Start(ctx)

	// test cleanup
	t.Cleanup(func() {
		cancel()
		_ = ioR.Close()
		_ = ioW.Close()
		im.wg.Wait()
	})

	// input "i 0" to line reader
	writeInput(t, ioW, "i 0\n")
	if got := readMenu(t, bus.Menu); got != "i 0" {
		t.Fatalf("first menu command = %q, want %q", got, "i 0")
	}

	// swap to shell reader
	im.SwapReader(NewByteReader(bus.Shell))
	// input "id" to shell reader
	writeInput(t, ioW, "id")

	for _, want := range []byte("id") {
		if got := readShellByte(t, bus.Shell); got != want {
			t.Fatalf("shell byte = %q, want %q", got, want)
		}
	}

	// swap back to line reader
	im.SwapReader(NewLineReader(bus.Menu))
	// input "exit" to line reader
	writeInput(t, ioW, "exit\n")
	// With the bug present it would only get "xit"
	if got := readMenu(t, bus.Menu); got != "exit" {
		t.Fatalf("second menu command = %q, want %q", got, "exit")
	}
}
