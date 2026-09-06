# Phase 22 - TL-Lite Triage Log

> Seeded by Tech Lead from the triage-log template at planning kickoff.
> Append-only allocation source for TRG-22 and ESC-22 identifiers.

No triage events have occurred. The first allocated identifiers will be
TRG-22-001 and ESC-22-001, in their separate series. Read this file before
allocation; never recall, reserve in chat, or reuse an identifier.

Only TL-Lite or Tech Lead allocates an ESC. Write its durable entry and the
current escalation-register row before reporting a BLOCK. Ordinary REPAIR
allocates no ESC and requests no Human attention.

Actual entries use `project-docs/templates/triage-log-template.md` and the
current schema-v2 exact typed fields. Do not paste placeholder TRG sections
into this file. Append corrections instead of editing earlier verdicts.
Runner-owned run-ledger state is created only after execution approval.

---

## TRG-22-001 - Task 2 retained live-output golden full-gate failure

- **Trigger:** task-2 review-repair exact full-gate failure in the retained cmd/fid live-output golden.
- **Originating task:** task-2
- **Trigger signature:** `task-2:latent-gate-failure:cmd-fid-live-output-golden`
- **Repair category:** in-scope-implementation
- **Repair authority:** minimum causal scope under DEC-175
- **Source:** `project-docs/docs/phases/phase-22/tasks/task-2-done.md`
- **Candidate:** RECOVERY_REQUIRED
- **Attempt count:** 2
- **Checkpoint:** verified phase HEAD `4dd41c97385f0655140bf56b094528cc4cafe6c1`; phase worktree clean except Runner-owned `project-docs/docs/phases/phase-22/run-ledger.md`; task-2 report preserves the failed review-repair receipt.
- **Verdict:** REPAIR
- **Reasoning:** The required repository-wide gate reproducibly failed in `cmd/fid/TestLiveOutputGoldenFixtures`; isolated/package diagnostics do not clear that receipt. Provenance supplied with the trigger establishes that the test and fixture were introduced by phase-21/task-14 before the Phase 22 base, while task 2 did not displace their live-output subject. DEC-175 nevertheless assigns a latent implementation defect first observable through task 2's required gate to task 2 with the minimum causal source/regression scope. Repairing the live-output golden's deterministic expectation can preserve its accepted framing assertion and the exact full-gate command, result, scope, and tolerance.
- **Local decisions made:**
  - Task 2 owns this latent gate-observed defect under DEC-175; R-14/DEC-186 does not apply because this is neither `cmd/fid/gate_run_test.go` nor a task-8 signal/timeout fixture.
  - The continuation may make only the minimum causal live-output golden repair; it must not alter product Stop/Attention behavior, public interfaces, gate definitions, or assertion semantics.
- **Resolution evidence:** The immutable failed receipt is `go test -race -coverprofile=coverage.out ./...` from `/home/mark/projects/fid`, observed `FAIL` in `cmd/fid/TestLiveOutputGoldenFixtures` (terminated `exit_code: 1`): the full run observed 200 `x` bytes where the retained fixture expected the coverage-child `GOCOVERDIR not set` warning and 134 bytes. The report's later isolated/package diagnostics are diagnostic only and do not replace this receipt.
- **If REPAIR:** owning task task-2; category in-scope-implementation; next action dispatch a same-task continuation limited to `cmd/fid/live_output_test.go`, `cmd/fid/testdata/val-036/live-frames.jsonl`, and the required evidence update in `project-docs/docs/phases/phase-22/tasks/task-2-done.md`; preserve the live-output framing subject, retain the failed receipt and add its mandatory deviation note, then run from `/home/mark/projects/fid` the byte-identical required command `go test -race -coverprofile=coverage.out ./...` to exit 0 with all packages passing, every remaining task-2 required command, and fresh independent review-1 re-review. No ESC.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** N/A
- **Exact verification:** N/A
- **Originating report:** N/A
- **Control repair base:** N/A
- **If BLOCK:** N/A
- **Human attention required:** no
- **References:** `project-docs/docs/phases/phase-22/tasks/task-2-done.md:5-9,73-76,95-105,110-118`; `project-docs/docs/phases/phase-22/tasks/task-2.md:32-34,36-64,94-114`; `project-docs/docs/phases/phase-22/execution-scope.md:8-27`; `project-docs/docs/current/decisions_active.md:49-54`; `project-docs/docs/phases/phase-22/decisions.md:45-54`; DEC-175; DEC-186.

---

## TRG-22-002 - Correction for TRG-22-001 task 2 retained live-output golden full-gate failure

