package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/alecthomas/kong"
)

//   ▄████  ▒█████   ██ ▄█▀
//  ██▒ ▀█▒▒██▒  ██▒ ██▄█▒
// ▒██░▄▄▄░▒██░  ██▒▓███▄░
// ░▓█  ██▓▒██   ██░▓██ █▄
// ░▒▓███▀▒░ ████▓▒░▒██▒ █▄
//  ░▒   ▒ ░ ▒░▒░▒░ ▒ ▒▒ ▓▒
//   ░   ░   ░ ▒ ▒░ ░ ░▒ ▒░
// ░ ░   ░ ░ ░ ░ ▒  ░ ░░ ░
//       ░     ░ ░  ░  ░
// Gok: Reverse shell handler

const VERSION = "0.0"

type Core struct {
    mu sync.RWMutex
    Config    Config
    listeners map[string]*Listener
    sessions  map[int]*Session

    // Event Channels
    newSession chan *Session
    input chan string
}

type Listener struct {
    // id int
    address string
    port     int
    listener net.Listener
}

type Session struct {
    ID int
    Conn net.Conn
    Addr string
}

// global core instance var ??
// var core *Core

func NewCore(cfg Config) *Core {
    return &Core{
	Config:    cfg,
	listeners: make(map[string]*Listener),
	sessions:  make(map[int]*Session),
	// channels
	newSession: make(chan *Session),
	input: make(chan string),
    }
}

func (c *Core) InitListeners() {
    fmt.Printf("[+] Initializing listeners:\n\t")

    for _, addr := range c.Config.bindIps {
	for _, port := range c.Config.PortRange.Ports {
	    fmt.Printf("%s:%d ", addr, port)

	    l, err := c.StartListener(addr,port)
	    if err != nil {
		log.Printf("[-] Failed to start listener on %s:%d: %v", addr, port, err)
		continue
	    }
	    c.mu.Lock()
	    c.listeners[fmt.Sprintf("%s:%d",addr,port)] = l
	    c.mu.Unlock()
	}

    }
    fmt.Printf("\n[*] Waiting for connections...\n")
}

func (c *Core) StartListener(addr string, port int) (*Listener, error) {

    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
    if err != nil {
	return nil, fmt.Errorf("Failed to start listener: %v", err)
    }

    l := &Listener {
	port: port,
	listener: listener,
    }

    go func() {
	defer listener.Close()
	for {
	    conn, err := listener.Accept()
	    if err != nil {
		log.Printf("[-] Error accepting connection: %v", err)
		continue
	    }

	    // create session
	    session := &Session{
		// TODO: FAKE-IT:
		ID: 42,
		Conn: conn,
		Addr: conn.RemoteAddr().String(),
	    }

	    c.mu.Lock()
	    c.sessions[session.ID] = session
	    c.mu.Unlock()

	    // announce to channel
	    c.newSession <- session

	}
    }()

    return l, nil
}

func InteractiveShell(conn net.Conn) {
    done := make(chan bool) // Channel for coordination

    // Goroutine 1: Remote -> Local
    go func() {
	io.Copy(os.Stdout, conn) // Copy everything from conn to stdout
	done <- true             // Signal we're done
    }()

    // Goroutine 2: Local <- Remote
    go func() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
	    line := scanner.Text() + "\n"
	    conn.Write([]byte(line)) // Send to remote shell
	}
	done <- true // Signal we're done
    }()

    <-done // Block until one goroutine finishes

}

func handleSession(conn net.Conn) {
    defer conn.Close()
    log.Printf("[+] Accepting connection from: %v", conn.RemoteAddr())

    // TODO: FAKE-IT

    InteractiveShell(conn)
}

// type Event struct {
//     Type string
//     Payload interface{}
// }

func (c *Core) RunREPL() {
    // read user input
    go func() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
	    c.input <- scanner.Text()
	}
    }()

    fmt.Print("GOK > ")
    for {
	select {

	// session Channels:
	case newSession := <-c.newSession:
	    fmt.Printf("\n[+] New session #%d from %s\n", newSession.ID, newSession.Addr)
	    fmt.Print("GOK > ")

	// case

	// User input channels
	case input := <-c.input:
	    c.handleCmd(input)

	    fmt.Print("GOK > ")
	}

    }

}

func (c *Core) handleCmd(input string) {
    // split on all whitespace
    args := strings.Fields(input)
    if len(args) == 0 {
	return	
    }

    cmd := args[0]
    switch cmd {
    // Management
    case "listeners", "lis", "l":
	c.mu.RLock()
	if len(c.listeners) == 0 {
	    fmt.Println("[!] No active listeners")
	} else {
	    fmt.Println("\nListeners:")
	    for lis := range c.listeners {
		fmt.Printf("%v\n", lis)
	    }
	}
	c.mu.RUnlock()

    case "sessions", "sesh","sess", "s":
	c.mu.RLock()
	if len(c.sessions) == 0 {
	    fmt.Println("[!] No active sessions")
	} else {
	    fmt.Println("\nActive Sessions:")
	    for id, sess := range c.sessions {
	    	fmt.Printf("\t[%d] %s\n", id, sess.Addr)
	    }
	}
	c.mu.RUnlock()

    case "interact","i":
	c.mu.RLock()
	if len(c.sessions) == 0 {
	    fmt.Println("[!] No active sessions")
	} else {
	    fmt.Println("\nActive Sessions:")
	    for id, sess := range c.sessions {
	    	fmt.Printf("\t[%d] %s\n", id, sess.Addr)
	    }
	}
	c.mu.RUnlock()

    default:
	fmt.Printf("[-] Unknown command: %s\n", cmd)
    }
}

func main() {
    fmt.Println()
    fmt.Println("\tGOK: Reverse Shell Handler")
    fmt.Println("")
    fmt.Println(
	"\t  ▄████  ▒█████   ██ ▄█▀\n",
	"\t ██▒ ▀█▒▒██▒  ██▒ ██▄█▒ \n",
	"\t▒██░▄▄▄░▒██░  ██▒▓███▄░ \n",
	"\t░▓█  ██▓▒██   ██░▓██ █▄ \n",
	"\t░▒▓███▀▒░ ████▓▒░▒██▒ █▄\n",
	"\t ░▒   ▒ ░ ▒░▒░▒░ ▒ ▒▒ ▓▒\n",
	"\t  ░   ░   ░ ▒ ▒░ ░ ░▒ ▒░\n",
	"\t░ ░   ░ ░ ░ ░ ▒  ░ ░░ ░ \n",
	"\t      ░     ░ ░  ░  ░   \n",
	"")
    fmt.Printf("\tVersion: %s (Alpha)", VERSION)
    fmt.Println()

    kong.Parse(&Flags)

    // log.Printf("Flags: %v", Flags)

    config := Config{
	PortRange: Flags.PortRange,
	bindIps:   Flags.BoundIPs,
    }

    core := NewCore(config)
    core.InitListeners()

    core.RunREPL()

    // done := make(chan bool)
    // <-done  // Will wait forever (nothing sends to it)
}
