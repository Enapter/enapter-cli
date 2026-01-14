package enaptercli

import (
	"io"
	"os"

	"golang.org/x/term"
)

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

func colorsSupported(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	if !term.IsTerminal(int(f.Fd())) {
		return false
	}
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	return true
}
