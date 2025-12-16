package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Listener struct {
	address  string
	port     int
	listener net.Listener
}

type ListenerManager interface {
	Init(config Config)
	Start(addr string, port int) (*Listener, error)
	GetAddresses() []string
}

type ShellListenerManager struct {
	mu             sync.RWMutex
	listeners      map[string]*Listener
	display        Display
	sessionManager *SessionManager
	eventChan      chan<- Event // Write-only channel for decoupling
}

// NewShellListenerManager creates a new listener manager.
// It requires its dependencies to be injected upon creation.
func NewShellListenerManager(sm *SessionManager, display Display, eventChan chan<- Event) *ShellListenerManager {
	return &ShellListenerManager{
		listeners:      make(map[string]*Listener),
		sessionManager: sm,
		display:        display,
		eventChan:      eventChan,
	}
}

func (lm *ShellListenerManager) Init(config Config) {
	lm.display.Message("[+] Initializing listeners:\n\t")

	for _, addr := range config.bindIps {
		for _, port := range config.PortRange.Ports {
			lm.display.Message(fmt.Sprintf("%s:%d ", addr, port))

			l, err := lm.Start(addr, port)
			if err != nil {
				log.Printf("[-] Failed to start listener on %s:%d: %v", addr, port, err)
				continue
			}
			lm.mu.Lock()
			// The key for the map is the address string.
			id := fmt.Sprintf("%s:%d", l.address, l.port)
			lm.listeners[id] = l
			lm.mu.Unlock()
		}
	}
	lm.display.Message("\n[*] Waiting for connections...\n")
}

// GetAddresses returns a slice of strings representing the addresses of active listeners.
func (lm *ShellListenerManager) GetAddresses() []string {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	addrs := make([]string, 0, len(lm.listeners))
	for addr := range lm.listeners {
		addrs = append(addrs, addr)
	}
	return addrs
}


func (lm *ShellListenerManager) Start(addr string, port int) (*Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil, fmt.Errorf("failed to start listener: %w", err)
	}

	l := &Listener{
		address:  addr,
		port:     port,
		listener: listener,
	}

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("[-] Error accepting connection: %v", err)
				return // Exit goroutine when listener is closed
			}

			// create session using the injected session manager
			session := lm.sessionManager.AddSession(conn)

			// announce to channel using the injected channel
			lm.eventChan <- NewSessionEvent{Session: session}
		}
	}()

	return l, nil
}
