# Phase 21 Final Review - One-Issue Operational Walking Skeleton

Reviewer: QA Lead
Date: 2026-09-03
Review type: INITIAL
Prior review: none

Phase HEAD reviewed: `5d8b93c715a07fc04f17f4529225d91e344cdab8`
Worktree state at review start: clean except the Phase Runner's own
`project-docs/docs/phases/phase-21/run-ledger.md` modification, which is the
runner's permitted dirt.
Session scratch: allocated fresh with `mktemp -d` per attempt; cleanup is not
recorded as evidence.

## Verdict

BLOCK

The phase product is sound: all seven `AGENTS.md` phase-gate commands passed in
this session with observed terminations, and 85 of 87 validation assertions
PASS on session-captured evidence. The phase does not advance because two
canonical specification documents contradict the shipped Phase 21 interface,
and one of those contradictions makes `VAL-002` uncertifiable against its own
cited source requirement. Reconciling `03_TRD_SCHEMA.md` §8.2 requires
specification authority and lies outside every existing phase task's
`Files To Create Or Modify` scope, so it is a blocker candidate rather than a
`REPAIR`. The second defect (MF-2) is bounded and decision-free inside task 13
and is recorded in the Repair Plan so a single corrective pass can close both.

No red test, failed gate, incomplete in-scope implementation, or report defect
contributed to this verdict.

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Validation contract:
  `project-docs/docs/phases/phase-21/qa/validation-contract.md` (VAL-001..087)
- QA plan: `project-docs/docs/phases/phase-21/qa/qa-plan.md` (Ready, §7.9)
- Task specs and completion reports:
  `project-docs/docs/phases/phase-21/tasks/` (22 tasks, 22 reports)
- Mid-phase reviews: `qa/review-1.md`..`qa/review-7b.md`;
  independent checkpoint `qa/group-1-checkpoint-review.md`
- Planning authorities: `../contracts.md`, `../public-interface.md`,
  `../v1-convergence-inventory.md`
- Current decisions: `project-docs/docs/current/decisions_active.md`
  (DEC-152..DEC-177); escalations: `tl_escalations_current.md` (none open)

### Ledger state

`project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-21/run-ledger.md`
returned `VALID state=ADVANCEABLE` (terminated `exit_code: 0`).

## Phase-Level Gate Results

