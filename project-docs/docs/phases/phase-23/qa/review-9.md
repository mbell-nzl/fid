# Phase 23 - Mid-Phase Review 9 (tasks 10-12)

Reviewer: qa-reviewer
Date: 2026-09-27
Tasks covered: phase-23/task-10, phase-23/task-11, phase-23/task-12
Review type: INITIAL
Prior review: none

## Verdict

REPAIR. Two required assertions have bounded evidence or fixture gaps in their owning tasks. The window does not advance.

## Scope Reviewed

- Task specs and completion reports: `tasks/task-{10,11,12}.md` and `tasks/task-{10,11,12}-done.md` in this phase.
- Phase validation contract, QA plan sections 4, 5, 7, 8 and 22, `08_QUALITY_PLAN.md`, current context, and DEC-249, DEC-250, DEC-258, DEC-259 and DEC-260.
- Task 10 implementation and repair through `a425ffe22ec12558e953bbbb5cb265d22f28779e`; task 11 resolution commit `466745f2dc4505d60eb67ac50193e004532b7239` and authorized test repair `3dff4574a71438900100449e36d7044969ba9162`; task 12 checkpoint `0e4bb33cbe32fa9d4d83955b3121a4135edf78cd` and implementation `1679230ee4dca4d1cbffbb5b2732b7d1a9dbd06d`.
- Relevant source and tests in `internal/workflow`, `internal/review`, `internal/journal`, `internal/links`, `cmd/fid` and `cmd/fid/integration`; committed reports and their histories. No held-out scenario material was read.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 2
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-23-005 and ESC-23-006 (both resolved)
- Covered validation assertions: VAL-003, VAL-009, VAL-010, VAL-025, VAL-043, VAL-044, VAL-045, VAL-046, VAL-047, VAL-048, VAL-049, VAL-050, VAL-051, VAL-052, VAL-053, VAL-054, VAL-078, VAL-079, VAL-080, VAL-081

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence and limit |
|---|---|---|---|
| VAL-003 | 10, 12 | ON TRACK | Task 10 preserves the converged `human-outcome.recorded` parity row; task 12 removes the last fence keys and tests rejected malformed events against unchanged projection. |
| VAL-009 | 10 | ON TRACK | The repaired real-process journey checks gate-satisfied/review-unsatisfied denial, the action-specific guard, and all-satisfied completion (`p23_human_outcome_test.go:131-213`). |
| VAL-010 | 10 | ON TRACK | The same journey checks `not-applicable` with linked Review and Human Outcome Evidence and `available` for gate Evidence (`:149-206`). |
| VAL-025 | 10 | ON TRACK | Real Human Outcome positive and two mismatch denials are reconciled with task 6's other six cells in the report; the journey checks the Human cells through the producer. |
| VAL-043 | 10 | ON TRACK | The real-process authority journey checks four rejected forms with unchanged Event counts and a finalized, provenanced successful record (`:45-88`). |
| VAL-044 | 10 | ON TRACK | Real link and admission tests distinguish matching approved from wrong type, class and result. |
| VAL-045 | 10 | ON TRACK | Evidence record and evaluation assertions check the Human producer, trust source and artifact-availability values. |
| VAL-046 | 10 | NOT ON TRACK | MF-1: the cited artifact test has no Human Obligation, and the required repeat-while-open observation is absent. |
| VAL-047 | 11 | ON TRACK | Task 11's false capability state and working resolution path are evidenced; the true state is assigned to task 17. |
| VAL-048 | 11 | ON TRACK | The built-binary test rejects wrong class and result with no Event, then accepts a qualifying Outcome. |
| VAL-049 | 11 | ON TRACK | The real-process journey rejects an added input field and reads the original raise facts and resolution history. |
| VAL-050 | 11 | ON TRACK | The Condition-matching producer and operator Stops resolve; unrelated Stops remain open. |
| VAL-051 | 12 | ON TRACK | The independently rerun real-process journey checks digest mismatch, retained predecessor and pointer, and one active row. |
| VAL-052 | 12 | ON TRACK | The journey reads back all three relations and checks unknown requirement and relation rejection with unchanged Event count. |
| VAL-053 | 12 | NOT ON TRACK | MF-2: the binary, production catalogs and declared absence search are clear, but the portable fake still recognizes `requirement.synced`. |
| VAL-054 | 12 | ON TRACK | Real-process and command tests reject requirement endpoints and both removed relations and retain Issue links. |
| VAL-078 | 11 | ON TRACK | Built-binary and workflow tests exercise qualifying and non-qualifying durable facts for all four Condition types, including producer and operator Stops for the two repaired types. |
| VAL-079 | 10-12 | ON TRACK | Reports record package times within 600/900 s and terminated task gates; phase-end measurement remains assigned to final QA. Task 12 uses the DEC-259 command text. |
| VAL-080 | 10, 11 | ON TRACK | Journal tests check canonical and malformed per-type discriminator shapes, including the bare `review-recorded` form. |
| VAL-081 | 10, 11 | ON TRACK | Journal tests distinguish the nullable-binding rejection from shape rejection and retain the Phase 22 canonical forms. |

