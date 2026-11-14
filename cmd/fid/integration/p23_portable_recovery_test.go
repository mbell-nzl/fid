//go:build integration

package integration

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

var p23PortableTables = []string{
	"review_records", "review_findings", "finding_dispositions", "concerns", "deviations",
	"decision_refs", "human_outcomes", "envelope_protected_decisions", "requirement_refs",
	"requirement_issue_links", "envelope_obligations", "evidence_records", "budget_samples", "activity_log",
}

var p23PortableEvents = []string{
	"review.recorded", "finding.addressed", "concern.recorded", "concern.disposed",
	"deviation.recorded", "decision.referenced", "human-outcome.recorded",
	"requirement.referenced", "requirement.linked",
}

func TestPortableRoundTripCoversEveryPhase23Table(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	root := t.TempDir()
	populated := map[string]bool{}
	fieldDigestFamilies := map[string]bool{}
	events := map[string]journal.Event{}
	var validationWorkspace string
	for _, fixture := range []struct {
		name string
		make func(*testing.T) string
	}{
		{"mixed-findings", func(t *testing.T) string { return p23MixedFindingFixture(t, binary, root).Project }},
		{"protected-decision", func(t *testing.T) string {
			checkpoint := p23ProtectedDecisionFixture(t, binary, root)
			step := checkpoint.Steps[1]
			stdout, stderr, err := closureRunCommand(step.CWD, binary, step.Args...)
			if err != nil || stderr != "" || stdout == "" {
				t.Fatalf("record Human Outcome: %v stdout=%q stderr=%q", err, stdout, stderr)
			}
			return checkpoint.Project
		}},
		{"review-exhaustion", func(t *testing.T) string { return p23ReviewExhaustionFixture(t, binary, root).Project }},
		{"references-and-judgment", func(t *testing.T) string { return p23PortableJudgmentFixture(t, binary, root) }},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			workspace := fixture.make(t)
			p23PortableCollectFieldDigestFamilies(t, workspace, fieldDigestFamilies)
			for table, rows := range p23PortableRoundTrip(t, binary, workspace) {
				if len(rows) > 0 {
					populated[table] = true
				}
			}
			for kind, event := range p23PortableJournalEvents(t, workspace) {
				events[kind] = event
			}
			validationWorkspace = workspace
			if fixture.name == "review-exhaustion" {
				p23PortableAssertSingleReviewCharge(t, workspace)
			}
		})
	}
	for _, table := range p23PortableTables {
		if !populated[table] {
			t.Errorf("Phase 23 table %s never populated", table)
		}
	}
	for _, family := range []string{"issue.", "claim.", "run.", "review."} {
		if !fieldDigestFamilies[family] {
			t.Errorf("no %s keyed row included in field-digest rebuild parity", family)
		}
	}
	for _, kind := range p23PortableEvents {
		if _, found := events[kind]; !found {
			t.Errorf("Phase 23 Event %s never exported", kind)
		}
	}
	p23PortableAssertCorruptEventsRejected(t, validationWorkspace, events)
	p23PortableAssertPriorSchemaRejected(t, binary, validationWorkspace)
}

func p23PortableCollectFieldDigestFamilies(t *testing.T, workspace string, families map[string]bool) {
	t.Helper()
	conn := p22FaultDB(t, workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	rows, err := conn.Query(`SELECT operation_scope, request_field_digests_json FROM idempotency_results`)
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	for rows.Next() {
		var scope, fields string
		p22FaultOK(t, rows.Scan(&scope, &fields))
		var named map[string]string
		if err := json.Unmarshal([]byte(fields), &named); err != nil {
			t.Errorf("%s invalid field digests: %v", scope, err)
		} else if len(named) < 2 || !strings.HasPrefix(named["--actor"], "sha256:") {
			t.Errorf("%s field digests = %q, want public input fields and --actor", scope, fields)
		}
		for _, family := range []string{"issue.", "claim.", "run.", "review."} {
			if strings.HasPrefix(scope, family) {
				families[family] = true
			}
		}
	}
	p22FaultOK(t, rows.Err())
}

func TestPortableAssertCorruptEventsRejected(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	root := t.TempDir()
	events := map[string]journal.Event{}
	workspaces := []string{
		p23MixedFindingFixture(t, binary, root).Project,
		p23ReviewExhaustionFixture(t, binary, root).Project,
		p23PortableJudgmentFixture(t, binary, root),
	}
	protected := p23ProtectedDecisionFixture(t, binary, root)
	step := protected.Steps[1]
	stdout, stderr, err := closureRunCommand(step.CWD, binary, step.Args...)
	if err != nil || stderr != "" || stdout == "" {
		t.Fatalf("record Human Outcome: %v stdout=%q stderr=%q", err, stdout, stderr)
	}
	workspaces = append(workspaces, protected.Project)
	for _, workspace := range workspaces {
		for kind, event := range p23PortableJournalEvents(t, workspace) {
			events[kind] = event
		}
	}
	p23PortableAssertCorruptEventsRejected(t, workspaces[2], events)
}

func p23PortableAssertSingleReviewCharge(t *testing.T, workspace string) {
	t.Helper()
	conn := p22FaultDB(t, workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var count, consumed int
	p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*), COALESCE(SUM(value),0) FROM budget_samples WHERE budget_name='review-cycles'`).Scan(&count, &consumed))
	if count != 1 || consumed != 1 {
		t.Fatalf("review-cycles charged %d rows consuming %d, want one Review charge", count, consumed)
	}
	var matched int
	p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*) FROM budget_samples b JOIN review_records r ON r.id=b.source_sample_id WHERE b.budget_name='review-cycles'`).Scan(&matched))
	if matched != 1 {
		t.Fatalf("review-cycles Review ID match = %d, want 1", matched)
	}
}

