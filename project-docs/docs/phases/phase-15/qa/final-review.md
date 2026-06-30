# Phase 15 Final Review - `fid dev *` Composites + Cleanup

Reviewer: QA Lead
Date: 2026-06-29

## Verdict

PASS WITH IMPROVEMENTS

Deterministic phase gate: PASS. VAL-001 through VAL-011 have reproducible QA
evidence. VAL-012 is carried into the phase-end human checklist because no human
dogfood outcome has been recorded yet.

Blocking findings: 0. Major findings: 0. Improvement findings: 2 carry-forward
notes already triaged as RESUME.

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Active decisions: `project-docs/docs/current/decisions_active.md`
- Current progress: `project-docs/docs/current/progress_current.md`
- Current TL escalations: `project-docs/docs/current/tl_escalations_current.md`
- Validation contract: `project-docs/docs/phases/phase-15/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-15/qa/qa-plan.md`
- Mid-phase review: `project-docs/docs/phases/phase-15/qa/review-1.md`
- Triage log: `project-docs/docs/phases/phase-15/triage-log.md`
- Task specs and completion reports: `project-docs/docs/phases/phase-15/tasks/task-1.md` through `task-4.md`, and `task-1-done.md` through `task-4-done.md`
- Commits verified:
  - `0906428adf9bb4d89b41dc94ad5b147b1b7c538f` - `refactor: split issue and replay cleanup paths`; `Refs: phase-15/task-1`
  - `07ac0ae653753f2a12921cf26bf0cde1c12dc6c5` - `feat(dev): add dev start and gate composites`; `Refs: phase-15/task-2`
  - `b9e650d5686e6fc6325448eefdc85d4145450e13` - `feat(dev): add finish and block composites`; `Refs: phase-15/task-3`
  - `64cd9ec0254d0b1bb00cc9305119b810e0787038` - `test(dev): add phase 15 composite smoke`; `Refs: phase-15/task-4`

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | Empty output. |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Output: `0 issues.` |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | `cmd/fid` 76.9%, `internal/dev` 75.4%, `internal/issues` 81.7%, `internal/journal` 80.5%, `internal/labels` 82.4%, `internal/policy` 84.6%; all packages passed. |
| Coverage report | `go tool cover -func=coverage.out` | PASS | Total coverage: `75.8%`. |
| Integration | `go test -tags integration -race ./...` | PASS | All packages passed, including `cmd/fid` and `internal/dev`. |
| Focused e2e smoke | `go test -tags integration -race ./cmd/fid -run TestV01EDevCompositeSmokeAndRebuildParity -count=1` | PASS | Output: `ok github.com/mbell-nzl/fid/cmd/fid 2.600s`. |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | All packages passed; `BenchmarkRebuildWorkspace1kEvents` and `BenchmarkRebuildWorkspace10kEvents` ran in `internal/journal`. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Empty output. |
| Full phase check | QA-plan phase-end command list | PASS | All required phase-end commands passed, including the benchmark command. |

