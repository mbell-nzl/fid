# Phase 20 - TL-Lite Triage Log

> **Purpose.** Append-only log of mid-phase triage events handled by TL-Lite.
> Phase Runner reads TL-Lite's `VERDICT:` line to decide whether to resume or
> stop. Tech Lead audits this log at phase wrap-up.
>
> **Numbering.** Identifiers are phase-scoped: `TRG-20-NNN` and `ESC-20-NNN`.
>
> **This file is the allocation source for both.** Read this file, take the
> highest `NNN` already recorded in that series, plus one; start at `001` when
> empty. Never recall a number from memory, never renumber a published entry,
> and never emit a bare placeholder.
>
> **Write the entry before reporting the verdict.** Newest entry at top.

| TRG-20-NNN | Date | Trigger | Source | Verdict | ESC-20-NNN (if BLOCK) | Summary |
|---|---|---|---|---|---|---|
| TRG-20-014 | 2026-08-16 | Task 16's fresh vendor-mode module-absence output includes a module-header line that its current contract declares a real VAL-032 failure | `tasks/task-16-done.md`; attached independent command and report-lint evidence | BLOCK | ESC-20-014 | DEC-144 corrects the parenthesized vendor-mode wording but leaves the mandatory module-header line outside Task 16's one-line complete-output pass condition; reconciling that protected specification/quality-gate contract is required before Task 16 can resume. |
| TRG-20-013 | 2026-08-16 | Task 16 cannot satisfy DEC-143's exact `go mod why -vendor` output literal; amended report-lint rejects valid declaration explanations | `tasks/task-16-done.md`; attached commit and command evidence | BLOCK | ESC-20-013 | DEC-143's replacement build-graph command has a documented exact-literal contradiction with the observed Go output; correcting that protected specification/quality-gate authority conflict is required before the separate parser regression can be a local repair. |
| TRG-20-010 | 2026-08-15 | Task 10 blocker report, focused suite failure, and report-lint failure | `tasks/task-10-done.md`; attached failure output | BLOCK | ESC-20-010 | Genuine unresolved product-facing command-test, required-gate, and report-lint failures; no verified corrective commit or accepted exception applies. |
| TRG-20-009 | 2026-08-12 | Review 3 returned three unresolved Task 9 evidence findings | `qa/review-3.md` | BLOCK | ESC-20-009 | Genuine unresolved QA Blocking findings; supplied review verification and integration PASS do not resolve missing deterministic evidence matrices. |
| TRG-20-008 | 2026-08-11 | Task 8 report marks six required gates PASS without nested exit code or explicit termination evidence | `tasks/task-8-done.md`; report-lint output | BLOCK | ESC-20-008 | Genuine unresolved evidence failure: commit and report-lint verification do not prove command termination, and DEC-138 does not apply. |
| TRG-20-007 | 2026-08-11 | Task 5 Developer blocker report: required dedicated integration gate failed because the fixed classification inventory omitted the newly added parity-fixture rebuild journey | `tasks/task-5-done.md` | BLOCK | ESC-20-007 | Genuine unresolved required-gate failure; no verified completion commit or accepted decision provides a RESUME ground. |
| TRG-20-006 | 2026-08-11 | Task 4 recovery completion report fails report-lint after verified recovery commit: three invalid go-test evidence tokens and pending commit fields | `tasks/task-4-done.md`; report-lint output | BLOCK | ESC-20-006 | Genuine unresolved report-lint hard stop in the committed report; attached verification proves the faulty commit, not a correction. |
| TRG-20-005 | 2026-08-11 | Task 4 Developer blocker report: required golangci-lint gate failed on unused `describeLockHolder`; no commit and downstream gates not run | `tasks/task-4-done.md` | BLOCK | ESC-20-005 | Genuine unresolved required-gate failure. No verified commit or accepted decision supplies a RESUME ground. |
| TRG-20-004 | 2026-08-10 | Infrastructure Review 1 found a retained integration-tagged in-process test outside both accepted gates | `qa/infra-review-1.md` | BLOCK | ESC-20-004 | Task 3f's direct-parent inventory omitted `internal/issues/transition_rules_test.go`; repository-wide classification and ordinary ownership are required before Task 4. |
| TRG-20-003 | 2026-08-09 | Task 3d blocker report: required full-suite coverage gate has no captured terminal output and report-lint fails on three go-test prose false matches | `tasks/task-3d-done.md` | BLOCK | ESC-20-003 | The required gate is UNVERIFIED, not PASS; the new report-lint control also has an unresolved three-failure false-match defect. |
| TRG-20-002 | 2026-08-09 | Bounded QA re-review review-1b returned BLOCK with two unresolved findings before Task 4 | `qa/review-1b.md` | BLOCK | ESC-20-002 | VAL-043 has an exhaustive evidence gap and the required integration gate fails in three init scenarios. |
| TRG-20-001 | 2026-08-09 | QA Review 1 has two unresolved Blocking evidence gaps before Task 4 can start | `qa/review-1.md` | BLOCK | ESC-20-001 | VAL-002 and VAL-043 lack required exhaustive boundary evidence across Tasks 1-3. |

