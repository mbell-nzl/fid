# Phase 23 Human fixture guide

Only the Human judges VAL-070, VAL-071, and VAL-072; VAL-073 remains a phase-end Human view review. The generator and recorded text make the checks runnable.

## Generate a fresh set

From `/home/mark/projects/fid`:

```bash
FID_P23_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-23-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration
```

Each run creates a separate retained `checkpoints-*` directory. Keep its binary and workspaces together. The manifest binds the binary digest, source commit, exact working directories, argument arrays, and expected Human stdout. Regenerate after the delivered binary changes.

## Recorded fixture set

- Manifest: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/manifest.json`
- Binary: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid`
- Binary SHA-256: `sha256:2b0eb89940e02d4ad84d971981dee85f4bdbcc9e190da4213d8ca71bc98e64ef`
- Source commit: `df24c0cf832e6c29a0f0a743150c91ff178d841a`; source dirty: `true`

## VAL-070 mixed Findings

- Project: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings`
- IDs: `{"first_review_id": "REV-367292", "issue_id": "P22-764499", "second_review_id": "REV-b90c1c"}`

### Read 1

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings`
- Arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings", "review", "get", "REV-367292"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings review get REV-367292
```

Verbatim Human stdout:

```text
Purpose: Inspect independent Reviews and their Findings.

Reviews
  Review REV-367292
    Issue: P22-764499; Envelope: ENV-e2b3c3 revision 1
    Run: RUN-0b264e; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-7d0d9392020521f7351a9197fad2ac4d; worktree sha256:a4936266d2d82b477ba646de792bf42fcbe02de77cca3c2481a83b1ab2758b10; HEAD 114d9451b3f51394221c0839145ca411370828b2; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base 114d9451b3f51394221c0839145ca411370828b2; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 3; minor 0; nit 0
  Finding FND-fixture-addressed
    Review: REV-367292; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
  Finding FND-fixture-blocking
    Review: REV-367292; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-resolved
    Review: REV-367292; Severity: major; Blocking: true; Disposition: resolved
```

- JSON counterpart for read 1: cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings", "--json", "review", "get", "REV-367292"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings --json review get REV-367292
```

### Read 3

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings`
- Arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings", "review", "list", "--issue", "P22-764499"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings review list --issue P22-764499
```

Verbatim Human stdout:

```text
Purpose: Inspect independent Reviews and their Findings.

Reviews
  Review REV-367292
    Issue: P22-764499; Envelope: ENV-e2b3c3 revision 1
    Run: RUN-0b264e; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-7d0d9392020521f7351a9197fad2ac4d; worktree sha256:a4936266d2d82b477ba646de792bf42fcbe02de77cca3c2481a83b1ab2758b10; HEAD 114d9451b3f51394221c0839145ca411370828b2; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base 114d9451b3f51394221c0839145ca411370828b2; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 3; minor 0; nit 0
  Finding FND-fixture-addressed
    Review: REV-367292; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
  Finding FND-fixture-blocking
    Review: REV-367292; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-resolved
    Review: REV-367292; Severity: major; Blocking: true; Disposition: resolved
  Review REV-b90c1c
    Issue: P22-764499; Envelope: ENV-e2b3c3 revision 1
    Run: RUN-7fb8a5; Reviewer: actor-reviewer-two; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-7d0d9392020521f7351a9197fad2ac4d; worktree sha256:a4936266d2d82b477ba646de792bf42fcbe02de77cca3c2481a83b1ab2758b10; HEAD 114d9451b3f51394221c0839145ca411370828b2; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base 114d9451b3f51394221c0839145ca411370828b2; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-fixture-new
    Review: REV-b90c1c; Severity: major; Blocking: true; Disposition: open
```

- JSON counterpart for read 3: cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings", "--json", "review", "list", "--issue", "P22-764499"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings --json review list --issue P22-764499
```

### Read 5

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings`
- Arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings", "finding", "list", "--issue", "P22-764499"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings finding list --issue P22-764499
```

Verbatim Human stdout:

```text
Purpose: Inspect Finding dispositions and current blocking state.

Findings
  Finding FND-fixture-addressed
    Review: REV-367292; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
  Finding FND-fixture-blocking
    Review: REV-367292; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-new
    Review: REV-b90c1c; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-resolved
    Review: REV-367292; Severity: major; Blocking: true; Disposition: resolved
```

- JSON counterpart for read 5: cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings", "--json", "finding", "list", "--issue", "P22-764499"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings --json finding list --issue P22-764499
```

### Read 7

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings`
- Arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings", "resume", "P22-764499", "--repository", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings resume P22-764499 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings
```

Verbatim Human stdout:

