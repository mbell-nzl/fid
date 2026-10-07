# Phase 24 - Mid-Phase Review 8 (tasks 9, 9a)

Reviewer: qa-reviewer
Date: 2026-10-07
Tasks covered: phase-24/task-9, phase-24/task-9a
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: project-docs/docs/phases/phase-24/qa/review-7.md

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-24-001, ESC-24-002, ESC-24-003 (resolved history; none open)
- Covered validation assertions: VAL-011, VAL-012, VAL-013, VAL-014, VAL-015, VAL-016, VAL-035, VAL-052, VAL-060, VAL-064, VAL-065; VAL-036 advanced only, with completion allocated to D2

**Window D1 PASS. Review-7 F1 is RESOLVED by fresh independent evidence.** Both Activation prepare identities now accept their Unicode scalar boundaries. The independent 16-case Human/JSON matrix accepts and persists all twelve qualifying inputs and refuses all four Unicode129 inputs with exact metadata and no effects. The repair preserves ordinary byte bounds and component variant behavior. Both reports independently pass report-lint, retain prior failed evidence and contain every required exact command.

There are zero remaining Must-fix findings and zero blocker candidates. This closes review-7's bounded D1 repair requirement. The Phase Runner owns advancement. D2 after task 11, the separate Human walkthrough after task 11, and E after task 14 remain mandatory; this review supplies neither Human judgment nor phase-end acceptance.

## Scope Reviewed

Reviewed branch `phase-24` in `/home/mark/projects/fid` at `0dbb3e45f5a993efbd4bd336c7e816f76c8e6ddd`. The only initial worktree dirt was the Runner-owned `project-docs/docs/phases/phase-24/run-ledger.md`. This reviewer does not edit or stage it and writes and commits only this review artifact. No product code, tests, reports, task specs, decisions, current documents, plans, contracts or tools were changed by this reviewer.

Loaded AGENTS.md, the qa-evidence skill, workflow sections 2.4/2.7/8, the quality plan and review template, all four requested current context files, phase decisions including DEC-284, validation contract, task-9/task-9a specs and reports, QA-plan evidence/gates and section 8, review-7 in full and preceding review-6. Inspected relevant Phase 21/22/24 public-interface sections and the input inventory. CLI observations use `05_DESIGN_PLAN.md` sections 2-3 and 6-7 and the delivered amendments in `06_DESIGN_PLAN_CLI.md`. No held-out scenario material was read.

| Checkpoint | Independently inspected identity and scope |
|---|---|
| Prior D1 review | `a5e5b79e8133f5550edc773a4c1797a6a3559c79`, review-7 REPAIR; F1 and its six-path Repair Plan are the governing review boundary. |
| Task-9 repair | `29c8682c98c612443ed85df503d818dacc70552c`, `fix(cli): honor activation prepare identity scalar bounds`, `Refs: phase-24/task-9`. Full source/test/report diff inspected. |
| Report finalization | `0dbb3e45f5a993efbd4bd336c7e816f76c8e6ddd`, `docs(task): finalize completion report for phase-24 task-9`, `Refs: phase-24/task-9`. Report-only finalization inspected. |
| Task-9a retained slice | `2d5e9d55c83086b46268ca69793f30f65031b91b`, `feat(cli): converge grammar failures and unknown commands`; finalized report at `a8ae5b3c86a80fff7c245fdf72a0cbf9a920be97`. Parser/domain/message changes and relevant behavioral tests inspected; neither slice nor report changes in the F1 repair. |
| Preceding Window C | `149f494a9eee202914c6784ed3cbfdb7ce4c01c4`, review-6 PASS. Its task-7/task-8 resolutions remain preceding-window evidence, not newly claimed full Window C verification. |

The aggregate diff `a5e5b79e..0dbb3e45` changes exactly these six allowed paths:

```text
cmd/fid/input_descriptors_activation.go
cmd/fid/input_descriptor_test.go
internal/cli/input_descriptor.go
internal/cli/input_validate.go
internal/cli/input_validate_test.go
project-docs/docs/phases/phase-24/tasks/task-9-done.md
```

Each is a regular `100644` blob with normal tracked index flags. There are no changes to the component renderer, domain inputs, accepted interface, dependencies, tools, task/QA contracts, decisions or ledger in these commits. `git diff --check a5e5b79e..HEAD` terminates exit 0 without output. Individual `git rev-parse --verify <hash>^{commit}` calls resolve all 28 distinct full report/provenance hashes and the additional review/finalization checkpoints, each with observed exit 0. Hash existence is supplemented by inspecting the repair and affected source/test diffs.

## Validation Contract Coverage

ON TRACK below concerns the completed D1 allocation. Positive and negative controls were inspected and freshly executed; later-window assertions retain their owners.

| VAL ID | Claimed by | Status | Fresh evidence |
|---|---|---|---|
| VAL-011 | 9 | ON TRACK | O1 registry/live-tree/reflection equality across all 39 input leaves, seeded mismatch rejection and actual selected ActorProfile/GateDefinition trees. No pending fence remains. |
| VAL-012 | 9 | ON TRACK | O1 whole-tree field/help conformance, component full-help pin, selected variant metadata and expectation Git commands. The repair retains documented 1..128 identity bounds. |
| VAL-013 | 9 | ON TRACK | O1 real-process templates for all 39 leaves retain one required-only object, documented key order and two-space indentation. |
| VAL-014 | 9 | ON TRACK | O1 strict decoding into actual request types and descriptor JSON-type checks, including selected raw arrays; seeded reflection drift rejects. |
| VAL-015 | 9 | ON TRACK | O1 standalone examples succeed; 234 ordered forbidden combinations refuse before project access with the specified details and no effects. |
| VAL-016 | 9 | ON TRACK; F1 resolved | O1 recursive structural and variant corpora plus B1's independent prepare boundary matrix. Qualifying Unicode127/128 succeeds; Unicode129 refuses at the exact identity path without effects. Ordinary byte/scalar and supplementary-plane controls pass. |
| VAL-035 | 9 | ON TRACK | I1 derives the fake/fid dispatch intersection and exercises all nine accepted-input/removed-member pairs, retained Run/Gate parity and dual-adapter recovery. No one-sided dispatch row. |
| VAL-052 | 9, 9a | ON TRACK | Inspected complete exact-gate receipts and all four required package timings for both tasks; below budgets and warning levels. Fresh focused timings are separate evidence below. |
| VAL-060 | 9 | ON TRACK | O1 all standalone examples outside a project exit 0 with empty stderr and unchanged directories; exclusivity controls retain no-effects assertions. |
| VAL-064 | 9a | ON TRACK | O1 complete class/domain/typed-flag and built-binary corpora in both modes, project/non-project states, valid/malformed controls, short/long/512/513/non-ASCII/escape-cut boundaries. O2 retained grammar refusal and state-integrity subjects pass. |
| VAL-065 | 9a | ON TRACK | O1 unknown root/group words and near misses return operation-unavailable, correct dotted receiving paths, sorted allowed words and exit 2. No-word help succeeds; full details and bounded paired messages remain. |
| VAL-036 | 9a advances | DEFERRED completion | O1 current Human `error:` payload agrees with JSON grammar messages. Task 10 owns additional reason/blocked-by/next lines and the complete pairing corpus, reviewed again at D2. |

## Gate And Evidence Checks

### Report-to-spec comparison and preservation

Compared each report's `## Tests run`, `## Exact quality-gate results`, timing records and deviations with its task's `## Exact Commands To Run` and applicable QA-plan sections 4-5/8. Each contains all six required exact commands in a final passing sequence with output tokens and observed termination. No missing, silently skipped or narrower substituted final command was found. Both contain:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

