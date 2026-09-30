# Phase 23 - Mid-Phase Review 12 (Tasks 13-15)

Reviewer: qa-reviewer
Date: 2026-09-30
Tasks covered: phase-23/task-13, phase-23/task-14, phase-23/task-15
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-23/qa/review-11.md
Re-review of: review-11

## Verdict

PASS. Both review-11 Must-fix findings are closed by committed, independently checked evidence. Human acceptance of VAL-070 through VAL-073 remains assigned to phase end.

## Scope Reviewed

- Task 13-15 specs, completion reports, final source and test diffs, review-11, the phase QA plan and validation contract, DEC-262/263/264, `schema.md` sections 5 and 6.2, `contracts.md` section 14, `public-interface.md` section 5, and the quality plan.
- Repair commits `b71ed1ba18a277da22a71c5fdcf18eadb181be26` (task 14) and `a777ee3` (task 15 report); task implementation commit `0dd4c57849f8be7779e1c2e43cf7b14093c945f0`. The Runner-owned `run-ledger.md` was already dirty; this review did not alter it.
- Working directory for all independent commands: `/home/mark/projects/fid`.

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 1
- Blocker candidates: 0
- Existing escalation references: ESC-23-007, ESC-23-008, ESC-23-009 (resolved)
- Covered validation assertions: VAL-004, VAL-055, VAL-056, VAL-057, VAL-058, VAL-059, VAL-060, VAL-061, VAL-062, VAL-063, VAL-070, VAL-071, VAL-072, VAL-073, VAL-079

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-004, VAL-061, VAL-062, VAL-063 | 15 | ON TRACK | `internal/db/migration0018_test.go:248-349` compares the 18-down schema with a fresh 17, preserves a Run, completes full rollback, and checks 17→18→17→18 integrity. Latest-schema, seeded-link/activity, maintenance and removed-leaf tests are recorded in the report. Focused migration rerun passed. |
| VAL-055 | 13, 15 | ON TRACK | `p23_activity_inventory_test.go:12-109` calls the producing leaves, checks subjects and counts, then compares the distinct activity names to exactly: `requirement-referenced`, `requirement-linked`, `decision-referenced`, `review-recorded`, `finding-addressed`, `finding-resolved`, `human-outcome-recorded`, `concern-recorded`, `concern-disposed`, `deviation-recorded`. The focused real-process rerun passed. Requirement/Decision subjects are delivered with 0018 per DEC-262. |
| VAL-056, VAL-057 | 13 | ON TRACK | The DEC-263 closed detail sets for both `protected-decision-required` and `human-outcome-required` contain their respective StopDetails plus `stop_ids`, `attention_ids`, `event:null`, `replayed:false` (`public-interface.md:270,275`; `internal/cli/codes.go`). Real-process tests compare the IDs with durable open Stop/Attention rows, prove no Event or idempotency write on denial/retry, and exercise guarded-Obligation precedence. The protected Decision journey passed independently. |
| VAL-058, VAL-060 | 14 | ON TRACK | The guide appendix at `qa/human-fixtures.md:655` has 18 distinct leaf headings, each with verbatim Human stdout and a JSON result-fact block. All 18 guide Human blocks match the retained manifest byte-for-byte; every recorded JSON fact recursively matches the manifest's complete JSON result (the guide deliberately omits private prose). `p23_human_fixtures_test.go:80-145` reruns all pairs from the retained binary, compares complete JSON result maps, reconstructs exact Human text from JSON facts, and checks non-proof qualifiers and prose exclusion. Ordinary and tagged focused reruns passed. |
| VAL-059 | 14 | ON TRACK | The report and review-11's independent mixed-state and overflow journeys show bounded resume sections and `result-limit-exceeded` on overflow. No new implementation change followed review-11. |
| VAL-070, VAL-071, VAL-072, VAL-073 | 14; Human at phase end | DEFERRED | The three retained Human workspaces and the every-leaf appendix are runnable. Human comprehension and presentation judgment remain phase-end checks under the QA plan. |
| VAL-079 | 13-15 | ON TRACK | Reports record final tagged times of 265.793, 240.423, and 240.222 seconds against 900; ordinary package times stay below 600. Review-11 independently ran the final tagged gate. |

## Gate And Evidence Checks

I compared each report's `## Tests run` and `## Exact quality-gate results` with its spec's `## Exact Commands To Run` and QA-plan EVID-25 through EVID-28. Tasks 13 and 14 each retain PASS evidence, output tokens, and observed `exit_code: 0` for all five exact commands. Task 15 retains the same five, the 100k rebuild benchmark, and all three declared absence searches with observed/expected `exit_code: 1`. Earlier failed attempts remain in the reports. The report-only task-15 repair correctly did not claim a new product gate run. All three two-argument report-lint commands returned `PASS report-lint completed with no failures` (each terminated `exit_code: 0`).

