package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/pflag"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/spf13/cobra"
)

func canonicalExactArgs(operation string, count int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(count)(cmd, args); err != nil {
			return positionalFailure(cmd, operation, count, args, err.Error())
		}
		return nil
	}
}

func canonicalNoArgs(operation string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return positionalFailure(cmd, operation, 0, args, rootCommandKey(cmd)+" accepts no arguments")
		}
		return nil
	}
}

func positionalFailure(cmd *cobra.Command, operation string, count int, args []string, message string) error {
	var value any
	if len(args) > count {
		value = args[count]
	}
	allowed := strings.Fields(cmd.Use)[1:]
	if count == 0 {
		allowed = []string{"no arguments"}
	}
	return writeGrammarFailure(cmd, operation, cli.CodeInvalidArgument, message,
		map[string]any{"argument": "arguments", "value": value, "allowed": allowed})
}

func canonicalFlagError(operation string) func(*cobra.Command, error) error {
	return func(cmd *cobra.Command, cause error) error {
		argument := "flags"
		var value any
		allowed := []string{}
		var unknown *pflag.NotExistError
		var missing *pflag.ValueRequiredError
		var invalid *pflag.InvalidValueError
		switch {
		case errors.As(cause, &unknown):
			argument = "--" + unknown.GetSpecifiedName()
			if unknown.GetSpecifiedShortnames() != "" {
				argument = "-" + unknown.GetSpecifiedName()
			}
			value = argument
			if !cmd.HasSubCommands() {
				cmd.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) {
					if flag.Name != "help" && flag.Name != "project" && flag.Name != "json" {
						allowed = append(allowed, "--"+flag.Name)
					}
				})
				sort.Strings(allowed)
			}
		case errors.As(cause, &missing):
			argument = "--" + missing.GetFlag().Name
			allowed = flagDomainFor(cmd, missing.GetFlag()).Allowed
		case errors.As(cause, &invalid):
			argument = "--" + invalid.GetFlag().Name
			value = invalid.GetValue()
			allowed = flagDomainFor(cmd, invalid.GetFlag()).Allowed
			var domainError *flagDomainError
			if errors.As(cause, &domainError) {
				cause = fmt.Errorf("invalid argument %q for %q flag", invalid.GetValue(), argument)
			}
		}
		return writeGrammarFailure(cmd, operation, cli.CodeInvalidArgument, cause.Error(),
			map[string]any{"argument": argument, "value": value, "allowed": allowed})
	}
}

func requireCanonicalFlag(cmd *cobra.Command, operation, name string) error {
	if cmd.Flags().Changed(name) {
		return nil
	}
	return writeGrammarFailure(cmd, operation, cli.CodeInvalidArgument, "--"+name+" is required",
		map[string]any{"argument": "--" + name, "value": nil, "allowed": []string{"--" + name}})
}

// CF-047: combination constraints retain their request-level envelope; these
// are not parser, positional, required-flag or flag-domain failures.
const flagCombinationArgument = "request"

func renderCanonicalRequestError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := cli.NewFailure(cli.CodeInvalidArgument, "invalid command request", nil,
		map[string]any{"argument": flagCombinationArgument, "value": nil, "allowed": []string{"accepted " + operation + " input"}})
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	if err := renderIssueError(cmd.ErrOrStderr(), cause); err != nil {
		return err
	}
	return cli.NewFailureError(failure, cause)
}

// Every stable leaf gets grammar validation before its workspace-aware handler.
// Descriptor examples retain their own ordered exclusivity checks.
func installLeafGrammar(cmd *cobra.Command) {
	if cmd.HasSubCommands() {
		cmd.SetFlagErrorFunc(canonicalFlagError(rootCommandOperation(cmd)))
		for _, child := range cmd.Commands() {
			installLeafGrammar(child)
		}
		return
	}
	if rootCommandKey(cmd) == "sql" || rootCommandKey(cmd) == "journal inspect" {
		return
	}
	operation := leafOperation(cmd)
	cmd.SetFlagErrorFunc(canonicalFlagError(operation))
	installFlagDomains(cmd)
	if cmd.DisableFlagParsing {
		return
	} // evidence.evaluate parses example-only flags in its handler.
	count := len(strings.Fields(cmd.Use)) - 1
	validator := canonicalExactArgs(operation, count)
	if count == 0 {
		validator = canonicalNoArgs(operation)
	}
	cmd.Args = func(cmd *cobra.Command, args []string) error {
		if example, _ := cmd.Flags().GetBool("example"); example {
			return nil
		}
		return validator(cmd, args)
	}
}

func leafOperation(cmd *cobra.Command) string {
	switch rootCommandKey(cmd) {
	case "init":
		return "project.init"
	case "rebuild":
		return "project.rebuild"
	case "capabilities":
		return "capabilities.get"
	case "ready":
		return "ready.get"
	case "resume":
		return "projection.resume"
	}
	return rootCommandOperation(cmd)
}

func writeGrammarFailure(cmd *cobra.Command, operation string, code cli.Code, message string, details map[string]any) error {
	message = cli.BoundedMessage(message)
	failure, err := cli.NewFailure(code, message, nil, details)
	if err != nil {
		return err
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	cause := fmt.Errorf("%s", message)
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	if _, err := fmt.Fprintf(cmd.ErrOrStderr(), "error: %s\n", message); err != nil {
		return fmt.Errorf("write grammar failure: %w", err)
	}
	return cli.NewFailureError(failure, cause)
}
