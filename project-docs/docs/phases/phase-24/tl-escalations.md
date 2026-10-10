# Phase 24 Escalations

## ESC-24-009 - Task 13 R13 required command contradicts its required result

- **Raised:** 2026-10-11 by TL-Lite in TRG-24-017 (BLOCK); ledger STOP
  2026-10-11 10:43.
- **Protected category:** protected required-command definition /
  specification contradiction.
- **Originating task:** task-13.
- **Trigger signature:** task-13:qa-rework:review-11-r13-command-contradiction
  (attempt 2). The attempt-1 chain
  `task-13:qa-rework:review-11-r13-inventory-receipt` stays in history.
- **Condition:** Review-11 F1 asked for Task 13's owned R13 receipt. The
  10:26 TASK-RECOVERY told the Developer to run
  `rg -n "p23CheckJudgmentQualifiers\|PRIVATE_" cmd` and observe exit 0
  with three qualifier references. Those bytes printed nothing and exited
  1 (report `f446a3be`). The plain-`|` spelling printed three references
  and exited 0. The directive forbade changing the command.
- **Status:** RESOLVED 2026-10-11 by DEC-291 (decision). There was no
  contradiction in the accepted command. The R13 row is a Markdown table
  cell, where `\|` is the escape for `|`; rendered, R13 is
  `rg -n "p23CheckJudgmentQualifiers|PRIVATE_" cmd`. The escaped bytes
  were copied out of the cell into the ledger directive. Every earlier
  R13 count, and the owned receipts for R2, R6 and R11, used the rendered
  form. The inventory header now states the rule. Review-11 is not
  edited. The Runner resumes Task 13 at attempt 3 through the clause 6
  route, then a fresh review-12 re-reviews F1. The Tech Lead dry run
  confirmed the route through to `ADVANCEABLE`, plus nine negative
  guards. No further Tech Lead or Human input is needed unless an append
  is rejected as written or the receipt does not show the required
  result.
- **Links:** [DEC-291](decisions.md); [TRG-24-017](triage-log.md);
  [run ledger](run-ledger.md); [Task 13 report](tasks/task-13-done.md);
  [review-11](qa/review-11.md); [inventory](input-inventory.md).

## ESC-24-008 - Task 14 terminal-success state has no permitted Human provenance

- **Raised:** 2026-10-11 by TL-Lite in TRG-24-016 (BLOCK); ledger STOP
  2026-10-11 08:18.
- **Protected category:** protected product/public-interface/provenance
  decision.
- **Originating task:** task-14.
- **Trigger signature:** task-14:public-interface-decision:terminal-state-provenance-qa-pass
  (attempt 1).
- **Condition:** The VAL-049 journey needs `issue transition --to qa_pass`.
  No Human output or help names any workflow state: `workflow list` prints
  names in Human mode and counts in JSON, and `Resume.issue` has no
  states. DEC-278 item 9 does not allow `setup` or `caller-authored` for
  that value. `TestP24TakeControlTerminalStateProvenance` and the required
  tagged gate failed (exit 1). The Task 14 report at `98d45617` had
  report-lint `OUTCOME BLOCKED` (exit 0). It also flagged a second gap:
  no Human output prints the Issue revision that `--expected-revision`
  needs.
- **Status:** RESOLVED 2026-10-11 by DEC-290 (decision). `resume` gains
  the Issue's revision and next workflow states in both modes. JSON
  `Resume.issue.next_states` is `{state, category}[]`. The Human Issue
  block adds `Revision:` and `Next states:` lines, and the
  `issue transition` help points `--to` and `--expected-revision` at
  `resume`. Task 14 owns the change, with a bounded scope and VAL-068.
  The journey takes both values as `fid-output`, and DEC-278's classes
  are unchanged. The Tech Lead dry run confirmed the clause 8 route
  through to `ADVANCEABLE`, plus five negative guards. The Human may
  override the choice by DEC before Task 14 closes. No further Tech Lead
  or Human input is needed unless a clause 8 append is rejected as
  written or the journey finds another value with no class.
- **Links:** [DEC-290](decisions.md); [TRG-24-016](triage-log.md);
  [run ledger](run-ledger.md); [Task 14 report](tasks/task-14-done.md);
  [Task 14 spec](tasks/task-14.md).

## ESC-24-007 - Window D2 Human walkthrough and interim judgment absent

- **Raised:** 2026-10-09 by TL-Lite in TRG-24-013 (BLOCK). Record-identity
  corrections TRG-24-014 and TRG-24-015 followed; TRG-24-015 governs.
  Ledger STOP 2026-10-09 16:28.
