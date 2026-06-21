# Phase 10 QA Plan - Permissions Matrix + Policy

> Phase: 10 - Permissions matrix + policy + CF-003
> Refreshed: 2026-06-21
> QA Lead: Codex
>
> Location:
> `project-docs/docs/phases/phase-10/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 10 introduces the v0.1 permissions matrix as a cross-cutting gate on
every implemented write command. It reconciles CF-003 role-discovery corruption
handling, adds the in-code matrix, `policy.Authorize`, read-only
`fid policy show`, and retrofits enforcement onto existing write surfaces. No
schema migration or new product entity is in scope.

Sources used:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/01_PRD.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/03_TRD_SCHEMA.md` section 13
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 11
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/docs/phases/phase-10/decisions.md`
- `project-docs/docs/phases/phase-10/qa/validation-contract.md`
- `project-docs/docs/phases/phase-10/tasks/task-1.md`
- `project-docs/docs/phases/phase-10/tasks/task-2.md`
- `project-docs/docs/phases/phase-10/tasks/task-3.md`

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: YES.
  PRD-ROLE-001/002 map to VAL-005; PRD-ROLE-003 maps to VAL-004/007;
  PRD-ROLE-004 maps to VAL-006; PRD-ROLE-005 maps to VAL-002; TRD-ARCH-010 maps
  to VAL-002/003/004/007/008; TRD-ARCH-002/CF-003 maps to VAL-001; CF-007 maps
  to VAL-004/005.
- Every `VAL-NNN` has one validator role: YES. VAL-001..008 are QA Lead;
  VAL-009 is Human.
- Every `VAL-NNN` has a concrete evidence method: YES, with CLI, unit,
  integration, source-state, command-output, or human-observation evidence named
  below.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: YES. VAL-009 is human-only and is fanned into
  HUMAN-09a..HUMAN-09d below.
- Conditional or binary `VAL-NNN` assertions require both qualifying and
  non-qualifying states: YES. VAL-001..008 require both positive and negative
  evidence where applicable; positive-only evidence cannot pass them.
- Parity or "matches reference" `VAL-NNN` assertions have deterministic
  conformance proof when available: YES. VAL-002 requires format-to-format data
  equivalence and implemented action coverage; VAL-008 requires show-vs-enforce
  evidence for both a permitted and denied cell.
- Revisions requested from Tech Lead: none. The contract is executable as
  drafted.

## 3. Requirement Coverage

| Requirement | Covered By | Verification |
|---|---|---|
| PRD-ROLE-001 - `--by` required on every write | VAL-005 | Missing `--by` reject with zero append. |
| PRD-ROLE-002 - unknown roles rejected clearly | VAL-001, VAL-005 | Sound role discovery; unknown-role reject listing known roles with zero append. |
| PRD-ROLE-003 - denied permissions rejected with role/action/not-permitted error | VAL-003, VAL-004, VAL-007, VAL-008, VAL-009 | Unit and CLI denial evidence, parity proof, human readability judgment. |
| PRD-ROLE-004 - `human` escape valve without verifier impersonation | VAL-006, VAL-009 | CLI success where non-human is denied; AC verification persisted as `human`; human judgment. |
| PRD-ROLE-005 - display active permissions matrix | VAL-002, VAL-008, VAL-009 | `fid policy show` formats, action coverage, show-vs-enforce parity, human readability. |
| TRD-ARCH-002 - malformed journal line is corruption | VAL-001 | Role-discovery corrupt-tail unit/integration and CLI evidence. |
| TRD-ARCH-010 - v0.1 code constant; every write action mapped; denied combos error | VAL-002, VAL-003, VAL-004, VAL-007, VAL-008 | Matrix completeness, `Authorize`, every write surface, derived action keys, parity. |
| Design CLI section 11 - `fid policy show [--format json\|toon\|markdown\|yaml]` | VAL-002 | All four formats render same action-role booleans. |
| CF-003 - corrupt-tail role scan debt | VAL-001 | Reconciled error behavior. |
| CF-007 - validate before append | VAL-004, VAL-005, VAL-007 | Every permission and role-identity denial proves non-zero exit plus zero new journal lines/rows. |
| DEC-A..E | VAL-001..008 | Task evidence and completion reports cite the applied decision. |

## 4. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | CF-003 corrupt-tail reconciliation | VAL-001 | Medium | Review must confirm corruption is reported for tail and first-line cases, while clean journals resolve. |
| Task 2 | Matrix, `Authorize`, and `fid policy show` | VAL-002, VAL-003 | Medium | Review must confirm one single-source matrix feeds both display and authorization. |
| Task 3 | Retrofit enforcement onto every write command | VAL-004, VAL-005, VAL-006, VAL-007, VAL-008 | High | Cross-cutting write-surface change; validate-before-append and show-vs-enforce parity are standing checks. |
| Phase end | Human dogfood and final deterministic audit | VAL-001..009 | High | QA Lead executes deterministic checklist; Human executes only HUMAN-09a..HUMAN-09d. |

## 5. Required Evidence Strategy

| Evidence ID | Required Evidence | Covers | Owner | Produced By |
|---|---|---|---|---|
| EVID-01 | Unit/integration output proving corrupt-tail role discovery returns corruption naming the journal file, first-line corruption still errors, and clean role journals resolve. | VAL-001 | Developer, QA Reviewer, QA Lead | Task 1 report, review-1, final review |
| EVID-02 | CLI output from a write command against a corrupt-tail workspace: non-zero exit and corruption error naming the journal file. | VAL-001 | Developer, QA Reviewer, QA Lead | Task 1 report, review-1, final review |
| EVID-03 | `fid policy show` output or tests for json, toon, markdown, and yaml proving the same action-role booleans and every implemented v0.1a-c write action. | VAL-002 | Developer, QA Reviewer, QA Lead | Task 2 report, review-2, final review |
| EVID-04 | Unit output proving every implemented action has exactly one matrix entry and `Authorize` permits one allowed cell while denying one denied cell with `PermissionDeniedError{Role, Action}`. | VAL-003 | Developer, QA Reviewer, QA Lead | Task 2 report, review-2, final review |
| EVID-05 | CLI/file-state proof for a known-role denied write: non-zero exit, error names role + action + not permitted, zero new journal lines and zero new rows; permitted counterpart succeeds with exactly one event. | VAL-004 | Developer, QA Reviewer, QA Lead | Task 3 report, review-3, final review |
| EVID-06 | CLI/file-state proof for missing `--by` and unknown role: each rejects before append, preserves distinct error text, and permitted counterpart proceeds. | VAL-005 | Developer, QA Reviewer, QA Lead | Task 3 report, review-3, final review |
| EVID-07 | CLI/read-back proof that `--by human` succeeds where a non-human role is denied, and human AC verification persists or reads back `_role = human`. | VAL-006 | Developer, QA Reviewer, QA Lead | Task 3 report, review-3, final review |
| EVID-08 | CLI proof for derived action keys: AC verify permits `qa_reviewer` and denies `workflow_runner`; transition to `dev_complete` permits `developer` and denies `qa_reviewer`; each denial has zero append. | VAL-007 | Developer, QA Reviewer, QA Lead | Task 3 report, review-3, final review |
| EVID-09 | Show-vs-enforce parity proof for one permitted and one denied cell: `fid policy show` boolean equals actual write outcome. | VAL-008 | Developer, QA Reviewer, QA Lead | Task 3 report, review-3, final review |
| EVID-10 | Completion-report command evidence matching the exact commands in each task spec; missing, skipped, or narrower commands are Blocking unless recorded BLOCKED with reason. | VAL-001..008, no VAL gates | Developer, QA Reviewer, QA Lead | Every task review and final review |
| EVID-11 | Coverage report excerpts showing COV-001 package thresholds and COV-002 `cmd/fid` threshold. | no VAL | Developer, QA Reviewer, QA Lead | Every task report and final review |
| EVID-12 | Human checklist results for policy matrix readability and denial-error actionability, one observable per row. | VAL-009 | Human | Human checklist |

## 6. VAL Verification Map

| VAL | Deterministic Verification | Mid-Review Verification | Human Verification |
|---|---|---|---|
| VAL-001 | QA-01, QA-02; Task 1 command evidence | review-1 | None |
| VAL-002 | QA-03, QA-04; Task 2 command evidence | review-2 | None |
| VAL-003 | QA-05; Task 2 command evidence | review-2 | None |
| VAL-004 | QA-06, QA-12; Task 3 command evidence | review-3 | None |
| VAL-005 | QA-07; Task 3 command evidence | review-3 | None |
| VAL-006 | QA-08; Task 3 command evidence | review-3 | None |
| VAL-007 | QA-09, QA-10, QA-12; Task 3 command evidence | review-3 | None |
| VAL-008 | QA-11; Task 3 command evidence | review-3 | None |
| VAL-009 | None; deterministic checks only support setup validity | Final review confirms checklist is generated | HUMAN-09a, HUMAN-09b, HUMAN-09c, HUMAN-09d |

## 7. Acceptance Criteria

- VAL-001..008 each have deterministic evidence in the relevant task report,
  mid-phase review, and final review.
- VAL-009 is not passed by code reading or QA inference; it requires completed
  HUMAN-09a..HUMAN-09d rows.
- `fid policy show` and `Authorize` remain cross-surface consistent. VAL-008
  cannot pass unless both a displayed permitted cell succeeds and a displayed
  denied cell fails.
- Every permission denial and role-identity denial checked under VAL-004/005/007
  must show non-zero exit, a role/action-named or identity-specific error as
  applicable, and zero new journal lines/rows.
- Task completion reports must record every exact command required by their task
  specs as PASS or BLOCKED with reason. Missing, silently skipped, or narrower
  substituted commands are Blocking.
- The phase-end final review records PASS, FAIL, BLOCKED, or CARRY-FORWARD for
  every VAL-001..009.

## 8. Expected Tests

| Test Scope | Expected Coverage | Related VAL |
|---|---|---|
| `internal/policy` role discovery tests | Corrupt first line, corrupt tail after valid role event, un-upgradeable tail event, clean multi-line role journal. | VAL-001 |
| CLI corrupt-tail workspace test | Write command exits non-zero and names journal corruption. | VAL-001 |
| `internal/policy` matrix tests | Implemented action set completeness and exactly-one-entry checks against v0.1a-c actions. | VAL-003 |
| `internal/policy` authorize tests | Allowed cell returns nil; denied cell returns `PermissionDeniedError`; missing/unknown role errors preserved. | VAL-003, VAL-005 |
| `cmd/fid` policy show tests | json default, toon, markdown, yaml render same rows; no `--by`; no journal/activity write. | VAL-002 |
| Write-surface enforcement tests | Existing write commands call `Authorize` with derived action keys and update fixtures to permitted combinations per DEC-E. | VAL-004, VAL-007 |
| Validate-before-append CLI tests | Denied permission, missing `--by`, unknown role, denied AC verify, and denied transition all append zero journal lines/rows. | VAL-004, VAL-005, VAL-007 |
| Human escape valve tests | Human succeeds on a write denied to a non-human; human AC verification persists as `human`. | VAL-006 |
| Show-vs-enforce parity tests | One permitted cell and one denied cell compare `fid policy show` data to actual write outcomes. | VAL-008 |
| Documentation/decision checks | `06_DESIGN_PLAN_CLI.md` and Phase 10 decisions reflect applied behavior or completion reports record verified-no-change. | VAL-002, VAL-004..008 |

## 9. Deterministic Gates

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task and phase end |
| Lint/static analysis | `golangci-lint run ./...` | Exit 0 | Every task and phase end |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | Exit 0 and writes `coverage.out` | Every task and phase end |
| Coverage report | `go tool cover -func=coverage.out` | Exit 0 and threshold evidence for COV-001/COV-002 | Every task and phase end |
| Integration | `go test -tags integration -race ./...` | Exit 0 | Required for Task 1, Task 3, and phase end; Task 2 if its test harness uses integration tags |
| Benchmarks | `go test -bench . -benchmem ./...` | Exit 0 | Phase end; Task 3 if any benchmark fixture role changes |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | Exit 0 and binary produced at `/tmp/fid` | Required for Task 1, Task 2, Task 3, and phase end |
| Full phase check | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go test -bench . -benchmem ./...`; `go build -o /tmp/fid ./cmd/fid` | All PASS or BLOCKED with reason | Phase-end final review |

