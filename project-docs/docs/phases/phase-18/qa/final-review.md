# Phase 18 Final Review - Workflow Tables Migration + Engine Data-Source Swap + Gate Seams

Reviewer: QA Lead
Date: 2026-07-23

## Verdict

**PASS**

This fresh rerun closes both findings from the 2026-07-20 blocked review.
The DEC-072 loader remediation terminates under the unchanged one-connection
fixture, every exact Task 5 gate command now terminates and passes in order,
the separate Quality Plan unit suite passes, coverage thresholds remain met,
and all performance controls remain below target. The Tech Lead-owned QA-30
specification synchronization is committed and verified.

Validation rollup: **24 PASS, 0 FAIL, 0 BLOCKED, 0 CARRY-FORWARD**. There is no
unresolved Blocking or Major finding. The four-item Human checklist is ready
for phase sign-off.

## Scope And Evidence

- Reviewed task/report HEAD:
  `37fb202a11f7adab28592f8a6002e001681be087`
  (`docs(task): finalize completion report for phase-18 task-5`).
- DEC-072 remediation implementation:
  `795a7df673708d21e558448e20edeb367f864078`
  (`fix(workflow): release registry rows before nested loads`).
- Tech Lead DEC-072/spec-sync commit:
  `4c029fbd129c277c2e2e5306ee5ff08dfd02649a`.
- Phase brief: `project-docs/docs/current/phase_brief.md`
- Progress and escalation state:
  `project-docs/docs/current/progress_current.md` and
  `project-docs/docs/current/tl_escalations_current.md`
- Decisions: `project-docs/docs/phases/phase-18/decisions.md`, including
  DEC-067 through DEC-072 and the canonical-suite verdict table.
- Validation contract:
  `project-docs/docs/phases/phase-18/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-18/qa/qa-plan.md`
- Task specs and committed completion reports:
  `project-docs/docs/phases/phase-18/tasks/`
- Mid-phase reviews:
  `review-1.md`, `review-1-remediation.md`, `review-2.md`, and
  `review-2-remediation.md`.
- Prior blocked final review: replaced by this fresh rerun; stale failure
  results were not reused as passing evidence.
- Pre-existing worktree change:
  `project-docs/docs/phases/phase-18/run-ledger.md` was not authored,
  modified, staged, or treated as QA evidence by this review.

### Reproducible Evidence Index

