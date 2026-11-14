package export

import (
	"os"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// TestMain guards git discovery, because these tests make git repositories.
func TestMain(m *testing.M) {
	if child, err := checkrun.LaunchChild(); child {
		if err != nil {
			os.Exit(125)
		}
		os.Exit(0)
	}
	os.Exit(gittest.Main(m))
}
