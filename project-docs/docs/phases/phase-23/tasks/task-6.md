# Phase 23 - Task 6: Review evidence, review cycles, and rework exhaustion

## Objective

Make a passing Review usable as Evidence for a `review` Obligation, charge
exactly one `review-cycles` sample per `changes-required` Review, and refuse the
next rework atomically when that budget is exhausted.

## Validation Contract Coverage

- VAL-024: `evidence.link` accepts a `pass` Review and rejects
  `changes-required` and `blocked`.
- VAL-025: An evidence type outside the Obligation's accepted list rejects with
  exactly one code in the accepted precedence: the two Human Outcome mismatches
  return `outcome-not-permitted` and the other four return
  `evidence-type-not-accepted` (DEC-226). This task supplies receipt/review
  positive links and all six mismatch classifications; the matching Human
  Outcome durable link and real-producer closure are Task 10 (DEC-243).
- VAL-077: Linking a Review whose recording Actor has since taken a mutating
  Claim on that Issue and revision rejects `review-not-independent` at link time.
- VAL-008: the leaf re-check - `evidence.evaluate` with a null Obligation ID
  returns one evaluation per Obligation in ID order - lands here, because this
  task owns the leaf's optional Obligation ID.
- VAL-026: An unresolved blocking Finding keeps the Obligation unsatisfied even
  after a later `pass` Review.
- VAL-027: Exactly one `review-cycles` charge per `changes-required` Review,
  keyed on the Review ID, never double-charged.
- VAL-028: `run.observe-budget` rejects the fid-owned `review-cycles` unit.
- VAL-029: Exhaustion commits the Review, charge, Stop, and Attention
  atomically and refuses the next rework.
- VAL-030: `review-cycles` and `review-rework-exhausted` require each other and
  a `review-cycles` budget requires a `review` Obligation.

## Binding Constraints

- `contracts.md` sections 6 and 13 and kit dispositions C-03 and C-04 are the
  exact authority. Only `review-cycles` and `review-rework-exhausted` are added;
  generic `retry-exhausted` stays Phase 24.
- The charge source sample identity is the Review ID, so replay, linked
  recovery, and import cannot charge the same Review twice. Use the existing
  `budget_samples` uniqueness; add no new table or column.
- One `changes-required` Review is one cycle. Never derive a charge from a
  Finding count, a rework Claim, or a Run.
- Admission of the last allowed cycle is atomic with charging it and raising its
  Stop and Attention. The Review still commits; the Stop blocks the next rework.
- The exhaustion Stop uses code `review-rework-exhausted`, subject Envelope,
  resolver capability `stop.resolve`, blocked actions exactly `claim.acquire`
  and `run.start`, and resolution condition `budget-observed` on that budget.
- A `review` Obligation with no `review-cycles` budget is valid: unlimited
  rework is a recorded policy choice, not an error.
- `evidence.link` appends one immutable record per (Obligation, subject) pair
  and never overwrites prior proof. `evidence.evaluate` gains its optional
  Obligation ID here.
- Do not implement staling or positive Human Outcome evidence linking; those
  are tasks 7 and 10. DEC-243 permits the production type-admission classifier
  and the two Human Outcome mismatch denials here. A pure or seeded fixture
  may exercise those cells before the Human Outcome writer exists; identify
  its evidence class and never claim it as a real recorded Outcome or a
  successful Human Outcome link. Add no temporary Human Outcome writer.
- `execution-scope.md` is incorporated on the usual convergence terms.

DEC-243 cell accounting, with otherwise qualifying subjects and trust:

| Evidence type | gate-evidence Obligation | review Obligation | human-outcome Obligation |
|---|---|---|---|
| receipt | durable link here | evidence-type-not-accepted | evidence-type-not-accepted |
| review | evidence-type-not-accepted | durable link here | evidence-type-not-accepted |
| human-outcome | outcome-not-permitted | outcome-not-permitted | type-compatible here; real durable link at Task 10 |

