# Phase 24 - Task 14: Human-takes-control journey and phase conformance - Completion Report

Refs: phase-24/task-14

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

Same-task continuation attempt 2, governed by `TRG-24-016 / DEC-290`.
The phase checkpoint was `314b2d038fc4bc13c801569153fba8d665e565a1`;
the prior implementation was `4b6a3837e52f24a5a8975a1d2abef94635c74823`
and its report commit was `98d45617b35c212683b4346d30a760eb8ed630d3`.
The latest TASK-RECOVERY tuple was checked as exact copies: trigger
`task-14:public-interface-decision:terminal-state-provenance-qa-pass`,
attempt count `2`, category `in-scope-implementation`, repair mode
`task-scope`, repair authority `TRG-24-016 / DEC-290`.
The original attempt-1 BLOCKED report and all its failed evidence are
preserved verbatim below as a quoted historical artifact.

## Summary of what was implemented

- Delivered DEC-290 clause 3: loaded registry edge enumeration includes implicit cancellation, returns non-nil empty lists, and sorts states bytewise with their categories. Resume JSON adds only `issue.next_states`; its Human Issue block adds Revision and Next states immediately after Status. Only the two authorized transition flag usages point to `fid resume <issue-id>`.
- Completed the retained real-process Human journey: protected refusal without an Event, caller-authored Decision Source, approved Decision Reference and Human Outcome, resolved Stops, qualifying Evidence, terminal-success transition, and successful `envelope.complete`. Every command and input value records one of DEC-278's five provenance classes and its source line. The final transition takes its state and Issue revision from the Human resume read after the last mutation; actor attribution stays setup.
- Proved live/rebuilt/imported semantic parity, including field-digest maps, observed repositories, three new detail reads, projection integrity, archive validation, and manifest/import projection schema 19.
- Added the recursive forbidden-member descriptor scan, retained Stop-rendering and absence proofs, and extended the retained generator to 59 manifest rows covering judgment views and three valid authoring workspaces alongside all prior denials and repository views.
- Added the retained Human-takes-control scenario to the quality plan. No clause 6 convergence was needed; retained assertions keep their names, subjects, and strength.

## Files created or modified

- `internal/workflow/registry.go`
- `internal/workflow/resume.go`
- `cmd/fid/resume_human.go`
- `cmd/fid/issue_transition.go`
- `internal/workflow/registry_next_states_test.go`
- `cmd/fid/p24_resume_issue_test.go`
- `cmd/fid/integration/p24_take_control_test.go`
- `cmd/fid/integration/p24_take_control_parity_test.go`
- `cmd/fid/integration/p24_human_fixtures_test.go`
- `cmd/fid/input_descriptor_test.go`
- `project-docs/specs/08_QUALITY_PLAN.md`
- `project-docs/docs/phases/phase-24/tasks/task-14-done.md`

## Documentation updated

- `project-docs/specs/08_QUALITY_PLAN.md`: critical scenario inventory now cites the Human-only journey and its rebuild/export/import parity, with state and revision sourced from the final Human resume.
- This completion report records current deterministic evidence and preserves the original BLOCKED report. The task spec, decisions, public-interface contract, validation contract, QA artifacts, workflow controls, and runner-owned ledger were not edited.

## Validation Contract Coverage

