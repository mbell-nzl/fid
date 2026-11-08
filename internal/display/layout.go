package display

import (
	"fmt"
	"strings"
)

// TopBorder is a rounded box's top edge, width columns wide, with a label at
// each end. Labels may be coloured: their escape sequences take no width.
// Labels too wide to fit are left out.
func TopBorder(left, right string, width int) string {
	inside := width - 2
	if inside <= 0 {
		return "╭╮"
	}
	left, right = borderLabel(left), borderLabel(right)
	fill := inside - VisibleWidth(left) - VisibleWidth(right)
	if fill < 0 {
		return "╭" + strings.Repeat("─", inside) + "╮"
	}
	return "╭" + left + strings.Repeat("─", fill) + right + "╮"
}

// BottomBorder is a rounded box's bottom edge, width columns wide.
func BottomBorder(width int) string {
	if width <= 2 {
		return "╰╯"
	}
	return "╰" + strings.Repeat("─", width-2) + "╯"
}

func borderLabel(label string) string {
	if label == "" {
		return ""
	}
	return " " + label + " "
}

// ProgressBar is a bar width columns wide, filled in proportion to
// value/total, followed by the count.
func ProgressBar(value, total, width int) string {
	width = max(width, 0)
	if total <= 0 {
		return strings.Repeat("░", width) + "  0/0"
	}
	value = min(max(value, 0), total)
	filled := value * width / total
	return fmt.Sprintf("%s%s  %d/%d", strings.Repeat("█", filled), strings.Repeat("░", width-filled), value, total)
}

// PadRight fits plain text to width columns, padding it on the right, or
// cutting it to end in an ellipsis.
func PadRight(text string, width int) string {
	text = clip(text, width)
	return text + strings.Repeat(" ", max(width-VisibleWidth(text), 0))
}

// PadLeft fits plain text to width columns, padding it on the left, or
// cutting it to end in an ellipsis.
func PadLeft(text string, width int) string {
	text = clip(text, width)
	return strings.Repeat(" ", max(width-VisibleWidth(text), 0)) + text
}

func clip(text string, width int) string {
	runes := []rune(text)
	switch {
	case width <= 0:
		return ""
	case len(runes) <= width:
		return text
	}
	return string(runes[:width-1]) + Ellipsis
}

// VisibleWidth counts the columns text takes on a terminal: one per rune,
// none for ANSI escape sequences.
func VisibleWidth(text string) int {
	width := 0
	state := ansiText
	pending := 0
	for _, r := range text {
		switch state {
		case ansiEscape:
			pending++
			if r == '[' {
				state = ansiCSI
				continue
			}
			width += pending
			state, pending = ansiText, 0
		case ansiCSI:
			pending++
			if r >= 0x40 && r <= 0x7e {
				state, pending = ansiText, 0
			}
		default:
			if r == '\x1b' {
				state, pending = ansiEscape, 1
				continue
			}
			width++
		}
	}
	if state != ansiText {
		width += pending
	}
	return width
}

const (
	ansiText = iota
	ansiEscape
	ansiCSI
)
