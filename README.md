Gok
===

```
   ▄████  ▒█████   ██ ▄█▀
  ██▒ ▀█▒▒██▒  ██▒ ██▄█▒
 ▒██░▄▄▄░▒██░  ██▒▓███▄░
 ░▓█  ██▓▒██   ██░▓██ █▄
 ░▒▓███▀▒░ ████▓▒░▒██▒ █▄
  ░▒   ▒ ░ ▒░▒░▒░ ▒ ▒▒ ▓▒
   ░   ░   ░ ▒ ▒░ ░ ░▒ ▒░
 ░ ░   ░ ░ ░ ░ ▒  ░ ░░ ░
       ░     ░ ░  ░  ░
```

Gok is a feature-rich (reverse) shell handler and offensive utility tool. 

Gok aims to replace the muscle memory of `nc -lvnp 9001` with `gok -p 9001`, and to fast-track the exploitation and/or enumeration process immediately after the shell lands. By providing essential utilities and features, it remains useful beyond landing and upgrading the shell.

# Features

**Core features**
- Manage multiple listeners and sessions
  - Specify multiple IP addresses and ports
  - Jump between multiple sessions (History persistance)
- Automatic shell upgrading (WIP)
- Prober: Automatically gather information
- (Somewhat) Asynchronous and event-driven architecture
- Cross-platform support (Linux, Windows, MacOS) (WIP)
- Multiplexing

## Interfactive Modes
- Menu mode - Management and overview
  - Command interface (Manage: sessions, listeners)
  - Manage multiple reverse shell sessions
  - Session interaction switching and escaping
- Shell mode - Raw shell interaction
  - Byte-by-byte shell interaction (no line-buffering)
  - Escape and resume sessions seamlessly.
- Meta mode (WIP)
  - Run scripts/modules
  - Utility: download, upload tools

## Prober
**Modes**: Default, Agressive (WIP), Stealth (WIP)

Default:
- OS detection
- Available binaries enumeration

Aggressive (WIP):
- Elaborate binary enumeration
- File system enumeration
- User enumeration
- SUID enumeration
- ...

Stealth (WIP):
- ... (TBD)

# Installation

## **Requirements:**
- `go`

Optional:
- `docker` (for integration tests)

## Running Gok

**Build from source:**
```bash
# Clone repo
git clone https://github.com/smavl/gok
cd gok
# build binary
go build cmd/gok/main.go
./main
# Or run directly with go
go run cmd/gok/main.go
```
**one-line: (WIP)**
```bash
#go install github.com/smavl/gok@latest
#gok -p 9001
```

# Using Gok
## Basic Usage

**Specify listener(s):**
```bash
$ gok -p 9001              # Single port
$ gok -p 9001-9005         # Port range
$ gok -p 9001 -b 0.0.0.0   # Specify bind address for listener(s)
```


**Menu commands:**
```
GOK > help
Available Commands:
  listeners, lis, l         - List all active listeners
  sessions, sesh, sess, s   - List all active sessions
  interact, int, i <id>     - Interact with a session
  kill, k <id>              - Kill a session
  help, h                   - Show this help message
  exit, quit, q             - Exit the application
```

**Session interaction:**
- Drop into the shell: `interact n`, `int n`, `i n`
- Escape by hitting `Ctrl+D` to background the session and return to menu mode


# Testing

You can run the tests with:
```bash
# run "unit" tests
$ go test ./...
# Run dockerized integration tests (also)
$ go test -tags=integration ./... -v
```

# Kudos

This was inspired by [penelope](https://github.com/brightio/penelope) which I encourage all to check out!
