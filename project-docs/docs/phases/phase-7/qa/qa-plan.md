# Phase 7 QA Plan - Acceptance Criteria Surface

> Phase: 7 - Acceptance criteria surface (v0.1c)
> Refreshed: 2026-06-14
> QA Lead: Codex
>
> Location: `project-docs/docs/phases/phase-7/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 7 ships the acceptance-criteria write/read surface:
`fid issue ac add`, `fid issue ac verify`, `fid issue ac list`, and
real-command evidence for the existing `all_acs_verified` transition rule plus
delete-DB `fid rebuild` parity.

Source artifacts used:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/01_PRD.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/docs/phases/phase-7/decisions.md`
- `project-docs/docs/phases/phase-7/qa/validation-contract.md`
- `project-docs/docs/phases/phase-7/qa/review-1.md`
- `project-docs/docs/phases/phase-7/tasks/task-1.md`
- `project-docs/docs/phases/phase-7/tasks/task-1-done.md`
- `project-docs/docs/phases/phase-7/tasks/task-2.md`
- `project-docs/docs/phases/phase-7/tasks/task-2b.md`
- `project-docs/docs/phases/phase-7/tasks/task-3.md`
- `project-docs/docs/phases/phase-7/tasks/task-4.md`

DEC-020 note: this plan follows the phase brief, decisions active file, and
validation contract: AC writes do not emit `activity_log` in Phase 7. The
current progress file still says DEC-020 awaited confirmation; that is a
current-state doc sync note, not a validation-contract discrepancy.

### Resumed Fixup-Only Run

This refresh preserves the original Phase 7 plan for Tasks 1-4. The resumed run
is scoped only to the Review 1 fixups captured by DEC-021 / ESC-005 and Task
2b:

- VAL-004: strengthened deterministic append-order proof for same-role
  fail->pass, with no `time.Sleep` or timestamp-boundary dependence.
- VAL-009: targeted rebuild parity for latest-per-role after the DEC-021
  ordering change.
- No VAL: committed Task 1 completion-report finalization check.
- No VAL: Task 2b file-size tripwire check is scoped only to the AC
  write-path split files, `internal/issues/ac.go` and
  `internal/issues/ac_verify.go`, each `<= 400` lines.
- No VAL: `internal/issues/issue.go` remains explicit pre-existing debt,
  tracked as CF-012, and is out of scope for Task 2b.

Task 2b does not reopen or replace Task 1, Task 2, Task 3, or Task 4 coverage.
Task 2b may run broader regression tests, but the fixup-only validation scope
above is the only new Phase Runner obligation before final QA.

Discrepancy note: the validation contract still lists VAL-009's owning window as
Task 4. Task 2b adds targeted rebuild-stability evidence for DEC-021 latest
ordering without weakening or replacing the original Task 4 VAL-009 obligation.

### Phase Requirement Coverage

| Source | Requirement | Covered by |
|---|---|---|
| PRD-AC-001 | Ordered ACs on tasks and subtasks. | VAL-001 |
| PRD-AC-002 | Reject ACs on epics. | VAL-002 |
| PRD-AC-003 | Multi-role AC verification with independent outcomes. | VAL-003 |
| PRD-AC-004 | Verification stores role, status, evidence fields, and timestamp. | VAL-005 |
| PRD-AC-005 | Verification history is append-only. | VAL-004 |
| PRD-AC-006 | AC listing includes latest verification status per verifier role. | VAL-006 |
| PRD-OUT-001 | Read commands output JSON by default. | VAL-007 |
| PRD-OUT-002 | Read commands support TOON. | VAL-007 |
| PRD-OUT-003 | Read commands support markdown. | VAL-007 |
| PRD-OUT-003a | Read commands support YAML. | VAL-007 |
| PRD-OUT-004 | Read commands accept `--format` override. | VAL-007 |
| PRD-WF-007 | `all_acs_verified` transition rule with required roles. | VAL-008 |
| ADR-0001 / TRD persistence | Journal is canonical and SQLite projection is derived. | VAL-009 |
| TRD schema §2.3 | `acceptance_criteria` and `ac_verifications` state, IDs, append-only semantics, and latest-per-role query behavior. | VAL-001, VAL-003, VAL-004, VAL-005, VAL-006, VAL-009 |
| TRD workflow | Task workflow uses `all_acs_verified{"roles":["developer"]}` and QA workflow uses `{"roles":["qa_lead"]}`. | VAL-008 |
| CLI design §3 | `fid issue ac add`, `fid issue ac verify`, `fid issue ac list` syntax, flags, stderr, and JSON shape. | VAL-001..VAL-007 |
| Phase brief / roadmap v0.1c | Human AC lifecycle dogfood over add/list, verify/list, and epic rejection. | VAL-010 |
| DEC-020 / CF-011 | AC writes do not emit `activity_log` in Phase 7; timeline reconciliation is deferred. | No VAL by design; verified as scope discipline in QA-13. |
| Phase brief out-of-scope | `fid issue complete`, `task_completions`, concerns, and derived concern flags are Phases 8-9. | No Phase 7 VAL by design. |

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: YES.
- Every `VAL-NNN` has one validator role: YES.
- Every `VAL-NNN` has a concrete evidence method: YES.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: YES. VAL-010 is intentionally fan-out only.
- Conditional or binary `VAL-NNN` assertions require both qualifying and
  non-qualifying states: YES for VAL-002, VAL-003, VAL-004, VAL-005, VAL-008.
