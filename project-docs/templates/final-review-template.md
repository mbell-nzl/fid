# Phase <N> Final Review - <Phase Title>

Reviewer: QA Lead
Date: <YYYY-MM-DD>
Review type: INITIAL / FRESH INDEPENDENT RE-REVIEW
Prior review: <path or none>

## Verdict

PASS / REPAIR / BLOCK

- `PASS`: no Must-fix finding remains; Improvements may remain.
- `REPAIR`: valid Must-fix work has bounded, decision-free ownership in an
  existing phase task. The phase does not advance.
- `BLOCK`: protected authority is required, repair exceeds existing phase task
  scope, or the bounded no-progress recovery budget is exhausted. A failed gate,
  incomplete in-scope work, or a valid Must-fix finding is not `BLOCK` by
  itself.

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
| Coverage/confidence | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -coverprofile=coverage.out ./...` | PASS/FAIL/N/A | <output excerpt or reason> |
| Coverage report | `go tool cover -func=coverage.out` | PASS/FAIL/N/A | <output excerpt or reason> |
| Integration | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 ./cmd/fid/integration/...` | PASS/FAIL/N/A | <output excerpt or reason> |
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

| Review | Outcome | Must-fix disposition | Fresh re-review evidence | Blocker candidate / existing ESC |
|---|---|---|---|---|
| review-<R>.md | PASS / REPAIR / BLOCK | <none, owning repair task, or protected boundary> | <resolving review and evidence or N/A> | <candidate summary, ESC-<phase>-<NNN>, or none> |

No mid-phase `REPAIR` is closed by implementation evidence alone. Cite the
fresh independent review that resolved each original Must-fix finding.

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

### Must-fix

- None.

### Improvement

- None.

## Prior Finding Resolution

Required for a fresh independent final re-review. Preserve every Must-fix
finding from the prior final review and resolve it only from newly inspected
evidence. Rerun the representative affected verification; do not modify or
append corrections to the prior review.

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| <finding ID> | task-<M> | RESOLVED / STILL OPEN / PROTECTED BLOCKER CANDIDATE | <command result, diff, commit, or artifact> |

## Repair Plan

Required for outcome `REPAIR`. Every Must-fix finding needs one row. The owning
task must already exist in the phase, and repair must not change accepted
product, architecture, data-model, public-interface, dependency, security, or
validation-contract decisions.

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| <finding ID> | task-<M> | VAL-NNN / no VAL | <bounded files/behavior> | <exact test, gate, report, or file-state proof> | <what the new QA Lead must inspect and rerun> |

## Blocker Candidates

- None.

Required for outcome `BLOCK`. Each candidate names a protected category,
evidence, why it is outside existing repair authority, the authority or external
action required, and any existing `ESC-<phase>-<NNN>`. QA roles do not allocate
`ESC` identities or edit escalation registers.

## Unresolved Risks And Carry-Forward Candidates

- None.

## Final Decision

```text
Must-fix: <count>  Improvement: <count>  Blocker candidates: <count>
VERDICT: PASS | REPAIR | BLOCK
```
