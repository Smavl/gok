package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

// TODO: string??

func (c *Core) EnableShellMode() {
	session, _ := c.SessionManager.Get(c.activeShellID)
	// TODO: error handling


	// Bring session to foreground (redirect buffer output to "stdout")
	session.Foreground()
	oldMenuMode := c.Mode
	c.Mode = &ShellMode{}

	// drop into shell (blocking)
	c.runShellReader()

	// When session is escaped: background the session
	session.Background()
	c.Mode = oldMenuMode
}

type Core struct {
	mu     sync.RWMutex
	Config Config
	listeners map[string]*Listener

	// sessions  map[int]*Session
	SessionManager *SessionManager

	Display Display
	Mode Mode

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
	terminalDisplay := NewTerminalDisplay()
	// NOTE: Do we need a proper constructor?
	menuMode := MenuMode{
		Display: terminalDisplay,
	}
	return &Core{
		Config:    cfg,
		// managers
		listeners: make(map[string]*Listener),
		SessionManager: NewSessionManager(terminalDisplay),

		Display: terminalDisplay,
		Mode: &menuMode,

		// channels
		newSession: make(chan *Session),
		// sessionDied: make(chan *Session),
	}
}

func (c *Core) InitListeners() {
	c.Message("[+] Initializing listeners:\n\t")

	for _, addr := range c.Config.bindIps {
		for _, port := range c.Config.PortRange.Ports {
			c.Message("%s:%d ", addr, port)

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
	c.Message("\n[*] Waiting for connections...\n")
}

func (c *Core) Prompt() {
	c.Display.Prompt()
}

func (c *Core) Message(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	c.Display.Message(msg)
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
	c.Prompt()
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
	session.conn.Write([]byte("\n"))

	// For now, simple line-based with escape check
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		// Check for escape
		if input == "exit" || input == "~~~" {
			return  // Exit this reader
		}

		// Send to remote
		session.conn.Write([]byte(input + "\n"))
	}
}


func (c *Core) RunREPL() {
	// Start reader
	go c.runMenuReader()

	for {
		select {

		// session Channels:
		case sess := <-c.newSession:
			c.Mode.OnNewSession(sess)

		// User input channels
		// case input := <-c.input:
		// 	c.handleMainMenu(input)

		}

	}

}

func (c *Core) handleMainMenu(input string) {
	// split on all whitespace
	args := strings.Fields(input)
	
	defer c.Prompt()

	if len(args) == 0 {
		return
	}

	subCmd := args[0]
	switch subCmd {
	// Management
	case "listeners", "lis", "l":
		c.mu.RLock()
		if len(c.listeners) == 0 {
			c.Message("[!] No active listeners\n")
		} else {
			c.Message("\nListeners:\n")
			for lis := range c.listeners {
				c.Message("- %v\n", lis)
			}
		}
		c.mu.RUnlock()

	case "sessions", "sesh", "sess", "s":
		if c.SessionManager.GetAmount() == 0 {
			c.Message("\n[!] No active sessions\n")
		} else {
			c.Message("\nActive Sessions:\n")
			for _, sess := range c.SessionManager.GetSessions() {
				c.Message("\t[%d] %s\n", sess.ID, sess.Addr)
			}
		}
	// c.Prompt()

	case "interact", "int", "i":
		// check if session arg is supplied
		if len(args) == 2 {
			id, err := strconv.Atoi(args[1])
			if err != nil {
				c.Message("Id: %v is not a number\n",id)
			}

			sessionExists := c.SessionManager.Exists(id)

			if sessionExists {
				c.activeShellID = id
				c.Message("[*] Session #%v: Dropping into shell..\n",id)
				c.EnableShellMode()
			} else {
				c.Message("Session #%d not found\n",id)
			}

		} else {
			c.Message("[!] No session chosen, or invalid argument\n")
		}

	case "exit", "quit":
		// TODO: Need to do more?
		os.Exit(0)

	default:
		// TODO: Add help suggestion
		c.Message("[!] Unknown command: %s\n", subCmd)
	}
}
