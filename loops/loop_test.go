package loops

import (
	"os/exec"
	"testing"
)

func TestV9LoopSelftest(t *testing.T) {
	t.Parallel()
	output, err := exec.Command("sh", "selftest.sh").CombinedOutput()
	if err != nil {
		t.Fatalf("loop selftest: %v\n%s", err, output)
	}
}