## Gate And Evidence Checks

Working directory for independent commands: `/home/mark/projects/fid`. The QA plan section 5 task gates match each task's exact block after the DEC-259 amendment takes effect at task 12. No additional window-specific full gate is required. These report rows were compared in both `## Tests run` and `## Exact quality-gate results`; focused tests do not replace the full gates.

| Task | Exact commands required by spec and QA plan | Report result |
|---|---|---|
| 10 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -timeout 30m -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` | All five final rows PASS with output token and observed `exit_code: 0` at report lines 512-532; earlier `-count=1` deviation remains disclosed. |
| 11 | The same five byte-identical task-11 commands | All five final DEC-258 rows PASS with token and observed `exit_code: 0` at report lines 55-71; earlier budget failures, cached attempt and midnight failure remain visible. DEC-258 authorizes the out-of-scope test-parallelism repair. |
| 12 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...`; `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/`; `rg -n --glob '!*_test.go' "requirement\.synced|last_synced_hash" cmd internal` | All seven final rows PASS with output token and observed termination at report lines 90-97. The declared search records expected and observed `exit_code: 1`; the benchmark self-asserts `<30000 ms/rebuild`. Earlier BLOCKED rows, tagged failure, format failure and report-lint failures remain visible. |

Independent mechanical checks: `report-lint.sh` returned `PASS report-lint completed with no failures` (observed `exit_code: 0`) for each of the three report/spec pairs, including `OUTCOME COMPLETE`. `git rev-parse --verify <hash>^{commit}` resolved the reported task 10, 11 and 12 implementation hashes and the task 12 report hash `b76312762867f18bf70798243991ef77567ba001`. The task 10 and 11 report commits `bc7e6f393450780137f3d89ad9763e382028ccd3` and `21f60f2e691b7ddc1a9f2440d23883c1a4f0eb8c` also resolve.

Reproduced named tests, each run uncached with `-count=1` (QA verification, not substituted task gates):

- Task 10: `go test -tags integration -count=1 -run '^TestP23HumanOutcomeRequiredRaisedByReviewRecord$' ./cmd/fid/integration` -> `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.481s` (terminated `exit_code: 0`). `go test -count=1 -run '^TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation$' ./cmd/fid` -> `ok  	github.com/mbell-nzl/fid/cmd/fid	0.916s` (terminated `exit_code: 0`); this green test does not create a Human Obligation.
- Task 11: `go test -tags integration -count=1 -run '^TestP23StopResolveHumanOutcomeConditionKeysOnConditionNotRaisingCode$' ./cmd/fid/integration` -> `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.577s` (terminated `exit_code: 0`).
- Task 12: `go test -tags integration -count=1 -run '^TestP23RequirementReferenceAndLinkRealProcess$' ./cmd/fid/integration` -> `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.949s` (terminated `exit_code: 0`). An earlier byte-identical focused invocation returned `(cached)` and was not used as fresh behavioral evidence.

Full phase gates and the phase-end `checkptr` control were left to the QA Lead because these findings concern narrow in-window assertions and fixtures.

## Per-Task Findings

### Task 10 - Provenanced Human Outcomes

#### Must-fix

