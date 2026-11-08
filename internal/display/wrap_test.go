package display

import (
	"slices"
	"strings"
	"testing"
)

func TestWrapBreaksProseAtSpaces(t *testing.T) {
	t.Parallel()
	path := "/home/mark/.config/fid/config.toml"
	for _, test := range []struct {
		name        string
		text        string
		first, rest int
		want        []string
	}{
		{"fits", "Backed up 3 hours ago", 66, 66, []string{"Backed up 3 hours ago"}},
		{"fits with its spaces kept", "two  spaces", 11, 11, []string{"two  spaces"}},
		{"breaks at the last space that fits", "Can't read user config " + path + "; last backed up 1 hour ago", 66, 66,
			[]string{"Can't read user config " + path + "; last", "backed up 1 hour ago"}},
		{"first line narrower", "one two three four", 7, 10, []string{"one two", "three four"}},
		{"a word wider than the line stands whole", "fix " + path + " now", 10, 10, []string{"fix", path, "now"}},
		{"a run of spaces at a break is dropped", "aaaa    bbbb", 6, 6, []string{"aaaa", "bbbb"}},
		{"empty", "", 10, 10, []string{""}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := Wrap(test.text, test.first, test.rest); !slices.Equal(got, test.want) {
				t.Fatalf("Wrap(%q, %d, %d) = %q, want %q", test.text, test.first, test.rest, got, test.want)
			}
		})
	}
}

// A command breaks between arguments, never between a flag and its value,
// and each broken line ends in a shell continuation that must fit too.
func TestWrapCommandContinuesTheCommand(t *testing.T) {
	t.Parallel()
	create := "fid issue create --kind task --title <title> --actor human:<name> --idempotency-key <key>"
	for _, test := range []struct {
		name        string
		command     string
		first, rest int
		want        []string
	}{
		{"fits", "fid backup", 62, 60, []string{"fid backup"}},
		{"keeps each flag with its value", create, 62, 60,
			[]string{"fid issue create --kind task --title <title> \\", "--actor human:<name> --idempotency-key <key>"}},
		{"the continuation must fit", "fid aaaa bbbb", 9, 9, []string{"fid \\", "aaaa bbbb"}},
		{"the last line needs no continuation", "fid aaaa bbbb cccc", 10, 9, []string{"fid aaaa \\", "bbbb cccc"}},
		{"a flag followed by a flag stands alone", "fid status --all --json", 12, 12, []string{"fid status \\", "--all --json"}},
		{"an argument wider than the line stands whole", "fid show " + strings.Repeat("x", 12), 10, 10,
			[]string{"fid show \\", strings.Repeat("x", 12)}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := WrapCommand(test.command, test.first, test.rest); !slices.Equal(got, test.want) {
				t.Fatalf("WrapCommand(%q, %d, %d) = %q, want %q", test.command, test.first, test.rest, got, test.want)
			}
		})
	}
}

func TestStateTextReadsAWorkflowState(t *testing.T) {
	t.Parallel()
	for state, want := range map[string]string{
		"todo": "to do", "in_progress": "in progress", "qa_review": "qa review", "done": "done", "closed": "closed",
	} {
		if got := StateText(state); got != want {
			t.Errorf("StateText(%q) = %q, want %q", state, got, want)
		}
	}
}
