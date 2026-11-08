package backup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

var healthNow = time.Date(2026, 11, 8, 9, 0, 0, 0, time.UTC)

// healthLedger is a project whose ledger holds only a genesis Event dated
// started, and a user config directory.
type healthLedger struct {
	t          *testing.T
	root, home string
	configHome string
}

func newHealthLedger(t *testing.T, started time.Time) healthLedger {
	t.Helper()
	root := gittest.NewWorkspace(t)
	journal := filepath.Join(config.LedgerDir(root), config.JournalDirName)
	if err := os.MkdirAll(journal, 0o755); err != nil {
		t.Fatal(err)
	}
	genesis := `{"actor":{"actor_id":"fid-bootstrap","actor_type":"service","producer":"fid","provenance_ref":"fid:init","session_ref":null},` +
		`"aggregate":{"id":"project-0123456789abcdef0123456789abcdef","sequence":1,"type":"project"},"causation_id":null,` +
		`"command":"project.init","correlation_id":null,"effects":[],"event_id":"EVT-` + started.UTC().Format("20060102") + `-000001",` +
		`"event_schema":2,"event_type":"project.initialized","idempotency":null,"occurred_at":"` + started.UTC().Format(time.RFC3339Nano) + `","repository":null}` + "\n"
	if err := os.WriteFile(filepath.Join(journal, started.UTC().Format("2006-01-02")+".jsonl"), []byte(genesis), 0o644); err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	return healthLedger{t: t, root: root, home: filepath.Join(base, "home"), configHome: filepath.Join(base, "config")}
}

func (ledger healthLedger) environment() []string {
	return []string{"HOME=" + ledger.home, "XDG_CONFIG_HOME=" + ledger.configHome}
}

func (ledger healthLedger) configPath() string {
	return filepath.Join(ledger.configHome, "fid", config.ConfigFileName)
}

func (ledger healthLedger) setTarget() {
	ledger.t.Helper()
	if err := config.SetBackupTarget(ledger.configPath(), "/srv/backups.git"); err != nil {
		ledger.t.Fatal(err)
	}
}

func (ledger healthLedger) writeConfig(contents string) {
	ledger.t.Helper()
	if err := os.MkdirAll(filepath.Dir(ledger.configPath()), 0o700); err != nil {
		ledger.t.Fatal(err)
	}
	if err := os.WriteFile(ledger.configPath(), []byte(contents), 0o600); err != nil {
		ledger.t.Fatal(err)
	}
}

func (ledger healthLedger) recordPath() string {
	return filepath.Join(config.LedgerDir(ledger.root), RecordFileName)
}

func (ledger healthLedger) writeRecord(lines ...string) {
	ledger.t.Helper()
	if err := os.WriteFile(ledger.recordPath(), []byte(strings.Join(lines, "")), 0o644); err != nil {
		ledger.t.Fatal(err)
	}
}

// healthLine is one record line, finished ago before healthNow.
func healthLine(t *testing.T, outcome Outcome, ago time.Duration, kind Kind) string {
	t.Helper()
	finished := stamp(healthNow.Add(-ago))
	record := Record{Schema: RecordSchema, StartedAt: finished, FinishedAt: finished, Outcome: outcome}
	if outcome == OutcomeFailed {
		record.Failure = &RecordedFailure{Kind: kind, Message: "it failed"}
	} else {
		commit := strings.Repeat("a", 40)
		record.Commit = &commit
	}
	line, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	return string(line) + "\n"
}

