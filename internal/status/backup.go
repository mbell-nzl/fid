package status

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/backup"
)

// Backup is the backup's health as fid status shows it. Its human text
// names failure kinds, never kernel nouns.
type Backup struct {
	State               string         `json:"state"`
	Reasons             []string       `json:"reasons"`
	Summary             string         `json:"summary"`
	NextAction          *string        `json:"next_action"`
	LastBackupAt        *string        `json:"last_backup_at"`
	LedgerStartedAt     string         `json:"ledger_started_at"`
	AgeSeconds          *int64         `json:"age_seconds"`
	Latest              *backup.Record `json:"latest"`
	ConsecutiveFailures int            `json:"consecutive_failures"`
	TargetSet           *bool          `json:"target_set"`
	UserConfigPath      *string        `json:"user_config_path"`
	UserConfigError     *string        `json:"user_config_error"`
	RecordPath          string         `json:"record_path"`
	RecordError         *string        `json:"record_error"`
}

func describeBackup(health backup.Health, now time.Time) Backup {
	described := Backup{State: string(health.State), Reasons: []string{}, Latest: health.Latest,
		ConsecutiveFailures: health.ConsecutiveFailures, TargetSet: health.TargetSet, RecordPath: health.RecordPath,
		UserConfigError: errorText(health.UserConfigError), RecordError: errorText(health.RecordError)}
	for _, reason := range health.Reasons {
		described.Reasons = append(described.Reasons, string(reason))
	}
	if health.LastBackup != nil {
		at := health.LastBackup.UTC().Format(time.RFC3339Nano)
		described.LastBackupAt = &at
	}
	if !health.Started.IsZero() {
		described.LedgerStartedAt = health.Started.UTC().Format(time.RFC3339Nano)
	}
	if health.Age != nil {
		seconds := int64(*health.Age / time.Second)
		described.AgeSeconds = &seconds
	}
	if health.UserConfigPath != "" {
		path := health.UserConfigPath
		described.UserConfigPath = &path
	}
	described.Summary = backupSummary(health, now)
	described.NextAction = backupNextAction(health)
	return described
}

// backupSummary leads with the backup's age once it's over a day, and
// otherwise with what's wrong, if anything.
func backupSummary(health backup.Health, now time.Time) string {
	has := func(reason backup.HealthReason) bool { return slices.Contains(health.Reasons, reason) }
	problems := []string{}
	if has(backup.ReasonLastBackupFailed) && health.Latest != nil && health.Latest.Failure != nil {
		if health.ConsecutiveFailures > 1 {
			problems = append(problems, fmt.Sprintf("last %d backups failed (%s)", health.ConsecutiveFailures, health.Latest.Failure.Kind))
		} else {
			problems = append(problems, fmt.Sprintf("last backup failed (%s)", health.Latest.Failure.Kind))
		}
	}
	if has(backup.ReasonNoTarget) {
		problems = append(problems, "no backup target is set")
	}
	if has(backup.ReasonUserConfigUnreadable) {
		if health.UserConfigPath != "" {
			problems = append(problems, "can't read user config "+health.UserConfigPath)
		} else {
			problems = append(problems, "can't find user config: "+innermost(health.UserConfigError))
		}
	}
	if has(backup.ReasonRecordUnreadable) {
		problems = append(problems, "can't read the backup record "+health.RecordPath)
	}
	parts := []string{}
	switch {
	case has(backup.ReasonNoRecentBackup) && health.LastBackup != nil:
		parts = append(parts, fmt.Sprintf("no backup for %s: last backed up %s", duration(*health.Age), local(now, *health.LastBackup)))
	case has(backup.ReasonNoRecentBackup):
		parts = append(parts, fmt.Sprintf("no backup for %s: none since the ledger began %s", duration(*health.Age), local(now, health.Started)))
	}
	parts = append(parts, problems...)
	if len(parts) == len(problems) && health.Age != nil {
		recent := "not backed up yet"
		if health.LastBackup != nil {
			recent = "backed up " + ago(*health.Age)
			if len(problems) > 0 {
				recent = "last " + recent
			}
		}
		parts = append(parts, recent)
	}
	summary := strings.Join(parts, "; ")
	if summary == "" {
		return ""
	}
	return strings.ToUpper(summary[:1]) + summary[1:]
}

// backupNextAction is the fix to make first: config before backing up.
func backupNextAction(health backup.Health) *string {
	has := func(reason backup.HealthReason) bool { return slices.Contains(health.Reasons, reason) }
	next := ""
	switch {
	case has(backup.ReasonUserConfigUnreadable) && health.UserConfigPath != "":
		next = "fix " + health.UserConfigPath + ", then fid backup"
	case has(backup.ReasonNoTarget):
		next = "add [backup] target to " + health.UserConfigPath + ", then fid backup"
	case has(backup.ReasonRecordUnreadable):
		next = "fix or remove " + health.RecordPath + ", then fid backup"
	case len(health.Reasons) > 0:
		next = "fid backup"
	default:
		return nil
	}
	return &next
}

// duration reads as minutes, hours or days, rounded down.
func duration(d time.Duration) string {
	count, unit := int(d/time.Minute), "minute"
	switch {
	case d >= 48*time.Hour:
		count, unit = int(d/(24*time.Hour)), "day"
	case d >= time.Hour:
		count, unit = int(d/time.Hour), "hour"
	}
	if count != 1 {
		unit += "s"
	}
	return fmt.Sprintf("%d %s", count, unit)
}

func ago(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	return duration(d) + " ago"
}

// local is a time in now's time zone.
func local(now, at time.Time) string {
	return at.In(now.Location()).Format("2006-01-02 15:04")
}

func errorText(err error) *string {
	if err == nil {
		return nil
	}
	text := err.Error()
	return &text
}

func innermost(err error) string {
	if err == nil {
		return "unknown"
	}
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}
	return err.Error()
}
