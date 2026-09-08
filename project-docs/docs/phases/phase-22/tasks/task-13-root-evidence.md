# Task 13 — Root boundary evidence

Root owns only `internal/export/export_boundary_test.go` and this contribution
record. The Developer owns production/archive semantics and the final five
required gates. All commands here run from `/home/mark/projects/fid` and
retain root-session attribution.

## Initial bounded red

The initial table uses a complete valid bootstrap archive as a control, changes
only ZIP file mode for special entries, and adds an accurately sized/digested,
valid-CRC 64 MiB + 1 non-Artifact entry for the byte limit. Before the production
scanner changed, the regular control passed and the validator incorrectly
accepted all three special file modes and the expanded oversized member.
Source-project snapshots were unchanged by every validation attempt.

`go test ./internal/export -run '^TestArchiveRejectsSpecialEntriesAndExpandedMember$' -count=1`

- Semantic red: `FAIL	github.com/mbell-nzl/fid/internal/export	0.702s`
  (terminated `exit_code: 1`). Each bad case reported `validate error = <nil>`
  instead of its expected path-unsafe/resource-limit check. No production
  change preceded this result.

The full numeric/path matrix and final root results remain in progress.

## Expanded root integration ownership

The Human's super-delegate authorization permits direct implementation and
bounded assistance. Root now owns archive writer/manifest/collector integration,
retained test convergence, this evidence, and final task-gate/report integration.
Task11's reused logical task13 session owns archive_scan.go/archive_limits_test.go;
task12's reused logical task13 session owns semantic validator/tests; review4's
new logical implementation-assistance session owns a journal-backed fixture helper
and will not independently judge that contribution. Prior completed tasks stay closed.

## Root observed focused receipts

- Byte-identical initial boundary command later GREEN: `ok  	github.com/mbell-nzl/fid/internal/export	0.518s` (terminated `exit_code: 0`).
- Mechanical split focused boundary/Source command: `ok  	github.com/mbell-nzl/fid/internal/export	0.832s` (terminated `exit_code: 0`).
- `go test ./internal/export -run '^TestCreateRecoveryManifestV2ConfigAndOrdinalJournal$' -count=1`: semantic RED `schema = fid.export-manifest.v1, want fid.export-manifest.v2`; `FAIL	github.com/mbell-nzl/fid/internal/export	0.075s` (terminated `exit_code: 1`).
- Combined retained-v1/Recovery focused run: retained-v1 semantic RED `retained v1 profile rejected`; `FAIL	github.com/mbell-nzl/fid/internal/export	0.141s` (terminated `exit_code: 1`). Its Recovery failure was unrelated empty-line fixture construction and is NOT semantic reference evidence.
- `go test ./internal/export -run '^(TestValidateRetainsPinnedV1Profile|TestCreateRecoveryManifestV2ConfigAndOrdinalJournal)$' -count=1`: GREEN `ok  	github.com/mbell-nzl/fid/internal/export	0.143s` (terminated `exit_code: 0`).
- `go test ./internal/export -run '^TestValidateRecoveryManifestClosedConfigurationAndInventory$' -count=1`: semantic RED for ordinal gap and unknown entry family, both `validate error = <nil>`; `FAIL	github.com/mbell-nzl/fid/internal/export	0.507s` (terminated `exit_code: 1`).
- `go test ./internal/export -run '^(TestValidateRecoveryManifestClosedConfigurationAndInventory|TestCreateRecoveryManifestV2ConfigAndOrdinalJournal|TestValidateRetainsPinnedV1Profile)$' -count=1`: GREEN `ok  	github.com/mbell-nzl/fid/internal/export	0.648s` (terminated `exit_code: 0`).
- After initial private semantic replay integration: `go test ./internal/export -run '^(TestArchiveWithValidHashesButInvalidRecoveryReferenceRejects|TestCreateRecoveryManifestV2ConfigAndOrdinalJournal|TestValidateRetainsPinnedV1Profile)$' -count=1` GREEN `ok  	github.com/mbell-nzl/fid/internal/export	0.285s` (terminated `exit_code: 0`). Strict ordered Effect-family completion and stronger real-recovery controls remain pending; this is not final VAL-055 closure.

