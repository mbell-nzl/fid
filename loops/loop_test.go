package loops

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestV9LoopSelftest(t *testing.T) {
	t.Parallel()
	output, err := exec.Command("sh", "selftest.sh").CombinedOutput()
	if err != nil {
		t.Fatalf("loop selftest: %v\n%s", err, output)
	}
}

func TestLoopLineBudget(t *testing.T) {
	t.Parallel()
	raw, err := os.ReadFile("loop.py")
	if err != nil {
		t.Fatal(err)
	}
	if lines := len(strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n")); lines >= 150 {
		t.Fatalf("loop.py has %d lines; must stay under 150", lines)
	}
}
