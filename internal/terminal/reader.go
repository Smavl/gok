package terminal

import (
	"os"

	"github.com/smavl/gok/internal/event"
)

type InputReader interface {
	Read() error
}

type LineReader struct {
	lineBuffer []byte
	buf        [1]byte
	outChan chan<- event.MenuCmdEvent
}

func NewLineReader(outChan chan<- event.MenuCmdEvent) *LineReader {
	return &LineReader{
		lineBuffer: []byte{},
		outChan: outChan,
	}
}

func (r *LineReader) Read() error {
	n, err := os.Stdin.Read(r.buf[:])
	if err != nil || n == 0 {
		return err
	}

	// Build line until we get \n
	if r.buf[0] == '\n' {
		line := string(r.lineBuffer)
		r.lineBuffer = []byte{}
		// return event.MenuCmdEvent{Input: line}, nil
		r.outChan <- event.MenuCmdEvent{Input: line}
		return  nil
	} else if r.buf[0] != '\r' { // Ignore carriage returns
		r.lineBuffer = append(r.lineBuffer, r.buf[0])
	}
	// No complete line yet
	return nil
}

type ByteReader struct {
	buf [1]byte
	outChan chan<- event.ShellByteEvent
	
}

func NewByteReader(outChan chan<- event.ShellByteEvent) *ByteReader {
	return &ByteReader{
		outChan: outChan,
	}
}

func (r *ByteReader) Read() error {
	n, err := os.Stdin.Read(r.buf[:])
	if err != nil || n == 0 {
		return  err
	}
	r.outChan <- event.ShellByteEvent{Byte: r.buf[0]}
	return nil
}
