package cli

import (
	"strings"
	"testing"
)

func TestBoundedMessage(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ name, input, want string }{
		{"controls", "a\n\r\t\x1b\x7f\u0085b", `a\n\r\t\u001b\u007f\u0085b`},
		{"512", strings.Repeat("a", 512), strings.Repeat("a", 512)},
		{"513", strings.Repeat("a", 513), strings.Repeat("a", 511) + "…"},
		{"scalars", strings.Repeat("界", 513), strings.Repeat("界", 511) + "…"},
		{"escaped 512", strings.Repeat("a", 506) + "\x1b", strings.Repeat("a", 506) + `\u001b`},
		{"escaped 513", "\x1b" + strings.Repeat("a", 507), `\u001b` + strings.Repeat("a", 505) + "…"},
		{"escape cut", strings.Repeat("a", 509) + "\x1bmore", strings.Repeat("a", 509) + "…"},
		{"already escaped", `\x1b\n\t\u0085\U0001f600\\`, `\x1b\n\t\u0085\U0001f600\\`},
		{"quoted cut", strings.Repeat("a", 509) + `\x1bmore`, strings.Repeat("a", 509) + "…"},
		{"backslash pair cut", strings.Repeat("a", 510) + `\\tail`, strings.Repeat("a", 510) + "…"},
		{"Unicode escape cut", strings.Repeat("a", 507) + `\U0001f600tail`, strings.Repeat("a", 507) + "…"},
		{"octal escape cut", strings.Repeat("a", 510) + `\033tail`, strings.Repeat("a", 510) + "…"},
		{"literal slash", strings.Repeat("a", 509) + `\users!`, strings.Repeat("a", 509) + `\u` + "…"},
	} {
		t.Run(row.name, func(t *testing.T) {
			if got := BoundedMessage(row.input); got != row.want {
				t.Fatalf("message = %q, want %q", got, row.want)
			}
		})
	}
}

// Escaped text cut to a column drops an escape that crosses the cut whole,
// as BoundedMessage does, so no escape is left half shown.
func TestCutEscapedText(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, input string
		width       int
		want        string
	}{
		{"fits", "short", 10, "short"},
		{"exactly fits", "abcde", 5, "abcde"},
		{"cut", "abcdefgh", 5, "abcd…"},
		{"scalars", "界界界界界界", 4, "界界界…"},
		{"escape crosses the cut", `ab\u001bcd`, 5, "ab…"},
		{"escape before the cut", `\nabcdef`, 5, `\nab…`},
		{"escape ends at the cut", `ab\u001bcd`, 9, `ab\u001b…`},
		{"one column", "abc", 1, "…"},
		{"no columns", "abc", 0, ""},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			if got := CutEscapedText(row.input, row.width); got != row.want {
				t.Fatalf("CutEscapedText(%q, %d) = %q, want %q", row.input, row.width, got, row.want)
			}
		})
	}
}
