# Phase 10 Final Review - Permissions Matrix + Policy

Reviewer: QA Lead
Date: 2026-06-21

## Verdict

PASS for deterministic QA. VAL-009 is CARRY-FORWARD to the generated human
checklist because it requires human dogfood judgment and cannot be passed by QA
inference.

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Validation contract:
  `project-docs/docs/phases/phase-10/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-10/qa/qa-plan.md`
- Task specs and completion reports:
  `project-docs/docs/phases/phase-10/tasks/`
- Mid-phase reviews:
  `project-docs/docs/phases/phase-10/qa/review-1.md`,
  `review-2.md`, `review-3.md`, `review-3-addendum.md`
- Current decisions and escalation state:
  `project-docs/docs/phases/phase-10/decisions.md`,
  `project-docs/docs/current/decisions_active.md`,
  `project-docs/docs/current/tl_escalations_current.md`

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | No output. |
| Lint/static analysis | `golangci-lint run ./...` | PASS | `0 issues.` |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | `ok github.com/mbell-nzl/fid/cmd/fid (cached) coverage: 77.8%`; `internal/issues` 81.9%, `internal/journal` 80.3%, `internal/policy` 81.2%, `internal/workflow` 95.5%. |
| Coverage report | `go tool cover -func=coverage.out` | PASS | `total: (statements) 76.3%`; threshold packages also shown in the race coverage output above. |
| Integration | `go test -tags integration -race ./...` | PASS | `ok` for `cmd/fid`, `internal/config`, `internal/db`, `internal/format`, `internal/issues`, `internal/journal`, `internal/policy`, `internal/workflow`; `internal/bootstrap [no test files]`. |
| End-to-end | N/A | N/A | Phase 10 QA plan records no accepted end-to-end or held-out scenario harness for this phase. |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | `BenchmarkRebuildWorkspace1kEvents-6 ... 388.1 ms/rebuild`; `BenchmarkRebuildWorkspace10kEvents-6 ... 3805 ms/rebuild`; all packages PASS. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | Exit 0; `/tmp/fid` was then used for manual CLI evidence. |
| Full phase check | Full command list from QA plan section 9 | PASS | All listed phase-level commands above passed in this final-review session. |