- Parity or "matches reference" `VAL-NNN` assertions have deterministic
  conformance proof when available: YES for VAL-007 and VAL-009.
- Revisions requested from Tech Lead: None.
- Validation contract file modified by QA Lead: NO.

### VAL Evidence Refinement

| VAL ID | Requirement coverage | Refined evidence method | Validator role | Task/window | Plan coverage |
|---|---|---|---|---|---|
| VAL-001 | PRD-AC-001; TRD schema `acceptance_criteria`; CLI design `ac add` | Task 1 tests and transcript show two AC adds on task/subtask, `AC-` ID shape, persisted text/ordinal, and ordinal-then-id read-back through `issue get --section acceptance` or `ac list` after Task 3 exists. | QA Lead | Task 1; Review 1; final QA | QA-01, QA-10; MPR-01 |
| VAL-002 | PRD-AC-002; CF-007 validate-before-append | Task 1 tests and transcript show epic rejection and task/subtask success; rejection evidence must include non-zero exit, kind-constraint error, and unchanged journal line count/projected AC rows. | QA Lead | Task 1; Review 1; final QA | QA-02, QA-10; MPR-01 |
| VAL-003 | PRD-AC-003 | Task 2 tests and transcript verify one AC as `developer=pass` and `qa_lead=fail`; latest-per-role evidence shows both roles retained independently. | QA Lead | Task 2; Review 1; final QA | QA-03, QA-10; MPR-01 |
| VAL-004 | PRD-AC-005 | Task 2b tests and transcript show same role fail-only latest is fail, then immediate fail->pass retains two rows and latest becomes pass by journal append order, without `time.Sleep` or timestamp-boundary dependence. | QA Lead | Task 2b; fixup re-review; final QA | QA-04, QA-10; MPR-01-FIX |
| VAL-005 | PRD-AC-004 | Task 2 tests and transcript show verifier_role, status, evidence_type, evidence_ref, notes, verified_at persisted; invalid status rejects before append with unchanged journal/rows. | QA Lead | Task 2; Review 1; final QA | QA-05, QA-10; MPR-01 |
| VAL-006 | PRD-AC-006; TRD latest-per-role semantics | Task 3 tests and transcript show `ac list --format json` shape with one latest entry per role, not raw history; superseded statuses are absent from list output. | QA Lead | Task 3; Review 2; final QA | QA-06, QA-10; MPR-02 |
| VAL-007 | PRD-OUT-001..004; CLI design format dispatch | Task 3 tests and transcripts run `json`, `toon`, `markdown`, `yaml`, and no `--format`; omitted flag must structurally match JSON output. | QA Lead | Task 3; Review 2; final QA | QA-07, QA-10; MPR-02 |
| VAL-008 | PRD-WF-007; Phase 5 fixture replacement | Task 4 integration/CLI evidence creates real ACs, rejects gated transition before required pass with unmet AC IDs/roles and no state/journal change, verifies via real `ac verify`, then transition succeeds. | QA Lead | Task 4; Review 2; final QA | QA-08, QA-10; MPR-02 |
| VAL-009 | ADR-0001; TRD journal canonical projection | Task 4 integration/CLI evidence snapshots AC and verification state, deletes SQLite DB, runs `fid rebuild`, and compares `acceptance_criteria`, `ac_verifications`, and latest-per-role read-back. Task 2b adds targeted rebuild parity for the DEC-021 append-order latest tie-break. | QA Lead | Task 4; Task 2b fixup evidence; Review 2; fixup re-review; final QA | QA-09, QA-10; MPR-02; MPR-01-FIX |
| VAL-010 | Phase brief; roadmap v0.1c dogfood | Human dogfood in a real `.fid/` workspace, split into ordered add+list, multi-role verify/list, and epic rejection rows. | Human | Phase end | HV-010a, HV-010b, HV-010c |

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | `fid issue ac add` | VAL-001, VAL-002 | Medium | New write path, journal event, projection, and validate-before-append epic rejection. |
| Task 2 | `fid issue ac verify` | VAL-003, VAL-004, VAL-005 | Medium | Append-only history and latest-per-role semantics affect later listing and workflow rules. |
| Task 2b | DEC-021 / ESC-005 fixup for deterministic latest-per-role ordering | VAL-004, targeted VAL-009 | High | Resumed fixup-only task. Also checks committed Task 1 report finalization and scoped AC split-file sizes for `internal/issues/ac.go` and `internal/issues/ac_verify.go` `<= 400`; those two checks carry no VAL. `internal/issues/issue.go` is CF-012 pre-existing debt and out of scope. |
| Task 3 | `fid issue ac list` | VAL-006, VAL-007 | Medium | Public output shape and format dispatch; must agree with existing issue acceptance read path. |
| Task 4 | Real `all_acs_verified` gating and rebuild parity | VAL-008, VAL-009 | High | End-to-end closure; proves Phase 5 fixture evidence is replaced by real writes and replay. |
| Phase end | Human dogfood | VAL-010 | Medium | Human-only judgment over real AC lifecycle clarity and behavior. |