- VAL-049: PASS — `TestP24HumanTakeControlJourney`, `cmd/fid/integration/p24_take_control_test.go:400`, calls `p24CompleteTakeControl`: exit-5 protected refusal leaves the journal unchanged; Human help/example discovery leads through Decision, Outcome, Evidence and transition to exact `completed envelope <id>` output. Its source recorder rejects unobserved or multiline fid sources and restricts setup to actor IDs/project directory. At line 381 the two transition preconditions are logged as `fid-output` from the final `Next states:` and `Revision:` lines. `TestP24JourneyProvenanceKeepsExactSourceLine` at line 158 covers a revision token also present in an unrelated Issue ID. The retained `TestP24TakeControlTerminalStateProvenance` at line 71 and forbidden-source guard remain and pass; JSON, SQL and journal never supply acting-caller values.
- VAL-050: PASS — `TestP24TakeControlRebuildExportImportParity`, `cmd/fid/integration/p24_take_control_parity_test.go:12`, completes that same Human journey, then compares all portable projection rows, field-digest maps, and Concern/Deviation/Human Outcome detail-read results against both rebuilt and imported projections. The portable row equality includes observed repository records; projection-integrity assertions run for both. Export validation passes and both archive manifest and import result explicitly assert `projection_schema == 19`.
- VAL-051: PASS — `TestP24DescriptorForbiddenMembers`, `cmd/fid/input_descriptor_test.go:21`, recursively checks every registry descriptor's children, element and variants for the five forbidden member names, including `repository`. `TestP24StopHumanRealCorpus`, `cmd/fid/integration/p24_human_fixtures_test.go:90`, retains the real Stop no-raw-JSON assertions. The disclaimer absence command below terminates with declared `exit_code: 1` and no matches.
- VAL-063: PASS — the exact `rg -n "PRIVATE_" cmd/fid/integration` absence command below terminates with declared `exit_code: 1`, no output.
- VAL-052: PASS — the final-tree ordinary and tagged gates report fresh, uncached elapsed values in the timing table below. Every ordinary package is below 600 s and tagged integration below 900 s. The earlier, superseded ordinary attempt's `cmd/fid` time of 530.567 s is preserved and flagged as headroom, not copied into the final timing evidence.
- VAL-068: PASS — `TestRegistryNextStates`, `internal/workflow/registry_next_states_test.go:9`, compares the several-edge `qa_review` result exactly against sorted cancelled/cancelled, qa_block/started and qa_pass/completed; it checks edge membership through `RulesForEdge`, no-edge states, unknown workflow, non-nil empty slices and JSON `[]`. `TestP24ResumeIssuePosition` and `TestP24ResumeIssueNextStates`, `cmd/fid/p24_resume_issue_test.go:10` and `:20`, check line placement plus full, empty and nil Human lists. `TestP24IssueTransitionPositionHelp` at `:42` checks both exact usages. Built-binary `TestP24TakeControlResumePair`, `cmd/fid/integration/p24_take_control_test.go:116`, checks JSON/Human revision 4 and all three categorized edges, exact Issue/edge JSON key sets, then revision 5 with JSON `[]` and Human `none.` after cancellation. Independent real CLI reads below confirm the changed Human entrypoint and help. Pairing/parity JSON is never a journey provenance source.
- VAL-053 preparation: complete — `TestP24HumanFixtures`, `cmd/fid/integration/p24_human_fixtures_test.go:41`, requires and emits the stuck `take-control-start` row and verifies its refusal. The phase-end Human takes-control judgment remains QA/Human-owned; no Human acceptance is claimed here.
- VAL-054 and VAL-055 preparation: complete — the same generator retains all five denial rows and four clean/dirty, same/different repository-view rows and their original assertions. Human judgment remains phase-end work.
- VAL-056 preparation: complete — the same generator adds every judgment mutation/list pairing, the Concern, Deviation and Human Outcome detail reads, and the judgment resume row. All appear in the 59-row manifest with their Human/JSON pair evidence. Human readability judgment remains phase-end work.
- VAL-057 preparation: complete — the same generator emits and rehearses valid `run close` with `expected_repository`, `decision record-reference` with digest omitted, and a third-leaf `concern record` workspace. Human first submissions run against retained pre-action workspace state before JSON replay; restoring pre-state permits independent repetition without weakening retained assertions. Human authoring judgment remains phase-end work.

## Tests run

All commands ran from `/home/mark/projects/fid`, except the independent
real-entrypoint reads whose working directory is the generated take-control
workspace named below. Focused commands are supporting evidence, not
substitutes for the seven required gates.

- `go test -run '^TestP24ResumeIssuePosition$' -count=1 ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`). Expected RED before product changes: missing Revision and Next states; package 0.181 s.
- `go test -run '^TestRegistryNextStates$' -count=1 ./internal/workflow` - FAIL - `FAIL` (terminated `exit_code: 1`). Expected RED before product changes: undefined NextState/NextStates.
- `go test -run '^(TestRegistryNextStates|TestP24ResumeIssue.*|TestP24IssueTransitionPositionHelp)$' -count=1 ./internal/workflow ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.169s` (terminated `exit_code: 0`)
  Workflow 0.245 s.
