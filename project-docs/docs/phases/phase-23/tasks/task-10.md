# Phase 23 - Task 10: Provenanced Human Outcomes

## Objective

Record authorized Human judgment with active Kit Activation provenance, allow it
to satisfy only a matching `human-outcome` Obligation, and never present it as
deterministic proof.

## Validation Contract Coverage

- VAL-043: Recording requires a human Actor type in a finalized Activation
  profile with the capability.
- VAL-044: A Human Outcome satisfies only a matching class with result
  `approved`.
- VAL-045: Human Outcome Evidence is never serialized as deterministic proof.
- VAL-046: A `human-outcome-required` Stop is raised exactly when the Human
  Obligation is the only unsatisfied one.
- VAL-003: task 1's fence assertion for `human-outcome.recorded` converges
  here without losing its parity subject.


## Binding Constraints

- `contracts.md` section 10 and `public-interface.md` sections 3 and 4 are the
  exact authority; `schema.md` section 2 fixes `human_outcomes`.
- Authority is established by active-version Activation provenance and the
  declared capability, never by a role string or a free-form `--by human` flag.
  A non-human actor type, an unfinalized Activation, or a missing capability all
  reject `authority-missing` with no Event.
- Standalone Issue inventory attribution can never reach this leaf
  (PRD-AUTH-007). Prove the negative case explicitly.
- A Human Outcome may satisfy only a `human-outcome` Obligation whose
  `outcome_class` equals its class, and only with result `approved`. Any other
  link rejects `outcome-not-permitted` naming the Obligation type.
- The Evidence record for a Human Outcome retains exactly `producer:"human"`,
  `trust_source:"activated-human"`, and
  `artifact_availability:"not-applicable"`. No field, reason, or text may
  present it as deterministic or fid-executed proof (PRD-HUM-001).
- Raise a `human-outcome-required` Stop exactly when the Human Obligation is the
  only unsatisfied Obligation. Routine status raises no Attention (PRD-HUM-002).
- DEC-192: delete exactly the `"human-outcome.recorded"` key and prove its
  live-apply and replay parity; converge task 1's fence assertion for that key.
- Add the `human-outcome.record` grant action exactly here. Replace the seeded
  Outcome fixtures used by tasks 7 and 9 with real recorded Outcomes and record
  that convergence with provenance.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/review/` - Human Outcome rules, provenance, effects.
- `internal/evidence/evidence.go`, `policy.go` - the Human Outcome evidence type.
- `internal/workflow/authority.go`, `evaluation.go`, `resolve.go`.
- `internal/journal/project_effects.go`, `current_event_catalog.go`.
- `internal/cli/codes.go` - `outcome-not-permitted`, `outcome-not-found`,
  `human-outcome-required`.
- `cmd/fid/human_outcome.go` - new leaf wiring for `record` and `list`.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change, including the task 7 and task 9 seeded fixtures.
- `project-docs/docs/phases/phase-23/tasks/task-10-done.md`.

### Files To Read For Context

- `contracts.md` sections 7, 9, 10, 13; `public-interface.md` sections 2, 3, 4,
  4.1, 5; `schema.md` sections 1, 2, 6.
- `task-9-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Attention resolution through a Human Outcome, which is task 11.
- The `human_outcome_resolution` capability limit, which is task 17.

## Acceptance Criteria

- [ ] Agent actor type, unfinalized Activation, and missing capability each
  reject `authority-missing` with no Event; the valid case records with its
  provenance bound. (VAL-043)
- [ ] Only the matching class with `approved` satisfies; every other link
  attempt rejects `outcome-not-permitted`. (VAL-044)
- [ ] The Evidence record and evaluation carry exactly the declared triple and
  no field or text claims deterministic proof. (VAL-045)
- [ ] The `human-outcome-required` Stop appears only when the Human Obligation
  is the single unsatisfied one, and never for routine status. (VAL-046)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row, including the seeded task 7 and task 9
  fixtures, records its provenance and executing command.

## First Failing Test

`TestHumanOutcomeRequiresFinalizedActivationProvenance` fails against the
absent Human Outcome writer.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-10-done.md` in the full format,
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
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
```

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

- This task belongs to QA Review window 4.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-10`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
