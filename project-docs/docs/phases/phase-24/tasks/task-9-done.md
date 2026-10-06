# Phase 24 - Task 9: Activation descriptors; whole-tree input conformance - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2
- Next action: None.
- Human attention required: no

OUTCOME COMPLETE

## Summary of what was implemented

- Completed all eight activation descriptors without changing accepted input shapes. `activation component` retains raw typed ActorProfile[] / GateDefinition[] content selected by component_type.
- Extended shared reflection, duplicate scanning, closed validation, help and deterministic examples for those selected raw arrays. Reflection checks both actual element structs; foreign-variant members reject at input.content[i]. Nested sensitive values remain redacted. Accepted Unicode scalar bounds and cwd/artifact path rules are scoped to the component contract; existing code-selected object callers preserve their behavior.
- Registered the last descriptor and removed the pending fence while retaining registry=live, reflection and seeded-mismatch assertions. Whole-tree help, typed templates, example form/order, exclusivity, structural validation and no-write checks derive all 39 live leaves; the variant corpus includes both component arrays and all ten Stop object variants.
- Assembled fake/fid input-admission parity over the exact nine-leaf intersection with the 17 changed leaves. Removed the remaining activation, Review and Envelope family decoder wrappers; retained output-only budget state rejection and converged displaced decoder fixtures to qualifying shared-reader input.
- Implementation checkpoints before long commands: `9d654e976792836cbdf7de8c114d0831b73f721e` (`feat(cli): complete activation descriptor conformance`) and final source checkpoint `69589747ca77a1a3dae7e4d6083229d0193ed52f` (`fix(cli): preserve accepted component string and path bounds`). Both contain Refs: phase-24/task-9. Historical failed implementation/report attempts remain recorded below.

## Files created or modified

Task-wide file inventory (including the earlier checkpoint and this continuation):

- `cmd/fid/input_descriptors_activation.go` (created)
- `cmd/fid/activation.go`
- `cmd/fid/input_descriptor_test.go`
- `cmd/fid/input_descriptors.go`
- `cmd/fid/input_descriptors_binding.go`
- `cmd/fid/input_descriptors_evidence.go`
- `cmd/fid/input_descriptors_judgment.go`
- `cmd/fid/input_descriptors_review.go`
- `cmd/fid/input_descriptors_run.go`
- `cmd/fid/review.go`
- `cmd/fid/finding_test.go`
- `cmd/fid/envelope.go`
- `cmd/fid/envelope_test.go`
- `cmd/fid/concern.go`
- `cmd/fid/decision.go`
- `cmd/fid/deviation.go`
- `cmd/fid/human_outcome.go`
- `cmd/fid/requirement.go`
- `cmd/fid/integration/p22_fake_test.go`
- `cmd/fid/integration/p22_semantic_fake_test.go`
- `internal/cli/input_descriptor.go`
- `internal/cli/input_render.go`
- `internal/cli/input_validate.go`
- `internal/cli/input_validate_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-24/tasks/task-9-done.md`

The shared extension is exactly the four DEC-284 files. This continuation did not edit or stage the runner-owned ledger, triage, authorities, contracts, plans, controls or generated files.

## Documentation updated

- `project-docs/specs/06_DESIGN_PLAN_CLI.md`: documents activation adoption across all 39 input leaves, both component content variants, sensitive fields, Unicode scalar bounds, unbounded strings and distinct cwd/artifact path rules.
- This report records both historical failures and the current red/green, provenance, runtime and full-gate receipts.

## Validation Contract Coverage

