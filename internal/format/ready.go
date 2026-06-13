package format

import "github.com/mbell-nzl/fid/internal/issues"

type ReadyView struct {
	Task *issues.ReadyTask `json:"task"`
}
