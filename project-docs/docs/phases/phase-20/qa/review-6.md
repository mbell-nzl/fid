# Phase 20 - Mid-Phase Review 6 (Tasks 16-17)

Reviewer: qa-reviewer
Date: 2026-08-16
Tasks covered: phase-20/task-16 and phase-20/task-17
Review type: INITIAL
Prior review: none (review-5d.md inspected as relevant Task 15 format-guard context)

## Verdict

PASS

No Must-fix or Blocking finding remains in Tasks 16-17. The one documentation
Improvement is explicitly owned by pending Task 18 and does not change this
review-window result.

## Scope Reviewed

- Task specifications and completion reports:
  - `project-docs/docs/phases/phase-20/tasks/task-16.md`
  - `project-docs/docs/phases/phase-20/tasks/task-16-done.md`
  - `project-docs/docs/phases/phase-20/tasks/task-17.md`
  - `project-docs/docs/phases/phase-20/tasks/task-17-done.md`
- QA and accepted-contract artifacts:
  - `project-docs/docs/phases/phase-20/qa/validation-contract.md`
  - `project-docs/docs/phases/phase-20/qa/qa-plan.md` (EVID-13, EVID-15,
    Review 6)
  - `project-docs/docs/phases/phase-20/contracts.md` §§2-7 and §9
  - `project-docs/docs/phases/phase-20/decisions.md` (DEC-143 through
    DEC-146)
  - `project-docs/docs/phases/phase-20/qa/review-5d.md`
- Current context, canonical-spec amendments, and documentation:
  - `AGENTS.md`; `project-docs/docs/current/phase_brief.md`,
    `decisions_active.md`, `progress_current.md`, and
    `tl_escalations_current.md`
  - `project-docs/specs/02_TRD.md`, `03_TRD_SCHEMA.md`,
    `04_TRD_CONFIG.md`, `06_DESIGN_PLAN_CLI.md`, and
    `08_QUALITY_PLAN.md`
  - `project-docs/docs/playbook.md`, `CHANGELOG.md`, and pending
    `project-docs/docs/phases/phase-20/tasks/task-18.md`
- Live source, command tree, and recovery evidence:
  - `cmd/fid/root.go`, `cmd/fid/root_test.go`, and
    `cmd/fid/contract_surface_test.go`
  - Task 16a/16b recovery reports, `run-ledger.md`, `triage-log.md`, and the
    committed recovery history through `603fdac9cfd1f09d255983407c71ae4186ee4c9b`

## Summary

- Review outcome: PASS
- Blocking findings: 0
- Must-fix findings: 0
- Improvement findings: 1
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-030, VAL-031, VAL-032, VAL-033, VAL-034,
  VAL-037 (documentation progress)

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-030 | task-16 | ON TRACK | `ls internal/format` produced the absent-directory diagnostic and observed `exit_code: 2`, matching Task 16's declared expected non-zero status; the exact import search produced no output with observed `exit_code: 1`. `TestProductSourceContainsNoUniversalSerializers` also passed in this review. |
| VAL-031 | tasks 14-16 | ON TRACK | The focused live-tree guard passed. Its source walks local, persistent, and inherited flags (`cmd/fid/contract_surface_test.go:126-142`), and its injected inherited and local negative controls assert `fid issue --format` and `fid ready --yaml` respectively (`:87-116`). |
| VAL-032 | task-16 | ON TRACK | `go mod tidy`, repository race/coverage test, and build each observed `exit_code: 0`; `go.mod` YAML search observed `exit_code: 1`; `go mod why -vendor -m gopkg.in/yaml.v3` observed exactly one DEC-145 stanza: `# gopkg.in/yaml.v3` followed by `(main module does not need to vendor module gopkg.in/yaml.v3)`, then `exit_code: 0`. |
| VAL-033 | task-16 | ON TRACK | The exact universal-serializer source search produced no output with observed `exit_code: 1`; the focused product-source guard passed. |
| VAL-034 | task-17 | ON TRACK | The source absence search produced no output with observed `exit_code: 1`; `TestPromisedGovernanceCommandsAreAbsent` passed and asserts unknown command plus empty stdout/stderr for all three former roots (`cmd/fid/root_test.go:9-35`). The root registers neither family (`cmd/fid/root.go:58-75`). |
| VAL-037 | tasks 16 and 18 | DEFERRED TO TASK 18 | The architecture portion is current: `AGENTS.md:32` lists `internal/cli` and `internal/project` and no `internal/format`. The full EVID-15 playbook/changelog synchronization remains Task 18's accepted scope (`tasks/task-18.md:27-31,47`). Its absence is recorded below as an Improvement, not upgraded to a pass. |

