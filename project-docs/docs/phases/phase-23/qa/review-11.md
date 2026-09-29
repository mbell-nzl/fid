# Phase 23 - Mid-Phase Review 11 (Tasks 13-15)

Reviewer: qa-reviewer
Date: 2026-09-30
Tasks covered: phase-23/task-13, phase-23/task-14, phase-23/task-15
Review type: INITIAL
Prior review: none (review-10 passed the preceding window)

## Verdict

REPAIR. The final tree passes the required task-15 commands and the focused behavior checks. Task 14 still lacks the accepted every-leaf Human fixture and JSON fact comparison; task 15's completion report lacks per-row retained-test provenance. Both repairs have existing, decision-free task ownership. The window does not advance.

## Scope Reviewed

- Task specs and completion reports: `tasks/task-{13,14,15}.md` and `tasks/task-{13,14,15}-done.md`; phase QA plan, validation contract, DEC-262/263/264, `schema.md` sections 5 and 6.2, `public-interface.md` sections 5 and 7, `contracts.md` section 14, the quality plan, and review-10.
- Final committed source at `cdc2f1472b5b375d12af4d7095396ed12aa9dc25`, including the task commits `bfee063f4cd1abb7a8d884b66802bb0136228db9`, `cdcc4f01fb2060d4fa5721a4108b9fda1dde328a`, and `0dd4c57849f8be7779e1c2e43cf7b14093c945f0`. Each reported hash resolved with `git rev-parse --verify <hash>^{commit}` (terminated `exit_code: 0`); subjects and `Refs: phase-23/task-N` matched. The Runner ledger was the only dirty path before this review.
- Inspected the denial adapters and closed-code tests, real-process activity and Human journeys, 0018 SQL and rollback tests, changed retained tests, and the fixture guide. No product or governing file was edited by QA.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 2
- Improvement findings: 1
- Blocker candidates: 0
- Existing escalation references: ESC-23-007, ESC-23-008, ESC-23-009 (all resolved)
- Covered validation assertions: VAL-004, VAL-055, VAL-056, VAL-057, VAL-058, VAL-059, VAL-060, VAL-061, VAL-062, VAL-063, VAL-070, VAL-071, VAL-072, VAL-073, VAL-079

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence and limit |
|---|---|---|---|
| VAL-004 | 15 | ON TRACK | Fresh 0018 migration tests passed; schema 18 and rebuilt-projection FK/integrity assertions are in `migration0018_test.go` and the Concern journey. Three phase-level fresh-workspace pragma receipts remain with QA Lead. |
| VAL-055 | 13, 15 | ON TRACK | Uncached built-binary activity inventory passed. `p23_activity_inventory_test.go:12-85` executes all ten producing leaves and compares each Event's subject/count; lines 86-109 compare the sorted distinct names for those Events with exactly the accepted ten. Requirement and Decision subjects and requirement-subject Concern/Deviation are also probed in their separate real-process journeys. Production Review, Finding, Outcome, Concern, Deviation, Requirement, and Decision writers use only those ten Phase 23 names. |
| VAL-056 | 13 | ON TRACK | `p23_failure_contract_test.go:18-77` preserves closed envelope, code, category, subject, and detail keys. Both uncached precedence journeys returned their Stop code over an unsatisfied guarded Obligation, and the exact-schema CLI test passed. |
| VAL-057 | 13 | ON TRACK | `p23_task13_fixture_replay_test.go:15-73` compares both protected-denial responses with durable open Stop/Attention rows, null Event, false replay, unchanged journal and idempotency counts. The built-binary replay/conflict and noncommitting exhaustion paths ran in the full tagged gate. |
| VAL-058, VAL-060 | 14 | NOT ON TRACK | The ordinary renderer test covers 18 operation names with synthetic expected text, and selected real-process Human inspections pass. EVID-26 and `public-interface.md` section 7 also require verbatim tagged real-process fixtures for **every** new leaf and a comparison with each JSON field set. The tagged generator and guide do not supply that inventory or comparison (MF-1). No proof/prose violation was found in the fixtures inspected. |
| VAL-059 | 14 | ON TRACK | Ordinary resume test and uncached mixed-state/overflow tagged journeys pass; the 101-Concern case returns `result-limit-exceeded`. |
| VAL-061 | 15 | ON TRACK | Latest-schema introspection covers removed tables, columns, rules, indexes and both implicit kind triggers; binary tests reject the removed leaves without a new Event. The view-column list is measured identical across 0018. Three production absence searches independently returned their declared exit 1. |
| VAL-062 | 15 | ON TRACK | Seeded 17-to-18 test keeps Issue rows and deletes non-Issue `links` and `activity_log` rows; valid and invalid probes cover the rebuilt reference triggers and Issue-only open-Concern count. Uncached migration tests passed. |
| VAL-063 | 15 | ON TRACK | The populated binary closure journey checks maintenance before and after rebuild; maintenance tests inventory retained check names and accept valid generic references. Both package and tagged journeys passed. |
| VAL-070, VAL-071, VAL-072, VAL-073 | 14; Human at phase end | DEFERRED | The three generated Human workspaces and their invocation steps pass the fixture generator. Human comprehension and view judgment remain assigned to phase end; MF-1 concerns the separate deterministic every-leaf text contract. Task 18 is assigned final-binary regeneration. |
| VAL-079 | 13-15; final at phase end | ON TRACK | Reports record tagged times 265.793/229.795/240.222 s (limit 900), `cmd/fid` 204.920/196.303/198.174 s (limit 600), and measured `internal/workflow` 62.384/117.669/101.597 s (limit 600). The independent final-tree tagged gate took 238.406 s. The independent 100k benchmark reported 22582 ms/rebuild, below 30000; none reaches the 80% headroom warning. Final uncached phase measurement remains with QA Lead. |

