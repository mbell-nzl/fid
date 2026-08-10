# Phase 20 - Task 3d: Evidence-token control and gate-record correction

## Summary of what was implemented

- Added a POSIX `sh`, dependency-free `report-lint.sh` control that rejects an
  impossible silent-success token and tokens that do not look like `go test`
  output, while accepting real package-output markers.
- Narrowed token validation to recorded command-result bullets so wrapped
  blocker explanations containing `go test` are not false matches; legitimate
  `BLOCKED` and `UNVERIFIED` records remain valid non-PASS states.
- Added DEC-131 correction notes to Tasks 1, 2, 3, 3a, and 3b. They mark the
  affected full-suite records UNVERIFIED, preserve their original history, and
  cite Task 3c's current-tree gate as the verified Phase 20 baseline.
- Repaired Codex long-running gate handling in the canonical Developer, QA
  Lead, QA Reviewer, and Phase Runner prompts, regenerated provider prompts,
  and added deterministic workflow-lint coverage.
- Recorded the Human-authorized ESC-20-003 recovery and its explicit exclusion
  of broader test-performance work in DEC-132.

## Files created or modified

- `AGENTS.md`
- `project-docs/tools/report-lint.sh`
- `project-docs/tools/workflow-lint.sh`
- `project-docs/agents/developer.md`
- `project-docs/agents/qa-lead.md`
- `project-docs/agents/qa-reviewer.md`
- `project-docs/agents/phase-runner.md`
- `.codex/agents/developer.toml`
- `.codex/agents/qa_lead.toml`
- `.codex/agents/qa_reviewer.toml`
- `.codex/agents/phase_runner.toml`
- `.claude/agents/phase-runner.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/templates/phase-runner-prompt.md`
- `project-docs/docs/phases/phase-20/decisions.md`
- `project-docs/docs/phases/phase-20/tasks/task-1-done.md`
- `project-docs/docs/phases/phase-20/tasks/task-2-done.md`
- `project-docs/docs/phases/phase-20/tasks/task-3-done.md`
- `project-docs/docs/phases/phase-20/tasks/task-3a-done.md`
- `project-docs/docs/phases/phase-20/tasks/task-3b-done.md`
- `project-docs/docs/phases/phase-20/tasks/task-3c.md`
- `project-docs/docs/phases/phase-20/tasks/task-3d.md`
- `project-docs/docs/phases/phase-20/tasks/task-3d-done.md`
- `project-docs/docs/phases/phase-20/triage-log.md`
- `project-docs/docs/phases/phase-20/run-ledger.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/reviews/2026-08-10-gate-execution-and-test-runtime-diagnosis.md`

## Documentation updated

- Corrected the five historical completion reports required by DEC-131.
- Clarified the general command-termination rule in `AGENTS.md` and
  `09_AI_WORKFLOW.md`, and repaired the canonical/generated role prompts.
- Captured the diagnosis for posterity in the requested review document.
- Preserved the Phase Runner's append-only STOP, triage, and escalation state;
  closure is recorded after the finalized report passes report-lint.

## Validation Contract Coverage

- No VAL-NNN applies. This is V1 delivery-process tooling and historical
  evidence correction governed by AGENTS.md (Commands), Phase 19 DEC-102,
  Phase 20 DEC-131, and the Human-authorized DEC-132 recovery, not product
  behavior.

## Tests run

- Pre-existing Task 3d red: `project-docs/tools/report-lint.sh --self-test` -
  EXPECTED RED before the original token control -
  `SELF-TEST FAIL report-lint self-test found 2 failure(s)`.
- Recovery red: `project-docs/tools/report-lint.sh --self-test` - EXPECTED RED
  after adding the wrapped-blocker regression -
  `SELF-TEST FAIL report-lint self-test found 1 failure(s)`.
- Recovery red: `project-docs/tools/workflow-lint.sh` - EXPECTED RED after
  adding canonical/generated prompt checks -
  `FAIL workflow-lint found 29 failure(s)`.
- `project-docs/tools/workflow-lint.sh` - PASS -
  `PASS workflow-lint completed with no failures`.
- The integration gate remained silent across the initial 30-second yield and
  three retained-session polls, then returned `exit_code: 0`; its package
  output included `ok  github.com/mbell-nzl/fid/cmd/fid  160.765s`.

## Exact quality-gate results

- `project-docs/tools/report-lint.sh --self-test` - PASS -
  `SELF-TEST PASS report-lint self-test completed with no failures`.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`.
- `golangci-lint run ./...` - PASS - `0 issues.`.
- `go test -race -coverprofile=coverage.out ./...` - PASS -
  `ok  github.com/mbell-nzl/fid/internal/workflow  (cached)  coverage: 80.5% of statements`.
- `go tool cover -func=coverage.out` - PASS -
  `total:  (statements)  77.7%`.
- `go test -tags integration -race ./...` - PASS -
  `ok  github.com/mbell-nzl/fid/internal/workflow  (cached)`.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`.
- Conditional correction-report lints for Tasks 1, 2, 3, 3a, and 3b are run
  after this task commit makes the corrected reports committed and clean; the
  finalized report records their results.

## Deviations from spec

- The Human explicitly authorized a recovery super-agent to cross the normal
  one-task/session and role boundaries for ESC-20-003. The recovery also fixes
  canonical long-running gate handling and records the diagnosis; it makes no
  product change and does not begin the broader performance follow-up pack.
- Conditional lints for the five corrected historical reports must run after
  the task commit because `report-lint.sh` intentionally rejects reports with
  uncommitted modifications. Their results are added to this report in the
  immediately following final-report commit.

## Concerns, technical debt, or notes for the Tech Lead

- The integration gate still spends 160.765 seconds in `cmd/fid`. DEC-132 keeps
  performance refactors out of this recovery; the Human requested a separate
  infrastructure follow-up pack before Task 4.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