func TestReadHealthAppliesTheBackupAgeRules(t *testing.T) {
	t.Parallel()
	hours := func(n int) time.Duration { return time.Duration(n) * time.Hour }
	type want struct {
		state    HealthState
		reasons  []HealthReason
		age      time.Duration
		last     time.Duration // ago; -1 for none
		failures int
		latest   Outcome
	}
	for _, test := range []struct {
		name    string
		started time.Duration // ago
		target  bool
		lines   func(*testing.T) []string
		want    want
	}{
		{"pushed 3 hours ago is current", hours(100), true,
			func(t *testing.T) []string { return []string{healthLine(t, OutcomePushed, hours(3), "")} },
			want{HealthCurrent, nil, hours(3), hours(3), 0, OutcomePushed}},
		{"unchanged counts as a backup", hours(100), true,
			func(t *testing.T) []string {
				return []string{healthLine(t, OutcomePushed, hours(30), ""), healthLine(t, OutcomeUnchanged, hours(2), "")}
			},
			want{HealthCurrent, nil, hours(2), hours(2), 0, OutcomeUnchanged}},
		{"exactly 24 hours is still current", hours(100), true,
			func(t *testing.T) []string { return []string{healthLine(t, OutcomePushed, hours(24), "")} },
			want{HealthCurrent, nil, hours(24), hours(24), 0, OutcomePushed}},
		{"over 24 hours warns", hours(100), true,
			func(t *testing.T) []string { return []string{healthLine(t, OutcomePushed, hours(25), "")} },
			want{HealthWarning, []HealthReason{ReasonNoRecentBackup}, hours(25), hours(25), 0, OutcomePushed}},
		{"exactly 72 hours still warns", hours(100), true,
			func(t *testing.T) []string { return []string{healthLine(t, OutcomePushed, hours(72), "")} },
			want{HealthWarning, []HealthReason{ReasonNoRecentBackup}, hours(72), hours(72), 0, OutcomePushed}},
		{"over 72 hours needs you", hours(100), true,
			func(t *testing.T) []string { return []string{healthLine(t, OutcomePushed, hours(73), "")} },
			want{HealthNeedsYou, []HealthReason{ReasonNoRecentBackup}, hours(73), hours(73), 0, OutcomePushed}},
		{"a failed latest backup warns however recent the last good one", hours(100), true,
			func(t *testing.T) []string {
				return []string{healthLine(t, OutcomePushed, hours(2), ""), healthLine(t, OutcomeFailed, hours(1), KindPushRejected)}
			},
			want{HealthWarning, []HealthReason{ReasonLastBackupFailed}, hours(2), hours(2), 1, OutcomeFailed}},
		{"consecutive failures are counted", hours(100), true,
			func(t *testing.T) []string {
				return []string{healthLine(t, OutcomeFailed, hours(9), KindPushRejected), healthLine(t, OutcomePushed, hours(8), ""),
					healthLine(t, OutcomeFailed, hours(3), KindLockUnavailable), healthLine(t, OutcomeFailed, hours(2), KindLockUnavailable),
					healthLine(t, OutcomeFailed, hours(1), KindLockUnavailable)}
			},
			want{HealthWarning, []HealthReason{ReasonLastBackupFailed}, hours(8), hours(8), 3, OutcomeFailed}},
		{"a success after failures clears them", hours(100), true,
			func(t *testing.T) []string {
				return []string{healthLine(t, OutcomeFailed, hours(3), KindLockUnavailable), healthLine(t, OutcomePushed, hours(1), "")}
			},
			want{HealthCurrent, nil, hours(1), hours(1), 0, OutcomePushed}},
		{"old and failing needs you", hours(100), true,
			func(t *testing.T) []string {
				return []string{healthLine(t, OutcomePushed, hours(80), ""), healthLine(t, OutcomeFailed, hours(1), KindTargetUnreachable)}
			},
			want{HealthNeedsYou, []HealthReason{ReasonNoRecentBackup, ReasonLastBackupFailed}, hours(80), hours(80), 1, OutcomeFailed}},
		{"a new ledger is not backed up yet", hours(2), true, nil,
			want{HealthNotBackedUpYet, nil, hours(2), -1, 0, ""}},
		{"a ledger never backed up warns after 24 hours", hours(25), true, nil,
			want{HealthWarning, []HealthReason{ReasonNoRecentBackup}, hours(25), -1, 0, ""}},
		{"a restored ledger with no record needs you", hours(24 * 40), true, nil,
			want{HealthNeedsYou, []HealthReason{ReasonNoRecentBackup}, hours(24 * 40), -1, 0, ""}},
		{"a new ledger whose first backup failed warns", hours(2), true,
			func(t *testing.T) []string { return []string{healthLine(t, OutcomeFailed, hours(1), KindPushRejected)} },
			want{HealthWarning, []HealthReason{ReasonLastBackupFailed}, hours(2), -1, 1, OutcomeFailed}},
		{"no target warns", hours(100), false,
			func(t *testing.T) []string { return []string{healthLine(t, OutcomePushed, hours(1), "")} },
			want{HealthWarning, []HealthReason{ReasonNoTarget}, hours(1), hours(1), 0, OutcomePushed}},
		{"no target on a new ledger warns", hours(1), false, nil,
			want{HealthWarning, []HealthReason{ReasonNoTarget}, hours(1), -1, 0, ""}},
		{"no target for over 72 hours needs you", hours(73), false, nil,
			want{HealthNeedsYou, []HealthReason{ReasonNoRecentBackup, ReasonNoTarget}, hours(73), -1, 0, ""}},
		{"a backup in the future counts as now", hours(100), true,
			func(t *testing.T) []string { return []string{healthLine(t, OutcomePushed, -hours(2), "")} },
			want{HealthCurrent, nil, 0, -hours(2), 0, OutcomePushed}},
		{"a ledger begun in the future counts as now", -hours(5), true, nil,
			want{HealthNotBackedUpYet, nil, 0, -1, 0, ""}},
		{"a torn last line is skipped", hours(100), true,
			func(t *testing.T) []string {
				return []string{healthLine(t, OutcomePushed, hours(2), ""), `{"schema":"fid.backup-record.v1","started_at":"2026-11-08T08:`}
			},
			want{HealthCurrent, nil, hours(2), hours(2), 0, OutcomePushed}},
		{"a last line that isn't a record is skipped", hours(100), true,
			func(t *testing.T) []string {
				return []string{healthLine(t, OutcomePushed, hours(2), ""), "{\"schema\":\"something-else\"}\n", "[]\n"}
			},
			want{HealthCurrent, nil, hours(2), hours(2), 0, OutcomePushed}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ledger := newHealthLedger(t, healthNow.Add(-test.started))
			if test.target {
				ledger.setTarget()
			}
			if test.lines != nil {
				ledger.writeRecord(test.lines(t)...)
			}
			health, err := ReadHealth(HealthOptions{WorkspaceRoot: ledger.root, Environment: ledger.environment(), Now: healthNow})
			if err != nil {
				t.Fatalf("ReadHealth() error = %v", err)
			}
			if health.State != test.want.state || !slices.Equal(health.Reasons, test.want.reasons) {
				t.Fatalf("health = %s %v, want %s %v", health.State, health.Reasons, test.want.state, test.want.reasons)
			}
			if health.Age == nil || *health.Age != test.want.age {
				t.Fatalf("age = %v, want %v", health.Age, test.want.age)
			}
			if test.want.last == -1 && health.LastBackup != nil || test.want.last != -1 && (health.LastBackup == nil || !health.LastBackup.Equal(healthNow.Add(-test.want.last))) {
				t.Fatalf("last backup = %v, want %v ago", health.LastBackup, test.want.last)
			}
			if health.ConsecutiveFailures != test.want.failures {
				t.Fatalf("consecutive failures = %d, want %d", health.ConsecutiveFailures, test.want.failures)
			}
			if test.want.latest == "" && health.Latest != nil || test.want.latest != "" && (health.Latest == nil || health.Latest.Outcome != test.want.latest) {
				t.Fatalf("latest = %+v, want outcome %q", health.Latest, test.want.latest)
			}
			if !health.Started.Equal(healthNow.Add(-test.started)) {
				t.Fatalf("started = %v, want %v ago", health.Started, test.started)
			}
			if health.TargetSet == nil || *health.TargetSet != test.target || health.UserConfigPath != ledger.configPath() || health.RecordPath != ledger.recordPath() {
				t.Fatalf("health = %+v, want target set %v at %s, record %s", health, test.target, ledger.configPath(), ledger.recordPath())
			}
		})
	}
}

