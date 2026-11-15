package doneeval

import "strings"

// hasRollback recognizes an ATX level-two heading, with Markdown's optional
// closing hashes, outside backtick/tilde fences. Four-space indented code and
// blockquoted examples are not headings. Fence length and character must match.
func hasRollback(body string) bool {
	var fence byte
	length := 0
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSuffix(line, "\r")
		spaces := len(line) - len(strings.TrimLeft(line, " "))
		if spaces > 3 {
			continue
		}
		line = line[spaces:]
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			char := line[0]
			n := len(line) - len(strings.TrimLeft(line, string(char)))
			if fence == 0 {
				if char == '`' && strings.Contains(line[n:], "`") {
					continue
				}
				fence, length = char, n
			} else if char == fence && n >= length && strings.TrimSpace(line[n:]) == "" {
				fence = 0
			}
			continue
		}
		if fence != 0 || !strings.HasPrefix(line, "##") {
			continue
		}
		rest := strings.TrimPrefix(line, "##")
		if rest == "" || rest[0] != ' ' && rest[0] != '\t' {
			continue
		}
		rest = strings.TrimSpace(rest)
		if end := strings.TrimRight(rest, "#"); end != rest && len(end) > 0 && (end[len(end)-1] == ' ' || end[len(end)-1] == '\t') {
			rest = strings.TrimSpace(end)
		}
		if rest == "Rollback" {
			return true
		}
	}
	return false
}