Task 9's sixth command is `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid`, with predeclared expected exit 1 and final recorded/observed no-match exit 1. Task 9a's sixth command is `rg -n '"argument": *"request"' cmd/fid -g '!*_test.go'`, with final exit 0 and exactly four permitted family fallback literals. Task 9a additionally records the unchanged inventory section-1 block and R17/R18 counts. No extra in-window full-gate command from the QA plan is omitted; benchmark, full-checkptr, archive-scale and Human-generator controls have separate later allocations.

| Report | Final Developer evidence inspected; not independently rerun full gates | Immutable evidence assessment |
|---|---|---|
| task-9-done.md | Recovery exact sequence: silent format exit 0, lint `0 issues.`, ordinary workflow `103.615s`, coverage `62.3%`, tagged `318.051s`, R5 expected/observed exit 1. | All 20 failed/blocked receipt lines present at review-7 HEAD remain byte-for-byte; the five original failed/blocked lines at `f680b20` also remain. New F1 RED `FAIL`/exit 1 and byte-identical GREEN command, causal deviation and six fresh required-gate receipts are added. |
| task-9a-done.md | Final exact sequence: silent format exit 0, lint `0 issues.`, ordinary workflow `115.306s`, coverage `62.3%`, tagged `340.634s`, four R19 matches/exit 0. | Report is byte-identical to its review-7 version. All 23 failed/blocked/no-tests receipt lines from implementation `2d5e9d55` remain, including the first required tagged `FAIL`/exit 1 at `367.898s`, later policy.show failure and explicitly non-evidentiary empty selection. Causal convergence notes remain. |

The independent preservation/scope/command audit ends `PASS bounded scope, required-command coverage and immutable failure evidence` (terminated `exit_code: 0`). It checks individual commit hashes, exact repair path set, regular files/index flags, unchanged task-9a report, historical failed receipt equality and all required commands. The report-lint results below independently establish committed, clean COMPLETE reports. They do not replace the fresh behavior checks.

| Inspected final exact-gate timing | cmd/fid | internal/workflow | Slowest other ordinary package | Tagged integration |
|---|---:|---:|---|---:|
| Task 9, F1 repair | 418.624 s | 103.615 s | internal/issues 188.589 s | 318.051 s |
| Task 9a | 399.177 s | 115.306 s | internal/issues 196.455 s | 340.634 s |

All are below ordinary 600 s/tagged 900 s budgets and 480 s/720 s warnings. These are inspected Developer full-gate timings, not extrapolations from the focused QA runs.

### Fresh reproduced evidence

Working directory is `/home/mark/projects/fid`, except each fresh binary fixture directory. Full nested command result objects were preserved; live sessions were polled to their own observed exit_code. Wrapper completion, silence and absence of new output were never treated as termination. Each Go selection uses `-count=1`.

