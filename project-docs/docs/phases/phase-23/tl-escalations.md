# Phase 23 Escalations

## ESC-23-010 - Task 18 Failure-code count contradicts the named public set

- **Raised:** 2026-09-30
- **Allocated by:** TL-Lite, TRG-23-011
- **Candidate source:** `project-docs/docs/phases/phase-23/tasks/task-18-done.md` (BLOCKED)
- **Protected category:** public-interface and specification contradiction.
- **Trigger:**
  - `public-interface.md` section 5 names 19 Failure codes under the lead
    "Added codes".
  - The same section's DEC-230 paragraph and DEC-230 say the added set is 18.
  - The task-18 corpus and `p23FailureContracts` both hold 19 rows.
  - Separately, the shared fake/fid audit compared result field sets and
    replay, not complete semantic values, so VAL-068 stayed open.
- **Trigger signature:** task-18:protected-contract:failure-code-cardinality-and-semantic-fake-equality
- **Attempt count:** 1
- **Checkpoint:** report `333514946`, corrected at `dac84a1` and `a156691`;
  report-lint `OUTCOME BLOCKED` (terminated `exit_code: 0`). All five exact
  task commands passed.
- **Question:** Record the exact accepted Failure-code members and count,
  and reconcile section 5 with DEC-230.
- **Resolution condition:** The decision and the task-18 continuation scope
  are recorded. The Runner resumes task 18 and keeps its attempt-1 evidence.
- **Links:** TRG-23-011; DEC-230; `public-interface.md` sections 5 and 5.1;
  `tasks/task-18-done.md`; VAL-056; VAL-068.
- **Human attention required:** yes. The Human asked the Tech Lead to
  resolve it on 2026-09-30.
- **Status:** RESOLVED 2026-09-30 under DEC-265.
  - **Tech Lead finding:** there is no member conflict. The section 5 table,
    the 5.1 producers, `p23FailureContracts`, and `failures.json` all name
    the same 19 codes. "18" counts codes new to `codes.go`, and
    `requirement-not-found` predates the phase (`ab4ac8f`, phase-20/task-1).
    QA plan 13.2 had already recorded "19 declared, 18 new". The defect was
    section 5's lead sentence, which called all 19 "added".
  - **Resolution:** Section 5 declares 19 codes: 18 added and 1 retained.
    Only the wording in the section 5 header, lead, and DEC-230 paragraph
    changes. DEC-230 gains a clarification line, and the validation contract
    gains a note-only row. Codes, contracts, VAL text, ownership, and
    commands are unchanged.
  - **Continuation:** same-task `in-scope-implementation` at attempt 2,
    no-progress 0/2. The work is:
    - pin the 19/18/1 partition;
    - complete the full semantic fake/fid parity and rejection classes for
      VAL-068;
    - record the VAL-067 at-or-after provenance;
    - rerun the five unchanged commands.

    QA window 16-18 then follows. Handoff: `qa/resume-after-esc-23-010.md`.
- **Second disposition (2026-09-30):** TL-Lite recorded a BLOCK under
  TRG-23-011 / ESC-23-010, deliberately allocating no new ID, after task-18
  attempts 2-6 (STOP 18:58).
  - **Trigger:** the fake expected the retained Phase 21 four-field
    `authority-denied` details from `claim.acquire`, while fid returned
    `{operation,actor}` (`claim.go:232`, phase-21/task-10), which is outside
    task-18 scope. The focused and tagged commands failed (terminated
    `exit_code: 1`) only on that subtest.
  - **Status:** RESOLVED 2026-09-30 under DEC-266.
  - **Tech Lead finding:** the registry (`codes.go:286-288`, phase-20/task-1)
    was never converged to Phase 21 section 7 for three of the four
    `authority-*` codes.
    - Eight producers emit the Phase 20 shape.
    - `gate.go:325/327` emit the contract shape, and the registry rejects
      it. A scratch probe showed this; those paths are latent defects.
    - The fake's other authority rows matched fid's non-conforming shape.
  - **Resolution:** the contract governs, and no contract, VAL, or command
    changes. A task-18 planning amendment converges the registry, the eight
    producers, `codes_test.go`, and the fake's two-field rows. It adds gate
    regression tests.
  - **Continuation:** same-task attempt 7, no-progress 0/2. Handoff:
    `qa/resume-after-dec-266.md`.

## ESC-23-009 - Migration 0018 down conflicts with the accepted 0016 rollback

- **Raised:** 2026-09-30
- **Allocated by:** TL-Lite, TRG-23-010
- **Candidate source:** `project-docs/docs/phases/phase-23/tasks/task-15-done.md` (BLOCKED)
- **Protected category:** schema/persistence contract and specification contradiction.
- **Trigger:**
  - `schema.md` section 5 required 0018 down to restore the byte-identical
    0002 `concerns` definition, including `escalated_to_esc_id`.
  - Migration 0016 down then runs
    `ALTER TABLE concerns ADD COLUMN escalated_to_esc_id`, so the literal
    down chain duplicates the column.
  - The focused rollback test returned `FAIL` (terminated `exit_code: 1`)
    because 0018 down is absent.
- **Trigger signature:** task-15:protected-contract:0018-down-concerns-column-duplication
- **Attempt count:** 1
- **Checkpoint:** implementation checkpoint `fa6355a`; report `7884725`,
  corrected at `b1f0028`; report-lint `OUTCOME BLOCKED` (terminated
  `exit_code: 0`).
- **Question / proposed direction:** Choose and record the exact 0018-down
  `concerns` definition. Reconcile section 5 with 0016 down and keep the
  full down chain executable.
- **Resolution condition:** Record the decision and task-15 continuation
  scope. The Runner resumes task 15 from the preserved checkpoint and keeps
  the failed evidence.
- **Links:** TRG-23-010; `schema.md` section 5;
  `migrations/0002_v01_tables.up.sql`;
  `migrations/0016_p22_remove_escalations.down.sql`; `tasks/task-15-done.md`.
- **Human attention required:** yes. The Human asked the Tech Lead to resolve
  it, approved DEC-264, and approved the `schema.md` edit on 2026-09-30.