---

## TRG-20-014 - Task 16's complete-output rule conflicts with the observed vendor-mode result

- **Trigger:** the second Task 16 continuation's fresh `go mod why -vendor`
  output has an extra module-header line, while the current Task 16/DEC-144
  contract treats any extra line beyond one parenthesized does-not-need line as
  a real VAL-032 failure.
- **Trigger signature:** `task-16:specification-contradiction:go-mod-why-vendor-complete-output`
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-16-done.md`;
  attached independent `go mod why -vendor -m gopkg.in/yaml.v3` and report-lint
  evidence; attached verification for `8d347a3a959716696472b5432241dcb4f75eb80e`
  and `afb262ddc3636f6a4c9a05d6b9e77b4aeac276eb`.
- **Candidate:** blocker candidate
- **Attempt count:** 2 same-task continuations for this Task 16 control;
  the complete-output contradiction is newly exposed by the second.
- **Checkpoint:** Product deletion is committed at `afb262d`; the report update
  is committed at `8d347a3`; independent report-lint exits 0 with
  `OUTCOME BLOCKED`. The fresh exact command exits 0 with two lines:
  `# gopkg.in/yaml.v3` followed by the DEC-144 vendor-mode does-not-need line.
- **Verdict:** BLOCK
- **Reasoning:** DEC-144 resolves the earlier parenthesized-line wording
  conflict only. Task 16 and the QA plan still require the *whole output* to be
  one parenthesized line and declare any extra line a real VAL-032 failure.
  The tool's module-header line therefore leaves a binding output-contract
  contradiction; accepting it, changing the command/evidence, or weakening
  the assertion requires protected specification/quality-gate authority.
- **Local decisions made:** none.
- **Resolution evidence:** No recorded decision accepts the documented
  two-line vendor-mode result. The attached commit verification establishes the
  blocker report and product state, not resolution of the contradiction.
- **If BLOCK:** category specification contradiction / quality-gate authority;
  ESC-20-014 raised in `tl_escalations_current.md`; resolution condition: full
  Tech Lead reconciles the Task 16, QA-plan, and DEC-144 evidence contract
  without weakening VAL-032, then directs the same-task continuation.
- **Human attention required:** yes.
- **References:** `tasks/task-16.md:121-128`; `tasks/task-16-done.md:5-10`,
  `73-87`, `96-100`; `decisions.md` DEC-144; `qa/qa-plan.md:167`;
  `qa/validation-contract.md:55,97`; `project-docs/specs/09_AI_WORKFLOW.md`
  §§2.3, 9.

## TRG-20-010 - Task 10 partial migration has unresolved command and report-gate failures

