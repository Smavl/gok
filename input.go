package main

import (
	"os"
)


type InputReader interface {
	Read() Event
	Stop()
}

type LineReader struct {
	lineBuffer []byte
	buf        [1]byte
}

func NewLineReader() *LineReader {
	return &LineReader{
		lineBuffer: []byte{},
	}
}


func (r *LineReader) Read() Event {
	n, err := os.Stdin.Read(r.buf[:])
	if err != nil || n == 0 {
		return nil
	}

	// Build line until we get \n
	if r.buf[0] == '\n' {
		line := string(r.lineBuffer)
		r.lineBuffer = []byte{}
		return MenuCmdEvent{input: line}
	} else if r.buf[0] != '\r' { // Ignore carriage returns
		r.lineBuffer = append(r.lineBuffer, r.buf[0])
	}
	return nil
}

type ByteReader struct {
	buf [1]byte
}

func NewByteReader() *ByteReader {
	return &ByteReader{}
}


func (r *ByteReader) Read() Event {
	n, err := os.Stdin.Read(r.buf[:])
	if err != nil || n == 0 {
		return nil
	}
	return ShellByteEvent{Byte: r.buf[0]}
}


func (r *LineReader) Stop() {
	// Nothing needed for line reader
}
func (r *ByteReader) Stop() {
	// Nothing needed for byte reader
}
