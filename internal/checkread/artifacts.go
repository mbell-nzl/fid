package checkread

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/artifactstore"
	"github.com/mbell-nzl/fid/internal/config"
)

// ArtifactProblem is a present observation, never a change to a check receipt.
type ArtifactProblem struct {
	Check, ResultID, ArtifactID, Path, Problem string
}

// VerifyKept checks only the required, counting passes selected by Evaluate.
// It reads bytes without creating directories or recording availability. Live
// hand-off and done share this check; historical replay must never call it.
func VerifyKept(root string, evaluation Evaluation) ([]ArtifactProblem, error) {
	problems := []ArtifactProblem{}
	store := ""
	for _, name := range evaluation.Required {
		if !evaluation.Satisfied(name) {
			continue
		}
		record := evaluation.Latest[name]
		for _, artifact := range record.Result.Artifacts {
			if store == "" {
				cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(root), config.ConfigFileName))
				if err != nil {
					return nil, fmt.Errorf("verify kept Artifacts: %w", err)
				}
				store = filepath.Join(config.LedgerDir(root), cfg.ArtifactRoot)
			}
			path := filepath.Join(store, strings.TrimPrefix(artifact.Digest, "sha256:"))
			if err := artifactstore.Verify(path, artifact.Digest); err != nil {
				problems = append(problems, ArtifactProblem{name, record.Start.ID, artifact.ID, artifact.Path, err.Error()})
			}
		}
	}
	return problems, nil
}
