# Phase 2 QA Plan - Widen The Schema Vertically

> Phase: 2 - Widen The Schema Vertically
> Refreshed: 2026-06-07
> QA Lead: Codex
> Status: READY. No open phase-start ESC IDs; no command placeholders found.
>
> Lives at `project-docs/docs/phases/phase-2/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 2 completes the v0.1 SQLite schema vertically: migration `0002` adds the
remaining v0.1 base and stub tables; `0003` adds `issues_enriched` and its
indexes and wires `fid issue get`; `0004` and `0005` add the v0.1 SQL backstop
triggers; phase-end proves deterministic rebuild and the adapter boundary.

Source artifacts used:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/docs/phases/phase-2/qa/validation-contract.md`
- `project-docs/docs/phases/phase-2/tasks/task-1.md`
- `project-docs/docs/phases/phase-2/tasks/task-2.md`
- `project-docs/docs/phases/phase-2/tasks/task-3.md`
- `project-docs/docs/phases/phase-2/tasks/task-4.md`
- Referenced spec sections from `01_PRD.md`, `02_TRD.md`, `03_TRD_SCHEMA.md`,
  `05_DESIGN_PLAN.md`, and `06_DESIGN_PLAN_CLI.md`.

In-scope requirement coverage:

| Requirement source | Covered by |
|---|---|
| PRD-OPS-003, PRD-OPS-004 | VAL-004, VAL-016 |
| PRD-ISSUE-006 | VAL-007 |
| PRD-DEV-007 | VAL-005, VAL-007 |
| PRD-OUT-001, PRD-OUT-004 | VAL-007 |
| TRD §4.1, §4.3, §4.5, §5 | VAL-001..004, VAL-007, VAL-009..015 |
| TRD SCHEMA §1.2, §2.2, §2.3, §5, §7, §9, §10, §10.1, §11, §14.1, §14.2, §14.3, §14.4, §14.7 | VAL-001..003, VAL-005..006, VAL-008..015 |
| TRD-ARCH-001, TRD-ARCH-002, TRD-ARCH-005, TRD-ARCH-006, TRD-ARCH-008 | VAL-003, VAL-016, VAL-017, VAL-009..014, VAL-008 |
| TRD-PERF-004 | VAL-008 |
| TRD-SEC-004 | VAL-009..015 |
| Design `fid issue get`, `05_DESIGN_PLAN.md` §6.1 JSON, `06_DESIGN_PLAN_CLI.md` `fid issue get` | VAL-007 |

Scope deferrals accepted by the validation contract and not re-validated here:
sequence-vs-scheduling trigger, `activity_log` projection-only guard, and the
dedicated `activity` effect op.

Task 2 references PRD-ISSUE-018, but its `--section` behavior is explicitly out
of scope for Phase 2. QA must not claim PRD-ISSUE-018 PASS from VAL-007 evidence;
final documentation review should record the deferral if the task completion
report repeats the requirement reference.

## 2. Validation Contract Refinement Check

| Check | Result | Notes |
|---|---|---|
| All in-scope requirements referenced by at least one `VAL-NNN` | YES | Mapped in section 1 and in `validation-contract.md`. |
| Every `VAL-NNN` has one validator role | YES | All rows name QA Lead. |
| Every `VAL-NNN` has a concrete evidence method | YES | This plan refines each method into deterministic evidence rows and checklist items. |
| Every `VAL-NNN` maps to a task/window or phase-end check | YES | VAL-001..004 task 1; VAL-005..008 task 2; VAL-009..012 task 3; VAL-013..015 task 4; VAL-016..017 phase-end. |
| Every `VAL-NNN` describes one observable, or has acceptable fan-out | YES | Multi-case assertions are one rule family and are split into explicit both-state evidence requirements below. No human fan-out is required to prove VAL-001..017. |
| Conditional or binary assertions require both qualifying and non-qualifying states | YES | VAL-002, VAL-004..015 require both accept/reject or stale/fresh states. |
| Parity or "matches reference" assertions have deterministic conformance proof | YES | VAL-003 and VAL-016 require structural row equality after journal replay/rebuild; VAL-007 requires JSON/value parity before and after rebuild. |
| Revisions requested from Tech Lead | None | Contract is ready for task execution. |

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | Migration `0002`, journal reachability, projection schema version bump | VAL-001..004 | High | Schema transcription and rebuild freshness affect all later tasks. Completion report must include package-level `internal/journal` evidence for VAL-003 and `ProjectionSchemaVersion` stale/current cases for VAL-004. |
| Task 2 | `issues_enriched`, indexes, enriched `fid issue get` | VAL-005..008 | High | User-facing JSON and COV-001 carry-in. Completion report must include package-level `internal/issues` evidence, CLI evidence, both concern states, all triage states, rebuild parity, and EXPLAIN index use. |
| Task 3 | Kind-invariant and kind-immutability triggers | VAL-009..012 | Medium | Raw-SQL RAISE evidence is mandatory. Mid-phase QA review occurs after this task. |
| Task 4 | Polymorphic ref/link/claim triggers and phase-end invariants | VAL-013..017 | High | Must include trigger both-state evidence, full-schema journal rebuild determinism, phase-level extras, and adapter-boundary evidence. |

