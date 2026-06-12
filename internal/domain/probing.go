package domain

import (
	"fmt"
	"strconv"

	"github.com/alecthomas/kong"
)

type ProbingOptions struct {
	ProbingMode   ProbingMode
	DisableProber bool
	// TimeoutPerOperation time.Duration
}

type ProbingMode int

const (
	Default ProbingMode = iota
	Agressive
	Stealth
)

func isValidProbeMode(pm int) bool {
	// check if int is a valid probing mode
	if pm < 0 || pm > 2 {
		return false
	}
	return true
}

func (p *ProbingMode) Decode(ctx *kong.DecodeContext) error {
	var value string

	if err := ctx.Scan.PopValueInto("probing-mode", &value); err != nil {
		return err
	}

	// cast string to int
	pmInt, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("Invalid probing mode value, cast to int: %v", err)
	}

	// check if valid mode
	if !isValidProbeMode(pmInt) {
		return fmt.Errorf("Invalid probing mode: %d. Valid modes are 0 (Default), 1 (Agressive), 2 (Stealth)", pmInt)
	}

	// cast to ProbingMode
	*p = ProbingMode(pmInt)
	return nil
}
