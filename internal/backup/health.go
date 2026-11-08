package backup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

const (
	// WarnAfter is how long without a backup before fid status warns.
	WarnAfter = 24 * time.Hour
	// NeedsYouAfter is how long without a backup before it needs a person.
	NeedsYouAfter = 72 * time.Hour
)

// HealthState is how far the ledger's backups need attention.
type HealthState string

const (
	HealthCurrent        HealthState = "current"
	HealthNotBackedUpYet HealthState = "not-backed-up-yet"
	HealthWarning        HealthState = "warning"
	HealthNeedsYou       HealthState = "needs-you"
)

// HealthReason is one thing wrong with the ledger's backups. Health lists
// them in this order.
type HealthReason string

const (
	ReasonNoRecentBackup       HealthReason = "no-recent-backup"
	ReasonLastBackupFailed     HealthReason = "last-backup-failed"
	ReasonNoTarget             HealthReason = "no-target"
	ReasonUserConfigUnreadable HealthReason = "user-config-unreadable"
	ReasonRecordUnreadable     HealthReason = "record-unreadable"
)

var reasonOrder = []HealthReason{ReasonNoRecentBackup, ReasonLastBackupFailed, ReasonNoTarget, ReasonUserConfigUnreadable, ReasonRecordUnreadable}

// HealthOptions names the project, the environment that locates user
// config, and the time to judge the backups' age at.
type HealthOptions struct {
	WorkspaceRoot string
	Environment   []string
	Now           time.Time
}

// Health is what the record and user config say about the ledger's backups.
type Health struct {
	State   HealthState
	Reasons []HealthReason
	// LastBackup is the latest pushed or unchanged backup's finish; nil
	// when the record has none.
	LastBackup *time.Time
	// Started is the genesis Event's time, which a ledger never backed up
	// counts its age from.
	Started time.Time
	// Age is how long since LastBackup, or since Started; nil when the
	// record can't be read. A time in the future counts as now.
	Age    *time.Duration
	Latest *Record
	// ConsecutiveFailures counts the failed lines that end the record.
	ConsecutiveFailures int
	// TargetSet is nil when user config can't be read.
	TargetSet       *bool
	UserConfigPath  string
	UserConfigError error
	RecordPath      string
	RecordError     error
}

// ReadHealth judges the ledger's backups from its record and user config,
// without changing either. Trouble with them is a reason in the result; only
// a journal whose first Event can't be read is an error.
func ReadHealth(options HealthOptions) (Health, error) {
	genesis, err := genesisEvent(options.WorkspaceRoot)
	if err != nil {
		return Health{}, fmt.Errorf("read backup health: %w", err)
	}
	health := Health{Started: genesis.OccurredAt, RecordPath: filepath.Join(config.LedgerDir(options.WorkspaceRoot), RecordFileName)}
	health.readRecord(options.Now)
	health.readTarget(lookup(options.Environment))
	slices.SortFunc(health.Reasons, func(left, right HealthReason) int {
		return slices.Index(reasonOrder, left) - slices.Index(reasonOrder, right)
	})
	switch {
	case health.Age != nil && *health.Age > NeedsYouAfter:
		health.State = HealthNeedsYou
	case len(health.Reasons) > 0:
		health.State = HealthWarning
	case health.LastBackup == nil:
		health.State = HealthNotBackedUpYet
	default:
		health.State = HealthCurrent
	}
	return health, nil
}

func (health *Health) readRecord(now time.Time) {
	records, err := readRecords(health.RecordPath)
	if err != nil {
		health.RecordError = err
		health.Reasons = append(health.Reasons, ReasonRecordUnreadable)
		return
	}
	from := health.Started
	for index := range records {
		record := &records[index]
		if record.Outcome == OutcomeFailed {
			health.ConsecutiveFailures++
			continue
		}
		finished, _ := time.Parse(time.RFC3339Nano, record.FinishedAt)
		health.LastBackup, from = &finished, finished
		health.ConsecutiveFailures = 0
	}
	if len(records) > 0 {
		health.Latest = &records[len(records)-1]
	}
	age := max(now.Sub(from), 0)
	health.Age = &age
	if age > WarnAfter {
		health.Reasons = append(health.Reasons, ReasonNoRecentBackup)
	}
	if health.ConsecutiveFailures > 0 {
		health.Reasons = append(health.Reasons, ReasonLastBackupFailed)
	}
}

func (health *Health) readTarget(getenv func(string) string) {
	path, err := config.UserConfigPathFrom(getenv)
	if err == nil {
		health.UserConfigPath = path
		var user config.UserConfig
		if user, err = config.ReadUserConfig(path); err == nil {
			set := strings.TrimSpace(user.BackupTarget) != ""
			health.TargetSet = &set
			if !set {
				health.Reasons = append(health.Reasons, ReasonNoTarget)
			}
			return
		}
	}
	health.UserConfigError = err
	health.Reasons = append(health.Reasons, ReasonUserConfigUnreadable)
}

// readRecords returns the record's lines in order, none when it doesn't
// exist. It skips a line that isn't a whole record: appends repair a torn
// last line, and a reader can meet one between backups.
func readRecords(path string) ([]Record, error) {
	payload, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read the backup record: %w", err)
	}
	records := []Record{}
	for line := range bytes.SplitSeq(payload, []byte{'\n'}) {
		var record Record
		if json.Unmarshal(line, &record) != nil || !wholeRecord(record) {
			continue
		}
		records = append(records, record)
	}
	return records, nil
}

func wholeRecord(record Record) bool {
	if _, err := time.Parse(time.RFC3339Nano, record.FinishedAt); err != nil || record.Schema != RecordSchema {
		return false
	}
	switch record.Outcome {
	case OutcomePushed, OutcomeUnchanged:
		return record.Failure == nil
	case OutcomeFailed:
		return record.Failure != nil
	}
	return false
}
