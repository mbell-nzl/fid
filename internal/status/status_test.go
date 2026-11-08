package status

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/backup"
)

// nzdt is 13 hours ahead of UTC, so local times differ from UTC dates too.
var (
	nzdt   = time.FixedZone("NZDT", 13*60*60)
	testAt = time.Date(2026, 11, 8, 9, 0, 0, 0, time.UTC).In(nzdt)
)

// kernelWords are the kernel nouns and IDs that human output never names
// (RESET 1.7).
var kernelWords = regexp.MustCompile(`(?i)\b(envelopes?|activations?|obligations?|attempts?|attention items?|runs?|claims?)\b|ENV-|STP-|ATT-`)

func assertNoKernelWords(t *testing.T, text string) {
	t.Helper()
	if found := kernelWords.FindAllString(text, -1); len(found) > 0 {
		t.Fatalf("human output names kernel words %q:\n%s", found, text)
	}
}

func before(d time.Duration) *time.Time {
	at := testAt.Add(-d).UTC()
	return &at
}

func failedLine(kind backup.Kind, d time.Duration) *backup.Record {
	return &backup.Record{Schema: backup.RecordSchema, FinishedAt: before(d).Format(time.RFC3339Nano), Outcome: backup.OutcomeFailed,
		Failure: &backup.RecordedFailure{Kind: kind, Message: "an Attempt failed"}}
}

func TestDescribeBackupWording(t *testing.T) {
	t.Parallel()
	hour := time.Hour
	set, unset := true, false
	health := func(state backup.HealthState, last *time.Time, age time.Duration, reasons ...backup.HealthReason) backup.Health {
		return backup.Health{State: state, Reasons: reasons, LastBackup: last, Started: *before(1000 * hour), Age: &age, TargetSet: &set,
			UserConfigPath: "/cfg/fid/config.toml", RecordPath: "/ledger/backups.jsonl"}
	}
	stale, failed, noTarget := backup.ReasonNoRecentBackup, backup.ReasonLastBackupFailed, backup.ReasonNoTarget
	for _, test := range []struct {
		name    string
		health  backup.Health
		summary string
		next    string
	}{
		{"hours", health(backup.HealthCurrent, before(3*hour), 3*hour), "Backed up 3 hours ago", ""},
		{"one hour", health(backup.HealthCurrent, before(hour), hour), "Backed up 1 hour ago", ""},
		{"minutes", health(backup.HealthCurrent, before(30*time.Minute), 30*time.Minute), "Backed up 30 minutes ago", ""},
		{"one minute", health(backup.HealthCurrent, before(time.Minute), time.Minute), "Backed up 1 minute ago", ""},
		{"seconds", health(backup.HealthCurrent, before(20*time.Second), 20*time.Second), "Backed up just now", ""},
		{"in the future", health(backup.HealthCurrent, before(-2*hour), 0), "Backed up just now", ""},
		{"not yet", health(backup.HealthNotBackedUpYet, nil, 2*hour), "Not backed up yet", ""},
		{"over 24 hours", health(backup.HealthWarning, before(26*hour), 26*hour, stale),
			"No backup for 26 hours: last backed up 2026-11-07 20:00", "fid backup"},
		{"over 72 hours", health(backup.HealthNeedsYou, before(100*hour), 100*hour, stale),
			"No backup for 4 days: last backed up 2026-11-04 18:00", "fid backup"},
		{"never, over 24 hours", func() backup.Health {
			h := health(backup.HealthWarning, nil, 30*hour, stale)
			h.Started = *before(30 * hour)
			return h
		}(), "No backup for 30 hours: none since the ledger began 2026-11-07 16:00", "fid backup"},
		{"failed once", func() backup.Health {
			h := health(backup.HealthWarning, before(2*hour), 2*hour, failed)
			h.Latest, h.ConsecutiveFailures = failedLine(backup.KindPushRejected, hour), 1
			return h
		}(), "Last backup failed (push-rejected); last backed up 2 hours ago", "fid backup"},
		{"failed three times", func() backup.Health {
			h := health(backup.HealthWarning, before(8*hour), 8*hour, failed)
			h.Latest, h.ConsecutiveFailures = failedLine(backup.KindLockUnavailable, hour), 3
			return h
		}(), "Last 3 backups failed (lock-unavailable); last backed up 8 hours ago", "fid backup"},
		{"failed, never backed up", func() backup.Health {
			h := health(backup.HealthWarning, nil, 2*hour, failed)
			h.Latest, h.ConsecutiveFailures = failedLine(backup.KindPushRejected, hour), 1
			return h
		}(), "Last backup failed (push-rejected); not backed up yet", "fid backup"},
		{"old and failing", func() backup.Health {
			h := health(backup.HealthNeedsYou, before(80*hour), 80*hour, stale, failed)
			h.Latest, h.ConsecutiveFailures = failedLine(backup.KindTargetUnreachable, hour), 1
			return h
		}(), "No backup for 3 days: last backed up 2026-11-05 14:00; last backup failed (target-unreachable)", "fid backup"},
		{"no target", func() backup.Health {
			h := health(backup.HealthWarning, nil, hour, noTarget)
			h.TargetSet = &unset
			return h
		}(), "No backup target is set; not backed up yet", "add [backup] target to /cfg/fid/config.toml, then fid backup"},
		{"no target for over 72 hours", func() backup.Health {
			h := health(backup.HealthNeedsYou, nil, 73*hour, stale, noTarget)
			h.Started, h.TargetSet = *before(73 * hour), &unset
			return h
		}(), "No backup for 3 days: none since the ledger began 2026-11-05 21:00; no backup target is set",
			"add [backup] target to /cfg/fid/config.toml, then fid backup"},
		{"unreadable user config", func() backup.Health {
			h := health(backup.HealthWarning, before(hour), hour, backup.ReasonUserConfigUnreadable)
			h.TargetSet, h.UserConfigError = nil, errors.New("parse config: bad")
			return h
		}(), "Can't read user config /cfg/fid/config.toml; last backed up 1 hour ago", "fix /cfg/fid/config.toml, then fid backup"},
		{"no user config to find", func() backup.Health {
			h := health(backup.HealthWarning, before(hour), hour, backup.ReasonUserConfigUnreadable)
			h.TargetSet, h.UserConfigPath = nil, ""
			h.UserConfigError = &wrapped{"locate user config", errors.New("$HOME is not defined")}
			return h
		}(), "Can't find user config: $HOME is not defined; last backed up 1 hour ago", "fid backup"},
		{"unreadable record", func() backup.Health {
			h := health(backup.HealthWarning, nil, 0, backup.ReasonRecordUnreadable)
			h.Age, h.RecordError = nil, errors.New("is a directory")
			return h
		}(), "Can't read the backup record /ledger/backups.jsonl", "fix or remove /ledger/backups.jsonl, then fid backup"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			described := describeBackup(test.health, testAt)
			next := ""
			if described.NextAction != nil {
				next = *described.NextAction
			}
			if described.Summary != test.summary || next != test.next {
				t.Fatalf("backup = %q, next %q\nwant      %q, next %q", described.Summary, next, test.summary, test.next)
			}
			assertNoKernelWords(t, described.Summary+"\n"+next)
		})
	}
}

