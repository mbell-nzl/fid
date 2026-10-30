package main

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/labels"
	"github.com/spf13/cobra"
)

func renderLabelError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var duplicate *labels.DuplicateLabelError
	if errors.As(cause, &duplicate) {
		failure, err := cli.NewFailure(cli.CodeDuplicateEntity, "entity already exists", &cli.Subject{Type: cli.SubjectLabel, ID: duplicate.Name}, map[string]any{"entity_type": "label", "identity": map[string]any{"fields": []string{"name"}, "values": []string{duplicate.Name}}})
		if err != nil {
			return err
		}
		if !jsonMode {
			return writeHumanFailure(cmd, operation, failure, cause)
		}
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	var limit *labels.ResultLimitExceededError
	if errors.As(cause, &limit) {
		failure, err := cli.NewFailure(cli.CodeResultLimitExceeded, "read result exceeds its collection limit", nil, map[string]any{"collection": limit.Collection, "limit": limit.Limit, "observed": limit.Observed})
		if err != nil {
			return err
		}
		if !jsonMode {
			return writeHumanFailure(cmd, operation, failure, cause)
		}
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	return renderIssueMutationError(cmd, operation, "", cause)
}