## 4. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | `internal/db` migration test output plus schema introspection proving all `0002` tables, representative columns, FKs, and static CHECKs exist. | VAL-001 | Developer, QA Lead verifies | Task 1 completion; final review |
| EVID-02 | `internal/db` up/down/up test output proving `0002` drops only its tables and preserves `projection_meta`, `roles`, and `issues`. | VAL-002 | Developer, QA Lead verifies | Task 1 completion; final review |
| EVID-03 | Package-level `internal/journal` test output proving journal effects populate `requirements`, `labels`, `epic_metadata`, `acceptance_criteria`, `task_completions`, and `concerns`; read-back equals effect values. | VAL-003 | Developer, QA Lead verifies | Task 1 completion; final review |
| EVID-04 | Freshness decision test output for prior-version `ProjectionRebuild` and current-version `ProjectionFresh`, plus rebuild evidence when a stored prior version exists. | VAL-004 | Developer, QA Lead verifies | Task 1 completion; final review |
| EVID-05 | View tests over journal-projected issues proving open concerns count and `concerns_raised` for open, addressed-only, and absent-concern states. | VAL-005 | Developer, QA Lead verifies | Task 2 completion; final review |
| EVID-06 | View tests proving latest triage verdict mapping for latest `RESUME`, latest `BLOCK`, and no triage. | VAL-006 | Developer, QA Lead verifies | Task 2 completion; final review |
| EVID-07 | Package-level `internal/issues` tests plus `cmd/fid` CLI JSON/integration output proving enriched fields are sourced from `issues_enriched` and survive delete-DB + rebuild. | VAL-007 | Developer, QA Lead verifies | Task 2 completion; final review |
| EVID-08 | `EXPLAIN QUERY PLAN` test output showing index search using `idx_concerns_issue_status` and `idx_triage_issue_created`, not table scans. | VAL-008 | Developer, QA Lead verifies | Task 2 completion; final review |
| EVID-09 | Raw-SQL trigger tests proving invalid `issues.parent_id` kind hierarchy writes RAISE and valid hierarchies insert. | VAL-009 | Developer, QA Lead verifies | Task 3 completion; mid-phase review; final review |
| EVID-10 | Raw-SQL trigger tests proving `acceptance_criteria`/`task_completions` task-or-subtask guards and `epic_metadata` epic guard, both reject and accept states. | VAL-010 | Developer, QA Lead verifies | Task 3 completion; mid-phase review; final review |
| EVID-11 | Raw-SQL trigger tests proving `triage_events.epic_id` and `decisions.epic_id` epic guards, including valid refs and allowed NULL decision epic. | VAL-011 | Developer, QA Lead verifies | Task 3 completion; mid-phase review; final review |
| EVID-12 | Raw-SQL trigger tests proving `UPDATE issues SET kind=...` RAISES and non-kind updates succeed. | VAL-012 | Developer, QA Lead verifies | Task 3 completion; mid-phase review; final review |
| EVID-13 | Raw-SQL trigger tests proving links, claims, and activity_log reject dangling polymorphic refs and accept resolvable refs. | VAL-013 | Developer, QA Lead verifies | Task 4 completion; final review |
| EVID-14 | Raw-SQL trigger tests proving link endpoint scope and relation compatibility for rejected decision/escalation endpoints, rejected incompatible combo, and accepted `implements`, `blocks`, and `relates_to` combos. | VAL-014 | Developer, QA Lead verifies | Task 4 completion; final review |
| EVID-15 | Raw-SQL trigger tests proving claims accept `ref_type='issue'` and reject unsupported ref types such as `qa_review`. | VAL-015 | Developer, QA Lead verifies | Task 4 completion; final review |
| EVID-16 | Full-schema journal fixture evidence: replay touches every new v0.1 table, delete-DB + `/tmp/fid rebuild` produces identical row/structural state, and all triggers remain active without false positives. | VAL-016 | Developer, QA Lead verifies | Task 4 completion; final review |
| EVID-17 | Import graph and source scan output proving no `internal/...` imports `cmd/fid`; source review notes showing enriched-read logic lives in `internal/issues` and `cmd/fid` delegates. | VAL-017 | Developer, QA Lead verifies | Task 4 completion; final review |
| EVID-18 | Task gate outputs for each task: format, lint, race coverage, coverage report, plus conditional commands required by the task spec. | VAL-001..017, no VAL for mechanical quality | Developer, QA Lead verifies | Every completion report; reviews |
| EVID-19 | Phase-level command gate output: lint, tests, coverage, integration, benchmarks, build, PRD compliance, security review, performance spot-check, documentation review, deterministic checklist. | VAL-001..017, no VAL for mechanical quality | QA Lead | Final review |

