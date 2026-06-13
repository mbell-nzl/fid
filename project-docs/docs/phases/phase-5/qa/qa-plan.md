# Phase 5 QA Plan - Mutation + Workflow Transitions

> Phase: 5 - Mutation + workflow transitions (v0.1b)
> Refreshed: 2026-06-13
> QA Lead: Codex
>
> Location:
> `project-docs/docs/phases/phase-5/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 5 proves that fid can mutate issue fields, record audited activity,
advance issues through hard-coded default workflows, reject invalid writes
without poisoning the canonical journal, and produce actionable transition
errors. The explicit deferrals remain out of scope: `fid activity list`,
`--epic`/`--inbox` reparent/promote, claim-gated transitions, AC/completion/
concern write commands, configurable workflows, and escalation workflows.

Sources used:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/01_PRD.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/docs/phases/phase-5/qa/validation-contract.md`
- `project-docs/docs/phases/phase-5/tasks/task-1.md`
- `project-docs/docs/phases/phase-5/tasks/task-2.md`
- `project-docs/docs/phases/phase-5/tasks/task-3.md`
- `project-docs/docs/phases/phase-5/tasks/task-4.md`
- `project-docs/docs/phases/phase-5/tasks/task-5.md`

Requirement coverage rollup: all in-scope PRD/TRD/design requirements are
referenced by at least one `VAL-NNN`: YES.

| Requirement | Phase 5 coverage |
|---|---|
| PRD-ISSUE-015 | `fid issue update` field mutation and invalid-update preservation: VAL-001. |
| PRD-ISSUE-021 | Owner role/actor update audit: VAL-003. |
| PRD-WF-001/005/006 | Default workflow allowed edges and status categories: VAL-006. |
| PRD-WF-002 | Invalid edge rejection naming current and requested states: VAL-007. |
| PRD-WF-003/004 | Rule evaluation before transition and rich unmet-rule errors: VAL-008, VAL-009, VAL-010. |
| PRD-WF-007/008 | `all_acs_verified` and `has_completion_report`: VAL-008, VAL-009. |
| PRD-WF-009 | `concerns_addressed`: VAL-012. |
| PRD-WF-010 | `all_child_tasks_done_or_cancelled`: VAL-010. |
| PRD-WF-013 | Cancel transition semantics: VAL-011. |
| PRD-OPS-001/002 | Journal append as commit point; one clean logical write: VAL-002, VAL-005, VAL-007. |
| PRD-OUT-017 | Write-side `activity_log` projection only; public read command deferred: VAL-004. |
| TRD §4.3 / TRD-ARCH-001/004/005 | Journal effects, validate-before-append, thin CLI adapter, rebuild parity: VAL-001, VAL-002, VAL-005, VAL-006, VAL-007. |
| TRD §4.5 / §4.6 / TRD-ARCH-009 | App-layer transition rules and default workflows: VAL-006..VAL-012. |
| TRD SCHEMA §14.2 / §14.6 | Polymorphic ref backstop and `activity_log` projection discipline: VAL-004. |
| CLI design `fid issue update` | Flags, field set, output/read-back parity: VAL-001, VAL-003. |
| CLI design `fid issue transition` | `--to`, `--by`, status/category read-back, rich errors: VAL-006..VAL-011. |
| CF-005 | `activity` op plus real emitters in update/transition: VAL-003, VAL-004, VAL-006. |
| CF-007 / DEC-014 | Rejected writes append zero lines and do not poison later reads: VAL-005, VAL-007, VAL-013. |

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: YES.
- Every `VAL-NNN` has one validator role: YES.
- Every `VAL-NNN` has a concrete evidence method: YES.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: YES. VAL-013 is human-only and is fanned out in
  section 9.
- Conditional or binary `VAL-NNN` assertions require both qualifying and
  non-qualifying states: YES. VAL-001, VAL-004, VAL-005, VAL-007, VAL-008,
  VAL-009, VAL-010, VAL-011, and VAL-012 all require both states.
- Parity or "matches reference" `VAL-NNN` assertions have deterministic
  conformance proof when available: YES. VAL-002 and VAL-004 require rebuild
  parity; VAL-006 requires structural state/category conformance to TRD §4.6.