- VAL-011: PASS - `TestInputDescriptorRegistryMatchesLiveTreeAndTypes` (`cmd/fid/input_descriptor_test.go:419`) compares all 39 live leaves with the registry, reflects their types, rejects overlap and preserves the seeded mismatch failure. `TestP24ActivationComponentRawArrayReflection` (line 470) reflects both ActorProfile and GateDefinition element trees. `TestInputRawArrayReflectionRejectsVariantDrift` (`internal/cli/input_validate_test.go:122`) rejects missing/extra fields in either variant and the untyped-raw-array exemption.
- VAL-012: PASS - `TestP24WholeTreeDescriptorProcessCorpus` (line 151) walks help metadata for every field in all 39 descriptors, including both component trees and element bounds; `TestP24ActivationComponentHelpPin` (line 195) pins the full component help. Existing representative help pins remain in the ordinary gate.
- VAL-013: PASS - `TestP24WholeTreeDescriptorProcessCorpus` checks built-binary example equality, required-only fields, documented nested key order via `p24TemplateOrder`, and two-space formatting across all 39 leaves. `TestP24DescriptorVariantProcessCorpus` (line 164) also checks deterministic real-template selection while varying structural fixtures across both component types and all Stop codes.
- VAL-014: PASS - the whole-tree corpus decodes each example into its actual request type with unknown fields disallowed, independently checks required JSON types, and validates qualifying expanded structures. Both raw-array element trees are checked against their real Go types; raw JSON does not bypass reflection.
- VAL-015: PASS - `TestAdoptedExampleProcessExclusivity` (line 46) derives all 39 live leaves and checks 234 ordered conflicts with --json, --project, --input, --actor, --idempotency-key and positional input, plus standalone examples. It checks exit 2, expected error details and unchanged empty directories for conflicts, and exit 0, valid output, empty stderr and unchanged directories when alone.
- VAL-016: PASS - the whole-tree structural corpus and `TestP24DescriptorVariantProcessCorpus` exercise unknown, duplicate, missing, null, type, enum, declared-format and bound failures, recursively through selected trees. `TestP24ActivationComponentVariantMismatchAndSensitive` (line 204) exercises foreign content in both directions at input.content[0] and both environment[].value and redact_literals failures with value:null, no sentinel disclosure, no Event and no workspace change.
- VAL-035: PASS - `TestP24ChangedFakeFIDAssembly` (`cmd/fid/integration/p22_fake_test.go:26`) derives the intersection of fake dispatch with the 17 changed leaves and requires exactly the nine two-sided corpus rows. Each row tests accepted new input and identical typed removed-member rejection, with unchanged Event count/fake sequence on rejection. `gate.run` passes input admission and reaches matching committed budget denials on both adapters. Existing Run and gate parity tests and dual-adapter recovery conformance also pass the tagged gate.
- VAL-060: PASS - the built-binary whole-tree corpus runs every --example from an empty non-project directory, asserting exit 0, empty stderr and an unchanged empty tree before structural checks. The exclusivity corpus independently repeats the standalone no-write check for all 39 leaves.
- I-17: PASS - `TestP24InputDescriptorsHaveNoTopLevelSensitiveFields` (`cmd/fid/input_descriptor_test.go:490`) retains its guard over the complete 39-descriptor registry; sensitive markings remain nested.
- VAL-052: PASS - the fresh exact ordinary gate reports cmd/fid 410.123 s, internal/workflow 102.835 s and slowest other ordinary package internal/issues 199.339 s; all tested packages are under 600 s. The exact tagged gate reports 278.829 s, under 900 s. No recorded package reaches the 80% headroom thresholds (480 s ordinary / 720 s tagged). Current gate receipts below include actual termination results; cached historical timings do not supply this evidence.

## Tests run

Historical implementation checkpoint and report-only recovery receipts are preserved below. These are the earlier sessions' records, not current-session verification of the finalized report. The original task outcome was RECOVERY_REQUIRED. Current attempt-2 receipts follow them.

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
- `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-9-done.md project-docs/docs/phases/phase-24/tasks/task-9.md` - PASS - `PASS report-lint completed with no failures` (terminated `exit_code: 0`)
  Byte-identical verification after the committed formatting repair printed literal `OUTCOME RECOVERY_REQUIRED`; this validates the recovery report, not task completion.

Historical timing finding retained: required ordinary cmd/fid 282.333 s (failed assertions); slowest other uncached package internal/export 103.313 s; required tagged integration 333.935 s; internal/workflow was cached in that gate, followed by its supplemental fresh 55.650 s run recorded above.


