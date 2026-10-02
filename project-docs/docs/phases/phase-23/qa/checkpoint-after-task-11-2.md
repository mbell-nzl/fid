# Phase 23 - Independent QA Lead Checkpoint After Task 11 - Re-review 2

Reviewer: QA Lead (fresh, independent session configured as QA Lead; no prior
involvement in this phase's reviews, repairs, or reports)
Date: 2026-09-26
Checkpoint scope: `qa-plan.md` section 8, "Independent QA Lead checkpoint after
task 11" - PASS required before task 12 is dispatched
Review type: FRESH INDEPENDENT RE-REVIEW (linked)
Re-review of: `project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11.md`
(ledger review-6, commit `148c33234b2f50ee856cada663c6af9d8ab09f66`). The
Runner ledger names this re-review review-7 (`run-ledger.md:1068`).
Attestation: this review attests to no prior review. Every result below was
reproduced in this session; nothing is transcribed from a report or from the
original checkpoint. The original checkpoint was written in a session
configured as Tech Lead (ledger `run-ledger.md:1045`), so each of its findings,
areas, and coverage rows was re-derived here and is marked confirmed,
overturned, or carried.

## Verdict

BLOCK

The task-10 repair does what the original Repair Plan asked: MF-1, MF-2, MF-3,
and I-4 are each verified against their Required evidence, the change set is
exactly the three authorised files, no assertion was deleted or narrowed, and
five scratch mutation runs prove the new assertions can fail. All gates pass at
HEAD, and the tasks 12-18 feasibility sweep is clean.

The checkpoint still cannot PASS, for three reasons:

- Independent re-derivation found three matters that need protected
  specification authority (BC-1..BC-3). Each changes what VAL-009, VAL-043,
  VAL-048, VAL-050, or VAL-078 must show. Each is backed by built-binary
  evidence.
- Two bounded Must-fix evidence gaps remain. The first (MF-4) is VAL-027's
  linked-recovery leg, which the original checkpoint marked ON TRACK. The
  second (MF-5) is VAL-009's evidence-method state. The task-10 report says
  that state is unreachable and already covered by "Task 2 probes". Both
  claims are false.
- `cmd/fid` is at 567.662 s of its 600 s budget. Headroom (32.3 s) is now
  smaller than the 39.3 s spread measured on identical code (I-1).

Task 12 stays held until the Tech Lead dispositions BC-1..BC-3. The two
Must-fix repairs and a fresh independent re-review then follow in the next
unused `checkpoint-after-task-11-<R>.md`.

## Scope Reviewed

- Context: `AGENTS.md`, `CLAUDE.md`; `project-docs/docs/current/phase_brief.md`,
  `decisions_active.md`, `progress_current.md`, `tl_escalations_current.md`.
- QA plan: `project-docs/docs/phases/phase-23/qa/qa-plan.md` sections 4
  (EVID-09, 13, 15, 17, 20, 21, 23), 5, 6, 7 (QA-10, 13, 15, 17, 20, 21, 23), 8.
- Validation contract: `project-docs/docs/phases/phase-23/qa/validation-contract.md`
  rows VAL-009, 010, 015, 022, 025, 027, 029, 039, 040, 041, 043, 047, 048,
  049, 050, 078, 079, 080, 081.
- Task specs and reports: `tasks/task-10.md`, `tasks/task-10-done.md` (read in
  full at HEAD), `tasks/task-11.md`, `tasks/task-11-done.md`,
  `tasks/task-6.md` and `tasks/task-6-done.md` (VAL-025 nine-cell table,
  VAL-027), `tasks/task-2.md` and `tasks/task-2-done.md` (VAL-009 blocking
  half), and `tasks/task-12.md` to `tasks/task-18.md` in full.
- Accepted contract set: `contracts.md` sections 2, 6, 9, 10, 11;
  `public-interface.md` sections 4 and 4.1; `execution-scope.md` section 6;
  Phase 21 `contracts.md` section 6 and `public-interface.md` AuthorityGrant.
- Decisions: `decisions.md` DEC-227, DEC-235, DEC-236, DEC-239, DEC-240,
  DEC-242, DEC-243, DEC-248.
- Template: `project-docs/templates/mid-phase-review-template.md`.
- Prior review: `project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11.md`.

### Commit identity

All resolved in this session with `git rev-parse`.

| Artifact | Commit |
|---|---|
| Phase HEAD at reproduction | `908080d10c46f4ca9736f2a914035a3191a4ed8e` (matches the expected HEAD) |
| Repair: tests | `af0cbb2c1448844ec7024cd823631533de655945` (`Refs: phase-23/task-10`) |
| Repair: report | `90d921b8084e58f2415c812bf789128d9a343195` |
| Repair: report conformance | `908080d10c46f4ca9736f2a914035a3191a4ed8e` |
| Original checkpoint (review-6) | `148c33234b2f50ee856cada663c6af9d8ab09f66` |
| Task 10 implementation | `0ff2eeef66d4633309a43ae6a341aecd38def4d0` |
| Task 11 implementation / report | `c5bf0f89c69c57003ece353ef005a7596dd266e4` / `777ee47d76990b73c88ff68bbdb9e6623b794de4` |
| Planning base (R-13 baselines) | `0195897a347dfe1872db0f7be3800c5c3ccdefcf` |

`git diff --name-status 777ee47..HEAD` lists only `cmd/fid/integration/p22_fake_test.go`,
`cmd/fid/integration/p23_human_outcome_test.go` (both `//go:build integration`),
`qa/checkpoint-after-task-11.md`, and `tasks/task-10-done.md`. Every package
compiled by the ordinary suite is therefore byte-identical at `777ee47` and HEAD.

### Reproduction environment

`go version go1.26.2 linux/amd64`; `uname -a`: `Linux DESKTOP-54BJRRP
6.6.114.1-microsoft-standard-WSL2 ... x86_64 GNU/Linux`; `nproc` `6`; `AMD Ryzen 7
5700G with Radeon Graphics`; `MemTotal: 8133236 kB`; `golangci-lint has version
2.12.2`; `ripgrep 15.1.0`. Each of the three commands above terminated
`exit_code: 0`.

Working directory: `/home/mark/projects/fid` for every repository command.

Pre-flight checks:

- `/tmp` tmpfs had 3.4G free.
- No `go test` process was running. One unrelated idle shell poller from an
  earlier session was present and was left untouched.
- `/tmp/fid-integration-binary-*` count was 0.

Heavy gates ran strictly serially. Only read-only file inspection overlapped
them.

Scratch evidence lives under
`/tmp/claude-1000/-home-mark-projects-fid/5751f0e0-86ce-4419-ba62-89a77982313b/scratchpad/qa-cp2/`.
Probes and mutations ran in a detached `git worktree` at HEAD inside that
directory. The worktree was removed afterwards, and no tracked repository file
was modified.

## Summary

- Review outcome: BLOCK
- Must-fix findings: 2 (MF-4 task 6; MF-5 task 10)
- Improvement findings: 4 (I-1 and I-2 carried; I-5 and I-6 new)
- Blocker candidates: 3 (BC-1, BC-2, BC-3)
- Existing escalation references: none
- Covered validation assertions: VAL-009, VAL-010, VAL-015, VAL-022, VAL-025,
  VAL-027, VAL-029, VAL-039, VAL-040, VAL-041, VAL-043, VAL-047, VAL-048,
  VAL-049, VAL-050, VAL-078, VAL-079, VAL-080, VAL-081

## Gate And Evidence Checks

| Check | Command or Evidence | Result | Notes |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`) |
| Lint | `golangci-lint run ./...` | PASS | `0 issues.` (terminated `exit_code: 0`) |
| Ordinary race suite | `go test -race -timeout 30m -count=1 ./...` | PASS | last line `ok  	github.com/mbell-nzl/fid/internal/workflow	221.016s`; 21 packages `ok`, 0 `FAIL` lines (terminated `exit_code: 0`, wall 569 s) |
| Tagged integration gate | `go test -tags integration -race -timeout 30m -count=1 ./cmd/fid/integration/...` | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	708.114s` (terminated `exit_code: 0`, wall 708 s); integration binary directories 0 -> 0 |
| Checkpoint journeys | the ten-journey `-run '^(TestP23ReviewRecordRealProcessJourney\|...\|TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen)$'` command, byte-identical to the caller's text | PASS | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	44.807s` (terminated `exit_code: 0`); 10 top-level `--- PASS`, 9 subtests `--- PASS`, 0 `--- FAIL` |
| Named domain fixtures | `go test -race -count=1 -timeout 5m -run '^(TestCheckIndependenceRejectsClaimAndRunHistoryDirections\|TestTerminalFindingDispositionAuthorityMatrix\|TestChangesRequiredReviewChargesExactlyOneReviewCycle)$' ./internal/review/ -v` | PASS | `ok  	github.com/mbell-nzl/fid/internal/review	1.150s` (terminated `exit_code: 0`); 3 top-level plus 12 subtests PASS |
| Rebuild benchmark | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` | PASS | `ok  	github.com/mbell-nzl/fid/internal/journal	24.790s`; in-run `22136 ms/rebuild` against the 30000 ms bar (terminated `exit_code: 0`) |
| Build | `go build -o <scratch>/fid-qa-checkpoint-2 ./cmd/fid` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`); 16825752-byte binary |
| Capability probe (VAL-047 first half) | `fid-qa-checkpoint-2 --json capabilities` in a fresh `git init` + `fid init --prefix QAC` project | PASS | `"capability_id":"attention.v1","limits":{"human_outcome_resolution":false,"state":true}` (terminated `exit_code: 0`). The same probe run from the repository root, against that checkout's own pre-existing `.fid` workspace, returned `internal-error` "capability discovery failed". Its cause was not investigated, and it is not used as evidence |
| Report lint, task 10 | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-10-done.md project-docs/docs/phases/phase-23/tasks/task-10.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`) |
| Report lint, task 11 | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-11-done.md project-docs/docs/phases/phase-23/tasks/task-11.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`) |
| Task-10 exact-command audit | `task-10-done.md` `## Exact quality-gate results` against `task-10.md` `## Exact Commands To Run` | PASS | The spec lists five commands and declares no conditional or expected-non-zero command. The report records each command twice, as the original rows and the QA-repair rows. All ten rows are `PASS` with a verbatim token and `(terminated exit_code: 0)`. The repair rows are byte-identical without `-count=1`; `internal/cli`, `internal/config`, and `internal/validation` replayed `(cached)` on unchanged inputs, and this session's `-count=1` runs cover every package fresh. The three DEC-227 package times are recorded (`cmd/fid` 564.337 s, `internal/workflow` 225.996 s, integration 693.638 s). |
| Report claim accuracy | `task-10-done.md` corrected VAL-043, VAL-025, VAL-009, VAL-010 entries against the cited lines | FAIL (partial) | VAL-043 and VAL-025 match the code line by line. The VAL-009 entry contains two false claims (MF-5), and the VAL-010 entry has one wrong citation (I-6). |
| Scratch probes and mutations (QA evidence only, detached worktree at `908080d`) | `GOFLAGS=-buildvcs=false go test -tags integration -race -timeout 30m -count=1 [-v] -run <selector> ./cmd/fid/integration/...` in the worktree | recorded as observed | Observed results:<br>- First probe attempt without `GOFLAGS`: `build integration fid binary: exit status 1` / `error obtaining VCS status: exit status 128` (terminated `exit_code: 1`). This is a linked-worktree VCS stamping limitation, and no test ran.<br>- Rerun with `-buildvcs=false`: probe `ok ... 7.965s` (terminated `exit_code: 0`).<br>- Mutations: M1 `:50` FAIL, M2 `:126` FAIL, M2b `:159` FAIL, M3 `:87` FAIL, M4 `:86` FAIL, each terminated `exit_code: 1` as the mutation intends.<br>Every mutation was reverted with `git checkout --`; the worktree was then removed. |
| Closed grant allowlists | `sortedUniqueCapabilities` (`internal/workflow/activation.go:832`) vs `const actions` (`internal/journal/current_event_value_relations.go:63`) | PASS | 30 and 30 names; `diff` of the sorted sets terminated `exit_code: 0` with no output |

