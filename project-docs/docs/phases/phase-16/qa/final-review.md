# Phase 16 Final Review - Activity Timeline + Read-Only Escape Hatches

Reviewer: QA Lead
Date: 2026-07-05

## Verdict

PASS WITH IMPROVEMENTS

Deterministic phase-end QA passes. VAL-001 through VAL-012 are PASS with
reproducible evidence. VAL-013 is CARRY-FORWARD to the human checklist because
the required dogfood evidence is human-only and has not yet been recorded.

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Run ledger: `project-docs/docs/phases/phase-16/run-ledger.md`
- Decisions: `project-docs/docs/phases/phase-16/decisions.md`
- Validation contract: `project-docs/docs/phases/phase-16/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-16/qa/qa-plan.md`
- Mid-phase reviews: `project-docs/docs/phases/phase-16/qa/review-1.md`, `review-2.md`
- Completion reports: `project-docs/docs/phases/phase-16/tasks/task-1-done.md` through `task-6-done.md`
- Triage log: `project-docs/docs/phases/phase-16/triage-log.md`

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | No files listed. |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Output: `0 issues.` |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | All packages passed; examples: `cmd/fid` 77.5%, `internal/activity` 81.1%, `internal/issues` 81.5%, `internal/journal` 78.9%, `internal/sqlquery` 65.3%. |
| Coverage report | `go tool cover -func=coverage.out` | PASS | Total statement coverage: 75.6%. |
| Integration | `go test -tags integration -race ./...` | PASS | All packages passed; `cmd/fid` completed in 64.951s and full package set completed without failure. |
| End-to-end | N/A | N/A | No accepted standalone E2E harness yet; Phase 16 lifecycle E2E is covered by integration test `TestV01EActivityLifecycleCoverageAndRebuildParity`. |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | All packages passed; `BenchmarkRebuildWorkspace1kEvents` 372.2 ms/rebuild and `BenchmarkRebuildWorkspace10kEvents` 3605 ms/rebuild. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Command exited 0 and produced `/tmp/fid`. |
| Full phase check | QA plan section 5 full phase gate | PASS | Every command above passed in this review window. |