## 4. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Task 1 completion report records exact required gate results, coverage output, and AC add/epic-reject transcripts. | VAL-001, VAL-002 | Developer | Task 1 completion report |
| EVID-02 | Task 2 completion report records exact required gate results, coverage output, pass verify, multi-role, fail->pass, and invalid-status transcripts. | VAL-003, VAL-004, VAL-005 | Developer | Task 2 completion report |
| EVID-03 | Task 2b completion report records exact required gate results, coverage output, deterministic no-sleep fail->pass transcript, rebuild latest-per-role parity transcript, Task 1 committed-report proof, and scoped file-size proof for `internal/issues/ac.go` and `internal/issues/ac_verify.go` only. | VAL-004, targeted VAL-009, no VAL for report/file-size hygiene | Developer | Task 2b completion report |
| EVID-04 | Task 3 completion report records exact required gate results, coverage output, populated `ac list` outputs in four formats, default JSON evidence, and empty-list transcript. | VAL-006, VAL-007 | Developer | Task 3 completion report |
| EVID-05 | Task 4 completion report records exact required gate results, coverage output, integration/build condition outcomes, reject->verify->pass transcript, and rebuild parity transcript. | VAL-008, VAL-009 | Developer | Task 4 completion report |
| EVID-06 | Review 1 checks task specs, reports, commits, in-window tests, command coverage, and validation tracking for Tasks 1-2. | VAL-001..VAL-005 | QA Reviewer | `qa/review-1.md` |
| EVID-07 | Review 2 checks task specs, reports, commits, in-window tests, command coverage, and validation tracking for Tasks 3-4. | VAL-006..VAL-009 | QA Reviewer | `qa/review-2.md` |
| EVID-08 | QA Lead final review reruns the phase-level deterministic gate and seeded fixture checklist with pasted command output/file-state evidence. | VAL-001..VAL-009 | QA Lead | `qa/final-review.md` |
| EVID-09 | Human checklist records one outcome each for ordered add+list, multi-role verify/list, and epic rejection in a real `.fid/` workspace. | VAL-010 | Human | `qa/human-checklist.md` |
| EVID-10 | Documentation review confirms `06_DESIGN_PLAN_CLI.md` examples and Phase 7 reports match implemented public output. | VAL-001..VAL-007, no VAL for doc hygiene | QA Lead | `qa/final-review.md` |
| EVID-11 | Fixup re-review checks Task 2b only: VAL-004, targeted VAL-009, committed Task 1 report proof, and scoped file-size evidence for `internal/issues/ac.go` and `internal/issues/ac_verify.go` only. | VAL-004, targeted VAL-009, no VAL for report/file-size hygiene | QA Reviewer | `qa/review-1-fixup.md` |

