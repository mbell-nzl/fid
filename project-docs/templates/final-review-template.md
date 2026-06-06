# Phase <N> Final Review - <Phase Title>

Reviewer: QA Lead
Date: <YYYY-MM-DD>

## Verdict

PASS / PASS WITH IMPROVEMENTS / BLOCK

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Validation contract:
  `project-docs/docs/phases/phase-<N>/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-<N>/qa/qa-plan.md`
- Task specs and completion reports:
  `project-docs/docs/phases/phase-<N>/tasks/`
- Mid-phase reviews: `project-docs/docs/phases/phase-<N>/qa/review-*.md`
- Current decisions: `project-docs/docs/current/decisions_active.md`

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS/FAIL/N/A | <output excerpt or reason> |
| Lint/static analysis | `golangci-lint run ./...` | PASS/FAIL/N/A | <output excerpt or reason> |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS/FAIL/N/A | <output excerpt or reason> |
| Coverage report | `go tool cover -func=coverage.out` | PASS/FAIL/N/A | <output excerpt or reason> |
| Integration | `go test -tags integration -race ./...` | PASS/FAIL/N/A | <output excerpt or reason> |
| End-to-end | `<E2E_COMMAND>` | PASS/FAIL/N/A | <output excerpt or reason> |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS/FAIL/N/A | <output excerpt or reason> |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS/FAIL/N/A | <output excerpt or reason> |
| Full phase check | <full command list from `AGENTS.md`> | PASS/FAIL/N/A | <output excerpt or reason> |

## Validation Contract Outcomes

Every `VAL-NNN` from the validation contract must appear here. `PASS` without
evidence is invalid.

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS / FAIL / BLOCKED / CARRY-FORWARD | <QA Lead / QA Reviewer / Human> | <literal output, file reference, review reference, screenshot, API trace, or carry-forward rationale> |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| <PRD/TRD/design/quality ref> | MET / PARTIAL / NOT MET / N/A | <VAL ID, task, test, review, or source reference> |

## Deterministic Checklist Execution

Run every `(qa-lead)` item from the QA plan. PASS requires reproducible
evidence.

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS / FAIL / BLOCKED | <literal output, file-state proof, log, API trace, screenshot, or source reference> |

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: YES/NO.
- Invalid or unvalidated command, test scope, fixture, artifact, or human
  checklist step found: <none or list>.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC References |
|---|---|---|---|---|
| review-<R>.md | CONTINUE / BLOCK | YES/NO/N/A | <fixed/carry-forward/cancelled> | <ESC-NNN or none> |

## Documentation Review

- Impacted docs updated: YES/NO.
- Specs needing Tech Lead wrap-up sync: <none or list>.
- Completion reports present for all tasks: YES/NO.

## Human Checklist Readiness

Generated at
`project-docs/docs/phases/phase-<N>/qa/human-checklist.md`.

- Contains only `(human)` items: YES/NO.
- Uses one observable per row: YES/NO.
- Human-only `VAL-NNN` assertions mapped: YES/NO/N/A.
- Product/perception/dogfooding risks to call out: <none or list>.

## Coverage Caveat

A deterministic PASS proves only the asserted mechanics, data states, command
outputs, and structural conformance checks. It does not prove unasserted
appearance, workflow feel, product-owner judgment, or dogfooding quality; those
remain unconfirmed until the Human checklist records outcomes.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | READY / BLOCKED / N/A | <evidence> |
| Automated gates | READY / DEFERRED / BLOCKED / N/A | <evidence or recorded deferral> |
| Held-out scenarios | READY / DEFERRED / BLOCKED / N/A | <evidence or recorded deferral> |
| Release readiness | READY / DEFERRED / BLOCKED / N/A | <evidence or recorded deferral> |

## Findings

### Blocking

- None.

### Major

- None.

### Improvement

- None.

## Unresolved Risks And Carry-Forward Candidates

- None.