Executed in this session in `AGENTS.md` order from the repository root
`/home/mark/projects/fid`, each as its own process with its exit status
observed from the command's own result.

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`) |
| Lint/static analysis | `golangci-lint run ./...` | PASS | `0 issues.` (terminated `exit_code: 0`) |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | `ok  	github.com/mbell-nzl/fid/internal/workflow	28.710s	coverage: 27.4% of statements` (terminated `exit_code: 0`); 21 packages, 0 `FAIL`, 0 `SKIP`, first line `ok  	github.com/mbell-nzl/fid/cmd/fid	349.018s	coverage: 77.3% of statements` |
| Coverage report | `go tool cover -func=coverage.out` | PASS | `total:										(statements)				61.2%` (terminated `exit_code: 0`) |
| Integration | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.711s` (terminated `exit_code: 0`) |
| End-to-end | `N/A` | N/A | No end-to-end or held-out scenario harness is accepted in the specs (`AGENTS.md` Commands). |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | `ok  	github.com/mbell-nzl/fid/internal/workflow	2.210s` (terminated `exit_code: 0`); includes `BenchmarkRebuildWorkspace100kEvents-6   	       1	26940902254 ns/op	     26941 ms/rebuild` with the retained per-sample `<30 s` `b.Fatalf` |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`) |
| Full phase check | the seven commands above in order | PASS | seven captured tokens, each terminated `exit_code: 0` (QA-87) |

No numeric coverage threshold is defined for this project
(`08_QUALITY_PLAN.md` §7 treats coverage as diagnostic); the observed `total:`
line is the required evidence. `internal/claims` and `internal/policy` report
`coverage: 0.0%` because they hold no in-package tests; both are live packages
exercised through `cmd/fid` and `internal/workflow`.

## Validation Contract Outcomes

Every `VAL-001..087` has one outcome. Mid-phase rows cite the fresh independent
review that closed their window plus this session's re-verification at phase
HEAD: the repository-wide `-race` gate and the tagged integration gate above
executed every named backing test, and each cited test name was confirmed
present at HEAD with `go test -list '.*' ./...` and
`go test -tags integration -list '.*' ./cmd/fid/integration/...`.

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Reviewer | `review-1b.md` (MF-1 RESOLVED: exhaustive canonical 0013 inventory with extra-object/CHECK/FK rejection). Session: `TestMigration0013OperationalKernelInventoryAndRoundTrip` and `TestOperationalKernelInventoryComparisonRejectsSchemaDrift` present at HEAD; `internal/db` `ok` under `-race`. Spec/migration diff: `03_TRD_SCHEMA.md` `0013` up block is line-identical to `migrations/0013_p21_operational_kernel.up.sql` (296 lines each) and `0013` down is identical (27 lines) |
| VAL-002 | BLOCKED | QA Reviewer | See MF-1. Shipped behavior conforms to task-20's accepted binding constraint and `v1-convergence-inventory.md` §§4, 6 (`TestMigration0014RetiresV1AuthorityExecutionWithoutReconstruction`, `internal/db` `ok`), and the observed final schema has `claims` with index `claims_one_active_mutation` and no `roles`/`gate_runs`/`work_claims`. It does **not** match the assertion's cited source `03_TRD_SCHEMA.md` §8.2 `0014` up/down SQL. Two adopted authorities disagree on "exactly", so PASS is not certifiable |
| VAL-003 | PASS | QA Reviewer | `review-1b.md`; `TestEventV2RejectsUnknownTopLevelKey`, `TestBuildProjectInitializedEventMatchesV2Golden` present at HEAD; `internal/journal` and `internal/bootstrap` `ok` under `-race` |
| VAL-004 | PASS | QA Reviewer | `review-1b.md`; `internal/journal/event_v2_test.go` sequence-gap append/replay rejection; `internal/journal` `ok` under `-race` |
| VAL-005 | PASS | QA Reviewer | `review-1b.md` plus checkpoint MF-1 resolution (`0444ccc fix(activation): replay stored activation results`) verified in the second independent checkpoint review; `internal/journal` and `cmd/fid` `ok` under `-race` |
| VAL-006 | PASS | QA Reviewer | `review-1b.md`; `internal/journal/idempotency_test.go:58` digest-conflict, no second Event; `internal/journal` `ok` |
| VAL-007 | PASS | QA Reviewer | `review-1b.md`; `internal/journal/idempotency_test.go:82` concurrent same-key, one Event identity, one replay under `-race` |
| VAL-008 | PASS | QA Reviewer | `review-1b.md`; `internal/bootstrap/init_test.go:15` and `cmd/fid/init_test.go:428`; `internal/bootstrap` and `cmd/fid` `ok` |
| VAL-009 | PASS | QA Reviewer | `review-1b.md`. Session: live `/tmp/fid --json capabilities` in a fresh `mktemp -d` workspace returned exactly 12 rows sorted by `capability_id`, byte-matching the `contracts.md` §2 matrix field by field (`available`, `limits`, `schema_min`, `schema_max`, `trust_source`); terminated `exit_code: 0` |
| VAL-010 | PASS | QA Reviewer | `review-1b.md` (MF-2 RESOLVED: seven named executable surface probes). Session: the same live matrix shows exactly the seven required `available:false` rows - `attention.v1`, `export.portable.v1`, `human.outcomes.v1`, `kit.activation.v1`, `reviews.coverage.v1`, `telemetry.trajectory.v1`, `work.claims.v1` |
| VAL-011 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/activation_test.go:200`; `cmd/fid` `ok` under `-race` |
| VAL-012 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/activation_test.go:233`, `:209`, `:213` |
| VAL-013 | PASS | QA Reviewer | `review-2b.md`. Session absence control rerun at HEAD: `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"recover|newActivationRecover|activation\.recover' cmd/fid internal` - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`) |
| VAL-014 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/activation_test.go:397` |
| VAL-015 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/activation_test.go:105` |
| VAL-016 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/gate_test.go:26` |
| VAL-017 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/gate_test.go:71` |
| VAL-018 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/gate_test.go:112` |
| VAL-019 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/envelope_test.go:13` |
| VAL-020 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/envelope_test.go:89` |
| VAL-021 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/envelope_test.go:97` |
| VAL-022 | PASS | QA Reviewer | `review-2b.md`, `cmd/fid/envelope_test.go:42` |
| VAL-023 | PASS | QA Reviewer | `review-3b.md`, `cmd/fid/authority_test.go` all-match and single-dimension failure matrix; `cmd/fid` `ok` |
| VAL-024 | PASS | QA Reviewer | `review-3b.md` seeded project-local role-attribution adversary |
| VAL-025 | PASS | QA Reviewer | `review-3b.md`. Session: the `--by` absence control terminates `exit_code: 1` at HEAD (see VAL-075) |
| VAL-026 | PASS | QA Reviewer | `review-3b.md` fixed-clock nonfuture Event-v2 ordering |
| VAL-027 | PASS | QA Reviewer | `review-3b.md` source-budget inventory; `08_QUALITY_PLAN.md` file/function budgets |
| VAL-028 | PASS | QA Reviewer | `review-3b.md` retained-family operation/read matrices |
| VAL-029 | PASS | QA Reviewer | `review-3b.md`; `TestDefineAppendsOneAtomicWorkflowDefinedEvent` present at HEAD, `internal/workflow` `ok` |
| VAL-030 | PASS | QA Reviewer | `review-3b.md` plus checkpoint MF-2 resolution (`1d89aa4 fix(cli): map failure categories to process exits`, `capability-unavailable` exit 5) verified in the second independent checkpoint review |
| VAL-031 | PASS | QA Reviewer | `review-3b.md`; `98217f4` removed `role` normalization; `cmd/fid/stats_test.go:37`; `internal/stats` `ok` |
| VAL-032 | PASS | QA Reviewer | `review-4b.md`, `cmd/fid/generic_claim_test.go:65` |
| VAL-033 | PASS | QA Reviewer | `review-4b.md`, `cmd/fid/generic_claim_test.go:201` |
| VAL-034 | PASS | QA Reviewer | `review-4b.md`, `cmd/fid/generic_claim_test.go:242` |
| VAL-035 | PASS | QA Reviewer | `review-4b.md`, `cmd/fid/generic_run_test.go:56` |
| VAL-036 | CARRY-FORWARD | Human | Sole Human assertion. Deterministic substrate verified this session: `cmd/fid/testdata/val-036/human-session.txt` and `live-frames.jsonl` exist and are byte-clean - `LC_ALL=C grep -P '[\x00-\x08\x0b\x0c\x0e-\x1f\x7f]'` over each terminated `exit_code: 1` with no output, so no raw control byte can reach a terminal; frames are sequence 1..9 across `stdout`/`stderr`/`control`/`terminal` with exactly one `output-truncated`, one `redaction-applied`, one `artifact-rejected`, and one last `terminal` frame. The four perception observables HUMAN-01a..01d are presented unanswered in `qa/human-checklist.md` and are not pre-answered, delegated, or inferred here |
| VAL-037 | PASS | QA Reviewer | `review-4b.md`, `cmd/fid/generic_run_test.go:149` |
| VAL-038 | PASS | QA Reviewer | `review-4b.md`, `cmd/fid/task12_execution_transition_test.go:114` |
| VAL-039 | PASS | QA Reviewer | `review-4b.md`, `cmd/fid/task12_execution_transition_test.go:53` and `:248` |
| VAL-040 | PASS | QA Reviewer | `review-4b.md` dual-call evaluator and terminal no-mutation assertions |
| VAL-041 | PASS | QA Reviewer | `review-4b.md`, `TestBuildIssueRevisionEffectsRejectsInvalidInputsWithoutEffects` |
| VAL-042 | PASS | QA Lead | `group-1-checkpoint-review.md` second independent review: `VAL-042: PASS`, `CHECKPOINT VERDICT: PASS` at Group 1 HEAD `5aeb19b`, before any task 13 dispatch (`c8ccc8b` precedes `d2feab0`) |
| VAL-043 | PASS | QA Lead | `group-1-checkpoint-review.md` second independent review: `VAL-043: PASS`, with FEAS-01..FEAS-19 all PASS from a fresh `mktemp -d` allocation |
| VAL-044 | PASS | QA Reviewer | `review-5b.md`; `TestExecuteAppliesAcceptedGatePolicy` rerun this session: `ok  	github.com/mbell-nzl/fid/internal/evidence	2.261s` (terminated `exit_code: 0`) |
| VAL-045 | PASS | QA Reviewer | `review-5b.md`, `cmd/fid/gate_test.go` one-field unsupported-policy matrix. Cross-paired with VAL-044 per DEC-158 R-14 |
| VAL-046 | PASS | QA Reviewer | `review-5b.md`. Session: both declared containment lexicons rerun at HEAD, each `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`) - one over `cmd internal`, one over `contracts.md`/`public-interface.md` |
| VAL-047 | PASS | QA Reviewer | `review-5b.md`; `TestRunLiveFramesInterleavedOutputControlsAndFailure` rerun this session, `internal/evidence` `ok` (terminated `exit_code: 0`) |
| VAL-048 | PASS | QA Reviewer | `review-5b.md`. Session: the committed `live-frames.jsonl` golden has exactly the three control kinds `output-truncated`, `redaction-applied`, `artifact-rejected` in required order and no unknown kind |
| VAL-049 | PASS | QA Reviewer | `review-5b.md`. Session canary scan: `rg -n "CANARY-val-036|CANARY-live-secret|literal-val-036|literal-live-secret|redaction-canary" cmd/fid/testdata cmd/fid/integration/testdata` - `exit 1 (no output)` (terminated `exit_code: 1`), so no credential or registered literal survives into any committed machine, Human, or conformance surface; the fixture retains `[control] redaction-applied replacements=2` |
| VAL-050 | PASS | QA Reviewer | `review-5b.md`; `TestRunLiveDrainsBothPipesAfterTruncationAndPreservesOutcome` rerun this session, `internal/evidence` `ok` |
| VAL-051 | PASS | QA Reviewer | `review-5b.md`; `TestGateRunAttemptAppendFailureDoesNotSpawn` present at HEAD, `cmd/fid` `ok` |
| VAL-052 | PASS | QA Reviewer | `review-5b.md` complete Receipt binding assertion, extended by DEC-177 durable rejection facts; `TestGateRunReplayRetainsRejectedArtifacts` rerun this session (`ok  	github.com/mbell-nzl/fid/cmd/fid	46.323s`, terminated `exit_code: 0`) |
| VAL-053 | PASS | QA Reviewer | `review-5b.md` failed/spawn-failed/timed-out/signaled/fid-error outcome and exit classification |
| VAL-054 | PASS | QA Reviewer | `review-5b.md` fail-closed rejection of caller-authored `producer:'caller'` Evidence |
| VAL-055 | PASS | QA Reviewer | `review-6c.md`; `TestGateRunPublishesSanitizedArtifactByDigest` rerun this session, `cmd/fid` `ok` |
| VAL-056 | PASS | QA Reviewer | `review-6c.md`, `cmd/fid/artifact_test.go:74`; `TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable` and `TestArtifactCandidatePathRejectsEveryUnsafeCandidate` rerun this session |
| VAL-057 | PASS | QA Reviewer | `review-6c.md`, `cmd/fid/artifact_test.go:252`; `TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation` rerun this session |
| VAL-058 | PASS | QA Reviewer | `review-6c.md`, `cmd/fid/evidence_test.go:180`, `:203` |
| VAL-059 | PASS | QA Reviewer | `review-6c.md`, `cmd/fid/evidence_test.go:18` |
| VAL-060 | PASS | QA Reviewer | `review-6c.md`, `cmd/fid/evidence_test.go:25` |
| VAL-061 | PASS | QA Reviewer | `review-6c.md`, `cmd/fid/evidence_test.go:38` |
| VAL-062 | PASS | QA Reviewer | `review-6c.md`, `cmd/fid/evidence_test.go:46` |
| VAL-063 | PASS | QA Reviewer | `review-6c.md`; `TestPhase21AttentionBoundary` present at HEAD and green in the `cmd/fid` package result (DEC-172 named test, not a zero-match search). Session live matrix confirms the required unavailable `attention.v1` row |
| VAL-064 | PASS | QA Reviewer | `review-6c.md`; `TestEnvelopeCompleteReusesQualifiedTerminalEvidenceAfterRunRelease` present at HEAD |
| VAL-065 | PASS | QA Reviewer | `review-6c.md`, `cmd/fid/task18_completion_resume_test.go:16` and `:54` |
| VAL-066 | PASS | QA Reviewer | `review-6c.md` positive completion path after Run finish and Claim release |
| VAL-067 | PASS | QA Reviewer | `review-6c.md`; `TestEnvelopeCompleteReturnsAllCurrentReasonsWithoutMutation` present at HEAD |
| VAL-068 | PASS | QA Reviewer | `review-6c.md`, `cmd/fid/task18_completion_resume_test.go:201` |
| VAL-069 | PASS | QA Reviewer | `review-7b.md`; `TestCreateIsByteIdenticalWithCompleteDeterministicInventory` present at HEAD, `internal/export` `ok` |
| VAL-070 | PASS | QA Reviewer | `review-7b.md` marks this `PASS`: `TestValidateRejectsEveryArchiveInvalidCheckWithoutMutation` at `internal/export/export_test.go:97` covers all eleven invalid-archive checks with pre/post workspace tree digests |
| VAL-071 | PASS | QA Reviewer | `review-7b.md`. Session: `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"import|newExportImport|export\.import' cmd/fid internal` - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`); live matrix shows `export.portable.v1` `available:false` |
| VAL-072 | PASS | QA Lead | QA-81 this session. All 57 `public-interface.md` §§4.1/4.2 operation rows have a corpus case file (`comm` of the extracted operation list against `testdata/p21_conformance/operations` produced no spec-only entries), plus 4 §5 fail-closed leaves; all 34 §7 failure codes have a case file; all 6 terminal outcomes, 3 control kinds, and both process streams have live fixtures. `ls cmd/fid/integration/testdata/p21_conformance` - `operations` (terminated `exit_code: 0`). `TestPhase21ConformanceCorpusIsCompleteAndByteExact` byte-compares every regenerated fixture and rejects unexpected or missing files; the tagged gate terminated `exit_code: 0` |
| VAL-073 | PASS | QA Lead | QA-82 this session, QA-plan §12 pipeline over the committed corpus in a fresh `mktemp -d`. Extraction produced 670 kept rows and 326 distinct values. Positive control: 142 closed members independently extracted from `public-interface.md` §§3, 7, 8 - **missing count 0**. Negative control: `grep -E $'\t[^\t]*[a-z0-9]_[a-z0-9]' "$SCRATCH/kept.tsv"` - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`) |
| VAL-074 | PASS | QA Reviewer | `review-7b.md`. Session QA-93 re-verification: both broad classified inventories terminate `exit_code: 0` at phase HEAD and every produced line maps to task-20's preserved per-row disposition (`UNCLASSIFIED=0`); per-file match distribution is byte-identical to task-20 HEAD `2d34575` and the Event-type literal set is identical (21 values, empty `comm` both ways) |
| VAL-075 | PASS | QA Reviewer | `review-7b.md`. Session: all ten declared expected-non-zero controls rerun at phase HEAD - five `rg` absence searches each `exit 1 (no output)` (expected `exit_code: 1`), and `ls internal/dev`, `ls internal/runs`, `ls cmd/fid/dev.go`, `ls cmd/fid/policy.go`, `ls cmd/fid/issue_complete.go` each `exit 2 (no output)` (expected `exit_code: 2`) |
| VAL-076 | PASS | QA Reviewer | `review-7b.md`; `TestFinalBinaryInitAndRebuildRemainRoleFree` present at HEAD and green in the tagged gate. Session: fresh `/tmp/fid init` in `mktemp -d` returned `exit_code: 0` and `SELECT name FROM sqlite_master WHERE type='table' AND name IN ('claims','work_claims','roles','gate_runs')` returned only `claims` |
| VAL-077 | PASS | QA Reviewer | `review-7b.md`; same tagged test deletes `.fid/fid.sqlite`, rebuilds, and reasserts role-free bootstrap and capability parity |
| VAL-078 | PASS | QA Reviewer | `review-7b.md`; `TestCheckReportsFinalSchemaIntegrityClassesWithoutRepair` present at HEAD, `internal/maintenance` `ok` |
| VAL-079 | PASS | QA Reviewer | `review-7b.md`; same test, `attempt-receipt-integrity` and `artifact-integrity` findings |
| VAL-080 | PASS | QA Reviewer | `review-7b.md`; same test, `frozen-state-present` and `removed-state-residue` with unchanged journal bytes and SQLite `total_changes` |
| VAL-081 | PASS | QA Reviewer | `review-7b.md`; `TestFinalSchemaNamedReadsAndExportMatchAfterRebuild` rerun this session (`ok  	github.com/mbell-nzl/fid/cmd/fid	46.323s`, terminated `exit_code: 0`) |
| VAL-082 | PASS | QA Reviewer | `review-7b.md`; same test compares export validation and exact export bytes across rebuild. Independently reinforced by QA-85's before/after `archive_digest` equality assertion in the tagged journey |
| VAL-083 | PASS | QA Lead | QA-80/QA-90 this session on the recorded reference instance. Three consecutive fresh one-iteration processes of `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal`, run with no concurrent load: `27018 ms/rebuild` (terminated `exit_code: 0`), `27233 ms/rebuild` (terminated `exit_code: 0`), `27448 ms/rebuild` (terminated `exit_code: 0`). All three below 30 s; no best-of or median sampling. Identity commands all terminated `exit_code: 0` and match `02_TRD.md` §9 exactly: `go version go1.26.2 linux/amd64`; `linux`/`amd64`; `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64`; `6`; `model name	: AMD Ryzen 7 5700G with Radeon Graphics`; `MemTotal:        8133236 kB` |
| VAL-084 | PASS | QA Lead | QA-83 this session. `TestPhase21RealProcessJourneyStartsOneActivatedEnvelope` drives the real built binary through, in order, `issue transition --to started` (nonterminal), `gate run`, `evidence link` then `evidence evaluate` (`status == "satisfied"` asserted), `issue transition --to completed` (terminal-success), `run finish`, `claim release`, `envelope complete`, `resume`, `export create`, `export validate`. No mock or in-process substitute. `go test -tags integration -race ./cmd/fid/integration/...` - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.711s` (terminated `exit_code: 0`) |
| VAL-085 | PASS | QA Lead | QA-84 this session. `p21AssertUndeclaredTransitionRejected` requires process exit `4`, stderr containing `"code":"transition-not-permitted"` and `workflow-transition-invalid`, empty stdout, and a zero Event delta, asserted both before and after rebuild. The 54-file `failures/` corpus covers all 34 accepted §7 failure codes with their category, channel, and exit; regenerated and byte-compared by `TestPhase21ConformanceCorpusIsCompleteAndByteExact`. Tagged gate terminated `exit_code: 0` |
| VAL-086 | PASS | QA Lead | QA-85 this session. `p21AssertJournalTypes` asserts an exact ordered 22-Event list (length and per-index equality) pre-rebuild and the byte-identical list post-rebuild, so post-rebuild commands appended zero Events. The list matches `contracts.md` §16 family by family and contains exactly one `envelope.created`, one `workflow.defined`, and one `workflow.assigned`. Post-rebuild command trace is limited to `rebuild`, `capabilities`, `envelope get`, `receipt get`, `receipt list`, `evidence evaluate`, `resume`, `export create`, `export validate`, and the negative transition check; `archive_digest` is asserted equal before and after. Tagged gate terminated `exit_code: 0`. Note: the accepted count is 22 per DEC-174, superseding the literal `20` in QA-plan §7.9 QA-85 |
| VAL-087 | PASS | QA Lead | QA-86 this session. Seven review artifacts plus the checkpoint exist and `git log --reverse` ordering places each after its window's last task report and the checkpoint between task 12 and task 13 (`11396b4`/`6caaa6f` -> `9a4ed61`/`c8ccc8b` -> `d2feab0`). The ledger holds 40 `TASK ... COMPLETE` events and 40 matching `- report-lint: PASS - OUTCOME COMPLETE` records, one per completion, before each next dispatch. All 22 committed reports pass a fresh two-argument `report-lint.sh` in this session, each `PASS report-lint completed with no failures` (terminated `exit_code: 0`), which mechanically verifies required-command coverage and observed-termination formatting per report. `workflow-sequence-lint.sh` returned `VALID state=ADVANCEABLE` (terminated `exit_code: 0`) |

Totals: 85 `PASS`, 1 `BLOCKED` (VAL-002), 1 `CARRY-FORWARD` (VAL-036, Human).

## Requirement Compliance

QA-88. Each row names the assertion and the artifact that satisfies the cited
requirement. Negative and unavailable requirements are listed explicitly.

| Requirement | Status | Evidence |
|---|---|---|
| PRD-FOUND-001, PRD-FOUND-005 (capability discovery) | MET | VAL-009, VAL-010; live 12-row matrix byte-matching `contracts.md` §2 |
| PRD-FOUND-002 (Event v2, sequences) | MET | VAL-003, VAL-004, VAL-026; `internal/journal` green under `-race` |
| PRD-FOUND-003 (role-free bootstrap, rebuild) | MET | VAL-008, VAL-076, VAL-077, VAL-081; fresh init/rebuild observation |
| PRD-FOUND-004 (wire conformance) | MET | VAL-072; 61-case operation corpus covering all 57 §4.1/4.2 rows |
| PRD-FOUND-005 (Human rendering) | PARTIAL | VAL-036 deterministic substrate verified; perception judgment carried to the Human checklist |
| PRD-FOUND-006 (Activation lifecycle) | MET | VAL-011, VAL-012 |
| PRD-AUTH-001, PRD-AUTH-002 (activated profile authority) | MET | VAL-014 |
| PRD-AUTH-003, PRD-AUTH-004 (authorization, no caller augmentation) | MET | VAL-015, VAL-023, VAL-024, VAL-074, VAL-085 |
| PRD-ISSUE-001..008 (retained Issue mutations) | MET | VAL-025, VAL-028; `--by` absence control `exit_code: 1` |
| PRD-ISSUE-007 (ordered acceptance criteria) | MET | VAL-028; AC add/list retained, AC verification fail-closed under VAL-030 |
| PRD-WF-003, PRD-WF-004, PRD-WF-005 (transitions, evaluator, reasons) | MET | VAL-029, VAL-038, VAL-039, VAL-040, VAL-062, VAL-064, VAL-065 |
| PRD-ENV-001..006, PRD-ENV-008 (Envelope, Obligation, budgets, completion) | MET | VAL-019..VAL-022, VAL-066, VAL-067 |
| PRD-RUN-001, PRD-RUN-002, PRD-RUN-004 (Claim binding, one active mutating Claim, versioned lifecycle) | MET | VAL-032, VAL-033, VAL-034; final schema index `claims_one_active_mutation` observed |
| PRD-RUN-003, PRD-RUN-005 (Run binding and terminalization) | MET | VAL-035, VAL-037 |
| PRD-RUN-014, PRD-PROJ-001, PRD-PROJ-002 (resume projection) | MET | VAL-068; `fid resume` exercised in the tagged journey |
| PRD-GATE-001..007, PRD-GATE-009 (Registration, executor, live output, Receipt, evidence sufficiency) | MET | VAL-016..VAL-018, VAL-044..VAL-054 |
| PRD-EVID-001..005 (Evidence linkage, validity, freshness, Artifacts) | MET | VAL-055..VAL-061 |
| PRD-OPS-003 (maintenance detection) | MET | VAL-078, VAL-079, VAL-080 |
| PRD-OPS-007 (redaction and unsafe-material handling) | MET | VAL-049, VAL-056; session canary scan `exit_code: 1` |
| PRD-REC-001, PRD-REC-002 (export create/validate) | MET | VAL-069, VAL-070, VAL-082 |
| PRD-HUM-002 (no routine Attention) | MET (negative) | VAL-063; `attention.v1` remains `available:false` |
| PRD-AC-002 | N/A - unresolvable citation | No `PRD-AC-*` identifier exists in `01_PRD.md`; see IMP-1. VAL-030's behavior is satisfied and its second citation `PRD-HUM-001` resolves |
| PRD-CLM-001..005 | N/A - unresolvable citation | No `PRD-CLM-*` identifier exists in `01_PRD.md`; see IMP-1. The live equivalents `PRD-RUN-001/002/004` are MET above |
| PRD-HUM-001 (Human Outcomes never deterministic proof) | MET (negative) | VAL-030; `human.outcomes.v1` `available:false`, judgment writers fail closed at exit 5 |
| PRD-REV-* (Review coverage) | DEFERRED | `reviews.coverage.v1` `available:false`, Phase 23 |
| PRD-TELEM-* (telemetry) | DEFERRED | `telemetry.trajectory.v1` `available:false`, Phase 25 |
| TRD §4 (idempotency, kebab-case public values) | MET | VAL-005, VAL-006, VAL-007, VAL-073 |
| TRD §3.1, DEC-153 DP-3 (no containment claim) | MET (negative) | VAL-046; both containment lexicons `exit_code: 1` |
| TRD §9 (100k rebuild target) | MET | VAL-083; three fresh samples 27018/27233/27448 ms on the reference instance |
| TRD Schema §8.2 (accepted slice DDL) | PARTIAL | `0013` up/down blocks are line-identical to the shipped migration; the `0014` up/down blocks are not - see MF-1 and VAL-002 |
| DEC-153 DP-5, PRD-REC (import absence) | MET (negative) | VAL-071; import absence control `exit_code: 1` |
| ADR-0011, DEC-157 (V1 removal) | MET | VAL-002 behavior, VAL-074, VAL-075; ten absence controls at their declared statuses |
| Quality Plan §8 (cadence and controls) | MET | VAL-087 |
| Quality Plan budgets (file/function size) | MET | VAL-027; `golangci-lint run ./...` `0 issues.` |

## Deterministic Checklist Execution

QA-01..QA-80 were executed by the named QA Reviewers in their review sessions
and are re-verified here at phase HEAD by (a) the repository-wide `-race` gate
and the tagged integration gate, which execute every backing test, (b)
`go test -list` confirmation that each cited test name still exists at HEAD,
and (c) fresh reruns of the expected-non-zero controls and the security-relevant
focused suites. QA-81..QA-93 were executed directly in this session.

| QA ID | Result | Evidence |
|---|---|---|
| QA-01..QA-09 | PASS | `review-1b.md` (`VERDICT: PASS`); this session `internal/db`, `internal/journal`, `internal/bootstrap`, `internal/cli`, `internal/workflow`, `cmd/fid` all `ok` under `-race`; live capability matrix byte-match |
| QA-10..QA-21 | PASS | `review-2b.md` (`VERDICT: PASS`); `cmd/fid` `ok`; VAL-013 absence control `exit_code: 1` |
| QA-22..QA-30 | PASS | `review-3b.md` (`VERDICT: PASS`); `cmd/fid`, `internal/stats` `ok`; checkpoint MF-2 exit-5 repair verified in the checkpoint rerun |
| QA-31..QA-40 | PASS | `review-4b.md` (`VERDICT: PASS`); `cmd/fid` `ok` |
| QA-41..QA-42 | PASS | `group-1-checkpoint-review.md` second independent review (`CHECKPOINT VERDICT: PASS`), VAL-042 and VAL-043 separately decided, FEAS-01..FEAS-19 PASS from fresh scratch |
| QA-43..QA-54 | PASS | `review-5b.md` (`VERDICT: PASS`); this session `go test -race -count=1 ./internal/evidence -run '^(TestArtifactCandidatePathRejectsEveryUnsafeCandidate|...|TestRejectedArtifactsResolvePathsUnderExecutionCWD)$'` - `ok  	github.com/mbell-nzl/fid/internal/evidence	2.261s` (terminated `exit_code: 0`), 8/8 `--- PASS` |
| QA-55..QA-68 | PASS | `review-6c.md` (`VERDICT: PASS`); this session `go test -race -count=1 ./cmd/fid -run '^(TestGateRunPublishesSanitizedArtifactByDigest|...|TestFinalSchemaNamedReadsAndExportMatchAfterRebuild)$'` - `ok  	github.com/mbell-nzl/fid/cmd/fid	46.323s` (terminated `exit_code: 0`) |
| QA-69..QA-71 | PASS | `review-7b.md` (`VERDICT: PASS`); `internal/export` `ok`; import absence control `exit_code: 1` |
| QA-72..QA-79 | PASS | `review-7b.md`; VAL-002 evidence is recorded but see MF-1 for its source-authority contradiction; VAL-074/075 controls rerun at HEAD |
| QA-80 | PASS | Three fresh consecutive 100k samples 27018/27233/27448 ms, each terminated `exit_code: 0`, three distinct captured result lines from three separate processes taken under no concurrent load (addresses RISK-02 and RISK-03); six identity commands captured and matching `02_TRD.md` §9 |
| QA-81 | PASS | VAL-072 row above: 57/57 §4.1/4.2 operation rows present in the corpus with no spec-only entry, 34/34 §7 failure codes present, all live-frame variants present, byte-exact regeneration test green, `ls cmd/fid/integration/testdata/p21_conformance` terminated `exit_code: 0` |
| QA-82 | PASS | VAL-073 row above: positive control missing-member count **0** over 142 independently extracted closed members; negative control `exit 1 (no output)` (terminated `exit_code: 1`); corpus non-vacuous at 670 kept rows / 326 distinct values |
| QA-83 | PASS | VAL-084 row above; ordered real-binary journey; tagged gate `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.711s` (terminated `exit_code: 0`) |
| QA-84 | PASS | VAL-085 row above; negative matrix with exact codes, categories, exits, and asserted absence of the forbidden effect |
| QA-85 | PASS | VAL-086 row above; exact ordered 22-Event pre- and post-rebuild assertion, one `envelope.created`, permitted post-rebuild command trace, equal `archive_digest`. Accepted count is 22 per DEC-174; the QA-plan literal `20` is superseded |
| QA-86 | PASS | VAL-087 row above; 40/40 ledger report-lint records, 22/22 fresh report lints `PASS` with `exit_code: 0`, review/checkpoint ordering confirmed by `git log --reverse` |
| QA-87 | PASS | Phase-Level Gate Results table; seven captured tokens, each terminated `exit_code: 0` |
| QA-88 | PASS | Requirement Compliance table above, including the explicit negative/unavailable rows and the two unresolvable citation families |
| QA-89 | PASS | Security review below |
| QA-90 | PASS | `go test -bench . -benchmem ./...` terminated `exit_code: 0` with `BenchmarkRebuildWorkspace100kEvents-6 ... 26941 ms/rebuild`, `BenchmarkList10kIssuesWithRegistry-6 ... 44.40 ms/registry-read`, `BenchmarkInitializeWorkspace-6 ... 62686633 ns/op`; per-sample `<30 s` `b.Fatalf` retained and satisfied. Deviation note: the three dedicated task-21 samples were 27804/27866/27906 ms; this session's three are 27018/27233/27448 ms, a 1.6-2.8% improvement, and the in-gate sample was 26941 ms. No regression against Phase 20 recorded values |
| QA-91 | **FAIL** | Documentation review below. `02_TRD.md`, `04_TRD_CONFIG.md` rules text, `07_DESIGN_PLAN_FLOWS.md`, `GLOSSARY.md`, and `AGENTS.md` phase-21 surfaces match the shipped interface; `03_TRD_SCHEMA.md` §8.2 `0014` and `04_TRD_CONFIG.md` §2 do not. See MF-1 and MF-2 |
| QA-92 | PASS | 22/22 committed reports pass a fresh two-argument `report-lint.sh` in this session (`PASS report-lint completed with no failures`, terminated `exit_code: 0`), which mechanically checks spec-command coverage, PASS/BLOCKED statuses, non-empty concerns and dependency-budget sentinels, commit-hash resolution, subject match, and committed-report cleanliness. Deviation notes for every failed-then-repaired attempt are present in the reports and preserved in the append-only ledger; no report records a dependency change |
| QA-93 | PASS | All four `v1-convergence-inventory.md` §7 recipes rerun fresh at phase HEAD with the DEC-158 status split. Role-policy recipe `exit 1 (no output)` (expected `exit_code: 1`); `--by` recipe `exit 1 (no output)` (expected `exit_code: 1`); broad table-reader recipe terminated `exit_code: 0` with 25 lines; `EventType` recipe terminated `exit_code: 0` with 21 lines. Every produced line maps to task-20's preserved per-row disposition: per-file match distribution is byte-identical to task-20 HEAD `2d34575` (`diff` of the per-file counts produced no output) and the Event-type literal value set is identical (21 values, `comm` empty in both directions). `UNCLASSIFIED=0`. Neither broad recipe was narrowed and neither produced a vacuous no-match |

### QA-89 security review detail

- **Canary scan.** `rg -n "CANARY-val-036|CANARY-live-secret|literal-val-036|literal-live-secret|redaction-canary" cmd/fid/testdata cmd/fid/integration/testdata` - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). No credential value or registered literal survives into any committed live-frame golden, Human fixture, or conformance corpus file. The Human fixture retains only the typed count `[control] redaction-applied replacements=2`.
- **Closed environment.** `internal/evidence/policy.go:47-48` admits exactly the ten documented `inherit` names `CI`, `GOCACHE`, `GOMODCACHE`, `GOPATH`, `HOME`, `LANG`, `LC_ALL`, `PATH`, `TMPDIR`, `XDG_CACHE_HOME`, matching `04_TRD_CONFIG.md` verbatim. `TestExecuteAppliesAcceptedGatePolicy` observes the assembled environment; green this session.
- **`credential-unavailable` before Attempt.** `TestExecuteRejectsMissingCredentialBeforeSpawn` green this session (`--- PASS`, `internal/evidence` `ok`, terminated `exit_code: 0`); `credential-unavailable.json` is present in the failure corpus with its category, channel, and exit.
- **Symlink and traversal rejection matrix.** `TestArtifactCandidatePathRejectsEveryUnsafeCandidate` and `cmd/fid/artifact_test.go:74` cover changed-during-copy, directory, missing, non-regular, outside-cwd, oversized, and symlink - the complete `artifact-rejected` reason enum in `public-interface.md` §8. `TestResolveRejectsSymlinkedFIDDirectory`, `TestReadProjectConfigRejectsArtifactRootSymlinkComponentEscape`, and `TestReadProjectConfigRejectsArtifactRootSymlinkEscape` cover workspace and Artifact-root escapes. All green in the `-race` gate.
- **`unsafe-path` details.** `TestRejectedArtifactsResolvePathsUnderExecutionCWD` green this session; the Human fixture's rejection notice renders only the declared relative path `missing-artifact.txt`, with no credential and no unrelated absolute path.
- **Credential values never reach journal, projection, error arguments, or export.** `internal/evidence/live_test.go:63` fails if either canary appears in any emitted frame; `public-interface.md` §8 requires replacement before framing, counting, excerpt, hash, Receipt, Artifact, error, export, or Human display. The committed goldens contain neither canary, and `TestFinalSchemaNamedReadsAndExportMatchAfterRebuild` compares exact export bytes across rebuild.
- **No containment claim.** VAL-046: both declared containment lexicons terminate `exit_code: 1`; accepted explicit values remain `sandbox:none`, `network:allowed`, `write:allowed`.

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: YES.
  QA-plan §6 records fresh validation at `463183a` and `b367a3b`; the Group 1
  checkpoint re-ran FEAS-01..FEAS-19 fresh at `5aeb19b`.
- Invalid or unvalidated command, test scope, fixture, artifact, or human
  checklist step found: **none**. Every command this review executed produced
  its declared status. Two QA-plan text items are superseded rather than
  invalid and are recorded for accuracy:
  - QA-plan §7.9 QA-85 states a pre-rebuild Event count of `20`. DEC-174 and the
    validation contract's amended `VAL-086` row fix the count at `22`, and
    `contracts.md` §16 enumerates the 22 families. The contract governs; the
    plan literal is stale.
  - QA-plan §5.1 prose says "the six task-20 absence searches" while listing
    five. Task 20's `## Expected Non-Zero Commands` declares five `rg` searches
    and five `ls` checks (ten total), matching §5.1's own "exactly the 10
    original declarations". All ten were rerun at their declared statuses.
  Both are QA-authored plan text, not contract or task changes.

