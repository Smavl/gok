package main

import (
	"context"
	"sync"
)

type InputManagerImpl struct {
	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
	reader  InputReader
	eventCh chan<- Event
	wg      sync.WaitGroup
}

type InputManager interface {
	Start(ctx context.Context)
	Stop()
	run()
	SwapReader(reader InputReader)
}

func NewInputManager(initialInputRead InputReader, evenCh chan<- Event) *InputManagerImpl {
	return &InputManagerImpl{
		reader:  initialInputRead,
		eventCh: evenCh,
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

func (im *InputManagerImpl) run() {
	// signal to WaitGroup when done
	defer im.wg.Done()
	for {
		select {
		case <-im.ctx.Done():
			// exits function (and wg.Done is called)
			return
		default:
			// Read event from InputReader (producer?)
			event := im.reader.Read()
			// in case of event
			if event != nil {
				select {
				// send event to
				case im.eventCh <- event:
				// Check also for Done signal inside event?
				case <-im.ctx.Done():
					return
				}
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
