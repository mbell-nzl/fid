# Phase 18 QA Plan - Workflow Tables Migration + Engine Data-Source Swap + Gate Seams

> Phase: 18 - Workflow Engine Migration (v0.2a, slice 1 of 3)
> Refreshed: 2026-07-19
> QA Lead: Codex
> Operating mode: Heavy
> Status: **READY FOR TASK 1 — validation contract Ready; Heavy Mode approved**

## Entry Readiness

Fresh command execution, package discovery, existing test harnesses, and CLI
help checks are feasible, and no required Phase 18 command is a placeholder.
DEC-067 applied the requested one-observable contract refinements, exact Task 4
schema-object inventory, and Task 5 performance controls. The validation
contract is Ready at VAL-001..024, every updated task spec cites its active
assertions, and the Human approved the Heavy Mode plan on 2026-07-19. No
unresolved entry finding or escalation remains; Task 1 may start under the
Phase Runner.

## 1. Scope And Source Documents

Phase 18 replaces hard-coded default-workflow data sources with migrated
workflow tables and DB-backed resolution while preserving v0.1 behavior,
projection rebuild determinism, and journal bytes. It also reserves the
DEC-066 gate-obligation schema/vocabulary seam without enabling receipt or gate
behavior. No new user-facing command is in scope.

