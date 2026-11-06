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

// renderRepositoryRefusal writes the refusal for a project command whose
// discovery refused, and returns its error: outside every git repository, or
// where the ledger's location is unsafe. It returns nil for any other cause,
// so each command's renderer handles those as before.
func renderRepositoryRefusal(cmd *cobra.Command, cause error) error {
	message := repositoryRefusalMessage(cause)
	var unavailable *project.RepositoryUnavailableError
	if message == "" || !errors.As(cause, &unavailable) {
		return nil
	}
	failure, err := cli.NewFailure(cli.CodeRepositoryUnavailable, message, nil,
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

// repositoryRefusalMessage returns the message for a discovery refusal, or ""
// when cause is not one.
func repositoryRefusalMessage(cause error) string {
	switch {
	case errors.Is(cause, project.ErrNotARepository):
		return "fid needs a git repository"
	case errors.Is(cause, project.ErrUnsafeLedger):
		return "the ledger location is unsafe"
	}
	return ""
}
