package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type versionGateError struct {
	entity  string
	version string
}

func (err *versionGateError) Error() string {
	return fmt.Sprintf("commands for %s ship in %s", err.entity, err.version)
}

func newGovernanceStubCommand(entity, use string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                use,
		Short:              "Command ships in v0.2.",
		Args:               cobra.ArbitraryArgs,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := &versionGateError{entity: entity, version: "v0.2"}
			if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
				return renderErr
			}
			return err
		},
	}

	return cmd
}

func newGovernanceStubFamily(entity string, commands []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   entity,
		Short: "Governance commands ship in v0.2.",
	}
	for _, use := range commands {
		cmd.AddCommand(newGovernanceStubCommand(entity, use))
	}
	return cmd
}

func newDecisionCommand() *cobra.Command {
	return newGovernanceStubFamily("decision", []string{"create", "get", "list"})
}

func newEscalationCommand() *cobra.Command {
	return newGovernanceStubFamily("escalation", []string{"create", "resolve", "get", "list"})
}

func newTriageCommand() *cobra.Command {
	return newGovernanceStubFamily("triage", []string{"record", "start", "verdict"})
}
