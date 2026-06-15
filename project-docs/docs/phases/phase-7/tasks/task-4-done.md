## Summary of what was implemented

- Added [`cmd/fid/e2e_v01c_ac_test.go`](/home/mark/projects/fid/cmd/fid/e2e_v01c_ac_test.go) as a new integration smoke for Phase 7 Task 4.
- The test proves `all_acs_verified` with real `fid issue ac add` and `fid issue ac verify` commands on the `qa_review -> qa_pass` edge, rejects before append when `qa_lead` passes are missing, then succeeds once the required real pass verifications exist.
- The same test snapshots `acceptance_criteria`, `ac_verifications`, `fid issue get --section acceptance --format json`, and `fid issue ac list --format json`, deletes `.fid/fid.sqlite`, runs `fid rebuild`, and proves parity after replay.
- No production defect surfaced. The required behavior was already present; this task closes the real-command end-to-end evidence gap left by the Phase 5 fixture-backed proof.
- To reach `qa_review` without changing the workflow graph or pulling Phase 8 into scope, the test inserts a minimal `task_completions` row directly into SQLite. No AC journal fixtures were used.
- Phase 5 fixture-based gating evidence is now superseded by this real-command transcript and integration test.

## Files created or modified

- `cmd/fid/e2e_v01c_ac_test.go`
- `project-docs/docs/phases/phase-7/tasks/task-4-done.md`

## Documentation updated

- Added this completion report at `project-docs/docs/phases/phase-7/tasks/task-4-done.md`.

## Validation Contract Coverage

- `VAL-008` - PASS. `TestV01CAllACsVerifiedGateAndRebuildParity` proves a real `qa_review -> qa_pass` transition rejects with unmet AC IDs and missing `qa_lead` role, appends no journal event, leaves issue state unchanged, then succeeds after real `fid issue ac verify ... --by qa_lead --status pass` writes.
- `VAL-009` - PASS. The same test snapshots projected `acceptance_criteria`, projected `ac_verifications`, and latest-per-role read-back via both `fid issue get --section acceptance --format json` and `fid issue ac list --format json`, deletes `.fid/fid.sqlite`, runs `fid rebuild`, and proves byte-for-byte parity on the read-back plus structural parity on both tables.
- Built-binary end-to-end transcript used for this report:

```text
+ /tmp/fid init
initialized fid project in .fid/
  prefix: FID-TASK4-TRANSCRIPT-CSKVCG
  config: .fid/config.toml
  pricing: .fid/pricing.toml
  journal: .fid/journal/2026-06-14.jsonl
  event: EVT-20260614-000001 project.initialized
  roles seeded: developer, qa_reviewer, qa_lead, tl_lite, tl, workflow_runner, human
+ /tmp/fid issue create --kind task --title Phase 7 AC gate parity --by developer
created issue FID-TASK4-TRANSCRIPT-CSKVCG-1a500b (task) "Phase 7 AC gate parity"
issue_id=FID-TASK4-TRANSCRIPT-CSKVCG-1a500b
+ /tmp/fid issue transition FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --to in_progress --by developer
transitioned FID-TASK4-TRANSCRIPT-CSKVCG-1a500b: todo → in_progress
+ sqlite3 insert completion row
+ /tmp/fid issue transition FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --to dev_complete --by developer
transitioned FID-TASK4-TRANSCRIPT-CSKVCG-1a500b: in_progress → dev_complete
+ /tmp/fid issue transition FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --to qa_review --by qa_lead
transitioned FID-TASK4-TRANSCRIPT-CSKVCG-1a500b: dev_complete → qa_review
+ /tmp/fid issue ac add FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --text First QA approval criterion --ordinal 1 --by developer
added AC AC-5d3267 on FID-TASK4-TRANSCRIPT-CSKVCG-1a500b (ordinal 1): "First QA approval criterion"
ac1=AC-5d3267
+ /tmp/fid issue ac add FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --text Second QA approval criterion --ordinal 2 --by developer
added AC AC-953a34 on FID-TASK4-TRANSCRIPT-CSKVCG-1a500b (ordinal 2): "Second QA approval criterion"
ac2=AC-953a34
journal_lines_before_reject=7
+ /tmp/fid issue get FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --format json
{"id":"FID-TASK4-TRANSCRIPT-CSKVCG-1a500b","kind":"task","title":"Phase 7 AC gate parity","status":"qa_review","status_category":"started","workflow_id":"task-default","parent_id":null,"priority":null,"sequence":null,"spec_md_path":null,"est_tokens":null,"due_date":null,"owner_role":"developer","owner_actor":null,"labels":[],"in_inbox":false,"concerns_raised":false,"open_concerns_count":0,"triage_status":null,"acceptance_criteria":[{"id":"AC-5d3267","ordinal":1,"text":"First QA approval criterion","latest_verifications":{"developer":null}},{"id":"AC-953a34","ordinal":2,"text":"Second QA approval criterion","latest_verifications":{"developer":null}}],"created_at":"2026-06-14T08:42:24Z","updated_at":"2026-06-14T08:42:27Z"}
+ /tmp/fid issue transition FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --to qa_pass --by qa_lead
error: transition issue "FID-TASK4-TRANSCRIPT-CSKVCG-1a500b" from "qa_review" to "qa_pass" blocked by workflow rules:
  - all_acs_verified: entities AC-5d3267; missing roles qa_lead
  - all_acs_verified: entities AC-953a34; missing roles qa_lead
reject_exit=1
journal_lines_after_reject=7
+ /tmp/fid issue get FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --format json
{"id":"FID-TASK4-TRANSCRIPT-CSKVCG-1a500b","kind":"task","title":"Phase 7 AC gate parity","status":"qa_review","status_category":"started","workflow_id":"task-default","parent_id":null,"priority":null,"sequence":null,"spec_md_path":null,"est_tokens":null,"due_date":null,"owner_role":"developer","owner_actor":null,"labels":[],"in_inbox":false,"concerns_raised":false,"open_concerns_count":0,"triage_status":null,"acceptance_criteria":[{"id":"AC-5d3267","ordinal":1,"text":"First QA approval criterion","latest_verifications":{"developer":null}},{"id":"AC-953a34","ordinal":2,"text":"Second QA approval criterion","latest_verifications":{"developer":null}}],"created_at":"2026-06-14T08:42:24Z","updated_at":"2026-06-14T08:42:27Z"}
issue_state_after_reject=UNCHANGED
+ /tmp/fid issue ac verify FID-TASK4-TRANSCRIPT-CSKVCG-1a500b AC-5d3267 --by qa_lead --status pass
verified AC-5d3267 on FID-TASK4-TRANSCRIPT-CSKVCG-1a500b by qa_lead: pass
+ /tmp/fid issue ac verify FID-TASK4-TRANSCRIPT-CSKVCG-1a500b AC-953a34 --by qa_lead --status pass
verified AC-953a34 on FID-TASK4-TRANSCRIPT-CSKVCG-1a500b by qa_lead: pass
+ /tmp/fid issue transition FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --to qa_pass --by qa_lead
transitioned FID-TASK4-TRANSCRIPT-CSKVCG-1a500b: qa_review → qa_pass
+ /tmp/fid issue get FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --format json
{"id":"FID-TASK4-TRANSCRIPT-CSKVCG-1a500b","kind":"task","title":"Phase 7 AC gate parity","status":"qa_pass","status_category":"completed","workflow_id":"task-default","parent_id":null,"priority":null,"sequence":null,"spec_md_path":null,"est_tokens":null,"due_date":null,"owner_role":"developer","owner_actor":null,"labels":[],"in_inbox":false,"concerns_raised":false,"open_concerns_count":0,"triage_status":null,"acceptance_criteria":[{"id":"AC-5d3267","ordinal":1,"text":"First QA approval criterion","latest_verifications":{"qa_lead":{"status":"pass","verified_at":"2026-06-14T08:42:24Z"}}},{"id":"AC-953a34","ordinal":2,"text":"Second QA approval criterion","latest_verifications":{"qa_lead":{"status":"pass","verified_at":"2026-06-14T08:42:24Z"}}}],"created_at":"2026-06-14T08:42:24Z","updated_at":"2026-06-14T08:42:28Z"}
+ sqlite3 snapshot acceptance_criteria before rebuild
[{"id":"AC-5d3267","issue_id":"FID-TASK4-TRANSCRIPT-CSKVCG-1a500b","ordinal":1,"text":"First QA approval criterion"},
{"id":"AC-953a34","issue_id":"FID-TASK4-TRANSCRIPT-CSKVCG-1a500b","ordinal":2,"text":"Second QA approval criterion"}]
+ sqlite3 snapshot ac_verifications before rebuild
[{"id":"ACV-858ead","ac_id":"AC-953a34","verifier_role":"qa_lead","status":"pass","evidence_type":"","evidence_ref":"","notes":"","completion_id":"","verified_at":"2026-06-14T08:42:24Z"},
{"id":"ACV-a91b3c","ac_id":"AC-5d3267","verifier_role":"qa_lead","status":"pass","evidence_type":"","evidence_ref":"","notes":"","completion_id":"","verified_at":"2026-06-14T08:42:24Z"}]
+ /tmp/fid issue get FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --section acceptance --format json
[{"id":"AC-5d3267","ordinal":1,"text":"First QA approval criterion","latest_verifications":{"qa_lead":{"status":"pass","verified_at":"2026-06-14T08:42:24Z"}}},{"id":"AC-953a34","ordinal":2,"text":"Second QA approval criterion","latest_verifications":{"qa_lead":{"status":"pass","verified_at":"2026-06-14T08:42:24Z"}}}]
+ /tmp/fid issue ac list FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --format json
[{"id":"AC-5d3267","ordinal":1,"text":"First QA approval criterion","latest_verifications":{"qa_lead":{"status":"pass","verified_at":"2026-06-14T08:42:24Z"}}},{"id":"AC-953a34","ordinal":2,"text":"Second QA approval criterion","latest_verifications":{"qa_lead":{"status":"pass","verified_at":"2026-06-14T08:42:24Z"}}}]
+ rm /tmp/fid-task4-transcript-CskvCG/.fid/fid.sqlite
+ /tmp/fid rebuild --by developer
rebuilding .fid/fid.sqlite from journal...
  replayed 10 events from 1 journal files
  schema version: 1
  elapsed: 14ms
rebuild complete
+ sqlite3 snapshot acceptance_criteria after rebuild
[{"id":"AC-5d3267","issue_id":"FID-TASK4-TRANSCRIPT-CSKVCG-1a500b","ordinal":1,"text":"First QA approval criterion"},
{"id":"AC-953a34","issue_id":"FID-TASK4-TRANSCRIPT-CSKVCG-1a500b","ordinal":2,"text":"Second QA approval criterion"}]
+ sqlite3 snapshot ac_verifications after rebuild
[{"id":"ACV-858ead","ac_id":"AC-953a34","verifier_role":"qa_lead","status":"pass","evidence_type":"","evidence_ref":"","notes":"","completion_id":"","verified_at":"2026-06-14T08:42:24Z"},
{"id":"ACV-a91b3c","ac_id":"AC-5d3267","verifier_role":"qa_lead","status":"pass","evidence_type":"","evidence_ref":"","notes":"","completion_id":"","verified_at":"2026-06-14T08:42:24Z"}]
+ /tmp/fid issue get FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --section acceptance --format json
[{"id":"AC-5d3267","ordinal":1,"text":"First QA approval criterion","latest_verifications":{"qa_lead":{"status":"pass","verified_at":"2026-06-14T08:42:24Z"}}},{"id":"AC-953a34","ordinal":2,"text":"Second QA approval criterion","latest_verifications":{"qa_lead":{"status":"pass","verified_at":"2026-06-14T08:42:24Z"}}}]
+ /tmp/fid issue ac list FID-TASK4-TRANSCRIPT-CSKVCG-1a500b --format json
[{"id":"AC-5d3267","ordinal":1,"text":"First QA approval criterion","latest_verifications":{"qa_lead":{"status":"pass","verified_at":"2026-06-14T08:42:24Z"}}},{"id":"AC-953a34","ordinal":2,"text":"Second QA approval criterion","latest_verifications":{"qa_lead":{"status":"pass","verified_at":"2026-06-14T08:42:24Z"}}}]
acceptance_criteria_parity=OK
ac_verifications_parity=OK
issue_get_acceptance_parity=OK
ac_list_parity=OK
```

