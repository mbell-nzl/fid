# Phase 24 - Task 9: Activation descriptors; whole-tree input conformance - Recovery Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: Seven activation leaves are adopted, but activation.component requires shared descriptor changes outside Task 9's named source-file scope. The preserved reflection regression and whole-tree registry assertion are red. Whole-tree corpora, fake/fid assembly and three decoder removals remain unfinished.
- Recovery attempts: 0
- Next action: Obtain a Task 9 scope amendment for `internal/cli/input_descriptor.go`, `internal/cli/input_render.go`, `internal/cli/input_validate.go`, and `internal/cli/input_validate_test.go`; dispatch a same-task continuation from this checkpoint. Implement the accepted component_type discriminator and raw typed array variants, finish the remaining accepted criteria, and rerun the byte-identical required gates.
- Human attention required: no

OUTCOME RECOVERY_REQUIRED

## Summary of what was implemented

- Added descriptors for activation prepare, commit, abort, finalize, resource, recover and recover-complete. The descriptors retain the accepted Phase 21 shapes and Phase 22 version/resource additions, with help, required-only examples, typed decoding, bounds and closed structural validation.
- Moved transition and recovery decoding ahead of workspace resolution. Structural failures now retain their exact field path through the activation failure mapper, without workspace writes.
- Added a built-process corpus for these seven leaves, strengthened the live registry assertion to require every live input leaf, and preserved a failing regression for the actual ActivationComponentRequest raw array.
- Added the top-level sensitive-field guard over the current registry (38 registered leaves). Removed the judgment-family decodeClosedJudgmentInput wrapper and made its callers use the existing shared descriptor reader directly.
- Preserved the component fence entry and its decoder until component adoption can be completed. The review and envelope decoders also remain. No whole-tree completion or fake/fid assembly is claimed.

## Files created or modified

- `cmd/fid/input_descriptors_activation.go` (new)
- `cmd/fid/activation.go`
- `cmd/fid/input_descriptor_test.go`
- `cmd/fid/concern.go`
- `cmd/fid/decision.go`
- `cmd/fid/deviation.go`
- `cmd/fid/human_outcome.go`
- `cmd/fid/requirement.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-24/tasks/task-9-done.md`

The existing runner-owned changes in run-ledger.md and triage-log.md were neither edited nor staged.

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: documents the seven adopted activation leaves and explicitly identifies this partial checkpoint and pending component adoption.

## Validation Contract Coverage

- VAL-011: BLOCKED - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes`, `cmd/fid/input_descriptor_test.go:237`, reports the one remaining unregistered live leaf, activation component. `TestP24ActivationComponentRawArrayReflection`, line 293, fails with `input.content[]: Go type uint8 is not an object`. The seeded mismatch assertion remains intact.
- VAL-012: BLOCKED - `TestP24ActivationSimpleDescriptorProcessCorpus`, line 279, passes descriptor-derived help field/type/requirement/meaning/constraint checks for all seven adopted activation leaves. Component help, both variants and its complete pinned text remain pending.
- VAL-013: BLOCKED - the same process corpus passes built-binary example equality against the documented-order, two-space-indented descriptor template for seven activation leaves. The 39-leaf form corpus remains pending.
- VAL-014: BLOCKED - the same corpus decodes each of those seven templates into its actual Go request type with unknown fields disallowed, then validates qualifying expanded input. Component raw-array type conformance and full-tree assembly remain pending.
- VAL-015: BLOCKED - whole-tree example exclusivity has not been assembled; the retained exclusivity corpus still lists the earlier 31 adopted leaves. No narrower success is substituted for this assertion.
- VAL-016: BLOCKED - the seven-leaf process corpus passes unknown, duplicate, missing, null, type, declared enum/format/bounds and nested resource/version failures, with exact typed details and unchanged empty workspaces. Both component variants, mismatched-variant paths and sensitive-field value:null cases remain pending.
- VAL-035: BLOCKED - Phase 22 fake/fid files were inspected but no nine-leaf converged-input conformance assembly was implemented in this checkpoint. The passing tagged gate is not asserted to prove the missing assembly.
- VAL-060: BLOCKED - the seven-leaf process corpus checks successful examples from empty non-project directories and verifies they remain empty; whole-tree side-effect coverage remains pending.
- VAL-052: BLOCKED - required ordinary gate: cmd/fid 282.333 s (failed assertions); slowest other **uncached** ordinary package internal/export 103.313 s. Required tagged integration: 333.935 s. Required gate reported internal/workflow and the other ordinary test packages as `(cached)`, without wall times. Supplemental fresh internal/workflow race run: 55.650 s. No observed sample reaches 480 s ordinary or 720 s tagged, but cached packages cannot supply current-session timing evidence or establish the slowest package across the entire ordinary suite.
- I-17: partial evidence - `TestP24InputDescriptorsHaveNoTopLevelSensitiveFields`, line 309, passes across the 38 registered descriptors. Component adoption must retain the guard and extend its coverage to all 39.

## Tests run

Working directory is `/home/mark/projects/fid` unless stated otherwise. Every termination below was observed from the nested command result.

- `go test ./cmd/fid -run '^(TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestP24ActivationSimpleDescriptorProcessCorpus|TestP24ActivationComponentRawArrayReflection)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). First red before product edits: all eight activation descriptors missing; actual raw content reflection reports uint8 instead of an object.
- `go test ./cmd/fid -run '^TestP24ActivationSimpleDescriptorProcessCorpus$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`). First implementation exposed workspace resolution preceding structural decoding, plus an element mutation violating the five-resource array length before reaching the element.
- `go test ./cmd/fid -run '^TestP24ActivationSimpleDescriptorProcessCorpus$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.598s` (terminated `exit_code: 0`)
  After decoding-order and qualifying-corpus repairs.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 ./cmd/fid -run '^(TestP24ActivationSimpleDescriptorProcessCorpus|TestP24InputDescriptorsHaveNoTopLevelSensitiveFields|TestP24JudgmentStructuralValidationCorpus|TestP24BindingDescriptorProcessCorpus|TestP24RunDescriptorProcessCorpus)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	9.418s` (terminated `exit_code: 0`)
  Includes the judgment callers after wrapper removal and retained binding/Run structural assertions.