- Revisions requested from Tech Lead: none.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | Validate-before-append seam for rejected writes | VAL-005 | High | Foundational CF-007 fix; mid-phase review required before update/transition depend on the seam. |
| Task 2 | `activity` effect op and `activity_log` projection discipline | VAL-004 | Medium | Must prove one row per activity effect, unresolved ref rejection, replay determinism, and no non-projection writer. |
| Task 3 | `fid issue update` with activity audit | VAL-001, VAL-002, VAL-003, advances VAL-005 | High | First real mutable CLI command and first activity emitter; requires read-back, journal-count, event-shape, DB, and rebuild evidence. |
| Task 4 | `internal/workflow` default graphs and pure rule functions | VAL-012, advances VAL-006, VAL-008, VAL-009, VAL-010 | Medium | Structural conformance to TRD §4.6 and both-state rule tests are required before Task 5 wiring. |
| Task 5 | `fid issue transition` with rule enforcement and activity audit | VAL-006, VAL-007, VAL-008, VAL-009, VAL-010, VAL-011, advances VAL-005 | High | Requires both allowed/rejected paths, rule fixtures, journal hygiene, one-event transition shape, and rebuild parity. |
| Phase end | Dogfood and final deterministic review | VAL-001..VAL-013 | High | QA Lead executes deterministic checklist; human executes only VAL-013 fan-out rows. |

## 4. Required Evidence Strategy

### VAL Routing Matrix

