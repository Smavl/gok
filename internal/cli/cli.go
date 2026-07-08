package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/smavl/gok/internal/domain"
)

var Flags struct {
	PortRange PortRange `help:"Ports to listen on" default:"9001" short:"p"`
	BoundIPs  []string  `help:"IPs to bind the listeners on" default:"[0.0.0.0]" short:"b"`
	// timeout flags
	ProbingCmdTimeout time.Duration `help:"Timeout for probing commands" default:"200ms" short:"t"`
	ProbingMode domain.ProbingMode `help:"Level of agressiveness for the prober" default:"0" short:"A"`
	DisableProber bool `help:"Disable the prober" default:"false" short:"D"`
	AutoInteract bool `help:"Automatically 'interact' with session upon successful connection" default:"false" short:"F"` // "F" for "fast" ?
}

type Config struct {
	// listener config
	BindIps				[]string
	PortRange			PortRange
	// probing config
	ProbingCmdTimeout	time.Duration
	ProbingMode			domain.ProbingMode
	// TODO: add logic validation
	DisableProber		bool

	// session settings
	AutoInteract		bool

	// misc
	// TODO: testmode/headless mode
	HeadlessMode bool
}

type PortRange struct {
	Ports []int
}


func isValidPort(port int) bool {
	minValidPort := 1
	maxValidPort := 65535

	if port < minValidPort {return false}
	if port > maxValidPort {return false}

	return true
}

// Custom parsing
// TODO: Add support for muliple?
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
			currentPort := start + i
			if !isValidPort(currentPort) {
				return fmt.Errorf("Invalid port in range: %d", currentPort)
			}
			p.Ports[i] = currentPort
		}
	} else {
		// single port
		port, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("Invalid non-integer value: %v", err)
		}
		if !isValidPort(port) {
			return fmt.Errorf("Invalid port: %d", port)
		}

		p.Ports = []int{port}

	}
	return nil
}