## Part (a) - Repair Verification

### Change set and diff audit

- `git diff --name-status 148c332..HEAD` lists exactly three files:
  `cmd/fid/integration/p22_fake_test.go`,
  `cmd/fid/integration/p23_human_outcome_test.go`, and
  `project-docs/docs/phases/phase-23/tasks/task-10-done.md`.
- `p22_fake_test.go` changes by +7/-0 lines: the `CoexistingReviewObligation`
  field and comment (`:294-297`), and one guarded append in the Obligation
  construction (`:381-383`). The existing replace-in-place branch (`:377-379`)
  and the `OBL-p22-corpus` ID are unchanged.
- Existing `ReviewObligation: true` consumers never set the new field, so the
  zero value skips the new branch and their Obligation lists are unchanged:
  `p23_freshness_supersession_test.go:13`,
  `p23_review_evidence_budget_test.go:12`, `:68`, and
  `p23_human_outcome_test.go:312`. All of them pass inside this session's full
  tagged gate.
- Every assertion-bearing line deleted from `p23_human_outcome_test.go` was
  accounted for with `git diff 148c332..HEAD | grep '^-'`. Each has an equal or
  stronger replacement:
  - the three `assertP23AuthorityMissing(..., before)` and `(..., 0)` calls
    became four-argument before/after checks;
  - the two `assertP23OutcomeNotPermitted(..., "gate-evidence")` calls became
    the bracketed `assertP23OutcomeLinkNotPermitted`, which calls the same
    assertion;
  - `p22AcceptedCall(run.observe-budget)` plus `evidence.RunGate` and its
    receipt check moved into `p23RunQualifyingGateSession` (`:723-743`), with
    the observe asserted and the receipt check retained at `:736-739`;
  - the missing-capability request previously named a nonexistent Issue
    (`ISS-p23-human-outcome`) and now uses a real one;
  - the removed comment asserting that the unfinalized path shares code with the
    `unregistered` case was replaced by an executed request.

  No assertion was deleted or narrowed.

