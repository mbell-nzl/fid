package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
	"github.com/mbell-nzl/fid/internal/waiverwork"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["waive"] = rootFlagCapability{Project: true, JSON: true} }

func newWaiveCommandWithTerminal(terminal func(any) bool) *cobra.Command {
	if terminal == nil {
		terminal = func(stream any) bool {
			return initIsTerminal(stream, func(fd uintptr) bool { return isatty.IsTerminal(fd) })
		}
	}
	var actor, reason, key, lane string
	cmd := &cobra.Command{Use: "waive <task>", Short: "Lower a handed-off change’s gates as an interactive person.", Args: canonicalExactArgs("task.waive", 1), RunE: func(cmd *cobra.Command, args []string) error {
		for _, flag := range []string{"actor", "reason", "lane"} {
			if err := requireCanonicalFlag(cmd, "task.waive", flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, "task.waive", err)
		}
		if !policy.Allows(actor, policy.PersonDecisions) {
			return renderHoldFailure(cmd, "task.waive", &persondecision.DeniedError{Actor: actor})
		}
		if strings.TrimSpace(reason) == "" {
			return renderHoldFailure(cmd, "task.waive", &persondecision.ReasonRequiredError{})
		}
		if !terminal(cmd.InOrStdin()) || !terminal(cmd.ErrOrStderr()) {
			return waiveFailure(cmd, fmt.Errorf("standard input and standard error must be terminals"))
		}
		result, err := waiverwork.Record(waiverwork.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), TaskID: args[0], Actor: actor, Reason: reason, Lane: lane, IdempotencyKey: key, Clock: commandDisplayEnvironment(cmd).now, Confirm: func(p waiverwork.Preview) (bool, error) {
			if err := waiverwork.WritePrompt(cmd.ErrOrStderr(), p); err != nil {
				return false, err
			}
			answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
			if errors.Is(err, io.EOF) {
				return false, nil
			}
			if err != nil {
				return false, fmt.Errorf("read waiver confirmation: %w", err)
			}
			return strings.TrimSuffix(strings.TrimSuffix(answer, "\n"), "\r") == args[0], nil
		}})
		if err != nil {
			return waiveFailure(cmd, err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "task.waive", map[string]any{"decision": result.Decision, "occurred_at": result.At}, &cli.Event{ID: result.EventID, OccurredAt: result.At})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ Waived %s · decision %s\n", cli.EscapeControlText(args[0]), cli.EscapeControlText(result.Decision.ID))
		return err
	}}
	cmd.Flags().StringVar(&lane, "lane", "", "Lower lane: routine or standard")
	cmd.Flags().StringVar(&actor, "actor", "", "Person making the decision")
	cmd.Flags().StringVar(&reason, "reason", "", "Why you waive this change")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	return cmd
}

func waiveFailure(cmd *cobra.Command, cause error) error {
	var cancelled *waiverwork.Cancelled
	if errors.As(cause, &cancelled) {
		_, err := fmt.Fprintln(cmd.ErrOrStderr(), "\n"+cancelled.Error())
		return err
	}
	var refused *waiverrecord.Refusal
	if errors.As(cause, &refused) {
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		paths := refused.Paths
		if paths == nil {
			paths = []string{}
		}
		return writeFailure(cmd, jsonMode, "task.waive", cause, cli.Code(refused.Code), refused.Message, nil, map[string]any{"paths": paths})
	}
	return reviewFailure(cmd, "task.waive", cause, "")
}
