# Phase 23 wrap-up audit

Date: 2026-10-03. Tech Lead, acting on the Human's instruction to complete the
Phase 23 wrap. Decision: DEC-271. Outcome: complete and ready for PR/merge,
with no active phase. [Final acceptance](qa/final-acceptance.md) reconciles
every validation outcome.

## Audit basis and evidence limits

Read: final reviews 0/1/2, both Human checklists, review-15/15a, the
89-row validation contract, all 21 completion reports' deviation sections,
phase decisions DEC-222..270, the triage log, the escalation archive, the run
ledger to RUN-END, the canonical specs, and the active carry-forward register.
The phase diff runs from merge base `0195897a` to `704eb03` (289 commits). The
last product change is `f5fdec0` (task 16 MF-3); final review 2 qualified
`6c94277`, and later commits change only documents and agent configuration.
Product gate results belong to their executing reports and reviews. This wrap
changes documentation only.

## Decision audit

| Decisions | Final disposition |
|---|---|
| DEC-222 | Heavy contract set and 18-task pack delivered; the six section 15 canonical amendments are verified present (03 §3.2, §5, §8.4; 06 §3, §5; 02 §4; GLOSSARY). |
| DEC-223..227 | Link-preserving 0017 rebuild, named Stop raisers, mutating-side independence, failure-code precedence and non-committing denials, and wall-time budgets all delivered. VAL-079 budgets pass with headroom. DEC-227's open project-default question is settled in DEC-271: no change. |
| DEC-228..237, 240 | Planning-pass dispositions and the Condition encoding (DEC-235 Option A, DEC-236 relation call, VAL-080..084) delivered. Final review 2 retains VAL-080..084 PASS. |
| DEC-238, 239 | ESC-23-002/003 recoveries completed. DEC-239's Run-holder fixture boundary recorded in 02 §5.2 and 08 §8. |
| DEC-241..248 | Task 5/6/9 scope and ownership corrections delivered; DEC-246 coverage baseline reconstruction is retained as a process deviation, not re-attested. |
| DEC-249..257 | ESC-23-004 dispositions delivered by tasks 2, 6, 10, 11. GLOSSARY "applicable" (DEC-249) and the 09 BLOCK re-review route (DEC-257) are synced. |
| DEC-258, 259, 261 | Parallel tests, `paralleltest`, the tiered loop, and the modernc `checkptr` exclusion are in AGENTS, 08 and 09. The phase-end full-checkptr control passed in final review 2. |
| DEC-260 | Requirement writes need a trusted Actor and no capability; the closed set stays at thirty. |
| DEC-262 | `activity_log` rows split at 0018; VAL-055 PASS. |
| DEC-263, 265, 266 | Denial fields, 19 Failure codes, and the Phase 21 `authority-*` details. 06 records them in its new Phase 23 amendment section. |
| DEC-264 | 0018 down restores schema 17. 03 §8.4 now incorporates `schema.md` sections 4-5. |
| DEC-267 | VAL-005 two-column correction and the task-18 MF-1 harness repair; closed by final review 2. |
| DEC-268, 269 | Checklist 1 FAIL; A1-A3 repaired by tasks 19-21; checklist 2 is 7/7 PASS. 06, 07 and GLOSSARY synced. The A4-A8 agenda becomes CF-045 (DEC-271). |
| DEC-270 | Pre-change Phase 23 journals fail closed. 03 §8.4 states it; 08 has no conflicting pre-adoption promise (verified). DEC-270 was flagged to the Human for possible reversal; no reversal was recorded. |

## Task deviations

Every report is historical evidence; failed attempts stay on record. All 21
reports pass report-lint (final review 2, `audits.report-1..21`).

