package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/mbell-nzl/fid/internal/approvalwork"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["approve"] = rootFlagCapability{Project: true, JSON: true} }

func newApproveCommandWithTerminal(terminal func(any) bool) *cobra.Command {
	if terminal == nil {
		terminal = func(stream any) bool {
			return initIsTerminal(stream, func(fd uintptr) bool { return isatty.IsTerminal(fd) })
		}
	}
	var actor, reason, key string
	cmd := &cobra.Command{Use: "approve <task>", Short: "Approve the full change as an interactive person.", Args: canonicalExactArgs("task.approve", 1), RunE: func(cmd *cobra.Command, args []string) error {
		for _, flag := range []string{"actor", "reason"} {
			if err := requireCanonicalFlag(cmd, "task.approve", flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, "task.approve", err)
		}
		if !policy.Allows(actor, policy.PersonDecisions) {
			return renderHoldFailure(cmd, "task.approve", &persondecision.DeniedError{Actor: actor})
		}
		if strings.TrimSpace(reason) == "" {
			return renderHoldFailure(cmd, "task.approve", &persondecision.ReasonRequiredError{})
		}
		if !terminal(cmd.InOrStdin()) || !terminal(cmd.ErrOrStderr()) {
			return approveFailure(cmd, fmt.Errorf("standard input and standard error must be terminals"))
		}
		result, err := approvalwork.Record(approvalwork.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), TaskID: args[0], Actor: actor, Reason: reason, IdempotencyKey: key, Clock: commandDisplayEnvironment(cmd).now, Confirm: func(p approvalwork.Preview) (bool, error) {
			if err := approvalwork.WritePrompt(cmd.ErrOrStderr(), p); err != nil {
				return false, err
			}
			answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
			if errors.Is(err, io.EOF) {
				return false, nil
			}
			if err != nil {
				return false, fmt.Errorf("read approval confirmation: %w", err)
			}
			return strings.TrimSuffix(strings.TrimSuffix(answer, "\n"), "\r") == args[0], nil
		}})
		if err != nil {
			return approveFailure(cmd, err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "task.approve", map[string]any{"decision": result.Decision, "occurred_at": result.At}, &cli.Event{ID: result.EventID, OccurredAt: result.At})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ Approved %s · decision %s\n", cli.EscapeControlText(args[0]), cli.EscapeControlText(result.Decision.ID))
		return err
	}}
	cmd.Flags().StringVar(&actor, "actor", "", "Person making the decision")
	cmd.Flags().StringVar(&reason, "reason", "", "Why you approve this change")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	return cmd
}

func approveFailure(cmd *cobra.Command, cause error) error {
	var cancelled *approvalwork.Cancelled
	if errors.As(cause, &cancelled) {
		_, err := fmt.Fprintln(cmd.ErrOrStderr(), "\n"+cancelled.Error())
		return err
	}
	return reviewFailure(cmd, "task.approve", cause, "")
}
