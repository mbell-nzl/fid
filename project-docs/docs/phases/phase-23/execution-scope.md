# Phase 23 Execution Scope And Gate Reach

> Tech Lead planning inspection, 2026-09-14. ACCEPTED under DEC-222.
> No build, test, lint, or feasibility gate was run by Tech Lead.
> Base `0195897a347dfe1872db0f7be3800c5c3ccdefcf`.
> Amended 2026-09-15: section 7 adds the test-time headroom disposition
> (DEC-227, R-13), and three file scopes are corrected under DEC-167/173
> (R-08, R-09) in the affected task specs.
> Amended 2026-09-19: section 6 gains the closed-registry owners (DEC-229),
> the unfencing-task standing rule (DEC-231), and its generalisation with the
> six-row strict-path registry inventory (DEC-232, R-17/R-18/I-31).

## 1. Ownership Rule

Every task owns its named implementation paths, the matching in-package tests,
and only those retained assertions or fixtures invalidated by its own accepted
change. The default gate compiles every `cmd/` and `internal/` package; its
required tagged companion also builds the CLI and runs real-process journeys.
Accordingly `cmd/`, `internal/`, their `testdata/`, root `testdata/`, and the
accepted migration chain are in scope for minimum causal convergence only, not
for unrelated behavior, cleanup, or redesign. This is the task pack's explicit
scope for those gate roots, never a waiver of narrow task ownership.

For each displaced retained row, the completion report records file and line,
introducing commit and phase/task, the accepted invalidating change, how the
row's subject is preserved or replaced, and the exact executing command. Labels
and filenames never assign ownership. A latent production defect first observed
by a new gate follows DEC-175: the Runner obtains the minimum causal scope
amendment before repair. This document grants no unrestricted repair authority.

## 2. Planning Inspection

Read-only evidence used to design the commands below:

- `cmd/fid/integration/main_test.go` carries `//go:build integration`; its
  `TestMain` builds `./cmd/fid` into fresh temporary storage and runs the
  tagged suite. Plain `go test ./...` excludes those files entirely.
- `internal/db/migrate.go` loads migrations from an explicit directory, so
  database and fresh-init tests exercise the actual SQL rather than a
  duplicated test-only schema. Every `schema.md` DDL row is therefore
  executable evidence, not documentation.
- `internal/journal/rebuild_bench_test.go` exposes
  `BenchmarkRebuildWorkspace1kEvents`, `...10kEvents`, and `...100kEvents`;
  the 100k case self-fails at 30 seconds or more. Tasks 1, 15, and 16 preserve
  that command and its accepted contract.
- Live packages at the base are `activity`, `bootstrap`, `claims`, `cli`,
  `config`, `db`, `evidence`, `export`, `issues`, `journal`, `labels`, `links`,
  `maintenance`, `policy`, `project`, `requirements`, `sqlquery`, `stats`,
  `validation`, and `workflow`. This phase adds exactly one, `internal/review`,
  and removes exactly one, `internal/requirements`.
- Hosted CI (`.github/workflows/ci.yml`) runs format, lint, build, the shard
  verifier, and four exhaustive ordinary race-and-coverage shards. It does
  **not** run the tagged integration suite or any benchmark. Those remain
  mandatory local task and phase gates and are never substituted by a green PR.
  Closing that hosted gap is workflow maintenance recorded as CF-041, not
  Phase 23 product scope.

This inspection confirms command selection and tag reach, not runtime success.
QA Lead performs fresh feasibility validation in Heavy mode before Ready.

## 3. Task-To-Gate Matrix

| Tasks | Owned behavior | Executing gate |
|---|---|---|
| 1, 15, 16 | Migration, DDL, rebuild, reset, current schema, portable schema | ordinary race and coverage, full tagged integration, plus the named fresh-process rebuild benchmark |
| 2, 3, 5, 6, 7, 8, 9, 10, 11 | Domain rules, Effects, SQLite projections, authority and failure branches | ordinary race and coverage plus full tagged integration; the tagged gate catches displaced adapter fixtures |
| 4, 13, 14, 17, 18 | CLI, JSON, Human text, shared semantic fixtures, process journeys, availability matrix | both gates; tagged integration is mandatory and never substituted |
| 12 | Requirement reference replacement and link-model removal | both gates plus the fresh-process rebuild benchmark, because it removes an Event family from the replay corpus, plus its own declared expected-non-zero absence search |
| 2, 6 | New branch coverage in `internal/workflow` and `internal/evidence` | named package `go test -race -cover` before and after edits |
| 15 | V1 removal absence proof | both gates plus the declared expected-non-zero absence searches in that task spec |

