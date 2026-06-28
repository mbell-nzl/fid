# Phase 13 QA Plan - Labels + v0.1d Closure

> Phase: 13 - Labels + v0.1d closure
> Refreshed: 2026-06-27
> QA Lead: Codex QA Lead
>
> Location:
> `project-docs/docs/phases/phase-13/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 13 closes v0.1d by adding the label CLI surface and proving the full
v0.1 lifecycle smoke. Scope is limited to the three task specs:

- Task 1: `fid label add` / `fid label list` and the `label.add` matrix row.
- Task 2: `fid issue label add` / `fid issue label remove`, matrix rows, and
  per-command `activity_log` projection for issue-label deltas.
- Task 3: an integration smoke proving create -> claim -> transition ->
  complete -> unclaim with permission enforcement, links, labels, one event per
  write, and rebuild-stable state.

Explicit scope boundaries from the phase brief:

- No schema migration. `labels` and `issue_labels` already exist in
  `migrations/0002` and `03_TRD_SCHEMA.md` section 9.6.
- No projector change. Existing generic insert/delete projection is the expected
  path.
- No scope beyond the task specs. Broad CF-011 `activity_log` reconciliation,
  CF-014 file splitting, label update/delete/rename, configurable permissions,
  and release-readiness scans remain out of scope unless a task records a
  spec-backed blocker.

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
- `project-docs/docs/phases/phase-13/qa/validation-contract.md`
- `project-docs/docs/phases/phase-13/tasks/task-1.md`
- `project-docs/docs/phases/phase-13/tasks/task-2.md`
- `project-docs/docs/phases/phase-13/tasks/task-3.md`
- `project-docs/templates/qa-plan-template.md`

## 2. Requirement Coverage

| Requirement | Phase coverage | VAL IDs |
|---|---|---|
| PRD-ISSUE-015 | Label metadata exists through label creation/listing. | VAL-002, VAL-003 |
| PRD-ISSUE-023 | Issue labels can be added/removed without replacing the full set and changes are journaled/projected to `activity_log`. | VAL-005, VAL-006, VAL-007, VAL-009, VAL-010 |
| PRD-OUT-001..004 | Label list supports json default plus toon, markdown, yaml with the same fields. | VAL-003 |
| TRD-ARCH-010 | New write actions are represented once in the v0.1 in-code matrix; unknown and denied roles reject. | VAL-001, VAL-004, VAL-008 |
| `03_TRD_SCHEMA.md` section 9.6 | Existing `labels` / `issue_labels` tables support the phase without migration. | VAL-002, VAL-005, VAL-006 |
| `03_TRD_SCHEMA.md` section 13 | `label.add`, `issue.label.add`, and `issue.label.remove` matrix booleans match the schema reference. | VAL-001, VAL-004, VAL-008 |
| `06_DESIGN_PLAN_CLI.md` label commands | `fid label add` / `fid label list` spellings and fields are implemented. | VAL-002, VAL-003 |
| `06_DESIGN_PLAN_CLI.md` issue-label commands | `fid issue label add` / `remove` are single-label delta commands. | VAL-005, VAL-006 |
| PRD-CLAIM-005; PRD-WF-001..010; PRD-NFR-008; ADR-0001 | Closure smoke proves gated lifecycle, write-lock discipline, and rebuild parity. | VAL-009 |
| `08_QUALITY_PLAN.md` sections 5.2, 5.6 | Deterministic phase evidence and human dogfood checklist are produced separately. | VAL-010, no VAL gates |

## 3. Validation Contract Refinement Check

| Check | Result | Notes |
|---|---|---|
| All in-scope requirements referenced by at least one `VAL-NNN` | YES | See requirement coverage table. |
| Every `VAL-NNN` has one validator role | YES | VAL-001..009: QA Lead; VAL-010: Human. |
| Every `VAL-NNN` has a concrete evidence method | YES | Refined in section 5. |
| Every `VAL-NNN` maps to a task/window or phase-end check | YES | Task/window mapping in section 4. |
| Every `VAL-NNN` describes one observable, or human-checklist fan-out is explicitly required | YES | VAL-001 and VAL-003 are structural parity observables; VAL-009 is one scripted closure-smoke trace with required sub-evidence; VAL-010 is human-only and fans out to HUMAN-13a..HUMAN-13d. |
| Conditional or binary `VAL-NNN` assertions require both states | YES | VAL-002, VAL-004, VAL-005, VAL-006, VAL-007, VAL-008, and VAL-009 require success plus rejection/no-append or presence/absence evidence. |
| Parity or "matches reference" assertions have deterministic conformance proof when available | YES | VAL-001 requires matrix/show/enforce parity; VAL-003 requires four-format structural parity; VAL-009 requires rebuild parity. |
| Revisions requested from Tech Lead | None | No validation-contract blocker found. If Task 3 evidence proves only part of VAL-009, final QA must fail or request a contract split rather than passing partial proof. |

## 4. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | Label CRUD and `label.add` matrix row | VAL-001 partial, VAL-002, VAL-003, VAL-004, VAL-007 partial | Medium | Validates first label creation surface, list parity, duplicate rejection, and no `activity_log` row for label creation. |
| Task 2 | Single-label issue assignment, matrix rows, `activity_log` | VAL-001 complete, VAL-005, VAL-006, VAL-007, VAL-008 | High | Must not reuse replace-all `--labels`; must prove validate-before-append and one-row delta semantics. |
| Task 3 | v0.1d lifecycle integration smoke | VAL-009 | High | Test-only task. Mandatory integration gate. Any product defect found must be escalated, not fixed inside Task 3. |
| Phase end | Human dogfood and final QA outcome | VAL-010; final outcomes for VAL-001..010 | Medium | Human checklist contains only judgment/dogfood rows. |

## 5. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Unit/CLI output proving exactly one `label.add`, `issue.label.add`, and `issue.label.remove` matrix row, show-vs-enforce parity, and allowed/denied role checks against schema section 13. | VAL-001 | Developer, QA Lead | Tasks 1-2 reports; mid-phase review; final review |
| EVID-02 | Domain/CLI test output and journal/DB state proving `fid label add` creates one `labels` row in one event, and duplicate add exits non-zero with zero new journal lines/rows. | VAL-002 | Developer, QA Lead | Task 1 report; mid-phase review; final review |
| EVID-03 | CLI output or parsed test evidence for `fid label list` in json, toon, markdown, and yaml, comparing every label's `name`, `colour`, and `description` to the JSON reference. | VAL-003 | Developer, QA Lead | Task 1 report; mid-phase review; final review |
| EVID-04 | CLI output proving `label.add` accepts developer/tl/workflow_runner/human and rejects a denied role plus an unknown role before append, with error text naming role and `label.add`. | VAL-004 | Developer, QA Lead | Task 1 report; mid-phase review; final review |
| EVID-05 | Domain/CLI test output proving `fid issue label add` adds exactly one existing label, preserves pre-existing labels, emits one event, and rejects nonexistent/already-present labels before append. | VAL-005 | Developer, QA Lead | Task 2 report; mid-phase review; final review |
| EVID-06 | Domain/CLI test output proving `fid issue label remove` removes exactly one assigned label, preserves remaining labels, emits one event, and rejects absent-label removal before append. | VAL-006 | Developer, QA Lead | Task 2 report; mid-phase review; final review |
| EVID-07 | DB/query or CLI-backed test evidence proving issue-label add and remove each project exactly one `activity_log` row with issue, label, and direction, while `fid label add` projects none. | VAL-007 | Developer, QA Lead | Tasks 1-2 reports; mid-phase review; final review |
| EVID-08 | CLI output proving `issue.label.add` and `issue.label.remove` accept permitted roles and reject a denied role plus an unknown role before append, naming role and action. | VAL-008 | Developer, QA Lead | Task 2 report; mid-phase review; final review |
| EVID-09 | Integration test output from `go test -tags integration -race ./...` proving the full lifecycle trace, started-state rejection without a live claim, one event per successful write, at least one denied-role rejection with zero append, no nested-lock failure, and rebuild parity for issues/claims/links/issue_labels. | VAL-009 | Developer, QA Lead | Task 3 report; final review |
| EVID-10 | Human dogfood checklist rows, one observable per row, confirming label creation, non-destructive add, non-destructive remove, and end-to-end lifecycle are understandable from CLI output. | VAL-010 | Human | Phase-end human checklist |
| EVID-11 | Exact task-gate command outputs and coverage percentages; required commands marked PASS or BLOCKED with reason. | no VAL | Developer, QA Reviewer, QA Lead | Every task report; reviews |
| EVID-12 | File-state evidence proving no migration file was added and no projector change was made, unless a task recorded a spec-backed blocker. | no VAL | QA Lead | Final review |

## 6. Acceptance Criteria

- VAL-001..009 have deterministic evidence slots and cannot pass from
  code-reading alone.
- VAL-010 is human-only and fans out into HUMAN-13a..HUMAN-13d.
- Every task report records the exact required commands as `PASS`, or `BLOCKED`
  with a blocker reason. A missing, silently skipped, or narrower substituted
  command is a blocking QA finding.
- Task 3 records `go test -tags integration -race ./...` as mandatory PASS or
  BLOCKED; `N/A` is not allowed for that command.
- Final QA confirms no schema migration, no projector change, and no
  implementation scope beyond the task specs unless an escalation/decision
  explicitly records the deviation.

## 7. Expected Tests

| Area | Expected tests | Covers |
|---|---|---|
| `internal/labels` | Add/list service tests for successful add, optional colour/description, duplicate rejection with zero append, list fields, and no label-create `activity_log` row. | VAL-002, VAL-003, VAL-007 |
| `cmd/fid/label.go` | CLI tests for `fid label add`, `fid label list --format json|toon|markdown|yaml`, denied role, and unknown role. | VAL-002, VAL-003, VAL-004 |
| `internal/policy` | Matrix completeness tests for `label.add`, `issue.label.add`, `issue.label.remove`; allowed/denied role checks match schema section 13. | VAL-001, VAL-004, VAL-008 |
| `internal/issues` | Single-label add/remove service tests for one-row delta, existing-label precheck, already-present/absent rejection, remaining-label preservation, and `activity_log` direction. | VAL-005, VAL-006, VAL-007 |
| `cmd/fid/issue_label.go` | CLI tests for `fid issue label add` / `remove`, permitted roles, denied/unknown roles, and zero-append rejection paths. | VAL-005, VAL-006, VAL-008 |
| `cmd/fid/e2e_v01d_closure_test.go` | Integration test with `//go:build integration` for the full lifecycle closure smoke and rebuild parity. | VAL-009 |

