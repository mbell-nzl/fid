# Phase 24 - Task 13: Judgment disclaimers once per section; recorded reasoning shown - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: gate-repair
- Condition: Task 13 A7 invalidated the retained hidden-reasoning assertion in cmd/fid/integration/p23_concern_list_human_test.go:29. This file is outside Task 13's named file scope; the required tagged gate remains red.
- Recovery attempts: 0
- Next action: Phase Runner obtains a bounded Task 13 scope amendment naming cmd/fid/integration/p23_concern_list_human_test.go, then dispatches a same-task continuation to converge its displayed-reasoning assertion, retain all other assertions, rerun the byte-identical full gate, and finalize this report.
- Human attention required: no

## Summary of what was implemented

- Added one shared section-disclaimer function for non-empty Concerns, Deviations and Human Outcomes sections, before their items. Empty sections retain `None.`; item-level disclaimers are removed.
- Showed all closed recorded-reasoning fields, including every Finding disposition reason. Lists and resume sections collapse whitespace and cut values longer than 80 Unicode scalars to the first 80 plus `…`. Detail reads, mutation outputs and review get retain full text and indent retained line breaks.
- Escaped ESC, CR, BEL and NUL using the existing control-text escape function. List tabs collapse as whitespace; detail tabs escape.
- Replaced the retained leak predicates with a referenced-file sentinel guard. Neutral `RECORDED_` values replace the old reasoning fixture markers and matching conformance expected facts.
- Repaired the sentinel fixture to use its existing ignored `*.json` convention, so seeding referenced content does not alter an already bound repository snapshot or stale Gate Evidence.
- Updated the glossary and CLI examples. The switch and layout redesign remain deferred; Task 11 Stop/repository lines were not changed.
- Implementation is checkpointed for the named retained-test scope gap. No out-of-scope test was edited.

## Files created or modified

- `cmd/fid/p23_human.go`
- `cmd/fid/resume_human.go` (judgment sections only)
- `cmd/fid/p23_human_test.go`
- `cmd/fid/integration/p23_human_leaf_pairs_test.go`
- `cmd/fid/integration/p23_human_fixtures_test.go`
- `cmd/fid/integration/testdata/p23_conformance/operations.json`
- `project-docs/specs/GLOSSARY.md`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-24/tasks/task-13-done.md`

## Documentation updated

- `project-docs/specs/GLOSSARY.md:21` defines Project-owned prose as referenced project-file content; `:32` defines Recorded reasoning with every field in the accepted closed list.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` documents section disclaimers, list/detail rendering, escaping and matching Concern examples.
- This report preserves failed attempts, convergence provenance and the exact remaining scope gap. The runner-owned run ledger was not edited, staged or committed; its observed SHA-256 remained `82b416d1df573d6872f551bbab412e5357bf53663029a47544c9769b40669d07`.

## Validation Contract Coverage