- **Trigger:** task-2 review-repair exact full-gate failure in the retained cmd/fid live-output golden.
- **Originating task:** task-2
- **Trigger signature:** `task-2:latent-gate-failure:cmd-fid-live-output-golden`
- **Repair category:** in-scope-implementation
- **Repair authority:** minimum causal scope
- **Source:** `project-docs/docs/phases/phase-22/tasks/task-2-done.md`
- **Candidate:** recovery correction for TRG-22-001; preserve the original malformed triage record and do not treat it as the governing ledger record.
- **Attempt count:** 2
- **Checkpoint:** verified phase HEAD `4dd41c97385f0655140bf56b094528cc4cafe6c1`; phase worktree clean except Runner-owned `project-docs/docs/phases/phase-22/run-ledger.md`; task-2 report preserves the failed required gate.
- **Verdict:** REPAIR
- **Reasoning:** This labelled correction preserves the real same-task latent gate failure and bounded owner recorded in TRG-22-001. The prior record's Repair authority contains DEC-175, but Phase 22's decisions file does not contain a `### DEC-175:` heading, so that authority value is inadmissible for the schema-v2 ledger. This record replaces only the governing triage disposition with the decision-free authority `minimum causal scope`; it preserves the accepted gate, evidence, repair scope, and protected semantics.
- **Local decisions made:**
  - TRG-22-001 remains unchanged as failed orchestration evidence and is not the governing ledger record.
  - Task 2 owns the minimum causal live-output golden repair; the continuation may not alter product Stop/Attention behavior, public interfaces, gate definitions, results, scopes, tolerances, or assertion semantics.
- **Resolution evidence:** The immutable failed receipt is `go test -race -coverprofile=coverage.out ./...` from `/home/mark/projects/fid`, observed `FAIL` in `cmd/fid/TestLiveOutputGoldenFixtures` (terminated `exit_code: 1`): the full run observed 200 `x` bytes where the retained fixture expectation in that run included a coverage-child `GOCOVERDIR-not-set` warning and 134 observed bytes. Isolated/package diagnostics do not clear the failed receipt.
- **If REPAIR:** owning task task-2; category in-scope-implementation; next action dispatch a same-task continuation limited exactly to `cmd/fid/live_output_test.go`, `cmd/fid/testdata/val-036/live-frames.jsonl`, and the required evidence update in `project-docs/docs/phases/phase-22/tasks/task-2-done.md`; preserve the original failed receipt and add the mandatory deviation note, then run from `/home/mark/projects/fid` the byte-identical required command `go test -race -coverprofile=coverage.out ./...` to exit 0 with all packages passing, every remaining task-2 required command, and fresh independent review-1 re-review. No ESC.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Product Stop/Attention behavior, public interfaces, gate command/result/scope/tolerance, and live-output assertion semantics.
- **Exact verification:** From `/home/mark/projects/fid`: `go test -race -coverprofile=coverage.out ./...`; all remaining task-2 required commands; fresh independent review-1 re-review.
- **Originating report:** `project-docs/docs/phases/phase-22/tasks/task-2-done.md`
- **Control repair base:** N/A
- **If BLOCK:** N/A
- **Human attention required:** no
- **References:** `project-docs/docs/phases/phase-22/triage-log.md:20-50`; `project-docs/docs/phases/phase-22/tasks/task-2-done.md:5-9,73-76,95-105,110-118`; `project-docs/docs/phases/phase-22/tasks/task-2.md:32-34,36-64,94-114`; `project-docs/docs/phases/phase-22/execution-scope.md:8-27`; `project-docs/docs/phases/phase-22/decisions.md:45-54`; DEC-175.

---

## TRG-22-003 - Admissible correction for task 2 retained live-output golden full-gate failure

- **Trigger:** task-2 review-repair exact full-gate failure in the retained cmd/fid live-output golden.
- **Originating task:** task-2
- **Trigger signature:** task-2:latent-gate-failure:cmd-fid-live-output-golden
- **Repair category:** in-scope-implementation
- **Repair authority:** minimum causal scope
- **Source:** `project-docs/docs/phases/phase-22/tasks/task-2-done.md`
- **Candidate:** admissible correction for TRG-22-002; preserve TRG-22-001 and TRG-22-002 as failed/corrective orchestration evidence and use this record as the governing ledger record.
- **Attempt count:** 2
- **Checkpoint:** verified phase HEAD `4dd41c97385f0655140bf56b094528cc4cafe6c1`; phase worktree clean except Runner-owned `project-docs/docs/phases/phase-22/run-ledger.md`; task-2 report preserves the failed required gate.
- **Verdict:** REPAIR
- **Reasoning:** TRG-22-001 and TRG-22-002 remain unchanged as failed/corrective orchestration evidence. The schema-v2 validator rejected TRG-22-002 because its external Trigger signature used Markdown backticks rather than the required raw scalar beginning with `task-2:`. This TRG-22-003 writes the same decision-free same-task authority and exact raw identity, making it the admissible governing record without changing the accepted gate, mismatch, repair scope, or protected semantics.
- **Local decisions made:**
  - Task 2 owns the minimum causal live-output golden repair within the exact bounded files below.
  - TRG-22-003 is the admissible governing record; TRG-22-001 and TRG-22-002 are preserved and must not be edited or used as the governing ledger record.
