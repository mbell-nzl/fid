# Phase 23 Human Fixtures

Reviewed on 2026-10-03 by gpt-6-astra QA Lead. This guide replaces the earlier
fixture snapshots with one fresh manifest at the reviewed HEAD, including
DEC-267 and tasks 19–21 (DEC-269/DEC-270).

## Generate a fresh set

Working directory: `/home/mark/projects/fid`.

```bash
FID_P23_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-23-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration
```

Wait for the command to finish. Its output names a new manifest with current
paths and IDs. Seeded active Claims last one hour, so generate immediately before a later
Human session and use that manifest's commands. Exact expected text below is
bound to this recorded generation. The original [human-checklist.md](human-checklist.md)
is the immutable historical Human result. Record only the seven re-checks in
[human-checklist-2.md](human-checklist-2.md).

## Reviewed fixture identity

- Reviewed source HEAD: `32a55a94efb2a89dbf50b29accb434b9c753418c`.
- Generated at: `2026-10-02T17:17:10.215903366Z`.
- Checkpoint root: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561`.
- Manifest: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/manifest.json`.
- Binary: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid`.
- Binary SHA-256: `sha256:9a94cc6405b2c64fa5a6e2eef027b9ab39e81461b1f179b4910e53e0ffec29da`.
- Source dirty flag: `true`; only workflow/QA documents are dirty. Product source is the reviewed HEAD.
- Durable manifest: [human-fixture-manifest.json](final-evidence-1/human-fixture-manifest.json).

Generation receipt, including the complete terminal result and wall-time
observation: [human-fixture-generation.json](final-evidence-1/human-fixture-generation.json).
Every invocation below was executed by QA; receipt:
[human-guide-verification.json](final-evidence-1/human-guide-verification.json).
These receipts establish runnable mechanics and exact text, not Human judgment.

The protected-Decision and rework steps change their fixtures. Run inspections
before the steps and use a fresh generation when repeating a complete sequence.
The leaf mutations repeat existing requests idempotently.

## VAL-070 mixed Findings

Fixture IDs:

```json
{
  "first_review_id": "REV-52abd2",
  "issue_id": "P22-d7b499",
  "second_review_id": "REV-ca6a86"
}
```

### Inspection 1

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings review get REV-52abd2
```

Expected Human output:

```text
Purpose: Inspect independent Reviews and their Findings.

Reviews
  Review REV-52abd2
    Issue: P22-d7b499; Envelope: ENV-861002 revision 1
    Run: RUN-ec8e47; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-5475e8060295e582dbf4aa4f3dd79ee4; worktree sha256:fffa66b52c81449881d8b0db618e00d3415798368634615951a7b709a6fa1131; HEAD f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 3; minor 0; nit 0
  Finding FND-fixture-addressed
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
  Finding FND-fixture-blocking
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-resolved
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: resolved
```

### Inspection 2

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings --json review get REV-52abd2
```

### Inspection 3

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings review list --issue P22-d7b499
```

Expected Human output:

```text
Purpose: Inspect independent Reviews and their Findings.

Reviews
  Review REV-52abd2
    Issue: P22-d7b499; Envelope: ENV-861002 revision 1
    Run: RUN-ec8e47; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-5475e8060295e582dbf4aa4f3dd79ee4; worktree sha256:fffa66b52c81449881d8b0db618e00d3415798368634615951a7b709a6fa1131; HEAD f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 3; minor 0; nit 0
  Finding FND-fixture-addressed
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
  Finding FND-fixture-blocking
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-resolved
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: resolved
  Review REV-ca6a86
    Issue: P22-d7b499; Envelope: ENV-861002 revision 1
    Run: RUN-29022d; Reviewer: actor-reviewer-two; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-5475e8060295e582dbf4aa4f3dd79ee4; worktree sha256:fffa66b52c81449881d8b0db618e00d3415798368634615951a7b709a6fa1131; HEAD f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-fixture-new
    Review: REV-ca6a86; Severity: major; Blocking: true; Disposition: open
```

### Inspection 4

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings --json review list --issue P22-d7b499
```

### Inspection 5

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings finding list --issue P22-d7b499
```

Expected Human output:

```text
Purpose: Inspect Finding dispositions and current blocking state.

Findings
  Finding FND-fixture-addressed
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
  Finding FND-fixture-blocking
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-new
    Review: REV-ca6a86; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-resolved
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: resolved
```

### Inspection 6

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings --json finding list --issue P22-d7b499
```

### Inspection 7

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings resume P22-d7b499 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings
```

Expected Human output:

