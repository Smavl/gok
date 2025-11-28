# Roadmap (Short term)

- [x] Minimal POC for single connection
- [ ] CLI flags
    - [ ] Switch to kong? [kong](https://github.com/alecthomas/kong) 
- [ ] Add features to Core 
    - [ ] Listeners
    - [ ] Sessions
- [ ] Add simple menu handling (Channels?)
    - [ ] 
- [ ] Testing
    - [ ] Parsing of CLI flags (add option to "show flags" to unit test)
    - [ ] More streamlined testing setup?
- [ ] Reflect on architechture 
    - [ ] Need more ..Handlers? (Session, Listener)
    - [ ] Refactor/simplify?
- [ ] Interactive shell handling
- [ ] 


# TODO (Long term)

## Shell 

- [ ] OS detection
- [ ] Prober
    - [ ] Find binaries
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

- [ ] Listen on range -> multiple sessions -> Change between sessions
- [ ] Upload file
- [ ] Upload file

### Utility

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
- [ ] 

## UX 

- [ ] Print all ips gok is listening on when ip is "0.0.0.0"
- [ ] Print Different payloads!

### CLI flags
- [ ] Ports: "-p 9001", "-p 9001-9009" (Maybe "-p 9001 9003"?)

