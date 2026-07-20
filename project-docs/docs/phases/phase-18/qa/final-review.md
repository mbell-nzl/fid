# Phase 18 Final Review - Workflow Tables Migration + Engine Data-Source Swap + Gate Seams

Reviewer: QA Lead
Date: 2026-07-20

## Verdict

**BLOCK**

Phase 18 cannot close. The final review independently reproduced a deterministic
deadlock in a required Task 5 maintenance package fixture. The exact coverage,
unit, integration, and benchmark commands do not complete on HEAD even though
the Task 5 completion report records them as PASS. A disposable real-entrypoint
check independently proves VAL-010 behavior, but the deadlock still invalidates
the seven-command gate as a complete passing gate.

The DEC-071 performance remediation itself is effective: both the Task 5 run
and this final review keep the 10k-read and 1k/10k/100k rebuild rows below their
unchanged TRD-PERF targets. The block is the registry loader's connection-pool
deadlock, not a performance-target regression.

A second phase-exit prerequisite remains open: QA-30's Tech Lead spec sync has
not yet been recorded, and the known `02_TRD.md` / `03_TRD_SCHEMA.md` wording is
still stale. No concrete `ESC-NNN` is open.

## Scope And Evidence

- Branch/HEAD: `phase-18-workflow-migration` at
  `d7f11995c77d8ed028b4ecaaa2cbd61e8688522c`.
- Phase brief: `project-docs/docs/current/phase_brief.md`
- Decisions: `project-docs/docs/phases/phase-18/decisions.md`, including
  DEC-067 through DEC-071.
- Validation contract:
  `project-docs/docs/phases/phase-18/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-18/qa/qa-plan.md`
- Task specs and completion reports:
  `project-docs/docs/phases/phase-18/tasks/`
- Mid-phase reviews and remediation reviews:
  `project-docs/docs/phases/phase-18/qa/review-1.md`,
  `review-1-remediation.md`, `review-2.md`, and
  `review-2-remediation.md`.
- Run ledger: `project-docs/docs/phases/phase-18/run-ledger.md` (read-only in
  this review).
- Current decisions/progress/escalations and Quality Plan / workflow authority
  loaded per the QA Lead role.

### Reproducible Evidence Index

