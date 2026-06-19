package upgrader

type PtySpawner interface{
	GetPayload() string
}

type Python3 struct {}

// type Python2 struct {}
//
// type Socat struct {}
//
// type Script struct {}


func (pty Python3) GetPayload() string {
	payload := `python -c 'import pty; pty.spawn("/bin/bash")'`
	return payload
}
