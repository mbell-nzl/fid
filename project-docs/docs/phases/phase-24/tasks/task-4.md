# Phase 24 - Task 4: Repository observation and SourceInput; judgment reference leaves

## Objective

Introduce the shared foundation for the rest of Part 1:
- one repository-observation function with typed `repository-unavailable`
  reasons;
- `SourceInput` resolution, which checks a supplied project-file digest,
  observes an omitted one, and enforces path safety;
- the amended `source-digest-mismatch` details.

Converge the judgment reference leaves: `decision record-reference`,
`requirement record-reference`, `human-outcome record`, `concern record` and
`deviation record`. `requirement link` and `concern dispose` get descriptors
only (seven leaves leave the fence). Human Outcome references become
checked.

## Validation Contract Coverage

- VAL-008 (Source digest case): supplied versus omitted digest conflicts.
- VAL-009 and VAL-066 (the seven observing leaves of this task:
  `concern.record`, `concern.dispose`, `deviation.record`,
  `human-outcome.record`, `decision.record-reference`,
  `requirement.record-reference`, `requirement.link`): observed values never
  enter the digest, and a replay never observes.
- VAL-011..016 and VAL-060 (these seven leaves).
- VAL-017 (five changed leaves): removed `repository` members rejected.
- VAL-018, VAL-019 (these leaves' Source positions): supplied digest
  checked; omitted digest observed.
- VAL-020, VAL-021: path safety; external-reference digest.
- VAL-022 (the seven observing leaves): the recorded repository is fid's
  observation.
- VAL-023 (code, classifier, and the seven leaves' rows):
  `repository-unavailable` with the mapped reason.
- VAL-035 (these leaves, where the fake implements them; the derived set may
  be empty).
- VAL-052.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract sections 2, 3, 4 (rows for the five changed leaves), 6 and 7.
- **Observation.** One shared function, in `internal/project` beside
  `ObserveRepository`, returns the full recorded value: config `project_id`,
  `base_oid` = HEAD, and `capture_source:"git"`. It returns a typed error
  carrying one of the five reasons. Classify by type, never by message text,
  so each `unsafeRepository` cause in `repository.go` gains a typed cause.
  Typed causes still satisfy `errors.Is(err, project.ErrUnsafeRepository)`,
  so callers outside the observing-leaf list keep their behaviour
  (`run.recover`, `run.checkpoint`, `run.inspect`, close/abandon until
  task 7). The `internal/review` judgment leaves switch to it. Workflow and
  evidence callers switch in their own tasks.
- **Reason mapping** (contract section 6, DEC-278), by capture stage:
  - `not-a-repository`: `rev-parse --show-toplevel` fails, or HEAD or its
    tree does not resolve (a repository with no commit);
  - `tracked-flags`: an assume-unchanged or skip-worktree entry;
  - `unmerged-paths`: an unmerged index entry or an unmerged Git image;
  - `unsafe-path`: an invalid or non-UTF-8 path, an unsupported file type,
    mode or object type, or the path-count or output-size limit;
  - `changed-during-capture`: the two reads differ, or any other Git or file
    read fails after the root resolved.

  `detail` names the specific cause. Failures that are not repository
  observation (for example canonical-JSON encoding) stay internal errors.
- **`repository-unavailable`** is added to `internal/cli/codes.go` exactly as
  contract section 6 states: precondition, retryable, exit 4, null subject,
  `{reason,detail}`.
- **`SourceInput`** is the Phase 21 `Source` minus `repository`, with
  `digest` optional only for `project-file`.
  - The path resolves under the project root. Absolute paths and `.`/`..`
    segments are structural (`invalid-argument`).
  - A missing file, unreadable file, directory, non-regular file or a
    symlink escaping the root returns `source-digest-mismatch` with
    `observed_digest:null`. `expected_digest` is null when the digest was
    omitted.
  - Resolution generalizes `verifyDecisionSource` (whole-file SHA-256).
  - The stored Source keeps the Phase 21 output shape, with `digest`
    (supplied or observed) and `repository` (this command's observation).
- **Order.** Closed input, then idempotency replay, then the leaf's existing
  checks. Source resolution sits where `verifyDecisionSource` runs today;
  for leaves that had no check, after authority and before commit.
  Observation happens once, and never on a replay.
- `source-digest-mismatch` producers become every SourceInput leaf.
- **Source positions** are enumerated from the descriptor registry (every
  `SourceInput`-typed field of this task's leaves), with a completeness check
  that fails when a position has no corpus row.

## Scope

### Files To Create Or Modify

- `internal/project/` observation function and Source resolution (new file),
  with tests; `internal/project/repository.go` (typed unsafe causes)
- `internal/review/{decision,requirement,human_outcome,concern,deviation,judgment,review}.go`
  (`observedRepository` callers in this package)
- `internal/cli/codes.go` and the code conformance tests
- `cmd/fid/{decision,requirement,human_outcome,concern,deviation}.go`, their
  renderers' mapping for the two codes, and descriptor files for the seven
  leaves
- Retained tests and integration corpora building these inputs
  (`cmd/fid/integration/p23_*`, fake adapter cases): convergence with
  provenance

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` sections 2, 3, 4, 6, 7
- `project-docs/docs/phases/phase-24/input-inventory.md` sections 1-3 (R6, R7, R8)
- `project-docs/docs/phases/phase-24/decisions.md` DEC-272, DEC-277
- `internal/review/decision.go` `verifyDecisionSource`; `internal/project/repository.go`
- `project-docs/docs/phases/phase-23/public-interface.md` sections 3-5

### Out Of Scope

- Review, Envelope, Claim, Run, gate and Evidence leaves (tasks 5-8), and
  `ExpectedRepository` (task 5).
- Human rendering changes beyond mapping the two codes through the
  existing family renderers (task 10).

## Acceptance Criteria

- [ ] For each changed leaf, every Source position enumerated from its
  descriptor: matching digest accepted and stored; mismatch rejected with
  both digests and no Event. (VAL-018)
- [ ] With the digest omitted: an existing file is stored with its observed
  SHA-256; a missing file is rejected with both digests null. (VAL-019)
- [ ] Path-safety fixtures (absolute, `..`, escaping symlink, directory,
  in-root symlink) behave as VAL-020 states. (VAL-020)
- [ ] External reference: an omitted digest returns `invalid-argument` naming
  it; a supplied digest is stored with no file read. (VAL-021)
- [ ] Supplying `repository` inside any Source returns `invalid-argument`
  naming its path. (VAL-017)
- [ ] Stored Source repositories and Event repositories equal the
  independent observation. (VAL-022)
- [ ] For each of the seven observing leaves, unmerged paths and a non-Git
  directory return `repository-unavailable` with reasons `unmerged-paths`
  and `not-a-repository`, exit 4, details exactly `{reason,detail}` and no
  Event. A normal worktree succeeds. (VAL-023)
- [ ] The classifier maps every cause class to its reason: real fixtures for
  a repository with no commit (`git init` without a commit), a tracked flag
  (`git update-index --assume-unchanged`) and an unsupported file type (for
  example an untracked FIFO); injected unit rows for a capture change and a
  size limit. (VAL-023)
- [ ] For each of the seven, a retry after HEAD advanced (and the file
  changed, digest omitted) replays. Supplied versus omitted digest under one
  key conflicts naming the top-level field. (VAL-009, VAL-008)
- [ ] For each of the seven, a retry in an unobservable worktree, and (where
  the leaf references a project file) after deleting that file, replays with
  no Event. New-key controls fail
  with `repository-unavailable` and `source-digest-mismatch` respectively.
  (VAL-066)
- [ ] The seven leaves leave the fence with VAL-011..016 and VAL-060 rows.
  (VAL-011..016, VAL-060)
- [ ] Fake/fid conformance rows agree. (VAL-035)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`human-outcome record` with a reference whose digest does not match the file returns `source-digest-mismatch`. Today it is accepted.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: SourceInput and the judgment
  leaf input examples.
- `project-docs/specs/GLOSSARY.md`: **Expected repository** and
  **SourceInput** terms, if absent.

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
reports for `cmd/fid/integration`, `cmd/fid`, `internal/workflow` and the
slowest other ordinary package. The budgets are tagged integration under
900 s and no ordinary package over 600 s (VAL-052); report any ordinary
package over the 480 s warning. A breach caused by this task's own new journeys is a bounded speed
repair owned by this task; otherwise stop and escalate. Never record a
timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands). Every required command
is reported `PASS` or `BLOCKED` with reason, with its own output token and
observed termination. Focused red/green selections never replace these
commands. Report lint is a Runner control, not a Developer gate.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-4-done.md`.
  Commit with `Refs: phase-24/task-4`.
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