## Scope

### Files To Create Or Modify

- `internal/evidence/evidence.go`, `policy.go`, `execution.go`.
- `internal/claims/generic.go` and matching in-package tests - enforce
  VAL-029's `review-rework-exhausted` admission denial inside the locked
  `claim.acquire` domain operation before any Claim Event or projection effect
  (DEC-244). The CLI adapter renders the domain denial; it is not the admission
  boundary. Retain existing Claim precedence and other Stop behavior.
- `internal/workflow/budget.go`, `run_budget.go`, `recovery_budget.go`,
  `evaluation.go`, `resolve.go`.
- `internal/workflow/envelope.go` - the authoring-validation path for VAL-030
  (DEC-241). Replace the existing `budget-exhausted`-only admission with the
  accepted exact unit/action pairing and require a `review` Obligation for a
  `review-cycles` budget; retain acceptance of an unbudgeted `review`
  Obligation and all other budget validation. Matching ordinary in-package
  tests exercise all five VAL-030 cases. This scope is executed by both the
  existing full race/coverage gate and the named `./internal/workflow/` gate;
  no command change or unrelated Envelope redesign is authorized.
- `internal/review/` - the verdict effects that charge and raise.
- **No `current_event_catalog.go` edit in this task** (I-28, DEC-229): task 1
  registered `review.recorded`'s complete `schema.md` 6.2 rule including every
  optional row this task's charging and raising paths use. If a needed optional
  row is missing, that is a task 1 defect - stop and escalate rather than
  editing a map task 1 owns. This task also registers no typed-result grammar:
  task 4 unfenced `review.recorded` and owns all three of its result registries
  (DEC-231).
- `internal/workflow/stop.go` - the closed Stop registries (R-14). Add
  `review-rework-exhausted` to `stopCodes` and its exact `validStopDetails` key
  set from the `public-interface.md` section 2 StopDetails table. Add no
  Condition type here: the resolution condition is the retained
  `budget-observed`.
- `internal/journal/current_event_value_catalog.go` - two additions. First, the
  matching `currentStopDetailShapes` entry (R-14). `StopCode` validation keys on
  that map (`current_event_value_types.go`), so without the entry the journal
  rejects the Stop this task raises. Second, widen
  `currentValueShapes["Budget"]` and `["ImmutableBudget"]` (R-16, DEC-231): add
  `review-cycles` to the `unit` domain and `review-rework-exhausted` to the
  `exhaustion_action` domain, per `public-interface.md` section 2. Without both,
  the first `review-cycles` Budget projected inside a `Run` or a
  `run.budget-observed` result is rejected at append time. No other Budget rule
  changes.
- **No journal-side budget-fact validation exists for this charge, and none is
  required of this task** (`execution-scope.md` section 6, inventory row 9;
  DEC-234, from QA R-21). `validateCurrentBudgetFacts`
  (`internal/journal/current_event_budget_values.go:34-50`) switches on
  `event.EventType` and returns `nil` at its `default` for all nine Phase 23
  families - measured, `review.recorded` included - so the sample/state
  consistency check, the automatic-charge rule, and the exhaustion-Stop
  derivation the journal enforces for `attempts` and `interventions` do **not**
  apply to `review-cycles`. The registry fails open, so nothing here is
  blocked. The consequence for this task is evidential: **VAL-027 and VAL-029
  are satisfied by producer-side enforcement only**, and the completion report
  must not present their evidence as journal-layer proof of the charge. If you
  conclude the journal-side guarantee is needed after all, extend the switch and
  say so explicitly in the report; do not assume it already holds.