- **Resolution evidence:** The immutable failed receipt is `go test -race -coverprofile=coverage.out ./...` from `/home/mark/projects/fid`, observed `FAIL` in `cmd/fid/TestLiveOutputGoldenFixtures` (terminated `exit_code: 1`): the full run expected 200 `x` bytes but retained the coverage-child `GOCOVERDIR-not-set` warning and 134 observed bytes. Isolated/package diagnostics do not clear the failed receipt.
- **If REPAIR:** owning task task-2; category in-scope-implementation; next action dispatch a same-task continuation limited exactly to `cmd/fid/live_output_test.go`, `cmd/fid/testdata/val-036/live-frames.jsonl`, and `project-docs/docs/phases/phase-22/tasks/task-2-done.md`; preserve the original failed receipt and add the mandatory deviation note, then from `/home/mark/projects/fid` re-run the byte-identical required command `go test -race -coverprofile=coverage.out ./...`, all remaining task-2 required commands, and fresh independent review-1 re-review. No ESC.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Product Stop/Attention behavior, public interfaces, gate command/result/scope/tolerance, and live-output assertion semantics.
- **Exact verification:** From `/home/mark/projects/fid`: `go test -race -coverprofile=coverage.out ./...`; all remaining task-2 required commands; fresh independent review-1 re-review.
- **Originating report:** `project-docs/docs/phases/phase-22/tasks/task-2-done.md`
- **Control repair base:** N/A
- **If BLOCK:** N/A
- **Human attention required:** no
- **References:** `project-docs/docs/phases/phase-22/triage-log.md:20-82`; `project-docs/docs/phases/phase-22/tasks/task-2-done.md:5-9,73-76,95-105,110-118`; `project-docs/docs/phases/phase-22/tasks/task-2.md:32-34,36-64,94-114`; `project-docs/docs/phases/phase-22/execution-scope.md:8-27`; `project-docs/docs/phases/phase-22/decisions.md:45-54`.

---

## TRG-22-004 - Same-task recovery binding correction for task-2 review-1 repair

- **Trigger:** Phase 22 recovery-sequencing overlap: the active task-2 review-1 recovery remains open after MATERIAL-PROGRESS, so its pending same-task dispatch must retain the review-1 trigger identity while using TRG-22-003 as the substantive minimum-causal live-output scope record.
- **Originating task:** task-2
- **Trigger signature:** review-1:qa-rework:task-2-stop-attention-evidence
- **Repair category:** in-scope-implementation
- **Repair authority:** minimum causal scope
- **Source:** `project-docs/docs/phases/phase-22/run-ledger.md` active task-2 review-1 recovery; substantive scope authority `project-docs/docs/phases/phase-22/triage-log.md` TRG-22-003.
- **Candidate:** typed RECOVERY_REQUIRED candidate returned by the active task-2 review-1 recovery; this direct same-task REPAIR replaces only that pending dispatch identity.
- **Attempt count:** 2
- **Checkpoint:** task-2 review-1 recovery attempt 1 ended MATERIAL-PROGRESS with implementation `0bf39eb2c54aadef2999bed41c654973673cfc57`, report `4dd41c97385f0655140bf56b094528cc4cafe6c1`, and the retained full-gate failure preserved in the report; Runner ledger remains the pending orchestration record.
- **Verdict:** REPAIR
- **Reasoning:** The earlier review-1 findings are already implemented. The active same-task recovery therefore retains its task identity, review-1 trigger signature, and monotonic next attempt; replacing it with TRG-22-003's distinct latent-gate trigger would overlap the pending dispatch. TRG-22-003 remains the substantive live-output minimum-causal scope record and supplies the exact permitted paths and verification without changing any accepted gate, result, scope, or tolerance.
- **Local decisions made:**
  - Bind the next task-2 TASK-RECOVERY to this TRG's raw review-1 identity and attempt 2.
  - Use TRG-22-003 solely as the substantive authority for the bounded retained live-output repair scope.
- **Resolution evidence:** The active review-1 recovery reached MATERIAL-PROGRESS and its report records that both task-2 Must-fix findings are satisfied; the byte-identical required full gate then failed in retained `TestLiveOutputGoldenFixtures` with `FAIL` (terminated `exit_code: 1`). The repair remains decision-free and causal within the exact TRG-22-003 scope.
- **If REPAIR:** owning task task-2; category in-scope-implementation; next action dispatch the active same-task review-1 continuation under this identity, limited exactly to `cmd/fid/live_output_test.go`, `cmd/fid/testdata/val-036/live-frames.jsonl`, and `project-docs/docs/phases/phase-22/tasks/task-2-done.md`; preserve the failed receipt and mandatory deviation evidence, then from `/home/mark/projects/fid` re-run the byte-identical required command `go test -race -coverprofile=coverage.out ./...`, all remaining task-2 required commands, and fresh independent review-1 re-review. No ESC.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Product Stop/Attention behavior, public interfaces, gate command/result/scope/tolerance, live-output assertion semantics, task specs, QA artifacts, decisions, dependencies, workflow controls, and run-ledger entries.
- **Exact verification:** From `/home/mark/projects/fid`: `go test -race -coverprofile=coverage.out ./...`; all remaining task-2 required commands; fresh independent review-1 re-review.
- **Originating report:** `project-docs/docs/phases/phase-22/tasks/task-2-done.md`
- **Control repair base:** N/A
- **If BLOCK:** N/A
- **Human attention required:** no
- **References:** `project-docs/docs/phases/phase-22/run-ledger.md:158-208`; `project-docs/docs/phases/phase-22/triage-log.md:87-117`; `project-docs/docs/phases/phase-22/tasks/task-2-done.md:5-9,95-118`; `project-docs/docs/phases/phase-22/tasks/task-2.md:94-114`; `project-docs/specs/09_AI_WORKFLOW.md:2.3,2.4`.