func p23PortableAssertCorruptEventsRejected(t *testing.T, workspace string, events map[string]journal.Event) {
	t.Helper()
	source := p22FaultDB(t, workspace)
	baseline := p23PortableAllRows(t, source)
	p22FaultOK(t, source.Close())
	seenConnections := map[*sql.DB]bool{}
	seenWorkspaces := map[string]bool{}
	for _, eventType := range append(append([]string{}, p23PortableEvents...), "unknown.phase23") {
		base := events[eventType]
		if eventType == "unknown.phase23" {
			base = events["requirement.linked"]
			base.EventType = eventType
		}
		for _, corruption := range []string{"unknown-effect-table", "missing-required-effect", "malformed-body"} {
			if eventType == "unknown.phase23" && corruption != "unknown-effect-table" {
				continue
			}
			t.Run(eventType+"/"+corruption, func(t *testing.T) {
				clone := gittest.NewWorkspace(t)
				copyParityFixture(t, config.LedgerDir(workspace), config.LedgerDir(clone))
				if seenWorkspaces[clone] {
					t.Fatal("corruption probe reused a workspace")
				}
				seenWorkspaces[clone] = true
				conn := p22FaultDB(t, clone)
				defer func() { p22FaultOK(t, conn.Close()) }()
				if seenConnections[conn] {
					t.Fatal("corruption probe reused a projection connection")
				}
				seenConnections[conn] = true
				before := p23PortableAllRows(t, conn)
				if !reflect.DeepEqual(baseline, before) {
					t.Fatal("corruption probe did not start from the clean projection")
				}
				encoded, err := json.Marshal(base)
				p22FaultOK(t, err)
				var damaged journal.Event
				p22FaultOK(t, json.Unmarshal(encoded, &damaged))
				switch corruption {
				case "unknown-effect-table":
					damaged.Effects[0].Table = "unknown_phase23_table"
				case "missing-required-effect":
					missing := 0
					if eventType == "requirement.referenced" {
						for index, effect := range damaged.Effects {
							if effect.Table == "requirement_refs" && effect.Op == "insert" {
								missing = index
								break
							}
						}
					}
					damaged.Effects = append(damaged.Effects[:missing], damaged.Effects[missing+1:]...)
				case "malformed-body":
					damaged.Effects[0].Values["id"] = []any{"invalid scalar"}
				}
				hash, err := journal.HashEvent(damaged)
				p22FaultOK(t, err)
				rejection := journal.ProjectEvent(conn, journal.ProjectionInput{Event: damaged, EventHash: hash, File: "corrupt.jsonl", Line: 1, Offset: 1})
				if rejection == nil || !strings.Contains(rejection.Error(), "current Event") {
					t.Fatalf("current Event validator failed to reject payload: %v", rejection)
				}
				after := p23PortableAllRows(t, conn)
				if !reflect.DeepEqual(before, after) {
					t.Fatal("rejected Event changed projection")
				}
				beforeBytes, err := json.Marshal(before)
				p22FaultOK(t, err)
				afterBytes, err := json.Marshal(after)
				p22FaultOK(t, err)
				t.Logf("workspace=%s validator=%q full projection SHA-256 before=%x after=%x tables=%d", clone, rejection, sha256.Sum256(beforeBytes), sha256.Sum256(afterBytes), len(before))
			})
		}
	}
}

