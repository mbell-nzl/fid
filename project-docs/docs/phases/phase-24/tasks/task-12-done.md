# Phase 24 - Task 12: Concern, Deviation and Human Outcome get leaves - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added `concern get <concern-id>`, `deviation get <deviation-id>` and `human-outcome get <outcome-id>`, with one positional ID and global flags only.
- Reads ensure projection freshness and reuse each list's item loader. JSON returns the matching record under `concern`, `deviation` or `outcome`, with `event:null` and exit 0.
- Unknown IDs return the existing not-found codes, requested subject ID and `{id}` details with exit 3. Grammar failures use the accepted invalid-argument envelope and exit 2. Human output uses the existing item renderer.
- Added real-process record/list/get equality and failure pairing tests. Converged retained inventories and unknown-command alternatives to DEC-274 while preserving exact comparisons and rejection/no-journal-change checks.
- Capability discovery rows and activated-action allowlists are unchanged. Full reasoning and section disclaimers remain Task 13's scope.

## Files created or modified

- `cmd/fid/concern.go`
- `cmd/fid/deviation.go`
- `cmd/fid/human_outcome.go`
- `internal/review/concern_list.go`
- `internal/review/deviation.go`
- `internal/review/human_outcome.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/grammar_corpus_test.go`
- `cmd/fid/integration/p24_judgment_get_test.go`
- `cmd/fid/integration/exit_contract_test.go`
- `cmd/fid/integration/testdata/p23_conformance/product-surface.json`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-23/public-interface.md`
- `project-docs/docs/phases/phase-24/tasks/task-12-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: documented the leaves, operation IDs, result members, global flags, not-found and grammar behavior.
- `project-docs/docs/phases/phase-23/public-interface.md` section 5.1: synchronized the accepted DEC-274 additions to the three existing not-found producer rows.
- This completion report records both exact gate attempts, failed focused/probe attempts and convergence provenance. The runner-owned run ledger was not edited, staged or committed.

## Validation Contract Coverage

- VAL-042: PASS — `TestP24JudgmentGetRecords`, `cmd/fid/integration/p24_judgment_get_test.go:14`, records all three entities through the real CLI, compares each get record with its list item, checks the sole result member, operation, explicit `event:null`, exit 0 and empty stderr. It also checks Human item output, project selection from outside the workspace and unchanged journal events. Passed in the final exact tagged gate.
- VAL-043: PASS — `TestP24JudgmentGetRejections`, `cmd/fid/integration/p24_judgment_get_test.go:66`, asserts each unknown-ID code, exit 3, category, retryability, exact subject and `{id}` details in JSON and matching Human failures. Missing ID, extra argument and unknown flag are checked in both modes inside and outside a project with exit 2, exact grammar details, null subject and no journal changes. Passed in the final exact tagged gate.
- VAL-052: PASS — the exact gates below terminated successfully on the final tree. Recorded package times are below the accepted 600 s ordinary / 900 s tagged budgets. The earlier cmd/fid warning is retained in Tech Lead notes.

| Package | First exact attempt | Final exact attempt | Budget |
|---|---|---|---|
| cmd/fid | 484.599 s | 376.609 s | 600 s |
| internal/workflow | 142.144 s | cached; no new duration reported | 600 s |
| internal/export (slowest other ordinary package with a reported duration) | 183.669 s | 122.143 s | 600 s |
| cmd/fid/integration | 370.826 s, failed retained assertions | 381.138 s, PASS | 900 s |

Other timed ordinary packages in the first attempt were internal/db 71.210 s, internal/project 2.289 s and internal/review 1.262 s. Unchanged packages reported cached results; no duration was reconstructed for them.

- VAL-010 retained coverage: `TestEveryLiveLeafHasOneLineShort` at `cmd/fid/input_descriptor_test.go:2144` and `TestLiveLeafCommandsMatchContractInventories` at `cmd/fid/contract_surface_test.go:128` pass in the exact ordinary gate with all three added descriptions and transport registrations.
- VAL-064 / VAL-036 retained coverage: `TestGrammarWholeTreeClasses` and `TestGrammarBuiltBinaryCorpus` at `cmd/fid/grammar_corpus_test.go:283` and `:523` derive rows from the complete 93-leaf tree and check both output modes without a new exclusion. `TestP24JudgmentGetRejections` adds explicit pairing for the three read producers.

## Tests run

Working directory is `/home/mark/projects/fid` unless specified below. All Go test results were watched through observed termination.

Focused get attempt 1: test scaffolding did not compile (undefined helper).
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24JudgmentGet' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

