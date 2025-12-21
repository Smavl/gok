package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

type SessionState int

// TODO: Add state transitions validation?
const (
	StateActive SessionState = iota
	StateBackgrounded
	StateDead
	// StateIdle
)

// TODO: Add flag option
const defaultHistoryMaxLines = 50

type LineBuffer interface {
	Feed(bytes []byte)
	AddLine(line string)
}

type HistoryLineBuffer struct {
	lines      []string
	maxLines   int
	partialBuf string
}

func (lb *HistoryLineBuffer) Feed(bytes []byte) {
	lb.partialBuf += string(bytes)

	for {
		idx := strings.Index(lb.partialBuf, "\n")
		if idx == -1 {
			break
		}

		line := lb.partialBuf[:idx+1]
		lb.AddLine(line)
		// trim partialBuf from the newly added line
		lb.partialBuf = lb.partialBuf[idx+1:]
	}
}

func (lb *HistoryLineBuffer) AddLine(line string) {
	// ring the buffer
	if len(lb.lines) >= lb.maxLines {
		lb.lines = lb.lines[1:]

	}
	lb.lines = append(lb.lines, line)
}

func CreateLineBuffer(maxLines int) *HistoryLineBuffer {
	return &HistoryLineBuffer{
		lines:    make([]string, 0, maxLines),
		maxLines: maxLines,
	}
}

type Session struct {
	ID   int
	conn net.Conn
	Addr string

	mu           sync.Mutex
	state        SessionState
	outputBuffer bytes.Buffer
	display      io.Writer
	history      *HistoryLineBuffer

	// context things
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type SessionManager struct {
	mu        sync.RWMutex
	currentID int
	sessions  map[int]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		currentID: 0,
		sessions:  make(map[int]*Session),
	}
}

func (sm *SessionManager) incID() int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	res := sm.currentID
	sm.currentID += 1
	return res
}

func (sm *SessionManager) AddSession(conn net.Conn, display io.Writer) (*Session, error) {
	session := Session{
		ID:      sm.incID(),
		conn:    conn,
		Addr:    conn.RemoteAddr().String(),
		display: display,
		history: CreateLineBuffer(defaultHistoryMaxLines),
	}

	sm.mu.Lock()
	sm.sessions[session.ID] = &session
	sm.mu.Unlock()

	err := session.Start(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to add session: %v", err)
	}
	return &session, nil
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
		return nil, ErrSessionNotFound
	}

	return sesh, nil
}

func (s *Session) Start(ctx context.Context) error {
	s.mu.Lock()
	s.state = StateBackgrounded
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.mu.Unlock()

	s.wg.Add(1)
	go s.outputLoop()
	return nil
}

func (s *Session) outputLoop() {
	defer s.wg.Done()
	buf := make([]byte, 4096)

	// done channel to cancel Read()
	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-s.ctx.Done():
			s.conn.SetDeadline(time.Now())
		case <-done:
		}
	}()

	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			// Check if it's from our cancellation
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				select {
				case <-s.ctx.Done():
					return // Clean shutdown
				default:
					// clear the deadline
					s.conn.SetReadDeadline(time.Time{})
					continue
				}
			}

			// other must be normal error:
			s.mu.Lock()
			s.state = StateDead
			s.mu.Unlock()
			return
		}

		if n > 0 {
			data := buf[:n]
			s.history.Feed(data)

			s.mu.Lock()
			switch s.state {
			// forward to the display
			case StateActive:
				s.display.Write(data)
				// Buffer to background
			case StateBackgrounded:
				s.outputBuffer.Write(data)
				// ???
			case StateDead:
				s.mu.Unlock()
				return
			}
			s.mu.Unlock()
		}
	}
}

func (s *Session) Foreground() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = StateActive

	// drain buffer when forgrounding
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
	if s.cancel != nil {
		s.cancel()
	}
	s.mu.Unlock()

	s.wg.Wait()
	s.conn.Close()
}