### MF-1 (VAL-009, EVID-09, R-07b) - RESOLVED

The Required evidence is met.

- **All three types on one Envelope.** `assertP23ThreeObligationTypesCoexist`
  (`p23_human_outcome_test.go:586-604`, called at `:28`) reads
  `envelope_obligations` and requires exactly `OBL-p22-corpus:gate-evidence`,
  `OBL-p23-review:review`, and `OBL-human-release`/`OBL-human-security:human-outcome`.
- **Denied while any one is unsatisfied.** Each denial stage checks journal
  length and the whole-database projection (`p22FaultProjection`) before and
  after:
  - Stage 1 (`:107`): `envelope.complete` exits 4 `obligation-unsatisfied`,
    naming exactly `OBL-human-security`, `OBL-p22-corpus`, and `OBL-p23-review`
    (`:611-642`). The terminal `qa_pass` is denied `transition-not-permitted`
    (`:649-689`).
  - Stage 2 (`:118`): the review Obligation is now satisfied, and its reason
    disappears from the denial.
  - Stage 3 (`:139-146`): only `OBL-human-security` is unsatisfied.
    `envelope.complete` returns `human-outcome-required` and `qa_pass` is
    denied.
- **All three satisfied.** `p23EvidenceEvaluate` reports all four Obligations
  `satisfied` (`:173-179`). Then `qa_pass` (`:184`) and `envelope.complete`
  (`:186`) both succeed.
- **The review Obligation is load-bearing (scratch mutation M2b).** The mutation
  deletes `:114-118`, `:131`, `:138-143`, `:146`, `:156-160`, and `:173-179`,
  so the review Obligation is never satisfied. It fails at the permitted action,
  on mutated line 159 (original `:184`):
  `transition-not-permitted ... "unsatisfied":["obligation-evidence-absent"]`
  (terminated `exit_code: 1`).
- **Mutation M2 (narrower deletion).** It fails earlier, at mutated `:126`
  (`sql: no rows in result set`): with the review unsatisfied, the gate link
  raises no readiness Stop. This matches DEC-224, where the readiness Stop is
  raised only when the Human Obligation is the sole unsatisfied one.

### MF-2 (VAL-043) - RESOLVED

The Required evidence is met. All five requests run through the built binary:

1. Agent actor type: `actor-denied` (`:41-43`).
2. Unfinalized Activation, run twice in isolated workspaces (`:49-51`,
   `p23UnfinalizedActivationAttempts` at `:527-533`): once prepared with staged
   components, and once committed but not finalized. In both, the sole human
   Actor holds `human-outcome.record`, and a real Issue exists.
3. Missing capability against a finalized Activation (`:56-57`).
4. Standalone attribution `human:unregistered` (`:61-63`).
5. The valid Kit-mode request (`:67-81`).

Requests 1-4 each assert exit 5, empty stdout, `authority-missing`, and
`after == before`. `assertP23AuthorityMissing` (`:403-418`) now compares
`after != before` first, at `:405-407`. Before/After values are captured
immediately around the request itself (`:574-580`). The valid case asserts
exactly one appended Event (`:70-72`) and
`provenance.activation_id == a.activation["id"]` (`:78-81`).

Mutation M1 (`:531`, `p23ActivationCommitted` changed to
`p23ActivationFinalized`) failed at `:50` with
`authority-missing denial appended Event(s): 12 -> 13`, and the record returned
`"ok":true` (terminated `exit_code: 1`). This proves two things:

- finalization alone decides the unfinalized denial;
- the before/after bracket is consumed.

That same success happened in a workspace with no Envelope at all, which is
evidence for BC-2.

### MF-3 (VAL-025, DEC-243) - RESOLVED

The Required evidence is met.

- `assertP23OutcomeLinkNotPermitted` (`:438-449`) brackets every real-producer
  mismatch with `p21JournalEventCount` equality and a whole-database
  `p22FaultProjection` `reflect.DeepEqual`.
- It calls `assertP23OutcomeNotPermitted` (`:451-466`), which checks exit 4,
  `outcome-not-permitted`, and the exact `details.obligation_type`.
- Call sites:
  - release Outcome to `OBL-p22-corpus` expects `"gate-evidence"` (`:86`);
  - release Outcome to `OBL-p23-review` expects `"review"` (`:87`);
  - the security Outcome repeats both mismatches after completion (`:192-201`).
- Mutation M3 changed `:87`'s expectation to `"gate-evidence"`. It failed with
  `denial obligation_type = "review", want "gate-evidence"`, so the real
  producer classifies the review cell.
- Mutation M4 injected a Human Outcome recording after the link, at `:441`. It
  failed with `release-vs-gate outcome-not-permitted denial appended Event(s):
  21 -> 22`, so the bracket bites. Both terminated `exit_code: 1`.
- The report's corrected nine-cell reconciliation (`task-10-done.md:275-295`)
  matches `task-6-done.md:44-54` cell for cell:
  - receipt->gate-evidence and review->review remain Task 6 durable links;
  - the four `evidence-type-not-accepted` cells remain Task 6
    production-classifier evidence, which DEC-243 permits;
  - the two Human mismatches and the positive Human link
    (`:88-102`, evaluation `satisfied`) are now real-producer evidence.
- The self-contradictory sentence is preserved only inside the entry marked
  SUPERSEDED (`task-10-done.md:202-212`). The corrected entry names it as
  self-contradictory (`:289-291`). It is fixed.

### I-4 (VAL-010) - RESOLVED

With linked review Evidence and linked Human Outcome Evidence both present:

- the linked review evaluation is `satisfied` and `not-applicable`
  (`:114-117`);
- the linked security Outcome evaluation is `not-applicable` (`:169-172`);
- in the all-satisfied evaluation, `OBL-p22-corpus` is `available` and the
  other three are `not-applicable` (`:173-179`).

### Report-to-code line audit (corrected entries)

| Entry | Result | Detail |
|---|---|---|
| VAL-043 (corrected) | matches | All cited lines match the code: `:38-43`, `:45-51`, `:527`, `:53-57`, `:59-63`, `:65-`, `:403`, `:540-580`. |
| VAL-025 (corrected) | matches | `:438`, `:86`, `:87`, and `:88-100` match. The loop cited as `:193-201` actually begins at `:192`, which is immaterial. |
| VAL-009 (corrected) | two false claims | The cited lines `:28`, `:586`, `:611`, `:107`, `:118`, `:139-146`, `:184`, and `:186` match. The entry makes two claims contradicted by evidence (see MF-5): that "the Task 2 probes remain the unsatisfied-review-with-satisfied-gate evidence", and that the isolated "review-only" state is "not reachable". It also says `qa_pass` denials name "exactly the still-unsatisfied Obligations"; they compare a multiset of reason codes (`:663-682`), not Obligation IDs (I-6). |
| VAL-010 (corrected) | one wrong citation | The security Outcome evaluation is cited as `:158-161`; it is at `:169-172` (I-6). |

## Part (b) - Original Findings, Areas, And Coverage, Re-derived

