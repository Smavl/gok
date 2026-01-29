package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/smavl/gok/internal/cli"
	"github.com/smavl/gok/internal/core"
)

//   ▄████  ▒█████   ██ ▄█▀
//  ██▒ ▀█▒▒██▒  ██▒ ██▄█▒
// ▒██░▄▄▄░▒██░  ██▒▓███▄░
// ░▓█  ██▓▒██   ██░▓██ █▄
// ░▒▓███▀▒░ ████▓▒░▒██▒ █▄
//  ░▒   ▒ ░ ▒░▒░▒░ ▒ ▒▒ ▓▒
//   ░   ░   ░ ▒ ▒░ ░ ░▒ ▒░
// ░ ░   ░ ░ ░ ░ ▒  ░ ░░ ░
//       ░     ░ ░  ░  ░
// Gok: Reverse shell handler

const VERSION = "0.1"

func main() {
	fmt.Println()
	fmt.Println("\tGOK: Reverse Shell Handler")
	fmt.Println("")
	fmt.Println(
		"\t  ▄████  ▒█████   ██ ▄█▀\n",
		"\t ██▒ ▀█▒▒██▒  ██▒ ██▄█▒ \n",
		"\t▒██░▄▄▄░▒██░  ██▒▓███▄░ \n",
		"\t░▓█  ██▓▒██   ██░▓██ █▄ \n",
		"\t░▒▓███▀▒░ ████▓▒░▒██▒ █▄\n",
		"\t ░▒   ▒ ░ ▒░▒░▒░ ▒ ▒▒ ▓▒\n",
		"\t  ░   ░   ░ ▒ ▒░ ░ ░▒ ▒░\n",
		"\t░ ░   ░ ░ ░ ░ ▒  ░ ░░ ░ \n",
		"\t      ░     ░ ░  ░  ░   \n",
		"")
	fmt.Printf("\tVersion: %s", VERSION)
	fmt.Println()

	kong.Parse(&cli.Flags)

	config := cli.Config{
		PortRange:         cli.Flags.PortRange,
		BindIps:           cli.Flags.BoundIPs,
		ProbingCmdTimeout: cli.Flags.ProbingCmdTimeout,
		ProbingMode:       cli.Flags.ProbingMode,
	}

	c := core.NewCore(config)

	c.Start()
}