| VAL | Acceptance signal | Deterministic `(qa-lead)` evidence | Mid-phase review evidence | Human evidence |
|---|---|---|---|---|
| VAL-001 | Valid update fields read back, untouched fields preserved, invalid updates reject unchanged. | QA-05, QA-06 | Review 2 verifies Task 3 report/tests and command coverage. | None. |
| VAL-002 | Exactly one update event with documented effects; rebuild parity is byte-identical. | QA-07, QA-08 | Review 2 verifies Task 3 journal/rebuild evidence. | None. |
| VAL-003 | Owner role/actor update reads back and activity row records actor/from/to; unknown role rejects. | QA-05, QA-06, QA-09 | Review 2 verifies Task 3 DB/activity evidence. | None. |
| VAL-004 | `activity` op writes one row, unresolved refs reject, rebuild deterministic, no independent writer. | QA-04, QA-19 | Review 2 verifies Task 2; Review 4 verifies transition emitter reuse. | None. |
| VAL-005 | Rejected writes append zero lines and later reads succeed; accepted writes append exactly one line. | QA-03, QA-06, QA-14, QA-19 | Review 1 verifies Task 1; Review 2 invalid update; Review 4 invalid transition. | Covered indirectly by HUMAN-013e/f as dogfood judgment only. |
| VAL-006 | Allowed task/epic/subtask edges advance status and category. | QA-10, QA-13, QA-19 | Review 3 verifies graph readiness; Review 4 verifies CLI transitions. | None. |
| VAL-007 | Disallowed edge rejects naming current and requested states; zero journal append. | QA-14 | Review 4 verifies Task 5 invalid-edge evidence. | None. |
| VAL-008 | `in_progress -> dev_complete` enforces developer AC verification and completion report, both unmet and satisfied states. | QA-12, QA-15 | Review 3 verifies pure rule readiness; Review 4 verifies CLI/fixture enforcement. | None. |
| VAL-009 | `qa_review -> qa_pass` enforces qa_lead verification role parameter, both states. | QA-12, QA-16 | Review 3 verifies parameterized rule readiness; Review 4 verifies CLI/fixture enforcement. | None. |
| VAL-010 | Epic wrap rejects open children with child IDs and succeeds when all children done/cancelled. | QA-12, QA-17 | Review 3 verifies pure rule readiness; Review 4 verifies CLI/fixture enforcement. | None. |
| VAL-011 | Cancel from non-cancelled states succeeds; cancelled -> cancelled rejects. | QA-10, QA-18 | Review 3 verifies graph rule; Review 4 verifies CLI enforcement. | None. |
| VAL-012 | `concerns_addressed` returns satisfied iff no open concerns. | QA-11 | Review 3 verifies Task 4 unit tests. | None. |
| VAL-013 | Human dogfood confirms rich errors are actionable and rejected writes do not disrupt later work. | None; deterministic journal hygiene is covered by VAL-005/007. | None required beyond final QA preparing the human checklist. | HUMAN-013a through HUMAN-013f. |

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Completion report for Task 1 records the exact four required commands as `PASS` or `BLOCKED`, plus tests proving rejected create line count unchanged, later `issue get`/rebuild succeeds, accepted create appends one line, and dry-run projection leaves no committed DB/projection_meta mutation. | VAL-005 | Developer, QA Reviewer | Task 1, review-1 |
| EVID-02 | Task 2 package tests or DB evidence showing one `activity_log` row per `activity` effect, unresolved `(ref_type, ref_id)` rejected, replay/rebuild reproduces rows, and `activity_log` writes occur only through the projection path. | VAL-004 | Developer, QA Reviewer | Task 2, review-2 |
| EVID-03 | Task 3 CLI/package evidence showing each supported update field reads back with untouched fields preserved, invalid update returns non-zero and leaves read-back/journal line count unchanged, and `--by` validates known roles. | VAL-001, VAL-005 | Developer, QA Reviewer, QA Lead | Task 3, review-2, final review |
| EVID-04 | Task 3 journal evidence: exactly one event per update; event includes `issues` update effect, label delta effects when labels change, one `activity` effect, updated_at advance, preconditions, and byte-identical post-rebuild issue JSON. | VAL-002 | Developer, QA Reviewer, QA Lead | Task 3, review-2, final review |
| EVID-05 | Task 3 DB/activity evidence proving owner role/actor read-back and exactly one `activity_log` row with actor and before/after details; unknown owner role rejected. | VAL-003 | Developer, QA Reviewer, QA Lead | Task 3, review-2, final review |
| EVID-06 | Task 4 workflow tests proving all task/epic/subtask states, status categories, allowed/disallowed edges, and universal cancel semantics match TRD §4.6. | VAL-006, VAL-011 | Developer, QA Reviewer | Task 4, review-3 |
| EVID-07 | Task 4 rule tests proving developer and qa_lead `all_acs_verified`, `has_completion_report`, `concerns_addressed`, and `all_child_tasks_done_or_cancelled` each return correct structured satisfied/unsatisfied results. | VAL-008, VAL-009, VAL-010, VAL-012 | Developer, QA Reviewer | Task 4, review-3 |
| EVID-08 | Task 5 CLI/package evidence for allowed transitions across task, epic, and subtask workflows; `fid issue get --format json` shows expected status and status_category. | VAL-006 | Developer, QA Reviewer, QA Lead | Task 5, review-4, final review |
| EVID-09 | Task 5 CLI/package evidence for disallowed edge rejection naming current and requested states, with journal line count unchanged and later read succeeding. | VAL-007, VAL-005 | Developer, QA Reviewer, QA Lead | Task 5, review-4, final review |
| EVID-10 | Task 5 fixture-backed evidence for `in_progress -> dev_complete`: unmet developer AC verification and missing completion report are named with AC IDs; satisfied fixture transitions. | VAL-008 | Developer, QA Reviewer, QA Lead | Task 5, review-4, final review |
| EVID-11 | Task 5 fixture-backed evidence for `qa_review -> qa_pass`: unmet qa_lead verification gaps reject; satisfied fixture transitions. | VAL-009 | Developer, QA Reviewer, QA Lead | Task 5, review-4, final review |
| EVID-12 | Task 5 fixture-backed evidence for epic `in_progress -> wrapped`: open child IDs reject; all children done/cancelled succeeds. | VAL-010 | Developer, QA Reviewer, QA Lead | Task 5, review-4, final review |
| EVID-13 | Task 5 evidence that at least two distinct non-cancelled states can transition to `cancelled`, and cancelled -> cancelled is rejected. | VAL-011 | Developer, QA Reviewer, QA Lead | Task 5, review-4, final review |
| EVID-14 | Phase-end deterministic gate output: format, lint, race coverage, coverage report, integration command, benchmarks, build, PRD compliance, security review, performance spot-check, documentation review, and deterministic checklist evidence excerpts. | VAL-001..VAL-012, no VAL quality controls | QA Lead | Final review |
| EVID-15 | Human dogfood against a real `.fid/` workspace, split into per-observable rows for update usability, transition usability, invalid transition error actionability, rejected update error actionability, journal cleanliness after rejected writes, and later-read recovery. | VAL-013 | Human | Human checklist |

