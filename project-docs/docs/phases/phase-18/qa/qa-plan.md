# Phase 18 QA Plan - Workflow Tables Migration + Engine Data-Source Swap + Gate Seams

> Phase: 18 - Workflow Engine Migration (v0.2a, slice 1 of 3)
> Refreshed: 2026-07-19
> QA Lead: Codex
> Operating mode: Heavy
> Status: **BLOCKED BEFORE TASK 1 — contract and task-spec refinements required**

## Blocking Entry Findings

Fresh command execution, package discovery, existing test harnesses, and CLI
help checks are feasible. No required Phase 18 command is a placeholder. Task
execution remains blocked for the following planning defects, plus the Human
approval already required by the phase brief.

| ID | Blocking finding | Evidence | Required resolution |
|---|---|---|---|
| QA-BLK-18-01 | The draft validation contract bundles independent deterministic observables in VAL-005, VAL-007, VAL-009, and VAL-013, and it has no assertion for Task 5's issue-read/ready parity, dev-composite parity, or removal of production reads from the frozen Go map. VAL-011 is also narrower than the PRD-OPS-003/004 claim that the query store is reconstructed from scratch. | `validation-contract.md` assertion/pass-criterion text compared with Tasks 1–5 and `08_QUALITY_PLAN.md` §5.4. Section 2 below gives one-observable replacements and ownership. | Tech Lead applies the proposed assertion splits/additions, widens VAL-011, refines VAL-008/012 evidence criteria, updates every affected task's coverage/acceptance rows, and marks the revised contract Ready. QA Lead then refreshes the mappings in this plan. |
| QA-BLK-18-02 | Task 4 requires recreation of “issues indexes,” but migrations 0001–0006 define no named index on `issues`; the two named indexes support `issues_enriched` on `concerns` and `triage_events`. The required dependent-object set is therefore not exact. | `rg -n '^CREATE ' migrations/*.up.sql` found `idx_concerns_issue_status`, `idx_triage_issue_created`, `issues_enriched`, five triggers directly on `issues`, and additional triggers whose SQL references `issues`. | Tech Lead replaces “issues indexes” in Task 4 and VAL-007's successor with an explicit pre/post `sqlite_schema` inventory rule. At minimum it must cover `issues_enriched`, both supporting indexes, all five triggers directly on `issues`, every retained trigger whose SQL references `issues`, and the PK/UNIQUE autoindexes verified structurally with PRAGMAs. |
| QA-BLK-18-03 | The exact benchmark command is valid, but the current benchmark scope cannot satisfy the phase-level performance spot-check for this migration: it has 1k/10k-event rebuild benchmarks only, with no 100k-event rebuild or 10k-issue registry-backed read benchmark. | `internal/journal/rebuild_bench_test.go` defines only `BenchmarkRebuildWorkspace1kEvents` and `BenchmarkRebuildWorkspace10kEvents`; no read/write benchmark matched repository search. The fresh 10k baseline is 3990 ms/rebuild against TRD-PERF-002's 5s target. | Tech Lead adds Task 5 test scope and acceptance for a 100k-event rebuild benchmark (<30s, TRD-PERF-003) and a 10k-issue registry-backed read spot-check (<500ms, TRD-PERF-001), or records an explicit scope decision with an equivalent reproducible performance control. The existing configurable `seedScaleReplayWorkspace` helper is a feasible starting fixture. |

No `ESC-NNN` is raised by this plan: these are pre-task QA refinement requests
within Tech Lead authority, not unresolved product or architecture decisions.
If the Tech Lead rejects a required split, exact object inventory, or
performance control, escalate before authorizing Task 1.

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
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
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
Task 1. The untracked Phase 18 run ledger is outside QA Lead write scope and was
not modified. Release-readiness scanners and held-out scenarios remain N/A:
AGENTS.md records no accepted end-to-end/held-out harness, and Phase 18 is not a
public-release sign-off. Ordinary Human dogfood remains planned in section 9.

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: **NO**.
  Task 5's issue-read/ready parity, dev-composite parity, and production-map
  isolation are uncovered; VAL-011 also does not yet prove the full
  PRD-OPS-003/004 projection claim.
