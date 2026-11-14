package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["check"] = rootFlagCapability{Project: true, JSON: true} }

func newCheckCommand() *cobra.Command {
	var hold, actor, key string
	cmd := &cobra.Command{Use: "check [name]...", Short: "Execute checks for the held task; defaults to the lane's required checks.", RunE: func(cmd *cobra.Command, args []string) error {
		for _, flag := range []string{"hold", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "check.execute", flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderCheckFailure(cmd, err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		options := checkrun.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), HoldID: hold, Actor: actor, IdempotencyKey: key, Names: args, Clock: commandDisplayEnvironment(cmd).now, Version: checkrun.Version()}
		if !jsonMode {
			options.Before = func(name string) error {
				_, err := fmt.Fprintf(cmd.OutOrStdout(), "● Checking %s…\n", cli.EscapeControlText(name))
				return err
			}
			options.After = func(record checkrecord.Record) error { return writeCheckResult(cmd.OutOrStdout(), record) }
			options.Warning = func(message string) error {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), cli.EscapeControlText(message))
				return err
			}
			if display.IsTerminal(cmd.OutOrStdout()) {
				options.Output = &checkOutput{writer: cmd.OutOrStdout()}
			}
		}
		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer cancel()
		response, err := checkrun.Run(ctx, options)
		if err != nil {
			return renderCheckFailure(cmd, err)
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "check.execute", map[string]any{"task_id": response.TaskID, "checks": response.Checks}, &cli.Event{ID: response.EventID, OccurredAt: response.OccurredAt})
		}
		if len(response.Checks) == 0 {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "✓ No checks required for %s\n", strings.ToUpper(response.Lane[:1])+response.Lane[1:])
			return err
		}
		return nil
	}}
	cmd.Flags().StringVar(&hold, "hold", "", "The hold ID for this session")
	cmd.Flags().StringVar(&actor, "actor", "", "Actor attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	return cmd
}

func renderCheckFailure(cmd *cobra.Command, cause error) error {
	var refused *checkrecord.Refusal
	var failure cli.Failure
	var err error
	if errors.As(cause, &refused) {
		var subject *cli.Subject
		if refused.TaskID != "" {
			subject = &cli.Subject{Type: cli.SubjectIssue, ID: refused.TaskID}
		}
		failure, err = cli.NewFailure(cli.Code(refused.Code), cli.BoundedMessage(refused.Message), subject, refused.Details)
		failure.NextAction = refused.Next
	} else {
		var holdRefusal *holdwork.Refusal
		var source *policy.SourceError
		var denied *persondecision.DeniedError
		if errors.As(cause, &holdRefusal) || errors.As(cause, &source) || errors.As(cause, &denied) {
			return renderHoldFailure(cmd, "check.execute", cause)
		}
		if refusal := renderRepositoryRefusal(cmd, cause); refusal != nil {
			return refusal
		}
		failure, err = issueMutationFailure(cmd, "check.execute", "", cause)
		if err == nil && failure.Code == cli.CodeInvalidArgument {
			failure, err = cli.NewFailure(cli.CodeInternalError, cli.BoundedMessage("check failed: "+cause.Error()), nil, map[string]any{})
		}
	}
	if err != nil {
		return err
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "check.execute", failure, cause)
	}
	if failure.Code == cli.CodeAuthorityDenied {
		failure.Details = map[string]any{"actor_id": failure.Details["actor_id"]}
	}
	if err := cli.WriteHumanFailure(cmd.ErrOrStderr(), "check.execute", cmd.CommandPath(), failure, false); err != nil {
		return err
	}
	return cli.NewFailureError(failure, cause)
}

var checkReasonText = map[string]string{"tree-dirty": "the tree had uncommitted changes", "tree-changed": "the tree changed", "ownership-lost": "ownership was lost", "artifact-not-kept": "a listed file wasn't kept"}

func writeCheckResult(writer io.Writer, record checkrecord.Record) error {
	result := record.Result
	message := "✗ " + record.Start.Name + " " + result.Outcome
	if result.Outcome == "passed" {
		message = "✓ " + record.Start.Name + " passed · counts for hand-off"
		if !result.Counts {
			reasons := []string{}
			for _, reason := range result.Reasons {
				reasons = append(reasons, checkReasonText[reason])
			}
			message = "⚠ " + record.Start.Name + " passed, but doesn't count: " + strings.Join(reasons, ", ")
		}
	}
	if _, err := fmt.Fprintln(writer, cli.EscapeControlText(message)); err != nil {
		return err
	}
	for _, file := range result.Artifacts {
		if _, err := fmt.Fprintln(writer, "  Kept "+cli.EscapeControlText(file.Path)); err != nil {
			return err
		}
	}
	return nil
}

// checkOutput retains split UTF-8 characters while escaping terminal controls.
// Only line feeds are formatting; the child cannot emit terminal instructions.
type checkOutput struct {
	writer  io.Writer
	pending []byte
}

func (output *checkOutput) Write(data []byte) (int, error) {
	output.pending = append(output.pending, data...)
	var text strings.Builder
	for len(output.pending) > 0 && utf8.FullRune(output.pending) {
		r, size := utf8.DecodeRune(output.pending)
		output.pending = output.pending[size:]
		if r == '\n' {
			text.WriteRune(r)
		} else {
			text.WriteString(cli.EscapeControlText(string(r)))
		}
	}
	_, err := io.WriteString(output.writer, text.String())
	return len(data), err
}
