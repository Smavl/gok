package domain

// InteractiveSession Defines the minimal interface that a session needs to implement
type InteractiveSession interface {
	GetID() int
	Foreground()
	Background()
	Write([]byte) (int, error)
}

// CommandSession defines the minimal interface that prober needs from a session
type CommandSession interface {
	Write([]byte) (int, error)
	GetProbingLines() []string
	ClearProbingBuffer()
	GetProbingDataChannel() <-chan struct{}
}