```text
Issue P22-764499
  Status: todo
  Completion: incomplete
  Envelope: ENV-e2b3c3 (revision 1, active)

Current Run: RUN-7fb8a5
  Status: active
  Actor ID: actor-reviewer-two
  Executor ID: executor-fixture-second-review

Claim CLM-85890b
  Lease status now: active
  Recorded status: active
  Expires: 2026-09-29T08:58:30.667478918Z
  An active lease alone does not authorize an operation; its authority and repository checks still apply.

Repository
  Project: project-7d0d9392020521f7351a9197fad2ac4d
  Current worktree: sha256:a4936266d2d82b477ba646de792bf42fcbe02de77cca3c2481a83b1ab2758b10
  HEAD: 114d9451b3f51394221c0839145ca411370828b2
  Uncommitted changes: false
  Run's recorded worktree: sha256:a4936266d2d82b477ba646de792bf42fcbe02de77cca3c2481a83b1ab2758b10

Next steps
  Candidate operations below are not an execution sequence.
  Ready means the recorded prerequisites permit a request; supplied inputs and authority are checked when it runs.
  Action run.inspect run:RUN-7fb8a5: ready
    Purpose: Compare this repository with the Run's recorded start or checkpoint; inspect changed paths without recording a checkpoint.
  Action run.observe run:RUN-7fb8a5: ready
    Purpose: Have the activated harness record a fresh observation of whether this Run's executor is live, terminated or unknown.
    Requires: run.observe
  Action gate.run run:RUN-7fb8a5: blocked
    Purpose: Execute the configured check to produce Evidence for the current work.
    Requires: gate.run
  authority-unavailable for actor:actor-reviewer-two: gate.run
  budget-data-unavailable for envelope:ENV-e2b3c3: tokens
  Action claim.renew claim:CLM-85890b: blocked
    Purpose: Extend an eligible work lease after checking authority.
    Requires: claim.renew
  authority-unavailable for actor:actor-reviewer-two: claim.renew
  Action claim.release claim:CLM-85890b: ready
    Purpose: Give up an eligible work lease after checking authority; this does not prove the work complete.
    Requires: claim.release
  Action concern.record issue:P22-764499: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: concern.record
  authority-unavailable for issue:P22-764499: concern.record
  Action deviation.record issue:P22-764499: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: deviation.record
  authority-unavailable for issue:P22-764499: deviation.record
  Action decision.record-reference issue:P22-764499: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: decision.record-reference
  authority-unavailable for issue:P22-764499: decision.record-reference
  Action human-outcome.record issue:P22-764499: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: human-outcome.record
  authority-unavailable for issue:P22-764499: human-outcome.record
  Action review.record run:RUN-7fb8a5: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: review.record
  review-not-independent for run:RUN-7fb8a5: Run already reviewed

Completion requirements
  These describe what is needed to finish the Issue, separately from recovery checks above.
  Obligation OBL-p22-corpus: unsatisfied
    Evidence validity: invalid; freshness: stale; Artifacts: missing
  obligation-evidence-absent for obligation:OBL-p22-corpus: qualifying Evidence is absent
  Overall completion: incomplete
  obligation-evidence-absent for obligation:OBL-p22-corpus: qualifying Evidence is absent
  issue-not-terminal-success for issue:P22-764499: completed workflow category
  Open Findings: 3

Review and judgment
Reviews
  Review REV-367292
    Issue: P22-764499; Envelope: ENV-e2b3c3 revision 1
    Run: RUN-0b264e; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-7d0d9392020521f7351a9197fad2ac4d; worktree sha256:a4936266d2d82b477ba646de792bf42fcbe02de77cca3c2481a83b1ab2758b10; HEAD 114d9451b3f51394221c0839145ca411370828b2; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base 114d9451b3f51394221c0839145ca411370828b2; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 3; minor 0; nit 0
  Finding FND-fixture-addressed
    Review: REV-367292; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
  Finding FND-fixture-blocking
    Review: REV-367292; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-resolved
    Review: REV-367292; Severity: major; Blocking: true; Disposition: resolved
  Review REV-b90c1c
    Issue: P22-764499; Envelope: ENV-e2b3c3 revision 1
    Run: RUN-7fb8a5; Reviewer: actor-reviewer-two; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-7d0d9392020521f7351a9197fad2ac4d; worktree sha256:a4936266d2d82b477ba646de792bf42fcbe02de77cca3c2481a83b1ab2758b10; HEAD 114d9451b3f51394221c0839145ca411370828b2; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base 114d9451b3f51394221c0839145ca411370828b2; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-fixture-new
    Review: REV-b90c1c; Severity: major; Blocking: true; Disposition: open
Open Findings
  Finding FND-fixture-addressed
    Review: REV-367292; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
  Finding FND-fixture-blocking
    Review: REV-367292; Severity: major; Blocking: true; Disposition: open
  Finding FND-fixture-new
    Review: REV-b90c1c; Severity: major; Blocking: true; Disposition: open
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

- JSON counterpart for read 7: cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings", "--json", "resume", "P22-764499", "--repository", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings --json resume P22-764499 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/mixed-findings
```

## VAL-071 protected Decision

- Project: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision`
- IDs: `{"envelope_id": "ENV-415cac", "human_actor_id": "actor-human-outcome", "issue_id": "P22-8c4994"}`

### Read 1

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision`
- Arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision", "resume", "P22-8c4994", "--repository", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision resume P22-8c4994 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision
```

Verbatim Human stdout:

```text
Issue P22-8c4994
  Status: qa_review
  Completion: incomplete
  Envelope: ENV-415cac (revision 1, active)

Current Run: RUN-160063
  Status: active
  Actor ID: actor-recovery
  Executor ID: executor-corpus-prior

Claim CLM-db623c
  Lease status now: active
  Recorded status: active
  Expires: 2026-09-29T08:58:34.681176109Z
  An active lease alone does not authorize an operation; its authority and repository checks still apply.

Repository
  Project: project-ef54671165fe5e7b33f42bcf99980981
  Current worktree: sha256:1577c216515eda0b4fbd6b99708bace576cf96a92efc79d7577b620473738ea9
  HEAD: 8c8270f04f0c3f465220c5b2b178ccca658ee5b9
  Uncommitted changes: false
  Run's recorded worktree: sha256:1577c216515eda0b4fbd6b99708bace576cf96a92efc79d7577b620473738ea9

Next steps
  To collect proof for unfinished work, run the configured Gate (gate.run).
  A successful execution can supply Evidence; it does not itself complete the Issue.
  Candidate operations below are not an execution sequence.
  Ready means the recorded prerequisites permit a request; supplied inputs and authority are checked when it runs.
  Action run.inspect run:RUN-160063: ready
    Purpose: Compare this repository with the Run's recorded start or checkpoint; inspect changed paths without recording a checkpoint.
  Action run.observe run:RUN-160063: ready
    Purpose: Have the activated harness record a fresh observation of whether this Run's executor is live, terminated or unknown.
    Requires: run.observe
  Action gate.run run:RUN-160063: ready
    Purpose: Execute the configured check to produce Evidence for the current work.
    Requires: gate.run
  Action claim.renew claim:CLM-db623c: blocked
    Purpose: Extend an eligible work lease after checking authority.
    Requires: claim.renew
  authority-unavailable for actor:actor-recovery: claim.renew
  Action claim.release claim:CLM-db623c: ready
    Purpose: Give up an eligible work lease after checking authority; this does not prove the work complete.
    Requires: claim.release
  Action concern.record issue:P22-8c4994: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: concern.record
  authority-unavailable for issue:P22-8c4994: concern.record
  Action deviation.record issue:P22-8c4994: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: deviation.record
  authority-unavailable for issue:P22-8c4994: deviation.record
  Action decision.record-reference issue:P22-8c4994: ready
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: decision.record-reference
  Action human-outcome.record issue:P22-8c4994: ready
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: human-outcome.record

