# Phase 24 - Mid-Phase Review 7 (tasks 9, 9a)

Reviewer: qa-reviewer
Date: 2026-10-07
Tasks covered: phase-24/task-9, phase-24/task-9a
Review type: INITIAL
Prior review: none for Window D1; preceding Window C reviews are project-docs/docs/phases/phase-24/qa/review-5.md and review-6.md

## Summary

- Review outcome: REPAIR
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 0
- Existing escalation references: ESC-24-001, ESC-24-002, ESC-24-003 (resolved history; none open)
- Covered validation assertions: VAL-011, VAL-012, VAL-013, VAL-014, VAL-015, VAL-016, VAL-035, VAL-052, VAL-060, VAL-064, VAL-065; VAL-036 advanced only, with completion allocated to D2

**Window D1 requires bounded task-9 repair.** The new Activation prepare descriptor rejects accepted Unicode identities because it applies byte lengths to fields whose accepted limits count Unicode scalar values. Fresh binary probes reject 127- and 128-scalar values in both identity fields, in both output modes. This is a descriptor conformance defect within task 9, not a request to change the accepted interface or shared default length policy.

Both reports independently pass report-lint. Fresh ordinary and real-process selections pass, including whole-tree descriptors, variants, redaction, grammar/domain/unknown-command corpora, fake/fid assembly and retained closure assertions. Those passing tests do not cover the failing prepare boundary and do not close F1. The window must not advance until repair and a fresh independent D1 re-review. There are no blocker candidates. D2 after task 11, its Human walkthrough, and E after task 14 remain mandatory.

## Scope Reviewed

Reviewed branch `phase-24` at `a8ae5b3c86a80fff7c245fdf72a0cbf9a920be97` in `/home/mark/projects/fid`. The only initial worktree dirt was the Runner-owned `project-docs/docs/phases/phase-24/run-ledger.md`; it was not edited. This reviewer writes and commits only this review.

Loaded AGENTS.md; the qa-evidence skill; workflow sections 2.4, 2.7 and 8; quality plan; review template; the four current context files; phase decisions including DEC-284; task-9/task-9a specs and reports; validation contract; QA-plan evidence, gates and section 8; preceding reviews 5/6 and relevant review-4 context; input inventory; and relevant Phase 21-24 public-interface amendments. CLI observations use the applicable design-plan input/help/error guidance. No held-out scenario material was read.

Independently inspected commit contents, diffs, tests and documentation rather than accepting Developer receipts as reproduced behavior. Main source/report checkpoints:

| Task | Commits inspected | Scope assessment |
|---|---|---|
| 9 | `785a22491f1ba07bb748d1d4fb85f88ea1a8d45d`, `9d654e976792836cbdf7de8c114d0831b73f721e`, `69589747ca77a1a3dae7e4d6083229d0193ed52f`; report `d6e03c9f9ec7a445f3af50451cb9b14d4d52dd5c`, close `e39bfcc11237e108cebb58d46aa5dc696ab44ca1` | Activation handlers/descriptors, descriptor tests, retained shared-reader tests, fake/fid corpus, CLI documentation and report. Shared extension is exactly DEC-284's `internal/cli/input_descriptor.go`, `input_render.go`, `input_validate.go`, `input_validate_test.go`. |
| 9a | Implementation `2d5e9d55c83086b46268ca69793f30f65031b91b`; final report `a8ae5b3c86a80fff7c245fdf72a0cbf9a920be97` | Central parser/grammar declarations, constructors, shared failure code/message support, grammar and retained tests, tagged fixtures, CLI documentation and report. |

The task commits do not modify tools, task contracts, decisions, the ledger, dependencies or unrelated business-rule packages. Checked descriptor artifacts are regular `100644` blobs with normal tracked index flags. Main implementation subjects/Refs and relevant introducing test commits were inspected. All 27 distinct report/provenance hashes checked through individual `git rev-parse --verify <hash>^{commit}` calls resolve with observed exit 0. This verifies commit existence; retained assertion diffs were also read to establish what actually changed.

## Validation Contract Coverage