| Evidence ID | Result | Reproducible evidence |
|---|---|---|
| FR-E01 | PASS for the completed packages | The non-cached targeted race suite printed PASS for the Phase 18 DB, workflow, issues, and journal tests. Literal package tails: `ok .../internal/db 2.685s`, `ok .../internal/workflow 1.324s`, `ok .../internal/issues 7.375s`, `ok .../internal/journal 1.739s`. Named PASS tests include migration 0007/0008 schema, seed, FK, down/up, projection freshness, registry parity, exhaustive transition/claim checks, golden replay, rebuild parity, and DEC-071 pragmas. |
| FR-E02 | PASS | `go test -tags integration -race -count=1 -v ./cmd/fid -run '^TestV01ClosureJourneyAndRebuildParity$'` printed `--- PASS ... (6.32s)` and `ok .../cmd/fid 7.339s`; this covers complete projection/schema/FK comparison, repeated rebuilds, journal bytes, and the closure CLI journey. |
| FR-E03 | FAIL | `go test -race -count=1 -v -timeout 5s ./internal/maintenance -run '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$'` exited 1: `panic: test timed out after 5s`. Stack: `database/sql.(*DB).conn` -> `workflow.loadStates` at `internal/workflow/registry.go:112` -> `workflow.Load` at line 43 -> `maintenance.checkWorkflowStates` at `internal/maintenance/check.go:208`. The fixture sets `db.SetMaxOpenConns(1)` at `internal/maintenance/test_helpers_test.go:54`; `Load` holds the outer `workflows` rows while issuing nested queries. |
| FR-E04 | FAIL | Exact final gates contradicted the Task 5 report: `go test -race -coverprofile=coverage.out ./...` was bounded after 45s and exited 124; `go test -race ./...` exited 124 after 30s; `go test -tags integration -race ./...` remained stuck after 50s and was interrupted (130); the benchmark command produced every required row, then remained stuck and was interrupted (130). All stop at the same maintenance test path proven by FR-E03. |
| FR-E05 | PASS | `gofmt -l cmd internal` exited 0 with no output; `golangci-lint run ./...` exited 0 with `0 issues.`; `go build -o /tmp/fid ./cmd/fid` exited 0. |
| FR-E06 | PASS targets / FAIL command completion | Final-review benchmark rows: `BenchmarkList10kIssuesWithRegistry-6 46792809 ns/op 46.79 ms/registry-read`; rebuild 1k `294919945 ns/op 294.9 ms`, 10k `2779543050 ns/op 2780 ms`, 100k `26382594727 ns/op 26383 ms`. All meet `<500ms`, `<5s`, and `<30s`; the exact all-package command did not exit because of FR-E03. |
| FR-E07 | PASS | Exact VAL-024 inventories: `rg -n 'definitions' internal cmd --glob '*.go'` found the frozen fixture and tests plus the registry documentation comment; no runtime consumer caller reads the fixture. `rg -n 'workflow\.' internal cmd --glob '*.go' --glob '!**/*_test.go'` found DB-backed `workflow.Load` consumers and rule types only. Runtime parity: FR-E01, FR-E02, and the independent dev test output in FR-E08. |
| FR-E08 | PASS | `go test -race -count=1 -v ./internal/dev -run '^(TestStartRejectsEdgeMissingFromProjectedWorkflowRegistry|TestFinishRejectsEdgeMissingFromProjectedWorkflowRegistry|TestBlockRejectsEdgeMissingFromProjectedWorkflowRegistry|TestGateUsesProjectedWorkflowCategory)$'` printed all four `--- PASS` rows and `ok .../internal/dev 2.708s`. |
| FR-E09 | PASS | Five `project-docs/tools/report-lint.sh` invocations, one for each concrete `task-1` through `task-5` report/spec pair, passed. They verified every exact command record, report headings, resolving implementation hash/subject, and committed/clean report. Git verified Task hashes/subjects: `127f268 feat(workflow): add seeded workflow tables`; `799acf0 feat(workflow): load workflow registry from database`; `77400ce test(issues): prove released claims stay absent after rebuild`; `a27fb5f test(db): cover valid issue workflow update`; `37b3fb5 perf(rebuild): relax projection durability`. |
| FR-E10 | PASS | DEC-071 focused test `TestOpenForRebuildRelaxesDurabilityWithoutChangingNormalOpen` passed. Task 5 records rebuild values `foreign_keys=1`, `recursive_triggers=1`, `synchronous=0`, `journal_mode=memory`, while normal `Open` remains `foreign_keys=1`, `recursive_triggers=1`, `synchronous=2`, `journal_mode=delete`. Source inventory shows `rebuildWorkspaceLocked` calls `OpenForRebuild`; DSN pragmas are in `internal/db/db.go`, not a pooled `Exec`. |
| FR-E11 | PASS thresholds / stale aggregate | Current independently completed coverage results: `cmd/fid 77.7%`, `internal/issues 81.5%`, `internal/claims 81.4%`, `internal/workflow 83.3%`, `internal/policy 86.3%`, `internal/journal 81.0%`. All meet COV-001/COV-002. `go tool cover -func=coverage.out` exits 0 and reports `total ... 76.7%`, but that retained profile was not freshly produced by the failed final-review coverage run and does not make FR-E04 pass. |
| FR-E12 | PASS | `TestIssuesEnrichedQueryPlanUsesConcernAndTriageIndexes` passed in FR-E01, proving the `idx_concerns_issue_status` and `idx_triage_issue_created` query-plan controls. Static SQL review found fixed SQL or bound `?` arguments for registry/maintenance reads; dynamic issue list/ready SQL returns a separate argument slice. DSN construction uses `url.Values`; no `go.mod`/`go.sum` change exists. |
| FR-E13 | BLOCKED | No Phase 18 canonical spec-sync diff or wrap sync table exists. `02_TRD.md` §4.6 still labels workflows hard-coded and `03_TRD_SCHEMA.md` §14.11 still says maintenance validates against the Go definitions. The phase brief/current progress also still describe Task 5 as uncommitted and the brief exit checklist stops at VAL-014. |
| FR-E14 | PASS | Disposable real-entrypoint VAL-010 check in `/tmp/fid-phase18-final-maint.kLl2V7`: initialized a workspace, created task `P18-c354c2`, deleted projected `task-default/todo` with FKs disabled only for fixture setup, then ran `/tmp/fid maintenance check --warn-as-error --format json`. Exit was 1; literal findings were `foreign_key_integrity` and `workflow_state` naming `P18-c354c2`; journal SHA-256 before/after was identical: `bfe82e5f...ce149`. Paired non-cached CLI tests `TestMaintenanceCheckCommandRendersHealthyAndFailingResults` and `TestMaintenanceCheckCommandPromotesWarningsWithoutRewritingThem` both PASS (`ok .../cmd/fid 1.784s`). |

