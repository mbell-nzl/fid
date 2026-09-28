# Phase 23 - Task 13: Phase 23 CLI surface, failure codes, and idempotency

## Objective

Complete the real-process CLI and JSON surface for every Phase 23 operation,
land every added failure code with its exact shape, and make every Phase 23
mutation replay its original typed snapshot, including committed denials.

## Validation Contract Coverage

- VAL-055: Every Phase 23 leaf is reachable through the real process with
  canonical JSON and closed request validation, and each of the ten added
  ActivityEvent names is emitted by its producing leaf with no name outside the
  accepted set.
- VAL-056: Every added failure code returns its exact category, exit, subject,
  and closed details from the producing operation named in
  `public-interface.md` section 5.1, and the two precedence probes return the
  single code the accepted order fixes (DEC-226). `finding-unresolved` is a
  ReasonCode, not a failure code, and is not registered in `internal/cli`.
- VAL-057: Every Phase 23 mutation is idempotent and replays its original typed
  snapshot; every Stop-blocked denial commits nothing - no Event, no idempotency
  result, null `EventRef`, `replayed:false` - and returns the same immutable Stop
  and Attention IDs on repeat (DEC-226). `idempotency-conflict` must be
  reproduced through the built binary, never claimed from a fixture inventory.

## Binding Constraints

- `GLOSSARY.md` (DEC-226) distinguishes a **Failure code** (names one refused
  operation; exactly one per failed response; carries category, retryability,
  and exit status) from an **Evaluation reason** (describes why one Obligation
  is unsatisfied; zero or more; no category, retryability, or exit status).
  Register an identifier in the set it belongs to. An identifier may belong to
  both sets only where the same fact both refuses an action and describes a
  state, as `protected-decision-required` does.

- `public-interface.md` sections 1, 4, and 5 are the exact authority. Add no
  leaf, flag, field, or code that is not listed there.
- Every mutation takes `--input`, `--actor`, and `--idempotency-key` and no
  positional argument. Every read uses its declared read grammar. All request
  objects stay closed recursively.
- List leaves keep default limit 100, accepted range 1..10000, and
  `result-limit-exceeded` on a bounded complete projection.
- DEC-216: store and replay the complete typed snapshot for each mutation.
  A committed `review-rework-exhausted` or `protected-decision-required` denial
  replays its original Budget or protected-Decision state, Stop and Attention
  IDs, and EventRef without re-deriving current state. A denial that commits
  nothing has a null EventRef.
- Reuse existing adapters and error rendering; do not create a second response
  construction path or a generic formatter.
- Human text obligations belong to task 14. This task may emit minimal Human
  output, but must not pin fixtures that task 14 will immediately replace.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `cmd/fid/review.go`, `finding.go`, `concern.go`, `deviation.go`,
  `decision.go`, `human_outcome.go`, `requirement.go`, `evidence.go`,
  `envelope.go`, `claim.go`, `root.go`.
- `internal/cli/codes.go`, `internal/cli` response and failure adapters.
- Vocabulary rename (DEC-226): `reasonCodeDefinitions` names the **failure
  code** registry, not the evaluation-reason set, and that collision is what let
  `finding-unresolved` be filed as a failure code. Rename it to
  `failureCodeDefinitions` at its declaration and its three other uses
  (`internal/cli/codes.go`, `internal/cli/envelope.go:38`), plus the test-local
  `reasonCodeExpectation` in `internal/cli/codes_test.go`. Rename nothing in
  `cmd/fid/evidence_test.go`: its `reasonCodes` is the `assertEvaluation`
  parameter naming evaluation reasons, which is the correct sense (I-20). This is a pure
  identifier rename of unexported names: no behavior, no JSON key, and no
  public value changes. `GLOSSARY.md` defines both terms.
- Idempotency snapshot writers in the owning domain packages.
- Matching in-package tests and `cmd/fid/integration/` fixtures.
- `project-docs/docs/phases/phase-23/tasks/task-13-done.md`.

### Files To Read For Context

- `public-interface.md` sections 1, 2, 3, 4, 5; `contracts.md` section 13;
  `schema.md` section 6.3.