## Validation Contract Outcomes

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | Task 1 report records corrupt-tail, first-line corruption, un-upgradeable tail, clean journal, and CLI evidence. QA rerun excerpt: `rc=1`; `error: malformed journal entry at .fid/journal/2026-06-21.jsonl:2`; `read stopped - manual repair required`. Review 1 verdict CONTINUE. |
| VAL-002 | PASS | QA Lead | Task 2 report and Review 2 confirm `fid policy show` json/toon/markdown/yaml parity and implemented-action coverage. QA excerpt: json row for `issue.create.epic` has `developer:false`, `tl:true`, `human:true`; markdown header renders `| action | developer | qa_reviewer | qa_lead | tl_lite | tl | workflow_runner | human |`; read-only counts stayed `journal_lines=1 activity_log=0` before and after. |
| VAL-003 | PASS | QA Lead | Task 2 report and Review 2 cite `TestMatrixContainsExactlyTheImplementedV01WriteActions`, `TestAuthorizePermitsAllowedCombinationAndRejectsDeniedCombination`, and preserved missing/unknown-role tests. Final phase `go test -race -coverprofile=coverage.out ./...` passed. |
| VAL-004 | PASS | QA Lead | QA-06 manual excerpt below: denied `issue.create.epic` by `developer` exited 1, named role/action/not permitted, and kept `journal_lines=1`, `issues=0`; permitted `tl` counterpart exited 0 and advanced to `journal_lines=2`, `issues=1`, `epic_metadata=1`. Review 3 blocker cleared by addendum and committed reports. |
| VAL-005 | PASS | QA Lead | QA-07 manual excerpt below: missing `--by` returned `error: --by is required`; unknown role listed known roles; both preserved counts; permitted `developer` task create advanced `journal_lines=2 -> 3` and `issues=1 -> 2`. |
| VAL-006 | PASS | QA Lead | QA-08 manual excerpt below: `--by human` epic create succeeded; `workflow_runner` AC verify denial preserved counts; `--by human` AC verify succeeded and DB read-back returned `verification_role=human|status=pass`. |
| VAL-007 | PASS | QA Lead | QA-09 and QA-10 excerpts below: AC verify permitted `qa_reviewer` and denied `workflow_runner` for `issue.ac.verify.workflow_runner` with zero count change; transition denied `qa_reviewer` for `issue.transition.dev_complete`, then permitted `developer` after setup. |
| VAL-008 | PASS | QA Lead | QA-11 excerpt below: `fid policy show --format json` displayed `issue.create.epic` with `developer:false`, `tl:true`; corresponding developer write failed and tl write succeeded. |
| VAL-009 | CARRY-FORWARD | Human | Human checklist generated at `project-docs/docs/phases/phase-10/qa/human-checklist.md` with HUMAN-09a..HUMAN-09d. QA did not mark this PASS because it requires human dogfood judgment. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| PRD-ROLE-001 - `--by` required on writes | MET | VAL-005; QA-07 missing-role CLI excerpt. |
| PRD-ROLE-002 - unknown roles rejected clearly | MET | VAL-001, VAL-005; QA-07 unknown-role excerpt lists known roles. |
| PRD-ROLE-003 - denied permissions rejected with role/action/not permitted | MET | VAL-004, VAL-007, VAL-008; QA-06/09/10/11 excerpts. |
| PRD-ROLE-004 - `human` escape valve without verifier impersonation | MET | VAL-006; QA-08 DB read-back `verification_role=human`. |
| PRD-ROLE-005 - display active permissions matrix | MET deterministically; human readability pending | VAL-002 PASS; VAL-009 CARRY-FORWARD for dogfood readability. |
| TRD-ARCH-002 / CF-003 - malformed journal line is corruption | MET | VAL-001; corrupt-tail CLI reproduction and Task 1 tests. |
| TRD-ARCH-010 - code matrix, complete write actions, denied combos error | MET | VAL-002, VAL-003, VAL-004, VAL-007, VAL-008. |
| CF-007 - validate before append | MET | QA-06, QA-07, QA-09, QA-10 row-count excerpts show zero append on denials. |
| Design CLI section 11 - `fid policy show` formats | MET | VAL-002; final `policy show` format excerpts and Task 2 tests. |

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS | Task 1 report records corrupt tail, first-line corruption, un-upgradeable tail, and clean role journal evidence; Review 1 marked VAL-001 ON TRACK and verdict CONTINUE. |
| QA-02 | PASS | Final CLI reproduction: append `{` to `.fid/journal/2026-06-21.jsonl`; `/tmp/fid issue create --kind task --title "Corrupt tail" --by developer`; `rc=1`; `error: malformed journal entry at .fid/journal/2026-06-21.jsonl:2`; `read stopped - manual repair required`. |
| QA-03 | PASS | Final `policy show` excerpts: json, toon, markdown, yaml each exited 0. Markdown header: `| action | developer | qa_reviewer | qa_lead | tl_lite | tl | workflow_runner | human |`. JSON includes all implemented rows through `maintenance.rebuild`; Task 2 tests assert format parity. |
| QA-04 | PASS | Final read-only check: before `fid policy show`, `journal_lines=1 activity_log=0`; after json/toon/markdown/yaml show commands, `journal_lines=1 activity_log=0`. `TestPolicyShowIsReadOnlyAndDoesNotAcceptByFlag` also proves `--by` is an unknown flag. |
| QA-05 | PASS | Task 2 report and Review 2 cite matrix completeness and `Authorize` permit/deny tests; final full race coverage gate passed. |
| QA-06 | PASS | Literal excerpt: `qa06-before-denied journal_lines=1 issues=0 ...`; denied command `rc=1`, `error: role "developer" not permitted for action "issue.create.epic"`; after denied counts unchanged. Permitted command `rc=0`, `created issue QA-9c3e37 (epic) "Permitted epic"`; after permitted `journal_lines=2 issues=1 epic_metadata=1`. |
| QA-07 | PASS | Literal excerpt: missing `--by` `rc=1`, `error: --by is required`, counts unchanged at `journal_lines=2 issues=1`. Unknown role `rc=1`, `known roles: developer, human, qa_lead, qa_reviewer, tl, tl_lite, workflow_runner`, counts unchanged. Permitted task create `rc=0`, `created issue QA-ebdb13 (task) "Permitted task"`, then `journal_lines=3 issues=2`. |
| QA-08 | PASS | Literal excerpt: human epic `rc=0`, `created issue QA-c46411 (epic) "Human epic"`. `workflow_runner` AC verify denied with `issue.ac.verify.workflow_runner`, counts unchanged at `journal_lines=5 ac_verifications=0`. Human AC verify `rc=0`, `verified AC-0aa322 on QA-ebdb13 by human: pass`; after `journal_lines=6 ac_verifications=1`; DB read-back `verification_role=human|status=pass`. |
| QA-09 | PASS | Literal excerpt: `qa_reviewer` AC verify `rc=0`, `verified AC-4486ed on QA-ebdb13 by qa_reviewer: fail`, counts `journal_lines=8 ac_verifications=2`. `workflow_runner` denial `rc=1`, `error: role "workflow_runner" not permitted for action "issue.ac.verify.workflow_runner"`; after denial counts unchanged at `journal_lines=8 ac_verifications=2`. |
| QA-10 | PASS | Setup created task, transitioned to `in_progress`, and completed it. Denied command `rc=1`, `error: role "qa_reviewer" not permitted for action "issue.transition.dev_complete"`; counts unchanged at `journal_lines=11 activity_log=1 task_completions=1`. Permitted command `rc=0`, `transitioned QA-2999e9: in_progress → dev_complete`; after `journal_lines=12 activity_log=2`. |
| QA-11 | PASS | `fid policy show --format json` excerpt: `"action":"issue.create.epic","developer":false,...,"tl":true,...,"human":true`. Denied developer cell `rc=1`, role/action permission error. Permitted tl cell `rc=0`, `created issue QA-04a76b (epic) "Parity permitted"`; final counts `journal_lines=13 issues=5 epic_metadata=3`. |
| QA-12 | PASS | Task 3 report and `decisions.md` Task 3 execution note list DEC-E fixture moves for rebuild, raw journal fixtures, transition fixtures, setup roles, and direct-prerequisite matrix rows. Review 3 marked fixture-role-change note PASS. |
| QA-13 | PASS | Compared task specs with completion reports: Task 1 required format/lint/race coverage/cover plus integration/build; Task 2 required format/lint/race coverage/cover plus build and conditional integration N/A; Task 3 required format/lint/race coverage/cover/integration/build plus conditional benchmark. Reports record PASS or justified N/A for each; Review 1/2/3 also confirmed command coverage. |
| QA-14 | PASS | Full phase gate executed in this final review: format, lint, race coverage, cover report, integration, benchmarks, and build all PASS. |
| QA-15 | PASS | COV-001/COV-002 from final gate: `cmd/fid` 77.8% >= 70%; `internal/issues` 81.9%, `internal/journal` 80.3%, `internal/policy` 81.2%, `internal/workflow` 95.5% all >= 80%. |
| QA-16 | PASS | Task reports record docs or verified-no-change notes for `06_DESIGN_PLAN_CLI.md`, `01_PRD.md`, Phase 10 decisions, and `03_TRD_SCHEMA.md` direct-prerequisite rows. |
| QA-17 | PASS | `git diff --name-only 716329cc4cbeda174ce97dcb4e6f3192720dc3b9^..5293fd2674b0afdc902feb2f7079d12b45e82551` lists code/docs/testdata files and no `migrations/` path. |

