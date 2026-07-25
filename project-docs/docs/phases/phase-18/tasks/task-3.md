# Phase 18 - Task 3: Transition-path data-source swap (behavior-preserving)

## Objective

Repoint the issue transition path to resolve workflow states, categories,
allowed edges, and per-transition rules from the DB-backed registry (Task 2)
instead of the package-global Go map, preserving v0.1 behavior exactly —
including rich error messages and started-category claim enforcement.

## Validation Contract Coverage

- VAL-008: every v0.1 transition behaves identically, proven by an **exhaustive** ordered-pair state matrix across all three frozen workflows plus rule satisfied/blocked fixtures for every ruled edge and exact built-CLI output/error for the allowed, invalid-edge, blocked-rule, and cancel categories (representative-only evidence is insufficient).
- VAL-009: entering a started-category state requires the caller's live claim — missing, wrong-holder, and stale claims are rejected and a live caller claim is accepted, using DB-resolved categories.
- VAL-020: leaving a started-category state auto-releases the caller's live claim within the same event/projection operation (claim absent in the projection and after rebuild).

## Binding Constraints

- **Behavior-preservation is the acceptance bar.** The transition path must
  produce byte-identical outcomes and error messages to v0.1. Add
  characterization tests that lock `InvalidTransitionError` and
  `BlockedTransitionError` text and the `entersStartedCategory` /
  `leavesStartedCategory` decisions before changing the source.
- **Swap only the transition path in this task:** `internal/issues/transition.go`
  and `internal/issues/transition_rules.go` (edge legality, rule references, and
  category resolution used for claim gating). Load the registry from the open
  connection already available inside the `WithLock` block. Do **not** repoint
  the other consumers yet (maintenance, issue get/list category display, dev
  composites, ready) — those are Task 5. The Go map and the DB seed are
  behavior-identical (proven in Task 2), so the mixed state is consistent.
- **Rule evaluation stays app-layer.** `evaluateTransitionRule` continues to
  dispatch named Go rule functions; only the *source of the rule list and its
  params* changes from the Go map to the registry. Preserve the parameterised
  `all_acs_verified` roles exactly.
- **Preconditions unchanged.** The `issue.transitioned` effect keeps its
  `status`/`workflow_id` preconditions and the same activity-log effect; this
  task changes resolution, not the journal event shape.
- **Cancel semantics come from the seeded explicit edges** (Task 1), not a code
  special-case. If the old `toState == "cancelled"` special-case in
  `RulesForEdge` is retired, prove via characterization that cancel-from-any
  and the `cancelled→cancelled` rejection are unchanged.

## Scope

### Files To Create Or Modify

- `internal/issues/transition.go`
- `internal/issues/transition_rules.go`
- `internal/issues/transition_test.go`, `internal/issues/transition_rules_test.go` (characterization + parity)
- Possibly `internal/workflow/workflow.go` (retire/guard the code-level cancel special-case if superseded by seeded edges)

### Files To Read For Context

- `internal/workflow/registry.go` (Task 2 loader)
- `internal/issues/transition.go`, `internal/issues/transition_rules.go`
- `internal/claims/` (`RequireLiveIssueClaimHeld`, `MaybeBuildIssueReleaseEffect`)
- `project-docs/specs/02_TRD.md` §4.6

### Out Of Scope

- `issues` FK migration (Task 4).
- Non-transition consumers of `workflow.*` (maintenance, get/list, dev, ready) — Task 5.
- Any new command or workflow authoring.

## Acceptance Criteria

- [ ] The transition path resolves edges/categories/rules from the DB registry
- [ ] An exhaustive ordered-pair state matrix across all three frozen workflows
  matches v0.1: allowed edges succeed; every disallowed edge is rejected with the
  same `InvalidTransitionError` message (VAL-008)
- [ ] Every ruled edge is exercised in both satisfied and blocked states;
  blocked rules produce the same `BlockedTransitionError` naming the same unmet
  conditions, entity IDs, and missing roles as v0.1, with exact built-CLI output
  for the allowed/invalid/blocked/cancel categories (VAL-008)
- [ ] Cancel-from-any-non-cancelled succeeds; `cancelled→cancelled` rejected (VAL-008)
- [ ] Entering a started-category state requires the caller's live claim:
  missing, wrong-holder, and stale claims are rejected and a live caller claim is
  accepted, using DB-resolved categories (VAL-009)
- [ ] Leaving a started-category state auto-releases the caller's live claim in
  the same event/projection operation; the claim is absent in the projection and
  after rebuild (VAL-020)
- [ ] Full suite green; no change to the `issue.transitioned` journal event shape

## Remediation (Review 1 — VAL-020 rebuild evidence)

Review 1 returned `BLOCK` (TRG-002; resolved by DEC-068): VAL-020 was marked
`PASS` without proof that the released caller claim stays absent after a
projection rebuild. **VAL-020 is unchanged.** Redo this task within scope:

- Add a characterization test that (a) enters a started-category state holding a
  live caller claim, (b) transitions out of started, (c) rebuilds the **exact**
  released-claim workspace/projection, and (d) asserts the caller claim is still
  absent after the rebuild. The existing general rebuild test rebuilds an
  *entering-started* fixture and does **not** satisfy this — the released-claim
  fixture must be the one rebuilt.
- Correct the `task-3-done.md` VAL-020 evidence to cite the new test; remove the
  "existing rebuild check remains covered by the transition suite" claim.
- Rerun the full `## Exact Commands To Run` gate below and re-record results.
- A QA Reviewer re-review of the VAL-020 remediation must pass before Task 4.

## First Failing Test

A characterization test asserting a blocked `in_progress→dev_complete`
transition (missing developer AC verification) yields the exact v0.1
`BlockedTransitionError` text while the engine reads from the DB registry.

## Documentation Updates

- None expected in specs. If the code-level cancel special-case is retired,
  note it in the completion report (spec §14.11/§4.6 wording is synced by the
  Tech Lead at wrap).

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates.

Conditional command required by this task:

- `go test -tags integration -race ./...` — transition e2e parity is primary
  evidence for VAL-008/009.

## Notes For Developer

- Depends on Tasks 1–2. This is the behavior-critical task; if any transition
  outcome or message would change, stop and escalate rather than adjust the
  assertion.
- Load the registry once per transition inside the existing lock/connection;
  do not add a second DB open.
- Completion report at `project-docs/docs/phases/phase-18/tasks/task-3-done.md`.