```text
Issue P22-d7b499
  Status: todo
  Completion: incomplete
  Envelope: ENV-861002 (revision 1, active)

Current Run: RUN-29022d
  Status: active
  Actor ID: actor-reviewer-two
  Executor ID: executor-fixture-second-review

Claim CLM-0d6730
  Lease status now: active
  Recorded status: active
  Expires: 2026-10-02T18:17:15.761350562Z
  An active lease alone does not authorize an operation; its authority and repository checks still apply.

Repository
  Project: project-5475e8060295e582dbf4aa4f3dd79ee4
  Current worktree: sha256:fffa66b52c81449881d8b0db618e00d3415798368634615951a7b709a6fa1131
  HEAD: f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442
  Uncommitted changes: false
  Run's recorded worktree: sha256:fffa66b52c81449881d8b0db618e00d3415798368634615951a7b709a6fa1131

Next steps
  Candidate operations below are not an execution sequence.
  Ready means the recorded prerequisites permit a request; supplied inputs and authority are checked when it runs.
  Action run.inspect run:RUN-29022d: ready
    Purpose: Compare this repository with the Run's recorded start or checkpoint; inspect changed paths without recording a checkpoint.
  Action run.observe run:RUN-29022d: ready
    Purpose: Have the activated harness record a fresh observation of whether this Run's executor is live, terminated or unknown.
    Requires: run.observe
  Action gate.run run:RUN-29022d: blocked
    Purpose: Execute the configured check to produce Evidence for the current work.
    Requires: gate.run
  authority-unavailable for actor:actor-reviewer-two: gate.run
  budget-data-unavailable for envelope:ENV-861002: tokens
  Action claim.renew claim:CLM-0d6730: blocked
    Purpose: Extend an eligible work lease after checking authority.
    Requires: claim.renew
  authority-unavailable for actor:actor-reviewer-two: claim.renew
  Action claim.release claim:CLM-0d6730: ready
    Purpose: Give up an eligible work lease after checking authority; this does not prove the work complete.
    Requires: claim.release
  Action concern.record issue:P22-d7b499: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: concern.record
  authority-unavailable for issue:P22-d7b499: concern.record
  Action deviation.record issue:P22-d7b499: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: deviation.record
  authority-unavailable for issue:P22-d7b499: deviation.record
  Action decision.record-reference issue:P22-d7b499: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: decision.record-reference
  authority-unavailable for issue:P22-d7b499: decision.record-reference
  Action human-outcome.record issue:P22-d7b499: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: human-outcome.record
  authority-unavailable for issue:P22-d7b499: human-outcome.record
  Action review.record run:RUN-29022d: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: review.record
  review-not-independent for run:RUN-29022d: Run already reviewed

Completion requirements
  These describe what is needed to finish the Issue, separately from recovery checks above.
  Obligation OBL-p22-corpus: unsatisfied
    Evidence validity: invalid; freshness: stale; Artifacts: missing
  obligation-evidence-absent for obligation:OBL-p22-corpus: qualifying Evidence is absent
  Overall completion: incomplete
  obligation-evidence-absent for obligation:OBL-p22-corpus: qualifying Evidence is absent
  issue-not-terminal-success for issue:P22-d7b499: completed workflow category
  Open Findings: 3

Review and judgment
Reviews
  Review REV-52abd2
    Issue: P22-d7b499; Envelope: ENV-861002 revision 1
    Run: RUN-ec8e47; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-5475e8060295e582dbf4aa4f3dd79ee4; worktree sha256:fffa66b52c81449881d8b0db618e00d3415798368634615951a7b709a6fa1131; HEAD f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 3; minor 0; nit 0
  Finding FND-fixture-addressed
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
  Finding FND-fixture-blocking
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-resolved
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: resolved
  Review REV-ca6a86
    Issue: P22-d7b499; Envelope: ENV-861002 revision 1
    Run: RUN-29022d; Reviewer: actor-reviewer-two; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-5475e8060295e582dbf4aa4f3dd79ee4; worktree sha256:fffa66b52c81449881d8b0db618e00d3415798368634615951a7b709a6fa1131; HEAD f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base f3a1a2e2f19fb6e64d033d31c2c53dcba65e1442; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-fixture-new
    Review: REV-ca6a86; Severity: major; Blocking: true; Disposition: open
Open Findings
  Finding FND-fixture-addressed
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
  Finding FND-fixture-blocking
    Review: REV-52abd2; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-new
    Review: REV-ca6a86; Severity: major; Blocking: true; Disposition: open
Concerns
  None.
Deviations
  None.
Decision References
  None.
Human Outcomes
  None.
Protected Decisions
  None.

Budgets
Budget attempts: remaining 3 attempts; consumed 0; maximum 3; state available
Budget tokens: remaining unknown input-tokens; consumed unknown; maximum 10; state unknown
```

### Inspection 8

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings --json resume P22-d7b499 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/mixed-findings
```

## VAL-071 protected Decision

Fixture IDs:

```json
{
  "envelope_id": "ENV-3f958e",
  "human_actor_id": "actor-human-outcome",
  "issue_id": "P22-504696"
}
```

### Inspection 1

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision resume P22-504696 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision
```

Expected Human output:

```text
Issue P22-504696
  Status: qa_review
  Completion: incomplete
  Envelope: ENV-3f958e (revision 1, active)

Current Run: RUN-134ab7
  Status: active
  Actor ID: actor-recovery
  Executor ID: executor-corpus-prior

Claim CLM-00f41a
  Lease status now: active
  Recorded status: active
  Expires: 2026-10-02T18:17:20.479180459Z
  An active lease alone does not authorize an operation; its authority and repository checks still apply.

Repository
  Project: project-8a8cfda890a2cec7623a946b7d5f5b5f
  Current worktree: sha256:fd3a069e2077680bef8375286d77679b1e78936b06a19afdd73749cf1c99de12
  HEAD: 3990d0ea31fd2a668bdc005630419d217271c022
  Uncommitted changes: false
  Run's recorded worktree: sha256:fd3a069e2077680bef8375286d77679b1e78936b06a19afdd73749cf1c99de12

Next steps
  To collect proof for unfinished work, run the configured Gate (gate.run).
  A successful execution can supply Evidence; it does not itself complete the Issue.
  Candidate operations below are not an execution sequence.
  Ready means the recorded prerequisites permit a request; supplied inputs and authority are checked when it runs.
  Action run.inspect run:RUN-134ab7: ready
    Purpose: Compare this repository with the Run's recorded start or checkpoint; inspect changed paths without recording a checkpoint.
  Action run.observe run:RUN-134ab7: ready
    Purpose: Have the activated harness record a fresh observation of whether this Run's executor is live, terminated or unknown.
    Requires: run.observe
  Action gate.run run:RUN-134ab7: ready
    Purpose: Execute the configured check to produce Evidence for the current work.
    Requires: gate.run
  Action claim.renew claim:CLM-00f41a: blocked
    Purpose: Extend an eligible work lease after checking authority.
    Requires: claim.renew
  authority-unavailable for actor:actor-recovery: claim.renew
  Action claim.release claim:CLM-00f41a: ready
    Purpose: Give up an eligible work lease after checking authority; this does not prove the work complete.
    Requires: claim.release
  Action concern.record issue:P22-504696: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: concern.record
  authority-unavailable for issue:P22-504696: concern.record
  Action deviation.record issue:P22-504696: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: deviation.record
  authority-unavailable for issue:P22-504696: deviation.record
  Action decision.record-reference issue:P22-504696: ready
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: decision.record-reference
  Action human-outcome.record issue:P22-504696: ready
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: human-outcome.record

Blockers and attention
Stop STP-8fa9af blocks issue.transition-terminal-success because protected-decision-required: {"affected_actions":["issue.transition-terminal-success"],"decision_class":"security","envelope_id":"ENV-3f958e","required_capabilities":["decision.record-reference"]}; qualifying resolution: decision-satisfied for envelope:ENV-3f958e decision class security (resolver capability stop.resolve).
Attention ATT-1619e3: open; protected-decision-required; linked Stops: STP-8fa9af

Completion requirements
  These describe what is needed to finish the Issue, separately from recovery checks above.
  Obligation OBL-p22-corpus: satisfied
    Evidence validity: valid; freshness: fresh; Artifacts: available
  Overall completion: incomplete
  issue-not-terminal-success for issue:P22-504696: completed workflow category
  Open Findings: 0

Review and judgment
Reviews
  None.
Open Findings
  None.
Concerns
  None.
Deviations
  None.
Decision References
  Decision Reference DREF-fixture-security
    Class: security; Outcome: approved; Authority Actor: actor-recovery
    Required capabilities: decision.record-reference
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
Human Outcomes
  None.
Protected Decisions
  Class security: unsatisfied; affected actions: issue.transition-terminal-success
    Resolution requires an approved Decision Reference with the required capabilities and an approved Human Outcome of the same class.

Budgets
Budget attempts: remaining 2 attempts; consumed 1; maximum 3; state available
Budget tokens: remaining 10 input-tokens; consumed 0; maximum 10; state available
```