Blockers and attention
Stop STP-b4168f blocks issue.transition-terminal-success because protected-decision-required: {"affected_actions":["issue.transition-terminal-success"],"decision_class":"security","envelope_id":"ENV-415cac","required_capabilities":["decision.record-reference"]}; qualifying resolution: decision-satisfied for envelope:ENV-415cac decision class security (resolver capability stop.resolve).
Attention ATT-ede39c: open; protected-decision-required; linked Stops: STP-b4168f

Completion requirements
  These describe what is needed to finish the Issue, separately from recovery checks above.
  Obligation OBL-p22-corpus: satisfied
    Evidence validity: valid; freshness: fresh; Artifacts: available
  Overall completion: incomplete
  issue-not-terminal-success for issue:P22-8c4994: completed workflow category
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

- JSON counterpart for read 1: cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision", "--json", "resume", "P22-8c4994", "--repository", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision --json resume P22-8c4994 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision
```

### Read 3

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision`
- Arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision", "decision", "list", "--issue", "P22-8c4994"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision decision list --issue P22-8c4994
```

Verbatim Human stdout:

```text
Purpose: Inspect Decision References and their authority.

Decision References
  Decision Reference DREF-fixture-security
    Class: security; Outcome: approved; Authority Actor: actor-recovery
    Required capabilities: decision.record-reference
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

- JSON counterpart for read 3: cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision", "--json", "decision", "list", "--issue", "P22-8c4994"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision --json decision list --issue P22-8c4994
```

### Read 5

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision`
- Arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision", "stop", "list", "--subject", "envelope:ENV-415cac"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision stop list --subject envelope:ENV-415cac
```

Verbatim Human stdout:

```text
stop.list
Stop STP-b4168f blocks issue.transition-terminal-success because protected-decision-required: {"affected_actions":["issue.transition-terminal-success"],"decision_class":"security","envelope_id":"ENV-415cac","required_capabilities":["decision.record-reference"]}; qualifying resolution: decision-satisfied for envelope:ENV-415cac decision class security (resolver capability stop.resolve).
```

- JSON counterpart for read 5: cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision", "--json", "stop", "list", "--subject", "envelope:ENV-415cac"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision --json stop list --subject envelope:ENV-415cac
```

### Human-run steps

Run these in order after recording the initial reading. The supplied Outcome reasoning is fixture data for Human review.

- blocked action before Outcome: expected exit `5`; cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision", "--json", "issue", "transition", "P22-8c4994", "--to", "qa_pass", "--actor", "agent:actor-recovery", "--expected-revision", "4", "--idempotency-key", "fixture-human-transition"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision --json issue transition P22-8c4994 --to qa_pass --actor agent:actor-recovery --expected-revision 4 --idempotency-key fixture-human-transition
```

- record authorized Human Outcome: expected exit `0`; cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision", "--json", "human-outcome", "record", "--input", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision/fixture-human-outcome.json", "--actor", "actor-human-outcome", "--idempotency-key", "fixture-human-outcome"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision --json human-outcome record --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision/fixture-human-outcome.json --actor actor-human-outcome --idempotency-key fixture-human-outcome
```

- retry guarded action: expected exit `0`; cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision", "--json", "issue", "transition", "P22-8c4994", "--to", "qa_pass", "--actor", "agent:actor-recovery", "--expected-revision", "4", "--idempotency-key", "fixture-human-transition"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/protected-decision --json issue transition P22-8c4994 --to qa_pass --actor agent:actor-recovery --expected-revision 4 --idempotency-key fixture-human-transition
```

## VAL-072 review rework exhaustion

The verbatim readback below is the original Human-checklist snapshot. For the
DEC-269 re-check, regenerate the fixture with `TestP23HumanFixtures`; its
`review-exhaustion` checkpoint has the amended `human-outcome-recorded`
Condition and the following supported resolution sequence. The generator
places `fixture-rework-outcome.json` in that checkpoint's project directory
and includes `actor-human-outcome` as an activated Human Actor.

- Project: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion`
- IDs: `{"issue_id": "P22-4cf573", "review_id": "REV-3d0854"}`

### Read 1

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion`
- Arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion", "resume", "P22-4cf573", "--repository", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion resume P22-4cf573 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion
```

Verbatim Human stdout:

```text
Issue P22-4cf573
  Status: todo
  Completion: incomplete
  Envelope: ENV-4fc37e (revision 1, active)

Current Run: RUN-6e8ecc
  Status: failed
  Actor ID: actor-denied
  Executor ID: executor-fixture-cycle

Current Claim: none.

Repository
  Project: project-22109211fb77fdabe2af451f37766460
  Current worktree: sha256:2d435a306d1fc326f44f7403e8ef32c218816ccf6689fc14b9a6544aa7c42762
  HEAD: d1a6fc27879b9a4d2557e5231d3a8f414ad0c939
  Uncommitted changes: false
  Run's recorded worktree: sha256:2d435a306d1fc326f44f7403e8ef32c218816ccf6689fc14b9a6544aa7c42762

Next steps
  Candidate operations below are not an execution sequence.
  Ready means the recorded prerequisites permit a request; supplied inputs and authority are checked when it runs.
  Action run.inspect run:RUN-6e8ecc: ready
    Purpose: Compare this repository with the Run's recorded start or checkpoint; inspect changed paths without recording a checkpoint.
  Action run.observe run:RUN-6e8ecc: blocked
    Purpose: Have the activated harness record a fresh observation of whether this Run's executor is live, terminated or unknown.
    Requires: run.observe
  run-inactive for run:RUN-6e8ecc: failed
  Action gate.run run:RUN-6e8ecc: blocked
    Purpose: Execute the configured check to produce Evidence for the current work.
    Requires: gate.run
  authority-unavailable for actor:actor-denied: gate.run
  run-inactive for run:RUN-6e8ecc: failed
  budget-exhausted for envelope:ENV-4fc37e: review-cycles
  budget-data-unavailable for envelope:ENV-4fc37e: tokens
  Action concern.record issue:P22-4cf573: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: concern.record
  authority-unavailable for issue:P22-4cf573: concern.record
  Action deviation.record issue:P22-4cf573: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: deviation.record
  authority-unavailable for issue:P22-4cf573: deviation.record
  Action decision.record-reference issue:P22-4cf573: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: decision.record-reference
  authority-unavailable for issue:P22-4cf573: decision.record-reference
  Action human-outcome.record issue:P22-4cf573: blocked
    Purpose: Inspect the operation's requirements before submitting a request.
    Requires: human-outcome.record
  authority-unavailable for issue:P22-4cf573: human-outcome.record

Blockers and attention
Stop STP-5b10b2 blocks claim.acquire, run.start because review-rework-exhausted: {"budget_name":"review-cycles","consumed":1,"envelope_id":"ENV-4fc37e","last_review_id":"REV-3d0854","maximum":1,"remaining":0,"unit":"review-cycles"}; qualifying resolution: budget-observed for envelope:ENV-4fc37e budget review-cycles (resolver capability stop.resolve).
Attention ATT-00c589: open; review-rework-exhausted; linked Stops: STP-5b10b2

Completion requirements
  These describe what is needed to finish the Issue, separately from recovery checks above.
  Obligation OBL-p22-corpus: unsatisfied
    Evidence validity: invalid; freshness: stale; Artifacts: not-applicable
  obligation-evidence-absent for obligation:OBL-p22-corpus: qualifying Evidence is absent
  Overall completion: incomplete
  obligation-evidence-absent for obligation:OBL-p22-corpus: qualifying Evidence is absent
  issue-not-terminal-success for issue:P22-4cf573: completed workflow category
  Open Findings: 1

Review and judgment
Reviews
  Review REV-3d0854
    Issue: P22-4cf573; Envelope: ENV-4fc37e revision 1
    Run: RUN-6e8ecc; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-22109211fb77fdabe2af451f37766460; worktree sha256:2d435a306d1fc326f44f7403e8ef32c218816ccf6689fc14b9a6544aa7c42762; HEAD d1a6fc27879b9a4d2557e5231d3a8f414ad0c939; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base d1a6fc27879b9a4d2557e5231d3a8f414ad0c939; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-fixture-cycle
    Review: REV-3d0854; Severity: major; Blocking: true; Disposition: open
Open Findings
  Finding FND-fixture-cycle
    Review: REV-3d0854; Severity: major; Blocking: true; Disposition: open
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

- JSON counterpart for read 1: cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion", "--json", "resume", "P22-4cf573", "--repository", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion --json resume P22-4cf573 --repository /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion
```

### Read 3

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion`
- Arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion", "stop", "list", "--subject", "envelope:ENV-4fc37e"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion stop list --subject envelope:ENV-4fc37e
```

Verbatim Human stdout:

```text
stop.list
Stop STP-5b10b2 blocks claim.acquire, run.start because review-rework-exhausted: {"budget_name":"review-cycles","consumed":1,"envelope_id":"ENV-4fc37e","last_review_id":"REV-3d0854","maximum":1,"remaining":0,"unit":"review-cycles"}; qualifying resolution: budget-observed for envelope:ENV-4fc37e budget review-cycles (resolver capability stop.resolve).
```

- JSON counterpart for read 3: cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion", "--json", "stop", "list", "--subject", "envelope:ENV-4fc37e"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion --json stop list --subject envelope:ENV-4fc37e
```

### Read 5

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion`
- Arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion", "attention", "list", "--subject", "envelope:ENV-4fc37e"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion attention list --subject envelope:ENV-4fc37e
```

Verbatim Human stdout:

```text
attention.list
Attention ATT-00c589 open for review-rework-exhausted; qualifying resolution requires dispositions for every linked Stop:
  Stop STP-5b10b2 blocks claim.acquire, run.start because review-rework-exhausted: {"budget_name":"review-cycles","consumed":1,"envelope_id":"ENV-4fc37e","last_review_id":"REV-3d0854","maximum":1,"remaining":0,"unit":"review-cycles"}; qualifying resolution: budget-observed for envelope:ENV-4fc37e budget review-cycles (resolver capability stop.resolve).
```

- JSON counterpart for read 5: cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion", "--json", "attention", "list", "--subject", "envelope:ENV-4fc37e"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion --json attention list --subject envelope:ENV-4fc37e
```

### Human-run steps

Run these in order after recording the initial reading. The supplied Outcome reasoning is fixture data for Human review.

- exhausted rework denial: expected exit `4`; cwd `/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion`; args `["--project", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion", "--json", "claim", "acquire", "--input", "/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion/fixture-rework-claim.json", "--actor", "actor-recovery", "--idempotency-key", "fixture-exhausted-claim"]`

```bash
/home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/fid --project /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion --json claim acquire --input /home/mark/projects/fid/.cache/phase-23-human/checkpoints-1365380666/review-exhaustion/fixture-rework-claim.json --actor actor-recovery --idempotency-key fixture-exhausted-claim
```

In the regenerated checkpoint, first inspect `fid stop list --subject
envelope:<envelope-id>`. Then record the supplied approved
`review-rework-extension` Outcome with `fid human-outcome record --actor
actor-human-outcome --input fixture-rework-outcome.json --idempotency-key
fixture-rework-outcome`. Retry the same `claim acquire` command above against
the regenerated checkpoint; it now exits `0` and returns an admitted Claim.

## Review prompts

- VAL-070: Identify the resolved Finding, the merely addressed Finding, and what still blocks completion from fid output.
- VAL-071: Run the blocked action, inspect the Stop, record the Human Outcome with the listed human Actor, and retry. Explain why permission changed.
- VAL-072: Inspect the rework denial, Stop, and Attention. Identify the exhausted budget and exact qualifying resolution.

Record Human judgments in the phase-end checklist. Keep `.fid/` internals out of the inspection.

## 2026-09-30 QA repair — review-11 MF-1 every-leaf Human/JSON pairs

The original three VAL-070..072 Human checkpoints above remain available. This additional retained workspace exercises all 18 Phase 23 leaves through the built `fid` process. Each Human block below is byte-for-byte stdout. Each JSON block contains the result facts returned by the paired `--json` command, with free-form project prose fields omitted from the guide. The tagged test compares the full JSON result, reconstructs every Human fact line from its JSON fields, and reruns both modes after the generator exits. The ordinary test pins exact text for all 18 operation names against JSON-decoded facts.