- **Status:** RESOLVED 2026-09-30 under DEC-264.
  - **Tech Lead finding:** section 5 restored objects as first created, not
    as they stood at schema 17. The defect has three instances; triage
    named only the first:
    - `concerns` with `escalated_to_esc_id`;
    - `issues_enriched` from 0003 instead of 0011;
    - the two implicitly dropped 0004 kind triggers, which were omitted.
  - **Resolution:** 0018 down restores each object's fresh-schema-17
    definition. The only accepted difference is SQLite's rename quoting.
    - Section 5 now lists nine ordered steps, including an `agent_runs`
      rebuild that keeps rows.
    - The Tech Lead ran the full chain on a scratch database: 17 -> 18 -> 17
      gave a `sqlite_master` identical to a fresh 17, then 0017 through 0001
      down ran to empty.
    - Unchanged: migrations 0001-0017, section 4, VAL text, ownership, and
      commands.
  - **Continuation:** same-task `in-scope-implementation` at attempt 2,
    no-progress 0/2. Attempt 1 could not author the down file under a
    contradictory definition. Handoff: `qa/resume-after-esc-23-009.md`.

## ESC-23-008 - Task-13 protected-decision denial response contradicts closed failure schema

- **Raised:** 2026-09-28
- **Allocated by:** TL-Lite, TRG-23-008
- **Candidate source:** `project-docs/docs/phases/phase-23/tasks/task-13-done.md` (BLOCKED)
- **Protected category:** public-interface/specification contradiction and validation assertion semantics.
- **Trigger:** `protected-decision-required` has four exact closed detail fields under accepted `public-interface.md` section 5 and retains the closed `fid-cli.v1` error envelope under section 1. The same accepted section 5.2 and VAL-057 require the repeated denial response to expose immutable Stop/Attention IDs, null EventRef, and `replayed:false`. The built binary returns the closed failure without those additional fields; the focused replay probe and full tagged task-13 gate each returned `FAIL` (terminated `exit_code: 1`).
- **Trigger signature:** task-13:protected-contract:protected-decision-denial-envelope-val-057
- **Attempt count:** 1
- **Checkpoint:** verified report commit `348a538167f94bd1d5a5f9a9ef798a90698565e7` with matching subject; report-lint `OUTCOME BLOCKED` (terminated `exit_code: 0`). DEC-262 closed the unrelated ActivityEvent scope gap; the issue-subject rows now pass, while task 15 retains the requirement/decision rows. DEC-227 package times are within budget.
- **Question / proposed direction:** Full Tech Lead: obtain the Human or accepted contract owner's decision on how the closed failure response represents the Stop/Attention IDs, null EventRef, and replay value, or explicitly amend VAL-057's accepted response requirement. Align `public-interface.md` sections 1, 5, and 5.2 with VAL-056/VAL-057 and task-13 acceptance without silently widening the public schema or weakening the failing assertion.
- **Resolution condition:** Record the protected contract decision and exact task-13 continuation scope; the Runner then resumes task 13, preserves failed evidence, reruns unchanged required verification, and requests any affected independent re-review.
- **Links:** TRG-23-008 in `project-docs/docs/phases/phase-23/triage-log.md`; `project-docs/docs/phases/phase-23/public-interface.md` sections 1, 5, 5.2; `project-docs/docs/phases/phase-23/qa/validation-contract.md` VAL-056 and VAL-057; DEC-226; `project-docs/docs/phases/phase-23/tasks/task-13-done.md`.
- **Human attention required:** yes. The Human delegated the resolution to the Tech Lead on 2026-09-29.
- **Status:** RESOLVED 2026-09-29 under DEC-263.
  - **Tech Lead finding:**
    - Section 5 copied the two codes' StopDetails verbatim as their failure
      details. It never added the denial fields that DEC-226 part 4
      (Human-approved), section 5.2, `schema.md` 6.3, and VAL-057 require.
      This is a drafting defect in the accepted set, not a task-13 defect.
    - `human-outcome-required` has the same defect under VAL-056. Triage
      named only the protected-Decision half.
  - **Resolution:** both codes' details become their StopDetails followed by
    `stop_ids`, `attention_ids`, `event` (always null), and `replayed`
    (always `false`). This is the `review-rework-exhausted` and Phase 22
    `budget-exhausted` form.
    - Unchanged: the envelope, StopDetails, Stop projections, codes,
      producers, precedence, VAL-056/057 text, and gates.
    - Aligned: `public-interface.md` 5 and 5.2 and `06_DESIGN_PLAN_CLI.md`
      section 5.
  - **Continuation:** same-task `in-scope-implementation` at attempt 2,
    no-progress 0/2. Attempt 1 made material progress, and it was not
    authorized to widen the closed schema. Handoff:
    `qa/resume-after-esc-23-008.md`.

## ESC-23-007 - Task-13 ActivityEvent journal and requirement-reference schema boundary

- **Raised:** 2026-09-28
- **Allocated by:** TL-Lite, TRG-23-007
- **Candidate source:** `project-docs/docs/phases/phase-23/tasks/task-13-done.md` (BLOCKED)
- **Protected category:** data-model/persistence contract and cross-task migration ownership.
- **Trigger:** A scoped `review-recorded` activity trial was rejected by the
  strict journal validator with `unknown table` for `activity_log/activity`.
  `requirement.record-reference` activity cannot resolve `requirement_refs`
  under the 0016 trigger. The replacement trigger is in task 15's 0018. The
  focused journey and the exact tagged gate each returned `FAIL`
  (terminated `exit_code: 1`).