- `gofmt -w cmd/fid/activation.go cmd/fid/input_descriptors_activation.go cmd/fid/input_descriptor_test.go cmd/fid/concern.go cmd/fid/decision.go cmd/fid/human_outcome.go cmd/fid/requirement.go cmd/fid/deviation.go` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `goimports -w cmd/fid/activation.go cmd/fid/input_descriptors_activation.go cmd/fid/input_descriptor_test.go cmd/fid/concern.go cmd/fid/decision.go cmd/fid/human_outcome.go cmd/fid/requirement.go cmd/fid/deviation.go` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `go build -o /tmp/fid-p24-task9-uoK52u/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid-p24-task9-uoK52u/fid activation prepare --example` - PASS - `}` (terminated `exit_code: 0`)
  Working directory `/tmp/fid-p24-task9-uoK52u`, a fresh non-project directory. Printed the required fields in descriptor order; stderr was empty.
- `/tmp/fid-p24-task9-uoK52u/fid --json activation abort --input /tmp/fid-p24-task9-uoK52u/rejected.json --idempotency-key task9-entrypoint` - observed expected structural rejection (terminated `exit_code: 2`); working directory `/tmp/fid-p24-task9-uoK52u`. Terminal output was exactly `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["activation_id","expected_status","expected_version","reason"],"argument":"input.unexpected","value":"unexpected"},"message":"invalid input.unexpected","retryable":false,"subject":null},"ok":false,"operation":"activation.abort"}`. No Event envelope and no stdout.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	55.650s` (terminated `exit_code: 0`)
  Supplemental fresh timing after the required gate's cached result; does not replace the required ordinary gate.
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-9-done.md project-docs/docs/phases/phase-24/tasks/task-9.md` - FAIL - `FAIL report-lint found 11 failure(s)` (terminated `exit_code: 11`). Report-only recovery attempt 1 reproduced the Runner's formatting failure.

## Exact quality-gate results

The following six commands ran exactly in task-spec order on the checkpoint source tree. The ordinary failure was preserved and the remaining commands still ran. No source or test edits followed this gate attempt.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Immutable failed attempt: exactly the missing component registry entry and raw-array reflection regression failed. cmd/fid 282.333 s, coverage 73.3%; no retained test failed. Shared-source scope is required to repair the remaining reds.
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.1%` (terminated `exit_code: 0`)
  This is a report from the failed ordinary attempt, not evidence that the ordinary gate passed.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	333.935s` (terminated `exit_code: 0`)
- `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` - BLOCKED - `cmd/fid/envelope.go:178:func decodeClosedEnvelopeInput(stdin io.Reader, path string, target any) error {` (expected `exit_code: 1`; terminated `exit_code: 0`). Three matches remain: review.go:135, activation.go:326 and envelope.go:178. The inventory is not an absence proof.

Report-lint is the Runner's control under DEC-158 and was not run by the Developer at the implementation checkpoint. This report-only recovery explicitly authorizes the exact report-lint verification recorded above.

## Deviations from spec

