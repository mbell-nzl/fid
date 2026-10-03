# Phase 24 - Task 5: Review leaves and repository expectations

## Objective

Introduce `ExpectedRepository` and its shared check. Converge
`review record`:
- `repository` is replaced by a required `expected_repository`;
- coverage items, exclusions and `evidence_refs` become `SourceInput`;
- `Review.repository` is fid's observation.

`repository-mismatch` gains its amended closed details. `finding address`
gets its descriptor (two leaves leave the fence).

## Validation Contract Coverage

- VAL-009, VAL-011..016, VAL-017, VAL-018, VAL-019, VAL-022, VAL-023
  (`review.record` rows; descriptors for both leaves).
- VAL-024 (`review.record`): expectation match, HEAD differs, dirty differs.
- VAL-025 (shape): `repository-mismatch` details are exactly
  `{expected:{head_oid,dirty},observed:Repository}`.
- VAL-035 (these leaves), VAL-052.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract sections 3, 4 (the `review.record` row), 6 and 7.
- **`ExpectedRepository`** is `{head_oid:oid,dirty:boolean}`, defined once
  for reuse by tasks 7 and 8. The check compares only `head_oid` and
  `dirty` with the single observation. Any difference returns
  `repository-mismatch` with `expected` echoing the input and `observed` the
  full observation. Nothing is committed.
- The existing exact `sameRepository(request.Repository, observed)` check in
  `internal/review/review.go` is replaced by the expectation check. The
  `RepositoryMismatchError` details change shape. `review.record` is the only
  producer at this point.
- Coverage rules stay unchanged, applied to resolved digests. Items and
  exclusions may not reference the same Source digest.
- `finding address` has no input shape change.

## Scope

### Files To Create Or Modify

- `internal/project/` (`ExpectedRepository` type and check, new),
  `internal/review/review.go`
- `cmd/fid/review.go` (`repositoryDetails` mapping), `cmd/fid/finding.go`,
  and descriptor files for both leaves
- `internal/cli/codes.go` details definition if validated; code
  conformance tests
- Retained Review tests and integration corpora (`p23_*`; helpers such as
  `p23CurrentRepository` and `cloneP23Repository`): convergence with
  provenance

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` sections 3, 4, 6
- `project-docs/docs/phases/phase-24/input-inventory.md` sections 1-3
- `project-docs/docs/phases/phase-24/decisions.md` DEC-273
- `project-docs/docs/phases/phase-23/contracts.md` sections 3-5

### Out Of Scope

- Run, gate and Evidence expectation leaves (tasks 7, 8).
- Human-mode rendering of the mismatch (task 10).

## Acceptance Criteria

- [ ] `review.record` with a matching expectation is accepted, and the
  stored Review repository equals the independent observation. HEAD-differs
  and dirty-differs both return `repository-mismatch` with no Event.
  (VAL-024, VAL-022)
- [ ] `repository-mismatch` JSON details are exactly the amended closed
  shape. (VAL-025)
- [ ] Coverage, exclusion and `evidence_refs` Sources pass the
  VAL-017/018/019 rows. (VAL-017, VAL-018, VAL-019)
- [ ] An unobservable repository returns `repository-unavailable`.
  (VAL-023)
- [ ] A retry after HEAD advanced replays. (VAL-009)
- [ ] Both leaves leave the fence with VAL-011..016 rows; fake/fid agree.
  (VAL-011..016, VAL-035)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`review record` with `expected_repository` whose `dirty` differs from the worktree returns `repository-mismatch` with `expected:{head_oid,dirty}`. Today the input has no `expected_repository`.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: the `review record` input
  example and the `repository-mismatch` details.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root the Runner records. No command is a placeholder.
Apply gofmt and goimports to edited Go files before the final gate. Every
command below is required before task close. The tagged integration gate is
mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

The `-timeout 30m` bounds follow DEC-227/259. Record the wall time `go test`
reports for `cmd/fid/integration`, `cmd/fid` and `internal/workflow`. The
budgets are tagged integration under 900 s and no ordinary package over 600 s
(VAL-052). A breach caused by this task's own new journeys is a bounded speed
repair owned by this task; otherwise stop and escalate. Never record a
timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands). Every required command
is reported `PASS` or `BLOCKED` with reason, with its own output token and
observed termination. Focused red/green selections never replace these
commands. Report lint is a Runner control, not a Developer gate.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-5-done.md`.
  Commit with `Refs: phase-24/task-5`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
- Remove this task's leaves from the pending-descriptor fence (task 3). Each
  adopted leaf gets a descriptor, help, `--example`, typed template,
  exclusivity and structural-validation corpus rows (VAL-011..016).