## Gate And Evidence Checks

Working directory for every independent command: `/home/mark/projects/fid`. I compared both `## Tests run` and `## Exact quality-gate results` against each task's exact block and the QA plan's EVID-25 through EVID-28. Task 13 and 14 each record all five byte-identical commands with final PASS tokens and observed `exit_code: 0`; task 15 records all nine, including the three declared expected-exit-1 searches. Earlier red attempts are retained. All three two-argument report-lint runs ended `PASS report-lint completed with no failures` (each terminated `exit_code: 0`); this mechanical check does not assess MF-1 or MF-2.

| Check | Command or evidence | Result | Observed token and termination |
|---|---|---|---|
| Task-15 exact 1 | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`). |
| Task-15 exact 2 | `golangci-lint run ./...` | PASS | `0 issues.` (terminated `exit_code: 0`). |
| Task-15 exact 3 | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` | PASS | `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t(cached)\tcoverage: 62.1% of statements` (terminated `exit_code: 0`). Cached package results do not supply fresh timing. |
| Task-15 exact 4 | `go tool cover -func=coverage.out` | PASS | `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t60.5%` (terminated `exit_code: 0`). |
| Task-15 exact 5 | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t238.406s` (terminated `exit_code: 0`). The silent running session was polled until exit. |
| Task-15 exact 6 | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` | PASS | `ok  \tgithub.com/mbell-nzl/fid/internal/journal\t25.255s` (terminated `exit_code: 0`); benchmark line reports `22582 ms/rebuild`. |
| Task-15 exact 7 | `rg -n --glob '!*_test.go' "task_completions|ac_verifications" cmd internal` | PASS | `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |
| Task-15 exact 8 | `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` | PASS | `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |
| Task-15 exact 9 | `rg -n --glob '!*_test.go' "has_completion_report|concerns_addressed" cmd internal` | PASS | `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |
| Uncached in-window named tests | `go test -tags integration -run '^(TestP23TenNameActivityInventoryFromProducingLeaves|TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23FailurePrecedenceWithProtectedDecisionAndObligation|TestP23FailurePrecedenceWithReadinessStopAndObligation|TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess|TestP23ConcernListHumanProjection|TestP23ResumeRejectsJudgmentSectionOverflow|TestV01ClosureJourneyAndRebuildParity|TestProcessExitContract|TestPhase23LeavesRejectUnknownRequestFields)$' -count=1 ./cmd/fid/integration/...` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t18.589s` (terminated `exit_code: 0`). |
| Uncached migration and maintenance | `go test -run '^(TestMigration0018.*|TestLatestSchemaHasNoV1JudgmentTables|TestCheckErrorsReportsAllSeededStructuralViolationFamilies|TestCheckAcceptsPopulatedGenericReferencesWithoutStructuralErrors)$' -count=1 ./internal/db ./internal/maintenance` | PASS | `ok  \tgithub.com/mbell-nzl/fid/internal/maintenance\t0.149s` (terminated `exit_code: 0`); DB package line was `0.651s`. |
| Uncached Human tests | `go test -run '^(TestPhase23HumanLeafTextMatchesRecordedFacts|TestResumeReportsRealFindingsAndJudgmentSections|TestConcernListGenericReadExcludesV1Rows|TestConcernListRetainsIssueScopedV1ReadUntilRemoval)$' -count=1 ./cmd/fid`; `go test -tags integration -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.537s` and `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t6.409s` (each terminated `exit_code: 0`). |
| CLI exact-detail schema | `go test -run '^TestReasonCodeDetailSchemaIsExactForEveryCode$' -count=1 ./internal/cli` | PASS | `ok  \tgithub.com/mbell-nzl/fid/internal/cli\t0.026s` (terminated `exit_code: 0`). |
| DEC-264 unchanged-file proof | `git diff --name-only fa6355a HEAD -- migrations/0001* ... migrations/0017* migrations/0018_p23_remove_v1_judgment.up.sql` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`); every 0001-0017 path and 0018 up was in the command's path set. |

## Per-Task Findings

### Task 13 - Phase 23 CLI surface, failure codes, and idempotency

#### Must-fix

- None.

#### Improvement

- None.

Observation: `cmd/fid/decision.go:22` and `cmd/fid/human_outcome.go:24` add sorted durable IDs, `event:nil`, and `replayed:false` to the decoded StopDetails, then use the retained `cli.NewFailure` path. `internal/cli/codes.go:257-264` closes both detail shapes. The real-process helper compares IDs against open Stop/Attention rows; the exact-schema and precedence tests remain green. `internal/workflow/stop.go` and the journal StopDetails value catalog have no diff since DEC-263's base, so Stop projections were not widened. The retained `fid-cli.v1` envelope, categories, exits, subjects, and precedence remain asserted. The unresolved question about a protected Decision becoming unsatisfied again after its Stop resolved stays a Tech Lead wrap-up note under DEC-263, outside VAL-057.

### Task 14 - Human review and judgment projections

#### Must-fix

- **MF-1** [`cmd/fid/integration/p23_human_fixtures_test.go:94`]: the tagged generator checks only selected reads from three Human workspaces and captures its own expected Human output at lines 130-142. `qa/human-fixtures.md` records those selected reads, without a verbatim Human/JSON fact pair for every one of the 18 new leaf operations enumerated by `cmd/fid/p23_human_test.go:45-66`. The generator checks JSON invocations for nonempty output but does not compare their field sets with Human facts. Thus EVID-26 and `public-interface.md` section 7 are unproven despite passing tests. Owning task: task-14. Affected assertions: VAL-058, VAL-060. Repair boundary: task-14's named ordinary/tagged fixture tests and `qa/human-fixtures.md`, with report evidence updated; retain the accepted Human/JSON contract. Required evidence: per-leaf verbatim expected Human output in both suites, matching JSON fact/field comparison, no-proof/no-prose checks, and the unchanged five exact gates. Re-review scope: a fresh independent reviewer inspects all leaf pairs against `contracts.md` section 14 and reruns representative tagged and ordinary cases.

#### Improvement

- **I-1** [`internal/workflow/resume.go:1`]: the file is 513 lines, above the quality plan's 400-line non-test target. Fix direction: split the six bounded judgment projection helpers on the next authorized causal edit; keep this as Tech Lead wrap-up debt.

Observation: `TestP23ConcernListHumanProjection` independently passed through the real CLI for JSON and Human output, filters, limit refusal, and mixed-route rejection. The generated VAL-070..072 workspaces remain usable for the later Human judgment; their selected command output cannot close VAL-058's every-leaf requirement.

### Task 15 - Migration 0018 and V1 judgment removal

#### Must-fix

- **MF-2** [`tasks/task-15-done.md:43`]: the VAL-061 table names replacement tests, and line 71 gives only grouped introducing hashes. Converged retained rows such as AC add/list (`internal/issues/ac_test.go`), Issue detail/list, maintenance, and the three removed-leaf process-exit cases lack individual file/line, introducing commit **and task**, the particular accepted invalidating change, and the executing command together. `report-lint` passes because it does not enforce this evidence rule. Owning task: task-15. Affected assertion: no VAL (completion-report provenance rule; VAL-061 evidence remains independently green). Repair boundary: `tasks/task-15-done.md` only, with one traceable entry per converged retained row and no changed tests, SQL, commands, or contract. Required evidence: provenance entries checked against `git blame`/introducing commits and the existing executing gate or named test; a committed report and passing two-argument report-lint. Re-review scope: a fresh independent reviewer compares every converged retained row in the task-15 diff to the new report entries and checks representative hashes and commands.

#### Improvement

- None.

Observation: `migrations/0018_p23_remove_v1_judgment.down.sql` follows `schema.md` section 5's nine ordered groups: teardown/rename, four unchanged V1 tables, post-0016 V1 Concern, two 0004 kind guards, old Concern index, row-copying `agent_runs` rebuild, 0011 view, four 0005 link guards, then the 0016 activity guard. `TestMigration0018DownMatchesFreshSchema17` compares each noninternal `sqlite_master` tuple (`type`, `name`, `tbl_name`, `sql`) to fresh 17 while normalizing only the accepted quoted `judgment_concerns` and `agent_runs` names. The seeded Run survives with its role/model and null restored completion ID; `MigrateDown` succeeds through the full chain; 17→18→17→18 ends with clean FK and integrity pragmas. The uncached migration package passed. `git diff` found no change to 0001-0017 or 0018 up from the pre-DEC-264 task-15 checkpoint. The activity inventory and production writer search agree with `schema.md` section 6.2: Requirement and Decision rows retain their own subjects, judgment subjects resolve as declared, and no additional Phase 23 name is produced.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-1 | task-14 | VAL-058, VAL-060 | Named Human fixture tests, generator, `qa/human-fixtures.md`, and task-14 report; implement the accepted every-leaf evidence without changing contract or gates. | Verbatim ordinary and tagged Human expectations for every new leaf, JSON fact/field comparisons, no-proof/no-prose checks, and all five original commands with observed termination. | Inspect all leaf pairs and rerun representative real-process and ordinary cases, then compare report and spec commands. |
| MF-2 | task-15 | no VAL | Task-15 completion report only; preserve product, SQL, tests, and exact commands. | Per-converged-row file/line, introducing commit and task, invalidating change, executing command; commit and passing report-lint. | Diff every retained converged test row against report provenance and independently verify representative introducing commits. |

## Blocker Candidates

- None. Both Must-fix findings have bounded owners in this window and require no protected definition change.

## Existing Escalation References

- ESC-23-007 (DEC-262), ESC-23-008 (DEC-263), and ESC-23-009 (DEC-264) are resolved. This review allocates no escalation identity.

## Final Lines

Review written: project-docs/docs/phases/phase-23/qa/review-11.md
Must-fix: 2  Improvement: 1  Blocker candidates: 0
VERDICT: REPAIR
