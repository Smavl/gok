package main

import (
	"fmt"
	"net"
	"sync"
)

type Session struct {
	ID   int
	Conn net.Conn
	Addr string
	started bool
	mu sync.Mutex
}

type SessionManager struct {
	mu sync.RWMutex
	id_counter int
	sessions  map[int]*Session
	activeShellID int
}


func NewSessionManager() *SessionManager {
	return &SessionManager{
		id_counter: 0,
		sessions: make(map[int]*Session),
	}
}


func (sm *SessionManager) incID() int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	res := sm.id_counter
	sm.id_counter += 1
	return res
}

func (sm *SessionManager) AddSession(conn net.Conn) *Session {
	session := Session{
		ID:   sm.incID(),
		Conn: conn,
		Addr: conn.RemoteAddr().String(),
	}

	sm.mu.Lock()
	sm.sessions[session.ID] = &session
	sm.mu.Unlock()

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

	sessions := make([]*Session, 0,len(sm.sessions))
	for _, sess := range sm.sessions {
		sessions = append(sessions,sess)
	}
	return sessions
}

func (sm *SessionManager) Exists(ID int) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, ok := sm.sessions[ID]
	return ok
}

func (sm *SessionManager) Get(ID int) (*Session,error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	sesh, ok := sm.sessions[ID]
	if !ok {
		// TODO: Add better message
		return nil, fmt.Errorf("Failed to get session")
	}

	return sesh, nil
}

func (sm *SessionManager) SetActive(ID int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.activeShellID = ID
}

// func (sm *SessionManager) GetActive(ID int) { 
// 	sm.Get()
// }