## Mid-Phase Review Disposition

| Review | Outcome | Must-fix disposition | Fresh re-review evidence | Blocker candidate / existing ESC |
|---|---|---|---|---|
| review-1.md | REPAIR | MF-1 task-1 / VAL-001; MF-2 task-3 / VAL-010 | `review-1b.md` `VERDICT: PASS`, both RESOLVED from fresh inspection; repair commits `5539651`, `a169b41`, `401f609`, `68cc7ab` | none; ESC-21-005 resolved by DEC-163 before the repair sequence |
| review-1b.md | PASS | none | N/A | none |
| review-2.md | REPAIR | task-4 | `review-2b.md` `VERDICT: PASS`; repair commits `cd83fe9`, `5606254` | none |
| review-2b.md | PASS | none | N/A | none |
| review-3.md | REPAIR | task-9 | `review-3b.md` `VERDICT: PASS`; repair commit `98217f4` (stats role grouping) | none |
| review-3b.md | PASS | none | N/A | none |
| review-4.md | REPAIR | tasks 10, 11, 12 | `review-4b.md` `VERDICT: PASS`; repair commits `59b079a`, `8c763e9`, `0ced511` | none |
| review-4b.md | PASS | none | N/A | none |
| group-1-checkpoint-review.md (first) | REPAIR | MF-1 task-4 / VAL-005 Activation replay; MF-2 task-9 / VAL-030 exit class | Second independent review in the same artifact: `CHECKPOINT VERDICT: PASS`, VAL-042 PASS, VAL-043 PASS, FEAS-01..19 PASS at Group 1 HEAD `5aeb19b`; repair commits `0444ccc`, `1d89aa4` | none; DEC-168/DEC-169 governed the attempt and repair ownership |
| review-5.md | REPAIR | task-15 | `review-5b.md` `VERDICT: PASS` | none |
| review-5b.md | PASS | none | N/A | none |
| review-6.md | REPAIR | tasks 16, 17, 18 | `review-6b.md` returned REPAIR again for task-17; `review-6c.md` `VERDICT: PASS` closes the window | none |
| review-6b.md | REPAIR | task-17 | `review-6c.md` `VERDICT: PASS`; repair commits `5417b02`, `74e1efc` | none |
| review-6c.md | PASS | none | N/A | none |
| review-7.md | REPAIR | task-19 / VAL-070 | `review-7b.md` `VERDICT: PASS`, VAL-070 recorded PASS from `internal/export/export_test.go:97`; repair commits `3e27365`, `dfb7be6`, `5918b17`, `47109dc` | none |
| review-7b.md | PASS | none | N/A | none |