---

## TRG-22-005 - Tech Lead root-cause correction for the task-2 retained live-output golden failure

- **Trigger:** task-2 review-repair exact full-gate failure in the retained cmd/fid live-output golden.
- **Originating task:** task-2
- **Trigger signature:** review-1:qa-rework:task-2-stop-attention-evidence
- **Repair category:** in-scope-implementation
- **Repair authority:** DEC-189
- **Source:** `project-docs/docs/phases/phase-22/tasks/task-2-done.md`; substantive authority `project-docs/docs/phases/phase-22/decisions.md` DEC-189.
- **Candidate:** Tech Lead labelled correction of TRG-22-003 and TRG-22-004; preserve TRG-22-001..004 unedited as orchestration evidence and use this record as the governing ledger record.
- **Attempt count:** 2
- **Checkpoint:** verified phase HEAD `4dd41c97385f0655140bf56b094528cc4cafe6c1`; phase worktree clean except Runner-owned `project-docs/docs/phases/phase-22/run-ledger.md` and this triage log; task-2 report preserves the failed required gate.
- **Verdict:** REPAIR
- **Reasoning:** TRG-22-001..003 corrected only the formal admissibility of the same disposition and never established a mechanism; TRG-22-004 corrected only the recovery-identity binding. DEC-189 supplies the missing substance: this is a nondeterministic retained test, not subject-preserving convergence and not a DEC-175 production defect under a new gate. Under a coverage-instrumented build the fixture's child inherits the closed executor environment without `GOCOVERDIR` and its exit hook injects a 54-byte warning into observed stderr; `RunLive` stops accumulating `OutputBytes` at the ceiling, so 80 emitted bytes plus that 54-byte warning is exactly the 134 the failed run recorded, against the golden's 200. The repair therefore belongs in the harness, and the goldens must not move.
- **Local decisions made:**
  - Task 2's scope is amended under DEC-189 to add exactly `cmd/fid/live_output_test.go`; its required gates already walk `cmd/fid`, so the assignment is executable.
  - Both goldens are excluded from the allowed files. Re-baselining `live-frames.jsonl` or `human-session.txt` is forbidden; the correct repair leaves them byte-identical.
  - The mechanism must be confirmed by the deterministic diagnostic below before any repair lands. If it is not reproduced, return `RECOVERY_REQUIRED` to Tech Lead instead of adjusting the fixture.
  - This record supersedes the substantive scope of TRG-22-003 and the pending dispatch bound by TRG-22-004 for the same task; the review-1 trigger identity and attempt 2 are retained byte-for-byte.
- **Resolution evidence:** The immutable failed receipt is `go test -race -coverprofile=coverage.out ./...` from `/home/mark/projects/fid`, observed `FAIL` in `cmd/fid/TestLiveOutputGoldenFixtures` (terminated `exit_code: 1`), recording 134 observed bytes and a coverage-child `GOCOVERDIR-not-set` warning against the committed golden's 200. Isolated and package-level diagnostics passed and do not clear that receipt; they are the same race landing the other way, not a package-versus-repository distinction.
- **If REPAIR:** owning task task-2; category in-scope-implementation; next action dispatch the active same-task review-1 continuation under this identity and attempt 2, limited exactly to the allowed files below; confirm the mechanism, apply the minimum causal harness repair, preserve the original failed receipt and add the mandatory deviation note, then run the exact verification below and obtain a fresh independent review-1 re-review. No ESC.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** cmd/fid/live_output_test.go,project-docs/docs/phases/phase-22/tasks/task-2-done.md
- **Forbidden changes:** `cmd/fid/testdata/val-036/live-frames.jsonl` and `cmd/fid/testdata/val-036/human-session.txt` (must remain byte-identical), all product code including `internal/evidence/`, Stop/Attention behavior, public interfaces, schema, gate command/result/scope/threshold/tolerance, live-output assertion semantics, task specs, QA artifacts, decisions, dependencies, workflow controls, and run-ledger entries.
- **Exact verification:** From `/home/mark/projects/fid`, in order: (1) mechanism diagnostic `B=$(mktemp -d)` then `go test -c -cover -o "$B/covbin" ./cmd/fid` then `env -i "$B/covbin" -test.run='^TestLiveGoldenChild$' -- --live-golden-child`, whose stderr must contain the verbatim `warning: GOCOVERDIR not set, no coverage data emitted` line and which terminates non-zero at the child's own `exit_code: 23`; (2) `go test -race -coverprofile=coverage.out ./cmd/fid -run '^TestLiveOutputGoldenFixtures$' -count=20` to exit 0; (3) the byte-identical required command `go test -race -coverprofile=coverage.out ./...` to exit 0 with all packages passing; (4) every remaining task-2 required command from `tasks/task-2.md` Exact Commands; (5) `git status --porcelain -- cmd/fid/testdata/val-036` recording `exit 0 (no output)` as proof the goldens did not move; (6) fresh independent review-1 re-review.
- **Expected non-zero commands:** the step (1) `env -i "$B/covbin" ...` mechanism diagnostic is expected to terminate `exit_code: 23`; record it as `PASS` with the observed status and the verbatim warning line. No other command in this directive has a non-zero pass condition.
- **Originating report:** `project-docs/docs/phases/phase-22/tasks/task-2-done.md`
- **Control repair base:** N/A
- **If BLOCK:** N/A
- **Human attention required:** no
- **References:** `project-docs/docs/phases/phase-22/decisions.md` DEC-189; `project-docs/docs/phases/phase-22/triage-log.md:20-149`; `project-docs/docs/phases/phase-22/run-ledger.md:183-209`; `project-docs/docs/phases/phase-22/tasks/task-2-done.md:5-9,73-76,95-118`; `project-docs/docs/phases/phase-22/tasks/task-2.md:32-34,94-116`; `project-docs/docs/phases/phase-22/qa/review-1.md:97-98,122-123`; `internal/evidence/live.go:96-115,175-213`; `internal/evidence/executor.go:59-80`; `internal/evidence/policy.go:147-174`; `cmd/fid/live_output_test.go:43-107`.

