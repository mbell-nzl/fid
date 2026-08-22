# Phase 20 - Mid-Phase Review 1c (DEC-131/DEC-132 recovery: tasks 3c-3d)

Reviewer: qa-reviewer
Date: 2026-08-10
Tasks covered: phase-20/task-3c, phase-20/task-3d

## Verdict

CONTINUE

## Scope Reviewed

- Task specs: `tasks/task-3c.md`, `tasks/task-3d.md`.
- Completion reports: `tasks/task-3c-done.md`, `tasks/task-3d-done.md`.
- Validation contract: `qa/validation-contract.md` (VAL-043 and VAL-044).
- QA plan: EVID-14, EVID-16, and deterministic-gate rows 143-150.
- Prior recovery records: DEC-131, DEC-132, `qa/review-1.md`, and
  `qa/review-1b.md`.
- Commits: `f978b6d382b3ae90ff44a7f4577799bb048329d4` and
  `d0bffbcd63d9a36474fde7a3da2c76facb604501`.

Both reported hashes resolve with `git rev-parse --verify <hash>^{commit}`.
Their subjects and `Refs: phase-20/task-3c` / `Refs: phase-20/task-3d`
footers match. `git diff --check` is clean for each recovery commit. This review
does not inspect or plan DEC-132's separate infrastructure performance pack.

## Summary

- Blocking findings: 0
- Improvement findings: 0
- Escalations raised: none
- Covered validation assertions: VAL-043, VAL-044

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-043 | task-3c | ON TRACK | `internal/config/config_test.go:158-162` supplies independent rejected `max_excerpt_bytes` zero, negative, overflow, float, and quoted rows alongside retained accepted values. The current-session focused test passed: `ok  github.com/mbell-nzl/fid/internal/config  1.035s`. This closes the Review 1b missing non-qualifying evidence. |
| VAL-044 | task-3c | ON TRACK — closes at Task 18 | The three binary integration initializations in `cmd/fid/e2e_v01_closure_test.go:19` and `cmd/fid/e2e_v02a_workflow_authoring_test.go:17,79` explicitly pass `--prefix FID`, preserving their assertions while respecting strict derived-prefix validation. Current-session focused integration evidence passed: `ok  github.com/mbell-nzl/fid/cmd/fid  9.755s`. The full 21-leaf request matrix remains Task 18 work. |
| no VAL | task-3d | ON TRACK | The POSIX workflow control is governed by DEC-131/DEC-132 and AGENTS.md rather than a product assertion. Its self-test covers accepted real `go test` output plus rejected silent and mismatched tokens; its workflow guardrail is independently green. |

## Gate And Evidence Checks

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Task 3c exact command coverage | `project-docs/tools/report-lint.sh tasks/task-3c-done.md tasks/task-3c.md` | PASS | All six exact task/QA-plan gates are present with PASS records; independent lint ended `PASS report-lint completed with no failures`. |
| Task 3d exact command coverage | `project-docs/tools/report-lint.sh tasks/task-3d-done.md tasks/task-3d.md` | PASS | All seven exact gates, including `report-lint.sh --self-test`, are present with plausible `go test` output tokens; independent lint ended `PASS report-lint completed with no failures`. |
| Task 3d control regression | `project-docs/tools/report-lint.sh --self-test` | PASS | The independent self-test passed the real-token acceptance, impossible-token rejection, mismatched-token rejection, corrected-history, and wrapped-blocker cases; terminal line: `SELF-TEST PASS report-lint self-test completed with no failures`. |
| Task 3d workflow guardrails | `project-docs/tools/workflow-lint.sh` | PASS | Canonical and generated Developer/QA Lead/QA Reviewer prompts require full-result preservation, `session_id` polling, and observed `exit_code`; Phase Runner has the wrapper-completion guardrail. Terminal line: `PASS workflow-lint completed with no failures`. |
| Focused VAL-043/044 reproduction | `go test -race ./internal/config ./cmd/fid -run 'TestParseProjectConfigBoundaryDomains|TestV01ClosureJourneyAndRebuildParity|TestV02aWorkflowAuthoringExitScenarioAndRebuildParity|TestMaintenanceCheckReportsIncompatibleWorkflowAssignment' -count=1 -tags integration` | PASS | Both packages passed in the current session: `internal/config 1.035s`; `cmd/fid 9.755s`. It exercises each newly rejected config input and all three repaired real-binary journeys. |
| Required independent integration reproduction | `go test -tags integration -race ./...` | PASS | Run exactly from `/home/mark/projects/fid`; complete nested result included `exit_code: 0`. Real output included `ok  github.com/mbell-nzl/fid/cmd/fid (cached)` and ended `ok  github.com/mbell-nzl/fid/internal/workflow (cached)`. |
| DEC-131 correction | `qa/review-1.md` `## Corrections` | PASS | This review appends the authorized correction, retracting Review 1's impossible independent-integration row and replacing it with the real current-session output tail. |

## Per-Task Findings

### Task 3c - Integration-journey convergence and honest gate baseline

#### Blocking

- None.

#### Improvement

- None.

The recovery is contract-aligned: it fixes stale fixture inputs instead of
relaxing the derived-prefix rule, retains all three journey assertions, and
adds the five missing rejected `max_excerpt_bytes` states. The report's cited
`cd53c43319a3301d3a8cc0089b6785d9afbb163f` provenance commit resolves and
predates tasks 3a/3b as stated.

### Task 3d - Evidence-token control and gate-record correction

#### Blocking

- None.

#### Improvement

- None.

The report-lint control detects the specific impossible-token failure mode
without treating wrapped BLOCKED/UNVERIFIED explanatory text as a command
record. Its independent self-test and the generated-prompt/workflow-lint
evidence substantiate the DEC-132 long-running-command guardrails. The
authorized scope extension is recorded in DEC-132; no product behavior or
accepted validation criterion changed.

## Cross-cutting observations

- All five historical completion-report corrections preserve the original
  history while marking their impossible full-suite `go test` records
  UNVERIFIED and cite Task 3c's current-tree baseline, as DEC-131 requires.
- The exact integration reproduction above terminated from an observed
  `exit_code`; no inference was made from silence or outer-wrapper completion.

## Escalations

- None. ESC-20-002 is resolved by DEC-131 and ESC-20-003 by DEC-132. This
  bounded review found no new cross-task, acceptance, architecture, security,
  or quality-gate anomaly requiring escalation.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-20/qa/review-1c.md
Blocking: 0  Improvement: 0  Escalations: none
VERDICT: CONTINUE (NON-BLOCK)
```
