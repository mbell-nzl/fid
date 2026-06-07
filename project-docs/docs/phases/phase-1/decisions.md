# Phase 1 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow
> notes, and wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-006: Enforce SQLite pragmas per physical connection via DSN

- **Date:** 2026-06-07
- **Context:** QA Review 1 / ESC-001 / TRG-003 (BLOCK). `internal/db.Open`
  applies `PRAGMA foreign_keys = ON` and `PRAGMA recursive_triggers = ON` by
  `Exec` on the pooled `*sql.DB` (`internal/db/db.go:23-28`). SQLite pragma
  state is per physical connection, and `database/sql` opens additional
  physical connections on demand (default `MaxOpenConns = 0`, unlimited).
  `foreign_keys` defaults OFF on each new connection, so VAL-014's "both
  pragmas ON on every connection path" is not met. `TestOpenEnablesRequiredPragmas`
  only inspects one connection and almost always reuses the same physical
  connection, so it passed while the guarantee was false. The DB helper is
  consumed by the journal writer, replay, init, and issue tasks, so the gap
  spans subsystems.
- **Decision:** Enforce both pragmas through the connection string using
  `modernc.org/sqlite`'s `_pragma` query parameters, e.g.
  `file:<abs-path>?_pragma=foreign_keys(1)&_pragma=recursive_triggers(1)`,
  so the driver applies them on every physical connection it opens. Keep the
  post-open `Ping`. Add a regression test that forces multiple distinct
  physical connections (via `SetMaxOpenConns(N)` and N concurrently held
  `*sql.Conn`) and asserts both pragmas report ON on each. Implemented as
  remediation Task 7. Rejected alternatives: pinning the pool to a single
  connection (`SetMaxOpenConns(1)`) satisfies VAL-014 only by degenerate
  accident and serializes all DB access; a global `RegisterConnectionHook`
  is heavier than DSN params and was kept only as a fallback.
- **Consequences:** `db.Open` keeps its public signature; this is an internal
  correctness change plus a stronger test. ESC-001 is dispositioned **do-now**
  in Phase 1 and closes when Task 7 lands with passing multi-connection
  evidence. No schema, data model, or other public-interface change.
- **Refs:** ESC-001 (`tl_escalations_current.md`); TRG-003 (`triage-log.md`);
  QA Review 1 Task 1 Blocking (`qa/review-1.md`); VAL-014
  (`qa/validation-contract.md`); QA-01 (`qa/qa-plan.md`); task-7
  (`tasks/task-7.md`); `internal/db/db.go`; TRD-ARCH-007; TRD §4.4.

### DEC-007: Close COV-001 with package-level domain tests (Task 8)

- **Date:** 2026-06-07
- **Context:** QA-11 / Quality Plan COV-001 was the sole blocking finding in the
  Phase 1 final review. The domain packages `internal/issues` (0.0%),
  `internal/policy` (0.0%), and `internal/journal` (68.3%) were below the 80%
  package-level threshold. The behavior was substantially exercised through
  `cmd/fid` integration tests, but Go attributes coverage to a package's own
  test binary, so the named domain packages registered no own-package coverage
  and several real branches (policy corrupt-tail tolerance, issues not-found /
  input-guard paths) were unproven at any layer. The human accepted that the
  threshold percentage is noisy at tracer-bullet scale but asked the Tech Lead
  to judge whether this reflected weak TDD or was genuinely acceptable.
- **Decision:** Remediate rather than defer. Phase 1 deliberately locks the
  hardest-to-change contracts, so thin domain-level tests here are worth the
  cost, and the untested branches included corruption/not-found behavior we
  care about (TRG-001/TRG-002 areas). Authorize a **test-only** Task 8 adding
  package-level unit tests for `internal/issues`, `internal/policy`, and
  `internal/journal` to meet COV-001 (≥80%) while holding COV-002 (`cmd/fid`
  ≥70%); no production behavior change. A blanket "coverage is meaningless at
  small scale" deferral was rejected as the wrong precedent on the foundation
  phase.
