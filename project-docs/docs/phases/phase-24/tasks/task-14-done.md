# Phase 24 - Task 14: Human-takes-control journey and phase conformance - Blocker Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: public-interface-decision
- Condition: VAL-049 needs the terminal-success workflow state for `issue transition --to`, but the stuck fixture's Human reads and the entire advertised help tree never print `qa_pass`. That value cannot be sourced from any permitted DEC-278 provenance class.
- Recovery attempts: 0
- Next action: Tech Lead disposes the missing terminal-state provenance as a product-gap candidate and supplies accepted authority and bounded implementation scope for exposing the necessary value. Then resume Task 14, retaining the closed provenance classes and all failed evidence.
- Human attention required: yes

## Summary of what was implemented

- Added a deterministic activated start fixture with a Builder Claim/Run, a security Human Outcome Obligation and an unsatisfied protected security Decision. It has no Receipt-backed Obligation. The Builder's terminal-success transition returns `protected-decision-required`, exit 5, with blocking IDs and Decision/Outcome guidance, without appending an Event.
- Extended `TestP24HumanFixtures` to retain a `take-control-start` manifest row alongside every Task 11 denial and view. The specified first failing manifest assertion is green after fixture construction.
- Added `TestP24TakeControlTerminalStateProvenance`, a real-process diagnostic that rediscovers the Issue from Human output, reads the stuck state's Human projections and recursively walks the advertised help tree. It fails because no permitted source supplies the terminal-success state. This is supporting blocker evidence, not a completed VAL-049 journey.
- No product behavior, accepted contract, workflow control, runner ledger or QA artifact was changed. Task 14 is not complete.

## Files created or modified

- `cmd/fid/integration/p24_take_control_test.go`
- `cmd/fid/integration/p24_human_fixtures_test.go`
- `project-docs/docs/phases/phase-24/tasks/task-14-done.md`

## Documentation updated

- This blocker report records the protected gap and durable partial work.
- `project-docs/specs/08_QUALITY_PLAN.md` has not been updated to claim a retained completed journey; that required update remains unfinished behind the blocker.

## Validation Contract Coverage

- VAL-049: BLOCKED — `TestP24TakeControlTerminalStateProvenance`, `cmd/fid/integration/p24_take_control_test.go:66`, proves the real protected refusal at lines 71-76 and the missing terminal-state source at lines 94-103. `issue transition --help` prints only `--to string Target state`; `workflow list` prints workflow names. The test never reads JSON-mode output, SQL or journal data to supply the acting caller's values.
- VAL-050: BLOCKED — the required completed VAL-049 workspace and its Phase 24 live/rebuilt/imported parity have not been produced. No parity PASS is claimed.
- VAL-051: BLOCKED — the registry-wide forbidden-member scan has not been added. The exact absence-command results are recorded below; the existing Stop renderer assertions are retained but do not alone close this assertion.
- VAL-063: PASS — the exact `rg -n "PRIVATE_" cmd/fid/integration` command terminated `exit_code: 1`, with no output, as predeclared by the task spec.
- VAL-052: BLOCKED for complete fresh timing evidence — freshly reported package times are below their budgets, but the ordinary gate returned `(cached)` for `internal/workflow` and most other packages. No prior-session elapsed value is copied or claimed as observed. See the timing table below.
- VAL-053 preparation: partial — `TestP24HumanFixtures`, `cmd/fid/integration/p24_human_fixtures_test.go:68`, requires the new `take-control-start` workspace. Human acceptance is not claimed.
- VAL-054/055 preparation: retained — Task 11 denial/view assertions remain byte-unchanged. Introducing commit: `ab05e496206809860081f941d6c1c171eb845cca`, `feat(cli): checkpoint task 11 Human views and fixture contract blocker`, `Refs: phase-24/task-11`. The diff only adds the start-state assertion and generator calls. No displaced assertion was converged, deleted or weakened.
- VAL-056/057 preparation: BLOCKED — judgment-view extension and the three authoring workspaces remain unfinished; no Human judgment or QA review was created.

## Tests run

All commands ran from `/home/mark/projects/fid`.