| Task(s) | Deviations and disposition |
|---|---|
| 1 | DEC-238 cross-role recovery; failed benchmark and runtime-budget attempts kept. Accepted. |
| 2 | Additive fixture option in retained `p22_fake_test.go` (DEC-249 repair), flagged for the Tech Lead. Accepted: it changes no behavior. |
| 4, 5, 8, 13, 19, 21 | Fixture ordering, compile, lint-string and report-format first failures, each repaired and rerun byte-for-byte. No contract change. |
| 6 | Missed before-edit coverage baseline reconstructed under DEC-246 (TDD-order/process deviation kept). DEC-244/245 scope additions. |
| 9 | Overlapping root invocation caused one `cmd/fid` failure; the serial rerun passed. Cause not established; noted for workflow review. |
| 10 | Required gates run with `-count=1` appended. Kept as a historical deviation; task 11's continuation ran the exact commands. |
| 11 | DEC-258 edited files outside task-11 scope (parallel annotations only); DEC-253 ratified `stop_conditions.go`. |
| 12, 14, 15, 17 | Retained-test convergence after accepted removals, with introducing-commit provenance recorded. TRG-23-009 extended task 14 by the minimum scope. |
| 16 | Final-review-1 MF-3: 60-second fixture Claim lease expired under a long gate. Repaired to one hour with a persisted-lease regression; no expiry rule changed. |
| 18 | Seven report attempts across DEC-266 scope gaps and two model-capacity interruptions (ESC-23-011). |
| 20 | DEC-270 replaced the impossible pre-change replay criterion. |

Every immediate QA recommendation is closed: final review 2 has no Must-fix.
I-1 (wrap synchronization and the A4-A8 disposition) is done by this wrap.
I-2 (`internal/workflow/resume.go` size) stays with its next authorized edit
under DEC-267; it is a size target, not a blocker.

## Canonical specification sync

| Spec | Verdict | Rationale |
|---|---|---|
| 00_PROJECT_CHARTER.md | verified-no-change | Review and Human judgment are within the chartered local-first operational backend; no goal or boundary moved. |
| 01_PRD.md | touched: DEC-242 | PRD-REV-004 addressed/resolved split. PRD-REV-006 (exhaustion returns the code plus Attention) still holds after DEC-269. PRD-ISSUE-011 already asks for stable requirement identity, which DEC-269 delivers. |
| 02_TRD.md | touched: DEC-238, DEC-239, DEC-242, DEC-271 | Earlier task syncs plus this wrap's package boundary paragraph (`internal/review` added, `internal/requirements` removed, `internal/links` restricted). |
| 03_TRD_SCHEMA.md | touched: DEC-222, DEC-242, DEC-264, DEC-270 | Section 8.4 now covers 0018 and its down migration, and fail-closed pre-change journals. |
| 04_TRD_CONFIG.md | verified-no-change | No configuration key, file, or consuming path changed; Human authority comes from the existing Activation profile (DEC-251). |
| 05_DESIGN_PLAN.md | verified-no-change | Section 4 hierarchy and section 6 ("state what failed, the subject, and the next legitimate action") still hold. Tasks 19/21 move toward section 6; A5's remaining gap is CF-045. |
| 06_DESIGN_PLAN_CLI.md | touched: DEC-242, DEC-263, DEC-265, DEC-266, DEC-269, DEC-270 | New "Phase 23 accepted amendment" section and requirement-identity paragraph. |
| 07_DESIGN_PLAN_FLOWS.md | touched: DEC-239, DEC-242, DEC-269 | Section 7 gains the rework-extension resolution step. |
| 08_QUALITY_PLAN.md | touched: DEC-239, DEC-259 | No pre-adoption journal-compatibility promise found (DEC-270 check). DEC-227 bound not made a project default. |
| 09_AI_WORKFLOW.md | touched: DEC-257, DEC-258, DEC-259 | No wrap-time workflow change; workflow findings go to the retrospective. |
| 10_IMPLEMENTATION_ROADMAP.md | touched: DEC-271 | Phase 23 complete; next work is the CF-045 scheduling decision, then Phase 24 planning. |
| GLOSSARY.md | touched: DEC-226, DEC-242, DEC-249, DEC-269 | Adds Requirement Reference (stable requirement ID versus reference ID) and widens Human Outcome to Stop conditions. |
| AGENTS.md | touched: DEC-258, DEC-259, DEC-271 | Adds `internal/review/` to the package list. The rest of the stale inventory stays with CF-033. Commands unchanged for DEC-227. |

