package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/spf13/cobra"
)

func writeCapabilityUnavailable(cmd *cobra.Command, operation, capabilityID string, ownerPhase int) error {
	failure, err := cli.NewFailure(
		cli.CodeCapabilityUnavailable,
		"capability is unavailable",
		&cli.Subject{Type: cli.SubjectCapability, ID: capabilityID},
		map[string]any{"capability_id": capabilityID, "available": false, "owner_phase": ownerPhase},
	)
	if err != nil {
		return err
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, fmt.Errorf("capability %q is unavailable", capabilityID))
	}
	return writeHumanFailure(cmd, operation, failure, fmt.Errorf("capability %q is unavailable", capabilityID))
}