No mid-phase `REPAIR` was closed by implementation evidence alone. Every
original Must-fix finding is resolved in a later, separately authored review
artifact by a fresh independent reviewer, and `git log --reverse` confirms each
repair landed before its closing review.

## Documentation Review

- Impacted docs updated: **NO** - two defects, MF-1 and MF-2.
- Verified matching the shipped interface:
  - `03_TRD_SCHEMA.md` `0013` up and down: line-identical to
    `migrations/0013_p21_operational_kernel.up.sql` (296 lines) and
    `.down.sql` (27 lines), including the DEC-177 `rejected_artifacts_json`
    column and its `json_valid`/`json_type` CHECK.
  - `02_TRD.md` §9 reference instance: all six identity commands match exactly.
  - `04_TRD_CONFIG.md` rules text: `timeout_ms:1..600000`,
    `output_bytes:1..10485760`, `excerpt_bytes:0..min(output_bytes,16384)`,
    `artifact_bytes:1..1073741824` match `internal/evidence/policy.go:71-80`
    exactly; the ten-name `inherit` allowlist matches
    `internal/evidence/policy.go:47-48` exactly.
  - `07_DESIGN_PLAN_FLOWS.md` §10: the project-authored-workflow and
    rebuild-repeat prose matches the shipped tagged journey.
  - `06_DESIGN_PLAN_CLI.md`: the Phase 21 amendment lists the accepted leaves
    including Receipt `get|list` and delegates exact leaf spelling to
    `public-interface.md`.
  - `GLOSSARY.md`: Work Envelope, Obligation, Receipt, Activation, Claim, Run,
    Artifact, Attempt, capability availability, terminal-success, Routine,
    public enum, and Gate Definition Registration are all defined.
  - `AGENTS.md`: `internal/evidence` and `internal/claims`, the packages this
    phase introduced or reshaped, are listed in the Architecture section.
