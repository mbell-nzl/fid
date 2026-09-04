# Phase 21 - TL Escalations

> Status: **CLOSED with the phase (DEC-179, 2026-09-04).** All seventeen
> allocated Phase 21 escalations are resolved. ESC-21-007 through ESC-21-011
> were resolved by DEC-165..DEC-169 during execution but were never written to
> this register while open; DEC-179 item 4 records them here from their
> governing decisions and reports the omission as a workflow defect.

## Open Escalations

None.

## Resolved Escalations

### ESC-21-017 - Accepted Receipt get/list leaves are not implemented

- **Raised by:** TRG-21-023 after the Task 22 attempt-4 tagged gate reached
  `fid receipt get` and received unknown-command exit 2; inventory also found no
  Receipt list or read API.
- **Protected category:** task-boundary and data-model authority for missing
  accepted public read leaves, durable rejected-Artifact facts, and protected
  contract-prose reconciliation.
- **Resolution:** DEC-177 assigns Task 22 the minimum Receipt get/list domain,
  CLI, persistence, golden, schema, and regression-test scope. Migration 0013
  stores canonical `rejected_artifacts_json` under ADR-0011's pre-baseline
  replacement policy; `contracts.md` §10 changes only the stale command spelling
  to exact `fid resume`. No operation, result, filter, ordering, error, Event,
  command block, corpus, VAL, budget, or prior evidence changes.
- **Recovery:** RUN-RESUME attempt 6 under TRG-21-023 / ESC-21-017 / DEC-177,
  then Task 22 `in-scope-implementation` in task scope with no-progress `0/2`.

### ESC-21-016 - Task-22 activation Event writer is outside accepted scope

- **Raised by:** TRG-21-022 after Task 22's attempt-2 tagged integration gate
  observed `activation.finalizeted` instead of accepted
  `activation.finalized`.
- **Protected category:** task-boundary authority for a production
  implementation repair affecting an accepted public Event contract.
- **Resolution:** DEC-175 gives Task 22 the bounded Event-type derivation repair
  and non-tagged regression assertion. DEC-176 preserves the immutable
  attempt-3 TRIAGE and reconciles the schema-v2 STOP resume as attempt 4. No
  contract, schema, interface, command, corpus, VAL, budget, or prior evidence
  changes.
- **Recovery:** RUN-RESUME attempt 4 with resolution authority DEC-176, then
  Task 22 `in-scope-implementation` in task scope under DEC-175/176 with
  no-progress `0/2`.

### ESC-21-015 - Task-22 exact journey cannot use a seeded reference workflow

- **Raised by:** TRG-21-021 after the accepted 20-Event direct-terminal journey
  conflicted with every migration-seeded task-applicable workflow.
- **Protected category:** specification contradiction, public interface, and
  workflow contract.
- **Resolution:** DEC-174 makes the journey define and assign its own accepted
  two-edge task workflow before Issue creation and amends the exact prefix to
  22 Events. The observed seeded-workflow rejection remains VAL-085 negative
  evidence.
- **Recovery:** Task 22 attempt 2 completed the authored-workflow journey to
  the independent activation Event-name defect and recorded
  `MATERIAL-PROGRESS`.

### ESC-21-014 - Task-20 V1 absence controls exceed its accepted source scope

- **Raised by:** TRG-21-020 after task 20's required absence controls still
  matched production in `internal/bootstrap/`, `internal/journal/`,
  `internal/labels/`, `internal/links/`, and `internal/requirements/`, none of
  which the task's Files To Create Or Modify list named.
- **Protected category:** specification contradiction and protected
  task-boundary authority.
- **Resolution:** DEC-173, 2026-09-02. The inventory (§1.1, §6) and DEC-160 §5
  already assign task 20 every residue: the `ActorRole`/`legacyRole` fallback
  readers, the role-bearing lock/bootstrap metadata, the bootstrap fixed-role
  registry, the `dev.*` golden, every `roles` fixture 0014 orphans, and the
  V1 closure journey. The file list was the incomplete artifact, so it is
  amended with bounded subject-preserving dispositions and task 20 gains the
  tagged integration gate under DEC-167. The ten absence declarations, both
  broad inventories, VAL-002/074/075, 0014, the frozen-table policy, and the
  public interface are unchanged. `08_QUALITY_PLAN.md` §3 records the converse
  executable-ownership rule.
