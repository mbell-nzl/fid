# Phase 23 - Task 10: Provenanced Human Outcomes

## Objective

Record authorized Human judgment with active Kit Activation provenance, allow it
to satisfy only a matching `human-outcome` Obligation, and never present it as
deterministic proof.

## Validation Contract Coverage

- VAL-043: Recording requires a human Actor type in a finalized Kit Activation
  profile with the capability, and standalone Issue-inventory attributed
  authority rejects (PRD-AUTH-007).
- VAL-044: A Human Outcome satisfies only a matching class with result
  `approved`.
- VAL-045: the Human Outcome Evidence record and its evaluation carry
  `producer:"human"`, `trust_source:"activated-human"`, and
  `artifact_availability:"not-applicable"`. The Human-text half of
  `contracts.md` section 14 belongs to task 14 (VAL-058/VAL-060); do not pin text
  task 14 replaces.
- VAL-046: A `human-outcome-required` Stop is raised by each of its four named
  Events exactly when the Human Obligation becomes the only unsatisfied one, at
  most once while unresolved, and never for routine status (DEC-224).
- VAL-009: the completing half - a guarded action is permitted only once the
  gate-evidence, review, and human-outcome Obligations are all satisfied. This
  is the first task where that state is reachable.
- VAL-010: re-verified here with linked review and Human Outcome evidence
  present, not only with evidence absent.
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
- The `artifact.availability-observed` raising sequence has a proven fixture
  technique (I-24, reproduced live at QA readiness): remove the retained
  content-addressed blob under `.fid/artifacts/<digest>` and observe `missing`,
  then restore it and observe `available`. `internal/evidence/artifact.go`
  derives availability by digest-verifying that blob and writes one
  `referenced_artifacts/update` effect with a precondition on the prior value,
  so `missing -> available` is exactly one qualifying Event and the Obligation
  transitions unsatisfied -> satisfied inside it.
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
- `internal/workflow/authority.go`, `activation.go`, `evaluation.go`,
  `resolve.go` - the new grant action belongs to the closed
  `sortedUniqueCapabilities` allowlist in `activation.go` (R-08, DEC-167/173).
- `internal/workflow/` readiness-Stop raising shared by the four producers, and
  its call sites for `envelope.activated`, `evidence.linked`, `review.recorded`,
  and `artifact.availability-observed` (DEC-224).
- `internal/journal/project_effects.go`, `current_event_catalog.go` - the fence
  key, plus the optional `stops/insert`, `attention_items/insert`, and
  `attention_stop_links/insert` effect rows added to `evidence.linked` and
  `artifact.availability-observed` under `schema.md` section 6.2 (DEC-224).
- `internal/workflow/stop.go` - the closed Stop registries (R-14). Add
  `human-outcome-required` to `stopCodes` with its exact `validStopDetails` key
  set, `human-outcome-recorded` to `conditionTypes`, and the `outcome_class`
  field to the `Condition` struct together with its exactly-one-non-null-field
  rule in `validateCondition`, all from `public-interface.md` section 2.
- `internal/journal/current_event_value_catalog.go` - the matching
  `currentStopDetailShapes` entry, and the closed `Condition` value grammar
  gaining `human-outcome-recorded` and `outcome_class:?ID` (R-14).
- `internal/export/archive_semantics_catalog.go` - mirror this task's
  `retainedCurrentRules` amendments for `evidence.linked` and
  `artifact.availability-observed` into `semanticEventRules` (R-14). Without
  the mirror, `export validate` and `import` reject the journal of any
  workspace whose readiness Stop was raised by either Event.
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

- This task belongs to QA Review window 4.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-10`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
