// Package display renders fid's human output: colour roles, symbols, boxes
// and columns. It is ported from whet's display package, which Mark Bell
// wrote and authorised copying into fid. fid always uses its Unicode
// symbols, and colours only a terminal.
package display

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

// Renderer holds the display settings for one output stream.
type Renderer struct {
	Color bool
}

// NewRenderer colours output only for a terminal, and only while NO_COLOR is
// unset or empty (https://no-color.org).
func NewRenderer(terminal bool, noColor string) Renderer {
	return Renderer{Color: terminal && noColor == ""}
}

// IsTerminal reports whether w is a terminal.
func IsTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok || file == nil {
		return false
	}
	return isatty.IsTerminal(file.Fd()) || isatty.IsCygwinTerminal(file.Fd())
}