## Task 5 Seven-Command Gate Audit

Task 5's report is committed and passes mechanical report-lint, but substantive
final-review reproduction contradicts four of its seven PASS records.

| # | Exact Task 5 command | Task 5 report | Final-review result | Evidence |
|---:|---|---|---|---|
| 1 | `gofmt -l cmd internal` | PASS, no output | PASS | FR-E05. |
| 2 | `golangci-lint run ./...` | PASS, `0 issues.` | PASS | FR-E05. |
| 3 | `go test -race -coverprofile=coverage.out ./...` | PASS | **FAIL** | Exit 124 after 45s; other package coverage printed, but maintenance never completed (FR-E03/04). |
| 4 | `go tool cover -func=coverage.out` | PASS, total `74.1%` | BLOCKED as a fresh gate | Command exits 0 on retained `coverage.out` with total `76.7%`, but command 3 did not produce a fresh complete profile (FR-E11). |
| 5 | `go test -tags integration -race ./...` | PASS | **FAIL** | Exact command did not complete after 50s and was interrupted; same deadlock (FR-E03/04). |
| 6 | `go test -bench . -benchmem ./...` | PASS with four rows | **FAIL as a complete command; performance rows PASS** | Final rows all meet targets, but the all-package command did not exit because package tests include the deadlocked maintenance test (FR-E06). |
| 7 | `go build -o /tmp/fid ./cmd/fid` | PASS | PASS | FR-E05. |

Task 5's recorded benchmark rows remain valid historical evidence:

| Control | Task 5 report | Final review | Target | Target result |
|---|---:|---:|---:|---|
| 10k-issue registry read | `46807508 ns/op` (`46.81 ms`) | `46792809 ns/op` (`46.79 ms`) | `<500 ms` | PASS |
| 1k-event rebuild diagnostic | `278769882 ns/op` (`278.8 ms`) | `294919945 ns/op` (`294.9 ms`) | diagnostic | PASS |
| 10k-event rebuild | `2647464126 ns/op` (`2647 ms`) | `2779543050 ns/op` (`2780 ms`) | `<5 s` | PASS |
| 100k-event rebuild | `25235973633 ns/op` (`25236 ms`) | `26382594727 ns/op` (`26383 ms`) | `<30 s` | PASS |

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | Exit 0; no files listed (FR-E05). |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Exit 0; `0 issues.` (FR-E05). |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | FAIL | Timed out; FR-E03/04. |
| Coverage report | `go tool cover -func=coverage.out` | BLOCKED | Existing profile parses at 76.7%, but no fresh complete profile was produced; FR-E11. |
| Unit | `go test -race ./...` | FAIL | Timed out; FR-E03/04. |
| Integration | `go test -tags integration -race ./...` | FAIL | Hung on the same maintenance path and was interrupted; FR-E03/04. |
| End-to-end | N/A | N/A | No accepted E2E/held-out command exists; focused closure integration passed in FR-E02. |
| Benchmarks | `go test -bench . -benchmem ./...` | FAIL | All required rows meet targets, but exact command did not exit; FR-E06. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Exit 0; FR-E05. |
| Full phase check | Seven Task 5 commands plus exact unit suite | FAIL | Coverage, unit, integration, and benchmark commands do not complete. |

