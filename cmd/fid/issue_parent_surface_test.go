package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestIssueUpdateAndTransitionDoNotRegisterParentMutationFlags(t *testing.T) {
	for _, command := range []*cobra.Command{
		newIssueUpdateCommand(),
		newIssueTransitionCommand(),
	} {
		for _, flagName := range []string{"epic", "parent"} {
			if command.Flags().Lookup(flagName) != nil {
				t.Fatalf("%s registers --%s after Issue creation", command.Name(), flagName)
			}
		}
	}
}
