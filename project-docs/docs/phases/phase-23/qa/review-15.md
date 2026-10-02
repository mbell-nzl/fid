# Phase 23 - Mid-Phase Review 15 (tasks 19-21)

Reviewer: qa-reviewer
Date: 2026-10-02
Tasks covered: phase-23/task-19, phase-23/task-20, phase-23/task-21
Review type: INITIAL
Prior review: none (review-14 covered tasks 16-18)

## Summary

- Review outcome: REPAIR
- Must-fix findings: 2
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-23-012 (existing phase-end work; no new blocker candidate)
- Covered validation assertions: VAL-085, VAL-086, VAL-087, VAL-088, VAL-089

The implementation evidence supports all five assertions. The window cannot advance because two completion reports omit byte-identical passing runs after recorded failing commands. These are report and verification repairs within their owning tasks. VAL-070..073 remain Human-owned judgments for the later regenerated-checklist review.

## Scope Reviewed

- Task specs and committed completion reports for tasks 19, 20, and 21; the phase QA plan, validation contract, DEC-268..270, contracts sections 6, 11, 12, and 14, and public-interface sections 1 and 2.
- Product commits `180f5168` (task 19), `0ba005d7` (task 20), and `c6e070b` (task 21), plus the bounded task-21 fake commit `343b458`. Their reported full hashes resolve with `git rev-parse --verify <hash>^{commit}` (each terminated `exit_code: 0`), and `git log` confirms their reported subjects and `Refs: phase-23/task-N` trailers. The final report commits are `d7f960d`, `bccec6f`, and `cbe2cc0`, respectively; all three reports are committed and clean.
- Current phase HEAD before this review: `cbe2cc0ef337d939a59897db61d552521d69ef32`. Only Runner-owned `run-ledger.md` and `triage-log.md` were dirty; this review did not edit them.

## Validation Contract Coverage

| Assertion | Window assessment | Evidence and limits |
|---|---|---|
| VAL-085 | Behavior ON TRACK; report trace incomplete (MF-1) | `cmd/fid/integration/p23_human_failure_test.go:41,168` runs 15 real-binary Human/JSON pairs, inventories the seven renderer codes, checks exact stderr, empty stdout, exit parity, repeatable JSON bytes and unchanged Event counts. The focused race selection passed; source diff leaves the JSON writer and exit mapping intact. |
| VAL-086 | ON TRACK | `p23_review_evidence_budget_test.go:77,204` exercises the exhausted Stop, wrong class, rejected result, other Issue, non-human authority denial, approved same-Issue Outcome, atomic resolution, admitted Claim/Run, and separate manual resolution. `budget_test.go:146` pins the Condition and its nullable fields. The focused race selection passed. |
| VAL-087 | ON TRACK | `internal/evidence/execution_test.go:86` and `internal/workflow/budget_test.go:168,181` exercise the isolated active-Run admission and resume call sites, excluding both unknown and exhausted rework budgets while retaining attempts/tokens controls. This is the task spec's allowed isolated-fixture route for a CLI state the Claim/Run rules cannot construct. The shared fake also tests both excluded states and both generic controls. |
| VAL-088 | Behavior ON TRACK; report trace incomplete (MF-2) | The real-process journey checks later consumed totals 2 and 3, remaining 0, and no new Stop. `p23_human_failure_test.go:14` asserts the exact two Human lines; the matrix checks the JSON twin. The focused race selection passed. |
| VAL-089 | ON TRACK | `p23_requirement_reference_test.go` checks stable/reference IDs, supersession, exact old-reference lookup and stable links, all/active/Issue-filtered lists. `p23_requirement_roundtrip_test.go:11` checks new-binary rebuild, unchanged digest, five byte-identical JSON reads, and idempotent retry. The independent preserved-journal replay below confirms DEC-270(b). Human fixture tests cover distinct headings, reference lines, sources and Issue links. |

VAL-070..073 are context only: their fixture/step commands passed, while comprehension and presentation outcomes stay with the Human. This review does not convert deterministic fixture readiness into a Human PASS.

## Gate And Evidence Checks

All three specs require the same five commands from the repository root. No QA-plan command is added specifically to tasks 19-21. QA-plan section 22 places the full phase gate, full-checkptr control, and VAL-079 fresh timing run at phase end; this bounded review reran focused behavior and two representative required commands. The table records each report's claimed token and observed termination; those full-gate results are report evidence mechanically checked by independent report-lint, not claimed as fresh QA executions.

