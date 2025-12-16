package main

// since Go doesnt have cool enums like rust
type Event interface {
	isEvent()
}

//# Events

// Triggered when a shell lands
type NewSessionEvent struct {
	Session *Session
}

// for menuMode
type UserInputEvent struct {
	io string 
}

func (e NewSessionEvent) isEvent() {}

func (e UserInputEvent) isEvent() {}