Coverage requirements:

- VAL-003 must have package-level `internal/journal` tests, not only `cmd/fid`
  or `internal/db` evidence. Final coverage evidence must show
  `github.com/mbell-nzl/fid/internal/journal` at or above 80%.
- VAL-007 must have package-level `internal/issues` tests plus CLI evidence.
  Final coverage evidence must show `github.com/mbell-nzl/fid/internal/issues`
  at or above 80% and `github.com/mbell-nzl/fid/cmd/fid` at or above 70%.
- `internal/policy` is a standing COV-001 package but is not changed by Phase 2.
  Final review still records whether it remains at or above 80% if it appears in
  the refreshed coverage profile.

Acceptance criteria and expected tests:

| Area | Acceptance criteria | Expected tests |
|---|---|---|
| Migration `0002` | All remaining v0.1 base/stub tables exist with required static constraints; down migration is clean; stale projection versions rebuild. | `internal/db` migration introspection, up/down/up, freshness decision tests; `internal/journal` package tests for representative insert effects. |
| Journal reachability | Representative new table groups are populated by journal effects and read back equal. | Package-level `internal/journal` replay/projection tests for `requirements`, `labels`, `epic_metadata`, `acceptance_criteria`, `task_completions`, and `concerns`. |
| `issues_enriched` and `fid issue get` | Concern count/flag and latest triage status are correct in the view and surfaced by CLI JSON before and after rebuild. | `internal/db` view tests, `internal/issues` package-level enriched `Get` tests, `cmd/fid` CLI JSON/rebuild tests. |
| Index usage | Concern and triage joins use the required indexes, not table scans. | `internal/db` `EXPLAIN QUERY PLAN` assertion tolerant of SQLite wording but strict on named index use. |
| Kind triggers | Invalid raw-SQL structural writes RAISE; valid writes continue to insert/update. | `internal/db` trigger tests for both reject and accept states for VAL-009..012. |
| Ref/link/claim triggers | Invalid polymorphic refs, endpoint scopes, relation combos, and claim scopes RAISE; valid rows insert. | `internal/db` trigger tests for both reject and accept states for VAL-013..015. |
| Full rebuild and boundary | Full-schema journal rebuild is deterministic with triggers active; `internal/` has no `cmd/fid` dependency. | `internal/journal` or CLI-backed rebuild determinism test, `/tmp/fid rebuild` smoke, `go list` import graph, `rg` source scan, source-review notes. |