## Validation Contract Outcomes

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | Task 1 report and Review-1 confirm behavior-preserving split. Final `wc -l`: `cmd/fid/issue.go` 22, largest split issue file `issue_helpers.go` 228, `internal/journal/replay.go` 205, `projection.go` 192, `inspect.go` 128. Phase gate passed. |
| VAL-002 | PASS | QA Lead | Final function-span command shows target label functions <=60 lines: `internal/labels/add.go` `Add` 18, `prepareLabelAdd` 20, `addLabelLocked` 26; `internal/issues/label.go` `AddLabel` 28, `RemoveLabel` 28, `prepareIssueLabelWrite` 35, `issueLabelWriteLocked` 26. Phase gate passed. |
| VAL-003 | PASS | QA Lead | Task 2 report cites `internal/dev/start_test.go` one `dev.started` event with claim insert, transition/activity, and open run; `cmd/fid/dev_test.go` verifies CLI confirmation and projected `in_progress`. Phase gate passed with `internal/dev` 75.4%. |
| VAL-004 | PASS | QA Lead | Task 2 report cites `internal/dev/start_test.go` valid unclaimed `todo` success plus zero-append rejection for conflicting live claim and invalid edge. Review-1 focused start/gate tests passed. |
| VAL-005 | PASS | QA Lead | Task 2 report cites `internal/dev/gate_test.go` one `dev.gate.recorded` event, active attempt derivation, open-run linking, and zero-append rejects. `cmd/fid/dev_test.go` verifies the command string is recorded and not executed. |
| VAL-006 | PASS | QA Lead | `rg` proof: `internal/policy/matrix.go` rows for `dev.start`, `dev.gate`, `dev.finish`, `dev.block`; `project-docs/specs/03_TRD_SCHEMA.md` section 13 has the same developer/human roles. `cmd/fid/dev_test.go` denies `qa_reviewer` for all four actions. Tasks 2-3 reports record PASS. |
| VAL-007 | PASS | QA Lead | Task 3 report cites `TestFinishClosesRunCompletesTransitionsAndReleasesClaimInOneEvent`, proving one `dev.finished` event closes run, writes completion snapshot, transitions to `dev_complete`, and releases claim. Phase gate passed. |
| VAL-008 | PASS | QA Lead | Task 3 report cites `TestFinishRejectsWhenDeveloperACsUnverifiedBeforeAppend` and `TestFinishSucceedsWithoutPreexistingCompletionOnceDeveloperACsPass`, proving both reject and success states. |
| VAL-009 | PASS | QA Lead | Task 3 report cites `TestBlockRecordsFailedGateCompletionTransitionsAndReleasesClaim`, proving one `dev.blocked` event with `failed_gate` completion, reason, transition, and claim release. |
| VAL-010 | PASS | QA Lead | `internal/issues/complete.go` has `BuildCompletionEffects`; `internal/dev/finish.go` and `internal/dev/block.go` call it with `dev_complete` and `failed_gate`. Final function spans: `Complete` 19 lines, `BuildCompletionEffects` 21, support helpers <=39. `go test -race -coverprofile=coverage.out ./...` passed with `internal/issues` 81.7%. |
| VAL-011 | PASS | QA Lead | Task 4 adds `cmd/fid/e2e_v01e_dev_test.go::TestV01EDevCompositeSmokeAndRebuildParity`. Fresh focused run passed. `rg` confirms assertions for `qa_reviewer` denial, one-line deltas for start/gate/finish/block, and rebuild parity. |
| VAL-012 | CARRY-FORWARD | Human | Generated `project-docs/docs/phases/phase-15/qa/human-checklist.md` with one human observable per row (`HV-15a`..`HV-15d`). No human dogfood result is recorded in this review, so PASS is not claimed. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| CF-014 file split | MET | VAL-001; final `wc -l`; Task 1 report; Review-1. |
| CF-015 label helper extraction | MET | VAL-002; final function spans; Task 1 report. |
| PRD-DEV-008 developer composite lifecycle | MET | VAL-003, VAL-007, VAL-009, VAL-011. |
| PRD-CLAIM-005 started-category live claim | MET | VAL-003, VAL-004, VAL-011. |
| PRD-DEV-004 / PRD-DEV-010 gate recording | MET | VAL-005, VAL-007, VAL-009, VAL-011. |
| PRD-DEV-009 completion reporting | MET | VAL-007, VAL-010. |
| PRD-WF-007..010 transition rule enforcement | MET | VAL-008. |
| TRD-ARCH-010 / schema section 13 policy matrix | MET | VAL-006. |
| CF-013 completion orchestration cleanup | MET | VAL-010. |
| Phase e2e smoke and rebuild determinism | MET | VAL-011. |
| Phase human dogfood | CARRY-FORWARD | VAL-012 checklist generated; human outcome pending. |

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS | Task reports 1-4 record every required task command as PASS, including integration and build. Task 4 correctly marked the task-level benchmark command not applicable because no benchmark was touched; phase-end benchmark was run and passed. |
| QA-02 | PASS WITH NOTE | File/function targets for VAL-001, VAL-002, and VAL-010 passed. Non-blocking note: `internal/dev/Finish` is 109 lines and `Block` is 97 lines; this is outside VAL-010 and is already triaged as TRG-002 RESUME. |
| QA-03 | PASS | Task 1 and Task 3 reports record behavior-preserving refactors; required gates passed. Review-1 found no behavior-change or test-weakening issue. |
| QA-04 | PASS | Task 2 report, Review-1 focused tests, and phase gate cover `dev start`/`dev gate` one-event success and zero-append reject paths. |
| QA-05 | PASS | Matrix rows in `internal/policy/matrix.go` match `03_TRD_SCHEMA.md` section 13; `cmd/fid/dev_test.go` covers show/enforce and denied-role behavior for all four composites. |
| QA-06 | PASS | Task 3 report cites finish success, unmet-AC rejection, and no-pre-existing-completion success after AC verification. |
| QA-07 | PASS | Task 3 report cites block success path with `failed_gate`, reason, transition, and claim release. |
| QA-08 | PASS | `go test -tags integration -race ./...` passed; fresh focused e2e smoke passed. `cmd/fid/e2e_v01e_dev_test.go` asserts finish path, block path, permission denial, and rebuild parity. |
| QA-09 | PASS | `project-docs/specs/02_TRD.md` and `06_DESIGN_PLAN_CLI.md` were updated by Task 2. `project-docs/specs/07_DESIGN_PLAN_FLOWS.md` was updated by Task 4, not intentionally left unchanged; `git show --stat 64cd9ec -- project-docs/specs/07_DESIGN_PLAN_FLOWS.md` shows 12 changed lines to add `--model`, `--exit-code`, `--output-excerpt`, and claim-release output. |

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: YES, in `qa-plan.md` section 6.
- Invalid or unvalidated command, test scope, fixture, artifact, or human checklist step found: none.
- Required phase-end benchmark command was feasible and passed.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC/TRG References |
|---|---|---|---|---|
| `review-1.md` | CONTINUE | N/A; 0 Blocking | Completion-report H1 consistency and Task 2 function-size debt accepted as non-blocking. | TRG-001 RESUME; no ESC. |
| Final review window | PASS WITH IMPROVEMENTS | N/A; 0 Blocking | Task 3 `internal/dev` `Finish`/`Block` function-size debt accepted as non-blocking technical debt. | TRG-002 RESUME; no ESC. |

