package display

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

// Screen is a terminal's size, in columns and rows.
type Screen struct {
	Width, Height int
}

// TerminalScreen is the size of the terminal w writes to, or nil when w
// isn't a terminal or its size is unknown.
func TerminalScreen(w io.Writer) *Screen {
	file, ok := w.(*os.File)
	if !ok || file == nil || !IsTerminal(w) {
		return nil
	}
	size, err := unix.IoctlGetWinsize(int(file.Fd()), unix.TIOCGWINSZ)
	if err != nil || size.Row == 0 || size.Col == 0 {
		return nil
	}
	return &Screen{Width: int(size.Col), Height: int(size.Row)}
}

// rows counts the rows text takes on the screen: a line wider than the
// screen wraps onto more than one.
func (screen Screen) rows(text string) int {
	rows := 0
	for _, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		rows += max(1, (VisibleWidth(line)+screen.Width-1)/screen.Width)
	}
	return rows
}

// shellCharacters are those that make git run a pager through the shell.
const shellCharacters = "|&;<>()$`\\\"' \t\n*?[#~=%"

// Page writes text to stdout, through a pager when stdout is a terminal,
// its screen, and text is taller than it, as git does. The pager is $PAGER
// when it's set and less otherwise, and there is none when that is empty
// or cat. A pager with arguments or shell syntax runs through sh -c. Its
// environment is environ, with LESS=FRX when LESS is unset. When the pager
// can't start, text is printed directly. Quitting the pager early isn't an
// error, and fid waits for it to exit.
func Page(stdout, stderr io.Writer, screen *Screen, environ []string, text string) error {
	pager, ok := pagerCommand(environ)
	if screen == nil || !ok || screen.rows(text) <= screen.Height {
		return write(stdout, text)
	}
	pager.Stdin, pager.Stdout, pager.Stderr, pager.Env = strings.NewReader(text), stdout, stderr, environ
	if _, set := lookupEnv(environ, "LESS"); !set {
		pager.Env = append(append([]string{}, environ...), "LESS=FRX")
	}
	// An interrupt typed in the pager reaches fid too: fid waits for the
	// pager rather than leaving it without the terminal.
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	defer signal.Stop(interrupts)
	if err := pager.Start(); err != nil {
		return write(stdout, text)
	}
	var exit *exec.ExitError
	if err := pager.Wait(); err != nil && !errors.As(err, &exit) {
		return fmt.Errorf("page: %w", err)
	}
	return nil
}

// pagerCommand is the pager to run, or false when there is none.
func pagerCommand(environ []string) (*exec.Cmd, bool) {
	command, set := lookupEnv(environ, "PAGER")
	if !set {
		command = "less"
	}
	if command == "" || command == "cat" {
		return nil, false
	}
	if strings.ContainsAny(command, shellCharacters) {
		return exec.Command("/bin/sh", "-c", command), true
	}
	path, ok := lookPath(command, environ)
	return exec.Command(path), ok
}

// lookPath finds command in the directories environ's PATH names, as the
// shell would. A command naming a directory is used as it is.
func lookPath(command string, environ []string) (string, bool) {
	if strings.Contains(command, "/") {
		return command, true
	}
	path, _ := lookupEnv(environ, "PATH")
	for _, directory := range filepath.SplitList(path) {
		candidate := filepath.Join(directory, command)
		if info, err := os.Stat(candidate); err == nil && info.Mode().IsRegular() && info.Mode()&0o111 != 0 {
			return candidate, true
		}
	}
	return "", false
}

// lookupEnv is name's last value in environ, and whether it is set.
func lookupEnv(environ []string, name string) (string, bool) {
	value, set := "", false
	for _, entry := range environ {
		if key, found, ok := strings.Cut(entry, "="); ok && key == name {
			value, set = found, true
		}
	}
	return value, set
}

func write(w io.Writer, text string) error {
	if _, err := io.WriteString(w, text); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