- **Trigger:** Task 10 Developer produced a blocker report because the focused
  command suite remains red; no task commit or report-finalization commit
  exists, and the full required gate was not run. The required report-lint also
  failed.
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-10-done.md`; prompt
  failure evidence for `project-docs/tools/report-lint.sh` and
  `go test ./cmd/fid -count=1`.
- **Verdict:** BLOCK
- **Reasoning:** The focused package command terminated `exit_code: 1` with
  `FAIL` because AC/label/e2e assertions still expect obsolete stderr-success
  output and `--format` behavior; this is a product-command and behavioral-test
  failure, not a stale verification inventory. Report-lint independently
  terminated `exit_code: 9` after nine failures, including missing
  PASS/BLOCKED required-command results and no commit hash. The report confirms
  that the full gate was not run and that no task or report-finalization commit
  exists; no attached verification evidence or accepted decision supplies a
  RESUME ground. Completing the contract migration and its DEC-140 matrix spans
  several Task 10 command/test surfaces, exceeding local authority.
- **Local decisions made:** none
- **References:** `tasks/task-10-done.md` (Tests run; Exact quality-gate
  results; Concerns; Commit hash and subject);
  `project-docs/docs/current/phase_brief.md` (Task 10; Active Constraints);
  `AGENTS.md` (Commands; Hard Stops); `project-docs/specs/09_AI_WORKFLOW.md`
  §2.3, §8, §9.
- **If BLOCK:** ESC-20-010 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve ESC-20-010 and direct bounded
  Task 10 recovery, including complete DEC-140 evidence, exact required-gate
  execution, commits, and clean report-lint before Phase Runner resumes.

## TRG-20-009 - Review 3 has unresolved Task 9 evidence gaps

- **Trigger:** Phase 20 Review 3 returned BLOCK with three bounded Task 9
  evidence findings: missing deterministic transition expected-revision
  absent/matching/stale coverage; incomplete canonical transition
  outcome/detail/zero-append/Event-identity matrix; and incomplete owner-flag
  plus attribution negative-state assertions.
- **Source:** `project-docs/docs/phases/phase-20/qa/review-3.md`; attached
  verification for review commit `0da39ddd89dd645474dbb337cb21ff2112c4a96a`
  and the independent integration result
  `ok  github.com/mbell-nzl/fid/cmd/fid/integration  8.949s` (terminated
  `exit_code: 0`).
- **Verdict:** BLOCK
- **Reasoning:** The three findings are explicit QA Blocking findings against
  Task 9's deterministic transition, closed-request, and negative-state
  evidence. They are not misclassified. The supplied verification proves the
  Review 3 record and a broad integration command, not that the missing
  contract-specific assertions have been added in a verified corrective commit;
  no such commit is supplied. No accepted recorded decision closes or weakens
  these requirements. The only permitted QA-blocking RESUME grounds therefore
  do not apply.
- **Local decisions made:** none
- **References:** `qa/review-3.md` (Task 9 Blocking findings; Gate And Evidence
  Checks); `project-docs/specs/09_AI_WORKFLOW.md` §2.3 (TL-Lite QA-review
  triage); `AGENTS.md` (Hard Stops; deterministic evidence rules).
- **If BLOCK:** ESC-20-009 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve ESC-20-009 and direct bounded
  Task 9 evidence recovery plus the required re-review before Phase Runner
  resumes.
- **Tech Lead adjudication (DEC-140, 2026-08-12):** verdict sustained as correct
  triage on all three findings, each independently confirmed against the tree.
  Unlike `TRG-20-008` this resolves against the evidence record, not the Tech
  Lead: the requirement was published in the task spec, the validation contract
  rows the Developer prompt requires loading, and the phase brief's DEC-130
  Known Risk. Task 9 closes at `86d927d` and is not reopened; `tasks/task-9a.md`
  recovers the missing assertions and `qa/review-3b.md` re-reviews them.

## TRG-20-008 - Task 8 has unverified gate-termination evidence

- **Trigger:** Task 8 completion report records all six required quality gates
  as PASS, but the exact quality-gate rows omit nested command `exit_code` or
  an explicit termination result, so required gate termination evidence is
  UNVERIFIED under `AGENTS.md`.
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-8-done.md`; report-lint
  output; attached commit-verification evidence for implementation commit
  `4df4d6690a0562873c570ddaeba2140ad3adef81` and report-finalization commit
  `ce282c25c617f103b4338b0de136318f82f1b8c3`.
- **Verdict:** BLOCK
- **Reasoning:** The report records PASS tokens but no nested `exit_code` or
  explicit termination signal for any required gate. The prompt's verified
  commits and clean report-lint result prove report/commit state, not that the
  six commands terminated; outer completion and attestation are explicitly
  insufficient. No accepted decision resolves it: DEC-138 applies only to a
  named stale verification inventory and expressly leaves uncaptured tokens as
  failed gates. This is therefore a genuine unresolved quality-gate failure.
- **Local decisions made:** none
- **References:** `tasks/task-8-done.md` (Exact quality-gate results);
  `AGENTS.md` (Commands: termination and UNVERIFIED rules);
  `project-docs/specs/09_AI_WORKFLOW.md` §2.3, §8 (nested termination and
  evidence rules); `decisions.md` DEC-138 decision 5.
- **If BLOCK:** ESC-20-008 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve ESC-20-008 and direct exact-
  command evidence recovery before Task 8 or later work resumes.
- **Tech Lead adjudication (DEC-139, 2026-08-11):** verdict sustained as correct
  triage — the rule was in force and unmet, and DEC-138's fourth RESUME ground
  correctly does not reach an evidence record. The escalation resolves against
  the Tech Lead, not the Developer: the `exit_code` requirement existed only in
  `project-docs/agents/phase-runner.md`, while `AGENTS.md`, `09_AI_WORKFLOW.md`
  §8, and the completion-report template required only the verbatim token, which
  Task 8 supplied for all six commands. Task 8 closes at `4df4d66` with no
  evidence recovery; the requirement is promoted to canon forward-effective.

## TRG-20-007 - Task 5 blocked by incomplete integration-journey inventory