- Manifest: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/manifest.json`
- Binary: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/fid`
- Binary SHA-256: `sha256:f11fa4c7873600e34ab5310fd2e059349eb279bbb02a196323dff6ae9b3a5a6d`
- Source commit: `44556095af20058b84d93957597456e7741d626f`; source dirty: `true`
- Leaf workspace: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Judgment checks: every Concern, Deviation, and Human Outcome view carries its non-proof qualifier; addressed Findings say that addressed is a Builder assertion; no Human block contains the `PRIVATE_` project prose sentinels seeded in the JSON facts.

### `requirement.record-reference`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "requirement", "record-reference", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-requirement.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-requirement"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "requirement", "record-reference", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-requirement.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-requirement"]`

Verbatim Human stdout:

```text
Purpose: Inspect Requirement References and Issue links.

Requirement References and links
  Requirement RREF-dad9dc
    Status: active
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON result facts for comparison (project prose omitted):

```json
{
  "idempotency_key": "leaf-requirement",
  "replayed": true,
  "requirement": {
    "created_at": "2026-09-29T23:39:27.851042559Z",
    "event_id": "EVT-20260929-000018",
    "id": "RREF-dad9dc",
    "source": {
      "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
      "display_label": "Task 9 decision",
      "external_id": null,
      "path": "base.txt",
      "repository": {
        "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
        "capture_source": "git",
        "dirty": false,
        "dirty_digest": null,
        "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
        "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
        "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
        "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
      },
      "section": null,
      "source_type": "project-file"
    },
    "status": "active",
    "superseded_by_id": null
  }
}
```

### `requirement.link`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "requirement", "link", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-link.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-link"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "requirement", "link", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-link.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-link"]`

Verbatim Human stdout:

```text
Purpose: Inspect Requirement References and Issue links.

Requirement References and links
  Requirement REQ-leaf
    Requirement: REQ-leaf; Issue: P22-b8865b; Relation: implements
```

JSON result facts for comparison (project prose omitted):

```json
{
  "idempotency_key": "leaf-link",
  "link": {
    "created_at": "2026-09-29T23:39:27.991208446Z",
    "event_id": "EVT-20260929-000019",
    "issue_id": "P22-b8865b",
    "relation": "implements",
    "requirement_id": "REQ-leaf"
  },
  "replayed": true
}
```

### `concern.record`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "concern", "record", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-concern.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-concern"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "concern", "record", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-concern.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-concern"]`

Verbatim Human stdout:

```text
Purpose: Inspect a raised concern and its disposition as judgment.

Concerns
  Concern CON-leaf
    Subject: issue:P22-b8865b
    Severity: major; Status: open; Raised by: actor-recovery
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
    This Concern is recorded judgment, not deterministic proof.
```

JSON result facts for comparison (project prose omitted):

```json
{
  "concern": {
    "actor_id": "actor-recovery",
    "created_at": "2026-09-29T23:39:28.149301929Z",
    "disposition": null,
    "event_id": "EVT-20260929-000020",
    "id": "CON-leaf",
    "references": [
      {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
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
      "id": "P22-b8865b",
      "type": "issue"
    }
  },
  "idempotency_key": "leaf-concern",
  "replayed": true
}
```

### `concern.dispose`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "concern", "dispose", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-dispose.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-dispose"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "concern", "dispose", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-dispose.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-dispose"]`

Verbatim Human stdout:

```text
Purpose: Inspect a raised concern and its disposition as judgment.

Concerns
  Concern CON-leaf
    Subject: issue:P22-b8865b
    Severity: major; Status: disposed; Raised by: actor-recovery
    Disposition: carried-forward; Disposed by: actor-recovery
    Work now owned by Issue: P22-b8865b
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
    This Concern is recorded judgment, not deterministic proof.
```

JSON result facts for comparison (project prose omitted):

```json
{
  "concern": {
    "actor_id": "actor-recovery",
    "created_at": "2026-09-29T23:39:28.149301929Z",
    "disposition": {
      "action": "carried-forward",
      "actor_id": "actor-recovery",
      "disposed_at": "2026-09-29T23:39:28.364091417Z",
      "event_id": "EVT-20260929-000021",
      "issue_id": "P22-b8865b",
      "supersedes_concern_id": null,
      "supersedes_finding_id": null
    },
    "event_id": "EVT-20260929-000020",
    "id": "CON-leaf",
    "references": [
      {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
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
      "id": "P22-b8865b",
      "type": "issue"
    }
  },
  "idempotency_key": "leaf-dispose",
  "replayed": true
}
```

### `deviation.record`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "deviation", "record", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-deviation.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-deviation"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "deviation", "record", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-deviation.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-deviation"]`

Verbatim Human stdout:

```text
Purpose: Inspect a recorded deviation as judgment.

Deviations
  Deviation DEV-leaf
    Subject: issue:P22-b8865b
    Disposition: pending-decision; Authority Actor: actor-recovery
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
    This Deviation is recorded judgment, not deterministic proof.
```

JSON result facts for comparison (project prose omitted):

```json
{
  "deviation": {
    "actor_id": "actor-recovery",
    "authority": {
      "actor_id": "actor-recovery",
      "reference": null
    },
    "created_at": "2026-09-29T23:39:28.528957533Z",
    "disposition": "pending-decision",
    "event_id": "EVT-20260929-000022",
    "id": "DEV-leaf",
    "references": [
      {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
        },
        "section": null,
        "source_type": "project-file"
      }
    ],
    "subject": {
      "id": "P22-b8865b",
      "type": "issue"
    },
    "supersedes_deviation_id": null
  },
  "idempotency_key": "leaf-deviation",
  "replayed": true
}
```

### `decision.record-reference`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "decision", "record-reference", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-decision.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-decision"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "decision", "record-reference", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-decision.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-decision"]`

Verbatim Human stdout:

```text
Purpose: Inspect Decision References and their authority.

Decision References
  Decision Reference DREF-leaf
    Class: security; Outcome: approved; Authority Actor: actor-recovery
    Required capabilities: decision.record-reference
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON result facts for comparison (project prose omitted):

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
    "created_at": "2026-09-29T23:39:28.741456268Z",
    "event_id": "EVT-20260929-000023",
    "id": "DREF-leaf",
    "outcome": "approved",
    "scope": {
      "issue_id": "P22-b8865b"
    },
    "source": {
      "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
      "display_label": "Task 9 decision",
      "external_id": null,
      "path": "base.txt",
      "repository": {
        "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
        "capture_source": "git",
        "dirty": false,
        "dirty_digest": null,
        "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
        "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
        "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
        "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
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

### `human-outcome.record`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "human-outcome", "record", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-outcome.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-outcome"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "human-outcome", "record", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-outcome.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-outcome"]`

