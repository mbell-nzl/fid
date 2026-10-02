# Phase 23 - Independent QA Lead Checkpoint After Task 11

Reviewer: QA Lead (independent)
Date: 2026-09-25
Checkpoint scope: `qa-plan.md` section 8, "Independent QA Lead checkpoint after
task 11" - PASS required before task 12 is dispatched
Review type: INITIAL
Prior reviews: `review-1.md` .. `review-5.md` (cited for context only; this
checkpoint reproduces its own scope fresh and attests to no prior review)

## Verdict

REPAIR

The seven required built-binary reproductions were all executed and all ten
tagged real-process journeys pass at HEAD. Five of the seven areas are fully
evidenced. Two are not: the VAL-009 all-satisfied state (EVID-09 / R-07b) is
never exercised with a `review` Obligation present, and VAL-043's required
request set and its "with no Event" clause are only partially asserted. A third
defect in the same task-10 fixture (VAL-025's second real-producer mismatch)
was found while reproducing those two. All three are bounded, decision-free
repairs inside task 10's own integration fixture, and all three are contradicted
by claims the task-10 completion report makes. Tasks 12-18 feasibility is
otherwise clean.

Task 12 does not dispatch until a fresh independent QA Lead re-review verifies
the named evidence in the next unused `checkpoint-after-task-11-<R>.md`.

## Scope Reviewed

- QA plan:
  - `project-docs/docs/phases/phase-23/qa/qa-plan.md` sections 4 (EVID-09,
    EVID-13, EVID-15, EVID-17, EVID-20, EVID-21, EVID-23), 5 (deterministic
    gates), 6 (command and feasibility validation), 7 (QA-10, QA-13, QA-15,
    QA-17, QA-20, QA-21, QA-23), 8 (mid-phase cadence).
- Validation contract:
  - `project-docs/docs/phases/phase-23/qa/validation-contract.md` rows VAL-009,
    VAL-010, VAL-015, VAL-022, VAL-025, VAL-027, VAL-029, VAL-039, VAL-040,
    VAL-041, VAL-043, VAL-047, VAL-048, VAL-049, VAL-050, VAL-078, VAL-079,
    VAL-080, VAL-081.
- Task specs and completion reports: `tasks/task-1.md` .. `task-11.md` and
  `task-1-done.md` .. `task-11-done.md` (tasks 10 and 11 read in full; tasks
  1-9 read for the assertions this checkpoint reproduces).
- Forward task specs: `tasks/task-12.md` .. `task-18.md` in full.
- Accepted contract set: `contracts.md`, `public-interface.md`, `schema.md`,
  `execution-scope.md` (sections named by tasks 12-18).
- Decisions: `decisions.md` DEC-227, DEC-235, DEC-236, DEC-239, DEC-240,
  DEC-242, DEC-243, DEC-248.
- Runner ledger `run-ledger.md` (read only; its uncommitted 200-line addition is
  the permitted Runner-owned dirt and was not changed).

### Commit identity

| Artifact | Commit |
|---|---|
| Phase HEAD at reproduction | `777ee47d76990b73c88ff68bbdb9e6623b794de4` |
| Task 10 implementation | `0ff2eeef66d4633309a43ae6a341aecd38def4d0` |
| Task 10 report | `2082098e01f2a66e301896b6703a6124842c636d`, conformed by `9f98168d34e2459b6806a20f1d1f0dfd38d68164` |
| Task 11 implementation | `c5bf0f89c69c57003ece353ef005a7596dd266e4` |
| Task 11 report | `777ee47d76990b73c88ff68bbdb9e6623b794de4` |
| Prior review 5 (tasks 7-9) | `2b33bb133a5712e39ef4c1c9d0614fe64ed627a8` |
| Planning base (for R-13 baselines) | `0195897a347dfe1872db0f7be3800c5c3ccdefcf` |

All hashes independently resolved with `git rev-parse`. Both implementation
commits carry their `Refs: phase-23/task-10` and `Refs: phase-23/task-11`
footers.

### Reproduction environment

`go version go1.26.2 linux/amd64`; `Linux 6.6.114.1-microsoft-standard-WSL2
x86_64`; `nproc` 6; AMD Ryzen 7 5700G. Working directory
`/home/mark/projects/fid` for every command. Heavy gates were run serially;
only read-only probes overlapped them.

## Summary

- Review outcome: REPAIR
- Must-fix findings: 3
- Improvement findings: 4
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-009, VAL-015, VAL-022, VAL-025, VAL-027,
  VAL-029, VAL-039, VAL-040, VAL-041, VAL-043, VAL-047, VAL-048, VAL-049,
  VAL-050, VAL-079

## Required Checkpoint Reproduction

Every row below was reproduced in this session against the built binary. The
tagged suite builds `./cmd/fid` into a temporary directory and drives it as a
real process (`cmd/fid/integration/main_test.go` `TestMain`,
`integrationFIDBinary`/`closureRunCommand`); no assertion below rests on an
in-process call.

Command (one invocation, all ten journeys):

```text
go test -tags integration -race -timeout 30m -count=1 -v ./cmd/fid/integration/... \
  -run '^(TestP23ReviewRecordRealProcessJourney|TestP23FindingDispositionRealProcessJourney|TestP23ReviewEvidenceRealProcessJourney|TestP23ReviewCyclesRealProcessJourney|TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23HumanOutcomeRequiredRaisedByReviewRecord|TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess|TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen)$'
```

Result: PASS - `ok  github.com/mbell-nzl/fid/cmd/fid/integration  37.798s`
(terminated `exit_code: 0`). Ten top-level journeys PASS, nine named subtests
PASS, zero FAIL lines.

| # | Area | Reproduced through | Result | Outcome |
|---|---|---|---|---|
| 1 | Independence both states (VAL-015) | `TestP23ReviewRecordRealProcessJourney` (built binary) + `TestCheckIndependenceRejectsClaimAndRunHistoryDirections` (DEC-239 isolated durable-history query) | PASS | ON TRACK |
| 2 | Resolution authority (VAL-022) | `TestP23FindingDispositionRealProcessJourney` (built binary, 6 subtests) + `TestTerminalFindingDispositionAuthorityMatrix` (12 subtests, both terminal actions) | PASS | ON TRACK |
| 3 | Cycle charging and exhaustion (VAL-027, VAL-029) | `TestP23ReviewCyclesRealProcessJourney` (built binary) | PASS | ON TRACK (import leg remains task 16 per the contract) |
| 4 | Protected-Decision blocking and satisfaction (VAL-039-041) | `TestP23DecisionReferenceProtectedLifecycleRealProcess` (built binary) | PASS | ON TRACK |
| 5 | Human Outcome provenance (VAL-043) | `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess` (built binary) | PASS, but under-asserted | **NOT ON TRACK - MF-2** |
| 6 | Attention resolution through an Outcome (VAL-048-050) | `TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess` and `...RejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen` (built binary) | PASS | ON TRACK |
| 7 | VAL-009 all-satisfied state (EVID-09 / R-07b) | `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess:144` (built binary) | PASS, but the required state is not the one exercised | **NOT ON TRACK - MF-1** |

### Area 1 - Independence both states (VAL-015)

Both required directions were observed.

- Real CLI rejection of a durable mutating-Claim holder:
  `TestP23ReviewRecordRealProcessJourney/mutating_claim_actor_is_rejected_without_an_event`
  PASS. `cmd/fid/integration/p23_review_record_test.go:24-36` asserts exit 4
  with `review-not-independent`, stderr carrying `"type":"claim"`, the exact
  `"id":"<claimID>"`, and `"conflicting_actor_id":"<actor>"`, with
  `p21JournalEventCount` byte-equal before and after.
- Run-holder-only direction through the DEC-239 isolated durable-history query:
  `go test -race -count=1 -timeout 5m -run '^(TestCheckIndependenceRejectsClaimAndRunHistoryDirections|TestTerminalFindingDispositionAuthorityMatrix|TestChangesRequiredReviewChargesExactlyOneReviewCycle)$' ./internal/review/ -v`
  - PASS - `ok  github.com/mbell-nzl/fid/internal/review  1.145s` (terminated
  `exit_code: 0`). `internal/review/review_test.go:263-303` asserts
  `NotIndependentError{ConflictingType:"claim", ConflictingID:"CLM-claim-holder"}`
  and then `{ConflictingType:"run", ConflictingID:"RUN-run-holder"}`, each with
  the exact conflicting Actor. The fixture is an in-memory `claims`/`runs`
  history only, correctly labelled diagnostic query evidence, not a CLI
  lifecycle claim.
- Same independent Reviewer's post-rework Review succeeds:
  `p23_review_record_test.go:98-110` closes the review Claim, drives a real
  mutating rework Run, then records a second Review as the same
  `actor-denied`, asserting the Event count moves by exactly one for the
  Review plus both Finding resolutions.

### Area 2 - Resolution authority (VAL-022, DEC-242)

- Built binary: `TestP23FindingDispositionRealProcessJourney` PASS with all six
  subtests PASS - `same_review_and_malformed_inapplicability_are_non_mutating`,
  `omitted_and_empty_optional_collection_replay_the_task4_digest`,
  `address_requires_mutating_claim_and_stays_unresolved`,
  `non_independent_and_blocked_terminal_reviews_reject`,
  `changed_inapplicability_reason_conflicts_without_an_event`,
  `second_terminal_and_unknown_findings_are_non_mutating`.
- Domain matrix for **both** terminal actions:
  `TestTerminalFindingDispositionAuthorityMatrix` PASS with twelve subtests -
  `resolved/{creating_Review, mutating_Claim_holder, blocked_verdict,
  different_independent_Actor, same_independent_Reviewer_later,
  second_terminal}` and the identical six under `inapplicable/`. The positive
  cases assert the stored `action`, `actor_id`, and
  `qualifying_review_id:"REV-current"`.

### Area 3 - Cycle charging and exhaustion (VAL-027, VAL-029)

`TestP23ReviewCyclesRealProcessJourney` PASS.
`cmd/fid/integration/p23_review_evidence_budget_test.go:66-128` observes,
through the built binary: a `blocked` verdict charging nothing
(`result.budget == nil`); the exhausting `changes-required` Review committing
in exactly one Event with `budget{name:"review-cycles", consumed:1,
remaining:0, state:"exhausted"}` plus exactly one Stop
(`review-rework-exhausted`) and one Attention in the same response; a full
`fid rebuild` followed by replay returning `replayed:true` with the journal
length unchanged (no second charge); `run.observe-budget` on the fid-owned unit
rejected with no new `budget_samples` row; and both
`claim.acquire` and `run.start` refused `review-rework-exhausted` with
`"event":null`, `"replayed":false`, and `claims`/`runs` row counts unchanged.
The import leg is explicitly deferred to task 16 by the contract row itself.

### Area 4 - Protected-Decision blocking and satisfaction (VAL-039, VAL-040, VAL-041)

`TestP23DecisionReferenceProtectedLifecycleRealProcess` PASS.
`cmd/fid/integration/p23_decision_reference_test.go` observes two Stops open at
activation (`:24`), per-class denial naming only the declared affected action
with a bracketed `p21JournalEventCount` no-Event check (`:29-30`, `:81-82`,
helper at `:133-151`), `rejected`/`deferred` reference outcomes and an
authority mismatch each leaving both Stops open (`:32-38`), the outcome-only
and reference-only states leaving the Stop open (`:55-59`), the both-approved
pair resolving exactly its own Stop while the other stays open (`:66-70`), and
the second class resolving last (`:84-86`). Source-digest mismatch and missing
source both reject with the journal length unchanged (`:108-111`).

### Area 5 - Human Outcome provenance (VAL-043) - NOT ON TRACK

`TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess` PASS, but the
assertions it contains do not cover the required set. See Must-fix MF-2.

What is proven: the agent-actor-type negative (`:34-36`), the
standalone/unregistered attributed-actor negative (`:41-43`), the
missing-capability negative against a dedicated finalized Activation
(`:47`, helper `:434-468`), and the valid case recording exactly one Event with
`provenance.activation_id` equal to the finalized Activation (`:51-65`).

What is not proven: the **unfinalized-Activation** request (the contract's and
the task spec's own second negative) is never executed - it is replaced by a
code comment at `:38-40` asserting it shares the `unregistered` case's code
path; and the "with no Event" clause is asserted for **none** of the three
negatives, because `assertP23AuthorityMissing`
(`cmd/fid/integration/p23_human_outcome_test.go:356-368`) accepts a `before int`
parameter and never uses it. `grep -n "p21JournalEventCount"` on that file
confirms the `before` values captured at `:34` and `:41` are never compared to
an `after`, and the third call site passes the literal `0` (`:47`).

### Area 6 - Attention resolution through an Outcome (VAL-048, VAL-049, VAL-050)

Both task-11 journeys PASS.

- VAL-048: `p23_attention_resolution_test.go:159-186` records an approved
  Outcome of an unrelated class and a `rejected` result for the matching class;
  each resolves zero Stops (`len(stops) != 0` guard) and
  `assertP23HumanOutcomeStopState(..., "security", "open")` re-reads the Stop
  from the projection. The matching class with `approved` resolves exactly one
  Stop and the state flips to `resolved`.
- VAL-049: `:88-114`. `attention.resolve` with one added unknown field
  `outcome_id` rejects with the journal length unchanged; the byte-identical
  request without that key succeeds; `p23StopIdentity` re-reads `subject`,
  `code`, and `raised_at` after resolution and asserts equality with the values
  captured immediately after the raise; `p23StopResolutionFact` asserts each
  Stop's stored `fact_event_id` is exactly the qualifying `review.recorded`
  Event ID.
- VAL-050: `:150-185`. A Stop raised on an unrelated `review-recorded`
  Condition stays open (`assertP23AttentionStopUnresolved`) through both
  non-qualifying attempts and after the qualifying `security` Outcome resolves
  exactly the matching `human-outcome-recorded` Stop.

### Area 7 - VAL-009 all-satisfied state (EVID-09 / R-07b) - NOT ON TRACK

The blocking half is proven: `p23_human_outcome_test.go:110-115` shows
`envelope.complete` rejected `human-outcome-required` with
`obligation_id:"OBL-human-security"` and `outcome_class:"security"`, journal
length unchanged.

The all-satisfied half is exercised (`:144`, `envelope.complete` succeeds;
`:142`, the terminal `qa_pass` transition succeeds) but **not with the required
Obligation set**. See Must-fix MF-1.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-009 | task 10 (completing half) | **NOT ON TRACK** | `envelope.complete` at `p23_human_outcome_test.go:144` succeeds against an Envelope holding one `gate-evidence` and two `human-outcome` Obligations and **no `review` Obligation** (fixture `{HumanObligations: true}` at `:22`; `review.record` is absent from the journey's action list at `:21`). MF-1. |
| VAL-010 | task 10 (re-verification) | NOT ON TRACK (secondary) | `:137-140` re-verifies `not-applicable` with linked Human Outcome Evidence only; the contract's "linked review **and** Human Outcome evidence" state is unreachable in that fixture for the same reason as MF-1. Improvement I-4. |
| VAL-015 | task 4 | ON TRACK | Reproduced, Area 1 above. |
| VAL-022 | task 5 | ON TRACK | Reproduced, Area 2 above. |
| VAL-025 | tasks 6, 10 (final closure) | **NOT ON TRACK** | Task 10 re-proves only the `human-outcome`→`gate-evidence` mismatch, twice (`:70-71` and the single-element loop at `:149-157`, both asserting `obligation_type:"gate-evidence"`). The `human-outcome`→`review` mismatch is never exercised with the real producer, and neither denial is bracketed by a journal-count check. MF-3. |
| VAL-027 | task 6 | ON TRACK | Reproduced, Area 3 above. Import leg is task 16 by contract. |
| VAL-029 | task 6 | ON TRACK | Reproduced, Area 3 above. |
| VAL-039 | task 9 | ON TRACK | Reproduced, Area 4 above. |
| VAL-040 | task 9 | ON TRACK | Reproduced, Area 4 above. |
| VAL-041 | task 9 | ON TRACK | Reproduced, Area 4 above. |
| VAL-043 | task 10 | **NOT ON TRACK** | Three of five required requests executed; unfinalized-Activation negative absent; "no Event" clause asserted for none. MF-2. |
| VAL-047 (first half) | task 11 | ON TRACK | Independently probed with the binary built at HEAD in this session: `/tmp/fid-qa-checkpoint --json capabilities` returns `"capability_id":"attention.v1"` with `"limits":{"human_outcome_resolution":false`. Source confirms `internal/workflow/capabilities.go:105` and `internal/export/archive_profiles.go:5` both still declare `false`. |
| VAL-048 | task 11 | ON TRACK | Reproduced, Area 6 above. |
| VAL-049 | task 11 | ON TRACK | Reproduced, Area 6 above. |
| VAL-050 | task 11 | ON TRACK | Reproduced, Area 6 above. |
| VAL-078 (resolution) | task 11 | ON TRACK with a note | `review-recorded` and `finding-resolved` each get a non-qualifying `stop resolve` rejecting `resolution-invalid` with no Event and a qualifying one succeeding (`:57-81`). Improvement I-3 records that `human-outcome-recorded`'s pair is evidenced through Outcome recording rather than the same `stop resolve` route. |
| VAL-079 | tasks 1-11 and phase end | ON TRACK | Fresh measurement below; all three named packages inside budget, with `cmd/fid` at 97.6% of its 600 s budget (Improvement I-1). |
| VAL-080 | tasks 9, 10, 11 | ON TRACK | Ten-row closed `currentConditionShapes` confirmed in `internal/journal/current_event_value_catalog.go:59-70`: six retained Phase 22 rows byte-unchanged, `decision-satisfied`/`human-outcome-recorded`/`finding-resolved` each carrying exactly their own discriminator, and `review-recorded` at the bare four-field form. |
| VAL-081 | task 9, re-confirmed 10 and 11 | ON TRACK | `internal/journal/current_event_value_relations.go` `case "Condition"` untouched by tasks 10 and 11 (`git show --name-only` on both implementation commits lists neither that file nor `current_event_value_types.go`); task 11's own `TestTask11FindingResolvedConditionDispatchPreservesNullableBinding` ran green inside the full ordinary gate below. |

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`) |
| Lint | `golangci-lint run ./...` | PASS | `0 issues.` (terminated `exit_code: 0`) |
| Ordinary race suite | `go test -race -timeout 30m -count=1 ./...` | PASS | last line `ok  github.com/mbell-nzl/fid/internal/workflow  228.577s`; zero `FAIL` lines; 21 packages `ok` (terminated `exit_code: 0`, wall 593 s) |
| Tagged integration gate | `go test -tags integration -race -timeout 30m -count=1 ./cmd/fid/integration/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration  688.252s` (terminated `exit_code: 0`, wall 689 s) |
| Build verification | `go build -o /tmp/fid-qa-checkpoint ./cmd/fid` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`) |
| Rebuild benchmark | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` | PASS | `ok  github.com/mbell-nzl/fid/internal/journal  24.501s`, in-run `21733 ms/rebuild` against the 30000 ms bar (terminated `exit_code: 0`) |
| Checkpoint journeys | the ten-journey `-run` invocation above | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration  37.798s` (terminated `exit_code: 0`) |
| Named domain fixtures | `go test -race -count=1 -timeout 5m -run '^(TestCheckIndependenceRejectsClaimAndRunHistoryDirections\|TestTerminalFindingDispositionAuthorityMatrix\|TestChangesRequiredReviewChargesExactlyOneReviewCycle)$' ./internal/review/ -v` | PASS | `ok  github.com/mbell-nzl/fid/internal/review  1.145s` (terminated `exit_code: 0`) |
| Report/spec mechanical control | `project-docs/tools/report-lint.sh` for tasks 10 and 11 | PASS | Both `OUTCOME COMPLETE` / `PASS report-lint completed with no failures` (terminated `exit_code: 0`) |
| Exact task-command audit | tasks 10 and 11 `## Exact quality-gate results` against each spec's `## Exact Commands To Run` | PASS | Both specs list the same five commands with no conditional and no expected-non-zero command; both reports record all five `PASS` with a token and observed termination. The `-count=1` addition is recorded as a deviation in both reports and is a cache-bypass flag, not a narrowing. |
| Commit identity | `git rev-parse` on all six task-10/11 commits; `Refs:` footers | PASS | All resolve; both implementation commits carry their `Refs: phase-23/task-N` footer. |
| Closed grant allowlist | `internal/workflow/activation.go` `sortedUniqueCapabilities` vs. the second copy in `internal/journal/current_event_value_relations.go:63` | PASS | Both hold the identical sorted 30-name set; the phase brief's "seven grant actions" over the 23-name planning base reconciles exactly. |
| Report claim accuracy | task-10-done.md VAL-043, VAL-025, VAL-009 entries vs. the cited source lines | **FAIL** | Three claims are contradicted by the code they cite. See MF-1, MF-2, MF-3. |

## Tasks 12-18 Fresh Feasibility Sweep

Every non-default command and fixture assumption tasks 12-18 prescribe was
re-validated at HEAD `777ee47`. Nothing carries forward from the planning-base
validation in `qa-plan.md` section 6.

### Non-default commands

| Task | Command or fixture | Validation method | Result | Evidence |
|---|---|---|---|---|
| 12, 15, 16 | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` | executed | PASS | `ok  github.com/mbell-nzl/fid/internal/journal  24.501s`, `21733 ms/rebuild` (terminated `exit_code: 0`). Carries no `-timeout`; 24.5 s against the 600 s default leaves ample margin. Benchmark declared at `internal/journal/rebuild_bench_test.go:19`. |
| 12 | `rg -n --glob '!*_test.go' "requirement\.synced\|last_synced_hash" cmd internal` (declared expected `exit_code: 1` at task close) | executed | PASS (syntax and scope) | `exit_code: 0` today with matches in exactly the three files the spec predicts - `internal/requirements/sync.go`, `internal/journal/current_event_catalog.go:119`, `internal/export/archive_semantics_catalog.go:91` - all three inside task 12's own file scope (the third assigned by `execution-scope.md` section 6 / R-09). The declared non-zero pass condition is reachable by this task alone. |
| 15 | `rg -n --glob '!*_test.go' "task_completions\|ac_verifications" cmd internal` | executed | PASS (syntax and scope) | `exit_code: 0`; matches only `internal/issues/detail.go` (3), `internal/issues/transition_rules.go:135`, `internal/maintenance/operational.go:185`, `internal/maintenance/errors.go:62`, `internal/maintenance/warnings.go` (3) - every file named in task 15's own scope. |
| 15 | `rg -n --glob '!*_test.go' "requirement\.synced\|last_synced_hash" cmd internal` | executed | PASS | Same three files as the task-12 row; task 12 removes all three before task 15 runs, so this search is already satisfied by then and is a genuine regression guard here. |
| 15 | `rg -n --glob '!*_test.go' "has_completion_report\|concerns_addressed" cmd internal` | executed | PASS (syntax and scope) | `exit_code: 0`; matches only `internal/workflow/rules.go:7-8`, inside task 15's scope. |
| 13, 14, 17, 18 | no conditional or expected-non-zero commands declared | spec inspection | PASS | Each of the four `## Exact Commands To Run` blocks is the five-command default plus the DEC-227 `-timeout 30m`; each `List any conditional command` reads `- None.` |

### Fixture, file-state, and contract assumptions

| Task | Assumption | Result | Evidence |
|---|---|---|---|
| 12 | Fence map holds exactly `requirement.referenced` and `requirement.linked`, so the DEC-192 "delete exactly two keys then remove the empty map" instruction is executable | PASS | `internal/journal/project_effects.go:12-15` holds exactly those two keys; the guard branch is at `:18-20`. |
| 12 | `validateCurrentUpdate`'s closed `switch effect.Table` is at `:14`; task 8's sibling `judgment_concerns` entry exists | PASS | `internal/journal/current_event_updates.go:14` is the switch; `judgment_concerns` entries at `:41` (values) and `:86` (guards). No `requirement_refs` entry yet, as the spec predicts. |
| 12 | `current_event_json.go`'s `effect.Table + "." + column` switch is at `:14` | PASS | Confirmed at `internal/journal/current_event_json.go:14`. |
| 12 | `RequirementRef`/`RequirementLink` shapes and the two result rows exist in the accepted interface | PASS | `public-interface.md:116`, `:118` (section 3) and `:180`, `:181` (section 4). |
| 12 | `internal/links/{registry,add,list,errors,types}.go` and `internal/requirements/` present | PASS | All present; `internal/requirements/` holds `parser.go`, `sync.go`, `attribution.go`, `errors.go`, `types.go` for deletion. |
| 13 | `reasonCodeDefinitions` declaration plus exactly three other uses, including `internal/cli/envelope.go:38` | PASS | `internal/cli/codes.go:232` (declaration), `:396`, `:407`, and `internal/cli/envelope.go:38` - four occurrences, matching the spec exactly. `reasonCodeExpectation` at `internal/cli/codes_test.go:30`; `cmd/fid/evidence_test.go:465`'s `reasonCodes` is the `assertEvaluation` parameter the spec correctly says not to rename. |
| 13 | Eleven named `cmd/fid/` leaf files present | PASS | `review.go`, `finding.go`, `concern.go`, `deviation.go`, `decision.go`, `human_outcome.go`, `requirement.go`, `evidence.go`, `envelope.go`, `claim.go`, `root.go` all exist. |
| 14 | Phase 22 precedent fixture and doc exist; the Phase 23 target does not yet | PASS | `cmd/fid/integration/p22_human_fixtures_test.go` and `project-docs/docs/phases/phase-22/qa/human-fixtures.md` present; `phase-22/public-interface.md:300` is section 6 "Human And Independent Review"; `phase-23/qa/human-fixtures.md` correctly absent. |
| 14 | `internal/workflow/{resume,resume_actions}.go`, `cmd/fid/{resume,resume_human,resume_human_guidance}.go` present | PASS | All five present. |
| 15 | Every file named for deletion or edit exists; migration 0018 does not | PASS | `internal/issues/{ac,ac_verify,capability_unavailable,concern_address,concern_list,concern_raise,detail,get,list,issue,read,transition_rules}.go`, `internal/workflow/{rules,registry,workflow}.go`, `internal/maintenance/{errors,warnings,operational,check}.go`, `internal/stats/stats.go`, `internal/journal/current_event_columns.go`, `cmd/fid/{issue_ac,concern}.go` all present. Removal precedents `migrations/0014_p21_retire_v1_authority_execution.up.sql` and `migrations/0016_p22_remove_escalations.up.sql` present; `migrations/0018_p23_remove_v1_judgment.up.sql` correctly absent. |
| 16 | `internal/export/{export,archive_profiles}.go` present and the `schema-unsupported` archive-check path exists | PASS | Both files present; `schema-unsupported` exercised today by `internal/export/archive_v2_test.go:111`, `:139`. |
| 17 | The false `human_outcome_resolution` row exists in both advertising surfaces and in the live binary | PASS | `internal/workflow/capabilities.go:105` and `internal/export/archive_profiles.go:5`; built-binary probe at HEAD returns `"human_outcome_resolution":false`. Task 17's flip is reachable. |
| 18 | `p22ClosedFieldGrammar`'s `Condition` row still holds six types and the hard counts are still 50 and 90 | PASS | `cmd/fid/integration/p22_closed_values_test.go:25` holds exactly the six Phase 22 types; `:108` asserts `accepted != 50 \|\| rejected != 90`. The spec's I-41 description of this self-contained wire oracle is still exactly accurate. |
| 18 | The wire authority to verify the extension against holds all ten types | PASS | `internal/journal/current_event_value_catalog.go:59-70`, ten rows. |
| 18 | Grant action set is already at the thirty names task 18 must audit | PASS | Counted 30 in `sortedUniqueCapabilities` and 30 in the journal copy, identical and sorted. `internal/workflow/authority.go` present for audit support. |
| 18 | `cmd/fid/integration/testdata/` and the Phase 22 audit precedent exist | PASS | Both present. |
| 12-18 | Every `## First Failing Test` name is genuinely new | PASS | `grep -rn "func <name>("` over `cmd internal` returns no match for `TestRequirementRecordReferenceSupersedesPriorActiveRow`, `TestPhase23LeavesRejectUnknownRequestFields`, `TestResumeReportsRealFindingsAndJudgmentSections`, `TestLatestSchemaHasNoV1JudgmentTables`, `TestPortableRoundTripCoversEveryPhase23Table`, `TestFinalCapabilityMatrixMatchesAcceptedTable`, `TestGrantActionSetIsExactlyThirtyNames`. |

### R-13 test-time headroom (DEC-227, VAL-079)

Measured fresh in this session at HEAD `777ee47`, gates run serially with
`-count=1` so no result is a cached replay.

| Package | Measured now | Budget | Headroom | Planning base (`0195897`) | Task 11 report |
|---|---|---|---|---|---|
| `cmd/fid/integration` (tagged) | **688.252 s** | 900 s | 211.7 s (23.5%) | 560.243 s | 700.168 s |
| `cmd/fid` (ordinary) | **585.401 s** | 600 s | **14.6 s (2.4%)** | 471.955 s | 546.059 s |
| `internal/workflow` (ordinary) | **228.577 s** | 600 s | 371.4 s (61.9%) | 458.353 s | 213.335 s |

Next-slowest ordinary packages: `internal/export` 322.340 s, `internal/journal`
237.085 s, `internal/issues` 234.029 s - all comfortably inside budget.

Two conclusions for tasks 12-18:

1. The tagged gate at 688.252 s is **88 s past Go's implicit 10-minute
   per-binary default**. DEC-227's explicit `-timeout 30m` is therefore load
   bearing, not defensive, for every one of tasks 12-18. Every one of the seven
   remaining specs carries it; none may drop it.
2. `cmd/fid` at 585.401 s sits at 97.6% of both the VAL-079 600 s budget and
   Go's implicit default. Tasks 13 (real-process conformance over every new and
   amended leaf) and 14 (Human renderers and fixtures) both add `cmd/fid`
   package tests. Recorded as Improvement I-1 rather than Must-fix, because
   VAL-079 is not breached today, but the margin is roughly one added
   integration-style `cmd/fid` test.

## Per-Task Findings

### Task 10 - Provenanced Human Outcomes

#### Must-fix

- **MF-1 (VAL-009, EVID-09, R-07b) - the all-satisfied state never contains a
  `review` Obligation.** `cmd/fid/integration/p23_human_outcome_test.go:22`
  builds its Envelope with `p22ExecutionFixture{HumanObligations: true}`. In
  `cmd/fid/integration/p22_fake_test.go:372-376` the single non-Human Obligation
  is assigned `"type": "gate-evidence"` at `:372` and **replaced in place** by
  `"type": "review"` at `:374` under the same `"id": "OBL-p22-corpus"` when
  `ReviewObligation` is set; `obligations := []any{obligation}` at `:376` then
  appends only `human-outcome` rows. The two types are therefore mutually
  exclusive in every fixture in the repository, and
  `grep -rn '"type": "review"' cmd/fid/integration/*.go` returns exactly that
  one construction site. The journey at `:21` does not even grant
  `review.record`, so no Review can be recorded in it. The `envelope.complete`
  at `:144` and the terminal `qa_pass` transition at `:142` therefore prove the
  all-satisfied state for **gate-evidence + human-outcome only**. VAL-009's
  evidence method requires "the qualifying state with gate-evidence, review, and
  human-outcome Obligations all satisfied"; its fail criterion is "the
  all-satisfied state is never exercised"; `task-10.md`'s own Validation
  Contract Coverage entry says "permitted only once the gate-evidence, review,
  and human-outcome Obligations are all satisfied". The task-10 report's claim
  that "Lines 142-144 confirm `envelope.complete` succeeds once every Obligation
  (gate-evidence, review, human-outcome) is satisfied" is contradicted by the
  fixture it cites.

- **MF-2 (VAL-043) - two of the five required requests are missing or
  unasserted.** (a) The contract's and `task-10.md`'s **unfinalized-Activation**
  negative is never executed. `p23_human_outcome_test.go:38-40` substitutes a
  comment ("This is the same `TrustedActorUnavailableError` code path an
  unfinalized Activation produces") for the state itself; the `unregistered`
  request at `:41-43` exercises standalone Issue-inventory attribution, which is
  the contract's *fourth* negative, not its second. No test anywhere else
  covers it - `internal/review/` has no human-outcome test file. (b) The "with
  no Event" clause of VAL-043's pass criterion is asserted for **none** of the
  three executed negatives: `assertP23AuthorityMissing`
  (`p23_human_outcome_test.go:356-368`) takes `before int` and never reads it,
  so the counts captured at `:34` and `:41` are dead, and `:47` passes the
  literal `0`. `task-10.md`'s acceptance criterion reads "Agent actor type,
  unfinalized Activation, and missing capability each reject `authority-missing`
  **with no Event**". The task-10 report's claim that each negative asserts "no
  Event via the `p21JournalEventCount` before/after check" is false.

- **MF-3 (VAL-025, DEC-243) - only one of the two required real-producer
  mismatches is re-proved.** DEC-243 states "Task 10 records a real approved
  Human Outcome, links it to a matching human-outcome Obligation, and **proves
  the two mismatches with that real producer**", and `task-10.md`'s acceptance
  criterion reads "attempts to link it to gate-evidence **and review** reject
  `outcome-not-permitted` with exact closed details and **no Event/projection
  change**". In practice both mismatch attempts target the *same*
  `gate-evidence` Obligation: `:70-71` links `releaseOutcomeID` to
  `OBL-p22-corpus` asserting `obligation_type:"gate-evidence"`, and the
  single-element loop at `:149-157` links `securityOutcomeID` to the same
  `OBL-p22-corpus` asserting the same `"gate-evidence"`. The
  `human-outcome`→`review` cell holds only task 6's explicitly-labelled
  intermediate production-classifier evidence
  (`task-6-done.md` VAL-025 table), which DEC-243 classifies as intermediate.
  Additionally, neither `assertP23OutcomeNotPermitted` call site is bracketed by
  a journal-count check and the helper (`:384`) performs none, so the "no
  Event/projection change" clause is unevidenced for both. The task-10 report's
  description of the second cell as "to review type" while asserting
  `obligation_type:"gate-evidence"` is self-contradictory.

All three share one root cause: no integration fixture in the repository can
produce an Envelope carrying `gate-evidence`, `review`, and `human-outcome`
Obligations simultaneously, and the two task-10 assertion helpers discard the
no-Event evidence they were given.

#### Improvement

- **I-4** VAL-010's task-10 clause ("re-verified here with linked review **and**
  Human Outcome evidence present") rides on the same missing `review`
  Obligation; `:134-136` is explicit that it substitutes "review-equivalent
  (Human Outcome) Evidence". The MF-1 fixture repair should close this in the
  same pass.

### Task 11 - Attention resolution through Human Outcomes

#### Must-fix

- None.

#### Improvement

- **I-2** `internal/workflow/stop_conditions.go` (43 added lines in
  `c5bf0f8`) is not literally named in `task-11.md`'s file list, which reads
  "`internal/workflow/resolve.go`, `evaluation.go`, `attention` readers" plus
  "`internal/workflow/stop.go` - the closed Stop registries". The report
  records this as an explicit deviation with sound reasoning (it is where every
  other Condition's qualifying-fact logic lives and the only implementable site
  for VAL-078's behaviour, which the spec assigns to this task), and the change
  is additive, in-package, and decision-free. It is not a Must-fix: no
  protected item changed and nothing in the repair is a Developer-authorable
  artifact. Fix direction: the Tech Lead should ratify it in `decisions.md` at
  wrap-up the way DEC-244..DEC-248 ratified the comparable task 6 and task 9
  scope readings, so the phase's file-ownership record stays complete.
- **I-3** VAL-078's `human-outcome-recorded` non-qualifying/qualifying pair is
  evidenced through `human-outcome record` responses rather than through an
  explicit `stop resolve` carrying a non-qualifying `fact_event_id`, which is
  the route used for `review-recorded` and `finding-resolved` (`:59-64`,
  `:74-79`) and for `decision-satisfied` at task 9. The current evidence is
  behaviourally correct - the readiness Stop resolves atomically inside the
  recording Event and has no caller-driven resolve path - so this is not a
  Must-fix. Fix direction: the window 10-12 reviewer or final review should
  record explicitly which of the four types is proved by which route, so
  VAL-078's "for each, a resolution attempt with a non-qualifying fact Event"
  is not read as uniform when it is not.

### Cross-task

#### Improvement

- **I-1** `cmd/fid` ordinary package time measured 585.401 s against the
  VAL-079 600 s budget - 97.6% of budget and 97.6% of Go's implicit 10-minute
  per-binary default. Trend from the planning base: 471.955 s → 539.822 s
  (task 10) → 546.059 s (task 11) → 585.401 s (this session, same HEAD as task
  11, so the delta is machine-load variance on an already-thin margin). Tasks 13
  and 14 both add `cmd/fid` package tests. Fix direction: tasks 13 and 14 should
  measure `cmd/fid` before and after their own additions and treat a breach as
  the bounded speed repair DEC-227 assigns to the causing task, rather than
  discovering it at the phase-end VAL-079 gate.

### Tasks 1-9

No new finding. Areas 1-4 of this checkpoint reproduced their assertions
independently and all passed; `review-1.md` .. `review-5.md` remain the
window-level records and are cited, not re-litigated.

## Prior Finding Resolution

Not applicable: this is the initial task-11 checkpoint. No prior
`checkpoint-after-task-11*.md` exists.

## Repair Plan

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-1 | task-10 | VAL-009 (and VAL-010 via I-4) | `cmd/fid/integration/p23_human_outcome_test.go` and, if the fixture builder must change, `cmd/fid/integration/p22_fake_test.go`'s Obligation construction (`:372-380`) only - give the review Obligation its own ID so it can coexist with `gate-evidence` instead of replacing it. No product, domain, adapter, schema, Event, code, or contract change. | A real-process journey in which one Envelope carries a `gate-evidence`, a `review`, and a `human-outcome` Obligation; the guarded action is denied while any one of the three is unsatisfied; and `envelope.complete` (and/or `issue.transition-terminal-success`) succeeds only after all three are satisfied, with the journal length unchanged on every denial. Correct the task-10 report's VAL-009 entry to cite the new evidence. | Rerun the ten-journey `-run` invocation above plus the full tagged gate; read the new fixture's Obligation list and confirm all three types are present and satisfied at the permitted action. |
| MF-2 | task-10 | VAL-043 | `cmd/fid/integration/p23_human_outcome_test.go` only: add the fourth negative driving a **prepared-but-not-finalized** Kit Activation whose human Actor holds `human-outcome.record`, and make `assertP23AuthorityMissing` actually compare `p21JournalEventCount` before and after (or bracket every call site). No product change; `internal/review/human_outcome.go` behaviour is already correct. | All five VAL-043 requests executed through the built binary: agent actor type, unfinalized Activation, missing capability, standalone Issue-inventory attribution, and the valid Kit-mode case; the first four each exit 5 with `authority-missing` **and** an unchanged journal length; the valid case appends exactly one Event with its `provenance.activation_id` bound. Correct the task-10 report's VAL-043 entry. | Rerun `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess` with `-v`; read `assertP23AuthorityMissing` and confirm the `before` value is consumed; confirm five distinct requests exist. |
| MF-3 | task-10 | VAL-025 (final closure), DEC-243 | `cmd/fid/integration/p23_human_outcome_test.go` only: link a real recorded approved Human Outcome to a **`review`**-type Obligation and assert `outcome-not-permitted` with `obligation_type:"review"`, and bracket both mismatch denials with a journal-count check. The `review` Obligation the MF-1 repair adds is the natural target, so MF-1 and MF-3 should be repaired together. | Both real-producer mismatches observed: `human-outcome`→`gate-evidence` returning `outcome_not_permitted` with `obligation_type:"gate-evidence"`, and `human-outcome`→`review` returning it with `obligation_type:"review"`, each with the journal length and projection unchanged; the nine-cell reconciliation in the task-10 report updated to name the real-producer route for both, and the self-contradictory "to review type … `obligation_type:"gate-evidence"`" sentence corrected. | Rerun the same journey; read both denial call sites and the helper; re-read the task-10 report's VAL-025 nine-cell reconciliation against task-6-done.md's table. |

All three repairs are inside task 10's own declared scope ("Matching in-package
tests and `cmd/fid/integration/` fixtures"), require no decision, and change no
accepted product, architecture, data-model, public-interface, dependency,
security, gate, or validation-contract item. The phase does not advance and
task 12 is not dispatched until a fresh independent QA Lead writes
`checkpoint-after-task-11-2.md` verifying the Required evidence column above.

## Blocker Candidates

- None. No finding requires protected product, architecture, data-model,
  public-interface, dependency, security, external-infrastructure, or
  specification authority, and none exceeds task 10's existing scope. The
  no-progress recovery budget is not engaged.

## Existing Escalation References

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11.md
Must-fix: 3  Improvement: 4  Blocker candidates: 0
VERDICT: REPAIR
```