| Evidence ID | Result | Fresh reproducible evidence |
|---|---|---|
| FRR-E01 | PASS | The seven Task 5 commands ran in their specified order. Format exited 0 with no output; lint printed `0 issues.`; race/coverage exited 0 and produced `coverage.out`; the coverage report printed `total: (statements) 76.8%`; integration exited 0 for every package; benchmarks completed for every package; build exited 0 and produced `/tmp/fid`. |
| FRR-E02 | PASS | Coverage package results: `cmd/fid 77.7%`, `internal/issues 81.5%`, `internal/claims 81.4%`, `internal/workflow 83.2%`, `internal/policy 86.3%`, and `internal/journal 81.0%`. These satisfy COV-001/COV-002. |
| FRR-E03 | PASS | Separate exact unit suite `go test -race ./...` terminated successfully. Fresh package results included `cmd/fid 81.700s`, `internal/db 9.801s`, `internal/issues 32.963s`, `internal/journal 25.208s`, `internal/maintenance 5.033s`, and `internal/workflow 2.440s`. |
| FRR-E04 | PASS | Former deadlock reproduction `go test -race -count=1 -v -timeout 5s ./internal/maintenance -run '^TestCheckWorkflowStatesReadsProjectedWorkflowTables$'` printed `--- PASS` and completed in `1.235s`. The bounded registry regression `TestLoadReleasesWorkflowRowsBeforeNestedQueries` also passed under `SetMaxOpenConns(1)` in `0.38s`. |
| FRR-E05 | PASS | Fresh non-cached race run over the Phase 18 DB/workflow/issues/journal scopes passed. It covered migration 0007 schema/seeds/seam/down-up; migration 0008 FK/data/object/rendered-byte parity; prior/current projection versions; registry/frozen-definition equality; all reserved-rule cases; the 122-pair transition matrix and rule branches; all claim states; golden replay; and rebuild parity. Package results: DB `2.714s`, workflow `1.739s`, issues `6.701s`, journal `1.977s`. |
| FRR-E06 | PASS | Fresh closure integration `go test -tags integration -race -count=1 -v ./cmd/fid -run '^TestV01ClosureJourneyAndRebuildParity$'` passed in `6.15s` (`ok .../cmd/fid 7.168s`). It checks the complete projected-table set, schema/FKs, first and second rebuilds, journal bytes, and closure CLI behavior. |
| FRR-E07 | PASS | Fresh maintenance CLI tests passed. A built real-entrypoint healthy workspace returned `{"status":"pass","errors":[],"warnings":[],"info":{"issue_count":0,"claim_count_live":0,"recent_activity_7d":0}}`. After creating a task and deleting only its projected `task-default/todo` seed row with FKs disabled in the disposable fixture, `maintenance check --warn-as-error --format json` exited 1 with ordered `foreign_key_integrity` and `workflow_state` errors naming the invalid pair. Journal SHA-256 stayed `84d1abcccd258b176e39fbe8feb7af51d1b1322fec351ccf1ecb6cb9ad6fe231`. |
| FRR-E08 | PASS | Fresh projected-registry consumer tests passed: get/list/ready parity; reserved-rule behavior absence; and dev start/finish/block/gate source controls. The dev package command printed all four named PASS rows and `ok .../internal/dev 2.672s`. |
| FRR-E09 | PASS | Exact VAL-024 inventories reran. `rg -n 'definitions' internal cmd --glob '*.go'` found only the frozen fixture and test/parity references. `rg -n 'workflow\\.' internal cmd --glob '*.go' --glob '!**/*_test.go'` found DB-backed `workflow.Load` consumers and pure rule types; no production caller invokes the frozen fixture accessors. A separate reserved-name inventory found only the constants, with no receipt/gate query or evaluator. |
| FRR-E10 | PASS | Five fresh `project-docs/tools/report-lint.sh` runs passed. They verified every required command record, heading, implementation hash/subject, and committed/clean report. Task 5 resolves to `795a7df` and its finalized report is committed at `37fb202`. |
| FRR-E11 | PASS | QA-30 spec sync is committed at `4c029fb`: `02_TRD.md` §§4.4/4.6 record rebuild-only relaxed durability and the DB-backed runtime source; `03_TRD_SCHEMA.md` §§3/14.11/15 record shipped workflow tables, projected-table maintenance auditing, frozen-map test-only status, and the rebuild pragma exception. DEC-072's verdict table records `touched` or `verified-no-change` for all `00`–`10` specs plus `GLOSSARY.md`; the phase brief/progress record the applied sync and VAL-001..024. |
| FRR-E12 | PASS | Security/data-loss and query-plan checks remain green: fixed SQL or bound `?` parameters are used in registry/maintenance reads; `TestIssuesEnrichedQueryPlanUsesConcernAndTriageIndexes` passed; rebuild/closure and maintenance evidence showed unchanged journal bytes; `go.mod` and `go.sum` are unchanged by the remediation. |

## Task 5 Seven-Command Gate Audit

This table supersedes the failed 2026-07-20 reproduction. Each command below
was executed afresh against report HEAD, in the exact Task 5 order.

| # | Exact command | Fresh rerun | Evidence |
|---:|---|---|---|
| 1 | `gofmt -l cmd internal` | PASS | Exit 0; no output. |
| 2 | `golangci-lint run ./...` | PASS | Exit 0; `0 issues.` |
| 3 | `go test -race -coverprofile=coverage.out ./...` | PASS | Exit 0; all packages passed and a fresh profile was produced (FRR-E01/02). |
| 4 | `go tool cover -func=coverage.out` | PASS | Exit 0; `total ... 76.8%` from the profile created by command 3. |
| 5 | `go test -tags integration -race ./...` | PASS | Exit 0; all packages terminated successfully. FRR-E06 separately reran the closure integration non-cached. |
| 6 | `go test -bench . -benchmem ./...` | PASS | Exit 0; all packages terminated and every required performance row passed. |
| 7 | `go build -o /tmp/fid ./cmd/fid` | PASS | Exit 0; built binary was used for FRR-E07 real-entrypoint evidence. |

Fresh benchmark comparison:

| Control | Fresh result | Target | Outcome |
|---|---:|---:|---|
| 10k-issue registry read | `49671627 ns/op` (`49.67 ms/registry-read`) | `<500 ms` | PASS |
| 1k-event rebuild diagnostic | `306815531 ns/op` (`306.8 ms/rebuild`) | diagnostic | PASS |
| 10k-event rebuild | `2918252293 ns/op` (`2918 ms/rebuild`) | `<5 s` | PASS |
| 100k-event rebuild | `27660882944 ns/op` (`27661 ms/rebuild`) | `<30 s` | PASS |

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | FRR-E01; no files listed. |
| Lint/static analysis | `golangci-lint run ./...` | PASS | FRR-E01; `0 issues.` |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | FRR-E01/02. |
| Coverage report | `go tool cover -func=coverage.out` | PASS | Fresh total `76.8%`; enforced package thresholds met. |
| Unit | `go test -race ./...` | PASS | FRR-E03; exact command terminated successfully. |
| Integration | `go test -tags integration -race ./...` | PASS | FRR-E01; FRR-E06 supplies non-cached closure evidence. |
| End-to-end | N/A | N/A | No accepted separate E2E/held-out command exists; integration-tagged CLI journeys are the documented substitute for this phase. |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | FRR-E01 and benchmark table; exact command terminated. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | FRR-E01/07. |
| Full phase check | Seven Task 5 commands plus exact unit suite | PASS | All eight commands terminated and passed; QA-25 is green. |

## Validation Contract Outcomes

