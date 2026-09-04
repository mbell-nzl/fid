package main

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/maintenance"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/spf13/cobra"
)

var errMaintenanceCheckFailed = errors.New("maintenance check found integrity divergence")

const legacyStartedWithoutClaimGraceSeconds = 300

type maintenanceCheckResult struct {
	Status   string                `json:"status"`
	Errors   []maintenance.Finding `json:"errors"`
	Warnings []maintenance.Finding `json:"warnings"`
	Info     maintenance.Info      `json:"info"`
}

func init() {
	rootFlagCapabilities["maintenance check"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["maintenance vacuum"] = rootFlagCapability{Project: true, JSON: true}
}

func newMaintenanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Inspect and maintain the derived projection.",
	}
	cmd.AddCommand(newMaintenanceCheckCommand())
	cmd.AddCommand(newMaintenanceVacuumCommand())
	return cmd
}

func newMaintenanceVacuumCommand() *cobra.Command {
	var actor string

	cmd := &cobra.Command{
		Use:   "vacuum",
		Short: "Compact the derived SQLite projection without changing the journal.",
		Args:  canonicalNoArgs("maintenance.vacuum"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return writeCapabilityUnavailable(cmd, "maintenance.vacuum", "telemetry.trajectory.v1", 25)
		},
	}
	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.SetFlagErrorFunc(canonicalFlagError("maintenance.vacuum"))
	return cmd
}

func newMaintenanceCheckCommand() *cobra.Command {
	var warnAsError bool

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Audit projected workspace integrity without changing the journal.",
		Args:  canonicalNoArgs("maintenance.check"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			workspaceRoot, err := maintenanceCheckWorkspaceRoot(cmd)
			if err != nil {
				return renderMaintenanceCheckError(cmd, err)
			}
			if _, err := journal.EnsureProjectionFresh(workspaceRoot, issueMigrationsDir()); err != nil {
				return renderMaintenanceCheckError(cmd, err)
			}

			paths := config.WorkspacePathsForRoot(workspaceRoot)
			conn, err := fiddb.Open(paths.DatabasePath)
			if err != nil {
				return renderMaintenanceCheckError(cmd, err)
			}
			defer func() { _ = conn.Close() }()

			result, err := maintenance.Check(maintenance.CheckOptions{
				DB:                              conn,
				JournalDir:                      paths.JournalDir,
				StartedWithoutClaimGraceSeconds: legacyStartedWithoutClaimGraceSeconds,
			})
			if err != nil {
				return renderMaintenanceCheckError(cmd, err)
			}

			response := maintenanceCheckResponse(result, warnAsError)
			if response.Status == "fail" {
				return renderMaintenanceDivergence(cmd, response)
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "maintenance.check", maintenanceResultMap(response), nil)
			}
			return writeIssueLine(cmd.OutOrStdout(), "maintenance check: %s\n", response.Status)
		},
	}

	cmd.Flags().BoolVar(&warnAsError, "warn-as-error", false, "Promote warning findings to a failing outcome")
	cmd.SetFlagErrorFunc(canonicalFlagError("maintenance.check"))
	return cmd
}

func maintenanceCheckWorkspaceRoot(cmd *cobra.Command) (string, error) {
	workingDirectory, err := commandWorkingDirectory(cmd)
	if err != nil {
		return "", err
	}
	projectPath, err := rebuildProjectCandidate(cmd, workingDirectory)
	if err != nil {
		return "", err
	}
	return project.Resolve(workingDirectory, projectPath)
}

func maintenanceCheckResponse(result maintenance.Result, warnAsError bool) maintenanceCheckResult {
	response := maintenanceCheckResult{Errors: result.Errors, Warnings: result.Warnings, Info: result.Info}
	if len(result.Errors) > 0 || warnAsError && len(result.Warnings) > 0 {
		response.Status = "fail"
		return response
	}
	if len(result.Warnings) > 0 {
		response.Status = "warn"
		return response
	}
	response.Status = "pass"
	return response
}

func maintenanceResultMap(result maintenanceCheckResult) map[string]any {
	return map[string]any{
		"status": result.Status, "errors": findingsMap(result.Errors), "warnings": findingsMap(result.Warnings),
		"info": map[string]any{"issue_count": result.Info.IssueCount, "claim_count_live": result.Info.ClaimCountLive, "recent_activity_7d": result.Info.RecentActivity},
	}
}

func findingsMap(findings []maintenance.Finding) []map[string]any {
	result := make([]map[string]any, 0, len(findings))
	for _, finding := range findings {
		result = append(result, map[string]any{"check": finding.Check, "message": finding.Message, "entity_id": finding.EntityID})
	}
	return result
}

func renderMaintenanceDivergence(cmd *cobra.Command, result maintenanceCheckResult) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := cli.NewFailure(cli.CodeProjectionDivergence, "maintenance check found integrity divergence", nil, maintenanceResultMap(result))
	if err != nil {
		return err
	}
	if !jsonMode {
		if err := writeIssueLine(cmd.ErrOrStderr(), "error: %s\n", errMaintenanceCheckFailed); err != nil {
			return err
		}
		return cli.NewFailureError(failure, errMaintenanceCheckFailed)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "maintenance.check", failure, errMaintenanceCheckFailed)
}

func renderMaintenanceCheckError(cmd *cobra.Command, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := maintenanceCheckFailure(cmd, cause)
	if err != nil {
		return err
	}
	if !jsonMode {
		if err := renderIssueError(cmd.ErrOrStderr(), cause); err != nil {
			return err
		}
		return cli.NewFailureError(failure, cause)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "maintenance.check", failure, cause)
}

func maintenanceCheckFailure(cmd *cobra.Command, cause error) (cli.Failure, error) {
	if errors.Is(cause, project.ErrNotFound) {
		workingDirectory, err := commandWorkingDirectory(cmd)
		if err != nil {
			return cli.Failure{}, err
		}
		projectPath, err := rebuildProjectCandidate(cmd, workingDirectory)
		if err != nil {
			return cli.Failure{}, err
		}
		return cli.NewFailure(cli.CodeProjectNotFound, "fid project was not found", &cli.Subject{Type: cli.SubjectProject, ID: projectPath}, map[string]any{"path": projectPath})
	}
	if failure, handled, err := journalReplayStoppedFailure(cause); err != nil {
		return cli.Failure{}, err
	} else if handled {
		return failure, nil
	}
	var limit *maintenance.ResultLimitExceededError
	if errors.As(cause, &limit) {
		return cli.NewFailure(cli.CodeResultLimitExceeded, "maintenance findings exceed their collection limit", nil, map[string]any{
			"collection": limit.Collection, "limit": limit.Limit, "observed": limit.Observed,
		})
	}
	return cli.NewFailure(cli.CodeInternalError, "maintenance check failed", nil, map[string]any{})
}
