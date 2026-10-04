package main

import (
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	for _, operation := range []string{"attention list", "attention resolve"} {
		rootFlagCapabilities[operation] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newAttentionCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "attention", Short: "Manage Attention items."}
	cmd.AddCommand(newAttentionListCommand(), newAttentionResolveCommand())
	return cmd
}

func newAttentionListCommand() *cobra.Command {
	var subject, status string
	var limit int
	cmd := &cobra.Command{Use: "list", Short: "List Attention items.", Args: canonicalNoArgs("attention.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderStopError(cmd, "attention.list", err)
		}
		items, err := workflow.ListAttention(workflow.StopListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), Subject: subject, Status: status, Limit: limit})
		if err != nil {
			return renderStopError(cmd, "attention.list", err)
		}
		return writeStopSuccess(cmd, "attention.list", map[string]any{"items": items}, "", "")
	}}
	cmd.Flags().StringVar(&subject, "subject", "", "Subject type and ID")
	cmd.Flags().StringVar(&status, "status", "", "Attention status")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum Attention items")
	cmd.SetFlagErrorFunc(canonicalFlagError("attention.list"))
	return cmd
}

func newAttentionResolveCommand() *cobra.Command {
	var input, actor, key string
	cmd := &cobra.Command{Use: "resolve", Short: "Resolve an Attention item.", Args: canonicalNoArgs("attention.resolve"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.AttentionResolveRequest
		if err := stopRequest(cmd, "attention.resolve", input, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderStopError(cmd, "attention.resolve", err)
		}
		result, err := workflow.ResolveAttention(workflow.StopOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderStopError(cmd, "attention.resolve", err)
		}
		return writeStopSuccess(cmd, "attention.resolve", map[string]any{"attention": result.Attention, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, result.EventID, result.Attention.RaisedAt)
	}}
	runMutationFlags(cmd, &input, &actor, &key)
	installInputDescriptor(cmd, inputDescriptorRegistry["attention resolve"])
	return cmd
}