Attempt 2 evidence, captured in this session. Focused red/green selections below precede the shared product edits or their corresponding repairs; they do not replace the exact gate.

- `go test ./internal/cli -run '^TestInputRawArrayVariants$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  First shared red before product edits: raw arrays reflected as uint8; both selected variant trees, nested duplicate redaction and help/example rendering failed.
- `go test ./internal/cli -run '^TestInputRawArrayVariants$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.003s` (terminated `exit_code: 0`)
  Same command after minimum typed-array selection, reflection, scanning, validation and rendering support.
- `go test ./cmd/fid -run '^(TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestP24ActivationComponentRawArrayReflection)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Own-session red preserved the live-registry mismatch and actual raw-array reflection failure.
- `go test ./cmd/fid -run '^(TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestP24ActivationComponentRawArrayReflection)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.160s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestP24WholeTreeDescriptorProcessCorpus$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  First draft failed to build: undefined commandLeaves; connected the corpus to the live tree.
- `go test ./cmd/fid -run '^TestP24WholeTreeDescriptorProcessCorpus$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Ordering mutation used an out-of-enum duplicate Stop action; retained the duplicate subject with a qualifying declared action.
- `go test ./cmd/fid -run '^TestP24WholeTreeDescriptorProcessCorpus$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  An empty optional array caused the ordering mutation to panic; constructed a qualifying duplicate instead.
- `go test ./cmd/fid -run '^TestP24WholeTreeDescriptorProcessCorpus$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  The template helper required runtime minimum array counts in required-only templates; DEC-276 templates intentionally contain one editable element. Runtime qualifying fixtures retain declared minimum counts.
- `go test ./cmd/fid -run '^TestP24WholeTreeDescriptorProcessCorpus$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	4.625s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestP24DescriptorVariantProcessCorpus$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  The test reordered selected enum values and used an invalid environment placeholder. Kept accepted allowed-value order and supplied a qualifying environment name; did not alter product enums.
- `go test ./cmd/fid -run '^TestP24DescriptorVariantProcessCorpus$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	14.221s` (terminated `exit_code: 0`)
  Both component arrays and all ten existing Stop object variants pass structural rows.
- `go test ./internal/cli -run '^TestInputArrayHelpElementBounds$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Array-element integer bounds were absent from help.
- `go test ./internal/cli -run '^TestInputArrayHelpElementBounds$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.002s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^TestP24ActivationComponentHelpPin$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Initial complete-text pin was pending; the full rendered help supplied the baseline after the failure.
- `go test ./cmd/fid -run '^TestP24ActivationComponentHelpPin$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Element bounds changed the complete help text; pin captured the resulting full text, not a subset.
- `go test ./cmd/fid -run '^TestP24ActivationComponentHelpPin$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.152s` (terminated `exit_code: 0`)
  Intermediate pin before the accepted-string boundary regression.
- `go test ./cmd/fid -run '^TestP24ActivationComponentHelpPin$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Boundary repair removed accidental 4KB limits and distinguished accepted path formats; converged the full-text pin to the accepted fields.
- `go test ./cmd/fid -run '^TestP24ActivationComponentHelpPin$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.158s` (terminated `exit_code: 0`)
  Final whole-text SHA-256: 301357ce9975f00a63dac43f6fd894f85da852544834801546b54c5694094168.
- `go test ./cmd/fid -run '^TestP24ActivationComponentAcceptedStringAndPathBounds$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Unbounded session/literal/redaction strings, scalar-counted Unicode bounds and artifact path bounds exposed seven failures.
- `go test ./cmd/fid -run '^TestP24ActivationComponentAcceptedStringAndPathBounds$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.160s` (terminated `exit_code: 0`)
- `go test ./internal/cli -run '^(TestInputRawArrayVariants|TestInputRawArrayReflectionRejectsVariantDrift|TestInputCodeObjectVariantsRemainSelected|TestInputRawArrayScalarBoundsPreserveOrdinaryCallers|TestInputArrayHelpElementBounds)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.006s` (terminated `exit_code: 0`)
  Both typed array variants reflect actual structs; missing/extra fields in either variant and untyped raw arrays fail. Existing object callers retain selection, duplicate handling, redaction, examples and ordinary byte bounds.