ON TRACK below refers to the D1 allocation, not final phase acceptance. Qualifying and non-qualifying cases were inspected and rerun. Structural negatives retain exact paths and no effects, rather than relying only on valid templates.

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-011 | 9 | ON TRACK | O1 registry/live-tree/type equality across all 39 input leaves, seeded mismatch rejection, actual ActorProfile/GateDefinition variant reflection and all Stop variants. No remaining adoption fence. |
| VAL-012 | 9 | ON TRACK | O1 per-leaf help corpus, complete component help pin, recursive field metadata and template commands, including the required expectation Git commands. Boundary implementation mismatch is F1. |
| VAL-013 | 9 | ON TRACK | O1 built-process examples are required-only ordered objects with two-space indentation across all 39 input leaves. |
| VAL-014 | 9 | ON TRACK | O1 strict request-type decoding, descriptor JSON type checks and selected variant templates; raw JSON array fields do not bypass reflection. |
| VAL-015 | 9 | ON TRACK | O1 standalone examples succeed; all six forbidden combinations per leaf reject before project access, with exact metadata and unchanged directories. |
| VAL-016 | 9 | NOT ON TRACK | O1 whole-tree/variant structural corpora pass, including nested keys, duplicates, required/null/type/enum/format controls and sensitive values. Independent prepare identity boundary probe fails accepted scalar-length controls: F1. |
| VAL-035 | 9 | ON TRACK | I1 regenerated changed fake/fid dispatch intersection, accepted input and removed-member refusal on both adapters, retained Run/Gate and dual-adapter recovery conformance. No one-sided dispatch case. |
| VAL-052 | 9, 9a | ON TRACK | Both reports contain all required final exact commands, observed termination, coverage and package timings below budgets/warnings. Fresh bounded QA timings are separate evidence below. |
| VAL-060 | 9 | ON TRACK | O1 all standalone examples outside a project exit 0 with empty stderr and unchanged directory snapshots; exclusivity refusal controls remain. |
| VAL-064 | 9a | ON TRACK | O1 complete grammar classes/domains/typed-flag registry and built-binary corpus in both modes, project/non-project snapshots, inherited string domains and malformed controls. Short/long, 512/513-scalar, non-ASCII and escape-cut controls reproduced. |
| VAL-065 | 9a | ON TRACK | O1 unknown root/group words and near misses return operation-unavailable with the owning dotted path, sorted allowed words and exit 2; no-word help succeeds. Both modes and bounded message controls reproduced. |
| VAL-036 | 9a advances | DEFERRED completion | O1 current Human `error:` and JSON grammar messages agree. Task 10's additional failure/next lines and pairing corpus remain D2; this review does not certify full VAL-036. |

## Gate And Evidence Checks

### Exact report-to-spec comparison

Compared each report's `## Tests run`, `## Exact quality-gate results`, timing tables and deviations against its task's `## Exact Commands To Run` and the applicable QA-plan commands. Both contain the five common commands below plus their distinct sixth inventory command. No missing, silently skipped or narrower substituted final required command was found.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...

Task 9 sixth command:
rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid

Task 9a sixth command:
rg -n '"argument": *"request"' cmd/fid -g '!*_test.go'
```

Task 9 predeclares expected exit 1 for its absence command; its final observed exit 1/no output is correctly recorded as PASS. Task 9a's command returns exit 0 and exactly four permitted family fallback map literals. The declarations and observed statuses agree.

| Report | Final Developer evidence inspected, not a fresh QA rerun | Failed evidence/provenance assessment |
|---|---|---|
| task-9-done.md, exact gates at line 190 | Format silent exit 0, lint `0 issues.`, ordinary workflow `102.835s`, coverage `61.8%`, tagged `278.829s`, declared R5 absence exit 1 | Earlier ordinary, R5 and report-lint failures remain. DEC-284 shared extension and component scalar/path correction are recorded. Comparing the earlier report blob at `f680b20` preserves all five original failed/blocked receipt lines byte-for-byte. |
| task-9a-done.md, final exact gates following line 328 | Format silent exit 0, lint `0 issues.`, ordinary workflow `115.306s`, coverage `62.3%`, tagged `340.634s`, four R19 matches exit 0 | First required tagged failure (`FAIL`, exit 1, `367.898s`) and subsequent focused policy.show failure remain with causal convergence notes. Comparing the earlier report at `2d5e9d55` preserves all 23 failed/blocked/no-tests receipt lines byte-for-byte. The no-tests selector is explicitly not evidence. |

The independent preservation/path audit terminated exit 0 with `PASS immutable failed receipts and task commit path inventories inspected`. Full final attempts are distinguishable from earlier failed, incomplete and supplemental selections. Passing report-lint does not establish behavioral acceptance of F1.

| Final report timing | cmd/fid | internal/workflow | Slowest other ordinary package | Tagged integration |
|---|---:|---:|---:|---:|
| Task 9 | 410.123 s | 102.835 s | internal/issues 199.339 s | 278.829 s |
| Task 9a | 399.177 s | 115.306 s | internal/issues 196.455 s | 340.634 s |

All are below ordinary 600 s/tagged 900 s budgets and 480 s/720 s warning levels. These are the reports' full-gate times, not timings extrapolated from focused QA selections.

### Fresh reproduced evidence

Working directory: `/home/mark/projects/fid`, except each binary scenario's fresh scratch directory. Full nested command result objects were preserved. Yielded nested sessions were polled until their own exit_code arrived; silence and wrapper completion were not treated as termination. All focused Go selections use `-count=1`.

```text
O1:
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid ./internal/cli -run '^(TestGrammar.*|TestInput.*|TestBoundedMessage|TestReasonCodeDetailSchemaIsExactForEveryCode|TestAdoptedExampleProcessExclusivity|TestP24WholeTreeDescriptorProcessCorpus|TestP24DescriptorVariantProcessCorpus|TestP24Activation.*|TestP24InputDescriptorsHaveNoTopLevelSensitiveFields|TestP24BindingDescriptorProcessCorpus|TestP24RunDescriptorProcessCorpus|TestP24JudgmentStructuralValidationCorpus|TestP24ReviewAndFindingDescriptorHelpAndStructure|TestDecodeClosedReviewInputRejectsMalformedInapplicability|TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState|TestLiveLeafCommandsMatchContractInventories)$'

I1:
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid/integration -run '^(TestP24ChangedFakeFIDAssembly|TestP24RunChangedFakeFIDInputParity|TestP24GateChangedFakeFIDInputParity|TestP22DualAdapterRecoveryConformance|TestProcessExitContract|TestV01ClosureJourneyAndRebuildParity|TestP23SharedAmendedLeafConformance|TestP23FailureCorpusHasEveryAcceptedCode)$'

O2, remaining named retained ordinary assertions:
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid -run '^(TestInitCommandEnforcesExplicitPrefixGrammarBeforeStaging|TestIssueACAddCommandRejectsBeforeAppend|TestIssueCreateAndUpdateRejectRemovedOwnerFlagsWithCanonicalFailures|TestIssueTransitionRequiresTargetStateAndKnownByRole|TestIssueUpdateRejectsInvalidInputAndLeavesStateUnchanged|TestLabelAddRejectsExplicitlyEmptyOptionalFlagsBeforeAppend|TestLinkJSONRequestAndTypedFailuresAreAtomic|TestPhase22AttentionReadBoundary|TestPolicyShowIsReadOnlyAndDoesNotAcceptByFlag|TestPromisedGovernanceCommandsAreAbsent|TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork|TestReceiptReadFailuresAreTypedAndClosed|TestRetainedInventoryJSONMissingActorsUseCanonicalEnvelope|TestRetainedInventoryJSONRequestErrorsUseCanonicalEnvelope|TestRetainedInventoryRequiredActorHumanErrorsRemainProse)$'