## 5. Deterministic Gates

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task and phase end |
| Lint/static analysis | `golangci-lint run ./...` | Exit 0 | Every task and phase end |
| Race coverage gate | `go test -race -coverprofile=coverage.out ./...` | Exit 0; `coverage.out` produced | Every task and phase end |
| Coverage report | `go tool cover -func=coverage.out` | Exit 0; package-level coverage thresholds recorded | Every task and phase end |
| Integration | `go test -tags integration -race ./...` | Exit 0, or explicit task report `N/A` only when no tagged tests apply before Task 4 | Conditional for Tasks 1-3; required for Task 4 and phase end |
| Benchmarks | `go test -bench . -benchmem ./...` | Exit 0; benchmark output recorded as spot-check evidence | Task 4 and phase end |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | Exit 0; `/tmp/fid` exists and is executable | Task 2, Task 4, phase end; recommended for Task 1; conditional for Task 3 |
| Phase-level unit/race | `go test -race ./...` | Exit 0 | Phase-end final review per Quality Plan §5.2 |
| Adapter boundary import graph | `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/...` | Output has no `cmd/fid` import | Task 4 and phase end |
| Adapter boundary source scan | `rg -n "github.com/mbell-nzl/fid/cmd/fid|cmd/fid" internal` | No matches; `rg` exit 1 is acceptable evidence for no matches | Task 4 and phase end |

Task gate alignment:

- Tasks 1-4 all require exactly:
  `gofmt -l cmd internal`, `golangci-lint run ./...`,
  `go test -race -coverprofile=coverage.out ./...`, and
  `go tool cover -func=coverage.out`.
- Task 1 conditionally requires `go test -tags integration -race ./...` if
  migration/rebuild tests are tagged, and recommends `go build -o /tmp/fid
  ./cmd/fid`.
- Task 2 requires `go build -o /tmp/fid ./cmd/fid` and conditionally requires
  integration tests if CLI/rebuild parity tests are tagged.
- Task 3 conditionally requires integration tests if trigger tests are tagged.
- Task 4 requires `go test -tags integration -race ./...`,
  `go test -bench . -benchmem ./...`, and `go build -o /tmp/fid ./cmd/fid`.