## Validation Contract Outcomes

Every active VAL-001 through VAL-024 is recorded below. `PASS` rows cite
reproducible evidence. VAL-010 uses the real-entrypoint proof rather than the
deadlocked package fixture.

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | FR-E01: `TestMigration0007CreatesWorkflowSchema` and all four table subtests PASS; Task 1/Review 1 schema evidence. |
| VAL-002 | PASS | QA Lead | FR-E01: `TestMigration0007SeedsFrozenV01WorkflowDefinitions`, `TestRegistryEqualsV01Definitions`, and all three TRD fixture subtests PASS. |
| VAL-003 | PASS | QA Lead | FR-E01 exhaustive seed/registry and 122-pair transition matrix PASS, including explicit cancel and absent edges. |
| VAL-004 | PASS | QA Lead | FR-E01 seed and registry deep-equality tests plus ruled-edge satisfied/blocked tests PASS. |
| VAL-005 | PASS | QA Lead | FR-E01 `TestMigration0007ReservesGateObligationShapeWithoutSeedingReservedRules` PASS; Task 1 report records canonical JSON loss-free round-trip. |
| VAL-006 | PASS | QA Lead | FR-E01 `TestMigration0008EnforcesWorkflowAndStateForeignKeys` PASS; Review 2 remediation proves rejected invalid INSERT/UPDATE, accepted valid INSERT/UPDATE with read-back, and clean FK state. |
| VAL-007 | PASS | QA Lead | FR-E01 `TestMigration0008PreservesIssuesAndDependentObjectsThroughDownUp` PASS with canonical row snapshots. |
| VAL-008 | PASS | QA Lead | FR-E01 exhaustive 122-pair matrix and all ruled-edge branch tests PASS; Task 3 report records literal CLI allowed/invalid/blocked/cancel text. |
| VAL-009 | PASS | QA Lead | FR-E01 missing, wrong-holder, stale, and live claim tests PASS, including DB-category source mutation. |
| VAL-010 | PASS | QA Lead | FR-E14 proves healthy/corrupt CLI states, deterministic workflow-pair finding, non-zero `--warn-as-error`, and unchanged journal bytes. FR-E03 remains a separate gate/test-harness failure. |
| VAL-011 | PASS | QA Lead | FR-E02 repeated full-projection closure/rebuild parity PASS; Task 5 report names every projected-table/schema/FK comparison. |
| VAL-012 | PASS | QA Lead | FR-E01 prior/current freshness decision and real prior-version rebuild tests PASS. |
| VAL-013 | PASS | QA Lead | FR-E01 `TestReplayJournalAppliesGoldenV1JournalInFileOrder` PASS; expected migrated projection asserted. |
| VAL-014 | PASS | QA Lead | FR-E01 `TestRegistryEqualsV01Definitions` PASS for all states/categories/edges/rules/params across three workflows. |
| VAL-015 | PASS | QA Lead | FR-E01 `TestLoadSurfacesReservedRulesAsUnimplementedWithoutWriting` PASS under SQLite `query_only` for both reserved names. |
| VAL-016 | PASS | QA Lead | FR-E01 migration 0007 gate-seam test PASS; zero seeded reserved-name rows asserted. |
| VAL-017 | PASS | QA Lead | FR-E01 `TestTransitionReservedGateRuleRemainsUnevaluated` PASS; production reserved-name inventory contains only three rule constants and no evaluator/query. |
| VAL-018 | PASS | QA Lead | FR-E01 migration 0008 dependent-object/down-up test and all retained trigger subtests PASS; query-plan control also PASS. |
| VAL-019 | PASS | QA Lead | FR-E01 `TestMigration0008PreservesIssueGetAndListRenderedBytes` PASS; Task 4 records real CLI pre/post rebuild byte parity. |
| VAL-020 | PASS | QA Lead | FR-E01 exact DEC-068 rebuilt released-claim test and same-event release-effect test PASS; Review 1 remediation independently passed. |
| VAL-021 | PASS | QA Lead | FR-E02 closure test PASS with per-file journal-byte equality at migration, stale refresh, rebuild, and second rebuild boundaries. |
| VAL-022 | PASS | QA Lead | FR-E01 get/list/ready projected-category test plus FR-E02 closure CLI journey PASS. |
| VAL-023 | PASS | QA Lead | FR-E08 projected-registry start/finish/block/gate tests and FR-E02 closure journey PASS. |
| VAL-024 | PASS | QA Lead | FR-E07 exact two-`rg` inventories show only approved frozen fixture/test map references and DB-backed production call sites; paired FR-E01/02/08 runtime parity PASS. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| TRD-ARCH-009; TRD_SCHEMA §§2-3 | MET | VAL-001..007, VAL-014..016, VAL-018..019 PASS. |
| PRD-WF-001–010; TRD §4.6 behavior parity | MET | VAL-002..004, VAL-008, VAL-014, VAL-022..024 PASS. |
| PRD-WF-014/015 seam; DEC-066 | MET | VAL-005, VAL-015..017 PASS; no receipt behavior found. |
| PRD-CLAIM-005/006 | MET | VAL-009/020 PASS, including DEC-068 remediation. |
| PRD-OPS-003/004; TRD-ARCH-001..003 | MET | VAL-011..013/021 PASS; FR-E02 and FR-E10. |
| PRD-OPS-005–008; TRD_SCHEMA §14.11 | MET behavior / gate defect open | VAL-010 PASS via FR-E14; package-level fixture deadlock remains Blocking under QA-25. |
| TRD-PERF-001/002/003 | MET | FR-E06 and Task 5 rows all below target after DEC-071. |
| TRD-PERF-004 / TRD-ARCH-008 | MET | `TestIssuesEnrichedQueryPlanUsesConcernAndTriageIndexes` PASS (FR-E12). |
| COV-001/COV-002 | MET percentages; gate blocked | FR-E11 percentages exceed thresholds, but the full coverage command fails to complete. |
| Heavy Mode phase gate / Quality Plan §5.2 | **NOT MET** | FR-E04; four required commands do not complete. |

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS | FR-E01 migration 0007 schema test and table subtests. |
| QA-02 | PASS | FR-E01 seed/registry frozen-state parity tests. |
| QA-03 | PASS | FR-E01 seed/registry edge parity and 122-pair matrix. |
| QA-04 | PASS | FR-E01 rule/params parity and branch tests. |
| QA-05 | PASS | FR-E01 canonical gate JSON round-trip test. |
| QA-06 | PASS | FR-E01 both reserved names loaded under `query_only`. |
| QA-07 | PASS | FR-E01 zero reserved seed rows. |
| QA-08 | PASS | FR-E01 runtime unimplemented-rule result plus reserved-name inventory. |
| QA-09 | PASS | FR-E01 and Review 2 remediation complete INSERT/UPDATE matrix. |
| QA-10 | PASS | FR-E01 canonical issues row preservation. |
| QA-11 | PASS | FR-E01 complete dependent-object/down-up fingerprint and retained trigger probes. |
| QA-12 | PASS | FR-E01 rendered get/list byte parity. |
| QA-13 | PASS | FR-E01 exhaustive transition/rule matrix; Task 3 literal CLI report evidence. |
| QA-14 | PASS | FR-E01 four-state started-entry claim matrix. |
| QA-15 | PASS | FR-E01 same-event release and exact released-fixture rebuild; DEC-068 review PASS. |
| QA-16 | PASS | FR-E14 literal healthy/corrupt CLI evidence, `--warn-as-error` exit 1, ordered finding payload, and identical journal hash. |
| QA-17 | PASS | FR-E02 complete projection, schema, FK, and two-rebuild closure test. |
| QA-18 | PASS | FR-E01 stale/current projection-schema branches. |
| QA-19 | PASS | FR-E01 golden v0.1 replay test. |
| QA-20 | PASS | FR-E02 four-boundary journal-byte parity. |
| QA-21 | PASS | FR-E01 registry-vs-frozen-definition equality. |
| QA-22 | PASS | FR-E01 get/list/ready parity plus FR-E02 closure CLI. |
| QA-23 | PASS | FR-E08 dev composites plus FR-E02 closure CLI. |
| QA-24 | PASS | FR-E07 inventories paired with FR-E01/02/08 runtime parity. |
| QA-25 | **FAIL** | FR-E04: coverage, unit, integration, and benchmark commands do not complete. |
| QA-26 | PASS | FR-E01 migration 0007 and 0008 down/up tests, schema/row checks, clean FK evidence. |
| QA-27 | PASS | FR-E02/10/12: no journal-byte change, DSN-safe pragmas, bound/static SQL, no dependency addition; no injection/data-loss finding beyond the availability deadlock. |
| QA-28 | PASS | FR-E06 all four rows within controls; FR-E12 query-plan test PASS. Exact benchmark command completion remains failed under QA-25. |
| QA-29 | PASS | FR-E09: five report-lint passes and Git-verified hashes/subjects. |
| QA-30 | **BLOCKED** | FR-E13: Tech Lead canonical spec/current-doc sync is not yet recorded. |

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: **YES**,
  as recorded in `qa-plan.md` §6.
