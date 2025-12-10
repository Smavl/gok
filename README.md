# Gok - (Reverse) Shell Handler

## Features

- Listen on multiple ports simultaneously
- Manage multiple reverse shell sessions
- Session switching and escaping
- Simple command interface (sessions, listeners, interact)

## Installation

**Build from source:**
```bash
git clone https://github.com/smavl/gok
cd gok
go build
./gok -p 9001
```

**Requirements:**
- `go`

## Basic usage

**Start listener:**
```bash
gok -p 9001              # Single port
gok -p 9001-9005         # Port range
gok -p 9001 -b 0.0.0.0   # Specify bind address
```

**Menu commands:**
- `sessions` (or `s`) - List active sessions
- `listeners` (or `l`) - List active listeners
- `interact <id>` (or `i <id>`) - Interact with session


**Shell interaction**
- Type `exit` or `~~~` to escape and return to menu




Inspired by [penelope](https://github.com/brightio/penelope) 
