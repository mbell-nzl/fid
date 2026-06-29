# Phase 15 - Task 3: CF-013 completion-effect extraction + `fid dev finish` + `fid dev block`

## Objective

Ship the two "closing" composites, reusing a shared completion-effect builder
extracted from `internal/issues/complete.go` (CF-013):

- **CF-013:** extract the `Complete` orchestration into a reusable,
  status-parameterised completion-effect builder. `fid issue complete` behavior
  stays identical (status `dev_complete`); `fid dev finish`/`fid dev block`
  consume the same builder.
- `fid dev finish <id> --by developer --completion-md <p> [--tokens --commits
  --files --human-intervention --intervention-note]`: one `dev.finished` event =
  close the issue's open `agent_runs` row + insert a `task_completion`
  (`dev_complete`, gate snapshot, supersede priors) + transition
  `in_progress → dev_complete` + release the caller's claim.
- `fid dev block <id> --by developer --reason-md <p>`: one `dev.blocked` event =
  insert a `task_completion` (`failed_gate`, gate snapshot, reason, supersede
  priors) + transition `in_progress → failed_gate` + release the caller's claim.

## Requirements

- PRD refs: PRD-DEV-008, PRD-DEV-009, PRD-DEV-010, PRD-WF-007..010
- TRD refs: `02_TRD.md` §4.7 (gate snapshot at finish/block); `03_TRD_SCHEMA.md`
  §2 (`task_completions`, `gate_runs`, `agent_runs`, `claims`, `issues`), §13
  (`dev.finish`/`dev.block`)