### Original Must-fix findings

All three were valid at `148c332` and are resolved at HEAD.

- **MF-1:** CONFIRMED. At `148c332`, `p22_fake_test.go:374` replaced
  `OBL-p22-corpus` in place. `git grep '"type": "review"' 148c332` returns only
  that site, and the journey's action list lacked `review.record`.
- **MF-2:** CONFIRMED. At `148c332`, `assertP23AuthorityMissing(t, response,
  before int)` never read `before`, and one call site passed the literal `0`.
- **MF-3:** CONFIRMED. At `148c332`, both mismatches targeted `OBL-p22-corpus`
  with `"gate-evidence"`, and neither was bracketed.

### Original Improvement findings

- **I-1 (`cmd/fid` time headroom):** CONFIRMED and carried, and stronger than
  first stated. Fresh measurement: 567.662 s, with 32.3 s (5.4%) of headroom.
  The ordinary-compiled code has not changed since `777ee47`, yet four runs of
  it measured 546.059 s, 585.401 s, 564.337 s, and 567.662 s. That 39.3 s
  spread is larger than the remaining headroom, so VAL-079 can breach through
  machine variance alone.
- **I-2 (task 11 `stop_conditions.go` scope):** CONFIRMED and carried.
  - `c5bf0f8` adds 43 lines to `internal/workflow/stop_conditions.go`.
  - That file was introduced by `f973bca` (`Refs: phase-22/task-2`) and is
    reached from `qualifyingFact` (`internal/workflow/stop.go:618-634`), which
    is in scope.
  - `task-11.md` names `resolve.go`, `evaluation.go`, "attention readers", and
    `stop.go`, not `stop_conditions.go`.
  - The deviation is disclosed at `task-11-done.md:257-270`.
  - The change is additive and decision-free. The Tech Lead should ratify it.
- **I-3 (VAL-078 route per Condition type):** CONFIRMED, but materially
  UNDER-STATED. The original called the atomic-only route "behaviourally
  correct". Built-binary probes show the `stop resolve` route rejects even the
  qualifying fact for `human-outcome-recorded`, and a Stop that carries that
  Condition can become permanently unresolvable. Superseded by BC-3.
- **I-4:** CONFIRMED and RESOLVED (see Part (a)).

### The seven reproduced areas

| # | Area | This session | Outcome versus original |
|---|---|---|---|
| 1 | VAL-015 independence | Built binary: `TestP23ReviewRecordRealProcessJourney` (`p23_review_record_test.go:24-36` Claim rejection with exact Claim and Actor and no Event; `:98-110` same-Reviewer post-rework Review). Isolated durable-history query: `TestCheckIndependenceRejectsClaimAndRunHistoryDirections` (`internal/review/review_test.go:263-303`). Both PASS. | Confirmed ON TRACK |
| 2 | VAL-022 resolution authority | 6 built-binary subtests PASS, plus `TestTerminalFindingDispositionAuthorityMatrix` with 12 subtests (6 per terminal action) PASS | Confirmed ON TRACK |
| 3 | VAL-027 / VAL-029 | VAL-029 is fully reproduced (`p23_review_evidence_budget_test.go:66-128`: atomic exhausting Event with budget, Stop, and Attention; non-committing `claim.acquire`/`run.start` denials). VAL-027's recording, replay, and rebuild legs are reproduced. Its **linked-recovery leg** (validation contract: "after linked recovery at task 6"; `task-6.md:194`) has no test anywhere. No test records a `changes-required` Review and then performs `run.recover`. The only files containing both `review-cycles` and `recover` are the fixture builder `p22_fake_test.go` and `p23_review_evidence_budget_test.go`, where the match is the `actor-recovery` Actor name. | VAL-029 confirmed. VAL-027 **overturned** to NOT ON TRACK; the original was under-stated. MF-4. |
| 4 | VAL-039/040/041 | `TestP23DecisionReferenceProtectedLifecycleRealProcess` PASS. Per-class denial with an unsatisfied Obligation also present (`:29-30`); rejected, deferred, and wrong-authority references plus the four satisfaction states (`:32-70`); the unaffected action proceeds (`:79`); one Stop resolved per satisfying pair (`:66-70`, `:85-86`) | Confirmed ON TRACK |
| 5 | VAL-043 provenance | MF-2 resolved (Part (a)). New: BC-2, an unresolved Envelope-grant contradiction in the accepted interface. | Confirmed ON TRACK against the row's five requests; BC-2 raised |
| 6 | VAL-048/049/050 | Both task-11 journeys PASS. VAL-049 is reproduced (`p23_attention_resolution_test.go:83-114`). VAL-048's "and reject" clause and VAL-078's "rejects with no Event" have no committed observation for `human-outcome-recorded`/`decision-satisfied`. A matching-Condition Stop that is not a producer code is never resolved (probe). | VAL-049 confirmed. VAL-048 and VAL-050 **overturned** to pending BC-3. |
| 7 | VAL-009 all-satisfied | Resolved (MF-1). New: MF-5 (the evidence-method state) and BC-1 (guarded-action semantics). | All-satisfied state confirmed; MF-5 and BC-1 raised |

### Original coverage table

- **VAL-047 (first half):** confirmed, using this session's own binary probe.
- **VAL-080:** confirmed. `internal/journal/current_event_value_catalog.go:59-70`
  holds ten rows; the six Phase 22 rows are unchanged, and `review-recorded`
  has the bare four-field form.
- **VAL-081:** the conclusion is confirmed, but the original's stated evidence
  is wrong. It said the implementation commits list neither
  `current_event_value_relations.go` nor the types file.
  `git show --name-only 0ff2eee` does list `current_event_value_relations.go`:
  task 10 changed the `ActorProfile` allowlist at `:63`, as `task-10.md`
  requires. `case "Condition"` is untouched, and `c5bf0f8` does not touch the
  file.
- **VAL-078:** overturned from "ON TRACK with a note" to NOT ON TRACK, pending
  BC-3.
- **VAL-027:** overturned (MF-4).
- **VAL-079:** confirmed ON TRACK with the fresh values in Part (c).

## Part (c) - Tasks 12-18 Fresh Feasibility Sweep And R-13

### Commands and fixture assumptions (re-validated at HEAD `908080d`)