## 6. Command And Feasibility Validation

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| Phase branch `phase-2/schema-widening` | Branch file-state check | PASS | `git branch --list phase-2/schema-widening` printed `* phase-2/schema-widening`. |
| Task specs and validation contract | File-state check | PASS | `find project-docs/docs/phases/phase-2/qa project-docs/docs/phases/phase-2/tasks -maxdepth 1 -type f` listed `validation-contract.md` and `task-1.md` through `task-4.md`. |
| QA/final/human review templates | File-state check | PASS | `find project-docs/templates -maxdepth 1 -type f` listed `qa-plan-template.md`, `mid-phase-review-template.md`, `final-review-template.md`, and `human-checklist-template.md`. |
| Go module and package pattern `./...` | Module check | PASS | `go list -m` returned `github.com/mbell-nzl/fid`; `find . -maxdepth 2 -type f -name go.mod` returned `./go.mod`. |
| `gofmt -l cmd internal` | Help output | PASS | `gofmt -h` printed `usage: gofmt [flags] [path ...]` and documented `-l` as listing files whose formatting differs. |
| `golangci-lint run ./...` | Help output and config file check | PASS | `golangci-lint --help` listed command `run` as "Lint the code"; `find` found `./.golangci.yml`. |
| `go test -race -coverprofile=coverage.out ./...` | Help output | PASS | `go help test` validates `go test [build/test flags] [packages]`; `go help build` documents `-race`; `go help testflag` documents `-coverprofile`. |
| `go tool cover -func=coverage.out` | Tool/version and usage output | PASS | `go tool cover -V` returned `cover version go1.26.2`; `go tool cover -h` documented `-func string output coverage profile information for each function`. |
| `go test -tags integration -race ./...` | Help output | PASS | `go help build` documents `-tags` and `-race`; `go help test` validates package-list mode over `./...`. |
| `go test -bench . -benchmem ./...` | Help output | PASS | `go help testflag` documents `-bench regexp` and `-benchmem`. |
| `go build -o /tmp/fid ./cmd/fid` | Help output and path feasibility | PASS | `go help build` documents `go build [-o output] [packages]`; `/tmp` is writable per workspace permissions; `cmd/fid/main.go` exists. |
| Adapter import graph command | Existing phase evidence and command shape | PASS | Phase 1 final review used `go list -f '{{.ImportPath}} {{join .Imports " "}}' ./internal/...`; `go list -m` confirms current module. |
| Adapter source scan command | Tool availability and path check | PASS | `rg` is approved and available; `find internal ...` listed `internal/issues`, `internal/journal`, and `internal/db` source files. |
| Migration/test file assumptions | File-state check | PASS | Existing `migrations/0001_v01a_base.*.sql`, `internal/db/*`, `internal/journal/*`, `internal/issues/*`, and `cmd/fid/*` files are present; Phase 2 tasks name new `0002`-`0005` files to be created by Developers. |
| Coverage threshold extraction | Prior evidence + command output shape | PASS | Phase 1 final review records package coverage extraction from `go tool cover -func=coverage.out`, including `internal/issues 85.2%`, `internal/journal 80.6%`, and `cmd/fid 70.5%`. Phase 2 final review must refresh these numbers from a new profile. |
| Held-out scenarios | Scope check | PASS | Task specs and phase brief state held-out scenario impact is none and "Do not read held-out scenarios"; this plan does not reference hidden scenario paths or commands. |
| Open escalations at phase start | Required context check | PASS | `tl_escalations_current.md` says `_None open._`; user confirmed baseline open ESC IDs: none. |
| Task 1 command scope | Task-spec inspection | PASS | `task-1.md` names exact required commands and conditional integration/build commands; no placeholders. |
| Task 2 command scope | Task-spec inspection | PASS | `task-2.md` names exact required commands, requires `go build -o /tmp/fid ./cmd/fid`, and conditionally names integration tests; no placeholders. |
| Task 3 command scope | Task-spec inspection | PASS | `task-3.md` names exact required commands and conditionally names integration tests; no placeholders. |
| Task 4 command scope and phase-end extras | Task-spec inspection | PASS | `task-4.md` names exact required commands plus required integration, benchmark, and build commands for phase-end evidence; no placeholders. |
| VAL-016 full-schema fixture assumption | File-state and task-spec inspection | PASS | Task 4 requires building the determinism fixture with journal events; existing `internal/journal/project.go`, `internal/journal/replay.go`, and `cmd/fid/rebuild.go` files are present. |
| VAL-017 boundary assumption | File-state and command validation | PASS | Existing `cmd/fid/issue.go`, `internal/issues/issue.go`, and import graph/source-scan commands are available for final boundary proof. |
| HUMAN-01 feasibility | Human-checklist scope review | PASS | Requires product-owner judgment over enriched JSON usefulness only; no hidden scenario or deterministic substitute needed. |
| HUMAN-02 feasibility | Human-checklist scope review | PASS | Requires dogfooding on a disposable workspace after `/tmp/fid` is built by final QA; uses visible CLI behavior only. |
| HUMAN-03 feasibility | Human-checklist scope review | PASS | Requires human acceptance/rejection of non-blocking QA risks after final review exists. |

No referenced command in `AGENTS.md`, task specs, or this plan is a placeholder.

