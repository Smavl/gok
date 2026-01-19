# Roadmap (Short term)

## Version 0.0
- [x] CLI flags
    - [x] Switch to kong? [kong](https://github.com/alecthomas/kong) 
- [x] Add features to Core 
    - [x] Listeners
    - [x] Sessions
- [x] Menu actions
    - [x] View Sessions, and listners
    - [x] Interact with session (+ switching)
- [x] testing "Setup" 
- [x] Testing
    - [x] Parsing of CLI flags (add option to "show flags" to unit test)
    - [x] Light Session testing
    - [x] basic Dockerized integration testing setup
- [x] Interactive shell handling
    - [x] Add escaping from session (to main menu)
- [x] Prober 
    - [x] OS detection (linux atm)
- [x] Prober (linux)
    - [x] `which`
    - [x] EnumerateBinaries
- [x] Headless mode (terminal) - (Only for testing rn)

## Version 0.1

- [x] Timeout refactor
- [ ] CLI flags
    - [x] Timeout flag(s)
    - [ ] Probing flags (modes)
- [ ] Menu actions
    - [ ] Session details (info from the prober?)
- [ ] Prober 
    - [ ] Modes (Default, Agressive, stealth)
    - [ ] Add delimiter to optimize performance of enumeration
- [ ] Prober (linux)
    - [ ] EnumerateUser
    - [ ] EnumerateUsers

## Version 0.2

- [ ] Debugging- and/or general logs
- [ ] CLI flags
- [ ] Meta mode
    - [ ] Add escaping
    - [ ] Implement line-buffered input handler
    - [ ] inside active shell to run payloads and 
- [ ] Automatic shell upgrading 
    - [ ] Simple shell upgrading implementation (python)
    - [ ] CLI flag: Automatically drop into shell (default)

# Roadmap - Long term

## Shell 

- [x] OS detection
    - [x] linux
    - [ ] windows
- [ ] Session handling
    - [x] print history after entering session again
- [ ] Automatic shell upgrader (uses Prober?)
    - [ ] python
    - [ ] script
    - [ ] perl
    - [ ] php
    - [ ] (socat thingy)
    - [ ] ...
- [ ] Dynamic
    - [ ] Terminal resizing
- [ ] Meta-mode 
    - [ ] Escape feature (from shell)

## Prober

- [x] Detect `which` binary
- [x] Find binaries
- [ ] EnumerateUser(s) 
- [ ] Parse env
- [ ] Modes: default, aggressive, (stealthy?)
- [ ] Fingerprinting

## Menu mode

- [ ] Command history (up and down arrows)
- [ ] Auto-completion/tab completion

## Meta mode

### Features 

- [ ] Command history (up and down arrows)
- [ ] Different prompt
- [ ] Auto-completion/tab completion
- [ ] Ability to open new Terminal window
- [ ] Notion of capabilties (using information from the prober?)
- [ ] Session features
    - [ ] Add another session (same target)
        - [ ] Same window/shell and gok instance
        - [ ] New session - New window
- [ ] Uploading files
    - [ ] methods/strategies:
    - [ ] nc 
    - [ ] in-memory base64?
    - [ ] curl, wget
    - [ ] FTP?
- [ ] Download files
    - [ ] -||-

### Modules

- [ ] Module architechture
    - [ ] Should each payload be its own module?

- [ ] Script/Payload runner 
    - [ ] Linpeas (in new terminal window?)
    - [ ] Support for home-brought scripts
- [ ] File system tool - Upload/download
- [ ] Port forwarding/proxy

## UX 

- [ ] Print all ips gok is listening on when ip is "0.0.0.0"
- [ ] Print Different payloads!
- [ ] Implement user config integration 
    - [ ] (like: log file location, flags?, custom module inport, )

## Manual Test Scenarios

- [x] Port range is parsed correctly
- [x] Listen on range -> multiple sessions -> Change between sessions
    - [x] Simple testing
    - [x] Test with multiple session on same machine
    - [ ] Test with multiple real sessions
- [ ] Upload file
- [ ] Download file


## Dockerized tests 

- [x] Simple rev shell. Assertions: sessions = 1, OS = Linux, binaries
- [ ] Better test utils
    - [ ] param for packages that the test container will have

## To be ordered features

- [ ] SSH 
    - [ ] SSH key enumerator
    - [ ] SSH key injector (`echo ... >> PATH/.ssh/authorized_keys`)
    - [ ] Port forwarding
- [ ] Logging
- [ ] Obfuscation?
    - [ ] download, upload, shell upgrade?
- [ ] Docker playground (no `testcontainers-go`)
    - [ ] Template docker file
    - [ ] Demos?