- Design refs: `06_DESIGN_PLAN_CLI.md` §5 `fid dev finish`, `fid dev block`
- Quality refs: `08_QUALITY_PLAN.md` (validate-before-append, size budgets);
  CF-013

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-15/qa/validation-contract.md`.

- VAL-007: `dev finish` emits one `dev.finished` event = close open run +
  completion(`dev_complete`, gate snapshot, supersede) + transition + claim
  release (PI-5/PI-6).
- VAL-008: `dev finish` rejects before append (zero new lines) when
  `all_acs_verified[developer]` is unmet; `has_completion_report` is satisfied by
  the in-event completion (PI-4).
- VAL-009: `dev block` emits one `dev.blocked` event = completion(`failed_gate`,
  gate snapshot, reason, supersede) + transition + claim release (PI-5).
- VAL-010: CF-013 — `Complete` orchestration extracted into a shared
  completion-effect builder reused by `internal/dev`; `fid issue complete`
  behavior unchanged; functions ≤60 lines.
- VAL-006 (finish/block authorization exercised here): denied role rejected
  before append for `dev.finish`/`dev.block`.

## Scope

### Files To Create Or Modify

- `internal/issues/complete.go` - CF-013: extract a reusable completion-effect
  builder (e.g. `BuildCompletionEffects(conn, params) ([]journal.Effect,
  CompleteResult, error)`) covering id generation, attempt number, gate snapshot,
  prior-supersession, and value assembly, **parameterised by completion status**
  (`dev_complete` for `Complete`/`dev finish`; `failed_gate` for `dev block`).
  `Complete` becomes a thin wrapper around the builder with status `dev_complete`
  — behavior identical. Keep `Complete` and the new helpers ≤60 lines each.
- `internal/issues/complete_test.go` - keep existing behavior assertions; add a
  direct test for the builder if helpful (no behavior change to `fid issue
  complete`).
- `internal/dev/finish.go` (new) - `Finish(options)`: authorize `dev.finish`,
  load issue/edge, evaluate `all_acs_verified[developer]` against existing
  verifications (reject before append when unmet); build one `dev.finished` event
  with effects: (1) conditional `agent_runs` close (set `ended_at` on the open
  run when one exists, PI-6), (2) the completion effects from the shared builder
  (status `dev_complete`), (3) `issues` status update `in_progress → dev_complete`
  + `status_changed` activity, (4) claim release via
  `claims.BuildIssueReleaseEffect` (PI-5).
- `internal/dev/block.go` (new) - `Block(options)`: authorize `dev.block`,
  validate the `in_progress → failed_gate` edge; build one `dev.blocked` event
  with effects: (1) the completion effects (status `failed_gate`, reason →
  `notes_md_path`), (2) `issues` status update `in_progress → failed_gate` +
  activity, (3) claim release (PI-5). Does not touch `agent_runs` (PI-6).
- `internal/dev/*_test.go` (new) - domain tests for finish (happy path, run
  close, unmet-AC reject, no-prior-completion success) and block (happy path,
  reason recorded, supersession).
- `cmd/fid/dev.go` - add `dev finish`/`dev block` subcommands with flags per the
  CLI spec; stderr confirmation output matching `06_DESIGN_PLAN_CLI.md`.
- `cmd/fid/dev_test.go` - CLI tests incl. allow/deny role and reject paths.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - align `fid dev finish`/`fid dev
  block` confirmation examples with the implementation if they diverge.

### Files To Read For Context

- `internal/issues/complete.go` (the orchestration to extract; gate snapshot,
  attempt, supersession)
- `internal/issues/transition.go`, `internal/issues/transition_rules.go`
  (`evaluateTransitionRules`, `all_acs_verified`/`has_completion_report`, status
  update + activity effect)
- `internal/claims/release.go` (`BuildIssueReleaseEffect`)
- `internal/dev/dev.go`, `internal/dev/start.go` (Task 2 scaffolding/effect
  patterns)
- `internal/workflow/workflow.go` (the `dev_complete`/`failed_gate` edges + rules)

### Out Of Scope

- `fid dev resume`/`abandon`, run lifecycle, git context (v0.2c).
- Any change to `fid issue complete` behavior or its CLI/JSON/journal contract.
- Raising concerns or escalations as part of `dev block`.

## Acceptance Criteria

- [ ] `fid dev finish` appends exactly one `dev.finished` event that closes the
  open run (`ended_at` set), inserts one `task_completions` row (`dev_complete`,
  gate snapshot + `gates_passed`/`gates_failed`, prior completions superseded),
  transitions `in_progress → dev_complete`, and releases the caller's claim. (VAL-007)
- [ ] `fid dev finish` rejects before append (non-zero exit, zero new journal
  lines, error naming `all_acs_verified`) when developer ACs are unverified, and
  succeeds with no pre-existing completion once developer ACs pass. (VAL-008)
- [ ] `fid dev block` appends exactly one `dev.blocked` event that inserts one
  `task_completions` row (`failed_gate`, gate snapshot, reason as `notes_md_path`,
  priors superseded), transitions `in_progress → failed_gate`, and releases the
  caller's claim. (VAL-009)
- [ ] `fid issue complete` behavior is unchanged after the CF-013 extraction
  (existing complete tests pass unmodified except for moves); the shared builder
  is reused by `internal/dev`; `Complete` and the extracted helpers are each
  ≤60 lines. (VAL-010)
- [ ] `dev.finish`/`dev.block` deny a non-developer role before append
  (non-zero exit, zero new journal lines, error names role + action). (VAL-006)

## TDD Plan

- First failing test: `internal/dev` `Finish` after start+verify+gate produces
  one event closing the run + completion + transition + release — fails red.
- Additional tests: finish reject on unverified developer ACs; finish success
  with no prior completion (has_completion_report satisfied in-event); block
  produces `failed_gate` completion + transition + release; supersession of prior
  completions; CF-013 — `Complete` still passes its existing tests with the
  builder underneath; CLI deny path for finish/block.
- Regression risks: ordering of effects so the projector applies them cleanly;
  do not let the `has_completion_report` rule block finish (the completion is
  in-event, PI-4); keep one event per composite; keep the release effect
  preconditioned (reuse `BuildIssueReleaseEffect`); functions ≤60 lines, files
  ≤400.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - align `fid dev finish`/`fid dev
  block` confirmation examples with the implementation if they diverge.
- Record the CF-013 extraction in the completion report (decision-backed
  behavior preservation).

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: `fid dev finish`/`block` reuse the same
  transition/claim-release/completion effect shapes as the single-command paths;
  keep them identical for rebuild determinism.
- Audit or release-readiness impact: None.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close. Record each command in
the completion report as `PASS`, or `BLOCKED` with reason. Record the
`internal/dev`, `internal/issues`, and `cmd/fid` coverage percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` — required: the package set under test
  (`cmd/fid`) contains integration-tagged tests; record PASS/BLOCKED.
- `go build -o /tmp/fid ./cmd/fid` — required: CLI build used by the new
  finish/block command tests.

## Notes For Developer

- CF-013 first: extract the completion-effect builder and prove `fid issue
  complete` is byte-for-byte unchanged before wiring finish/block to it. This is
  a decision-backed extraction (CF-013), not a behavior change.
- PI-4: build the completion effect, then evaluate only `all_acs_verified[
  developer]` against the DB for the `dev_complete` edge; treat
  `has_completion_report` as satisfied by the in-event completion. Do **not**
  require a pre-existing completion for `dev finish`.
- PI-5: both composites release the caller's claim even though `dev_complete`/
  `failed_gate` stay in the started category; reuse
  `claims.BuildIssueReleaseEffect` (preconditioned delete).
- PI-6: `dev finish` closes the issue's open run only when one exists; `dev
  block` does not touch runs.
- Keep one event per composite and the CF-007 validate-before-append seam intact.
- This is the heaviest task in the phase (CF-013 extraction + two composites). If
  it grows beyond a single session, stop and escalate to split block into its
  own follow-up rather than rushing the extraction.
