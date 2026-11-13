// Package status assembles fid status: what needs a person, the work in
// progress and to do, policy health, and how current the ledger's backup is.
// Text and JSON are both rendered from one Result.
package status

import (
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/backup"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/next"
)

// NeedsYouBackup is the kind of item that says the backup needs a person.
// Later tasks add kinds for approvals, epic plans, untriaged work and tasks
// waiting for a person.
const NeedsYouBackup = "backup"

// CreateTaskCommand is the command the empty project's hint names.
const CreateTaskCommand = "fid issue create --kind task --title <title> --actor human:<name> --idempotency-key <key>"

type Options struct {
	WorkspaceRoot string
	MigrationsDir string
	// Environment locates user config, as fid backup's does.
	Environment []string
	// Now is when status is judged, in the time zone it shows times in.
	Now time.Time
}

// Result is everything fid status shows, uncapped.
type Result struct {
	GeneratedAt time.Time
	Backup      Backup
	NeedsYou    []NeedsYouItem
	Warnings    []Warning
	InProgress  issues.OverviewSection
	Ready       issues.OverviewSection
	Next        *string
	entries     map[string]next.Entry
	// Hint is set only when the project has no epics or tasks at all.
	Hint *Hint
}

// NeedsYouItem is something waiting on a person, with what to do next.
// Subject names the issue it's about, when it is about one.
type NeedsYouItem struct {
	Kind       string   `json:"kind"`
	Summary    string   `json:"summary"`
	NextAction *string  `json:"next_action"`
	Subject    *Subject `json:"subject"`
}

type Subject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Hint struct {
	Summary    string `json:"summary"`
	NextAction string `json:"next_action"`
}

// Build reads the project's work and its backup health.
func Build(options Options) (Result, error) {
	return build(options, nil)
}

func build(options Options, afterOverview func()) (Result, error) {
	overview, planned, err := readWork(options, afterOverview)
	if err != nil {
		return Result{}, fmt.Errorf("read status: %w", err)
	}
	health, err := backup.ReadHealth(backup.HealthOptions{WorkspaceRoot: options.WorkspaceRoot, Environment: options.Environment, Now: options.Now})
	if err != nil {
		return Result{}, fmt.Errorf("read status: %w", err)
	}
	result := Result{GeneratedAt: options.Now, Backup: describeBackup(health, options.Now), NeedsYou: []NeedsYouItem{},
		Warnings: readPolicyWarnings(options.WorkspaceRoot), InProgress: overview.InProgress, Ready: overview.ToDo, entries: map[string]next.Entry{}}
	for _, entry := range planned.Queue {
		result.entries[entry.ID] = entry
	}
	if planned.Kind == next.KindNew && planned.Task != nil {
		result.Next = &planned.Task.ID
	}
	result.addWaiting(planned)
	result.checkWorktrees(options.WorkspaceRoot, planned)
	if health.State == backup.HealthNeedsYou {
		result.NeedsYou = append(result.NeedsYou, NeedsYouItem{Kind: NeedsYouBackup, Summary: result.Backup.Summary, NextAction: result.Backup.NextAction})
	}
	if overview.Empty {
		result.Hint = &Hint{Summary: "No tasks yet", NextAction: CreateTaskCommand}
	}
	return result, nil
}

// Fields is the result as the project.status JSON result.
func (result Result) Fields() map[string]any {
	return map[string]any{
		"generated_at": result.GeneratedAt.UTC().Format(time.RFC3339),
		"backup":       result.Backup,
		"needs_you":    result.NeedsYou,
		"warnings":     result.Warnings,
		"in_progress": mapSection(result.InProgress, func(task issues.OverviewTask) workingTask {
			return workingTask{OverviewTask: task, Working: result.entries[task.ID].Working, Kept: result.entries[task.ID].Kept}
		}),
		"ready": result.readyFields(),
		"hint":  result.Hint,
	}
}