Sources loaded for this plan:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/decisions_principles.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/docs/phases/phase-18/decisions.md` (DEC-067)
- `project-docs/docs/phases/phase-18/qa/validation-contract.md`
- `project-docs/docs/phases/phase-18/tasks/task-1.md` through `task-5.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md` §§2.6, 4, 6.1, and 8
- Requirement excerpts named by the brief, contract, and tasks from
  `01_PRD.md`, `02_TRD.md` §§4.3, 4.6, 7 and TRD-ARCH-002/003/008/009,
  `03_TRD_SCHEMA.md` §§1.2, 2, 3, and 14, and
  `10_IMPLEMENTATION_ROADMAP.md` v0.2a/Phase 18
- DEC-066 in
  `project-docs/docs/reviews/2026-07-18-roadmap-update-decisions.md`
- Relevant workflow, evidence-obligation, journal, operating-mode, and
  validation-assertion entries in `project-docs/specs/GLOSSARY.md`
- `project-docs/templates/qa-plan-template.md`,
  `validation-contract-template.md`, and `mid-phase-review-template.md`
- Current migration/test/fixture evidence under `migrations/`,
  `internal/db/`, `internal/workflow/`, `internal/issues/`,
  `internal/journal/`, `internal/maintenance/`, and `cmd/fid/`

There are no Phase 18 completion reports yet; their absence is expected before
Task 1. The Phase 18 run ledger is outside QA Lead write scope and was not
modified. Release-readiness scanners and held-out scenarios remain N/A:
AGENTS.md records no accepted end-to-end/held-out harness, and Phase 18 is not a
public-release sign-off. Ordinary Human dogfood remains planned in section 9.

Non-blocking current-doc drift remains outside this artifact's write scope:
`progress_current.md` still describes Human approval as pending, and the phase
brief's exit checklist still names VAL-001..014 even though its status and
contract index correctly record approval and VAL-001..024. The Ready contract,
DEC-067, and the phase brief's approved status are the execution authorities;
the Tech Lead should synchronize those stale summary lines during the next
current-doc refresh.

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: **YES**.
  The Ready contract covers workflow schema/seeds/seam, registry and transition
  parity, issue FKs/data/schema/read parity, claim entry/release, maintenance,
  full-projection rebuild, projection freshness, replay/journal immutability,
  remaining read/dev consumers, and production Go-map isolation.
- Every `VAL-NNN` has one validator role: **YES**; VAL-001..024 each assign QA
  Lead.
- Every `VAL-NNN` has a concrete evidence method: **YES**; section 4 translates
  each method into reproducible tests, CLI output, SQL/PRAGMA/file state, or
  static inventory evidence.
- Every `VAL-NNN` maps to a task/window or phase-end check: **YES**; the final
  task/window and `QA-NN` map appears below and in sections 3 and 7.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: **YES**. All 24 assertions are deterministic and no
  Human fan-out is needed.
- Conditional or binary assertions require qualifying and non-qualifying
  states: **YES**. VAL-006, VAL-008..010, VAL-012, VAL-015..017, and VAL-024
  name the required accepted/rejected, healthy/corrupt, stale/current,
  implemented/unimplemented, or production/test states.
- Parity assertions have deterministic conformance proof: **YES**. They use
  deep-equal tuple/matrix comparisons, canonical row/schema snapshots, literal
  CLI bytes, normalized full-projection snapshots, and per-file journal hashes.
- Every active VAL-001..024 has a deterministic `QA-NN` mapping: **YES**.
- Revisions requested from Tech Lead: **none**. DEC-067 applied the prior
  refinements and marked the contract Ready.
- One-line coverage rollup: **all in-scope requirements are referenced by at
  least one validation assertion: YES.**

### Final Contract To QA Mapping

| Owning window | Active assertions | Deterministic QA steps | Primary evidence rows |
|---|---|---|---|
| Task 1 | VAL-001, VAL-002, VAL-003, VAL-004, VAL-005, VAL-012, VAL-016 | QA-01..05, QA-07, QA-18 | EVID-01..05, EVID-12 |
| Task 2 | VAL-014, VAL-015; advances VAL-002..004 | QA-06, QA-21 | EVID-02..04, EVID-05, EVID-14 |
| Task 3 | VAL-008, VAL-009, VAL-020 | QA-13..15 | EVID-08, EVID-09 |
| Task 4 | VAL-006, VAL-007, VAL-018, VAL-019 | QA-09..12 | EVID-06, EVID-07 |
| Task 5 / phase end | VAL-010, VAL-011, VAL-013, VAL-017, VAL-021, VAL-022, VAL-023, VAL-024 | QA-08, QA-16, QA-17, QA-19, QA-20, QA-22..24 | EVID-05, EVID-10, EVID-11, EVID-13, EVID-15 |

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA notes |
|---|---|---|---|---|
| Task 1 | Workflow tables, exact seeds, gate seam, projection version | VAL-001..005, VAL-012, VAL-016 | High | Exact DDL/seed set equality, migration down/up, zero reserved-name seeds, loss-free gate JSON round-trip, and stale/current projection-version branches. |
| Task 2 | DB-backed registry and frozen-definition parity | VAL-014, VAL-015; advances VAL-002..004 | High | Decoder preserves role order/params, surfaces both reserved names as known-but-unimplemented, performs no evaluation/write, and deep-equals all three frozen definitions. |
| Task 3 | Transition-path source swap and claim parity | VAL-008, VAL-009, VAL-020 | Critical | Exhaustive ordered-pair/rule matrix, exact errors, all four started-entry claim states, and atomic auto-release; **Review 1 occurs after Task 3 and before Task 4**. |
| Task 4 | SQLite `issues` table rebuild and FKs | VAL-006, VAL-007, VAL-018, VAL-019 | Critical | Exact row/CLI snapshots; insert/update FK probes; down/up and `foreign_key_check`; named dependent-object inventory below; **Review 2 occurs after Task 4 and before Task 5**. |
| Task 5 | Remaining consumers, maintenance, rebuild/journal closure | VAL-010, VAL-011, VAL-013, VAL-017, VAL-021..024 | Critical | Full-projection parity, healthy/corrupt audit states, remaining consumer matrices, no production Go-map reads, and the exact 10k-read/10k-rebuild/100k-rebuild controls below. |

### Task 4 Exact Dependent-Object Inventory

Task 4's pre/post migration-0008 and `.down`/`.up` fingerprint must account for
all of the following; generic “issues indexes” evidence is insufficient:

- Recreate the five triggers defined `ON issues`:
  `issues_parent_kind_guard_insert`, `issues_parent_kind_guard_update`,
  `issues_kind_immutable`, `issues_sequence_scheduling_guard_insert`, and
  `issues_sequence_scheduling_guard_update`.
- Drop/recreate and prove definition/output parity for the `issues_enriched`
  view.
- Preserve unchanged, without dropping/recreating, the two supporting indexes
  on other tables: `idx_concerns_issue_status` and
  `idx_triage_issue_created`.
- Preserve and functionally verify the five migration-0004 other-table kind
  triggers whose bodies query `issues`:
  `acceptance_criteria_issue_kind_guard_insert`,
  `task_completions_issue_kind_guard_insert`,
  `epic_metadata_issue_kind_guard_insert`,
  `triage_events_epic_kind_guard_insert`, and
  `decisions_epic_kind_guard_insert`.
- Preserve and functionally verify the five migration-0005 ref-integrity
  triggers whose bodies query `issues`:
  `links_from_ref_integrity_guard_insert`,
  `links_to_ref_integrity_guard_insert`,
  `claims_ref_integrity_guard_insert`,
  `claims_ref_integrity_guard_update`, and
  `activity_log_ref_integrity_guard_insert`.
- Verify the `issues(id)` PK autoindex and `UNIQUE(parent_id, sequence)`
  autoindex structurally through `PRAGMA index_list`/`index_info` rather than
  relying on generated autoindex names.

### Task 5 Exact Performance Controls

The Task 5 and phase-end `go test -bench . -benchmem ./...` gate must register,
run, and report fixture size plus `ns/op`/the custom duration metric for:

- the retained 10k-event full rebuild at **<5s** (TRD-PERF-002);
- a new 100k-event full rebuild at **<30s** (TRD-PERF-003), directly scaling
  the configurable `seedScaleReplayWorkspace` fixture; and
- a new 10k-issue registry-backed single read at **<500ms**
  (TRD-PERF-001), exercising the migrated DB registry read path rather than the
  frozen Go map.

The existing 1k-event rebuild remains a diagnostic baseline. Crossing any
target is Blocking, even if the benchmark process exits zero.

## 4. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Canonically ordered schema fingerprint from `sqlite_schema`, `PRAGMA table_info`, `index_list/index_info`, and `foreign_key_list` for all four workflow tables | VAL-001 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-02 | Deep-equal ordered `(workflow_id,state,category)` tuples against frozen v0.1 definitions for exactly three workflows | VAL-002 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-03 | Deep-equal transition tuple set for every known-state pair, including every explicit non-cancelled→cancelled edge and absence of disallowed/cancelled→cancelled edges | VAL-003 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-04 | Deep-equal rule tuple set, semantic `params_json` comparison, role-order/error-order preservation, and no extra/missing rule | VAL-004 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-05 | Gate-obligation JSON insert/load/parse/round-trip, reserved-name recognition, seed exclusion, and phase-end static/runtime evidence of no receipt/gate evaluator; recorded separately | VAL-005, VAL-015..017 | Developer, QA Reviewer, QA Lead | Tasks 1, 2, 5 / reviews / final |
| EVID-06 | FK metadata plus raw INSERT and UPDATE probes for unknown workflow, known-workflow/invalid-state, and valid pair; clean `PRAGMA foreign_key_check` | VAL-006 | Developer, QA Reviewer, QA Lead | Task 4 / Review 2 / final review |
| EVID-07 | Canonical pre/post issue-row snapshot, the complete named dependent-object inventory from section 3, and frozen `issue get/list` output matrix, as three separately reported results | VAL-007, VAL-018, VAL-019 | Developer, QA Reviewer, QA Lead | Task 4 / Review 2 / final review |
| EVID-08 | Exhaustive frozen-v0.1-vs-registry transition matrix for every known state pair and every rule success/failure branch, plus built-CLI exact output/error checks | VAL-008 | Developer, QA Reviewer, QA Lead | Task 3 / Review 1 / final review |
| EVID-09 | Missing/wrong/stale/live claim entry cases and separate same-event out-of-started claim release proof | VAL-009, VAL-020 | Developer, QA Reviewer, QA Lead | Task 3 / Review 1 / final review |
| EVID-10 | Healthy migrated workspace with no workflow finding; test-only corrupt pair with deterministic finding; same violation under `--warn-as-error` with non-zero exit and unchanged finding payload | VAL-010 | Developer, QA Lead | Task 5 / final review |
| EVID-11 | Normalized per-table snapshots of the complete projection before deletion, after rebuild, and after second rebuild; equal workflow rows/issues; clean FK/schema evidence | VAL-011 | Developer, QA Lead | Task 5 / final review |
| EVID-12 | Prior-version projection triggers rebuild/catch-up to current; otherwise-fresh current-version projection remains fresh | VAL-012 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-13 | Golden v0.1 `issue.created`/`issue.transitioned` replay result and separate hashes/bytes for every journal file before/after migrate, refresh, rebuild, and second rebuild | VAL-013, VAL-021 | Developer, QA Lead | Task 5 / final review |
| EVID-14 | Canonical registry-vs-frozen-definition deep equality across every state/category/edge/rule/params tuple for all three workflows | VAL-014 | Developer, QA Reviewer, QA Lead | Task 2 / Review 1 / final review |
| EVID-15 | Post-swap `issue get/list/ready` output matrix, dev start/finish/block effects, and production-source reference inventory, reported as separate evidence | VAL-022..024 | Developer, QA Lead | Task 5 / final review |
| EVID-16 | Exact per-task command results, coverage package percentages, committed report-lint results, verified hashes/subjects, review verdicts, phase gates, and spec-sync disposition | no VAL | Developer, QA Reviewer, QA Lead | Every task/review/final review |
| EVID-17 | Exact benchmark output for the 10k-issue registry-backed single read (<500ms), retained 10k-event rebuild (<5s), and new 100k-event rebuild (<30s), with fixture sizes and duration metrics, plus the `issues_enriched` query-plan regression result | no VAL | Developer, QA Lead | Task 5 / final review |

## 5. Deterministic Gates

No required command contains a placeholder. Each completion report must record
every exact command in its task spec verbatim as `PASS`, or a blocker report
must record it `BLOCKED` with reason. Missing, skipped, or narrower substituted
commands fail the gate.

| Gate | Command | Expected result | Required when |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | Exit 0; no files listed | Every task and phase end |
| Lint/static analysis | `golangci-lint run ./...` | Exit 0; `0 issues.` | Every task and phase end |
| Coverage execution | `go test -race -coverprofile=coverage.out ./...` | All packages PASS; profile produced | Every task and phase end |
| Coverage report | `go tool cover -func=coverage.out` | Exit 0; COV-001 packages `internal/issues`, `claims`, `workflow`, `policy`, and `journal` each >=80%; `cmd/fid` >=70%; no threshold decrease | Every task and phase end |
| Unit suite | `go test -race ./...` | All packages PASS | Phase end per Quality Plan §5.2 |
| Integration/CLI E2E | `go test -tags integration -race ./...` | All packages PASS, including Phase 18 migration/rebuild/CLI parity tests | Every task and phase end; explicitly required by Tasks 1–4 |
| Separate E2E harness | N/A | AGENTS.md records no accepted separate E2E/held-out command; integration-tagged CLI journeys are used | Phase 18 |
| Benchmarks | `go test -bench . -benchmem ./...` | All benchmarks execute; retained 10k-event rebuild <5s, new 100k-event rebuild <30s, and new 10k-issue registry-backed single read <500ms; target breaches fail the gate even on exit 0 | Task 5 and phase end |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | Exit 0; `/tmp/fid` exists | Task 5 and phase end |
| Full phase gate | The seven executable commands above in documented order (format, lint, coverage, coverage report, integration, benchmarks, build), plus the exact unit suite | Every command PASS with literal output excerpts | Task 5/final review |

At phase end QA also performs PRD/TRD compliance, SQL/input security review,
performance spot-check, documentation/spec-sync review, deterministic
checklist execution, and Human-checklist generation per Quality Plan §5.2.
Release-readiness scans are not silently substituted into this non-release
phase.

## 6. Command And Feasibility Validation

Fresh Heavy Mode validation was performed on 2026-07-19. No result was carried
forward from Phase 17.

### 6.1 Command Validation

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Toolchain | Exact version commands | PASS | `go version go1.26.2 linux/amd64`; `golangci-lint 2.12.2` built with Go 1.26.2. |
| `gofmt -l cmd internal` | Exact command plus `gofmt -h` | PASS | Exit 0, no files listed; help confirms `-l` and path operands. |
| `golangci-lint run ./...` | Exact command plus `run --help` | PASS | Exit 0, `0 issues.`; help confirms the `run` syntax. |
| `go test -race -coverprofile=coverage.out ./...` | Exact command | PASS | Exit 0. Package totals: `issues` 81.5%, `claims` 81.1%, `workflow` 95.5%, `policy` 86.3%, `journal` 81.0%, `cmd/fid` 77.7%. |
| `go tool cover -func=coverage.out` | Exact command plus help | PASS | Exit 0; function rows and `total ... 76.6%` emitted. `go tool cover -help` exits 2 after printing usage, which confirms `-func`; package thresholds use the preceding `go test` package totals. |
| `go test -race ./...` | Exact command plus `go help testflag` | PASS | Exit 0; `cmd/fid` completed in 63.332s and all packages passed. Help confirms `-race`. |
| `go test -tags integration -race ./...` | Exact command, build-tag inventory, and test listing | PASS | Exit 0; `cmd/fid` completed in 76.292s. Thirteen integration-tagged files exist, including closure, transition-rule, rebuild, and CLI journeys. |
| `go test -bench . -benchmem ./...` | Exact command plus help, current benchmark inventory, configurable fixture inspection, and Task 5 acceptance review | PASS | Exit 0. Baseline: 1k rebuild 458.0 ms; 10k rebuild 3990 ms. Help confirms `-bench`/`-benchmem`; `seedScaleReplayWorkspace` accepts arbitrary event counts for the new 100k case; Task 5 requires registration of the 100k-event rebuild and 10k-issue registry-backed read under this exact gate before task close. |
| `go build -o /tmp/fid ./cmd/fid` | Exact command plus `go help build` and file state | PASS | Exit 0; `/tmp/fid` produced at 12,893,890 bytes; help confirms `-o`. `/tmp` is writable. |
| Package scopes | `go list ./...` and `go test -list` on relevant packages | PASS | `cmd/fid`, `internal/db`, `workflow`, `issues`, `journal`, and `maintenance` resolve; existing migration, transition, claim, maintenance, projection-version, rebuild, and benchmark tests list successfully. |
| Completion-report lint | Script usage inspection and `project-docs/tools/report-lint.sh --self-test` | PASS | Self-test passed clean, missing/skipped/narrower command, blocker, fabricated/wrong hash, subject, uncommitted, and dirty-report cases. Git 2.53.0 is available. |
| Phase 18 placeholder scan | `rg` over AGENTS.md, five task specs, and the contract | PASS | No command placeholder found; only AGENTS.md's explanatory `<TBD>` sentence matched. |

### 6.2 Fixture, Artifact, And Checklist Feasibility

| Referenced assumption or step | Validation method | Result | Evidence |
|---|---|---|---|
| Frozen v0.1 workflow reference | Existing test listing plus successful full suite; source/TRD excerpt comparison | PASS | `TestWorkflowDefinitionsMatchTRDExactly` exists and passed; `internal/workflow/workflow.go` contains all three definitions and `02_TRD.md` §4.6 names the same states/categories/rules. Task 1 must freeze, not rewrite, this reference. |
| Migration version targeting and down/up | Existing helper/source inspection | PASS | `internal/db` uses golang-migrate; `newTestMigrator` exposes `Migrate(version)`, and existing tests target versions 1 and 5. `MigrateDown` also exists. |
| Workflow schema/seed introspection | Existing DB helpers and SQLite PRAGMAs | PASS | `assertTableColumns`, `assertForeignKeysContain`, `sqlite_schema`, and temp-DB migration patterns exist under `internal/db`; Task 1 outputs are declared. |
| Exact Task 4 dependent-object set | `CREATE`/body inventory across migrations 0001–0006 plus Task 4/VAL-018 comparison | PASS | Section 3 names the five `ON issues` triggers, `issues_enriched`, both surviving named indexes, all five 0004 kind triggers and all five 0005 ref-integrity triggers whose bodies query `issues`, and both structural autoindexes. Task 4 requires pre/post/down/up `sqlite_schema` + PRAGMA fingerprints. |
| Pre/post migration row and CLI output fixture | Target-version migrator plus existing CLI migration/rebuild helpers | PASS | Existing helpers can create a version-7 DB, snapshot rows/output, migrate to version 8, and compare. Task 4 must commit the frozen fixture/expected output; a same-binary auto-migrating “before” call without version control is insufficient. |
| Invalid/valid FK pairs | Existing raw-SQL fixture approach | PASS | Maintenance tests already disable FKs/checks only in test setup and restore them, while normal DB opens enforce FKs. Task 4 can test ordinary rejection; Task 5 can seed corruption through this test-only bypass. |
| Exhaustive transition and claim parity | Existing tests and integration test listing | PASS | Current tests cover allowed/invalid/rule/cancel and into/out-of-started claim behavior. Frozen definitions provide finite state/rule sets for exhaustive generated matrices. |
| Maintenance healthy/corrupt and `--warn-as-error` | Existing fixtures, tests, and built CLI help | PASS | `newValidStructuralFixture`/`newStructuralFixture` provide both states; `fid maintenance check --help` confirms `--warn-as-error` and `--format`. |
| Full projection and journal-byte parity | Existing closure/rebuild helpers | PASS | `closureProjectedTableSnapshots`, `assertClosureProjectedTablesEqual`, `assertClosureJournalBytesEqual`, and repeated rebuild scale tests exist. Task 5 must add the four workflow tables and any omitted projected tables to the snapshot inventory. |
| Projection-version stale/current branches | Existing test harness | PASS | `TestEnsureProjectionFreshRebuildsStoredPriorProjectionSchemaVersion` and freshness-decision tests exist; Task 1 can extend them from version 3→4 and add the fresh counterpart. |
| No production reads of frozen map | `rg -n 'definitions' internal cmd --glob '*.go'` and `rg -n 'workflow\.' internal cmd --glob '*.go' --glob '!**/*_test.go'` | PASS as reproducible future checks | Current definition/production call sites are enumerated in `internal/workflow`, `internal/issues`, `internal/dev`, and `internal/maintenance`; after Task 5, the same inventory must return only approved fixture/test references. Source-state proof is paired with runtime matrices. |
| CLI steps for Human dogfood | Built binary help | PASS | `issue transition/get/list`, `ready`, `maintenance check`, and `rebuild` help all exit 0 with the required flags. Human steps require a disposable copy, never the sole real workspace. |
| 10k-read / 10k- and 100k-rebuild performance fixtures | Benchmark/helper inspection, package test-fixture inventory, and Task 5 acceptance review | PASS | `seedScaleReplayWorkspace(testing.TB, workspace, targetEvents)` directly supports the new 100k rebuild scale and retained 10k scale. Existing temp-DB/issue fixtures make a 10k-issue registry-backed read benchmark feasible. Task 5 must add both missing benchmark registrations and report fixture size plus duration metrics under the exact benchmark command. |
| Completion reports and reviews | Task specs, template/script usage, and file checks | PASS as future artifacts | `task-1-done.md` through `task-5-done.md` are intentionally absent before work. Each task spec declares its path. Reviews are feasible from the committed reports and template after Tasks 3 and 4. |
| Held-out/Mode B scenario | AGENTS.md and Quality Plan §7 | N/A | No accepted held-out location/command exists; Phase 18 is not release readiness. No hidden scenario material was read. |

## 7. Deterministic Checklist

Every item is executed by the QA Lead at phase end, after the relevant
Developer/Reviewer evidence is available. Literal command output or exact
file/DB state is required; code reading alone cannot produce a PASS. All
VAL-001..024 IDs below are active in the Ready contract.

| ID | Check | Covers | Evidence required |
|---|---|---|---|
| QA-01 | Compare the exact four-table schema fingerprint with TRD_SCHEMA §3, including clarified PK/UNIQUE/FK expectations `(qa-lead)` | VAL-001 | Test output plus canonical PRAGMA/`sqlite_schema` rows. |
| QA-02 | Compare the complete ordered state/category tuple set for all three defaults with the frozen v0.1 reference `(qa-lead)` | VAL-002 | Deep-equality test output showing no extra/missing tuple. |
| QA-03 | Compare every known-state edge pair, including all cancel and absent-edge states, with the frozen v0.1 reference `(qa-lead)` | VAL-003 | Exhaustive matrix/set-diff output, empty diff. |
| QA-04 | Compare every edge rule and parsed params value with the frozen v0.1 reference `(qa-lead)` | VAL-004 | Deep-equality output with rule/role/params tuple counts. |
| QA-05 | Round-trip the canonical gate-obligation JSON without loss `(qa-lead)` | VAL-005 | Test output and loaded canonical JSON/typed value. |
| QA-06 | Prove both reserved rule names are recognized as known-but-unimplemented by the loader `(qa-lead)` | VAL-015 | Test output for both names; no panic/evaluation/write. |
| QA-07 | Prove no seed row uses either reserved rule name `(qa-lead)` | VAL-016 | Ordered SQL query output with zero rows. |
| QA-08 | Prove no production evaluator can satisfy or query receipt/gate state for either reserved rule `(qa-lead)` | VAL-017 | Runtime injected-rule outcome plus static production reference/call inventory. |
| QA-09 | Exercise INSERT and UPDATE with unknown workflow, invalid known pair, and valid pair; run `foreign_key_check` `(qa-lead)` | VAL-006 | Rejection/acceptance errors and zero-row PRAGMA output. |
| QA-10 | Compare every `issues` row/column before and after migration 0008 `(qa-lead)` | VAL-007 | Canonical row snapshots, counts, and empty diff. |
| QA-11 | Compare the complete section 3 dependent-object fingerprint before/after migration and after down/up `(qa-lead)` | VAL-018 | Named `sqlite_schema`/PRAGMA snapshots, functional retained-trigger probes, and empty diff. |
| QA-12 | Compare frozen `fid issue get` and `list` outputs before/after migration 0008 `(qa-lead)` | VAL-019 | Literal pre/post outputs or hashes with byte-equality result. |
| QA-13 | Execute the exhaustive registry-backed transition matrix and exact built-CLI parity cases `(qa-lead)` | VAL-008 | Matrix counts, empty diff, and literal allowed/invalid/blocked/cancel output excerpts. |
| QA-14 | Exercise missing, wrong-holder, stale, and live claims when entering started `(qa-lead)` | VAL-009 | Exact errors and accepted transition/DB state. |
| QA-15 | Prove leaving started releases the live caller claim atomically `(qa-lead)` | VAL-020 | Event effect plus post-projection and rebuilt claim absence. |
| QA-16 | Run maintenance against healthy and test-only corrupt workflow pairs, then rerun corrupt state with `--warn-as-error` `(qa-lead)` | VAL-010 | Literal JSON output, ordered finding, exit statuses, unchanged journal bytes. |
| QA-17 | Delete and rebuild the full projection, compare every table, rebuild again, and compare again `(qa-lead)` | VAL-011 | Per-table snapshot hashes/diffs, schema/FK state, and both rebuild outputs. |
| QA-18 | Exercise prior and current projection-schema versions `(qa-lead)` | VAL-012 | Rebuild/catch-up proof for prior and no version-triggered rebuild for current. |
| QA-19 | Replay the frozen v0.1 journal into the Phase 18 schema and compare the expected projection `(qa-lead)` | VAL-013 | Golden replay output and projection snapshot diff. |
| QA-20 | Compare every journal file before/after migrate, refresh, rebuild, and second rebuild `(qa-lead)` | VAL-021 | File list, byte counts/hashes, and byte-equality output at each boundary. |
| QA-21 | Compare the DB-loaded registry with the frozen definitions for all states/categories/edges/rules/params `(qa-lead)` | VAL-014 | Canonical deep-equality output for all three workflows. |
| QA-22 | Compare post-swap issue get/list status-category output and ready selection/order with frozen v0.1 expectations `(qa-lead)` | VAL-022 | Matrix output and literal CLI excerpts. |
| QA-23 | Compare `fid dev start`, `finish`, and `block` results/effects with frozen v0.1 expectations `(qa-lead)` | VAL-023 | Built-CLI/integration output plus event/projection effect snapshots. |
| QA-24 | Verify no production file reads the frozen definitions map after Task 5 `(qa-lead)` | VAL-024 | Output from the two exact `rg` inventory commands in section 6, limited to approved fixture/test references, paired with QA-22/23 runtime evidence. |
| QA-25 | Run every exact phase gate and enforce COV-001/COV-002 `(qa-lead)` | no VAL | Literal command outputs and package percentages. |
| QA-26 | Verify migrations 0007 and 0008 down/up cleanly restore their documented prior/current states `(qa-lead)` | no VAL | Migration test output, schema/row snapshots, and clean FK check. |
| QA-27 | Review registry/maintenance SQL for parameter binding, dynamic identifier confinement, test-only corruption bypass, and no journal/data-loss path `(qa-lead)` | no VAL | Test output plus reproducible query/reference inventory and findings. |
| QA-28 | Run the 10k-issue registry-backed read, retained 10k-event rebuild, new 100k-event rebuild, and `issues_enriched` query-plan controls against TRD-PERF-001..004 `(qa-lead)` | no VAL | Exact benchmark output with fixture sizes, `ns/op`/duration metrics and target comparisons (<500ms/<5s/<30s), plus EXPLAIN index evidence. |
| QA-29 | Verify all five committed completion reports against exact task commands and verified commit metadata `(qa-lead)` | no VAL | Five report-lint PASS outputs plus hashes/subjects copied from Git. |
| QA-30 | Verify Tech Lead spec sync records `touched`, `verified-no-change`, or `N/A` for the full canonical spec suite and updates the known TRD §4.6/Schema §14.11 wording `(qa-lead)` | no VAL | Committed spec/decision diff and wrap-up sync table. |

## 8. Mid-Phase Review Cadence

The Heavy Mode cadence is stricter than the default after-every-three-tasks
rule because Tasks 3 and 4 are independent critical swaps.

- **Review 1 after Task 3, before Task 4:** write and commit
  `qa/review-1.md`, covering Tasks 1–3 and VAL-001..005, VAL-008, VAL-009,
  VAL-012, VAL-014..016, and VAL-020. Independently rerun at least one
  representative required command, run report-lint for all three committed
  reports, and verify the exhaustive transition/claim evidence. `VERDICT:
  BLOCK` stops Task 4.
- **Review 2 after Task 4, before Task 5:** write and commit
  `qa/review-2.md`, covering Task 4 and VAL-006, VAL-007, VAL-018, and
  VAL-019. Independently rerun the migration/FK/complete named-object/CLI-read
  parity scope and report-lint Task 4. Any row loss, undefined/missing object,
  FK anomaly, or narrower command evidence is Blocking. `VERDICT: BLOCK` stops
  Task 5.
- Task 5 is followed by the comprehensive QA Lead final review rather than a
  third mid-phase review. Any cross-task defect found in Task 5 is triaged
  before final-review PASS.
- Reviews use `project-docs/templates/mid-phase-review-template.md`, compare
  completion-report commands verbatim with task specs and this plan, cite each
  in-window VAL, and are committed by the QA Reviewer under its own authorship.

## 9. Human Verification Scope

All validation assertions remain QA Lead-owned and machine-checkable; Human
observations cannot pass a VAL. Heavy Mode still requires a Human checklist and
ordinary dogfood. Run these only after deterministic QA on a disposable copy
of a v0.1 workspace, never the sole real workspace.

| ID | Check | Covers | Human evidence expected |
|---|---|---|---|
| HUMAN-01a | Judge whether representative allowed and rejected transitions feel behaviorally unchanged after migration `(human)` | no VAL | Command notes identifying the disposable workspace and observed success/error behavior. |
| HUMAN-01b | Judge whether issue get/list/ready presentation and ordering feel unchanged and sensible `(human)` | no VAL | Observed outputs and concise product judgment. |
| HUMAN-01c | Judge whether `fid maintenance check` remains clear and actionable on the migrated copy `(human)` | no VAL | Output excerpt and clarity judgment. |
| HUMAN-01d | Decide whether the demonstrated backup/copy, migration, and rebuild evidence is sufficient to accept Phase 18 for continued dogfood `(human)` | no VAL | Explicit accept/request-changes decision with rationale. |

No Mode B held-out scenario is listed because no held-out scenario location or
command is accepted for this project. The final `human-checklist.md` contains
only the four rows above.

## 10. Phase-End Decision Rule

Entry is **READY**: DEC-067 is applied, the contract and all five task specs are
aligned on active VAL-001..024, this plan's evidence/task mappings are final,
all referenced commands and planned controls are feasible, the Human approved
Heavy Mode, and no unresolved entry blocker or escalation remains. The Phase
Runner may begin Task 1.

After implementation, the phase can proceed to final-review PASS only when:

- All five task implementations and finalized completion reports are committed;
  every task's exact command is recorded `PASS`, or the phase remains stopped
  on an accepted blocker report.
- Review 1 after Task 3 and Review 2 after Task 4 are committed with no
  unresolved Blocking finding and `VERDICT: CONTINUE`.
- Every executable deterministic gate in section 5 passes with literal output;
  coverage thresholds and applicable TRD performance targets are met.
- Every active validation assertion VAL-001..024 has a final `PASS`, `FAIL`,
  `BLOCKED`, or `CARRY-FORWARD` outcome with cited evidence. `PASS` without
  evidence is invalid.
- Every `(qa-lead)` item has reproducible execution/file/DB evidence; source
  reading alone is not used to pass behavior.
- Human checks remain human-judgment-only, are copied to the Human checklist,
  and are accepted or converted into follow-up work.
- No unresolved Blocking defect, data-loss/journal-integrity concern, spec
  deviation, command anomaly, or required escalation remains.
- Tech Lead records full canonical spec-suite sync, including the known
  `02_TRD.md` §4.6 and `03_TRD_SCHEMA.md` §14.11 source-of-truth updates, and
  records any migration/rollback documentation required for dogfood.