## 8. Exact Task Gates

Every command below is exact. Focused commands may be used during TDD, but they
do not substitute for the close gate.

| Task | Required gate commands | Conditional commands |
|---|---|---|
| Task 1 | `gofmt -l cmd internal`<br>`golangci-lint run ./...`<br>`go test -race -coverprofile=coverage.out ./...`<br>`go tool cover -func=coverage.out`<br>`go build -o /tmp/fid ./cmd/fid` | `go test -tags integration -race ./...` only if an integration fixture is touched; otherwise record `N/A` with reason. |
| Task 2 | `gofmt -l cmd internal`<br>`golangci-lint run ./...`<br>`go test -race -coverprofile=coverage.out ./...`<br>`go tool cover -func=coverage.out`<br>`go build -o /tmp/fid ./cmd/fid` | `go test -tags integration -race ./...` only if an integration fixture is touched; otherwise record `N/A` with reason. |
| Task 3 | `gofmt -l cmd internal`<br>`golangci-lint run ./...`<br>`go test -race -coverprofile=coverage.out ./...`<br>`go tool cover -func=coverage.out`<br>`go test -tags integration -race ./...` | `go build -o /tmp/fid ./cmd/fid` only if non-test wiring changed; otherwise record `N/A` with reason. |

Task 3 note: `go test -tags integration -race ./...` is mandatory because the
task deliverable is the integration closure smoke.

