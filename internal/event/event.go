package event

// Events 
type NewSessionEvent struct {
	SessionID int
	SessionAddr string
	SystemOS string
}
type ShellByteEvent struct {
	Byte byte
}

type MenuCmdEvent struct {
	Input string
}

type EventBus struct {
	Session chan NewSessionEvent
	Shell chan ShellByteEvent
	Menu chan MenuCmdEvent
}

func NewEventBus() *EventBus {
	return &EventBus{
		// TODO: what is `100` magic number? and what should it be?
		Session: make(chan NewSessionEvent, 100),
		Shell:   make(chan ShellByteEvent, 100),
		Menu:    make(chan MenuCmdEvent, 100),
	}
}