- **Trigger signature:** task-13:scope-gap:review-list-and-activity-producers
- **Checkpoint:** `838fc32` (implementation), `007c948` (report).
- **Links:** TRG-23-007; `tasks/task-13.md` VAL-055; `schema.md` sections 4 and 6.2; `execution-scope.md` section 6.
- **Human attention required:** yes. The Human delegated the resolution to the Tech Lead on 2026-09-28.
- **Status:** RESOLVED 2026-09-28 under DEC-262.
  - **Tech Lead finding:**
    - The `unknown table` rejection is a planning gap of the DEC-167/173
      family. Execution-scope section 6 registered only the ten new tables,
      but all nine strict-path families declare `activity_log/activity`.
    - The trigger half also affects `decision-referenced`. Triage named only
      the requirement half.
  - **Resolution:** no DDL, migration, order, rule, or name changes.
    - Task 13 gains the `activity_log` registry entry (row 11) and the five
      Review-package producers, and emits only `issue`-subject rows.
    - Task 15 emits the `requirement`/`decision`-subject rows with 0018 and
      closes VAL-055's ten-name inventory.
    - VAL-055's owner cell names both tasks; its pass and fail text is
      unchanged.
    - The activity subjects are fixed in `schema.md` section 6.2.
  - **Continuation:** same-task `in-scope-implementation` at attempt 2,
    no-progress 0/2. Attempt 1 made material progress, and it was not
    authorized to touch the registry that caused its remaining failure.
    Handoff: `qa/resume-after-esc-23-007.md`.

## ESC-23-006 - Task-12 requirement authority and report-lint contract conflict

- **Raised:** 2026-09-27
- **Allocated by:** TL-Lite, TRG-23-006
- **Candidate source:** `project-docs/docs/phases/phase-23/tasks/task-12-done.md` (RECOVERY_REQUIRED)
- **Protected category:** a public-interface/specification contradiction, plus authority over the task and control definitions.
- **Trigger:** Two texts conflicted.
  - `public-interface.md` section 4 required Actor-profile capabilities for
    `requirement.record-reference` and `requirement.link`.
  - Section 4.1, execution-scope section 6 row 7 and R-19 exclude both names
    from the closed 30-name set.

  The real-process journey then failed twice at Activation staging with
  `invalid activation input` (each terminated `exit_code: 1`). Separately,
  task 12's expected-nonzero declaration was wrapped over two lines, and
  report-lint rejected it with `FAIL expected non-zero declaration has
  invalid grammar` (terminated `exit_code: 1`).
- **Trigger signature:** task-12:protected-contract:requirement-capability-registry-and-report-lint-grammar
- **Checkpoint:** implementation commit `0e4bb33`, report commit `3284869`.
- **Links:** TRG-23-006; `public-interface.md` sections 4 and 4.1; `execution-scope.md` section 6; R-19; `tasks/task-12.md`; `project-docs/tools/report-lint.sh:104`.
- **Human attention required:** yes. The Human delegated the resolution to the Tech Lead on 2026-09-27.
- **Status:** RESOLVED 2026-09-27 under DEC-260.
  - **Part 1 (Option B):** requirement writes need a trusted Actor from the
    finalized Kit Activation, but no Actor-profile capability and no Envelope
    grant. Only the section 4 paragraph is corrected; the thirty-name set,
    section 4.1, R-19, VAL-067 and task 18 are unchanged.
  - **Part 2:** the report-lint grammar is intended. Task 12's declaration was
    a Tech Lead authoring defect and is rewritten on one line. Report-lint on
    the retained report and the corrected spec now passes (terminated
    `exit_code: 0`).
  - **Continuation:** same-task `in-scope-implementation` at attempt 2.
    Handoff: `qa/resume-after-esc-23-006.md`.

## ESC-23-005 - Task-11 ordinary gate budget has no authorized task-scope repair

- **Raised:** 2026-09-27
- **Allocated by:** TL-Lite, TRG-23-005
- **Candidate source:** `project-docs/docs/phases/phase-23/tasks/task-11-done.md`; Phase Runner triage trigger
- **Protected category:** gate-repair scope and ownership authority; an accepted gate-definition change would require protected decision authority
- **Trigger:** Task-11's unchanged ordinary race/coverage gate passed tests but exceeded the accepted 600-second `cmd/fid` package budget twice, with no demonstrated task-11 causal speed repair.
- **Trigger signature:** review-7:qa-rework:task-11-BC-3
- **Recovery attempts:** next mediated attempt 4; measured ordinary-gate attempts 686.386 s and 630.927 s; both remain failed against VAL-079.
- **Checkpoint:** attached verification of implementation commit `466745f2dc4505d60eb67ac50193e004532b7239` with subject match; attached report-lint returned `OUTCOME RECOVERY_REQUIRED` and `PASS report-lint completed with no failures` (terminated `exit_code: 0`); Runner supplied report commit `c1588fc4389c7dc8210247fc48ede2618c783799`.
- **Evidence:** Exact `go test -race -timeout 30m -coverprofile=coverage.out ./...` terminated `exit_code: 0`, but `cmd/fid` took 686.386 s and then 630.927 s. A fresh same-session pre-repair parent took 614.601 s. Task 11 added no ordinary `cmd/fid` test file; focused TestStop timing was 109.999 s on the parent and 112.324 s current, while unrelated packages also rose. The tagged integration gate passed at 766.070 s, under 900 s. Neither a task-11 cause nor a passing ordinary budget result is established.
- **Question / proposed direction:** Full Tech Lead: determine bounded ownership and a closure route for the ordinary-gate speed cause under DEC-227 and VAL-079, accounting for the separate DEC-254 task-13 headroom authorization. If the accepted command, 600-second budget, or pass criterion must change, seek the authority that owns that definition. Do not treat passing tests or an isolated later green as resolution of the retained failed attempts.
- **Resolution condition:** Recorded full Tech Lead disposition assigning an authorized bounded repair and fresh exact verification under the unchanged gate, or an accepted protected decision changing the gate definition; then the required fresh independent checkpoint review must PASS before task 12.
- **Links:** TRG-23-005 in `project-docs/docs/phases/phase-23/triage-log.md`; `project-docs/docs/phases/phase-23/tasks/task-11.md` Exact Commands To Run; `project-docs/docs/phases/phase-23/tasks/task-11-done.md`; `project-docs/docs/phases/phase-23/decisions.md` DEC-227 and DEC-254; `project-docs/docs/phases/phase-23/qa/validation-contract.md` VAL-079.
- **Human attention required:** yes
- **Status:** RESOLVED 2026-09-27 under DEC-258 (Human-delegated repair). The cause was serial test structure plus measurement spread, not task 11. Test-only `t.Parallel()` across the suites, with the gate definition unchanged. The final exact gates at `3dff457` passed: `cmd/fid` 531.533 s, `internal/workflow` 165.920 s, tagged 326.421 s. Evidence: `qa/esc-23-005-repair.md`; handoff: `qa/resume-after-esc-23-005.md`.

