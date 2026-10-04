# Phase 24 - Task 3: Input descriptor framework, help, templates and field-named validation - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: gate-repair
- Condition: The tagged integration gate rejects retained `stop.raise` fixtures that omit required nullable `resolution_condition.budget_name` and `capability_id`. The shared fixture is `cmd/fid/integration/p22_fake_test.go:514`, outside Task 3's named file scope.
- Recovery attempts: 1
- Next action: Phase Runner obtains bounded same-task scope for the retained integration fixture. Continue Task 3 by adding its required nullable members, preserve the sequential journey assertions, run focused red/green, then rerun all five exact gates on the final tree. Repair any later revealed in-scope steps under the same task.
- Human attention required: no

## Summary of what was implemented

- Added an `internal/cli` descriptor model, Go type reflection check, closed structural validator, help renderer, and editable JSON template renderer. The validator enforces size, UTF-8, one object, duplicate keys, ordered unknown/member checks, formats and bounds, and field-path errors with sensitive-value redaction.
- Registered descriptors for `stop raise`, `stop resolve`, and `attention resolve`, with the other 36 live `--input` leaves in a disjoint pending fence. Connected the three leaves to descriptor validation and `--example` before workspace access. Preserved their domain validation after input decoding.
- Added a one-line `Short` to every previously bare leaf in scope and replaced the shared `Run JSON input` flag text. Updated the CLI design specification.

## Files created or modified

- `internal/cli/input_descriptor.go`, `internal/cli/input_validate.go`, `internal/cli/input_render.go`, `internal/cli/input_validate_test.go`
- `cmd/fid/input_descriptors.go`, `cmd/fid/input_descriptor_command.go`, `cmd/fid/input_descriptor_test.go`
- `cmd/fid/stop.go`, `cmd/fid/attention.go`, `cmd/fid/run.go`
- `cmd/fid/concern.go`, `cmd/fid/decision.go`, `cmd/fid/deviation.go`, `cmd/fid/finding.go`, `cmd/fid/human_outcome.go`, `cmd/fid/requirement.go`, `cmd/fid/review.go`
- `cmd/fid/contract_surface_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-24/tasks/task-3-done.md`

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` records the three adopted help, template, exclusivity and field-path conventions.

## Validation Contract Coverage

- VAL-010: PASS for the current tree - `TestEveryLiveLeafHasOneLineShort` at `cmd/fid/input_descriptor_test.go:84`; the live inventory regeneration printed `0 .../no-short` among 92 leaves.
- VAL-011: PASS for Task 3's three leaves - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` at `cmd/fid/input_descriptor_test.go:42` proves registry/fence disjointness and union with all 39 live input leaves, reflects the three Go types, and rejects a seeded mismatch. The remaining 36 descriptors belong to later adoption steps in the same phase assertion.
- VAL-012: PASS for Task 3's three leaves - `TestAdoptedInputHelpAndExamples` at `cmd/fid/input_descriptor_test.go:99` checks every top-level field and template command; SHA-256 `b3a71078d99b3047a9ffdba35e599bd8298da2e6f6e9b06f709ba3283af5bf75` pins the full `stop raise --help` text. Nested field help is emitted by `internal/cli/input_render.go:16`.
- VAL-013: PASS for Task 3's three leaves - exact full template strings for all three are asserted in `TestAdoptedInputHelpAndExamples` at `cmd/fid/input_descriptor_test.go:99`; the real `fid stop raise --example` printed one indented object ending `}`.
- VAL-014: PASS for Task 3's three leaves - the same test decodes each exact template into its registered Go request type with unknown fields disallowed at `cmd/fid/input_descriptor_test.go:163`.
- VAL-015: PASS for Task 3's rule - `TestExampleRejectsOtherInputsInFixedOrder` at `cmd/fid/input_descriptor_test.go:180` covers all six conflicts and exit 2; the real `--json attention resolve --example --input /definitely-missing` returned an `invalid-argument` envelope for `--json` first, with `value:null`, `allowed:["--example"]`, operation `attention.resolve`, and exit 2. `cmd/fid/input_descriptor_command.go:47` constructs the ordered details before any workspace read.
- VAL-016: PASS for Task 3's adopted leaves, pending integration fixture convergence - `TestAdoptedStructuralValidationPaths` at `cmd/fid/input_descriptor_test.go:218`, `TestStructuralFailuresCoverEachAdoptedLeaf` at `:248`, and `TestInputValidationOrderAndSensitiveValues` at `internal/cli/input_validate_test.go:8` cover top-level and nested paths, duplicate/unknown/missing/null/type/enum/format/bound failures and redaction. The tagged integration gate exposed an out-of-scope retained fixture whose omitted required nullable members are now rejected; domain failures after a valid structure remain in the existing flow.
- VAL-052: BLOCKED - `cmd/fid` was 196.604 s, below the 480 s warning; tagged integration was 190.815 s, below 900 s but failed. `internal/workflow` and other ordinary packages were cached in the exact ordinary gate, so their current wall times cannot yet be asserted. Recheck on the completed final tree.
- VAL-059: PASS - `attention resolve --help` uses `Resolve an Attention item.` and `Resolve an Attention item JSON input file (or - for stdin)`; `runMutationFlags` now uses neutral `JSON input file` text at `cmd/fid/run.go:208`. The three non-Run adopted leaves use their own descriptor text.
- VAL-060: PASS for Task 3's three leaves - `TestAdoptedInputHelpAndExamples` at `cmd/fid/input_descriptor_test.go:99` runs each `--example` in an empty non-project directory, checks exit 0 and empty stderr, and confirms no directory entry was created.