- `go test -tags integration -run '^TestP24(HumanTakeControlJourney|TakeControlResumePair|TakeControlTerminalStateProvenance)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`). First construction attempt 3.631 s: nil template value misclassified as placeholder.
- `go test -tags integration -run '^TestP24(HumanTakeControlJourney|TakeControlResumePair|TakeControlTerminalStateProvenance)$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.703s` (terminated `exit_code: 0`)
  Same command after nil handling repair.
- `go test -tags integration -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`). Expected extension RED, 14.634 s: seven required extension names missing from manifest.
- `go test -tags integration -run '^TestP24TakeControlRebuildExportImportParity$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`). Setup 3.542 s: trusted Human needed the existing concern.record grant.
- `go test -tags integration -run '^TestP24(HumanFixtures|TakeControlRebuildExportImportParity)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`). Generator construction 16.241 s: read result.items instead of the existing result.issues contract.
- `go test -tags integration -run '^TestP24(HumanFixtures|TakeControlRebuildExportImportParity|TakeControlResumePair)$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	28.583s` (terminated `exit_code: 0`)
  Fixtures, parity and pairing pass after causal setup/generator repairs.
- `go test -run '^TestP24DescriptorForbiddenMembers$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.181s` (terminated `exit_code: 0`)

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	81.135s` (terminated `exit_code: 0`)
  cmd/fid 343.934 s, workflow 81.135 s.
- `FID_P24_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/p24-task14-entrypoint-YqRkoM go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP24HumanFixtures|TestP24StopHumanRealCorpus|TestP24TakeControl.*|TestP24HumanTakeControlJourney)$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	70.800s` (terminated `exit_code: 0`)
  Fresh retained-fixture root; also used for independent CLI reads.
- `go test -tags integration -run '^TestP24JourneyProvenanceKeepsExactSourceLine$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`). Expected RED 1.348 s: source recorder could select an unrelated line containing the same revision token.
- `go test -tags integration -run '^TestP24JourneyProvenanceKeepsExactSourceLine$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.579s` (terminated `exit_code: 0`)
  Same command after requiring the supplied exact observed source line.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24(JourneyProvenanceKeepsExactSourceLine|HumanTakeControlJourney|TakeControlRebuildExportImportParity)$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.270s` (terminated `exit_code: 0`)
  Corrected recorder, journey and parity before restarting the full final gate.
- `go build -o /tmp/fid-p24-task14-attempt2 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

- `python3 -c 'import os; print(os.getenv("FID_P24_HUMAN_FIXTURE_ROOT", "<unset>"))'` - PASS - `<unset>` (terminated `exit_code: 0`)
  The final gate uses this unset environment; the tagged suite runs the generator's ordinary t.TempDir branch.

Independent entrypoint workspace:
`/home/mark/projects/fid/.cache/p24-task14-entrypoint-YqRkoM/fixtures-2139482666/take-control`.
Manifest: `/home/mark/projects/fid/.cache/p24-task14-entrypoint-YqRkoM/fixtures-2139482666/manifest.json`, 59 rows.

- `/home/mark/projects/fid/.cache/p24-task14-entrypoint-YqRkoM/fixtures-2139482666/fid --project /home/mark/projects/fid/.cache/p24-task14-entrypoint-YqRkoM/fixtures-2139482666/take-control issue list` - PASS - `P24-eabb13 | qa_review | Human takes control` (terminated `exit_code: 0`)
- `/home/mark/projects/fid/.cache/p24-task14-entrypoint-YqRkoM/fixtures-2139482666/fid --project /home/mark/projects/fid/.cache/p24-task14-entrypoint-YqRkoM/fixtures-2139482666/take-control resume P24-eabb13 --repository /home/mark/projects/fid/.cache/p24-task14-entrypoint-YqRkoM/fixtures-2139482666/take-control` - PASS - `    Resolution requires an approved Decision Reference with the required capabilities and an approved Human Outcome of the same class.` (terminated `exit_code: 0`)
- `/home/mark/projects/fid/.cache/p24-task14-entrypoint-YqRkoM/fixtures-2139482666/fid --project /home/mark/projects/fid/.cache/p24-task14-entrypoint-YqRkoM/fixtures-2139482666/take-control issue transition --help` - PASS - `      --project string   fid project root` (terminated `exit_code: 0`)