- `go test ./cmd/fid -run '^TestP24ActivationComponentVariantMismatchAndSensitive$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.715s` (terminated `exit_code: 0`)
  Real subprocesses prove mismatch in both directions at input.content[0], and both sensitive fields redact invalid values with value:null and no workspace change.
- `go test ./cmd/fid -run '^(TestDecodeClosedReviewInputRejectsMalformedInapplicability|TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.158s` (terminated `exit_code: 0`)
  Shared reader preserves malformed review rejection, valid empty optional arrays and authored budget state rejection for both empty-string and null states.
- `go test ./cmd/fid -run '^(TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestP24ActivationComponentRawArrayReflection|TestAdoptedExampleProcessExclusivity)$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Fence removal exposed an unused loop variable; removed the now-empty loop, retaining registry=live and all exclusivity assertions.
- `go test ./cmd/fid -run '^(TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestP24ActivationComponentRawArrayReflection|TestAdoptedExampleProcessExclusivity)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.896s` (terminated `exit_code: 0`)
  39 live leaves and 234 conflict combinations, plus standalone examples.
- `go test -tags integration ./cmd/fid/integration -run '^TestP24ChangedFakeFIDAssembly$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  The dispatch-derived shared leaf set initially lacked five setup operation rows in the fake.
- `go test -tags integration ./cmd/fid/integration -run '^TestP24ChangedFakeFIDAssembly$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`)
  The converged run.observe-budget fixture omitted its scenario Name, producing an empty source_sample_id; supplied the qualifying fixture identity.
- `go test -tags integration ./cmd/fid/integration -run '^TestP24ChangedFakeFIDAssembly$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.064s` (terminated `exit_code: 0`)
  All nine dispatch-intersection rows have both accepted input and removed-member rejection on both adapters; rejections preserve state. gate.run accepts the input and produces the same committed domain denial on both adapters.
- `go test -tags integration ./cmd/fid/integration -run '^(TestP24ChangedFakeFIDAssembly|TestP24RunChangedFakeFIDInputParity|TestP24GateChangedFakeFIDInputParity)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.272s` (terminated `exit_code: 0`)
- `go test ./internal/cli -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	0.024s` (terminated `exit_code: 0`)
- `go test ./cmd/fid -run '^(TestP24WholeTreeDescriptorProcessCorpus|TestP24DescriptorVariantProcessCorpus|TestP24ActivationComponentHelpPin|TestP24ActivationComponentAcceptedStringAndPathBounds|TestInputDescriptorRegistryMatchesLiveTreeAndTypes|TestP24ActivationComponentRawArrayReflection|TestP24ActivationComponentVariantMismatchAndSensitive)$' -count=1` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	18.277s` (terminated `exit_code: 0`)
  Final tree focused green.
- `gofmt -w internal/cli/input_descriptor.go internal/cli/input_validate.go internal/cli/input_validate_test.go cmd/fid/input_descriptor_test.go cmd/fid/input_descriptors_activation.go` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `goimports -w internal/cli/input_descriptor.go internal/cli/input_validate.go internal/cli/input_validate_test.go cmd/fid/input_descriptor_test.go cmd/fid/input_descriptors_activation.go` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.115s` (terminated `exit_code: 0`)
  First touched-package race tier, cmd/fid 279.960 s, before final boundary regression.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid ./internal/cli` - PASS - `ok  	github.com/mbell-nzl/fid/internal/cli	1.109s` (terminated `exit_code: 0`)
  Final touched-package race tier, cmd/fid 276.472 s; implementation checkpoint 69589747 preceded this long command.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP24ChangedFakeFIDAssembly|TestP22DualAdapterRecoveryConformance|TestP22FakeReviewReworkStopCondition|TestP24RunChangedFakeFIDInputParity|TestP24GateChangedFakeFIDInputParity)$'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	11.322s` (terminated `exit_code: 0`)
  Final touched tagged tier covers the new assembly and retained recovery conformance.
