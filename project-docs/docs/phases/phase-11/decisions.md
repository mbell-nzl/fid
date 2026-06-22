# Phase 11 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow notes, and
> wrap-up decisions. DEC numbering is global and never resets.

## Decisions

### DEC-028: Claim writes must refresh stale projections without nested locks

- **Date:** 2026-06-21
- **Context:** ESC-008 / TRG-001 found that `internal/claims.Acquire`,
  `Renew`, and `Release` call `journal.EnsureProjectionFresh` from inside their
  claim write-lock callback. When the journal is ahead of SQLite, that helper
  can catch up or rebuild by trying to acquire `.fid/write.lock` again, which
  produced nested-lock failures in Task 4 fixture work.
- **Decision:** The current lock ordering is not acceptable for Phase 11 final
  QA. Add do-now Task 5 before final review. Claim write paths must validate
  against a fresh projection while avoiding nested acquisition of the same
  checkout-scoped write lock. A lock-aware journal helper or equivalent narrow
  refactor is acceptable; public claim CLI flags/output and claim semantics must
  not change.
- **Consequences:** Validation contract gains VAL-015 and QA plan gains
  corresponding evidence. QA review may resume after Task 5 is complete and
  reviewed. No schema migration or spec change is planned because this enforces
  the existing TRD write-path/concurrency contract.
- **Refs:** ESC-008; phase-11/task-5; PRD-NFR-008; `02_TRD.md` §4.3 and §4.9;
  ADR-0001; VAL-015.

### DEC-027: `fid claim list --by` is a read filter, not an actor flag

- **Date:** 2026-06-21
- **Context:** `06_DESIGN_PLAN_CLI.md` defines `fid claim list [--stale] [--by
  <role>] [--format ...]` and describes `--by` as a filter by claiming role.
  Phase 10 permission enforcement applies to write commands via acting roles;
  `claim list` is read-only.
- **Decision:** Treat `fid claim list --by` as a filter over stored
  `claimed_by`, not as an acting role. The command does not call
  `policy.Authorize`, does not require the role to be known, and appends no
  journal/activity rows. Unknown filter strings return zero rows unless an
  existing read-filter pattern forces a contrary decision, in which case
  Developer must escalate before implementation.
- **Consequences:** This keeps read surfaces inspectable by agents before they
  act and avoids overloading `--by` with two meanings on the same command.
- **Refs:** PRD-CLAIM-007; `06_DESIGN_PLAN_CLI.md` `fid claim list`;
  `decisions_principles.md` read-ordering rule; phase-11/task-3.

### DEC-026: Claim schema and matrix handling at kickoff

- **Date:** 2026-06-21
- **Context:** Phase 10 intentionally kept the in-code matrix limited to
  implemented write actions. Phase 11 now implements claim writes. The claims
  table/triggers already exist from v0.1a, and `03_TRD_SCHEMA.md` section 13
  already contains the `issue.claim.*` default-profile rows.
- **Decision:** Phase 11 adds no schema migration by plan. Task 1 appends the
  four claim actions to the in-code matrix by copying the existing section 13
  role booleans. `03_TRD_SCHEMA.md` is edited only if implementation proves the
  current spec wrong; otherwise it receives a verified-no-change verdict at wrap.
- **Consequences:** Matrix growth remains per-implemented-action, while the
  schema/spec remains the durable source for default cell values.
- **Refs:** TRD-ARCH-010; TRD-ARCH-011; `03_TRD_SCHEMA.md` sections 11.1 and 13;
  phase-11/task-1.

## Carry-Forward Review At Kickoff

| ID | Disposition | Rationale |
|---|---|---|
| CF-008 | Re-defer to v0.2 | Resolution reason is a product/data-model enhancement unrelated to claims. |
| CF-011 | Re-defer broad reconciliation to v0.1e; schedule claim renewal projection now | PRD-CLAIM-008a explicitly requires renewal to project to `activity_log`, so Phase 11 handles that per-command case without widening the broader timeline policy. |
| CF-013 | Re-defer to v0.1e | `complete.go` orchestration cleanup belongs with `fid dev finish`/`block`, not primitive claims. |

## Wrap-Up Dispositions (Tech Lead, 2026-06-22)

### Validation contract outcome audit

All fifteen assertions reach a terminal state.

- VAL-001..013 and VAL-015: PASS in `qa/final-review.md` with cited test/CLI
  evidence. No disposition needed.
- VAL-014: was CARRY-FORWARD to the human checklist in final review; the
  generated `qa/human-checklist.md` (HV-01a..HV-01e) is now complete with
  overall **PASS** (completed 2026-06-22). Final outcome: **PASS**. No
  FAIL/BLOCKED/CARRY-FORWARD remains; no defer/cancel/escalation required.

### QA recommendation disposition

Final review had no Blocking or Major findings. One Improvement finding:

- **File-size tripwire** — `cmd/fid/issue.go` (984 lines) and
  `internal/journal/replay.go` (485 lines) exceed the 400-line file target
  (AGENTS.md / `08_QUALITY_PLAN.md`). **Disposition: defer** to new carry-forward
  **CF-014** (v0.1e structural cleanup, alongside CF-013). Not do-now: behaviour
  gates passed, the QA plan defines no file-size gate, and the split is a
  behaviour-preserving refactor better batched with the `fid dev *` composites
  that also touch these paths. No spec impact.

### DEC / deviation spec-impact walk

