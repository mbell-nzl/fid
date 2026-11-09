package main

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/display"
	"github.com/spf13/cobra"
)

// displayEnvironmentKey carries a displayEnvironment in the root command's
// context, so tests fix the time, its zone, HOME, the XDG directories, the
// pager and the terminal.
type displayEnvironmentKey struct{}

// displayEnvironment is what status and show read beyond the ledger: the
// environment, the clock, and the size of the terminal stdout is, if any.
type displayEnvironment struct {
	environ []string
	now     func() time.Time
	screen  func(io.Writer) *display.Screen
}

func commandDisplayEnvironment(cmd *cobra.Command) displayEnvironment {
	if environment, ok := cmd.Root().Context().Value(displayEnvironmentKey{}).(displayEnvironment); ok {
		return environment
	}
	return displayEnvironment{environ: os.Environ(), now: time.Now, screen: display.TerminalScreen}
}

// writeHuman renders human text, coloured only when stdout is a terminal,
// and writes it through a pager when that terminal is shorter than the
// text.
func writeHuman(cmd *cobra.Command, environment displayEnvironment, write func(io.Writer, display.Renderer) error) error {
	renderer := display.NewRenderer(display.IsTerminal(cmd.OutOrStdout()), environmentValue(environment.environ, "NO_COLOR"))
	var text strings.Builder
	if err := write(&text, renderer); err != nil {
		return err
	}
	return display.Page(cmd.OutOrStdout(), cmd.ErrOrStderr(), environment.screen(cmd.OutOrStdout()), environment.environ, text.String())
}

// environmentValue is name's value in environ, or "" when it is unset.
func environmentValue(environ []string, name string) string {
	value := ""
	for _, entry := range environ {
		if key, found, ok := strings.Cut(entry, "="); ok && key == name {
			value = found
		}
	}
	return value
}