- `go build -o /tmp/fid-p24-task9-recovery2-KLEQP2/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `/tmp/fid-p24-task9-recovery2-KLEQP2/fid activation component --example` - PASS - `}` (terminated `exit_code: 0`)
  Working directory /tmp/fid-p24-task9-recovery2-KLEQP2/nonproject. Printed the raw authority-profile template with empty stderr.
- `/tmp/fid-p24-task9-recovery2-KLEQP2/fid --json activation component --input /tmp/fid-p24-task9-recovery2-KLEQP2/mismatch.json --idempotency-key task9-recovery2-entrypoint` - observed expected rejection (terminated `exit_code: 2`). Working directory `/tmp/fid-p24-task9-recovery2-KLEQP2/nonproject`; terminal output token was exactly `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["argv","artifact_bytes","artifact_paths","credential_names","cwd","definition_id","definition_version","environment","excerpt_bytes","network","output_bytes","redact_literals","sandbox","shell","success_exit_codes","timeout_ms","write"],"argument":"input.content[0]","value":null},"message":"invalid input.content[0]","retryable":false,"subject":null},"ok":false,"operation":"activation.component"}`. No Event or workspace writes.
- `ls -A /tmp/fid-p24-task9-recovery2-KLEQP2/nonproject` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  Observed the unchanged empty runtime working directory.

## Exact quality-gate results

Historical first gate attempt, retained verbatim. These receipts belong to the earlier failed checkpoint, not the current completion attempt.

The following six commands ran exactly in task-spec order on the checkpoint source tree. The ordinary failure was preserved and the remaining commands still ran. No source or test edits followed this gate attempt.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - `FAIL` (terminated `exit_code: 1`). Immutable failed attempt: exactly the missing component registry entry and raw-array reflection regression failed. cmd/fid 282.333 s, coverage 73.3%; no retained test failed. Shared-source scope is required to repair the remaining reds.
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.1%` (terminated `exit_code: 0`)
  This is a report from the failed ordinary attempt, not evidence that the ordinary gate passed.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	333.935s` (terminated `exit_code: 0`)
- `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` - BLOCKED - `cmd/fid/envelope.go:178:func decodeClosedEnvelopeInput(stdin io.Reader, path string, target any) error {` (expected `exit_code: 1`; terminated `exit_code: 0`). Three matches remain: review.go:135, activation.go:326 and envelope.go:178. The inventory is not an absence proof.

Report-lint is the Runner's control under DEC-158 and was not run by the Developer at the implementation checkpoint. This report-only recovery explicitly authorizes the exact report-lint verification recorded above.

Current attempt 2: the final source checkpoint `69589747ca77a1a3dae7e4d6083229d0193ed52f` preceded the final touched-package race tiers and all six required commands. Tests were uncached for fresh timing. The following exact commands ran serially in task-spec order from `/home/mark/projects/fid`; each terminated and these receipts are captured in this session.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	102.835s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					61.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	278.829s` (terminated `exit_code: 0`)
- `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` - PASS - `exit 1 (no output)` (expected `exit_code: 1`; terminated `exit_code: 1`)

Required ordinary package timing: cmd/fid 410.123 s (coverage 73.6%), internal/workflow 102.835 s (62.1%), slowest other ordinary package internal/issues 199.339 s (79.6%). Overall coverage is 61.8%. Required tagged integration: 278.829 s.

The last command's no-match condition is the task-spec-declared non-zero success: observed exit_code 1 with no output. All structural input decoding uses the shared descriptor reader; no family-local decoder remains.

The Phase Runner independently runs report-lint against this committed report. It was not run here, and no historical or remembered control result is claimed as current evidence.

## Deviations from spec

- Attempt 2 used only the ordinary Task 9 scope plus the four shared paths explicitly authorized by TRG-24-008 and DEC-284 clauses 4-6. The recovery identity remains task-9, chain trigger `task-9:report-lint:noncanonical-pass-evidence-termination`, attempt 2, authority `task spec; DEC-284`; these were checked against the ledger and governing triage entry as exact values.
- The first required ordinary failure remains `FAIL` (terminated `exit_code: 1`); its causal defect was missing typed raw-array reflection and registration for the accepted component content. Attempt 2 added selected typed element trees and the final component descriptor under the authorized shared extension, preserving registry=live and seeded mismatch. The byte-identical ordinary command now passes; its final token and observed termination are in the current gate receipts below.
- The original R5 failure remains `cmd/fid/envelope.go:178:func decodeClosedEnvelopeInput(stdin io.Reader, path string, target any) error {` (terminated `exit_code: 0`, expected 1). Attempt 2 converged retained Review/Envelope fixtures to qualifying shared-reader input and removed all remaining family wrappers. The byte-identical absence command now returns the declared exit 1 with no output.
- The accepted Task 9 component strings exposed accidental generic 4KB bounds and byte-counted Unicode values. The new boundary test failed with `FAIL` (terminated `exit_code: 1`), then passed the same command with `ok  	github.com/mbell-nzl/fid/cmd/fid	0.160s` (terminated `exit_code: 0`). Raw variant child bounds now use the accepted Unicode scalar rules; ordinary callers retain their original behavior. The whole-text help pin follows the corrected accepted fields and retains its full assertion.
- Newly reached corpus failures were fixture-conformance defects, not permission to weaken controls: supplied qualifying enum duplicates, array lengths, environment names and observation identity; kept exact error paths, types, allowed members, no-write checks and rejection strength. Required-only templates retain DEC-276's one editable array element, while qualifying runtime fixtures satisfy declared minimum counts. Every failing attempt and subsequent same-command green is recorded under Tests run.
- No accepted input/Event shape, public interface, descriptor meaning, gate, assertion, coverage policy, status/scope, threshold, budget or tolerance changed. Report-lint was not run in this implementation continuation: it is the Phase Runner's independent control per the task spec. Historical report-repair receipts below do not attest to this finalized report.

Attempt 2 retained-assertion convergence (introducing commits inspected with git history and Refs; labels and filenames are not ownership evidence):

| Retained assertion / fixture | Introducing commit | Phase / task | Accepted cause and preserved subject |
|---|---|---|---|
| Registry=live and seeded mismatch | `04b6b0992fe20a02df402a183269a9847ddae5ab` | phase-24/task-3 | Task 9's accepted activation adoption requires the final live leaf; registration completes the assertion without a fence/exemption. |
| p24DescriptorProcessCorpus / qualifying mutation assumptions | `8d1c83855c0f681515e44e2c3a61bbe543ef515b` | phase-24/task-7 | Task 9 broadens accepted conformance to all 39 leaves and selected raw content. Duplicate mutations use qualifying declared values and preserve required counts, exact paths/details and no-write assertions. |
| Required-only typed template helper | `226cdc94e6460cc254d1b4e916344b8ea2e2cdb3` | phase-24/task-5 | Task 9 adds activation arrays with runtime minima 2/5. DEC-276 template semantics still require one editable element; strict actual types, required keys and order remain checked, and expanded runtime baselines retain the minima. |
| TestDecodeClosedReviewInputRejectsMalformedInapplicability | `0be596e9068795ceb3a0a89bdac50d538f8e8791` | phase-23/task-5 | Task 9's accepted single structural source displaces the family decoder call. Uses complete qualifying input via the shared reader; retains malformed null/duplicate/nested-unknown rejection and empty optional-array acceptance. |
| TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState | `2339dc9ebd8d0e6b5502ca32e1562892cf704ce6` | phase-22/task-12 | Task 9's decoder convergence displaces the sparse family fixture. Real CLI still rejects both empty-string and null authored state, with typed failures and unchanged workspace. |
| Raw ActivationComponentRequest reflection regression | `785a22491f1ba07bb748d1d4fb85f88ea1a8d45d` | phase-24/task-9 | Kept the regression; added actual ActorProfile and GateDefinition type trees, with no raw-byte escape from reflection. |
| Complete component help pin | `9d654e976792836cbdf7de8c114d0831b73f721e` | phase-24/task-9 | Corrected accidentally imposed bounds to the already accepted component contract; reconverged the entire help-text digest rather than narrowing the assertion. |

The generalized p24DescriptorProcessCorpus function was introduced by Task 7; its retained binding mutation rows originated in Task 6, as recorded in the historical table below. The introducing function and the older row subjects have separate provenance.

Historical deviations and failed attempts, preserved verbatim; references to incomplete implementation or pending scope below describe the earlier checkpoint only:

- Attempt 2 found and repaired accidentally imposed generic text bounds on unbounded component strings. The focused accepted-string/path test failed with `FAIL` (terminated `exit_code: 1`) for session/literal/redaction strings, Unicode bounds and artifact paths. Component descriptors now follow the accepted Phase 21 fields; raw variant bounds count Unicode scalars while existing ordinary callers retain their length behavior. The same command passed with `ok  	github.com/mbell-nzl/fid/cmd/fid	0.160s` (terminated `exit_code: 0`). The component help pin is converged to these corrected declared fields, preserving its whole-text assertion.

Historical deviations and failed attempts (preserved):

- Report-only recovery attempt 1 reproduced `FAIL report-lint found 11 failure(s)` (terminated `exit_code: 11`). Root cause: the 11 historical PASS rows had punctuation or explanatory text after their termination fields, but report-lint requires those fields at the end of each evidence line. Repair: remove only trailing punctuation and move explanations to continuation lines, preserving every historical command, status, token, observed termination, provenance, scope finding and the typed RECOVERY_REQUIRED outcome. The later byte-identical verification ended with `PASS report-lint completed with no failures` (terminated `exit_code: 0`) and printed `OUTCOME RECOVERY_REQUIRED`.
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

Historical checkpoint scope findings (now resolved by the authorized implementation and final gates; preserved as provenance):

- Scope gap: `internal/cli/input_descriptor.go:88` treats raw content as a slice of uint8 and recursively requires an object element. `internal/cli/input_validate.go:99` and `:181` select object variants only through root `code`; `internal/cli/input_render.go:52` labels variants as code and `:109` selects example variants through code. These shared files and their regression tests require ordinary Task 9 source scope, not a tools-only control-conformance directive. The accepted component_type shapes need no new product or architecture decision; the missing authority is the file scope.
- Continue with all eight activation descriptors, both typed content variants and exact mismatch paths; remove the final pending fence while retaining registry=live. Complete all six 39-leaf corpora, the sensitive nested failure cases, component help pin and fake-dispatch-derived nine-leaf assembly. Remove all three remaining family decoders without retaining another structural source.
- Pending decoder-fixture provenance was inspected but no convergence is claimed: TestDecodeClosedReviewInputRejectsMalformedInapplicability was introduced by `0be596e9068795ceb3a0a89bdac50d538f8e8791` (phase-23/task-5); TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState by `2339dc9ebd8d0e6b5502ca32e1562892cf704ce6` (phase-22/task-12). Preserve valid empty optional-array acceptance and both output-only-state rejection states when converging them to descriptor validation.
- Observed timing samples are below their warning thresholds. Cached package results do not establish fresh wall times; complete the remaining VAL-052 timing evidence in the continuation.
- No grammar/unknown-command implementation (Task 9a), Human failure-line implementation (Task 10), runner control, ledger or triage edits were made.

## Concerns, technical debt, or notes for the Tech Lead

- Fresh gate wall times are below both budgets and their 80% warning thresholds: cmd/fid 410.123 s, internal/workflow 102.835 s, slowest other ordinary internal/issues 199.339 s, tagged integration 278.829 s. No headroom warning is required.
- Per DEC-284 clause 6, the Runner performs the late window C immediately after Task 9 completion and before dispatching Task 9a, reviewing tasks 7-8 at their close commits plus current HEAD and feasibility for tasks 9a-14. D1 remains after Task 9a and reviews the shared regression evidence.
- Report-lint remains the Runner's independent control. Grammar/unknown commands and Human failure lines were kept outside Task 9; no workflow control, authority or ledger edit was made by this continuation.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