## ESC-23-004 - Post-task-11 checkpoint re-review: three protected blocker candidates (BC-1, BC-2, BC-3)

- **Raised by:** TRG-23-004
- **Status:** RESOLVED 2026-09-26.
  - BC-1 by DEC-249 and BC-3 by DEC-250 (Tech Lead).
  - BC-2 by DEC-251 (Human, Option B).
  - Carry items: I-1 by DEC-254 (Part 2 Human-approved), I-2 by DEC-253, I-5 by DEC-252.
  - DEC-255/256 release ordered repairs: task 6 MF-4+I-5, task 2 BC-1, task 10 MF-5+I-6, task 11 BC-3. The first three owners are complete in the ledger.
  - Next: a fresh `checkpoint-after-task-11-3.md` re-review.
  - Task 12 stays held until that re-review returns PASS with Must-fix 0 (qa-plan section 8).
- **Originating task:** task-10
- **Trigger signature:** task-10:qa-rework:review-7-BC-1-BC-2-BC-3
- **Source:** `project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11-2.md` (review-7, commit 9cb9420516c461ddd9b4d22d1944759a416f7bb7)
- **Category:** specification, product semantics and public interface. Each sub-item is separately resolvable.
- **Sub-item BC-1 (product semantics, VAL-009):** contracts.md:91-92 says a guarded action requires "every applicable Obligation" to be satisfied, but "applicable" is undefined. contracts.md:60-63 and :408-409 imply per-Obligation guarded_actions, while the evaluator (internal/workflow/evaluation.go:436,496-512; envelope.go:521,546-556) blocks terminal success on any unsatisfied Obligation regardless of its guarded_actions. Decision needed: does an Obligation block only its declared guarded_actions? Option A (yes): a decision-free task-2 evaluator conformance repair, then VAL-009 (both halves), the journey stage 1-3 `qa_pass` assertions and task 13 VAL-056 are re-derived. Option B (no, global): record the reading in a DEC and a GLOSSARY entry, including what guarded_actions governs. The task-10 stage 3 assertion then stands and MF-5's `qa_pass` expectation is fixed. Affects tasks 2, 10, 13, 14 (VAL-059, VAL-071).
- **Sub-item BC-2 (public interface / authorization, VAL-043):** public-interface.md:243 (section 4.1) requires an Envelope grant for `human-outcome.record`. public-interface.md:222-224 and contracts.md:391-396 do not, and the product and task-9 journey need none (internal/review/human_outcome.go:79-89). Decision needed: which text governs. Option A (grant required): a task-10 authorization repair plus a VAL-043 negative case, and the task-9 journey fixture gains a grant for `actor-human-outcome`. Option B (grant not required): correct the section 4.1 cell under a DEC. Human ruling needed if the interface changes (DEC-222). Affects tasks 9, 10, 14 (VAL-071) and 18 (VAL-067/068).
- **Sub-item BC-3 (resolution semantics, VAL-048, VAL-050, VAL-078):** internal/workflow/stop.go `eventMatchesCondition` and `qualifyingFact` have no rule for `human-outcome-recorded` or `decision-satisfied`, so a Stop carrying such a Condition raised by a non-producer code can never be resolved. contracts.md:453-455 and :461-462 fix that a `human-outcome-recorded` Stop resolves with a class-and-result-matching Outcome Event and that matching Stops resolve. They do not fix the `decision-satisfied` qualifying fact (which fact of the pair, Envelope scope), subject scoping for `human-outcome-recorded`, or whether recording resolves every Condition-matching Stop regardless of raising code. Decision needed: the Tech Lead fixes both qualifying rules. Task 11 then implements them in its stop.go scope and evidences the VAL-048 reject clause and VAL-078 per type. Affects task 11 and task 17 (`human_outcome_resolution:true`).
- **Independent bounded Must-fix items (not blocked by the candidates unless noted):** MF-4 (task-6, VAL-027 linked-recovery leg, test-only, independent of all three candidates, may dispatch once the Tech Lead releases the stop). MF-5 (task-10, VAL-009 evidence-method state plus I-6 citations, test/report-only, its `qa_pass` assertion waits on BC-1). Each has an owning task, bounded scope, required evidence and a fresh independent re-review in the review's Repair Plan.
- **Carry items:** I-1 `cmd/fid` headroom, I-2 task-11 scope ratification, I-5 multi-grant `gate.run` owner, I-6 report citations. No action here.
- **Constraints:** Do not choose any reading in code, tests or reports before the decision is recorded. Do not narrow VAL-009, VAL-043, VAL-048, VAL-050 or VAL-078. A fresh independent checkpoint re-review must PASS before task 12.
- **References:** TRG-23-004 in `project-docs/docs/phases/phase-23/triage-log.md`; review-7 lines 661-776.

## ESC-23-003 - Unreachable Run-holder-only CLI evidence

- **Raised by:** TRG-23-003, Task 4 attempt 8.
- **Original condition:** VAL-015 required a CLI Run-holder-only independence
  rejection although Run start binds the Run Actor to the Claim Actor and
  retains durable Claim history. A diagnostic Run-history query existed, but
  accepting it required protected evidence authority. The historical BLOCK
  and STOP remain in `triage-log.md` and `run-ledger.md`.
- **Status:** RESOLVED 2026-09-20 under DEC-239 and the Human's explicit
  delegation through Task 9 and its mid-phase review.
- **Resolution:** Preserve Actor binding, durable history and both independence
  query directions. Accept the isolated Run-holder query as diagnostic evidence
  only; retain real CLI Claim rejection and same-Reviewer rework. Reconcile the
  canonical TRD, Review flow and quality plan with that proof boundary.
- **Residual work:** Task 4 still needs its same-Reviewer implementation repair,
  multiple-prior-Finding resolution proof, built-binary evidence, VAL-003 parity
  closure and current-session exact gates. The decision resolves the protected
  contradiction; it does not mark those observations passed.