### Inspection 2

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision --json resume P22-504696 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision
```

### Inspection 3

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision decision list --issue P22-504696
```

Expected Human output:

```text
Purpose: Inspect Decision References and their authority.

Decision References
  Decision Reference DREF-fixture-security
    Class: security; Outcome: approved; Authority Actor: actor-recovery
    Required capabilities: decision.record-reference
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

### Inspection 4

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision --json decision list --issue P22-504696
```

### Inspection 5

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision stop list --subject envelope:ENV-3f958e
```

Expected Human output:

```text
stop.list
Stop STP-8fa9af blocks issue.transition-terminal-success because protected-decision-required: {"affected_actions":["issue.transition-terminal-success"],"decision_class":"security","envelope_id":"ENV-3f958e","required_capabilities":["decision.record-reference"]}; qualifying resolution: decision-satisfied for envelope:ENV-3f958e decision class security (resolver capability stop.resolve).
```

### Inspection 6

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision --json stop list --subject envelope:ENV-3f958e
```

### Step 1: blocked action before Outcome

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision`. Expected exit: `5`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision --json issue transition P22-504696 --to qa_pass --actor agent:actor-recovery --expected-revision 4 --idempotency-key fixture-human-transition
```

### Step 2: record authorized Human Outcome

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision --json human-outcome record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision/fixture-human-outcome.json --actor actor-human-outcome --idempotency-key fixture-human-outcome
```

### Step 3: retry guarded action

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/protected-decision --json issue transition P22-504696 --to qa_pass --actor agent:actor-recovery --expected-revision 4 --idempotency-key fixture-human-transition
```

## VAL-072 review rework exhaustion

Fixture IDs:

```json
{
  "envelope_id": "ENV-f39156",
  "human_actor_id": "actor-human-outcome",
  "issue_id": "P22-add6ba",
  "review_id": "REV-8abebf"
}
```

Read the denial and Stop/Attention before consulting the resolution step. Record the three rework judgments from what fid displays.


### Human rework denial

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `4`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion claim acquire --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion/fixture-rework-claim.json --actor actor-recovery --idempotency-key fixture-exhausted-claim
```

### Inspection 1

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion resume P22-add6ba --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion
```

Expected Human output:

```text
Issue P22-add6ba
  Status: todo
  Completion: incomplete
  Envelope: ENV-f39156 (revision 1, active)

Current Run: RUN-c3557d
  Status: failed
  Actor ID: actor-denied
  Executor ID: executor-fixture-cycle

Current Claim: none.

Repository
  Project: project-4c535c06f3bce304b4a98c743e1e43d8
  Current worktree: sha256:148cb9983ef5818e12aecc536655bdf07bf3988003afd4cfd205591b97808d89
  HEAD: 288a00221fec8b08567ca0227b29e4c06be5d9ae
  Uncommitted changes: false
  Run's recorded worktree: sha256:148cb9983ef5818e12aecc536655bdf07bf3988003afd4cfd205591b97808d89

Next steps
  Candidate operations below are not an execution sequence.
  Ready means the recorded prerequisites permit a request; supplied inputs and authority are checked when it runs.
  Action run.inspect run:RUN-c3557d: ready
    Purpose: Compare this repository with the Run's recorded start or checkpoint; inspect changed paths without recording a checkpoint.
  Action run.observe run:RUN-c3557d: blocked
    Purpose: Have the activated harness record a fresh observation of whether this Run's executor is live, terminated or unknown.
    Requires: run.observe
  run-inactive for run:RUN-c3557d: failed
  Action gate.run run:RUN-c3557d: blocked
    Purpose: Execute the configured check to produce Evidence for the current work.
    Requires: gate.run
  authority-unavailable for actor:actor-denied: gate.run
  run-inactive for run:RUN-c3557d: failed
  budget-data-unavailable for envelope:ENV-f39156: tokens
  Action concern.record issue:P22-add6ba: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: concern.record
  authority-unavailable for issue:P22-add6ba: concern.record
  Action deviation.record issue:P22-add6ba: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: deviation.record
  authority-unavailable for issue:P22-add6ba: deviation.record
  Action decision.record-reference issue:P22-add6ba: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: decision.record-reference
  authority-unavailable for issue:P22-add6ba: decision.record-reference
  Action human-outcome.record issue:P22-add6ba: ready
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: human-outcome.record

Blockers and attention
Stop STP-c3de71 blocks claim.acquire, run.start because review-rework-exhausted: {"budget_name":"review-cycles","consumed":1,"envelope_id":"ENV-f39156","last_review_id":"REV-8abebf","maximum":1,"remaining":0,"unit":"review-cycles"}; qualifying resolution: human-outcome-recorded for envelope:ENV-f39156 outcome class review-rework-extension (resolver capability stop.resolve).
Attention ATT-9a4ab8: open; review-rework-exhausted; linked Stops: STP-c3de71

Completion requirements
  These describe what is needed to finish the Issue, separately from recovery checks above.
  Obligation OBL-p22-corpus: unsatisfied
    Evidence validity: invalid; freshness: stale; Artifacts: not-applicable
  obligation-evidence-absent for obligation:OBL-p22-corpus: qualifying Evidence is absent
  Overall completion: incomplete
  obligation-evidence-absent for obligation:OBL-p22-corpus: qualifying Evidence is absent
  issue-not-terminal-success for issue:P22-add6ba: completed workflow category
  Open Findings: 1

Review and judgment
Reviews
  Review REV-8abebf
    Issue: P22-add6ba; Envelope: ENV-f39156 revision 1
    Run: RUN-c3557d; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-4c535c06f3bce304b4a98c743e1e43d8; worktree sha256:148cb9983ef5818e12aecc536655bdf07bf3988003afd4cfd205591b97808d89; HEAD 288a00221fec8b08567ca0227b29e4c06be5d9ae; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base 288a00221fec8b08567ca0227b29e4c06be5d9ae; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-fixture-cycle
    Review: REV-8abebf; Severity: major; Blocking: true; Disposition: open
Open Findings
  Finding FND-fixture-cycle
    Review: REV-8abebf; Severity: major; Blocking: true; Disposition: open
Concerns
  None.
Deviations
  None.
Decision References
  None.
Human Outcomes
  None.
Protected Decisions
  None.

Budgets
Budget attempts: remaining 3 attempts; consumed 0; maximum 3; state available
Budget review-cycles: remaining 0 review-cycles; consumed 1; maximum 1; state exhausted
Budget tokens: remaining unknown input-tokens; consumed unknown; maximum 10; state unknown
```