### Expected Tests By Task

| Task | Expected test coverage |
|---|---|
| Task 1 | Red/green regression for rejected create journal line count unchanged; later `Get` and rebuild recovery; accepted create one-line append; dry-run helper rollback on success and failure; `projection_meta` unchanged. |
| Task 2 | `activity` effect inserts one row; unresolved polymorphic ref rejected; activity rows replay deterministically; `details` serializes as JSON text; multiple activity effects produce one row each. |
| Task 3 | `internal/issues.Update` and CLI tests for every supported field; labels replace-all delta; invalid inputs unchanged and zero-append; owner change activity row; one-event update; rebuild parity; `--by` role validation. |
| Task 4 | Workflow graph table tests for allowed/disallowed edges, categories, and cancel; rule function tests for both satisfied and unsatisfied states, including developer vs qa_lead role parameters and ID-carrying unmet details. |
| Task 5 | Transition service and CLI tests for allowed edges, disallowed edge errors, developer and qa_lead verification rules, completion report rule, epic child rule, cancel semantics, zero-append rejected transitions, one-event successful transition, preconditions, activity effect, and rebuild parity. |

## 5. Deterministic Gates

If a required command is unavailable, still a placeholder, or infeasible when
execution is requested, the QA plan is BLOCKED. No Phase 5 task command is a
placeholder.

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task and phase end |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Every task and phase end |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | Every task and phase end |
| Coverage report | `go tool cover -func=coverage.out` | PASS; final review cites package/file evidence for COV-001/COV-002 scope touched by the phase | Every task and phase end |
| Integration | `go test -tags integration -race ./...` | PASS or BLOCKED with reason | Phase end |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | Phase end |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Phase end; conditional task evidence for Tasks 3 and 5 if manual CLI smoke is used |
| Full phase check | `gofmt -l cmd internal` + `golangci-lint run ./...` + `go test -race -coverprofile=coverage.out ./...` + `go tool cover -func=coverage.out` + `go test -tags integration -race ./...` + `go test -bench . -benchmem ./...` + `go build -o /tmp/fid ./cmd/fid` | PASS or explicit BLOCKED entries in final review | Phase-end gate |

Each task completion report must record the task-spec exact commands below as
`PASS`, or `BLOCKED` with a reason in a blocker report. Missing, silently
skipped, or narrower substituted commands are Blocking findings.

| Task | Required exact commands | Conditional commands |
|---|---|---|
| Task 1 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | None |
| Task 2 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | None |
| Task 3 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go build -o /tmp/fid ./cmd/fid` if manual CLI smoke of `fid issue update` is used as evidence for VAL-001/003 |
| Task 4 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | None |
| Task 5 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out` | `go build -o /tmp/fid ./cmd/fid` if manual CLI smoke of `fid issue transition` is used as evidence for VAL-006/007/008/010/011 |