## 5. Task Command And Evidence Coverage

Every task completion report must record each required command as `PASS`, or
`BLOCKED` with reason in a blocker report. Missing, silently skipped, or
narrower substituted commands are Blocking for Phase Runner/QA Reviewer.

| Task | Required commands | Conditional commands | Required evidence coverage |
|---|---|---|---|
| Task 1 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go build -o /tmp/fid ./cmd/fid` if real CLI transcripts are generated for evidence. Because Task 1 requires pasted task and epic-reject transcripts, omission must explain the alternate transcript source. | Gate output; coverage thresholds `internal/issues >= 80%` and `cmd/fid >= 70%`; AC add success; ordinal read-back; epic rejection non-zero; unchanged journal/rows on rejection; stderr confirmation. |
| Task 2 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go build -o /tmp/fid ./cmd/fid` if real CLI transcripts are generated for evidence. Because Task 2 requires pasted verification transcripts, omission must explain the alternate transcript source. | Gate output; coverage thresholds `internal/issues >= 80%` and `cmd/fid >= 70%`; valid verify persisted fields; multi-role independence; fail->pass two rows and latest pass; invalid status zero append; invalid evidence-type and AC/issue mismatch tests. |
| Task 2b | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | Required: `go build -o /tmp/fid ./cmd/fid` to regenerate deterministic fail->pass and rebuild transcripts for the completion report. Integration command: none required by Task 2b unless the Developer places the rebuild/latest test behind `integration`; if so, `go test -tags integration -race ./...` becomes required and must be recorded as `PASS` or `BLOCKED`. | Gate output; coverage thresholds `internal/issues >= 80%` and `cmd/fid >= 70%`; same-AC same-role fail-only latest=fail then immediate fail->pass latest=pass with two retained rows; no `time.Sleep` ordering proof; delete DB + `fid rebuild` preserves latest-per-role; committed Task 1 report has hash/subject and no `Pending git commit`; `internal/issues/ac.go` and `internal/issues/ac_verify.go` are each `<= 400` lines. `internal/issues/issue.go` is CF-012 pre-existing debt and out of scope for Task 2b. |
| Task 3 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go build -o /tmp/fid ./cmd/fid` if real CLI transcripts are generated for evidence. Because Task 3 requires pasted format transcripts, omission must explain the alternate transcript source. | Gate output; coverage thresholds `internal/issues >= 80%` and `cmd/fid >= 70%`; JSON shape; latest-per-role only; four formats render; omitted `--format` equals JSON; empty list; issue-not-found error. |
| Task 4 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go test -tags integration -race ./...` if the end-to-end test is behind the `integration` build tag; `go build -o /tmp/fid ./cmd/fid` if real CLI transcripts are generated for evidence. Because Task 4 requires full end-to-end transcript, omission of build must explain the alternate transcript source. | Gate output; coverage thresholds `cmd/fid >= 70%`, `internal/issues >= 80%`, `internal/workflow >= 80%`; reject path names AC IDs/roles and leaves status/journal unchanged; success path after real `ac verify`; no journal fixtures; delete DB + `fid rebuild` parity for AC rows, verification rows, and latest-per-role read-back. |

## 6. Expected Tests

Task 1 expected tests:

- Engine test for successful `AddAcceptanceCriterion` projection with `AC-`
  ID, text, ordinal, and one `issue.ac_added` event.
- CLI tests for required flags, positive ordinal, unknown role, issue-not-found,
  stderr confirmation, and epic rejection before append.
- Read-back test for multiple ACs ordered by ordinal then id.

Task 2 expected tests:

- Engine test for successful `VerifyAcceptanceCriterion` projection with
  `ACV-` ID and all stored fields.