Phase-end QA gate:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

## 9. Command And Feasibility Validation

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Help output and path existence. | PASS | `gofmt -h` shows `usage: gofmt [flags] [path ...]` and `-l`; `ls cmd internal` found both paths. |
| `golangci-lint run ./...` | Help output. | PASS | `golangci-lint run --help` shows usage `golangci-lint run [flags]`. |
| `go test -race -coverprofile=coverage.out ./...` | Go help output. | PASS | `go help testflag` documents `-coverprofile`; `go help build` documents `-race`; package pattern `./...` is standard Go package-list syntax. |
| `go tool cover -func=coverage.out` | Tool help output. | PASS | `go tool cover -help` documents `-func string` as coverage percentage output. |
| `go test -tags integration -race ./...` | Go help output plus existing integration-tag files. | PASS | `go help build` documents `-tags` and `-race`; `rg -n "go:build integration" cmd/fid` found existing integration tests including `e2e_v01d_link_test.go`. |
| `go test -bench . -benchmem ./...` | Go help output. | PASS | `go help testflag` documents `-bench regexp` and `-benchmem`. |
| `go build -o /tmp/fid ./cmd/fid` | Go help output and path existence. | PASS | `go help build` documents `-o output`; `ls cmd/fid` found the CLI main package path; `/tmp` is an approved writable root. |
| Task specs | File-state check. | PASS | `rg --files project-docs/docs/phases/phase-13` found `tasks/task-1.md`, `tasks/task-2.md`, and `tasks/task-3.md`. |
| QA templates | File-state check. | PASS | `ls project-docs/templates/mid-phase-review-template.md project-docs/templates/human-checklist-template.md` found both templates. |
| Existing label schema | Schema inspection. | PASS | `sed -n '640,660p' project-docs/specs/03_TRD_SCHEMA.md` shows `CREATE TABLE labels` and `CREATE TABLE issue_labels`. |
| Section-13 permission reference | Schema inspection. | PASS | `sed -n '880,915p' project-docs/specs/03_TRD_SCHEMA.md` shows `issue.label.add`, `issue.label.remove`, and `label.add` rows with developer, tl, workflow_runner, and human permitted. |
| Label CLI spellings | Design inspection. | PASS | `sed -n '960,982p' project-docs/specs/06_DESIGN_PLAN_CLI.md` shows `fid label add <name> [--colour <hex>] [--description <text>] --by <role>` and `fid label list [--format json|toon|markdown|yaml]`. |
| Issue-label CLI spellings | Design inspection. | PASS | `sed -n '220,250p' project-docs/specs/06_DESIGN_PLAN_CLI.md` shows `fid issue label add <issue-id> <label> --by <role>` and `fid issue label remove <issue-id> <label> --by <role>`. |
| Placeholder command scan for phase inputs | Text search. | PASS | `rg -n "<E2E_COMMAND>|<TBD>|placeholder|N/A" project-docs/docs/phases/phase-13 project-docs/docs/current/phase_brief.md project-docs/specs/08_QUALITY_PLAN.md` found only expected `N/A` conditional-command notes; no placeholder command blocks phase planning. |
| Temp-workspace lifecycle fixture assumption | Existing integration-test path check. | PASS | `cmd/fid/e2e_v01d_link_test.go` exists and the Task 3 spec names its helpers as the pattern for temp workspace, journal-line, nested-lock, and rebuild-parity evidence. |
| Human checklist feasibility | Template and VAL inspection. | PASS | `human-checklist-template.md` exists; VAL-010 explicitly requires one row per observable after deterministic QA passes. |