---

## TRG-22-006 - Workflow-sequence validator rejects DEC-189 same-task supersession

- **Trigger:** The Phase Runner cannot append the user-authorized TRIAGE TRG-22-005 and its bound task-2 recovery through the canonical sequence validator.
- **Originating task:** task-2
- **Trigger signature:** task-2:control-conformance:workflow-sequence-lint-supersession-and-review-trigger
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** `project-docs/docs/phases/phase-22/run-ledger.md` and the rejected TRG-22-005 append attempt.
- **Candidate:** protected workflow-control blocker candidate.
- **Attempt count:** 1
- **Checkpoint:** runner-attested verified phase HEAD `3e06f12fe4ae42e2694958043b95741bbb2c5c8a`; persisted ledger and phase artifacts committed; worktree clean; the ledger remains unchanged and standalone validation remains `VALID state=REPAIR_PENDING`.
- **Verdict:** BLOCK
- **Reasoning:** The reproduced canonical append rejects the accepted DEC-189 same-task supersession and the accepted review-form trigger identity with `INVALID TRIAGE REPAIR must be followed by its bound TASK-RECOVERY before TRIAGE; INVALID schema v2 TRIAGE Originating task does not match Trigger signature` (terminated `exit_code: 1`). Restoring the accepted behavior requires changing `project-docs/tools/workflow-sequence-lint.sh`. The accepted workflow expressly reserves any repair of that verifier to separate protected decision authority, so a normal control-conformance directive cannot authorize it.
- **Local decisions made:**
  - Preserve the current valid ledger state; do not bypass, directly edit, or submit an alternative sequence to the canonical validator.
  - Do not authorize a product repair, task-2 recovery, or verifier self-repair under this TRG.
- **Resolution evidence:** The user-attached reproducer ran from `/home/mark/projects/fid`: the exact user-authorized TRG-22-005 TRIAGE candidate piped to `project-docs/tools/workflow-sequence-lint.sh project-docs/docs/phases/phase-22/run-ledger.md --append-stdin` was rejected at `exit_code: 1`, leaving the ledger unchanged. `project-docs/templates/run-ledger-template.md` and DEC-189 require the opposite result: a later schema-v2 TRIAGE REPAIR may replace a pending same-task dispatch and retain the review-form identity.
- **If REPAIR:** N/A; `project-docs/tools/workflow-sequence-lint.sh` is excluded from normal control-conformance repair.
- **Control trigger signature:** N/A
- **Control attempt count:** N/A
- **Control contract:** N/A
- **Control mismatch:** N/A
- **Allowed files:** N/A
- **Forbidden changes:** Any product code, tests, task specs, QA artifacts, product docs, phase decisions, run-ledger entries, accepted command, ledger semantics, trigger identity, task scope, DEC-189 authority, workflow verifier, or command/assertion/status/semantics/scope/threshold/budget/tolerance/policy change without separate protected workflow-control authority.
- **Exact verification:** N/A pending separate protected workflow-control authority.
- **Originating report:** `project-docs/docs/phases/phase-22/tasks/task-2-done.md`
- **Control repair base:** `3e06f12fe4ae42e2694958043b95741bbb2c5c8a`
- **If BLOCK:** category protected workflow-control definition/policy; ESC-22-001 raised in `project-docs/docs/current/tl_escalations_current.md`; resolution condition separate protected workflow-control authority reconciles `project-docs/tools/workflow-sequence-lint.sh` with the accepted same-task supersession and review-form trigger contract before Phase 22 resumes.
- **Human attention required:** yes
- **References:** `project-docs/specs/09_AI_WORKFLOW.md:427-458,586-614,703-713`; `project-docs/templates/run-ledger-template.md`; `project-docs/docs/phases/phase-22/decisions.md:222-270` (DEC-189); `project-docs/docs/phases/phase-22/triage-log.md:153-185` (TRG-22-005); `project-docs/docs/phases/phase-22/run-ledger.md:158-209`; `project-docs/tools/workflow-sequence-lint.sh`.