- **Protected category:** protected Human product judgment / external
  Human authority.
- **Originating task:** task-11 (checkpoint attribution only; review-9
  names no owning task).
- **Trigger signature:** task-11:qa-window-d2:human-walkthrough:val-054-055-pending
  (attempt 1).
- **Condition:** Review-9 (`b9083f3b`) was BLOCK with one finding, B1-H.
  The deterministic D2 checks and the fixture regeneration at Task 11
  close `93a1f792` passed. QA plan section 8 requires the interim
  VAL-054/055 Human walkthrough at window D2, but
  `qa/human-walkthrough-after-task-11.md` and its nine Human rows
  (HUMAN-02a..02e, HUMAN-03a..03d) were missing. Deterministic evidence
  cannot stand in for Human judgment. Two Runner TRIAGE candidates were
  rejected before append and are preserved in TRG-24-015.
- **Status:** RESOLVED 2026-10-10 by DEC-288 (decision). The Human recorded
  his own answers at `a23b9f33`, clarified on his instruction at
  `fcc60def`, bound to manifest
  `.cache/phase-24-human/fixtures-2268063598/manifest.json`. All five
  refusals were clear from the output alone, and all four views were
  accepted with reasons. No row is rejected or unclear. The verifier has
  no direct route from a STOP to a REVIEW, so the Runner resumes through
  an evidence-only Task 11 recovery: attempt 2, `qa-rework`,
  `task-scope`, no Developer dispatch. That runs RUN-RESUME,
  TASK-RECOVERY, REPAIR RESOLVED and TASK COMPLETE, then a fresh QA
  `review-10` re-review of review-9 closes B1-H. The Tech Lead's dry run
  confirmed every receipt. The interim answers don't replace the
  phase-end HUMAN-02/03 rows. No further Tech Lead or Human input is
  needed unless a DEC-288 append is rejected as written or review-10
  does not pass.
- **Links:** [DEC-288](decisions.md); [TRG-24-013..015](triage-log.md);
  [run ledger](run-ledger.md); [review-9](qa/review-9.md);
  [Human walkthrough](qa/human-walkthrough-after-task-11.md).

## ESC-24-006 - Plain-copy run.inspect success conflicts with retained repository refusal

- **Raised:** 2026-10-09 by TL-Lite in TRG-24-012 (BLOCK); ledger STOP
  2026-10-09 08:55.
- **Protected category:** protected public-interface/contract decision;
  specification contradiction.
- **Originating task:** task-11.
- **Trigger signature:** task-11:public-interface:plain-copy-run-inspect-repository-view
  (attempt 1). The separate session-continuation recovery
  (`task-11:session-continuation:developer-session-interrupted`, attempt 1)
  is unchanged.
- **Condition:** Task 11 and VAL-041 required every recorded-versus-current
  Human view, `run.inspect` included, to print `different worktree` for a
  plain copy. Contract section 2 keeps `run.inspect`'s Phase 21-23
  behaviour, and section 10 keeps JSON unchanged. Built-binary plain-copy
  `run inspect --repository <copy>` refused in both modes with
  `recovery-state-conflict` (`repository identity changed`, exit 4).
  `InspectRun` reconciles first (`run_inspect.go:45`;
  `reconciliation.go:53`; `ac5e6f9`, phase-22/task-4). Both focused
  `TestP24HumanFixtures` runs failed (exit 1). The Task 11 report at
  `9ce8896` had report-lint `OUTCOME BLOCKED` (exit 0).
- **Status:** RESOLVED 2026-10-09 by DEC-287 (decision). Option A: keep the
  accepted refusal, with no behaviour or JSON change. Under accepted
  contracts, a successful `run.inspect` only ever says `same worktree`.
  VAL-041 proves `different worktree` on `resume`. The
  `<state>-different-run.inspect` fixture row asserts the exact clause 2
  refusal pair (exit 4). Option B (inspect a plain copy) is CF-046(c) and
  CF-026a. Tech Lead dry runs confirm that the clause 7 RUN-RESUME and
  TASK-RECOVERY each append `APPEND_CONFIRMED state=REPAIR_PENDING`. No
  further Tech Lead or Human input is needed unless a clause 7 append is
  rejected as written.
- **Links:** [DEC-287](decisions.md); [TRG-24-012](triage-log.md);
  [run ledger](run-ledger.md); [Task 11 report](tasks/task-11-done.md);
  [Task 11 spec](tasks/task-11.md).