## 10. Deterministic Checklist

Every item is executed by QA Lead during final review unless the item is
explicitly assigned to the mid-phase review window. PASS requires reproducible
evidence in the review artifact.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify each task completion report records every exact required command for that task as `PASS`, or `BLOCKED` with a reason; Task 3 must include mandatory `go test -tags integration -race ./...` `(qa-lead)` | no VAL | File-state proof from `project-docs/docs/phases/phase-13/tasks/task-*-done.md` plus copied command-result excerpts. |
| QA-02 | Run the phase-end gate exactly as listed in section 8 `(qa-lead)` | no VAL | Literal command output for all seven phase-end commands. |
| QA-03 | Verify the in-code matrix has exactly one row each for `label.add`, `issue.label.add`, `issue.label.remove`; verify `fid policy show` and `Authorize` use the same role booleans as schema section 13 `(qa-lead)` | VAL-001 | Test output and CLI output excerpt for policy show; schema-reference comparison. |
| QA-04 | Verify `fid label add` success and duplicate rejection, including one event on success and zero new journal lines/rows on duplicate `(qa-lead)` | VAL-002 | Test output plus journal-line/DB-state excerpt or task evidence trace. |
| QA-05 | Verify `fid label list` in json, toon, markdown, and yaml exposes the same label set and `name`/`colour`/`description` fields `(qa-lead)` | VAL-003 | Parsed output comparison or literal CLI excerpts for all four formats. |
| QA-06 | Verify `label.add` permitted, denied, and unknown role behavior, including zero append on rejection and error text naming role plus action `(qa-lead)` | VAL-004 | CLI/test output and journal-line count excerpt. |
| QA-07 | Verify `fid issue label add` adds exactly one existing label, preserves other labels, emits one event, and rejects nonexistent/already-present labels before append `(qa-lead)` | VAL-005 | CLI/test output, issue-label state, and journal-line count excerpt. |
| QA-08 | Verify `fid issue label remove` removes exactly one assigned label, preserves remaining labels, emits one event, and rejects absent-label removal before append `(qa-lead)` | VAL-006 | CLI/test output, issue-label state, and journal-line count excerpt. |
| QA-09 | Verify issue-label add/remove each project exactly one `activity_log` row with issue, label, and direction; verify `fid label add` projects no `activity_log` row `(qa-lead)` | VAL-007 | DB/query or test output excerpt for all three states. |
| QA-10 | Verify `issue.label.add` and `issue.label.remove` permitted, denied, and unknown role behavior, including zero append on rejection and error text naming role plus action `(qa-lead)` | VAL-008 | CLI/test output and journal-line count excerpt. |
| QA-11 | Verify the Task 3 integration smoke proves the full lifecycle trace, rejected started transition without live claim, denied-role rejection, one event per successful write, no nested-lock failure, and rebuild parity `(qa-lead)` | VAL-009 | `go test -tags integration -race ./...` output and task evidence trace. |
| QA-12 | Verify Phase 13 did not add schema migrations and did not change `internal/journal/project_effects.go` unless a task records a spec-backed blocker `(qa-lead)` | no VAL | `git diff --name-only main...HEAD` file-state output and any recorded blocker/decision references. |
| QA-13 | Verify documentation updates or verified-no-change notes required by Tasks 1-3 are present in completion reports and phase decisions `(qa-lead)` | no VAL | File-state proof from task reports and `project-docs/docs/phases/phase-13/decisions.md` if created. |

