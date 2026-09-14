# Task 14 Delegate Notes

Root preparation while Task 13 gates run. This is guidance under the accepted
Task14 spec/DEC-202, not a replacement contract or completion evidence.
Implementation must wait for Task13 COMPLETE in the ledger.

## Ownership And Execution

One Developer owns the whole Task14 implementation, matching tests and report.
No subordinate agents or concurrent root source edits. Complete a coherent task;
do not return after tiny implementation chunks. Preserve red/green and every
failed gate, then run all five exact commands in order. A silent full Go gate
may take minutes: observe termination and never kill it as presumed hung.
Root handles administrative docs, report lint and any later scope question.

## Existing Building Blocks

- Task13 `Validate(ValidateOptions{ArchivePath,MigrationsDir})` is read-only.
  V2 validation uses trusted migrations and normal appliers in private SQLite.
  Import must reject v1 rather than convert it.
- Task13 ZIP scanner enforces actual decompression, count, path and line limits.
  Keep validation and staged extraction tied to the same archive content; do
  not validate one path snapshot then publish different bytes after replacement.
- `archiveProjectConfig.projectConfig` yields a strict portable config.
  `config.WriteBootstrapFiles` can materialize the normal staged layout.
- `journal.RebuildWorkspace` applies normal migrations and ordered Event files,
  creates/closes its DB and holds the normal lock inside the staged `.fid`.
  `db.Open(":memory:")` is a file path; Task13 deliberately uses sql.Open for its
  private in-memory semantic projection. Import itself needs a real staged DB.
- Ordinal chunk names sort before subsequently appended dated journal files.
  Preserve canonical Event bytes/order/IDs/hash/idempotency; physical filenames
  are delivery layout. No import domain Event or SQLite archive image.
- Artifact bodies must use imported configuration's storage root, preserve
  digest/availability and retain unavailable metadata without invented bytes.

## Publication

Target is an existing safe directory with no `.fid`; Git/source files are
permitted and must remain byte-identical. Reject symlinked/unsafe destinations
and every existing `.fid`, including an empty directory, except a verified
same-import response-loss retry. Use an exclusive target-specific lock, stage
on the target filesystem, sync files/directories, atomically publish without
replacement, then sync the parent. Inject faults on both sides of each boundary.

Plain POSIX os.Rename may replace an existing empty directory; account for this
when implementing no-replace publication. Existing go.mod includes x/sys v0.42.0
indirectly; if its platform primitives are used directly, record the need,
stdlib alternatives and license rationale before promotion. No new module or
version is authorized merely by this note. Check accepted platform constraints.

Marker is private delivery metadata, excluded from canonical export. Bind the
installed project/history, exact archive digest/key and actual destination.
Same digest/key retry replays installed facts after response loss; changed
requests conflict. Do not use a marker's mere existence as authority.

## Imported Destination Qualification (DEC-202)

- Preserve historical Project ID, Envelope/Claim/Run repository values and
  Events. New destination observations/Claim/Run are new facts.
- `ReconcileRepository` currently rejects every changed Worktree ID and copies
  the historical Worktree ID into its result. Ordinary callers must retain that
  rejection. Only the verified installed-import association may qualify the
  current destination, while all ancestry/path/dirty-state checks remain.
- `InspectRun` composes reconciliation and liveness. `RecoverRun` currently
  copies priorClaim.WorktreeID into successorClaim and invokes reconciliation
  directly. Successor facts must bind the actual qualified destination.
- Imported observation idempotency replay returns original historical facts.
  It cannot satisfy the fresh local liveness requirement. Guard recovery using
  the installed history boundary, not caller timestamps alone.
- `run.go`, `claims/generic.go`, `evaluation.go` and `evidence/execution.go`
  currently compare claims/registrations/runs to historical Envelope Worktree
  bindings. Change only the DEC-202 imported route; preserve Project ID,
  ordinary foreign Worktree rejection and exact current target/freshness checks.
- Import alone must deny historical Claim/Run mutation or execution; recovery
  needs trusted current liveness, exact reconciliation, valid grants, existing
  Stop/Budget checks and current Claim conditions. No spawn on import.

Task15 owns CLI import and shared fake/fid/parity corpus; Task14 still proves
its complete domain/public API behavior and retained integration gate. Task15
DEC-203 adds only its already accepted import failure registry files. Reviews
remain serialized: Review5 after Task15, Review6 after Task18, final independent
QA after the final corpus and retained Human fixtures.
