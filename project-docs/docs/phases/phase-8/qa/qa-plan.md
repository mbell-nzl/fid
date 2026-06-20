# Phase 8 QA Plan - Task Completions

> Phase: 8 - Task completions (v0.1c)
> Refreshed: 2026-06-20
> QA Lead: Codex

## 1. Scope And Source Documents

Phase 8 ships `fid issue complete`, completion attempt numbering and
supersession, gate snapshotting, real `has_completion_report` gating evidence,
delete-DB rebuild parity, and the CF-012 `internal/issues/issue.go` split.

Sources:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/docs/phases/phase-8/qa/validation-contract.md`
- `project-docs/docs/phases/phase-8/tasks/task-1.md`
- `project-docs/docs/phases/phase-8/tasks/task-2.md`
- `project-docs/docs/phases/phase-8/tasks/task-3.md`
- `project-docs/docs/phases/phase-8/tasks/task-4.md`

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: YES.
- Every `VAL-NNN` has one validator role: YES.
- Every `VAL-NNN` has a concrete evidence method: YES.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: YES. `VAL-008` is fanned out below.
- Conditional or binary assertions require both states: YES for `VAL-002`,
  `VAL-003`, `VAL-004`, `VAL-005`, and `VAL-006`.
- Parity assertions have deterministic conformance proof where available: YES
  for `VAL-007` and `VAL-009`.
- Revisions requested from Tech Lead: none.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | `fid issue complete` core write | VAL-001, VAL-002, VAL-003, VAL-004 | High | New public CLI write path; validate-before-append and public transcript evidence required. |
| Task 2 | Gate snapshot | VAL-005 | Medium | Fixture-seeded `gate_runs` allowed by DEC-C; must also prove no-gate state. |
| Task 3 | Real `has_completion_report` gating + rebuild parity | VAL-006, VAL-007 | High | Replaces Phase 5 fixture evidence; reject path must show zero state and journal change. |
| Task 4 | CF-012 `issue.go` split | VAL-009 | Medium | Refactor-only; must prove size, behavior preservation, and no weakened tests. |

## 4. VAL Coverage Map

| VAL | Planned verification | Owner/window | Evidence expectation |
|---|---|---|---|
| VAL-001 | Deterministic QA step QA-02; Review 1 confirms task evidence | QA Lead; mid-phase review after Task 3 | CLI/test output showing `TC-` row, stored fields, and JSON read-back. |
| VAL-002 | QA-03; Review 1 confirms task evidence | QA Lead; mid-phase review after Task 3 | Two real completions with attempts 1/2, prior row retained as `superseded`, new row `dev_complete`. |
| VAL-003 | QA-04; Review 1 confirms task evidence | QA Lead; mid-phase review after Task 3 | Both intervention and non-intervention states with stored flag/note values. |
| VAL-004 | QA-05; Review 1 confirms task evidence | QA Lead; mid-phase review after Task 3 | Epic rejection non-zero with clear error and unchanged journal/rows; task/subtask success. |
| VAL-005 | QA-06; Review 1 confirms task evidence | QA Lead; mid-phase review after Task 3 | Fixture-seeded pass+fail snapshot, source `gate_runs` preserved, no-gate `[]` and 0/0. |
| VAL-006 | QA-07; Review 1 confirms task evidence | QA Lead; mid-phase review after Task 3 | Real command reject -> complete -> pass transcript/test output; rejected transition names issue and appends nothing. |
| VAL-007 | QA-08; Review 1 confirms task evidence | QA Lead; mid-phase review after Task 3 | Before/after structural parity after deleting SQLite DB and running `fid rebuild`. |
| VAL-008 | HUMAN-01a..HUMAN-01e | Human phase-end dogfood | One human checklist row per observable. |
| VAL-009 | QA-09; Review 2 confirms task evidence | QA Lead; mid-phase review on Task 4 | `wc -l`, diff/API-change review, test-count/no-weaken proof, full gate output. |

## 5. Task Command Gates

Each task completion report must record every required command as `PASS`, or
`BLOCKED` with reason. Missing, silently skipped, or narrower substituted
commands fail QA review.

| Task | Required commands | Conditional commands | Coverage expectation |
|---|---|---|---|
| Task 1 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go build -o /tmp/fid ./cmd/fid` to produce the binary for pasted CLI transcripts. | `cmd/fid >= 70%`; `internal/issues >= 80%`. |
| Task 2 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go test -tags integration -race ./...` if any snapshot test is behind `integration`; `go build -o /tmp/fid ./cmd/fid` if CLI transcripts are pasted. | `internal/issues >= 80%`. |
| Task 3 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go test -tags integration -race ./...` if the e2e test is behind `integration`; `go build -o /tmp/fid ./cmd/fid` for real CLI transcripts. | `cmd/fid >= 70%`; `internal/issues >= 80%`; `internal/workflow >= 80%`. |
| Task 4 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go build -o /tmp/fid ./cmd/fid` to confirm the binary still builds after the split. | `internal/issues >= 80%`. |

Phase-end QA also runs the full phase gate:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

## 6. Command And Feasibility Validation

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | `gofmt -h` plus `cmd/` and `internal/` existence check | PASS | Help lists `-l`; directories exist. |
| `golangci-lint run ./...` | `golangci-lint run --help` | PASS | Help shows `run [flags]` and package-wide lint command. |
| `go test -race -coverprofile=coverage.out ./...` | `go help test` and `go help testflag` | PASS | Help documents package list mode, `-race`, and `-coverprofile`. |
| `go tool cover -func=coverage.out` | `go tool cover -help` | PASS | Help documents `-func string` for coverage profile output. |
| `go test -tags integration -race ./...` | `go help build` and existing `cmd/fid/e2e_v01c_ac_test.go` pattern | PASS | Help documents `-tags`; existing integration e2e test exists. |
| `go test -bench . -benchmem ./...` | `go help testflag` | PASS | Help documents `-bench` and `-benchmem`. |
| `go build -o /tmp/fid ./cmd/fid` | `go help build` and `/tmp`, `cmd/` existence check | PASS | Help documents `-o`; `/tmp` and `cmd/` exist. |
| Completion reports | Template/path existence check | PASS | `completion-report-template.md` and phase task directory exist. |
| Mid-phase reviews | Template/path existence check | PASS | `mid-phase-review-template.md` and phase QA directory exist. |
| Final review and human checklist | Template/path existence check | PASS | `final-review-template.md` and `human-checklist-template.md` exist. |
| Gate snapshot fixture assumption | Contract/task/spec inspection | PASS | DEC-C and Task 2 explicitly allow fixture-seeded `gate_runs`; `go test` supports `testdata`. |
| Rebuild parity assumption | Existing command/test pattern check | PASS | `cmd/fid/rebuild_test.go` and `cmd/fid/e2e_v01c_ac_test.go` exist; Task 3 requires `fid rebuild`. |
| VAL-009 starting file-size state | `wc -l internal/issues/issue.go` | PASS | Current `issue.go` is 677 lines, matching CF-012 trigger. |
| Human dogfood checklist feasibility | Phase brief and contract inspection | PASS | VAL-008 explicitly assigns human phase-end dogfood with five observables. |

## 7. Deterministic Checklist

| ID | Check | Covers | Evidence required |
|---|---|---|---|
| QA-01 | Compare Task 1..4 completion reports to task specs and this QA plan, including required and triggered conditional commands `(qa-lead)` | VAL-001..VAL-007, VAL-009 | Report excerpts showing each command `PASS` or `BLOCKED`; conditional command rationale. |
| QA-02 | Verify `fid issue complete` creates a `TC-<hash>` completion with role, tokens, commits, files, notes path, and JSON read-back `(qa-lead)` | VAL-001 | CLI transcript or test output plus JSON/read-model excerpt. |
| QA-03 | Verify first and second completions produce attempt 1 then 2, retain both rows, mark prior `superseded`, and leave new attempt `dev_complete` `(qa-lead)` | VAL-002 | DB/read-back output or test output proving both rows and statuses. |
| QA-04 | Verify human-intervention flagged and unflagged states, including note storage and absence of note when omitted `(qa-lead)` | VAL-003 | Test or CLI/DB output for both states. |
| QA-05 | Verify epic completion is rejected before append and task/subtask completion succeeds `(qa-lead)` | VAL-004 | Non-zero error output plus unchanged journal/row counts for epic; success output for task/subtask. |
| QA-06 | Verify gate snapshot pass+fail, source `gate_runs` preservation, matching attempt filtering, and no-gate `[]` 0/0 `(qa-lead)` | VAL-005 | Fixture/test output or DB snapshots covering both seeded and no-gate states. |
| QA-07 | Verify real `has_completion_report` gated transition rejects without completion and succeeds after real `fid issue complete` `(qa-lead)` | VAL-006 | Real command or e2e test output showing reject, unchanged state/journal, complete, and successful transition. |
| QA-08 | Verify deleting SQLite DB and running `fid rebuild` reproduces completion and gate projected state `(qa-lead)` | VAL-007 | Before/after structural comparison output for `task_completions` and `gate_runs`. |
| QA-09 | Verify `internal/issues/issue.go` and new same-package split files are each <=400 non-test lines, no public API/signature change is introduced, no tests are deleted/weakened, and full task gate is green `(qa-lead)` | VAL-009 | `wc -l`; diff/test-count evidence; full gate output. |
| QA-10 | Run phase-end full command gate and record coverage threshold evidence `(qa-lead)` | no VAL | Literal command output for format, lint, coverage, integration, benchmarks, and build. |
| QA-11 | Review docs/spec-sync notes in completion reports for public CLI output, DEC-C fixture scope, Phase 5 fixture evidence replacement, and CF-012 closure input `(qa-lead)` | no VAL | Completion-report excerpts and final-review documentation verdict. |

## 8. Mid-Phase Review Cadence

| Review | Window | Artifact | Required focus |
|---|---|---|---|
| Review 1 | Immediately after Task 3 closes, before Task 4 starts | `project-docs/docs/phases/phase-8/qa/review-1.md` | Validate Task 1..3 evidence for VAL-001..VAL-007, command coverage, conditional commands, validate-before-append, gate fixture disclosure, real `has_completion_report` commands, and rebuild parity. |
| Review 2 | On Task 4 after its completion report and commit, before phase-end final review | `project-docs/docs/phases/phase-8/qa/review-2.md` | Validate VAL-009 evidence, file-size proof, behavior-preserving/no-public-API-change proof, no deleted/weakened tests, full command gate, and conditional build. |

Both reviews use `project-docs/templates/mid-phase-review-template.md` and record
`VERDICT: CONTINUE`, `VERDICT: RESUME`, or `VERDICT: BLOCK`.

## 9. Human Verification Scope

These rows are the `VAL-008` human-only fan-out and must be copied into the
phase-end human checklist.

| ID | Check | Covers | Human evidence expected |
|---|---|---|---|
| HUMAN-01a | Record a completion attempt in a real `.fid/` workspace and confirm the fields read as expected: role, tokens, commits, files, notes path `(human)` | VAL-008 | Human observation or transcript note. |
| HUMAN-01b | Record a second attempt and confirm attempt numbering plus supersession are understandable and expected `(human)` | VAL-008 | Human observation or transcript note. |
| HUMAN-01c | Inspect a completion with a gate snapshot and confirm the rolled-up gate result reads as expected `(human)` | VAL-008 | Human observation or transcript note. |
| HUMAN-01d | Record a completion with human intervention and confirm the flag/note presentation is clear and expected `(human)` | VAL-008 | Human observation or transcript note. |
| HUMAN-01e | Satisfy a real `has_completion_report`-gated transition and confirm the reject/pass lifecycle behaves as expected `(human)` | VAL-008 | Human observation or transcript note. |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Task 1..4 completion reports exist and record all required/conditional command
  outcomes.
- Review 1 and Review 2 are complete with no unresolved `VERDICT: BLOCK`.
- QA-01..QA-11 have reproducible evidence or a recorded blocker.
- HUMAN-01a..HUMAN-01e are present in `human-checklist.md`.
- Every `VAL-001..VAL-009` has a final outcome slot in `final-review.md`.
