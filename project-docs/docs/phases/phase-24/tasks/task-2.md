# Phase 24 - Task 2: Input-as-written digests and named conflict fields, all families

## Objective

Replace the 11 family-local request-digest helpers with one shared helper.
It computes the request digest from the input as written plus one digest per
top-level public input name, and it is adopted at all 43 call sites. Every
family's `idempotency-conflict` JSON details gain `differing_fields`
(DEC-275). This pins the digest rule before the input shapes change.

## Validation Contract Coverage

- VAL-002 (every keyed leaf): rebuild parity of the field-digest maps.
- VAL-005: the 100k rebuild stays under 30000 ms.
- VAL-006: conflicts name exactly the differing names, for every keyed leaf.
- VAL-007: canonically equal retries replay, for every keyed leaf.
- VAL-008 (review case): optional-member presence; the DEC-242 exception.
- VAL-009 and VAL-066 (the seven `internal/review` leaves that observe today:
  `concern.record`, `concern.dispose`, `deviation.record`,
  `human-outcome.record`, `decision.record-reference`,
  `requirement.record-reference` and `requirement.link`): observed values
  never enter the digest, and a replay never observes.
- VAL-052: gate wall-time budgets.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract section 7. The digest keeps each family's existing components
  (operation where used, project ID, resolved Actor or standalone
  attribution, request). The request is the closed canonical input as the
  caller wrote it. An optional member is present exactly when written. The
  only exception is DEC-242 (`inapplicable_findings: []` equals omission).
- Field digests: one entry per top-level public input name, plus `--actor`.
  - For `--input` leaves the names are the input's top-level field names.
  - For flag and positional leaves (issue, label, link, workflow) the names
    are the flag spellings with `--`, and positional names as shown in the
    leaf's `Use` line without brackets.
  - Each digest is `sha256:` over the canonical JSON of the member's value as
    written. An absent member has no entry.
- Compute digests from the decoded input before any observation or domain
  lookup. Never mutate the request afterwards. The idempotency replay check
  runs before any observation (contract section 2). Move it earlier in any
  of the seven leaves above that observes first today.
- `differing_fields` is the sorted unique set of names whose digest differs,
  or that are present on only one side. It is never empty. Every family's
  conflict renderer adds it to the closed details. Human conflict lines stay
  unchanged until task 10.
- `gate.run` records idempotency on `gate.execution-started`, the Receipt
  update and the committed `run.execution-blocked` denial. All of these
  carry real maps.

## Scope

### Files To Create Or Modify

- `internal/journal/` shared helper (new file) and `idempotency.go`
- The 32 files with `RequestDigest:` sites (inventory R1). Remove the 11
  helpers (R2) in `internal/{issues,labels,links,workflow,evidence,review}`
- `cmd/fid` conflict branches mapping `IdempotencyConflictError` in every
  family renderer; `internal/cli/codes.go` details definition if validated
- Tests: a per-keyed-leaf conflict and retry corpus (new), plus convergence of
  retained conflict-detail assertions

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` sections 6, 7
- `project-docs/docs/phases/phase-24/input-inventory.md` R1, R2
- `project-docs/docs/phases/phase-24/decisions.md` DEC-275
- `project-docs/docs/phases/phase-23/public-interface.md` "Task 5 input
  addition" (DEC-242 canonicalization)

### Out Of Scope

- Human-mode conflict output (task 10). Input shape changes (tasks 4-8).
- New observation sites. VAL-009 rows for leaves that start observing in
  tasks 4-8 belong to those tasks.

## Acceptance Criteria

- [ ] The corpus has a one-field-change row for every keyed leaf: the 50
  live leaves with `--idempotency-key` minus `export import`, so 49 at
  planning (inventory R16). It also has a two-field row and an `--actor` row in each of the 20
  families. `differing_fields` is exact and sorted, the exit is 4, and the
  journal line count is unchanged. A completeness test walks the live tree
  and fails when a keyed leaf has no row. (VAL-006)
- [ ] For every keyed leaf, a canonically equal retry under the same key
  replays byte-identically except `replayed:true`, with the original
  EventRef. For `--input` leaves, the retry reorders keys and changes
  whitespace. For flag and positional leaves, it permutes the flag order.
  (VAL-007)
- [ ] `review.record`: nonempty `inapplicable_findings` versus omitted
  conflicts naming it; `[]` versus omitted replays. (VAL-008)
- [ ] For each of the seven observing leaves above: commit, advance HEAD (and
  change the referenced file where the leaf has one), then retry. The retry
  replays byte-identically with no Event and no failure. (VAL-009)
- [ ] For each of the seven: make the worktree unobservable (unmerged paths),
  then retry. Where the leaf references a project file, also delete the file
  and retry. Each retry replays with no Event. A control request under a new
  key in the same state fails with no Event (any code until task 4 adds
  `repository-unavailable`). (VAL-066)
- [ ] Rebuild parity of the maps over a journal with one keyed mutation per
  keyed leaf. (VAL-002)
- [ ] `rg` over R2 shows only the shared helper.
- [ ] Benchmark and wall times inside budget. (VAL-005, VAL-052)

## First Failing Test

A `run.start` conflict test changes `claim_version` under the same key and expects `differing_fields:["claim_version"]`. Today the details have no such member.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: the idempotency-conflict
  details and the input-as-written rule.

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
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
rg -n "RequestDigest:" internal cmd -g '!*_test.go'
rg -n "func [a-zA-Z]*(RequestDigest|requestDigest|judgmentDigest)\(" internal -g '!*_test.go'
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

- The benchmark: every keyed write computes field digests (VAL-005).
- The two `rg` lines regenerate R1/R2. The second must list only the shared
  helper, so the helper's name must match its pattern
  (`func [a-zA-Z]*(RequestDigest|requestDigest|judgmentDigest)\(`). Both are
  expected to terminate `exit_code: 0`.
- R16 (keyed leaves) is enforced by the corpus completeness test, which walks
  the live tree in-process.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-2-done.md`.
  Commit with `Refs: phase-24/task-2`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
