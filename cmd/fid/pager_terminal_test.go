//go:build linux

// Opening a pseudo-terminal takes Linux's ioctls; CI runs on Linux.

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/display"
	"golang.org/x/sys/unix"
)

// fidProcess is what fid, run as a process of its own, wrote to stdout,
// with a terminal's "\r\n" read as "\n", and, when its pager ran, what the
// pager read and the pager's environment.
type fidProcess struct {
	stdout, paged, pagerEnviron string
	pagerRan                    bool
}

// On a real terminal, fid measures the screen itself: status and show page
// text taller than the terminal, and print text that fits, as they do
// through a pipe. The commands run as fid, with fid's own environment.
func TestStatusAndShowPageOnARealTerminal(t *testing.T) {
	t.Parallel()
	fixture, record := newPagerFixture(t)
	// fid reads the real clock, so the backup's age is fixed from it.
	fixture.writeRecordAt(time.Now(), statusRecordLine{"pushed", 3*time.Hour + 30*time.Minute, ""})
	colour := regexp.MustCompile("\033\\[[0-9;]*m")
	for _, args := range [][]string{{"status"}, {"show", "FID-e3a1c0"}} {
		piped := runFIDProcess(t, fixture, record, nil, args...)
		if piped.pagerRan || piped.stdout == "" || strings.Contains(piped.stdout, "\033") {
			t.Fatalf("%v piped printed %q, pager ran %v; want plain text printed", args, piped.stdout, piped.pagerRan)
		}
		short, tall := display.Screen{Width: 80, Height: 8}, display.Screen{Width: 80, Height: 60}
		paged := runFIDProcess(t, fixture, record, &short, args...)
		environ := strings.Split(paged.pagerEnviron, "\n")
		if !paged.pagerRan || paged.stdout != "" || colour.ReplaceAllString(paged.paged, "") != piped.stdout ||
			!colour.MatchString(paged.paged) || !slices.Contains(environ, "LESS=FRX") {
			t.Fatalf("%v on an 80x8 terminal printed %q and paged %q (pager ran %v) with environment %q; want the text in colour paged with LESS=FRX, as piped:\n%s",
				args, paged.stdout, paged.paged, paged.pagerRan, environ, piped.stdout)
		}
		fits := runFIDProcess(t, fixture, record, &tall, args...)
		if fits.pagerRan || colour.ReplaceAllString(fits.stdout, "") != piped.stdout || !colour.MatchString(fits.stdout) {
			t.Fatalf("%v on an 80x60 terminal printed %q, pager ran %v; want the text in colour printed, as piped:\n%s",
				args, fits.stdout, fits.pagerRan, piped.stdout)
		}
	}
}

// runFIDProcess runs fid in fixture's project as a process of its own,
// stdout a terminal of screen's size or, when screen is nil, a pipe. fid's
// environment is the test's, with the fixture's HOME, user config, PATH and
// PAGER, and without LESS or NO_COLOR.
func runFIDProcess(t *testing.T, fixture *statusFixture, record string, screen *display.Screen, args ...string) fidProcess {
	t.Helper()
	for _, path := range []string{record, record + ".env"} {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, executable, args...)
	command.Dir = fixture.workspace
	command.Env = []string{runAsFIDVariable + "=1"}
	for _, entry := range os.Environ() {
		if key, _, _ := strings.Cut(entry, "="); key != "LESS" && key != "NO_COLOR" {
			command.Env = append(command.Env, entry)
		}
	}
	// Later entries win: the fixture's PATH and PAGER among them.
	command.Env = append(command.Env, append([]string{"HOME=" + fixture.home, "XDG_CONFIG_HOME=" + fixture.configHome}, fixture.environ...)...)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if screen == nil {
		err = command.Run()
	} else {
		leader, follower := openTerminal(t, *screen)
		command.Stdout = follower
		reached := make(chan []byte, 1)
		go func() {
			// Reading fails, with EIO, once fid and its pager have closed
			// the terminal and everything they wrote has been read.
			read, _ := io.ReadAll(leader)
			reached <- read
		}()
		err = command.Start()
		// fid holds the terminal now; closing ours lets the read end.
		_ = follower.Close()
		if err == nil {
			err = command.Wait()
		}
		stdout.Write(bytes.ReplaceAll(<-reached, []byte("\r\n"), []byte("\n")))
	}
	if err != nil || stderr.Len() != 0 {
		t.Fatalf("fid %v = %v, stderr %q", args, err, stderr.String())
	}
	result := fidProcess{stdout: stdout.String()}
	if paged, err := os.ReadFile(record); err == nil {
		environ, err := os.ReadFile(record + ".env")
		if err != nil {
			t.Fatal(err)
		}
		result.paged, result.pagerEnviron, result.pagerRan = string(paged), string(environ), true
	}
	return result
}

// openTerminal opens a pseudo-terminal of screen's size: a process writes
// to follower, and what reaches the screen is read from leader.
func openTerminal(t *testing.T, screen display.Screen) (leader, follower *os.File) {
	t.Helper()
	leader, err := os.OpenFile("/dev/ptmx", os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		t.Fatalf("open a pseudo-terminal: %v", err)
	}
	t.Cleanup(func() { _ = leader.Close() })
	if err := unix.IoctlSetPointerInt(int(leader.Fd()), unix.TIOCSPTLCK, 0); err != nil {
		t.Fatalf("unlock the pseudo-terminal: %v", err)
	}
	number, err := unix.IoctlGetUint32(int(leader.Fd()), unix.TIOCGPTN)
	if err != nil {
		t.Fatalf("number the pseudo-terminal: %v", err)
	}
	follower, err = os.OpenFile(fmt.Sprintf("/dev/pts/%d", number), os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		t.Fatalf("open the pseudo-terminal's follower: %v", err)
	}
	t.Cleanup(func() { _ = follower.Close() })
	size := unix.Winsize{Row: uint16(screen.Height), Col: uint16(screen.Width)}
	if err := unix.IoctlSetWinsize(int(follower.Fd()), unix.TIOCSWINSZ, &size); err != nil {
		t.Fatalf("size the pseudo-terminal: %v", err)
	}
	return leader, follower
}