- VAL-044: PASS focused evidence — `TestP24JudgmentSectionDisclaimers`, `cmd/fid/p23_human_test.go:177`, checks zero, one and three items in every Concern, Deviation and Human Outcome leaf view and their resume sections: no disclaimer with none, exactly one before non-empty items, no old item sentences. The view set at `:138` follows the callers of `writeP23Human`/`writeP23MutationHuman` and `resumeJudgmentHumanLines`. The production absence command also passed with its declared exit 1.
- VAL-045: PASS focused evidence — `TestP24RecordedReasoningEveryFieldAndView`, `cmd/fid/p23_human_test.go:213`, covers each closed field at 80/81 non-ASCII scalars, multiline and whitespace cases in all list, detail, mutation and resume views, including nested Review Findings and a second Finding disposition. `TestPhase23HumanLeafTextMatchesRecordedFacts` retains exact text/fact comparisons. `TestP24JudgmentDetailAndResumeSentinel`, `cmd/fid/integration/p23_human_leaf_pairs_test.go:455`, checks the real three get leaves and resume. The manual CLI record/list/get below demonstrates full versus cut output.
- VAL-046: PASS focused evidence — the `controls` cases of `TestP24RecordedReasoningEveryFieldAndView` assert escaped ESC/CR/BEL/NUL and absence of raw bytes in every reasoning view. The real CLI mutation/get retain the indented literal `\\u001b\\r\\u0007\\u0000` continuation.
- VAL-047: PASS focused evidence — `p24SeedJudgmentSource`, `cmd/fid/integration/p23_human_leaf_pairs_test.go:367`, seeds a real referenced file with `PROJECT_FILE_CONTENT_SENTINEL_P24_JUDGMENT`. `TestP23HumanFixtures` scans every leaf-pair Human output and checkpoint Human inspection on generation and post-exit replay; `TestP24JudgmentDetailAndResumeSentinel` adds all three get leaves and resume and asserts reference metadata remains visible. Both passed in the final focused tagged run. `TestP24SentinelPreservesRepositoryBinding` at `:502` also proves the sentinel does not change the repository snapshot.
- VAL-048: PASS documentation evidence — `GLOSSARY.md:21` and `:32` explicitly distinguish referenced file content from Concern summary/reasoning/disposition reason, Finding summary/every disposition reason, Deviation expected/observed/reasoning and Human Outcome reasoning, matching contract section 10. QA Lead's review remains at its prescribed checkpoint.
- VAL-052: PASS timing evidence — the exact ordinary and tagged executions reported the times below; all recorded ordinary packages are under 600 s and tagged integration is under 900 s. The tagged execution failed assertions, recorded separately below.

| Package | Exact gate reported time | Budget |
|---|---|---|
| cmd/fid | 530.445 s | 600 s |
| internal/workflow | 136.519 s | 600 s |
| internal/issues (slowest other ordinary package) | 238.583 s | 600 s |
| cmd/fid/integration (failed exact attempt) | 358.688 s | 900 s |

The touched-package ordinary pre-gate reported 345.812 s; the initial focused tagged pre-gate reported 24.632 s and the final repair selection reported 20.353 s.

- VAL-063: PASS — the exact `rg -n "PRIVATE_" cmd/fid/integration` command terminated exit 1 with no output. Both old leak predicates were replaced and all four matching expected facts were renamed.
- VAL-056: DEFERRED — accepted Human phase-end review; this task supplies deterministic evidence and fixtures, not a Human acceptance result.
- Full task closure remains pending the required integration gate's retained-test convergence below.

## Tests run

Working directory: `/home/mark/projects/fid` unless the runtime workspace is stated. Every recorded test execution was watched to termination.

Expected red/green for the first summary slice:

- `go test ./cmd/fid -run '^TestP24ConcernListSummary$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestP24ConcernListSummary$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.183s` (terminated `exit_code: 0`)

Expected red for disclaimer placement and remaining fields, then byte-identical green after implementation:

- `go test ./cmd/fid -run '^TestP24(JudgmentSectionDisclaimers|RecordedReasoningEveryFieldAndView)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestP24(JudgmentSectionDisclaimers|RecordedReasoningEveryFieldAndView)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.917s` (terminated `exit_code: 0`)

Expected red/green for list whitespace collapse:

- `go test ./cmd/fid -run '^TestP24RecordedReasoningEveryFieldAndView/concern.list/summary/whitespace$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test ./cmd/fid -run '^TestP24RecordedReasoningEveryFieldAndView/concern.list/summary/whitespace$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.187s` (terminated `exit_code: 0`)

Additional focused green checks:

- `go test ./cmd/fid -run '^TestP24(ConcernListSummary|JudgmentSectionDisclaimers|RecordedReasoningEveryFieldAndView)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.201s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestPhase23HumanLeafTextMatchesRecordedFacts$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.172s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^Test(P24(ConcernListSummary|JudgmentSectionDisclaimers|RecordedReasoningEveryFieldAndView)|Phase23HumanLeafTextMatchesRecordedFacts)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.177s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP23(HumanFixtures|SharedLeafConformance)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	16.522s` (terminated `exit_code: 0`)
- `go test -tags integration ./cmd/fid/integration -run '^TestP24JudgmentDetailAndResumeSentinel$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	5.620s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	345.812s` (terminated `exit_code: 0`)
- `FID_P23_HUMAN_FIXTURE_ROOT=/tmp/fid-task13-human.zZI1np go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^Test(P23(HumanFixtures|HumanFixtureStepsRemainRunnable|SharedLeafConformance)|P24JudgmentDetailAndResumeSentinel)$' -v` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	24.632s` (terminated `exit_code: 0`)

Required full tagged attempt 1 failed both retained assertions; this is the same captured execution classified BLOCKED under the exact quality-gate section:

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)

The focused reproduction first failed both assertions:

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP23ConcernListHumanProjection|TestP24All49KeyedLeavesRebuildInOneJournal)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

The new sentinel-binding regression failed before the causal helper repair and passed after it:

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24SentinelPreservesRepositoryBinding$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24SentinelPreservesRepositoryBinding$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.657s` (terminated `exit_code: 0`)

The byte-identical two-test selection then reached the later steps of the 49-leaf fixture and failed only the out-of-scope Concern assertion:

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP23ConcernListHumanProjection|TestP24All49KeyedLeavesRebuildInOneJournal)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

Final in-scope repair verification, including unchanged 49-leaf/rebuild assertions:

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^Test(P23(HumanFixtures|HumanFixtureStepsRemainRunnable|SharedLeafConformance)|P24(JudgmentDetailAndResumeSentinel|SentinelPreservesRepositoryBinding|All49KeyedLeavesRebuildInOneJournal))$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	20.353s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)

Real entrypoint exercise used the retained CLI binary in `/tmp/fid-task13-human.zZI1np/checkpoints-486060414/leaf-contract`. The request in `/tmp/fid-task13-human.zZI1np/runtime-concern.json` has 81 copies of `界`, then LF and `Recorded continuation: ` followed by ESC/CR/BEL/NUL. The record/get show full text with indentation and escapes; the two-item list has one disclaimer and the first 80 scalars plus ellipsis. Referenced file content never appeared.

- `/tmp/fid-task13-human.zZI1np/checkpoints-486060414/fid --project /tmp/fid-task13-human.zZI1np/checkpoints-486060414/leaf-contract concern record --actor actor-recovery --idempotency-key task13-runtime --input /tmp/fid-task13-human.zZI1np/runtime-concern.json` - PASS - `      Recorded continuation: \u001b\r\u0007\u0000` (terminated `exit_code: 0`)
- `/tmp/fid-task13-human.zZI1np/checkpoints-486060414/fid --project /tmp/fid-task13-human.zZI1np/checkpoints-486060414/leaf-contract concern list` - PASS - `    Reasoning: 界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界界…` (terminated `exit_code: 0`)
- `/tmp/fid-task13-human.zZI1np/checkpoints-486060414/fid --project /tmp/fid-task13-human.zZI1np/checkpoints-486060414/leaf-contract concern get CON-runtime` - PASS - `      Recorded continuation: \u001b\r\u0007\u0000` (terminated `exit_code: 0`)

## Exact quality-gate results

The seven required commands were executed in their declared order. Commands 1-4 passed; command 5 failed, and its remaining out-of-scope assertion prevents closure. Commands 6-7 passed with their declared non-zero exits. Post-repair format/lint and focused tests above do not substitute for a green full tagged gate.

### Expected Non-Zero Commands

| Command | Declared expectation | Reason |
|---|---|---|
| `rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'` | `exit_code: 1`, no output | No per-item disclaimer remains |
| `rg -n "PRIVATE_" cmd/fid/integration` | `exit_code: 1`, no output | Old leak checks and recorded-reasoning markers are replaced |