The tag selection is case-sensitive and remains exactly
`go test -tags integration -race -timeout 30m ./cmd/fid/integration/...`
(the `-timeout` is the DEC-227 bound in section 7; the tag selection itself is
unchanged). No other environment-gated suite is assigned. New helpers and fake
adapters must be used by real assertions; compiling an unused helper is not
behavioral evidence.

Every task runs the tagged integration gate from task 1. There is no delayed
first integration gate and therefore no planned convergence debt. If
implementation reveals a build tag or environment switch outside these
commands, amend the task before assigning that scope; never narrow a gate.

## 4. Ordered Delivery And Fence Ownership

`schema.md` section 6.1 assigns each progressive fence key to exactly one task:
task 4 `review.recorded`; task 5 `finding.addressed`; task 8
`concern.recorded`, `concern.disposed`, `deviation.recorded`; task 9
`decision.referenced`; task 10 `human-outcome.recorded`; task 12
`requirement.referenced`, `requirement.linked` plus the now-empty map and its
guard branch. Every other key stays fenced. Task 1's fence assertion is
scheduled convergence for each later owner, substituting a still-fenced Event
and preserving the parity subject; it is never deleted.

`public-interface.md` section 4.1 assigns each new grant action to its first
owning task. Adding an action before its owning task, or leaving one
unimplemented at task 18, both fail the final audit.

## 5. Review And Command Accounting

Three-task review windows: 1-3, 4-6, 7-9, 10-12, 13-15, 16-18. QA Reviewer
reviews each completed window.

An independent QA Lead checkpoint after task 11 must PASS before task 12. It
checks real Review independence, Finding resolution authority, review-cycle
charging and exhaustion, protected-Decision blocking and satisfaction, Human
Outcome provenance, and Attention resolution through a Human Outcome.

Acting on the Phase 22 workflow note that Human-readable guidance was checked
too late, the window 13-15 review additionally performs a deterministic
Human-text conformance check against `contracts.md` section 14, using the
verbatim fixtures task 14 records in
`project-docs/docs/phases/phase-23/qa/human-fixtures.md`. A Human-owned
validation assertion still covers product judgment in the human checklist.

