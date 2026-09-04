package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/evidence"
	"github.com/spf13/cobra"
)

func init() {
	for _, operation := range []string{"evidence link", "evidence get", "evidence evaluate"} {
		rootFlagCapabilities[operation] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newEvidenceCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "evidence", Short: "Link and evaluate trusted gate Evidence."}
	cmd.AddCommand(newEvidenceLinkCommand(), newEvidenceGetCommand(), newEvidenceEvaluateCommand())
	return cmd
}

func newEvidenceLinkCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "link", Short: "Link one exact fid-executed Receipt.", Args: canonicalNoArgs("evidence.link"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "evidence.link", flag); err != nil {
				return err
			}
		}
		var request evidence.EvidenceLinkRequest
		if err := decodeClosedGateInput(cmd.InOrStdin(), inputPath, &request); err != nil {
			return renderEvidenceError(cmd, "evidence.link", err)
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderEvidenceError(cmd, "evidence.link", err)
		}
		result, err := evidence.LinkEvidence(evidence.EvidenceLinkOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key, Request: request})
		if err != nil {
			return renderEvidenceError(cmd, "evidence.link", err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "evidence.link", map[string]any{"evidence": result.Evidence, "evaluation": result.Evaluation, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.EventOccurredAt})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", result.Evidence.ID, result.Evaluation.Status)
		return err
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Evidence link JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("evidence.link"))
	return cmd
}

func newEvidenceGetCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "get <evidence-id>", Short: "Read one Evidence record.", Args: canonicalExactArgs("evidence.get", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderEvidenceError(cmd, "evidence.get", err)
		}
		record, err := evidence.GetEvidence(workspace, issueMigrationsDir(), args[0])
		if err != nil {
			return renderEvidenceError(cmd, "evidence.get", err)
		}
		evaluation, err := evidence.EvaluateEvidence(workspace, issueMigrationsDir(), evidence.EvidenceEvaluateRequest{ObligationID: record.ObligationID, TargetRepository: record.TargetRepository})
		if err != nil {
			return renderEvidenceError(cmd, "evidence.get", err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "evidence.get", map[string]any{"evidence": record, "evaluation": evaluation}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", record.ID, evaluation.Status)
		return err
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("evidence.get"))
	return cmd
}

func newEvidenceEvaluateCommand() *cobra.Command {
	var inputPath string
	cmd := &cobra.Command{Use: "evaluate", Short: "Evaluate one Obligation without mutation.", Args: canonicalNoArgs("evidence.evaluate"), RunE: func(cmd *cobra.Command, _ []string) error {
		if err := requireCanonicalFlag(cmd, "evidence.evaluate", "input"); err != nil {
			return err
		}
		var request evidence.EvidenceEvaluateRequest
		if err := decodeClosedGateInput(cmd.InOrStdin(), inputPath, &request); err != nil {
			return renderEvidenceError(cmd, "evidence.evaluate", err)
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderEvidenceError(cmd, "evidence.evaluate", err)
		}
		result, err := evidence.EvaluateEvidence(workspace, issueMigrationsDir(), request)
		if err != nil {
			return renderEvidenceError(cmd, "evidence.evaluate", err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "evidence.evaluate", map[string]any{"evaluation": result}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", result.ObligationID, result.Status)
		return err
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Evidence evaluation JSON input")
	cmd.SetFlagErrorFunc(canonicalFlagError("evidence.evaluate"))
	return cmd
}

func renderEvidenceError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var notFound *evidence.EvidenceNotFoundError
	if errors.As(cause, &notFound) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeEvidenceNotFound, "Evidence was not found", &cli.Subject{Type: cli.SubjectEvidence, ID: notFound.ID}, map[string]any{"id": notFound.ID})
	}
	return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Evidence request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}