- **Trigger:** Phase 20 Task 5 Developer blocker report says the required
  dedicated integration gate failed. TestIntegrationClassification has a fixed
  journey inventory that omitted the newly added
  TestV01ParityFixtureRebuildPreservesEveryProjectionTable journey. The first
  four exact Task 5 gates passed; benchmark and build were not run; no task
  commit was created.
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-5-done.md`.
- **Verdict:** BLOCK
- **Reasoning:** The dedicated required integration gate genuinely failed; the
  report identifies the omitted journey in the fixed classification inventory.
  This is neither a misclassification nor covered by an accepted decision, and
  the prompt supplies no verified completion commit (the report confirms none
  exists). Correcting the test/gate classification apparatus is outside
  local-decision authority.
- **Local decisions made:** none
- **References:** `project-docs/docs/phases/phase-20/tasks/task-5-done.md`
  (Exact quality-gate results and blocker); `AGENTS.md` (Commands: required
  gate failures are hard stops); `project-docs/specs/09_AI_WORKFLOW.md` §2.3
  and §9.
- **If BLOCK:** ESC-20-007 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve ESC-20-007 before Task 5
  recovery or any subsequent task resumes.

## TRG-20-006 - Task 4 recovery report has invalid committed lint evidence

- **Trigger:** Task 4 recovery completion report fails
  `project-docs/tools/report-lint.sh` after recovery commit
  `e56b3a0ce0433f1044f70b7969c14e19234dfdf8`: the focused test token
  `0.009s`, ordinary race/coverage token `coverage: 80.5% of statements`, and
  dedicated integration token `4.625s` do not look like `go test` output; its
  Commit hash and Commit subject sections still say `Pending commit`.
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-4-done.md`; Phase
  Runner report-lint output and attached commit-verification evidence.
- **Verdict:** BLOCK
- **Reasoning:** This is a genuine unresolved report-lint hard stop. The
  attached evidence verifies that commit
  `e56b3a0ce0433f1044f70b7969c14e19234dfdf8` exists and has subject
  `feat(rebuild): map typed authority and lock failures`, but the prompt also
  establishes that its committed completion report retains all four defects.
  It therefore does not establish the only applicable RESUME ground—an
  already-resolved failure in a verified commit. No accepted decision permits
  invalid evidence tokens or pending commit metadata.
- **Local decisions made:** none
- **References:** `tasks/task-4-done.md` (Tests run; Exact quality-gate
  results; Commit hash; Commit subject); `AGENTS.md` (Commands: evidence-token
  rules); `project-docs/specs/09_AI_WORKFLOW.md` §2.3 (report-lint triage and
  verified-commit rule), §2.4 (report-lint hard stop).
- **If BLOCK:** ESC-20-006 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve the committed Task 4 report
  correction and recovery path before the Phase Runner resumes.

## TRG-20-015 - Task 17 report-lint false-rejects byte-identical absence evidence

- **Trigger:** Task 17's truthful expected-non-zero absence command is
  byte-identical in the task spec and committed report, but report-lint's
  `awk -v` comparison normalizes the required backslash-dot escape and rejects
  it.
- **Trigger signature:** `task-17:quality-gate-control-defect:awk-v-backslash-dot-normalization`
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-17-done.md`; the
  prompt's independent `project-docs/tools/report-lint.sh` output and attached
  commit-verification evidence.
- **Candidate:** `RECOVERY_REQUIRED` / report-lint control defect.
- **Attempt count:** 1 mediated Task 17 report-recovery attempt for this
  signature.
- **Checkpoint:** product commit
  `95ee13d090b424ad541971d7fafe2ce295f09143`, report-finalization commit
  `26fabda6c944a2f3835f5b75ab80c5e2ec68efc3`, and report-recovery evidence
  commit `7c74edcdbd7e2ff53e84c353b89ce4b8db83d7ce` are attached as verified;
  independent report-lint still terminates `exit_code: 2` with
  `OUTCOME RECOVERY_REQUIRED` and two false missing-command failures.
- **Verdict:** BLOCK
- **Reasoning:** The command in Task 17's declaration and report row is
  literally equal at 104 bytes. The control's `awk -v` handling turns its
  required `\\.` escape into `.`, so the rejection is real and cannot be
  repaired by rewording without violating exact command coverage. Repairing
  `report-lint.sh` is outside Task 17 authority and changes a canonical
  quality-gate/workflow control; it requires full Tech Lead authority.
- **Local decisions made:** none.
- **Resolution evidence:** The supplied commit verification proves the product
  and report-recovery checkpoints, not resolution: the independently attached
  report-lint failure remains. The prompt expressly prohibits changing the
  truthful command spelling, command, threshold, assertion strength, or report
  evidence.
- **If BLOCK:** category quality-gate definition / workflow authority;
  ESC-20-015 raised in `tl_escalations_current.md`; resolution condition: full
  Tech Lead authorizes and directs a byte-exact report-lint control repair,
  then routes the resulting recovery without changing Task 17 evidence.
- **Human attention required:** yes.
- **References:** `tasks/task-17.md:45-62`; `tasks/task-17-done.md:45-54,
  62-65`; `project-docs/specs/09_AI_WORKFLOW.md` §§2.3, 9;
  `project-docs/templates/triage-log-template.md`.

## TRG-20-005 - Task 4 blocked by unused lock-holder helper

- **Trigger:** Phase 20 Task 4 Developer blocker report says the required
  `golangci-lint run ./...` gate failed on an unused function at
  `internal/journal/lock.go:185` (`describeLockHolder`); the Developer did not
  commit, and the remaining required gates were not run.
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-4-done.md`
- **Verdict:** BLOCK
- **Reasoning:** This is a genuine unresolved required quality-gate failure,
  which is a Phase Runner hard stop. It is not misclassified; no commit exists,
  so no attached verification evidence can establish resolution; and no
  accepted decision authorizes a lint failure or skipped downstream gates. The
  partial Task 4 recovery must be directed by the full Tech Lead before the
  Phase Runner resumes.