- **Recovery:** RUN-RESUME retains TRG-21-020 / ESC-21-014 with DEC-173 at
  attempt 2, then task-20 `in-scope-implementation` in ordinary task scope.
  The Developer continues from `aef03cf`, converges the enumerated residues and
  fixtures with per-row provenance, reruns the complete amended command block,
  finalizes the report, and completes task 20 before task 21 dispatch.

### ESC-21-013 - Task-17 Attention absence control versus unavailable capability matrix

- **Raised by:** TRG-21-019 after the required broad production-source search
  matched the exact accepted unavailable `attention.v1` capability row.
- **Protected category:** specification contradiction and public interface.
- **Resolution:** DEC-172, 2026-08-30. The row remains exact and public with
  `available:false`, `{"state":false}`, and `phase-22-unavailable`; negative
  discovery metadata is not Attention state/effect. The later QA evidence
  recipe was defective because zero matches would also conceal the required
  row. VAL-063, QA-61, and task 17 now use
  `TestPhase21AttentionBoundary`, which requires that row as the sole production
  Attention match and rejects any additional table, Event, Effect, CLI leaf,
  state, operation, or mutation. The old failed command and blocker evidence
  remain unchanged.
- **Recovery:** RUN-RESUME retains TRG-21-019 / ESC-21-013 with DEC-172 at
  attempt 3, then task-17 `gate-repair` in ordinary task scope.
  The Developer adds the boundary test, preserves and commits the existing
  implementation, reruns the revised six-command block, finalizes the report,
  and completes task 17 before task 18 dispatch.

### ESC-21-012 - Task-15 Window 5 QA recovery budget exhausted

- **Raised by:** TRG-21-017, corrected for identity by TRG-21-018, after the
  Window 5 task-15 QA rework for MF-15-01 and MF-15-02 recorded two consecutive
  no-progress attempts.
- **Protected category claimed:** recovery-budget-exhausted.
- **Resolution:** DEC-170, 2026-08-30. **Not sustained on that category.** All
  five attempts under `task-15:qa-rework:review-5:MF-15-01-MF-15-02` ended with
  the Developer child shut down by the orchestrator while running
  (`run-ledger.md:2730`, `:2773`, `:2816`, `:2859`, `:2902`); none returned a
  typed result, so the NO-PROGRESS verdicts measured interruption rather than
  repair capability. The corrected category is orchestration session-liveness,
  as in DEC-168. The durable tree already contains MF-15-01's causal fix and
  MF-15-02's four negative cases, so the repair is at verification and closure.
  DEC-170 makes an externally interrupted attempt never consume the no-progress
  budget, restates that no budget bounds a Developer child session's lifetime
  while fixing per-command budgets for the unchanged task-15 block, orders the
  durable work committed before the long commands with each result appended as
  it terminates, treats the recorded attempt-5 gate rows as unverified until
  rerun in session, and replaces MF-15-02's history-rewriting provenance route
  with a forward-only one because `79e3a93` and `8491cde` are ancestors of the
  control repair base `09a2a7b`.
- **Recovery:** `RUN-RESUME` retaining TRG-21-018 and ESC-21-012 with resolution
  authority `DEC-170`, then one `TASK-RECOVERY task-15`, category
  `session-continuation`, mode `task-scope`, attempt 6, consecutive no-progress
  `0/2`. Task 16 waits for a fresh independent `review-5b` PASS. No finding,
  assertion, command, VAL, QA artifact, contract, schema, or ownership changes.

### ESC-21-011 - Group 1 checkpoint repair needs completed-task owner routing

- **Recorded at wrap-up (DEC-179 item 4).** Allocated by TRG-21-015 and retained
  by the identity-only successor TRG-21-016; resolved 2026-08-29 but never
  written to this register while open.