- Required command placeholders: **none**.
- Invalid command syntax: **none**.
- Final-review feasibility defect found: **yes**. The commands are valid, but
  implementation commit `37b3fb5` leaves a required test path deadlocked under
  the documented one-connection fixture, making four exact all-package commands
  non-terminating. Actual CLI maintenance behavior is separately green in
  FR-E14. This is a product/test regression, not a QA-plan syntax defect.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC References |
|---|---|---|---|---|
| `review-1.md` | BLOCK | YES | Task 2 report-finalization commit `7edfff6` lacks `Refs: phase-18/task-2`; DEC-068 carries the process reminder into CF-023. | none |
| `review-1-remediation.md` | CONTINUE | YES | DEC-068 test `TestTransitionOutOfStartedReleasedCallerClaimRemainsAbsentAfterRebuild` and commit `77400ce` passed independently; VAL-020 PASS. | none |
| `review-2.md` | BLOCK | YES | No improvement finding. | none |
| `review-2-remediation.md` | CONTINUE | YES | DEC-069 valid UPDATE/read-back commit `a27fb5f` passed independently; VAL-006 PASS. | none |

DEC-070's incomplete benchmark/build gate and DEC-071's genuine 33.3846063s
100k failure were both resolved without weakening targets. DEC-071's rebuilt
connection behavior and final performance rows are verified in FR-E06/10.

