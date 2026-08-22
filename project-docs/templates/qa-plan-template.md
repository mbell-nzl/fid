# Phase <N> QA Plan - <Phase Title>

> Phase: <N> - <Phase Title>
> Refreshed: <YYYY-MM-DD>
> QA Lead: <name or tool>
>
> Create this file at
> `project-docs/docs/phases/phase-<N>/qa/qa-plan.md`.

## 1. Scope And Source Documents

Summarize the phase goal and list the source artifacts used to build this QA
plan.

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/docs/phases/phase-<N>/qa/validation-contract.md`
- `project-docs/docs/phases/phase-<N>/tasks/task-<M>.md`

## 2. Validation Contract Refinement Check

Confirm the validation contract is ready to be the verification target before
task execution begins. Gaps are QA-plan defects.

- All in-scope requirements referenced by at least one `VAL-NNN`: YES/NO.
- Every `VAL-NNN` has one validator role: YES/NO.
- Every `VAL-NNN` has a concrete evidence method: YES/NO.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES/NO.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: YES/NO.
- Conditional or binary `VAL-NNN` assertions require both qualifying and
  non-qualifying states: YES/NO/N/A.
- Parity or "matches reference" `VAL-NNN` assertions have deterministic
  conformance proof when available: YES/NO/N/A.
- Revisions requested from Tech Lead: <none or list>.

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task <M> | <task focus> | VAL-NNN | <low/medium/high> | <coverage or review note> |

## 4. Required Evidence Strategy

Each row states how QA expects evidence to prove the validation contract. Rows
that do not map to a validation assertion use `no VAL`.

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | <test output, command output, file-state proof, screenshot, API trace, or human observation> | VAL-NNN | <Developer / QA Reviewer / QA Lead / Human> | <task/window/final review> |

## 5. Deterministic Gates

Use concrete project commands after stack selection. If any required command is
still a placeholder when execution is requested, mark the QA plan BLOCKED.

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Every task |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | Task-level gate |
| Coverage report | `go tool cover -func=coverage.out` | PASS and threshold evidence | Task-level gate |
| Integration | `go test -tags integration -race ./cmd/fid/integration/...` | PASS or N/A | Phase-end when applicable |
| End-to-end | `<E2E_COMMAND>` | PASS or N/A | Surface/user-journey phases |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | Phase-end or performance-sensitive tasks |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Phase-end or release candidate |
| Full phase check | `gofmt -l cmd internal` + `golangci-lint run ./...` + coverage gate + integration + benchmarks + build | PASS | Phase-end gate |

## 6. Command And Feasibility Validation

Validate every referenced command, test scope, fixture setup, artifact
assumption, and human-checklist feasibility using help output, dry-run mode,
project documentation, existing file-state checks, or equivalent proof. Invalid
or unvalidated steps block the QA plan.

In a normal-mode phase, results may carry forward from the prior phase's QA
plan for commands that are unchanged; cite the prior plan as evidence and
revalidate any changed or new command. Heavy-mode phases validate fresh with no
carry-forward (`09_AI_WORKFLOW.md` section 4).

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | <help/dry-run/docs/script inspection> | PASS/FAIL/BLOCKED | <output excerpt or file reference> |

## 7. Deterministic Checklist

Items here are executed or verified by QA Lead and must end with `(qa-lead)`.
Every `QA-NN` cites covered `VAL-NNN` IDs or `no VAL`. PASS without evidence is
invalid.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | <deterministic check> `(qa-lead)` | VAL-NNN | <command output, file-state proof, log, API trace, screenshot, or source reference> |

## 8. Mid-Phase Review Cadence

- Default cadence: QA Reviewer after every 3 completed tasks.
- Exceptions: <none or named checkpoints>.
- Review template:
  `project-docs/templates/mid-phase-review-template.md`.
- Each mid-phase review verifies in-window `VAL-NNN` assertions and records
  gate evidence or why a gate was not run.
- Each review records exactly one outcome: `PASS`, `REPAIR`, or `BLOCK`.
- `REPAIR` is required when a valid Must-fix finding has bounded,
  decision-free ownership inside an existing reviewed task. The review names
  the owning task, affected assertion, exact repair boundary, required evidence,
  and fresh re-review scope. The window does not advance.
- After repair, the Phase Runner dispatches a Developer continuation for the
  owning task and then a fresh independent QA Reviewer for the same window. The
  new review uses the next unused review path, cites the prior review, and
  verifies each original finding; the original reviewer does not append a
  correction or resolve its own finding.
- `BLOCK` is reserved for a protected authority boundary, work outside the
  reviewed task/window scope, or exhaustion of the bounded no-progress recovery
  budget. QA records a blocker candidate but does not allocate an `ESC` ID or
  edit the escalation register.

## 9. Human Verification Scope

Items here require human judgment and must end with `(human)`. Include only
perception, product-owner judgment, dogfooding, tone, visual hierarchy, or other
human-only checks. Use one observable per row. If a parent `VAL-NNN` spans
several observables, use sub-IDs that cite the same parent assertion.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-01a | <single observable manual check> `(human)` | VAL-NNN or no VAL | <observation, decision, screenshot, or notes> |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Required deterministic gates pass with observed termination evidence.
- Every `VAL-NNN` has planned evidence and a final-review outcome slot.
- Every `(qa-lead)` item has reproducible evidence.
- Human checks are human-judgment-only and mapped to the human checklist.
- Every mid-phase review is `PASS`; each earlier `REPAIR` has a fresh
  independent `PASS` re-review that resolves the original Must-fix findings.
- No unresolved Must-fix finding or blocker candidate remains.
- Required documentation updates are named.

An in-scope failed gate or incomplete evidence triggers repair; it is not made
acceptable by recording `BLOCKED`. A protected missing or contradictory
contract prevents final review and is emitted as a blocker candidate for
TL-Lite or the full Tech Lead, which alone may allocate an `ESC` identity.