## 6. Command And Feasibility Validation

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Help output plus path existence check | PASS | `gofmt -h` shows `-l` lists files whose formatting differs; `ls` confirmed `cmd` and `internal` exist. |
| `golangci-lint run ./...` | Help output plus config existence check | PASS | `golangci-lint --help` lists `run`; `golangci-lint run --help` documents the command; `ls` confirmed `.golangci.yml`. |
| `go test -race -coverprofile=coverage.out ./...` | `go help test` and quality-plan command definition | PASS | `go help test` documents package-list mode and test flags; `08_QUALITY_PLAN.md` §5.1 defines this task-level gate. |
| `go tool cover -func=coverage.out` | Cover tool help and coverage command dependency | PASS | `go tool cover -h` documents `-func string` for profile percentages; command depends on prior `coverage.out` from `go test -coverprofile`. |
| `go test -tags integration -race ./...` | `go help build` tag support plus integration scan | PASS | `go help build` documents `-tags`; `rg -n "^//go:build integration" -g "*_test.go"` found no current integration-tagged tests, so the command is syntactically feasible and currently runs the normal package set under the tag. |
| `go test -bench . -benchmem ./...` | `go help test` benchmark support plus benchmark scan | PASS | `go help test` documents benchmark output; `rg -n "func Benchmark" -g "*_test.go"` found benchmark functions in `internal/journal/rebuild_bench_test.go`. |
| `go build -o /tmp/fid ./cmd/fid` | `go help build` plus writable output path and main package path | PASS | `go help build` documents `-o`; sandbox allows `/tmp`; `ls` confirmed `cmd/fid` exists. |
| Task 1 scope paths | File existence check | PASS | `find` confirmed `internal/journal/project.go`, `internal/journal/writer.go`, `internal/issues/issue.go`, `internal/issues/issue_test.go`, and Task 1 spec exist. |
| Task 2 scope paths | File existence check | PASS | `find` confirmed `internal/journal/project.go`, `internal/journal/event.go`, `internal/db/ref_triggers_test.go`, specs, and Task 2 spec exist. |
| Task 3 scope paths and CLI design | File existence and spec search | PASS | `find` confirmed `cmd/fid/issue.go`, `cmd/fid/issue_test.go`, `internal/issues/issue.go`, `internal/issues/detail.go`; `rg` confirmed `06_DESIGN_PLAN_CLI.md` has `fid issue update`. |
| Task 4 scope paths and workflow references | File/spec search | PASS | `find` confirmed `internal/issues/issue.go`; `rg` confirmed TRD/PRD references for `task-default`, `epic-default`, `subtask-default`, and rule names. `internal/workflow` is intentionally created by Task 4. |
| Task 5 scope paths and CLI design | File/spec search | PASS | `find` confirmed current issue/journal/CLI files; `rg` confirmed `06_DESIGN_PLAN_CLI.md` has `fid issue transition`. `internal/workflow` is a Task 4 output and review-3 must verify it before Task 5 begins. |
| Journal fixture assumption for Task 5 rule evidence | Quality-plan fixture convention plus task scope check | PASS | `08_QUALITY_PLAN.md` §3.4 allows committed `testdata/` fixtures; Task 5 also permits generated journal fixtures in tests because v0.1c write commands are deferred. Review-4 must verify fixture paths/helpers cited in the completion report. |
| Completion report paths | Workflow and template check | PASS | `09_AI_WORKFLOW.md` requires `project-docs/docs/phases/phase-5/tasks/task-<n>-done.md`; `completion-report-template.md` exists. |
| Mid-phase review artifact paths | Template check | PASS | `mid-phase-review-template.md` exists; Phase Runner can enforce `review-1.md` through `review-4.md` at the checkpoints in section 8. |
| Final review and human checklist artifact paths | Template check | PASS | `final-review-template.md` and `human-checklist-template.md` exist; QA Lead writes final review and extracts only section 9 human rows. |
| Human dogfood with real `.fid/` workspace | Feasibility from product surface and phase scope | PASS | Phase scope delivers `fid issue update` and `fid issue transition`; human checklist uses a real initialized `.fid/` workspace after Task 5 and avoids held-out scenarios. |

## 7. Deterministic Checklist