- `task-12-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Human text content, resume projections, capability advertisement, and the
  final conformance audit.

## Acceptance Criteria

- [ ] Every Phase 23 leaf returns its declared result through the real process,
  and unknown-field, wrong-nullability, trailing-JSON, and duplicate-key
  negatives reject with no Event. (VAL-055)
- [ ] Every added code matches its accepted category, exit, subject, and closed
  details with no extra key and no reused generic code. (VAL-056)
- [ ] Every mutation replays its original typed result with `replayed:true` and
  no second Event, and both committed denials replay their original Stop,
  Attention, and EventRef. (VAL-057)
- [ ] The failure-code registry and its uses are renamed to name the failure
  set, with no behavior, JSON key, or public value change. (DEC-226)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestPhase23LeavesRejectUnknownRequestFields` fails against the incomplete leaf
surface.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` section 5: record the full Phase 23
  closed code additions accepted by the slice contract.
- `project-docs/docs/phases/phase-23/tasks/task-13-done.md` in the full format.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
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

- This task belongs to QA Review window 5.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-13`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

## Planning Amendment - DEC-249 (2026-09-26)

- In VAL-056's second precedence probe, the unsatisfied Obligation that
  accompanies the open readiness Stop must guard the probed action. An
  Obligation that does not guard it no longer contributes a reason.
- **First step (DEC-254 Part 2, Human-approved):** make `cmd/fid` time
  headroom. This step is test-only and committed separately, before any new
  task-13 test is added.
  - Scope: existing `cmd/fid/*_test.go` files, for `t.Parallel()` on tests
    shown to share no process-global state (cwd, env, package vars), and for
    fixture or binary reuse.
  - Target: 480 s or less on
    `go test -race -timeout 30m -coverprofile=coverage.out ./...` for the
    `cmd/fid` package.
  - Not allowed: removing, weakening, skipping or relocating any assertion.
  - Then run `go test -race -timeout 30m -count=2 ./cmd/fid/` once as a flake
    guard.
  - Report the before and after times. If the target is missed, record the
    achieved time and escalate before adding new `cmd/fid` tests.
- DEC-254 Part 1's breach-attribution rule still applies to the rest of the
  task.

## Planning Amendment - DEC-258 (2026-09-27)

- The DEC-254 Part 2 first step above is **already satisfied** by the
  ESC-23-005 repair. Do not repeat it. Record `cmd/fid`'s time from this task's
  own exact gate as usual; `qa/esc-23-005-repair.md` holds the before and after
  times.
- New top-level tests call `t.Parallel()` unless they touch process-global
  state; `golangci-lint`'s `paralleltest` enforces this. A long subtest matrix
  whose cases build independent workspaces parallelizes its subtests.

## Planning Amendment - DEC-262 (2026-09-28, ESC-23-007)

The `activity_log` gap is a planning defect. It was not caused by this task.
Continue as attempt 2 of `task-13:scope-gap:review-list-and-activity-producers`.

- **Scope added:**
  - `internal/journal/current_event_columns.go`: add the `activity_log`
    entry, derived from `migrations/0002` DDL and proven by
    `TestCurrentEventColumnDomainsMatchAcceptedDDL`.
  - If the `activity` Effect reaches them on the strict path, also:
    - `current_event_cardinality.go`: `review.recorded` alone may carry more
      than one `activity_log/activity`;
    - the `current_event_json.go` column switch, for `activity_log.details`;
    - any other `internal/journal/current_event_*.go` check that rejects it.

    This is `execution-scope.md` section 6 row 11.
  - `internal/review/review.go`, `finding.go`, `concern.go`, `deviation.go`,
    and `human_outcome.go`, which are the activity producers.
  - Matching tests.
- **Emit** the `issue`-subject rows in `schema.md` section 6.2:
  - `review-recorded`;
  - `finding-resolved` (one per `resolved` disposition; none for
    `inapplicable`);
  - `finding-addressed`;
  - `human-outcome-recorded`;
  - `concern-recorded`, `concern-disposed`, and `deviation-recorded` for
    `issue`, `envelope`, `run`, `review`, and `finding` subjects.
- **Do not emit** any `requirement`- or `decision`-subject row. Under the
  0016 trigger it aborts the Event. Leave `decision.go` and `requirement.go`
  unchanged.
  - The following rows are **task 15's**:
    - `requirement-referenced`;
    - `requirement-linked`;
    - `decision-referenced`;
    - the `requirement`-subject Concern/Deviation rows.
  - Report them under VAL-055 as "deferred to task 15 (DEC-262)", never as
    `PASS` and never as `BLOCKED`.
  - The activity-log inventory at this task asserts the seven `issue`-subject
    names. It also asserts that no name outside the accepted ten appears.