- `internal/cli/codes.go` - `review-changes-required`,
  `review-rework-exhausted`, `evidence-type-not-accepted`, and the already
  accepted `outcome-not-permitted` for the two Human Outcome mismatch denials
  (DEC-243; Task 10 reuses this registration for the positive-Outcome path and
  its further denials). `finding-unresolved`
  is a ReasonCode, **not** a failure code (DEC-226): add it, together with
  `human-outcome-absent`, `review-not-independent`, and `outcome-superseded`, to
  the evaluation reason precedence order in `internal/workflow/evaluation.go`
  immediately after `obligation-evidence-absent`, in that order.
- `cmd/fid/evidence.go` and matching tests - the accepted typed Evidence link
  request/result, optional-Obligation evaluation result, and denial adapters
  owned by this task (VAL-008/024/025/077). `cmd/fid/run.go` and `claim.go`
  own denial rendering only; domain operations enforce admission.
- `cmd/fid/review.go` and matching existing tests - replace Task 4's
  `budget:null`, empty `stops`, and empty `attention` placeholders with the
  domain result produced by this task's atomic charge/exhaustion operation
  (DEC-245). Preserve the accepted Review result shape and replay semantics.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-6-done.md`.

### Files To Read For Context

- `contracts.md` sections 2, 4, 5, 6, 13; `public-interface.md` sections 2, 4, 5;
  `schema.md` sections 1, 6.2, 6.3.
- `project-docs/docs/phases/phase-22/contracts.md` section 6 for the retained
  budget arithmetic and denial rules.
- `task-5-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Freshness and staling, Human Outcome recording and positive evidence links,
  protected Decisions, and Human text. The two Human Outcome type-mismatch
  denials and production classifier evidence above are in scope (DEC-243).

## Acceptance Criteria

- [ ] `pass` links and can satisfy the Obligation; `changes-required` and
  `blocked` reject `review-changes-required` with the blocking Finding IDs.
  (VAL-024)
- [ ] The receipt/gate-evidence and review/review pairs link durably. All six
  mismatches return the exact DEC-226 codes: Human Outcome against gate-evidence
  or review returns `outcome-not-permitted`; the other four return
  `evidence-type-not-accepted` with the accepted list. Exercise the actual
  production admission path/classifier, never a test-only matrix. For Human
  cells, a named pure or seeded fixture is permitted and must be labelled as
  intermediate evidence. A nine-cell classifier acceptance for the matching
  Human pair is not a successful durable link. (VAL-025 task-6 evidence;
  DEC-243; final PASS also requires Task 10)
- [ ] Report each VAL-025 cell's expected code/admission result, observed
  result, evidence route and producing test. Observe no Event/projection change
  for actual rejected mutation attempts. Explicitly leave the positive Human
  link and real-producer Human-row closure to Task 10; do not mark the full
  phase assertion PASS from classifier or seeded-history evidence. (DEC-243)
- [ ] A blocking unresolved Finding keeps the Obligation unsatisfied with reason
  `finding-unresolved` even after a later `pass` Review. (VAL-026)
- [ ] Consumption rises by exactly one per distinct `changes-required` Review
  and by zero on replay, recovery, and import. (VAL-027)
- [ ] A harness `run.observe-budget` against `review-cycles` rejects with no
  sample. (VAL-028)
- [ ] The exhausting Review, its charge, its Stop, and its Attention commit in
  one Event, and the next `claim.acquire` and `run.start` return
  `review-rework-exhausted` with no Claim or Run. (VAL-029)
- [ ] The three invalid budget and Obligation pairings reject; the matched pair
  and the unbudgeted `review` Obligation are accepted. (VAL-030)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestChangesRequiredReviewChargesExactlyOneReviewCycle` fails against the
absent unit and charging path.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-6-done.md` in the full format.
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
go test -race -timeout 30m -cover ./internal/evidence/
go test -race -timeout 30m -cover ./internal/workflow/
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

- `go test -race -timeout 30m -cover ./internal/evidence/` and `go test -race -timeout 30m -cover ./internal/workflow/` before and after the edits, because this task adds branch-heavy evaluation and budget paths to two existing packages.

## Notes For Developer

- This task belongs to QA Review window 2.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-6`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