## TRG and ESC audit

Seventeen unique TRGs: 13 BLOCK, 2 REPAIR, 2 RESUME. Every BLOCK names an ESC,
and all twelve ESCs are resolved and archived in [tl-escalations.md](tl-escalations.md).
TRG-23-011 has an appended disposition section, not a second identity.

| Trigger | Hindsight audit |
|---|---|
| 001, 002 | BLOCK correct (protected budget/ownership); 002 is a corrective schema record for 001, same ESC-23-002. |
| 003 | BLOCK correct: real contradiction between VAL-015 evidence and accepted lifecycle (DEC-239). |
| 004 | BLOCK correct: three protected blocker candidates (DEC-249..257). |
| 005 | BLOCK correct: budget breach with no task-scope cause; Human-delegated DEC-258. |
| 006, 007, 008, 010, 011 | BLOCK correct, but each exposed a **planning** gap (authority row, journal registry, closed failure details, rollback chain, code count). Correct escalation of Tech Lead defects. |
| 009 | REPAIR correct: minimum same-task convergence. |
| 012 | REPAIR correct: one bounded retry after model capacity. |
| 013 | BLOCK correct after the bounded retry also failed (ESC-23-011). |
| 014 | RESUME correct in hindsight: capacity was restored and the retry ran. Its identity defect was corrected by 015. |
| 015 | RESUME correct; governing record for the retry. |
| 016 | BLOCK correct: protected spec contradiction plus pending Human judgment (ESC-23-012). |
| 017 | BLOCK correct, cited the already-open ESC-23-012 rather than minting one. Acceptable because the blocker was inside that escalation's repair scope; noted for the workflow retrospective. The blocker was a Tech Lead-authored criterion (DEC-270). |

No RESUME should have escalated. No verdict lacks a durable record. DEC-002
autonomous resume stays permanent.

## Carry-forward trigger audit

| ID | Wrap disposition |
|---|---|
| CF-008, CF-026a, CF-026b | Retain; Phase 23 adds no Issue lifecycle or inspection projection decision. |
| CF-018 | Retain; `internal/sqlquery/sqlquery.go` is unchanged in the phase diff. |
| CF-023, 029, 030, 036, 039 | Retain for dedicated workflow maintenance. Phase 23's report-lint churn (seven task-18 report attempts, repeated row-grammar failures) and 13 BLOCKs reinforce CF-030/036/039. |
| CF-033 | Partial: `internal/review/` added to AGENTS; the nonexistent-package list and quality locator remain. |
| CF-037 | Retain for Phase 25. |
| CF-038 | Retain. The only edit to `live_output_test.go` was DEC-258's parallel annotation (`3dff457`, 3 lines), not a causal edit of coalescing. |
| CF-040 | Phase 23 portion delivered (VAL-073 PASS). Pre-existing view polish stays with Phase 25. A4 (raw Stop JSON) moves to CF-045, since it degraded a Phase 23 judgment. |
| CF-041 | Retain. The Makefile changed only for DEC-259 flags; hosted CI still omits the tagged suite. |
| CF-042 | Partial delivered by task 1 (children reuse the parent binary; 24-hour prune). Cleanup-on-kill stays in maintenance. |
| CF-043 | Retain for a future TL; no trigger. |
| CF-044 | Retain; `mustJournalPath` is unchanged. |
| CF-045 | NEW: Human operating experience (A4-A8 and residual PARTIALs). Timing is decided in the pre-Phase-24 discussion. |

## Wrap verification

Administrative checks only; product gates belong to final review 2.

- `project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-23/run-ledger.md`
  — PASS — `VALID state=ADVANCEABLE` (terminated `exit_code: 0`)
- `git diff --check` — PASS — no output (terminated `exit_code: 0`)
- Local link scan of the wrap and current documents — PASS — no broken links.