| Task | Command or assumption | Validation method | Result | Evidence |
|---|---|---|---|---|
| 12, 15, 16 | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` | executed | PASS | `22136 ms/rebuild`; `ok ... internal/journal 24.790s` (terminated `exit_code: 0`) |
| 12, 15 | `rg -n --glob '!*_test.go' "requirement\.synced\|last_synced_hash" cmd internal` (declared expected `exit_code: 1` at task close) | executed | PASS (syntax and scope) | Terminated `exit_code: 0`. Matches only `internal/export/archive_semantics_catalog.go:91`, `internal/requirements/sync.go` (6 lines), and `internal/journal/current_event_catalog.go:119`. All three are in task 12's scope, so the declared `exit_code: 1` is reachable. |
| 15 | `rg -n --glob '!*_test.go' "task_completions\|ac_verifications" cmd internal` | executed | PASS (syntax and scope) | Terminated `exit_code: 0`. Matches only `internal/maintenance/{operational,errors,warnings}.go` and `internal/issues/{detail,transition_rules}.go`, all in task 15's list. |
| 15 | `rg -n --glob '!*_test.go' "has_completion_report\|concerns_addressed" cmd internal` | executed | PASS | Terminated `exit_code: 0`; `internal/workflow/rules.go:7-8` only |
| 12-18 | `## Exact Commands To Run` and declarations | spec inspection plus `grep` | PASS | Every spec carries `-timeout 30m` on both long gates. Task 12 has 7 commands, with the benchmark and `rg` listed as conditional and the `rg` under `## Expected Non-Zero Commands`. Task 15 has 9, with the benchmark conditional and three `rg` declared `exit_code: 1`. Task 16 has 6, with the benchmark conditional. Tasks 13, 14, 17, and 18 have 5 each and state `- None.` |
| 12-18 | `## First Failing Test` names are new | `rg -n "func <name>\(" cmd internal`, run for all seven names | PASS | Each terminated `exit_code: 1 (no output)`; no match is the pass condition. |
| 12 | fence map, guard, update switch, JSON switch | file-state inspection | PASS | `internal/journal/project_effects.go:12-15` holds exactly `requirement.referenced` and `requirement.linked`, with the guard at `:18-20`. `current_event_updates.go:14` and `:63` hold the switches, with `judgment_concerns` at `:41` and `:86` and no `requirement_refs` entry yet. `current_event_json.go:14` holds its switch. |
| 12 | interface rows; `internal/links`, `internal/requirements`, CLI files | file-state inspection | PASS | `public-interface.md:116`, `:118`, `:180-183`. `internal/links/{registry,add,list,errors,types}.go` and `internal/requirements/{parser,sync,attribution,errors,types}.go` are present, as are `cmd/fid/{requirement,link}.go`. `requirement-not-found` is at `internal/cli/codes.go:29`. |
| 13 | rename targets | `rg` | PASS | `reasonCodeDefinitions` is declared at `internal/cli/codes.go:232` and used at `:396`, `:407`, and `internal/cli/envelope.go:38`. `reasonCodeExpectation` is at `codes_test.go:30`. The `reasonCodes` parameter at `cmd/fid/evidence_test.go:465` is correctly left alone. All eleven `cmd/fid` leaf files are present. |
| 14 | precedent and targets | `ls` | PASS | `p22_human_fixtures_test.go`, the Phase 22 `qa/human-fixtures.md`, and the resume files are present. The Phase 23 `human-fixtures.md` and `p23_human_fixtures_test.go` are correctly absent. `resume.go:173` still has the `Findings{Available:false}` placeholder and a single-element `Obligations`. |
| 15 | deletion/edit targets and precedents | `test -f` | PASS | All 29 named files are present. `0014` and `0016` are present; `0018` is absent. `judgment_concerns` appears in `current_event_columns.go:20`, `current_event_json.go:79,81`, `current_event_catalog.go:29-30`, and `archive_semantics_catalog.go:41-42`. |
| 16 | archive files and `schema-unsupported` path | inspection | PASS | `internal/export/{export,archive_profiles}.go` and `cmd/fid/export.go` are present. `schema-unsupported` is emitted at `import.go:83` and `archive_validation.go:127,131,241`. |
| 17 | false rows to flip | inspection plus built-binary probe | PASS | `internal/workflow/capabilities.go:105,110,112`, `internal/export/archive_profiles.go:5`, and the binary all report `"human_outcome_resolution":false`. |
| 18 | wire oracle, authority, grant set | inspection | PASS | `p22_closed_values_test.go:25` has six Condition types and `:108` has the `50`/`90` counts. The catalog at `:59-70` has ten rows. The grant set is 30 = 30, identical in both copies. `cmd/fid/integration/testdata` and the Phase 22 `task-18.md` precedent are present. |

Nothing in the sweep blocks a later task mechanically. Some forward specs depend
on the outcomes of BC-1..BC-3, as noted under each candidate.

### R-13 / DEC-227 / VAL-079 measurements (fresh, serial, `-count=1`)

| Package | Measured now | Budget | Headroom | Planning base | Task 11 report | Original checkpoint (`777ee47`) | Repair report |
|---|---|---|---|---|---|---|---|
| `cmd/fid/integration` (tagged) | **708.114 s** | 900 s | 191.9 s (21.3%) | 560.243 s | 700.168 s | 688.252 s | 693.638 s |
| `cmd/fid` (ordinary) | **567.662 s** | 600 s | **32.3 s (5.4%)** | 471.955 s | 546.059 s | 585.401 s | 564.337 s |
| `internal/workflow` (ordinary) | **221.016 s** | 600 s | 379.0 s (63.2%) | 458.353 s | 213.335 s | 228.577 s | 225.996 s |

Next slowest ordinary packages: `internal/export` 311.392 s,
`internal/issues` 228.178 s, and `internal/journal` 227.503 s.

VAL-079 is not breached at HEAD. However, `cmd/fid` shows 546.059 s to
585.401 s (a 39.3 s spread) on byte-identical code, which is larger than its
remaining 32.3 s of headroom (I-1).

Two consequences follow:

- The repair touched only build-tagged files, so it cannot have moved
  `cmd/fid`.
- The tagged gate at 708 s is 108 s past Go's implicit 600 s per-binary
  default. DEC-227's `-timeout 30m` is therefore load-bearing for tasks 12-18.

## Required Judgments

### 1. The multi-grant `gate.run` concern - classification (i)

It is an ordinary concern for Tech Lead disposition. Precisely, it is a latent
Phase 21 conformance defect.

**What the code does.**

- `internal/evidence/execution.go:381`
  (`validateGateRunReferences`) calls `QueryRow` on
  `envelope_authority_grants` with no actor filter and no `ORDER BY`. It then
  requires that single row's selector to name the running Actor.
- Grants are stored by ordinal in actor-ID order (`envelope.go:857`), so only
  the first-sorted Actor can run a gate.
- `git blame` attributes the line to `b69499116e70d4c88423a30daa2010a4ff221f14`
  (2026-08-30, `Refs: phase-21/task-15`), before the Phase 23 planning base.
- Every other post-creation authorization iterates all grants
  (`internal/workflow/authority.go:36-74`, via `loadEnvelopeGrants` at
  `envelope.go:929-944`).

**What the contract says.** The accepted Phase 21 contract fixes the behaviour
uniquely:

- `contracts.md` section 6: "requires an exact action, capability, Issue,
  Envelope, and revision Grant match. Missing is deny."
- `public-interface.md`: `authority_grants:AuthorityGrant[1..64]`, and "grants
  sort by actor ID".

**Why it is not the other two classes.**

- It is not (ii), a task-10 Must-fix. The fix is product code outside the
  test-only repair boundary.
- It is not (iii), a blocker candidate. It needs no protected decision; the
  failure is fail-closed (it denies legitimate Actors and admits no
  illegitimate one); and `execution.go` already lies in tasks 2 and 6 file
  scopes (`task-2.md:62`, `task-6.md:74`).