- Specs needing Tech Lead wrap-up sync: `03_TRD_SCHEMA.md` §8.2 (MF-1),
  `04_TRD_CONFIG.md` §2 (MF-2), `qa/validation-contract.md` Source column
  (IMP-1), `06_DESIGN_PLAN_CLI.md` projection grouping (IMP-2),
  `AGENTS.md` Architecture package list (IMP-3). The phase decisions file
  already schedules a `Spec Sync Verdict` pass at wrap-up; these are its inputs.
- Completion reports present for all tasks: **YES** - 22 of 22, each committed
  and each passing a fresh two-argument `report-lint.sh` in this session.

## Human Checklist Readiness

Generated at `project-docs/docs/phases/phase-21/qa/human-checklist.md`.

- Contains only `(human)` items: YES.
- Uses one observable per row: YES - HUMAN-01a..01d, one per observable.
- Human-only `VAL-NNN` assertions mapped: YES - VAL-036 is the sole Human
  parent assertion and all four sub-IDs cite it.
- Product/perception/dogfooding risks to call out: the fixture is a
  non-success execution (`outcome=failed`, `exit=23`, `receipt=none`,
  `obligation=unsatisfied`). The Human should confirm that the absence of a
  Receipt and the unsatisfied Obligation read as an honest negative result and
  are not mistaken for a passing gate.

