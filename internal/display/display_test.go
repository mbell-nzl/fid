package display

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewRendererColorsOnlyATerminalWithoutNoColor(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name     string
		terminal bool
		noColor  string
		want     bool
	}{
		{"terminal, NO_COLOR unset", true, "", true},
		{"terminal, NO_COLOR set", true, "1", false},
		{"terminal, NO_COLOR any non-empty value", true, "false", false},
		{"not a terminal", false, "", false},
		{"not a terminal, NO_COLOR set", false, "1", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := NewRenderer(test.terminal, test.noColor).Color; got != test.want {
				t.Fatalf("NewRenderer(%v, %q).Color = %v, want %v", test.terminal, test.noColor, got, test.want)
			}
		})
	}
}

func TestIsTerminalRefusesBuffersAndFiles(t *testing.T) {
	t.Parallel()
	file, err := os.Create(filepath.Join(t.TempDir(), "output"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	for _, test := range []struct {
		name string
		got  bool
	}{
		{"buffer", IsTerminal(&bytes.Buffer{})},
		{"regular file", IsTerminal(file)},
		{"nil", IsTerminal(nil)},
	} {
		if test.got {
			t.Errorf("IsTerminal(%s) = true, want false", test.name)
		}
	}
}

func TestRendererColorize(t *testing.T) {
	t.Parallel()
	on, off := Renderer{Color: true}, Renderer{}
	for _, test := range []struct {
		name     string
		renderer Renderer
		role     Role
		want     string
	}{
		{"header", on, RoleHeader, "\033[1mtext\033[0m"},
		{"id", on, RoleID, "\033[36mtext\033[0m"},
		{"error", on, RoleError, "\033[31mtext\033[0m"},
		{"warning", on, RoleWarning, "\033[33mtext\033[0m"},
		{"success", on, RoleSuccess, "\033[32mtext\033[0m"},
		{"muted", on, RoleMuted, "\033[2mtext\033[0m"},
		{"label", on, RoleLabel, "\033[35mtext\033[0m"},
		{"normal", on, RoleNormal, "text"},
		{"unknown role", on, Role("unknown"), "text"},
		{"color off", off, RoleError, "text"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := test.renderer.Colorize(test.role, "text"); got != test.want {
				t.Fatalf("Colorize(%q) = %q, want %q", test.role, got, test.want)
			}
		})
	}
}

func TestRendererLine(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name     string
		renderer Renderer
		role     Role
		want     string
	}{
		{"colored", Renderer{Color: true}, RoleWarning, "\033[33m⚠\033[0m \033[33mlate\033[0m"},
		{"normal role leaves the text alone", Renderer{Color: true}, RoleNormal, "⚠ late"},
		{"color off", Renderer{}, RoleWarning, "⚠ late"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := test.renderer.Line(test.role, SymbolWarning, "late"); got != test.want {
				t.Fatalf("Line() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestSymbolsAreUnicode(t *testing.T) {
	t.Parallel()
	for symbol, want := range map[Symbol]string{
		SymbolSuccess: "✓", SymbolError: "✗", SymbolWarning: "⚠", SymbolStarted: "●",
		SymbolUnstarted: "○", SymbolEpic: "◆", SymbolNext: "→",
	} {
		if string(symbol) != want {
			t.Errorf("symbol = %q, want %q", symbol, want)
		}
	}
}

func TestTopBorder(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, left, right string
		width             int
		want              string
	}{
		{"labels", "fid", "Wed", 20, "╭ fid ──────── Wed ╮"},
		{"colored label keeps its width", "\033[1mfid\033[0m", "Wed", 20, "╭ \033[1mfid\033[0m ──────── Wed ╮"},
		{"labels too wide", "a long left label", "right", 12, "╭──────────╮"},
		{"no labels", "", "", 6, "╭────╮"},
		{"too narrow", "fid", "Wed", 2, "╭╮"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := TopBorder(test.left, test.right, test.width)
			if got != test.want {
				t.Fatalf("TopBorder() = %q, want %q", got, test.want)
			}
			if test.width > 2 && VisibleWidth(got) != test.width {
				t.Fatalf("TopBorder() is %d wide, want %d", VisibleWidth(got), test.width)
			}
		})
	}
}

func TestBottomBorder(t *testing.T) {
	t.Parallel()
	for width, want := range map[int]string{6: "╰────╯", 2: "╰╯", 0: "╰╯"} {
		if got := BottomBorder(width); got != want {
			t.Errorf("BottomBorder(%d) = %q, want %q", width, got, want)
		}
	}
}

func TestProgressBar(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name              string
		value, max, width int
		want              string
	}{
		{"partly done", 3, 8, 10, "███░░░░░░░  3/8"},
		{"all done", 8, 8, 10, "██████████  8/8"},
		{"none done", 0, 5, 4, "░░░░  0/5"},
		{"nothing to do", 0, 0, 4, "░░░░  0/0"},
		{"over the total", 9, 8, 4, "████  8/8"},
		{"negative width", 1, 2, -1, "  1/2"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := ProgressBar(test.value, test.max, test.width); got != test.want {
				t.Fatalf("ProgressBar(%d, %d, %d) = %q, want %q", test.value, test.max, test.width, got, test.want)
			}
		})
	}
}

func TestPadding(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, text string
		width      int
		right      string
		left       string
	}{
		{"shorter", "ab", 4, "ab  ", "  ab"},
		{"exact", "abcd", 4, "abcd", "abcd"},
		{"longer ends in an ellipsis", "abcdef", 4, "abc…", "abc…"},
		{"multibyte runes count once", "héllo", 6, "héllo ", " héllo"},
		{"width one", "abc", 1, "…", "…"},
		{"width zero", "abc", 0, "", ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := PadRight(test.text, test.width); got != test.right {
				t.Errorf("PadRight(%q, %d) = %q, want %q", test.text, test.width, got, test.right)
			}
			if got := PadLeft(test.text, test.width); got != test.left {
				t.Errorf("PadLeft(%q, %d) = %q, want %q", test.text, test.width, got, test.left)
			}
		})
	}
}

func TestVisibleWidth(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, text string
		want       int
	}{
		{"plain", "fid", 3},
		{"unicode", "◆ épic", 6},
		{"colored", "\033[1;36mfid\033[0m", 3},
		{"bare escape", "\033x", 2},
		{"unterminated sequence", "a\033[3", 4},
		{"empty", "", 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := VisibleWidth(test.text); got != test.want {
				t.Fatalf("VisibleWidth(%q) = %d, want %d", test.text, got, test.want)
			}
		})
	}
}
