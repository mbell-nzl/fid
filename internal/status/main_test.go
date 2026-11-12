package status

import (
	"os"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestMain(m *testing.M) {
	os.Exit(gittest.Main(m))
}
