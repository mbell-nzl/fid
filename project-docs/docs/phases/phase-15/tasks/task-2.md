# Phase 15 - Task 2: `internal/dev` foundation + `fid dev start` + `fid dev gate`

## Objective

Stand up the developer-composite surface and ship the two "opening" composites:

- New `internal/dev` domain package (PI-1) that orchestrates the existing
  `claims`/`issues`/`runs`/`workflow` effect builders into **one journal event
  per composite** (PI-2), validate-before-append.
- Add the `dev.start`, `dev.gate`, `dev.finish`, `dev.block` rows to the in-code
  permissions matrix (copy `03_TRD_SCHEMA.md` section 13 verbatim: `developer`
  and `human` only).
- `fid dev start <id> --by developer --model <m>`: one `dev.started` event =
  acquire the caller's claim + transition `todo → in_progress` + open one
  `agent_runs` row (`ended_at` NULL).
- `fid dev gate <id> --cmd <c> --exit-code <n> [--output-excerpt <s>] --by
  developer`: one `dev.gate.recorded` event = append one `gate_runs` row on the
  derived active attempt, recorded only when the issue is started-category and
  the caller holds the live claim.

`fid dev finish`/`fid dev block` are Task 3.

## Requirements

- PRD refs: PRD-DEV-008, PRD-DEV-004, PRD-DEV-010, PRD-CLAIM-005
- TRD refs: TRD-ARCH-010; `02_TRD.md` §§2-3 (package boundaries), §4.7 (gate run
  persistence); `03_TRD_SCHEMA.md` §2 (`claims`, `agent_runs`, `gate_runs`),
  §7 (`agent_runs`), §13 (`dev.*` matrix rows)
- Design refs: `06_DESIGN_PLAN_CLI.md` §5 `fid dev start`, `fid dev gate`
- Quality refs: `08_QUALITY_PLAN.md` (validate-before-append, size budgets)

## Validation Contract Coverage

Assertions live in
`project-docs/docs/phases/phase-15/qa/validation-contract.md`.

- VAL-003: `dev start` emits one `dev.started` event = claim + `in_progress`
  transition + one open `agent_runs` row; issue/claim/run state correct.
- VAL-004: `dev start` rejects before append (zero new lines) on a conflicting
  live claim or a disallowed edge; valid unclaimed `todo` path succeeds via
  claim-in-same-event (PI-3).
- VAL-005: `dev gate` appends one `gate_runs` row on the derived active attempt
  with reported `command`/`exit_code`/`output_excerpt`; rejects before append
  without a started state + live caller claim (DN-A, PI-2).
- VAL-006 (this task adds the four matrix rows; finish/block authorization is
  exercised in Task 3): matrix has `dev.start`/`dev.gate`/`dev.finish`/
  `dev.block` rows matching section 13; `fid policy show` reflects them;
  developer allowed, denied role rejected before append.

## Scope

### Files To Create Or Modify

- `internal/dev/dev.go` (new package; PI-1) - shared composite scaffolding:
  one `journal.WithLock` + `EnsureProjectionFreshWithLock` + open conn +
  validate-before-append + `PrepareEvent` + assemble effects +
  `DryRunProjectEvent` → `AppendEvent` → `ProjectEvent`. Action constants
  `dev.start`/`dev.gate`/`dev.finish`/`dev.block`; `fidVersion`.
- `internal/dev/start.go` (new) - `Start(options)`: authorize `dev.start`,
  validate issue exists + kind task/subtask + `todo → in_progress` edge allowed
  + no live claim by another role; build one `dev.started` event with effects:
  (1) `claims` insert (caller, default TTL), (2) `issues` status update
  `todo → in_progress` + the `status_changed` `activity_log` effect (same shape
  as `internal/issues/transition.go`), (3) `agent_runs` insert
  (`role`=`--by`, `model`=`--model`, `started_at`=now, `ended_at` NULL).
- `internal/dev/gate.go` (new) - `Gate(options)`: authorize `dev.gate`, require
  started-category state + live caller claim, derive active attempt
  (`MAX(task_completions.attempt_number)+1`, 0 when none), generate
  `GATE-<hash>`, build one `dev.gate.recorded` event with one `gate_runs`
  insert (`command`, `exit_code`, `output_excerpt`, `ran_at`,
  `recorded_by_role`, `attempt_number`, optional `agent_run_id` = the open run).
- `internal/dev/*_test.go` (new) - domain tests for both composites incl.
  one-event, effect shapes, attempt derivation, reject paths.
- `internal/policy/matrix.go` - append `dev.start`, `dev.gate`, `dev.finish`,
  `dev.block` rows (`developer: true, human: true`) copying section 13 verbatim
  (DEC-026 pattern).
- `internal/policy/matrix_test.go` - extend completeness/coverage test for the
  four rows.
- `cmd/fid/dev.go` (new) - `newDevCommand()` + `dev start`/`dev gate` cobra
  subcommands with flags per the CLI spec + the `--exit-code`/`--output-excerpt`
  additions (DN-A) and `--model` on start (DN-B); stderr confirmation output
  matching `06_DESIGN_PLAN_CLI.md`; `renderIssueError` on failure.
- `cmd/fid/dev_test.go` (new) - CLI tests incl. allow/deny role, reject paths,
  warning/confirmation lines.
- `cmd/fid/root.go` - register `newDevCommand()`.
- `project-docs/specs/02_TRD.md` - §§2-3: add the `internal/dev` domain package
  (developer composites orchestrating claim/transition/run/gate effects into one
  event). Bump "Last updated".
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - `fid dev start`: add `--model`
  (required) to the flag table; `fid dev gate`: add `--exit-code` (required) and
  `--output-excerpt` (optional) to the flag table and note fid records — does
  not execute — the command (DN-A/DN-B). Align confirmation examples.