- **Raised by:** TRG-21-015 after the independent Group 1 checkpoint returned
  REPAIR with MF-1 owned by completed task 4 and MF-2 owned by completed task 9,
  which the open DEC-168 task-13 orchestration wrapper could not re-own.
- **Protected category:** cross-task ownership and phase-order authority under
  DEC-168 item 10(a).
- **Resolution:** DEC-169. MF-1 routing is correct and widens to the whole
  Activation family, because `assertActivationAvailable` runs before
  `journal.CommitIdempotent` and rejects an exact same-key replay. MF-2 keeps
  task 9 as owner but widens to the complete section-7 category-to-exit mapping,
  because `cmd/fid/main.go` was unconditionally `os.Exit(1)` and no Group 1
  command ran a real binary - the fifth instance of the DEC-167 executable
  ownership gap. VAL-042/043, FEAS-01..19, reviewer independence, and the
  checkpoint contract are unchanged, and the complete checkpoint reruns.
- **Recovery:** TL-Lite allocates identity-only successor TRG-21-016
  (`Originating task: task-4`), retaining ESC-21-011. The ledger runs
  `REVIEW infra-review-1` REPAIR with owners `task-4, task-9`, TRIAGE, STOP,
  RUN-RESUME at attempt 2 under DEC-169, the task-4 and task-9 recovery triples,
  then `REVIEW infra-review-1b` PASS on the complete rerun.

### ESC-21-010 - Group 1 checkpoint unavailable after interrupted QA sessions

- **Recorded at wrap-up (DEC-179 item 4).** Allocated by TRG-21-013, corrected
  for identity by TRG-21-014; resolved 2026-08-29 but never written to this
  register while open.
- **Raised by:** TRG-21-013 after three fresh bounded QA Lead checkpoint
  sessions were each interrupted before returning a verdict, leaving no
  `qa/group-1-checkpoint-review.md`.
- **Protected category claimed:** unavailable external capability.
- **Resolution:** DEC-168. **Not sustained on that category.** The QA capability
  was available; orchestration was applying the 900-second FEAS-set command
  budget as the child session's total lifetime and restarting from zero. The
  binding reading is 120s per FEAS probe, 900s for the FEAS set, 300s for the
  tagged integration command, and **no budget on a child session's lifetime**.
  `qa-plan.md` and `group-1-checkpoint-contract.md` are not edited; VAL-042,
  VAL-043, FEAS-01..19, commands, statuses, and reviewer independence are
  unchanged; preserved scratch and diagnostic runs are inadmissible, so every
  probe reruns fresh and the artifact is written incrementally.
- **Recovery:** RUN-RESUME retaining TRG-21-014 / ESC-21-010 under DEC-168,
  continuing task-13 at attempt 4, then one `session-continuation` task-scope
  recovery. The runner dispatches the QA Lead first, must not interrupt it at
  900 seconds, and must not dispatch the task-13 Developer before a committed
  checkpoint PASS.

### ESC-21-009 - Task-12 recovery budget claimed exhausted on owner-routed gates

- **Recorded at wrap-up (DEC-179 item 4).** Allocated by TRG-21-012; resolved
  2026-08-28 but never written to this register while open.
- **Raised by:** TRG-21-012 after two consecutive runner-mediated task-12
  continuations failed the unchanged required gates on failures attributed to
  task 9 and future task 20.
- **Protected category claimed:** recovery-budget-exhausted plus task-boundary
  and phase-order authority.
- **Resolution:** DEC-167. **Not sustained.** The `internal/dev` failure was
  task 12's own collateral, the `issue ac verify` failure was VAL-030's
  fail-closed leaf rather than a transport gap, and the budget did not apply
  because the ledger directive had forbidden task 12 from repairing the only two
  failing rows - both no-progress outcomes were guaranteed by directive. The
  root cause is gate coverage, not ownership: `cmd/fid/integration/` is
  build-tagged and DEC-165 had assigned task 9 that package without a command
  that runs it. Task 12's scope is amended to `internal/dev/` and
  `cmd/fid/integration/` with per-row provenance, and `08_QUALITY_PLAN.md` §3
  gains both standing rules.