func p23PortableAllRows(t *testing.T, conn *sql.DB) map[string][]string {
	t.Helper()
	names, err := conn.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	p22FaultOK(t, err)
	var tables []string
	for names.Next() {
		var table string
		p22FaultOK(t, names.Scan(&table))
		tables = append(tables, table)
	}
	p22FaultOK(t, names.Err())
	p22FaultOK(t, names.Close())
	result := make(map[string][]string, len(tables))
	for _, table := range tables {
		result[table] = p23PortableTableRows(t, conn, table)
	}
	return result
}

func p23PortableAssertPriorSchemaRejected(t *testing.T, binary, workspace string) {
	t.Helper()
	archive := filepath.Join(t.TempDir(), "current.zip")
	p21Run(t, binary, workspace, "--json", "export", "create", archive)
	entries := p22ArchiveEntries(t, archive)
	var manifest map[string]any
	p22FaultOK(t, json.Unmarshal(entries["manifest.json"], &manifest))
	if manifest["schema"] != "fid.export-manifest.v2" || manifest["projection_schema"] != float64(23) {
		t.Fatalf("current manifest profile = %#v", manifest)
	}
	manifest["projection_schema"] = 17
	entries["manifest.json"], _ = json.Marshal(manifest)
	priorArchive := filepath.Join(t.TempDir(), "prior.zip")
	p22WritePortableArchive(t, priorArchive, entries)
	target := filepath.Join(t.TempDir(), "imported")
	closureRun(t, filepath.Dir(target), "git", "clone", "-q", workspace, target)
	stdout, stderr, err := closureRunCommand(target, binary, "--json", "export", "import", priorArchive, "--repository", target, "--idempotency-key", "prior-schema")
	if err == nil || stdout != "" || !strings.Contains(stderr, `"code":"archive-invalid"`) || !strings.Contains(stderr, `"code":"schema-unsupported"`) {
		t.Fatalf("prior schema import = err:%v stdout:%q stderr:%q", err, stdout, stderr)
	}
	if _, err := os.Stat(config.LedgerDir(target)); !os.IsNotExist(err) {
		t.Fatalf("rejected import published .fid: %v", err)
	}
}

