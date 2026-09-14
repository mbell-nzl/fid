# Phase 22 Task 19 Independent Review

Reviewer: independent review `/root/human_repair_review`  
Date: 2026-09-12 Pacific/Auckland  
Branch: `phase-22`  
Reviewed commit: `eff28fb1b7ea0a543017bef00ccb183cc79fc396`  
Baseline parent: `818685f5ac86e0cf5201ded6bc0205da772fb829`

## Verdict

**REPAIR REQUIRED before fresh Human retest.**

Must-fix findings: **1**. Improvement findings: **2**. Human outcomes remain
unresolved: this review does not mark VAL-076, VAL-077 or VAL-078 PASS and does
not replace the six Human judgments required by DEC-219.

Adoption, authority, JSON preservation and activation guidance are factually
consistent with task 19, DEC-219 and the Phase 22 contracts. Imported liveness
guidance has one remaining state-dependent contradiction described in MF-01.
The root agent's full task gates were still pending during this review and are
not claimed as this reviewer's PASS evidence.

## Scope And Evidence

The review inspected the task-19 spec, DEC-219, immutable Human checklist
attempt 1, VAL-053/076/077/078, contracts sections 3/8/10, public-interface
section 6, the complete product/test diff from the baseline parent, and the
directly involved workflow import-binding, observation, reconciliation, resume
action and activation types. No hidden scenario material was read. Product
source remained frozen; this file is the reviewer's only authored artifact.

Observed evidence, all from `/home/mark/projects/fid` unless a fixture path is
named:

| Evidence | Result |
|---|---|
| `git diff --check 818685f5ac86e0cf5201ded6bc0205da772fb829 eff28fb1b7ea0a543017bef00ccb183cc79fc396` | PASS — exit 0 (no output); terminated `exit_code: 0` |
| `go build -o /tmp/fid-task19-review ./cmd/fid` | PASS — exit 0 (no output); terminated `exit_code: 0` |
| `go test ./cmd/fid -run '^TestHuman(RecoveryGuidance\|GuidanceDoesNotRecommendBlockedGate\|ActivationTerminalStateDoesNotOfferRecovery)$' -count=1` | PASS — `ok github.com/mbell-nzl/fid/cmd/fid 0.141s`; terminated `exit_code: 0` |
| `go test ./internal/workflow -run '^TestRecoveryReadContextWithoutImportPreservesRecordedLiveness$' -count=1` | PASS — `ok github.com/mbell-nzl/fid/internal/workflow 0.434s`; terminated `exit_code: 0` |
| `go test -tags integration ./cmd/fid/integration -run '^TestP22HumanGuidance$' -count=1` | PASS — `ok github.com/mbell-nzl/fid/cmd/fid/integration 3.026s`; terminated `exit_code: 0` |

The review binary exercised all six retained public Human reads in
`.cache/phase-22-human/checkpoints-4132086272`: recovery resume/inspect,
imported resume/inspect, prepared activation status and committed activation
status. Each terminated `exit_code: 0`. The observed text exposed the adopted
`base.txt`, the source and current worktree identities, unknown destination
liveness, verified imported project/digest/destination, all five activation
resources, restore-prior before commit and complete-forward after commit.

For each of those same six fixtures and argv, the retained attempt-1 binary and
`/tmp/fid-task19-review` produced byte-identical JSON with both commands
terminating `exit_code: 0`:

| Read | Retained/current bytes | Byte equal |
|---|---:|---|
| Recovery resume | 9707 / 9707 | yes |
| Recovery run inspect | 2248 / 2248 | yes |
| Imported resume | 8215 / 8215 | yes |
| Imported run inspect | 1877 / 1877 | yes |
| Prepared activation status | 1430 / 1430 | yes |
| Committed activation status | 1677 / 1677 | yes |

This is direct evidence for unchanged JSON on the reviewed surfaces. It does
not substitute for the task's complete required command set.

## Contract Review

### Adoption And Completion Proof

The retained Recovery projection remains the source of the adopted paths. The
Human renderer now explains that adoption keeps existing uncommitted changes
for the successor Run while establishing no testing, review, approval or
correctness. It also avoids inventing a prior Gate interruption from missing
Evidence. Immediate actions appear before separately labelled completion
requirements, so the repeated obligation reason retains its two factual
contexts instead of appearing as the recommended recovery operation.

### Candidate Operations And Authority