End-to-end or held-out scenario validation remains N/A for Phase 10 because the
accepted specs have not selected a scenario harness or command.

## 10. Command And Feasibility Validation

| Referenced Step | Validation Method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Tool help plus directory existence. | PASS | `gofmt -h` documents `-l` and path args; `cmd/` and `internal/` exist. |
| `golangci-lint run ./...` | Tool help. | PASS | `golangci-lint --help` lists `run`; `golangci-lint run --help` documents package-wide linting flags. |
| `go test -race -coverprofile=coverage.out ./...` | Go test help plus source/test presence. | PASS | `go help testflag` documents `-race`, `-coverprofile`, and package-list mode; `cmd/` and `internal/` contain `*_test.go` files. |
| `go tool cover -func=coverage.out` | Tool usage output. | PASS | `go tool cover -h` documents `-func` for coverage profiles produced by `go test -coverprofile`. Help exits non-zero after printing usage, which is acceptable for syntax validation only. |
| `go test -tags integration -race ./...` | Go test/build flag help plus task specs. | PASS | `go help build` documents `-tags`; `go help testflag` documents `-race`; Task 1 and Task 3 require integration checks. |
| `go test -bench . -benchmem ./...` | Go test help. | PASS | `go help testflag` documents `-bench` and `-benchmem`; `internal/journal/rebuild_bench_test.go` exists. |
| `go build -o /tmp/fid ./cmd/fid` | Go build help plus package directory existence. | PASS | `go help build` documents `-o output` and package args; `cmd/fid/` exists; `/tmp` is writable in the sandbox profile. |
| Phase 10 QA artifact path | File-state check. | PASS | `project-docs/docs/phases/phase-10/qa/` exists and is writable under the workspace root. |
| Phase 10 task specs | File-state check. | PASS | `project-docs/docs/phases/phase-10/tasks/task-1.md`, `task-2.md`, and `task-3.md` exist. |
| Phase 10 validation contract | File-state check. | PASS | `project-docs/docs/phases/phase-10/qa/validation-contract.md` exists and defines VAL-001..009. |
| `fid policy show` format contract | Spec inspection. | PASS | `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 11 names `--format json|toon|markdown|yaml` and the JSON row shape. |
| Matrix cell reference | Spec inspection. | PASS | `project-docs/specs/03_TRD_SCHEMA.md` section 13 defines canonical action keys, role columns, and `human` behavior. |
| Validate-before-append proof | Task/contract inspection. | PASS | VAL-004/005 and Task 3 require zero new journal lines/rows for denials; QA checklist below requires file-state proof. |
| Show-vs-enforce parity proof | Task/contract inspection. | PASS | VAL-008 and Task 3 require one permitted and one denied cell compared between `fid policy show` and actual write outcome. |
| Human checklist feasibility | Template and contract inspection. | PASS | `project-docs/templates/human-checklist-template.md` supports one observable per row; VAL-009 is human-only and fanned below. |
| HUMAN-09a | Template/contract inspection. | PASS | The step is a perception check on matrix readability in a real `.fid/` workspace; no deterministic substitute is claimed. |
| HUMAN-09b | Template/contract inspection. | PASS | The step is a perception check on JSON scanability only; deterministic JSON validity remains QA Lead evidence under VAL-002. |
| HUMAN-09c | Template/contract inspection. | PASS | The step is a human actionability judgment on a permission-denied error; deterministic stderr content remains QA Lead evidence under VAL-004/007. |
| HUMAN-09d | Template/contract inspection. | PASS | The step is a human judgment that the error suggests a usable next step; deterministic command success/failure cannot prove this. |
| Completion-report command comparison | Workflow/template inspection. | PASS | `08_QUALITY_PLAN.md` section 5.1 and mid-review template require comparing reports against exact task-spec commands. |
| Placeholder command audit | Text search. | PASS | `rg` over AGENTS, 08_QUALITY_PLAN, Phase 10 tasks, and validation contract found no placeholder required command for Phase 10; only AGENTS explanatory TBD marker text appeared. |

## 11. Deterministic Checklist

Every row must be executed or verified by QA Lead at final review unless it is
explicitly marked BLOCKED with reason. PASS without cited output, file state, or
artifact evidence is invalid.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Confirm Task 1 completion evidence includes corrupt-tail role discovery reported as corruption naming the journal file, first-line corruption still errors, and a clean role journal resolves `(qa-lead)` | VAL-001 | Task 1 report excerpts, test output, and review-1 status. |
| QA-02 | Reproduce or cite CLI evidence for a write against a corrupt-tail workspace: non-zero exit and corruption error naming the journal file `(qa-lead)` | VAL-001 | CLI output excerpt and journal fixture/setup notes. |
| QA-03 | Confirm `fid policy show` renders json, toon, markdown, and yaml with the same action-role booleans and every implemented v0.1a-c write action `(qa-lead)` | VAL-002 | Command/test output for all four formats plus implemented-action coverage evidence. |
| QA-04 | Confirm `fid policy show` is read-only: no `--by`, no permission check, no new journal line, and no `activity_log` row `(qa-lead)` | VAL-002 | CLI output or test log plus before/after journal and DB row counts. |
| QA-05 | Confirm every implemented write action has exactly one matrix entry and `Authorize` permits allowed/denies denied with `PermissionDeniedError{Role, Action}` `(qa-lead)` | VAL-003 | Unit test output and matrix/action-set evidence. |
| QA-06 | Verify a known-role denied write exits non-zero, names role + action + not permitted, and appends zero journal lines/rows; verify the permitted counterpart appends exactly one event `(qa-lead)` | VAL-004 | CLI output plus before/after journal line counts and DB row counts. |
| QA-07 | Verify missing `--by` rejects with `--by is required`, unknown role rejects with known roles listed, both append zero journal lines/rows, and a permitted role proceeds `(qa-lead)` | VAL-005 | CLI output plus before/after journal line counts and DB row counts. |
| QA-08 | Verify `--by human` succeeds on a write denied to a non-human role and human AC verification persists/read backs `_role = human` `(qa-lead)` | VAL-006 | CLI output and read-back JSON/toon/DB evidence. |
| QA-09 | Verify AC-verify action derivation: `qa_reviewer` permitted and `workflow_runner` denied, with the denied path naming `issue.ac.verify.workflow_runner` or the derived denied action and appending zero journal lines/rows `(qa-lead)` | VAL-007 | CLI output plus before/after journal and DB counts. |
| QA-10 | Verify transition action derivation: `developer` may transition to `dev_complete`, `qa_reviewer` is denied for `issue.transition.dev_complete`, and the denied path appends zero journal lines/rows `(qa-lead)` | VAL-007 | CLI output plus before/after journal and DB counts. |
| QA-11 | Execute the standing show-vs-enforce parity check: one `fid policy show` permitted cell must succeed when invoked, and one displayed denied cell must fail when invoked `(qa-lead)` | VAL-008 | `fid policy show` output excerpt, corresponding write output, and outcome comparison. |
| QA-12 | Confirm Task 3 completion report lists every fixture role change made under DEC-E and that each moved fixture still asserts the original behavior with a permitted role `(qa-lead)` | VAL-004, VAL-007, no VAL | Completion-report excerpt and review-3 assessment. |
| QA-13 | Compare every completion-report command result against each task spec's exact command list; flag missing, silently skipped, or narrower substituted commands as Blocking `(qa-lead)` | VAL-001..008, no VAL | Task reports, review files, and command result excerpts. |
| QA-14 | Run or cite the full phase-level command gate exactly as listed in section 9 `(qa-lead)` | no VAL | Final-review command output excerpts for format, lint, coverage, integration, benchmarks, and build. |
| QA-15 | Verify COV-001 and COV-002 thresholds from `go tool cover -func=coverage.out`: `internal/issues`, `internal/journal`, `internal/policy`, `internal/workflow` >= 80%, and `cmd/fid` >= 70% `(qa-lead)` | no VAL | Coverage report excerpts with package totals or equivalent function/package aggregation. |
| QA-16 | Review documentation updates or verified-no-change notes for `06_DESIGN_PLAN_CLI.md`, `01_PRD.md`, and Phase 10 decisions `(qa-lead)` | VAL-002, VAL-004..008, no VAL | Completion-report docs section and diff/file-state evidence. |
| QA-17 | Confirm no schema migration was added for Phase 10 `(qa-lead)` | no VAL | Git diff/file-state evidence for `migrations/` and schema docs. |

## 12. Mid-Phase Review Cadence

This phase uses a review-after-each-task cadence instead of the default
after-three-tasks cadence because each task is a distinct validation window and
Task 3 has high cross-surface blast radius.

| Review | Timing | Covered Tasks | VAL Focus | Required Review Actions |
|---|---|---|---|---|
| review-1 | After Task 1 completion, before Task 2 starts | Task 1 | VAL-001 | Compare Task 1 report with exact commands; confirm corrupt-tail and clean-journal evidence. |
| review-2 | After Task 2 completion, before Task 3 starts | Task 2 | VAL-002, VAL-003 | Confirm matrix completeness, `Authorize`, `fid policy show` formats, and no duplicated display/enforcement source. |
| review-3 | After Task 3 completion, before final review | Task 3 | VAL-004, VAL-005, VAL-006, VAL-007, VAL-008 | Confirm every write-surface gate, validate-before-append denials, human escape valve, action derivation, and show-vs-enforce parity. |

Each review uses `project-docs/templates/mid-phase-review-template.md`, writes
`project-docs/docs/phases/phase-10/qa/review-1.md`, `review-2.md`, or
`review-3.md`, and records CONTINUE or BLOCK. A BLOCK must stop the runner
unless triaged under the documented workflow.

## 13. Human Verification Scope

Only these rows require human perception, dogfooding, or product-owner judgment.
They must be copied into
`project-docs/docs/phases/phase-10/qa/human-checklist.md` at final review.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-09a | In a real `.fid/` workspace, inspect `fid policy show --format markdown` or the preferred human-readable format and judge whether the role x action matrix is readable enough to find an action and compare roles `(human)` | VAL-009 | Human observation and notes naming the format used. |
| HUMAN-09b | In the same workspace, inspect `fid policy show --format json` only for scanability of action names and role booleans, not JSON validity `(human)` | VAL-009 | Human observation; deterministic JSON validity remains QA Lead work. |
| HUMAN-09c | Trigger a permission-denied write and judge whether the error is actionable for an operator: it names the role, names the action, and says the role is not permitted `(human)` | VAL-009 | Human observation with the command attempted and stderr excerpt if practical. |
| HUMAN-09d | Judge whether the denied-write error gives enough signal to choose a next step, such as switching to a permitted role, checking `fid policy show`, or escalating a permission mismatch `(human)` | VAL-009 | Human notes on the next step the message implies. |

## 14. Phase-End Decision Rule

The phase can proceed to final review only when:

- Required deterministic gates pass or are recorded BLOCKED with escalation.
- Every VAL-001..009 has planned evidence and a final-review outcome slot.
- Every `(qa-lead)` item has reproducible evidence or is marked FAIL/BLOCKED.
- HUMAN-09a..HUMAN-09d are the only human-judgment checks and are present in the
  generated human checklist.
- The two standing checks are evidenced: VAL-008 show-vs-enforce parity, and
  VAL-004/005 validate-before-append on every permission or identity denial.
- No unresolved blocking findings or untriaged command anomalies remain.
- Documentation updates or verified-no-change notes are recorded for impacted
  specs and Phase 10 decisions.
