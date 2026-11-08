package display

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// fakePagerDir holds the fake pager, as "pager" and as "less", written
// before any test runs so no test's fork holds it open for writing while
// another runs it.
var fakePagerDir string

// fakePager records its arguments, environment and stdin in the directory
// FAKE_PAGER_RECORD names, or quits at once, reading nothing, when
// FAKE_PAGER_QUIT is set. With FAKE_PAGER_INTERRUPT set it first
// interrupts its parent, as an interrupt typed in a pager reaches fid too.
const fakePager = `#!/bin/sh
if [ -n "$FAKE_PAGER_QUIT" ]; then exit 3; fi
if [ -n "$FAKE_PAGER_INTERRUPT" ]; then kill -INT "$PPID"; fi
printf '%s\n' "$@" > "$FAKE_PAGER_RECORD/args"
env > "$FAKE_PAGER_RECORD/env"
cat > "$FAKE_PAGER_RECORD/stdin"
`

func TestMain(m *testing.M) {
	code := func() int {
		directory, err := os.MkdirTemp("", "fid-pager-")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		defer func() { _ = os.RemoveAll(directory) }()
		for _, name := range []string{"pager", "less"} {
			if err := os.WriteFile(filepath.Join(directory, name), []byte(fakePager), 0o755); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
		}
		fakePagerDir = directory
		return m.Run()
	}()
	os.Exit(code)
}

// lines is count numbered lines, each width columns wide.
func lines(count, width int) string {
	var text strings.Builder
	for index := range count {
		line := fmt.Sprintf("line %d ", index+1)
		text.WriteString(line + strings.Repeat("x", max(width-len(line), 0)) + "\n")
	}
	return text.String()
}

// pageResult is what a Page call wrote and what the fake pager recorded.
type pageResult struct {
	err            error
	stdout, stderr string
	paged          bool
	args, environ  []string
	stdin          string
}

func page(t *testing.T, text string, screen *Screen, environ ...string) pageResult {
	t.Helper()
	record := t.TempDir()
	environ = append([]string{"PATH=" + os.Getenv("PATH"), "FAKE_PAGER_RECORD=" + record}, environ...)
	var stdout, stderr bytes.Buffer
	result := pageResult{err: Page(&stdout, &stderr, screen, environ, text)}
	result.stdout, result.stderr = stdout.String(), stderr.String()
	stdin, err := os.ReadFile(filepath.Join(record, "stdin"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	result.paged, result.stdin = err == nil, string(stdin)
	for name, list := range map[string]*[]string{"args": &result.args, "env": &result.environ} {
		payload, err := os.ReadFile(filepath.Join(record, name))
		if err == nil {
			*list = strings.Split(strings.TrimSuffix(string(payload), "\n"), "\n")
		}
	}
	return result
}

func TestPagePagesOnlyTextTallerThanTheScreen(t *testing.T) {
	t.Parallel()
	pager := "PAGER=" + filepath.Join(fakePagerDir, "pager")
	screen := &Screen{Width: 40, Height: 10}
	for _, test := range []struct {
		name   string
		text   string
		screen *Screen
		paged  bool
	}{
		{"fits", lines(10, 40), screen, false},
		{"taller", lines(11, 40), screen, true},
		{"lines the terminal wraps count as more than one row", lines(6, 41), screen, true},
		{"colour codes take no width", strings.Repeat("\033[36m"+strings.Repeat("x", 40)+"\033[0m\n", 10), screen, false},
		{"not a terminal", lines(50, 40), nil, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := page(t, test.text, test.screen, pager)
			if result.err != nil || result.paged != test.paged {
				t.Fatalf("Page() = %v, paged %v; want paged %v", result.err, result.paged, test.paged)
			}
			if test.paged && (result.stdin != test.text || result.stdout != "") || !test.paged && result.stdout != test.text {
				t.Fatalf("pager read %q and stdout has %q, want the text once", result.stdin, result.stdout)
			}
		})
	}
}

// The pager is $PAGER when it's set, less otherwise, and none when PAGER is
// empty or cat. less gets LESS=FRX, unless the user has a LESS of their own.
func TestPageChoosesThePagerAsGitDoes(t *testing.T) {
	t.Parallel()
	script := filepath.Join(fakePagerDir, "pager")
	tall := lines(20, 10)
	for _, test := range []struct {
		name    string
		environ []string
		paged   bool
		args    []string
		less    string
	}{
		{"PAGER", []string{"PAGER=" + script}, true, []string{""}, "LESS=FRX"},
		{"PAGER with arguments, through sh", []string{"PAGER=" + script + " -S --shift 4"}, true, []string{"-S", "--shift", "4"}, "LESS=FRX"},
		{"less when PAGER is unset", []string{"PATH=" + fakePagerDir + ":" + os.Getenv("PATH")}, true, []string{""}, "LESS=FRX"},
		{"the user's own LESS", []string{"PAGER=" + script, "LESS=-S"}, true, []string{""}, "LESS=-S"},
		{"the user's own empty LESS", []string{"PAGER=" + script, "LESS="}, true, []string{""}, "LESS="},
		{"PAGER=cat", []string{"PAGER=cat"}, false, nil, ""},
		{"empty PAGER", []string{"PAGER="}, false, nil, ""},
		{"PAGER that can't start", []string{"PAGER=" + filepath.Join(t.TempDir(), "missing")}, false, nil, ""},
		{"less missing", []string{"PATH=" + t.TempDir()}, false, nil, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := page(t, tall, &Screen{Width: 80, Height: 10}, test.environ...)
			if result.err != nil || result.paged != test.paged {
				t.Fatalf("Page() = %v, paged %v; want paged %v (stderr %q)", result.err, result.paged, test.paged, result.stderr)
			}
			if !test.paged {
				if result.stdout != tall {
					t.Fatalf("stdout = %q, want the text printed directly", result.stdout)
				}
				return
			}
			if result.stdin != tall || !slices.Equal(result.args, test.args) || !slices.Contains(result.environ, test.less) {
				t.Fatalf("pager read %q with arguments %q and environment %q; want %q and %s", result.stdin, result.args, result.environ, test.args, test.less)
			}
		})
	}
}

// Quitting the pager before it reads everything isn't an error.
func TestPageAcceptsAPagerThatQuitsEarly(t *testing.T) {
	t.Parallel()
	text := lines(20000, 79)
	result := page(t, text, &Screen{Width: 80, Height: 24}, "PAGER="+filepath.Join(fakePagerDir, "pager"), "FAKE_PAGER_QUIT=1")
	if result.err != nil || result.stdout != "" {
		t.Fatalf("Page() = %v, stdout %d bytes; want no error and nothing printed", result.err, len(result.stdout))
	}
}

// An interrupt typed in the pager doesn't stop fid, which waits for the
// pager. Were it not caught, it would end this test's process.
func TestPageWaitsForThePagerThroughAnInterrupt(t *testing.T) {
	t.Parallel()
	text := lines(30, 10)
	result := page(t, text, &Screen{Width: 80, Height: 24}, "PAGER="+filepath.Join(fakePagerDir, "pager"), "FAKE_PAGER_INTERRUPT=1")
	if result.err != nil || !result.paged || result.stdin != text {
		t.Fatalf("Page() = %v, paged %v; want the text paged", result.err, result.paged)
	}
}