These are the task's two required conditional absence commands.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	136.519s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`)
- Integration reason: the above is the observed first failed attempt, not a second execution. The sentinel-caused stale Evidence was repaired in scope, but `TestP23ConcernListHumanProjection` still forbids Task 13's accepted recorded summary. Full byte-identical rerun is pending the bounded scope amendment; no passing gate is claimed.
- `rg -n "This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment" cmd/fid -g '!*_test.go'` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)
- `rg -n "PRIVATE_" cmd/fid/integration` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)

## Deviations from spec

- No product or accepted-contract deviation. A6/A7 implementation and named documentation are complete; the task remains RECOVERY_REQUIRED for one retained-test scope gap.
- First full tagged failure: final token `FAIL` (terminated `exit_code: 1`). Root causes were (1) Task 13's accepted A7 invalidating a retained hidden-summary predicate and (2) this task's sentinel helper committing a new file after an existing caller had recorded Gate Evidence, changing the bound tree and causing `evidence-stale`. The focused two-test command reproduced both with `FAIL`/exit 1.
- Causal fixture repair: use the already accepted fixture `*.json` ignore convention for the referenced sentinel file, remove the helper's Git commit/re-observation, and assert exact before/after repository equality. The new regression progressed from `FAIL`/exit 1 to `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.657s`/exit 0. The byte-identical two-test selection then failed only the Concern predicate, retaining `FAIL`/exit 1. The final focused selection passed at 20.353 s with every 49-leaf inventory, later completion step and rebuild parity assertion retained. Neither Evidence freshness nor a gate assertion was weakened.
- Expected implementation reds are preserved above: absent Concern summary, missing remaining reasoning/disclaimer placement, and uncollapsed list tabs. Each ended with `FAIL`/exit 1; their byte-identical green commands are recorded above.
- Scope gap: `cmd/fid/integration/p23_concern_list_human_test.go:29` is absent from the task's named paths. Its obsolete `strings.Contains(human, "private concern prose")` rejection was invalidated by this task's accepted A7 change. No repair is assigned to another task. The continuation needs that exact path authorized so it can assert displayed summary/reasoning while preserving purpose, identity, filtering, bounds, malformed-filter rejection and no-journal-change assertions.

Retained-test convergence provenance was established with introducing-commit Git log/show/blame, not filenames or embedded task numbers:

| Retained subject | Introducing commit / phase / task | Task 13 disposition |
|---|---|---|
| Unit Human exact text and free-text suppression | `59d2653d`, phase-23/task-14 | Retain exact fact/text comparisons, add accepted recorded fields, move disclaimers to section expectations, retain non-proof meaning |
| Leaf-pair inputs, exact expected Human output, leak predicates and post-exit replay | `b71ed1ba`, phase-23/task-14 | Rename reasoning markers, retain JSON/replay equality, replace suppression predicates with actual referenced-file sentinel scans |
| `operations.json` and its conformance reader | File/reader added by `33351494`, phase-23/task-18; the four reasoning expected facts introduced by `38efa3b8`, phase-23/task-18 | Rename only corresponding expected values; reader, shape, historical ownership fields and semantic assertions retained |
| 49-leaf single-journal/rebuild fixture | `95c738f2`, phase-24/task-2 | Test unchanged; repair this task's shared helper that made its earlier Evidence stale |
| Generic Concern Human projection's hidden-summary predicate | `59d2653d`, phase-23/task-14 | Unmodified because outside named scope; exact path and accepted invalidating change recorded for same-task recovery |

## Concerns, technical debt, or notes for the Tech Lead

- Required scope amendment: authorize `cmd/fid/integration/p23_concern_list_human_test.go` for Task 13 retained-assertion convergence. The accepted contract uniquely requires displaying its recorded summary and reasoning; no product decision or protected gate-definition change is needed.
- Headroom warning: the exact ordinary gate reported cmd/fid at 530.445 s, 88.4% of its 600 s budget, above the 480 s warning and leaving 69.555 s headroom. Other ordinary packages and tagged integration remained inside their budgets.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
