package main

import "errors"


var (
	ErrSessionNotFound = errors.New("Session not found")
	ErrSessionDed = errors.New("Session is dead")
)