## Command And Feasibility Validation Result

- Command and feasibility validation completed before implementation: YES.
- Invalid or unvalidated command, test scope, fixture, artifact, or human
  checklist step found: none.
- Placeholder command defect found: none. Phase commands in AGENTS.md and the QA
  plan are concrete for Phase 10.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC References |
|---|---|---|---|---|
| review-1.md | CONTINUE | N/A | One reproducibility improvement for CLI setup carried into final QA by reproducing corrupt-tail CLI evidence. | none |
| review-2.md | CONTINUE | N/A | Matrix row-cell hard-coding improvement deferred to future matrix changes; not a Phase 10 blocker. | none |
| review-3.md | BLOCK | YES | No improvement findings. Blocking commit-field drift cleared by addendum and committed reports. | ESC-007 |
| review-3-addendum.md | CONTINUE | YES | Review 3's blocker cleared. | none |

Review 3 blocker clearance was independently verified in this final review:
`git rev-parse --verify eeaf9faed754bea9fe24d2ce7c8616d4b554959d^{commit}`
returned the reconciliation commit; `git show -s --format=%H%n%s` matched the
three task implementation commits and the reconciliation commit; `rg -n
"Pending commit|Commit hash|Commit subject" ...task-*-done.md` returned only
the concrete `Commit hash`/`Commit subject` headings and no `Pending commit`.

