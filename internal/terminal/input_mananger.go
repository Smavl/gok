package terminal

import (
	"context"
	"io"
	"os"
	"sync"
)

type InputManagerImpl struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	reader InputReader
	input  io.Reader
	wg     sync.WaitGroup
}

type InputManager interface {
	Start(ctx context.Context)
	Stop()
	run()
	SwapReader(reader InputReader)
}

func NewInputManager(initialInputRead InputReader) *InputManagerImpl {
	return NewInputManagerWithInput(initialInputRead, os.Stdin)
}

// To make testing easier, allow to supply
func NewInputManagerWithInput(initialInputRead InputReader, input io.Reader) *InputManagerImpl {
	return &InputManagerImpl{
		reader: initialInputRead,
		input:  input,
	}
}

func (im *InputManagerImpl) Start(ctx context.Context) {
	im.mu.Lock()
	defer im.mu.Unlock()

	im.ctx, im.cancel = context.WithCancel(ctx)
	im.wg.Add(1)

	go im.run()
}

func (im *InputManagerImpl) RequestStop() {
	im.mu.Lock()

	if im.cancel != nil {
		im.cancel()
	}

	im.mu.Unlock()
}

// TODO: Error handling?
// Run loop
func (im *InputManagerImpl) run() {
	// signal to WaitGroup when done
	defer im.wg.Done()
	var buf [1]byte

	for {
		select {
		case <-im.ctx.Done():
			// exits function (and wg.Done is called)
			return
		default:
		}

		n, err := im.input.Read(buf[:])
		if err != nil || n == 0 {
			continue
		}

		im.mu.Lock()
		reader := im.reader
		im.mu.Unlock()

		// Dispatch the byte using the current mode handler.
		_ = reader.HandleByte(buf[0])
	}
}

func (im *InputManagerImpl) SwapReader(reader InputReader) {
	im.mu.Lock()
	im.reader = reader
	im.mu.Unlock()
}