- Report-only recovery attempt 1 reproduced `FAIL report-lint found 11 failure(s)` (terminated `exit_code: 11`). Root cause: the 11 historical PASS rows had punctuation or explanatory text after their termination fields, but report-lint requires those fields at the end of each evidence line. Repair: remove only trailing punctuation and move explanations to continuation lines, preserving every historical command, status, token, observed termination, provenance, scope finding and the typed RECOVERY_REQUIRED outcome.
- Task 9 is incomplete. The empty-fence removal, component descriptor, both discriminated variants, complete whole-tree corpora, fake/fid assembly, and remaining decoder removals are deferred to the same task's continuation, with their failures retained. No accepted shape, enum, threshold, gate command, status or assertion semantics changed.
- Required ordinary attempt failed with token `FAIL` (terminated `exit_code: 1`). Root cause: the accepted component array cannot be reflected by the current shared descriptor model (`input.content[]: Go type uint8 is not an object`), and the component is consequently unregistered. The necessary shared implementation paths are outside the assigned named file scope. No shared-source repair or later passing attempt is claimed.
- Required R5 attempt ended with token `cmd/fid/envelope.go:178:func decodeClosedEnvelopeInput(stdin io.Reader, path string, target any) error {` (terminated `exit_code: 0`, expected 1). Root cause: component still uses its old decoder, and two retained sparse-input decoder fixtures have not yet been converged. The judgment wrapper was removed, but no passing absence rerun is claimed.
- The touched-package pre-gate tier could not be green while the scope-bound component regressions remained red. A focused race tier passed and the exact full gate was still executed; it is explicitly failed, not substituted or used to close this task.
- Retained fixture convergence follows the introducing-commit provenance below. These changes enable the accepted Task 9 length-two commit-digest and length-five resource shapes while preserving each original failure subject and assertion strength.

| Retained assertion or fixture row | Introducing commit | Phase / task | Convergence in this checkpoint |
|---|---|---|---|
| Live input leaf accounting in TestInputDescriptorRegistryMatchesLiveTreeAndTypes | `04b6b0992fe20a02df402a183269a9847ddae5ab` | phase-24/task-3 | Strengthened the completion check to require a registered descriptor for every live leaf; kept reflection, overlap and seeded-mismatch checks. Component stays red. |
| Qualifying object-array expansion in p24ExpandDescriptorObjects | `1baaecc7983872b06e2a47e6df0801bfb1370ea0` | phase-24/task-6 | Creates at least the declared minimum number of objects; retains recursive qualifying fields and validation. |
| Qualifying digest fixture generation in p24ExpandDescriptorObjects | `1baaecc7983872b06e2a47e6df0801bfb1370ea0` | phase-24/task-6 | Extends the retained scalar-digest fixture generator to array elements, supplying valid distinct sorted digests at the declared length so accepted commit input reaches each structural check. |
| Array element-type structural row in p24DescriptorProcessCorpus | `1baaecc7983872b06e2a47e6df0801bfb1370ea0` | phase-24/task-6 | Replaces the first element while preserving accepted array length; exact input.path[0], typed details and no-write checks remain. |
| Array element-format structural row in p24DescriptorProcessCorpus | `1baaecc7983872b06e2a47e6df0801bfb1370ea0` | phase-24/task-6 | Replaces the first element with the same invalid format value while preserving accepted array length and exact path/details. |
| Sorted-unique structural row in p24DescriptorProcessCorpus | `1baaecc7983872b06e2a47e6df0801bfb1370ea0` | phase-24/task-6 | Duplicates an otherwise valid element at the accepted length, retaining the array-path ordering failure and no-write assertions. |

Provenance was obtained with `git log -S` for the live-accounting text and helper/row introductions, then verified with `git show -s` including each commit's Refs footer. No introducing phase/task was inferred from a filename or historical test label.

## Concerns, technical debt, or notes for the Tech Lead

- Scope gap: `internal/cli/input_descriptor.go:88` treats raw content as a slice of uint8 and recursively requires an object element. `internal/cli/input_validate.go:99` and `:181` select object variants only through root `code`; `internal/cli/input_render.go:52` labels variants as code and `:109` selects example variants through code. These shared files and their regression tests require ordinary Task 9 source scope, not a tools-only control-conformance directive. The accepted component_type shapes need no new product or architecture decision; the missing authority is the file scope.
- Continue with all eight activation descriptors, both typed content variants and exact mismatch paths; remove the final pending fence while retaining registry=live. Complete all six 39-leaf corpora, the sensitive nested failure cases, component help pin and fake-dispatch-derived nine-leaf assembly. Remove all three remaining family decoders without retaining another structural source.
- Pending decoder-fixture provenance was inspected but no convergence is claimed: TestDecodeClosedReviewInputRejectsMalformedInapplicability was introduced by `0be596e9068795ceb3a0a89bdac50d538f8e8791` (phase-23/task-5); TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState by `2339dc9ebd8d0e6b5502ca32e1562892cf704ce6` (phase-22/task-12). Preserve valid empty optional-array acceptance and both output-only-state rejection states when converging them to descriptor validation.
- Observed timing samples are below their warning thresholds. Cached package results do not establish fresh wall times; complete the remaining VAL-052 timing evidence in the continuation.
- No grammar/unknown-command implementation (Task 9a), Human failure-line implementation (Task 10), runner control, ledger or triage edits were made.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `785a22491f1ba07bb748d1d4fb85f88ea1a8d45d`

## Commit subject

- `feat(cli): checkpoint activation input descriptors`
