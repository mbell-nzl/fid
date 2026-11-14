package checkrecord

import "github.com/mbell-nzl/fid/internal/config"

// DefinitionFor is shared by execution and result selection so the digest has
// exactly the same scope in both places. Scheduling metadata is excluded.
func DefinitionFor(check config.Check) Definition {
	return Definition{Run: check.Run, Timeout: check.Timeout, Environment: EnvironmentNames(), Artifacts: append([]string{}, check.Artifacts...)}
}