| Exact required command | Task 19 report | Task 20 report | Task 21 final attempt report |
|---|---|---|---|
| `gofmt -l cmd internal` | `exit 0 (no output)`; `exit_code: 0` | same; `exit_code: 0` | same; `exit_code: 0` |
| `golangci-lint run ./...` | `0 issues.`; `exit_code: 0` | `0 issues.`; `exit_code: 0` | `0 issues.`; `exit_code: 0` (earlier `* staticcheck: 1`, `exit_code: 1`, retained; byte-identical retry) |
| `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` | `ok .../internal/workflow (cached) coverage: 62.1%`; `exit_code: 0` | `ok .../internal/workflow 117.580s coverage: 62.1%`; `exit_code: 0` | `ok .../internal/workflow (cached) coverage: 62.7%`; `exit_code: 0` |
| `go tool cover -func=coverage.out` | `60.5%`; `exit_code: 0` | `total: ... 60.5%`; `exit_code: 0` | `total: ... 60.7%`; `exit_code: 0` |
| `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` | `ok .../integration 215.530s`; `exit_code: 0` | `ok .../integration 253.239s`; `exit_code: 0` | `ok .../integration 291.867s`; `exit_code: 0` |

Task 19 records ordinary `cmd/fid` 215.387s, uncached `internal/workflow` 55.842s, and tagged integration 215.530s. Task 20 records 344.249s, 117.580s, and 253.239s. Task 21's final attempt records 215.093s, the prior uncached workflow 113.517s (final cached), and tagged integration 291.867s. The reported values are below the unchanged 600s ordinary and 900s tagged budgets. Task 20's attempt-1 gates were explicitly BLOCKED at the protected replay contradiction, then all five ran after DEC-270. No required gate was silently skipped or narrowed in the final reports.

Independent controls and reproduction, with the command's own token and observed termination:

| Check | Result |
|---|---|
| `project-docs/tools/report-lint.sh <task-N-done.md> <task-N.md>` for N=19,20,21 | Each printed `OUTCOME COMPLETE` and `PASS report-lint completed with no failures` (each terminated `exit_code: 0`). This checks committed report paths, subjects, required commands and row grammar; it does not supply missing byte-identical TDD retry rows. |
| `gofmt -l cmd internal` | `exit 0 (no output)` (terminated `exit_code: 0`). |
| `golangci-lint run ./...` | `0 issues.` (terminated `exit_code: 0`). |
| Focused tagged race selection for the Stop matrix, requirement identity/round trip/Human fixtures, rework journey/manual resolution, and both shared-fake tests plus dual-adapter corpus | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t21.608s` (terminated `exit_code: 0`). Command: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestP23RequirementReferenceAndLinkRealProcess|TestP23RequirementNewShapeJournalRoundTrip|TestP23HumanFixtures|TestP23ReviewCyclesRealProcessJourney|TestP23ReviewReworkManualResolutionWithPriorOutcome|TestP22FakeReviewReworkStopCondition|TestP22FakeGenericAdmissionExcludesReviewCycles|TestP22DualAdapterRecoveryConformance)$' ./cmd/fid/integration/...`. |
| Focused workflow and evidence race selections | `ok .../internal/workflow 2.375s` and `ok .../internal/evidence 1.022s` (each terminated `exit_code: 0`); they reach the actual two generic admission call sites and the preserved attempts control. |
| `go test -tags integration -run '^TestP23HumanFixtureStepsRemainRunnable$' -count=1 ./cmd/fid/integration` | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t4.791s` (terminated `exit_code: 0`). The committed guide marks its old snapshot historical and instructs regeneration; the generator supplies the activated Human and qualifying Outcome input. |
| DEC-270(b) handoff commands, current HEAD binary | `go build -o /tmp/fid-task20-candidate ./cmd/fid`: `exit 0 (no output)` (`exit_code: 0`). `/tmp/fid-task20-candidate --project /tmp/fid-task20-prechange.NjCRbu/judgment --json rebuild`: JSON code `journal-replay-stopped`, `.fid/journal/2026-10-01.jsonl` line 18 (`exit_code: 6`, expected). Exact `sha256sum` command before and after printed `a8fdbdccfc7fd39e27b80f256811d8b70113a0e30285122a166148455d1ba16e` (both `exit_code: 0`). The old stored result lacks `requirement_id`; refusal leaves source bytes unchanged. |
| Four previously red focused commands, rerun verbatim by QA | Task 19 parity, matrix and `gate/open-stop`: `ok .../integration 4.131s`, `6.285s`, `5.108s` (each `exit_code: 0`). Task 21 guidance parity: `ok .../integration 2.980s` (`exit_code: 0`). These independently show the current behavior; they do not retroactively add Developer-session PASS rows to completion reports. |

`git diff 3bf54ba..HEAD` shows no change to `current_event_json.go`, request/DDL/Event definitions, `go.mod`, `go.sum`, task gate commands or thresholds. The validation-contract change in that span is DEC-270's authorized VAL-089 fail-closed criterion. `current_event_value_catalog.go` adds only `RequirementRef.requirement_id`; `current_event_result_projection.go` stops deleting that field. Exact result validation remains in place. No legacy, tolerant or versioned journal path, environment-gated seed test, or default-skipped test was added. The accepted contracts and public-interface files did not change during implementation.

## Task 19 - Human-mode failures from the Stop-error family

### Must-fix

- **MF-1** [`project-docs/docs/phases/phase-23/tasks/task-19-done.md:34`]: The report preserves three expected-red commands (`TestP23StopFailureHumanAndJSONParity`, `TestP23StopFailureCodeMatrix`, and its `gate/open-stop` subtest), but records only a broader combined green command afterward. The requested byte-identical later result for each failed attempt is absent. Owning task: task-19. Affected assertion: VAL-085 (verification trace). Repair boundary: the task-19 completion report and fresh execution of those three exact commands; no product, test, spec or gate edit. Required evidence: append their verbatim passing output tokens and observed termination from a same-task continuation, retain the original failed rows, and pass report-lint on the committed report. Re-review scope: compare each exact command byte-for-byte, verify the new commit and report-lint, rerun a representative affected command. Current QA reruns pass but cannot be transcribed as Developer-session evidence.

### Improvement

- None.

Observation: `writeStopFailure` prints the JSON message to Human stderr and keeps stdout empty; the renderer's invalid-argument fallback uses that writer. The JSON branch and typed failure return remain unchanged. The real-binary matrix tests every calling leaf and the closed code inventory.

## Task 20 - Stable requirement identity on every reference

### Must-fix

- None.

### Improvement

- None.

Observation: the new stored result includes `requirement_id` immediately after `id`, exact old-reference lookup retains the old row with stable links, and the Human renderer uses `Requirement <stable ID>`, a separate `Reference ID` line, and an `Issue links` section. The committed new-shape test and independent old-journal refusal cover both DEC-270 states. The report preserves the attempt-1 blocked criterion and the accepted amendment; its red requirement and Human-fixture commands each have byte-identical later PASS rows. Retained assertion provenance names introducing commits and the DEC-269 displacement.

## Task 21 - Resolvable review-rework exhaustion

### Must-fix

- **MF-2** [`project-docs/docs/phases/phase-23/tasks/task-21-done.md:51`]: The guidance test's `go test -tags integration -run '^TestP23StopFailureHumanAndJSONParity$' -count=1 ./cmd/fid/integration` failed as recorded, but the report's later passing selections use different flags and test scopes. The byte-identical passing retry is missing. Owning task: task-21. Affected assertion: VAL-088 (verification trace). Repair boundary: the task-21 completion report and fresh execution of that exact command only; no product, fake, test, spec or gate edit. Required evidence: retain the first `FAIL` and termination, append the exact-command PASS token and observed termination from a same-task continuation, and pass report-lint on the committed report. Re-review scope: compare command bytes, verify commit/report-lint, rerun the guidance test and a representative affected verification.

### Improvement

- None.

Observation: the real-process tests cover all four non-qualifying Outcome states, atomic and manual resolution, later charges, and exact Human guidance. The single `GenericBudgetAdmissionBlocks` helper is used by both gate admission and resume reasons. The task-21 scope correction is bounded: `p22_semantic_fake_test.go` and `p22_semantic_denial_test.go` are the actual shared memory model and gate-denial mirror; commit `343b458` changes only those two fake paths plus the report. No product path changed in the continuation. The existing token-budget denial remains in the dual-adapter corpus, and the fake tests pin the changed Stop and both generic-budget controls.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh independent re-review scope |
|---|---|---|---|---|---|
| MF-1 | task-19 | VAL-085 | Run the three original focused commands byte-for-byte; edit only `task-19-done.md`. | Retained failed attempts; three new PASS tokens with observed exit codes, causal note, committed report and report-lint COMPLETE. | Compare all three pairs, verify report commit, rerun one affected command and report-lint. |
| MF-2 | task-21 | VAL-088 | Run the original guidance command byte-for-byte; edit only `task-21-done.md`. | Retained failed attempt; new PASS token with observed exit code, causal note, committed report and report-lint COMPLETE. | Compare the pair, verify report commit, rerun guidance verification and report-lint. |

Repair order: task 19, then task 21. The Phase Runner must keep this window open and dispatch a fresh independent QA Reviewer after both same-task report repairs. This review does not mark either finding resolved.

## Blocker Candidates

- None. Both findings have exact, decision-free report ownership. DEC-269 and DEC-270 already supply protected product and journal authority; the present review found no new contract, architecture, security, data-model, dependency or gate-definition question.

## Existing Escalation References

- ESC-23-012 remains open for the separate phase-end Human re-check and fresh final review. TRG-23-017's task-20 contradiction was disposed by DEC-270. Neither is a new blocker candidate for this window.

## Final Lines

Review written: project-docs/docs/phases/phase-23/qa/review-15.md
Must-fix: 2  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
OUTCOME REPAIR