## Validation Contract Outcomes

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead / QA Reviewer | Task 1 report records activity shape/order/all-format coverage; Review 1 marks ON TRACK; phase gate passed. |
| VAL-002 | PASS | QA Lead / QA Reviewer | Task 1 report records `--on` include/exclude tests; Review 1 verifies parameterized `ref_type`/`ref_id` filtering. |
| VAL-003 | PASS | QA Lead / QA Reviewer | Task 1 report records exact-boundary and date-only `--since` tests; Review 1 verifies `occurred_at >= ?` semantics. |
| VAL-004 | PASS | QA Lead / QA Reviewer | Task 1 report records `--event` matching/non-matching tests; Review 1 verifies event-name filtering. |
| VAL-005 | PASS | QA Lead / QA Reviewer | Task 2 report records one-event/one-row behavior for create, AC add, AC verify, and complete plus rejection preservation; Review 1 reran focused `internal/issues` and `internal/dev` tests. |
| VAL-006 | PASS | QA Lead / QA Reviewer | Task 3 report records concern raise/address, unclaim, link add, and `dev gate` `gate_recorded` rows plus zero-row rejection paths; Review 2 reran `TestV01EActivityLifecycleCoverageAndRebuildParity`. |
| VAL-007 | PASS | QA Lead / QA Reviewer | Task 3 report and Review 2 record lifecycle coverage including `gate_recorded`, delete `fid.sqlite`, `fid rebuild`, and ordered `activity_log` equality; final review reran `go test -tags integration ./cmd/fid -run TestV01EActivityLifecycleCoverageAndRebuildParity` PASS. |
| VAL-008 | PASS | QA Lead / QA Reviewer | Task 4 report records SELECT, EXPLAIN, read PRAGMA, `WITH ... SELECT`, and all formats; Review 2 and final focused SQL rerun passed. |
| VAL-009 | PASS | QA Lead / QA Reviewer | Task 4 report records classifier rejection plus read-only SQLite backstop; Review 2 and final reruns passed `TestSQLCommandRejectsMutationAttemptsWithoutChangingProjectedCounts` and classifier/backstop tests. |
| VAL-010 | PASS | QA Lead | Task 5 report records known/missing event coverage, all formats, and `HashEvent` parity; final review reran `go test ./internal/journal ./cmd/fid -run 'TestFindEventReturnsStoredEventLocationAndHash|TestFindEventReturnsTypedNotFoundError|TestRevertVerdictClassifiesKnownEventTypes|TestRevertVerdictDefaultsUnknownTypesToSafeNotEligible|TestJournalInspect'` PASS. |
| VAL-011 | PASS | QA Lead | Task 5 report records inverse, retraction, and `eligible:false` correction verdict table; same final journal-focused rerun passed. |
| VAL-012 | PASS | QA Lead | Task 6 report records raw `gate_runs` history, finished-task history, empty-list behavior, all formats, and completion-snapshot separation; final review reran `go test ./internal/issues ./cmd/fid ./internal/format -run 'Test(GetReturnsCentralizedDetailShape|SectionValue|IssueGetSectionsAndDueDateReadBack|WriteMarkdownRendersIssueGateLists)'` PASS. |
| VAL-013 | CARRY-FORWARD | Human | Human-only dogfood has not been run. Generated `project-docs/docs/phases/phase-16/qa/human-checklist.md` with HUMAN-16a..d rows covering activity timeline, SQL SELECT, SQL mutation rejection, and journal inspect. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| PRD-OUT-017 activity timeline | MET | VAL-001..007 PASS; Review 1 and Review 2 evidence; final lifecycle rerun PASS. |
| CF-011 broad `activity_log` reconciliation | MET | VAL-005/006/007 PASS; DEC-043 included `dev gate` and `gate_recorded`; rebuild parity passed. |
| PRD-OUT-008 SQL escape hatch | MET | VAL-008/009 PASS; accepted read statements and all-format output tested. |
| TRD-SEC-002 SQL read-only security | MET | VAL-009 PASS; classifier rejection and read-only connection backstop tested. |
| PRD-OPS-013 journal inspect | MET | VAL-010/011 PASS; Task 5 evidence includes event location, hash, effects, and revert verdicts. |
| CF-016 gate history read surface | MET | VAL-012 PASS; Task 6 evidence includes in-progress/finished/empty states and all formats. |
| Heavy-mode QA cadence | MET | DEC-045 repaired missed cadence; Review 1 and Review 2 are recorded in the run ledger with 0 blocking findings. |
| Human dogfood | CARRY-FORWARD | VAL-013 is human-only; checklist generated and awaiting human execution. |

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS | `project-docs/tools/report-lint.sh` passed for all six completion reports. `rg` confirmed task specs and reports record required command coverage; Tasks 2-4 record required integration PASS; Tasks 1/5/6 record integration not required with reasons; all record build PASS. |
| QA-02 | PASS | Task specs and completion reports cite covered VAL IDs; completion reports list documentation updates or explicit unchanged/deviation notes. |
| QA-03 | PASS | Review 1 validates VAL-005 and PI-6; Review 2 validates VAL-006. Task 3 final report and final lifecycle rerun include `gate_recorded` and rebuild parity. |
| QA-04 | PASS | `TestV01EActivityLifecycleCoverageAndRebuildParity` passed at final review, proving ordered `activity_log` equality after delete-and-rebuild. |
| QA-05 | PASS | Final review reran SQL classifier/backstop tests and `TestSQLCommandRejectsMutationAttemptsWithoutChangingProjectedCounts`; both passed. |
| QA-06 | PASS | Final review reran journal inspect/hash/verdict tests; Task 5 report records the PI-5 event-type verdict table and built CLI smoke using event `EVT-20260705-000004`. |
| QA-07 | PASS | Final review reran gates section tests; Task 6 report records built CLI smoke after `dev start` -> `dev gate` -> `dev finish`, returning raw gate history from `gate_runs`. |
| QA-08 | PASS | Task reports name updates to `02_TRD.md`, `03_TRD_SCHEMA.md`, `06_DESIGN_PLAN_CLI.md`, `AGENTS.md`, and task reports. Phase wrap still needs to replace pending spec-sync verdicts in `decisions.md`. |

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: YES, QA plan section 6.
- Invalid or unvalidated command, test scope, fixture, artifact, or human checklist step found: none.
- Placeholder audit: `rg -n "<TBD>|<E2E_COMMAND>|PLACEHOLDER|TODO" AGENTS.md project-docs/docs/phases/phase-16/qa project-docs/docs/phases/phase-16/tasks` found only the standing AGENTS.md warning text and the QA plan's already-recorded placeholder-audit note.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC References |
|---|---|---|---|---|
| `review-1.md` | CONTINUE | N/A | Activity CLI fixture issue resolved by DEC-044 / Task 4; Task 2 report self-contained rejection-evidence note remains a process improvement only. | none |
| `review-2.md` | CONTINUE | N/A | `internal/sqlquery/sqlquery.go` 595-line file-size finding remains a carry-forward candidate for the next SQL package reopen. | none |