- **Local decisions made:** none
- **References:** `project-docs/docs/phases/phase-20/tasks/task-4-done.md`
  (Exact quality-gate results; Concerns, technical debt, or notes for the Tech
  Lead); `AGENTS.md` (Commands; Hard Stops); `project-docs/specs/09_AI_WORKFLOW.md`
  §2.3 and §9.
- **If BLOCK:** ESC-20-005 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve ESC-20-005 before any Task 4
  recovery or subsequent task resumes.

## TRG-20-004 - Infrastructure Review 1 found an orphaned tagged test

- **Trigger:** The fresh infrastructure review found
  `internal/issues/transition_rules_test.go` still tagged while Task 3f's
  canonical dedicated command runs only `cmd/fid/integration`; the ordinary
  command excludes the file and the classification regression scanned only
  the direct `cmd/fid` parent.
- **Source:** `project-docs/docs/phases/phase-20/qa/infra-review-1.md`
- **Verdict:** BLOCK
- **Reasoning:** This is a cross-task gate-definition regression against
  DEC-133's complete tagged-test accounting. The test's real persistence
  assertion is valid, but no accepted command executes it. The Human-approved
  pack explicitly requires blocking findings to be corrected and returned to
  the same reviewer before Task 4.
- **Local decisions made:** none; correction follows DEC-133's existing rule
  that in-process tests belong to the ordinary suite and that every tagged
  test must be inventoried.
- **References:** `qa/infra-review-1.md`; `decisions.md` DEC-133;
  `infra-improvement-pack.md` Invariants and Final Evidence Contract.
- **If BLOCK:** ESC-20-004 raised in `tl_escalations_current.md`; correction
  commit `75c856a` awaits complete post-fix gates and same-reviewer re-review.

## TRG-20-003 - Task 3d blocked by uncaptured coverage gate and report-lint false matches