### Inspection 2

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion --json resume P22-add6ba --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion
```

### Inspection 3

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion stop list --subject envelope:ENV-f39156
```

Expected Human output:

```text
stop.list
Stop STP-c3de71 blocks claim.acquire, run.start because review-rework-exhausted: {"budget_name":"review-cycles","consumed":1,"envelope_id":"ENV-f39156","last_review_id":"REV-8abebf","maximum":1,"remaining":0,"unit":"review-cycles"}; qualifying resolution: human-outcome-recorded for envelope:ENV-f39156 outcome class review-rework-extension (resolver capability stop.resolve).
```

### Inspection 4

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion --json stop list --subject envelope:ENV-f39156
```

### Inspection 5

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion attention list --subject envelope:ENV-f39156
```

Expected Human output:

```text
attention.list
Attention ATT-9a4ab8 open for review-rework-exhausted; qualifying resolution requires dispositions for every linked Stop:
  Stop STP-c3de71 blocks claim.acquire, run.start because review-rework-exhausted: {"budget_name":"review-cycles","consumed":1,"envelope_id":"ENV-f39156","last_review_id":"REV-8abebf","maximum":1,"remaining":0,"unit":"review-cycles"}; qualifying resolution: human-outcome-recorded for envelope:ENV-f39156 outcome class review-rework-extension (resolver capability stop.resolve).
```

### Inspection 6

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion --json attention list --subject envelope:ENV-f39156
```

### Step 1: exhausted rework denial

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `4`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion --json claim acquire --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion/fixture-rework-claim.json --actor actor-recovery --idempotency-key fixture-exhausted-claim
```

### Step 2: record approved review-rework-extension Outcome

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion --json human-outcome record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion/fixture-rework-outcome.json --actor actor-human-outcome --idempotency-key fixture-rework-outcome
```

Task 21 resolution: this activated Human records an approved `review-rework-extension` Outcome. Its Event resolves the matching rework Stop and Attention atomically. No internal file editing is involved.

### Step 3: retry rework Claim

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion --json claim acquire --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/review-exhaustion/fixture-rework-claim.json --actor actor-recovery --idempotency-key fixture-exhausted-claim
```

## Every new leaf — Human and JSON pairs

### 01. requirement.record-reference

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract requirement record-reference --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-requirement.json --actor actor-recovery --idempotency-key leaf-requirement
```

Expected Human output:

```text
Purpose: Inspect Requirement References and Issue links.

Requirement References
  Requirement REQ-leaf
    Reference ID: RREF-150ca7; Status: active
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json requirement record-reference --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-requirement.json --actor actor-recovery --idempotency-key leaf-requirement
```

Expected JSON result facts:

```json
{
  "idempotency_key": "leaf-requirement",
  "replayed": true,
  "requirement": {
    "created_at": "2026-10-02T17:17:31.152056319Z",
    "event_id": "EVT-20261002-000018",
    "id": "RREF-150ca7",
    "requirement_id": "REQ-leaf",
    "source": {
      "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
      "display_label": "Task 9 decision",
      "external_id": null,
      "path": "base.txt",
      "repository": {
        "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
        "capture_source": "git",
        "dirty": false,
        "dirty_digest": null,
        "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
        "project_id": "project-302aace960d053da046da917ed52841b",
        "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
        "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
      },
      "section": null,
      "source_type": "project-file"
    },
    "status": "active",
    "superseded_by_id": null
  }
}
```

### 02. requirement.link

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract requirement link --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-link.json --actor actor-recovery --idempotency-key leaf-link
```

Expected Human output:

```text
Purpose: Inspect Requirement References and Issue links.

Issue links
  REQ-leaf implements P22-4338c8
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json requirement link --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-link.json --actor actor-recovery --idempotency-key leaf-link
```

Expected JSON result facts:

```json
{
  "idempotency_key": "leaf-link",
  "link": {
    "created_at": "2026-10-02T17:17:31.304799438Z",
    "event_id": "EVT-20261002-000019",
    "issue_id": "P22-4338c8",
    "relation": "implements",
    "requirement_id": "REQ-leaf"
  },
  "replayed": true
}
```

### 03. concern.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract concern record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-concern.json --actor actor-recovery --idempotency-key leaf-concern
```

Expected Human output:

```text
Purpose: Inspect a raised concern and its disposition as judgment.

Concerns
  Concern CON-leaf
    Subject: issue:P22-4338c8
    Severity: major; Status: open; Raised by: actor-recovery
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
    This Concern is recorded judgment, not deterministic proof.
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json concern record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-concern.json --actor actor-recovery --idempotency-key leaf-concern
```

Expected JSON result facts:

```json
{
  "concern": {
    "actor_id": "actor-recovery",
    "created_at": "2026-10-02T17:17:31.480113205Z",
    "disposition": null,
    "event_id": "EVT-20261002-000020",
    "id": "CON-leaf",
    "reasoning": "PRIVATE_CONCERN_REASONING",
    "references": [
      {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "project_id": "project-302aace960d053da046da917ed52841b",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
        },
        "section": null,
        "source_type": "project-file"
      }
    ],
    "scope": {
      "paths": []
    },
    "severity": "major",
    "status": "open",
    "subject": {
      "id": "P22-4338c8",
      "type": "issue"
    },
    "summary": "PRIVATE_CONCERN_SUMMARY"
  },
  "idempotency_key": "leaf-concern",
  "replayed": true
}
```

### 04. concern.dispose

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract concern dispose --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-dispose.json --actor actor-recovery --idempotency-key leaf-dispose
```

