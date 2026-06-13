# Phase 4 QA Plan - Issue Authoring And Discovery

> Phase: 4 - v0.1b Issue authoring and discovery
> Refreshed: 2026-06-09
> QA Lead: Codex

## 1. Scope And Sources

Phase 4 verifies the ordered authoring/discovery slice: sequence scheduling
backstop, issue authoring domain service, `fid issue create`, centralised
`fid issue get --section`, and `fid issue list`.

Sources inspected:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/01_PRD.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/specs/05_DESIGN_PLAN.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/docs/phases/phase-4/qa/validation-contract.md`
- `project-docs/docs/phases/phase-4/qa/review-3.md`
- `project-docs/docs/phases/phase-4/tasks/task-1.md` through `task-5.md`

## 2. Validation Contract Refinement Check

The Phase 4 validation contract is acceptable as the verification target. No
Tech Lead bounce is required before implementation.

- All in-scope Phase 4 requirements are mapped to VAL-001..VAL-010: YES.
- Every VAL has a validator role, evidence method, and owning task/window: YES.
- VAL-001..VAL-009 are deterministic and observable: YES.
- VAL-010 is human-only and must fan out into one dogfood checklist row per
  observable command/error flow: YES.
- Conditional/binary assertions require both passing and rejecting states:
  YES, especially VAL-001, VAL-002, VAL-004, VAL-007, VAL-008, VAL-009.
- Parity/reference assertions require deterministic conformance evidence:
  YES, especially VAL-005, VAL-006, VAL-009.

QA-plan refinements within approved scope:

- Filter evidence for VAL-008 must show both included and excluded issues for
  every filter.
- Unknown label, unknown section, and unknown format checks must prove rejection
  and supported-value messaging, not just non-zero exit.
- ESC-002 is resolved-pending, not a contract bounce: VAL-003 still requires
  due-date read-back, and Task 4 closes that gap without schema or journal
  changes by proving `due_date` through `internal/issues.Get` and
  `fid issue get --format json`.
- VAL-010 must remain in the human checklist only; QA Lead must not pass it by
  code reading.

## 3. Ordered Task Windows

| Task | Focus | Covered VAL IDs | Required Evidence |
|---|---|---|---|
| 1 | Sequence scheduling backstop | VAL-001 | Raw SQL trigger accept/reject tests; journal projection/rebuild proof; task gate output. |
| 2 | Issue authoring domain service | VAL-002, VAL-003, VAL-004, VAL-005 | `internal/issues` and journal tests for hierarchy, fields, metadata, due-date persistence, one-event/rebuild behavior; task gate output. Due-date public read-back remains open for Task 4 per ESC-002. |
| 3 | `fid issue create` CLI | VAL-001, VAL-002, VAL-003, VAL-004, VAL-005 | CLI smoke tests, invalid create errors, journal line count/rebuild proof, `go build`; task gate output. Due-date public read-back remains open for Task 4 per ESC-002. |
| 4 | Centralised detail read model, sections, and due-date read-back closeout | VAL-003, VAL-006, VAL-007 | `internal/issues.Get` populated/null `due_date` read-back tests, `fid issue get --format json` populated/null `due_date` read-back tests, `internal/issues` populated/empty detail tests, `cmd/fid` section tests, adapter-boundary scan; `go build`; task gate output. |
| 5 | `fid issue list` filters/order/formats | VAL-008, VAL-009 | Include/exclude filter tests, stable ordering tests, JSON/TOON/YAML/markdown format tests, unknown format rejection; `go build`; task gate output. |
| Phase end | Full authoring/discovery loop | VAL-001..VAL-010 | Full phase gate, deterministic checklist evidence, final review outcomes, human checklist for VAL-010. |

## 4. Command Evidence Plan

Every task completion report must record the exact required commands from that
task spec as `PASS`, or `BLOCKED` with reason. Missing, silently skipped, or
narrower substituted commands fail the gate.

| Window | Required Commands |
|---|---|
| Tasks 1-2 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`. |
| Tasks 3-5 | Task 1-2 commands plus `go build -o /tmp/fid ./cmd/fid`. |
| Phase end | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race ./...`; `go test -bench . -benchmem ./...`; `go build -o /tmp/fid ./cmd/fid`. |

Coverage evidence must include the `go tool cover -func=coverage.out` output
and confirm COV-001 domain package thresholds and COV-002 `cmd/fid` threshold
where applicable.

## 5. Feasibility Validation Evidence

| Referenced Step | Validation Method | Result | Evidence |
|---|---|---|---|
| Go toolchain | Version probe | PASS | `go version` returned `go1.26.2 linux/amd64`. |
| Lint toolchain | Version probe | PASS | `golangci-lint version` returned `2.12.2`. |
| `gofmt -l cmd internal` | Help output and path existence | PASS | `gofmt -h` documents `-l`; `cmd/` and `internal/` exist. |
| `golangci-lint run ./...` | Help output | PASS | `golangci-lint run --help` documents the `run` command. |
| `go test -race -coverprofile=coverage.out ./...` | Go test help and existing tests | PASS | `go help testflag` documents `-race` and `-coverprofile`; repo has `*_test.go` under `cmd/` and `internal/`. |
| `go tool cover -func=coverage.out` | Tool help | PASS | `go tool cover -h` documents `-func`. It requires the prior coverage command to create `coverage.out`. |
| `go test -tags integration -race ./...` | Go test help and project gate docs | PASS | `go help testflag` supports test flags; `AGENTS.md` and `08_QUALITY_PLAN.md` require the phase-end integration command. |
| `go test -bench . -benchmem ./...` | Go test help and benchmark existence | PASS | `go help testflag` documents `-bench` and `-benchmem`; `internal/journal/rebuild_bench_test.go` contains benchmarks. |
| `go build -o /tmp/fid ./cmd/fid` | Go build help and writable output path | PASS | `go help build` documents `-o`; `/tmp` and `cmd/fid/` exist. |
| Task completion artifacts | File-state check | PASS | `project-docs/docs/phases/phase-4/tasks/` exists and task specs require `task-N-completion.md`. |
| QA artifacts | File-state check | PASS | `project-docs/docs/phases/phase-4/qa/` exists; templates for QA plan, final review, mid-phase review, and human checklist exist. |
| Human dogfood checklist | Scope review | PASS | VAL-010 is executable from the built `/tmp/fid` binary after Tasks 1-5 and is limited to human judgment. |

No placeholder required command was found in the Phase 4 task specs. End-to-end
scenario validation remains `N/A` per `AGENTS.md` until a harness is accepted.

## 6. Deterministic QA Checklist

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify Task 1 completion report includes exact task gates and VAL-001 trigger/projection evidence `(qa-lead)` | VAL-001 | Completion report excerpts plus referenced test output. |
| QA-02 | Verify Task 2 completion report includes exact task gates and VAL-002..VAL-005 domain/journal evidence `(qa-lead)` | VAL-002..VAL-005 | Completion report excerpts plus referenced test output. |
| QA-03 | Verify Task 3 completion report includes exact task gates, build, CLI create smoke, invalid error, one-event, and rebuild evidence `(qa-lead)` | VAL-001..VAL-005 | Completion report excerpts, command output, and journal/rebuild proof. |
| QA-04 | Verify Task 4 completion report includes exact task gates, build, due-date read-back, section tests, and adapter-boundary scan evidence `(qa-lead)` | VAL-003, VAL-006, VAL-007 | Completion report excerpts showing `internal/issues.Get` populated/null `due_date`, `fid issue get --format json` populated/null `due_date`, section outputs/rejections, command output, and scan result. |
| QA-05 | Verify Task 5 completion report includes exact task gates, build, filter include/exclude, ordering, format parity, and unknown-format evidence `(qa-lead)` | VAL-008, VAL-009 | Completion report excerpts and command output. |
| QA-06 | Run the full phase-end command gate `(qa-lead)` | VAL-001..VAL-009 | Literal command output in `final-review.md`. |
| QA-07 | Record a final outcome for every VAL-001..VAL-010 `(qa-lead)` | VAL-001..VAL-010 | `PASS`, `FAIL`, `BLOCKED`, or `CARRY-FORWARD` with cited evidence/disposition. |

## 7. Mid-Phase Review Cadence

Default cadence is QA Reviewer after every 3 completed tasks. Phase 4 therefore
requires the recommended review after Task 3 at:

- `project-docs/docs/phases/phase-4/qa/review-3.md`

The Task 3 review covers VAL-001..VAL-005, completion-report command coverage
for Tasks 1-3, early Task 4/5 scope checks, and ESC-002. ESC-002 is
resolved-pending after TL disposition on 2026-06-09 and closes when Task 4 lands
the due-date read-back evidence required above. No additional mid-phase review
is required solely for ESC-002; if Task 4 or Task 5 materially changes public
output shape beyond the task specs, QA Lead may add an extra review after Task 4
before list work continues.

## 8. Human Verification Scope

The phase-end human checklist must contain only human-judgment checks. Fan out
VAL-010 into one observable per row:

- HUMAN-01a: create an epic/task/subtask loop feels usable from the built
  binary `(human)`.
- HUMAN-01b: `fid issue get` and supported `--section` outputs are useful for
  real inspection `(human)`.
- HUMAN-01c: `fid issue list` filters/order/formats are useful for scanning a
  working set `(human)`.
- HUMAN-01d: one invalid hierarchy or scheduling error is clear and actionable
  `(human)`.

## 9. Phase-End Decision Rule

Phase 4 can proceed to final approval only when all five tasks are complete in
order, required task gates are evidenced, the phase-end gate is rerun, every
VAL-NNN has a final outcome, deterministic QA items have reproducible evidence,
and VAL-010 has a human checklist outcome or accepted carry-forward.
