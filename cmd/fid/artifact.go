package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/evidence"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["artifact verify"] = rootFlagCapability{Project: true, JSON: true}
}

func newArtifactCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "artifact", Short: "Inspect retained gate Artifacts."}
	cmd.AddCommand(newArtifactVerifyCommand())
	return cmd
}

func newArtifactVerifyCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "verify", Short: "Observe retained Artifact availability.", Args: canonicalNoArgs("artifact.verify"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "artifact.verify", flag); err != nil {
				return err
			}
		}
		var request evidence.ArtifactVerifyRequest
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry["artifact verify"], &request); err != nil {
			return renderGateError(cmd, "artifact.verify", err)
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderGateError(cmd, "artifact.verify", err)
		}
		result, err := evidence.VerifyArtifact(evidence.ArtifactVerifyOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key, Request: request})
		if err != nil {
			return renderGateError(cmd, "artifact.verify", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "artifact.verify", map[string]any{"artifact": result.Artifact, "changed": result.Changed, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.EventOccurredAt})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", result.Artifact.ID, result.Artifact.Availability)
		return err
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Artifact verify JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("artifact.verify"))
	installInputDescriptor(cmd, inputDescriptorRegistry["artifact verify"])
	return cmd
}