## ESC-24-005 - Root-flag refusal needs public-interface mapping authority

- **Raised:** 2026-10-08 by TL-Lite in TRG-24-011 (BLOCK); ledger STOP
  2026-10-08 05:24.
- **Protected category:** protected public-interface/contract decision.
- **Originating task:** task-10.
- **Trigger signature:** task-10:public-interface:unavailable-root-flag-failure-mapping
  (attempt 1). The DEC-285 activity-mapping recovery stayed at attempt 2.
- **Condition:** Built-binary
  `fid [--json] --project <dir> export import … --repository <dir>` printed
  `error: --project is not available for export import until its
  fid-cli.v1 transport is migrated` and exited 7 in both modes. stdout was
  empty and no file changed, but there was no `fid-cli.v1` failure and no
  `reason:` line. The guard is `rejectUnavailableRootFlags`
  (`cmd/fid/root.go`; `cd53c43`, phase-20/task-2), and `export import`
  refuses `--project` by registry (`5ec69e6`, phase-22/task-15). The Task 10
  report at `e0495b5` had report-lint `OUTCOME BLOCKED` (exit 23). The
  seventeen missing registry rows and unfinished fixtures were ordinary
  Task 10 work, not part of this escalation.
- **Status:** RESOLVED 2026-10-09 by DEC-286 (decision). Phase 21 section 2
  makes a forbidden flag `invalid-argument` (invalid-input, exit 2, subject
  null) with the offending spelling and the row's alternatives. Task 9a
  supplies the operation and the unknown-flag `allowed` computation, and
  section 9 supplies the Human lines. DEC-286 fixes the remaining values:
  `value` is the flag spelling, and the message is `--project is not
  available for export import` (the false Phase 20 clause is dropped). The
  scope is `root.go`, a pure `canonical_request.go` extraction and three
  test files. `sql` and `journal inspect` keep their raw refusal. Tech Lead
  dry runs confirm that the clause 8 RUN-RESUME and TASK-RECOVERY each
  append `APPEND_CONFIRMED state=REPAIR_PENDING`. No further Tech Lead or
  Human input is needed unless a clause 8 append is rejected as written.
- **Links:** [DEC-286](decisions.md); [TRG-24-011](triage-log.md);
  [run ledger](run-ledger.md); [Task 10 report](tasks/task-10-done.md);
  [Task 10 spec](tasks/task-10.md).

## ESC-24-004 - Activity list failure mapping needs public-interface authority

- **Raised:** 2026-10-07 by TL-Lite in TRG-24-010 (BLOCK); ledger STOP
  2026-10-07 20:50.
- **Protected category:** protected public-interface/contract decision.
- **Originating task:** task-10.
- **Trigger signature:** task-10:public-interface:activity-list-failure-envelope-mapping
  (attempt 1).
- **Condition:** In both modes, built-binary `activity list --since invalid`
  and `--on invalid` printed raw prose and exited 7. With no project the
  command was silent and exited 7. Neither mode emitted a `fid-cli.v1`
  failure, so TRG-24-009's renderer adoption had no JSON failure to
  preserve. The Task 10 report at `9c3e003` had report-lint
  `OUTCOME BLOCKED` (exit 11). The report also found `stats.go` and
  `issue_ac.go` forcing JSON failures in Human mode, with no edit scope.
- **Status:** RESOLVED 2026-10-07 by DEC-285 (decision). No mapping is
  invented. `activity.list` is a Phase 21 section 4.1 leaf, so the raw
  exit-7 output is a Phase 21 section 1/7 nonconformance, not an accepted
  contract. It now fails through `issueReadFailure`, the retained-read
  mapper its section 4.1 siblings already use, and DEC-285 clause 2 gives
  the exact rows. Task 10 also gains `stats.go` and `issue_ac.go`; there,
  the forced `true` becomes the live JSON mode. The recovery identity moves
  to TRG-24-010's trigger at attempt 2. Tech Lead dry runs confirm the
  clause 8 RUN-RESUME and TASK-RECOVERY each append
  `APPEND_CONFIRMED state=REPAIR_PENDING`. No further Tech Lead or Human
  input is needed unless a clause 8 append is rejected as written.
- **Links:** [DEC-285](decisions.md); [TRG-24-009, TRG-24-010](triage-log.md);
  [run ledger](run-ledger.md); [Task 10 report](tasks/task-10-done.md);
  [Task 10 spec](tasks/task-10.md).

## ESC-24-003 - Protected sequence validator rejects corrected Task 9 recovery