## 11. Mid-Phase Review Cadence

- Cadence for this 3-task phase: one mid-phase review after Task 2 completes.
- Review file: `project-docs/docs/phases/phase-13/qa/review-1.md`.
- Review window: Tasks 1-2.
- Covered assertions: VAL-001 through VAL-008, with VAL-001 completed by Task 2
  and VAL-007 requiring both Task 1's no-activity label creation proof and Task
  2's issue-label activity projection proof.
- No mid-phase review after Task 1. Task 1 leaves VAL-001 partial and does not
  yet prove the issue-label side of VAL-007.
- Task 3 is reviewed at phase-end final review, not in a second mid-phase
  review, because it is the final closure smoke and requires the final
  deterministic gate.

Mid-phase review must compare completion-report command results with the exact
Task 1 and Task 2 gates in section 8. Missing, skipped, narrower, or silently
substituted required commands are blocking unless recorded as `BLOCKED` with a
reason in a blocker report.

## 12. Human Verification Scope

Items here require human dogfooding or product judgment and must be generated
into `project-docs/docs/phases/phase-13/qa/human-checklist.md` at phase end.
Each item is one observable from VAL-010.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-13a | Confirm from real CLI output that a label was created and listed with its name, colour, and description visible `(human)` | VAL-010 | Human note or pasted CLI excerpt from a real `.fid` workspace. |
| HUMAN-13b | Confirm from real CLI output that adding one label to an issue left the issue's existing labels in place `(human)` | VAL-010 | Human note or pasted before/after CLI excerpt. |
| HUMAN-13c | Confirm from real CLI output that removing one label from an issue left the remaining labels in place `(human)` | VAL-010 | Human note or pasted before/after CLI excerpt. |
| HUMAN-13d | Confirm from real CLI output that create -> claim -> transition -> complete -> unclaim ran end-to-end with a link and label on the issue `(human)` | VAL-010 | Human note or pasted command transcript. |

## 13. Phase-End Decision Rule

The phase can proceed to final sign-off only when:

- All task completion reports exist and record exact gate results.
- The deterministic phase-end gate passes or is explicitly BLOCKED with
  escalation.
- Final review records `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD` for every
  VAL-001..010.
- VAL-001..009 have cited deterministic evidence, not code-reading-only claims.
- VAL-010 has a generated human checklist containing only HUMAN-13a..HUMAN-13d
  or narrower human-only rows.
- No schema migration, projector change, or out-of-task-scope work is present
  without a recorded spec-backed decision.
- No unresolved blocking defects remain.