Runner report lint is a post-report-commit, pre-next-dispatch control using
`project-docs/tools/report-lint.sh` with exactly two positional arguments: the
completed `task-N-done.md` followed by its `task-N.md`, both under
`project-docs/docs/phases/phase-23/tasks/`. Concrete N is 1..18. This is Runner
orchestration and never a Developer `Exact Commands` obligation. The first
concrete invocation is:

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-1-done.md project-docs/docs/phases/phase-23/tasks/task-1.md
```

Required commands report their own output token and observed termination. Tasks
12 and 15 declare expected non-zero commands; task 12's declaration was added
2026-09-15 with its `requirement.synced` absence search (R-07d). Red tests retain their failing
evidence and never authorize weakening a gate. No cleanup is a gate. Do not
impose a parent session lifetime as a command timeout. External interruption
retains no-progress state and follows current schema-v2 recovery.

## 6. Corrected File Scopes (2026-09-15, extended 2026-09-19)

Three assignments were incomplete under DEC-167/173 and are corrected in the
task specs before dispatch, not after a failure:

- The closed 23-name activated-capability and grant-action allowlist is
  `sortedUniqueCapabilities` in `internal/workflow/activation.go`; Envelope
  grants reuse it through `sortedUniqueEnvelopeActions` in `envelope.go`, and
  `authority.go` holds no action list. Tasks 4, 5, 8, 9, and 10 each add
  `internal/workflow/activation.go` to their file scope for their own new
  action (R-08).
- `internal/export/archive_semantics_catalog.go` carries the
  `requirement.synced` archive semantic rule. Task 12 removes the Event, so that
  file joins task 12's scope; without it, task 12's own declared absence search
  fails with its cause outside the task (R-09).
- Task 10 owns the readiness-Stop effect rows added to `evidence.linked` and
  `artifact.availability-observed` in `internal/journal/current_event_catalog.go`
  (DEC-224), alongside its fence key.

Four further closed registries had no owning task. Corrected under R-14
(DEC-229); same rule, same reason - a task owns only what it names, and a
registry nobody names is rejected at runtime with its cause outside the failing
task's scope:

- `internal/workflow/stop.go` holds the closed `stopCodes`, `conditionTypes`,
  `Condition`, and `validStopDetails` registries, and
  `internal/journal/current_event_value_catalog.go` holds their journal-side
  grammar (`currentStopDetailShapes`, which `StopCode` validation itself keys
  on, and the `Condition` value shape). Both files join the scope of every task
  that adds to them: task 6 (`review-rework-exhausted`), task 9
  (`protected-decision-required`, `decision-satisfied`, `decision_class`),
  task 10 (`human-outcome-required`, `human-outcome-recorded`,
  `outcome_class`), and task 11 (`review-recorded`, `finding-resolved`,
  `finding_id`). A Condition type is registered by the task that first raises a
  Stop carrying it, never deferred to task 11; task 11 completes the closed
  four-type set and owns the resolution behavior of all four (VAL-078).
- `internal/journal/current_event_columns.go` joins task 1, which registers all
  ten new Phase 23 tables in `currentEventTables` alongside the nine effect
  rules it already registers. `validateCurrentEffect` returns `unknown table`
  for an effect on an unregistered table, and
  `TestCurrentEventColumnDomainsMatchAcceptedDDL` proves the registration
  against the 0017 DDL in the same task that ships it. Task 15 renames the
  `judgment_concerns` key with its table; no other task edits this map.
- `internal/export/archive_semantics_catalog.go` mirrors the journal catalogs
  for `export validate` and `import`, and no test asserts parity between them.
  Standing rule: **a task that adds, amends, or removes a journal Event rule
  mirrors that same change in `semanticEventRules` in the same task.** Task 1
  mirrors the nine new rules, task 9 its two amendments, task 10 its two, task
  12 the `requirement.synced` deletion (R-09, already assigned), and task 15
  the concerns rename. Without this, a divergence stays silent until VAL-064 at
  task 16, where the cause sits outside that task's ordinary scope.

A fifth gap, larger than the four above, was found while confirming them
(R-16, DEC-231), a sixth and seventh on the same path while confirming that
(R-17 and the Tech Lead's own re-derivation, DEC-232), and an eighth and ninth
off that path while confirming *that* (R-19 and the Tech Lead's own
re-derivation, DEC-233). They share one cause, so they are stated once as a
general rule with a closed inventory:

- Registering the nine families in `currentEventRules` routes them onto the
  journal's **strict** validation path - `validateCurrentEvent`
  (`current_event_validation.go:26`) for the append checks and
  `validateCurrentAppliedResults` (`current_event_result_projection.go:12`,
  called from `project_effects.go:20`) for the typed snapshot, which
  `schema.md` 6.3 makes mandatory for every Phase 23 mutation. That path
  consults several closed registries that fail closed on an unregistered
  family, table, or column.
- **Frame (widened by I-33). The obligation is not "registries the strict path
  consults for a Phase 23 family". It is every closed journal registry that a
  Phase 23 value or Effect reaches, whichever family carries it.** There are
  exactly two entry points, and a sweep that walks only the first is incomplete:
  1. a **Phase 23 family** on the strict path - the six registries in the table
     below, owned by the task that unfences the family; and
  2. a **Phase 23 value inside a live or retained Phase 21/22 family** - rows 7
     and 8 below. R-19 is the proof that this second entry point exists: the
     capability allowlist is reached through `activation.component-staged`, a
     never-fenced Phase 21/22 family, and three consecutive sweeps missed it
     because all three walked only the nine families' own path.
- **Standing rule (generalised by I-31): a task that deletes a family's fence
  key owns that family's entry in every closed journal registry the strict path
  consults - never task 1.** That covers the family's typed-result grammar, its
  `currentExpectedResult` branch, its value shapes, any `currentProjectedValue`
  fill helper, its repeated-effect cardinality, its `*_json` column grammars,
  and its update field-and-guard sets. Every one of these executes only when an
  Event of the family is appended, so the unfencing task is the earliest point
  at which any of them is provable, and it is the point at which the task's own
  first successful Event proves them together. This is the same principle
  DEC-229 applied to Condition types, which are registered by the task that
  first raises a Stop carrying them. Task 1 is deliberately **not** the owner
  here, unlike `current_event_columns.go`: that catalog is proved by
  `TestCurrentEventColumnDomainsMatchAcceptedDDL`, which introspects the 0017
  DDL task 1 ships and needs no Event at all, whereas all of these would sit
  behind task 1's own fence unexecuted. Ownership is only assignable where it is
  executable.
- **The inventory known today, all under `internal/journal/`. A registry that
  is not on this list and rejects - or silently fails to constrain - a Phase 23
  value is a Developer escalation, not a silent workaround.** Rows 1-6 are
  reached by a Phase 23 family on the strict path. Rows 7 and 8 are reached by a
  Phase 23 **value** inside a live Phase 21/22 family and are the reason the
  frame above is stated the way it is.

| # | Registry | File | Rejection message | Phase 23 owners |
|---|---|---|---|---|
| 1 | `fields` grammar map in `validateCurrentResultJSON` | `current_event_json.go:135` | `unregistered current result family` | 4, 5, 8, 9, 10, 12 |
| 2 | `currentExpectedResult` switch | `current_event_result_projection.go:69` | `unknown current result` | 4, 5, 8, 9, 10, 12 |
| 3 | `currentValueShapes` | `current_event_value_catalog.go:6` | `missing projected <Kind>.<field>` | 4, 8, 9, 10, 12 |
| 4 | `validateCurrentCardinality` table switch | `current_event_cardinality.go:14` | `too many <table>/<op> Effects` | 4 only (R-17) |
| 5 | `validateCurrentJSON` `table.column` switch | `current_event_json.go:14` | `unregistered current JSON column <t>.<c>` | 4, 8, 9, 10, 12 |
| 6 | `validateCurrentUpdate` table switch | `current_event_updates.go:14` | `current Event update table not admitted` | 8, 12; renamed by 15 |
| 7 | `const actions` capability allowlist in `validateCurrentValueRelations` `case "ActorProfile"` | `current_event_value_relations.go:63` | `invalid closed enum artifact.verify,...` | 4, 5, 8, 9, 10 (R-19) |
| 8 | `case "Condition"` nullable-binding relation in the same function | `current_event_value_relations.go:48` | none - **fails open** | 9, 10, 11 (DEC-233) |

- **Row 7 (R-19), the second copy of the capability allowlist.** The closed
  23-name set exists **twice** in non-test code: `sortedUniqueCapabilities`
  (`internal/workflow/activation.go:832`) and this `const actions` string. They
  are byte-for-byte the same names and nothing cross-checks them. The journal
  copy is reached with no Phase 23 family involved at all:
  `activation_components.evidence_json` is validated as `[]ActorProfile` when
  `component_type` is `authority-profile` (`current_event_json.go:52-63`), on
  `activation.component-staged`, which is in `currentEventRules`
  (`current_event_catalog.go:28`) and behind no fence. `public-interface.md`
  section 4 requires each Phase 23 leaf's capability in the Actor profile, and a
  profile reaches the journal only through Kit Activation, so **the first
  Activation staging a Phase 23-capable profile is rejected before task 4
  records a single Review**. Each of tasks 4, 5, 8, 9, and 10 therefore adds
  `internal/journal/current_event_value_relations.go` beside the
  `internal/workflow/activation.go` entry R-08 already gives it, and adds its
  own name to **both** copies in one change. 23 becomes 30: exactly the seven
  grant actions of `public-interface.md` 4.1 and **not**
  `requirement.record-reference` or `requirement.link`, which stay operation IDs
  only. Task 18's audit compares the two copies against each other (VAL-067), so
  a future divergence fails a test instead of staying silent - the same
  protection DEC-229 gave `semanticEventRules`.
- **Row 8, the `Condition` nullable-binding relation. Tech Lead finding under
  DEC-233, same file, found while confirming R-19.** `case "Condition"`
  enforces the accepted "only the field matching the type is non-null" rule for
  the two Phase 22 fields and nothing else: `(type == "budget-observed") !=
  (budget_name != nil) || (type == "capability-available") != (capability_id !=
  nil)`. `public-interface.md` section 2 extends that rule to `outcome_class`,
  `decision_class`, and `finding_id`, and tasks 9, 10, and 11 already add each
  field to `validateCondition` in `internal/workflow/stop.go` - the **producer**
  side only. Reproduced under overlay with the shape widened and the working
  tree untouched: a `human-outcome-recorded` Condition with all three new fields
  non-null, one with none of them non-null, and a `decision-satisfied` Condition
  carrying `outcome_class` instead of `decision_class` are **all accepted**,
  while the Phase 22 control still returns `condition nullable binding
  mismatch`. This registry fails **open**, so it blocks nothing and is not a
  dispatch blocker; it is a silent loss of an accepted journal-layer invariant,
  the class R-09 and R-18 were accepted to prevent. Each of tasks 9, 10, and 11
  extends this branch for its own field in the same change that adds the field
  to `currentValueShapes["Condition"]` and to `validateCondition`. See ESC-23-001
  for the separate, blocking question of **how** that shape may be widened.

- `validateCurrentTerminalValues`'s second switch in that same file
  (`current_event_updates.go:61`) fails **open**, so it blocks nothing. The
  owning task may add a status-transition guard there and should say in its
  completion report whether it did; the required change is the admission switch
  above.
- Two registries on the path were checked and need **no** Phase 23 entry; this
  is a positive finding, recorded so no later pass re-derives it.
  `validateCurrentBindings`'s `entityTable` map
  (`current_event_validation.go:282`) and `currentResultEntity`'s map
  (`current_event_idempotency.go:65`) are both keyed by Event type and both
  fail **open** to the aggregate ID, which is the correct `result_entity_id`
  for all nine Phase 23 families. Do not add rows to either.
- **Six further checks on the widened frame returned no work. Recorded as
  positives under DEC-233 so the next pass does not re-derive them.**
  1. `validateCurrentValueRelations` reaches `return nil` for any kind not in
     its switch, verified under overlay for `Review`. The nine new Phase 23
     value kinds therefore need **no** relations entry, and rows 7 and 8 are the
     file's complete Phase 23 surface. That is what makes the file closable.
  2. `validateCurrentBlockedActions` (`current_event_json.go:117`) closes the
     six-name action set. `public-interface.md` section 2 says the Stop
     blocked-action subset "adds nothing", and all four names the three new Stop
     codes use are already in it. The `budget-` prefix relation
     (`current_event_value_relations.go:76`) keys on a prefix no Phase 23 code
     has. **No entry, either place.**
  3. `StopCode` is not a separate list: `current_event_value_types.go:88` keys
     it on `currentStopDetailShapes`, which tasks 6, 9, 10, and 11 already own.
     Registering the details shape registers the code. **No second edit.**
  4. Envelope grant actions reach no journal registry.
     `sortedUniqueEnvelopeActions` reuses the `activation.go` allowlist, and
     `envelope.created` / `envelope.activated` are in `retainedCurrentRules`,
     not `currentEventRules`. The retained branch returns at
     `current_event_validation.go:64` before any JSON or value validation. Row 7
     is the **only** second copy.
  5. Export needs no value-side registration. `validateSemanticMap`
     (`archive_semantics_values.go:10`) validates by DDL column type and checks
     `_json` columns only for well-formed JSON, and `loadSemanticTables`
     (`archive_semantics.go:125`) introspects the installed schema, so the ten
     new tables and every new value arrive from migration 0017 for free. The
     `semanticEventRules` mirror DEC-229 already assigns is the whole export-side
     obligation - with one exception, next.
  6. **Exception, and the one piece of new work from this sweep:
     `semanticPhase22Events` (`archive_semantics_catalog.go:85`) has no Phase 23
     analogue.** `validateEvent` (`archive_semantics.go:328`) rejects an Event
     when `!found || (profile == 14 && semanticPhase22Events[type])`. Once task 1
     mirrors the nine families into `semanticEventRules`, `found` becomes true
     for them at **every** profile, so a legacy projection-14 archive carrying
     `review.recorded` would validate clean - the exact hole Phase 22 added that
     map to close. Task 1 adds the nine families to it in the same file and the
     same change as the mirror, and the map's name is corrected to read
     "families that do not exist at projection 14". This is provable with no
     Event, like `current_event_columns.go`, so task 1 is the right owner.
     VAL-065 is unaffected and already holds by a different route: the manifest
     check `decoded.ProjectionSchema != fiddb.ActiveMigrationVersion`
     (`archive_validation.go:130`) rejects a cross-version archive with
     `schema-unsupported` before any Event is read.
- Owners, one-to-one with the nine fence keys: task 4 `review.recorded`; task 5
  `finding.addressed`; task 8 `concern.recorded`, `concern.disposed`,
  `deviation.recorded`; task 9 `decision.referenced`; task 10
  `human-outcome.recorded`; task 12 `requirement.referenced`,
  `requirement.linked`. A value shape is registered once, by the first of those
  tasks to need it, and reused unchanged afterwards - so task 5 adds no shape,
  and task 10 reuses task 9's `ProtectedDecision`.
- Two shared shapes are widened outside that set, each by the task that makes
  the wider domain reachable: `currentValueShapes["Claim"]`'s `mode:=mutation`
  becomes `mutation,review` at **task 3**, reachable because `run.abandoned` and
  `run.recovered` both project a `Claim`; and `["Budget"]` / `["ImmutableBudget"]`
  gain `review-cycles` in `unit` and `review-rework-exhausted` in
  `exhaustion_action` at **task 6**.
- **The journal `Subject` value shape is not widened by any task (R-18,
  DEC-232).** DEC-231 directed task 4 to widen
  `currentValueShapes["Subject"]`'s `type:=run,envelope,activation` "to the
  complete accepted SubjectType set"; that direction is withdrawn as wrong on
  three counts. `Subject` is referenced by the `Condition`, `Stop`, and
  `Attention` shapes, which serve `stop.raised`, `stop.resolved`,
  `attention.resolved`, `run.budget-observed`, and every `stops/insert` in any
  family - all live from task 1 and none behind the fence, so no widening of it
  is "monotonic behind the fence". The two domains are not nested either:
  `stops`/`attention_items` admit `activation` (`0015:56,82`) and the four
  judgment tables do not, while the judgment tables admit `issue`, `review`,
  `finding`, and `requirement` and Stop does not, so no single shape can serve
  both. And `public-interface.md` section 2's `SubjectType` row is the CLI
  **error** subject vocabulary - `internal/cli/codes.go:88`, which is why
  `requirement` can be "retained" there - not the Stop or judgment subject; a
  literal reading would let a Stop resolution condition name `type:"path"`.
- **Task 4 instead registers a separate `JudgmentSubject` shape**,
  `type:=issue,envelope,run,review,finding,requirement id:ID`, from the closed
  set already accepted in `contracts.md` section 8 and enforced identically by
  the `subject_type` CHECK on `review_findings`, `judgment_concerns`,
  `deviations`, and `human_outcomes`. `Finding` uses it at task 4; `Concern` and
  `Deviation` reuse it unchanged at task 8 and `HumanOutcome` at task 10.
  `Subject` stays exactly as it is, and no task edits it.
- Note for every owner: `idempotency_results` is encoded as an optional stage in
  `currentEventRules` for 15 of the 22 existing families and a required stage
  for the other 7 (I-30), but `validateCurrentAppliedResults` requires the
  result row unconditionally for all of them. "Optional" in `schema.md` 6.2 is
  the stage encoding, not permission to omit the typed snapshot.

Not corrected, and deliberately so: tasks 4 and 5 do **not** gain
`internal/journal/current_event_catalog.go`. Task 1's binding constraints
already register all nine Phase 23 Event types, their operation mappings, and
their complete `schema.md` section 6.2 effect rules, so `review.recorded` and
`finding.addressed` need no catalog edit at tasks 4 and 5. Those tasks own only
their fence key in `internal/journal/project_effects.go`. That rejection is
unaffected by R-16 and R-17: every registry in the inventory above lives in a
different file, and tasks 4 and 5 gain the rows they need there, not
`current_event_catalog.go`.

## 7. Test-Time Headroom (DEC-227)

Fresh QA Lead measurements at the planning base: the tagged integration suite
runs 560.2 s against Go's implicit 600 s per-binary default, with `cmd/fid` at
472.0 s and `internal/workflow` at 458.4 s under the race gate, and
`t.Parallel()` appears in none of the 58 `cmd/fid/integration` files. Phase 23
adds real-process journeys to exactly those packages at tasks 4, 6, 11, 14, 17,
and 18, so an unbounded default would convert ordinary growth into a killed
command - which is `BLOCKED`, never `PASS`, and whose only cure inside a task
would be a protected command change.

The disposition:

- Every Phase 23 task runs the two long gates with an explicit `-timeout 30m`.
  The bound is deliberately generous: it makes the harness limit a declared
  constant instead of an accidental cliff, and it weakens no assertion, changes
  no scope, and removes no gate.
- Wall time stays a measured budget, not an aspiration. Every completion report
  records the times `go test` reports for `cmd/fid/integration`, `cmd/fid`, and
  `internal/workflow`. Budgets: tagged integration under 900 s, no ordinary
  package over 600 s, asserted at phase end by VAL-079.
- Speed repair follows causation. A task whose own new journeys push a package
  past its budget owns the bounded repair - reuse of seeded workspaces, fewer
  redundant end-to-end setups, and `t.Parallel()` in **new** Phase 23 journey
  files where the fixture is self-contained. A retrofit of the 58 existing files
  is not in any task's scope; if the cause is pre-existing, stop and escalate
  rather than absorbing it.
- Task 1 additionally stops the per-run integration binary leak in
  `cmd/fid/integration/main_test.go` (CF-042, cause corrected under R-15 /
  DEC-229). The dominant cause is not timeout kills: five tests re-execute the
  test binary, each child re-enters `runIntegrationSuite` and rebuilds its own
  ~16 MB binary, and the crash-process children never reach `RemoveAll`. One
  clean run - exit 0, no kill - leaks 20 directories and 306 MB into a 3.9 GB
  RAM-backed `/tmp`, so this phase's ~26 tagged runs would need about 7.8 GB and
  fill `/tmp` around the twelfth. Because the leak is per-run, not per-kill, an
  age bound prunes nothing within a working day. The mitigation is therefore
  causal: `runIntegrationSuite` publishes the built binary path, a process that
  finds it reuses it and neither builds nor removes, and the 24-hour prune stays
  only as the residual safety net for a parent killed mid-run. Reuse also
  removes about twenty redundant `go build ./cmd/fid` invocations per run, which
  helps the budgets above. A full `/tmp` is a gate-blocking environment failure,
  not a test failure.

## 8. Incorporating Earlier Delivery Lessons

- `contract-coverage.md` maps every binding contract clause to its schema row,
  owning task, and validation assertion. It exists because Phase 22 shipped a
  cancellation clause with no schema row (DEC-195) and an undelivered checkpoint
  allocation (DEC-194). QA Lead re-derives that matrix independently before
  execution approval; a clause with no schema row or no owning task is a
  planning defect, not an implementation surprise.
- Cross-cutting validators reach many producers. The Obligation generalization
  in task 2 and the Evidence widening in task 6 are the two known wide-reach
  changes; both are isolated to their own task rather than bundled with a leaf.
- End-to-end journeys enter with each implemented wire surface at tasks 4, 6,
  11, and 14; task 18 audits their completeness, not their first behavior.
- Before a BLOCK is emitted, its allocated ESC must already exist in the
  current escalation register. Mis-scoped directives and externally interrupted
  sessions never consume a no-progress recovery budget.