- Tests for multi-role independence, fail-only latest, fail->pass two-row
  history, invalid status, invalid evidence-type, AC not found, AC/issue
  mismatch, unknown role, and one `issue.ac_verified` event.

Task 2b expected tests:

- Deterministic same-AC same-role fail-only then immediate fail->pass coverage,
  with no `time.Sleep` and latest-per-role resolved by journal append order.
- Rebuild latest-per-role parity for a multi-write fixture after deleting the
  SQLite DB and running `fid rebuild`.
- No remaining `time.Sleep` or timestamp-boundary dependence in AC verification
  ordering tests.
- Regression tests for multi-role independence and field persistence may remain
  in the gate, but they do not expand the fixup-only validation scope beyond
  VAL-004 and targeted VAL-009.

Task 3 expected tests:

- Read-path test for design JSON shape with latest-per-role map.
- Tests for repeated verification latest-only output, empty AC list,
  issue-not-found, all four format dispatches, and default JSON parity.

Task 4 expected tests:

- End-to-end test over real commands for create -> ac add -> transition reject
  -> ac verify -> transition pass.
- Assertions for unmet AC IDs/roles in reject error, unchanged status and
  unchanged journal line count on rejection.
- Delete SQLite DB -> `fid rebuild` -> structural parity for AC rows,
  verification rows, and latest-per-role output.

## 7. Deterministic Gates

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task and phase-end |
| Lint/static analysis | `golangci-lint run ./...` | Exit 0 | Every task and phase-end |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | Exit 0 and writes `coverage.out` | Every task and phase-end |
| Coverage report | `go tool cover -func=coverage.out` | Exit 0 with package/function coverage; thresholds evidenced | Every task and phase-end after coverage profile exists |
| Integration | `go test -tags integration -race ./...` | Exit 0, or recorded N/A if no Phase 7 test is tagged integration | Task 4 when tagged; Task 2b only if the rebuild/latest test is tagged integration; always phase-end |
| Benchmarks | `go test -bench . -benchmem ./...` | Exit 0 | Phase-end |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | Exit 0 and `/tmp/fid` exists | Required for Task 2b transcripts; conditional task transcripts otherwise; always phase-end |
| Full phase check | `gofmt -l cmd internal` + `golangci-lint run ./...` + `go test -race -coverprofile=coverage.out ./...` + `go tool cover -func=coverage.out` + `go test -tags integration -race ./...` + `go test -bench . -benchmem ./...` + `go build -o /tmp/fid ./cmd/fid` | All applicable commands pass | Phase-end final review |

