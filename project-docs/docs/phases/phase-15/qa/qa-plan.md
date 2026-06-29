# Phase 15 QA Plan - `fid dev *` Composites + Cleanup

> Phase: 15 - `fid dev *` composites + CF-013/CF-014/CF-015 cleanup
> Refreshed: 2026-06-29
> QA Lead: Codex QA Lead
>
> Location: `project-docs/docs/phases/phase-15/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 15 ships the `fid dev start/gate/finish/block` composite commands and the
scheduled cleanup items on the same paths. QA verifies only the phase contract:
one-event composite writes, validate-before-append rejection paths,
permissions-matrix parity, behavior-preserving cleanup, deterministic rebuild,
and phase-end human dogfood. No new product scope is added here.

Sources used:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/docs/phases/phase-15/qa/validation-contract.md`
- `project-docs/docs/phases/phase-15/tasks/task-1.md`
- `project-docs/docs/phases/phase-15/tasks/task-2.md`
- `project-docs/docs/phases/phase-15/tasks/task-3.md`
- `project-docs/docs/phases/phase-15/tasks/task-4.md`

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: YES.
- Every `VAL-NNN` has one validator role: YES.
- Every `VAL-NNN` has a concrete evidence method: YES.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: YES. `VAL-012` fans out into `HUMAN-15a`..`HUMAN-15d`.
- Conditional or binary assertions require both qualifying and non-qualifying
  states: YES. `VAL-004`, `VAL-005`, `VAL-006`, and `VAL-008` require success
  and reject-path evidence.
- Parity or "matches reference" assertions have deterministic conformance proof
  when available: YES. `VAL-006` requires section-13 matrix/show/enforce
  parity; `VAL-011` requires rebuild determinism.
- Revisions requested from Tech Lead: none.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | CF-014 file split and CF-015 label helper extraction | VAL-001, VAL-002 | Medium | Refactor-only; evidence must prove size targets and unchanged behavior. |
| Task 2 | `internal/dev`, matrix rows, `dev start`, `dev gate` | VAL-003, VAL-004, VAL-005, VAL-006 partial | High | Public CLI additions already confirmed; require one-event and zero-append reject evidence. |
| Task 3 | CF-013 builder, `dev finish`, `dev block` | VAL-006 remainder, VAL-007, VAL-008, VAL-009, VAL-010 | High | Heaviest behavior window; finish/block denial evidence completes VAL-006. |
| Task 4 | Integration smoke and dogfood prep | VAL-011, VAL-012 prep | Medium | Adds integration-tagged e2e smoke; human run remains phase-end. |
| Phase end | Final deterministic review and human checklist | VAL-001..VAL-012 | High | QA Lead records final outcome for every VAL; human executes only `VAL-012` checks. |

## 4. Validation Coverage Mapping

| VAL | Required proof | Window | Evidence owner |
|---|---|---|---|
| VAL-001 | `cmd/fid/issue.go` and `internal/journal/replay.go` plus resulting files are <=400 lines; pre-existing command/journal behavior gates pass unchanged. | Task 1, review 1, final | Developer, QA Reviewer, QA Lead |
| VAL-002 | Target label write functions are <=60 lines; label add/remove tests and behavior pass unchanged. | Task 1, review 1, final | Developer, QA Reviewer, QA Lead |
| VAL-003 | `dev start` emits one `dev.started` event with claim, `todo -> in_progress`, and open run. | Task 2, review 1, final | Developer, QA Reviewer, QA Lead |
| VAL-004 | `dev start` proves both valid unclaimed `todo` success and zero-append rejection for conflicting claim and invalid edge. | Task 2, review 1, final | Developer, QA Reviewer, QA Lead |
| VAL-005 | `dev gate` records reported result on active attempt and proves zero-append rejection without started state plus live claim. | Task 2, review 1, final | Developer, QA Reviewer, QA Lead |
| VAL-006 | Four `dev.*` matrix rows match section 13; `policy show` and `Authorize` use same source; developer allowed, `qa_reviewer` denied before append for all composites. | Tasks 2-3, review 1, final | Developer, QA Reviewer, QA Lead |
| VAL-007 | `dev finish` emits one `dev.finished` event with run close, completion snapshot, transition, and claim release. | Task 3, review 1, final | Developer, QA Reviewer, QA Lead |
| VAL-008 | `dev finish` rejects unverified developer ACs before append and succeeds after developer AC verification with no pre-existing completion. | Task 3, review 1, final | Developer, QA Reviewer, QA Lead |
| VAL-009 | `dev block` emits one `dev.blocked` event with `failed_gate` completion, reason, transition, and claim release. | Task 3, review 1, final | Developer, QA Reviewer, QA Lead |
| VAL-010 | Shared completion-effect builder is reused by `internal/dev`; `fid issue complete` behavior unchanged; target functions <=60 lines. | Task 3, review 1, final | Developer, QA Reviewer, QA Lead |
| VAL-011 | Integration smoke drives finish and block paths, denial, one-event effects, and rebuild determinism. | Task 4, final | Developer, QA Lead |
| VAL-012 | Human can identify start/gate/finish/block outcomes from CLI output without reading implementation. | Phase end | Human |

