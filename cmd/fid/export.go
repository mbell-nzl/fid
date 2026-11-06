package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/export"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["export"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["export create"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["export validate"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["export import"] = rootFlagCapability{JSON: true}
}

func newExportCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "export", Short: "Create and validate portable recovery archives.", RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	}}
	cmd.AddCommand(newExportCreateCommand(), newExportValidateCommand(), newExportImportCommand())
	return cmd
}

func newExportImportCommand() *cobra.Command {
	var repository, key string
	cmd := &cobra.Command{Use: "import <archive-path>", Short: "Import a portable recovery archive.", Args: canonicalExactArgs("export.import", 1), RunE: func(cmd *cobra.Command, args []string) error {
		for _, flag := range []string{"repository", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "export.import", flag); err != nil {
				return err
			}
		}
		target, err := project.Discover(repository)
		if repositoryRefusalMessage(err) != "" {
			return renderExportError(cmd, "export.import", args[0], err)
		} else if err != nil {
			target = repository // Import reports an unusable target itself.
		}
		result, err := export.Import(export.ImportOptions{ArchivePath: args[0], RepositoryRoot: target, MigrationsDir: issueMigrationsDir(), IdempotencyKey: key})
		if err != nil {
			return renderExportError(cmd, "export.import", args[0], err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "export.import", map[string]any{"import": result}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", result.TargetRoot, result.ProjectID)
		return err
	}}
	cmd.Flags().StringVar(&repository, "repository", "", "Existing empty target repository directory")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Import idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("export.import"))
	return cmd
}

func newExportCreateCommand() *cobra.Command {
	var repository string
	cmd := &cobra.Command{Use: "create <archive-path>", Short: "Create a deterministic export archive.", Args: canonicalExactArgs("export.create", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderExportError(cmd, "export.create", args[0], err)
		}
		_ = repository
		result, err := export.Create(export.CreateOptions{WorkspaceRoot: workspace, ArchivePath: args[0]})
		if err != nil {
			return renderExportError(cmd, "export.create", args[0], err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "export.create", map[string]any{"export": result}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", result.ArchivePath, result.ArchiveDigest)
		return err
	}}
	cmd.Flags().StringVar(&repository, "repository", "", "Repository path for the export observation")
	cmd.SetFlagErrorFunc(canonicalFlagError("export.create"))
	return cmd
}

func newExportValidateCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "validate <archive-path>", Short: "Validate an export archive without mutation.", Args: canonicalExactArgs("export.validate", 1), RunE: func(cmd *cobra.Command, args []string) error {
		result, err := export.Validate(export.ValidateOptions{ArchivePath: args[0], MigrationsDir: issueMigrationsDir()})
		if err != nil {
			return renderExportError(cmd, "export.validate", args[0], err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "export.validate", map[string]any{"validation": result}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", result.ArchivePath, result.Status)
		return err
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("export.validate"))
	return cmd
}

func renderExportError(cmd *cobra.Command, operation, archivePath string, cause error) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var invalid *export.InvalidArchiveError
	if errors.As(cause, &invalid) {
		checks := make([]map[string]any, 0, len(invalid.Checks))
		for _, item := range invalid.Checks {
			var entry any
			if item.Entry != nil {
				entry = *item.Entry
			}
			checks = append(checks, map[string]any{"code": item.Code, "entry": entry, "detail": item.Detail})
		}
		failure, err := cli.NewFailure(cli.CodeArchiveInvalid, "export archive is invalid", &cli.Subject{Type: cli.SubjectExport, ID: archivePath}, map[string]any{"archive_path": archivePath, "checks": checks})
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
		}
		return writeHumanFailure(cmd, operation, failure, cause)
	}
	var targetConflict *export.ImportTargetConflictError
	if errors.As(cause, &targetConflict) {
		failure, err := cli.NewFailure(cli.CodeImportTargetConflict, "import target conflicts with portable recovery archive", &cli.Subject{Type: cli.SubjectExport, ID: archivePath}, map[string]any{"target_root": targetConflict.TargetRoot, "reason": targetConflict.Reason})
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
		}
		return writeHumanFailure(cmd, operation, failure, cause)
	}
	failure, err := cli.NewFailure(cli.CodeInternalError, "export operation failed", nil, map[string]any{})
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	return writeHumanFailure(cmd, operation, failure, cause)
}
