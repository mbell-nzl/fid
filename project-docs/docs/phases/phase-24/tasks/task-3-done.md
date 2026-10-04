# Phase 24 - Task 3: Input descriptor framework, help, templates and field-named validation - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2
- Next action: None.
- Human attention required: no

OUTCOME COMPLETE

## Summary of what was implemented

- Added an `internal/cli` descriptor model, Go type reflection check, closed structural validator, help renderer, and editable JSON template renderer. The validator enforces size, UTF-8, one object, duplicate keys, ordered unknown/member checks, formats and bounds, and field-path errors with sensitive-value redaction.
- Registered descriptors for `stop raise`, `stop resolve`, and `attention resolve`, with the other 36 live `--input` leaves in a disjoint pending fence. Connected the three leaves to descriptor validation and `--example` before workspace access. Preserved their domain validation after input decoding.
- Added a one-line `Short` to every previously bare leaf in scope and replaced the shared `Run JSON input` flag text. Updated the CLI design specification.
- Converged the retained Phase 22 `stop.raise` fixture with the two explicit nullable `resolution_condition` members, allowing its previously rejected raise to pass structural validation without changing journey assertions.
- Adapted structural error details to the accepted `invalid-argument` envelope schema: raw numeric validation values become strings at the error boundary, while missing or redacted values remain null. `attention resolve` now reports an integer-bound failure with exit 2 and no Event. No change to `stop.go` was needed.
- QA review-1 repair adds process corpora for the three adopted leaves: 18 ordered `--example` refusals and 42 structural failures with exact details, plus standalone template and well-typed controls.

## Files created or modified

- `internal/cli/input_descriptor.go`, `internal/cli/input_validate.go`, `internal/cli/input_render.go`, `internal/cli/input_validate_test.go`
- `cmd/fid/input_descriptors.go`, `cmd/fid/input_descriptor_command.go`, `cmd/fid/input_descriptor_test.go`
- `cmd/fid/stop.go`, `cmd/fid/attention.go`, `cmd/fid/run.go`
- `cmd/fid/concern.go`, `cmd/fid/decision.go`, `cmd/fid/deviation.go`, `cmd/fid/finding.go`, `cmd/fid/human_outcome.go`, `cmd/fid/requirement.go`, `cmd/fid/review.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/integration/p22_fake_test.go` (bounded TRG-24-001 fixture repair)
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-24/tasks/task-3-done.md`
- QA repair changed only `cmd/fid/input_descriptor_test.go` and this report; the retained `internal/cli/input_validate_test.go` supplies the sensitive-value and order controls.

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` records the three adopted help, template, exclusivity and field-path conventions.

## Validation Contract Coverage

- VAL-010: PASS for the current tree - `TestEveryLiveLeafHasOneLineShort` at `cmd/fid/input_descriptor_test.go:276`; the live inventory regeneration printed `0 .../no-short` among 92 leaves.
- VAL-011: PASS for Task 3's three leaves - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` at `cmd/fid/input_descriptor_test.go:234` proves registry/fence disjointness and union with all 39 live input leaves, reflects the three Go types, and rejects a seeded mismatch. The remaining 36 descriptors belong to later adoption steps in the same phase assertion.
- VAL-012: PASS for Task 3's three leaves - `TestAdoptedInputHelpAndExamples` at `cmd/fid/input_descriptor_test.go:291` checks every top-level field and template command; SHA-256 `b3a71078d99b3047a9ffdba35e599bd8298da2e6f6e9b06f709ba3283af5bf75` pins the full `stop raise --help` text. Nested field help is emitted by `internal/cli/input_render.go:16`.
- VAL-013: PASS for Task 3's three leaves - exact full template strings for all three are asserted in `TestAdoptedInputHelpAndExamples` at `cmd/fid/input_descriptor_test.go:291`; the real `fid stop raise --example` printed one indented object ending `}`.
- VAL-014: PASS for Task 3's three leaves - the same test decodes each exact template into its registered Go request type with unknown fields disallowed at `cmd/fid/input_descriptor_test.go:358`.
- VAL-015: PASS for all three adopted leaves - `TestAdoptedExampleProcessExclusivity` at `cmd/fid/input_descriptor_test.go:44` runs six ordered conflict combinations per leaf through a built CLI, checks exit 2 and the exact first Human conflict or `--json` envelope, and checks operation plus typed `{argument,value,allowed}` in the shared failure writer, no Event, and an empty non-project workspace. A standalone `--example` succeeds for each leaf with one JSON object, empty stderr, and no effects. The earlier exclusivity test remains intact.
- VAL-016: PASS for the three adopted leaves - `TestAdoptedStructuralProcessDetails` at `cmd/fid/input_descriptor_test.go:494` checks 42 explicit failures: every applicable top-level class per leaf, all seven classes in both `input.subject.*` and `input.resolution_condition.subject.*`, Attention array cases, and variant details. Each row compares exact typed path, value, and allowed list, then checks the built CLI `invalid-argument` envelope, exit 2, no Event field, and an untouched workspace. Three valid-structure controls reach domain validation. `TestInputValidationOrderAndSensitiveValues` at `internal/cli/input_validate_test.go:8` retains ordered and redacted-value controls.
- VAL-052: PASS - this continuation's exact ordinary gate reports `cmd/fid` 223.497 s and `internal/export` 98.884 s; the touched uncached packages report `cmd/fid` 195.175 s and `internal/cli` 1.116 s. The prior uncached `internal/workflow` 93.837 s and `internal/issues` 166.391 s are preserved. The exact tagged integration gate reports 298.059 s. All ordinary packages were below 600 s, tagged integration below 900 s, and no ordinary package reached the 480 s warning.
- VAL-059: PASS - `attention resolve --help` uses `Resolve an Attention item.` and `Resolve an Attention item JSON input file (or - for stdin)`; `runMutationFlags` now uses neutral `JSON input file` text at `cmd/fid/run.go:208`. The three non-Run adopted leaves use their own descriptor text.
- VAL-060: PASS for Task 3's three leaves - `TestAdoptedInputHelpAndExamples` at `cmd/fid/input_descriptor_test.go:291` runs each `--example` in an empty non-project directory, checks exit 0 and empty stderr, and confirms no directory entry was created.

