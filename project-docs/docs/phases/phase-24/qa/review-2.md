# Phase 24 - Mid-Phase Review 2 (Tasks 1-3)

Reviewer: qa-reviewer
Date: 2026-10-05
Tasks covered: phase-24/task-1, phase-24/task-2, phase-24/task-3
Review type: FRESH INDEPENDENT RE-REVIEW
Prior review: `project-docs/docs/phases/phase-24/qa/review-1.md` (F1-F5 and Repair Plan)

## Verdict

PASS. All five Must-fix findings in review-1's Repair Plan have fresh, committed evidence and passed affected independent reruns. Window A may advance.

## Scope Reviewed

- Task specs and committed completion reports: `project-docs/docs/phases/phase-24/tasks/task-{1,2,3}{,-done}.md`; Phase 24 validation contract, QA plan (Window A and section 12.7), current context, run ledger, `AGENTS.md`, `08_QUALITY_PLAN.md`, and `09_AI_WORKFLOW.md` sections 2.7 and 8.
- Repair commits: task 1 `64a8228295095a686a171d081a8ab4a8b9ff2052`, task 2 `ed1ce980e52ba4eeb6af00b58dbd6670c3b70731`, task 3 `545803dc2185f565d2b02258bd18ee0fd1d7d2ed`; report commit chain ends at `67615ef90bf94a5e6eaa65f983eef4e42a524270`, `74c8f60ee04ba5055878b57c3b3c53ef7008c4a4`, and `34e550ff3bc60530fbc1e8a7b8061eb6cc1ed0ac`. Each task report's stated implementation hash passed a separate `git rev-parse --verify '<hash>^{commit}'` (observed `exit_code: 0`).
- `git diff --name-only a5052db..HEAD` contains four test files and the three task reports only; `git diff --check a5052db..HEAD` terminated `exit_code: 0`. No accepted product, architecture, data-model, public-interface, dependency, validation, or gate definition changed during QA repair. Held-out material was not opened.

## Summary

- Review outcome: PASS
- Must-fix findings: 0
- Improvement findings: 1
- Blocker candidates: 0
- Existing escalation references: none
- Covered validation assertions: VAL-001-016, VAL-052, VAL-058-060, and task 2's seven rows of VAL-066

## Validation Contract Coverage

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-003 | task-1 | ON TRACK | Valid receipt rebuild and malformed insert-map/receipt-guard CLI cases at `internal/export/export_test.go:125`; projection snapshot check at `internal/journal/current_event_validation_test.go:14`. F1 resolved below. |
| VAL-004 | task-1 | ON TRACK | Schema-19 validate/import control and exact schema-18 `manifest.json` refusal with unchanged target at `internal/export/export_test.go:75`. F2 resolved below. |
| VAL-066 (task-2 leaves) | task-2 | ON TRACK | Seven unmerged-state pairs and five project-file Source deletion-state pairs at `cmd/fid/integration/p24_idempotency_test.go:310`; observable new-key preflight remains. F3 resolved below. Later observing leaves remain assigned to tasks 4-8. |
| VAL-015 | task-3 | ON TRACK | Three adopted leaves times six ordered built-process conflict cases, plus three standalone success controls at `cmd/fid/input_descriptor_test.go:44`. F4 resolved below. Later leaves remain assigned to tasks 4-9. |
| VAL-016 | task-3 | ON TRACK | Forty-two exact structural failures and three well-typed controls at `cmd/fid/input_descriptor_test.go:494`. F5 resolved below. Later leaves remain assigned to tasks 4-9. |
| VAL-001, 002, 005-014, 052, 058-060 | tasks 1-3 as mapped in review-1 | ON TRACK | Review-1 records their in-window tests, benchmark, inventory and timing evidence. The QA repair diff changed no behavior, contract, or test outside the five repaired evidence areas; the committed reports preserve final exact gates. Phase-end samples and later-task leaf completion remain due where the contract schedules them. |

## Gate And Evidence Checks

Working directory for all commands below: `/home/mark/projects/fid`. The task specs require seven, eight, and five exact commands for tasks 1, 2, and 3. Their final `## Exact quality-gate results` each record all commands byte-identically, with `PASS`, a command-output token, and observed `exit_code: 0`; `## Tests run` preserves focused failed attempts and later green results. QA-plan section 5 requires these same per-task gates; Window A adds no other task command. The task-1 and task-2 reports retain their benchmark and registry searches; task 3 has no conditional command. The reports preserve prior ordinary/tagged failures and byte-identical later passing gates; no report's earlier `BLOCKED` label is used as a protected blocker. The latest tagged integration results are 293.970 s, 288.286 s, and 298.059 s (each below 900 s); the reported ordinary package times remain below 600 s. This review did not rerun the full phase-level gate.

