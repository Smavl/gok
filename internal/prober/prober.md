Prober
===

## Features

- Dynamically-built information gathering phases built on prior gathered intel
- OS Detection
- "Binary" enumeration to be used in payloads and shell upgrader
- Extendable and 
- (WIP) Capabilties (+ derived capabilties)
- (WIP) Agressive enumeration with explotation
- (WIP) 

## Design

Prober is attached to newly spawned session.

The prober is defined by it's phases. 
- First a Genesis phase (which for now only has OS Detection). 
- Later phases are: Initial, Recon, Deep

The phases are defined by their configuration. For now there is _Default_,_Agressive_, _Stealth_ (only _Default_ is implemented). 
These configurations determine what operations that are run in each phase. E.g. in the initial phase there is an enumeration of binaries on the target.

The phases are designed such that the information that they may gather could prove useful in later phases (or operations). Thus they are dynamically built using just that information, and in theory `gok` behaves different for every revshell landed.

The operations should be general and the specifics of the actual operation run the concrete target will depend on what strategies are selected for that operation.

At the moment the prober is only triggered when the shell lands and the session is started. This might be changed to be more adjustable to accommodate for stability and usability concerns.

The prober should also provide the user with some capabilties of the specific session to convey limitations and opportunities. Like how can exfiltration and file transfer be done? Only via in-process base64 conversion, or also via a curl-based approach, which could be better for bigger files. Capabilties also should help `gok` know which utilities and execution of payloads are possible, and improve flexibility in terms of turning to fallbacks if needed.

## Purpose

The purpose of the prober is to gather information about the target system.

This may be relevant for later commands and other utilities provided by `gok`. Such as the automatic shell-upgrader, file transfer capabilties, and more.

A basic implementation of the prober is necessary for providing some of the core features above. However, a more involved implementation should allow for more advanced usage of `gok`.

The prober was initial meant to be delegated the work of OS detection and detecting tools for auto-upgrading the shell (e.g. is Python installed). 
However I wanted the prober to provide insights that could help other features of `gok` such as transfering files, running payloads and more.