- **Trigger:** Task 3d Developer produced a blocker report because the required
  `go test -race -coverprofile=coverage.out ./...` terminated in the session
  without captured terminal output, making its gate UNVERIFIED rather than
  PASS; independent report-lint then failed three checks because the new
  go-test token control matched explanatory `go test` prose in the blocker
  report. Task 3d has no commit.
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-3d-done.md`
- **Verdict:** BLOCK
- **Reasoning:** This is a genuine unresolved required-gate and report-lint
  hard stop. AGENTS.md and DEC-131 require a real terminal token for a
  full-suite `go test` PASS and classify uncaptured output as UNVERIFIED, not
  success. The report-lint failure is a second genuine hard stop; no attached
  verified commit resolves either condition, and DEC-131 requires the new
  control to land rather than remain a promise. The resulting partial-work
  recovery changes delivery tooling and evidence across Task 3d and the five
  corrected completion reports, exceeding local-decision authority.
- **Local decisions made:** none
- **References:** `project-docs/docs/phases/phase-20/tasks/task-3d-done.md`
  (Exact quality-gate results and blocker); `AGENTS.md` (Commands, evidence
  token and UNVERIFIED rules); `project-docs/docs/phases/phase-20/decisions.md`
  (DEC-131, decisions 4-6 and 9); `project-docs/specs/09_AI_WORKFLOW.md` §2.3
  and §9.
- **If BLOCK:** ESC-20-003 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve ESC-20-003 before any Task 3d
  recovery or review-1c work resumes.

## TRG-20-002 - Review 1b unresolved VAL-043 evidence and integration-gate failure

- **Trigger:** bounded QA re-review review-1b returned BLOCK with two new
  unresolved findings before Task 4: VAL-043 lacks explicit rejected
  lower-bound/type/overflow cases for `executor_limits.max_excerpt_bytes`, and
  the exact required integration command fails in three `cmd/fid init`
  scenarios with invalid project initialization argument, contradicting
  task-3b's recorded PASS gate.
- **Source:** `project-docs/docs/phases/phase-20/qa/review-1b.md`
- **Verdict:** BLOCK
- **Reasoning:** These are QA Reviewer Blocking findings and a required-gate
  failure. They are not misclassified, no supplied verified commit resolves
  them, and no accepted decision covers them: DEC-130 applies only to the
  earlier ESC-20-001 recovery and does not weaken or close VAL-043. The repair
  affects strict-config evidence and the integration journey across subsystems,
  so it exceeds local-decision authority.
- **Local decisions made:** none
- **References:** `project-docs/docs/phases/phase-20/qa/review-1b.md`
  (VAL-043 and Task 3b Blocking findings);
  `project-docs/docs/current/tl_escalations_current.md` (ESC-20-001 archive
  note); `project-docs/docs/phases/phase-20/decisions.md` (DEC-130);
  `project-docs/specs/09_AI_WORKFLOW.md` §2.3.
- **If BLOCK:** ESC-20-002 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve ESC-20-002 before Task 4
  starts.

## TRG-20-001 - Review 1 unresolved cross-task evidence gaps

- **Trigger:** QA Review 1 has two unresolved Blocking evidence gaps before
  Task 4 can start: Task 1 does not exhaustively prove exact detail
  value/null/nested schemas and redaction behavior for every registered reason
  code; Tasks 2-3 do not prove the required strict-config and init boundary
  matrix, including one-character/128-byte/MaxInt64 and overflow
  just-inside/outside cases.
- **Source:** `project-docs/docs/phases/phase-20/qa/review-1.md`
- **Verdict:** BLOCK
- **Reasoning:** These are QA Reviewer Blocking findings that are neither
  misclassified, resolved by the supplied verified review commit, nor covered
  by an accepted decision. They leave VAL-002 and VAL-043 without the required
  qualifying/non-qualifying evidence and span completed Tasks 1-3.
- **Local decisions made:** none
- **References:** `project-docs/docs/phases/phase-20/qa/review-1.md` (Task 1
  and Task 2 Blocking findings); `project-docs/docs/phases/phase-20/qa/validation-contract.md`
  (VAL-002, VAL-043); `project-docs/specs/09_AI_WORKFLOW.md` \u00a72.3.
- **If BLOCK:** ESC-20-001 raised in `tl_escalations_current.md`; recommended
  next action: invoke full Tech Lead to resolve the cross-task evidence
  recovery before Task 4 starts.

## TRG-20-010 Recovery adjudication (DEC-141)

- **Original verdict preserved:** `BLOCK (escalated as ESC-20-010)` was the
  literal outcome available under the then-current binary router.
- **Typed reclassification:** `REPAIR` / `in-scope-implementation`.
- **Evidence:** The Developer session ended normally after 616 seconds with no
  kill or timeout; the remaining failures were contract-obsolete Task 10 tests
  and incomplete accepted-scope implementation. No protected decision or
  external capability was required. Stash `779beab1` preserved the exact tree.
- **Resolution:** DEC-141 authorizes the same-task continuation. The v2-marked
  run ledger carries the repair lifecycle and exact completion evidence.
- **Human attention required:** no.

<!--
Phase wrap-up checklist:

- [ ] Walked every TRG-20-NNN above.
- [ ] Confirmed each RESUME verdict was correct in hindsight.
- [ ] Flagged any TRG that should have escalated.
- [ ] Confirmed every entry carries a real, unique, non-placeholder ID, and
      that every BLOCK names the ESC it raised.
- [ ] Confirmed autonomous TL-Lite resume in effect (DEC-002; no warm-up tally).
-->

---

## TRG-20-011 - Task 14 gate cannot close within the assigned Task 14/15 split

- **Trigger:** Task 14 required exact full gate cannot pass after the approved
  eight legacy handlers remove `--format` because remaining `--format` callers
  and formatter-test helpers are explicitly assigned to Task 15.
- **Trigger signature:** `task-14:specification-contradiction:task-14-task-15-format-convergence-split`
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-14-done.md`;
  `project-docs/docs/phases/phase-20/tasks/task-15.md`.
- **Candidate:** blocker candidate.
- **Attempt count:** 1.
- **Checkpoint:** uncommitted Task 14 edits and the Phase Runner-owned
  run-ledger append; no Task 14 commit-verification evidence was supplied.
- **Verdict:** BLOCK.
- **Reasoning:** The Task 14 report records that its exact lint, race/coverage,
  and integration commands fail because legacy formatter helpers and test
  callers became invalid when the eight in-scope handlers removed `--format`.
  Task 15 explicitly owns those remaining caller and test-convergence changes,
  so repairing the failed Task 14 gate would consume another task's assigned
  scope; alternatively changing the Task 14 gate/disposition would alter a
  protected gate definition. This is a binding specification contradiction,
  not an in-scope failed-gate repair.
