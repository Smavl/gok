package terminal

import (
	"context"
	"sync"
)

type InputManagerImpl struct {
	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
	reader  InputReader
	wg      sync.WaitGroup
}

type InputManager interface {
	Start(ctx context.Context)
	Stop()
	run()
	SwapReader(reader InputReader)
}

func NewInputManager(initialInputRead InputReader) *InputManagerImpl {
	return &InputManagerImpl{
		reader: initialInputRead,
	}
}

func (im *InputManagerImpl) Start(ctx context.Context) {
	im.mu.Lock()
	defer im.mu.Unlock()

	im.ctx, im.cancel = context.WithCancel(ctx)
	im.wg.Add(1)

	go im.run()
}

func (im *InputManagerImpl) Stop() {
	im.mu.Lock()

	if im.cancel != nil {
		im.cancel()

	}
	im.mu.Unlock()

	// wait gracefully
	im.wg.Wait()
}

// TODO: Error handling?
// Run loop
func (im *InputManagerImpl) run() {
	// signal to WaitGroup when done
	defer im.wg.Done()
	for {
		select {
		case <-im.ctx.Done():
			// exits function (and wg.Done is called)
			return
		default:
			// Read event from InputReader (reader dispatches internally)
			_ = im.reader.Read()
		}
	}
}

func (im *InputManagerImpl) SwapReader(reader InputReader) {
	im.Stop()

	im.mu.Lock()
	im.reader = reader
	im.mu.Unlock()

	im.Start(context.Background())
}