## Gate And Evidence Checks

All commands below were run by this Reviewer from `/home/mark/projects/fid` in
this session. Tokens are verbatim command output tokens; every termination is
the observed nested command `exit_code`.

| Check | Command or evidence | Result | Notes |
|---|---|---|---|
| Review 6 required reproduction | `go test -race -coverprofile=coverage.out ./...` | PASS | `ok  github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 80.6% of statements` (terminated `exit_code: 0`). |
| Task 16/17 formatting | `gofmt -l cmd internal` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`). |
| Task 16/17 lint | `golangci-lint run ./...` | PASS | `0 issues.` (terminated `exit_code: 0`). |
| Coverage report | `go tool cover -func=coverage.out` | PASS | `total: (statements) 78.7%` (terminated `exit_code: 0`). |
| Task 16/17 integration | `go test -tags integration -race ./cmd/fid/integration/...` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid/integration 4.932s` (terminated `exit_code: 0`). |
| Task 16/17 build | `go build -o /tmp/fid ./cmd/fid` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`). |
| Task 16 module tidy | `go mod tidy` | PASS | `exit 0 (no output)` (terminated `exit_code: 0`). No `go.mod` or `go.sum` working-tree diff resulted. |
| VAL-030 directory absence | `ls internal/format` | PASS | `ls: cannot access 'internal/format': No such file or directory` (expected `exit_code: 2`; terminated `exit_code: 2`). |
| VAL-030 import absence | `rg -n "internal/format" cmd internal` | PASS | `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |
| VAL-033 serializer absence | `rg -ni "writeTOON|writeYAML|writeMarkdown|yaml.Marshal" cmd internal` | PASS | `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |
| VAL-032 direct module absence | `rg -ni "yaml" go.mod` | PASS | `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |
| VAL-032 build-graph absence | `go mod why -vendor -m gopkg.in/yaml.v3` | PASS | `# gopkg.in/yaml.v3` followed by `(main module does not need to vendor module gopkg.in/yaml.v3)`; exactly one two-line stanza and no further non-empty output (terminated `exit_code: 0`). |
| VAL-034 stub/promise absence | `rg -n "newDecisionCommand|newEscalationCommand|newTriageCommand|versionGateError|ships in v0\.2" cmd/fid` | PASS | `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`). |
| Observable removal controls | `go test ./cmd/fid -run '^(TestProductSourceContainsNoUniversalSerializers|TestLiveCobraTreeContainsNoUniversalFormatFlags|TestLiveCobraTreeUniversalFormatGuardDetectsInjectedFlags|TestPromisedGovernanceCommandsAreAbsent)$' -count=1` | PASS | `ok  github.com/mbell-nzl/fid/cmd/fid 0.069s` (terminated `exit_code: 0`). This covers clean and injected non-qualifying flag states as well as all three unknown governance roots. |
| Task 16 completion-report mechanics | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-16-done.md project-docs/docs/phases/phase-20/tasks/task-16.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). It found all twelve exact declared commands, Task 16's resolving hash/subject, and a committed clean report. |
| Task 17 completion-report mechanics | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-20/tasks/task-17-done.md project-docs/docs/phases/phase-20/tasks/task-17.md` | PASS | `PASS report-lint completed with no failures` (terminated `exit_code: 0`). It found all seven exact declared commands, Task 17's resolving hash/subject, and a committed clean report. |

The Task 16 and Task 17 reports each record every command declared by their
task spec, including all QA-plan task-specific absence checks; the QA plan
states there are no further commands for either report. Their exact-gate rows
match the current declarations. Task 16's expected `exit_code: 2` for `ls` and
its three declared no-match searches remain correctly distinguished from a
zero-exit PASS.

## Per-Task Findings

### Task 16 - Delete universal formatters and YAML

#### Blocking

- None.

#### Must-fix

- None.

#### Improvement

- [CHANGELOG.md:8-11,38-39,84-94] The changelog has no Phase 20 removal entry
  yet and its historical `0.1.0` section still describes the retired formatter
  and future governance-stub promise. The claims are labelled historical, so
  this is not a product or removal contradiction; however, VAL-037/EVID-15 is
  not complete until the live changelog clearly records the delivered Phase 20
  removals. Fix direction: Task 18 should add the required Phase 20 changelog
  entry and preserve historical context without presenting it as the current
  surface.

### Task 17 - Delete promised-command stubs

#### Blocking

- None.

#### Must-fix

- None.

#### Improvement

- None.

## Consistency, Removal, And Recovery Assessment

- **Accepted CLI/config/schema/Event/removal consistency: PASS for this
  window.** The live root supplies only `--project` and `--json`
  (`cmd/fid/root.go:52-75`), consistent with `contracts.md` §§2.1-2.3;
  `go.mod` has no YAML requirement; and the source tree has no
  `internal/format` directory. The retained governance tables remain the
  correct non-removal boundary under `contracts.md:661-663`; this Task 17
  removal only deletes their unsupported adapters. The repository-wide race
  suite and real-process integration suite passed after the removals, so the
  retained Task 1-15 functionality still compiles and its exercised behavior
  remains intact. `contracts.md` §§4 and 7 retain the accepted DDL and
  `event_schema: 1`/revision-event contract, and neither Task 16 nor 17 changes
  migrations or journal event code.
- **Task 15 guard carry-forward: PASS.** Review 5d's repaired inherited/local
  flag controls remain live at `cmd/fid/contract_surface_test.go:87-145` and
  passed in this fresh focused run. Task 16's deletion evidence therefore does
  not merely show positive absence: it also proves the flag walker fails when
  either an inherited persistent or local universal flag is injected.
- **Task 16 recoveries: PASS as committed evidence.** `git rev-parse --verify`
  resolved `afb262d`, `d9d6add`, `d8388d5`, `20813fc`, and `d0f9cea` to
  commits with their recorded subjects/footers. The recovery reports and
  `task-16-done.md:56-92` preserve the withdrawn `go.sum` command, both
  pre-DEC-145 vendor-output failed attempts, all report-lint failures, and their
  observed terminations. DEC-143/144/145 altered only the accepted evidence
  notation/literal protocol under protected authority; the current required
  command, YAML dependency boundary, assertion, and failed evidence are not
  silently replaced. The fresh exact stanza result matches the final DEC-145
  contract.
- **Task 17 control-conformance recovery: PASS as committed evidence.**
  `git rev-parse --verify` resolved `95ee13d`, `88b0377`, `0c1094c`, and
  `92cb0f4`. `git diff --name-status 9edb337..88b0377` reports only
  `project-docs/tools/report-lint.sh`; its diff replaces escape-interpreting
  `awk -v` equality with opaque shell comparisons and adds red/green `\.` and
  `\t` fixtures. That is the recorded DEC-146 scope, not a change to Task 17's
  command, expected exit, evidence token, VAL-034, or product behavior. The
  recorded scope receipt in `run-ledger.md:1414-1423` binds
  `9edb337 -> 88b0377`, and `task-17-done.md:58-59` retains the original failed
  lint attempt before documenting the unchanged successful invocation. Fresh
  report-lint for the final report passed.
- **Documentation progress: PARTIAL, deferred.** AGENTS and the playbook match
  the removal boundary (`AGENTS.md:32`; `project-docs/docs/playbook.md:266-267`).
  The changelog work is intentionally in Task 18's file scope and acceptance
  criterion (`tasks/task-18.md:27-31,47`), so the full VAL-037/EVID-15 check is
  not passed early and does not reopen either reviewed task.

## Prior Finding Resolution

Not applicable: this is the initial Review 6, not a re-review of a Review 6
finding. Relevant `review-5d.md` was inspected because its Task 15
inherited-format guard is a direct prerequisite of Task 16; its prior Must-fix
is confirmed resolved by the fresh focused guard result above.

## Repair Plan

Not required: there are no Must-fix findings.

## Blocker Candidates

- None. The current escalation register reports no open ESC
  (`project-docs/docs/current/tl_escalations_current.md:10-15`). The historical
  Task 16 specification/control contradictions and Task 17 report-lint defect
  are closed under DEC-143 through DEC-146, with their immutable evidence
  preserved; none is reopened by the fresh command evidence.

## Existing Escalation References

- None. `ESC-20-012` through `ESC-20-015` are historical resolved references,
  not current escalation coverage.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-20/qa/review-6.md
Must-fix: 0  Improvement: 1  Blocker candidates: 0
VERDICT: CONTINUE
```
