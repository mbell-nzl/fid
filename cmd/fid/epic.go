package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/epicwork"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
	"github.com/spf13/cobra"
)

func init() {
	for _, leaf := range []string{"epic approve", "epic validation add", "epic validation update", "epic validation remove"} {
		rootFlagCapabilities[leaf] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newEpicCommandWithTerminal(terminal func(any) bool) *cobra.Command {
	if terminal == nil {
		terminal = func(stream any) bool {
			return initIsTerminal(stream, func(fd uintptr) bool { return isatty.IsTerminal(fd) })
		}
	}
	cmd := &cobra.Command{Use: "epic", Short: "Approve epic plans and record validation evidence."}
	validation := &cobra.Command{Use: "validation", Short: "Maintain an epic's validation list."}
	for _, action := range []string{"add", "update", "remove"} {
		validation.AddCommand(newEpicValidationChange(action))
	}
	validation.AddCommand(newEpicEvidenceCommand("pass"), newEpicEvidenceCommand("verify"))
	cmd.AddCommand(validation, newEpicApproveCommand(terminal))
	return cmd
}

func newEpicValidationChange(action string) *cobra.Command {
	var o epicwork.Options
	op := "epic.validation." + action
	use, n := action+" <epic>", 1
	if action != "add" {
		use += " <id>"
		n = 2
	}
	cmd := &cobra.Command{Use: use, Short: "Change an epic validation.", Args: canonicalExactArgs(op, n), RunE: func(cmd *cobra.Command, args []string) error {
		required := []string{"actor", "idempotency-key"}
		if action != "remove" {
			required = append(required, "claim", "verifier")
		}
		for _, flag := range required {
			if err := requireCanonicalFlag(cmd, op, flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, op, err)
		}
		o.WorkspaceRoot, o.MigrationsFS, o.EpicID, o.Clock = root, fidmigrations.FS, args[0], commandDisplayEnvironment(cmd).now
		if len(args) > 1 {
			o.ValidationID = args[1]
		}
		o.RecipeSet = cmd.Flags().Changed("recipe")
		r, err := epicwork.Change(o, action)
		if err != nil {
			return epicFailure(cmd, op, o.Actor, err)
		}
		return writeEpicResult(cmd, op, r)
	}}
	cmd.Flags().StringVar(&o.Actor, "actor", "", "Who changes the validation")
	cmd.Flags().StringVar(&o.IdempotencyKey, "idempotency-key", "", "Retry key")
	if action != "remove" {
		cmd.Flags().StringVar(&o.Validation.Claim, "claim", "", "Claim to validate")
		cmd.Flags().StringVar(&o.Validation.Verifier, "verifier", "", "check:<name>, review:<profile> or human")
		cmd.Flags().StringVar(&o.Validation.Recipe, "recipe", "", "Project-relative fixture recipe")
	}
	return cmd
}

func newEpicApproveCommand(terminal func(any) bool) *cobra.Command {
	var o epicwork.Options
	const op = "epic.approve"
	cmd := &cobra.Command{Use: "approve <epic>", Short: "Approve the previewed epic plan as a person.", Args: canonicalExactArgs(op, 1), RunE: func(cmd *cobra.Command, args []string) error {
		for _, flag := range []string{"actor", "reason"} {
			if err := requireCanonicalFlag(cmd, op, flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, op, err)
		}
		if !policy.Allows(o.Actor, policy.PersonDecisions) {
			return epicFailure(cmd, op, o.Actor, &persondecision.DeniedError{Actor: o.Actor})
		}
		if !terminal(cmd.InOrStdin()) || !terminal(cmd.ErrOrStderr()) {
			return epicFailure(cmd, op, o.Actor, fmt.Errorf("standard input and standard error must be terminals"))
		}
		o.WorkspaceRoot, o.MigrationsFS, o.EpicID, o.Clock = root, fidmigrations.FS, args[0], commandDisplayEnvironment(cmd).now
		cancelled := false
		o.Confirm = func(p epicplan.Plan) (bool, error) {
			if err := writeEpicPreview(cmd, p); err != nil {
				return false, err
			}
			line, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				return false, fmt.Errorf("read epic approval confirmation: %w", err)
			}
			yes := err == nil && strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r") == p.EpicID
			cancelled = !yes
			return yes, nil
		}
		r, err := epicwork.Approve(o)
		if cancelled {
			_, err = fmt.Fprintln(cmd.ErrOrStderr(), "Epic approval cancelled; nothing recorded")
			return err
		}
		if err != nil {
			return epicFailure(cmd, op, o.Actor, err)
		}
		return writeEpicResult(cmd, op, r)
	}}
	cmd.Flags().StringVar(&o.Actor, "actor", "", "Person approving the epic")
	cmd.Flags().StringVar(&o.Reason, "reason", "", "Why you approve the plan")
	cmd.Flags().StringVar(&o.IdempotencyKey, "idempotency-key", "", "Retry key")
	return cmd
}

func writeEpicPreview(cmd *cobra.Command, p epicplan.Plan) error {
	var b strings.Builder
	fmt.Fprintf(&b, "Approve epic %s\nPlan:\n", cli.EscapeControlText(p.EpicID))
	for _, t := range p.Tasks {
		accept := ""
		if t.Accept {
			accept = " · will accept"
		}
		lane := ""
		if t.Kind == "task" {
			lane = t.Lane
			if lane == "" {
				lane = "not set"
			}
			lane = " · " + lane
		}
		fmt.Fprintf(&b, "  %s %s · %s%s%s\n", t.ID, cli.EscapeControlText(t.Title), t.State, lane, accept)
	}
	fmt.Fprintln(&b, "VALIDATIONS:")
	for _, v := range p.Validations {
		fmt.Fprintf(&b, "  %s %s · %s · %s\n", v.ID, cli.EscapeControlText(v.Claim), cli.EscapeControlText(v.Verifier), cli.EscapeControlText(v.Recipe))
	}
	fmt.Fprintf(&b, "Type %s to approve: ", p.EpicID)
	_, err := fmt.Fprint(cmd.ErrOrStderr(), b.String())
	return err
}

func epicFailure(cmd *cobra.Command, op, actor string, cause error) error {
	var source *policy.SourceError
	if errors.As(cause, &source) {
		return renderHoldFailure(cmd, op, cause)
	}
	if !policy.Allows(actor, policy.WorkSet) || (op == "epic.approve" || op == "epic.validation.pass") && !policy.Allows(actor, policy.PersonDecisions) {
		return renderHoldFailure(cmd, op, &persondecision.DeniedError{Actor: actor})
	}
	return reviewFailure(cmd, op, cause, "")
}

func writeEpicResult(cmd *cobra.Command, op string, r epicwork.Result) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), op, map[string]any{"epic_id": r.EpicID, "validation": r.Validation, "decision_id": r.DecisionID, "accepted": r.Accepted, "checks": r.Checks, "warnings": r.Warnings, "occurred_at": r.At}, &cli.Event{ID: r.EventID, OccurredAt: r.At})
	}
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ %s · %s\n", op, cli.EscapeControlText(r.EpicID))
	for _, check := range r.Checks {
		if err == nil {
			mark := "✓"
			if check.Outcome != "passed" {
				mark = "✗"
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s · %s\n", mark, check.ValidationID, check.Outcome, check.Commit)
		}
	}
	for _, warning := range r.Warnings {
		if err == nil {
			_, err = fmt.Fprintln(cmd.OutOrStdout(), "⚠ "+cli.EscapeControlText(warning))
		}
	}
	return err
}
