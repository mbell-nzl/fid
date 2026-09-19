# Phase 23 - Task 4: Immutable Review Records and computed independence

## Objective

Implement `review.record` as one atomic mutation that creates an immutable
Review Record with its coverage, verdict, and zero or more Findings, and
enforce independence from durable Claim and Run history rather than from a role
name.

## Validation Contract Coverage

- VAL-014: The Review binds Run, Issue, Envelope revision, re-observed
  repository identity, profile, coverage, verdict, and evidence references, and
  rejects a Receipt produced by another Run and a coverage item and exclusion
  sharing one `Source` digest.
- VAL-015: Independence is enforced in both directions from Claim and Run
  history, over the mutating side only (DEC-225): a mutating-Claim holder and a
  Run holder under such a Claim reject; a distinct Actor succeeds, and so does
  that same Actor's second Review after a rework cycle.
- VAL-003: task 1's fence assertion for `review.recorded` converges here
  without losing its parity subject.
- VAL-016: A repository identity mismatch rejects with no Event.
- VAL-017: Verdict and Finding consistency is enforced for all three verdicts.
- VAL-018: Exactly one Review exists per review Run.
- VAL-019: Review recording is idempotent under its key and digest.

## Binding Constraints

- `contracts.md` sections 3 and 4 and `public-interface.md` sections 3 and 4
  are the exact authority. `schema.md` section 2 fixes the storage.
- Re-observe the repository under the fid write lock and compare it with the
  supplied identity before any durable effect. Never silently correct it.
- Independence is computed: the recording Actor must not equal any Actor that
  held a mutating Claim or a Run on the reviewed Issue within the same Envelope
  revision. The rejection names the exact conflicting subject and Actor. No
  role string, actor-type heuristic, or attribution label participates.
- Coverage is stored and reported, never validated as true. fid asserts nothing
  about whether the Reviewer actually examined a declared subject.
- Referenced Receipt IDs must belong to the recording Run.
- DEC-192: delete exactly the `"review.recorded"` key from
  `unavailableJudgmentEvents` as part of implementing this Event, and prove
  live-apply and replay parity for it. Leave every other key fenced.
- DEC-192 scheduled convergence: this change displaces task 1's
  `review.recorded` fence rejection assertion. Converge it by substituting a
  still-fenced Event, preserving the VAL-003 subject. Never delete the parity
  assertion; record introducing-commit provenance and the executing command.
- Add the `review.record` grant action exactly here, per
  `public-interface.md` section 4.1. Add no other action.