- Every `VAL-NNN` has one validator role: **YES**; all current rows assign QA
  Lead.
- Every `VAL-NNN` has a concrete evidence method: **NO**. VAL-005's
  “structurally accepted/recognized/no behavior” boundary and VAL-007's
  undefined “issues indexes” inventory are not independently reproducible as
  written.
- Every `VAL-NNN` maps to a task/window or phase-end check: **NO**. VAL-005 is
  assigned to Task 1 although loader recognition first exists in Task 2 and
  phase-wide behavior absence cannot be closed until Task 5.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: **NO**. VAL-005, VAL-007, VAL-009, and VAL-013 bundle
  independent machine-checkable outcomes; Human fan-out cannot repair them.
- Conditional or binary assertions require qualifying and non-qualifying
  states: **NO**. VAL-006/009/010 do name both operational states, but VAL-012
  proves only stale-version behavior and not the current-version counterpart.
- Parity assertions have deterministic conformance proof: **PARTIAL**.
  VAL-002/003/004/014 are structurally comparable; VAL-008 says “every” while
  proposing representative evidence, and VAL-011 compares only workflow tables
  and issues rather than the full projection.
- Every current VAL-001..014 has a provisional deterministic `QA-NN` mapping:
  **YES**, in section 7. Provisional subchecks do not cure bundled contract
  assertions.
- One-line coverage rollup: **all in-scope requirements are referenced by at
  least one validation assertion: NO.**

### Required Tech Lead Contract Refinements

Proposed IDs after VAL-014 are recommendations; the Tech Lead owns numbering
and edits. If IDs change, update task specs and this plan together.

| Proposal | Required contract/task change | Reason |
|---|---|---|
| REF-18-01 | Narrow VAL-005 to one Task 1 observable: a canonical gate-obligation JSON value such as `{"gates":[{"id":"gate-test","version":1}]}` can be inserted, loaded, parsed, and round-tripped without losing ID/version. Add proposed VAL-015 for Task 2 loader recognition of `required_gates_current` and `has_completion` as known-but-unimplemented names; VAL-016 for Task 1 seed exclusion of both names; and VAL-017 for Task 5/phase-end proof that no production evaluator can satisfy or query receipt/gate state for either name. | Separates schema shape, vocabulary recognition, seed absence, and behavior absence; assigns each to the first honest window. |
| REF-18-02 | Narrow VAL-007 to exact `issues` row/value preservation. Add proposed VAL-018 for exact dependent-schema-object preservation and VAL-019 for byte-identical `fid issue get`/`list` output over a frozen v0.1 fixture before/after migration 0008. Amend Task 4 accordingly. | Row parity, schema-object parity, and CLI-output parity are independently falsifiable. |
| REF-18-03 | Narrow VAL-009 to “entering a started-category state requires the caller's live claim,” with missing/wrong/stale/live states. Add proposed VAL-020 for “leaving started auto-releases the caller's live claim in the same event/projection operation.” | Claim admission and automatic release are separate behaviors. |
| REF-18-04 | Narrow VAL-013 to golden v0.1 event replay producing the expected migrated projection. Add proposed VAL-021 for every journal file remaining byte-identical across migration, stale-version refresh, delete-and-rebuild, and a second rebuild. | Replay compatibility and file immutability are independent invariants. |
| REF-18-05 | Add proposed VAL-022 for the post-Task-5 `issue get/list` status-category plus `ready` result/order parity matrix; VAL-023 for `fid dev start/finish/block` outcome/effect parity; and VAL-024 for the static invariant that production files no longer read the frozen `definitions` map while tests retain it. Amend Task 5 coverage and acceptance. | These are explicit Task 5/phase-brief requirements with no current VAL owner. A single “remaining consumers” row would bundle distinct surfaces. |
| REF-18-06 | Widen VAL-011 to normalized snapshots of every projected table, including all four workflow tables and all pre-existing tables, after delete-and-rebuild and a second rebuild. Require equal `sqlite_schema`/FK integrity where schema state is part of the claim; document any intentionally normalized volatile field. | PRD-OPS-003/004 promise reconstruction of the query store, not only workflow tables and issues. Repeated parity proves determinism. |
| REF-18-07 | Replace VAL-008's “representative” evidence with an exhaustive structural matrix over every ordered pair of known states in all three frozen workflows, plus rule satisfied/blocked fixtures and exact CLI error/output characterization for allowed, invalid, blocked, and cancel categories. | “Every v0.1 transition” cannot pass from representative-only evidence. |
| REF-18-08 | Add the non-qualifying current-version case to VAL-012: the prior version triggers rebuild/catch-up, while an otherwise-fresh current-version projection is not rebuilt because of schema version. | Required both-state proof for a conditional assertion. |
| REF-18-09 | Clarify VAL-001's “PKs” as the keys specified by TRD_SCHEMA §3, explicitly stating whether `workflow_transition_requires` intentionally has no PK/UNIQUE constraint. | The draft DDL gives PKs to the first three tables but none to `workflow_transition_requires`; the current pass wording is ambiguous. |