Expected Human output:

```text
Purpose: Inspect a raised concern and its disposition as judgment.

Concerns
  Concern CON-leaf
    Subject: issue:P22-4338c8
    Severity: major; Status: disposed; Raised by: actor-recovery
    Disposition: carried-forward; Disposed by: actor-recovery
    Work now owned by Issue: P22-4338c8
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
    This Concern is recorded judgment, not deterministic proof.
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json concern dispose --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-dispose.json --actor actor-recovery --idempotency-key leaf-dispose
```

Expected JSON result facts:

```json
{
  "concern": {
    "actor_id": "actor-recovery",
    "created_at": "2026-10-02T17:17:31.480113205Z",
    "disposition": {
      "action": "carried-forward",
      "actor_id": "actor-recovery",
      "disposed_at": "2026-10-02T17:17:31.704936997Z",
      "event_id": "EVT-20261002-000021",
      "issue_id": "P22-4338c8",
      "reason": "PRIVATE_DISPOSITION_REASON",
      "supersedes_concern_id": null,
      "supersedes_finding_id": null
    },
    "event_id": "EVT-20261002-000020",
    "id": "CON-leaf",
    "reasoning": "PRIVATE_CONCERN_REASONING",
    "references": [
      {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "project_id": "project-302aace960d053da046da917ed52841b",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
        },
        "section": null,
        "source_type": "project-file"
      }
    ],
    "scope": {
      "paths": []
    },
    "severity": "major",
    "status": "disposed",
    "subject": {
      "id": "P22-4338c8",
      "type": "issue"
    },
    "summary": "PRIVATE_CONCERN_SUMMARY"
  },
  "idempotency_key": "leaf-dispose",
  "replayed": true
}
```

### 05. deviation.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract deviation record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-deviation.json --actor actor-recovery --idempotency-key leaf-deviation
```

Expected Human output:

```text
Purpose: Inspect a recorded deviation as judgment.

Deviations
  Deviation DEV-leaf
    Subject: issue:P22-4338c8
    Disposition: pending-decision; Authority Actor: actor-recovery
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
    This Deviation is recorded judgment, not deterministic proof.
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json deviation record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-deviation.json --actor actor-recovery --idempotency-key leaf-deviation
```

Expected JSON result facts:

```json
{
  "deviation": {
    "actor_id": "actor-recovery",
    "authority": {
      "actor_id": "actor-recovery",
      "reference": null
    },
    "created_at": "2026-10-02T17:17:31.884176128Z",
    "disposition": "pending-decision",
    "event_id": "EVT-20261002-000022",
    "expected": "PRIVATE_EXPECTED",
    "id": "DEV-leaf",
    "observed": "PRIVATE_OBSERVED",
    "reasoning": "PRIVATE_DEVIATION_REASONING",
    "references": [
      {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "project_id": "project-302aace960d053da046da917ed52841b",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
        },
        "section": null,
        "source_type": "project-file"
      }
    ],
    "subject": {
      "id": "P22-4338c8",
      "type": "issue"
    },
    "supersedes_deviation_id": null
  },
  "idempotency_key": "leaf-deviation",
  "replayed": true
}
```

### 06. decision.record-reference

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract decision record-reference --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-decision.json --actor actor-recovery --idempotency-key leaf-decision
```

Expected Human output:

```text
Purpose: Inspect Decision References and their authority.

Decision References
  Decision Reference DREF-leaf
    Class: security; Outcome: approved; Authority Actor: actor-recovery
    Required capabilities: decision.record-reference
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json decision record-reference --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-decision.json --actor actor-recovery --idempotency-key leaf-decision
```

Expected JSON result facts:

```json
{
  "attention": [],
  "decision": {
    "authority": {
      "actor_id": "actor-recovery",
      "required_capabilities": [
        "decision.record-reference"
      ]
    },
    "class": "security",
    "created_at": "2026-10-02T17:17:32.114461714Z",
    "event_id": "EVT-20261002-000023",
    "id": "DREF-leaf",
    "outcome": "approved",
    "scope": {
      "issue_id": "P22-4338c8"
    },
    "source": {
      "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
      "display_label": "Task 9 decision",
      "external_id": null,
      "path": "base.txt",
      "repository": {
        "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
        "capture_source": "git",
        "dirty": false,
        "dirty_digest": null,
        "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
        "project_id": "project-302aace960d053da046da917ed52841b",
        "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
        "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
      },
      "section": null,
      "source_type": "project-file"
    }
  },
  "idempotency_key": "leaf-decision",
  "protected_decisions": [],
  "replayed": true,
  "stops": []
}
```

### 07. human-outcome.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract human-outcome record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-outcome.json --actor actor-recovery --idempotency-key leaf-outcome
```

Expected Human output:

```text
Purpose: Inspect authorized Human judgments and their provenance.

Human Outcomes
  Human Outcome OUT-leaf
    Actor: actor-recovery; Class: security; Result: approved
    Activation: ACT-f1d9ce revision 10; Profile: activation:p22-corpus
    This Outcome is authorized Human judgment, not deterministic proof.
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json human-outcome record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-outcome.json --actor actor-recovery --idempotency-key leaf-outcome
```

Expected JSON result facts:

```json
{
  "attention": [],
  "idempotency_key": "leaf-outcome",
  "outcome": {
    "actor_id": "actor-recovery",
    "class": "security",
    "created_at": "2026-10-02T17:17:32.32174547Z",
    "event_id": "EVT-20261002-000024",
    "id": "OUT-leaf",
    "provenance": {
      "activation_id": "ACT-f1d9ce",
      "activation_version": 10,
      "profile_ref": "activation:p22-corpus"
    },
    "reasoning": "PRIVATE_OUTCOME_REASONING",
    "references": [
      {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "project_id": "project-302aace960d053da046da917ed52841b",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
        },
        "section": null,
        "source_type": "project-file"
      }
    ],
    "result": "approved",
    "scope": {
      "issue_id": "P22-4338c8"
    },
    "subject": {
      "id": "P22-4338c8",
      "type": "issue"
    }
  },
  "protected_decisions": [],
  "replayed": true,
  "stops": []
}
```

### 08. review.record

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract review record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-review.json --actor actor-denied --idempotency-key leaf-review
```

