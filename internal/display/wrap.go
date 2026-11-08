package display

import (
	"fmt"
	"strings"
)

// Wrap breaks text at spaces into lines, the first at most first columns
// wide and the rest at most rest. Text that fits is returned as it is. A
// word wider than its line stands whole on a line of its own: fid never
// breaks inside a word.
func Wrap(text string, first, rest int) []string {
	if VisibleWidth(text) <= first {
		return []string{text}
	}
	return fill(strings.Split(text, " "), first, rest, "")
}

// WrapCommand breaks a shell command between its arguments, never between a
// flag and its value, into lines the first at most first columns wide and
// the rest at most rest. Each broken line ends in " \", so the lines run as
// one command; copying only the first leaves the shell waiting.
func WrapCommand(command string, first, rest int) []string {
	if VisibleWidth(command) <= first {
		return []string{command}
	}
	arguments := []string{}
	for _, word := range strings.Fields(command) {
		if last := len(arguments) - 1; last >= 0 && isFlag(arguments[last]) && !isFlag(word) {
			arguments[last] += " " + word
			continue
		}
		arguments = append(arguments, word)
	}
	return fill(arguments, first, rest, " \\")
}

// isFlag says whether an argument is a flag still waiting for its value.
func isFlag(argument string) bool {
	return strings.HasPrefix(argument, "-") && !strings.Contains(argument, " ")
}

// fill packs words into lines, joined by spaces, ending each broken line in
// suffix, which must fit too; the last line needs none.
func fill(words []string, first, rest int, suffix string) []string {
	lines, line, width := []string{}, "", first
	for index, word := range words {
		room := width - VisibleWidth(suffix)
		if index == len(words)-1 {
			room = width
		}
		switch {
		case line == "":
			line = word
		case VisibleWidth(line)+1+VisibleWidth(word) <= room:
			line += " " + word
		default:
			lines = append(lines, strings.TrimRight(line, " ")+suffix)
			line, width = word, rest
		}
	}
	return append(lines, line)
}

// StateText is how a workflow state reads: underscores as spaces, and todo
// as "to do".
func StateText(state string) string {
	if state == "todo" {
		return "to do"
	}
	return strings.ReplaceAll(state, "_", " ")
}

// BlockedText names an unstarted issue's first open blocker and how many
// more there are, or is empty when nothing blocks it.
func BlockedText(blockers []string) string {
	switch len(blockers) {
	case 0:
		return ""
	case 1:
		return "blocked by " + blockers[0]
	}
	return fmt.Sprintf("blocked by %s +%d", blockers[0], len(blockers)-1)
}
