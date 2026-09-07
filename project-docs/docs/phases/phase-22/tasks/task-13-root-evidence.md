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