The run ledger contains out-of-chronological-order and duplicated entries noted
by DEC-070 and tracked by CF-023. It was not modified. Those hygiene defects do
not create an `ESC-NNN` or alter the evidence verdict.

## Documentation Review

- Completion reports present, committed, clean, and report-lint valid for all
  five tasks: **YES** (FR-E09).
- Impacted canonical/current docs fully updated: **NO**.
- Tech Lead wrap-up sync required before PASS:
  - `project-docs/specs/02_TRD.md` §§4.3/4.6: migrated DB-backed source of truth
    and rebuild-scoped durability note.
  - `project-docs/specs/03_TRD_SCHEMA.md` §14.11: maintenance now audits against
    projected workflow tables, not the canonical Go map.
  - Phase 18 spec-sync table for the full canonical suite.
  - Current phase brief/progress status and the phase-brief exit list
    (`VAL-001..024`, not `VAL-001..014`).
- No task spec, decision, run-ledger, or product-code file was modified by this
  review.

## Human Checklist Readiness

Generated at
`project-docs/docs/phases/phase-18/qa/human-checklist.md`.

- Contains only `(human)` items: **YES**.
- Uses one observable per row: **YES**.
- Human-only `VAL-NNN` assertions mapped: **N/A**; every Phase 18 VAL is QA
  Lead-owned and deterministic.