Expected Human output:

```text
Purpose: Record an independent Review of a Run and its bound Issue.

Reviews
  Review REV-959494
    Issue: P22-4338c8; Envelope: ENV-192626 revision 1
    Run: RUN-1cf746; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-302aace960d053da046da917ed52841b; worktree sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150; HEAD d4d8902533118db91365b0535458d10fd36f98ce; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base d4d8902533118db91365b0535458d10fd36f98ce; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-leaf
    Review: REV-959494; Severity: major; Blocking: true; Disposition: open
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json review record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-review.json --actor actor-denied --idempotency-key leaf-review
```

Expected JSON result facts:

```json
{
  "attention": [],
  "budget": null,
  "idempotency_key": "leaf-review",
  "replayed": true,
  "review": {
    "actor_id": "actor-denied",
    "coverage": {
      "exclusions": [],
      "items": [
        {
          "assessed": true,
          "source": {
            "digest": "sha256:0000000000000000000000000000000000000000000000000000000000000000",
            "display_label": "Phase 22 corpus",
            "external_id": null,
            "path": "project-docs/p22-corpus.md",
            "repository": {
              "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
              "capture_source": "git",
              "dirty": false,
              "dirty_digest": null,
              "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
              "project_id": "project-302aace960d053da046da917ed52841b",
              "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
              "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
            },
            "section": null,
            "source_type": "project-file"
          }
        }
      ]
    },
    "envelope_id": "ENV-192626",
    "envelope_revision": 1,
    "event_id": "EVT-20261002-000030",
    "evidence_refs": [],
    "findings": [
      {
        "blocking": true,
        "created_at": "2026-10-02T17:17:33.176610983Z",
        "dispositions": [],
        "id": "FND-leaf",
        "review_id": "REV-959494",
        "scope": {
          "paths": []
        },
        "severity": "major",
        "state": "open",
        "subject": {
          "id": "P22-4338c8",
          "type": "issue"
        },
        "summary": "real-process Review finding"
      }
    ],
    "id": "REV-959494",
    "independence": "distinct-actor",
    "issue_id": "P22-4338c8",
    "profile": "review-default",
    "receipt_ids": [],
    "recorded_at": "2026-10-02T17:17:33.176610983Z",
    "repository": {
      "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
      "capture_source": "git",
      "dirty": false,
      "dirty_digest": null,
      "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
      "project_id": "project-302aace960d053da046da917ed52841b",
      "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
      "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
    },
    "resolved_finding_ids": [],
    "run_id": "RUN-1cf746",
    "verdict": "changes-required"
  },
  "stops": []
}
```

### 09. finding.address

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract finding address --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-address.json --actor actor-recovery --idempotency-key leaf-address
```

Expected Human output:

```text
Purpose: Record a Builder assertion that a Finding was addressed.

Findings
  Finding FND-leaf
    Review: REV-959494; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json finding address --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract/leaf-leaf-address.json --actor actor-recovery --idempotency-key leaf-address
```

Expected JSON result facts:

```json
{
  "finding": {
    "blocking": true,
    "created_at": "2026-10-02T17:17:33.176610983Z",
    "dispositions": [
      {
        "action": "addressed",
        "actor_id": "actor-recovery",
        "disposed_at": "2026-10-02T17:17:33.915543183Z",
        "event_id": "EVT-20261002-000035",
        "id": "FDP-3d9cc3",
        "qualifying_review_id": null,
        "reason": "PRIVATE_FINDING_REASON"
      }
    ],
    "id": "FND-leaf",
    "review_id": "REV-959494",
    "scope": {
      "paths": []
    },
    "severity": "major",
    "state": "addressed",
    "subject": {
      "id": "P22-4338c8",
      "type": "issue"
    },
    "summary": "real-process Review finding"
  },
  "idempotency_key": "leaf-address",
  "replayed": true
}
```

### 10. review.get

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract review get REV-959494
```

Expected Human output:

```text
Purpose: Inspect independent Reviews and their Findings.

Reviews
  Review REV-959494
    Issue: P22-4338c8; Envelope: ENV-192626 revision 1
    Run: RUN-1cf746; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-302aace960d053da046da917ed52841b; worktree sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150; HEAD d4d8902533118db91365b0535458d10fd36f98ce; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base d4d8902533118db91365b0535458d10fd36f98ce; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-leaf
    Review: REV-959494; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json review get REV-959494
```

Expected JSON result facts:

```json
{
  "review": {
    "actor_id": "actor-denied",
    "coverage": {
      "exclusions": [],
      "items": [
        {
          "assessed": true,
          "source": {
            "digest": "sha256:0000000000000000000000000000000000000000000000000000000000000000",
            "display_label": "Phase 22 corpus",
            "external_id": null,
            "path": "project-docs/p22-corpus.md",
            "repository": {
              "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
              "capture_source": "git",
              "dirty": false,
              "dirty_digest": null,
              "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
              "project_id": "project-302aace960d053da046da917ed52841b",
              "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
              "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
            },
            "section": null,
            "source_type": "project-file"
          }
        }
      ]
    },
    "envelope_id": "ENV-192626",
    "envelope_revision": 1,
    "event_id": "EVT-20261002-000030",
    "evidence_refs": [],
    "findings": [
      {
        "blocking": true,
        "created_at": "2026-10-02T17:17:33.176610983Z",
        "dispositions": [
          {
            "action": "addressed",
            "actor_id": "actor-recovery",
            "disposed_at": "2026-10-02T17:17:33.915543183Z",
            "event_id": "EVT-20261002-000035",
            "id": "FDP-3d9cc3",
            "qualifying_review_id": null,
            "reason": "PRIVATE_FINDING_REASON"
          }
        ],
        "id": "FND-leaf",
        "review_id": "REV-959494",
        "scope": {
          "paths": []
        },
        "severity": "major",
        "state": "addressed",
        "subject": {
          "id": "P22-4338c8",
          "type": "issue"
        },
        "summary": "real-process Review finding"
      }
    ],
    "id": "REV-959494",
    "independence": "distinct-actor",
    "issue_id": "P22-4338c8",
    "profile": "review-default",
    "receipt_ids": [],
    "recorded_at": "2026-10-02T17:17:33.176610983Z",
    "repository": {
      "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
      "capture_source": "git",
      "dirty": false,
      "dirty_digest": null,
      "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
      "project_id": "project-302aace960d053da046da917ed52841b",
      "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
      "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
    },
    "resolved_finding_ids": [],
    "run_id": "RUN-1cf746",
    "verdict": "changes-required"
  }
}
```

