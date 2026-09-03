package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/evidence"
	"github.com/spf13/cobra"
)

func init() {
	for _, operation := range []string{"receipt get", "receipt list"} {
		rootFlagCapabilities[operation] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newReceiptCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "receipt", Short: "Read fid-executed Receipts."}
	cmd.AddCommand(newReceiptGetCommand(), newReceiptListCommand())
	return cmd
}

func newReceiptGetCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "get <receipt-id>", Short: "Read one Receipt.", Args: canonicalExactArgs("receipt.get", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderReceiptError(cmd, "receipt.get", err)
		}
		receipt, err := evidence.GetReceipt(workspace, issueMigrationsDir(), args[0])
		if err != nil {
			return renderReceiptError(cmd, "receipt.get", err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "receipt.get", map[string]any{"receipt": receipt}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", receipt.ID, receipt.Outcome)
		return err
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("receipt.get"))
	return cmd
}

func newReceiptListCommand() *cobra.Command {
	var runID, obligationID, outcome string
	var limit int
	cmd := &cobra.Command{Use: "list", Short: "List Receipts.", Args: canonicalNoArgs("receipt.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderReceiptError(cmd, "receipt.list", err)
		}
		items, err := evidence.ListReceipts(evidence.ReceiptListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), RunID: runID, ObligationID: obligationID, Outcome: outcome, Limit: limit})
		if err != nil {
			return renderReceiptError(cmd, "receipt.list", err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "receipt.list", map[string]any{"items": items}, nil)
		}
		for _, receipt := range items {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", receipt.ID, receipt.Outcome); err != nil {
				return err
			}
		}
		return nil
	}}
	cmd.Flags().StringVar(&runID, "run", "", "Run ID")
	cmd.Flags().StringVar(&obligationID, "obligation", "", "Obligation ID")
	cmd.Flags().StringVar(&outcome, "outcome", "", "Receipt outcome")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum Receipts")
	cmd.SetFlagErrorFunc(canonicalFlagError("receipt.list"))
	return cmd
}

func renderReceiptError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var notFound *evidence.ReceiptNotFoundError
	if errors.As(cause, &notFound) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeReceiptNotFound, "Receipt was not found", &cli.Subject{Type: cli.SubjectReceipt, ID: notFound.ID}, map[string]any{"id": notFound.ID})
	}
	return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Receipt request", nil, map[string]any{"argument": "request", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}
