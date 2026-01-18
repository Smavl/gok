// //go:build !integration
package main

import (
	"io"
	"net"
	"testing"
)


func init()  {
	cmdTimeout = TestingTimeout
}


func TestIdIncrement(t *testing.T) {
	sm  := NewSessionManager()

	before := sm.currentID
	sm.incID()
	after := sm.currentID

	if before + 1 != after {
		t.Errorf(`Increment ID, got %d, wanted %d`, after, (before +1 ))
	}
}

func TestGetNon_existentSession(t *testing.T) {
	sm  := NewSessionManager()

	// test 10 first ids
	for i := range 10 {
		session, err := sm.Get(i)
		// session, err := sm.Get(42)

		want := ErrSessionNotFound 

		// TEST: check that error is returned
		if err == nil {
			t.Errorf(`SessionManager: failed successfully to get Non-existent, got %v`, session)
		}

		// TEST: check that error is correct
		if err != want {
			t.Errorf(`SessionManager: Get() returned wrong error: got %v, wanted %v`, err, want)
		}
	}
}


func TestAddFakeSession(t *testing.T) {
	sm  := NewSessionManager()

	// Fake connection
	fakeConn, client := net.Pipe()
	go io.Copy(io.Discard, client) // AddSession expects output from session-workaround

	d := NewFakeDisplay()

	_, err := sm.AddSession(fakeConn,d)

	// TEST: check for error
	if err != nil {
		t.Errorf(`SessionManager: AddSession() returned error: %v`, err)
	}

}

func TestSessionPopulate(t *testing.T) {

	sm  := NewSessionManager()

	cmdTimeout = TestingTimeout

	for sessions := range 12 {

		// Fake connection
		fakeConn, client := net.Pipe()
		go io.Copy(io.Discard, client) // AddSession expects output from session-workaround

		d := NewFakeDisplay()

		_, err := sm.AddSession(fakeConn,d)

		// TEST: check for error
		if err != nil {
			t.Errorf(`SessionManager: AddSession() returned error: %v`, err)
		}

		gotAmountOfSessions := sm.GetAmountOfSessions()
		wantSessions := sessions +1
		// TEST: Check that amount of sessions is correct
		if gotAmountOfSessions != wantSessions {
			t.Errorf(`SessionManager: GetAmountOfSessions() returned wrong amount: got %d, wanted %d`, gotAmountOfSessions, wantSessions)
		}

	}

}


func TestGetExistingSession(t *testing.T) {
	sm := NewSessionManager()

	// Add a session
	fakeConn, client := net.Pipe()
	go io.Copy(io.Discard, client)
	d := NewFakeDisplay()

	session, _ := sm.AddSession(fakeConn, d)

	// Try to get it back
	got, err := sm.Get(session.ID)

	// TEST: check for error
	if err != nil {
		t.Errorf("Get() returned error for existing session: %v", err)
	}

	// TEST: check that the correct session is returned
	if got.ID != session.ID {
		t.Errorf("Get() returned wrong session: got ID %d, want %d", got.ID, session.ID)
	}
}

func TestIDUnique(t *testing.T) {
	sm := NewSessionManager()
	
	cmdTimeout = TestingTimeout

	// Add sessions
	const count = 20
	for range count {
		fakeConn, client := net.Pipe()
		go io.Copy(io.Discard, client)
		sm.AddSession(fakeConn, NewFakeDisplay())
	}

	sessions := sm.GetSessions()
	ids := make(map[int]int) 

	for _, s := range sessions {
		ids[s.ID]++
	}

	// Check for duplicates
	for id, count := range ids {
		if count > 1 {
			// TEST: duplicate ID found
			t.Errorf("ID %d appears %d times (should be unique)", id, count)
		}
	}

	// Check for gaps/expected range
	if len(ids) != count {
		// TEST: unexpected number of unique IDs
		t.Errorf("Expected %d unique IDs, got %d", count, len(ids))
	}
}

