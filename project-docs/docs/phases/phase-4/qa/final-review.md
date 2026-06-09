# Phase 4 Final Review - Issue Authoring And Discovery

Reviewer: QA Lead
Date: 2026-06-09

## Verdict

BLOCK

Phase 4 behavior evidence is complete for VAL-001..VAL-009, and ESC-002 is
closed by Task 4 due-date read-back evidence. The phase cannot pass the
automated gate because `internal/issues` coverage is `79.5%`, below the COV-001
domain-package threshold of `80%`.

## Scope And Evidence

- Phase brief: `project-docs/docs/current/phase_brief.md`
- Validation contract: `project-docs/docs/phases/phase-4/qa/validation-contract.md`
- QA plan: `project-docs/docs/phases/phase-4/qa/qa-plan.md`
- Mid-phase review: `project-docs/docs/phases/phase-4/qa/review-3.md`
- Triage log: `project-docs/docs/phases/phase-4/triage-log.md`
- Completion reports: `project-docs/docs/phases/phase-4/tasks/task-1-completion.md` through `task-5-completion.md`
- Task commits verified by `git show --no-patch --format='%H %s' bce37eb 397cb77 2bfe94f c62c6aa 2a281f0`.

## Phase-Level Gate Results

| Gate | Command | Result | Evidence |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | PASS | No output; exit 0. |
| Lint/static analysis | `golangci-lint run ./...` | PASS | `0 issues.` |
| Coverage/confidence | `go test -race -coverprofile=coverage.out ./...` | PASS | Exit 0; `cmd/fid` `79.5%`, `internal/issues` `79.5%`, `internal/journal` `80.7%`, `internal/policy` `85.0%`. |
| Coverage threshold | `go tool cover -func=coverage.out` plus COV-001/COV-002 check | FAIL | Total `75.3%`; COV-002 passes with `cmd/fid 79.5% >= 70%`; COV-001 fails because `internal/issues 79.5% < 80%`. |
| Integration | `go test -tags integration -race ./...` | PASS | Exit 0; packages passed including `cmd/fid 12.548s`, `internal/issues 5.512s`, `internal/journal 18.731s`. |
| End-to-end | N/A | N/A | AGENTS.md records E2E as `N/A` until a held-out scenario harness is accepted. |
| Benchmarks | `go test -bench . -benchmem ./...` | PASS | Exit 0; `BenchmarkRebuildWorkspace1kEvents` `442.4 ms/rebuild`; `BenchmarkRebuildWorkspace10kEvents` `3812 ms/rebuild`. |
| Build/package | `go build -o /tmp/fid ./cmd/fid` | PASS | No output; exit 0. |
| Full phase check | Full Phase 4 gate from QA plan | FAIL | Only failing component is COV-001 coverage for `internal/issues`. |

## Validation Contract Outcomes

| VAL ID | Outcome | Validator | Evidence |
|---|---|---|---|
| VAL-001 | PASS | QA Lead | Task 1 raw SQL/projection matrix in `task-1-completion.md`; Task 3 CLI errors for sequential/prioritised/unattached/metadata-less states; phase gate `go test -race -coverprofile=coverage.out ./...` passed. |
| VAL-002 | PASS | QA Lead | Task 2 `TestCreateSupportsEpicTaskAndSubtaskAuthoring` and `TestCreateRejectsInvalidHierarchyAndUnknownLabels`; Task 3 CLI create/get smoke for epic, task, unattached task, and subtask. |
| VAL-003 | PASS | QA Lead | Task 2/3 persistence evidence plus Task 4 closeout: `TestGetReturnsCentralizedDetailShape` proves populated `due_date` `2026-06-30` and omitted `nil` through `internal/issues.Get`; `TestIssueGetSectionsAndDueDateReadBack` proves `fid issue get --format json` returns `"due_date":"2026-06-30"` and omitted `null`. |
| VAL-004 | PASS | QA Lead | Task 2 proves exactly one epic metadata row for epics and rejects non-epic scheduling; Task 3 proves CLI scheduling metadata and missing-metadata error behavior. |
| VAL-005 | PASS | QA Lead | Task 2 proves one `issue.created` journal event per create and delete-DB rebuild parity; Task 3 proves one CLI journal line per create and `fid rebuild --by developer` parity. |
| VAL-006 | PASS | QA Lead | Task 4 `TestGetReturnsCentralizedDetailShape`; markdown parity tests; adapter-boundary scan `rg -n 'SELECT|QueryRow|Query\(' cmd/fid/issue.go` returned no matches. |
| VAL-007 | PASS | QA Lead | Task 4 `TestIssueGetSectionsAndDueDateReadBack` covers `acceptance`, `concerns`, `completions`, `links`, `claims`, empty sections, and unknown-section rejection with supported names. |
| VAL-008 | PASS | QA Lead | Task 5 `TestListFiltersByKind`, `TestListFiltersAndOrdering/filters`, and `TestIssueListFiltersAndOrdering` cover include/exclude evidence for kind, status, category, epic, label, priority, inbox, open concerns, and empty result exit 0. |
| VAL-009 | PASS | QA Lead | Task 5 `TestListFiltersAndOrdering/ordering` and `TestIssueListDispatchesFormatsAndRejectsUnknownFormats` cover stable order, JSON/TOON/YAML parity, markdown structural output, and `--format xml` rejection. |
| VAL-010 | CARRY-FORWARD | Human | Human-only dogfood judgment is not executed by QA Lead. Checklist generated at `project-docs/docs/phases/phase-4/qa/human-checklist.md`; `/tmp/fid` build passed, so the checklist is executable. |