## 8. Command And Feasibility Validation

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | `gofmt -h`; path existence check for `cmd/` and `internal/`. | PASS | Help output includes `usage: gofmt [flags] [path ...]` and `-l list files whose formatting differs`; `find` shows `cmd/fid` and `internal/issues`, `internal/journal`, `internal/workflow`. |
| `golangci-lint run ./...` | `golangci-lint --version`; `golangci-lint run --help`; config file existence. | PASS | Version output: `golangci-lint has version 2.12.2`; help output says `Usage: golangci-lint run [flags]`; `.golangci.yml` exists. |
| `go test -race -coverprofile=coverage.out ./...` | `go version`; `go help testflag`; `go help build`; module file existence. | PASS | `go version go1.26.2 linux/amd64`; test help documents `-coverprofile`; build help documents `-race`; `go.mod` exists. |
| `go tool cover -func=coverage.out` | `go tool cover -h`; coverage artifact assumption check. | PASS | Help output documents `-func string output coverage profile information for each function`; `coverage.out` currently exists, and the required preceding `go test -race -coverprofile=coverage.out ./...` regenerates it. |
| `go test -tags integration -race ./...` | `go help build`; existing integration test file check. | PASS | Build help documents `-tags`; `cmd/fid/e2e_v01b_test.go` begins with `//go:build integration`, proving the tag pattern is in use. |
| `go test -bench . -benchmem ./...` | `go help testflag`; benchmark existence check. | PASS | Test help documents `-bench regexp` and `-benchmem`; `internal/journal/rebuild_bench_test.go` defines `BenchmarkRebuildWorkspace1kEvents` and `BenchmarkRebuildWorkspace10kEvents`. |
| `go build -o /tmp/fid ./cmd/fid` | `go help build`; `/tmp` writable root from execution environment; `cmd/fid/main.go` exists. | PASS | Build help documents `usage: go build [-o output] [build flags] [packages]` and `-o`; `/tmp` is writable per sandbox profile; `cmd/fid/main.go` exists. |
| Task 1 transcript fixture | Task spec and CLI design check. | PASS | Task 1 requires transcripts for add and epic-reject; design §3 defines `fid issue ac add` syntax and stderr confirmation. |
| Task 2 transcript fixture | Task spec and CLI design check. | PASS | Task 2 requires transcripts for pass verify, multi-role, fail->pass, invalid-status; design §3 defines `fid issue ac verify` syntax and allowed statuses/evidence types. |
| Task 2b deterministic transcript fixture | Task 2b spec, DEC-021, and CLI design check. | PASS | Task 2b requires deterministic fail->pass and rebuild latest-per-role transcripts; DEC-021 defines append-order latest semantics; design §3 keeps `fid issue ac verify` syntax unchanged. |
| Task 2b Task 1 committed-report check | `git show` command feasibility and Review 1 finding check. | PASS | `git show HEAD:project-docs/docs/phases/phase-7/tasks/task-1-done.md` is valid git object-path syntax after the fixup commit; Review 1 requires the committed artifact to contain the verified hash/subject and no `Pending git commit`. |
| Task 2b scoped AC split file-size check | `wc` command feasibility with path existence; DEC-022 scope check. | PASS | `wc -l internal/issues/ac.go internal/issues/ac_verify.go` uses existing AC write-path split files; DEC-022 scopes Task 2b acceptance to those files only, each `<= 400` lines. `internal/issues/issue.go` is explicitly excluded as CF-012 pre-existing debt. |
| Task 3 transcript fixture | Task spec and CLI design check. | PASS | Task 3 requires one populated `ac list` per format plus empty-list; design §3 defines `fid issue ac list` JSON shape and format flag. |
| Task 4 real-data/rebuild fixture | Task spec, existing e2e pattern, and TRD workflow rule check. | PASS | Task 4 points to `cmd/fid/e2e_v01b_test.go`; TRD task workflow lists `all_acs_verified{"roles":["developer"]}` on `in_progress -> dev_complete` and `{"roles":["qa_lead"]}` on `qa_review -> qa_pass`. |
| Human dogfood workspace | Human-checklist template and normal `fid init` workflow feasibility. | PASS | Human checklist template exists; Phase 7 commands are ordinary CLI operations against a real `.fid/` workspace after tasks complete. |
| Hidden/held-out scenario boundary | Phase brief and quality plan review. | PASS | Phase brief says held-out scenario harness remains a scheduled commitment and Task 4 may recommend only; no hidden scenario content is needed for this plan. |

No referenced command is a placeholder. The QA plan is not blocked.

## 9. Deterministic Checklist

