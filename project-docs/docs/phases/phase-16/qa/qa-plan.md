# Phase 16 QA Plan - Activity Timeline + Read-Only Escape Hatches

> Phase: 16 - Activity Timeline + Read-Only Escape Hatches
> Refreshed: 2026-07-05
> QA Lead: phase-runner assigned
>
> Location: `project-docs/docs/phases/phase-16/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 16 is Heavy Mode. It verifies the activity timeline, CF-011 projection
reconciliation, `fid sql`, `fid journal inspect`, and CF-016 gate history read
surface.

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md` section 4
- `project-docs/docs/phases/phase-16/qa/validation-contract.md`
- `project-docs/docs/phases/phase-16/tasks/task-1.md` through `task-6.md`

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: YES.
- Every `VAL-NNN` has one validator role: YES.
- Every `VAL-NNN` has a concrete evidence method: YES.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: YES. VAL-013 fans out into HUMAN-16a..c.
- Conditional or binary `VAL-NNN` assertions require both qualifying and
  non-qualifying states: YES.
- Parity or "matches reference" `VAL-NNN` assertions have deterministic
  conformance proof when available: YES. VAL-007 rebuild equality and VAL-010
  `HashEvent` parity are explicit.
- Revisions requested from Tech Lead: none after DEC-040..042.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | `fid activity list` read package and command | VAL-001..004 | Medium | Filter and format surface; verify shape before CF-011 backfill. |
| Task 2 | CF-011 create/ac/complete activity rows | VAL-005 | Medium | Decision-backed behavior change; confirm one event and zero-on-reject. |
| Task 3 | CF-011 concerns/unclaim/link add + e2e rebuild | VAL-006/007 | High | Integration/e2e and rebuild determinism; integration gate required. |
| Task 4 | `fid sql` read-only escape hatch | VAL-008/009 | High | TRD-SEC-002 security evidence; integration gate required. |
| Task 5 | `fid journal inspect` | VAL-010/011 | Medium | Canonical hash parity and revert verdict mapping. |
| Task 6 | `issue get --section gates` | VAL-012 | Medium | Public read surface; raw history distinct from completion snapshot. |

## 4. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Unit/CLI logs for activity list shape, ordering, filters, and formats | VAL-001..004 | Developer / QA Lead | Task 1 / Review 1 |
| EVID-02 | Journal-line and `activity_log` deltas for create/ac/complete success and rejection paths | VAL-005 | Developer / QA Lead | Task 2 / Review 1 |
| EVID-03 | Journal-line and `activity_log` deltas for concern/unclaim/link add, including link-add parity | VAL-006 | Developer / QA Lead | Task 3 / Review 2 |
| EVID-04 | Full lifecycle activity timeline and delete-and-rebuild equality output | VAL-007 | Developer / QA Lead | Task 3 / Review 2 |
| EVID-05 | `fid sql` accepted read statements and all-format rendering | VAL-008 | Developer / QA Lead | Task 4 / Review 2 |
| EVID-06 | TRD-SEC-002 integration evidence for rejected mutations and unchanged DB state | VAL-009 | Developer / QA Lead | Task 4 / Review 2 |
| EVID-07 | Journal inspect known/missing event output plus `HashEvent` parity | VAL-010 | Developer / QA Lead | Task 5 / Final review |
| EVID-08 | Journal inspect revert verdict examples for inverse and not-eligible correction | VAL-011 | Developer / QA Lead | Task 5 / Final review |
| EVID-09 | `--section gates` in-progress, finished, empty, all-format output | VAL-012 | Developer / QA Lead | Task 6 / Final review |
| EVID-10 | Human dogfood notes for timeline, SQL read/reject, and journal inspect | VAL-013 | Human | Human checklist |