Construction attempts during the coordinated signature migration terminated
exit_code 1 with `too many arguments in call to journalArchiveEntries`,
`undefined: archiveEntryByteLimit`, and `cannot use archivePath (variable of type
string) as ValidateOptions value in argument to Validate`. These were corrected
without changing assertions; none are semantic red evidence. The initial
Recovery test also inserted an extra empty line and failed `journal-invalid`;
its author corrected the fixture before capturing the true acceptance red.
Agent-owned focused receipts are retained in their messages and must be
attributed if copied, never recorded as root-captured PASS.

## Serial root continuation after usage reset

The Human approved resumption with one implementation owner at a time. Root
owns all remaining Task13 code/tests/report; no subagent is active. Earlier
ownership statements above describe historical dispatches and remain intact.

Fresh root command receipts from this resumed session:

- Initial focused construction failure: `non-name entries[journalPath] on left side of :=`, unused encoding/json and undefined err; terminated exit_code 1. The interrupted semantic test edit was finished; no assertion changed.
- Broad focused attempt then failed `release fixture predecessor Claim: claim release input is invalid`, `FAIL	github.com/mbell-nzl/fid/internal/export	7.187s` (terminated exit_code 1). The fixture supplied unsupported `interrupted`; the accepted enum is completed/failed/operator-released. Using operator-released repairs setup, not product policy.
- `go test ./internal/export -run '^TestArchiveWithValidHashesButInvalidRecoveryReferenceRejects$' -count=1` passed: `ok  	github.com/mbell-nzl/fid/internal/export	0.795s` (terminated exit_code 0). It first validates a complete real recovered history, then changes only the Recovery observation reference and recomputes hashes.
- Count/early-stop focused command passed: `ok  	github.com/mbell-nzl/fid/internal/export	10.135s` (terminated exit_code 0). Actual central records are checked before ZIP File allocation and by directory inspection; overflow consumes exactly L+1 bytes and does not read the next member.
- New semantic-corruption matrix initially rejected an invented test column, twice (event_id then committed_event_id; actual schema is result_event_id). Those construction receipts were failures 1.161s and 1.196s, both exit_code 1, not reference proof.
- Corrected semantic RED: `go test ./internal/export -run '^TestArchiveRejectsIntegrityValidSemanticCorruptions$' -count=1` returned `FAIL	github.com/mbell-nzl/fid/internal/export	1.193s` (terminated exit_code 1): unresolved idempotency Event and foreign Recovery aggregate were both accepted. Event-reference and Recovery-ownership checks repaired the causes.
- Retired run.finished semantic RED: `FAIL	github.com/mbell-nzl/fid/internal/export	0.760s` (terminated exit_code 1), before removing that retired v2 family from the catalog.
- Boolean-as-version semantic RED: `FAIL	github.com/mbell-nzl/fid/internal/export	0.766s` (terminated exit_code 1), before restricting booleans to accepted boolean columns.
- Executor-identity rewrite semantic RED: `FAIL	github.com/mbell-nzl/fid/internal/export	0.823s` (terminated exit_code 1), before restricting operational update columns to mutable contract fields.
- Combined semantic/full valid controls passed: `go test ./internal/export -run '^(TestArchiveRejectsIntegrityValidSemanticCorruptions|TestArchiveWithValidHashesButInvalidRecoveryReferenceRejects|TestSemantic|TestCreate|TestValidate)' -count=1`: `ok  	github.com/mbell-nzl/fid/internal/export	9.340s` (terminated exit_code 0).
- Actual ZIP path matrix RED: dot/NUL names reached consume (`checks=[] consumed=1`); `FAIL	github.com/mbell-nzl/fid/internal/export	0.006s` (terminated exit_code 1). Added path rejection; identical command GREEN `ok  	github.com/mbell-nzl/fid/internal/export	0.006s` (terminated exit_code 0).