### Files To Read For Context

- `internal/issues/transition.go` (started-category enter/leave, status update +
  activity effect, claim integration)
- `internal/claims/acquire.go`, `internal/claims/shared.go`,
  `internal/claims/release.go` (claim effect builders;
  `BuildIssueReleaseEffect`/`MaybeBuildIssueReleaseEffect`,
  `RequireLiveIssueClaimHeld`, `loadCurrentClaim`)
- `internal/runs/record.go` (`agent_runs` insert effect shape, `RUN-<hash>`)
- `internal/issues/complete.go` (gate snapshot + attempt derivation;
  `nextCompletionAttemptNumber`, `completionGateSnapshot`)
- `internal/workflow/workflow.go` (edges, `CategoryForState`, `RulesForEdge`)
- `internal/policy/matrix.go`, `internal/policy/authorize.go`
- `cmd/fid/run.go`, `cmd/fid/claim.go` (cobra + error-render pattern)

### Out Of Scope

- `fid dev finish` / `fid dev block` (Task 3) and CF-013 extraction.
- `fid dev resume` / `fid dev abandon`, run `status`/lifecycle, git context (all
  v0.2c).
- Computing cost/tokens on the dev-opened run (DN-B: model only; cost/tokens
  NULL).
- Executing the gate command (DN-A: reported exit code only).

## Acceptance Criteria

- [ ] `fid dev start <id> --by developer --model <m>` on a fresh `todo` task
  appends exactly one `dev.started` event with a claim insert, a
  `todo → in_progress` transition (+ `status_changed` activity), and one open
  `agent_runs` row; issue is `in_progress`, `fid claim list` shows the live
  claim, the run is open. (VAL-003)
- [ ] `dev start` rejects before append (non-zero exit, zero new journal lines,
  actionable error) on a live claim held by another role and on a non-`todo`
  issue; the valid unclaimed `todo` path succeeds via claim-in-same-event. (VAL-004)
- [ ] `fid dev gate <id> --cmd <c> --exit-code <n> [--output-excerpt <s>] --by
  developer` appends one `dev.gate.recorded` event + one `gate_runs` row on the
  derived active attempt recording the reported fields; rejects before append
  (zero new lines) when the issue is not started-category or the caller lacks
  the live claim. fid does not execute the command. (VAL-005)
- [ ] In-code matrix has exactly one row each for `dev.start`/`dev.gate`/
  `dev.finish`/`dev.block` (developer + human) matching section 13;
  `fid policy show` includes them; developer allowed, `qa_reviewer` rejected
  before append naming role + action. (VAL-006)
- [ ] `02_TRD.md` records the `internal/dev` boundary; `06_DESIGN_PLAN_CLI.md`
  records the `--model`/`--exit-code`/`--output-excerpt` flags. (doc)

## TDD Plan

- First failing test: `internal/dev` `Start` for a fresh `todo` task returns a
  result and writes the three effects — fails red (package absent).
- Additional tests: `dev start` reject on conflicting claim / disallowed edge;
  exactly one journal line per composite; `dev gate` attempt-number derivation
  (no prior completion ⇒ attempt 1; after a completion ⇒ next) and reject
  without started+claim; matrix completeness for the four rows; CLI allow
  (developer) and deny (`qa_reviewer`) reject-before-append.
- Regression risks: lock ordering — refresh projection inside `WithLock`
  (DEC-028), never nest `.fid/write.lock`; do not break the CF-007
  validate-before-append seam; keep new files ≤400 lines and functions ≤60;
  reuse the existing claim/transition/run effect shapes byte-for-byte so rebuild
  stays deterministic.

## Documentation Updates

- `project-docs/specs/02_TRD.md` - add the `internal/dev` domain package (§§2-3).
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` - `fid dev start --model`;
  `fid dev gate --exit-code`/`--output-excerpt` + "records, does not execute"
  note; align confirmation examples.

## Scenario / Parity / Audit Notes

- Held-out scenario impact: None.
- Cross-surface parity impact: `fid policy show` must reflect the four new
  `dev.*` rows (show-vs-enforce parity).
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
`internal/dev`, `internal/policy`, and `cmd/fid` coverage percentages.

Conditional commands:

```text
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- `go test -tags integration -race ./...` — required: the package set under test
  (`cmd/fid`) contains integration-tagged tests; record PASS/BLOCKED.
- `go build -o /tmp/fid ./cmd/fid` — required: CLI build used by the new
  `dev` command tests.

## Notes For Developer

- PI-1: `internal/dev` is a new inward-pointing domain boundary; record it in
  `02_TRD.md` and the completion report.
- PI-2: each composite emits exactly **one** event with all its effects;
  validate everything before the first append (CF-007). The projector is
  effect-driven — confirm no event-type registry needs the new
  `dev.started`/`dev.gate.recorded` types; if one exists, update it and note it.
- PI-3: in `dev start`, the claim is acquired in the same event, so do **not**
  call the DB-based "require live claim held" check that the plain transition
  uses for the started-category enter — the claim insert in this event satisfies
  PRD-CLAIM-005. Still reject when another role holds a live claim.
- DN-A/DN-B are public-interface additions confirmed at kickoff; reflect them in
  the CLI spec in this task. If the human has not confirmed at task start, stop
  and escalate rather than guessing.
- Reuse the existing effect shapes (claims insert/release, issues status update +
  `status_changed` activity, agent_runs insert) rather than re-deriving them, so
  replay stays deterministic.
- `GATE-<hash>` and `RUN-<hash>` ids via `journal.GenerateID`.
