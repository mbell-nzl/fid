# Phase 6 QA Plan - Ready Queue + v0.1b Closure

> Phase: 6 - Ready Queue + v0.1b Closure
> Refreshed: 2026-06-14
> QA Lead: Codex
>
> Source contract: `project-docs/docs/phases/phase-6/qa/validation-contract.md`.

## Scope

Phase 6 proves `fid ready` for v0.1b plus the two Phase 5 carry-forward fixes:
CF-009 update-output readability and CF-010 file-size refactor. DEC-017 and
DEC-018 are locked by human approval on 2026-06-14: readiness is status-category
based only, and empty bare `fid ready` exits 0 with stdout `{"task": null}` plus
a stderr `note:`.

Grounding files read for this plan:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/docs/phases/phase-6/qa/validation-contract.md`
- `project-docs/docs/phases/phase-6/tasks/task-1.md` through `task-5.md`

## Validation Refinement

| Check | Result |
|---|---|
| Every in-scope PRD/TRD/design/carry-forward requirement maps to `VAL-NNN` | PASS: VAL-001..012 cover ready queue, formats, empty contract, CF-009, CF-010, v0.1b smoke, and dogfood. |
| Every `VAL-NNN` has validator, evidence method, and task/window | PASS: table below. |
| Bundled assertions split by observable | PASS: VAL-012 fans out to HUMAN-012a/b/c. VAL-003, VAL-007, VAL-008, and VAL-009 require both states. |
| Conditional/binary assertions name both states | PASS: positive-only evidence cannot pass VAL-003, VAL-005, VAL-007, VAL-008, or VAL-009. |
| Parity assertions require deterministic conformance evidence | PASS: VAL-010 and VAL-011 require rebuild/parity evidence, not code reading. |
| DEC-017/DEC-018 alignment | PASS: VAL-003 and VAL-007 are evaluated against locked status-category readiness and frozen empty stdout/stderr contract. |
| Open escalations baseline | PASS: `project-docs/docs/current/tl_escalations_current.md` says no open escalations. |

## Task Inventory

| Task | Focus | Covered VAL IDs | QA risk | Notes |
|---|---|---|---|---|
| Task 1 | CF-009 update confirmation one field per line | VAL-009 | Low | Presentation-only, but stderr line counting must be exact. |
| Task 2 | CF-010 split oversized files | VAL-010 | Medium | Refactor-only; requires line-count and behavior/parity proof. |
| Task 3 | Internal ready selection engine | VAL-002..006; supports VAL-001 | High | Core scheduling semantics and DEC-017 live here. |
| Task 4 | `fid ready` CLI, formats, empty state | VAL-001..008 | High | Public CLI contract; DEC-018 and format dispatch are frozen surfaces. |
| Task 5 | v0.1b smoke and held-out recommendation | VAL-011; supports VAL-012 | High | Cross-command evidence plus required integration gate. |
| Phase end | Human dogfood | VAL-012 | Medium | Human-only expectation checks split by sequential/prioritised/empty. |

## Evidence Map

| Evidence ID | Covers | Validator | Task/window | Evidence command(s) or human check | Pass rule | Fail/block rule |
|---|---|---|---|---|---|---|
| EVID-001 | VAL-001 | QA Lead | Task 4 / final review | Completion report and final review cite `cmd/fid/ready_test.go` output from `go test -race -coverprofile=coverage.out ./...`; QA final also runs `fid ready --format json` or equivalent test fixture transcript. | One selected task appears under exactly `task` with `id`, `title`, `status`, `epic_id`, `epic_title`, `sequence`, `priority`, `est_tokens`; no extra selected tasks. | Missing/renamed key, wrong task, multiple tasks, no runnable evidence, or command omitted. |
| EVID-002 | VAL-002 | QA Lead | Tasks 3-4 / final review | Task 3 unit evidence from `internal/issues/ready_test.go`; Task 4 CLI evidence from `cmd/fid/ready_test.go`; required gate `go test -race -coverprofile=coverage.out ./...`. | Sequential epic returns the lowest-sequence unfinished task when that task is `todo`, and skips every higher-sequence task. | Higher-sequence task selected, current actionable task skipped, no CLI confirmation by Task 4, or missing gate evidence. |
| EVID-003 | VAL-003 | QA Lead | Tasks 3-4 / final review | Unit and CLI fixtures for both DEC-017 states: current sequential task `in_progress` and current sequential task `todo`; required gate plus final CLI transcript. | `in_progress` current task returns no task plus blocker/note naming current task; same task as `todo` returns the task. | Positive-only evidence, wrong selected task, missing blocker/note, or semantics based on claims/links instead of status category. |
| EVID-004 | VAL-004 | QA Lead | Tasks 3-4 / final review | Unit and CLI fixtures with prioritised epic containing actionable, started, and lower-priority tasks; `go test -race -coverprofile=coverage.out ./...`. | Highest-priority actionable `todo` task is selected; started and lower-priority tasks are not chosen ahead. | Started/lower-priority task selected, no mixed-state fixture, or missing CLI confirmation by Task 4. |
| EVID-005 | VAL-005 | QA Lead | Tasks 3-4 / final review | Unit and CLI fixture with unattached `todo` plus unattached non-actionable task; required gate output. | Unattached `todo` is eligible project-wide; unattached started/done/cancelled task is excluded. | Non-actionable unattached task selected, actionable unattached task never eligible, or positive-only fixture. |
| EVID-006 | VAL-006 | QA Lead | Tasks 3-4 / final review | Mixed fixture with sequential, prioritised, and unattached candidates; deterministic winner asserted in unit and CLI tests. | Exactly one project-wide top task is returned; each sequential epic contributes at most its current actionable task; tie-break is documented. | Multiple tasks, non-current sequential contribution, non-deterministic ordering, or undocumented tie-break. |
| EVID-007 | VAL-007 | QA Lead | Task 4 / final review | CLI tests and final transcript for both DEC-018 states: empty and non-empty; `go test -race -coverprofile=coverage.out ./...`; `go build -o /tmp/fid ./cmd/fid`. | Empty exits 0, stdout is exactly `{"task": null}` or structurally identical frozen JSON, stderr has `note:`; non-empty returns task with no stderr note. | Non-zero empty exit, missing null task, stdout shape changed, note printed on stdout, spurious non-empty note, or no build evidence. |
| EVID-008 | VAL-008 | QA Lead | Task 4 / final review | CLI tests for `--format json`, `toon`, `markdown`, `yaml` over selected and empty fixtures; default-without-flag compared to JSON. | All four formats render selected and empty states without error; omitted `--format` equals JSON. | Any format fails/malformed, empty state unsupported in a format, or default is not JSON. |
| EVID-009 | VAL-009 | QA Lead | Task 1 / review-1 | `cmd/fid/issue_update_test.go` output from required gate; completion report includes stderr examples for multi-field and single-field updates plus updated `06_DESIGN_PLAN_CLI.md` example. | Multi-field update prints one changed field per line; single-field update prints exactly one field line; header remains `updated <id>:`. | Comma-joined fields, line count mismatch, missing design-doc update, or update behavior/audit regression. |
| EVID-010 | VAL-010 | QA Lead | Task 2 / review-1 | `wc -l internal/issues/transition.go internal/issues/update.go internal/journal/project.go`; task gate output; rebuild-parity test or CLI transcript. | All three named files are <=400 lines; full gate passes; rebuild parity proves projected behavior unchanged. | Any named file >400 lines, exported/API behavior changes without spec, suite/parity failure, or line counts omitted. |
| EVID-011 | VAL-011 | QA Lead | Task 5 / final review | Integration-tagged v0.1b smoke from `go test -tags integration -race ./...`; task gate output; post-rebuild parity assertion; held-out recommendation file state. | Create -> list -> get -> update -> transition -> ready read-back is consistent; delete DB + `fid rebuild` reproduces identical projected state for touched issues. | Any command read-back mismatch, rebuild parity failure, required integration gate omitted, or recommendation artifact missing. |
| HUMAN-012a | VAL-012 | Human | Phase end | Human checklist: run/dogfood a real sequential epic where the current task is ready or in flight. | Human records that `fid ready` matches expected sequential scheduling and explanation. | Human expectation differs, explanation is unclear/wrong, or check cannot be run. |
| HUMAN-012b | VAL-012 | Human | Phase end | Human checklist: run/dogfood a real prioritised epic with multiple candidate priorities. | Human records that `fid ready` picks the expected highest-priority actionable task. | Human expectation differs or priority ordering is unclear. |
| HUMAN-012c | VAL-012 | Human | Phase end | Human checklist: run/dogfood an empty/all-blocked workspace or epic. | Human records that empty output and `note:` explanation are understandable and expected. | Human finds empty result surprising, ambiguous, or inconsistent with DEC-018. |

## Required Gate Coverage

Every completion report must record each required command as `PASS`, or
`BLOCKED` with reason in a blocker report. Missing, silently skipped, or narrower
substituted commands are blocking QA findings.

| Task | Required command evidence | Additional task evidence |
|---|---|---|
| Task 1 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | Update stderr transcript and design example diff for VAL-009. |
| Task 2 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `wc -l` for the three named files and rebuild-parity proof for VAL-010. |
| Task 3 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | Unit test names/output covering VAL-002..006 and documented tie-break. |
| Task 4 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go build -o /tmp/fid ./cmd/fid` | CLI transcripts or tests for VAL-001..008, including DEC-018 empty/non-empty states. |
| Task 5 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go build -o /tmp/fid ./cmd/fid` | E2E smoke transcript/assertions, rebuild parity, and `held-out-recommendation.md`. |
| Phase end | Recommended phase-level gate from AGENTS.md: task gate plus integration, benchmarks, and build | Final review records outcomes for VAL-001..012 and emits human checklist rows HUMAN-012a/b/c. |

## Feasibility And Scope Checks

| Check | Evidence proving practical now | Blocks phase execution if |
|---|---|---|
| Branch and baseline | `git branch --show-current` returned `phase-6`; `git status --short` was clean before plan authoring. | Wrong branch, dirty unrelated task work, or unresolved merge/conflict state. |
| Phase artifacts exist | `find project-docs/docs/phases/phase-6 -maxdepth 2 -type f` showed validation contract and task-1..task-5 specs. | Missing validation contract, missing task spec, or no QA plan. |
| No open escalations | `project-docs/docs/current/tl_escalations_current.md` states `Open Escalations: None`. | New open ESC appears before/during a task, or DEC-017/DEC-018 is reopened. |
| Commands are concrete | Task specs list concrete Go/golangci/gofmt commands; no placeholders found. `gofmt -h`, `golangci-lint --version`, `go help test`, `go tool cover -h`, and `go help build` validate syntax/tool availability. | Any required command is unavailable, syntactically invalid, still a placeholder, or needs an unrecorded dependency/tool decision. |
| Test and CLI scaffolding exists | `cmd/fid/*_test.go`, `internal/issues/*_test.go`, and `internal/journal/*_test.go` exist; current tests already use temp workspaces, `runFIDCommand`, format assertions, and rebuild parity patterns. | No practical harness for a VAL row, or evidence would require hidden held-out scenarios. |
| CF-010 target is real | `wc -l` baseline: `transition.go` 426, `update.go` 551, `project.go` 524. | Files are not reducible within task scope without behavior/API/schema changes. |
| Task 5 integration scope | `go help test` confirms `-tags` support; repo already has temp-workspace and rebuild helpers. | Integration test cannot be isolated in `t.TempDir()`, requires network, or requires held-out scenario content. |
| Human checks are feasible | VAL-012 uses only a real `.fid/` workspace and public `fid ready` behavior; no hidden evaluator prompt required. | Human cannot run the workspace scenarios, or the check requires non-human deterministic proof only. |

## Mid-Phase Review Cadence

`qa_reviewer` runs twice:

| Review | Window | Output | Required focus |
|---|---|---|---|
| Review 1 | After Tasks 1-3 complete | `project-docs/docs/phases/phase-6/qa/review-1.md` | Validate completion reports, gate coverage, VAL-009, VAL-010, and Task 3 progress toward VAL-002..006. |
| Review 2 | After Tasks 4-5 complete, before QA final review | `project-docs/docs/phases/phase-6/qa/review-2.md` | Validate public `fid ready` contract, DEC-017/DEC-018 evidence, VAL-001..008, VAL-011, required Task 4 build, and required Task 5 integration gate. |

Each review compares completion-report command results with the exact task-spec
commands and this plan. Missing or narrower command evidence is Blocking.

## Phase-End Rule

QA final review may declare PASS only when every deterministic row above has
reproducible evidence, every required gate is recorded, VAL-001..012 each has a
final `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD` outcome, and the generated
human checklist contains only HUMAN-012a, HUMAN-012b, and HUMAN-012c plus any
new human-only judgment items discovered during execution.