```text
O1:
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid ./internal/cli -run '^(TestGrammar.*|TestInput.*|TestBoundedMessage|TestReasonCodeDetailSchemaIsExactForEveryCode|TestAdoptedExampleProcessExclusivity|TestP24WholeTreeDescriptorProcessCorpus|TestP24DescriptorVariantProcessCorpus|TestP24Activation.*|TestP24InputDescriptorsHaveNoTopLevelSensitiveFields|TestP24BindingDescriptorProcessCorpus|TestP24RunDescriptorProcessCorpus|TestP24JudgmentStructuralValidationCorpus|TestP24ReviewAndFindingDescriptorHelpAndStructure|TestDecodeClosedReviewInputRejectsMalformedInapplicability|TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState|TestLiveLeafCommandsMatchContractInventories)$'

I1:
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid/integration -run '^(TestP24ChangedFakeFIDAssembly|TestP24RunChangedFakeFIDInputParity|TestP24GateChangedFakeFIDInputParity|TestP22DualAdapterRecoveryConformance|TestProcessExitContract|TestV01ClosureJourneyAndRebuildParity|TestP23SharedAmendedLeafConformance|TestP23FailureCorpusHasEveryAcceptedCode|TestP22FakeReviewReworkStopCondition)$'

O2:
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid -run '^(TestInitCommandEnforcesExplicitPrefixGrammarBeforeStaging|TestIssueACAddCommandRejectsBeforeAppend|TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures|TestIssueTransitionRequiresTargetStateAndKnownByRole|TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged|TestLabelAddRejectsExplicitlyEmptyOptionalFlagsBeforeAppend|TestLinkJSONRequestAndTypedFailuresAreAtomic|TestPhase22AttentionReadBoundary|TestPolicyShowIsReadOnlyAndDoesNotAcceptByFlag|TestPromisedGovernanceCommandsAreAbsent|TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork|TestReceiptReadFailuresAreTypedAndClosed|TestRetainedInventoryJSONMissingActorsUseCanonicalEnvelope|TestRetainedInventoryJSONRequestErrorsUseCanonicalEnvelope|TestRetainedInventoryRequiredActorHumanErrorsRemainProse)$'

B1:
python3 /tmp/fid-review8-b7to7a1h/boundaries.py

P1 (future sensitive-fallback corpus feasibility):
python3 /tmp/fid-review8-b7to7a1h/fallback.py
```

| Fresh check | Result | Verbatim output token and observed termination |
|---|---|---|
| Representative exact required command `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`) |
| O1 affected ordinary/shared/whole-tree/grammar checks | PASS | cmd/fid `64.741s`; last `ok  \tgithub.com/mbell-nzl/fid/internal/cli\t1.171s` (terminated `exit_code: 0`) |
| I1 tagged real-process parity and retained journeys | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t22.495s` (terminated `exit_code: 0`) |
| O2 retained ordinary refusal/state-integrity assertions | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t4.907s` (terminated `exit_code: 0`) |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-9-done.md project-docs/docs/phases/phase-24/tasks/task-9.md` | PASS | `OUTCOME COMPLETE`; last `PASS report-lint completed with no failures` (terminated `exit_code: 0`) |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-9a-done.md project-docs/docs/phases/phase-24/tasks/task-9a.md` | PASS | `OUTCOME COMPLETE`; same last token (terminated `exit_code: 0`) |
| Exact inventory section-1 Bash block, unchanged on stdin | PASS | ` 271 total` (nested Bash observed exit 0; enclosing inventory audit terminated `exit_code: 0`) |
| R1-R20 regeneration and scope/preservation audits | PASS | `PASS R1-R20 inventory; scratch=/tmp/fid-review8-b7to7a1h; binary=/tmp/tmp.x7qO2y5IBe/fid`; preservation token above (each enclosing audit terminated `exit_code: 0`; nested statuses retained) |
| B1 independent boundary matrix | PASS | `PASS 16-case Human/JSON prepare scalar matrix: 12 accepted and persisted, 4 refused with exact metadata and no effects` (terminated `exit_code: 0`) |
| Sensitive generic-fallback feasibility | FEASIBLE | `FEASIBLE activation sensitive generic fallback in both modes without disclosure or effects` (terminated `exit_code: 0`) |

`\t` denotes captured tab separators. No required or focused Go check failed in this review. Full repository lint, coverage, full ordinary/tagged suites, benchmarks, scale controls and full-checkptr were not independently rerun: representative affected D1 checks are required here, while full phase checks belong to final QA. The full exact Developer gate records were independently compared with their specifications.

Two supplemental fallback-harness diagnostics are retained rather than relabelled PASS. The first Python-stdin probe terminated `exit_code: 1` with `AssertionError`: its first component call validly staged with exit 0, so it did not yet drive the intended fallback. The next `/tmp/fid-review8-b7to7a1h/fallback.py` run reached the correct refusal but terminated `exit_code: 1` because the harness expected `idempotency result` rather than observed `idempotency_results/insert`. The completed harness uses fresh directories, stages a qualifying component, then submits a second component with the same content digest and a new key/ID/current expected version; its Human assertion matches the observed spelling. No product/test/report/control was altered. These were supplemental fixture/assertion errors, not required-command failures or a change to any accepted gate.