Items here are executed or verified by QA Lead at phase-end unless assigned to
a mid-phase review. Every item must end with `(qa-lead)`. PASS without cited
reproducible evidence is invalid.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify each Task 1..5 completion report exists, has the required headings, cites covered VAL IDs, records commit hash/subject, includes dependency-budget note, and records every exact task command as `PASS` or `BLOCKED` with reason `(qa-lead)` | no VAL | File-state evidence from `task-1-done.md`..`task-5-done.md`; `git rev-parse --verify <hash>^{commit}` evidence recorded by Phase Runner or final review. |
| QA-02 | Execute or verify the full phase gate: format, lint, race coverage, coverage report, integration, benchmarks, and build `(qa-lead)` | no VAL | Literal output excerpts for all commands in section 5. |
| QA-03 | Verify rejected create uses validate-before-append: journal line count unchanged, later unrelated read succeeds, rebuild succeeds, accepted create appends exactly one line, and dry-run projection does not advance projection state `(qa-lead)` | VAL-005 | Test output or CLI/file-state excerpts from Task 1 plus final-review seeded fixture rerun or cited reproducible log. |
| QA-04 | Verify `activity` effects project exactly one `activity_log` row per effect, unresolved refs reject, replay/rebuild reproduces activity rows, and only the projection path writes `activity_log` `(qa-lead)` | VAL-004 | Package test output and DB/query or fixture evidence from Task 2. |
| QA-05 | Verify `fid issue update` changes title, priority, labels, spec_md, est_tokens, due_date, owner_role, and owner_actor, while untouched fields are preserved in `fid issue get --format json` `(qa-lead)` | VAL-001, VAL-003 | CLI JSON excerpts, package test output, or final seeded CLI run. |
| QA-06 | Verify invalid update cases for priority out of 0-5, blank title, unknown label, unknown owner role, malformed due_date, and missing/unknown `--by` return non-zero actionable errors, append zero journal lines, and leave read-back unchanged `(qa-lead)` | VAL-001, VAL-003, VAL-005 | CLI/test output with before/after read-back and line-count evidence. |
| QA-07 | Verify one successful update appends exactly one journal event with `issues` update effect, label delta effects when labels change, one `activity` effect, updated_at advancement, and preconditions on changed prior values `(qa-lead)` | VAL-002 | Journal JSONL excerpt or structural test output; read-back evidence. |
| QA-08 | Verify deleting/rebuilding the SQLite DB after update reproduces byte-identical `fid issue get --format json` output from before rebuild `(qa-lead)` | VAL-002 | Pre/post rebuild JSON excerpts or deterministic test log. |
| QA-09 | Verify owner role/actor update creates exactly one `activity_log` row with actor and before/after details, and unknown owner role is rejected `(qa-lead)` | VAL-003 | DB/query evidence or package test output plus CLI read-back. |
| QA-10 | Verify workflow definitions structurally conform to TRD §4.6 for task, epic, and subtask states, categories, allowed edges, disallowed edges, and universal cancel edge semantics `(qa-lead)` | VAL-006, VAL-011 | Task 4 test output and/or structural test table excerpts. |
| QA-11 | Verify `concerns_addressed` returns unsatisfied with open concerns and satisfied with no open concerns `(qa-lead)` | VAL-012 | Task 4 unit test output for both states. |
| QA-12 | Verify `all_acs_verified`, `has_completion_report`, and `all_child_tasks_done_or_cancelled` rule functions return structured unmet details naming AC IDs or child IDs as applicable `(qa-lead)` | VAL-008, VAL-009, VAL-010 | Task 4/5 test output showing both states and IDs. |
| QA-13 | Verify allowed `fid issue transition` edges for task, epic, and subtask advance status and read back the expected status_category `(qa-lead)` | VAL-006 | CLI JSON output or package test output for all three kinds. |
| QA-14 | Verify disallowed `fid issue transition` edge rejects with current and requested state names, appends zero journal lines, and later read succeeds `(qa-lead)` | VAL-007, VAL-005 | CLI/test stderr, line-count, and read-back evidence. |
| QA-15 | Verify task `in_progress -> dev_complete` rejects when developer AC verification or completion report is missing, names unmet conditions and AC IDs, and succeeds when both are seeded `(qa-lead)` | VAL-008 | Fixture-backed test/CLI output for unmet and satisfied states. |
| QA-16 | Verify task `qa_review -> qa_pass` uses the qa_lead verifier role set: unmet qa_lead verification rejects and satisfied fixture transitions `(qa-lead)` | VAL-009 | Fixture-backed test/CLI output for both states. |
| QA-17 | Verify epic `in_progress -> wrapped` rejects with blocking child IDs when a child remains open and succeeds when all children are done/cancelled `(qa-lead)` | VAL-010 | Fixture-backed test/CLI output for both states. |
| QA-18 | Verify cancellation is allowed from at least two distinct non-cancelled states and rejected from already-cancelled `(qa-lead)` | VAL-011 | CLI/test output for both allowed states and rejected state. |
| QA-19 | Verify one successful transition appends exactly one journal event with status update effect, status and workflow_id preconditions, one `activity` effect, and rebuild parity `(qa-lead)` | VAL-006, VAL-004, VAL-005 | Journal excerpt plus pre/post rebuild read-back evidence. |
| QA-20 | Verify documentation touched or intentionally recorded as `verified-no-change` for `06_DESIGN_PLAN_CLI.md`, `07_DESIGN_PLAN_FLOWS.md`, `03_TRD_SCHEMA.md`, and `02_TRD.md` as named by task specs `(qa-lead)` | no VAL | Completion-report doc sections and git/file-state evidence. |
| QA-21 | Perform phase-level PRD compliance, security review for SQL/command injection and unsafe input handling, performance spot-check against TRD-PERF expectations, and documentation review `(qa-lead)` | no VAL | Final-review sections with findings and evidence references. |

