# Phase 9 Final Review - Concerns + v0.1c Closure

Reviewer: QA Lead
Date: 2026-06-20

## Verdict

BLOCK

Deterministic phase-end review found two blocking issues:

1. `go test -bench . -benchmem ./...` fails in `cmd/fid`
   `TestConcernListFiltersByIssueAndStatus`.
2. VAL-004 / QA-12 fails public read-back: after `fid concern address`, both
   `fid concern list --format json` and
   `fid issue get <issue> --section concerns --format json` show
   `status="addressed"` but do not expose `addressed_in_completion_id`.

VAL-011 is BLOCKED, not PASS, because no completed human dogfood evidence exists
yet. The human checklist has been generated for that follow-up.

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Validation contract:
  `project-docs/docs/phases/phase-9/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-9/qa/qa-plan.md`
- Mid-phase review:
  `project-docs/docs/phases/phase-9/qa/review-1.md`
- Task specs and reports:
  `project-docs/docs/phases/phase-9/tasks/task-1.md`,
  `task-1-done.md`, `task-2.md`, `task-2-done.md`, `task-3.md`,
  `task-3-done.md`
- Current decisions and escalations:
  `project-docs/docs/current/decisions_active.md`,
  `project-docs/docs/current/tl_escalations_current.md`
- Triage:
  `project-docs/docs/phases/phase-9/triage-log.md`

## Findings

### Blocking

- **VAL-004 / QA-12 public read-back omits `addressed_in_completion_id`.**
  Evidence from `/tmp/fid` CLI transcript in workspace
  `/tmp/fid-phase9-final-ae_70yx3`: after addressing `CON-7d528c` in
  completion `TC-9ce1a3`, `fid concern list --issue FID-61c1b9 --format json`
  returned keys `body_md_path, created_at, id, issue_id, raised_by_role,
  raised_in_completion_id, status, summary, updated_at`; `fid issue get
  FID-61c1b9 --section concerns --format json` returned keys `body_md_path,
  created_at, id, raised_by_role, raised_in_completion_id, status, summary,
  updated_at`. Neither surface exposed `addressed_in_completion_id`.
  Source check confirms `internal/issues/concern_list.go` `ConcernListItem` and
  `internal/issues/detail.go` `IssueConcern` omit that field.
- **Phase benchmark gate fails.** `go test -bench . -benchmem ./...` exited 1:
  `--- FAIL: TestConcernListFiltersByIssueAndStatus/all_concerns`;
  `concern_test.go:331: first concern issue_id = 001-bc4a28, want "001-b1c4f0"`.
  The required phase gate is therefore not green.

### Major

- None separate from the blocking findings.

### Improvement