## Resolution For TRG-22-006 - ESC-22-001 closed under DEC-190

- **Date:** 2026-09-06.
- **Status:** RESOLVED. Original BLOCK verdict and failed append evidence above
  remain unchanged.
- **Authority:** DEC-190, separate Human-authorized protected verifier
  maintenance; DEC-191 covers the discovered canonical/provider config drift.
- **Repair commit:** `3bb42b4efd391b216b752db2c31fe793f3538bac`.
- **Observed resolution:** All 30 recovery boundary checks, the full mechanical
  drill suite and workflow lint passed, each terminating `exit_code: 0`.
  The real Phase 22 normal RUN-RESUME, correcting TRIAGE TRG-22-005, and bound
  task-2 TASK-RECOVERY each returned `APPEND_CONFIRMED state=REPAIR_PENDING`
  (terminated `exit_code: 0`). Standalone validation returned
  `VALID state=REPAIR_PENDING` (terminated `exit_code: 0`). Original ledger
  bytes remain an unchanged prefix, followed by exactly those three events.
- **Reconciliation:** The rejected BLOCK and STOP were never persisted. This
  resolution invents neither event and is not a control-conformance amendment.
  The existing pending ordinary triage was corrected through the canonical
  append path; DEC-189/TRG-22-005 continue to govern product task 2, attempt 2.
- **Next action:** Resume the existing task-2 recovery entry, then task-3 MF-3
  and fresh independent review-1 re-review. Product repair remains pending;
  workflow maintenance does not claim its test gates or QA outcomes.
- **Evidence:** `project-docs/docs/reviews/2026-09-06-phase-22-workflow-recovery.md`.

---

## TRG-22-007 - Task 4 checkpoint projection and retained Stop fixture boundary

- **Trigger:** Task 4 completion report exposes two unresolved obligations outside the task's recorded source boundary: the accepted run.checkpoint behavior cannot admit or project run.checkpointed because internal/journal/project_effects.go rejects it, and the byte-identical required repository-wide race/coverage gate fails in the retained task-2 TestEquivalentStopsShareAttentionButKeepDistinctFacts fixture because its fixed 2026-09-05 Event is behind the journal date 2026-09-06.
- **Originating task:** task-4
- **Trigger signature:** task-4:scope-gap:checkpoint-projection-and-retained-gate
- **Attempt count:** 2
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** `project-docs/docs/phases/phase-22/tasks/task-4-done.md`.
- **Candidate:** protected task-boundary and retained-test ownership blocker candidate.
- **Checkpoint:** Commit-verification evidence supplied by the Runner: implementation `ac5e6f9bec80e06bfe5b7b33e3611a894d16038a` verified as a commit with subject `feat(recovery): add repository reconciliation observation`; completion-report commit `79c4bd10f84ba7fa43c467b10952055dcc1c03d9` is current HEAD with subject `docs(task): record task-4 gate recovery`.
- **Verdict:** BLOCK
- **Reasoning:** The accepted public and schema contracts assign `run.checkpoint` and the `run.checkpointed` effect to task 4, but task 4's recorded source list excludes the journal projection file that currently rejects that Event. Extending the task into that canonical journal-effect implementation changes an accepted task boundary and can affect event replay/projection behavior, so it requires Tech Lead authority. Separately, the failed Stop fixture has supplied Task 2 introducing provenance and its task-4 report states no accepted task-4 change displaced it. The mandatory gate executes it, but gate reach does not establish task-4 ownership. No current authority assigns its repair to task 4, and this triage cannot invent an owner or authorize a task-2 test edit.
- **Local decisions made:**
  - Preserve both immutable conditions: the rejected `run.checkpointed` projection path and the failed full-gate receipt.
  - Do not dispatch task-4 recovery, alter the required gate, change journal event/effect semantics, or edit the retained Task 2 fixture under task-4 authority.
