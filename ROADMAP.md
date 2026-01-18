# Roadmap (Short term)

- [x] Minimal POC for single connection
- [ ] CLI flags
    - [x] Switch to kong? [kong](https://github.com/alecthomas/kong) 
    - [ ] Automatically drop into shell (default)
- [x] Add features to Core 
    - [x] Listeners
    - [x] Sessions
- [x] Add simple menu handling (Channels?)
    - [x] View Sessions, and listners
    - [x] Interact with session
    - [x] Kill sessions
    - [x] Add help menu
- [x] testing "Setup" 
    - [ ] Docker setup?
- [ ] Testing
    - [x] Parsing of CLI flags (add option to "show flags" to unit test)
    - [ ] More streamlined testing setup
- [ ] Interactive shell handling
    - [x] Add escaping from session
    - [ ] Add "meta mode" inside active shell to run payloads and 
- [ ] Menu actions
    - [ ] Session details (info from the prober?)


# TODO (Long term)

## Shell 

- [x] OS detection
    - [x] linux
    - [ ] windows
- [ ] Session handling
    - [x] print history after entering session again
- [ ] Prober
    - [x] Detect `which` binary
    - [x] Find binaries
    - [ ] EnumerateUser(s)
    - [ ] Parse env
- [ ] Automatic shell upgrader (uses Prober?)
    - [ ] python
    - [ ] script
    - [ ] perl
    - [ ] php
    - [ ] ...
- [ ] Dynamic
    - [ ] Terminal resizing
- [ ] Obfuscation?
    - [ ] download, upload, shell upgrade?
- [ ] 

## Test Scenarios

- [x] Port range is parsed correctly
- [ ] Listen on range -> multiple sessions -> Change between sessions
    - [x] Simple testing
    - [ ] Test with multiple real sessions
- [ ] Upload file

## Utility

- [ ] Uploading files
    - [ ] curl, wget
    - [ ] nc 
    - [ ] in-memory base64?
- [ ] Running payloads 
    - [ ] Linpeas
    - [ ] Support for home-brought scripts
- [ ] Download files
    - [ ] -||-
- [ ] Logging?

## UX 

- [ ] Print all ips gok is listening on when ip is "0.0.0.0"
- [ ] Print Different payloads!

### CLI flags
- [x] Ports: "-p 9001", "-p 9001-9009" (Maybe "-p 9001 9003"?)