## 7. Deterministic Checklist

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify Task 1 completion report records the exact required gate commands as PASS or BLOCKED, with no narrowed substitutes `(qa-lead)` | VAL-001..004 | `task-1-completion.md` command table plus command excerpts. |
| QA-02 | Verify migration `0002` up/schema introspection evidence includes every required table, representative columns, FKs, and static CHECKs `(qa-lead)` | VAL-001 | `internal/db` test output or logs naming all required tables and checks. |
| QA-03 | Verify migration `0002` down and up/down/up evidence preserves Phase-1 tables and recreates new tables `(qa-lead)` | VAL-002 | Test output and schema-state excerpts. |
| QA-04 | Verify journal-effect evidence populates representative new table groups and read-back equals effect values through `internal/journal` package-level tests `(qa-lead)` | VAL-003 | `internal/journal` test names/output; no direct-SQL happy-path-only evidence accepted. |
| QA-05 | Verify projection freshness evidence covers prior-version stale rebuild and current-version fresh states `(qa-lead)` | VAL-004 | Test output for both `ProjectionRebuild` and `ProjectionFresh`; completion report records final `ProjectionSchemaVersion` after Task 1. |
| QA-06 | Verify Task 2 completion report records the exact required gate commands and required build command as PASS or BLOCKED `(qa-lead)` | VAL-005..008 | `task-2-completion.md` command table plus command excerpts. |
| QA-07 | Verify `issues_enriched` concern-count evidence covers open concerns, addressed-only concerns, and absent concerns `(qa-lead)` | VAL-005 | View test output and expected/actual values. |
| QA-08 | Verify `issues_enriched.triage_status` evidence covers latest `RESUME`, latest `BLOCK`, and no triage `(qa-lead)` | VAL-006 | View test output and expected/actual values. |
| QA-09 | Verify `fid issue get` enriched output evidence includes package-level `internal/issues` tests, CLI JSON evidence, zero/nonzero states, and delete-DB rebuild parity `(qa-lead)` | VAL-007 | `internal/issues` and `cmd/fid` test output; JSON excerpts showing `concerns_raised`, `open_concerns_count`, and `triage_status`. |
| QA-10 | Verify query-plan evidence shows both required indexes are used and no full scan occurs on `concerns` or `triage_events` `(qa-lead)` | VAL-008 | `EXPLAIN QUERY PLAN` test output showing `SEARCH ... USING INDEX` or equivalent index-use assertion. |
| QA-11 | Verify Task 3 completion report records the exact required gate commands as PASS or BLOCKED and records the `task_completions` task-or-subtask resolution `(qa-lead)` | VAL-009..012 | `task-3-completion.md` command table and documented resolution. |
| QA-12 | Verify raw-SQL `issues.parent_id` trigger evidence covers each invalid hierarchy and each valid hierarchy `(qa-lead)` | VAL-009 | Trigger test output with RAISE excerpts and valid insert control cases. |
| QA-13 | Verify raw-SQL kind guard evidence covers wrong-kind rejection and valid-kind insertion for `acceptance_criteria`, `task_completions`, and `epic_metadata` `(qa-lead)` | VAL-010 | Trigger test output with RAISE excerpts and valid insert control cases. |
| QA-14 | Verify raw-SQL `epic_id` guard evidence covers non-epic rejection, valid epic references, and NULL `decisions.epic_id` acceptance `(qa-lead)` | VAL-011 | Trigger test output with RAISE excerpts and valid insert control cases. |
| QA-15 | Verify raw-SQL kind immutability evidence covers rejected kind update and accepted non-kind update `(qa-lead)` | VAL-012 | Trigger test output with RAISE excerpt and valid update control. |
| QA-16 | Run mid-phase review after Task 3 before Task 4 begins; verify VAL-001..012 have evidence or documented blockers `(qa-lead)` | VAL-001..012 | `project-docs/docs/phases/phase-2/qa/review-3.md` using the mid-phase review template. |
| QA-17 | Verify Task 4 completion report records the exact required gate commands plus integration, benchmarks, and build as PASS or BLOCKED `(qa-lead)` | VAL-013..017 | `task-4-completion.md` command table plus command excerpts. |
| QA-18 | Verify raw-SQL polymorphic ref evidence covers dangling rejection and resolvable acceptance for links, claims, and activity_log `(qa-lead)` | VAL-013 | Trigger test output with RAISE excerpts and valid insert controls. |
| QA-19 | Verify raw-SQL link scope/relation evidence covers out-of-scope endpoint rejection, incompatible relation rejection, and accepted `implements`, `blocks`, and `relates_to` cases `(qa-lead)` | VAL-014 | Trigger test output with RAISE excerpts and valid insert controls. |
| QA-20 | Verify raw-SQL claim scope evidence covers `ref_type='issue'` acceptance and unsupported ref-type rejection `(qa-lead)` | VAL-015 | Trigger test output with RAISE excerpt and valid insert control. |
| QA-21 | Execute phase-level format, lint, race coverage, coverage report, integration, benchmark, and build gates before final verdict `(qa-lead)` | VAL-001..017, no VAL | Literal command output excerpts in `final-review.md`; no PASS without output. |
| QA-22 | Verify final coverage from the refreshed `coverage.out` keeps `internal/journal` >= 80%, `internal/issues` >= 80%, and `cmd/fid` >= 70% `(qa-lead)` | VAL-003, VAL-007, no VAL | `go tool cover -func=coverage.out` output or reproducible package aggregation. |
| QA-23 | Verify full-schema rebuild determinism over a journal touching every new v0.1 table with all triggers active `(qa-lead)` | VAL-016 | `/tmp/fid` rebuild command output plus row/structural equality evidence before/after delete-DB. |
| QA-24 | Verify adapter boundary with import graph and source scan, then source-review `cmd/fid` for thin delegation `(qa-lead)` | VAL-017 | `go list` output, `rg` no-match evidence, and concise source-review notes. |
| QA-25 | Verify PRD/TRD/design compliance, security review, performance spot-check, and documentation review sections are completed in final review `(qa-lead)` | VAL-001..017, no VAL | `final-review.md` sections with evidence and findings. |
| QA-26 | Record final outcome for every VAL-001..VAL-017 as PASS, FAIL, BLOCKED, or CARRY-FORWARD, with cited evidence for every PASS `(qa-lead)` | VAL-001..017 | Final review outcome table. |