- **Resume:** `qa/resume-after-esc-23-003.md`; governing TRG-23-003/ESC-23-003,
  DEC-239, same trigger, attempt 9. Tasks 10-11 remain held.

## ESC-23-001 - Widening the journal `Condition` value breaks journal replay and Attention deduplication

- **Phase:** 23 (planning; raised before any task dispatch)
- **Status:** **RESOLVED 2026-09-20 - Human chose Option A.** Recorded as
  DEC-235; `public-interface.md` section 2's `Condition` row amended, inventory
  row 8 restated and row 10 created, tasks 9, 10, and 11 rewritten. Nothing is
  blocked; all eighteen tasks are dispatchable. The record below is kept as
  raised, with the resolution at the end.
- **Raised by:** Tech Lead, while dispositioning R-19 under the widened
  inventory frame the QA Lead asked for in I-33
- **Category:** accepted public-value encoding, with migration and
  durable-replay consequences
- **Decision needed from:** Human

### Summary

`public-interface.md` section 2 amends the `Condition` value to add
`outcome_class:?ID`, `decision_class:?ID`, and `finding_id:?ID`, with "only the
field matching the type is non-null". Tasks 9, 10, and 11 each add one of those
fields to `internal/workflow/stop.go`'s `Condition` struct and to
`currentValueShapes["Condition"]`.

`Condition` is the **only** field-adding widening of a live journal value shape
in this phase - every other amendment in section 2 either widens a closed enum
(`Claim.mode`, `Budget.unit`, `Budget.exhaustion_action`, `StopCode`) or lands
on a family in `retainedCurrentRules`, which the retained branch returns from at
`current_event_validation.go:64` before any value validation. Adding a **field**
is not monotonic, for three independent reasons, none addressed anywhere in the
accepted contract set or the eighteen task specs: a search of every Phase 23
document for `semantic_key`, `EventSchema`, and `omitempty` returns nothing.

### Evidence

Reproduced under `go test -overlay`, working tree untouched, with only
`currentValueShapes["Condition"]` widened.

1. **Journal replay of every pre-widening Event fails.**
   `validateCurrentShape` (`current_event_value_types.go:112`) requires an
   **exact** field-count match: `len(object) != len(fields)` rejects. A Condition
   in today's four-field form returns `Condition: closed object has wrong shape`
   against the widened shape. That value is reached on the strict path through
   `validateCurrentJSON`'s `case "stops.resolution_condition_json"`
   (`current_event_json.go:27`), and `stop.raised`, `stop.resolved`, and
   `attention.resolved` are all in `currentEventRules`
   (`current_event_catalog.go:46-48`) and behind no fence. `ReplayJournal`
   (`replay.go:153`) calls `ProjectEvent` -> `applyEventEffects`
   (`project_effects.go:11`) -> `validateCurrentEvent` for **every historical
   Event**, so rebuild-from-journal revalidates old Events against current
   shapes. Phase 23 does not bump `EventSchema`, so `UpgradeToCurrent`
   (`upgrade.go:50`) passes a Phase 22 V2 Event through unchanged and no hook
   normalizes it. Any workspace holding a Stop raised before the widening can no
   longer rebuild - which also puts VAL-004 and VAL-064 at risk at tasks 15/16.
   **Four reach points, not one** (QA confirmation, `qa-plan.md` 17.2). Besides
   the `stops.resolution_condition_json` column, `Condition` is revalidated as a
   nested field of the `Stop` value in the typed results of `stop.raised`,
   `stop.resolved`, and `run.budget-observed` (`validateCurrentResultJSON`,
   `current_event_json.go:134-141`): a Phase 22 `Stop` carrying a four-field
   Condition is rejected with `Stop: resolution_condition: Condition: closed
   object has wrong shape`. Confirmed by execution rather than by reading: the
   same recorded journal bytes rebuild `REBUILD OK: events=1 stops=1
   attention=1` on the baseline and fail under the widening as a
   `CorruptJournalError` - the operator is told the journal is **malformed**,
   not that the binary is too new. Those four sites are what Option B's upgrader
   must rewrite.
2. **Attention deduplication silently breaks across the upgrade.**
   `attentionSemanticKey` (`stop.go:422`) hashes the whole `Condition` struct
   via `MarshalCanonicalJSON`, which preserves explicit nulls. Measured: the same
   `budget-observed` condition hashes
   `sha256:8aa57e2a0374b532…` before and `sha256:3407ff3b2d279107…` after the
   three nulls are added. `attention_items.semantic_key` is a stored column, so
   after the widening `loadOpenAttentionByKey` (`stop.go:557`) no longer
   matches an Attention that was open at upgrade time, and an identical Stop
   raises a **second** open Attention. Confirmed end to end across two binaries
   sharing one workspace: an identical second raise joins the open Attention on
   the baseline (`openAttention=1`) and opens a second one under the widening
   (`openAttention=2`), both keys stored side by side. Phase 22's
   one-open-Attention-per-semantic-key invariant is lost with nothing recording
   that it was given up.
3. **Every live `stop.raise` idempotency record is invalidated.** Measured by
   the QA Lead and not in the original record (`qa-plan.md` 17.2).
   `StopRaiseRequest` (`stop.go:72-79`) embeds `ResolutionCondition Condition`,
   and `stopDigest` (`stop.go:431`) canonicalises the whole request into
   `idempotency_results.request_digest`, so widening the struct changes the
   digest of every `stop.raise`. On a workspace whose Stop was raised by the
   baseline binary, replaying the **same** idempotency key under the widened
   struct returns `idempotency conflict for stop.raise:project-… key
   "retained-test-1"` instead of the stored result: a retry that was idempotent
   before the upgrade becomes a hard refusal after it. That is a third stored
   identity the widening breaks, and a third column Option B's migration must
   recompute.

None of the three is caught by any current assertion, and none is a defect in
the plan's ownership: tasks 9, 10, and 11 own exactly the right files. The
accepted encoding itself is what does not fit.

### Why this needs the Human