- **Local decisions made:** none.
- **Resolution evidence:** no accepted decision or verified corrective commit
  clears the contradiction; the prompt supplies no commit-verification
  evidence, so RESUME is unavailable.
- **If BLOCK:** category specification contradiction / protected task-boundary
  and gate-definition authority; ESC-20-011 raised in
  `tl_escalations_current.md`; resolution condition: full Tech Lead reconciles
  the Task 14/15 ordering and scope authority without weakening the six exact
  commands, thresholds, or behavioral assertions.
- **Human attention required:** yes.
- **References:** `tasks/task-14-done.md` (Workflow outcome; Exact
  quality-gate results; Deviations from spec); `tasks/task-15.md` (Objective;
  Scope; Acceptance Criteria; Exact Commands To Run);
  `project-docs/specs/09_AI_WORKFLOW.md` §2.3 and §9; `AGENTS.md` (Commands;
  Hard Stops).

---

## TRG-20-012 - Task 16 YAML absence and report-lint rules cannot both close truthfully

- **Trigger:** Phase 20 Task 16 cannot truthfully close because the QA-plan
  YAML-module absence command and report-lint applicability conflict with the
  accepted no-match evidence rules.
- **Trigger signature:** `task-16:specification-contradiction:yaml-module-absence-and-expected-nonzero-report-lint`
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-16-done.md`;
  `project-docs/docs/phases/phase-20/qa/qa-plan.md` YAML module absence row;
  attached commit and report-lint evidence.
- **Candidate:** typed `RECOVERY_REQUIRED` / report-lint failure.
- **Attempt count:** 1.
- **Checkpoint:** Product commit
  `afb262ddc3636f6a4c9a05d6b9e77b4aeac276eb` is verified by the supplied
  `git rev-parse --verify` result (`exit_code: 0`) and subject match
  `refactor(cli): remove universal formatter framework`; the Task 16 recovery
  evidence is committed through `72e8683`; report-lint returned
  `OUTCOME RECOVERY_REQUIRED` (`exit_code: 3`).
- **Verdict:** BLOCK.
- **Reasoning:** The Task 16 report records that `go mod tidy` retains
  `gopkg.in/yaml.v3` checksum entries because the accepted retained
  `github.com/golang-migrate/migrate/v4@v4.19.1` dependency requires it.
  Therefore the QA-plan command `rg -n "gopkg.in/yaml.v3" go.mod go.sum`
  cannot produce its required exit 1/no-output state without changing a
  dependency or the accepted QA assertion. Separately, the report-lint result
  rejects truthful PASS rows for expected nonzero absence commands, directly
  conflicting with `AGENTS.md`'s required real `exit 1 (no output)` evidence.
  Neither correction is a decision-free Task 16 repair: it requires dependency
  approval or QA/gate-definition and workflow authority.
- **Local decisions made:** none.
- **Resolution evidence:** The verified product commit establishes that direct
  YAML use and the direct `go.mod` requirement were removed, but it does not
  resolve the still-failing accepted QA-plan command or report-lint conflict;
  the supplied report-lint recovery result confirms the issue remains open.
- **If BLOCK:** category specification contradiction / dependency and
  quality-gate-definition authority; ESC-20-012 raised in
  `tl_escalations_current.md`; resolution condition: full Tech Lead reconciles
  VAL-032 and the Task 16 YAML absence evidence with retained migration-module
  requirements, and authorizes a compatible report-lint treatment of explicit
  expected-nonzero absence commands without weakening their real-status rule.
- **Human attention required:** yes.
- **References:** `tasks/task-16.md` (Acceptance Criteria; Exact Commands To
  Run; Notes For Developer); `tasks/task-16-done.md` (VAL-032; Exact
  quality-gate results; Deviations from spec); `qa/qa-plan.md` §5 (YAML module
  absence); `AGENTS.md` (Commands: no-match evidence rule);
  `project-docs/specs/09_AI_WORKFLOW.md` §2.3 (BLOCK and protected triggers).
- **Tech Lead adjudication (DEC-143, 2026-08-16):** verdict **sustained on both
  halves**, and both are correctly protected rather than repairable. The QA-plan
  YAML command was unsatisfiable while an accepted retained dependency requires
  the module, so any local fix would have needed dependency authority; the
  report-lint conflict needed gate-definition and workflow authority. Neither
  belonged to a Developer session. This escalation resolves **against the Tech
  Lead**, like ESC-20-008 and unlike ESC-20-009: planning validated only that
  the search pattern currently matched, never that its finished-work state was
  reachable, and the report-lint rule demanded evidence that `AGENTS.md` and §8
  forbid producing. VAL-032 keeps its assertion with a stronger proof, the
  control is amended rather than the evidence, Task 16 resumes as a same-task
  continuation after `task-16a`, and its failed attempts stay immutable.

---

## TRG-20-013 - DEC-143's replacement YAML build-graph literal is unsatisfiable

- **Trigger:** Task 16 cannot satisfy the protected DEC-143 exact output
  literal for `go mod why -vendor -m gopkg.in/yaml.v3`; the attached independent
  result exits 0 but ends `(main module does not need to vendor module
  gopkg.in/yaml.v3)`, while DEC-143, the Task 16 spec, and the QA plan require
  `(main module does not need module gopkg.in/yaml.v3)`. The amended
  report-lint also exits 6 because its expected-nonzero declaration parser
  treats valid instructional and wrapped explanation lines in Task 16's section
  as grammar declarations.
- **Trigger signature:**
  `task-16:specification-contradiction:dec-143-go-mod-why-vendor-literal-and-parser-regression`
- **Source:** `project-docs/docs/phases/phase-20/tasks/task-16-done.md`;
  attached commit-verification and independent command evidence.
- **Candidate:** typed `BLOCKED` / report-lint failure.
- **Attempt count:** 1.
- **Checkpoint:** Product commit
  `afb262ddc3636f6a4c9a05d6b9e77b4aeac276eb` is verified by the supplied
  `git rev-parse --verify` result (`exit_code: 0`) and subject match
  `refactor(cli): remove universal formatter framework`; documentation commits
  `e09b98c7dc9900d50372da8a5234858ef00d5566` and
  `1f93513000586dfd339e4ccd3ac03f57f8fc36fe` are likewise supplied as verified.
  The independent `go mod why -vendor` result terminates `exit_code: 0` with
  the contradictory vendor-mode literal, and independent report-lint terminates
  `exit_code: 6` with `OUTCOME BLOCKED` and six declaration-grammar failures.
- **Verdict:** BLOCK.
- **Reasoning:** DEC-143 correctly established that YAML absence is a
  requirement/build-graph proof and forbids any dependency change, but its
  replacement's required exact literal omits the Go tool's vendor-mode wording.
  The attached successful command proves that the prescribed literal cannot be
  truthfully recorded. Altering the exact expected output, command, or
  quality-gate acceptance requires protected specification and quality-gate
  authority; TL-Lite cannot substitute a literal or direct a dependency change.
  The report-lint parser failure is a separate local workflow-control defect,
  but the trigger expressly allows that repair only after the protected mismatch
  has been resolved by recorded authority. It therefore cannot independently
  make this event REPAIR or permit advancement.
- **Local decisions made:** none.
- **Resolution evidence:** The supplied verified commits establish committed
  deletion work and documentation, not resolution of the new exact-literal
  contradiction. DEC-143's authority does not cover its own divergent required
  literal, and the supplied report-lint result confirms the parser is still
  failing. RESUME is unavailable.
- **If BLOCK:** category specification contradiction / quality-gate-definition
  authority; ESC-20-013 raised in `tl_escalations_current.md`; resolution
  condition: full Tech Lead records the lawful accepted vendor-mode output (or
  another authority-approved exact proof) without changing the command, VAL-032
  assertion strength, or dependency policy, then directs bounded Task 16a parser
  repair and the Task 16 same-task continuation with fresh verification and any
  affected independent re-review.
- **Human attention required:** yes.
- **References:** `tasks/task-16-done.md` (Workflow outcome; VAL-032; Exact
  quality-gate results; Deviations); `tasks/task-16.md` (Exact Commands To Run;
  Expected Non-Zero Commands; Notes For Developer);
  `decisions.md` DEC-143 decisions 2-7; `project-docs/specs/09_AI_WORKFLOW.md`
  §8-9; `AGENTS.md` (Commands; Hard Stops).
- **Tech Lead adjudication (DEC-144, 2026-08-16):** verdict **sustained on both
  halves**, and the sequencing judgment — protected literal first, bounded parser
  repair second — was right. Both defects are Tech Lead drafting, like ESC-20-008
  and ESC-20-012: DEC-143 transcribed `go mod why`'s non-vendor message for a
  command it required to run with `-vendor`, and the Task 16a fixture used a
  section shape no real spec uses. The observed vendor-mode line at `exit_code: 0`
  is a lawful pass under the corrected condition; the command, `-vendor`, VAL-032,
  and the dependency policy are unchanged. `task-16b` repairs the parser
  red-first, then Task 16 continues in a fresh session and re-records every
  command, including the one whose `BLOCKED` row this triage correctly protected.