Verbatim Human stdout:

```text
Purpose: Inspect authorized Human judgments and their provenance.

Human Outcomes
  Human Outcome OUT-leaf
    Actor: actor-recovery; Class: security; Result: approved
    Activation: ACT-baf0ec revision 10; Profile: activation:p22-corpus
    This Outcome is authorized Human judgment, not deterministic proof.
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON result facts for comparison (project prose omitted):

```json
{
  "attention": [],
  "idempotency_key": "leaf-outcome",
  "outcome": {
    "actor_id": "actor-recovery",
    "class": "security",
    "created_at": "2026-09-29T23:39:28.940351293Z",
    "event_id": "EVT-20260929-000024",
    "id": "OUT-leaf",
    "provenance": {
      "activation_id": "ACT-baf0ec",
      "activation_version": 10,
      "profile_ref": "activation:p22-corpus"
    },
    "references": [
      {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
        },
        "section": null,
        "source_type": "project-file"
      }
    ],
    "result": "approved",
    "scope": {
      "issue_id": "P22-b8865b"
    },
    "subject": {
      "id": "P22-b8865b",
      "type": "issue"
    }
  },
  "protected_decisions": [],
  "replayed": true,
  "stops": []
}
```

### `review.record`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "review", "record", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-review.json", "--actor", "actor-denied", "--idempotency-key", "leaf-review"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "review", "record", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-review.json", "--actor", "actor-denied", "--idempotency-key", "leaf-review"]`

Verbatim Human stdout:

```text
Purpose: Record an independent Review of a Run and its bound Issue.

Reviews
  Review REV-93ace4
    Issue: P22-b8865b; Envelope: ENV-c45f4a revision 1
    Run: RUN-8664f6; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-c38e61c97b0200742aa4a3e766c8c02e; worktree sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600; HEAD 02e5a80e8429bc8f939cae2e8fdb691ad7571947; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base 02e5a80e8429bc8f939cae2e8fdb691ad7571947; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-leaf
    Review: REV-93ace4; Severity: major; Blocking: true; Disposition: open
```

JSON result facts for comparison (project prose omitted):

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
              "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
              "capture_source": "git",
              "dirty": false,
              "dirty_digest": null,
              "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
              "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
              "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
              "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
            },
            "section": null,
            "source_type": "project-file"
          }
        }
      ]
    },
    "envelope_id": "ENV-c45f4a",
    "envelope_revision": 1,
    "event_id": "EVT-20260929-000030",
    "evidence_refs": [],
    "findings": [
      {
        "blocking": true,
        "created_at": "2026-09-29T23:39:29.730679427Z",
        "dispositions": [],
        "id": "FND-leaf",
        "review_id": "REV-93ace4",
        "scope": {
          "paths": []
        },
        "severity": "major",
        "state": "open",
        "subject": {
          "id": "P22-b8865b",
          "type": "issue"
        }
      }
    ],
    "id": "REV-93ace4",
    "independence": "distinct-actor",
    "issue_id": "P22-b8865b",
    "profile": "review-default",
    "receipt_ids": [],
    "recorded_at": "2026-09-29T23:39:29.730679427Z",
    "repository": {
      "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
      "capture_source": "git",
      "dirty": false,
      "dirty_digest": null,
      "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
      "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
      "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
      "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
    },
    "resolved_finding_ids": [],
    "run_id": "RUN-8664f6",
    "verdict": "changes-required"
  },
  "stops": []
}
```

### `finding.address`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "finding", "address", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-address.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-address"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "finding", "address", "--input", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract/leaf-leaf-address.json", "--actor", "actor-recovery", "--idempotency-key", "leaf-address"]`

Verbatim Human stdout:

```text
Purpose: Record a Builder assertion that a Finding was addressed.

Findings
  Finding FND-leaf
    Review: REV-93ace4; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
```

JSON result facts for comparison (project prose omitted):

```json
{
  "finding": {
    "blocking": true,
    "created_at": "2026-09-29T23:39:29.730679427Z",
    "dispositions": [
      {
        "action": "addressed",
        "actor_id": "actor-recovery",
        "disposed_at": "2026-09-29T23:39:30.408258957Z",
        "event_id": "EVT-20260929-000035",
        "id": "FDP-0baf38",
        "qualifying_review_id": null
      }
    ],
    "id": "FND-leaf",
    "review_id": "REV-93ace4",
    "scope": {
      "paths": []
    },
    "severity": "major",
    "state": "addressed",
    "subject": {
      "id": "P22-b8865b",
      "type": "issue"
    }
  },
  "idempotency_key": "leaf-address",
  "replayed": true
}
```

### `review.get`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "review", "get", "REV-93ace4"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "review", "get", "REV-93ace4"]`

Verbatim Human stdout:

```text
Purpose: Inspect independent Reviews and their Findings.

Reviews
  Review REV-93ace4
    Issue: P22-b8865b; Envelope: ENV-c45f4a revision 1
    Run: RUN-8664f6; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-c38e61c97b0200742aa4a3e766c8c02e; worktree sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600; HEAD 02e5a80e8429bc8f939cae2e8fdb691ad7571947; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base 02e5a80e8429bc8f939cae2e8fdb691ad7571947; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-leaf
    Review: REV-93ace4; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
```