type wrapped struct {
	context string
	err     error
}

func (err *wrapped) Error() string { return err.context + ": " + err.err.Error() }
func (err *wrapped) Unwrap() error { return err.err }

func TestDescribeBackupFields(t *testing.T) {
	t.Parallel()
	age := 26 * time.Hour
	health := backup.Health{State: backup.HealthWarning, Reasons: []backup.HealthReason{backup.ReasonNoRecentBackup, backup.ReasonLastBackupFailed},
		LastBackup: before(26 * time.Hour), Started: *before(100 * time.Hour), Age: &age, Latest: failedLine(backup.KindPushRejected, time.Hour),
		ConsecutiveFailures: 1, UserConfigPath: "/cfg/fid/config.toml", RecordPath: "/ledger/backups.jsonl",
		UserConfigError: errors.New("bad config")}
	described := describeBackup(health, testAt)
	if described.State != "warning" || strings.Join(described.Reasons, ",") != "no-recent-backup,last-backup-failed" ||
		described.LastBackupAt == nil || *described.LastBackupAt != "2026-11-07T07:00:00Z" || described.LedgerStartedAt != "2026-11-04T05:00:00Z" ||
		described.AgeSeconds == nil || *described.AgeSeconds != 26*60*60 || described.Latest != health.Latest || described.ConsecutiveFailures != 1 ||
		described.TargetSet != nil || described.UserConfigPath == nil || *described.UserConfigPath != "/cfg/fid/config.toml" ||
		described.UserConfigError == nil || *described.UserConfigError != "bad config" || described.RecordPath != "/ledger/backups.jsonl" || described.RecordError != nil {
		t.Fatalf("described = %+v", described)
	}
	empty := describeBackup(backup.Health{State: backup.HealthWarning, Reasons: []backup.HealthReason{backup.ReasonRecordUnreadable}, RecordError: errors.New("gone")}, testAt)
	if empty.Reasons == nil || empty.LastBackupAt != nil || empty.AgeSeconds != nil || empty.UserConfigPath != nil || empty.RecordError == nil {
		t.Fatalf("described = %+v, want nulls for what's unknown", empty)
	}
}