Every active VAL-001 through VAL-024 was re-confirmed. Each PASS cites fresh
command evidence from this rerun.

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | FRR-E05: `TestMigration0007CreatesWorkflowSchema` and all four table subtests passed, including PK/UNIQUE/FK introspection. |
| VAL-002 | PASS | QA Lead | FRR-E05: seeded state/category parity and `TestRegistryEqualsV01Definitions` passed for all three workflows. |
| VAL-003 | PASS | QA Lead | FRR-E05: complete seeded edge parity plus the 122-pair transition matrix passed, including explicit cancel and absent-edge cases. |
| VAL-004 | PASS | QA Lead | FRR-E05: seed/registry rule and parameter parity plus every satisfied/blocked ruled-edge test passed. |
| VAL-005 | PASS | QA Lead | FRR-E05: canonical gate-obligation JSON round-trip test passed without loss. |
| VAL-006 | PASS | QA Lead | FRR-E05: migration 0008 accepted valid INSERT/UPDATE pairs, rejected invalid INSERT/UPDATE pairs, and returned a clean `foreign_key_check`. |
| VAL-007 | PASS | QA Lead | FRR-E05: canonical pre/post and down/up `issues` row snapshots were identical. |
| VAL-008 | PASS | QA Lead | FRR-E05: exhaustive ordered-pair matrix, rule branches, DB-source controls, and exact-error characterization passed. |
| VAL-009 | PASS | QA Lead | FRR-E05: missing, wrong-holder, stale, and live caller claim states all passed using DB-resolved categories. |
| VAL-010 | PASS | QA Lead | FRR-E04 proves the formerly hung one-connection path completes; FRR-E07 proves healthy, invalid-pair, non-zero `--warn-as-error`, ordered finding, and unchanged-journal states. |
| VAL-011 | PASS | QA Lead | FRR-E06: every projected table, schema, FK state, first rebuild, and second rebuild matched. |
| VAL-012 | PASS | QA Lead | FRR-E05: prior-version rebuild/catch-up and current-version no-rebuild branches both passed. |
| VAL-013 | PASS | QA Lead | FRR-E05: `TestReplayJournalAppliesGoldenV1JournalInFileOrder` passed with the expected migrated projection. |
| VAL-014 | PASS | QA Lead | FRR-E04/05: registry equality passed after the loader lifetime fix, proving unchanged contents and order. |
| VAL-015 | PASS | QA Lead | FRR-E05: both reserved names loaded as known-but-unimplemented under `query_only`, with no write/evaluation. |
| VAL-016 | PASS | QA Lead | FRR-E05: migration seam test found zero seeded reserved-name rows. |
| VAL-017 | PASS | QA Lead | FRR-E05/09: injected reserved rule stayed unevaluated and static production inventory found no receipt/gate evaluator or query. |
| VAL-018 | PASS | QA Lead | FRR-E05: complete dependent-object/down-up fingerprint, retained trigger probes, autoindexes, and query-plan control passed. |
| VAL-019 | PASS | QA Lead | FRR-E05: pre/post migration issue get/list rendered bytes were identical. |
| VAL-020 | PASS | QA Lead | FRR-E05: same-event release and exact released-claim rebuild-absence tests both passed. |
| VAL-021 | PASS | QA Lead | FRR-E06: every journal file remained byte-identical across migration, stale refresh, rebuild, and second rebuild. |
| VAL-022 | PASS | QA Lead | FRR-E05/06/08: projected get/list/ready category/result/order parity and closure journey passed. |
| VAL-023 | PASS | QA Lead | FRR-E06/08: projected-registry dev start/finish/block/gate outcomes and effects passed. |
| VAL-024 | PASS | QA Lead | FRR-E09 exact inventories show fixture-only frozen-map access; FRR-E05/06/08 supply paired runtime parity. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| TRD-ARCH-009; TRD_SCHEMA §§2–3 | MET | VAL-001..007, VAL-014..016, VAL-018..019 PASS. |
| PRD-WF-001–010; TRD §4.6 behavior parity | MET | VAL-002..004, VAL-008, VAL-014, VAL-022..024 PASS. |
| PRD-WF-014/015 seam; DEC-066 | MET | VAL-005, VAL-015..017 PASS; no receipt/gate behavior found. |
| PRD-CLAIM-005/006 | MET | VAL-009/020 PASS. |
| PRD-OPS-003/004; TRD-ARCH-001..003 | MET | VAL-011..013/021 PASS; FRR-E05/06. |
| PRD-OPS-005–008; TRD_SCHEMA §14.11 | MET | VAL-010 PASS with package, real-entrypoint, corrupt-pair, exit-code, and journal evidence. |
| TRD-PERF-001/002/003 | MET | Fresh read/rebuild rows all remain below target. |
| TRD-PERF-004 / TRD-ARCH-008 | MET | Query-plan and issues-rebuild controls passed in FRR-E05/12. |
| COV-001/COV-002 | MET | FRR-E02; every named domain package is at least 80%, CLI is at least 70%. |
| Heavy Mode / Quality Plan §5.2 | MET | FRR-E01–12 and the completed deterministic checklist below. |

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS | FRR-E05 migration 0007 schema introspection. |
| QA-02 | PASS | FRR-E05 complete state/category seed and registry equality. |
| QA-03 | PASS | FRR-E05 complete edge/cancel/absent-edge parity and 122-pair matrix. |
| QA-04 | PASS | FRR-E05 complete rule/roles/params equality and branch tests. |
| QA-05 | PASS | FRR-E05 canonical gate JSON round-trip. |
| QA-06 | PASS | FRR-E05 both reserved names under read-only loading. |
| QA-07 | PASS | FRR-E05 zero seeded reserved-name rows. |
| QA-08 | PASS | FRR-E05 injected-rule result plus FRR-E09 static inventory. |
| QA-09 | PASS | FRR-E05 valid/invalid INSERT and UPDATE matrix plus clean FK check. |
| QA-10 | PASS | FRR-E05 canonical pre/post issue rows and empty diff. |
| QA-11 | PASS | FRR-E05 complete dependent-object/down-up fingerprint and probes. |
| QA-12 | PASS | FRR-E05 rendered get/list byte parity. |
| QA-13 | PASS | FRR-E05 exhaustive transition/rule matrix and exact-error characterization. |
| QA-14 | PASS | FRR-E05 missing/wrong/stale/live claim states. |
| QA-15 | PASS | FRR-E05 same-event release and rebuilt claim absence. |
| QA-16 | PASS | FRR-E04/07 literal healthy/corrupt maintenance states, exit 1 under `--warn-as-error`, ordered errors, and identical journal SHA-256. |
| QA-17 | PASS | FRR-E06 full projection/schema/FK parity across two rebuilds. |
| QA-18 | PASS | FRR-E05 prior/current projection-schema branches. |
| QA-19 | PASS | FRR-E05 golden v0.1 replay. |
| QA-20 | PASS | FRR-E06 four-boundary journal-byte parity. |
| QA-21 | PASS | FRR-E04/05 registry equality after DEC-072. |
| QA-22 | PASS | FRR-E05/06/08 get/list/ready parity. |
| QA-23 | PASS | FRR-E06/08 dev composite parity. |
| QA-24 | PASS | FRR-E09 exact inventories paired with runtime parity. |
| QA-25 | PASS | FRR-E01–04: seven-command gate and exact unit suite all terminate; coverage thresholds pass. |
| QA-26 | PASS | FRR-E05 migration 0007/0008 down/up tests, row/schema checks, and FK state. |
| QA-27 | PASS | FRR-E06/07/12: fixed/bound SQL, no dependency change, and no journal/data-loss finding. |
| QA-28 | PASS | Fresh benchmark table plus FRR-E05 query-plan control. |
| QA-29 | PASS | FRR-E10: five report-lint passes and Git-resolved implementation/report metadata. |
| QA-30 | PASS | FRR-E11: committed DEC-072 TRD wording, current-doc sync note, and full canonical-suite disposition table. |

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: **YES**,
  as recorded in `qa-plan.md` §6.