- `go test -tags integration -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` — attempt 1, expected RED — `FAIL` (terminated `exit_code: 1`). Package time 12.540 s; failure `manifest missing take-control-start workspace`.
- The same byte-identical focused fixture command — attempt 2, FAIL — `FAIL` (terminated `exit_code: 1`). Package time 9.145 s; new setup incorrectly supplied Activation phase 24 instead of the accepted phase 21.
- The same byte-identical focused fixture command — attempt 3, FAIL — `FAIL` (terminated `exit_code: 1`). Package time 9.375 s; setup staged only one of the two required Activation components.
- The same byte-identical focused fixture command — attempt 4, FAIL — `FAIL` (terminated `exit_code: 1`). Package time 9.643 s; setup's Issue transitions lacked the Claim, active Run and expected revision required by the activated context.
- `go test -tags integration -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	14.669s` (terminated `exit_code: 0`)
- `go test -tags integration -run '^TestP24TakeControlTerminalStateProvenance$' -count=1 ./cmd/fid/integration` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Package time 3.172 s. Final diagnostic: `BLOCKED provenance: issue.transition --to requires terminal-success state qa_pass; no earlier Human fid output or help supplies it, and DEC-278 does not permit setup workflow states or caller-authored state values`.
- `go build -o /tmp/fid-p24-task14 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- Real entrypoint: `/tmp/fid-p24-task14 issue transition --help`, `/tmp/fid-p24-task14 workflow define --help`, and `/tmp/fid-p24-task14 --help` each terminated with observed `exit_code: 0`. The first printed `--to string Target state`; none printed the needed state. Their complete output was inspected in this session. This corroborates the diagnostic and does not establish journey completion.

The passing fixture command ran with the fixture-root variable unset. Setup
now stages both accepted Activation components and supplies the Builder's
Claim, Run and expected Issue revisions. Every existing denial/view assertion
is preserved.

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- `rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n "PRIVATE_" cmd/fid/integration` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)

The tagged gate's only failed test was
`TestP24TakeControlTerminalStateProvenance` (5.37 s). Its final diagnostic
is the same protected provenance gap as the focused run. No required gate
was narrowed, substituted, changed, interrupted or rerun after failure.
The QA-plan gates applicable to this Developer task are the same exact
format, lint, race/coverage, coverage-report, tagged integration and absence
commands above. Phase-end benchmarks, archive-scale, full-checkptr and the
QA-owned fixture-retention command remain phase-end QA work.

| Package | Go-reported elapsed/result | Budget evidence |
|---|---|---|
| `cmd/fid` | 337.575 s | Below 480 s warning and 600 s budget |
| `internal/workflow` | `(cached)`; no elapsed value printed | Fresh wall time UNVERIFIED; no inherited value asserted |
| Slowest other ordinary package with a freshly printed time: `internal/export` | 107.982 s | Below 480 s warning and 600 s budget; other cached packages cannot be ranked by fresh time |
| Tagged `cmd/fid/integration` | 362.182 s, test failure | Below 900 s budget; test outcome remains failed |

## Deviations from spec

- Task 14 explicitly requires escalation when the journey needs a value outside the closed provenance classes. Implementation stopped at that protected gap. The journey, parity, forbidden-member scan, judgment/authoring fixture extension and quality-plan update are unfinished, not deferred to another task.
- The initial fixture setup failures above were causal test-construction errors. No required-gate attempt was repaired or rerun to mine a green. The first red token was `FAIL` (terminated `exit_code: 1`); setup repairs aligned Activation and Builder prerequisites to existing contracts, and the same command later printed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	14.669s` (terminated `exit_code: 0`). All failed attempts are preserved.
- The pre-gate green tier cannot be satisfied while the protected provenance diagnostic is red. The exact full commands are nevertheless run in their specified order to record their actual results, not to claim completion.
- Required tagged-gate attempt 1 printed `FAIL` (terminated `exit_code: 1`). Root cause: the default stuck fixture needs `qa_pass` for the terminal transition, but Human output/help supplies no permitted source. Repair: none authorized by Task 14's test/documentation scope; changing the product or provenance/assertion semantics is explicitly out of scope. Later result: none; the failure remains preserved and requires Tech Lead disposition.

## Concerns, technical debt, or notes for the Tech Lead

- Product-gap candidate: before the first terminal-success request or its retry, the acting caller needs the workflow's terminal state and a permitted provenance source. The default fixture's only Issue is `qa_review`. `resume` identifies the required completed workflow category, not its concrete state. `workflow list` shows only `Epic default`, `Subtask default`, and `Task default`. Every advertised help page was inspected by the diagnostic; none supplied `qa_pass`.
- Current Task 14 scope is tests plus the quality-plan documentation update; its out-of-scope section prohibits new product behavior and requires escalation of journey gaps. Choosing a new Human read/help behavior or allowing setup/caller-authored workflow states needs protected authority. Do not weaken VAL-049 or widen DEC-278's classes to close this task.
- The start fixture uses setup-only JSON while constructing accepted state. Its acting-caller diagnostic uses Human mode only; setup values are not claimed as Human provenance. The blocked Builder request is the fixture precondition, not a provenance-compliant caller step.
- A separate source for the activated Issue's current revision may also be needed: the setup denial uses revision 4, while `issue get` and `resume` do not print that Issue revision. This is a follow-up inspection concern, not an independently completed journey finding.
- QA Window E remains scheduled after Task 14. The Phase Runner owns triage/escalation allocation and the ledger.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `4b6a3837e52f24a5a8975a1d2abef94635c74823`

## Commit subject

- `test(integration): checkpoint Task 14 terminal-state provenance blocker`
