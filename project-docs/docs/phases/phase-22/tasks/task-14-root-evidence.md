# Task14 root continuation evidence

Root took sole product ownership after the Developer's clean checkpoint
0a4d00de3581c05f5f1c1b9748da87916e28bc1c, with no live commands. This is
same-task recovery under task14 / DEC202 / DEC204, with no new product policy.
Working directory for every command here: `/home/mark/projects/fid`.

## Root red/green and causal repairs

- `go test ./internal/export -run '^(TestImportLockExcludesRepeatedContenders|TestImportBindingRejectsAlteredHistoryIdentity)$' -count=1` failed with `contender 2 entered while owner still held import lock`, `altered imported Event identity still qualifies installed binding`, and `altered imported Event identity still replays installation`; final token `FAIL` (terminated `exit_code: 1`).
- Root repaired exclusive locking by locking the target directory inode. No contender can unlink a lock pathname; process exit releases ownership and no new lock file changes the repository observation. Imported marker verification now checks immutable ordinal journal chunk size/digest and its ordered Event identity inventory. Replay additionally compares the entire recorded result against the validated archive's expected result.
- The byte-identical focused rerun passed: `ok  	github.com/mbell-nzl/fid/internal/export	0.158s` (terminated `exit_code: 0`).
- `go test ./internal/export -run '^TestImport' -count=1` passed: `ok  	github.com/mbell-nzl/fid/internal/export	5.390s` (terminated `exit_code: 0`). It includes five actual subprocess exits without deferred cleanup, followed by correct install/replay; changed digest, altered journal bytes, missing journal, altered result count and symlinked installation reject without changing destination state.
- The Developer's snapshot/stale-file test was strengthened at Task14 provenance: replacement now occurs after validation and before staging, rather than after staging was already complete. A preexisting `.fid-import.lock` file remains untouched, as required by source preservation; actual process death and repeated contender tests replace the old physical lock-file deletion assertion.
- `go test -race ./internal/export -run '^(TestImport|TestArchiveRejectsIntegrityValidSemanticCorruptions)' -count=1` passed: `ok  	github.com/mbell-nzl/fid/internal/export	97.510s` (terminated `exit_code: 0`). The destination journey replays the old observation key without gaining authority, rejects old observation Close/Abandon/recovery and historical Gate execution, then uses new local facts for destination-bound successor Claim/Run and real Gate/Receipt; import replay after those new Events preserves history.
- `go test ./internal/export -run '^TestValidateRetainsPinnedV1Profile$' -count=1` passed: `ok  	github.com/mbell-nzl/fid/internal/export	0.083s` (terminated `exit_code: 0`). Retained v1 validates, import returns the existing `schema-unsupported` archive check, and the rejected target remains empty.

## Final exact gate receipts

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)

- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	101.200s	coverage: 52.9% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					63.2%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	52.659s` (terminated `exit_code: 0`)

All five final exact commands were observed by root to terminate successfully.
Product checkpoint: `8c322b5`, on top of the retained Developer implementation.
The full race suite includes export at 438.056s and CLI at 449.577s; no test
timeout, budget, limit or gate definition was changed.

## Convergence provenance

- Pinned-v1 validation test: `acb569bb9447a3259991520669167fdc819c1706`,
  Phase22/task13. Task14 retains every validation assertion and adds rejected
  import/schema-unsupported plus unchanged target evidence. It executes under
  `go test -race -coverprofile=coverage.out ./...`.
- Nested observation and aggregate identity regressions: same Task13 commit,
  repaired in Task14 by the Developer and retained in root's focused race and
  full exact gate. ROOT-SEM-01/02 remain independently reviewable in Review5.
- DEC204 Close/Abandon consumer: original closure producer
  `b965b24748f57a0f5e42b8dd3054dfe6fd5dbb74`, Phase22/task8. New imported-history
  fencing preserves normal closure tests; the full ordinary and tagged suites
  execute both the retained semantics and new negative/positive import cases.
- Journal writer original `b0f36534963ae4768629d9850aca36e3f46c89ff`,
  Phase1/task3. Imported ordinal chunks invalidate assumptions about dated
  physical filenames, not ordered Event identity. Existing writer/replay tests
  remain in the full ordinary gate; the destination journey writes new Events
  after import without collision and verifies subsequent replay preservation.
- Developer's new Task14 physical lock-file assertion originated in checkpoint
  `0a4d00de3581c05f5f1c1b9748da87916e28bc1c`. Its accepted subject is crash
  recovery and preservation of preexisting source files. Root replaces deletion
  of that preexisting file with preservation and actual abrupt-exit/repeated
  contender evidence; it does not weaken exclusivity or retry behavior.