I2, remaining named retained fake recovery assertion:
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 ./cmd/fid/integration -run '^TestP22FakeReviewReworkStopCondition$'
```

| Fresh check | Result | Output token and observed termination |
|---|---|---|
| Representative exact required command `gofmt -l cmd internal` | PASS | `exit 0 (no output)`; terminated `exit_code: 0` |
| O1 ordinary/shared/whole-tree/grammar corpus | PASS | cmd/fid `65.139s`; final `ok  \tgithub.com/mbell-nzl/fid/internal/cli\t1.155s`; terminated `exit_code: 0` |
| I1 real-process parity and retained journeys | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t13.113s`; terminated `exit_code: 0` |
| O2 retained ordinary assertions | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid\t5.113s`; terminated `exit_code: 0` |
| I2 retained fake recovery assertion | PASS | `ok  \tgithub.com/mbell-nzl/fid/cmd/fid/integration\t1.604s`; terminated `exit_code: 0` |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-9-done.md project-docs/docs/phases/phase-24/tasks/task-9.md` | PASS | `OUTCOME COMPLETE`; last `PASS report-lint completed with no failures`; terminated `exit_code: 0` |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-9a-done.md project-docs/docs/phases/phase-24/tasks/task-9a.md` | PASS | `OUTCOME COMPLETE`; same last token; terminated `exit_code: 0` |
| Current inventory section-1 Bash block, supplied unchanged on stdin | PASS | ` 271 total`; terminated `exit_code: 0`; built a fresh binary at `/tmp/tmp.BktXGSPuY2/fid` |
| `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` | PASS, expected absence | `exit 1 (no output)`; expected `exit_code: 1`, terminated `exit_code: 1` |
| R19 query shown above | PASS | Last match starts `cmd/fid/issue_get_list.go:129:` in the final fresh query; four matches; terminated `exit_code: 0` |
| Prepare scalar-boundary harness described under F1 | FAIL | `FAIL Unicode scalar boundaries: 8 accepted cases rejected`; terminated `exit_code: 1` |
| Fresh sensitive activation-component fallback feasibility probe | FEASIBLE | `FEASIBLE sensitive guard reaches activation generic fallback in both modes without effects or disclosure`; terminated `exit_code: 0` |

`\t` represents captured tab separators. The failed boundary attempt remains failed evidence. Full repository lint, coverage, full ordinary/tagged suites, benchmarks, archive-scale and full-checkptr controls were not independently rerun: D1 calls for representative in-window ordinary/real-process reruns, not phase-end QA. O1/I1 cover the changed shared adapter, descriptor variants, grammar surfaces and displaced retained journeys; O2/I2 complete the remaining named retained tests in the two reports. The full exact Developer gates were inspected separately. The supplemental binary probe expanded boundary risk coverage and found F1 despite the green corpora.

## Task 9 - Activation descriptors; whole-tree input and help conformance

### Must-fix

- **F1 — [cmd/fid/input_descriptors_activation.go:36]: accepted prepare identities are measured in bytes rather than Unicode scalar values.** `activationInputString` at lines 11-14 sets 1..128 bounds without opting into scalar counting. `internal/cli/input_validate.go:154-160` therefore uses `len([]byte(...))`; the scalar marker in `input_descriptor.go` is currently applied to typed-array variant children, not these prepare fields. Phase 21 `public-interface.md:75` defines scalar counting and line 167 gives `kit_identity`/`backend_identity` the 1..128 limits; Phase 22 extends prepare while retaining those original fields. Task 9 requires preserving the accepted input shape/constraints. The new descriptor was introduced in `785a22491f1ba07bb748d1d4fb85f88ea1a8d45d`; this is not merely a claim of an inherited domain failure. `internal/workflow/activation.go:679-690` does not impose a competing byte limit on the identities. Owning task: **phase-24/task-9**, exactly one owner. Affected assertion: **VAL-016**, including its qualifying boundary controls. Repair boundary: correct only the accepted prepare identity scalar semantics in the Activation descriptor and, if necessary, the DEC-284 allowlisted shared descriptor/validator support; add permanent boundary regressions and update task-9's report. Preserve existing ordinary callers' default length behavior and all gate/spec definitions. Required evidence: both fields accept ASCII128 and Unicode127/128 and reject Unicode129; Human/JSON binary fixtures, exact path metadata and no effects for rejections; red/green regression, existing variant/shared impact checks, all six byte-identical task-9 gates and independent report-lint COMPLETE. Re-review scope: fresh independent D1 reviewer inspects the bounded repair and retained evidence, reproduces affected ordinary and real-process checks, rechecks both reports/grammar impact and cites this review before any advancement.

Independent reproduction used the inventory-built binary, SHA-256 `9c3b297f9c08f6e1d731718bba3ae59ba7ad7888c6e22f926925c3f350370f6a`. Each of 16 rows initializes its own directory with `fid init --prefix QA`; each mode therefore starts with no existing Activation. Invoke `fid activation prepare --input - --idempotency-key qa-scalar`, adding root `--json` for JSON mode. Use this accepted request, replacing only the tested identity:

```json
{
  "kit_identity": "kit",
  "backend_identity": "backend",
  "phase": 21,
  "direction": "activate",
  "expected_prior_activation_id": null,
  "expected_prior_version": null,
  "resources": [
    {"component": "backend", "old_identity": null, "new_identity": "new"},
    {"component": "kit", "old_identity": null, "new_identity": "new"},
    {"component": "overlay", "old_identity": null, "new_identity": "new"},
    {"component": "project", "old_identity": null, "new_identity": "new"},
    {"component": "provider", "old_identity": null, "new_identity": "new"}
  ]
}
```

| Replacement, independently for each identity | Expected exit, both modes | Observed exit, both modes | Observation |
|---|---:|---:|---|
| `"a"` repeated 128 times | 0 | 0 | Qualifying ASCII control succeeds. |
| `"é"` repeated 127 times | 0 | 2 | Accepted 127-scalar identity is incorrectly rejected. |
| `"é"` repeated 128 times | 0 | 2 | Accepted upper scalar boundary is incorrectly rejected. |
| `"é"` repeated 129 times | 2 | 2 | Over-limit control rejects. |

All 12 rejecting calls have empty stdout and unchanged pre/post file-content snapshots (including DB/journal). JSON reports `invalid-argument`, the exact `input.kit_identity` or `input.backend_identity` path and `allowed:["length 1..128"]`; Human reports the same field refusal. There are eight incorrect qualifying-state rejections. The harness command `python3 /tmp/fid-review7-boundaries-pw57jd79/activation_scalar_boundaries.py` terminates exit 1; its sibling `.json` contains commands, cwd, inputs, expected/observed exits, both streams and snapshots, SHA-256 `b3841dc8ed4fd3eadbb5120739fb61340959ee392c65ecfcad0f02faf0aed71b`. Scratch receipt availability is supplemental; the payload, matrix, commands and observed defect above make reproduction independent of those temporary files.

### Improvement

- None.

Observation: the four-file shared extension otherwise has a narrow structural purpose. Raw component arrays declare their discriminator, typed variants and actual reflection types. Recursive rendering/template generation and selected-variant validation cover both Actor profiles and Gate Definitions; the component correction removes invented string bounds, counts its declared bounds in scalars and distinguishes cwd/artifact path rules. O1 reproduces accepted and refusing variants, all ten Stop Condition variants, full help pins and sensitive suppression. A sensitive ancestor redacts nested error values; the no-top-level-sensitive scan remains. Reflection/template tests do not prove scalar boundaries for every simple descriptor, which is the evidence gap exposed by F1.

I1 derives the changed fake/fid intersection from dispatch rather than asserting a hand-picked subset: nine leaves within the 17 changed leaves. Each has accepted structural input and removed-member refusal on both adapters; Gate run's accepted input reaches the same committed domain denial. Retained Review malformed inapplicability and Envelope authored-budget-state tests keep their observable refusal subjects and no-effects assertions after conversion to the shared reader. CLI documentation describes delivered descriptor/example behavior and accepted component constraints; the prepare implementation must conform to the already accepted constraint rather than changing that documentation.

## Task 9a - Grammar failures and unknown commands across the live tree

### Must-fix

- None.

### Improvement

- None.

Observation: grammar classification is installed recursively at the adapter boundary. Unknown flag, required flag, malformed typed/domain flag value and positional-count errors name the offending spelling/value before project access. Root/group unknown commands use operation-unavailable with the correct dotted owning path and sorted available words; generated help/completion are allowed where present without becoming contract inventory rows. A group without a word still prints help and exits 0.

O1 reruns the class table, domain declaration/completeness checks, inherited text/ID/opaque-state constraints and built-binary rows in Human/JSON modes inside/outside projects, with file snapshots. The current 61 typed flags are complete: 17 integers and 44 booleans, including 39 example flags. Valid controls retain native integer syntax and accepted amended members, including superseded Claim and cancelled Receipt states; malformed values refuse. Domain/business cross-field failures remain in owning handlers. CF-047 live-row disagreements and combination semantics are explicitly enumerated in the report rather than silently adopted as new grammar criteria.

The shared message helper escapes control text before its 512-scalar bound, preserves details uncut, handles already-quoted values and does not split an escape. O1 reruns short/long, exact 512/513, non-ASCII and escape-cut states; Human's error message agrees with JSON. Task 10 still owns additional next/after-the-fact lines. R19 has exactly four residual family fallbacks, and R20 still accounts for the combination helper through `flagCombinationArgument = "request"`; that is an expressly preserved semantic combination path, not an omitted grammar class.

Retained-test provenance was inspected beyond hash existence. Removed AC/Concern invocations retain their original arguments and refusal/journal subjects; only displaced silence/generic-output expectations change to the accepted typed grammar result. The closure journey retains its calls and journal checks and adds both modes for removed requirement.sync/policy.show behavior. I1 reproduces these journeys, the process exit contract and Phase 23 shared/failure contracts. Review malformed-list cases and authored-budget-state refusals remain under O1. Relevant introducing commits predate these conversions and have the recorded task Refs; no deletion or weakening of their behavioral subjects was found. Documentation distinguishes delivered grammar/error behavior from pending task-10 Human lines. No dependency or new architectural boundary is introduced.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| F1 | phase-24/task-9 | VAL-016; accepted Activation prepare scalar constraints | `cmd/fid/input_descriptors_activation.go`, `cmd/fid/input_descriptor_test.go`, `project-docs/docs/phases/phase-24/tasks/task-9-done.md`; only causally necessary support/tests in DEC-284's `internal/cli/input_descriptor.go`, `internal/cli/input_validate.go`, `internal/cli/input_validate_test.go`. Correct the two prepare identity bounds without globally changing ordinary byte-bound fields, component variant behavior, domain inputs, public-interface specs or gate definitions. | Permanent RED/GREEN tests for each identity's ASCII128/Unicode127/Unicode128 acceptance and Unicode129 refusal; initialized fresh-per-mode real-process cases with exact metadata/no effects; preserved component scalar and ordinary byte-caller regression results; all six exact task-9 commands above with original failed evidence retained and new deviation/repair receipts; independent report-lint COMPLETE. | A new independent reviewer writes the next unused D1 review, cites review-7, checks F1 against fresh source/diff/positive/negative evidence, reruns affected ordinary/shared and tagged parity plus binary boundary rows, verifies current report-lint for both tasks, inventories and task-9a grammar impact. No advance on Developer receipts alone. |

This is decision-free conformance to existing accepted semantics. It authorizes no product/spec/contract/tool/ledger change and no broad shared length-policy rewrite. F1 remains open in this artifact.

## Prior Finding Resolution

Not a fresh re-review of Window C. Review-5's task-7/task-8 evidence findings were independently resolved by review-6; this review inspected those artifacts and retained current shared/real-process impact checks. Review-6 expressly reserved full task-9 acceptance for D1. F1 here is a new task-9 prepare boundary finding, not a self-correction of an earlier D1 verdict. A future D1 re-review must preserve and explicitly disposition F1 from this review.

## Cross-cutting observations

### Current inventory and feasibility

R1-R20 and the exact section-1 recursive Bash block were regenerated from the repository root at reviewed HEAD. Markdown-escaped alternations were rendered as actual regex alternation. R4's migration listing/schema query and R12's three named file listings were separately checked. R8 used `rg -n 'definition\(' internal/cli/codes.go`. All query/list commands terminated exit 0 except R5, whose no-match exit 1 is the declared pass condition. Fresh outputs were retained under `/tmp/fid-review7-inventory-cxy9grnf/`; the unchanged Bash block built a fresh scratch binary and generated its lists under `/tmp/tmp.BktXGSPuY2`.

| Registry | Fresh result and assessment |
|---|---|
| R1/R2 | 45 call-site matches; one shared ComputeRequestDigest helper. No duplicate request-digest implementation is introduced. |
| R3/R4 | 126 field-digest matches; 38 migration files through 0019 and 17 schema references. ProjectionSchemaVersion 5 is a separate internal projection contract from physical migration/export schema 19. |
| R5 | No decodeClosed input functions; expected/observed exit 1. Shared descriptors now own structural parsing. |
| R6/R7 | Four helpers; 16 resolved-fact references. Private resolved Registration/Envelope facts remain, without restoring removed public caller identity fields. |
| R8 | 85 lines: 84 registered code definitions plus helper declaration. The added operation-unavailable matches task 9a. |
| R9/R10 | 43 failure functions and 29 mapped error assertions. Current writer inventory includes the central grammar writer; retained conversions explain the changed assertion matches. Task 10 must use this current set. |
| R11/R12 | 19 output/help/fixture references; all three Stop/Attention renderer files exist. The task-11 output corpus remains a later obligation. |
| R13/R14 | 15 judgment/PRIVATE references; live inventory completeness test present and passing in O1. |
| R15/R16 | 92 live leaves, zero missing Short; 50 flag-keyed leaves, 49 canonical keyed leaves after excluding export import. Stable grammar scope excludes diagnostics/generated helpers as specified. |
| R17/R18 | 29 root/groups; 61 typed flags (17 integer, 44 boolean), including all 39 example flags. Corpus completeness checks and generated binary lists agree. |
| R19 | Exactly four literal map sites: issue_get_list.go:129, link_response.go:89, receipt.go:87, issue_mutation_response.go:221. They are permitted semantic/persistence family fallbacks, not grammar-class writers. |
| R20 | 23 accounted invalid-argument sites versus preceding 25. Three generic grammar sites were removed, central grammar classification added, and the retained combination fallback remains visible under its named constant. Task 10 must preserve the combination semantics and cover current sites/exclusions. |

The block produced 92 leaves, 29 groups, 39 input leaves, 50 keyed leaves, zero no-Short leaves and 61 typed flags (` 271 total`, exit 0). Inventory/corpus commands are executable; current deltas match descriptor adoption and task-9a central grammar/code changes, not missing rows.

Fresh fallback feasibility also exercises a valid initialized Activation with a staged Gate Definition containing literal environment/redaction sentinel `QA-REVIEW7-SENSITIVE-SENTINEL`. A subsequent qualifying component input reaches the real generic activation fallback in both modes. Both calls exit 2, have empty stdout, unchanged file-content snapshots and no sentinel in either stream. JSON retains `invalid activation input`, `argument:"input"`, `value:null`, and `allowed:["accepted activation input"]`; Human reports the actual idempotency-result JSON array sorted/unique error. This establishes that the sensitive domain-fallback branch is reachable with no disclosure, not final task-10 VAL-067 message acceptance. Prior C feasibility remains preceding evidence; D1 does not claim later corpora have already been implemented or accepted.

### Remaining cadence

The Phase Runner must dispatch the task-9 repair and fresh independent D1 re-review before advancing. D2 remains after task 11 and covers tasks 10/11, including current grammar/Human pairing, fallbacks, Stop and repository lines. Its QA Reviewer must regenerate the task-11-close Human fixture output with the accepted command:

```text
FID_P24_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-24-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration
```

The Human walkthrough is recorded separately at `qa/human-walkthrough-after-task-11.md`; interim VAL-054/055 judgment does not replace phase-end HUMAN-02/HUMAN-03. E remains after task 14 for tasks 12-14. Final phase deterministic gates, complete validation outcomes and Human judgment remain required. This review waives none of those checkpoints.

## Blocker Candidates

- None. F1 has exactly one existing task owner and a bounded implementation/report repair under accepted scalar semantics. No protected decision or external capability is needed.

## Existing Escalation References

- ESC-24-001 and ESC-24-002: resolved history under DEC-283.
- ESC-24-003: resolved shared-descriptor authority history under DEC-284; its accepted four-file boundary applies here.
- No open escalation is used to excuse a missing gate or acceptance criterion. No escalation ID was allocated.

## Final Lines

Review written: project-docs/docs/phases/phase-24/qa/review-7.md
Must-fix: 1  Improvement: 0  Blocker candidates: 0
VERDICT: REPAIR
