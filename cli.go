package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

var Flags struct {
	PortRange PortRange `help:"Ports to listen on" default:"9001" short:"p"`
	BoundIPs  []string  `help:"IPs to bind the listeners on" default:"[0.0.0.0]" short:"b"`
}

type Config struct {
	bindIps   []string
	PortRange PortRange
}

type PortRange struct {
	Ports []int
}

// Custom parsing
// TODO: Add support for muliple
func (p *PortRange) Decode(ctx *kong.DecodeContext) error {
	var value string

	if err := ctx.Scan.PopValueInto("ports", &value); err != nil {
		return err
	}

	if strings.Contains(value, "-") {
		split := strings.SplitN(value, "-", 2)

		// check if range parsed!
		if len(split) != 2 {
			return fmt.Errorf("Invalid port range: %s", value)
		}

		//
		start, err := strconv.Atoi(split[0])
		if err != nil {
			return fmt.Errorf("Invalid range: %s", split[0])
		}

		end, err := strconv.Atoi(split[1])
		if err != nil {
			return fmt.Errorf("Invalid end port: %s", split[1])
		}

		if start > end {
			return fmt.Errorf("Start port %d is greater than end port %d", start, end)
		}
		count := end - start + 1
		p.Ports = make([]int, count)

		for i := range p.Ports {
			p.Ports[i] = start + i
		}
	} else {
		// single port
		port, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("Invalid port value: %v", err)
		}

		p.Ports = []int{port}

	}
	return nil

}