**Requested action.** The Tech Lead assigns an owner under the DEC-175/DEC-194
pattern before task 14's VAL-071 Human fixture or task 18's corpus needs a
non-first-grant Actor to run a gate. Recorded as I-5.

**Effect of the workaround on evidence.** It does not weaken VAL-009, VAL-043,
or VAL-025:

- The gate-evidence Obligation's satisfaction rule is `fid-executed`
  independence and a `fid-executed` trust source, with no running-Actor
  constraint. The Receipt is produced by production `evidence.RunGate` against
  the Envelope's Registration and is linked by the builder through the real
  `evidence.link` leaf, which accepted it.
- VAL-043 and VAL-025 do not involve `gate.run`.

**Disclosed limitation.** The gate step runs `evidence.RunGate` in-process
(`p23_human_outcome_test.go:736`) rather than through `fid gate run`. That was
already true at `148c332`. Every guarded-action and linking observation goes
through the built binary.

### 2. VAL-009 "monotone" argument

The report's premise is false:

- Mutation M2b reached the state in which the review Obligation is the only
  unsatisfied one, with the gate and both human Obligations satisfied, on one
  Envelope.
- The contract itself makes satisfaction non-monotone: a tree change stales
  review and gate proof (VAL-031), and a later Outcome supersedes an earlier one
  (VAL-033).

**Evidence method.** The Task 10 half of VAL-009's evidence method is met: the
qualifying all-satisfied state with all three types.

**Both-states rule.** The rule is met on one Envelope:

- Stage 3 is a single-unsatisfied state in which both guarded actions are
  denied.
- The all-satisfied state permits both.

**"Denied while any one of the three is unsatisfied."** This is met on its
natural reading, "while at least one remains unsatisfied":

- every pre-completion stage is denied;
- `envelope.complete` names each unsatisfied Obligation;
- the stage 1->2 and stage 2->3 differentials isolate the review and gate
  contributions.

**What is not met.** VAL-009's own evidence method also names guarded-action
probes "with a satisfied gate-evidence Obligation and an unsatisfied review
Obligation". No such probe exists in the repository:

- `task-2-done.md:34` cites
  `TestEnvelopeActivationEvaluationIncludesEveryObligation`
  (`internal/workflow/envelope_test.go:115-126`), which only checks ordered
  activation evaluation.
- The repaired journey satisfies the review before the gate.

The report's reliance on "Task 2 probes" is therefore unsupported. This is
MF-5. The `qa_pass` expectation in that state depends on BC-1.

### 3. Additional defects

Recorded as MF-4, MF-5, I-5, I-6, and BC-1..BC-3 below.

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-009 | task 2 blocking half; task 10 completion | NOT ON TRACK | The all-satisfied state is evidenced (MF-1 resolved). The evidence-method state (satisfied gate, unsatisfied review) is unprobed and the report claim about it is false (MF-5). The meaning of "applicable Obligation" is unresolved (BC-1). |
| VAL-010 | task 2; task 10 re-verification | ON TRACK | `:114-117`, `:169-172`, `:173-179` (I-4 resolved) |
| VAL-015 | task 4 | ON TRACK | Area 1 |
| VAL-022 | task 5 | ON TRACK | Area 2 |
| VAL-025 | tasks 6 and 10 | ON TRACK | The Task 10 half is closed through the real producer (`:86-102`, `:192-201`), and the nine cells reconcile with `task-6-done.md:44-54` (MF-3 resolved). Final PASS is decided at final review. |
| VAL-027 | task 6; import leg task 16 | NOT ON TRACK | The linked-recovery leg is unevidenced (MF-4). QA probe: after a `changes-required` Review, `review-cycles` samples were count 1 / sum 1, and stayed 1 / 1 after `run.recover` (`RUN-2a2e59` -> `RUN-d43c4b`) and after `fid rebuild`. The behaviour is correct; the committed evidence is missing. |
| VAL-029 | task 6 | ON TRACK | Area 3 |
| VAL-039 / 040 / 041 | task 9 | ON TRACK | Area 4 |
| VAL-043 | task 10 | ON TRACK | All five requests are evidenced (MF-2 resolved). The separate Envelope-grant clause in `public-interface.md` 4.1 is BC-2. |
| VAL-047 (first half) | task 11 | ON TRACK | Built-binary probe returns `false` |
| VAL-048 | task 11 | NOT ON TRACK | The Stop stays open in the wrong-class and wrong-result states. The "reject" clause has no committed observation, and the resolution route is unresolved (BC-3). |
| VAL-049 | task 11 | ON TRACK | `p23_attention_resolution_test.go:83-114` |
| VAL-050 | task 11 | NOT ON TRACK | The unrelated-Condition Stop stays open (`:176`, `:185`). However, a Stop whose Condition matches but whose code is not a producer code is never resolved (BC-3). |
| VAL-078 | task 11 | NOT ON TRACK | `review-recorded` and `finding-resolved` are evidenced (`:52-81`). For `human-outcome-recorded` and `decision-satisfied`, `stop resolve` rejects the qualifying fact (BC-3). |
| VAL-079 | tasks 1-18, phase end | ON TRACK | 708.114 / 567.662 / 221.016 s against 900 / 600 / 600 s; see I-1 |
| VAL-080 | tasks 9, 10, 11 | ON TRACK | Catalog `:59-70` |
| VAL-081 | task 9, re-confirmed 10 and 11 | ON TRACK | `case "Condition"` is untouched; this session's full ordinary gate ran green |

## Per-Task Findings

### Task 6 - Review evidence, cycles, exhaustion

#### Must-fix

- **MF-4 (VAL-027): the linked-recovery leg has no evidence.**
  - The validation contract's evidence method is "Charge inspection after
    recording, after replay, and after linked recovery at task 6".
  - Its pass criterion is "by zero on replay, recovery, and import".
  - `task-6.md:194` is Task 6's own acceptance criterion; `contracts.md`
    section 6 carries the invariant.
  - Nothing in the repository recovers a Run after a `review-cycles` charge.
    The only `run.recover` tests with review Claims are
    `p23_claim_mode_test.go:11-66`, which record no Review.
    `TestP23ReviewCyclesRealProcessJourney` covers recording, rebuild, and
    replay only. `task-6-done.md` VAL-027 claims no recovery leg.
  - Window 4-6 (`review-4.md:68`) and the original checkpoint both marked the
    row ON TRACK.
  - This session's scratch probe shows the product behaves correctly, so this
    is an evidence-only gap.

#### Improvement

- None.

### Task 10 - Provenanced Human Outcomes

#### Must-fix

- **MF-5 (VAL-009): the evidence-method state is unprobed, and the report
  claims it is covered.**
  - VAL-009's evidence method requires guarded-action probes "with a satisfied
    gate-evidence Obligation and an unsatisfied review Obligation". None exists
    (Required Judgment 2).
  - `task-10-done.md:311-315` claims that "an isolated 'gate-evidence-only' or
    'review-only' unsatisfied state is not reachable on one Envelope" because
    satisfaction is monotone. Mutation M2b and VAL-031/033 contradict this.
  - The same entry claims that "the Task 2 probes remain the
    unsatisfied-review-with-satisfied-gate evidence". `task-2-done.md:34`
    contradicts this.
  - These are deterministic claims contradicted by evidence: the same defect
    class as the original MF-1..MF-3.
  - The state is constructible inside the existing task-10 fixture by
    satisfying the gate before the review.