- Do not implement or assert review-cycle charging, Evidence linking, Finding
  dispositions beyond creation, or staling. Those are tasks 5, 6, and 7. In
  particular, do not assert the absence of a `review-cycles` charge.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/review/` - new package: Review recording, independence, validation.
- `internal/workflow/authority.go` and `internal/workflow/activation.go` - the
  one new grant action. The closed allowlist is `sortedUniqueCapabilities` in
  `activation.go`; Envelope grants reuse it through `sortedUniqueEnvelopeActions`
  in `envelope.go`, and `authority.go` holds no action list (R-08, DEC-167/173).
- `internal/journal/project_effects.go` - the one fence key and its effects.
- `internal/journal/current_event_json.go`,
  `internal/journal/current_event_result_projection.go`, and
  `internal/journal/current_event_value_catalog.go` - the typed-result
  registration for `review.recorded` (R-16, DEC-231). Deleting the fence key
  routes this family onto the journal's strict result path, so all three are
  required in this task or the first successful `review.record` is rejected at
  append time. Add: the `review.recorded` entry in the `fields` grammar map in
  `current_event_json.go`, derived from the `review.record` result row in
  `public-interface.md` section 4 (`validateCurrentResultJSON` appends
  `idempotency_key` and `replayed` itself); the matching
  `currentExpectedResult` branch that re-derives that result from the committed
  rows; and the `Review`, `Coverage`, `CoverageItem`, `CoverageExclusion`,
  `Finding`, and `FindingDisposition` entries in `currentValueShapes` from
  `public-interface.md` section 3. `Review` and `Finding` are composite, so add
  `currentProjectedValue` fill helpers for `Review.findings[]` /
  `Review.coverage` and `Finding.dispositions[]` in the same style as
  `currentFillRun` and `currentFillStopAttention`. This task owns
  `FindingDisposition` even though task 5 is the disposition task, because
  `review.recorded`'s own rule already inserts `finding_dispositions` rows for
  `resolves_finding_ids`; task 5 then reuses both shapes unchanged.
  - Also in `current_event_json.go`, and distinct from the `fields` map above:
    register this family's four `*_json` columns in the `effect.Table +
    "." + column` switch at `:14`, which otherwise returns `unregistered
    current JSON column review_records.repository_json` on the very first
    Review. `review_records.repository_json` is `Repository`,
    `coverage_json` is `Coverage`, `evidence_refs_json` is `[]Source`,
    `receipt_ids_json` is `[]ID`, and `review_findings.scope_json` is the
    Finding `scope` object, all from `public-interface.md` section 3
    (DEC-232).
  - Do **not** widen the `Subject` shape. DEC-231 directed that and DEC-232
    withdraws it: `Subject` is referenced by the live `Condition`, `Stop`, and
    `Attention` shapes, and the Stop and judgment domains are not nested.
    Instead register a new `JudgmentSubject` shape,
    `type:=issue,envelope,run,review,finding,requirement id:ID`, from
    `public-interface.md` section 3, and have `Finding` reference it.
    `currentValueShapes["Subject"]` must be byte-identical before and after
    this task.
- `internal/journal/current_event_cardinality.go` - `validateCurrentCardinality`
  defaults to one Effect per `table/op` key and widens it only through a closed
  `switch effect.Table` that lists no Phase 23 table, so a Review with two
  Findings is rejected with `current Event review.recorded: too many
  review_findings/insert Effects` (R-17, DEC-232). Admit repeated
  `review_findings/insert` and `finding_dispositions/insert` **for
  `review.recorded` only**: `public-interface.md` section 4 accepts
  `findings:FindingInput[0..64]` and `resolves_finding_ids:ID[0..64]`, while
  `finding.addressed` inserts exactly one disposition and must keep its
  implicit maximum of one, so condition the widening on the Event type in the
  style the existing `case "stops", ...` branch already uses. This is the only
  task that needs an entry in this file; every other Phase 23 family writes one
  row per key, and `requirement.referenced`'s insert and update are distinct
  keys.
- `internal/cli/codes.go` - `review-not-independent`, `review-verdict-invalid`,
  `repository-mismatch`, `review-not-found`, and `idempotency-conflict`
  (I-21). The last is an accepted Phase 21 code with no `reasonCodeDefinitions`
  entry and no adapter mapping today, while VAL-019 and VAL-057 both require a
  real-process reproduction of it; registering it is conformance inside this
  task's declared scope, not a new code.
- `cmd/fid/review.go` - new leaf wiring for `record` and `get`.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-4-done.md`.

### Files To Read For Context

- `contracts.md` sections 3, 4, 5; `public-interface.md` sections 2, 3, 4, 4.1,
  5; `schema.md` sections 2, 6.
- `task-3-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Finding dispositions, Evidence linking, budgets, freshness, Human text
  beyond the minimum required to return a result, and `review list`.

## Acceptance Criteria

- [ ] A recorded Review is immutable and every bound field matches the request
  and fid's own observation. (VAL-014)
- [ ] A conflicted Actor rejects `review-not-independent` naming the conflicting
  subject; a genuinely independent Actor succeeds. (VAL-015)
- [ ] A stale supplied repository identity rejects `repository-mismatch` with
  expected and observed values and no Event. (VAL-016)
- [ ] All three invalid verdict and Finding combinations reject
  `review-verdict-invalid`; all three valid ones succeed. (VAL-017)
- [ ] A second `review.record` for the same Run rejects with no Event. (VAL-018)
- [ ] Replay returns the original typed result with `replayed:true` and no
  second Event; a different digest returns `idempotency-conflict`. (VAL-019)
- [ ] A Review carrying more than one Finding, and one resolving more than one
  earlier Finding, both commit in a single Event and project every Finding and
  disposition; `finding.addressed` still admits exactly one
  `finding_dispositions/insert`. (VAL-017, VAL-020; R-17, DEC-232)
- [ ] `currentValueShapes["Subject"]` is unchanged by this task, and `Finding`
  resolves through the new `JudgmentSubject` shape. (DEC-232)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestReviewRecordRejectsNonIndependentActor` fails against the absent
`review.record` operation. Capture the failure before product edits; a build
failure alone does not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-4-done.md` in the full format,
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
- Commit using Conventional Commits with `Refs: phase-23/task-4`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