- **Recovery:** RUN-RESUME retaining TRG-21-012 / ESC-21-009 under DEC-167, then
  one same-task task-12 `in-scope-implementation` recovery at attempt 6 with
  no-progress `0/2`; all six task-12 commands rerun byte-identically.

### ESC-21-008 - Retained `TestTask10*` fixtures treated as cross-task property

- **Recorded at wrap-up (DEC-179 item 4).** Allocated by TRG-21-010; resolved
  2026-08-26 but never written to this register while open.
- **Raised by:** TRG-21-010 after task 9 stopped on red rows in
  `cmd/fid/task10_contract_test.go` that still expected the Phase 20 `--by`
  transport, reading the filename and `TestTask10` prefix as Phase 21 task-10
  ownership.
- **Protected category claimed:** scope boundary / accepted-contract
  contradiction.
- **Resolution:** DEC-166. **Not sustained on that category.** Commit
  provenance places the file and its assertions in **phase-20** task-10, and
  Phase 21 task 10 is generic Claim work that owns none of those leaves; the
  implementation's rejection of `--by` is correct. A filename, test name, task
  number, or prose label is not current-phase ownership evidence. Task 9 owns
  and may rename the stale rows to subject-based names, preserving every
  assertion subject and strength. Recurrence of a classification is not evidence
  that it is correct and creates no protected category.
- **Recovery:** RUN-RESUME retaining TRG-21-010 / ESC-21-008 under DEC-166,
  continuing task 9 at attempt 12, then one same-task `in-scope-implementation`
  recovery with all five task-9 commands rerun byte-identically.

### ESC-21-007 - Task-9 closure attributed to task-10 Claim/Run fixtures

- **Recorded at wrap-up (DEC-179 item 4).** Allocated by TRG-21-009; resolved
  2026-08-25 but never written to this register while open.
- **Raised by:** TRG-21-009 after eight task-9 attempts, reporting that the
  required gates stayed red on replacement-bound task-10 Claim/Run fixtures its
  Out Of Scope forbade it to touch.
- **Protected category claimed:** protected task-boundary / cross-task
  authority.
- **Resolution:** DEC-165. **Not sustained on that category.** No generic
  Claim/Run fixture existed in the tree; every red row failed on an incidental
  read of a surface task 9 itself removed or renamed (`policy show`, the `--by`
  to `--group-by` stats rename, the fail-closed vacuum). Convergence ownership
  follows the invalidating accepted change, not the fixture's headline subject.
  Task 9's Out Of Scope is clarified to exclude Claim/Run implementation and
  semantics only, incidental `policy show` reads are repointed onto the retained
  legacy matrix helper with byte-identical assertions, and the V1 Claim, Run,
  dev, and closure behavioural assertions keep full strength until their own
  owner replaces them.
- **Recovery:** RUN-RESUME retaining TRG-21-009 / ESC-21-007 under DEC-165, then
  one same-task task-9 `in-scope-implementation` recovery at attempt 9 with
  no-progress `0/2`.

### ESC-21-006 - Schema-v2 ordinary repair authority cannot bind TASK-RECOVERY

- **Raised by:** TRG-21-008 after the exact task-8 recovery governed by
  TRG-21-007 was rejected.
- **Resolution:** DEC-164. In ordinary schema-v2 repair, TASK-RECOVERY `Source`
  references the governing TRG while the five external identity fields,
  including `Repair authority`, remain byte-identical. The verifier removes
  only its redundant requirement that the already exact-bound ordinary
  authority also contain the TRG ID and admits the bound same-task recovery as
  replacement for the pending dispatch whose blocker candidate was classified.
  Other overlaps plus control-conformance, legacy recovery, STOP/RUN-RESUME,
  decision, and exact external-binding guards are unchanged.