| Independent check | Reproduced result |
|---|---|
| `gofmt -l cmd internal` | `exit 0 (no output)`; terminated `exit_code: 0`. This is a byte-identical required task command. |
| `go test -count=1 -run '^(TestCurrentIdempotencyFieldDigestsReplayFailClosed|TestGateReceiptFieldDigestGuardAndMalformedInsertFailClosedThroughCLI|TestSchema19ArchiveImportsAndSchema18ManifestFailsClosed)$' ./internal/journal ./internal/export` | `ok  	github.com/mbell-nzl/fid/internal/journal	0.394s` and final `ok  	github.com/mbell-nzl/fid/internal/export	2.593s`; terminated `exit_code: 0`. F1/F2. |
| `go test -tags integration -count=1 -run '^TestP24ObservingJudgmentReplayNeverObserves$' ./cmd/fid/integration/` | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	11.127s`; terminated `exit_code: 0`. F3. |
| `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAdopted(ExampleProcessExclusivity|StructuralProcessDetails)$'` | `ok  	github.com/mbell-nzl/fid/cmd/fid	3.922s`; terminated `exit_code: 0`. F4/F5. |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-1-done.md project-docs/docs/phases/phase-24/tasks/task-1.md` | `PASS report-lint completed with no failures`; terminated `exit_code: 0`. |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-2-done.md project-docs/docs/phases/phase-24/tasks/task-2.md` | `PASS report-lint completed with no failures`; terminated `exit_code: 0`. |
| `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-3-done.md project-docs/docs/phases/phase-24/tasks/task-3.md` | `PASS report-lint completed with no failures`; terminated `exit_code: 0`. |

Focused review commands supplement the reports' required gates; they do not replace them.

## Per-Task Findings

### Task 1 - Idempotency field-digest storage

#### Must-fix

- None.

#### Improvement

- None. The new real Gate journal fixture exercises a valid receipt update, three malformed insert maps, and a mismatched stored receipt guard. Each invalid process run asserts `journal-replay-stopped`, integrity category, exit 6, exact journal file/line, and unchanged workspace digest. The archive test pins `schema-unsupported` to `manifest.json` while retaining the schema-19 control and unchanged-target check.

### Task 2 - Input-as-written digests and named conflict fields

#### Must-fix

- None.

#### Improvement

- None. The five source-bearing leaf names are explicitly checked against fixture Source presence. Each deleted-source same-key call compares the original response bytes and Event count; a distinct new-key request refuses with typed JSON and no Event in the same unmerged/deleted state. The preflight new-key request succeeds before the state becomes unobservable. Isolated missing-file observation semantics remain assigned to task 4.

### Task 3 - Input descriptor framework, help, templates and field validation

#### Must-fix

- None.

#### Improvement

- None. The process exclusivity corpus checks exact first Human conflict for each ordered combination, the JSON `--json` envelope, exit 2, empty output/workspace effects and three standalone `--example` successes outside a project. The shared writer asserts exact operation and typed details for each conflict. The structural corpus checks all applicable top-level classes, both nested subject depths, arrays and variants with exact path/value/allowed, exit 2, no Event, empty workspace, and well-typed controls reaching domain validation.

## Prior Finding Resolution

| Prior finding and review-1 repair obligation | Owning task | Result | Fresh evidence |
|---|---|---|---|
| F1 / VAL-003: valid receipt update; mismatched guard and three malformed insert maps stop with code/exit/line, no mutation | task-1 | RESOLVED | `64a8228` changes `internal/export/export_test.go:125` and `internal/journal/current_event_validation_test.go:14`; named rerun passed for journal and export. Valid fixture invokes a real Gate Receipt; rejected cases compare workspace digest or projection snapshot. Seven exact task gates and report-lint passed in the report. |
| F2 / VAL-004: schema-18 archive check names `manifest.json`, with codes and unchanged target | task-1 | RESOLVED | `internal/export/export_test.go:104-120` retains `assertArchiveCheck`, asserts the `schema-unsupported` check's `Entry == manifest.json`, and compares target digest. Named rerun passed. |
| F3 / VAL-066: five deleted-source replay/control pairs plus observable preflight | task-2 | RESOLVED | `ed1ce98` adds explicit five-leaf Source enumeration and deleted-state new-key calls at `cmd/fid/integration/p24_idempotency_test.go:310-394`. Named integration rerun passed; eight exact task gates and report-lint passed in the report. The report retains the first F3 focused `FAIL` (`exit_code: 1`) and later byte-identical `ok` (`exit_code: 0`). |
| F4 / VAL-015: 18 ordered conflicts with exact details, exit/no effects and standalone success | task-3 | RESOLVED | `545803d` adds `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44`; named race rerun passed. Three leaves × six cases and three standalone controls are explicit. Report retains the first corpus `FAIL` (`exit_code: 1`) and later `ok` (`exit_code: 0`). |
| F5 / VAL-016: applicable structural classes with exact details, exit/no Event and qualifying controls | task-3 | RESOLVED | Same commit adds `TestAdoptedStructuralProcessDetails` at `cmd/fid/input_descriptor_test.go:494`; named race rerun passed. Forty-two failure rows and three valid-structure controls are explicit. Five exact task gates and report-lint passed in the report. |

## Cross-cutting observations

- **Improvement; no VAL** — `project-docs/docs/phases/phase-24/run-ledger.md:384` cites task-3 repair implementation hash `545803d9f55b7a1e957a4d27b0a9a7703b6e5c2b`, which failed `git rev-parse --verify '<hash>^{commit}'` (`fatal: Needed a single revision`, terminated `exit_code: 128`). The actual committed test change is `545803dc2185f565d2b02258bd18ee0fd1d7d2ed` (`git rev-parse --verify` terminated `exit_code: 0`). Fix direction: the Phase Runner should correct its ledger checkpoint. The committed task-3 report cites a valid hash and passed report-lint; this bookkeeping typo does not change the five finding resolutions.

## Blocker Candidates

- None.

## Existing Escalation References

- None.

## Final Lines

```text
Review written: project-docs/docs/phases/phase-24/qa/review-2.md
Must-fix: 0  Improvement: 1  Blocker candidates: 0
VERDICT: PASS
```
