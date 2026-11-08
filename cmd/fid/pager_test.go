package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/display"
)

// newPagerFixture is showWork's project whose PAGER records what it reads
// in record, and its environment in record.env. The pager runs through sh,
// as $PAGER with arguments does, so no test runs a file just written.
func newPagerFixture(t *testing.T) (fixture *statusFixture, record string) {
	t.Helper()
	fixture = newStatusFixture(t, 100*time.Hour)
	fixture.setTarget()
	fixture.writeRecord(statusRecordLine{"pushed", 3 * time.Hour, ""})
	showWork(fixture)
	directory := t.TempDir()
	script, record := filepath.Join(directory, "pager.sh"), filepath.Join(directory, "paged")
	if err := os.WriteFile(script, []byte("env > \"$1.env\"\ncat > \"$1\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fixture.environ = []string{"PATH=" + os.Getenv("PATH"), "PAGER=/bin/sh " + script + " " + record}
	return fixture, record
}

// On a terminal, status and show page text taller than the screen, as it
// would be printed; JSON, text that fits, piped output and other commands
// never page.
func TestStatusAndShowPageLongTextOnATerminal(t *testing.T) {
	t.Parallel()
	short, tall := &display.Screen{Width: 80, Height: 8}, &display.Screen{Width: 80, Height: 60}
	for _, test := range []struct {
		name   string
		args   []string
		screen *display.Screen
		paged  bool
	}{
		{"status taller than the screen", []string{"status"}, short, true},
		{"show taller than the screen", []string{"show", "FID-e3a1c0"}, short, true},
		{"show that fits", []string{"show", "FID-e3a1c0"}, tall, false},
		{"show piped", []string{"show", "FID-e3a1c0"}, nil, false},
		{"status JSON", []string{"--json", "status"}, short, false},
		{"show JSON", []string{"--json", "show", "FID-e3a1c0"}, short, false},
		{"another command", []string{"issue", "list"}, short, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture, record := newPagerFixture(t)
			printed := fixture.run(test.args...)
			fixture.screen = test.screen
			result := fixture.run(test.args...)
			paged, err := os.ReadFile(record)
			if result.err != nil || result.stderr != "" || printed.err != nil || (err == nil) != test.paged {
				t.Fatalf("%v = %v, stderr %q; paged %v, want %v", test.args, result.err, result.stderr, err == nil, test.paged)
			}
			if !test.paged {
				if result.stdout != printed.stdout {
					t.Fatalf("%v printed %q, want %q", test.args, result.stdout, printed.stdout)
				}
				return
			}
			environ, err := os.ReadFile(record + ".env")
			if err != nil || result.stdout != "" || string(paged) != printed.stdout || !slices.Contains(strings.Split(string(environ), "\n"), "LESS=FRX") {
				t.Fatalf("%v printed %q and paged %q with environment %q (%v), want %q paged with LESS=FRX", test.args, result.stdout, paged, environ, err, printed.stdout)
			}
		})
	}
}
