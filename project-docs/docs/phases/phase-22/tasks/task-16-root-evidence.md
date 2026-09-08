# Task 16 root evidence and removal inventory

Root implemented Task16 directly after independent Review5 PASS (8125d8d,
ledger d562e5d). No concurrent product writer. Working directory for all
commands: `/home/mark/projects/fid`.

## First behavioral red and focused verification

- `go test ./internal/db -run '^TestFinalSchemaHasNoEscalationsOrEscalationReferences$' -count=1` — FAIL — `table "escalations" exists` (terminated `exit_code: 1`).
- After copying both exact schema.md section4 SQL blocks and updating the
  current version/profile, the byte-identical command — PASS — `ok  	github.com/mbell-nzl/fid/internal/db	0.050s` (terminated `exit_code: 0`).
- `go test ./internal/db ./internal/journal ./internal/maintenance ./internal/workflow ./internal/sqlquery ./cmd/fid -run 'Test(FinalSchema|LatestProjection|RebuildWorkspacePreservesFullPhaseTwo|MigrateUpCreatesRetained|RegistryEquals|LoadRetainsProjectDecision|PolymorphicRef|LinkScope|Check|RunReturnsStructuredRows|SQLCommand)' -count=1` — PASS — `ok  	github.com/mbell-nzl/fid/cmd/fid	1.417s` (terminated `exit_code: 0`). All six named packages executed matching tests; this is focused evidence, not the final gate.

## Complete removal-root classification

Root searched `rg -l 'escalation|esc_id|escalated_to' cmd internal migrations`
and classified all21 files present at the pre-edit checkpoint:

| Paths | Disposition and executing gate |
|---|---|
| migrations/0002 up/down, 0005 up, 0010 up/down (five files) | Retained historical producers exactly unchanged; new0016 is the accepted removal. Earlier migration tests and full reset run in ordinary race gate. |
| internal/db/escalation_workflow_test.go, workflow_assignments_test.go, operational_kernel_migration_test.go | Explicit historical13, 8→10→8, and13→14 boundaries retain seed assertions unchanged; ordinary gate. |
| internal/db/migrate_v01_tables_test.go | Latest schema inventory removes exactly escalation and two FK columns; retained triage/Concern rows remain; ordinary gate. |
| internal/db/ref_triggers_test.go | Latest fixture no longer creates impossible escalation row; retired endpoint rejection and retained resolvable references remain; ordinary gate. |
| internal/maintenance/check.go, errors.go, test_helpers_test.go | Current activity reference domain/resolver removes escalation. Missing-Event and healthy fixture rows use retained subjects, preserving counts and no-mutation checks; ordinary gate. |
| internal/workflow/rules.go, registry_test.go | Remove retired seed rationale/count, retain generic project Decision reserved rule and query-only registry checks; ordinary gate. |
| internal/workflow/authoring_test.go, cmd/fid/root_test.go | Existing negative assertions reject escalation applicability and CLI leaf; retained unchanged, ordinary gate. |
| internal/journal/phase2_rebuild_test.go | Full historical fixture now explicitly uses trusted migrations1..5 for original/rebuilt projections, preserving every original historical row and trigger. New latest tests independently reject retired table/columns and roll back preceding Effects; ordinary gate. |
| cmd/fid/sql_test.go, internal/sqlquery/sqlquery_test.go | LIMIT2 remains two ordered rows: epic-default and subtask-default. Only tests converge; production sqlquery.go untouched (CF018). Ordinary gate. |
| cmd/fid/integration/e2e_v01_closure_helpers_test.go | Latest snapshot excludes retired table; all other query subjects retained; full tagged gate. |

The new0016 up/down are exact copies of the accepted blocks; partial down is
not a supported product path. The latest-schema absence test examines all
sqlite_schema definitions, retained tables, default-workflow rows and removed
activity reference rejection; a separate full zero/down/up test succeeds.
No old-data conversion or replacement escalation alias exists.

Root's byte-comparison script extracted both SQL fences from schema.md section4
and compared each new migration's complete bytes. It reported
`EXACT 0016 up/down match accepted schema section4 byte-for-byte`
(terminated `exit_code: 0`). This is direct file-state evidence in addition to
the executing schema/reset tests.

