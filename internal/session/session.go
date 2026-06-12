package session

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/misc"
	"github.com/smavl/gok/internal/prober"
	"github.com/smavl/gok/internal/prober/types"
)

type SessionState int

// TODO: Add state transitions validation?
const (
	StateActive SessionState = iota
	StateBackgrounded
	StateDead
	StateProbing
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

type SystemInfo struct {
	OS types.OS
}

type Session struct {
	ID   int
	conn net.Conn
	Addr string

	mu      sync.Mutex
	state   SessionState
	display io.Writer
	history *HistoryLineBuffer

	// probing
	probingBuffer      *HistoryLineBuffer
	probingDataArrived chan struct{}
	SystemInfo         SystemInfo
	Prober             *prober.Prober
	ProbingOptions domain.ProbingOptions

	// context things
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (s *Session) GetID() int {
	return s.ID
}

type SessionManager struct {
	mu          sync.RWMutex
	currentID   int
	sessions    map[int]*Session
	probOpts domain.ProbingOptions 
}

func NewSessionManager(probingOpts domain.ProbingOptions) *SessionManager {
	return &SessionManager{
		currentID:   0,
		sessions:    make(map[int]*Session),
		probOpts: probingOpts,
		// ProbingOpTimout: probingOpTimeout,
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
		ID:                 sm.incID(),
		conn:               conn,
		Addr:               conn.RemoteAddr().String(),
		display:            display,
		history:            CreateLineBuffer(defaultHistoryMaxLines),
		probingBuffer:      CreateLineBuffer(defaultHistoryMaxLines),
		probingDataArrived: make(chan struct{}),
		SystemInfo:         SystemInfo{},
		// probingMode:        sm.probOpts.ProbingMode,
		ProbingOptions:   sm.probOpts,
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

func (sm *SessionManager) GetAmountOfSessions() int {
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
		return nil, misc.ErrSessionNotFound
	}

	return sesh, nil
}


// At this point the shell has landed
func (s *Session) Start(ctx context.Context) error {
	s.mu.Lock()
	s.state = StateBackgrounded
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.mu.Unlock()

	s.wg.Add(1)
	go s.outputLoop()

	// NOTE: probing session has to happen after outputLoop is initialized
	if !s.ProbingOptions.DisableProber {
		s.probeSession()
	}
	return nil
}

func (s *Session) probeSession() error {
	s.mu.Lock()
	s.state = StateProbing
	s.mu.Unlock()

	// TODO: Add OS detection as an operation in PhaseInitial
	// For now, assume Linux
	s.SystemInfo.OS = types.LinuxOs

	// Create prober with configured mode
	newProber, err := prober.NewProber(s, s.ProbingOptions)
	if err != nil {
		return fmt.Errorf("failed to create prober: %w", err)
	}

	s.mu.Lock()
	s.Prober = newProber
	s.mu.Unlock()

	// TODO: Is this needed? Does this kill valid sessions?
	// Use a timeout to prevent hanging on unresponsive shells
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = newProber.Run(ctx)
	if err != nil {
		return fmt.Errorf("probing failed: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = StateBackgrounded

	return nil
}

func (s *Session) GetProbingLines() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.probingBuffer.lines
}

func (s *Session) GetProbingDataChannel() <-chan struct{} {
	return s.probingDataArrived
}

func (s *Session) IsProberDone() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state == StateBackgrounded && s.Prober != nil
}

func (s *Session) ClearProbingBuffer() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// slices are well behaved so nil should be fine
	s.probingBuffer.lines = nil
	s.probingBuffer.partialBuf = ""
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

			s.mu.Lock()
			switch s.state {
			// forward to the display
			case StateActive:
				s.history.Feed(data)
				s.display.Write(data)
			// Buffer to background
			// case StateBackgrounded:
			//
			case StateProbing:
				s.probingBuffer.Feed(data)
				s.mu.Unlock()
				// signal that probing data is incomming (after releasing lock to avoid blocking)
				select {
				case s.probingDataArrived <- struct{}{}:
				default:
				}
				continue
			case StateDead:
				s.mu.Unlock()
				return
			}
			s.mu.Unlock()
		}
	}
}

// send data to the remote session
func (s *Session) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state == StateDead {
		return 0, misc.ErrSessionDed
	}
	return s.conn.Write(p)
}

// foregound the session to the user
func (s *Session) Foreground() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = StateActive

	if len(s.history.lines) > 0 {
		s.display.Write([]byte(string("[*] Resuming session...\n")))
	}

	for _, l := range s.history.lines {
		s.display.Write([]byte(l))
	}
}

func (s *Session) Background() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state != StateDead {
		s.state = StateBackgrounded
	}
}

func (s *Session) Stop() {
	s.mu.Lock()
	if s.state == StateDead {
		s.mu.Unlock()
		return
	}
	s.state = StateDead
	if s.cancel != nil {
		s.cancel()
	}
	s.mu.Unlock()

	s.wg.Wait()
	s.conn.Close()
}