- **Resolution evidence:** `project_effects.go` rejects `run.checkpointed` before its required effects; the accepted schema assigns that Event's checkpoint observation and runtime-version effects to task 4. The required `go test -race -coverprofile=coverage.out ./...` attempt failed in `TestEquivalentStopsShareAttentionButKeepDistinctFacts` with `system clock is behind latest journal date 2026-09-06; refusing to append out of order` because the fixture uses 2026-09-05. The report supplies introducing commit `ac827959f8359473cbd6546ac6bea1ff1fd20751` for task 2 and supplies no invalidating task-4 change.
- **If REPAIR:** N/A pending a full Tech Lead disposition that either grants a bounded task-4 journal-projection scope amendment conforming to the accepted `run.checkpointed` contract and separately establishes the retained fixture's causal owner, or records another accepted authority.
- **Control trigger signature:** N/A; no control-conformance repair is proposed.
- **Control attempt count:** N/A; no control-conformance repair is proposed.
- **Control contract:** N/A; no control-conformance repair is proposed.
- **Control mismatch:** N/A; no control-conformance repair is proposed.
- **Allowed files:** N/A; no control-conformance repair is proposed.
- **Forbidden changes:** Do not change the required gate command, accepted pass/fail semantics, scope, threshold, tolerance, product contract, task boundaries, or workflow-control definitions. Do not authorize editing task-2-owned retained tests as task-4 work without provenance and authority.
- **Exact verification:** N/A; no control-conformance repair is proposed.
- **Originating report:** project-docs/docs/phases/phase-22/tasks/task-4-done.md
- **Control repair base:** 79c4bd10f84ba7fa43c467b10952055dcc1c03d9
- **If BLOCK:** category protected task boundary and ambiguous retained-test ownership; ESC-22-002 raised in `project-docs/docs/current/tl_escalations_current.md`. Full Tech Lead must establish an authorized bounded repair route before task 4 can resume.
- **Human attention required:** yes
- **References:** `project-docs/docs/phases/phase-22/tasks/task-4.md:20-44,55-67`; `project-docs/docs/phases/phase-22/tasks/task-4-done.md:5-10,34-37,60-71`; `project-docs/docs/phases/phase-22/schema.md:268-274`; `project-docs/docs/phases/phase-22/public-interface.md:145-151`; `project-docs/docs/phases/phase-22/execution-scope.md:70-75`; `project-docs/docs/phases/phase-22/tasks/task-2.md:77-84`; `project-docs/docs/current/decisions_active.md` DEC-162/165/166/167/173; `project-docs/specs/09_AI_WORKFLOW.md` §2.3.

---

## Resolution For TRG-22-007 - ESC-22-002 closed under DEC-192/193

- **Date:** 2026-09-06.
- **Status:** RESOLVED. The original BLOCK verdict, its reasoning, and the
  preserved failed receipts above remain unchanged and correct on their own
  terms: TL-Lite was right to refuse to extend a task boundary or to infer
  retained-test ownership from gate reach.
- **Authority:** DEC-192 (progressive recovery-Event fence) and DEC-193
  (calendar-dependent retained fixtures are DEC-189 nondeterminism), in
  `project-docs/docs/phases/phase-22/decisions.md`.
- **Obligation 1 - `run.checkpointed` projection.** Not a protected boundary
  change. `schema.md` section 5 already assigns `run.checkpointed` to task 4,
  and `unavailableRecoveryEvents` is task 1's progressive "not yet available"
  fence, whose per-key removal is how the owning task implements its Event.
  Tasks 5, 7, 8 and 10 already carry `internal/journal/` in scope; task 4's
  file-level list omitted it. That is a Tech Lead planning omission under
  DEC-167/173, now corrected. Task 4 is amended to delete exactly the
  `"run.checkpointed"` key and to add `internal/journal/checkpoint_projection_test.go`.
- **Obligation 2 - retained fixed-date Stop fixture.** Provenance confirmed
  from the introducing commit `ac827959f8359473cbd6546ac6bea1ff1fd20751`
  (`feat(workflow): add typed Stop and Attention facts`, phase-22/task-2), not
  from a filename or label. Mechanism confirmed: `stopRaisedEvent` pins
  `time.Date(2026, 9, 5, ...)` while `initializedWorkflowWorkspace` writes a
  live-clock initialization Event, so the writer's out-of-order guard correctly
  refuses the backdated append from 2026-09-06 onward. Task 4's commit
  `ac5e6f9bec80e06bfe5b7b33e3611a894d16038a` touches only
  `internal/project/repository{,_test}.go`,
  `internal/workflow/reconciliation{,_test}.go` and one `activation.go`
  capability line, none of which reaches the writer, the guard, or the
  initializer. So this is not task-4 convergence, as TRG-22-007 held. It is
  instead the already accepted DEC-189 class: a nondeterministic retained test
  failing under an unchanged required gate, assigned to the observing task by
  scope amendment with the nondeterminism repaired at its source. Time
  dependence is nondeterminism; the same bytes and command pass or fail by the
  wall clock. Task 2 retains Stop/Attention subject ownership and its review-1b
  PASS is untouched.
- **Executability check:** Both amendments are executed by task 4's unchanged
  `go test -race -coverprofile=coverage.out ./...`, which walks
  `internal/journal` and `internal/workflow`. DEC-193 adds one repeat command
  proving the repaired fixture is stable rather than newly date-lucky.
- **Preserved:** Every required gate command, threshold, tolerance, and pass/fail
  semantic; the `2/1/2` assertion subject; the `run.checkpointed` Event and
  Effect contract; all VAL rows and validators; the task allocation; and both
  immutable failed receipts in the task-4 report.
