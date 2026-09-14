# Phase 22 QA Lead Handoff

> Planning handoff, not a QA plan or evidence receipt.
> Tech Lead: 2026-09-05. Heavy mode and product amendments await Human approval.

## Inputs And Required Output

Read `../contracts.md`, `../public-interface.md`, `../schema.md`,
`../execution-scope.md`, `../v1-convergence-inventory.md`, the 18 task specs,
and `validation-contract.md`. Author `qa-plan.md` from the QA plan template;
propose assertion/command refinements to Tech Lead before marking it Ready.
Tech Lead has run no gates; there are no current-session PASS receipts to reuse.

Map all VAL-001..080. QA Lead is the one final validator for 77 deterministic
rows; Human owns VAL-076..078. Mid-phase QA Reviewer evidence does not change
validator ownership. Fresh feasibility is required because this is Heavy mode.
Test commands for future behavior must name the producing task, and existing
environment/command assumptions must be verified now. Any added Developer
command goes into that task's Exact Commands section before dispatch.

## Independent Contract Review

The accessible kit source is
`/home/mark/projects/agentic-coding-project-starter-kit`, planning HEAD
`c5229617220f97ff9228df261e24718b1d1d6d89`.
Read only normal specification/port material, never excluded scenarios.
Primary sources: `project-docs/specs/trd/02_INTERFACES.md` sections 3.2-3.7,
`project-docs/specs/FID_COMPATIBILITY.md`, and decision 0005 (linked recovery).
Its compatibility matrix says no production candidate is qualified; preserve
that distinction. Check current version before using it as final evidence.

Record independent consistency findings in `consistency-review.md` before
Ready. In particular compare closed terminated/live/unknown liveness, exact
Human force grant binding, budget measurement evidence, exhaustion/Attention,
immutable history, resource agreement, and validate-only import semantics.
The fid CLI may differ in spelling; semantic disagreement needs a decision.
This session's TL comparison is input to the review, not independent evidence.

## Cadence And Early Journey

Six three-task reviews follow tasks 3,6,9,12,15,18. Task 6 must already have a
real subprocess recovery journey and reusable dual-adapter case structure.
The independent task-12 checkpoint writes `checkpoint-after-task-12.md` and
must PASS before task 13. It reruns recovery, force/unknown liveness, budget
atomicity, Stop disposition, activation directions and freshness; it also
validates import fault injection and environment feasibility before import
implementation begins. Ordinary deterministic checks are not held-out scenarios.

## Required Gate And Evidence Strategy

All commands run from the phase repository/worktree root. Final QA independently
runs the full phase gate:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

Task 17 lists the unchanged reference-information commands and three fresh
100k benchmark processes; reproduce all three for final VAL-071. Every sample
must self-assert parity and <30 seconds on the recorded reference instance.
No best-of, median, extra tolerance, stale receipt, or narrowed scope. If the
reference is unavailable, record BLOCKED with the actual mismatch.

Final QA also performs requirement compliance, security review, performance
spot-check, documentation/spec comparison, retained-test provenance audit,
deterministic checklist, and Human checklist. Pay attention to archive expansion
limits, path containment, force authority, unknown liveness, import publication,
Event effect validation, and old executor fencing. Use real process fault
injection and controlled synchronization; do not rely on sleeps to win races.

All parity comparisons enumerate projections and legitimate environmental
normalizations. Check every recovery Event family, committed denial, missing
Receipt, unavailable Artifact, null budget and open Attention. Require a
non-vacuous operation/error/live-frame inventory; do not rerun mutations to
hide a post-rebuild read difference.

## Human And Release Scope

Prepare three independently markable Human checks: determine the next recovery
action from a fresh session; distinguish pre/post-commit activation recovery;
recognize imported identity and remaining intervention. Task 18 produces
disposable fixtures and invocation instructions; Human records judgment.

No accepted held-out harness exists in this repository's current maturity
contract. Held-out/Mode B external scenario execution is N/A for this phase;
do not search for or read evaluator material. The three ordinary Human checks
remain required. Release threat model, external SAST selection, vulnerability,
license and secret release gates, migration rehearsal, and full twelve-capability
qualification retain Phase 25 ownership. This phase authorizes no live adoption.

## Feasibility Questions To Resolve Before Ready

- Does the archive limit/fault-injection strategy exercise actual expanded bytes
  without excessive machine memory or a new dependency?
- Can a bounded trusted observer and explicit force fact satisfy the kit's
  exact authority semantics on the recorded environment?
- Do schema-15 intermediate tests and schema-16 final import each execute their
  own documented paths without moving a gate later?
- Is every task's closed request/result/Event contract complete and every owned
  fixture executed by its listed gates?

These are QA design/feasibility work, not permission questions to the Human.
Escalate only an actual unresolved protected contract choice or capability gap.