JSON result facts for comparison (project prose omitted):

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
              "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
              "capture_source": "git",
              "dirty": false,
              "dirty_digest": null,
              "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
              "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
              "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
              "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
            },
            "section": null,
            "source_type": "project-file"
          }
        }
      ]
    },
    "envelope_id": "ENV-c45f4a",
    "envelope_revision": 1,
    "event_id": "EVT-20260929-000030",
    "evidence_refs": [],
    "findings": [
      {
        "blocking": true,
        "created_at": "2026-09-29T23:39:29.730679427Z",
        "dispositions": [
          {
            "action": "addressed",
            "actor_id": "actor-recovery",
            "disposed_at": "2026-09-29T23:39:30.408258957Z",
            "event_id": "EVT-20260929-000035",
            "id": "FDP-0baf38",
            "qualifying_review_id": null
          }
        ],
        "id": "FND-leaf",
        "review_id": "REV-93ace4",
        "scope": {
          "paths": []
        },
        "severity": "major",
        "state": "addressed",
        "subject": {
          "id": "P22-b8865b",
          "type": "issue"
        }
      }
    ],
    "id": "REV-93ace4",
    "independence": "distinct-actor",
    "issue_id": "P22-b8865b",
    "profile": "review-default",
    "receipt_ids": [],
    "recorded_at": "2026-09-29T23:39:29.730679427Z",
    "repository": {
      "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
      "capture_source": "git",
      "dirty": false,
      "dirty_digest": null,
      "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
      "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
      "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
      "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
    },
    "resolved_finding_ids": [],
    "run_id": "RUN-8664f6",
    "verdict": "changes-required"
  }
}
```

### `review.list`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "review", "list", "--issue", "P22-b8865b"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "review", "list", "--issue", "P22-b8865b"]`

Verbatim Human stdout:

```text
Purpose: Inspect independent Reviews and their Findings.

Reviews
  Review REV-93ace4
    Issue: P22-b8865b; Envelope: ENV-c45f4a revision 1
    Run: RUN-8664f6; Reviewer: actor-denied; Profile: review-default
    Verdict: changes-required; Independence: distinct-actor
    Repository: project project-c38e61c97b0200742aa4a3e766c8c02e; worktree sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600; HEAD 02e5a80e8429bc8f939cae2e8fdb691ad7571947; tree 720e5c7e58f0ed980bdde2d8743e3755ef4a8503
    Repository capture: git; base 02e5a80e8429bc8f939cae2e8fdb691ad7571947; dirty false; dirty digest <nil>
    Findings by severity: blocker 0; major 1; minor 0; nit 0
  Finding FND-leaf
    Review: REV-93ace4; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
```

JSON result facts for comparison (project prose omitted):

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
                "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
                "capture_source": "git",
                "dirty": false,
                "dirty_digest": null,
                "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
                "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
                "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
                "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
              },
              "section": null,
              "source_type": "project-file"
            }
          }
        ]
      },
      "envelope_id": "ENV-c45f4a",
      "envelope_revision": 1,
      "event_id": "EVT-20260929-000030",
      "evidence_refs": [],
      "findings": [
        {
          "blocking": true,
          "created_at": "2026-09-29T23:39:29.730679427Z",
          "dispositions": [
            {
              "action": "addressed",
              "actor_id": "actor-recovery",
              "disposed_at": "2026-09-29T23:39:30.408258957Z",
              "event_id": "EVT-20260929-000035",
              "id": "FDP-0baf38",
              "qualifying_review_id": null
            }
          ],
          "id": "FND-leaf",
          "review_id": "REV-93ace4",
          "scope": {
            "paths": []
          },
          "severity": "major",
          "state": "addressed",
          "subject": {
            "id": "P22-b8865b",
            "type": "issue"
          }
        }
      ],
      "id": "REV-93ace4",
      "independence": "distinct-actor",
      "issue_id": "P22-b8865b",
      "profile": "review-default",
      "receipt_ids": [],
      "recorded_at": "2026-09-29T23:39:29.730679427Z",
      "repository": {
        "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
        "capture_source": "git",
        "dirty": false,
        "dirty_digest": null,
        "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
        "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
        "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
        "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
      },
      "resolved_finding_ids": [],
      "run_id": "RUN-8664f6",
      "verdict": "changes-required"
    }
  ]
}
```

### `finding.list`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "finding", "list", "--issue", "P22-b8865b"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "finding", "list", "--issue", "P22-b8865b"]`

Verbatim Human stdout:

```text
Purpose: Inspect Finding dispositions and current blocking state.

Findings
  Finding FND-leaf
    Review: REV-93ace4; Severity: major; Blocking: true; Disposition: addressed
    Addressed is a Builder assertion, not a resolution.
```

JSON result facts for comparison (project prose omitted):

```json
{
  "items": [
    {
      "blocking": true,
      "created_at": "2026-09-29T23:39:29.730679427Z",
      "dispositions": [
        {
          "action": "addressed",
          "actor_id": "actor-recovery",
          "disposed_at": "2026-09-29T23:39:30.408258957Z",
          "event_id": "EVT-20260929-000035",
          "id": "FDP-0baf38",
          "qualifying_review_id": null
        }
      ],
      "id": "FND-leaf",
      "review_id": "REV-93ace4",
      "scope": {
        "paths": []
      },
      "severity": "major",
      "state": "addressed",
      "subject": {
        "id": "P22-b8865b",
        "type": "issue"
      }
    }
  ]
}
```

### `concern.list`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "concern", "list", "--subject", "issue:P22-b8865b"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "concern", "list", "--subject", "issue:P22-b8865b"]`

Verbatim Human stdout:

```text
Purpose: Inspect a raised concern and its disposition as judgment.

Concerns
  Concern CON-leaf
    Subject: issue:P22-b8865b
    Severity: major; Status: disposed; Raised by: actor-recovery
    Disposition: carried-forward; Disposed by: actor-recovery
    Work now owned by Issue: P22-b8865b
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
    This Concern is recorded judgment, not deterministic proof.
```

JSON result facts for comparison (project prose omitted):

```json
{
  "items": [
    {
      "actor_id": "actor-recovery",
      "created_at": "2026-09-29T23:39:28.149301929Z",
      "disposition": {
        "action": "carried-forward",
        "actor_id": "actor-recovery",
        "disposed_at": "2026-09-29T23:39:28.364091417Z",
        "event_id": "EVT-20260929-000021",
        "issue_id": "P22-b8865b",
        "supersedes_concern_id": null,
        "supersedes_finding_id": null
      },
      "event_id": "EVT-20260929-000020",
      "id": "CON-leaf",
      "references": [
        {
          "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
          "display_label": "Task 9 decision",
          "external_id": null,
          "path": "base.txt",
          "repository": {
            "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
            "capture_source": "git",
            "dirty": false,
            "dirty_digest": null,
            "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
            "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
            "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
            "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
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
        "id": "P22-b8865b",
        "type": "issue"
      }
    }
  ]
}
```

### `deviation.list`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "deviation", "list", "--subject", "issue:P22-b8865b"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "deviation", "list", "--subject", "issue:P22-b8865b"]`