- **MF-1 [cmd/fid/artifact_test.go:277; cmd/fid/envelope_test.go:490; tasks/task-10-done.md:255]** The report cites `TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation` as the `artifact.availability-observed` VAL-046 raising sequence. Its envelope comes from `envelopeCreateInput`, whose only Obligation is `gate-evidence`; the test checks artifact availability and never checks a `human-outcome-required` Stop or Attention. The contract also requires a repeat qualifying Event while the Stop is open, with no duplicate; the named task-10 journeys do not make that observation. Thus the four-sequence/both-state claim is incomplete despite passing tests. Owning task: task-10. Affected assertion: VAL-046. Repair boundary: add or converge in-scope task-10 tests and the task-10 report to prove the artifact flip on an Envelope with a sole remaining Human Obligation, its preceding non-qualifying state, exactly one Stop and Attention, a repeat qualifying Event without duplication, and resolution by the matching Outcome. Keep the other three producer proofs and accepted semantics. Required evidence: red/green focused output, direct Stop/Attention and Event observations, corrected per-sequence citations, and all five byte-identical task-10 gates with tokens and observed termination. Re-review scope: inspect the new fixture, report and diffs; rerun the artifact-focused proof and a representative affected gate.

#### Improvement

- None.

Observation: the repaired VAL-009 journey asserts the gate-satisfied/review-unsatisfied stage and the all-satisfied terminal path on one Envelope; its focused review-recorded test reproduced successfully.

### Task 11 - Attention resolution through Human Outcomes

#### Must-fix

- None.

#### Improvement

- None.

Observation: the DEC-250 diff routes manual and atomic judgment resolution through the same Condition predicate. The rerun checks a producer Stop and Condition-matching operator Stops, wrong facts, and unrelated open Stops. The report keeps the earlier VAL-079 timing failures and the DEC-258 repair authority; its final five exact gates pass.

### Task 12 - Requirement references and a single traceability model

#### Must-fix

- **MF-2 [cmd/fid/integration/p22_portable_fake_test.go:267; tasks/task-12.md:42; tasks/task-12-done.md:36]** `p22PortableEventTypes` still lists `requirement.synced`, and the fake archive reader at `p22_portable_fake_test.go:73` treats that list as its known Event families. The task requires removal of the Event and its replay fixtures, while VAL-053 and the report claim source-tree absence. The declared `rg` gate excludes `*_test.go`, so it cannot detect this retained fake entry. Provenance: `git log -S'requirement.synced'` shows this entry was introduced with the portable fake in `5ec69e6`, before Phase 23; task 12's removal invalidates it. Owning task: task-12. Affected assertion: VAL-053. Repair boundary: converge the in-scope integration fake's Event-family list to the accepted catalog and add a focused assertion that the removed family is rejected while retained families remain accepted; update the task-12 report's source-tree claim. Required evidence: introducing-commit citation, focused red/green fake-reader test output, source inventory including test files, and all seven byte-identical task-12 gates with tokens and observed termination. Re-review scope: inspect the fake-reader diff and catalog parity, rerun its focused test and the declared absence search, then compare the seven report rows to the spec.

#### Improvement

- None.

Observation: DEC-260's `internal/review/requirement.go` diff removes only the two capability checks while retaining trusted Actor resolution and digest verification under the lock. The independently rerun real-process journey covers mismatch, supersession, traceability and legacy-link rejection.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-1 | task-10 | VAL-046 | In-scope task-10 test/report changes for the missing artifact producer and repeat-while-open sequence; causal production repair only if the new test exposes a task-10 bug. | Focused red/green, four producer sequences with both states and no duplicate, corrected citations, five exact gates. | Inspect test behavior and report; rerun artifact proof and representative affected verification. |
| MF-2 | task-12 | VAL-053 | In-scope integration fake and task-12 report convergence; retain the accepted production/Event catalog semantics. | Provenance, removed-family rejection plus retained-family positive, source inventory including tests, seven exact gates. | Inspect fake and catalog parity; rerun focused fake test and absence search; audit report gates. |

## Blocker Candidates

- None. Both repairs fit their existing task test and report scope and require no new product, interface, gate or validation-contract decision.

## Existing Escalation References

- ESC-23-005: resolved under DEC-258/259; task 11 report preserves the failed timing history and records its authorized repair.
- ESC-23-006: resolved under DEC-260; task 12 report preserves the earlier failures and the corrected one-line expected-nonzero declaration. No escalation is open.

## Final Lines

Review written: project-docs/docs/phases/phase-23/qa/review-9.md
Must-fix: 2  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