Every remedy either changes an accepted public value's JSON encoding or adds
unaccepted machinery to the phase's highest-risk area. That is the same class as
DEC-223 to DEC-227, all of which the Human approved. Tech Lead planning
authority covers ownership, file scope, and recorded cause; it does not cover
changing what a Phase 23 value looks like on the wire.

### Who this reaches, and what is actually being decided

Recorded so the choice is made on the real population rather than an assumed
one (QA confirmation, `qa-plan.md` 17.4).

`schema.md` section 1 accepts task 1's edits to the shipped `0013` and `0015`
files as "a deliberate pre-adoption baseline replacement requiring fresh init
and rebuild; it is not an upgrade of old fid data", and `golang-migrate`
(`internal/db/migrate.go:8`) records only version and dirty state - it never
re-runs or checksums an applied file. Under that accepted policy there is no
released population of Phase 22 workspaces expected to run Phase 23 features,
so Option B's migration machinery is work for a population the policy says does
not exist. That does **not** dissolve this escalation:

1. **The incompatibility is in-phase, not only post-release.** The widening
   lands at tasks 9-11, and task 6 raises `review-rework-exhausted` Stops
   carrying ordinary `budget-observed` Conditions. Any workspace - developer,
   gate, or fixture - whose journal was written by a tasks 1-8 binary and is
   rebuilt by a task 9+ binary hits the measured `CorruptJournalError`. That is
   the VAL-004 and VAL-064 risk at tasks 15/16, and it is real rather than
   theoretical. One mitigating fact, checked: **no** pinned journal fixture in
   `testdata/` carries a `resolution_condition` today, so the exposure is live
   workspaces plus any fixture task 6 chooses to pin.
2. **The Attention and idempotency consequences need no rebuild at all.** Both
   are stored-column comparisons, so any workspace carrying an open Attention or
   a live `stop.raise` idempotency record across the upgrade is affected whether
   or not it ever replays.

The decision is therefore better framed as **"what encoding do we want for V1,
given that either choice costs in-phase churn"** than as "how do we protect
existing user data". Framed that way, Option A costs 31 lines and buys the
inventory row 8 invariant.

### Options

- **Option A (recommended) - the three new Condition fields are absent, not
  null, on the types they do not apply to.** `omitempty` on the **three new
  `*string` fields only** - `outcome_class`, `decision_class`, `finding_id` -
  plus per-type Condition shapes in the journal keyed on `type`, as
  `currentStopDetailShapes` is already keyed on `code` and already owned by
  tasks 6, 9, 10, and 11.

  **Read the encoding exactly; one reading of it destroys the remedy.**
  `budget_name` and `capability_id` keep today's present-and-null treatment, so
  the six Phase 22 types marshal byte-for-byte as they do now and a new type
  carries **five** fields - `type`, `subject`, `budget_name:null`,
  `capability_id:null`, and its own field - not three. The alternative reading,
  "a new type carries literally only its own field", would put `omitempty` on
  `budget_name` and `capability_id` as well; that reading was measured and it
  changes **every** existing semantic key (`38cda009…` -> `cf865c69…`,
  `da2f2a11…` -> `ef2f1a6a…`, `28607d77…` -> `9fb954bc…`), because `work-ended`
  drops from `{"budget_name":null,"capability_id":null,…}` to
  `{"subject":…,"type":"work-ended"}`. That is precisely the harm Option A
  exists to prevent. The correct reading was implemented under overlay and
  measured: the recorded pre-widening journal replays `REBUILD OK: events=1
  stops=1 attention=1`, and all three Phase 22 semantic keys stay
  **byte-identical to the baseline** (`da2f2a11…`, `38cda009…`, `28607d77…`).

  **Benefit.** Historical Events replay unchanged, and every stored
  `semantic_key` and `stop.raise` `request_digest` stays valid, so all three
  measured consequences above are eliminated. The exclusivity rule also becomes
  **structural**: measured, the wrong field is rejected with `Condition: missing
  required field "outcome_class"`, the own field null with `Condition:
  outcome_class: invalid ID`, two fields with `Condition: closed object has
  wrong shape`, and the own field alone is accepted - all four cases inventory
  row 8 leaves open, closed by the shape itself, which **removes** the relation
  work row 8 assigns to tasks 9, 10, and 11. Against the shipped corpus
  `internal/journal`, `internal/workflow`, and `internal/export` are all `ok`
  with **zero failures and zero fixture edits**.

  **Cost.** The public JSON for a Condition omits the new fields that do not
  apply to its type instead of showing them as null. That is a real amendment to
  the `public-interface.md` section 2 row, which presupposes present-and-null.
  Size measured at 16 lines of catalog map plus a 15-line dispatch branch - and
  it is **one file larger than this record first stated**. Unlike
  `currentStopDetailShapes`, which is selected at the call site from the row's
  `code` column (`current_event_json.go:35-40`), a Condition's discriminant
  lives **inside** the value, and the value is also reached generically as a
  field of `Stop` through `if shape, ok := currentValueShapes[kind]; ok`
  (`current_event_value_types.go:39`). Per-type selection therefore needs a
  `Condition` branch in `validateCurrentValue`, in
  `internal/journal/current_event_value_types.go` - a file the inventory names
  only as a *reader* of `currentStopDetailShapes` and that **no task's file
  scope lists**. Under the phase's own ownership discipline it must join the
  scope of whichever task implements the remedy.

  **One implementation constraint to carry into the decision: the lookup must
  use a checked type assertion.** The dispatch reads `object["type"]` **before**
  `validateCurrentShape` has established it is a string. Measured: the baseline
  degrades cleanly to `Condition: type: invalid closed enum
  =executor-terminated,…`, while the obvious `object["type"].(string)` panics
  with `interface conversion: interface {} is int, not string`. This code
  validates untrusted journal and archive bytes during replay and `import`, so
  the branch must return the ordinary rejection instead. One line of
  implementation; worth one sentence in the decision, because a panic here is a
  crash on a crafted archive.