## 8. Mid-Phase Review Cadence

- Cadence: one QA Reviewer mid-phase review after Task 3 completes and before
  Task 4 begins. This matches the phase runner expectation and avoids reviewing
  the final phase-end assertions before their task exists.
- Review artifact: `project-docs/docs/phases/phase-2/qa/review-3.md`.
- Template: `project-docs/templates/mid-phase-review-template.md`.
- Review scope: Tasks 1-3, VAL-001..012, exact task-gate evidence, coverage
  trend for `internal/journal` and `internal/issues`, and trigger RAISE evidence
  quality for Task 3.
- Final review: after Task 4 completion and before phase wrap-up. Final review
  writes `project-docs/docs/phases/phase-2/qa/final-review.md` and generates
  `project-docs/docs/phases/phase-2/qa/human-checklist.md`.

If Task 1 or Task 2 reports a blocking quality-gate failure, missing completion
report, missing VAL evidence, or a silently substituted command, QA Lead may add
an earlier review. That is stricter than the default and must be recorded in the
new review artifact.

## 9. Human Verification Scope

VAL-001..VAL-017 are deterministic in this phase. Human checks are not required
to prove any VAL, but the final human checklist should retain phase sign-off
judgment and dogfooding items.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-01 | Human reviews whether the enriched `fid issue get` JSON remains useful for real agent workflow consumption `(human)` | no VAL | Product-owner note or approval in human checklist. |
| HUMAN-02 | Human dogfoods `/tmp/fid issue get` and `/tmp/fid rebuild` on a disposable workspace after QA final review `(human)` | no VAL | Manual command notes, including any confusing output or approval. |
| HUMAN-03 | Human reviews unresolved QA findings and decides whether any non-blocking risks are acceptable for phase wrap-up `(human)` | no VAL | Acceptance note or requested follow-up. |

## 10. Phase-End Decision Rule

The phase can proceed to wrap-up only when:

- Task 1-4 completion reports exist and record every exact required command as
  PASS or BLOCKED with reason.
- The mid-phase review after Task 3 exists and has no unresolved blocking
  finding before Task 4 starts.
- Phase-level gates in section 5 pass, or failures are recorded as BLOCKED with
  escalation.
- Every VAL-001..VAL-017 has a final outcome and cited deterministic evidence.
- Package-level coverage evidence shows `internal/journal` >= 80% and
  `internal/issues` >= 80%; CLI coverage shows `cmd/fid` >= 70%.
- Deterministic checklist items QA-01..QA-26 have reproducible evidence.
- Human checklist contains only human-judgment items and does not substitute for
  deterministic VAL evidence.
- No unresolved blocking defects, missing artifacts, command placeholders, or
  unvalidated command assumptions remain.
