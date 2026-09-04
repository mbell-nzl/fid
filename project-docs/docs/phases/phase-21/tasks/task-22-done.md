# Phase 21 - Task 22: Final Wire Conformance And Real-Process Journey - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 3
- Next action: End-of-phase QA review.
- Human attention required: no

## Summary of what was implemented

- Generated and byte-compared the committed `p21_conformance` corpus for every
  accepted/fail-closed operation, failure code, live frame, and closed public
  value; the executable leaf inventory now includes all 63 retained leaves.
- Built the real-binary one-Issue journey with its project-authored workflow,
  Gate-before-terminal Evidence chain, exact ordered 22 Events, undeclared-edge
  negative checks, deterministic export, projection deletion/rebuild, and
  mutation-free post-rebuild repeat.
- Replaced the activation Event suffix synthesis with an exhaustive typed
  mapping so commit, abort, and finalize emit their accepted Event names.
- Added accepted read-only `receipt get/list` domain and CLI surfaces with exact
  filters, ordering, limits, complete JSON, typed not-found handling, and no
  Event append. Receipt rejections are stored as canonical JSON in migration
  0013 and survive idempotent replay and projection rebuild.
- Corrected the journey's invocation to the accepted root `fid resume` command
  and corrected only the matching stale prose in `contracts.md` section 10.

## Files created or modified

- `cmd/fid/activation_test.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/integration/p21_conformance_test.go`
- `cmd/fid/integration/testdata/p21_conformance/`
- `cmd/fid/receipt.go`
- `cmd/fid/receipt_test.go`
- `cmd/fid/root.go`
- `internal/cli/codes.go`
- `internal/cli/codes_test.go`
- `internal/db/operational_kernel_migration_test.go`
- `internal/evidence/execution.go`
- `internal/evidence/receipt_read.go`
- `internal/workflow/activation.go`
- `migrations/0013_p21_operational_kernel.up.sql`
- `project-docs/docs/phases/phase-21/tasks/task-22-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/contracts.md`
- `project-docs/docs/phases/phase-21/tasks/task-22-walking-skeleton.md`
- `project-docs/specs/03_TRD_SCHEMA.md`
- DEC-174 through DEC-177 and their current/archived workflow artifacts record
  each protected recovery without rewriting the prior attempts.

## Validation Contract Coverage

- VAL-072: PASS - `TestPhase21ConformanceCorpusIsCompleteAndByteExact` at
  `cmd/fid/integration/p21_conformance_test.go:25` regenerates and byte-compares
  every committed operation, failure, live-frame, and manifest fixture; the
  exact tagged gate passed with terminated `exit_code: 0`.
- VAL-073: PASS - the same test at
  `cmd/fid/integration/p21_conformance_test.go:56` requires every extracted
  closed member and rejects non-kebab or underscore values; corpus inventory
  returned `failures`, `live`, `manifest.json`, and `operations`.
- VAL-084: PASS - `TestPhase21RealProcessJourneyStartsOneActivatedEnvelope` at
  `cmd/fid/integration/p21_conformance_test.go:73` runs the final binary through
  one Envelope, Gate/Receipt/Evidence, the single terminal-success transition,
  Run finish, Claim release, and Envelope completion.
- VAL-085: PASS - `p21AssertUndeclaredTransitionRejected` at
  `cmd/fid/integration/p21_conformance_test.go:188` requires exit 4,
  `transition-not-permitted`, `workflow-transition-invalid`, and zero Event
  delta both before and after rebuild.
- VAL-086: PASS - the exact Event lists at
  `cmd/fid/integration/p21_conformance_test.go:139` and `:178` require the same
  ordered 22 Events; lines 147-177 perform only accepted reads, Receipt
  get/list, evaluation, resume, export/validate, and the negative check across
  projection deletion/rebuild.

## Tests run

- `go test -tags integration ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Attempt-1 evidence: the seeded `task-default` path rejected the requested
  terminal transition with `transition-not-permitted` /
  `workflow-transition-invalid`; DEC-174 preserved this as the VAL-085 negative
  and authorized the project-authored journey workflow.
- `go test -tags integration ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Attempt-2 red evidence: the retained target `unstarted -> in_progress` was
  outside the authored workflow; the repaired target is `unstarted -> started`
  and the undeclared edge remains a separate no-mutation negative assertion.
- `go test -race ./cmd/fid` - FAIL - `FAIL` (terminated `exit_code: 1`)
  DEC-175 red evidence: `activation_test.go:478: activation finalized Event
  type = "activation.finalizeted", want "activation.finalized"`.
- `go test -race ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	344.360s` (terminated `exit_code: 0`)
  DEC-175 green evidence after the exhaustive activation Event mapping.
- `go test -tags integration -race ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Attempt-4 first failure: `p21_conformance_test.go:148` invoked stale
  `fid projection resume` and observed Cobra exit 2; the test was corrected to
  the already-accepted root `fid resume` row.
- `go test -tags integration -race ./cmd/fid/integration/...` - FAIL - `FAIL` (terminated `exit_code: 1`)
  Attempt-4 byte-identical rerun: `p21_conformance_test.go:162` reached
  `fid receipt get` and observed unknown-command exit 2, producing TRG-21-023.
- `go test -race ./cmd/fid --run 'TestReceipt|TestGateRunReplayRetainsRejectedArtifacts|TestLiveLeafCommandsMatchContractInventories'` - FAIL - `FAIL` (terminated `exit_code: 1`)
  DEC-177 red evidence: inventory lacked Receipt leaves, both commands were
  absent, missing Receipt had no typed failure, and replay returned a nil
  rejection array.
- `go test -race ./cmd/fid --run 'TestReceipt|TestGateRunReplayRetainsRejectedArtifacts|TestLiveLeafCommandsMatchContractInventories'` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	7.981s` (terminated `exit_code: 0`)
- `go test -race ./internal/cli ./internal/db ./internal/evidence` - PASS - `ok  	github.com/mbell-nzl/fid/internal/evidence	2.274s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	32.404s	coverage: 27.4% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:									(statements)			61.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.927s` (terminated `exit_code: 0`)
- `ls cmd/fid/integration/testdata/p21_conformance` - PASS - `operations` (terminated `exit_code: 0`)
- `go test -bench . -benchmem ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	2.318s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- The immutable failed attempts exposed three planning/implementation gaps:
  the original journey assumed an unavailable seeded direct-terminal workflow;
  activation finalize emitted `activation.finalizeted`; and the accepted
  Receipt reads plus durable rejection facts were absent. DEC-174, DEC-175, and
  DEC-177 respectively authorized the causal repairs. Each later command was
  rerun byte-identically, and no accepted command, result, failure, Event,
  assertion, status, scope, threshold, budget, tolerance, or prior evidence was
  weakened.

## Concerns, technical debt, or notes for the Tech Lead

- None. The first partial post-activation gate attempt and every prior failed
  tagged/red result remain recorded in the append-only ledger and this report.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `ac6270ad34d51f5d39dedd2bff1fddeb7f9e5632`

## Commit subject

- `feat(evidence): expose durable receipt reads`