Verbatim Human stdout:

```text
Purpose: Inspect a recorded deviation as judgment.

Deviations
  Deviation DEV-leaf
    Subject: issue:P22-b8865b
    Disposition: pending-decision; Authority Actor: actor-recovery
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
    This Deviation is recorded judgment, not deterministic proof.
```

JSON result facts for comparison (project prose omitted):

```json
{
  "items": [
    {
      "actor_id": "actor-recovery",
      "authority": {
        "actor_id": "actor-recovery",
        "reference": null
      },
      "created_at": "2026-09-29T23:39:28.528957533Z",
      "disposition": "pending-decision",
      "event_id": "EVT-20260929-000022",
      "id": "DEV-leaf",
      "references": [
        {
          "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
          "display_label": "Task 9 decision",
          "external_id": null,
          "path": "base.txt",
          "repository": {
            "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
            "capture_source": "git",
            "dirty": false,
            "dirty_digest": null,
            "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
            "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
            "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
            "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
          },
          "section": null,
          "source_type": "project-file"
        }
      ],
      "subject": {
        "id": "P22-b8865b",
        "type": "issue"
      },
      "supersedes_deviation_id": null
    }
  ]
}
```

### `decision.list`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "decision", "list", "--issue", "P22-b8865b"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "decision", "list", "--issue", "P22-b8865b"]`

Verbatim Human stdout:

```text
Purpose: Inspect Decision References and their authority.

Decision References
  Decision Reference DREF-leaf
    Class: security; Outcome: approved; Authority Actor: actor-recovery
    Required capabilities: decision.record-reference
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON result facts for comparison (project prose omitted):

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
      "created_at": "2026-09-29T23:39:28.741456268Z",
      "event_id": "EVT-20260929-000023",
      "id": "DREF-leaf",
      "outcome": "approved",
      "scope": {
        "issue_id": "P22-b8865b"
      },
      "source": {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
        },
        "section": null,
        "source_type": "project-file"
      }
    }
  ]
}
```

### `human-outcome.list`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "human-outcome", "list", "--issue", "P22-b8865b"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "human-outcome", "list", "--issue", "P22-b8865b"]`

Verbatim Human stdout:

```text
Purpose: Inspect authorized Human judgments and their provenance.

Human Outcomes
  Human Outcome OUT-leaf
    Actor: actor-recovery; Class: security; Result: approved
    Activation: ACT-baf0ec revision 10; Profile: activation:p22-corpus
    This Outcome is authorized Human judgment, not deterministic proof.
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON result facts for comparison (project prose omitted):

```json
{
  "items": [
    {
      "actor_id": "actor-recovery",
      "class": "security",
      "created_at": "2026-09-29T23:39:28.940351293Z",
      "event_id": "EVT-20260929-000024",
      "id": "OUT-leaf",
      "provenance": {
        "activation_id": "ACT-baf0ec",
        "activation_version": 10,
        "profile_ref": "activation:p22-corpus"
      },
      "references": [
        {
          "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
          "display_label": "Task 9 decision",
          "external_id": null,
          "path": "base.txt",
          "repository": {
            "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
            "capture_source": "git",
            "dirty": false,
            "dirty_digest": null,
            "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
            "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
            "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
            "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
          },
          "section": null,
          "source_type": "project-file"
        }
      ],
      "result": "approved",
      "scope": {
        "issue_id": "P22-b8865b"
      },
      "subject": {
        "id": "P22-b8865b",
        "type": "issue"
      }
    }
  ]
}
```

### `requirement.get`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "requirement", "get", "REQ-leaf"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "requirement", "get", "REQ-leaf"]`

Verbatim Human stdout:

```text
Purpose: Inspect Requirement References and Issue links.

Requirement References and links
  Requirement RREF-dad9dc
    Status: active
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
  Requirement REQ-leaf
    Requirement: REQ-leaf; Issue: P22-b8865b; Relation: implements
```

JSON result facts for comparison (project prose omitted):

```json
{
  "links": [
    {
      "created_at": "2026-09-29T23:39:27.991208446Z",
      "event_id": "EVT-20260929-000019",
      "issue_id": "P22-b8865b",
      "relation": "implements",
      "requirement_id": "REQ-leaf"
    }
  ],
  "requirement": {
    "created_at": "2026-09-29T23:39:27.851042559Z",
    "event_id": "EVT-20260929-000018",
    "id": "RREF-dad9dc",
    "source": {
      "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
      "display_label": "Task 9 decision",
      "external_id": null,
      "path": "base.txt",
      "repository": {
        "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
        "capture_source": "git",
        "dirty": false,
        "dirty_digest": null,
        "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
        "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
        "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
        "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
      },
      "section": null,
      "source_type": "project-file"
    },
    "status": "active",
    "superseded_by_id": null
  }
}
```

### `requirement.list`

- Working directory: `/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract`
- Human arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "requirement", "list"]`
- JSON arguments: `["--project", "/home/mark/projects/fid/.cache/phase-23-human-repair/checkpoints-944862472/leaf-contract", "--json", "requirement", "list"]`

Verbatim Human stdout:

```text
Purpose: Inspect Requirement References and Issue links.

Requirement References and links
  Requirement RREF-dad9dc
    Status: active
    Reference: Task 9 decision; path base.txt; digest sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac
```

JSON result facts for comparison (project prose omitted):

```json
{
  "items": [
    {
      "created_at": "2026-09-29T23:39:27.851042559Z",
      "event_id": "EVT-20260929-000018",
      "id": "RREF-dad9dc",
      "source": {
        "digest": "sha256:f34848ca92665c342abd5816c9e3eda0e82180671195362bcd0080544a3bc2ac",
        "display_label": "Task 9 decision",
        "external_id": null,
        "path": "base.txt",
        "repository": {
          "base_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "capture_source": "git",
          "dirty": false,
          "dirty_digest": null,
          "head_oid": "02e5a80e8429bc8f939cae2e8fdb691ad7571947",
          "project_id": "project-c38e61c97b0200742aa4a3e766c8c02e",
          "tree_oid": "720e5c7e58f0ed980bdde2d8743e3755ef4a8503",
          "worktree_id": "sha256:6f63dd6a121835ceb46aeca0827c0baa6a3794ffe387f33b9d9c9011e928a600"
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