Focused get attempt 2: expected behavioral red; all three leaves returned operation-unavailable.
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24JudgmentGet' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

Focused get attempt 3: readback passed; the new test incorrectly expected the ID instead of the first extra argument.
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24JudgmentGet' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

Focused get attempt 4: corrected the expectation from Task 9a's accepted grammar class table.
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24JudgmentGet' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.242s` (terminated `exit_code: 0`)

Retained live inventory red, then green after adding the three accepted leaves and changing the exact count to 93:
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestLiveLeafCommandsMatchContractInventories|TestGrammarWholeTreeClasses)$' -count=1 ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestLiveLeafCommandsMatchContractInventories|TestGrammarWholeTreeClasses)$' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	8.479s` (terminated `exit_code: 0`)

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/review ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	386.064s` (terminated `exit_code: 0`)
- That touched-package run also reported `ok  	github.com/mbell-nzl/fid/internal/review	1.117s`.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24JudgmentGet' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.598s` (terminated `exit_code: 0`)

After the first required integration failure, the two displaced retained assertions reproduced red and then passed after causal convergence:
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP23CancelledFamiliesAbsentFromProduct|TestV01ClosureJourneyAndRebuildParity)$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
- `go test -tags integration -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP23CancelledFamiliesAbsentFromProduct|TestV01ClosureJourneyAndRebuildParity)$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.779s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^(TestP24JudgmentGetRecords|TestP24JudgmentGetRejections|TestP23CancelledFamiliesAbsentFromProduct|TestV01ClosureJourneyAndRebuildParity|TestProcessExitContract)$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.246s` (terminated `exit_code: 0`)

- `rg -n "func TestLiveLeafCommandsMatchContractInventories" cmd/fid` - PASS - `cmd/fid/contract_surface_test.go:128:func TestLiveLeafCommandsMatchContractInventories(t *testing.T) {` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

Real entrypoint checks used the newly built `/tmp/fid`, from `/tmp/fid-task12-entrypoint.WJPS14`, freshly initialized with `/tmp/fid --json init --prefix GET` (observed exit_code: 0). These are observed contract refusals, not non-zero quality-gate PASS claims:
- `/tmp/fid --json concern get ID-missing` - OBSERVED - `{"contract_version":"fid-cli.v1","error":{"category":"not-found","code":"concern-not-found","details":{"id":"ID-missing"},"message":"Concern was not found","retryable":false,"subject":{"id":"ID-missing","type":"concern"}},"ok":false,"operation":"concern.get"}` (terminated `exit_code: 3`)
- `/tmp/fid --json deviation get ID-missing` - OBSERVED - `{"contract_version":"fid-cli.v1","error":{"category":"not-found","code":"deviation-not-found","details":{"id":"ID-missing"},"message":"Deviation was not found","retryable":false,"subject":{"id":"ID-missing","type":"deviation"}},"ok":false,"operation":"deviation.get"}` (terminated `exit_code: 3`)
- `/tmp/fid --json human-outcome get ID-missing` - OBSERVED - `{"contract_version":"fid-cli.v1","error":{"category":"not-found","code":"outcome-not-found","details":{"id":"ID-missing"},"message":"Human Outcome was not found","retryable":false,"subject":{"id":"ID-missing","type":"human-outcome"}},"ok":false,"operation":"human-outcome.get"}` (terminated `exit_code: 3`)

Earlier exploratory fixture probes failed and are preserved here; neither is acceptance evidence. Working directory: repository root.
- `python3 -` - FAIL - `sqlite3.DatabaseError: malformed database schema (issues_parent_kind_guard_insert) - near "CASE": syntax error` (terminated `exit_code: 1`)
- The Python stdin diagnostic attempted to enumerate Concern, Deviation and Human Outcome IDs from the cached Human fixture databases using Python's SQLite adapter.
- `/tmp/fid --json --project /home/mark/projects/fid/.cache/phase-24-human/fixtures-2268063598/clean-same concern list` - FAIL - `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["accepted concern.list input"],"argument":"input","value":null},"message":"SQL logic error: no such column: subject_type (1)","retryable":false,"subject":null},"ok":false,"operation":"concern.list"}` (terminated `exit_code: 2`)

- `git diff --cached --check` - FAIL - `project-docs/docs/phases/phase-24/tasks/task-12-done.md:158: new blank line at EOF.` (terminated `exit_code: 2`)