- **Option B - keep the accepted encoding literally, and pay for it.** Every
  Condition always carries all seven fields. Cost, corrected upward by
  measurement (`qa-plan.md` 17.2-17.3): an `EventSchema` bump to V3, plus an
  upgrader that rewrites historical Condition JSON at **four** payload sites
  rather than one - the `stops.resolution_condition_json` column and the nested
  `Stop` inside the typed results of `stop.raised`, `stop.resolved`, and
  `run.budget-observed`; **and** a migration-0017 recompute of
  `attention_items.semantic_key` for every open Attention **and** of
  `idempotency_results.request_digest` for every live `stop.raise` record;
  **and** an accepted rule for what happens to an Attention whose key changes
  mid-life; **and** repair of the **11 failing tests and subtests** across 7
  top-level tests that the present-and-null encoding breaks in the shipped
  corpus before any Option B machinery is written
  (`TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay`,
  `TestActivationRecoveryEventsProjectIdenticallyLiveAndReplay/activation.recovery-selected`,
  `TestCurrentEventRejectsConditionalBindingCorruption` - three cases x
  live/replay, `TestTerminalSuccessEvaluationIncludesEnvelopeStop`,
  `TestRecoverRunRejectsApplicableStopWithoutEffects`,
  `TestStopAndAttentionListsFilterLimitAndRemainReadOnly`,
  `TestEquivalentStopsShareAttentionButKeepDistinctFacts`). It also **keeps**
  inventory row 8's relation work, which Option A removes. In its favour: it
  preserves the accepted value exactly, it keeps one shape per value kind as
  every other `currentValueShapes` entry does, and it is genuinely
  implementable - `ReplayJournal` hashes the stored Event before upgrading it
  (`replay.go:141`, then `:147`), so a payload-rewriting upgrader does not break
  hash chains. It simply costs four to five times Option A, in tasks 1 and 15,
  the two tasks already carrying the most migration risk.
- **Option C - defer the three fields to Phase 24.** Ship the four new Condition
  **types** (a safe enum widening) with no new fields, and carry the class and
  Finding discriminators in Stop details, where `decision_class` and
  `outcome_class` already exist for `protected-decision-required` and
  `human-outcome-required`. Rejected on inspection: the Attention semantic key is
  derived from `(subject, code, condition)` only, so with no class field in the
  Condition, two human-outcome Attentions of different classes on the same
  Envelope produce the same key, `loadOpenAttentionByKey` joins them into one
  Attention, and resolving one class closes the Attention covering both.
  Verified by the QA Lead, who adds that the same collision applies a second
  time to `decision-satisfied` and `decision_class`. Recorded here only so it is
  not re-derived.

### Resolution as decided

**The Human chose Option A on 2026-09-20, in the corrected reading**: `omitempty`
on the three new fields only, the six Phase 22 types keeping `budget_name` and
`capability_id` present-and-null, with per-type journal shapes keyed on `type`.
Recorded as DEC-235, which carries the full disposition. Applied in this batch:

- `public-interface.md` section 2's `Condition` row now states the encoding, the
  five-field judgment Condition, the untouched four-field Phase 22 Condition,
  and that exclusivity is structural. (*Corrected by DEC-236 / R-26: only the
  three judgment types that have a discriminator carry five fields;
  `review-recorded` carries four. `public-interface.md` itself is correct - this
  summary line is not, and is historical.*)