### 11. review.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract review list --issue P22-4338c8
```

Expected Human output:

```text
Purpose: Inspect independent Reviews and their Findings.

Reviews
  Review REV-959494
    Issue: P22-4338c8; Envelope: ENV-192626 revision 1
    Run: RUN-1cf746; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-302aace960d053da046da917ed52841b; worktree sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150; HEAD d4d8902533118db91365b0535458d10fd36f98ce; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base d4d8902533118db91365b0535458d10fd36f98ce; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-leaf
    Review: REV-959494; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json review list --issue P22-4338c8
```

Expected JSON result facts:

```json
{
  "items": [
    {
      "actor_id": "actor-denied",
      "coverage": {
        "exclusions": [],
        "items": [
          {
            "assessed": true,
            "source": {
              "digest": "sha256:0000000000000000000000000000000000000000000000000000000000000000",
              "display_label": "Phase 22 corpus",
              "external_id": null,
              "path": "project-docs/p22-corpus.md",
              "repository": {
                "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
                "capture_source": "git",
                "dirty": false,
                "dirty_digest": null,
                "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
                "project_id": "project-302aace960d053da046da917ed52841b",
                "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
                "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
              },
              "section": null,
              "source_type": "project-file"
            }
          }
        ]
      },
      "envelope_id": "ENV-192626",
      "envelope_revision": 1,
      "event_id": "EVT-20261002-000030",
      "evidence_refs": [],
      "findings": [
        {
          "blocking": true,
          "created_at": "2026-10-02T17:17:33.176610983Z",
          "dispositions": [
            {
              "action": "addressed",
              "actor_id": "actor-recovery",
              "disposed_at": "2026-10-02T17:17:33.915543183Z",
              "event_id": "EVT-20261002-000035",
              "id": "FDP-3d9cc3",
              "qualifying_review_id": null,
              "reason": "PRIVATE_FINDING_REASON"
            }
          ],
          "id": "FND-leaf",
          "review_id": "REV-959494",
          "scope": {
            "paths": []
          },
          "severity": "major",
          "state": "addressed",
          "subject": {
            "id": "P22-4338c8",
            "type": "issue"
          },
          "summary": "real-process Review finding"
        }
      ],
      "id": "REV-959494",
      "independence": "distinct-actor",
      "issue_id": "P22-4338c8",
      "profile": "review-default",
      "receipt_ids": [],
      "recorded_at": "2026-10-02T17:17:33.176610983Z",
      "repository": {
        "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
        "capture_source": "git",
        "dirty": false,
        "dirty_digest": null,
        "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
        "project_id": "project-302aace960d053da046da917ed52841b",
        "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
        "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
      },
      "resolved_finding_ids": [],
      "run_id": "RUN-1cf746",
      "verdict": "changes-required"
    }
  ]
}
```

### 12. finding.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract finding list --issue P22-4338c8
```

Expected Human output:

```text
Purpose: Inspect Finding dispositions and current blocking state.

Findings
  Finding FND-leaf
    Review: REV-959494; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json finding list --issue P22-4338c8
```

Expected JSON result facts:

```json
{
  "items": [
    {
      "blocking": true,
      "created_at": "2026-10-02T17:17:33.176610983Z",
      "dispositions": [
        {
          "action": "addressed",
          "actor_id": "actor-recovery",
          "disposed_at": "2026-10-02T17:17:33.915543183Z",
          "event_id": "EVT-20261002-000035",
          "id": "FDP-3d9cc3",
          "qualifying_review_id": null,
          "reason": "PRIVATE_FINDING_REASON"
        }
      ],
      "id": "FND-leaf",
      "review_id": "REV-959494",
      "scope": {
        "paths": []
      },
      "severity": "major",
      "state": "addressed",
      "subject": {
        "id": "P22-4338c8",
        "type": "issue"
      },
      "summary": "real-process Review finding"
    }
  ]
}
```

### 13. concern.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract concern list --subject issue:P22-4338c8
```

Expected Human output:

```text
Purpose: Inspect a raised concern and its disposition as judgment.

Concerns
  Concern CON-leaf
    Subject: issue:P22-4338c8
    Severity: major; Status: disposed; Raised by: actor-recovery
    Disposition: carried-forward; Disposed by: actor-recovery
    Work now owned by Issue: P22-4338c8
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
    This Concern is recorded judgment, not deterministic proof.
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json concern list --subject issue:P22-4338c8
```

Expected JSON result facts:

```json
{
  "items": [
    {
      "actor_id": "actor-recovery",
      "created_at": "2026-10-02T17:17:31.480113205Z",
      "disposition": {
        "action": "carried-forward",
        "actor_id": "actor-recovery",
        "disposed_at": "2026-10-02T17:17:31.704936997Z",
        "event_id": "EVT-20261002-000021",
        "issue_id": "P22-4338c8",
        "reason": "PRIVATE_DISPOSITION_REASON",
        "supersedes_concern_id": null,
        "supersedes_finding_id": null
      },
      "event_id": "EVT-20261002-000020",
      "id": "CON-leaf",
      "reasoning": "PRIVATE_CONCERN_REASONING",
      "references": [
        {
          "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
          "display_label": "Task 9 decision",
          "external_id": null,
          "path": "base.txt",
          "repository": {
            "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
            "capture_source": "git",
            "dirty": false,
            "dirty_digest": null,
            "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
            "project_id": "project-302aace960d053da046da917ed52841b",
            "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
            "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
          },
          "section": null,
          "source_type": "project-file"
        }
      ],
      "scope": {
        "paths": []
      },
      "severity": "major",
      "status": "disposed",
      "subject": {
        "id": "P22-4338c8",
        "type": "issue"
      },
      "summary": "PRIVATE_CONCERN_SUMMARY"
    }
  ]
}
```

### 14. deviation.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract deviation list --subject issue:P22-4338c8
```

