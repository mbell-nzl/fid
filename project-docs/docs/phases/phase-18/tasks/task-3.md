# Phase 18 - Task 3: Transition-path data-source swap (behavior-preserving)

## Objective

Repoint the issue transition path to resolve workflow states, categories,
allowed edges, and per-transition rules from the DB-backed registry (Task 2)
instead of the package-global Go map, preserving v0.1 behavior exactly —
including rich error messages and started-category claim enforcement.

## Validation Contract Coverage

- VAL-008: every v0.1 transition behaves identically (allowed edges succeed; invalid edges rejected with the same `InvalidTransitionError`; blocked rules produce the same `BlockedTransitionError` naming the same unmet conditions/entities/roles; cancel-from-any-non-cancelled works).
- VAL-009: started-category claim enforcement parity (into-started requires a live claim; out-of-started auto-releases), using DB-resolved categories.

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
- [ ] Allowed edges succeed; invalid edges rejected with the same
  `InvalidTransitionError` message as v0.1 (VAL-008)
- [ ] Blocked rules produce the same `BlockedTransitionError` naming the same
  unmet conditions, entity IDs, and missing roles as v0.1 (VAL-008)
- [ ] Cancel-from-any-non-cancelled succeeds; `cancelled→cancelled` rejected (VAL-008)
- [ ] Into-started transition requires a live claim; out-of-started auto-releases,
  using DB-resolved categories (VAL-009)
- [ ] Full suite green; no change to the `issue.transitioned` journal event shape

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