- **Recovery:** a Human-directed workflow-infrastructure change adds the
  ordinary-authority positive drill and Source negative drill, runs the full
  drill/workflow checks, and commits the corrected verifier. The Phase Runner
  then appends the pending TRG-21-007 task-8 recovery directly; no RUN-RESUME,
  new TRIAGE, historical rewrite, or task-8 contract change applies.

### ESC-21-005 - Multi-owner QA review recovery cannot bind the single-owner sequence control

- **Raised by:** TRG-21-005, corrected for external identity by TRG-21-006,
  after QA Review 1 of tasks 1-3.
- **Resolution:** DEC-163. A QA `REPAIR` review carries an ordered
  comma-separated owner list in exactly one `Owning task` field, matching the
  review artifact's Repair Plan row order. Every listed owner completes
  TASK-RECOVERY, REPAIR, and TASK complete before a linked re-review may close
  or replace the review. `CORRECTION` remains evidence-only and is rejected as
  an owner-selection mechanism. `workflow-sequence-lint.sh` is corrected under
  protected decision authority to implement this, with new drills, because it
  accepted an append that left the ledger unrecoverable. No historical record is
  rewritten; MF-1 and MF-2 keep their exact boundaries and evidence.
- **Recovery:** after the verifier change lands, RUN-RESUME attempt 2 under
  TRG-21-006 / ESC-21-005 / DEC-163 into `task-1 | qa-rework | task-scope |
  task-1:orchestration-integrity:review-1:multi-owner-ledger-recovery-binding |
  2` for review-1 MF-1, then an ordinary task-3 `qa-rework` recovery for MF-2,
  then fresh re-review `review-1b`. Task 4 stays gated.

### ESC-21-004 - Retained rebuild and role-row tests versus accepted Event-v2 contracts

- **Raised by:** TRG-21-004 after the DEC-160 policy-bridge recovery.
- **Resolution:** DEC-162 keeps the no-leaf-flags `project.rebuild` interface and
  the exact role-free Event-v2 bootstrap. Retained V1 tests are not protected
  contracts: the task whose accepted change invalidates a retained test adapts it
  in that same task, preserving subject and assertion strength. Task 2 gains the
  test-file scope, ordered adaptation rules, and a closed no-weakening list;
  `journal.RebuildWorkspace` drops its vestigial role parameter.
- **Recovery:** same-task, task-scope `in-scope-implementation` continuation of
  task 2 at attempt 6 under the TRG-21-004 trigger signature, with both
  previously failed required commands rerun byte-identically.

### ESC-21-001 - Role-free bootstrap versus retained V1 policy

- **Raised by:** TRG-21-001 after task 2.
- **Resolution:** DEC-160 keeps the exact zero-Effect Event-v2 bootstrap and
  authorizes a closed, nonpersistent seven-name V1 policy vocabulary only for
  inventoried legacy callers. Explicit legacy role Effects remain
  authoritative; new generic paths cannot consult it; task 20 removes it.
- **Recovery:** same-task task-2 continuation with focused policy red/green
  tests and byte-identical reruns of both failed required commands.

### ESC-21-002 - Markdown-delimited external trigger identity

- **Raised by:** TRG-21-002 after rejected TRIAGE candidates.
- **Resolution:** DEC-161 makes schema-v2 identity values plain single-line
  scalars. TRG-21-001/TRG-21-002 remain immutable malformed evidence; neither
  governs recovery.
- **Recovery:** none against the verifier or historical records.

### ESC-21-003 - Protected recovery route

- **Raised by:** TRG-21-003 and the valid STOP ledger.
- **Resolution:** DEC-161 uses already-valid TRG-21-003 as the governing record
  and supplies exact RUN-RESUME attempt 3 authority. No protected verifier edit
  or control-conformance repair is required.
- **Recovery:** Phase Runner appends RUN-RESUME, then task-scope TASK-RECOVERY
  for task 2 under DEC-160/DEC-161. Task 3 stays dependency-gated.

## Archive Rule

Any Phase 21 `BLOCK` must allocate an `ESC-21-NNN` through `triage-log.md` and
record its protected resolution condition here before Human attention.