| Independent check | Result | Reproduced evidence |
|---|---|---|
| Required `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`). |
| `go test -tags integration -run '^(TestP23HumanFixtures|TestP23TenNameActivityInventoryFromProducingLeaves|TestP23DecisionReferenceProtectedLifecycleRealProcess)$' -count=1 ./cmd/fid/integration` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t12.571s` (terminated `exit_code: 0`). This is a focused re-review run, not a substitute for the recorded full tagged gate. |
| `go test -run '^TestPhase23HumanLeafTextMatchesRecordedFacts$' -count=1 ./cmd/fid` | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t0.152s` (terminated `exit_code: 0`). |
| `go test -run '^(TestMigration0018DownMatchesFreshSchema17|TestMigration0018DownPreservesAgentRunAndFullRollback|TestMigration0018RoundTripIntegrity)$' -count=1 ./internal/db` | PASS | `ok  \tgithub.com/mbell-nzl/fid/internal/db\t0.308s` (terminated `exit_code: 0`). |
| Guide/manifest and provenance inventories | PASS | Independent read-only checks found 18/18 unique Human and JSON blocks matching the retained manifest; 66 provenance rows across 43 directly listed test files, with valid lines and 43 resolving introducing hashes. Of 47 test files changed by task 15, three are new tests and the remaining closure helper is cited with its journey row. |

## Per-Task Findings

### Task 13 - CLI surface, failure codes, and idempotency

#### Must-fix

- None.

#### Improvement

- None.

Observation: `p23_task13_fixture_replay_test.go:18-77` queries durable open Stops and linked Attentions independently of the response and compares the exact denial fields; both positive Stop-blocked cases and noncommitting retries are exercised. The failure schema and subjects match DEC-263 without widening StopDetails.

### Task 14 - Human review and judgment projections

#### Must-fix

- None.

#### Improvement

- **I-1, carried from review-11** [`internal/workflow/resume.go:1`]: the non-test file exceeds the quality plan's 400-line target. Fix direction: split the bounded judgment projection helpers during a later authorized causal edit.

Observation: The ordinary test pins exact text for all 18 operations; the tagged test uses the real CLI and reconstructs each Human view from JSON fields. The appendix's 18 verbatim pairs match the retained manifest, and the `PRIVATE_` prose and judgment qualifier checks align with `contracts.md` section 14. This closes review-11 MF-1 for VAL-058 and VAL-060.

### Task 15 - Migration 0018 and V1 judgment removal

#### Must-fix

- None.

#### Improvement

- None.

Observation: The down SQL follows `schema.md` section 5's nine ordered groups, including post-0016 V1 Concern shape, two restored kind triggers, the row-copying `agent_runs` rebuild, the 0011 view, and the four 0005 plus one 0016 reference guards. The independent DB rerun proves the specified fresh-17 schema comparison, Run preservation, full rollback and clean round trip. The 66-row provenance table in `task-15-done.md:79-156` supplies file/line, introducing commit and task, the specific accepted task-15 invalidation, and the exact executing command. `git blame 1d89aa4 -L 23,25` confirms the three old process-exit rows; `git blame 437d093 -L 586,592` confirms the old AC verification row; `git show 457f480:internal/db/ref_triggers_test.go` confirms the historical trigger cases. Their commit trailers match the reported tasks. This closes review-11 MF-2 (no VAL).

## Prior Finding Resolution

| Review-11 finding | Status | Fresh evidence |
|---|---|---|
| Task-14 MF-1 (VAL-058, VAL-060) | RESOLVED | 18/18 unique guide and manifest pairs; verbatim Human text, complete retained JSON fact maps with guide field subsets, exact Human reconstruction and no-proof/no-prose checks; independent ordinary and tagged reruns passed. |
| Task-15 MF-2 (no VAL) | RESOLVED | 66 per-row entries cover the retained converged tests, with path/line, introducing hash and task, specific invalidation and exact executing command; all cited hashes resolve, sampled historical rows match blame, and committed report-lint passes. |

Both prior Must-fix findings are closed. No Must-fix finding remains.

## Blocker Candidates

- None.

## Existing Escalation References

- ESC-23-007 (DEC-262), ESC-23-008 (DEC-263), and ESC-23-009 (DEC-264) are resolved. No escalation ID was allocated.

## Final Lines

Review written: project-docs/docs/phases/phase-23/qa/review-12.md
Must-fix: 0  Improvement: 1  Blocker candidates: 0
VERDICT: PASS