## 3. Task Inventory

| Task | Focus | Current covered VAL IDs | Proposed successor coverage | Risk | QA notes |
|---|---|---|---|---|---|
| Task 1 | Workflow tables, exact seeds, gate seam, projection version | VAL-001..005, VAL-012 | VAL-016 | High | Exact DDL/seed set equality, migration down/up, no reserved seeds, and stale/current projection-version branches. |
| Task 2 | DB-backed registry and frozen-definition parity | VAL-014; advances VAL-002..004 | VAL-015 | High | Decoder must preserve role order/params, tolerate known-unimplemented vocabulary structurally, and perform no writes. |
| Task 3 | Transition-path source swap and claim parity | VAL-008, VAL-009 | VAL-020 | Critical | Exhaustive edge matrix and exact errors; Review 1 is mandatory immediately after this task. |
| Task 4 | SQLite `issues` table rebuild and FKs | VAL-006, VAL-007 | VAL-018, VAL-019 | Critical | Exact row/schema/CLI snapshots, insert and update FK probes, down/up, `foreign_key_check`; Review 2 is mandatory immediately after this task. |
| Task 5 | Remaining consumers, maintenance, rebuild/journal closure | VAL-010, VAL-011, VAL-013 | VAL-017, VAL-021..024 | Critical | Full projection parity, seeded healthy/corrupt audit states, all remaining consumer matrices, performance controls, and full phase gate. |