#### Improvement

- **I-6: completion-report citation drift and clarity.**
  - The VAL-010 corrected entry cites `:158-161`; the code is at `:169-172`.
  - The non-superseded entries still cite pre-repair lines:
    - VAL-044 "line 70-71" is now `:86-87`;
    - VAL-045 "76,81,84" is now `:92`, `:97`, `:100`;
    - VAL-046 "190-250"/"252-335" is now `:234-294`/`:310-379`, and
      "88-103" is now `:124-131`;
    - VAL-081 "468-485 ... line 173" is now `:762-773` and `:217`.
  - "gate-evidence and review at stages 1-2" is imprecise: the review blocks
    at stage 1 only.
  - The `qa_pass` stages compare reason-code multisets, not Obligation IDs.
  - The original-session `-count=1` paragraph (`:410-419`) sits after the
    repair rows and reads as if it applied to them.
  - Fix direction: refresh the citations in the same continuation as MF-5.

### Task 11 - Attention resolution

#### Must-fix

- None in bounded scope. The route gap is BC-3, because its repair needs a
  qualifying rule the Tech Lead must fix.

#### Improvement

- **I-2 (carried):** ratify the `stop_conditions.go` scope reading in
  `decisions.md`.

### Cross-task

#### Improvement

- **I-1 (carried, strengthened): `cmd/fid` 567.662 s, 32.3 s headroom, smaller
  than the 39.3 s same-code spread.**
  - Tasks 13 and 14 add `cmd/fid` tests.
  - DEC-227 assigns a breach to the task whose new journeys caused it, and says
    a pre-existing cause "is escalated, not absorbed".
  - A breach driven by variance on this margin will be hard to attribute.
  - The Tech Lead should decide before tasks 13 and 14 dispatch whether to
    pre-assign `cmd/fid` speed headroom work. The budget itself is protected;
    this changes no threshold.
- **I-5: multi-grant `gate.run` (Required Judgment 1).** A Tech Lead owner
  assignment is required.

## Prior Finding Resolution

Prior review: `project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11.md`
(review-6, `148c332`). It is cited only; it is not edited, appended to, or
rewritten.

| Prior finding | Owning task | Result | Fresh evidence |
|---|---|---|---|
| MF-1 (VAL-009 all-satisfied without review) | task-10 | RESOLVED | Fixture `p22_fake_test.go:381-383`; `p23_human_outcome_test.go:28`, `:107`, `:118`, `:139-146`, `:173-186`. The ten-journey run PASSed (`44.807s`) and the full tagged gate PASSed (`708.114s`). Mutation M2b proves the review is load-bearing. |
| MF-2 (VAL-043 requests and "no Event") | task-10 | RESOLVED | Five requests at `:41-81`; helper `:403-418`. Mutation M1 failed at `:50` with `12 -> 13`. |
| MF-3 (VAL-025 review mismatch and brackets) | task-10 | RESOLVED | `:86`, `:87`, `:192-201`; helper `:438-449`. Mutations M3 (`obligation_type = "review"`) and M4 (`21 -> 22`). Nine-cell reconciliation checked against `task-6-done.md:44-54`. |
| I-1 (`cmd/fid` headroom) | cross-task | STILL OPEN | 567.662 s against 600 s, with a 39.3 s same-code spread; carried to the Tech Lead. |
| I-2 (task 11 file-scope reading) | task-11 | STILL OPEN | `c5bf0f8` stat; `task-11-done.md:257-270`; Tech Lead ratification pending. |
| I-3 (VAL-078 route per type) | task-11 | PROTECTED BLOCKER CANDIDATE | Superseded by BC-3 on built-binary probe evidence; the original under-stated it. |
| I-4 (VAL-010 with review and Human Outcome evidence) | task-10 | RESOLVED | `:114-117`, `:169-172`, `:173-179` |

Additional corrections to the original checkpoint, recorded without editing it:

- Area 3 / VAL-027 omitted the linked-recovery leg (now MF-4).
- The VAL-081 row's stated evidence is factually wrong; the conclusion stands.
- Areas 5-7 did not examine the specification points now raised as
  BC-1..BC-3.

## Repair Plan

The two Must-fix rows are bounded and decision-free. They are listed in repair
order. They proceed only after the Tech Lead or TL-Lite dispositions the
blocker candidates. MF-4 is independent of all three candidates and may be
dispatched as soon as the Tech Lead releases the stop. MF-5's `qa_pass`
expectation waits on BC-1.

| Finding | Owning task | Affected assertion | Exact repair boundary | Required evidence | Fresh re-review scope |
|---|---|---|---|---|---|
| MF-4 | task-6 | VAL-027 | Test-only. Extend `cmd/fid/integration/p23_review_evidence_budget_test.go` (the cycle journey or a sibling), reusing the unchanged helpers `p22Observe`, `p22Expire`, and `p22Recover` (`p22_fake_test.go:643-667`) and the review-Run recovery pattern (`p23_claim_mode_test.go:53-66`). Update `task-6-done.md`'s VAL-027 entry with a dated repair section. If the observation exposes a double charge, the causal repair stays inside `task-6.md:80`'s `internal/workflow/{budget,run_budget,recovery_budget}.go`. No contract, schema, Event, or gate change. | A built-binary journey records a `changes-required` Review in a review Run, giving exactly one `review-cycles` `budget_samples` row keyed on the Review ID with `consumed:1`. That Run is recovered through `run.recover` into a linked successor. The `review-cycles` sample count and consumed value are unchanged after recovery and after `fid rebuild`, and no second `review-cycles` sample appears. `task-6.md`'s exact commands are rerun with the three DEC-227 package times. | Rerun the named journey with `-v` and the full tagged gate. Read the new assertions and confirm they compare count and consumed across recovery and rebuild. Confirm the report cites them. |
| MF-5 | task-10 | VAL-009 | `cmd/fid/integration/p23_human_outcome_test.go` ordering and assertions only (the existing `CoexistingReviewObligation` fixture already constructs the state), plus `task-10-done.md`'s VAL-009 entry and the I-6 citations. No product or fixture-builder change. | A stage on the three-type Envelope where `OBL-p22-corpus` is satisfied and `OBL-p23-review` is unsatisfied, preferably with both Human Obligations satisfied: `envelope.complete` is denied naming `OBL-p23-review`, with journal and projection unchanged. `issue.transition-terminal-success` at that stage asserts the outcome BC-1's disposition fixes; under current product behaviour that is denied with exactly `["obligation-evidence-absent"]`, as mutation M2b observed. The report's "Task 2 probes" and "monotone/not reachable" claims are removed and the new stage cited. The five `task-10.md` commands are rerun. | Rerun `TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess -v` and the full tagged gate. Read the new stage. Confirm every VAL-009 citation resolves to existing code. |