The real Human resume printed `Revision: 4` and
`Next states: cancelled (cancelled), qa_block (started), qa_pass (completed)`.
The real help printed both exact pointers to `fid resume <issue-id>`.

## Exact quality-gate results

Final-tree sequence, run byte-identically in the task-specified order.
Each result below has the command's own last output line and observed
termination. No gate was substituted or narrowed.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	120.646s	coverage: 62.2% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	431.049s` (terminated `exit_code: 0`)
- `rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n "PRIVATE_" cmd/fid/integration` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)

| Package | Fresh Go-reported elapsed | Budget/headroom |
|---|---:|---|
| `cmd/fid` | 488.166 s | Below 600 s; above 480 s warning |
| `internal/workflow` | 120.646 s | Below 600 s |
| Slowest other ordinary: `internal/issues` | 228.100 s | Below 600 s |
| Tagged `cmd/fid/integration` | 431.049 s | Below 900 s |

All values are freshly printed in this continuation, not cached or inherited.
The QA-plan commands applicable to this Developer task are the same seven
commands. The fixture-retention command at the close commit, Human acceptance,
benchmarks, archive-scale and full-checkptr controls remain the QA plan's
explicitly QA-owned phase-end work, rather than omitted Developer gates.

## Deviations from spec

- No accepted contract, provenance class, assertion semantics, command, threshold, budget, or coverage policy was changed. Product edits are exactly DEC-290 clause 3; parity/setup JSON is outside the Human journey.
- Mandatory causal-repair note for the preserved attempt-1 required tagged gate: first failure token `FAIL` (terminated `exit_code: 1`). Root cause: Human output lacked both the concrete terminal-success state and usable current Issue revision. Accepted resolution `TRG-24-016 / DEC-290` supplies bounded resume fields and help usages. This continuation first added and ran the clause 3 RED tests, implemented that contract, and sourced both transition arguments from the final Human resume. The byte-identical required tagged command's later final-tree PASS result is recorded above; the old failure remains in the historical report below.
- Focused construction failures remain listed in Tests run. The first journey failure token `FAIL` (terminated `exit_code: 1`) came from treating nil template values as placeholders; the recorder now accepts exact unchanged template values at their original paths, and the same focused journey command later printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.703s` (terminated `exit_code: 0`). The fixture-extension RED token `FAIL` (terminated `exit_code: 1`) showed seven missing names; generator rows now cover all required judgment/authoring workspaces. Parity setup's `FAIL` (terminated `exit_code: 1`) showed its Human needed the existing concern.record grant; setup now grants concern.record/deviation.record and seeds their records before the Human journey. Generator construction's `FAIL` (terminated `exit_code: 1`) came from result.items; it now reads result.issues. The broader focused fixture/parity/pairing command later printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	28.583s` (terminated `exit_code: 0`), and the final byte-identical tagged gate passes all of these tests.
- A final provenance audit found that matching a revision token against an entire output could attach a different line containing the same token. The new exact-line regression first printed `FAIL` (terminated `exit_code: 1`, 1.348 s), before changing the recorder. The recorder now validates the supplied single source line against earlier observed output and retains that exact line. The byte-identical regression later printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.579s` (terminated `exit_code: 0`); the affected race selection printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.270s` (terminated `exit_code: 0`). No retained assertion was altered by this repair.
- The first full seven-command run in this continuation terminated successfully, but its tagged binary had compiled before that final recorder repair. It is superseded evidence for completion, preserved below. The full exact sequence was restarted on the corrected final tree, without concurrent heavy testing, narrowing any command, changing protected definitions, or mining a flaky green.
- Retained-test provenance: `p24_human_fixtures_test.go` assertions originated in `ab05e496206809860081f941d6c1c171eb845cca`, `feat(cli): checkpoint task 11 Human views and fixture contract blocker`, `Refs: phase-24/task-11`. This task's changes to that file are additions only. The terminal-state diagnostic originated in `4b6a3837e52f24a5a8975a1d2abef94635c74823`, `test(integration): checkpoint Task 14 terminal-state provenance blocker`, `Refs: phase-24/task-14`; it keeps its original checks, fatal diagnostic and forbidden-source guard and gains the exact two-line position check. No assertion displaced by clause 3 required clause 6 convergence.

### Superseded first seven-command run in this continuation

These are observed results from this session, retained as historical evidence;
the final-tree sequence above is the completion evidence.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	162.186s	coverage: 62.2% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	390.024s` (terminated `exit_code: 0`)
- `rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n "PRIVATE_" cmd/fid/integration` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)