## Tests run

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAdoptedExampleProcessExclusivity$'` - FAIL - `FAIL` (terminated `exit_code: 1`); initial edit did not compile because `os/exec` was unused.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAdoptedExampleProcessExclusivity$'` - FAIL - `FAIL` (terminated `exit_code: 1`); expected corpus red: 6 of 18 rows.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAdoptedExampleProcessExclusivity$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.449s` (terminated `exit_code: 0`); 18 conflict refusals and three standalone successes.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAdoptedStructuralProcessDetails$'` - FAIL - `FAIL` (terminated `exit_code: 1`); expected corpus red: zero structural rows.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAdoptedStructuralProcessDetails$'` - FAIL - `FAIL` (terminated `exit_code: 1`); array sorted-unique expectation used decoded values instead of retained `json.RawMessage` entries.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAdoptedStructuralProcessDetails$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.477s` (terminated `exit_code: 0`); corrected raw-value expectation.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAdopted(ExampleProcessExclusivity|StructuralProcessDetails)$'` - FAIL - `FAIL` (terminated `exit_code: 1`); well-typed controls reached domain validation, which also uses exit 2.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAdopted(ExampleProcessExclusivity|StructuralProcessDetails)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	3.186s` (terminated `exit_code: 0`); controls now assert the exact domain boundary.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.116s` (terminated `exit_code: 0`); `cmd/fid` also passed at 195.175 s.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAttentionResolveIntegerBoundWritesJSONFailure$'` - FAIL - `FAIL` (terminated `exit_code: 1`); first red observed exit 7 and empty stderr.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAttentionResolveIntegerBoundWritesJSONFailure$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.357s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.118s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	93.837s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid-p24-task3-a2 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `printf '%s' '{"attention_id":"ATT-1","expected_version":0,"fact_event_ids":["EVT-1"]}' | /tmp/fid-p24-task3-a2 --json attention resolve --input - --actor ACT-1 --idempotency-key bound-a2` - expected exit 2; real binary printed `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["integer range 1..unbounded"],"argument":"input.expected_version","value":"0"},"message":"invalid input.expected_version","retryable":false,"subject":null},"ok":false,"operation":"attention.resolve"}` with no Event (terminated `exit_code: 2`).
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24StopRaiseConflictNamesChangedDetails$' ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`); before the two null members, `input.resolution_condition.budget_name` was rejected.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP24StopRaiseConflictNamesChangedDetails$' ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.474s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`); package 248.184 s, later failure at `TestP23StopFailureCodeMatrix/attention/invalid-argument`.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestP23StopFailureCodeMatrix/attention/invalid-argument$' ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`); same empty JSON stderr in isolation.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAttentionResolveUnknownFieldNamesInputPath$'` - FAIL - `FAIL` (terminated `exit_code: 1`); expected first red: `argument:"input"` rather than `input.bogus`.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^TestAttentionResolveUnknownFieldNamesInputPath$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.368s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.112s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^Test(StructuralFailuresCoverEachAdoptedLeaf|AdoptedInputHelpAndExamples)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.370s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/cli -run '^TestInputValidationOrderAndSensitiveValues$'` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.014s` (terminated `exit_code: 0`)
- `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` - observed seven original decoders and no new decoder matching the pattern; last line `cmd/fid/activation.go:328:func decodeClosedActivationInput(stdin io.Reader, path string, target any) error {` (terminated `exit_code: 0`).
- Input-inventory section 1 regeneration block - observed `92` leaves, `39` input leaves, `0` no-Short leaves; last line `235 total` (terminated `exit_code: 0`).
- `go build -o /tmp/fid-p24-task3 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid-p24-task3 stop raise --example` from `/tmp` printed the exact template ending `}` with empty stderr (terminated `exit_code: 0`).
- `/tmp/fid-p24-task3 --json attention resolve --example --input /definitely-missing` - observed expected rejection: final JSON envelope named `--json`, operation `attention.resolve`, `invalid-argument` (terminated `exit_code: 2`).