- **Systemic prevention:** DEC-192 pre-authorizes the identical fence removal
  for tasks 5, 7, 8 and 10 and schedules the VAL-003 convergence for tasks 5
  and 10, so this cannot re-escalate four more times. DEC-193 bans absolute
  calendar dates on writer-appended Events; `internal/workflow/stop_test.go` is
  the only current instance.
- **Next action:** Runner resumes with a normal RUN-RESUME, appends an ordinary
  TRIAGE REPAIR carrying this authority, then dispatches the bound task-4
  same-task continuation. No Human decision is required.
- **Human attention required:** no. Nothing accepted was changed; a planning
  omission was corrected and an existing decision was applied to a second
  instance of its own class.

---

## TRG-22-008 - Task 6 recovery CLI lacks canonical checkpoint and inspection domain APIs

- **Trigger:** The task-6 real built-binary journey reached fid run inspect and terminated with exit_code: 1 and final token FAIL because the current CLI has no canonical checkpoint or inspection behavior; the Developer reports internal/workflow contains no checkpoint mutation/Event producer or inspection API, but task-6's accepted source scope is CLI and integration only.
- **Originating task:** task-6
- **Trigger signature:** task-6:scope-gap:recovery-cli-domain-producers
- **Attempt count:** 2
- **Repair category:** N/A
- **Repair authority:** N/A
- **Source:** `project-docs/docs/phases/phase-22/tasks/task-6-done.md`.
- **Candidate:** protected task-boundary and missing canonical domain behavior blocker candidate.
- **Checkpoint:** Runner verified implementation commit `c69456fc9bf42f811e4999f1968a48c4a82d3cc5` with subject `test(recovery): checkpoint task-6 CLI gap`; report commit `7f91325ee2c96d50205dca9564806080dc4e3059` is current HEAD with the finalized task-6 `RECOVERY_REQUIRED` report. Independent report-lint returned `OUTCOME RECOVERY_REQUIRED` (terminated `exit_code: 0`).
- **Verdict:** BLOCK
- **Reasoning:** The built-binary journey demonstrates that `fid run inspect` is absent, and task 6 cannot truthfully expose the required recovery CLI leaves while `internal/workflow` has neither a canonical checkpoint mutation/Event producer nor an inspection read API. Implementing those domain operations is outside task 6's accepted CLI-and-integration scope and changes a protected task boundary. The supplied material names no accepted decision or bounded authority that assigns the missing domain behavior to task 6 or another task, and no control-conformance mismatch is proposed.
- **Local decisions made:**
  - Preserve the immutable built-binary failure and the finalized `RECOVERY_REQUIRED` report.
  - Do not fabricate checkpoint or inspection behavior in the CLI, widen task 6, alter the public recovery leaf contract, or modify domain Event/API behavior without Tech Lead authority.
- **Resolution evidence:** `TestP22EarlyRecoveryJourney` reaches the real `fid run inspect` invocation and the built binary rejects the absent leaf, causing the focused integration test to terminate `exit_code: 1` with final token `FAIL`. The task report establishes that the canonical workflow producer and inspection API required to make the CLI conform are absent.
- **If REPAIR:** N/A pending a full Tech Lead disposition that identifies the accepted owner and grants a bounded domain-and-task-scope repair route or records an accepted authority; no ESC-free same-task repair is authorized.
- **Control trigger signature:** N/A; no control-conformance repair is proposed.
- **Control attempt count:** N/A; no control-conformance repair is proposed.
- **Control contract:** N/A; no control-conformance repair is proposed.
- **Control mismatch:** N/A; no control-conformance repair is proposed.
- **Allowed files:** N/A; no control-conformance repair is proposed.
- **Forbidden changes:** Do not change the required gate commands, task allocation, accepted public interface, schema, VAL rows, validator, workflow-control definitions, or protected semantics. Do not invent authority for `internal/workflow` domain changes under task-6's CLI-only scope.
- **Exact verification:** N/A; no control-conformance repair is proposed.
- **Originating report:** project-docs/docs/phases/phase-22/tasks/task-6-done.md
- **Control repair base:** 7f91325ee2c96d50205dca9564806080dc4e3059
- **If BLOCK:** category protected task boundary and ambiguous owner for required canonical workflow behavior; ESC-22-003 raised in `project-docs/docs/current/tl_escalations_current.md`; full Tech Lead must establish an accepted owner and bounded authority for the missing checkpoint mutation/Event producer and inspection API before task 6 can resume.
- **Human attention required:** yes
- **References:** `project-docs/docs/phases/phase-22/tasks/task-6.md:5-45,54-80`; `project-docs/docs/phases/phase-22/tasks/task-6-done.md:5-10,21-31,50-64,82-84`; `project-docs/docs/current/phase_brief.md:41`; `project-docs/docs/current/decisions_active.md:109-121`; `project-docs/docs/phases/phase-22/decisions.md:397-455`; `project-docs/specs/09_AI_WORKFLOW.md` §2.3.
