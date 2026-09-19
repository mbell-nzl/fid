# Phase 23 - Task 5: Finding dispositions and resolution authority

## Objective

Implement the Finding lifecycle: Builder-recorded `addressed` assertions,
resolution and inapplicability only through a later qualifying independent
Review, and a derived Finding state.

## Validation Contract Coverage

- VAL-020: Findings are created only with their Review, atomically, and are
  immutable.
- VAL-021: `finding.address` requires an active mutating Claim and leaves the
  Finding unresolved and still blocking; the review-Claim case rejects
  `claim-mode-invalid` with `required_mode:"mutation"`. At this task the review
  Obligation is unsatisfied regardless, so the non-vacuous "still unsatisfied"
  proof is VAL-026's sequence at task 6; this row's final PASS needs both.
- VAL-022: Only a later qualifying independent Review resolves or marks a
  Finding inapplicable, and no Review resolves a Finding it created. Six cases,
  including a `blocked`-verdict resolving Review (rejects) and a resolution by
  the same independent Reviewer's later Review (succeeds, DEC-225).
- VAL-023: `Finding.state` is derived from the latest disposition.
- VAL-003: task 1's fence assertion for `finding.addressed` converges here
  without losing its parity subject.


## Binding Constraints

- `contracts.md` section 5 is the exact authority; `schema.md` section 2 fixes
  `finding_dispositions` including its discriminator CHECK and the partial
  unique index that caps terminal dispositions at one.
- `addressed` is an assertion of rework, never a resolution. It must not change
  blocking state, satisfy an Obligation, or close a Stop.
- Resolution arrives only through `resolves_finding_ids` on a later Review
  Record that binds the same Issue and revision, satisfies task 4 independence,
  and has verdict `pass` or `changes-required`. There is no `finding.resolve`
  leaf and none may be added.
- A Review may not resolve a Finding it created, nor one already resolved or
  inapplicable. Both rejections name the exact Finding.
- DEC-192: delete exactly the `"finding.addressed"` key and prove its
  live-apply and replay parity. Converge task 1's fence assertion for that key
  by substituting a still-fenced Event.
- Add the `finding.address` grant action exactly here.
- Do not change Obligation evaluation in this task; the
  `finding-unresolved` reason is delivered with the review Evidence path in
  task 6.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/review/` - Finding creation effects, dispositions, derived state.
- `internal/workflow/authority.go` and `internal/workflow/activation.go` - the
  one new grant action. The closed allowlist is `sortedUniqueCapabilities` in
  `activation.go`; Envelope grants reuse it through `sortedUniqueEnvelopeActions`
  in `envelope.go`, and `authority.go` holds no action list (R-08, DEC-167/173).
- `internal/journal/current_event_value_relations.go` - the **second** copy of
  that same closed allowlist (R-19, DEC-233). `validateCurrentValueRelations`'s
  `case "ActorProfile"` holds a `const actions` string with the identical 23
  names at `:63`, and nothing cross-checks the two. Add
  `finding.address`
  there in the same change that adds it to `sortedUniqueCapabilities`, keeping
  both copies sorted and identical. This is not on any Phase 23 family's strict
  path: it is reached through `activation.component-staged`, a live Phase 21/22
  family, so without it the first Activation that stages a profile carrying
  `finding.address`
  is rejected with `invalid closed enum` before this task can append a single
  Event, with its cause outside this task's scope.
- `internal/journal/project_effects.go` - the one fence key and its effects.
- `internal/journal/current_event_json.go` and
  `internal/journal/current_event_result_projection.go` - the typed-result
  registration for `finding.addressed` (R-16, DEC-231). Deleting the fence key
  routes this family onto the journal's strict result path, so both are required
  in this task or the first successful `finding.address` is rejected at append
  time. Add the `finding.addressed` entry in the `fields` grammar map
  (`finding:Finding`, from the `finding.address` result row in
  `public-interface.md` section 4) and the matching `currentExpectedResult`
  branch. **No new value shape and no `current_event_value_catalog.go` edit:**
  task 4 registered `Finding` and `FindingDisposition` and their fill helper,
  and this task reuses both unchanged. If either is missing or wrong, that is a
  task 4 defect - stop and escalate rather than editing shapes this task does
  not own. This family also needs **no** entry in the remaining three registries
  of the `execution-scope.md` section 6 inventory - cardinality, `*_json`
  columns, and update admission - checked deliberately rather than assumed
  (DEC-232): `finding_dispositions` carries no `*_json` column, its one
  Effect is an insert so `validateCurrentUpdate` is never reached, and
  `finding.address` disposes exactly one Finding, so the default maximum of one
  `finding_dispositions/insert` per Event is the correct bound and must not be
  widened here.
- `internal/cli/codes.go` - `finding-not-found`, `finding-already-disposed`.
- `cmd/fid/finding.go` - new leaf wiring for `address` and `list`.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-5-done.md`.

### Files To Read For Context

- `contracts.md` sections 4 and 5; `public-interface.md` sections 3, 4, 4.1, 5;
  `schema.md` section 2.
- `task-4-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Obligation evaluation changes, Evidence linking, budgets, and Human text
  beyond the minimum required to return a result.

## Acceptance Criteria

- [ ] Findings have no independent create or update path and appear only inside
  their Review Event. (VAL-020)
- [ ] Address succeeds only under an active mutating Claim, and afterwards the
  Finding is `addressed`, unresolved, and still blocking. (VAL-021)
- [ ] Self-resolution, non-independent resolution, and a second terminal
  disposition all reject; a valid independent resolution records its qualifying
  Review. (VAL-022)
- [ ] A supplied `state` rejects and the derived state follows the dispositions.
  (VAL-023)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestFindingAddressDoesNotResolveOrUnblock` fails against the absent
disposition writer.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-5-done.md` in the full format,
  including the fence-key deletion and its parity proof.
- Flag spec impact; only Tech Lead accepts a contract or schema change.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -timeout 30m ./cmd/fid/integration/...
```

Two gates run for minutes and are explicitly bounded (DEC-227): the `-timeout
30m` above replaces Go's implicit 10-minute per-binary default, which the tagged
suite already approaches at 560 s. Record in the completion report the wall time
`go test` reports for `cmd/fid/integration`, `cmd/fid`, and `internal/workflow`,
so a regression is attributed to the task that caused it. Budgets: tagged
integration under 900 s and no ordinary package over 600 s (VAL-079). A breach is
a bounded speed repair owned by this task when this task's own new journeys
caused it; otherwise stop and escalate. Never record a timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. No exception to the full task
gate. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination. An
interrupted or hung command is never `PASS`; preserve all failed attempts and
rerun the byte-identical command after an authorized causal repair. Report lint
is a Runner control after the report commit, not a Developer gate.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- This task belongs to QA Review window 2.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-5`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
