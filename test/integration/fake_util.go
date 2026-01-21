package main

import (
	"os"
	"sync"
)

// type FakeConn struct {
// 	readBuf  *bytes.Buffer
// 	writeBuf *bytes.Buffer
// }
//
// func NewFakeConn() FakeConn{
// 	return FakeConn{}
// }
//
// func (fc *FakeConn) Close()


type FakeDisplay struct{
	mu sync.Mutex
}

func (d *FakeDisplay) Write(b []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return os.Stdout.Write(b)
}

func NewFakeDisplay() *FakeDisplay {
	return &FakeDisplay{}
}