## 8. Mid-Phase Review Cadence

Phase Runner must enforce this Phase 5-specific cadence instead of the default
every-third-task cadence. Each review uses
`project-docs/templates/mid-phase-review-template.md`, writes the named review
artifact, verifies in-window `VAL-NNN` evidence, compares completion-report
commands against task specs and this QA plan, and ends with
`VERDICT: CONTINUE` or `VERDICT: BLOCK`.

| Review | Enforced checkpoint | Output path | Window | Required focus |
|---|---|---|---|---|
| Review 1 | After Task 1 completes, before Task 2 starts | `project-docs/docs/phases/phase-5/qa/review-1.md` | Task 1 | VAL-005, CF-007 closure evidence, exact command coverage, dry-run seam test quality. |
| Review 2 | After Task 3 completes, before Task 4 starts | `project-docs/docs/phases/phase-5/qa/review-2.md` | Tasks 2-3 | VAL-001, VAL-002, VAL-003, VAL-004, and VAL-005 invalid-update journal hygiene; activity op + first emitter evidence. |
| Review 3 | After Task 4 completes, before Task 5 starts | `project-docs/docs/phases/phase-5/qa/review-3.md` | Task 4 | VAL-012 final evidence; VAL-006/008/009/010 readiness for Task 5; TRD §4.6 structural conformance. |
| Review 4 | After Task 5 completes, before final QA Lead review | `project-docs/docs/phases/phase-5/qa/review-4.md` | Task 5 | VAL-006..VAL-011 and VAL-005 transition journal hygiene; fixture feasibility and exact command coverage. |

Blocking review findings stop Phase Runner. Soft blockers, new concerns, or
new ESC entries route through TL-Lite per `09_AI_WORKFLOW.md`.

## 9. Human Verification Scope

Items here require human judgment or real dogfooding and must end with
`(human)`. They are the only rows QA Lead should copy into
`human-checklist.md`; deterministic proof belongs in `final-review.md`.
VAL-013 is intentionally fanned into one observable per row.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-013a | In a real `.fid/` workspace, update issue fields used in normal work and judge whether the command flow and confirmation are clear enough to trust `(human)` | VAL-013 | Notes on the update attempted, whether the result was understandable, and any confusing output. |
| HUMAN-013b | In the same workspace, advance at least one allowed workflow edge and judge whether the transition flow communicates the new state clearly `(human)` | VAL-013 | Notes naming the edge tried and whether the status/category outcome was understandable. |
| HUMAN-013c | Attempt an invalid transition and judge whether the rejection names the current/requested states and gives an actionable next step `(human)` | VAL-013 | Error excerpt or paraphrase plus pass/fail judgment on actionability. |
| HUMAN-013d | Attempt a rejected update, such as invalid priority or unknown owner role, and judge whether the rejection explains what to fix `(human)` | VAL-013 | Error excerpt or paraphrase plus pass/fail judgment on actionability. |
| HUMAN-013e | After the rejected update and rejected transition, run a later valid read or unrelated command and judge whether the workspace still behaves normally `(human)` | VAL-013 | Command attempted and observation that later reads/work continues or fails. |
| HUMAN-013f | Inspect dogfood friction from the combined mutate/transition session and decide whether any error wording or workflow behavior should block phase close despite deterministic PASS results `(human)` | VAL-013 | Product-owner judgment: pass/fail/partial with specific wording or workflow concerns. |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Required deterministic gates pass or are recorded BLOCKED with escalation.
- Every `VAL-NNN` has planned evidence and a final-review outcome slot.
- Every `(qa-lead)` item has reproducible evidence.
- All four mid-phase reviews exist or are explicitly BLOCKED by Phase Runner
  with reason.
- Human checks are human-judgment-only and mapped to the human checklist.
- No unresolved blocking defects remain.
- Required documentation updates or verified-no-change decisions are named in
  task completion reports.

Final review records every `VAL-NNN` outcome as `PASS`, `FAIL`, `BLOCKED`, or
`CARRY-FORWARD`. `PASS` without cited evidence is invalid.