Expected Human output:

```text
Purpose: Inspect a recorded deviation as judgment.

Deviations
  Deviation DEV-leaf
    Subject: issue:P22-4338c8
    Disposition: pending-decision; Authority Actor: actor-recovery
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
    This Deviation is recorded judgment, not deterministic proof.
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json deviation list --subject issue:P22-4338c8
```

Expected JSON result facts:

```json
{
  "items": [
    {
      "actor_id": "actor-recovery",
      "authority": {
        "actor_id": "actor-recovery",
        "reference": null
      },
      "created_at": "2026-10-02T17:17:31.884176128Z",
      "disposition": "pending-decision",
      "event_id": "EVT-20261002-000022",
      "expected": "PRIVATE_EXPECTED",
      "id": "DEV-leaf",
      "observed": "PRIVATE_OBSERVED",
      "reasoning": "PRIVATE_DEVIATION_REASONING",
      "references": [
        {
          "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
          "display_label": "Task 9 decision",
          "external_id": null,
          "path": "base.txt",
          "repository": {
            "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
            "capture_source": "git",
            "dirty": false,
            "dirty_digest": null,
            "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
            "project_id": "project-302aace960d053da046da917ed52841b",
            "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
            "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
          },
          "section": null,
          "source_type": "project-file"
        }
      ],
      "subject": {
        "id": "P22-4338c8",
        "type": "issue"
      },
      "supersedes_deviation_id": null
    }
  ]
}
```

### 15. decision.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract decision list --issue P22-4338c8
```

Expected Human output:

```text
Purpose: Inspect Decision References and their authority.

Decision References
  Decision Reference DREF-leaf
    Class: security; Outcome: approved; Authority Actor: actor-recovery
    Required capabilities: decision.record-reference
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json decision list --issue P22-4338c8
```

Expected JSON result facts:

```json
{
  "items": [
    {
      "authority": {
        "actor_id": "actor-recovery",
        "required_capabilities": [
          "decision.record-reference"
        ]
      },
      "class": "security",
      "created_at": "2026-10-02T17:17:32.114461714Z",
      "event_id": "EVT-20261002-000023",
      "id": "DREF-leaf",
      "outcome": "approved",
      "scope": {
        "issue_id": "P22-4338c8"
      },
      "source": {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "project_id": "project-302aace960d053da046da917ed52841b",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
        },
        "section": null,
        "source_type": "project-file"
      }
    }
  ]
}
```

### 16. human-outcome.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract human-outcome list --issue P22-4338c8
```

Expected Human output:

```text
Purpose: Inspect authorized Human judgments and their provenance.

Human Outcomes
  Human Outcome OUT-leaf
    Actor: actor-recovery; Class: security; Result: approved
    Activation: ACT-f1d9ce revision 10; Profile: activation:p22-corpus
    This Outcome is authorized Human judgment, not deterministic proof.
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json human-outcome list --issue P22-4338c8
```

Expected JSON result facts:

```json
{
  "items": [
    {
      "actor_id": "actor-recovery",
      "class": "security",
      "created_at": "2026-10-02T17:17:32.32174547Z",
      "event_id": "EVT-20261002-000024",
      "id": "OUT-leaf",
      "provenance": {
        "activation_id": "ACT-f1d9ce",
        "activation_version": 10,
        "profile_ref": "activation:p22-corpus"
      },
      "reasoning": "PRIVATE_OUTCOME_REASONING",
      "references": [
        {
          "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
          "display_label": "Task 9 decision",
          "external_id": null,
          "path": "base.txt",
          "repository": {
            "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
            "capture_source": "git",
            "dirty": false,
            "dirty_digest": null,
            "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
            "project_id": "project-302aace960d053da046da917ed52841b",
            "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
            "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
          },
          "section": null,
          "source_type": "project-file"
        }
      ],
      "result": "approved",
      "scope": {
        "issue_id": "P22-4338c8"
      },
      "subject": {
        "id": "P22-4338c8",
        "type": "issue"
      }
    }
  ]
}
```

### 17. requirement.get

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract requirement get REQ-leaf
```

Expected Human output:

```text
Purpose: Inspect Requirement References and Issue links.

Requirement References
  Requirement REQ-leaf
    Reference ID: RREF-150ca7; Status: active
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac

Issue links
  REQ-leaf implements P22-4338c8
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json requirement get REQ-leaf
```

Expected JSON result facts:

```json
{
  "links": [
    {
      "created_at": "2026-10-02T17:17:31.304799438Z",
      "event_id": "EVT-20261002-000019",
      "issue_id": "P22-4338c8",
      "relation": "implements",
      "requirement_id": "REQ-leaf"
    }
  ],
  "requirement": {
    "created_at": "2026-10-02T17:17:31.152056319Z",
    "event_id": "EVT-20261002-000018",
    "id": "RREF-150ca7",
    "requirement_id": "REQ-leaf",
    "source": {
      "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
      "display_label": "Task 9 decision",
      "external_id": null,
      "path": "base.txt",
      "repository": {
        "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
        "capture_source": "git",
        "dirty": false,
        "dirty_digest": null,
        "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
        "project_id": "project-302aace960d053da046da917ed52841b",
        "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
        "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
      },
      "section": null,
      "source_type": "project-file"
    },
    "status": "active",
    "superseded_by_id": null
  }
}
```

### 18. requirement.list

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract requirement list
```

Expected Human output:

```text
Purpose: Inspect Requirement References and Issue links.

Requirement References
  Requirement REQ-leaf
    Reference ID: RREF-150ca7; Status: active
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON counterpart:

Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract`. Expected exit: `0`.

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/leaf-contract --json requirement list
```

Expected JSON result facts:

```json
{
  "items": [
    {
      "created_at": "2026-10-02T17:17:31.152056319Z",
      "event_id": "EVT-20261002-000018",
      "id": "RREF-150ca7",
      "requirement_id": "REQ-leaf",
      "source": {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "d4d8902533118db91365b0535458d10fd36f98ce",
          "project_id": "project-302aace960d053da046da917ed52841b",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:0c8bfc0f39e219f6a50f7efa556147fbdf56c4a0637b44ae1d8fd32f29adc150"
        },
        "section": null,
        "source_type": "project-file"
      },
      "status": "active",
      "superseded_by_id": null
    }
  ]
}
```