## Documentation Review

- Impacted docs updated: YES. Phase decisions and `03_TRD_SCHEMA.md` were
  updated where behavior or direct-prerequisite matrix rows changed; PRD and CLI
  docs were recorded as verified-no-change where already accurate.
- Specs needing Tech Lead wrap-up sync: `project-docs/docs/current/phase_brief.md`
  and `progress_current.md` still describe Phase 10 as planned; phase wrap should
  refresh current context and archive resolved ESC-007.
- Completion reports present for all tasks: YES.

## Human Checklist Readiness

Generated at
`project-docs/docs/phases/phase-10/qa/human-checklist.md`.

- Contains only `(human)` items: YES.
- Uses one observable per row: YES.
- Human-only `VAL-NNN` assertions mapped: YES, VAL-009 via HV-09a..HV-09d.
- Product/perception/dogfooding risks to call out: VAL-009 remains unconfirmed
  until the human checklist is executed.

## Coverage Caveat

A deterministic PASS proves only the asserted mechanics, data states, command
outputs, and structural conformance checks. It does not prove matrix readability
or denial-message actionability as human product judgment; those remain
unconfirmed until the Human checklist records outcomes.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | READY WITH HUMAN CARRY-FORWARD | QA plan, validation contract, task reports, reviews, final review, and human checklist exist; VAL-009 awaits human execution. |
| Automated gates | READY | Final phase gate PASS. |
| Held-out scenarios | DEFERRED | No accepted held-out scenario harness or command for Phase 10. |
| Release readiness | DEFERRED | Security scan, vulnerability audit, license check, secret scan, threat model, and release checklist are not required before external release per current quality plan. |

## Findings

### Blocking

- None.

### Major

- None.

### Improvement

- None for Phase 10 final QA. Review 1 and Review 2 improvement notes were
  non-blocking and are either covered by final evidence or deferred to future
  matrix changes.

## Escalations

- No unresolved escalation remains for Phase 10 final QA.
- ESC-007 is still present in `project-docs/docs/current/tl_escalations_current.md`
  under Open Escalations, but its own status says resolved pending the mechanical
  reconciliation commit. The reconciliation commit now exists and Review 3
  addendum plus final-review verification clear the blocker. Phase wrap should
  archive ESC-007 to the phase escalation history and remove it from current
  open escalations.

## Unresolved Risks And Carry-Forward Candidates

- VAL-009: CARRY-FORWARD to human dogfood checklist.
- Current-context drift: `phase_brief.md` and `progress_current.md` still frame
  Phase 10 as planned; Tech Lead wrap should refresh them after final/human
  outcomes are accepted.
