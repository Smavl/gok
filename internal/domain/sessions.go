package domain

type Session interface {
	GetID() int
	Foreground()
	Background()
	Write([]byte) (int, error)
}