The renderer derives ready versus blocked solely from the existing
`NextAction.reasons` contract: an empty set is ready and any reason is blocked.
It calls operations candidates rather than a sequence, explains their purpose,
and states that request inputs and authority are checked at execution. The
Claim section explicitly says an active lease alone does not authorize an
operation. A ready Gate is recommended only when its reason set is empty; the
focused blocked-Gate case confirms a budget-blocked Gate is not recommended.
No mutation predicate or authority projection changed.

### Imported And Destination State

Imported history is labelled only after `ReadImportBinding` verifies the local
marker, target root, Project identity and imported journal material. The view
identifies the preserved Project, original Event count/order, archive digest
and destination, while the Run's recorded worktree and current destination
worktree remain separately labelled. It states that import grants no mutation
authority and separates destination inspection/liveness from completion
Evidence.

Liveness provenance follows the observation Event identity: an Event listed in
the verified import binding is historical; a later local Event is a destination
observation. Absence remains unknown. The focused integration case covered all
three states, verified that the Human reads append no Event, and rejected an
invalid import marker rather than presenting it as verified history. Dated
observations are described as facts at their observation time, not promises of
current executor state. That classification is correct, but the surrounding
imperative is not conditional on the classified state or Run status; see MF-01.

### Activation Direction And Resources

Nonterminal direction follows the accepted point-of-no-return field:
restore-prior before commit and complete-forward after commit. Finalized and
aborted states take precedence and offer no recovery direction. The output
distinguishes the activation/install commit from Git, names the kit/backend and
all five resources, shows prior/target/observed identity and state, and explains
that an absent prior identity means returning to absence. The public prepared
and committed fixtures showed the expected direction and consequence without
claiming post-commit rollback.

## Findings

### Must-fix

**MF-01 — Import-specific intervention guidance is unconditional and becomes
false or contradictory after the destination state advances.**

`recoveryLocalContextLines` always emits “Before continuing work here ...
establish fresh destination liveness ... (run.observe)” whenever verified import
metadata exists. It does not inspect `context.Liveness`,
`context.LivenessHistorical`, or the Run status. This creates two incorrect
states:

- after a fresh destination observation, the view simultaneously directs the
  Human to establish fresh destination liveness and says that a destination
  observation is already recorded below;
- when no active Run remains, `resumeRun` still returns the latest terminal Run,
  so the import section directs `run.observe` while `resumeNextActions` correctly
  marks that operation blocked with `run-inactive`. Terminal imported work may
  have no remaining recovery mutation at all.

This is a Human-visible factual conflict in the exact imported-versus-fresh
destination-liveness boundary assigned to task 19. Condition the guidance on
both Run status and provenance: an active Run with no destination observation
should request one; an active Run with a destination observation should identify
it and tell the Human to check its time/applicability without implying absence;
a terminal Run should not direct the Human to call `run.observe`. Add assertions
that reject the stale imperative in fresh-observation and terminal states. The
existing focused integration case proves the labels but does not assert absence
of the contradictory instruction.

### Improvement

1. `resumeNextStepLines` says “Available operations below” even though the list
   intentionally includes blocked operations. The following ready/blocked
   labels make the state unambiguous, so this is not a factual blocker; “Candidate
   operations considered” would be slightly more precise.
2. Terminal activation output still ends with the generic explanation of what
   restoring an absent prior resource means after stating that no activation
   recovery remains. The terminal precedence is clear and no direction is
   offered, so this does not violate the contract; conditionally omitting that
   sentence would reduce irrelevant text.

## Validation Disposition

| Assertion | Review disposition | Basis |
|---|---|---|
| VAL-053 | REPAIR REQUIRED | JSON remained byte identical, but the Human import prose can contradict current destination/Run state under MF-01. |
| VAL-076 | READY FOR HUMAN RETEST | Adoption, next action, candidates, authority and completion proof are explained; no Human outcome inferred. |
| VAL-077 | READY FOR HUMAN RETEST | Both directions, point-of-no-return consequence, resource identities and terminal behavior are explained; no Human outcome inferred. |
| VAL-078 | REPAIR REQUIRED | Restored identity/history and provenance are explained, but destination intervention text must reflect fresh and terminal states under MF-01. |

Fresh fixture verification and the six Human responses remain outside this
review's evidence. Attempt 1 stays immutable.

Review written: `project-docs/docs/phases/phase-22/qa/task-19-review.md`  
Must-fix: 1  Improvement: 2  Blocker candidates: 0  
VERDICT: REPAIR REQUIRED BEFORE HUMAN RETEST

## Closure Re-review