## Tests run

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

Working directory: `/home/mark/projects/fid`. The first lint attempt is retained, followed by the repaired attempt. The tagged integration result is an immutable failed attempt.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - FAIL - `* staticcheck: 1` (terminated `exit_code: 1`); first attempt.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.3%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - `FAIL` (terminated `exit_code: 1`); `cmd/fid/integration` 190.815 s. The retained fixture omits required nullable condition members and lies outside the named task file scope.

## Deviations from spec

- No accepted input shape, domain message, gate definition, or assertion was changed. The first lint attempt failed with verbatim last token `* staticcheck: 1` (terminated `exit_code: 1`) because the new validator used `bytes.IndexAny` for a membership check. Replacing it with `bytes.ContainsAny` in `internal/cli/input_validate.go:218` made the byte-identical lint command return `0 issues.` (terminated `exit_code: 0`).
- The tagged integration attempt failed with last token `FAIL` (terminated `exit_code: 1`). The accepted required-nullable validation newly exposes an omitted field in the retained Phase 22 fixture. No out-of-scope fixture edit or permissive input-shape change was made. The same exact integration command must be rerun after authorized convergence.

## Concerns, technical debt, or notes for the Tech Lead

- Scope gap: `cmd/fid/integration/p22_fake_test.go:514` builds the `stop.raise` input without `resolution_condition.budget_name` or `capability_id`. That line was introduced by `25caa7d40712455c3d1dc39dcdea77b4cc743dd1` (`Refs: phase-22/task-9`), before Phase 24. Task 3's accepted closed structural validation displaces its omission; the fixture feeds Phase 22 and 24 sequential journeys and needs same-task convergence without weakening them. The downstream blank `stop_id`/`attention_id` failures are consequences of the earlier rejected raise.
- Retained exact production-path inventory `TestPhase22AttentionReadBoundary` initially failed when `input_descriptors.go` joined the Attention source set. Its expected list at `cmd/fid/contract_surface_test.go:78` was introduced by `59d2653da7cfe14699dc4d1cd85ef032bec26a5c` (`Refs: phase-23/task-14`). Added the new path to that exact list, preserving the inventory assertion; the focused test then passed.
- Inventory R18 now counts 25 typed flags versus its planning count of 22 because the three adopted leaves each add `--example`. The later grammar corpus should regenerate that live inventory.
- The exact ordinary gate cached `internal/workflow` and every other ordinary package; their current wall times remain unverified for VAL-052. `cmd/fid` has over 400 s of headroom, and no observed package crossed the 480 s warning.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `04b6b0992fe20a02df402a183269a9847ddae5ab`

## Commit subject

- `feat(cli): checkpoint input descriptors and validation`