### Required gate attempt history

- gofmt -l cmd internal: silent success, exit_code 0 (before each lint attempt).
- golangci-lint run ./... attempt1: `* staticcheck: 4`, exit_code 1. Uppercase error strings and a boolean-expression simplification remained from interrupted work.
- Identical lint attempt2: `* staticcheck: 3`, exit_code 1. More uppercase strings became visible after the first corrections (per-file reporting limits).
- Identical lint attempt3: `0 issues.`, exit_code 0. All error strings corrected together; predicates unchanged.
- Full race/coverage gate is running at this record update; no passing outcome is claimed until observed termination.

### Contract interpretation and provenance

Independent QA examined Phase21 contracts section11 (lines333–337): manifestv1
validation explicitly does not establish identity/history parity. Phase22 task13
and public-interface section5 retain that behavior. Accordingly v1 uses its
fixed schema14/capability profile plus integrity/safety/bounds checks; v2 alone
uses full semantic replay. No fallback from a v2 failure to v1 exists.

The original export fixture/helper and date-based inventory assertions were
introduced at 35a75baf7b07d8f5eaace36ac8142ea6ac7c4ddd, Phase21/task19.
V2 requires ordinal chunks and journal-backed Sources/Artifacts, so those
fixtures now use normal public writers and actual Git. Both available and
unavailable Artifact cases remain; no former negative assertion was removed.

The original named first-failing test actually corrupted Source metadata.
It is retained under an accurate Source-metadata name. The semantic assistant
later captured the true public-Validate acceptance RED before its SQL replay
implementation, after unrelated scanner/writer work had begun. This is a TDD
sequencing/initial-test-quality deviation, not a claim that the initial test
proved Recovery references. Root subsequently strengthened it to a complete
real Recovery history with one missing observation reference.

### Root real CLI smoke, 2026-09-08 17:17 NZ

- `go build -o /tmp/fid ./cmd/fid` from repository root terminated `exit_code: 0`, `exit 0 (no output)`.
- Fresh scratch `/tmp/fid-p22-task13-smoke-1T0ybJ`; `/tmp/fid --project /tmp/fid-p22-task13-smoke-1T0ybJ --json init --prefix SMK` terminated `exit_code: 0`, returned `project.init`, `ok:true`, Project `project-75342dc407fa554b9469b675443ec293`.
- `/tmp/fid --project /tmp/fid-p22-task13-smoke-1T0ybJ --json export create /tmp/fid-p22-task13-smoke-1T0ybJ/recovery.zip` terminated `exit_code: 0`, returned `fid.export-manifest.v2`, projection15, one Event, two entries and archive digest `sha256:b4d4af87e360dfb0e65796205769006eec701b0e5b2e522141811bbc70c68e15`.
- `/tmp/fid --json export validate /tmp/fid-p22-task13-smoke-1T0ybJ/recovery.zip` from `/tmp` (outside any workspace) terminated `exit_code: 0`, returned `export.validate`, `ok:true`, `status:valid`, matching digest/schema/count. This is a Task13 smoke, not a retained Human delivery fixture.

### Required race/coverage attempt 1: timed out

`go test -race -coverprofile=coverage.out ./...` terminated `exit_code: 1`, final token `FAIL`. Export package emitted `panic: test timed out after 10m0s`, `TestArchiveCumulativeBoundary/8589934591 (5m6s)` and `FAIL	github.com/mbell-nzl/fid/internal/export	600.087s`. Stack was in `compress/flate.(*compressor).findMatch` through `writeZeroZIP`, generating the first 8GiB fixture. The CLI package passed at560.974s, but the required command FAILED; partial coverage is not closure evidence. Root will reuse a stdlib-generated, byte-aligned non-final DEFLATE block for zero fixture generation, retaining actual CRC-correct ZIP expansion and every original boundary, then rerun the byte-identical gate. No timeout/command/limit change.