## 4. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Canonically ordered schema fingerprint from `sqlite_schema`, `PRAGMA table_info`, `index_list/index_info`, and `foreign_key_list` for all four workflow tables | VAL-001 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-02 | Deep-equal ordered `(workflow_id,state,category)` tuples against frozen v0.1 definitions for exactly three workflows | VAL-002 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-03 | Deep-equal transition tuple set for every known-state pair, including every explicit non-cancelled→cancelled edge and absence of disallowed/cancelled→cancelled edges | VAL-003 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-04 | Deep-equal rule tuple set, semantic `params_json` comparison, role-order/error-order preservation, and no extra/missing rule | VAL-004 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-05 | Gate-obligation JSON insert/load/parse/round-trip, reserved-name recognition, seed exclusion, and phase-end static/runtime evidence of no receipt/gate evaluator; recorded separately | VAL-005; proposed VAL-015..017 | Developer, QA Reviewer, QA Lead | Tasks 1, 2, 5 / reviews / final |
| EVID-06 | FK metadata plus raw INSERT and UPDATE probes for unknown workflow, known-workflow/invalid-state, and valid pair; clean `PRAGMA foreign_key_check` | VAL-006 | Developer, QA Reviewer, QA Lead | Task 4 / Review 2 / final review |
| EVID-07 | Canonical pre/post issue-row snapshot, exact dependent-object fingerprint, and frozen `issue get/list` output matrix, as three separately reported results | VAL-007; proposed VAL-018/019 | Developer, QA Reviewer, QA Lead | Task 4 / Review 2 / final review |
| EVID-08 | Exhaustive frozen-v0.1-vs-registry transition matrix for every known state pair and every rule success/failure branch, plus built-CLI exact output/error checks | VAL-008 | Developer, QA Reviewer, QA Lead | Task 3 / Review 1 / final review |
| EVID-09 | Missing/wrong/stale/live claim entry cases and separate same-event out-of-started claim release proof | VAL-009; proposed VAL-020 | Developer, QA Reviewer, QA Lead | Task 3 / Review 1 / final review |
| EVID-10 | Healthy migrated workspace with no workflow finding; test-only corrupt pair with deterministic finding; same violation under `--warn-as-error` with non-zero exit and unchanged finding payload | VAL-010 | Developer, QA Lead | Task 5 / final review |
| EVID-11 | Normalized per-table snapshots of the complete projection before deletion, after rebuild, and after second rebuild; equal workflow rows/issues; clean FK/schema evidence | VAL-011 | Developer, QA Lead | Task 5 / final review |
| EVID-12 | Prior-version projection triggers rebuild/catch-up to current; otherwise-fresh current-version projection remains fresh | VAL-012 | Developer, QA Reviewer, QA Lead | Task 1 / Review 1 / final review |
| EVID-13 | Golden v0.1 `issue.created`/`issue.transitioned` replay result and separate hashes/bytes for every journal file before/after migrate, refresh, rebuild, and second rebuild | VAL-013; proposed VAL-021 | Developer, QA Lead | Task 5 / final review |
| EVID-14 | Canonical registry-vs-frozen-definition deep equality across every state/category/edge/rule/params tuple for all three workflows | VAL-014 | Developer, QA Reviewer, QA Lead | Task 2 / Review 1 / final review |
| EVID-15 | Post-swap `issue get/list/ready` output matrix, dev start/finish/block effects, and production-source reference inventory, reported as separate evidence | proposed VAL-022..024 | Developer, QA Lead | Task 5 / final review |
| EVID-16 | Exact per-task command results, coverage package percentages, committed report-lint results, verified hashes/subjects, review verdicts, phase gates, and spec-sync disposition | no VAL | Developer, QA Reviewer, QA Lead | Every task/review/final review |
| EVID-17 | 10k-issue read and 10k/100k-event rebuild benchmark output against TRD-PERF-001..003, plus `issues_enriched` query-plan regression result | no VAL | Developer, QA Lead | Task 5 / final review |

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
| Benchmarks | `go test -bench . -benchmem ./...` | All benchmarks PASS; 10k rebuild <5s, 100k rebuild <30s, 10k-issue read <500ms once QA-BLK-18-03 is closed | Task 5 and phase end |
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
| `go test -bench . -benchmem ./...` | Exact command plus help and benchmark inventory | PASS for command; evidence scope BLOCKED by QA-BLK-18-03 | Exit 0. 1k rebuild: 458.0 ms; 10k rebuild: 3990 ms. Help confirms `-bench`/`-benchmem`; repository inventory confirms no 100k/read benchmark yet. |
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
| Exact Task 4 dependent-object set | Migration inventory | BLOCKED by QA-BLK-18-02 | Existing objects are discoverable, but “issues indexes” is not an exact set. Five triggers are directly on `issues`; two named indexes support the view on other tables; other retained triggers reference `issues`. |
| Pre/post migration row and CLI output fixture | Target-version migrator plus existing CLI migration/rebuild helpers | PASS once REF-18-02 is applied | Existing helpers can create a version-7 DB, snapshot rows/output, migrate to version 8, and compare. Task 4 must commit the frozen fixture/expected output; a same-binary auto-migrating “before” call without version control is insufficient. |
| Invalid/valid FK pairs | Existing raw-SQL fixture approach | PASS | Maintenance tests already disable FKs/checks only in test setup and restore them, while normal DB opens enforce FKs. Task 4 can test ordinary rejection; Task 5 can seed corruption through this test-only bypass. |
| Exhaustive transition and claim parity | Existing tests and integration test listing | PASS | Current tests cover allowed/invalid/rule/cancel and into/out-of-started claim behavior. Frozen definitions provide finite state/rule sets for exhaustive generated matrices. |
| Maintenance healthy/corrupt and `--warn-as-error` | Existing fixtures, tests, and built CLI help | PASS | `newValidStructuralFixture`/`newStructuralFixture` provide both states; `fid maintenance check --help` confirms `--warn-as-error` and `--format`. |
| Full projection and journal-byte parity | Existing closure/rebuild helpers | PASS | `closureProjectedTableSnapshots`, `assertClosureProjectedTablesEqual`, `assertClosureJournalBytesEqual`, and repeated rebuild scale tests exist. Task 5 must add the four workflow tables and any omitted projected tables to the snapshot inventory. |
| Projection-version stale/current branches | Existing test harness | PASS | `TestEnsureProjectionFreshRebuildsStoredPriorProjectionSchemaVersion` and freshness-decision tests exist; Task 1 can extend them from version 3→4 and add the fresh counterpart. |
| No production reads of frozen map | `rg -n 'definitions' internal cmd --glob '*.go'` and `rg -n 'workflow\.' internal cmd --glob '*.go' --glob '!**/*_test.go'` | PASS as reproducible future checks | Current definition/production call sites are enumerated in `internal/workflow`, `internal/issues`, `internal/dev`, and `internal/maintenance`; after Task 5, the same inventory must return only approved fixture/test references. Source-state proof is paired with runtime matrices. |
| CLI steps for Human dogfood | Built binary help | PASS | `issue transition/get/list`, `ready`, `maintenance check`, and `rebuild` help all exit 0 with the required flags. Human steps require a disposable copy, never the sole real workspace. |
| 10k/100k performance fixture | Benchmark/helper inspection | BLOCKED by QA-BLK-18-03 | `seedScaleReplayWorkspace` is configurable and therefore feasible, but no 100k benchmark or 10k-issue read benchmark is currently registered for the exact benchmark command. |
| Completion reports and reviews | Task specs, template/script usage, and file checks | PASS as future artifacts | `task-1-done.md` through `task-5-done.md` are intentionally absent before work. Each task spec declares its path. Reviews are feasible from the committed reports and template after Tasks 3 and 4. |
| Held-out/Mode B scenario | AGENTS.md and Quality Plan §7 | N/A | No accepted held-out location/command exists; Phase 18 is not release readiness. No hidden scenario material was read. |