// Trouble reading user config or the record is a warning, never a failure.
func TestReadHealthWarnsOnUnreadableConfigAndRecord(t *testing.T) {
	t.Parallel()
	recent := func(t *testing.T, ledger healthLedger) {
		ledger.writeRecord(healthLine(t, OutcomePushed, time.Hour, ""))
	}
	for _, test := range []struct {
		name        string
		prepare     func(*testing.T, healthLedger)
		environment func(healthLedger) []string
		reasons     []HealthReason
		configPath  bool
		targetSet   bool
		age         bool
	}{
		{"malformed user config", func(t *testing.T, ledger healthLedger) {
			recent(t, ledger)
			ledger.writeConfig("[backup]\ntarget = unquoted\n")
		}, healthLedger.environment, []HealthReason{ReasonUserConfigUnreadable}, true, false, true},
		{"user config that isn't a file", func(t *testing.T, ledger healthLedger) {
			recent(t, ledger)
			if err := os.MkdirAll(ledger.configPath(), 0o700); err != nil {
				t.Fatal(err)
			}
		}, healthLedger.environment, []HealthReason{ReasonUserConfigUnreadable}, true, false, true},
		{"no HOME to find user config by", recent, func(healthLedger) []string { return nil },
			[]HealthReason{ReasonUserConfigUnreadable}, false, false, true},
		{"a record that isn't a file", func(t *testing.T, ledger healthLedger) {
			ledger.setTarget()
			if err := os.Mkdir(ledger.recordPath(), 0o755); err != nil {
				t.Fatal(err)
			}
		}, healthLedger.environment, []HealthReason{ReasonRecordUnreadable}, true, true, false},
		{"both, listed in order", func(t *testing.T, ledger healthLedger) {
			ledger.writeConfig("[elsewhere]\n")
			if err := os.Mkdir(ledger.recordPath(), 0o755); err != nil {
				t.Fatal(err)
			}
		}, healthLedger.environment, []HealthReason{ReasonUserConfigUnreadable, ReasonRecordUnreadable}, true, false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ledger := newHealthLedger(t, healthNow.Add(-100*time.Hour))
			test.prepare(t, ledger)
			health, err := ReadHealth(HealthOptions{WorkspaceRoot: ledger.root, Environment: test.environment(ledger), Now: healthNow})
			if err != nil {
				t.Fatalf("ReadHealth() error = %v", err)
			}
			if health.State != HealthWarning || !slices.Equal(health.Reasons, test.reasons) {
				t.Fatalf("health = %s %v, want warning %v", health.State, health.Reasons, test.reasons)
			}
			if (health.UserConfigPath != "") != test.configPath || (health.TargetSet != nil) != test.targetSet || (health.Age != nil) != test.age {
				t.Fatalf("health = %+v, want config path %v, target known %v, age known %v", health, test.configPath, test.targetSet, test.age)
			}
			problem := health.UserConfigError
			if slices.Contains(test.reasons, ReasonRecordUnreadable) {
				problem = health.RecordError
			}
			if problem == nil {
				t.Fatalf("health = %+v, want the problem's error", health)
			}
		})
	}
}

func TestReadHealthFailsWithoutAJournal(t *testing.T) {
	t.Parallel()
	ledger := newHealthLedger(t, healthNow)
	if err := os.RemoveAll(filepath.Join(config.LedgerDir(ledger.root), config.JournalDirName)); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadHealth(HealthOptions{WorkspaceRoot: ledger.root, Environment: ledger.environment(), Now: healthNow}); err == nil {
		t.Fatal("ReadHealth() succeeded without a journal")
	}
}