- Items: transition feel, get/list/ready presentation, maintenance clarity, and
  Phase 18 dogfood acceptance.
- Execution recommendation: wait until the Blocking maintenance/gate defect is
  remediated and final deterministic QA is rerun.

## Coverage Caveat

A deterministic PASS proves only the asserted mechanics, data states, command
outputs, and structural conformance checks. It does not prove unasserted
appearance, workflow feel, product-owner judgment, or dogfooding quality; those
remain unconfirmed until the Human checklist records outcomes. In this review,
the deterministic phase gate itself is also not green.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | BLOCKED | Final review/checklist exist and all VALs have outcomes, but the phase gate fails; Human checklist remains blank. |
| Automated gates | BLOCKED | FR-E04; format/lint/build and performance targets pass, but test-bearing commands do not complete. |
| Held-out scenarios | DEFERRED | No accepted harness/location; Phase 18 is not release readiness. |
| Release readiness | DEFERRED | Non-release phase; threat model/scans/scenarios are not substituted here. |

## Findings

### Blocking

- **QA-BLK-18-FINAL-01 — registry loader deadlocks the required maintenance
  fixture and all exact test-bearing gates.** Reproduce with:
  `go test -race -count=1 -v -timeout 5s ./internal/maintenance -run
  '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$'`. Expected: healthy /
  projected-corrupt maintenance evidence completes. Actual: timeout with
  `database/sql.(*DB).conn` waiting from `workflow.loadStates` while
  `workflow.Load` retains the outer `workflows` rows on a pool limited to one
  connection. Impact: Task 5's coverage, integration, and benchmark PASS
  records are contradicted; unit suite also cannot complete. VAL-010 behavior
  remains PASS through the real-entrypoint evidence in FR-E14.
  Fix direction: materialize workflow IDs and close the outer rows before
  nested queries, or load the registry with a non-nested query shape; preserve
  read-only behavior and ordering. Then rerun the focused maintenance healthy /
  corrupt / `--warn-as-error` evidence, all seven Task 5 commands, the exact
  unit suite, VAL-024 inventory/runtime parity, and final QA. QA Lead must not
  implement the fix.
- **QA-BLK-18-FINAL-02 — required Tech Lead spec sync is absent.** QA-30 and
  the QA-plan phase-end rule require the canonical-suite disposition and known
  TRD wording updates before final-review PASS. Evidence: FR-E13. This is a
  documentation/phase-exit block, not product-code authority for QA.

### Major

- None beyond the Blocking findings.

### Improvement

- Preserve the existing CF-023 dispositions: Task 2's report-finalization
  commit lacks its `Refs` footer, and the Phase 18 ledger has duplicate /
  out-of-time-order append entries. Do not rewrite published history or the
  append-only ledger in this QA review.

## Escalations

- Open concrete `ESC-NNN` entries: **none**.
- `project-docs/docs/current/tl_escalations_current.md` says `None` under Open
  Escalations.
- TRG-001 through TRG-005 were resolved by DEC-067 through DEC-071 without a
  minted `ESC-NNN`; placeholder strings `ESC-NNN` / `ESC-NEW` in the run ledger
  are not concrete open escalations.
- QA-BLK-18-FINAL-01 now requires Tech Lead/Phase Runner triage before phase
  closure; this review does not mint or edit escalation artifacts.

## Unresolved Risks And Carry-Forward Candidates

- Blocking availability/gate risk: QA-BLK-18-FINAL-01.
- Blocking documentation exit prerequisite: QA-BLK-18-FINAL-02.
- Existing improvement carry-forward: CF-023 workflow/report/ledger hygiene.
- No VAL is accepted as carry-forward; all 24 assertions have evidenced PASS
  outcomes, but phase closure remains blocked by the command gate and spec sync.

## Final Decision

**VERDICT: BLOCK**

Validation rollup: **24 PASS, 0 FAIL, 0 BLOCKED, 0 CARRY-FORWARD**.
Blocking findings: **2**. Open concrete `ESC-NNN`: **none**. Human checklist:
**4 items, all pending**.
