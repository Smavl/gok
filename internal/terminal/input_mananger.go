package terminal

import (
	"context"
	"sync"

	"github.com/smavl/gok/internal/event"
)

type InputManagerImpl struct {
	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
	reader  InputReader
	wg      sync.WaitGroup
	// Event channels
	shellChan chan<- event.ShellByteEvent
	menuChan  chan<- event.MenuCmdEvent
}

type InputManager interface {
	Start(ctx context.Context)
	Stop()
	run()
	SwapReader(reader InputReader)
}

// func NewInputManager(initialInputRead InputReader, evenCh chan<- event.Event) *InputManagerImpl {
func NewInputManager(initialInputRead InputReader, shellCh chan<- event.ShellByteEvent, menuCh chan<- event.MenuCmdEvent) *InputManagerImpl {
	return &InputManagerImpl{
		reader:  initialInputRead,
		shellChan: shellCh,
		menuChan:  menuCh,
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
			// Read event from InputReader
			val := im.reader.Read()

			// WARN: Dirty type cast
			switch e := val.(type) {
			case event.ShellByteEvent:
				// send to shell channel
				select {
				case im.shellChan <- e:
				case <-im.ctx.Done(): return
				}
			case event.MenuCmdEvent:
				// send to menu channel
				select {
				case im.menuChan <- e:
				case <-im.ctx.Done(): return
				}
			default:
				// unknown event type, ignore
			}
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