## 7. Deterministic Checklist

Every item is executed by the QA Lead at phase end, after the relevant
Developer/Reviewer evidence is available. Literal command output or exact
file/DB state is required; code reading alone cannot produce a PASS. Proposed
VAL IDs are inactive until the Tech Lead updates the contract.

| ID | Check | Covers | Evidence required |
|---|---|---|---|
| QA-01 | Compare the exact four-table schema fingerprint with TRD_SCHEMA §3, including clarified PK/UNIQUE/FK expectations `(qa-lead)` | VAL-001 | Test output plus canonical PRAGMA/`sqlite_schema` rows. |
| QA-02 | Compare the complete ordered state/category tuple set for all three defaults with the frozen v0.1 reference `(qa-lead)` | VAL-002 | Deep-equality test output showing no extra/missing tuple. |
| QA-03 | Compare every known-state edge pair, including all cancel and absent-edge states, with the frozen v0.1 reference `(qa-lead)` | VAL-003 | Exhaustive matrix/set-diff output, empty diff. |
| QA-04 | Compare every edge rule and parsed params value with the frozen v0.1 reference `(qa-lead)` | VAL-004 | Deep-equality output with rule/role/params tuple counts. |
| QA-05 | Round-trip the canonical gate-obligation JSON without loss `(qa-lead)` | VAL-005 (proposed narrowed) | Test output and loaded canonical JSON/typed value. |
| QA-06 | Prove both reserved rule names are recognized as known-but-unimplemented by the loader `(qa-lead)` | proposed VAL-015 | Test output for both names; no panic/write. |
| QA-07 | Prove no seed row uses either reserved rule name `(qa-lead)` | proposed VAL-016 | Ordered SQL query output with zero rows. |
| QA-08 | Prove no production evaluator can satisfy or query receipt/gate state for either reserved rule `(qa-lead)` | proposed VAL-017 | Runtime injected-rule outcome plus static production reference/call inventory. |
| QA-09 | Exercise INSERT and UPDATE with unknown workflow, invalid known pair, and valid pair; run `foreign_key_check` `(qa-lead)` | VAL-006 | Rejection/acceptance errors and zero-row PRAGMA output. |
| QA-10 | Compare every `issues` row/column before and after migration 0008 `(qa-lead)` | VAL-007 (proposed narrowed) | Canonical row snapshots, counts, and empty diff. |
| QA-11 | Compare the exact dependent schema-object fingerprint before/after migration and after down/up `(qa-lead)` | proposed VAL-018 | Named `sqlite_schema`/PRAGMA snapshots and empty diff. |
| QA-12 | Compare frozen `fid issue get` and `list` outputs before/after migration 0008 `(qa-lead)` | proposed VAL-019 | Literal pre/post outputs or hashes with byte-equality result. |
| QA-13 | Execute the exhaustive registry-backed transition matrix and exact built-CLI parity cases `(qa-lead)` | VAL-008 | Matrix counts, empty diff, and literal allowed/invalid/blocked/cancel output excerpts. |
| QA-14 | Exercise missing, wrong-holder, stale, and live claims when entering started `(qa-lead)` | VAL-009 (proposed narrowed) | Exact errors and accepted transition/DB state. |
| QA-15 | Prove leaving started releases the live caller claim atomically `(qa-lead)` | proposed VAL-020 | Event effect plus post-projection and rebuilt claim absence. |
| QA-16 | Run maintenance against healthy and test-only corrupt workflow pairs, then rerun corrupt state with `--warn-as-error` `(qa-lead)` | VAL-010 | Literal JSON output, ordered finding, exit statuses, unchanged journal bytes. |
| QA-17 | Delete and rebuild the full projection, compare every table, rebuild again, and compare again `(qa-lead)` | VAL-011 | Per-table snapshot hashes/diffs, schema/FK state, and both rebuild outputs. |
| QA-18 | Exercise prior and current projection-schema versions `(qa-lead)` | VAL-012 | Rebuild/catch-up proof for prior and no version-triggered rebuild for current. |
| QA-19 | Replay the frozen v0.1 journal into the Phase 18 schema and compare the expected projection `(qa-lead)` | VAL-013 (proposed narrowed) | Golden replay output and projection snapshot diff. |
| QA-20 | Compare every journal file before/after migrate, refresh, rebuild, and second rebuild `(qa-lead)` | proposed VAL-021 | File list, byte counts/hashes, and byte-equality output at each boundary. |
| QA-21 | Compare the DB-loaded registry with the frozen definitions for all states/categories/edges/rules/params `(qa-lead)` | VAL-014 | Canonical deep-equality output for all three workflows. |
| QA-22 | Compare post-swap issue get/list status-category output and ready selection/order with frozen v0.1 expectations `(qa-lead)` | proposed VAL-022 | Matrix output and literal CLI excerpts. |
| QA-23 | Compare `fid dev start`, `finish`, and `block` results/effects with frozen v0.1 expectations `(qa-lead)` | proposed VAL-023 | Built-CLI/integration output plus event/projection effect snapshots. |
| QA-24 | Verify no production file reads the frozen definitions map after Task 5 `(qa-lead)` | proposed VAL-024 | Output from the two exact `rg` inventory commands in section 6, limited to approved fixture/test references, paired with QA-22/23 runtime evidence. |
| QA-25 | Run every exact phase gate and enforce COV-001/COV-002 `(qa-lead)` | no VAL | Literal command outputs and package percentages. |
| QA-26 | Verify migrations 0007 and 0008 down/up cleanly restore their documented prior/current states `(qa-lead)` | no VAL | Migration test output, schema/row snapshots, and clean FK check. |
| QA-27 | Review registry/maintenance SQL for parameter binding, dynamic identifier confinement, test-only corruption bypass, and no journal/data-loss path `(qa-lead)` | no VAL | Test output plus reproducible query/reference inventory and findings. |
| QA-28 | Run relevant read/rebuild benchmarks and the `issues_enriched` query-plan test against TRD-PERF-001..004 `(qa-lead)` | no VAL | Benchmark durations, allocations, fixture sizes, and EXPLAIN index evidence; blocked until QA-BLK-18-03 closes. |
| QA-29 | Verify all five committed completion reports against exact task commands and verified commit metadata `(qa-lead)` | no VAL | Five report-lint PASS outputs plus hashes/subjects copied from Git. |
| QA-30 | Verify Tech Lead spec sync records `touched`, `verified-no-change`, or `N/A` for the full canonical spec suite and updates the known TRD §4.6/Schema §14.11 wording `(qa-lead)` | no VAL | Committed spec/decision diff and wrap-up sync table. |