## 5. Deterministic Gates

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Every task |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | Every task |
| Coverage report | `go tool cover -func=coverage.out` | PASS and package percentages recorded | Every task |
| Integration | `go test -tags integration -race ./...` | PASS or BLOCKED with reason | Required for Tasks 3 and 4; required for any task adding integration-tagged tests; phase end |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | Phase end |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Every task that touches `cmd/fid`; phase end |
| Full phase check | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go test -bench . -benchmem ./...`; `go build -o /tmp/fid ./cmd/fid` | PASS | Final review |

## 6. Command And Feasibility Validation

Heavy Mode requires fresh validation; do not carry forward prior command
validation without checking the current task specs and this plan.

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| Task command lists | Reviewed `task-1.md` through `task-6.md`; each lists the full task gate and conditional integration/build commands. | PASS | This QA plan; task specs. |
| Required commands are concrete | Reviewed `AGENTS.md` command section; no placeholder command is referenced by Phase 16 task specs. | PASS | `AGENTS.md`; task specs. |
| Integration gate discipline | Tasks 3 and 4 explicitly require `go test -tags integration -race ./...`; other tasks require it when integration-tagged tests are added. | PASS | Task 3, Task 4, this plan section 5. |
| Human checklist feasibility | VAL-013 has three CLI-visible observables and can fan out into HUMAN-16a..c. | PASS | validation-contract.md VAL-013; this plan section 9. |
| Heavy-mode cadence | Six tasks and shared public CLI/security surfaces justify reviews after Tasks 2 and 4. | PASS | DEC-040; this plan section 8. |

## 7. Deterministic Checklist

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify every completion report records all required task commands as PASS or BLOCKED `(qa-lead)` | no VAL | Completion reports and command output excerpts. |
| QA-02 | Verify task docs cite their covered VAL IDs and required documentation updates `(qa-lead)` | no VAL | Task spec/report review. |
| QA-03 | Verify CF-011 behavior changes preserve one event per writer and zero append on rejection `(qa-lead)` | VAL-005/006 | Test logs, journal-line deltas, `activity_log` row deltas. |
| QA-04 | Verify VAL-007 rebuild equality uses a deterministic ordered comparison `(qa-lead)` | VAL-007 | Integration/e2e output and comparison description. |
| QA-05 | Verify `fid sql` defense in depth: classifier rejection plus read-only connection `(qa-lead)` | VAL-009 | Integration logs and completion-report security note. |
| QA-06 | Verify journal inspect hash equals `journal.HashEvent` and verdict mapping covers current event types `(qa-lead)` | VAL-010/011 | Unit/CLI logs and event-type table. |
| QA-07 | Verify `--section gates` reads raw `gate_runs`, not completion snapshots `(qa-lead)` | VAL-012 | CLI/domain test output. |
| QA-08 | Verify docs named by task specs were updated or explicitly justified as unchanged `(qa-lead)` | no VAL | Diff review and completion reports. |

## 8. Mid-Phase Review Cadence

- Review 1 after Tasks 1-2: activity read surface plus first CF-011 behavior
  change window.
- Review 2 after Tasks 3-4: CF-011 e2e/rebuild and `fid sql` security window.
- Final review after Tasks 5-6 and human checklist preparation.
- Review template: `project-docs/templates/mid-phase-review-template.md`.

## 9. Human Verification Scope

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-16a | Read an activity timeline from CLI output and confirm it reflects recent writes `(human)` | VAL-013 | Human note from a real `.fid` workspace. |
| HUMAN-16b | Run a read-only `fid sql` query and a mutation attempt; confirm the read works and mutation is rejected `(human)` | VAL-013 | Human note with command references. |
| HUMAN-16c | Inspect a journal event by ID and confirm the output explains the event without reading implementation `(human)` | VAL-013 | Human note with command reference. |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Required deterministic gates pass or are recorded BLOCKED with escalation.
- Every `VAL-NNN` has planned evidence and a final-review outcome slot.
- Every `(qa-lead)` item has reproducible evidence.
- Human checks are human-judgment-only and mapped to the human checklist.
- No unresolved blocking defects remain.
- Required documentation updates are named in completion reports.