Closure date: 2026-09-12 Pacific/Auckland  
MF-01 repair: `94afb3b27b2174078e19229d6771d89e40d740d4` — `fix(cli): condition import guidance on current run state`  
Final product source: `5232d9cdfa32e71b1a54fe7a73025eac86a97a7c` — `fix(cli): clarify forward-only activation recovery`

This bounded re-review preserves the initial REPAIR verdict and MF-01 above.
It inspected only the two repair commits and reran the directly affected cases.
The root agent's full gate rerun and fixture generation remain separate evidence.

### Closure Evidence

| Evidence | Result |
|---|---|
| `git show --format=fuller --stat 94afb3b27b2174078e19229d6771d89e40d740d4` | PASS — exact six-file Human renderer/test repair and `Refs: phase-22/task-19` observed; terminated `exit_code: 0` |
| `git show --format=fuller --stat 5232d9cdfa32e71b1a54fe7a73025eac86a97a7c` | PASS — exact one-file activation wording repair and `Refs: phase-22/task-19` observed; terminated `exit_code: 0` |
| `git diff --check eff28fb1b7ea0a543017bef00ccb183cc79fc396 94afb3b27b2174078e19229d6771d89e40d740d4` | PASS — exit 0 (no output); terminated `exit_code: 0` |
| `git diff --check 94afb3b27b2174078e19229d6771d89e40d740d4 5232d9cdfa32e71b1a54fe7a73025eac86a97a7c` | PASS — exit 0 (no output); terminated `exit_code: 0` |
| `go test ./cmd/fid -run '^TestHumanImported(FreshObservationDoesNotRequestItAgain\|InactiveRunDoesNotRequestObservation)$' -count=1` | PASS — `ok github.com/mbell-nzl/fid/cmd/fid 0.099s`; terminated `exit_code: 0` |
| `go test -tags integration ./cmd/fid/integration -run '^TestP22HumanGuidance$' -count=1` | PASS — `ok github.com/mbell-nzl/fid/cmd/fid/integration 3.336s`; terminated `exit_code: 0` |
| `go test ./cmd/fid -run '^Test(HumanActivationTerminalStateDoesNotOfferRecovery\|ActivationStatusHumanAndJSONDescribeRecoveryDirectionAtCommitBoundary)$' -count=1` | PASS — `ok github.com/mbell-nzl/fid/cmd/fid 0.283s`; terminated `exit_code: 0` |

### Finding Closure

**MF-01 — RESOLVED.** Import identity/history and completion separation remain
unconditional facts. Intervention guidance now depends on the current Run and
liveness provenance:

- no Run or a terminal Run receives no continue-work or `run.observe`
  directive;
- an active Run with absent, imported-historical or unknown liveness requests a
  fresh destination observation;
- an active Run with a local live/terminated observation acknowledges the
  recorded destination observation and asks the Human to check its time and
  applicability without claiming it is missing.

The focused integration test now exercises a real imported terminal Run through
both resume and run inspection and rejects the stale directive. It also rejects
that directive after a fresh destination observation. The focused unit test
covers nil and succeeded/failed/abandoned Run states. These are the negative
state assertions absent from the initial implementation.

**Improvement 1 — RESOLVED.** The action introduction now says “Candidate
operations” and no longer calls blocked entries available.

**Improvement 2 — RESOLVED.** The absent-prior restoration explanation is now
limited to precommit nonterminal activation. It is omitted for terminal and
complete-forward states. The final wording still explains that complete-forward
must finish installation/finalization and that returning to the prior state is
unavailable, without suggesting a post-commit rollback.

### Closure Validation Disposition

| Assertion | Closure disposition | Basis |
|---|---|---|
| VAL-053 | ON TRACK | MF-01's Human-state contradiction is resolved; the initial six-surface JSON byte-parity evidence remains unchanged by Human-only renderer edits. |
| VAL-076 | READY FOR HUMAN RETEST | Adoption, proof, candidate-action and authority guidance remains intact; no Human outcome inferred. |
| VAL-077 | READY FOR HUMAN RETEST | Direction/resource/terminal guidance remains intact and the retained commit-boundary case passes; no Human outcome inferred. |
| VAL-078 | READY FOR HUMAN RETEST | Imported intervention guidance now distinguishes absent/historical/unknown, fresh destination, terminal and absent-Run states; no Human outcome inferred. |

Closure must-fix: 0  Improvement: 0  Blocker candidates: 0  
**CLOSURE VERDICT: PASS FOR INDEPENDENT TECHNICAL REVIEW — HUMAN RETEST REQUIRED.**
