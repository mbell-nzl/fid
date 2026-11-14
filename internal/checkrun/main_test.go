package checkrun

import (
	"os"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestMain(m *testing.M) {
	if child, err := LaunchChild(); child {
		if err != nil {
			os.Exit(125)
		}
		os.Exit(0)
	}
	os.Exit(gittest.Main(m))
}