## 8. Mid-Phase Review Cadence

The Heavy Mode cadence is stricter than the default after-every-three-tasks
rule because Tasks 3 and 4 are independent critical swaps.

- **Review 1 after Task 3, before Task 4:** write and commit
  `qa/review-1.md`, covering Tasks 1–3 and current VAL-001..005, VAL-008/009,
  VAL-012/014 plus their accepted successor IDs. Independently rerun at least
  one representative required command, run report-lint for all three committed
  reports, and verify the exhaustive transition/claim evidence. `VERDICT:
  BLOCK` stops Task 4.
- **Review 2 after Task 4, before Task 5:** write and commit
  `qa/review-2.md`, covering Task 4 and VAL-006/007 plus accepted successor
  IDs. Independently rerun the migration/FK/object-parity scope and report-lint
  Task 4. Any row loss, undefined/missing object, FK anomaly, or narrower
  command evidence is Blocking. `VERDICT: BLOCK` stops Task 5.
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

Task execution cannot start until QA-BLK-18-01..03 are resolved or formally
escalated, the Tech Lead-applied contract is Ready, affected task specs are
aligned, this plan's IDs/mappings are refreshed, and the Human approves the
Heavy Mode phase plan.

After implementation, the phase can proceed to final-review PASS only when:

- All five task implementations and finalized completion reports are committed;
  every task's exact command is recorded `PASS`, or the phase remains stopped
  on an accepted blocker report.
- Review 1 after Task 3 and Review 2 after Task 4 are committed with no
  unresolved Blocking finding and `VERDICT: CONTINUE`.
- Every executable deterministic gate in section 5 passes with literal output;
  coverage thresholds and applicable TRD performance targets are met.
- Every active validation assertion, including Tech Lead-added successor IDs,
  has a final `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD` outcome with cited
  evidence. `PASS` without evidence is invalid.
- Every `(qa-lead)` item has reproducible execution/file/DB evidence; source
  reading alone is not used to pass behavior.
- Human checks remain human-judgment-only, are copied to the Human checklist,
  and are accepted or converted into follow-up work.
- No unresolved Blocking defect, data-loss/journal-integrity concern, spec
  deviation, command anomaly, or required escalation remains.
- Tech Lead records full canonical spec-suite sync, including the known
  `02_TRD.md` §4.6 and `03_TRD_SCHEMA.md` §14.11 source-of-truth updates, and
  records any migration/rollback documentation required for dogfood.
