package terminal

import (
	"github.com/smavl/gok/internal/event"
)

type InputReader interface {
	HandleByte(byte) error
}

type LineReader struct {
	lineBuffer []byte
	outChan    chan<- event.MenuCmdEvent
}

func NewLineReader(outChan chan<- event.MenuCmdEvent) *LineReader {
	return &LineReader{
		lineBuffer: []byte{},
		outChan:    outChan,
	}
}

func (r *LineReader) HandleByte(b byte) error {
	// Build line until we get \n
	if b == '\n' {
		line := string(r.lineBuffer)
		r.lineBuffer = []byte{}
		// return event.MenuCmdEvent{Input: line}, nil
		r.outChan <- event.MenuCmdEvent{Input: line}
		return nil
	} else if b != '\r' { // Ignore carriage returns
		r.lineBuffer = append(r.lineBuffer, b)
	}
	// No complete line yet
	return nil
}

type ByteReader struct {
	outChan chan<- event.ShellByteEvent
}

func NewByteReader(outChan chan<- event.ShellByteEvent) *ByteReader {
	return &ByteReader{
		outChan: outChan,
	}
}

func (r *ByteReader) HandleByte(b byte) error {
	r.outChan <- event.ShellByteEvent{Byte: b}
	return nil
}
