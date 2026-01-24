package session

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/smavl/gok/internal/cli"
	"github.com/smavl/gok/internal/event"
)

type Listener struct {
	address  string
	port     int
	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// TerminalController is a minimal interface to avoid circular imports
type TerminalController interface {
	Message(format string, a ...any)
	Write([]byte) (int, error)
}

// type SessionConnectedEvent struct {
// 	Session *Session
// }

type ListenerManager interface {
	Init(config cli.Config)
	Start(addr string, port int) (*Listener, error)
	GetAddresses() []string
}

type ShellListenerManager struct {
	mu             sync.RWMutex
	listeners      map[string]*Listener
	terminal       TerminalController
	sessionManager *SessionManager
	eventChan      chan <- event.NewSessionEvent
}

func NewShellListenerManager(sm *SessionManager, terminal TerminalController, eventChan chan<- event.NewSessionEvent) *ShellListenerManager {
	return &ShellListenerManager{
		listeners:      make(map[string]*Listener),
		sessionManager: sm,
		terminal:       terminal,
		eventChan:      eventChan,
	}
}

func (lm *ShellListenerManager) Init(ctx context.Context, config cli.Config) {
	lm.terminal.Message("[+] Initializing listeners:\n\t")

	for _, addr := range config.BindIps {
		for _, port := range config.PortRange.Ports {
			lm.terminal.Message("%s:%d ", addr, port)

			l, err := lm.Start(ctx, addr, port)
			if err != nil {
				log.Printf("[-] Failed to start listener on %s:%d: %v", addr, port, err)
				continue
			}
			lm.mu.Lock()

			id := fmt.Sprintf("%s:%d", l.address, l.port)
			lm.listeners[id] = l
			lm.mu.Unlock()
		}
	}
	lm.terminal.Message("\n[*] Waiting for connections...\n")
}

func (lm *ShellListenerManager) GetAddresses() []string {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	addrs := make([]string, 0, len(lm.listeners))
	for addr := range lm.listeners {
		addrs = append(addrs, addr)
	}
	return addrs
}

func (lm *ShellListenerManager) Start(ctx context.Context, addr string, port int) (*Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil, fmt.Errorf("failed to start listener: %w", err)
	}

	listenerCtx, cancel := context.WithCancel(ctx)

	l := &Listener{
		address:  addr,
		port:     port,
		listener: listener,
		ctx:      listenerCtx,
		cancel:   cancel,
	}

	l.wg.Add(1)

	go l.acceptLoop(lm.sessionManager, lm.terminal, lm.eventChan)

	return l, nil
}
func (l *Listener) acceptLoop(sm *SessionManager, terminal TerminalController, eventChan chan<- event.NewSessionEvent) {
	defer l.wg.Done()
	defer l.listener.Close()

	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			// Set deadline; make Accept cancellable
			l.listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second))

			conn, err := l.listener.Accept()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // Check context again
				}
				log.Printf("[-] Error accepting connection: %v", err)
				return
			}

			session, err := sm.AddSession(conn, terminal)
			if err != nil {
				log.Printf("Failed to add session: %v", err)
				continue
			}

			eventChan <- event.NewSessionEvent{
				SessionID:   session.ID,
				SessionAddr: session.Addr,
				SystemOS:    session.SystemInfo.OS.String(),
			}
		}
	}
}

func (l *Listener) Stop() {
	if l.cancel != nil {
		l.cancel()
	}
	l.wg.Wait()
}