Items here are executed or verified by QA Lead during final review and must
receive reproducible evidence in `final-review.md`.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Add two or more ACs to a task/subtask and prove persisted `AC-` IDs, text, ordinal values, and ordinal-then-id read-back `(qa-lead)` | VAL-001 | Seeded CLI transcript plus JSON/read-back excerpt. |
| QA-02 | Run the same `ac add` shape against an epic and a task/subtask, proving epic rejection is non-zero with kind-constraint error and zero journal/row append while task/subtask succeeds `(qa-lead)` | VAL-002 | CLI transcript, pre/post journal line counts or DB row counts, and read-back. |
| QA-03 | Verify one AC as `developer=pass` and `qa_lead=fail`, proving latest-per-role retains both roles independently `(qa-lead)` | VAL-003 | CLI transcript and JSON/read-back excerpt showing both roles and statuses. |
| QA-04 | Verify the same AC and same role as fail then immediate pass, proving two retained verification rows and latest status pass by append order; also show fail-only latest before the pass and no `time.Sleep`/timestamp-boundary dependence `(qa-lead)` | VAL-004 | DB/query or test output showing two rows plus latest-per-role excerpts before and after pass, plus test/transcript evidence that the writes are immediate. |
| QA-05 | Prove valid verify stores role, status, evidence_type, evidence_ref, notes, verified_at, and invalid `--status` rejects before append `(qa-lead)` | VAL-005 | CLI transcript plus DB/query or JSON output for stored fields and pre/post append counts for invalid status. |
| QA-06 | Run `fid issue ac list --format json` after multi-role and repeated verifications, proving one latest entry per role in the documented shape and no raw history entries `(qa-lead)` | VAL-006 | JSON output excerpt and, where needed, DB row count proving hidden retained history is not duplicated in list output. |
| QA-07 | Run `fid issue ac list` with omitted format and with `json`, `toon`, `markdown`, and `yaml`; prove all render and omitted format is JSON `(qa-lead)` | VAL-007 | CLI output excerpts for each format and structural/default JSON comparison evidence. |
| QA-08 | Exercise real-command `all_acs_verified` gating: transition rejects before required pass with unmet AC IDs/roles and no state/journal change, then succeeds after real `ac verify` pass `(qa-lead)` | VAL-008 | Integration test output or seeded CLI transcript with state/journal pre/post proof. |
| QA-09 | Snapshot AC and verification state, delete SQLite DB, run `fid rebuild`, and prove projected AC rows, verification rows, and latest-per-role read-back are identical; Task 2b must include targeted latest-per-role parity for the DEC-021 append-order tie-break `(qa-lead)` | VAL-009 | Test output or scripted transcript plus before/after structural comparison. |
| QA-10 | Compare each task completion report against its task spec and this QA plan command table, including conditional command handling and transcript obligations `(qa-lead)` | VAL-001..VAL-009 | Report excerpts showing PASS/BLOCKED for every required and triggered conditional command. |
| QA-11 | Verify coverage thresholds from `go tool cover -func=coverage.out`: `cmd/fid >= 70%`, `internal/issues >= 80%`, `internal/journal >= 80%`, `internal/policy >= 80%`, and `internal/workflow >= 80%` where applicable `(qa-lead)` | no VAL | Coverage output excerpt. |
| QA-12 | Run the full phase-level command gate in order `(qa-lead)` | no VAL | Command output excerpts for format, lint, coverage test, coverage report, integration, benchmarks, and build. |
| QA-13 | Review documentation updates for CLI examples and Phase 5 fixture supersession note `(qa-lead)` | VAL-001..VAL-009, no VAL for doc hygiene | File-state references to updated docs/reports and any no-change rationale. |
| QA-14 | Review security/data-integrity risks for SQL injection, command injection, unsafe input handling, journal poisoning, and append-only history preservation `(qa-lead)` | VAL-002, VAL-004, VAL-005, VAL-008, VAL-009 | Review notes with file references and command/test evidence where available. |
| QA-15 | Verify the committed Task 1 completion report contains the verified commit hash and subject, and does not contain `Pending git commit` `(qa-lead)` | no VAL | Output or excerpt from `git show HEAD:project-docs/docs/phases/phase-7/tasks/task-1-done.md`. |
| QA-16 | Verify the Task 2b AC write-path split files `internal/issues/ac.go` and `internal/issues/ac_verify.go` are each `<= 400` lines; confirm `internal/issues/issue.go` remains excluded as CF-012 pre-existing debt `(qa-lead)` | no VAL | Output from `wc -l internal/issues/ac.go internal/issues/ac_verify.go` with both files at or below 400 lines, plus review note citing DEC-022 / CF-012 for the `issue.go` exclusion. |

## 10. Mid-Phase Review Cadence

Default project cadence is every three completed tasks. Phase 7 overrides this
with two bounded QA Reviewer windows because the four tasks form two natural
contracts: write paths first, read/closure second.

Resumed fixup cadence: because Review 1 already returned `BLOCK`, Task 2b is a
single fixup window. The Phase Runner must trigger the fixup re-review
immediately after the Task 2b completion report and commit verification, and
before QA Lead final review. The fixup re-review is limited to Task 2b:
VAL-004, targeted VAL-009 rebuild parity, Task 1 committed report
finalization, and the scoped file-size check for `internal/issues/ac.go` and
`internal/issues/ac_verify.go` only. `internal/issues/issue.go` is CF-012
pre-existing debt and remains out of scope.