## Task 9 - Activation descriptors; whole-tree input and help conformance

### Must-fix

- None. Review-7 F1 is RESOLVED; see Prior Finding Resolution.

### Improvement

- None.

Observation (design sections 2-3, 6-7): `cmd/fid/input_descriptors_activation.go:23` opts into scalar bounds through the small `activationPrepareIdentity` helper, called only for `kit_identity` and `backend_identity` at line 42. The old private marker becomes exported `InputField.ScalarBounds` at `internal/cli/input_descriptor.go:28`; `input_validate.go:154` still defaults to byte count and switches to `utf8.RuneCountInString` only when the marker is set. The existing recursive raw-array variant selection at `input_descriptor.go:164` performs the same operation as before. Object variants, ordinary field formats, reflection, renderers and domain handling are unchanged. This is conformance to Phase 21 public-interface section 2's scalar definition and section 2.2's identity bounds, retained by Phase 22 section 3's extended prepare input.

The permanent regression at `cmd/fid/input_descriptor_test.go:2544` exercises both identities and output modes in separately initialized projects, verifies identity persistence on success, full JSON envelope/Human prose on refusal and unchanged filesystem/database/journal snapshots. `internal/cli/input_validate_test.go:188` adds empty/ASCII/Unicode/astral scalar versus byte controls with exact refusal values/paths/bounds. Existing component scalar/path/unbounded-string, variant mismatch, sensitive redaction and ordinary-caller tests are retained and pass O1. Assertions test observable acceptance, streams, metadata and persistence rather than the marker's spelling.

B1 independently builds on the fresh inventory binary, SHA-256 `9f29053fb282fe9f62de49434d5d576810cbdd3afe74e29097cb1738a1da42c8`. For each field/boundary/mode it initializes a separate directory with `fid init --prefix QA`, then invokes `fid activation prepare --input - --idempotency-key qa-scalar`, adding root `--json` for JSON. Input is exactly review-7's request: default identities `kit`/`backend`, phase 21, direction activate, null prior ID/version, five sorted backend/kit/overlay/project/provider resources with null old identity and new identity `new`; replace only the tested identity.

| Replacement, independently for each identity | Expected exit in Human and JSON | Fresh observed exit in Human and JSON | Evidence |
|---|---:|---:|---|
| `a` repeated 128 times | 0 | 0 | ASCII upper boundary persisted; stdout nonempty, stderr empty. |
| `é` repeated 127 times | 0 | 0 | Qualifying scalar boundary persisted; stdout nonempty, stderr empty. |
| `é` repeated 128 times | 0 | 0 | Scalar upper boundary persisted; stdout nonempty, stderr empty. |
| `é` repeated 129 times | 2 | 2 | Empty stdout, exact field refusal and unchanged file/mode/content snapshots, including DB and journal. |

JSON refusal is the exact `fid-cli.v1` activation.prepare envelope: invalid-input/invalid-argument, message `invalid input.<field>`, retryable false, subject null, details `{argument:"input.<field>",value:<full identity>,allowed:["length 1..128"]}`. Human refusal is exactly `error: invalid input.<field>` plus LF. Success compares the JSON Activation identity to the submitted value and verifies the value in journal bytes in both modes. B1's receipt contains all argv, cwd, input, streams, observed subprocess exits and before/after snapshots; SHA-256 `fa0c7e33e0b334b940f4d9e30f0c486974fdf5e590b4a6d21f532b7441bb76da`. The commands, payload and matrix above permit reproduction without relying on scratch retention.

