package main

import (
	"errors"
	"io"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/show"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["show"] = rootFlagCapability{Project: true, JSON: true}
}

func newShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show an issue: its state, what blocks it, its criteria and its children.",
		Long: "Show an issue by its exact ID: its state, its epic and task, what blocks it and what it blocks, " +
			"its acceptance criteria, and its tasks or subtasks.\n\n" +
			"The text lists only blockers that aren't completed or cancelled; --json lists every link.",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := rebuildWorkspaceRoot(cmd)
			if err != nil {
				return renderShowFailure(cmd, err)
			}
			result, err := show.Build(show.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), ID: args[0]})
			if err != nil {
				return renderShowFailure(cmd, err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "issue.show", result.Fields(), nil)
			}
			return writeHuman(cmd, commandDisplayEnvironment(cmd), func(w io.Writer, renderer display.Renderer) error {
				return show.WriteText(w, result, renderer)
			})
		},
	}
}

func renderShowFailure(cmd *cobra.Command, cause error) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := showFailure(cmd, cause)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "issue.show", failure, cause)
	}
	return writeHumanFailure(cmd, "issue.show", failure, cause)
}

func showFailure(cmd *cobra.Command, cause error) (cli.Failure, error) {
	if errors.Is(cause, project.ErrNotFound) {
		return issueMutationFailure(cmd, "issue.show", "", cause)
	}
	var notFound *issues.NotFoundError
	if errors.As(cause, &notFound) {
		return cli.NewFailure(cli.CodeIssueNotFound, "issue was not found", &cli.Subject{Type: cli.SubjectIssue, ID: notFound.IssueID}, map[string]any{"id": notFound.IssueID})
	}
	if failure, handled, err := journalReplayStoppedFailure(cause); err != nil || handled {
		return failure, err
	}
	return cli.NewFailure(cli.CodeInternalError, cli.BoundedMessage("show failed: "+cause.Error()), nil, map[string]any{})
}
