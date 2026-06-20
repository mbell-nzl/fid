# Phase 9 QA Plan - Concerns + v0.1c Closure

> Phase: 9 - Concerns + derived flag + v0.1c closure
> Refreshed: 2026-06-20
> QA Lead: Codex QA Lead
> Status: READY for Phase Runner before Task 1
>
> Location:
> `project-docs/docs/phases/phase-9/qa/qa-plan.md`.

## 1. Scope And Source Documents

Phase 9 ships the concern write/read surface on the existing concern schema and
derived issue read model:

- `fid concern raise`
- `fid concern address`
- `fid concern list`
- real-write proof for `concerns_raised`, `open_concerns_count`, and the
  `concerns_addressed` rule seam
- end-to-end v0.1c closure smoke over ACs, completions, and concerns

Source artifacts loaded and honored:

- `AGENTS.md`
- `project-docs/docs/current/phase_brief.md`
- `project-docs/docs/current/decisions_active.md`
- `project-docs/docs/current/progress_current.md`
- `project-docs/docs/current/tl_escalations_current.md`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/specs/09_AI_WORKFLOW.md`
- `project-docs/specs/01_PRD.md`
- `project-docs/specs/02_TRD.md`
- `project-docs/specs/03_TRD_SCHEMA.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md`
- `project-docs/docs/phases/phase-9/qa/validation-contract.md`
- `project-docs/docs/phases/phase-9/tasks/task-1.md`
- `project-docs/docs/phases/phase-9/tasks/task-2.md`
- `project-docs/docs/phases/phase-9/tasks/task-3.md`

Binding decisions and carry-forward constraints:

- DEC-A: prove `concerns_addressed` at the rule-evaluation seam; do not add a
  v0.1 default workflow edge.
- DEC-B: concern domain logic belongs in `internal/issues`; `fid concern` is a
  CLI namespace only.
- DEC-C: concern raise/address reject invalid input before any journal append.
- DEC-D: address is an in-place `open` -> `addressed` mutation only.
- DEC-E: concern writes do not populate `activity_log`; CF-011 remains deferred.
- DEC-002: TL-Lite autonomous resume is enabled; Phase Runner records RESUME and
  continues, while BLOCK still stops.
- Current open carry-forward items are not closed by Phase 9: CF-003, CF-008,
  CF-011, and CF-013. Phase 9 must not introduce role-matrix semantics, a
  resolution-reason model, activity timeline population, or unrelated
  `complete.go` orchestration work.

## 2. Requirement Coverage

| Requirement or decision | Phase 9 interpretation | Covered by |
|---|---|---|
| PRD-DEV-005 | A developer concern is structured and linked to both issue and completion. | VAL-001, VAL-002 |
| PRD-DEV-006 | Concern lifecycle supports the Phase 9 writer path `open` -> `addressed`; `escalated` and `wontfix` remain schema-valid but producer-deferred. | VAL-004, VAL-005, VAL-008 |
| PRD-DEV-007 | Open concerns surface as derived flags and counts on issue reads. | VAL-003, VAL-006, VAL-008 |
| PRD-WF-009 | `concerns_addressed` is supported and satisfied when no open concerns exist. | VAL-007 |
| TRD architecture `internal/issues` | Concern business rules stay in the domain package, not the CLI adapter. | QA-19 |
| TRD write-path seam | Validation happens before append/project for concern writes. | VAL-002, VAL-005 |
| TRD schema `concerns` | Concern row fields, statuses, and completion links are preserved. | VAL-001, VAL-004, VAL-009 |
| TRD `issues_enriched` | Open concern count and flag derive from projected DB state. | VAL-003, VAL-006, VAL-009 |
| TRD default workflow graph | No new default edge is added for `concerns_addressed`. | VAL-007, QA-19 |
| ADR-0001 | Journal remains canonical; rebuild reproduces DB projection. | VAL-009 |
| CLI design section 6 | `fid concern raise`, `address`, and `list` flags/output shapes are honored. | VAL-001, VAL-004, VAL-008, VAL-010 |
| Roadmap v0.1c exit | Full AC -> completion -> concern lifecycle works end to end. | VAL-010, VAL-011 |
| DEC-E / CF-011 | No Phase 9 `activity_log` assertion or implementation expectation. | QA-19 |

## 3. Validation Contract Refinement Check

- All in-scope PRD/TRD/design requirements referenced by at least one
  `VAL-NNN`: YES.
- Every `VAL-NNN` has one validator role: YES.
- Every `VAL-NNN` has a concrete evidence method: YES.
- Every `VAL-NNN` maps to a task/window or phase-end check: YES.
- Every `VAL-NNN` describes one observable, or human-checklist fan-out is
  explicitly required: YES. VAL-011 is fanned out below as HUMAN-011a..e.
- Conditional or binary `VAL-NNN` assertions require both qualifying and
  non-qualifying states: YES for VAL-002, VAL-003, VAL-005, VAL-006, VAL-007,
  VAL-008, VAL-009, and VAL-010.
- Parity or "matches reference" assertions have deterministic conformance proof
  where available: YES. VAL-009 requires pre-delete and post-rebuild structural
  comparison.
- Revisions requested from Tech Lead: None.

QA refinement notes:

- VAL-001 cites `fid concern list` read-back, but Task 1 may not ship the top
  level list command. Task 1 can prove read-back through
  `fid issue get --section concerns`; final QA re-proves VAL-001 through both
  `fid issue get --section concerns` and `fid concern list` after Task 3.
- VAL-011 remains human-only for product dogfood judgment. The deterministic
  QA step for VAL-011 is artifact enforcement: generate and verify the
  phase-end human checklist contains one row per observable and carries the
  final human result back into `final-review.md`.

## 4. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| Task 1 | `fid concern raise`, validate-before-append, derived flag lights | VAL-001, VAL-002, VAL-003 | High | First real concern writer; reject-path journal poisoning and derived read parity are the main risks. |
| Task 2 | `fid concern address`, derived flag clears, rule seam proof | VAL-004, VAL-005, VAL-006, VAL-007 | High | Must preserve in-place mutation, reject re-address, and keep DEC-A no-edge scope. |
| Task 3 | `fid concern list`, rebuild parity, v0.1c smoke | VAL-008, VAL-009, VAL-010 | High | Public read shape, rebuild parity, and full version closure all land here. |
| Phase end | Human dogfood checklist and final QA review | VAL-001..011 | High | QA Lead must execute deterministic checklist before any PASS outcome. |

## 5. Required Evidence Strategy

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | Task 1 completion report with exact required commands marked PASS or BLOCKED with reason; coverage percentages for `internal/issues` and `cmd/fid`. | VAL-001, VAL-002, VAL-003 | Developer, verified by Phase Runner and QA Reviewer | task-1 |
| EVID-02 | Task 1 tests and CLI output showing a valid raise, issue-section read-back, ID shape, status `open`, body path, role, issue link, and completion link. | VAL-001 | Developer, QA Lead | task-1 and final review |
| EVID-03 | Task 1 reject-path evidence for unknown role, missing issue, missing completion, and completion-on-other-issue: non-zero exit, clear stderr, and unchanged journal/row counts. | VAL-002 | Developer, QA Lead | task-1 and final review |
| EVID-04 | Task 1 before/after issue get and issue list JSON showing false/0 before raise and true/incremented after raise. | VAL-003 | Developer, QA Lead | task-1 and final review |
| EVID-05 | Task 2 completion report with exact required commands marked PASS or BLOCKED with reason; coverage percentages for `internal/issues` and `cmd/fid`. | VAL-004, VAL-005, VAL-006, VAL-007 | Developer, verified by Phase Runner and QA Reviewer | task-2 |
| EVID-06 | Task 2 tests and CLI output showing address mutates the same row to `addressed`, sets `addressed_in_completion_id`, updates timestamp, and retains the row. | VAL-004 | Developer, QA Lead | task-2 and final review |
| EVID-07 | Task 2 reject-path evidence for unknown role, nonexistent concern, already-addressed concern, and missing completion: non-zero exit, clear stderr, and unchanged journal/row counts. | VAL-005 | Developer, QA Lead | task-2 and final review |
| EVID-08 | Task 2 partial and final clear evidence: two open concerns -> address one leaves true/remaining count; address last clears false/0 on get and list. | VAL-006 | Developer, QA Lead | task-2 and final review |
| EVID-09 | Integration test output from the rule seam showing unsatisfied with open `CON-` IDs, then satisfied after all concerns are addressed. | VAL-007 | Developer, QA Lead | task-2 and final review |
| EVID-10 | Task 3 completion report with all required commands marked PASS or BLOCKED with reason, including integration and build; coverage percentages recorded. | VAL-008, VAL-009, VAL-010 | Developer, verified by Phase Runner and QA Reviewer | task-3 |
| EVID-11 | List command evidence across `json`, `toon`, `markdown`, and `yaml` for `--issue`, `--status`, combined filters, and empty match exit 0 with an empty collection. | VAL-008 | Developer, QA Lead | task-3 and final review |
| EVID-12 | Rebuild parity proof: snapshot concern rows and derived issue flags, delete SQLite DB, run `fid rebuild --by developer`, compare post-rebuild state to snapshot. | VAL-009 | Developer, QA Lead | task-3 and final review |
| EVID-13 | v0.1c smoke test output over `ac add` -> `ac verify` -> `ac list` -> `issue complete` -> `concern raise` -> flag lights -> `concern address` -> flag clears -> `concern list`. | VAL-010 | Developer, QA Lead | task-3 and final review |
| EVID-14 | Mid-phase review artifact `review-1.md` reviewing Tasks 1..3 and their reports, gates, commits, and in-window VAL evidence. | VAL-001..010 | QA Reviewer | after task-3 |
| EVID-15 | Final review with deterministic command output excerpts and one final outcome for every VAL-001..011. | VAL-001..011 | QA Lead | phase-end |
| EVID-16 | Human checklist with one human-only row per VAL-011 observable and completed human outcome notes. | VAL-011 | Human, prepared by QA Lead | phase-end |
| EVID-17 | Architecture and decision conformance proof: no new workflow edge for `concerns_addressed`, concern domain code remains under `internal/issues`, no concern `activity_log` assertion or implementation expectation. | DEC-A, DEC-B, DEC-E, CF-011 | QA Lead | final review |

## 6. Deterministic Gates

No required command is a placeholder. All task and phase commands are concrete.
Missing, silently skipped, or narrower substituted required commands are failed
gates unless a blocker report records `BLOCKED` with a reason.

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | No files listed | Every task and phase final |
| Lint/static analysis | `golangci-lint run ./...` | PASS | Every task and phase final |
| Task coverage gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Every task |
| Coverage report | `go tool cover -func=coverage.out` | PASS and threshold evidence for DEC-007 | Every task and phase final |
| Task 1 build | `go build -o /tmp/fid ./cmd/fid` | PASS | Required for Task 1 because it adds a command surface |
| Task 1 integration | `go test -tags integration -race ./...` | PASS or documented N/A | Required only if Task 1 adds an integration-tagged test |
| Task 2 build | `go build -o /tmp/fid ./cmd/fid` | PASS | Required for Task 2 because it changes a command surface |
| Task 2 integration | `go test -tags integration -race ./...` | PASS or documented N/A | Required only if Task 2 adds an integration-tagged test |
| Task 3 integration | `go test -tags integration -race ./...` | PASS | Required by Task 3 |
| Task 3 build | `go build -o /tmp/fid ./cmd/fid` | PASS | Required by Task 3 |
| Task 3 benchmark | `go test -bench . -benchmem ./...` | PASS or documented N/A | Required only if Task 3 adds a hot query path; otherwise N/A |
| Phase unit suite | `go test -race ./...` | PASS | Phase final |
| Phase coverage gate | `go test -race -coverprofile=coverage.out ./...` | PASS | Phase final |
| Phase integration suite | `go test -tags integration -race ./...` | PASS | Phase final |
| Phase benchmark suite | `go test -bench . -benchmem ./...` | PASS | Phase final |
| Phase build | `go build -o /tmp/fid ./cmd/fid` | PASS | Phase final |

Coverage thresholds to record from `go tool cover -func=coverage.out`:

- COV-001: `internal/issues`, `internal/journal`, `internal/policy`, and
  `internal/workflow` each remain >= 80%.
- COV-002: `cmd/fid` remains >= 70%.

Security, dependency audit, license check, and secret scan remain release
readiness controls and are not Phase 9 gates unless an implementation task adds
a dependency or security-sensitive surface. New dependencies are not expected.

## 7. Command And Feasibility Validation

| Referenced Step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | Help output and path existence. | PASS | `gofmt -h` lists `-l`; `ls cmd/fid internal/issues` succeeded. |
| `golangci-lint run ./...` | Installed binary/version check. | PASS | `golangci-lint --version` returned `2.12.2`, matching DEC-003 v2 policy. |
| `go test -race -coverprofile=coverage.out ./...` | Go tool and flag help. | PASS | `go version` returned `go1.26.2 linux/amd64`; `go help testflag` lists `-race` behavior through build flags and `-coverprofile`. |
| `go tool cover -func=coverage.out` | Tool help output. | PASS | `go tool cover -h` lists `-func string` for stdout coverage percentages. |
| `go test -tags integration -race ./...` | Go test/build flag help. | PASS | `go help build` lists `-tags`; `go help testflag` confirms test command flag handling. |
| `go test -bench . -benchmem ./...` | Go test flag help. | PASS | `go help testflag` lists `-bench regexp` and `-benchmem`. |
| `go build -o /tmp/fid ./cmd/fid` | Go build help and target existence. | PASS | `go help build` lists `-o output`; `ls cmd/fid` shows the main command package. |
| Phase 9 task specs | File-state check. | PASS | `ls project-docs/docs/phases/phase-9/tasks` shows `task-1.md`, `task-2.md`, and `task-3.md`. |
| Phase validation contract | File-state check. | PASS | `ls project-docs/docs/phases/phase-9/qa` shows `validation-contract.md`. |
| QA templates | File-state check. | PASS | `ls project-docs/templates` shows `qa-plan-template.md`, `mid-phase-review-template.md`, `final-review-template.md`, and `human-checklist-template.md`. |
| Existing v0.1c e2e precedent | File-state and source search. | PASS | `ls cmd/fid` shows `e2e_v01c_ac_test.go` and `e2e_v01c_complete_test.go`; Task 3 adds the concern closure smoke beside them. |
| Existing rebuild precedent | File-state and source search. | PASS | `ls cmd/fid internal/journal` shows `rebuild.go`, `rebuild_test.go`, and rebuild parity tests. |
| Existing concern schema/read-path assumption | Source search. | PASS | `rg` found `IssueConcern`, `issueConcerns`, `concerns_raised`, `open_concerns_count`, `loadConcernStatus`, and `RuleConcernsAddressed` in existing code. |
| Future `fid concern raise/address/list` commands | Spec and task validation. | PASS after Task 3, BLOCKED before implementation if invoked. | CLI syntax is specified in `06_DESIGN_PLAN_CLI.md` section 6 and in task specs; Phase Runner must not expect these commands before their owning tasks complete. |
| Task completion reports `task-1-done.md`..`task-3-done.md` | Workflow artifact expectation. | PASS as future required artifacts. | `09_AI_WORKFLOW.md` requires each Developer to write a report; Phase Runner verifies each report after the task. |
| Human checklist feasibility | Spec and task validation. | PASS after Task 3. | VAL-011 dogfood uses the same commands implemented by Tasks 1..3 on a real `.fid/` workspace. QA Lead generates the checklist at phase final. |

## 8. Deterministic Checklist

Every item here is owned by QA Lead and must end with `(qa-lead)`. Final-review
PASS requires literal command output, file-state proof, logs, or an equivalent
reproducible excerpt.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Verify Task 1 completion report records exact required gates and DEC-007 coverage percentages; missing or narrowed commands are Blocking `(qa-lead)` | VAL-001, VAL-002, VAL-003 | File-state check of `project-docs/docs/phases/phase-9/tasks/task-1-done.md`; compare its `Tests run` and `Exact quality-gate results` against Task 1 commands. |
| QA-02 | Run the final phase gate `gofmt -l cmd internal` and require empty output `(qa-lead)` | no VAL | Command output from `gofmt -l cmd internal`. |
| QA-03 | Run the final phase gate `golangci-lint run ./...` and require PASS `(qa-lead)` | no VAL | Command output from `golangci-lint run ./...`. |
| QA-04 | Run `go test -race ./...` at phase final and require PASS `(qa-lead)` | VAL-001..010 | Command output from `go test -race ./...`. |
| QA-05 | Run `go test -race -coverprofile=coverage.out ./...` followed by `go tool cover -func=coverage.out`; verify COV-001/COV-002 thresholds `(qa-lead)` | VAL-001..010 | Command output plus coverage rows for `internal/issues`, `internal/journal`, `internal/policy`, `internal/workflow`, and `cmd/fid`. |
| QA-06 | Run `go test -tags integration -race ./...` at phase final and require PASS `(qa-lead)` | VAL-007, VAL-009, VAL-010 | Command output from the integration suite. |
| QA-07 | Run `go test -bench . -benchmem ./...` at phase final and require PASS `(qa-lead)` | no VAL | Command output from the benchmark suite. |
| QA-08 | Run `go build -o /tmp/fid ./cmd/fid` at phase final and require PASS `(qa-lead)` | VAL-001, VAL-004, VAL-008, VAL-010 | Command output and `/tmp/fid` build artifact existence. |
| QA-09 | In a temp workspace with `/tmp/fid`, execute a valid concern raise against a real completion, then read it through `fid issue get --section concerns --format json` and `fid concern list --format json`; assert `CON-` ID shape, issue/completion/role/summary/body/status fields `(qa-lead)` | VAL-001 | CLI transcript with commands: `fid init`, issue/AC/complete setup, `fid concern raise --on <issue> --by developer --in-completion <TC> --summary <text> --body-md <path>`, `fid issue get <issue> --section concerns --format json`, `fid concern list --issue <issue> --format json`. |
| QA-10 | In the same or a fresh temp workspace, execute each invalid `fid concern raise` case and prove zero append/row change before a valid append adds exactly one journal line and one concern row `(qa-lead)` | VAL-002 | CLI stderr/non-zero evidence for unknown role, nonexistent issue, nonexistent completion, completion-not-on-issue; file-state proof from journal line counts and concern-list/DB row counts before and after. |
| QA-11 | Capture issue get/list JSON before and after a real raise; require false/0 before and true/incremented after on both surfaces `(qa-lead)` | VAL-003 | CLI transcript from `fid issue get <issue> --format json` and `fid issue list --format json` before and after `fid concern raise`. |
| QA-12 | Execute a valid address on an open concern and prove the same row is retained with `status="addressed"` and `addressed_in_completion_id` set `(qa-lead)` | VAL-004 | CLI transcript from `fid concern address <CON> --in-completion <TC> --by developer`, `fid issue get <issue> --section concerns --format json`, and `fid concern list --issue <issue> --format json`; row count unchanged. |
| QA-13 | Execute each invalid `fid concern address` case and prove zero append/row change; verify valid address applies exactly one status mutation `(qa-lead)` | VAL-005 | CLI stderr/non-zero evidence for unknown role, nonexistent concern, already-addressed concern, nonexistent completion; file-state proof from journal line counts and concern-list/DB row counts. |
| QA-14 | Seed two open concerns, address one, then address the last; require partial true/remaining count and final false/0 on both issue get and list `(qa-lead)` | VAL-006 | CLI transcript from issue get/list JSON after each address. |
| QA-15 | Verify the `concerns_addressed` rule seam with real concern data: open concern returns unsatisfied naming `CON-` IDs, all addressed returns satisfied `(qa-lead)` | VAL-007 | Integration test output from `go test -tags integration -race ./...` or targeted `go test -race ./internal/issues ./internal/workflow -run ConcernsAddressed -v`; final review must cite the actual test name/output. |
| QA-16 | Verify `fid concern list` filters by `--issue`, `--status`, combined issue+status, all formats, and empty match exit 0 with empty collection `(qa-lead)` | VAL-008 | CLI transcript from `/tmp/fid concern list --format json`, `toon`, `markdown`, `yaml`, with seeded concerns across two issues and statuses. |
| QA-17 | Verify delete-DB plus `fid rebuild --by developer` reproduces concern rows and derived flags exactly `(qa-lead)` | VAL-009 | Pre-delete snapshots from `fid concern list --format json` and `fid issue get/list --format json`; command output from `fid rebuild --by developer`; post-rebuild snapshots showing structural equality. |
| QA-18 | Verify v0.1c closure smoke covers `ac add`, `ac verify`, `ac list`, `issue complete`, `concern raise`, flag lights, `concern address`, flag clears, and `concern list` `(qa-lead)` | VAL-010 | Output from `go test -race ./cmd/fid -run V01C -v` or a CLI transcript over `/tmp/fid`; final review must cite every step's observable result. |
| QA-19 | Verify architecture/decision constraints: no new default workflow edge for `concerns_addressed`, concern domain files stay in `internal/issues`, CLI wiring stays in `cmd/fid`, and no `activity_log` assertion is introduced `(qa-lead)` | DEC-A, DEC-B, DEC-E, CF-011 | File-state/source-search proof using `rg` over `internal/workflow`, `internal/issues`, `cmd/fid`, and task reports; cite matching lines or absence where relevant. |
| QA-20 | Verify Task 2 completion report records exact required gates and DEC-007 coverage percentages; missing or narrowed commands are Blocking `(qa-lead)` | VAL-004, VAL-005, VAL-006, VAL-007 | File-state check of `project-docs/docs/phases/phase-9/tasks/task-2-done.md`; compare report commands against Task 2 commands. |
| QA-21 | Verify Task 3 completion report records exact required gates, including integration and build, plus coverage percentages; missing or narrowed commands are Blocking `(qa-lead)` | VAL-008, VAL-009, VAL-010 | File-state check of `project-docs/docs/phases/phase-9/tasks/task-3-done.md`; compare report commands against Task 3 commands. |
| QA-22 | Verify `project-docs/docs/phases/phase-9/qa/human-checklist.md` is generated from the human subset and contains HUMAN-011a..e only for VAL-011 observables `(qa-lead)` | VAL-011 | File-state check of human checklist artifact and final-review outcome row for VAL-011. |
| QA-23 | Record a final outcome for every VAL-001..011 in `final-review.md`; PASS without cited evidence is invalid `(qa-lead)` | VAL-001..011 | File-state check of `final-review.md` validation outcome table with PASS/FAIL/BLOCKED/CARRY-FORWARD and evidence links/excerpts. |

## 9. Mid-Phase Review Cadence

Phase 9 has exactly three implementation tasks, so the routine review cadence is
one three-task window:

- `review-1.md` after Task 3 completion, before QA Lead final review.
- Review path:
  `project-docs/docs/phases/phase-9/qa/review-1.md`.
- Review template:
  `project-docs/templates/mid-phase-review-template.md`.
- Reviewer role: QA Reviewer, distinct from QA Lead.
- Review scope: Task 1, Task 2, Task 3 specs, commits, completion reports,
  tests, docs, exact command evidence, and VAL-001..010 progress.
- Phase Runner also routes any Developer blocker, soft blocker, or QA Reviewer
  Blocking finding through TL-Lite according to DEC-002 and `09_AI_WORKFLOW.md`.

The QA Reviewer must explicitly check:

- Each task report contains all required headings.
- Each task report includes commit hash and subject copied from git output.
- Each reported command matches the task spec and this QA plan.
- No required command is missing, silently skipped, or substituted with a
  narrower command.
- VAL evidence is sufficient for the task window or explicitly scheduled for
  phase final.
- DEC-A..E and the current carry-forward register are not contradicted.

## 10. Human Verification Scope

Items here require human dogfooding or product judgment and must end with
`(human)`. QA Lead copies these rows into
`project-docs/docs/phases/phase-9/qa/human-checklist.md` at phase final using
the human checklist template. VAL-011 PASS requires completed human evidence.

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-011a | Raise a concern against a real completion and confirm the visible open concern has the expected summary, role, issue link, and completion link `(human)` | VAL-011 | Human notes from a real `.fid/` workspace, including the command used and whether the visible concern matched expectation. |
| HUMAN-011b | After raising the concern, confirm issue reads make the open concern obvious enough through `concerns_raised` and `open_concerns_count` `(human)` | VAL-011 | Human observation of `fid issue get` and `fid issue list` output before/after raise. |
| HUMAN-011c | Address the concern from a later completion and confirm the lifecycle reads as an in-place resolution, not a new unrelated concern `(human)` | VAL-011 | Human observation of address command output and subsequent concern read-back. |
| HUMAN-011d | After addressing the last open concern, confirm issue reads no longer imply unresolved concern work remains `(human)` | VAL-011 | Human observation of `fid issue get` and `fid issue list` after address. |
| HUMAN-011e | Use `fid concern list` with issue and status filters and confirm the filtered views match the way a developer would look up concern state during dogfooding `(human)` | VAL-011 | Human notes on filter usefulness, empty-result behavior, and any confusing output. |

## 11. Phase-End Final Review Expectations

QA Lead final review creates:

- `project-docs/docs/phases/phase-9/qa/final-review.md`
- `project-docs/docs/phases/phase-9/qa/human-checklist.md`

The final review must include:

- Findings ordered by severity.
- Gate results with command output excerpts.
- Deterministic checklist outcomes QA-01..QA-23.
- One final outcome for every VAL-001..011: PASS, FAIL, BLOCKED, or
  CARRY-FORWARD.
- Requirement compliance against PRD-DEV-005..007, PRD-WF-009, TRD schema/read
  path/write path/default workflow graph, CLI design section 6, ADR-0001, and
  roadmap v0.1c closure.
- Test quality review, including red/green evidence from completion reports,
  reject-path tests, rebuild parity tests, and e2e smoke coverage.
- Documentation review for `06_DESIGN_PLAN_CLI.md` and task completion reports.
- Risk review for journal poisoning, rebuild divergence, public CLI/output
  shape drift, workflow-edge scope creep, and deferred CF-011 activity timeline.

The phase can proceed to final review only when:

- All three task completion reports exist or a blocker report stops the phase.
- Required task commands are PASS or explicitly BLOCKED with reason.
- `review-1.md` exists with `VERDICT: CONTINUE`, or a BLOCK has been triaged and
  resolved according to `09_AI_WORKFLOW.md`.
- No open ESC blocks execution.

The phase can receive a deterministic PASS only when:

- Every QA Lead deterministic item has reproducible evidence.
- Human checks are separated into `human-checklist.md` and VAL-011 is not marked
  PASS before human evidence exists.
- No unresolved blocking findings remain.
- Required documentation updates are named and verified.