O1 also rechecks all 39 descriptors, strict templates, recursive structural refusals, both ActorProfile/GateDefinition variants and all ten Stop variants. Sensitive values remain suppressed and no top-level field is sensitive. I1 regenerates the nine-leaf fake/fid intersection and exercises acceptance and removed-member refusals on both adapters. The repair adds no new architecture, dependency, input shape or documentation requirement; accepted help/documentation already specifies these bounds.

## Task 9a - Grammar failures and unknown commands across the live tree

### Must-fix

- None.

### Improvement

- None.

Observation (design sections 2-3, 6-7): the F1 diff does not change parser/domain declarations, grammar renderers or shared message construction. Fresh O1 still exercises 90 stable leaves and 29 root/groups, class failures before project resolution, every typed flag declaration plus seeded missing-declaration rejection, declared string domains and real-process Human/JSON rows inside and outside projects. Native integer syntax, amended cancelled/superseded members, opaque Unicode states, accepted label clearing and no-word help controls remain. O2 and I1 preserve the retained refusal, journal-integrity and removed-command journey subjects reviewed in review-7.

`canonical_request.go` keeps grammar classification at the adapter boundary; semantic combinations remain outside the four grammar classes under CF-047. R19 retains exactly the permitted four family literals. `internal/cli/message.go` escapes controls then bounds to 512 scalars without splitting existing or newly constructed escapes. Short/long, exact 512/513, non-ASCII and escape-cut controls reproduce Human/JSON message equality while preserving full details. Documentation records delivered grammar behavior and leaves task-10 Human lines pending. The report correctly advances rather than certifies full VAL-036. No new business rule, public enum, dependency or architectural boundary is introduced by this repair.

## Repair Plan

- None. The sole original Must-fix has been resolved by fresh independent source/diff, behavior and report evidence. No additional repair dispatch is required by this review.

## Prior Finding Resolution

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| **Review-7 F1:** prepare kit_identity/backend_identity used byte lengths and rejected accepted Unicode127/128; VAL-016. | phase-24/task-9 | **RESOLVED** | Exact six-path repair at `29c8682c`/`0dbb3e45`; explicit scalar opt-in only for the two prepare identity fields, ordinary default and recursive component semantics preserved. Permanent behavioral boundary/byte/astral regressions inspected. O1/I1/O2 pass with observed exit 0; independent B1 passes all 16 Human/JSON cases with twelve persisted acceptances and four exact no-effects refusals. Both reports independently lint COMPLETE/exit 0. All six current task-9 exact gates present, original failures preserved, causal RED/GREEN/deviation documented. |

Review-7 remains unchanged and continues to record its eight incorrect qualifying-state rejections. This separate independent artifact closes F1; Developer receipts alone were not treated as QA verification. Review-6's preceding Window C acceptance does not substitute for these D1 checks.

## Cross-cutting observations

R1-R20 and the exact section-1 recursive Bash block were freshly regenerated from repository root. Markdown-escaped alternations were interpreted as actual regex alternation. R4's migration listing/schema search and R12's file listings were separate nested commands; R8 inspected `definition(` rows. All nested query/list commands exited 0 except R5, whose no-match exit 1 is the declared absence condition. Full outputs and statuses are retained in `/tmp/fid-review8-b7to7a1h/registries.json` and `inventory-block.json`.

