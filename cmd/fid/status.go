package main

import (
	"errors"
	"io"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/status"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["status"] = rootFlagCapability{Project: true, JSON: true}
}

func newStatusCommand() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show what needs you, what's in progress and what's ready.",
		Long: "Show what needs you, the epics and tasks in progress and ready, and how current the ledger's backup is.\n\n" +
			"READY shows at most 10 tasks; --all shows every one. --json is always complete.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := rebuildWorkspaceRoot(cmd)
			if err != nil {
				return renderStatusFailure(cmd, err)
			}
			environment := commandDisplayEnvironment(cmd)
			result, err := status.Build(status.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(),
				Environment: environment.environ, Now: environment.now()})
			if err != nil {
				return renderStatusFailure(cmd, err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "project.status", result.Fields(), nil)
			}
			return writeHuman(cmd, environment, func(w io.Writer, renderer display.Renderer) error {
				return status.WriteText(w, result, renderer, all)
			})
		},
	}
	cmd.Flags().BoolVar(&all, "all", false, "Show every task to do")
	return cmd
}

func renderStatusFailure(cmd *cobra.Command, cause error) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := statusFailure(cmd, cause)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "project.status", failure, cause)
	}
	return writeHumanFailure(cmd, "project.status", failure, cause)
}

func statusFailure(cmd *cobra.Command, cause error) (cli.Failure, error) {
	if errors.Is(cause, project.ErrNotFound) {
		return issueMutationFailure(cmd, "project.status", "", cause)
	}
	if failure, handled, err := journalReplayStoppedFailure(cause); err != nil || handled {
		return failure, err
	}
	return cli.NewFailure(cli.CodeInternalError, cli.BoundedMessage("status failed: "+cause.Error()), nil, map[string]any{})
}