- **Consequences:** Task 8 landed (commit `f851b6a`, test files only) with
  independently verified coverage `internal/issues` 85.2%, `internal/policy`
  85.0%, `internal/journal` 80.6%, `cmd/fid` 70.5%; all gates green under
  `-race`. QA re-check (`f009131`) flipped the final-review verdict to PASS and
  COV-001 to MET. `internal/db` (36.2%), `internal/bootstrap`, and
  `internal/format` remain low but are not COV-001 domain packages, so they do
  not gate this phase.
- **Planning note (carry-in for future contracts):** The root cause was a
  planning seam — the Phase 1 validation contract prescribed CLI *integration*
  evidence for VAL-006/VAL-016 while COV-001 demands *package-level* coverage of
  the same packages. Future validation contracts should name package-level test
  evidence when a VAL's behavior lives in a COV-001 domain package.
- **Refs:** QA-11, COV-001/COV-002 (`qa/final-review.md`,
  `project-docs/specs/08_QUALITY_PLAN.md` §4.3); task-8 (`tasks/task-8.md`,
  `tasks/task-8-completion.md`); TRG-001/TRG-002 (`triage-log.md`).

### DEC-008: Authors commit their own completion reports and QA artifacts

- **Date:** 2026-06-07
- **Context:** The human observed during Phase 1 that Developers were not always
  committing their completion reports and QA authors were not committing their
  reviews, leaving artifacts to be committed by a later actor (e.g. `ba5ee18`
  finalized the Task 7 completion report in a separate commit by a different
  actor). This blurs authorship: git history stops reliably recording who
  attested to what.
- **Decision:** Make "authors commit their own artifacts under their own
  authorship before handoff" an explicit workflow rule. A task is not done while
  its completion report is uncommitted; QA reviews and final reviews are
  committed by their authoring agent.
- **Consequences:** Encoded in `.claude/agents/developer.md` (workflow step +
  constraint), `.claude/agents/qa-reviewer.md` (new commit step + constraint),
  `.claude/agents/qa-lead.md` (commit requirement + constraint), and
  `project-docs/specs/08_QUALITY_PLAN.md` (Definition of Done now requires the
  completion report be committed by its author; phase-level gate step 12 adds
  the QA-artifact commit rule). Committed in `78c3540`.
- **Refs:** `.claude/agents/developer.md`, `.claude/agents/qa-reviewer.md`,
  `.claude/agents/qa-lead.md`, `project-docs/specs/08_QUALITY_PLAN.md` §10 and
  §5.2.

### DEC-009: Reconcile TRD §4.3 last-event-hash freshness wording (TRG-001)

- **Date:** 2026-06-07
- **Context:** TRG-001 flagged that `02_TRD.md` §4.3 said a recorded
  event ID/hash mismatch at the recorded position "stops before answering and
  reports projection inconsistency," while the Phase 1 validation contract
  (VAL-012, `hash_mismatch_rebuilds`), the QA plan, and the Task 4
  implementation treat a `projection_meta` checkpoint hash mismatch as a
  *freshness* signal that triggers a rebuild/catch-up before serving. The
  durable TRD text contradicted the shipped, validated behavior.
- **Decision:** Reconcile the durable spec to the validated behavior rather than
  reopen a passed phase. The journal is the canonical source of truth
  (ADR-0001), so a projection checkpoint that diverges from the journal is
  untrustworthy and is rebuilt from the journal. The "never silently skip,
  truncate, or reorder journal entries" guarantee and the hard-stop on
  *malformed or precondition-conflicting* journal entries during replay are
  preserved. Updated `02_TRD.md` §4.3 wording accordingly.
- **Consequences:** `02_TRD.md` §4.3 now distinguishes projection-checkpoint
  divergence (→ rebuild) from in-journal corruption (→ stop). No code change;
  VAL-012 behavior is unchanged and now matches the TRD.