| Registry | Fresh result and D1 assessment |
|---|---|
| R1/R2 | 45 digest call-site matches; one shared helper. Unchanged by repair. |
| R3/R4 | 126 idempotency matches; 38 migration files through 0019 and 17 schema references. Unchanged physical/projection contracts. |
| R5 | Zero decoder matches, observed exit 1. Descriptor-only structural source remains. |
| R6/R7 | Four repository helpers and 16 resolved-fact references; no caller-identity revival. |
| R8 | 85 definition/declaration lines, including operation-unavailable; unchanged code inventory. |
| R9/R10 | 43 failure functions; 30 Human error-assertion matches. R10 increases from review-7's 29 solely because the new permanent boundary test adds its Human refusal assertion. |
| R11/R12 | 19 repository-output/help/fixture references; all three Stop/Attention renderer paths exist. |
| R13/R14 | 15 judgment/PRIVATE references; live inventory completeness test present and passing O1. |
| R15/R16 | 92 leaves, zero missing Short; 50 flag-keyed leaves, 49 canonical keyed leaves excluding export import. |
| R17/R18 | 29 root/groups; 61 typed flags, 17 integers and 44 booleans including all 39 example flags. Corpus completeness checks and binary lists agree. |
| R19 | Exactly four literal family fallbacks: issue_get_list.go:129, issue_mutation_response.go:221, link_response.go:89, receipt.go:87. |
| R20 | 23 accounted invalid-argument sites. Combination-helper semantics remain explicitly retained under CF-047; task 10 owns fallback convergence and corpus completeness. |

The unchanged block generates 92 leaves, 29 groups, 39 input leaves, 50 keyed leaves, zero no-Short leaves and 61 typed flags (` 271 total`). Its fresh binary is `/tmp/tmp.x7qO2y5IBe/fid`. Corpus foundations and inventory commands are executable; the repair introduces no new leaf, flag, fallback or future fixture requirement.

The completed sensitive-fallback probe uses a separately initialized Activation per output mode, stages a valid GateDefinition with literal environment sentinel `QA-REVIEW8-SENSITIVE-SENTINEL`, then submits matching content under a second component identity/key and the current expected version. The duplicate content digest reaches `idempotency_results/insert` validation, not a structural or typed descriptor refusal. Both calls exit 2 with empty stdout, unchanged file-content snapshots and no sentinel in either stream. JSON retains the current fixed `invalid activation input` message and `{argument:"input",value:null,allowed:["accepted activation input"]}`; Human carries `array must be sorted and unique`. This is fresh feasibility for the task-10 guard, not final VAL-067 message acceptance or new repair ownership for this existing domain cause.

Later corpus feasibility queries remain executable: task-13/14's exact disclaimer query currently returns three matches/exit 0, and `rg -n "PRIVATE_" cmd/fid/integration` returns twelve matches/exit 0. These are expected present-day inputs to later removal work, not passing absence gates; their final predeclared success remains exit 1. `rg -n '^func TestP24HumanFixtures\(' cmd/fid/integration` returns no output/exit 1: the task-11 generator remains pending, so its selector was not run or passed on a no-tests result. Preceding review-6/review-7 feasibility remains attributed historical context; current shared/grammar/real-process checks and inventories were reproduced here.

Current brief/progress snapshots retain their previously disclosed earlier recovery status. Verified commits, committed reports and the linked reviews establish the reviewed state. This review does not edit those snapshots or infer a changed cadence from them.

D2 remains after task 11 for tasks 10/11, including grammar/Human pairing after renderer adoption, generic fallbacks, Stop/Attention and repository lines. At task-11 close its QA Reviewer must regenerate fixtures with:

```text
FID_P24_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-24-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration
```

The separate Human walkthrough remains `project-docs/docs/phases/phase-24/qa/human-walkthrough-after-task-11.md`; interim VAL-054/055 does not replace phase-end HUMAN-02/HUMAN-03. E remains after task 14 for tasks 12-14. Complete phase-end deterministic gates, validation outcomes and Human judgments remain required.

## Blocker Candidates

- None. No protected authority, out-of-window repair or exhausted recovery budget is required. Zero blocker candidates remain.

## Existing Escalation References

- ESC-24-001 and ESC-24-002: resolved history under DEC-283.
- ESC-24-003: resolved shared-descriptor/orchestration history under DEC-284. Its shared authority and review-7's narrower six-path repair boundary were preserved.
- None open; no escalation identity was allocated or register modified.

## Final Lines

Review written: project-docs/docs/phases/phase-24/qa/review-8.md
Must-fix: 0  Improvement: 0  Blocker candidates: 0
VERDICT: PASS