- After fixture-generation repair, `go test -race ./internal/export -run '^(TestArchiveExpandedByteBoundaries|TestArchiveCumulativeBoundary|TestArchiveStopsAtFirstExpandedOverflow)$' -count=1` passed `ok  	github.com/mbell-nzl/fid/internal/export	168.333s` (terminated `exit_code: 0`). This streams all original 64MiB/1GiB/8GiB below/equal/above cases, with full CRC verification for admitted members.
- Final-source `gofmt -l cmd internal` terminated `exit_code: 0`, `exit 0 (no output)`; `golangci-lint run ./...` terminated `exit_code: 0`, `0 issues.`. Byte-identical full race/coverage attempt2 launched after these results.

- Full race/coverage attempt2 passed, final `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 52.9% of statements` (terminated `exit_code: 0`). CLI461.712s/72.8%, export367.930s/78.6%.
- `go tool cover -func=coverage.out` passed, final `total:										(statements)					63.2%` (terminated `exit_code: 0`). Tagged integration gate then launched.

- Required `go test -tags integration -race ./cmd/fid/integration/...` attempt1 FAILED, final `FAIL` (terminated `exit_code: 1`), package55.349s. `TestPhase21RealProcessJourneyStartsOneActivatedEnvelope` atp21_conformance_test.go:179 rejected its freshly created v2 archive: `err:exit status 7 stdout:"" stderr:""`. Diagnostic investigation pending; assertion and journey retained.

- Focused tagged diagnostic `go test -tags integration -race ./cmd/fid/integration/... -run '^TestPhase21RealProcessJourneyStartsOneActivatedEnvelope$' -count=1` FAILED, final `FAIL` (terminated `exit_code: 1`),3.217s. Private error was `Effect 4 table/operation idempotency_results/insert is not accepted for "run.execution-blocked"`. Source `internal/evidence/execution.go` uses normal CommitIdempotent after this Event's Stop/Attention/runtime effects; validator catalog omitted its final row. Added that required stage without changing the writer/contract.
- `go test ./cmd/fid -run '^TestExportValidateRendersInvalidArchiveChecks$' -count=1` RED `FAIL` (terminated `exit_code: 1`),0.082s. Both nil and named Check.Entry pointers caused empty stderr because the closed CLI checker accepts strings/null, not Go pointers. Corrected the scoped adapter's conversion. Byte-identical GREEN `ok  	github.com/mbell-nzl/fid/cmd/fid	0.088s` (terminated `exit_code: 0`). Introducing adapter provenance: `35a75baf7b07d8f5eaace36ac8142ea6ac7c4ddd`, `feat(export): create deterministic archive validation`, Phase21/task19.
- Byte-identical focused tagged GREEN `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.996s` (terminated `exit_code: 0`), including actual CLI export validation. Removed temporary direct-Validate diagnostic; original public journey/assertions are unchanged. Production repair requires new final format/lint/full race/coverage/tagged gates; earlier passing full gate remains historical, not final-source closure.

- Rebuilt `/tmp/fid` with `go build -o /tmp/fid ./cmd/fid` (terminated `exit_code: 0`, no output). From `/tmp`, `/tmp/fid --json export validate /tmp/fid-p22-task13-smoke-1T0ybJ/invalid.zip` now returns the complete `archive-invalid` integrity JSON, check `manifest-invalid`, `entry:null`, and observed `exit_code: 6`. This is the expected negative runtime behavior, not a silently successful command. The same command targeting `recovery.zip` returns `status:valid` and matching v2 digest (terminated `exit_code: 0`).
- Final-source format and lint again passed (silent exit0 / `0 issues.` exit0), then full race/coverage attempt3 launched. Source remains frozen during gates.
