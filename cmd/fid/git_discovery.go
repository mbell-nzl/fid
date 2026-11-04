package main

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/spf13/cobra"
)

// discoverWorkspaceRoot returns the top level of the git worktree that holds
// --project, or the working directory, without requiring a ledger there.
func discoverWorkspaceRoot(cmd *cobra.Command) (string, error) {
	workingDirectory, err := commandWorkingDirectory(cmd)
	if err != nil {
		return "", err
	}
	start, err := rebuildProjectCandidate(cmd, workingDirectory)
	if err != nil {
		return "", err
	}
	return project.Discover(start)
}

// renderRepositoryRefusal writes the refusal for a project command run outside
// every git repository and returns its error. It returns nil for any other
// cause, so each command's renderer handles those as before.
func renderRepositoryRefusal(cmd *cobra.Command, cause error) error {
	var unavailable *project.RepositoryUnavailableError
	if !errors.Is(cause, project.ErrNotARepository) || !errors.As(cause, &unavailable) {
		return nil
	}
	failure, err := cli.NewFailure(cli.CodeRepositoryUnavailable, "fid needs a git repository", nil,
		map[string]any{"reason": unavailable.Reason, "detail": unavailable.Detail})
	if err != nil {
		return err
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), leafOperation(cmd), failure, cause)
	}
	return writeHumanFailure(cmd, leafOperation(cmd), failure, cause)
}
