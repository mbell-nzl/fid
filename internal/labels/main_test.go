package labels

import (
	"os"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

// TestMain guards git discovery, because these tests make git repositories.
func TestMain(m *testing.M) {
	os.Exit(gittest.Main(m))
}