- `git add project-docs/docs/phases/phase-24/tasks/task-12-done.md` - FAIL - `fatal: Unable to create '/home/mark/projects/fid/.git/index.lock': Read-only file system` (terminated `exit_code: 128`)
- `git diff --cached --check` - FAIL - `project-docs/docs/phases/phase-24/tasks/task-12-done.md:158: new blank line at EOF.` (terminated `exit_code: 2`)
- `git add project-docs/docs/phases/phase-24/tasks/task-12-done.md` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `git diff --cached --check` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Exact quality-gate results

Every required command ran in the task's order, from `/home/mark/projects/fid`.

**Attempt 1 — before the final retained-test convergence**

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	142.144s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)

**Attempt 2 — final tree after causal convergence**

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	381.138s` (terminated `exit_code: 0`)

- Conditional commands: None.

## Deviations from spec

- No product or accepted contract deviation.
- Required integration attempt 1 failed with final token `FAIL` (terminated `exit_code: 1`), reporting `TestP23CancelledFamiliesAbsentFromProduct` and `TestV01ClosureJourneyAndRebuildParity`. Task 12's accepted DEC-274 leaves invalidated the pinned product command paths and Concern unknown-command alternatives. The causal repair added exactly three declared paths to the fixture and `get` to the Concern alternatives; cancelled-family predicates, exact envelope comparisons and journal immutability assertions were preserved. The focused byte-identical command progressed from `FAIL` (terminated `exit_code: 1`) to `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.779s` (terminated `exit_code: 0`). The full byte-identical integration command then passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	381.138s` (terminated `exit_code: 0`). Both full gate attempts remain above.
- Focused get attempts preserve three `FAIL` results (each terminated `exit_code: 1`): an undefined test helper, the expected missing-leaf red, and an incorrect new-test extra-argument expectation. The local process-exit helper fixed compilation; minimal read/CLI implementations fixed the behavioral red; Task 9a's class table determined `value:"extra"`. The identical focused command then passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.242s` (terminated `exit_code: 0`).
- The retained live inventory command first returned `FAIL` (terminated `exit_code: 1`): the three added leaves displaced the exact old inventory/count. Adding the accepted leaves and changing 90 to 93 preserved the complete equality check. The identical command passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	8.479s` (terminated `exit_code: 0`).
- Exploratory cached-fixture probes failed: Python's adapter rejected a trigger schema, and the rebuilt CLI reported `no such column: subject_type`. These inputs were unsuitable for current read acceptance. Fresh initialization and current CLI-created records supplied the later evidence; no fixture or product semantics were changed to accommodate those probes.

- The report draft's staged whitespace check failed with `project-docs/docs/phases/phase-24/tasks/task-12-done.md:158: new blank line at EOF.` (terminated `exit_code: 2`). Removing the extra EOF blank line fixed the file. The initial restage was blocked by the sandbox's read-only Git index (`exit_code: 128`), so the second check still saw the old staged draft and failed with the same token and exit 2. Staging through the approved Git escalation succeeded; the byte-identical whitespace check then passed with `exit 0 (no output)` (terminated `exit_code: 0`). These were report/index repairs only.

Retained-test convergence provenance, established with git log/blame rather than task labels:

| Displaced assertion | Introducing commit and phase/task | Accepted Task 12 convergence |
|---|---|---|
| Live leaf and root-flag equality | `37b59d62`, phase-20/task-18; the 90-leaf retained count was introduced by `2d5e9d55`, phase-24/task-9a | Add the three DEC-274 leaves; keep exact equality and the two diagnostic exclusions |
| Whole-tree grammar count / coverage | `2d5e9d55`, phase-24/task-9a | 93 leaves, still 29 groups; retain enumeration, all grammar classes and both-mode checks |
| Product command-surface fixture used by TestP23CancelledFamiliesAbsentFromProduct | `33351494`, phase-23/task-18 | Add only concern get, deviation get and human-outcome get; all table names and cancelled-family rejection predicates unchanged |
| Concern removed-command alternatives in assertRemovedCommandFailure | `2d5e9d55`, phase-24/task-9a; the closure journey originated in `7119d5b3`, phase-17/task-9 | Add get to the exact sorted allowed list; raise/address remain rejected in both modes with no journal append |

## Concerns, technical debt, or notes for the Tech Lead

- Headroom warning: the first exact ordinary gate reported cmd/fid at 484.599 s, 80.8% of the 600 s budget, above the 480 s warning and leaving 115.401 s headroom. The final attempt reported 376.609 s. Both observations are retained; no budget or command was changed.
- Final tagged integration time is 381.138 s, below the 900 s budget. internal/workflow's observed 142.144 s belongs to the first exact attempt; the second returned cached success, not a new timing sample.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `815a6b66af0381adee77077e826657511ac62292`

## Commit subject

- `feat(cli): add judgment get leaves`