func p23PortableJudgmentFixture(t *testing.T, binary, root string) string {
	t.Helper()
	workspace := filepath.Join(root, "judgment")
	if err := os.MkdirAll(workspace, 0o700); err != nil {
		t.Fatal(err)
	}
	actions := []string{"claim.acquire", "concern.dispose", "concern.record", "deviation.record", "envelope.activate", "envelope.create", "gate.register", "run.start"}
	a := p22NewFIDRecoveryAdapterAt(t, binary, workspace, "agent", 3600, []string{}, actions)
	a.task13ReplayCorpus = true
	for name, response := range map[string]p22ConformanceResponse{
		"requirement": p23RunInput(t, a, "requirement", "record-reference", a.actor, "portable-requirement", map[string]any{"id": "REQ-portable", "source": p23DecisionSource(a, true), "status": "active"}),
		"link":        p23RunInput(t, a, "requirement", "link", a.actor, "portable-link", map[string]any{"requirement_id": "REQ-portable", "issue_id": a.issueID, "relation": "verifies"}),
		"concern":     p23RunInput(t, a, "concern", "record", a.actor, "portable-concern", map[string]any{"id": "CON-portable", "subject": map[string]any{"type": "issue", "id": a.issueID}, "scope": map[string]any{"paths": []string{}}, "severity": "minor", "summary": "portable concern", "reasoning": "preserve judgment", "references": []any{}}),
		"disposal":    p23RunInput(t, a, "concern", "dispose", a.actor, "portable-disposal", map[string]any{"concern_id": "CON-portable", "action": "rejected", "issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "handled"}),
		"deviation":   p23RunInput(t, a, "deviation", "record", a.actor, "portable-deviation", map[string]any{"id": "DEV-portable", "subject": map[string]any{"type": "issue", "id": a.issueID}, "expected": "accepted", "observed": "different", "authority": map[string]any{"actor_id": "human-reviewer", "reference": nil}, "disposition": "pending-decision", "reasoning": "preserve deviation", "references": []any{}, "supersedes_deviation_id": nil}),
	} {
		p23AssertReviewSuccess(t, name, response)
	}
	return a.workspace
}

func p23PortableRoundTrip(t *testing.T, binary, workspace string) map[string][]string {
	t.Helper()
	live := p23PortableRows(t, workspace)
	archive := filepath.Join(t.TempDir(), "portable.zip")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "export", "create", archive)
	before := p23PortableTreeDigest(t, config.LedgerDir(workspace))
	p21Run(t, binary, workspace, "--project", workspace, "--json", "export", "validate", archive)
	if after := p23PortableTreeDigest(t, config.LedgerDir(workspace)); after != before {
		t.Fatal("export validate mutated source workspace")
	}
	liveEntries := p22ArchiveEntries(t, archive)
	p21Run(t, binary, workspace, "--project", workspace, "--json", "rebuild")
	p23PortableAssertRows(t, "rebuilt", live, p23PortableRows(t, workspace))
	p23AssertRebuiltProjectionIntegrity(t, workspace)
	rebuiltArchive := filepath.Join(t.TempDir(), "rebuilt.zip")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "export", "create", rebuiltArchive)
	if !reflect.DeepEqual(liveEntries, p22ArchiveEntries(t, rebuiltArchive)) {
		t.Fatal("rebuild changed canonical archive Event bytes")
	}
	target := filepath.Join(t.TempDir(), "imported")
	closureRun(t, filepath.Dir(target), "git", "clone", "-q", workspace, target)
	imported := p21Run(t, binary, target, "--json", "export", "import", archive, "--repository", target, "--idempotency-key", "p23-portable")
	if got := p21Path(t, imported, "result", "import", "projection_schema"); got != float64(23) {
		t.Fatalf("imported schema = %v", got)
	}
	p23PortableAssertRows(t, "imported", live, p23PortableRows(t, target))
	p23AssertRebuiltProjectionIntegrity(t, target)
	importedArchive := filepath.Join(t.TempDir(), "imported.zip")
	p21Run(t, binary, target, "--json", "export", "create", importedArchive)
	if !reflect.DeepEqual(liveEntries, p22ArchiveEntries(t, importedArchive)) {
		t.Fatal("import changed canonical archive Event bytes or identities")
	}
	return live
}

func p23PortableRows(t *testing.T, workspace string) map[string][]string {
	t.Helper()
	conn := p22FaultDB(t, workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	result := p23PortableAllRows(t, conn)
	delete(result, "projection_meta")
	delete(result, "schema_migrations")
	return result
}

func p23PortableTableRows(t *testing.T, conn *sql.DB, table string) []string {
	t.Helper()
	rows, err := conn.Query(`SELECT * FROM "` + table + `"`)
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	columns, err := rows.Columns()
	p22FaultOK(t, err)
	result := []string{}
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		p22FaultOK(t, rows.Scan(pointers...))
		encoded, err := json.Marshal(values)
		p22FaultOK(t, err)
		result = append(result, string(encoded))
	}
	p22FaultOK(t, rows.Err())
	sort.Strings(result)
	return result
}

func p23PortableAssertRows(t *testing.T, phase string, before, after map[string][]string) {
	t.Helper()
	if len(before) != len(after) {
		t.Errorf("%s table inventory changed: live=%d current=%d", phase, len(before), len(after))
	}
	for table := range before {
		if !reflect.DeepEqual(before[table], after[table]) {
			t.Errorf("%s %s rows changed: live=%v current=%v", phase, table, before[table], after[table])
		}
	}
}

func p23PortableTreeDigest(t *testing.T, root string) string {
	t.Helper()
	hash := sha256.New()
	err := filepath.WalkDir(root, func(path string, item fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if item.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, _ = hash.Write([]byte(strings.TrimPrefix(path, root)))
		_, _ = hash.Write(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func p23PortableJournalEvents(t *testing.T, workspace string) map[string]journal.Event {
	t.Helper()
	archive := filepath.Join(t.TempDir(), "events.zip")
	p21Run(t, integrationFIDBinary(t), workspace, "--json", "export", "create", archive)
	result := map[string]journal.Event{}
	for name, data := range p22ArchiveEntries(t, archive) {
		if !strings.HasPrefix(name, "journal/") {
			continue
		}
		for _, line := range bytes.Split(bytes.TrimSpace(data), []byte{'\n'}) {
			var event journal.Event
			if err := event.UnmarshalJSON(line); err != nil {
				t.Fatal(err)
			}
			result[event.EventType] = event
		}
	}
	return result
}