## Tests run

- `go test -tags integration ./cmd/fid -run TestV01CAllACsVerifiedGateAndRebuildParity -count=1`
  - FAIL before helper implementation (`undefined: insertTask4CompletionRow`, snapshot helper compile failures)
- `go test -tags integration ./cmd/fid -run TestV01CAllACsVerifiedGateAndRebuildParity -count=1`
  - PASS
- `gofmt -w cmd/fid/e2e_v01c_ac_test.go`
  - PASS
- `go test -tags integration ./cmd/fid -run TestV01CAllACsVerifiedGateAndRebuildParity -count=1`
  - PASS after formatting
- Built-binary transcript run with `/tmp/fid` plus a direct SQLite `task_completions` setup insert
  - PASS

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` - PASS (`cmd/fid` 78.9%, `internal/issues` 82.0%, `internal/workflow` 95.5%)
- `go tool cover -func=coverage.out` - PASS (`total: 76.7%`; required package thresholds satisfied)
- `go test -tags integration -race ./...` - PASS (required because `cmd/fid/e2e_v01c_ac_test.go` is behind the `integration` build tag)
- `go build -o /tmp/fid ./cmd/fid` - PASS (required because this report includes a real CLI transcript)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `a7cf2bd`

## Commit subject

- `test(e2e): add AC gating rebuild parity smoke`
