package output

import "golang.org/x/term"

func isTTY(fd int) bool {
	return term.IsTerminal(fd)
}