## Requirement Compliance

| Requirement | Status | Evidence |
|---|---|---|
| PRD-ISSUE-001..004, PRD-WF-005..006 | MET | VAL-002, VAL-004. |
| PRD-ISSUE-013..016, PRD-ISSUE-019..020, CF-004 | MET | VAL-001, VAL-003, VAL-005. |
| CF-006, PRD-ISSUE-006, PRD-ISSUE-018 | MET | VAL-006, VAL-007. |
| PRD-ISSUE-007, PRD-OUT-004 | MET | VAL-008, VAL-009. |
| 08_QUALITY_PLAN COV-001 | NOT MET | `internal/issues` coverage `79.5%`; threshold is `80%`. |
| 08_QUALITY_PLAN COV-002 | MET | `cmd/fid` coverage `79.5%`; threshold is `70%`. |
| Human dogfood goal | CARRY-FORWARD | VAL-010 checklist generated, not yet executed. |

## Deterministic Checklist Execution

| QA ID | Result | Evidence |
|---|---|---|
| QA-01 | PASS | Task 1 report records exact gates and VAL-001 trigger/projection evidence; phase gate rerun passed. |
| QA-02 | PASS | Task 2 report records exact gates and VAL-002..VAL-005 domain/journal evidence. Due-date public read-back was intentionally closed in Task 4 per ESC-002. |
| QA-03 | PASS | Task 3 report records exact gates, build, CLI create smoke, invalid errors, one-event evidence, and rebuild proof. |
| QA-04 | PASS | Task 4 report records exact gates, build, due-date read-back through `internal/issues.Get` and `fid issue get --format json`, section tests, and adapter-boundary scan. |
| QA-05 | PASS | Task 5 report records exact gates, build, include/exclude filters, ordering, format parity, and unknown-format rejection. Its COV satisfaction statement is contradicted by phase-end coverage and is handled under QA-06. |
| QA-06 | FAIL | Full gate rerun completed; COV-001 failed with `internal/issues` at `79.5% < 80%`. |
| QA-07 | PASS | Final outcomes recorded for VAL-001..VAL-010 in this review. |

## ESC-002 Disposition

Resolved. `review-3.md` blocked on `due_date` being persisted but absent from
`internal/issues.Get` and `fid issue get --format json`. TL disposition kept
VAL-003 unchanged and assigned the closeout to Task 4. Task 4 evidence proves:

- `internal/issues.Get`: populated `due_date` reads back as `2026-06-30`; omitted
  due date reads back as `nil`.
- `fid issue get --format json`: created issue with `--due-date 2026-06-30`
  returns `"due_date":"2026-06-30"`; omitted due date returns `null`.

No ESC-002 carry-forward remains.

## Mid-Phase Review Disposition

| Review | Verdict | Blocking Fixed? | Improvement Disposition | ESC References |
|---|---|---|---|---|
| `review-3.md` | BLOCK | YES | No improvements raised. | ESC-002 resolved by Task 4 due-date read-back evidence. |

## Documentation Review

- Impacted docs updated: YES. Task 5 updated `project-docs/specs/06_DESIGN_PLAN_CLI.md` for list markdown output; all task completion reports are present.
- Current-context wrap-up still needed: phase wrap-up should update `phase_brief.md`, `progress_current.md`, and archive/close ESC-002 from `tl_escalations_current.md`.
- Placeholder scan: `rg -n "<TBD>|TODO|PLACEHOLDER|PASS /|FAIL /|phase-<N>|<YYYY" project-docs/docs/phases/phase-4/qa project-docs/docs/phases/phase-4/tasks` returned no matches.

## Security And Performance Review

- Security spot-check: CLI issue adapter contains no direct SQL assembly; `rg -n 'SELECT|QueryRow|Query\(' cmd/fid/issue.go` returned no matches. Phase 4 list ordering is covered by whitelist tests in VAL-009.
- Performance spot-check: benchmarks passed. Rebuild benchmarks recorded `442.4 ms/rebuild` for 1k events and `3812 ms/rebuild` for 10k events.

## Human Checklist Readiness

Generated at `project-docs/docs/phases/phase-4/qa/human-checklist.md`.

- Contains only `(human)` items: YES.
- Uses one observable per row: YES.
- Human-only `VAL-NNN` assertions mapped: YES, VAL-010 only.
- Product/perception/dogfooding risks to call out: command wording, output usefulness, and one invalid hierarchy/scheduling error.

## Release Maturity Status

| Maturity area | Status | Evidence or deferral |
|---|---|---|
| Baseline artifacts | BLOCKED | Final review and human checklist now exist, but deterministic coverage gate blocks phase approval. |
| Automated gates | BLOCKED | COV-001 fails for `internal/issues` at `79.5%`. |
| Held-out scenarios | DEFERRED | AGENTS.md records E2E/scenario validation as `N/A` until a harness is accepted. |
| Release readiness | DEFERRED | Not a release-readiness phase. |

## Findings

### Blocking

- `internal/issues` coverage is `79.5%`, below the COV-001 `80%` threshold. Evidence: phase-end `go test -race -coverprofile=coverage.out ./...` and `go tool cover -func=coverage.out`.

### Major

- None.

### Improvement

- Task 5 completion report states `internal/issues` `79.5%` satisfies the task COV references. That statement is inaccurate against COV-001 and should be corrected during the coverage closeout.

## Unresolved Risks And Carry-Forward Candidates

- Blocking: raise `internal/issues` coverage to at least `80%` and rerun the full phase gate.
- Human: execute VAL-010 checklist after the deterministic block is cleared.
