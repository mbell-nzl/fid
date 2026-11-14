package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/spf13/cobra"
)

func installInputDescriptor(cmd *cobra.Command, descriptor cli.InputDescriptor) {
	cmd.Long = cli.InputHelp(descriptor)
	cmd.Flags().Bool("example", false, "Print an editable JSON input template")
	if flag := cmd.Flags().Lookup("input"); flag != nil {
		flag.Usage = strings.TrimSuffix(cmd.Short, ".") + " JSON input file (or - for stdin)"
	}
	originalArgs := cmd.Args
	cmd.Args = func(cmd *cobra.Command, args []string) error {
		if example, _ := cmd.Flags().GetBool("example"); example {
			return nil // RunE reports the ordered conflict, including arguments.
		}
		return originalArgs(cmd, args)
	}
	originalRun := cmd.RunE
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if example, _ := cmd.Flags().GetBool("example"); example {
			if err := rejectExampleConflicts(cmd, descriptor.Operation, args); err != nil {
				return err
			}
			data, err := cli.InputExample(descriptor)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(data)
			return err
		}
		return originalRun(cmd, args)
	}
}

func rejectExampleConflicts(cmd *cobra.Command, operation string, args []string) error {
	for _, spelling := range []string{"json", "project", "input", "actor", "idempotency-key"} {
		flag := cmd.Flags().Lookup(spelling)
		if spelling == "json" || spelling == "project" {
			flag = cmd.Root().PersistentFlags().Lookup(spelling)
		}
		if flag == nil || !flag.Changed {
			continue
		}
		var value any = flag.Value.String()
		if spelling == "json" {
			value = nil
		}
		return writeExampleConflict(cmd, operation, "--"+spelling, value)
	}
	if len(args) > 0 {
		return writeExampleConflict(cmd, operation, "arguments", args[0])
	}
	return nil
}

func writeExampleConflict(cmd *cobra.Command, operation, spelling string, value any) error {
	message := "--example cannot be combined with " + spelling
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	return writeStopFailure(cmd, jsonMode, operation, fmt.Errorf("%s", message), cli.CodeInvalidArgument, message, nil,
		map[string]any{"argument": spelling, "value": value, "allowed": []string{"--example"}})
}

func readDescriptorInput(stdin io.Reader, path string, descriptor cli.InputDescriptor, target any) error {
	var payload []byte
	var err error
	if path == "-" {
		payload, err = io.ReadAll(io.LimitReader(stdin, 1_048_577))
	} else {
		payload, err = os.ReadFile(path)
	}
	if err != nil {
		return fmt.Errorf("read %s input: %w", descriptor.Operation, err)
	}
	if err := cli.ValidateInput(payload, descriptor); err != nil {
		return err
	}
	if err := json.Unmarshal(payload, target); err != nil {
		return fmt.Errorf("decode %s input: %w", descriptor.Operation, err)
	}
	return nil
}
