package misc

import "errors"

var (
	ErrSessionNotFound       = errors.New("Session not found")
	ErrSessionDed            = errors.New("Session is dead!")
	CouldNotDetermineOSError = errors.New("Could Not Determine OS!")
	NoProberForOs			 = errors.New("No Prober for OS!")
	ErrNoProberMode			 = errors.New("No Prober for Mode!")
)