## Coverage Caveat

A deterministic PASS proves only the asserted mechanics, data states, command
outputs, and structural conformance checks. It does not prove unasserted
appearance, workflow feel, product-owner judgment, or dogfooding quality; those
remain unconfirmed until the Human checklist records outcomes. In particular,
VAL-036's four perception observables are unanswered at the time of this review.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | BLOCKED | Specs, TDD, 22 task reports, Ready QA plan, deterministic QA evidence, and the Human checklist are all present, but two canonical specs contradict the shipped interface (MF-1, MF-2) |
| Automated gates | DEFERRED | Format, lint, race test, coverage report, tagged integration, benchmarks, and build all pass. SAST is recorded `N/A` before external release; `govulncheck`, `go-licenses`, and `gitleaks` are named but not required by this phase's gate list |
| Held-out scenarios | DEFERRED | `AGENTS.md` records the harness location and command as `N/A` until accepted; no phase-21 task reads held-out material |
| Release readiness | DEFERRED | Threat model, pre-release checklist, changelog, and rollback path are not phase-21 scope; Phase 21 is a pre-baseline slice under ADR-0011 |

## Findings

### Must-fix

#### MF-1 - `03_TRD_SCHEMA.md` §8.2 states a migration `0014` that the phase did not ship

- Classification: Must-fix; **blocker candidate** (protected specification
  authority, outside every phase task's write scope).
- Affected assertion: VAL-002 (recorded `BLOCKED`); QA-91.
- Reproduction, from the repository root at phase HEAD:
  - `sed -n '/^`0014` up is exactly:/,/^`0014` down is/p' project-docs/specs/03_TRD_SCHEMA.md`
    versus `cat migrations/0014_p21_retire_v1_authority_execution.up.sql`.
  - Spec `0014` up contains five statements: a `DELETE FROM
    workflow_transition_requires ... WHERE rule IN (...) OR (params_json IS NOT
    NULL AND json_type(params_json, '$.roles') = 'array')`, then
    `DROP TABLE gate_runs`, `DROP TABLE roles`, `DROP TABLE claims`,
    `ALTER TABLE work_claims RENAME TO claims`.
  - Shipped `0014` up contains seven statements in a different order and adds
    `DROP INDEX work_claims_one_active_mutation` and
    `CREATE UNIQUE INDEX claims_one_active_mutation ON claims(issue_id) WHERE
    status = 'active'`, and its `DELETE` predicate is
    `... WHERE rule IN (...) AND workflow_id IN (SELECT id FROM workflows WHERE
    applies_to IN ('epic','task','subtask'))`.
  - Spec `0014` down omits the corresponding `DROP INDEX
    claims_one_active_mutation` and `CREATE UNIQUE INDEX
    work_claims_one_active_mutation` that the shipped down performs.
- Observable consequence, verified this session in a fresh `mktemp -d`
  workspace initialized with the built binary:
  `SELECT name, tbl_name FROM sqlite_master WHERE type='index' AND name LIKE
  '%claims%'` returned `claims_one_active_mutation` on `claims`. Implementing
  §8.2's `0014` literally would leave that index named
  `work_claims_one_active_mutation`, so the spec-exact and shipped final
  schemas differ in a named schema object.
- Scope note, established rather than assumed: on the shipped seed data the two
  `DELETE` predicates produce the **same** end state. After `fid init` at final
  schema, `SELECT workflow_id, rule, params_json FROM
  workflow_transition_requires` returned exactly
  `epic-default/all_child_tasks_done_or_cancelled` and
  `escalation-default/linked_to_decision`, both with `params_json` null and
  neither matching the three named rules. This finding is therefore a
  specification-accuracy defect with a schema-object-name consequence, not an
  observed data-loss defect.
- Contract conflict: DEC-154 adopts `03_TRD_SCHEMA.md` §8.2 as the accepted
  planning contract, and VAL-002 cites `TRD Schema §8.2` as its Source
  requirement with fail criterion "any extra/missing action fails". §8.2 also
  designates `v1-convergence-inventory.md` as the live removal authority, and
  the shipped migration conforms to that inventory (§§4, 6) and to task 20's
  accepted binding constraint "deletes the displaced role/judgment **Issue**
  rows". Two adopted authorities disagree about what `0014` is "exactly", so
  VAL-002 cannot be certified `PASS` against its own cited source.
- Why this is outside existing repair authority: no phase task lists
  `project-docs/specs/03_TRD_SCHEMA.md` under `Files To Create Or Modify`.
  Tasks 1, 20, and 21 list `03_TRD_SCHEMA.md §8.2` only under
  `Files To Read For Context`; task 22's access to the file is bounded by
  DEC-177 to the Receipt persistence detail and the Resume command spelling.
  Amending an accepted planning contract, or amending a task's scope to permit
  the edit, is Tech Lead specification authority.
- Requested action: reconcile §8.2's `0014` up and down blocks with the shipped
  migration and the convergence inventory, or record which authority governs
  "exactly" for VAL-002 and annotate §8.2 accordingly.

#### MF-2 - `04_TRD_CONFIG.md` §2's documented `config.toml` is rejected by the shipped binary

- Classification: Must-fix; bounded and decision-free inside task 13.
- Affected assertion: no VAL; QA-91.
- Provenance: the two offending lines are phase-21 additions.
  `git diff main..HEAD -- project-docs/specs/04_TRD_CONFIG.md` shows
  `+max_timeout_ms = 600000` and `+max_artifact_bytes = 1073741824` added to the
  §2 example under `[executor_limits]`.
- Reproduction, each case in its own fresh `mktemp -d` workspace initialized
  with `/tmp/fid init --prefix FID`, then one config edit, then
  `/tmp/fid --json issue create --kind task --title T --actor agent:dev
  --idempotency-key k1`:
  - baseline generated config - terminated `exit_code: 0`;
  - add only `max_timeout_ms = 600000` - terminated `exit_code: 2`;
  - add only `max_artifact_bytes = 1073741824` - terminated `exit_code: 2`.
  The rejection matches an intentionally bogus key
  (`bogus_unknown_key = 1` - terminated `exit_code: 2`), confirming both
  documented keys are treated as unknown.
- Cause: `internal/config/project_config.go:141` admits exactly
  `max_output_bytes` and `max_excerpt_bytes` under `[executor_limits]`, and
  `project_config.go:111-113` writes only those two on `init`. The same
  document's own rules paragraph names only those two as configured executor
  ceilings; `timeout_ms` and `artifact_bytes` are per-Gate-Definition values
  validated in `internal/evidence/policy.go:71-80`, not config keys.
- Why the repair is decision-free: the document already states the correct
  rule; only its illustrative example contradicts it and the implementation.
  Removing the two lines from the §2 example restores internal consistency
  without changing any accepted bound, interface, or behavior.
- Owning task: **task 13**, whose `Files To Create Or Modify` list contains both
  `internal/config/` and `project-docs/specs/04_TRD_CONFIG.md`.

### Improvement

- **IMP-1 - unresolvable requirement citations in the validation contract.**
  `VAL-030` cites `PRD-AC-002` and `VAL-032`/`VAL-033`/`VAL-034` cite
  `PRD-CLM-001/002`, `PRD-CLM-003/004`, and `PRD-CLM-005`. Neither `PRD-AC-*`
  nor `PRD-CLM-*` exists anywhere in `project-docs/specs/`
  (`rg -n "PRD-CLM" project-docs/specs/` and `rg -n "PRD-AC-" project-docs/specs/`
  each terminated `exit_code: 1`); both are legacy identifiers still present only
  in phase-7 and phase-17 archives. The live equivalents are `PRD-RUN-001`
  (Claim binding), `PRD-RUN-002` (one active mutating Claim), `PRD-RUN-004`
  (versioned idempotent Claim lifecycle), and `PRD-ISSUE-007` with
  `PRD-HUM-001` for acceptance-criteria judgment. All four assertions'
  observables are satisfied; only the traceability column is stale. The
  validation contract is Tech-Lead-owned; QA proposes the remap without editing
  it.
- **IMP-2 - `06_DESIGN_PLAN_CLI.md` target grouping still spells the resume
  leaf as a `projection` subcommand.** Line 142 lists
  `fid projection context|resume|changes|completion|human-brief`, while the
  shipped leaf is root `fid resume <issue-id>`
  (`cmd/fid/resume.go:17`, `Use: "resume <issue-id>"`, canonical operation
  `projection.resume`), which DEC-177 confirms is "already-authoritative".
  Authority is not ambiguous - the same file's phase-21 amendment delegates
  exact leaf spelling to `public-interface.md`, which specifies
  `fid resume <issue-id>` - so this is presentation drift in the forward-looking
  target grouping, not a shipped-surface defect.
- **IMP-3 - `AGENTS.md` Architecture package list is out of date, pre-existing.**
  It omits shipped `internal/bootstrap`, `internal/config`, `internal/labels`,
  `internal/links`, `internal/requirements`, and `internal/sqlquery`, and lists
  seven packages that do not exist yet. This is inherited, not phase-21 drift:
  `git log --diff-filter=A` dates every omitted package to 2026-06-06 through
  2026-07-05, all before the phase-21 merge base `f8a720b` dated 2026-08-22.
  Phase 21's own new package `internal/evidence` **is** listed.
- **IMP-4 - low package-level coverage in two Phase 21 packages.**
  `internal/evidence` reports 24.6% and `internal/workflow` 27.4% of statements
  in the phase gate, against a 61.2% repository total. Behavior is covered
  end-to-end through `cmd/fid` (77.3%) and the tagged journey, and
  `08_QUALITY_PLAN.md` §7 treats coverage as diagnostic with no threshold, so
  this is informational only.

## Prior Finding Resolution

Not applicable; this is the initial final review for phase 21.

## Repair Plan

MF-2 has bounded, decision-free ownership in an existing phase task and is
recorded here so a single corrective pass can close it alongside the MF-1
authority decision. MF-1 itself is **not** repairable under this plan; it is a
blocker candidate below.

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-2 | task-13 | no VAL | `project-docs/specs/04_TRD_CONFIG.md` §2 example block only: remove the `max_timeout_ms` and `max_artifact_bytes` lines from `[executor_limits]`. Change no rules text, no bound, no `internal/config` code, no `internal/evidence/policy.go` bound, and no other spec | The documented §2 example, written verbatim into a fresh `mktemp -d` workspace's `.fid/config.toml`, followed by `fid --json issue create --kind task --title T --actor agent:dev --idempotency-key k1` terminating `exit_code: 0`; plus task 13's exact required command block rerun with captured terminations, a committed report, and a passing two-argument `report-lint.sh` | Re-run the documented-example acceptance reproduction above and its bogus-key control; confirm the §2 example and the §2 rules paragraph now name the same two `[executor_limits]` keys; confirm `internal/config/project_config.go:141` is unchanged |

## Blocker Candidates

- **BC-F1 - contradictory accepted specification for migration `0014`.**
  - Protected category: **specification authority** (accepted planning contract
    adopted by DEC-154; cited Source requirement of an accepted validation
    assertion).
  - Evidence: MF-1 above - the `03_TRD_SCHEMA.md` §8.2 `0014` up/down blocks
    differ from `migrations/0014_p21_retire_v1_authority_execution.up.sql` and
    `.down.sql` in statement set, statement order, and `DELETE` predicate scope;
    the observed final schema index is `claims_one_active_mutation`, which
    §8.2's SQL cannot produce.
  - Authority boundary: no phase task lists `project-docs/specs/03_TRD_SCHEMA.md`
    in `Files To Create Or Modify`; tasks 1, 20, and 21 hold it only as read
    context, and task 22's access is bounded by DEC-177 to unrelated content.
    Editing an adopted planning contract, or widening a task's scope to permit
    it, exceeds Developer and QA authority.
  - Requested action: the full Tech Lead reconciles §8.2's `0014` up and down
    blocks with the shipped migration and `v1-convergence-inventory.md` §§4/6,
    or records which authority governs "exactly" for VAL-002 and annotates
    §8.2. Then VAL-002 can be re-decided by a fresh independent QA Lead.
  - Existing escalation reference: none. `tl_escalations_current.md` records no
    open escalation. QA does not allocate an `ESC` identity.

## Unresolved Risks And Carry-Forward Candidates

- **VAL-036 is unanswered.** The Human checklist is generated and presented but
  not executed. The phase cannot record VAL-036 as `PASS` until HUMAN-01a..01d
  are answered by the Human.
- **RISK-01 headroom.** The 100k rebuild now runs at 27.0-27.4 s against a 30 s
  budget, about 2.6-3.0 s of headroom on the reference instance. Phase 22 adds
  recovery and import work to the same replay path. Carry-forward candidate:
  keep the three-consecutive-fresh-sample rule and watch the margin.
- **Coverage attribution.** IMP-4 above; `internal/evidence` and
  `internal/workflow` are thinly covered in-package and rely on `cmd/fid` and
  the tagged journey. Worth a targeted unit-coverage pass if either package
  grows in Phase 22.
- **Config unknown-key rejection is path-dependent.** A `.fid/config.toml`
  containing an unknown key is rejected on config-reading mutation paths
  (`issue create`, exit 2) but not on `issue list` or `capabilities`, which do
  not read the project config. `04_TRD_CONFIG.md`'s "Unknown keys ... reject"
  rule and the `isKnownKey` implementation both date to 2026-08-09, before the
  phase-21 merge base, so this is inherited behavior and out of scope for this
  phase. Recorded so a later phase can decide whether uniform validation is
  wanted.

## Final Decision

```text
Must-fix: 2  Improvement: 4  Blocker candidates: 1
VERDICT: BLOCK
```