| Package | Fresh elapsed in superseded run | Budget/headroom |
|---|---:|---|
| `cmd/fid` | 530.567 s | Below 600 s; above 480 s warning |
| `internal/workflow` | 162.186 s | Below 600 s |
| Slowest other ordinary: `internal/export` | 266.660 s | Below 600 s |
| Tagged `cmd/fid/integration` | 390.024 s | Below 900 s |

### Preserved attempt 1 BLOCKED report

Verbatim content from report commit
`98d45617b35c212683b4346d30a760eb8ed630d3` follows, quoted to keep
historical evidence distinct from the current outcome. Its historical PASS
labels are not attestations of commands executed in this continuation.

> # Phase 24 - Task 14: Human-takes-control journey and phase conformance - Blocker Report
>
> ## Workflow outcome
>
> - Task outcome: BLOCKED
> - Category: public-interface-decision
> - Condition: VAL-049 needs the terminal-success workflow state for `issue transition --to`, but the stuck fixture's Human reads and the entire advertised help tree never print `qa_pass`. That value cannot be sourced from any permitted DEC-278 provenance class.
> - Recovery attempts: 0
> - Next action: Tech Lead disposes the missing terminal-state provenance as a product-gap candidate and supplies accepted authority and bounded implementation scope for exposing the necessary value. Then resume Task 14, retaining the closed provenance classes and all failed evidence.
> - Human attention required: yes
>
> ## Summary of what was implemented
>
> - Added a deterministic activated start fixture with a Builder Claim/Run, a security Human Outcome Obligation and an unsatisfied protected security Decision. It has no Receipt-backed Obligation. The Builder's terminal-success transition returns `protected-decision-required`, exit 5, with blocking IDs and Decision/Outcome guidance, without appending an Event.
> - Extended `TestP24HumanFixtures` to retain a `take-control-start` manifest row alongside every Task 11 denial and view. The specified first failing manifest assertion is green after fixture construction.
> - Added `TestP24TakeControlTerminalStateProvenance`, a real-process diagnostic that rediscovers the Issue from Human output, reads the stuck state's Human projections and recursively walks the advertised help tree. It fails because no permitted source supplies the terminal-success state. This is supporting blocker evidence, not a completed VAL-049 journey.
> - No product behavior, accepted contract, workflow control, runner ledger or QA artifact was changed. Task 14 is not complete.
>
> ## Files created or modified
>
> - `cmd/fid/integration/p24_take_control_test.go`
> - `cmd/fid/integration/p24_human_fixtures_test.go`
> - `project-docs/docs/phases/phase-24/tasks/task-14-done.md`
>
> ## Documentation updated
>
> - This blocker report records the protected gap and durable partial work.
> - `project-docs/specs/08_QUALITY_PLAN.md` has not been updated to claim a retained completed journey; that required update remains unfinished behind the blocker.
>
> ## Validation Contract Coverage
>
> - VAL-049: BLOCKED — `TestP24TakeControlTerminalStateProvenance`, `cmd/fid/integration/p24_take_control_test.go:66`, proves the real protected refusal at lines 71-76 and the missing terminal-state source at lines 94-103. `issue transition --help` prints only `--to string Target state`; `workflow list` prints workflow names. The test never reads JSON-mode output, SQL or journal data to supply the acting caller's values.
> - VAL-050: BLOCKED — the required completed VAL-049 workspace and its Phase 24 live/rebuilt/imported parity have not been produced. No parity PASS is claimed.
> - VAL-051: BLOCKED — the registry-wide forbidden-member scan has not been added. The exact absence-command results are recorded below; the existing Stop renderer assertions are retained but do not alone close this assertion.
> - VAL-063: PASS — the exact `rg -n "PRIVATE_" cmd/fid/integration` command terminated `exit_code: 1`, with no output, as predeclared by the task spec.
> - VAL-052: BLOCKED for complete fresh timing evidence — freshly reported package times are below their budgets, but the ordinary gate returned `(cached)` for `internal/workflow` and most other packages. No prior-session elapsed value is copied or claimed as observed. See the timing table below.
> - VAL-053 preparation: partial — `TestP24HumanFixtures`, `cmd/fid/integration/p24_human_fixtures_test.go:68`, requires the new `take-control-start` workspace. Human acceptance is not claimed.
> - VAL-054/055 preparation: retained — Task 11 denial/view assertions remain byte-unchanged. Introducing commit: `ab05e496206809860081f941d6c1c171eb845cca`, `feat(cli): checkpoint task 11 Human views and fixture contract blocker`, `Refs: phase-24/task-11`. The diff only adds the start-state assertion and generator calls. No displaced assertion was converged, deleted or weakened.
> - VAL-056/057 preparation: BLOCKED — judgment-view extension and the three authoring workspaces remain unfinished; no Human judgment or QA review was created.
>
> ## Tests run
>
> All commands ran from `/home/mark/projects/fid`.
>
> - `go test -tags integration -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` — attempt 1, expected RED — `FAIL` (terminated `exit_code: 1`). Package time 12.540 s; failure `manifest missing take-control-start workspace`.
> - The same byte-identical focused fixture command — attempt 2, FAIL — `FAIL` (terminated `exit_code: 1`). Package time 9.145 s; new setup incorrectly supplied Activation phase 24 instead of the accepted phase 21.
> - The same byte-identical focused fixture command — attempt 3, FAIL — `FAIL` (terminated `exit_code: 1`). Package time 9.375 s; setup staged only one of the two required Activation components.
> - The same byte-identical focused fixture command — attempt 4, FAIL — `FAIL` (terminated `exit_code: 1`). Package time 9.643 s; setup's Issue transitions lacked the Claim, active Run and expected revision required by the activated context.
> - `go test -tags integration -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	14.669s` (terminated `exit_code: 0`)
> - `go test -tags integration -run '^TestP24TakeControlTerminalStateProvenance$' -count=1 ./cmd/fid/integration` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Package time 3.172 s. Final diagnostic: `BLOCKED provenance: issue.transition --to requires terminal-success state qa_pass; no earlier Human fid output or help supplies it, and DEC-278 does not permit setup workflow states or caller-authored state values`.
> - `go build -o /tmp/fid-p24-task14 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - Real entrypoint: `/tmp/fid-p24-task14 issue transition --help`, `/tmp/fid-p24-task14 workflow define --help`, and `/tmp/fid-p24-task14 --help` each terminated with observed `exit_code: 0`. The first printed `--to string Target state`; none printed the needed state. Their complete output was inspected in this session. This corroborates the diagnostic and does not establish journey completion.
>
> The passing fixture command ran with the fixture-root variable unset. Setup
> now stages both accepted Activation components and supplies the Builder's
> Claim, Run and expected Issue revisions. Every existing denial/view assertion
> is preserved.
>
> ## Exact quality-gate results
>
> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
> - `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
> - `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.8%` (terminated `exit_code: 0`)
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
> - `rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
> - `rg -n "PRIVATE_" cmd/fid/integration` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
>
> The tagged gate's only failed test was
> `TestP24TakeControlTerminalStateProvenance` (5.37 s). Its final diagnostic
> is the same protected provenance gap as the focused run. No required gate
> was narrowed, substituted, changed, interrupted or rerun after failure.
> The QA-plan gates applicable to this Developer task are the same exact
> format, lint, race/coverage, coverage-report, tagged integration and absence
> commands above. Phase-end benchmarks, archive-scale, full-checkptr and the
> QA-owned fixture-retention command remain phase-end QA work.
>
> | Package | Go-reported elapsed/result | Budget evidence |
> |---|---|---|
> | `cmd/fid` | 337.575 s | Below 480 s warning and 600 s budget |
> | `internal/workflow` | `(cached)`; no elapsed value printed | Fresh wall time UNVERIFIED; no inherited value asserted |
> | Slowest other ordinary package with a freshly printed time: `internal/export` | 107.982 s | Below 480 s warning and 600 s budget; other cached packages cannot be ranked by fresh time |
> | Tagged `cmd/fid/integration` | 362.182 s, test failure | Below 900 s budget; test outcome remains failed |
>
> ## Deviations from spec
>
> - Task 14 explicitly requires escalation when the journey needs a value outside the closed provenance classes. Implementation stopped at that protected gap. The journey, parity, forbidden-member scan, judgment/authoring fixture extension and quality-plan update are unfinished, not deferred to another task.
> - The initial fixture setup failures above were causal test-construction errors. No required-gate attempt was repaired or rerun to mine a green. The first red token was `FAIL` (terminated `exit_code: 1`); setup repairs aligned Activation and Builder prerequisites to existing contracts, and the same command later printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	14.669s` (terminated `exit_code: 0`). All failed attempts are preserved.
> - The pre-gate green tier cannot be satisfied while the protected provenance diagnostic is red. The exact full commands are nevertheless run in their specified order to record their actual results, not to claim completion.
> - Required tagged-gate attempt 1 printed `FAIL` (terminated `exit_code: 1`). Root cause: the default stuck fixture needs `qa_pass` for the terminal transition, but Human output/help supplies no permitted source. Repair: none authorized by Task 14's test/documentation scope; changing the product or provenance/assertion semantics is explicitly out of scope. Later result: none; the failure remains preserved and requires Tech Lead disposition.
>
> ## Concerns, technical debt, or notes for the Tech Lead
>
> - Product-gap candidate: before the first terminal-success request or its retry, the acting caller needs the workflow's terminal state and a permitted provenance source. The default fixture's only Issue is `qa_review`. `resume` identifies the required completed workflow category, not its concrete state. `workflow list` shows only `Epic default`, `Subtask default`, and `Task default`. Every advertised help page was inspected by the diagnostic; none supplied `qa_pass`.
> - Current Task 14 scope is tests plus the quality-plan documentation update; its out-of-scope section prohibits new product behavior and requires escalation of journey gaps. Choosing a new Human read/help behavior or allowing setup/caller-authored workflow states needs protected authority. Do not weaken VAL-049 or widen DEC-278's classes to close this task.
> - The start fixture uses setup-only JSON while constructing accepted state. Its acting-caller diagnostic uses Human mode only; setup values are not claimed as Human provenance. The blocked Builder request is the fixture precondition, not a provenance-compliant caller step.
> - A separate source for the activated Issue's current revision may also be needed: the setup denial uses revision 4, while `issue get` and `resume` do not print that Issue revision. This is a follow-up inspection concern, not an independently completed journey finding.
> - QA Window E remains scheduled after Task 14. The Phase Runner owns triage/escalation allocation and the ledger.
>
> ## Dependency-budget note
>
> - No dependency changes.
>
> ## Commit hash
>
> - `4b6a3837e52f24a5a8975a1d2abef94635c74823`
>
> ## Commit subject
>
> - `test(integration): checkpoint Task 14 terminal-state provenance blocker`

## Concerns, technical debt, or notes for the Tech Lead

- Headroom warning: final-tree `cmd/fid` is 488.166 s, at or above 80% of its 600 s budget (480 s). The superseded current-session run was 530.567 s and carries the same warning. Both terminated successfully and are below the budget; no timing threshold was changed.
- VAL-053..057 deterministic fixture preparation is complete; the Human judgments and QA Window E remain scheduled phase-end work. The retained manifest and binary above are available for inspection; QA's declared command regenerates them at the task close commit.
- The earlier public-interface gap is resolved only by the bounded DEC-290 contract. Setup actors/project directory remain the only setup provenance; JSON pairing/parity does not supply journey values.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `e97835b9f9983e357b45e80c3dbe4e4b02cedb6d`

## Commit subject

- `feat(workflow): expose resume position for Human take-control`
