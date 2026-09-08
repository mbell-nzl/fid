# Task 15 root continuation evidence

Root is the sole product writer after the Developer's incomplete uncommitted
hand-back at phase HEAD `bab537b`. No final gate or completion report existed
at hand-back. Earlier Developer observations are attributed in the report;
the final exact gate receipts will be observations from this root session.

## Behavioral failures and repairs

Working directory for every command: `/home/mark/projects/fid`.

- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22PortableCorpusCoversRecoveryStates$' -count=1` — FAIL — `portable recovery states = [], want exact [abandoned active deduplicated-attention exhausted-budget interrupted missing-artifact partial-activation unknown-budget]` (terminated `exit_code: 1`). Root added the missing eight actual portable fixtures and their shared expectations.
- `go test -tags integration ./cmd/fid/integration/... -run '^TestP22Portable' -count=1` — FAIL — `panic: interface conversion: interface {} is nil, not float64` (terminated `exit_code: 1`). The independent fake initially read Budget maximum as a direct column; accepted policy is encoded in `budget_json`. Root corrected the model's decoding; expected Budget outcomes stayed fixed.
- The byte-identical command — FAIL — `FAIL` (terminated `exit_code: 1`). The partial activation fixture received CLI `invalid-argument` / exit 2 because it supplied the delivery phase and omitted prior resource identities. Root bound the upgrade to the finalized predecessor, its version/resources and accepted phase 21. No activation policy changed.
- The byte-identical command — PASS — `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	12.265s` (terminated `exit_code: 0`). All eight states passed fake/fid expectations and live/rebuilt/imported named-view and entry-byte parity.
- After adding a second independent fake reconstruction and shared valid-hash negative controls, the byte-identical command — PASS — `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	13.011s` (terminated `exit_code: 0`). Both fake and fid refuse unknown tables, columns, Effect operations, retired Event types and duplicate Event IDs.

## Interpretation and scope

The first final `go test -race -coverprofile=coverage.out ./...` attempt failed:
`NewFailure("import-target-conflict") error = new failure: import-target-conflict reason "reason-value" is not accepted`
(terminated `exit_code: 1`; final token `FAIL`). The generic schema test's
synthetic valid string was outside Task 15's new closed reason domain. Root
changed only that row's valid sample to `populated`; missing/extra/wrong-type
and unknown-reason rejection checks remain. This is a Task 15 test-construction
repair, not a product enum relaxation or inherited-failure claim. The exact
full command is rerun below in the completion report after formatting/lint.

The fake receives immutable archive entries as fixture input and interprets
logical Effects in memory. It never opens fid SQLite, journal paths or `.fid`,
and does not use fid domain appliers. Its schema literals are a pinned accepted
schema-15 column/default inventory; they are not response fixtures. This is a
portable semantic model, not a replacement for Task 13's full archive validator.

The one-time column inventory inspection initially failed under the older
Python SQLite parser (`sqlite3.OperationalError: near "CASE": syntax error`,
exit 1). Excluding trigger bodies from that column-only inspection succeeded;
no migration, production validator or gate was changed.

Full named views use the same observed repository input across all three
paths. This keeps live environmental observation controlled without deleting
Worktree IDs or any durable fields. The separate built CLI journey observes
the real destination, reconciles recovery, and executes a successor Gate with
a Receipt. Archive entry comparison includes every entry byte; there is no
opaque database comparison or weakened status/identity normalization.

## Convergence provenance

- `cmd/fid/export_test.go` import absence, `capabilities_test.go` deferred import
  invocation and the export leaf registry originated at
  `35a75baf7b07d8f5eaace36ac8142ea6ac7c4ddd`, Phase 21/task 19 (`git log -S
  'export import'` inspected). Task 15's accepted CLI import invalidates absence;
  current tests retain closed errors/no-side-effects and exact command/flag
  inventories. The final ordinary race command executes all three files.
- Existing Phase 22 fake/corpus helpers are extended in their current task
  lineage; their original recovery, activation and error assertions remain.
  The full tagged command executes original and added cases.

No dependency, public contract, VAL or protected command changed. Required
final gates and the product commit are recorded in `task-15-done.md`.