`project-docs/docs/current/tl_escalations_current.md` still records no open
escalations.

## Documentation Review

- Impacted docs updated: YES.
- `project-docs/specs/02_TRD.md`: updated by Task 2 for the `internal/dev`
  boundary.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: updated by Task 2 for
  `dev start --model`, `dev gate --exit-code`, and `--output-excerpt`.
- `project-docs/specs/07_DESIGN_PLAN_FLOWS.md`: updated by Task 4 to reconcile
  the developer execution flow with the implemented composite CLI. It was not
  intentionally left unchanged.
- Completion reports present for all tasks: YES.
- Specs needing Tech Lead wrap-up sync: none identified by QA beyond normal
  phase archive/current-context refresh.

## Human Checklist Readiness

Generated at
`project-docs/docs/phases/phase-15/qa/human-checklist.md`.

- Contains only `(human)` items: YES.
- Uses one observable per row: YES.
- Human-only `VAL-NNN` assertions mapped: YES, VAL-012 to `HV-15a`..`HV-15d`.
- Product/perception/dogfooding risks to call out: the CLI may still pass
  deterministic tests while a human finds the composite output unclear; this is
  exactly what VAL-012 is designed to check.

## Coverage Caveat

A deterministic PASS proves only the asserted mechanics, data states, command
outputs, and structural conformance checks. It does not prove unasserted workflow
feel or dogfooding clarity; those remain unconfirmed until the human checklist
records outcomes.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | READY | Validation contract, QA plan, task specs, task reports, review, final review, and human checklist exist. |
| Automated gates | READY | Full phase gate passed, including integration, benchmarks, and build. |
| Held-out scenarios | DEFERRED | Phase brief keeps held-out scenario harness out of scope until later release-readiness work. |
| Release readiness | DEFERRED | Threat model, release checklist, SAST, dependency audit, license scan, secret scan, and held-out scenario decision remain later maturity work per `08_QUALITY_PLAN.md`. |

## Findings

### Blocking

- None.

### Major

- None.

### Improvement

- TRG-001: Review-1 improvements remain RESUME: future completion reports should include the template H1; Task 2 structure debt was accepted as non-blocking.
- TRG-002: Task 3 `internal/dev` `Finish` and `Block` remain above the 60-line non-test function target; accepted as non-blocking technical debt because required behavior and gates passed.

## Unresolved Risks And Carry-Forward Candidates

- VAL-012 remains pending human execution in `human-checklist.md`.
- `internal/dev` function-size cleanup is a reasonable follow-up candidate, but it is not a phase blocker.