- **Refs:** TRG-001 (`triage-log.md`); VAL-012 (`qa/validation-contract.md`);
  `tasks/task-4-completion.md`; `02_TRD.md` §4.3.

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is
a positive assertion after review, not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | verified-no-change | Phase 1 implemented existing charter scope (tracer bullet); no identity/scope change. |
| `01_PRD.md` | verified-no-change | PRD-INIT/OPS/ROLE/OUT/NFR requirements MET as written (final review Requirement Compliance); no requirement reworded. |
| `02_TRD.md` | updated | §4.3 last-event-hash freshness wording reconciled to validated VAL-012 behavior (DEC-009 / TRG-001). |
| `03_TRD_SCHEMA.md` | verified-no-change | Minimal `projection_meta`/`roles`/`issues` schema matches migration `0001`; remaining tables deferred to Phase 2 as planned. |
| `04_TRD_CONFIG.md` | verified-no-change | `config.toml`/`pricing.toml`/prefix behavior matches VAL-001/VAL-002 evidence; `[export]` deferred as noted. |
| `05_DESIGN_PLAN.md` | verified-no-change | Journal→projection→rebuild design realized as specified; no design-language change. |
| `06_DESIGN_PLAN_CLI.md` | verified-no-change | `init`/`rebuild`/`issue create`/`issue get` shapes match; TOON/markdown/YAML deferred to Phase 3 as planned. |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | Tracer-bullet flow matches documented persistence flow; no flow change. |
| `08_QUALITY_PLAN.md` | updated | Author-commits-own-artifacts rule added to §10 Definition of Done and §5.2 phase-level gate step 12 (DEC-008). |
| `09_AI_WORKFLOW.md` | verified-no-change | Author-commit rule is encoded in the Quality Plan DoD and agent definitions, which are its canonical homes; the workflow narrative already references commit/report requirements and needed no change. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 1 delivered the v0.1a tracer-bullet slice per the roadmap slicing note; no re-slice required. |

## Triage-Log Audit

- [x] Walked every TRG-NNN in `triage-log.md` (TRG-001, TRG-002, TRG-003).
- [x] Confirmed each RESUME verdict was correct in hindsight. TRG-001 (RESUME)
  correctly continued the phase under VAL-012 behavior and is now reconciled in
  the TRD by DEC-009. TRG-002 (RESUME) correctly identified acceptable Phase 1
  bootstrap debt; carried forward as CF-003. TRG-003 (BLOCK→ESC-001) correctly
  stopped on the pragma gap; resolved by DEC-006 / Task 7.
- [x] Flagged any TRG that should have escalated. None mis-triaged: TRG-003 did
  escalate (ESC-001); TRG-001/TRG-002 were correctly local RESUME with
  carry-forward dispositions recorded here.
- [x] Confirmed autonomous TL-Lite resume remains in effect (DEC-002; no warm-up tally).

## Workflow Notes

Answer in 1-3 sentences each. Completed at phase wrap-up.

- **What worked well?** The three-tier quality model caught two real foundation
  defects before sign-off: the pooled-connection pragma gap (ESC-001→DEC-006,
  Task 7) and the domain-package coverage shortfall (QA-11→DEC-007, Task 8).
  Both were remediated with focused, test-backed tasks rather than waved
  through. VAL-by-VAL evidence made the final review auditable.
- **What was inefficient?** A planning seam cost a remediation cycle: the
  validation contract prescribed CLI integration evidence for VAL-006/VAL-016
  while COV-001 requires package-level coverage of the same packages, so
  domain packages read 0% until Task 8 backfilled unit tests. Catching this at
  contract-authoring time would have avoided the extra task.
- **Were tasks the right size?** Mostly yes. The original six tasks were
  well-scoped vertical slices. The two remediation tasks (7, 8) were correctly
  bounded and single-purpose.
- **Did the QA review cadence feel right?** Yes. Two mid-phase reviews plus a
  phase-end review surfaced both blockers; Review 1 caught the pragma gap early
  enough to fix in-phase.
- **Were there token-efficiency opportunities missed?** Minor. The Task 8
  blocker was knowable from the coverage model at planning time; resolving it in
  the validation contract would have saved a full dev/QA loop.

### Deep Retrospective Recommendation

Lightweight retro only. Two blockers, both cleanly remediated, no recurring
or cross-cutting failure pattern. The one durable lesson — align validation-
contract evidence methods with COV-001's package-level coverage loci — is
already captured as a carry-in note in DEC-007 and should be applied at the
next phase's contract-authoring step. No full deep retrospective warranted.