- **Raised:** 2026-10-07 by TL-Lite in TRG-24-007 (triage log only; the
  ledger admits no TRIAGE BLOCK while a TRIAGE REPAIR is pending, so no
  STOP was recorded).
- **Protected category:** protected canonical workflow verifier authority /
  orchestration integrity.
- **Originating task:** task-9.
- **Trigger signature:** task-9:workflow-control:append-only-triage-category-override
- **Condition:** TRG-24-006 (REPAIR, attempt 2) recorded category
  `scope-gap`, which is not a ledger recovery category, and minted a new
  trigger while Task 9's report-repair continuation was active. An
  in-section correction changed only the category. Both rejected candidates
  are below; the ledger stayed unchanged at `VALID state=REPAIR_PENDING`.
  - Corrected candidate: `does not preserve its governing schema v2 TRIAGE
    REPAIR authority` and `continuing recovery must preserve task and
    Trigger signature` (`exit_code: 1`).
  - Original candidate: `missing or invalid Category`, a Runner-side
    disposition error, and the same continuation message.
- **Status:** RESOLVED 2026-10-07 by DEC-284 (decision). No verifier change
  is needed. The accepted DEC-190 correcting-TRIAGE route binds the repair.
  TRG-24-008 (Tech Lead) keeps the chain trigger
  `task-9:report-lint:noncanonical-pass-evidence-termination` at attempt 2,
  with category `in-scope-implementation` and authority
  `task spec; DEC-284`. Product scope is unchanged. Tech Lead dry runs
  confirm each step appends `APPEND_CONFIRMED state=REPAIR_PENDING`. The
  Runner follows DEC-284 clause 4, then runs the missed QA window C after
  Task 9 (clause 6). No further Tech Lead or Human input is needed unless
  a clause 4 append is rejected as written.
- **Links:** [DEC-284](decisions.md); [TRG-24-006..008](triage-log.md);
  [run ledger](run-ledger.md); [Task 9 report](tasks/task-9-done.md).

## ESC-24-002 - Accepted Task 6 RUN-RESUME binds an unsatisfiable recovery tuple

- **Raised:** 2026-10-06 by the Phase Runner (reported to the Human).
  Allocated by the Tech Lead in DEC-283. The ledger admits no TRIAGE event
  while resume-pending, so no TRG exists (as with DEC-149 clause 6).
- **Protected category:** protected verifier (`workflow-sequence-lint.sh`)
  defect; orchestration integrity.
- **Originating task:** task-6.
- **Trigger signature:** task-6:session-liveness:developer-role-capacity-after-bounded-retry
- **Condition:** the 2026-10-05 20:48 RUN-RESUME records trigger
  `task-6:session-liveness:...` at attempt 3. Its `Next recovery` names
  `task-6:session-continuation:...` at attempt 1. The verifier binds the
  next TASK-RECOVERY to both, so no candidate validates. Rejected
  candidates: a REPAIR (`RUN-RESUME must be followed by its bound
  TASK-RECOVERY`) and a TASK-RECOVERY (tuple mismatch). The ledger is
  unchanged and still `VALID state=REPAIR_PENDING`.
- **Status:** RESOLVED 2026-10-06 by DEC-283 (decision). The Runner then
  carries out the clause 4 verifier maintenance and the clause 6 ledger
  sequence; no further Tech Lead or Human input is needed unless the
  verification items cannot be produced as stated.
- **Links:** [DEC-283](decisions.md); [run ledger](run-ledger.md);
  [Task 6 report](tasks/task-6-done.md).

## ESC-24-001 - Task 6 Developer role unavailable after bounded model-capacity retry

- **Raised:** 2026-10-05 by TL-Lite in TRG-24-003.
- **Protected category:** unrepairable external infrastructure/capability
  failure after bounded retry.
- **Originating task:** task-6.
- **Trigger signature:** task-6:session-liveness:developer-role-capacity-after-bounded-retry
- **Attempt count:** 2 interrupted spawn attempts. Neither returned a typed
  result, so neither consumes the no-progress budget.
- **Condition:** the selected Developer model reported capacity exhaustion
  on two spawns.
- **Status:** RESOLVED 2026-10-05 20:48. The Runner verified the restored
  capability:
  `restored-capability:developer-role-spawn-01a10b08-efe6-7050-8f88-65c524e62519`.
  The resumed session committed `1baaecc` and `7298964`. DEC-283 records
  the closure and reconciles the resume's recovery tuple.
- **Links:** [TRG-24-003](triage-log.md); [DEC-283](decisions.md).
