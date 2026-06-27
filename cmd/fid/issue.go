package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/claims"
	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/spf13/cobra"
)

func newIssueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues.",
	}

	cmd.AddCommand(newIssueACCommand())
	cmd.AddCommand(newIssueClaimCommand())
	cmd.AddCommand(newIssueCreateCommand())
	cmd.AddCommand(newIssueCompleteCommand())
	cmd.AddCommand(newIssueUnclaimCommand())
	cmd.AddCommand(newIssueLabelCommand())
	cmd.AddCommand(newIssueUpdateCommand())
	cmd.AddCommand(newIssueTransitionCommand())
	cmd.AddCommand(newIssueGetCommand())
	cmd.AddCommand(newIssueListCommand())
	return cmd
}

func newIssueACCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ac",
		Short: "Manage issue acceptance criteria.",
	}

	cmd.AddCommand(newIssueACAddCommand())
	cmd.AddCommand(newIssueACListCommand())
	cmd.AddCommand(newIssueACVerifyCommand())
	return cmd
}

func newIssueACAddCommand() *cobra.Command {
	var (
		text      string
		ordinal   int
		actorRole string
	)

	cmd := &cobra.Command{
		Use:   "add <issue-id>",
		Short: "Add an acceptance criterion to a task or subtask.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("text") {
				err := fmt.Errorf("--text is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			if !cmd.Flags().Changed("ordinal") {
				err := fmt.Errorf("--ordinal is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			if !cmd.Flags().Changed("by") {
				err := fmt.Errorf("--by is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.AddAcceptanceCriterion(issues.AddAcceptanceCriterionOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				Text:          text,
				Ordinal:       ordinal,
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(
				cmd.ErrOrStderr(),
				"added AC %s on %s (ordinal %d): %q\n",
				result.ID,
				result.IssueID,
				result.Ordinal,
				result.Text,
			)
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "Acceptance criterion text")
	cmd.Flags().IntVar(&ordinal, "ordinal", 0, "Acceptance criterion display order")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newIssueACVerifyCommand() *cobra.Command {
	var (
		actorRole    string
		status       string
		evidenceType string
		evidenceRef  string
		notes        string
	)

	cmd := &cobra.Command{
		Use:   "verify <issue-id> <ac-id>",
		Short: "Verify an acceptance criterion.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("by") {
				err := fmt.Errorf("--by is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			if !cmd.Flags().Changed("status") {
				err := fmt.Errorf("--status is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.VerifyAcceptanceCriterion(issues.VerifyAcceptanceCriterionOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				CriterionID:   args[1],
				ActorRole:     actorRole,
				Status:        status,
				EvidenceType:  issueCreateOptionalString(cmd, "evidence-type", evidenceType),
				EvidenceRef:   issueCreateOptionalString(cmd, "evidence-ref", evidenceRef),
				Notes:         issueCreateOptionalString(cmd, "notes", notes),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeIssueLine(
				cmd.ErrOrStderr(),
				"verified %s on %s by %s: %s\n",
				result.CriterionID,
				result.IssueID,
				result.ActorRole,
				result.Status,
			); err != nil {
				return err
			}
			if result.EvidenceType != nil && result.EvidenceRef != nil {
				return writeIssueLine(cmd.ErrOrStderr(), "  evidence: %s → %s\n", *result.EvidenceType, *result.EvidenceRef)
			}
			if result.EvidenceType != nil {
				return writeIssueLine(cmd.ErrOrStderr(), "  evidence: %s\n", *result.EvidenceType)
			}
			if result.EvidenceRef != nil {
				return writeIssueLine(cmd.ErrOrStderr(), "  evidence: %s\n", *result.EvidenceRef)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Verifying role")
	cmd.Flags().StringVar(&status, "status", "", "Verification status")
	cmd.Flags().StringVar(&evidenceType, "evidence-type", "", "Evidence type")
	cmd.Flags().StringVar(&evidenceRef, "evidence-ref", "", "Evidence reference")
	cmd.Flags().StringVar(&notes, "notes", "", "Verification notes")
	return cmd
}

func newIssueACListCommand() *cobra.Command {
	var formatName string

	cmd := &cobra.Command{
		Use:   "list <issue-id>",
		Short: "List acceptance criteria for an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.ListAcceptanceCriteria(issues.ListAcceptanceCriteriaOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, result); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}

func newIssueClaimCommand() *cobra.Command {
	var (
		actorRole string
		ttl       string
		force     bool
		renew     bool
		note      string
	)

	cmd := &cobra.Command{
		Use:   "claim <issue-id>",
		Short: "Acquire or renew an issue claim.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("by") {
				err := fmt.Errorf("--by is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			ttlSeconds, err := issueClaimTTLSeconds(cmd, ttl)
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			if force && renew {
				err := fmt.Errorf("--force and --renew are mutually exclusive")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if renew {
				result, err := claims.Renew(claims.RenewOptions{
					WorkspaceRoot: workspaceRoot,
					MigrationsDir: issueMigrationsDir(),
					IssueID:       args[0],
					ActorRole:     actorRole,
					TTLSeconds:    ttlSeconds,
				})
				if err != nil {
					if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
						return renderErr
					}
					return err
				}

				return writeIssueLine(cmd.ErrOrStderr(), "renewed claim on %s by %s (expires %s)\n", result.RefID, result.ClaimedBy, result.ExpiresAt)
			}

			result, err := claims.Acquire(claims.AcquireOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				ActorRole:     actorRole,
				TTLSeconds:    ttlSeconds,
				Note:          issueCreateOptionalString(cmd, "note", note),
				Force:         force,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if force && result.Previous != nil {
				return writeIssueLine(
					cmd.ErrOrStderr(),
					"force-acquired claim on %s by %s (previously held by %s until %s)\n",
					result.RefID,
					result.ClaimedBy,
					result.Previous.ClaimedBy,
					result.Previous.ExpiresAt,
				)
			}
			return writeIssueLine(cmd.ErrOrStderr(), "claimed %s by %s (expires %s)\n", result.RefID, result.ClaimedBy, result.ExpiresAt)
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Claiming role")
	cmd.Flags().StringVar(&ttl, "ttl", "30m", "Lease duration")
	cmd.Flags().BoolVar(&force, "force", false, "Force-acquire a live claim")
	cmd.Flags().BoolVar(&renew, "renew", false, "Renew the current live claim")
	cmd.Flags().StringVar(&note, "note", "", "Optional claim context")
	return cmd
}

func newIssueCreateCommand() *cobra.Command {
	var kind string
	var title string
	var actorRole string
	var epicID string
	var parentID string
	var priority int
	var sequence int
	var specMDPath string
	var estTokens int
	var dueDate string
	var labels []string
	var inbox bool
	var ownerRole string
	var ownerActor string
	var scheduling string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an issue.",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			issueParentID, err := issueCreateParentID(kind, epicID, parentID)
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			result, err := issues.Create(issues.CreateOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				Kind:          kind,
				Title:         title,
				ActorRole:     actorRole,
				ParentID:      issueParentID,
				Priority:      issueCreateOptionalInt(cmd, "priority", priority),
				Sequence:      issueCreateOptionalInt(cmd, "sequence", sequence),
				SpecMDPath:    issueCreateOptionalString(cmd, "spec-md", specMDPath),
				EstTokens:     issueCreateOptionalInt(cmd, "est-tokens", estTokens),
				DueDate:       issueCreateOptionalString(cmd, "due-date", dueDate),
				Labels:        labels,
				InInbox:       inbox,
				OwnerRole:     issueCreateOptionalString(cmd, "owner-role", ownerRole),
				OwnerActor:    issueCreateOptionalString(cmd, "owner-actor", ownerActor),
				Scheduling:    scheduling,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "created issue %s (%s) %q\n", result.ID, result.Kind, result.Title)
		},
	}

	cmd.Flags().StringVar(&kind, "kind", "", "Issue kind")
	cmd.Flags().StringVar(&title, "title", "", "Issue title")
	cmd.Flags().StringVar(&epicID, "epic", "", "Parent epic ID")
	cmd.Flags().StringVar(&parentID, "parent", "", "Parent task ID")
	cmd.Flags().IntVar(&priority, "priority", 0, "Issue priority")
	cmd.Flags().IntVar(&sequence, "sequence", 0, "Issue sequence")
	cmd.Flags().StringVar(&specMDPath, "spec-md", "", "Spec markdown path")
	cmd.Flags().IntVar(&estTokens, "est-tokens", 0, "Estimated tokens")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date")
	cmd.Flags().StringSliceVar(&labels, "label", nil, "Issue label")
	cmd.Flags().BoolVar(&inbox, "inbox", false, "Add issue to inbox")
	cmd.Flags().StringVar(&ownerRole, "owner-role", "", "Owner role")
	cmd.Flags().StringVar(&ownerActor, "owner-actor", "", "Owner actor")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.Flags().StringVar(&scheduling, "scheduling", "", "Epic scheduling")
	return cmd
}

func newIssueUnclaimCommand() *cobra.Command {
	var actorRole string

	cmd := &cobra.Command{
		Use:   "unclaim <issue-id>",
		Short: "Release the caller's current issue claim.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("by") {
				err := fmt.Errorf("--by is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := claims.Release(claims.ReleaseOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "released claim on %s by %s\n", result.RefID, result.ReleasedBy)
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role releasing its claim")
	return cmd
}

func newIssueCompleteCommand() *cobra.Command {
	var (
		actorRole         string
		notesMDPath       string
		actualTokens      int
		commits           string
		filesTouched      []string
		humanIntervention bool
		interventionNote  string
	)

	cmd := &cobra.Command{
		Use:   "complete <issue-id>",
		Short: "Record a task completion attempt.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("by") {
				err := fmt.Errorf("--by is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.Complete(issues.CompleteOptions{
				WorkspaceRoot:         workspaceRoot,
				MigrationsDir:         issueMigrationsDir(),
				IssueID:               args[0],
				ActorRole:             actorRole,
				NotesMDPath:           issueCreateOptionalString(cmd, "notes-md", notesMDPath),
				ActualTokens:          issueCreateOptionalInt(cmd, "tokens", actualTokens),
				CommitHashes:          issueSplitCommaSeparated(commits),
				FilesTouched:          filesTouched,
				HumanIntervention:     humanIntervention,
				HumanInterventionNote: issueCreateOptionalString(cmd, "intervention-note", interventionNote),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeIssueLine(cmd.ErrOrStderr(), "completed %s (attempt %d)\n", result.IssueID, result.AttemptNumber); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  completion: %s\n", result.ID); err != nil {
				return err
			}
			if result.ActualTokens != nil {
				if err := writeIssueLine(cmd.ErrOrStderr(), "  tokens: %d\n", *result.ActualTokens); err != nil {
					return err
				}
			}
			if len(result.CommitHashes) > 0 {
				if err := writeIssueLine(cmd.ErrOrStderr(), "  commits: %s\n", joinCommaSeparated(result.CommitHashes)); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.Flags().StringVar(&notesMDPath, "notes-md", "", "Path to completion notes")
	cmd.Flags().IntVar(&actualTokens, "tokens", 0, "Actual tokens used")
	cmd.Flags().StringVar(&commits, "commits", "", "Comma-separated commit hashes")
	cmd.Flags().StringSliceVar(&filesTouched, "files", nil, "Files touched")
	cmd.Flags().BoolVar(&humanIntervention, "human-intervention", false, "Mark that a human intervened during the attempt")
	cmd.Flags().StringVar(&interventionNote, "intervention-note", "", "Short note stored when --human-intervention is set")
	return cmd
}

func newIssueGetCommand() *cobra.Command {
	var formatName string
	var sectionName string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an issue by ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			issue, err := issues.Get(issues.GetOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ID:            args[0],
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			outputValue := any(issue)
			if strings.TrimSpace(sectionName) != "" {
				outputValue, err = issues.SectionValue(issue, sectionName)
				if err != nil {
					if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
						return renderErr
					}
					return err
				}
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, outputValue); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	cmd.Flags().StringVar(&sectionName, "section", "", "Issue detail section")
	return cmd
}

func newIssueUpdateCommand() *cobra.Command {
	var (
		title      string
		priority   int
		labels     []string
		ownerRole  string
		ownerActor string
		specMDPath string
		estTokens  int
		dueDate    string
		actorRole  string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.Update(issues.UpdateOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ID:            args[0],
				ActorRole:     actorRole,
				Title:         issueCreateOptionalString(cmd, "title", title),
				Priority:      issueCreateOptionalInt(cmd, "priority", priority),
				Labels:        issueOptionalStringSlice(cmd, "labels", labels),
				SpecMDPath:    issueCreateOptionalString(cmd, "spec-md", specMDPath),
				EstTokens:     issueCreateOptionalInt(cmd, "est-tokens", estTokens),
				DueDate:       issueCreateOptionalString(cmd, "due-date", dueDate),
				OwnerRole:     issueCreateOptionalString(cmd, "owner-role", ownerRole),
				OwnerActor:    issueCreateOptionalString(cmd, "owner-actor", ownerActor),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "updated %s:\n%s\n", result.ID, issueUpdateSummary(result.Changes))
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New title")
	cmd.Flags().IntVar(&priority, "priority", 0, "New priority")
	cmd.Flags().StringSliceVar(&labels, "labels", nil, "Replacement issue labels")
	cmd.Flags().StringVar(&ownerRole, "owner-role", "", "New owner role")
	cmd.Flags().StringVar(&ownerActor, "owner-actor", "", "New owner actor")
	cmd.Flags().StringVar(&specMDPath, "spec-md", "", "New spec markdown path")
	cmd.Flags().IntVar(&estTokens, "est-tokens", 0, "New estimated tokens")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "New due date")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newIssueTransitionCommand() *cobra.Command {
	var (
		toState   string
		actorRole string
	)

	cmd := &cobra.Command{
		Use:   "transition <id>",
		Short: "Transition an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.Transition(issues.TransitionOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ID:            args[0],
				ToState:       toState,
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "transitioned %s: %s \u2192 %s\n", result.ID, result.FromState, result.ToState)
		},
	}

	cmd.Flags().StringVar(&toState, "to", "", "Target state")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newIssueListCommand() *cobra.Command {
	var (
		kind           string
		status         string
		statusCategory string
		epicID         string
		label          string
		priorityFilter string
		inboxOnly      bool
		hasConcerns    bool
		order          string
		formatName     string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			priorities, err := issues.ParsePriorityFilter(priorityFilter)
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			listOptions := issues.ListOptions{
				WorkspaceRoot:  workspaceRoot,
				MigrationsDir:  issueMigrationsDir(),
				Kind:           kind,
				Status:         status,
				StatusCategory: statusCategory,
				EpicID:         epicID,
				Label:          label,
				Priorities:     priorities,
				Order:          order,
			}
			if cmd.Flags().Changed("inbox") {
				listOptions.InInbox = &inboxOnly
			}
			if cmd.Flags().Changed("has-concerns") {
				listOptions.HasOpenConcerns = &hasConcerns
			}

			items, err := issues.List(listOptions)
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, items); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&kind, "kind", "", "Filter by kind")
	cmd.Flags().StringVar(&status, "status", "", "Filter by exact status")
	cmd.Flags().StringVar(&statusCategory, "status-category", "", "Filter by status category")
	cmd.Flags().StringVar(&epicID, "epic", "", "Filter by parent epic ID")
	cmd.Flags().StringVar(&label, "label", "", "Filter by label")
	cmd.Flags().StringVar(&priorityFilter, "priority", "", "Filter by priority list")
	cmd.Flags().BoolVar(&inboxOnly, "inbox", false, "Show only inbox items")
	cmd.Flags().BoolVar(&hasConcerns, "has-concerns", false, "Show only issues with open concerns")
	cmd.Flags().StringVar(&order, "order", "", "Ordering")
	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}

func issueMigrationsDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "migrations"
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func renderIssueError(writer io.Writer, err error) error {
	var missingRoleErr *policy.MissingRoleError
	if errors.As(err, &missingRoleErr) {
		return writeIssueLine(writer, "error: %s\n", missingRoleErr.Error())
	}

	var unknownRoleErr *policy.UnknownRoleError
	if errors.As(err, &unknownRoleErr) {
		if err := writeIssueLine(writer, "error: %s\n", unknownRoleErr.Error()); err != nil {
			return err
		}
		return writeIssueLine(writer, "  known roles: %s\n", joinCommaSeparated(unknownRoleErr.KnownRoles))
	}

	var permissionDeniedErr *policy.PermissionDeniedError
	if errors.As(err, &permissionDeniedErr) {
		if err := writeIssueLine(writer, "error: %s\n", permissionDeniedErr.Error()); err != nil {
			return err
		}
		return writeIssueLine(writer, "  check `fid policy show` for permitted role/action combinations\n")
	}

	handled, renderErr := renderJournalReplayError(func(format string, args ...any) error {
		return writeIssueLine(writer, format, args...)
	}, err, "read stopped - manual repair required")
	if handled {
		return renderErr
	}

	var unsupportedFormatErr *fidformat.UnsupportedFormatError
	if errors.As(err, &unsupportedFormatErr) {
		return writeIssueLine(writer, "error: %s\n", unsupportedFormatErr.Error())
	}

	var unsupportedKindErr *issues.UnsupportedKindError
	if errors.As(err, &unsupportedKindErr) {
		return writeIssueLine(writer, "error: %s\n", unsupportedKindErr.Error())
	}

	var notFoundErr *issues.NotFoundError
	if errors.As(err, &notFoundErr) {
		return writeIssueLine(writer, "error: %s\n", notFoundErr.Error())
	}

	return writeIssueLine(writer, "error: %v\n", err)
}

func writeIssueLine(writer io.Writer, format string, args ...any) error {
	if _, err := fmt.Fprintf(writer, format, args...); err != nil {
		return fmt.Errorf("write issue output: %w", err)
	}
	return nil
}

func joinCommaSeparated(values []string) string {
	if len(values) == 0 {
		return ""
	}
	result := values[0]
	for _, value := range values[1:] {
		result += ", " + value
	}
	return result
}

func issueCreateParentID(kind, epicID, parentID string) (*string, error) {
	trimmedEpicID := strings.TrimSpace(epicID)
	trimmedParentID := strings.TrimSpace(parentID)

	if trimmedEpicID != "" && trimmedParentID != "" {
		return nil, fmt.Errorf("create issue: --epic and --parent are mutually exclusive")
	}

	switch strings.TrimSpace(kind) {
	case "task":
		if trimmedParentID != "" {
			return nil, fmt.Errorf("create issue: --parent is only valid for subtasks")
		}
		if trimmedEpicID == "" {
			return nil, nil
		}
		return &trimmedEpicID, nil
	case "subtask":
		if trimmedEpicID != "" {
			return nil, fmt.Errorf("create issue: --epic is only valid for tasks")
		}
		if trimmedParentID == "" {
			return nil, nil
		}
		return &trimmedParentID, nil
	default:
		if trimmedEpicID != "" {
			return nil, fmt.Errorf("create issue: --epic is only valid for tasks")
		}
		if trimmedParentID != "" {
			return nil, fmt.Errorf("create issue: --parent is only valid for subtasks")
		}
		return nil, nil
	}
}

func issueCreateOptionalInt(cmd *cobra.Command, flagName string, value int) *int {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	result := value
	return &result
}

func issueCreateOptionalString(cmd *cobra.Command, flagName, value string) *string {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	return &value
}

func issueClaimTTLSeconds(cmd *cobra.Command, ttl string) (*int, error) {
	if !cmd.Flags().Changed("ttl") {
		return nil, nil
	}

	duration, err := time.ParseDuration(ttl)
	if err != nil {
		return nil, fmt.Errorf("parse --ttl: %w", err)
	}
	if duration <= 0 {
		return nil, fmt.Errorf("--ttl must be greater than 0s")
	}
	if duration%time.Second != 0 {
		return nil, fmt.Errorf("--ttl must resolve to whole seconds")
	}

	ttlSeconds := int(duration / time.Second)
	return &ttlSeconds, nil
}

func issueSplitCommaSeparated(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	return result
}

func issueOptionalStringSlice(cmd *cobra.Command, flagName string, values []string) *[]string {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	result := slices.Clone(values)
	return &result
}

func issueUpdateSummary(changes []issues.FieldChange) string {
	if len(changes) == 0 {
		return "no changes"
	}

	parts := make([]string, 0, len(changes))
	for _, change := range changes {
		parts = append(parts, fmt.Sprintf("%s %s -> %s", issueUpdateFieldLabel(change.Field), issueUpdateValue(change.From), issueUpdateValue(change.To)))
	}
	return strings.Join(parts, "\n")
}

func issueUpdateFieldLabel(field string) string {
	switch field {
	case "spec_md_path":
		return "spec-md"
	case "est_tokens":
		return "est-tokens"
	case "due_date":
		return "due-date"
	case "owner_role":
		return "owner-role"
	case "owner_actor":
		return "owner-actor"
	default:
		return field
	}
}

func issueUpdateValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return "(none)"
	case string:
		return fmt.Sprintf("%q", typed)
	case []string:
		return "[" + strings.Join(typed, ", ") + "]"
	default:
		return fmt.Sprintf("%v", typed)
	}
}
