package upgrader

import (
	"os"

	"golang.org/x/term"
)


func Upgrade() {
	// get pty upgrade payload (e.g. python, socat ...)
	// ptyUpgradePayload := Python3{}.GetPayload()
	
}


// run pty spawner


// export env's 

// func exportENV(bashPath string) []string {
// 	res := []string{
// 		"export SHELL=" + bashPath,
// 		"export TERM=xterm-256color",
// 	}
// 	return res
// }

// get cols and rows

func GetTTYSize() (int, int, error) {
	cols,rows,err := term.GetSize(int(os.Stdin.Fd()))

	return rows, cols, err
}

// set cols and rows

// listener on window change
// requires developing a custom protocol for the best outcome