## Exact quality-gate results

Working directory: `/home/mark/projects/fid`. The original lint failure and both tagged integration failures are retained as immutable attempts. The third five-command gate is recorded in order below them.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `* staticcheck: 1` (terminated `exit_code: 1`); first attempt.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`); `cmd/fid/integration` 190.815 s. The retained fixture omits required nullable condition members and lies outside the named task file scope.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`); `cmd/fid/integration` 290.647 s. `TestP23StopFailureCodeMatrix/attention/invalid-argument` received empty JSON stderr; product-path repair is outside this continuation's allowlist.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	302.158s` (terminated `exit_code: 0`)

QA review-1 repair gate, run in order from `/home/mark/projects/fid` on the test checkpoint commit:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	298.059s` (terminated `exit_code: 0`)

## Deviations from spec

- QA review-1 repair changed tests only. Its initial focused failures (`FAIL`, terminated `exit_code: 1`) reflected an unused import, incomplete corpus rows, a raw-message test expectation, and an overstrict process control. Removing the unused import, completing the corpora, retaining the raw-value expectation, and checking that valid controls reached the accepted domain boundary produced focused `ok` results (terminated `exit_code: 0`). No accepted command, assertion, or product behavior changed.
- No accepted input shape, domain message, gate definition, or assertion was changed. The first lint attempt failed with verbatim last token `* staticcheck: 1` (terminated `exit_code: 1`) because the new validator used `bytes.IndexAny` for a membership check. Replacing it with `bytes.ContainsAny` in `internal/cli/input_validate.go:218` made the byte-identical lint command return `0 issues.` (terminated `exit_code: 0`).
- The first tagged integration attempt failed with last token `FAIL` (terminated `exit_code: 1`) because the retained Phase 22 fixture omitted required nullable members. Under TRG-24-001, the one-line fixture change added `budget_name:nil` and `capability_id:nil`; the byte-identical focused test changed from `FAIL` (terminated `exit_code: 1`) to `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	2.474s` (terminated `exit_code: 0`). The later full tagged attempts reached a distinct failure and remain failed evidence.
- The second exact tagged integration attempt failed with last token `FAIL` (terminated `exit_code: 1`) after 290.647 s because a numeric bound value reached `cli.NewFailure` for `invalid-argument`, whose accepted detail schema permits only a nullable string. Under TRG-24-002, `InputError.Details()` now renders non-string JSON values as strings at the structural error boundary, preserving null redaction and raw `InputError.Value`. The focused test first failed with `FAIL` (terminated `exit_code: 1`), reporting exit 7 and empty stderr; after repair it returned `ok  	github.com/mbell-nzl/fid/cmd/fid	2.357s` (terminated `exit_code: 0`). The byte-identical tagged command later returned `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	302.158s` (terminated `exit_code: 0`). No accepted schema, input shape, domain message, assertion, or gate changed.

## Concerns, technical debt, or notes for the Tech Lead

- The retained `cmd/fid/integration/p22_fake_test.go:514` fixture line was introduced by `25caa7d40712455c3d1dc39dcdea77b4cc743dd1` (`Refs: phase-22/task-9`), before Phase 24. Task 3's closed validation displaced its omission; TRG-24-001 authorized exactly the two null members. All sequential journey assertions remain intact, and the downstream blank `stop_id`/`attention_id` failures from the first gate no longer occur.
- The retained Attention invalid-argument matrix row at `cmd/fid/integration/p23_human_failure_test.go:80` was introduced by `180f5168e3f1338cb758f8388f1c26fd6f8428d3` (`Refs: phase-23/task-19`). Task 3's structural integer-bound validation displaced the row's prior domain failure path. TRG-24-002 repaired only the structural error adapter in `internal/cli/input_validate.go:24`; the row and its JSON/Human parity, repeatability, and no-Event assertions remained unchanged. `cmd/fid/stop.go` did not need an edit.
- Retained exact production-path inventory `TestPhase22AttentionReadBoundary` initially failed when `input_descriptors.go` joined the Attention source set. Its expected list at `cmd/fid/contract_surface_test.go:78` was introduced by `59d2653da7cfe14699dc4d1cd85ef032bec26a5c` (`Refs: phase-23/task-14`). Added the new path to that exact list, preserving the inventory assertion; the focused test then passed.
- Inventory R18 now counts 25 typed flags versus its planning count of 22 because the three adopted leaves each add `--example`. The later grammar corpus should regenerate that live inventory.
- VAL-052 timing evidence: `cmd/fid` 200.072 s in the exact ordinary gate, `internal/workflow` 93.837 s and `internal/issues` 166.391 s in the uncached internal-package check, and tagged integration 302.158 s in the exact gate. No ordinary package crossed the 480 s warning.
- QA repair timing: exact ordinary `cmd/fid` 223.497 s, `internal/export` 98.884 s, and tagged integration 298.059 s. No package reached 80% of its budget.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
