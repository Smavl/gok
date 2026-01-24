package core

import (
	"os"
	"strconv"
	"strings"

	"github.com/smavl/gok/internal/domain"
	"github.com/smavl/gok/internal/session"
	"github.com/smavl/gok/internal/terminal"
)

type CommandHandler struct {
	sessions  *session.SessionManager
	listeners *session.ShellListenerManager
	terminal  domain.TerminalController
	shellMode *terminal.RawShellMode
}

func NewCommandHandler(sessions *session.SessionManager, listeners *session.ShellListenerManager, term domain.TerminalController, shellMode *terminal.RawShellMode) *CommandHandler {
	return &CommandHandler{
		sessions:  sessions,
		listeners: listeners,
		terminal:  term,
		shellMode: shellMode,
	}
}

func (ch *CommandHandler) listSessions() {
	if ch.sessions.GetAmountOfSessions() == 0 {
		ch.terminal.Message("\n[!] No active sessions\n")
	} else {
		ch.terminal.Message("\nActive Sessions:\n")
		for _, sess := range ch.sessions.GetSessions() {
			ch.terminal.Message("\t[%d] %s\n", sess.ID, sess.Addr)
		}
	}
	defer ch.terminal.Prompt()

}
func (ch *CommandHandler) listListeners() {
	addresses := ch.listeners.GetAddresses()
	if len(addresses) == 0 {
		ch.terminal.Message("[!] No active listeners\n")
	} else {
		ch.terminal.Message("\nListeners:\n")
		for _, addr := range addresses {
			ch.terminal.Message("- %v\n", addr)
		}
	}
	defer ch.terminal.Prompt()
}
func (ch *CommandHandler) interact(args []string) {
	// check if session arg is supplied
	if len(args) == 2 {
		id, err := strconv.Atoi(args[1])
		if err != nil {
			ch.terminal.Message("Id: %v is not a number\n", id)
			ch.terminal.Prompt()
			return
		}

		session, err := ch.sessions.Get(id)

		if err == nil {
			ch.shellMode.Enter(session)
		} else {
			ch.terminal.Message("Session #%d not found\n", id)
			ch.terminal.Prompt()
		}

	} else {
		ch.terminal.Message("[!] No session chosen, or invalid argument\n")
		ch.terminal.Prompt()
	}
}

func (ch *CommandHandler) killSession(args []string) {
	if len(args) != 2 {
		ch.terminal.Message("[!] Usage: kill <session_id>\n")
		return
	}

	id, err := strconv.Atoi(args[1])
	if err != nil {
		ch.terminal.Message("[!] Invalid session ID: %v\n", args[1])
		return
	}

	session, err := ch.sessions.Get(id)
	if err != nil {
		ch.terminal.Message("[!] Session #%d not found\n", id)
		return
	}

	session.Stop()
	ch.terminal.Message("[*] Killed session #%d\n", id)
}

func (ch *CommandHandler) showHelp() {
	helpText := `Available Commands:
  listeners, lis, l         - List all active listeners
  sessions, sesh, sess, s   - List all active sessions
  interact, int, i <id>     - Interact with a session
  kill, k <id>              - Kill a session
  help, h                   - Show this help message
  exit, quit, q             - Exit the application
`
	ch.terminal.Message("%s", helpText)
}

func (ch *CommandHandler) Execute(input string) {
	args := strings.Fields(input)

	if len(args) == 0 {
		ch.terminal.Prompt()
		return
	}

	subCmd := args[0]
	switch subCmd {
	// Management
	case "listeners", "lis", "l":
		ch.listListeners()

	case "sessions", "sesh", "sess", "s":
		ch.listSessions()

	case "interact", "int", "i":
		ch.interact(args)

	case "kill", "k":
		ch.killSession(args)

	case "help", "h":
		ch.showHelp()

	case "exit", "quit", "q":
		os.Exit(0)

	default:
		// TODO: Add help suggestion
		ch.terminal.Message("[!] Unknown command: %s\n", subCmd)
		ch.terminal.Message(`[+] Type "help" to see available commands.`)
	}

	ch.terminal.Prompt()
}
