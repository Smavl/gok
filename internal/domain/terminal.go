package domain

type TerminalController interface {
	SetRaw()
	Restore()
	Write(b []byte) (int, error)
	Message(format string, args ...any)
	Prompt()
	ShowError(err error)
}
