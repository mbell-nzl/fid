# Phase 21 - TL Escalations

## Open Escalations

None. DEC-160 through DEC-173 resolve ESC-21-001..014. The run is stopped after
TRG-21-020; DEC-173 authorizes the task-20 attempt-2 continuation that clears it.

## Resolved Escalations

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
