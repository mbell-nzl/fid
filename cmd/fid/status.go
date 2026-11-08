package main

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/status"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["status"] = rootFlagCapability{Project: true, JSON: true}
}

// statusEnvironmentKey carries a statusEnvironment in the root command's
// context, so tests fix the time, its zone, HOME and the XDG directories.
type statusEnvironmentKey struct{}

type statusEnvironment struct {
	environ []string
	now     func() time.Time
}

func commandStatusEnvironment(cmd *cobra.Command) statusEnvironment {
	if environment, ok := cmd.Root().Context().Value(statusEnvironmentKey{}).(statusEnvironment); ok {
		return environment
	}
	return statusEnvironment{environ: os.Environ(), now: time.Now}
}

func newStatusCommand() *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show what needs you, what's in progress and what's to do.",
		Long: "Show what needs you, the epics and tasks in progress and to do, and how current the ledger's backup is.\n\n" +
			"TO DO shows at most 10 tasks; --all shows every one. --json is always complete.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := rebuildWorkspaceRoot(cmd)
			if err != nil {
				return renderStatusFailure(cmd, err)
			}
			environment := commandStatusEnvironment(cmd)
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
			renderer := display.NewRenderer(display.IsTerminal(cmd.OutOrStdout()), environmentValue(environment.environ, "NO_COLOR"))
			return status.WriteText(cmd.OutOrStdout(), result, renderer, all)
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

// environmentValue is name's value in environ, or "" when it is unset.
func environmentValue(environ []string, name string) string {
	value := ""
	for _, entry := range environ {
		if key, found, ok := strings.Cut(entry, "="); ok && key == name {
			value = found
		}
	}
	return value
}