- `review-1.md` recorded three non-blocking report-format improvements:
  Task 1 and Task 2 completion reports omit the template top-level title, and
  Task 3 VAL entries cite evidence without explicit PASS labels.

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | Empty output. |
| Lint/static analysis | `golangci-lint run ./...` | PASS | `0 issues.` |
| Unit race suite | `go test -race ./...` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid 24.303s`; `internal/issues 15.585s`; `internal/journal 18.040s`; all listed packages passed. |
| Coverage gate | `go test -race -coverprofile=coverage.out ./...` | PASS | `cmd/fid coverage: 77.8%`; `internal/issues 81.8%`; `internal/journal 80.3%`; `internal/policy 85.0%`; `internal/workflow 95.5%`. |
| Coverage report | `go tool cover -func=coverage.out` | PASS | Total `76.5%`; threshold rows: `cmd/fid` package 77.8% from coverage gate, `internal/issues` 81.8%, `internal/journal` 80.3%, `internal/policy` 85.0%, `internal/workflow` 95.5%. |
| Integration | `go test -tags integration -race ./...` | PASS | All listed packages passed, including `cmd/fid` and `internal/issues`. |
| End-to-end targeted | `go test -tags integration -race ./cmd/fid -run 'TestV01C(ConcernRebuildParity\|ClosureSmoke)' -count=1 -v` | PASS | `=== RUN TestV01CConcernRebuildParity` PASS; `=== RUN TestV01CClosureSmoke` PASS; package `ok .../cmd/fid 2.799s`. |
| Rule seam targeted | `go test -tags integration -race ./internal/issues -run TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites -count=1 -v` | PASS | Test ran and passed; package `ok .../internal/issues 1.590s`. |
| Benchmarks | `go test -bench . -benchmem ./...` | FAIL | Fails in `cmd/fid` test: `first concern issue_id = 001-bc4a28, want "001-b1c4f0"`; journal benchmarks themselves ran, but command exit was 1. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Command exited 0; `ls -l /tmp/fid` shows executable size `11981588`. |
| Full phase check | Full command list from `AGENTS.md` and QA plan | FAIL | Blocked by benchmark command failure and QA-12 read-back failure. |

## Validation Contract Outcomes

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | QA-09 CLI transcript: valid raise produced `CON-7d528c`, `issue_id=FID-61c1b9`, `raised_in_completion_id=TC-911010`, `raised_by_role=developer`, `summary="QA concern with body"`, `body_md_path=/tmp/fid-phase9-final-ae_70yx3/concern-body.md`, `status=open`; read-back via issue section and concern list. |
| VAL-002 | PASS | QA Lead | QA-10 CLI transcript: unknown role, missing issue, missing completion, and completion-on-other-issue each exited 1 with clear errors and left journal lines at 5 and concern rows at 0; valid raise added one journal line and one row. |
| VAL-003 | PASS | QA Lead | QA-09/QA-11 transcript: before raise, issue get/list both showed `[false,0]`; after raise, both showed `[true,1]`. |
| VAL-004 | FAIL | QA Lead | QA-12 transcript: address retained the row and status changed to `addressed`, but public read-back omitted `addressed_in_completion_id` on both `fid concern list` and `fid issue get --section concerns`. |
| VAL-005 | PASS | QA Lead | QA-13 transcript: unknown role, missing concern, missing completion, and already-addressed concern each exited 1 with clear errors and unchanged journal/row counts; valid address added one journal line and no new concern row. |
| VAL-006 | PASS | QA Lead | QA-14 transcript: two open concerns showed get/list `[true,2]`; after one address, get/list `[true,1]`; after last address, get/list `[false,0]`. |
| VAL-007 | PASS | QA Lead | Targeted integration output: `TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites` PASS under `-tags integration -race`; source test asserts unsatisfied named `CON-` IDs for open/partial states and satisfied after all addressed. |
| VAL-008 | PASS | QA Lead | QA-16 CLI transcript: filters returned expected IDs; `--issue FID-0dae6f --status open` returned only `CON-3b7f10` in `json`, `toon`, `markdown`, and `yaml`; empty issue+status returned empty collection with exit 0 in all formats. Benchmark gate failure remains blocking but is tracked separately. |
| VAL-009 | PASS | QA Lead | QA-17 transcript: pre-delete snapshots of concern list and issue get/list JSON structurally equaled post-`fid rebuild --by developer` snapshots; rebuild output ended `rebuild complete`; concern count 4. Targeted `TestV01CConcernRebuildParity` also PASS. |
| VAL-010 | PASS | QA Lead | Targeted e2e output: `TestV01CClosureSmoke` PASS under `-tags integration -race`; test covers AC add, AC verify, AC list, issue complete, concern raise, flag lights, concern address, flag clears, concern list. |
| VAL-011 | BLOCKED | Human | `human-checklist.md` generated with HUMAN-011a..e, but no actual human dogfood evidence exists. Per request and QA plan, VAL-011 is not marked PASS. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| PRD-DEV-005 structured concern linked to issue and completion | MET | VAL-001/002 PASS; QA-09/QA-10 CLI evidence. |
| PRD-DEV-006 concern lifecycle to addressed | PARTIAL | Address mutation works and reject paths pass (VAL-005), but VAL-004 FAILS because the addressing completion is not visible in public read-back. |
| PRD-DEV-007 open concerns surface as derived flags/counts | MET | VAL-003/006 PASS; QA-09/QA-14 evidence. |
| PRD-WF-009 `concerns_addressed` rule | MET | VAL-007 PASS at the rule-evaluation seam; no default edge added. |
| TRD schema `concerns` and write-path seam | PARTIAL | Validate-before-append passes for raise/address, but the public read model omits a lifecycle field required by VAL-004. |
| TRD `issues_enriched` derived columns | MET | Issue get/list false/0, true/count, and final false/0 states verified. |
| TRD default workflow graph / DEC-A | MET | `internal/workflow/workflow.go` contains no `concerns_addressed` edge; `RuleConcernsAddressed` remains evaluated through `internal/issues/transition_rules.go`. |
| ADR-0001 journal canonical, DB derived | MET | VAL-009 rebuild parity PASS. |
| CLI design section 6 concern commands | PARTIAL | Raise/address/list commands exist and most shapes pass; read-back after address omits `addressed_in_completion_id` required by the phase validation contract. |
| Roadmap v0.1c closure | BLOCKED | VAL-010 deterministic smoke PASS, but benchmark gate fails and VAL-011 human dogfood is not completed. |
| DEC-E / CF-011 no concern `activity_log` scope | MET | Source search over concern files and Phase 9 reports found no concern `activity_log` implementation or assertion. |

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS | Task 1 report records exact required gates: `gofmt`, `golangci-lint`, coverage test, `go tool cover`, build PASS, integration N/A with reason; coverage `cmd/fid 78.9%`, `internal/issues 81.8%`, `journal 80.3%`, `policy 85.0%`, `workflow 95.5%`. |
| QA-02 | PASS | `gofmt -l cmd internal` empty output. |
| QA-03 | PASS | `golangci-lint run ./...` output `0 issues.` |
| QA-04 | PASS | `go test -race ./...` all packages passed. |
| QA-05 | PASS | Coverage gate and `go tool cover -func=coverage.out` passed; COV-001/COV-002 thresholds satisfied. |
| QA-06 | PASS | `go test -tags integration -race ./...` passed. |
| QA-07 | FAIL | `go test -bench . -benchmem ./...` exited 1 due `TestConcernListFiltersByIssueAndStatus/all_concerns` failure in `cmd/fid`. |
| QA-08 | PASS | `go build -o /tmp/fid ./cmd/fid` exited 0; `/tmp/fid` executable exists. |
| QA-09 | PASS | CLI transcript proves valid raise, `CON-` ID shape, issue/completion/role/summary/body/status fields through issue section and concern list. |
| QA-10 | PASS | CLI transcript proves four invalid raise cases exit non-zero with clear errors and unchanged journal/row counts; valid raise adds exactly one journal line and one concern row. |
| QA-11 | PASS | CLI transcript proves issue get/list false/0 before raise and true/1 after raise. |
| QA-12 | FAIL | CLI transcript proves status changed to `addressed` and row count unchanged, but neither public read-back surface includes `addressed_in_completion_id`. |
| QA-13 | PASS | CLI transcript proves invalid address cases exit non-zero with unchanged journal/row counts; valid address applies one mutation. |
| QA-14 | PASS | CLI transcript proves partial and final derived flag states on both issue get and issue list. |
| QA-15 | PASS | Targeted integration output for `TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites` PASS; source assertions name open `CON-` IDs and satisfied final state. |
| QA-16 | PASS | CLI transcript proves issue/status/combined filters, all four formats, and empty-match exit 0 with empty collection. |
| QA-17 | PASS | CLI transcript proves delete DB + `fid rebuild --by developer` structural equality for concern rows and derived issue states. |
| QA-18 | PASS | Targeted e2e output for `TestV01CClosureSmoke` PASS. |
| QA-19 | PASS | Source search: concern domain functions in `internal/issues/concern_raise.go`, `concern_address.go`, `concern_list.go`; CLI wiring in `cmd/fid/concern.go`; `workflow.go` has no `concerns_addressed` transition edge; no concern `activity_log` implementation or assertion found. |
| QA-20 | PASS | Task 2 report records exact required gates, integration PASS, build PASS, and coverage thresholds. |
| QA-21 | PASS WITH IMPROVEMENT | Task 3 report records exact required gates including integration/build PASS and benchmark N/A with reason; coverage thresholds recorded. Review-1 improvement remains: VAL rows should carry explicit PASS labels. |
| QA-22 | PASS | `human-checklist.md` generated with only HUMAN-011a..e rows for VAL-011 observables. |
| QA-23 | PASS | This final review records outcomes for VAL-001..VAL-011; VAL-011 is BLOCKED, not PASS. |

## Task Report Command Coverage

| Task | Required commands | Reported coverage | QA result |
|---|---|---|---|
| Task 1 | `gofmt -l cmd internal`; `golangci-lint run ./...`; `go test -race -coverprofile=coverage.out ./...`; `go tool cover -func=coverage.out`; build required; integration conditional | All required commands reported PASS; build PASS; integration N/A because no integration-tagged tests were added | PASS |
| Task 2 | Same base gate; build required; integration required if tagged test added | All required commands reported PASS; build PASS; integration PASS for rule-seam test | PASS |
| Task 3 | `gofmt`; `golangci-lint`; coverage; `go tool cover`; integration; build; benchmark conditional | Required commands reported PASS; benchmark N/A because no new hot query path was added | PASS WITH IMPROVEMENT. Phase-level benchmark still fails and blocks this final review. |

Commit verification from git:

- `f35b3febd7a9976c21be0370d3d701c6aec8e008` -
  `feat(concern): add concern raise command`
- `3ea2d4f89002dfe14620d2ddedfe2eb108fef7ee` -
  `feat(concern): add concern address command`
- `a172634ceb47ed55ddcaf8f01d40e88ad0046129` -
  `feat: add concern list and v0.1c closure tests`

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: YES,
  recorded in `qa-plan.md` section 7.
- Placeholder commands found: none.
- Invalid command syntax found during final review: none.
- Feasible but failing required command: `go test -bench . -benchmem ./...`
  failed due a `cmd/fid` test failure.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC References |
|---|---|---|---|---|
| `review-1.md` | CONTINUE | N/A - no blocking findings | Three report-format improvements remain deferrable | none |

Implication: final review was allowed to proceed, but the fresh phase-end gate
and deterministic checklist supersede the mid-phase ON TRACK statuses where they
found failures.

## Documentation Review

- `06_DESIGN_PLAN_CLI.md` was updated by Tasks 1-3 for concern raise/address/list.
- Task completion reports exist for all three tasks.
- `triage-log.md` records TRG-001 RESUME for Task 3 expanded write set.
- Current `tl_escalations_current.md` lists no open escalations.
- Documentation gap: the validation contract requires addressed completion
  read-back, but the public concern read models do not expose that field.

## Test Quality Review

- Red/green evidence is present in Task 3 for concern list and format support.
- Reject-path tests and final CLI transcript cover raise/address journal
  poisoning risks.
- Rebuild parity and v0.1c smoke have integration coverage.
- Test quality defect: `TestConcernListFiltersByIssueAndStatus/all_concerns`
  assumes first-row ordering that failed under the required benchmark command.
  This must be resolved before the phase gate can pass.
- Coverage thresholds remain above DEC-007 floors.

## Risk Review

| Risk | Status | Evidence |
|---|---|---|
| Journal poisoning on rejected concern writes | Controlled | QA-10 and QA-13 show unchanged journal/row counts on invalid raise/address cases. |
| Rebuild divergence | Controlled | QA-17 and `TestV01CConcernRebuildParity` PASS. |
| Public CLI/output shape drift | Blocking | QA-12 shows missing `addressed_in_completion_id` in address read-back. |
| Workflow-edge scope creep | Controlled | `workflow.go` default graph has no `concerns_addressed` edge; rule seam only. |
| Deferred CF-011 activity timeline | Controlled | No concern `activity_log` assertion or implementation introduced. |
| Test determinism / phase benchmark gate | Blocking | Benchmark command fails due a concern-list test ordering expectation. |

## Human Checklist Readiness

Generated at:
`project-docs/docs/phases/phase-9/qa/human-checklist.md`.

- Contains only `(human)` items: YES.
- Uses one observable per row: YES.
- Human-only VAL assertions mapped: YES, VAL-011.
- VAL-011 final outcome: BLOCKED until a human completes HUMAN-011a..e.

## Coverage Caveat

The deterministic PASS items prove asserted mechanics, data states, command
outputs, and structural conformance checks only. They do not prove dogfooding
quality or product fit. VAL-011 remains unproven until the human checklist is
completed.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | BLOCKED | Required QA artifacts now exist, but final review has blocking findings and VAL-011 is not completed. |
| Automated gates | BLOCKED | Benchmark phase gate fails. |
| Held-out scenarios | DEFERRED | Held-out scenario harness remains a scheduled commitment before release readiness. |
| Release readiness | DEFERRED | Security scan, dependency audit, license check, secret scan, threat model, and release checklist are release-readiness controls, not Phase 9 gates. |

## Final Lines

```text
Review written: project-docs/docs/phases/phase-9/qa/final-review.md
Human checklist written: project-docs/docs/phases/phase-9/qa/human-checklist.md
Blocking: 2  Major: 0  Improvement: 3
VAL outcomes: PASS=9 FAIL=1 BLOCKED=1 CARRY-FORWARD=0
VERDICT: BLOCK
```