- `execution-scope.md` section 6: inventory **row 8 restated, not deleted** -
  the relation branch keeps the Phase 22 `budget_name` / `capability_id`
  binding and **no Phase 23 task extends it**. New **row 10** is the per-type
  shape map and its dispatch, failing closed, with task 9 creating it and owning
  `internal/journal/current_event_value_types.go`, and tasks 10 and 11 adding
  rows. *(Corrected by DEC-236: the reason first recorded here - "which the
  shapes cannot express" - was wrong, and the branch must still be **reached**
  by task 9's new dispatch or the binding is lost silently.)*
- Tasks 9, 10, and 11: the `omitempty`-on-new-fields-only rule in all three
  (each edits the struct), the checked type assertion in task 9, byte-identity
  of the six Phase 22 rows as an acceptance bar, and the relation branch marked
  do-not-touch.
- VAL-078 keeps its assertion, method, pass and fail criteria, validator, and
  owning task; only its **evidence class** is annotated, from relation rejection
  to shape rejection.
- The three held revisions are applied with it: R-22 (proving Event named in
  tasks 9 and 10, as a direct strict-path assertion rather than a raise), R-23
  (`task-11.md` now has a VAL-078 criterion), I-37 (`schema.md` section 1).

### Original required resolution, as raised

Human chooses A, B, or C (or directs another route). Whichever is chosen, record
it as a DEC, amend `public-interface.md` section 2's `Condition` row to match,
and update tasks 9, 10, and 11 plus inventory row 8 in `execution-scope.md`
section 6 before those tasks are dispatched. Three specifics for the Option A
branch, so the disposition is not rederived: the implementing task also gains
`internal/journal/current_event_value_types.go` in its file scope and the
checked-assertion requirement in its acceptance criteria; inventory row 8 is
**restated**, not deleted, because the invariant survives and only its enforcer
moves from the relation branch to the shape; and VAL-078's evidence class
becomes shape rejection rather than relation rejection, with the assertion
itself unchanged.

Three bounded QA revisions are deliberately held for the same disposition
batch, because each is phrased against whichever encoding is chosen and each
touches only tasks that are already blocked: R-22 (name the proving Event for
the `currentStopDetailShapes` and `Condition` entries tasks 9 and 10 add, since
three of DEC-224's four raising families are retained and do not exercise them),
R-23 (`task-11.md` has no acceptance criterion citing VAL-078), and I-37 (one
recorded sentence in `schema.md` section 1 on the silent in-place-migration
divergence). None of them is a dispatch blocker.

### Scope of the block

Tasks 1 through 8 are unaffected and may dispatch on schedule: none of them
touches the `Condition` shape or struct. Only tasks 9, 10, and 11 are blocked,
and task 12 onward only through the strict task ordering. The R-19 disposition
and everything else in DEC-233 stand independently of this escalation.

Independently confirmed by the QA Lead across all eighteen specs, not only the
three this record names (`qa-plan.md` 17.5): no task before 9 touches the shape
or struct, and task 6 - the closest, because it owns the same Stop registries
and co-edits the same catalog file on unrelated lines - explicitly forbids
adding a Condition type. One consequence worth the Runner's attention: **task 12
is substantively independent of the `Condition` shape**, so if the decision
takes time, task 12 is available out of order at the Runner's discretion.

### References

- `project-docs/docs/phases/phase-23/public-interface.md` section 2, `Condition`
  row; `contracts.md` section 11.
- `project-docs/docs/phases/phase-23/execution-scope.md` section 6, inventory
  row 8.
- `project-docs/docs/phases/phase-23/decisions.md` DEC-229 (Condition types are
  registered by the task that first raises a Stop carrying them), DEC-233, and
  DEC-234 (this record's R-20 amendments and what is held for the disposition
  batch).
- `project-docs/docs/phases/phase-23/qa/qa-plan.md` sections 17.2, 17.3, 17.4 -
  the QA Lead's independent reproduction of all three consequences, the measured
  Option A implementation, and the harm-population framing.
- `internal/journal/current_event_value_types.go:39,112`;
  `internal/journal/current_event_json.go:27,35-40,134-141`;
  `internal/journal/current_event_catalog.go:46-48`;
  `internal/journal/replay.go:141,147,153`;
  `internal/journal/project_effects.go:11`;
  `internal/journal/upgrade.go:50`; `internal/db/migrate.go:8`;
  `internal/workflow/stop.go:72-79,422,431,557`.

## Blocker candidates raised without an identifier

Recorded here so the phase register is complete; neither became an escalation.

- **QA pass 7, `qa-plan.md` 18.12 - conditional on branch B.** Withdrawn by the
  QA Lead once DEC-236 chose branch A, and independently re-checked on three
  legs in `qa-plan.md` 19.11: the discriminator pairing was accepted before the
  escalation existed, the Human approved a mechanism rather than a per-type
  field assignment, and the contradicting sentence was new in the `932fbea`
  batch. Restoring the accepted reading is Tech Lead planning authority, so no
  Human decision was required.
- **QA pass 8, `qa-plan.md` 20.5 - R-28, VAL-078's granularity.** Correctly
  raised as a candidate with no identifier allocated: the protected category is
  validation-contract authority, the repair is a Tech Lead contract edit, and no
  product, encoding, or Human judgment is requested. Dispositioned under DEC-237
  by splitting the row into VAL-078 and VAL-080. A Developer could not have
  repaired it under task 11's product-file scope, which is why it was a
  candidate and not a Developer deviation.

## ESC-23-002 - Task 1 gate closure needs protected ownership and budget authority

- **Phase:** 23
- **Status:** RESOLVED under DEC-238, 2026-09-20
- **Raised by:** TL-Lite via TRG-23-001
- **Category:** protected gate/control authority and cross-task retained-assertion ownership
- **Decision needed from:** full Tech Lead, with Human authority if the accepted gate definition must change

### Historical condition as raised

Task 1 cannot close. Its required 100k rebuild benchmark measures
32.500159404 s/rebuild against the strict under-30-second self-assertion, and
the exact repository-wide race gate fails with `cmd/fid` at 619.168 s and
`internal/workflow` at 614.196 s, each over DEC-227's accepted 600-second
ordinary-package budget. The same race run retains
`cmd/fid/task12_execution_transition_test.go:217`, whose assertion accepts only
`workflow-transition-invalid` but observed `claim-expired`, `run-inactive`, and
`workflow-transition-invalid`.

Task 1 may repair speed only where its own new journeys caused the breach. No
evidence establishes that causal ownership for the retained cmd/fid assertion
or both budget breaches. Changing the benchmark threshold, package budget,
command, or retained assertion semantics is protected and is not authorized by
the task. The full Tech Lead must establish bounded causal ownership and repair
authority, or obtain the authority required for any accepted-definition change.

### Evidence

- Benchmark command `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` failed: `FAIL github.com/mbell-nzl/fid/internal/journal 35.245s` (terminated `exit_code: 1`); measured rebuild 32.500159404 s.
- Race command `go test -race -timeout 30m -coverprofile=coverage.out ./...` failed (terminated `exit_code: 1`); `cmd/fid` took 619.168 s and `internal/workflow` 614.196 s.
- The task's temporary performance experiments were removed because they either failed to reach the budget or broke accepted rebuild parity; report lint returned `OUTCOME RECOVERY_REQUIRED` (terminated `exit_code: 0`).

### Original required resolution

Full Tech Lead: identify, from accepted ownership and causation evidence, whether
one or more bounded repairs exist and assign the correct task authority,
including a fresh independent review where required. If the only remedy changes
the accepted command, timeout, threshold, budget, assertion semantics, accepted
status/scope, or policy, obtain the owning protected authority before any
change. Do not assign the retained `cmd/fid` assertion to task 1 without that
authority.

### References

- TRG-23-001 in `project-docs/docs/phases/phase-23/triage-log.md`.
- `project-docs/docs/phases/phase-23/tasks/task-1.md` (Exact Commands To Run).
- `project-docs/docs/phases/phase-23/tasks/task-1-done.md` (Workflow outcome and gate evidence).
- `project-docs/docs/phases/phase-23/execution-scope.md:44-45,501-508`.
- `project-docs/docs/phases/phase-23/decisions.md` DEC-227 and `project-docs/docs/phases/phase-23/qa/validation-contract.md` VAL-079.


### Resolution under DEC-238

ESC-23-002 is **RESOLVED 2026-09-20 under DEC-238**. The Human delegated
bounded repair authority above the fixed agent roles. Rebuild statement reuse
and fresh/isolated test setup repair the causes without changing any gate or
assertion. Fresh exact gates pass: CLI **527.099s**, workflow **231.010s**,
tagged integration **669.586s**; rebuild samples **22.728s / 22.623s / 22.489s**.
Integration binary directories stay **0 → 0**. Task 1 is complete in its report;
the Runner must reconcile the retained STOP before advancing. Tasks 9-11 stay
held. See [repair evidence](qa/esc-23-002-repair.md) and
[resume handoff](qa/resume-after-esc-23-002.md).