## Documentation Review

- Impacted docs updated: YES, per task reports.
- Completion reports present for all tasks: YES, Tasks 1-6.
- Run ledger consistency: YES, Tasks 5 and 6 are complete with commits `12c5107` / `91ee35b` and `31fb3dd` / `56ab49f`.
- DEC-046 / TRG-046 disposition: PASS, triage log records RESUME and Task 5 report records the constrained `internal/activity/activity_test.go` remediation preserving VAL-001..004 and PRD-ERR-014.
- Specs needing Tech Lead wrap-up sync: phase wrap must replace pending spec-sync verdicts in `project-docs/docs/phases/phase-16/decisions.md`; current docs may also need refresh because the active phase brief/progress still describe the pre-Task-5 blocker state.

## Human Checklist Readiness

Generated at `project-docs/docs/phases/phase-16/qa/human-checklist.md`.

- Contains only `(human)` items: YES.
- Uses one observable per row: YES.
- Human-only `VAL-NNN` assertions mapped: YES, VAL-013 to HUMAN-16a..d.
- Product/perception/dogfooding risks to call out: VAL-013 remains unproven until a human records outcomes.

## Coverage Caveat

A deterministic PASS proves only the asserted mechanics, data states, command
outputs, and structural conformance checks. It does not prove dogfooding
quality or product-owner judgment; those remain unconfirmed until the human
checklist records outcomes.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | READY | Validation contract, QA plan, Reviews 1-2, final review, and human checklist exist. |
| Automated gates | READY | Full phase-level command gate passed in this review. |
| Held-out scenarios | DEFERRED | AGENTS.md records E2E/scenario validation as N/A until harness location and command are accepted. |
| Release readiness | DEFERRED | Security scan, vulnerability audit, license check, and secret scan are release-readiness gates, not Phase 16 gates. |

## Findings

### Blocking

- None.

### Major

- None.

### Improvement

- `internal/sqlquery/sqlquery.go` is 595 lines, above the 400-line non-test file target. This is not blocking because VAL-008/009 and the phase gate pass, but the scanner/tokenizer/read execution code should be split when `internal/sqlquery` is next reopened.
- `internal/format/markdown.go` is 661 lines after Task 6 markdown support. This is not Phase 16 blocking because behavior is covered and the file was already a broad renderer, but future formatter work should consider a behavior-preserving split.
- Task 2's report could have enumerated each VAL-005 rejection fixture more explicitly. Review 1 treated this as a process improvement only; report-lint and focused tests passed.

## Unresolved Risks And Carry-Forward Candidates

- VAL-013 is CARRY-FORWARD to the human checklist and must not be marked PASS until human dogfood evidence is recorded.
- Carry-forward candidate: split `internal/sqlquery/sqlquery.go` on the next SQL package reopen.
- Carry-forward candidate: split `internal/format/markdown.go` on the next formatter-heavy task.
- Existing CF-011 and CF-016 are ready for closure at phase wrap because VAL-005/006/007 and VAL-012 passed.