Consequential current-profile changes: ActiveMigrationVersion16, current v2
semantic validator16 and its three existing semantic cases. Legacy v1 remains
schema14. Task15's fake column inventory already excludes the removed columns
and table; its44 accepted Effect tables are unchanged, so its provenance comment
alone converges to16. Task17 separately owns capability qualification.

Activity remains Event-derived inspection: journal/project_effects.go requires
op activity for activity_log, and internal/activity reads that projection.
No Run/Stop/recovery authority is introduced there. Legacy triage status readers,
Concern/Review/history and their recorded Phase23/24/25 expiry stay intact.

## Retained assertion provenance

Root inspected `git blame --line-porcelain` at the prior HEAD for the changed
rows below and read each introducing commit's task footer. These are row
origins, not inferred from filenames. Task16's accepted removal/profile change
invalidates each current assertion; the historical fixture is explicitly scoped
instead of weakened. Every row executes under the ordinary final race command,
except the tagged closure snapshot which executes under the exact integration
command. Fake schema comment provenance is Task15/5ec69e6 and is not an assertion.

| Changed retained row (prior HEAD line) | Introducing row commit / task | Displaced assertion or fixture |
|---|---|---|
| `cmd/fid/integration/e2e_v01_closure_helpers_test.go:471` | `7119d5b38d79663f544f093d55f676c375746419` / Refs: phase-17/task-9 | `name:  "escalations",` |
| `cmd/fid/integration/e2e_v01_closure_helpers_test.go:472` | `7119d5b38d79663f544f093d55f676c375746419` / Refs: phase-17/task-9 | `query: "SELECT * FROM escalations ORDER BY id",` |
| `cmd/fid/sql_test.go:28` | `c0fd2cd8343826fdcd1a597f8f2f3a671e0e15dc` / Refs: phase-21/task-2 | `{"id": "escalation-default", "active": int64(1)},` |
| `internal/db/migrate_v01_tables_test.go:57` | `5b1467f6da2397ab78152c32d3ab828389e51d96` / Refs: phase-2/task-1 | `columns:            []string{"id", "issue_id", "raised_in_completion_id", "status", "escalated_to_esc_id", "triaged_in_event_id"},` |
| `internal/db/migrate_v01_tables_test.go:58` | `5b1467f6da2397ab78152c32d3ab828389e51d96` / Refs: phase-2/task-1 | `foreignKeyContains: []string{"issue_id->issues.id", "raised_in_completion_id->task_completions.id", "escalated_to_esc_id->escalations.id", "triaged_in_event_id->triage_events.id"},` |
| `internal/db/migrate_v01_tables_test.go:97` | `5b1467f6da2397ab78152c32d3ab828389e51d96` / Refs: phase-2/task-1 | `table:              "escalations",` |
| `internal/db/migrate_v01_tables_test.go:103` | `5b1467f6da2397ab78152c32d3ab828389e51d96` / Refs: phase-2/task-1 | `columns:            []string{"id", "epic_id", "issue_id", "triager_role", "source", "source_ref_id", "verdict", "esc_id", "reasoning_md", "created_at"},` |
| `internal/db/migrate_v01_tables_test.go:104` | `5b1467f6da2397ab78152c32d3ab828389e51d96` / Refs: phase-2/task-1 | `foreignKeyContains: []string{"epic_id->issues.id", "issue_id->issues.id", "esc_id->escalations.id"},` |
| `internal/db/ref_triggers_test.go:323` | `457f480812e01555b9505726f14dc67649f5d1ab` / Refs: phase-2/task-4 | `func insertRefTriggerEscalation(t *testing.T, conn *sql.DB, escalationID string) {` |
| `internal/db/ref_triggers_test.go:327` | `457f480812e01555b9505726f14dc67649f5d1ab` / Refs: phase-2/task-4 | `INSERT INTO escalations (id, status, trigger_class, raised_by_role, raised_at)` |
| `internal/db/ref_triggers_test.go:329` | `457f480812e01555b9505726f14dc67649f5d1ab` / Refs: phase-2/task-4 | ```, escalationID, "open", "policy", "qa_lead", "2026-06-07T15:09:00Z"); err != nil {`` |
| `internal/db/ref_triggers_test.go:330` | `457f480812e01555b9505726f14dc67649f5d1ab` / Refs: phase-2/task-4 | `t.Fatalf("insert escalations(%q) error = %v", escalationID, err)` |
| `internal/export/archive_semantics_test.go:122` | `acb569bb9447a3259991520669167fdc819c1706` / Refs: phase-22/task-13 | `validator, err := newSemanticValidator(testMigrationsDir(t), 15)` |
| `internal/export/archive_semantics_test.go:135` | `acb569bb9447a3259991520669167fdc819c1706` / Refs: phase-22/task-13 | `validator, err := newSemanticValidator(testMigrationsDir(t), 15)` |
| `internal/export/archive_semantics_test.go:179` | `acb569bb9447a3259991520669167fdc819c1706` / Refs: phase-22/task-13 | `validator, err := newSemanticValidator(testMigrationsDir(t), 15)` |
| `internal/journal/phase2_rebuild_test.go:27` | `457f480812e01555b9505726f14dc67649f5d1ab` / Refs: phase-2/task-4 | `conn := openReplayDB(t, dbPath)` |
| `internal/journal/phase2_rebuild_test.go:44` | `c0fd2cd8343826fdcd1a597f8f2f3a671e0e15dc` / Refs: phase-21/task-2 | `rebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t))` |
| `internal/maintenance/test_helpers_test.go:132` | `d2a087d858fe12a1b33d9380e03e4497d9194350` / Refs: phase-17/task-3 | `seedActivity(t, db, "ACT-missing-event-second", "EVT-20260715-999999", "escalation", "ESC-valid")` |
| `internal/maintenance/test_helpers_test.go:162` | `d2a087d858fe12a1b33d9380e03e4497d9194350` / Refs: phase-17/task-3 | `seedActivity(t, db, "ACT-escalation", "EVT-20260715-000001", "escalation", "ESC-valid")` |
| `internal/maintenance/test_helpers_test.go:250` | `d2a087d858fe12a1b33d9380e03e4497d9194350` / Refs: phase-17/task-3 | ``if _, err := db.Exec(`INSERT INTO escalations (id, status) VALUES (?, 'open')`, id); err != nil {`` |
| `internal/maintenance/test_helpers_test.go:251` | `d2a087d858fe12a1b33d9380e03e4497d9194350` / Refs: phase-17/task-3 | `t.Fatalf("insert escalation %q: %v", id, err)` |
| `internal/sqlquery/sqlquery_test.go:95` | `c0fd2cd8343826fdcd1a597f8f2f3a671e0e15dc` / Refs: phase-21/task-2 | `{"id": "escalation-default", "active": int64(1), "applies_to": "escalation"},` |
| `internal/workflow/registry_test.go:26` | `a509d59efcc826df417f64ea72a43d743936e669` / Refs: phase-19/task-1 | `if len(registry.workflows) != len(frozen)+1 {` |
| `internal/workflow/registry_test.go:27` | `a509d59efcc826df417f64ea72a43d743936e669` / Refs: phase-19/task-1 | `t.Fatalf("loaded workflow count = %d, want %d frozen defaults plus escalation-default", len(registry.workflows), len(frozen)+1)` |
| `internal/workflow/registry_test.go:116` | `a509d59efcc826df417f64ea72a43d743936e669` / Refs: phase-19/task-1 | `rules, ok := registry.RulesForEdge("escalation-default", "open", "resolved")` |
| `internal/workflow/registry_test.go:118` | `a509d59efcc826df417f64ea72a43d743936e669` / Refs: phase-19/task-1 | `t.Fatal("loaded escalation-default open -> resolved edge is missing")` |
| `internal/workflow/registry_test.go:122` | `a509d59efcc826df417f64ea72a43d743936e669` / Refs: phase-19/task-1 | `t.Fatalf("escalation rule = %#v, want %#v", rules, want)` |