| Review | Window | Trigger | Output | Required focus |
|---|---|---|---|---|
| Review 1 | Tasks 1-2 | Immediately after Task 2 completion report and commit verification, before Task 3 starts. | `project-docs/docs/phases/phase-7/qa/review-1.md` | Validate write-path evidence for VAL-001..VAL-005, command coverage, validate-before-append discipline, append-only verification history, and public transcript obligations. |
| Review 1 fixup | Task 2b only | Immediately after Task 2b completion report and commit verification, before QA Lead final review. | `project-docs/docs/phases/phase-7/qa/review-1-fixup.md` | Validate deterministic no-sleep VAL-004 proof, targeted VAL-009 rebuild latest-per-role parity, exact command coverage including required build, any triggered integration command, Task 1 committed-report finalization, and scoped file-size evidence showing only `internal/issues/ac.go` and `internal/issues/ac_verify.go` are each `<= 400` lines. |
| Review 2 | Tasks 3-4 | Immediately after Task 4 completion report and commit verification, before QA Lead final review. | `project-docs/docs/phases/phase-7/qa/review-2.md` | Validate read/format evidence for VAL-006..VAL-007 and closure evidence for VAL-008..VAL-009, including conditional integration/build commands and rebuild parity. |

Each review uses `project-docs/templates/mid-phase-review-template.md`, compares
completion reports to task-spec and QA-plan commands, and records whether each
in-window VAL is `ON TRACK`, `NOT ON TRACK`, or `DEFERRED`. Missing evidence for
an in-window VAL is Blocking unless the task report records a valid `BLOCKED`
reason.

## 11. Human Verification Scope

Items here require human judgment and must be copied into
`project-docs/docs/phases/phase-7/qa/human-checklist.md` during final review.
One row is used per VAL-010 observable.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HV-010a | In a real `.fid/` workspace, add at least two ordered ACs to a task and list/read them back; judge whether the add confirmation and ordered list are clear and behave as expected `(human)` | VAL-010 | Human outcome plus notes on ordering, command clarity, and any friction. |
| HV-010b | In the same real workspace, verify one AC from at least two roles and list latest-per-role; judge whether the multi-role latest status is understandable and matches expectation `(human)` | VAL-010 | Human outcome plus notes on role labels, latest status clarity, and any ambiguity. |
| HV-010c | In the same real workspace, attempt to add an AC to an epic; judge whether the rejection is clear, actionable, and matches the product expectation that epics cannot hold ACs `(human)` | VAL-010 | Human outcome plus error text excerpt or notes on clarity/actionability. |

## 12. Acceptance Criteria For This QA Plan

- Every VAL-001..VAL-010 has at least one planned evidence row.
- VAL-010 is split into exactly three human-checklist observables:
  ordered add+list, multi-role verify/list, and epic rejection.
- Mid-phase review windows are bounded and named for all four tasks.
- The resumed run is explicitly fixup-only and adds Task 2b review coverage
  without reopening Task 1-4 validation scope.
- Every task's exact required commands and conditional commands are named.
- Task 2b names required build evidence, conditional integration handling,
  Task 1 committed-report finalization, and the scoped file-size check for
  `internal/issues/ac.go` and `internal/issues/ac_verify.go` only, with
  `internal/issues/issue.go` recorded as CF-012 pre-existing debt out of scope.
- Command feasibility is validated with concrete local help output, file-state
  checks, or documented project usage.
- Deterministic and human checks are split by judgment type.
- No product code, task specs, or validation-contract edits are required.

## 13. Phase-End Decision Rule

The phase can proceed to final review only when:

- Task 1 through Task 4 completion reports exist, contain required headings,
  and record required command outcomes as `PASS` or `BLOCKED`.
- If running the resumed fixup path, Task 2b completion report exists, contains
  required headings, and records required command outcomes as `PASS` or
  `BLOCKED`.
- Review 1 and Review 2 exist with `VERDICT: CONTINUE`, or any `BLOCK` verdict
  has been resolved through the project escalation/triage flow.
- Review 1 fixup exists with `VERDICT: CONTINUE` before final QA when Task 2b
  is the Review 1 blocker resolution path.
- Every deterministic QA item has reproducible evidence available for
  `final-review.md`.
- Every human item remains human-only and is carried into
  `human-checklist.md`.
- No unresolved blocking defects or open escalations remain.
- Required documentation updates or no-change rationales are recorded.
