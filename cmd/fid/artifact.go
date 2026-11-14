package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/artifactverify"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["artifact verify"] = rootFlagCapability{Project: true, JSON: true}
}

func newArtifactCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "artifact", Short: "Inspect retained check Artifacts."}
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
		var request artifactverify.Request
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry["artifact verify"], &request); err != nil {
			return renderArtifactError(cmd, "artifact.verify", err)
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderArtifactError(cmd, "artifact.verify", err)
		}
		result, err := artifactverify.Verify(artifactverify.Options{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key, Request: request})
		if err != nil {
			return renderArtifactError(cmd, "artifact.verify", err)
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
	cmd.Flags().StringVar(&actor, "actor", "", "Who verifies, agent:<id> or human:<id>")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("artifact.verify"))
	installInputDescriptor(cmd, inputDescriptorRegistry["artifact verify"])
	return cmd
}

func renderArtifactError(cmd *cobra.Command, operation string, cause error) error {
	var denied *artifactverify.DeniedActorError
	if errors.As(cause, &denied) {
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		return writeFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, denied.Error(), nil, map[string]any{"actor_id": denied.Actor, "operation": operation, "required_actions": []string{operation}, "required_capabilities": []string{operation}})
	}
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var structural *cli.InputError
	if errors.As(cause, &structural) {
		return writeFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, structural.Error(), nil, structural.Details())
	}
	var conflict *journal.IdempotencyConflictError
	if errors.As(cause, &conflict) {
		return writeFailure(cmd, jsonMode, operation, cause, cli.CodeIdempotencyConflict, "idempotency key is bound to another request", nil, idempotencyConflictDetails(conflict))
	}
	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid gate request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted gate registration input"}})
}
