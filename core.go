package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

// TODO: string??

func (c *Core) EnableMainMenuMode() {
	// very nice go, no options!
	c.activeShellID = -1
}

func (c *Core) EnableShellMode() {
	session, _ := c.SessionManager.Get(c.activeShellID)
	// TODO: error handling

	// start reader (only once)
	session.mu.Lock()
	if !session.started {
		session.started = true
		// TODO: Revisit. Might look funny if active session dies in the background
		go io.Copy(os.Stdout,session.Conn)
	}
	session.mu.Unlock()
	// drop into shell (blocking)
	c.runShellReader()
}

type Core struct {
	mu     sync.RWMutex
	Config Config
	listeners map[string]*Listener

	// sessions  map[int]*Session
	SessionManager *SessionManager

	// shell
	activeShellID int

	// Event Channels
	newSession chan *Session
	// sessionDied chan *Session
}

type Listener struct {
	// id int
	address  string
	port     int
	listener net.Listener
}



func NewCore(cfg Config) *Core {
	return &Core{
		Config:    cfg,
		// managers
		listeners: make(map[string]*Listener),
		SessionManager: NewSessionManager(),

		// channels
		newSession: make(chan *Session),
		// sessionDied: make(chan *Session),
	}
}

func (c *Core) InitListeners() {
	fmt.Printf("[+] Initializing listeners:\n\t")

	for _, addr := range c.Config.bindIps {
		for _, port := range c.Config.PortRange.Ports {
			fmt.Printf("%s:%d ", addr, port)

			l, err := c.StartListener(addr, port)
			if err != nil {
				log.Printf("[-] Failed to start listener on %s:%d: %v", addr, port, err)
				continue
			}
			c.mu.Lock()
			c.listeners[fmt.Sprintf("%s:%d", addr, port)] = l
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

	l := &Listener{
		port:     port,
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
			session := c.SessionManager.AddSession(conn)


			// announce to channel
			c.newSession <- session

		}
	}()

	return l, nil
}

// read handlers
// Menu reader (line-buffered)
func (c *Core) runMenuReader() {
	fmt.Print("GOK > ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		c.handleMainMenu(scanner.Text())
	}
}

// Shell reader (raw terminal ):
// TODO: add x/term later
func (c *Core) runShellReader() {
	session, _ := c.SessionManager.Get(c.activeShellID)
	// TODO: error handling

	// workaround for resuming session
	session.Conn.Write([]byte("\n"))

	// For now, simple line-based with escape check
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		// Check for escape
		if input == "exit" || input == "~~~" {
			c.EnableMainMenuMode()
			return  // Exit this reader
		}

		// Send to remote
		session.Conn.Write([]byte(input + "\n"))
	}
}



func (c *Core) RunREPL() {
	// Start reader
	go c.runMenuReader()

	for {
		select {

		// session Channels:
		case sess := <-c.newSession:
			fmt.Printf("\n[+] New session #%d from %s\n", sess.ID, sess.Addr)
			fmt.Print("GOK > ")

		// User input channels
		// case input := <-c.input:
		// 	c.handleMainMenu(input)

		}

	}

}

func (c *Core) handleMainMenu(input string) {
	// split on all whitespace
	args := strings.Fields(input)
	lenArgs := len(args)
	if lenArgs == 0 {
		return
	}
	defer fmt.Print("\nGOK > ")

	subCmd := args[0]
	switch subCmd {
	// Management
	case "listeners", "lis", "l":
		c.mu.RLock()
		if len(c.listeners) == 0 {
			fmt.Println("[!] No active listeners")
		} else {
			fmt.Println("\nListeners:")
			for lis := range c.listeners {
				fmt.Printf("- %v\n", lis)
			}
		}
		c.mu.RUnlock()

	case "sessions", "sesh", "sess", "s":
		if c.SessionManager.GetAmount() == 0 {
			fmt.Println("\n[!] No active sessions")
		} else {
			fmt.Println("\nActive Sessions:")
			for _, sess := range c.SessionManager.GetSessions() {
				fmt.Printf("\t[%d] %s\n", sess.ID, sess.Addr)
			}
		}
	// fmt.Print("\nGOK > ")

	case "interact", "int", "i":
		// check if session arg is supplied
		if len(args) == 2 {
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Printf("Id: %v is not a number",id)
			}

			sessionExists := c.SessionManager.Exists(id)

			if sessionExists {
				c.activeShellID = id
				fmt.Printf("[*] Session #%v: Dropping into shell..\n",id)
				c.EnableShellMode()
			} else {
				fmt.Printf("Session #%d not found",id)
			}

		} else {
			fmt.Println("[!] No session chosen, or invalid argument")
		}

	case "exit", "quit":
		// TODO: Need to do more?
		os.Exit(0)

	default:
		// TODO: Add help suggestion
		fmt.Printf("[!] Unknown command: %s", subCmd)
	}
}