## Blocker Candidates

QA emits candidates only. No `ESC` identity was allocated, and no escalation
register was edited.

- **BC-1: "applicable Obligation" and `guarded_actions` semantics for guarded
  actions.**
  - *Protected category:* specification and product semantics, plus
    ambiguous acceptance criteria. VAL-009's "every applicable Obligation" is
    undefined in `GLOSSARY.md`.
  - *Evidence, product behaviour:*
    - The terminal-success transition evaluates every Envelope Obligation
      (`internal/workflow/evaluation.go:436`, `:496-512`), as do completion
      (`envelope.go:521`, `:546-556`) and `run_close.go:81`.
    - `guarded_actions` is consulted only in authoring (`envelope.go:630`,
      `:664`) and for the readiness Stop's blocked actions
      (`human_outcome_readiness.go:68`).
  - *Evidence, observed:*
    - Journey stage 3 (`p23_human_outcome_test.go:146`) asserts that
      `OBL-human-security` blocks `qa_pass`, although its `guarded_actions` is
      `["envelope.complete"]` (`p22_fake_test.go:386`). The comment at `:29-32`
      states that the terminal transition "evaluates every required Obligation
      regardless of its own guarded_actions".
    - Mutation M2b showed `OBL-p23-review` (`guarded_actions`
      `["envelope.complete"]`) alone denying `qa_pass`.
  - *Evidence, contract:* `contracts.md:60-63` (per-type limits on who may
    guard terminal success) and `:409` (the readiness Stop blocks "exactly that
    Obligation's guarded actions") imply per-Obligation applicability.
    `:91-92` ("every applicable Obligation") and PRD-ENV-010 are silent on the
    point. The behaviour dates from `710db23` (`Refs: phase-23/task-2`).
  - *Authority boundary:* QA and the Developer cannot choose between the two
    readings.
  - *Requested action:* the Tech Lead (Human if a product rule) decides whether
    an Obligation blocks only its declared `guarded_actions`.
    - If yes: a decision-free task-2 evaluator conformance repair, with VAL-009
      (both halves), the journey's stage 1-3 `qa_pass` assertions, and task 13's
      VAL-056 precedence probes re-derived.
    - If no: record the reading in a DEC and a `GLOSSARY.md` entry, including
      what `guarded_actions` governs.
  - *Forward impact:* task 13 VAL-056 and task 14 VAL-059 and VAL-071 fixtures.

- **BC-2: `human-outcome.record` Envelope-grant requirement contradiction.**
  - *Protected category:* public interface and authority/security rule.
  - *Evidence:*
    - `public-interface.md:243` (section 4.1) states "Profile with
      `actor_type:"human"` plus Envelope grant and active finalized
      Activation".
    - `public-interface.md:222-224` (section 4 body) and `contracts.md:391-392`
      (section 10) omit any Envelope grant.
    - The product checks only Actor type and capability
      (`internal/review/human_outcome.go:79-89`).
    - Built binary: mutation M1 recorded an approved Outcome
      (`EVT ... "ok":true`, `12 -> 13`) in a workspace containing no Envelope.
    - The committed task-9 journey has `actor-human-outcome` record Outcomes
      and atomically resolve the security protected-Decision Stop
      (`p23_decision_reference_test.go:55-56`, `:66-70`). That Actor has a
      profile (`p22_fake_test.go:343-345`) but no Envelope grant: the grants at
      `:364-375` name only `actor-recovery`, `actor-denied`, and
      `actor-reviewer-two`.
    - The section 4.1 row dates from the planning commit
      `f94e3541fed134852c3860e1a6c684991d86f696`, and no DEC addresses it.
  - *Authority boundary:* Human-accepted interface (DEC-222). Neither reading
    can be chosen by QA or the Developer.
  - *Requested action:* the Tech Lead (and Human if the interface changes)
    rules which text governs.
    - If the grant is required: a task-10 authorization repair and a
      validation-contract update adding the negative to VAL-043.
    - Otherwise: correct the section 4.1 cell under a DEC.
  - *Forward impact:* task 14 VAL-071 fixture and task 18 VAL-067/068 authority
    cases.

- **BC-3: No `stop resolve` qualifying rule for `human-outcome-recorded` and
  `decision-satisfied`, and unresolvable matching Stops.**
  - *Protected category:* specification (resolution semantics), ambiguous
    acceptance criteria (VAL-048 "reject", the VAL-078 per-type route), and a
    closed registry with no owning rule. This is the phase brief's recurring
    registry family; a registry rejection is to be escalated.
  - *Evidence, code:*
    - `internal/workflow/stop.go:638` (`eventMatchesCondition`) has no entry
      for either type, and `:618-634` (`qualifyingFact`) has no case.
    - Atomic resolution scans only the producer codes
      (`internal/review/human_outcome.go:239` for `human-outcome-required`,
      `decision.go:243` for `protected-decision-required`).
  - *Evidence, built-binary probe (scratch worktree at HEAD):*
    - `stop resolve` on the readiness Stop with a wrong-class Outcome Event
      and with a foreign Event each returned exit 4 `resolution-invalid`, with
      Events `24->24`. The rejection exists, but no committed test observes it.
    - A generic `stop raise` of an `operator-stop` carrying
      `human-outcome-recorded` (class `security`) was accepted with exit 0.
      This is the route `task-11.md` sanctions for per-type probes.
    - The matching approved security Outcome resolved only the readiness Stop
      `STP-971064`, not the generic Stop `STP-7f5468`.
    - `stop resolve` of `STP-7f5468` with that Outcome's Event
      `EVT-20260925-000026` returned `resolution-invalid` ("does not qualify"),
      with Events `26->26`. The Stop can never resolve.
    - By code reading, `decision-satisfied` behaves the same.
  - *Evidence, contract:* this conflicts with `contracts.md:455` (it "resolves
    only with the Event ID of a Human Outcome whose class and result match"),
    `:462` ("Only Stops whose recorded condition matches are resolved"),
    `:375` and `:410` (resolver capability `stop.resolve`), and VAL-078's
    evidence method.
  - *Authority boundary:* the qualifying rule for `decision-satisfied` (which
    fact of the pair, and Envelope scope) and any subject scoping for
    `human-outcome-recorded` are not uniquely fixed. Task 11 owns "the
    resolution behavior of all four" (`execution-scope.md:173-177`).
  - *Requested action:* the Tech Lead fixes the `stop resolve` qualifying rules
    for both types, and whether recording must resolve every open Stop whose
    Condition matches regardless of code. Task 11 then implements them in its
    `stop.go` scope and evidences VAL-048's rejection clause and VAL-078 per
    type.
  - *Forward impact:* task 17's `human_outcome_resolution:true` flip requires
    the complete path.

These candidates do not touch task 12's own work. Waiving the section 8 rule
that the checkpoint must PASS before task 12 is a Tech Lead or Runner decision,
not a QA decision.

## Existing Escalation References

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-23/qa/checkpoint-after-task-11-2.md
Must-fix: 2  Improvement: 4  Blocker candidates: 3
VERDICT: BLOCK
```
