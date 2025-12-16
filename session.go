package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
)

type SessionState int

const (
	StateIdle SessionState = iota
	StateActive
	StateBackgrounded
	StateDead
)

type Session struct {
	ID int

	conn net.Conn
	Addr string

	mu           sync.Mutex
	state        SessionState
	outputBuffer bytes.Buffer
	stopChan     chan struct{}
	display      io.Writer
}

type SessionManager struct {
	mu         sync.RWMutex
	currentID int
	sessions   map[int]*Session
	display    Display
}

func NewSessionManager(display Display) *SessionManager {
	return &SessionManager{
		currentID: 0,
		sessions:   make(map[int]*Session),
		display:    display,
	}
}

func (sm *SessionManager) incID() int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	res := sm.currentID
	sm.currentID += 1
	return res
}

func (sm *SessionManager) AddSession(conn net.Conn) *Session {
	session := Session{
		ID:      sm.incID(),
		conn:    conn,
		Addr:    conn.RemoteAddr().String(),
		display: sm.display,
	}

	sm.mu.Lock()
	sm.sessions[session.ID] = &session
	sm.mu.Unlock()

	session.Start()
	return &session
}

func (sm *SessionManager) GetAmount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.sessions)
}

func (sm *SessionManager) GetSessions() []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*Session, 0, len(sm.sessions))
	for _, sess := range sm.sessions {
		sessions = append(sessions, sess)
	}
	return sessions
}

func (sm *SessionManager) Exists(ID int) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, ok := sm.sessions[ID]
	return ok
}

func (sm *SessionManager) Get(ID int) (*Session, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	sesh, ok := sm.sessions[ID]
	if !ok {
		// TODO: Add better message
		return nil, fmt.Errorf("Failed to get session")
	}

	return sesh, nil
}

func (s *Session) Start() {
	s.mu.Lock()
	s.state = StateBackgrounded
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	go s.outputLoop()
}

func (s *Session) outputLoop() {
	buf := make([]byte, 4096)

	for {
		select {
		// stop
		case <-s.stopChan:
			return
		default:
			n, err := s.conn.Read(buf)
			if err != nil {
				s.mu.Lock()
				s.state = StateDead
				s.mu.Unlock()
				return
			}

			if n > 0 {
				data := buf[:n]

				s.mu.Lock()
				// Defines behavior of Background, Foreground
				switch s.state {
				case StateActive:
					s.display.Write(data)
				case StateBackgrounded:
					s.outputBuffer.Write(data)
				case StateDead:
					s.mu.Unlock()
					return
				}
				s.mu.Unlock()
			}
		}
	}
}

func (s *Session) Foreground() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = StateActive

	if s.outputBuffer.Len() > 0 {
		s.display.Write(s.outputBuffer.Bytes())
		s.outputBuffer.Reset()
	}

}

func (s *Session) Background() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = StateBackgrounded
}

func (s *Session) Write(data []byte) error {
	_, err := s.conn.Write(data)
	return err
}

func (s *Session) Stop() {
	s.mu.Lock()
	s.state = StateDead
	s.mu.Unlock()

	close(s.stopChan)
	s.conn.Close()
}