## 5. Command Evidence Requirements

Every task completion report must record each required command below as `PASS`
with output excerpt, or `BLOCKED` with reason and escalation. Missing, silently
skipped, or narrower substituted commands are failed gates.

| Window | Required commands | Required evidence |
|---|---|---|
| Task 1 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go build -o /tmp/fid ./cmd/fid` | Full command output/status; coverage percentages for `internal/labels`, `internal/issues`, `internal/journal`, and `cmd/fid`; `wc -l` proof for split files. |
| Task 2 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go build -o /tmp/fid ./cmd/fid` | Full command output/status; coverage percentages for `internal/dev`, `internal/policy`, and `cmd/fid`; CLI/domain evidence for VAL-003..006. |
| Task 3 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go build -o /tmp/fid ./cmd/fid` | Full command output/status; coverage percentages for `internal/dev`, `internal/issues`, and `cmd/fid`; CLI/domain evidence for VAL-006..010. |
| Task 4 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go build -o /tmp/fid ./cmd/fid` | Full command output/status; `cmd/fid` coverage; integration smoke output proving both composite paths, denial, and rebuild determinism. |
| Task 4 conditional | `go test -bench . -benchmem ./...` | Required only if a benchmark is added or touched; otherwise completion report states not applicable. |
| Phase end | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go test -bench . -benchmem ./...`; `go build -o /tmp/fid ./cmd/fid` | Final-review gate output; package coverage; deterministic checklist evidence; final outcome for every VAL. |

Integration gate discipline: because Tasks 1-3 touch `cmd/fid`, and Task 4 adds
an integration-tagged test, `go test -tags integration -race ./...` is required
for every task in this phase. It may be recorded `BLOCKED` only with a concrete
reason and escalation.

## 6. Command And Feasibility Validation

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Phase artifacts and template exist | File-state check | PASS | `ls` found task specs 1-4, `validation-contract.md`, and `qa-plan-template.md`. |
| `gofmt -l cmd internal` | Local help output | PASS | `gofmt -h` documents `-l` and path arguments. |
| `golangci-lint run ./...` | Tool version check | PASS | `golangci-lint --version` reports `2.12.2`, matching the v2 quality-plan target. |
| `go test -race -coverprofile=coverage.out ./...` | Local help output | PASS | `go help test` documents package testing; `go help build` documents `-race`; `go test` help documents `-coverprofile`. |
| `go tool cover -func=coverage.out` | Local usage output | PASS | `go tool cover -h` prints `-func string` coverage summary support. Usage exits non-zero, but syntax is documented. |
| `go test -tags integration -race ./...` | Local help plus file-state check | PASS | `go help build` documents `-tags`; `rg -n "//go:build integration" cmd internal` found existing integration tests in `cmd/fid` and `internal/issues`. |
| `go test -bench . -benchmem ./...` | Local help output | PASS | `go help test` supports test binary flags; `08_QUALITY_PLAN.md` names this exact benchmark command for phase-end. |
| `go build -o /tmp/fid ./cmd/fid` | Local help output | PASS | `go help build` documents `-o output` and package arguments; `/tmp` is writable. |
| File-size proof for VAL-001 | File-state baseline | PASS | `wc -l` baseline: `cmd/fid/issue.go` 985, `internal/journal/replay.go` 485. |
| Function-size proof for VAL-002/010 | File-state baseline and final line-span inspection | PASS | Baseline files exist: `internal/labels/add.go`, `internal/issues/label.go`, `internal/issues/complete.go`; final review must cite function spans or line-count report. |
| Human dogfood checklist feasibility | Contract fan-out | PASS | `VAL-012` has four CLI-visible observables: start, gate, finish, block. Human checklist can be generated from these without implementation access. |

## 7. Deterministic Checklist

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify task completion reports ran every required command exactly, including integration and build where required `(qa-lead)` | VAL-001..VAL-011 | Completion-report command table plus output excerpts. |
| QA-02 | Verify file-size and function-size targets after cleanup and CF-013 extraction `(qa-lead)` | VAL-001, VAL-002, VAL-010 | `wc -l` output for split files and reproducible function line-span evidence for target functions. |
| QA-03 | Verify behavior-preserving cleanup claims are backed by unchanged or mechanically moved tests, not weakened tests `(qa-lead)` | VAL-001, VAL-002, VAL-010 | `git diff`/review evidence plus passing relevant package tests from required gates. |
| QA-04 | Verify `dev start` and `dev gate` one-event success paths and zero-append reject paths `(qa-lead)` | VAL-003, VAL-004, VAL-005 | Test/CLI logs showing journal-line deltas and projected DB state. |
| QA-05 | Verify `dev.*` policy rows match section 13 and show/enforce parity holds for allowed and denied roles `(qa-lead)` | VAL-006 | Matrix test output plus CLI evidence for `developer` allowed and `qa_reviewer` denied before append. |
| QA-06 | Verify `dev finish` success path, unmet-AC rejection, and no-pre-existing-completion success after AC verification `(qa-lead)` | VAL-007, VAL-008 | Test/CLI logs showing one event, run close, completion, transition, release, and zero-append rejection. |
| QA-07 | Verify `dev block` success path records `failed_gate` completion, reason, transition, and claim release `(qa-lead)` | VAL-009 | Test/CLI logs showing one event and projected completion/state/claim values. |
| QA-08 | Verify integration smoke covers finish path, block path, permission denial, and delete-and-rebuild determinism `(qa-lead)` | VAL-011 | `go test -tags integration -race ./...` output and e2e assertion summary or test names. |
| QA-09 | Verify docs changed only where task specs require and no product scope was added `(qa-lead)` | no VAL | Diff review of `02_TRD.md`, `06_DESIGN_PLAN_CLI.md`, and optional `07_DESIGN_PLAN_FLOWS.md`. |

## 8. Mid-Phase Review Cadence

- Cadence: QA review after every second completed task, which is stricter than
  the required every-third-task cadence.
- Required review windows:
  - `review-1.md` after Tasks 1-2 complete: verify VAL-001..006 evidence,
    required command coverage, and integration gate discipline.
  - `review-2.md` after Task 4 complete or fold into final review if final QA
    starts immediately: verify VAL-007..011 evidence and Task 4 gate output.
- Review template: `project-docs/templates/mid-phase-review-template.md`.
- Each review must record `CONTINUE` or `BLOCK`, cite covered `VAL-NNN` IDs, and
  call out missing/narrower/substituted required commands as Blocking.

## 9. Human Verification Scope

Human checks are generated at phase end into
`project-docs/docs/phases/phase-15/qa/human-checklist.md`. They must contain only
human dogfood observations and use one observable per row.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-15a | Confirm CLI output makes it clear that `fid dev start` claims the task, moves it to `in_progress`, and opens a run `(human)` | VAL-012 | Human note from a real `.fid` workspace run. |
| HUMAN-15b | Confirm CLI output makes it clear that `fid dev gate` records a reported gate result and does not appear to execute the command `(human)` | VAL-012 | Human note from a real `.fid` workspace run. |
| HUMAN-15c | Confirm CLI output makes it clear that `fid dev finish` completes the task, closes the run, and releases the claim `(human)` | VAL-012 | Human note from a real `.fid` workspace run. |
| HUMAN-15d | Confirm CLI output makes it clear that `fid dev block` records a failed gate reason and releases the claim `(human)` | VAL-012 | Human note from a real `.fid` workspace run. |

## 10. Phase-End Decision Rule

The phase can proceed to final PASS only when:

- Every required task command is recorded `PASS`, or `BLOCKED` with an accepted
  escalation.
- The phase-level command gate passes.
- Every `VAL-001`..`VAL-011` has reproducible deterministic evidence.
- `VAL-012` is transferred to the human checklist and completed by the human, or
  explicitly marked `BLOCKED`/`CARRY-FORWARD` with reason.
- Integration gate discipline was followed for every task touching `cmd/fid` or
  adding/moving integration-tagged tests.
- No unresolved Blocking finding, scope addition, or required documentation gap
  remains.