- DEC-026 (no schema migration; matrix copies §13 rows): spec-impacting? No.
  `03_TRD_SCHEMA.md` §11.1 claims table and §13 claim rows already correct and
  match the in-code matrix booleans — `verified-no-change`.
- DEC-027 (`claim list --by` is a read filter): spec-impacting? No.
  `06_DESIGN_PLAN_CLI.md` already documents `--by` as "Filter by claiming role"
  and the JSON shape matches — `verified-no-change`.
- DEC-028 (lock-aware projection freshness): spec-impacting? No. It enforces the
  existing `02_TRD.md` §4.3/§4.9 write-path/lock contract and ADR-0001; Task 5
  report records verified-no-change for `02_TRD.md`, `03_TRD_SCHEMA.md`, and
  `06_DESIGN_PLAN_CLI.md`.
- Developer deviations: none reported across Tasks 1-5.
- Validation-contract deviations: none; VAL-015 was an additive ESC-008
  assertion, not a weakening.

## Spec Sync Verdict

Record a verdict for every spec at phase wrap-up. `verified-no-change` is a
positive assertion after review, not the default.

| Spec | Verdict | Evidence / Reference |
|---|---|---|
| `00_PROJECT_CHARTER.md` | N/A | No charter-level scope/constraint change in a claims slice. |
| `01_PRD.md` | verified-no-change | PRD-CLAIM-001..008a, PRD-ERR-003, PRD-OPS-003, PRD-NFR-008 all MET as written (final-review Requirement Compliance). |
| `02_TRD.md` | verified-no-change | DEC-028 enforces existing §4.3/§4.9 write-path/lock contract; TRD-ARCH-010/011 satisfied by VAL-001/005/013/015. |
| `03_TRD_SCHEMA.md` | verified-no-change | DEC-026; §11.1 claims table + §13 claim rows match the in-code matrix (VAL-001 PASS). |
| `04_TRD_CONFIG.md` | N/A | No `.fid/` directory or config change. |
| `05_DESIGN_PLAN.md` | verified-no-change | Cross-surface claim fields consistent (VAL-008 PASS). |
| `06_DESIGN_PLAN_CLI.md` | verified-no-change | `claim`/`unclaim`/`claim list` flags, `--by` filter (DEC-027), default TTL 1800, and JSON shape match implementation (VAL-002/007/008). |
| `07_DESIGN_PLAN_FLOWS.md` | verified-no-change | create→claim→transition→complete→unclaim lifecycle integration matches documented flow (VAL-011/012). |
| `08_QUALITY_PLAN.md` | verified-no-change | Phase gate run as defined; no new gate added. File-size finding deferred to CF-014, not a spec change. |
| `09_AI_WORKFLOW.md` | N/A | No agent-workflow change. |
| `10_IMPLEMENTATION_ROADMAP.md` | verified-no-change | Phase 11 delivered exactly the v0.1d Phase 11 entry; no re-slice. |

## Triage-Log Audit

- [x] Walked every TRG-NNN in `triage-log.md`.
- [x] Confirmed each RESUME verdict was correct in hindsight.
- [x] Flagged any TRG that should have escalated.
- [x] Confirmed autonomous TL-Lite resume remains in effect (DEC-002; no warm-up tally).

Audit note: TRG-001 is the only triage entry. Its verdict was **BLOCK**, not
RESUME, and it escalated to ESC-008 → resolved by DEC-028 (do-now Task 5 +
VAL-015). In hindsight the BLOCK was correct: the nested-lock concern spanned
`internal/claims`, journal/projection freshness, and concurrency, satisfying
trigger 4 (more than one task/subsystem) and trigger 8 (resumed cross-subsystem
work). No RESUME verdicts to confirm; no TRG was under-escalated.

## Workflow Notes

- **What worked well?** TDD plus validate-before-append discipline held across
  every claim reject path. The Phase 10 matrix seam let Task 1 add only the four
  implemented claim rows with no rework. Mid-phase QA caught the lock-ordering
  risk early enough to convert it into a clean do-now Task 5 rather than a
  post-phase defect, and the additive VAL-015 kept the contract honest.
- **What was inefficient?** ESC-008 surfaced only in the Task 4 completion
  report, after Task 4 was implemented; the nested-lock risk could have been
  anticipated at planning when claims first call `EnsureProjectionFresh` under
  the write lock. The five-task plan also grew to six artifacts (Task 5) mid
  phase, adding a review cycle.
- **Were tasks the right size?** Yes for Tasks 1-4 (one coherent claim concern
  each). Task 5 was correctly scoped as a narrow lock-ordering fix and stayed
  one session.
- **Did the QA review cadence feel right?** Yes. One review per task plus a
  dedicated Review 5 for the ESC-008 fix gave per-task signal; all five reviews
  were CONTINUE with only improvement notes, and final QA closed each note.
- **Were there token-efficiency opportunities missed?** Minor. Reviews 1-3
  flagged sparse task-report gate excerpts that final QA had to re-derive; richer
  inline gate output in the completion reports would have saved a re-run.

### Deep Retrospective Recommendation

No deep retrospective required. Phase 11 closed PASS with no Blocking/Major
findings, all VAL outcomes terminal, and a single deferred improvement (CF-014).
Carry the lock-ordering lesson forward into Phase-12 planning: any new write
command that refreshes projection under the write lock must use the held-lock
freshness helper from the start (DEC-028).
