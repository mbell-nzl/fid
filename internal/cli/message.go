package cli

import (
	"fmt"
	"strings"
	"unicode"
)

// EscapeControlText makes caller text safe on one Human output line. Existing
// Go-quoted escapes are literal text and pass through unchanged.
func EscapeControlText(text string) string {
	var result strings.Builder
	for _, value := range text {
		switch value {
		case '\n':
			result.WriteString(`\n`)
		case '\r':
			result.WriteString(`\r`)
		case '\t':
			result.WriteString(`\t`)
		default:
			if unicode.IsControl(value) {
				_, _ = fmt.Fprintf(&result, `\u%04x`, value)
			} else {
				result.WriteRune(value)
			}
		}
	}
	return result.String()
}

// BoundedMessage implements DEC-282: escape first, then retain at most 512
// Unicode scalars, dropping a whole escape if it crosses the 511-scalar cut.
func BoundedMessage(text string) string {
	values := []rune(EscapeControlText(text))
	if len(values) <= 512 {
		return string(values)
	}
	return string(values[:escapedPrefix(values, 511)]) + "…"
}

// CutEscapedText fits escaped text to width columns. Text that fits is
// returned as it is; longer text is cut to end in an ellipsis, dropping
// whole an escape that crosses the cut, as BoundedMessage does.
func CutEscapedText(text string, width int) string {
	values := []rune(text)
	switch {
	case len(values) <= width:
		return text
	case width <= 0:
		return ""
	}
	return string(values[:escapedPrefix(values, width-1)]) + "…"
}

// escapedPrefix is how many of values' leading runes fit in limit without
// cutting an escape. values must be longer than limit.
func escapedPrefix(values []rune, limit int) int {
	cut := 0
	for cut < limit {
		width := messageEscapeWidth(values[cut:])
		if cut+width > limit {
			break
		}
		cut += width
	}
	return cut
}

func messageEscapeWidth(values []rune) int {
	if len(values) < 2 || values[0] != '\\' {
		return 1
	}
	width := 1
	switch values[1] {
	case '\\', '"', '\'', 'a', 'b', 'f', 'n', 'r', 't', 'v':
		width = 2
	case 'x':
		width = 4
	case 'u':
		width = 6
	case 'U':
		width = 10
	case '0', '1', '2', '3', '4', '5', '6', '7':
		width = 4
	}
	if width > len(values) {
		return 1
	}
	if width > 2 {
		for _, value := range values[2:width] {
			if values[1] >= '0' && values[1] <= '7' {
				if value < '0' || value > '7' {
					return 1
				}
			} else if !strings.ContainsRune("0123456789abcdefABCDEF", value) {
				return 1
			}
		}
	}
	return width
}