- **Forbidden:**
  - loosening any journal check;
  - editing any effect rule, stage, or `semanticEventRules` row;
  - editing migrations or triggers;
  - editing ActivityEvent names;
  - using an invented Issue or a V1 row.

  All retained writers already declare `activity_log/activity`. If one is
  rejected once the table is registered, or a registry outside
  `internal/journal` rejects, return `RECOVERY_REQUIRED`. Do not edit the
  rule.
- **Retain** the failed trial, the zero-row probes, the failed tagged gate,
  and the failed report-lint in `task-13-done.md`. Disclose DEC-262 under
  Deviations. Then complete VAL-056 and VAL-057 as specified and run all five
  exact commands unchanged.
- **Acceptance criterion added:** the seven `issue`-subject names appear from
  their producing leaves with the `schema.md` section 6.2 subject, and no
  `requirement` or `decision` activity row is written. (VAL-055, DEC-262)

## Planning Amendment - DEC-263 (2026-09-29, ESC-23-008)

The contradiction is a drafting defect in the accepted set, not a task-13
defect. Continue as attempt 2 of
`task-13:protected-contract:protected-decision-denial-envelope-val-057`.

- **Accepted shape** (`public-interface.md` sections 5 and 5.2). Put the four
  denial fields inside `details`, in the same form as
  `review-rework-exhausted`. The envelope does not change.
  - `protected-decision-required`:
    `{envelope_id,decision_class,required_capabilities,affected_actions,stop_ids,attention_ids,event,replayed}`
  - `human-outcome-required`:
    `{envelope_id,obligation_id,outcome_class,stop_ids,attention_ids,event,replayed}`
  - Values:
    - `stop_ids`: the blocking Stops the details describe, sorted by ID.
    - `attention_ids`: the Attention linked to those Stops, sorted by ID.
    - `event`: always `null`.
    - `replayed`: always `false`.
- **Scope added:**
  - `internal/cli/codes.go`: extend both definitions with the existing
    `stop_ids`/`attention_ids`/`event`/`replayed` field helpers.
  - `cmd/fid/decision.go` `protectedDecisionFailure` and
    `cmd/fid/human_outcome.go` `humanOutcomeRequiredFailure`: add the four
    fields from the `StopBlockedError` the adapters already receive (`Stops`
    and `AttentionIDs`). Do not create a second response path.
  - `cmd/fid/issue_mutation_response.go` and `cmd/fid/envelope.go`: edit
    only if a caller must pass that `StopBlockedError` through unchanged.
  - Matching tests in `internal/cli`, `cmd/fid`, and `cmd/fid/integration`.
- **Converge the retained exact-details rows** that pin the old four-key or
  three-key shapes. These are rows from tasks 9, 10, 11, and 13, for example
  `assertP23ProtectedDecisionDenial` and the `human-outcome-required` and
  precedence probes.
  - Keep every existing key and value check.
  - Add `stop_ids` and `attention_ids` equal to the durable Stop and its
    linked Attention, plus `event:null` and `replayed:false`.
  - Record each row's provenance (introducing commit and task) and its
    executing command in `task-13-done.md`.
- **Strengthen VAL-057.** The protected-Decision repeat-denial probe must
  also assert that `stop_ids` and `attention_ids` equal the durable blocking
  Stop and its Attention. A byte-identical repeat alone is not enough.
- **Forbidden:**
  - editing StopDetails, Stop rows, or
    `internal/journal/current_event_value_catalog.go`;
  - editing the StopDetails keys in `internal/workflow/stop.go`;
  - changing any category, exit, subject, producer, or precedence;
  - changing the `fid-cli.v1` envelope;
  - removing or weakening any assertion;
  - changing any command.

  If any of these is needed, return `RECOVERY_REQUIRED`.
- **Retain** the attempt-2 BLOCKED report, the failing
  `TestP23DecisionReferenceProtectedLifecycleRealProcess` run, and the failed
  tagged gate in `task-13-done.md`. Disclose DEC-263 under Deviations. Then
  run the focused command from `qa/resume-after-esc-23-008.md` and all five
  exact commands unchanged.
- **Acceptance criterion added:** both Stop-blocked codes return their
  DEC-263 details exactly, with the IDs equal to the durable Stop and
  Attention, and no StopDetails or Stop projection changes. (VAL-056,
  VAL-057, DEC-263)