- Required command placeholders: **none**.
- Invalid command, test scope, fixture, artifact assumption, or Human step:
  **none**.
- DEC-072 feasibility recheck: **PASS**. The original one-connection command,
  bounded regression, seven-command gate, unit suite, focused maintenance
  fixture, VAL-024 inventories, and performance controls all executed as
  specified.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC References |
|---|---|---|---|---|
| `review-1.md` | BLOCK | YES | VAL-020 rebuilt-claim gap remediated by DEC-068; Task 2 report-footer hygiene remains in CF-023. | none |
| `review-1-remediation.md` | CONTINUE | YES | Exact released-claim rebuild evidence passed again in FRR-E05. | none |
| `review-2.md` | BLOCK | YES | VAL-006 valid UPDATE gap remediated by DEC-069. | none |
| `review-2-remediation.md` | CONTINUE | YES | Full valid/invalid INSERT/UPDATE matrix passed again in FRR-E05. | none |

Task 5 is the final task. Per the QA plan and DEC-072, this comprehensive
rerun replaces a third mid-phase review; no further mid-phase review is
scheduled.

## Documentation Review

- Completion reports present, committed, clean, and report-lint valid for all
  five tasks: **YES** (FRR-E10).
- Impacted canonical specifications synchronized: **YES** (FRR-E11).
- QA-30 full canonical-suite disposition recorded: **YES**.
- Task 5 report corrected with fresh terminating gates, `76.8%` coverage, and
  the DEC-072 implementation metadata: **YES**.
- Current phase brief/progress record DEC-072 and this rerun as the next action.
  The Phase Runner/Tech Lead should refresh their execution-status lines after
  this committed QA outcome; that post-review handoff is not a missing
  prerequisite for QA-30.

## Human Checklist Readiness

Generated at
`project-docs/docs/phases/phase-18/qa/human-checklist.md`.

- Contains only `(human)` items: **YES**.
- Uses one observable per row: **YES**.
- Human-only `VAL-NNN` assertions mapped: **N/A**; all Phase 18 assertions are
  deterministic and QA Lead-owned.
- Items: transition feel, get/list/ready presentation, maintenance clarity,
  and Phase 18 dogfood acceptance.
- Status: **ready for Human execution after this PASS**.

## Coverage Caveat

This deterministic PASS proves the asserted mechanics, data states, command
outputs, structural conformance, parity matrices, rebuild behavior, journal
immutability, and performance controls. It does not prove unasserted workflow
feel, product-owner judgment, or dogfooding quality; those remain pending in
the Human checklist.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | READY | Five committed reports, Ready contract, QA plan, reviews, 24 evidenced outcomes, PASS final review, and Human checklist ready. |
| Automated gates | READY | Seven-command gate, unit suite, coverage thresholds, integration, benchmarks, and build all pass. |
| Held-out scenarios | DEFERRED | No accepted harness/location; Phase 18 is not release readiness. |
| Release readiness | DEFERRED | Non-release phase; threat model/scanners/Mode B scenario are not silently substituted here. |

## Findings

### Blocking

- None.

### Major

- None.

### Improvement

- Preserve CF-023 for the already-recorded workflow/report/ledger hygiene
  issues. Do not rewrite published history or the append-only run ledger in
  this QA review.

## Escalations

- Phase 18 ESC-001 is resolved by DEC-072 and its remediation is verified by
  FRR-E01–04 and FRR-E11.
- It remains listed in `tl_escalations_current.md` only until this PASS is
  committed and the Tech Lead archives it at wrap.
- No unresolved escalation is raised by this rerun.

## Unresolved Risks And Carry-Forward Candidates

- Human product/dogfood judgment remains pending in the four-item checklist.
- Existing CF-023 remains the only workflow-hygiene carry-forward relevant to
  this phase.
- No validation assertion, gate failure, parity risk, data-loss risk, or
  performance breach is carried forward.

## Final Decision

**VERDICT: PASS**

Validation rollup: **24 PASS, 0 FAIL, 0 BLOCKED, 0 CARRY-FORWARD**.
QA-25: **PASS**. QA-30: **PASS**. Blocking findings: **0**. Major findings:
**0**. Human checklist: **4 items, ready and pending**.
