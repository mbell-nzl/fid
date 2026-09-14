package main

import (
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["resume"] = rootFlagCapability{Project: true, JSON: true}
}

func newResumeCommand() *cobra.Command {
	var repository string
	cmd := &cobra.Command{Use: "resume <issue-id>", Short: "Read bounded durable resume state.", Args: canonicalExactArgs("projection.resume", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderEnvelopeError(cmd, "projection.resume", err, "")
		}
		if repository == "" {
			repository, err = commandWorkingDirectory(cmd)
			if err != nil {
				return renderEnvelopeError(cmd, "projection.resume", err, "")
			}
		}
		resume, err := workflow.ResumeIssueState(workspace, issueMigrationsDir(), args[0], repository)
		if err != nil {
			return renderEnvelopeError(cmd, "projection.resume", err, "")
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "projection.resume", map[string]any{"resume": resume}, nil)
		}
		runID := ""
		if resume.Run != nil {
			runID = resume.Run.ID
		}
		context, err := workflow.ReadRecoveryContext(workspace, runID)
		if err != nil {
			return renderEnvelopeError(cmd, "projection.resume", err, "")
		}
		return writeResumeHuman(cmd, resume, context)
	}}
	cmd.Flags().StringVar(&repository, "repository", "", "Repository path for the durable observation")
	cmd.SetFlagErrorFunc(canonicalFlagError("projection.resume"))
	return cmd
}
