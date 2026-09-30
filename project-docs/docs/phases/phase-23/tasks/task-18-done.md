# Phase 23 - Task 18: Conformance assembly and wire-surface audit - Recovery report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: Attempt 3 now compares complete accepted result values for all 23 operations and full unknown-field and conflicting-payload failures for every mutation. VAL-068 remains open because per-operation authority, stale version/binding, and not-found cases are still missing; the two existing shared probes do not cover that matrix. The attempt-3 focused test and all five exact gates passed, but passing gates do not close the missing assertion.
- Recovery attempts: 3
- Next action: Continue task 18 from this checkpoint. Add each applicable authority, stale version/binding, and not-found case to the shared 23-operation case list; compare code, category, exit, subject, complete details, and unchanged Event count through both adapters. Compare each amended leaf's first response with its replay, strengthen generated ID/EventRef/digest syntax checks, and record contract-cited N/A for every non-applicable class. Verify provenance for any added fixture owned by an earlier task, then rerun the focused command and all five exact gates before a COMPLETE report. Preserve the current independent expectations and all attempt-1/2/3 evidence.
- Human attention required: no

## Summary of what was implemented

- Added an independently pinned thirty-action corpus. `TestGrantActionSetIsExactlyThirtyNames` compares it separately with the workflow `sortedUniqueCapabilities` map and the journal `const actions`, then compares the two implementations directly. A real activated profile and Envelope grant accept all thirty; five unknown or operation-only names reject without an Event.
- Extended the retained Condition fixture from six to ten types, with 54 accepted and 90 invalid/missing member cases. The per-type shape oracle checks the journal catalog, including the four-field `review-recorded` shape and each exclusive judgment discriminator. The Phase 22 document-snapshot assertion now retains its six historical types and adds the four accepted Phase 23 members.
- Attempt 1 added the shared leaf audit, field-set and replay comparison, and mutation unknown-field probes. Attempt 2 expanded `operations.json` to all 23 operations and made both shared tests consume it. The fake now pins selected contract-authored scalar values for all 23 results. These are partial semantic checks; complete result values, normalization and rejection parity remain open.
- Attempt 3 added independent full result models for all 23 operations from the accepted value grammar and each caller request. One normalization traversal handles only generated row IDs, timestamps, EventRefs, activation/project IDs, and request digests. Cross-case assertions compare raw Requirement, Concern, Review/Finding, judgment, Envelope/Claim, and Evidence identities. Every mutation now compares unknown-field and conflicting-key failures with full category, subject, details, exit, and no new Event. A second valid Receipt makes the `evidence.link` conflict input valid before its idempotency check.
- DEC-265's independent literal marker in `failures.json` classifies every declared row. `TestP23FailureCorpusHasEveryAcceptedCode` asserts 19 declared = 18 added + 1 retained, with `requirement-not-found` the sole retained member.
- Pinned the latest product table names and the built command tree. The test traverses real `--help` output, including commands without a short description, and inspects aliases; it compares the complete table and command inventories with the fixture. No cancelled product family was found. The task's `v1-convergence-inventory.md` section 1 already assigns this proof to task 18, so no correction was needed.
- Regenerated task 14's Human fixture manifest with the current binary. Its eighteen Human/JSON leaf pairs and all checkpoint inspections ran after generation; a fresh-state test ran all four recorded interactive steps at their expected exits. Human comprehension and presentation judgments remain for the Human.

## Files created or modified

- `cmd/fid/integration/p22_closed_values_test.go`, `cmd/fid/integration/p22_conformance_test.go`, `cmd/fid/integration/testdata/p22_conformance/wire-inventory.json`.
- `cmd/fid/integration/p23_conformance_test.go`, `p23_grant_action_audit_test.go`, `p23_cancelled_family_audit_test.go`, and `p23_human_fixtures_test.go`.
- `cmd/fid/integration/testdata/p23_conformance/operations.json`, `grant-actions.json`, `failures.json`, and `product-surface.json`.
- Attempt 2 changed `cmd/fid/integration/p23_conformance_test.go` and the `operations.json` and `failures.json` fixtures; no product file changed.
- Attempt 3 changed `cmd/fid/integration/p23_conformance_test.go` and added seven `p23_*_oracle_test.go` files for contract-authored accepted results; no product file changed.
- This report. The Runner-owned `project-docs/docs/phases/phase-23/run-ledger.md` was not edited.

## Documentation updated

- This report preserves attempts 1 and 2, records attempt 3's accepted-result and two rejection-class evidence, and retains DEC-265. `v1-convergence-inventory.md` section 1 was already correct (I-10).

## Validation Contract Coverage

- VAL-067: PASS for the task-18 audit - `TestGrantActionSetIsExactlyThirtyNames` (`cmd/fid/integration/p23_grant_action_audit_test.go:23`) verifies independent thirty-member profile/journal equality and real allow/reject behavior. The seven Phase 23 names have introducing commits and matching `Refs:` trailers in the provenance table below, all at their owning task. The final tagged gate executed the test.
- VAL-068: PARTIAL - `TestP23SharedLeafConformance` and `TestP23SharedAmendedLeafConformance` in `cmd/fid/integration/p23_conformance_test.go` consume the same 23-row inventory. The new oracle files compare complete accepted result maps under generated-value normalization, with raw cross-case identity checks; both tests also compare every mutation's full unknown-field and conflicting-payload failures and unchanged Event counts. The final focused selection and tagged gate pass. Per-operation authority, stale version/binding, and not-found comparisons, amended first-response/replay equality, and stronger generated-value syntax checks remain unfinished. The 19 declared = 18 added + 1 retained Failure-code partition remains proved by `TestP23FailureCorpusHasEveryAcceptedCode`.
- VAL-069: PASS - `TestP23CancelledFamiliesAbsentFromProduct` (`cmd/fid/integration/p23_cancelled_family_audit_test.go:22`) compares 59 latest-schema tables and 125 real command paths with a pinned full inventory and checks cancelled names and aliases. The exact tagged gate below executes it.
- VAL-070: FIXTURE READY, HUMAN PENDING - regenerated manifest `checkpoints-3135624237/manifest.json` has the mixed-Finding Review, Finding, and resume inspections; `TestP23HumanFixtures` reran each recorded inspection and leaf pair. Only a Human can judge whether resolved, addressed, and still-blocking are understood.
- VAL-071: FIXTURE READY, HUMAN PENDING - the same manifest records blocked action, Human Outcome, and retry steps. `TestP23HumanFixtureStepsRemainRunnable` (`p23_human_fixtures_test.go:165`) ran all three on a fresh protected-Decision checkpoint at exits 5, 0, 0. Human explanation judgment remains pending.
- VAL-072: FIXTURE READY, HUMAN PENDING - the manifest records the exhaustion Stop and Attention inspections; the fresh-state step test ran the exhausted rework denial at exit 4. Human comprehension judgment remains pending.
- VAL-073: FIXTURE READY, HUMAN PENDING - `TestP23HumanFixtures` reran exact Human text for eighteen leaves against the JSON facts, including purpose and qualifiers. Human acceptance of the presentation remains pending.
- VAL-079: PASS for the measured gate and report-inventory part. Attempt 1 recorded ordinary `cmd/fid` 214.615s, uncached `internal/workflow` 55.279s, and tagged `cmd/fid/integration` 229.859s. Attempt 2 recorded `cmd/fid` 202.293s, uncached `internal/workflow` 57.091s, and tagged `cmd/fid/integration` 273.836s. Attempt 3 recorded 199.423s, 57.209s, and 236.466s respectively. All are below the 600/600/900s budgets and below 80% headroom. The attempt-1 read-only inventory over task 1-17 reports found all three times in each report.

The attempt-2 per-operation results are retained below in the historical details block. Their executing commands remain in Tests run and Exact quality-gate results. The earlier producer tests and owners for the 19 declared Failure codes remain pinned in `failures.json`; its literal partition marks 18 `added` and only `requirement-not-found` `retained`.

VAL-067 introduction evidence: for each row, `git log --reverse --format='%H %s%n%b' -S'<name>' -- internal/workflow/activation.go` used the literal name in the first column. Each result printed the listed commit and `Refs:` trailer. The same complete seven-name set is in `public-interface.md` section 4.1 and `grant-actions.json`.

| Grant action (`-S` argument) | First introducing commit | Trailer / owner |
|---|---|---|
| `review.record` | `b500ab45406dd1cbd79153232d2add3ca8a74021` | `Refs: phase-23/task-4` / 4 |
| `finding.address` | `0be596e9068795ceb3a0a89bdac50d538f8e8791` | `Refs: phase-23/task-5` / 5 |
| `concern.record` | `3a0aed9d35ca2d43540e686a5cbbaa245177ea0f` | `Refs: phase-23/task-8` / 8 |
| `concern.dispose` | `3a0aed9d35ca2d43540e686a5cbbaa245177ea0f` | `Refs: phase-23/task-8` / 8 |
| `deviation.record` | `3a0aed9d35ca2d43540e686a5cbbaa245177ea0f` | `Refs: phase-23/task-8` / 8 |
| `decision.record-reference` | `828eb18f422ebbfacbe4e65c62d53a061b7e9b9f` | `Refs: phase-23/task-9` / 9 |
| `human-outcome.record` | `0ff2eeef66d4633309a43ae6a341aecd38def4d0` | `Refs: phase-23/task-10` / 10 |

Attempt 3 executes every row through the same `operations.json` case and the same fake/fid call sequence. `F3` is the final focused command in Tests run; `G3` is the final tagged gate in Exact quality-gate results. The owner column is pinned by `operations.json` and `public-interface.md` section 4. The oracle column names the contract-authored expected-value source. `Full` compares the complete result map, including counts and array order. `U` compares the complete unknown-field failure and unchanged Event count. `D/C` compares replay against the original request where captured, and a different valid payload against the complete idempotency-conflict details; the amended fixture still needs explicit first-response/replay byte equality. `A`, `S`, and `N` mean authority, stale version/binding, and not-found. “Pending” is incomplete work, never N/A.

| Operation | Owner | Accepted value / oracle (F3, G3) | U | D/C | A | S | N |
|---|---:|---|---|---|---|---|---|
| `review.record` | 4 | Full; Review/Finding oracle | PASS | PASS | Pending | Pending | Pending |
| `review.get` | 13 | Full; Review/Finding oracle | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | Pending |
| `review.list` | 13 | Full; Review/Finding oracle, ordered list | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | Pending filter classification |
| `finding.address` | 5 | Full; Review/Finding oracle | PASS | PASS | Pending | Pending | Pending |
| `finding.list` | 5 | Full; Review/Finding oracle, ordered list | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | Pending filter classification |
| `concern.record` | 8 | Full; Concern oracle | PASS | PASS | Pending | N/A §4 (no version/binding input) | Pending |
| `concern.dispose` | 8 | Full; Concern oracle | PASS | PASS | Pending | N/A §4 (no version/binding input) | Pending |
| `concern.list` | 8 | Full; Concern oracle, ordered list | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | Pending filter classification |
| `deviation.record` | 8 | Full; judgment oracle | PASS | PASS | Pending | N/A §4 (no version/binding input) | Pending |
| `deviation.list` | 8 | Full; judgment oracle, ordered list | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | Pending filter classification |
| `decision.record-reference` | 9 | Full; judgment oracle | PASS | PASS | Pending | N/A §4 (no version/binding input) | Pending |
| `decision.list` | 9 | Full; judgment oracle, ordered list | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | Pending filter classification |
| `human-outcome.record` | 10 | Full; judgment oracle | PASS | PASS | Pending | N/A §4 (no version/binding input) | Pending |
| `human-outcome.list` | 10 | Full; judgment oracle, ordered list | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | Pending filter classification |
| `requirement.record-reference` | 12 | Full; Requirement oracle | PASS | PASS | Pending | N/A §4 (no version/binding input) | N/A §4, contracts §12 (new external ID) |
| `requirement.link` | 12 | Full; Requirement oracle | PASS | PASS | Pending | N/A §4 (no version/binding input) | Pending |
| `requirement.get` | 13 | Full; Requirement oracle | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | Pending |
| `requirement.list` | 12 | Full; Requirement oracle, ordered list | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | N/A §§1,4 | Pending filter classification |
| `envelope.create` | 2 | Full; amended Envelope oracle | PASS | Partial first-response replay | Pending | N/A §4 (no version/binding input) | Pending |
| `envelope.activate` | 2,9 | Full; amended Envelope oracle | PASS | Partial first-response replay | Pending | Pending | Pending |
| `claim.acquire` | 3 | Full; amended Claim oracle | PASS | Partial first-response replay | Pending | Pending | Pending |
| `evidence.link` | 6,10 | Full; amended Evidence oracle | PASS | Partial first-response replay | Pending | Pending | Pending |
| `evidence.evaluate` | 6,13 | Full; amended Evidence oracle, ordered evaluations | N/A §§1,4; contracts §13 | N/A §§1,4; contracts §13 | N/A contracts §13 (no Actor) | N/A contracts §13 (no version/binding input) | Pending optional ID |

The Review/Finding, Concern, judgment, Requirement, amended Envelope, Claim, and Evidence oracles are the seven new `cmd/fid/integration/p23_*_oracle_test.go` files. They copy caller values from the fixture requests and add only contract-defined derivations; no fake expectation reads a fid result. Generated values are normalized by `p23NormalizeLeafValue` in `p23_conformance_test.go`; raw mutation/read identities are compared by the `p23Check*Identity` helpers. The duplicate different-payload probe for `evidence.link` creates another valid Receipt before measuring the rejected call. The existing two shared authority/stale probes remain code/exit-only and do not fill the pending per-row cells.

<details>
<summary>Attempt 2 per-operation evidence retained verbatim</summary>

For the table, `U` means an unknown-field rejection compared by operation/code/exit and unchanged Event count; `D` means same-key replay compared by field set and replay flag; `A` and `S` are the two shared fresh-state rejection probes, each compared by operation/code/exit and unchanged Event count. `R` means the necessary full result/failure equality remains unimplemented. A read's mutation-only `U` and `D` are N/A under `public-interface.md` sections 1 and 4; its mutation authority/stale-version checks are likewise N/A. The continuation must classify list filters and every other ID-taking input for not-found behavior under section 4 and the retained Phase 21/22 baseline. This table does not claim full details or generated-value normalization.

| Operation | Owner | Accepted semantic values compared beyond fields | Rejection comparison / remaining |
|---|---:|---|---|
| `review.record` | 4 | verdict `changes-required`, independence `distinct-actor`, Finding `FND-leaf` open and blocking | U, D, S shared; A/not-found/full details R |
| `review.get` | 13 | verdict `changes-required`, Finding addressed after later disposition | U/D/A/S N/A §§1,4; not-found/full details R |
| `review.list` | 13 | first Review verdict and addressed Finding | U/D/A/S N/A §§1,4; ID-filter/full values R |
| `finding.address` | 5 | `FND-leaf` addressed with reason `PRIVATE_FINDING_REASON` | U, D; A/S/not-found/full details R |
| `finding.list` | 5 | `FND-leaf` addressed and blocking | U/D/A/S N/A §§1,4; ID-filter/full values R |
| `concern.record` | 8 | `CON-leaf` open, major, null disposition | U, D; A/not-found/full details R |
| `concern.dispose` | 8 | `CON-leaf` disposed, carried forward with caller reason | U, D; A/not-found/full details R |
| `concern.list` | 8 | `CON-leaf` disposed, carried forward | U/D/A/S N/A §§1,4; ID-filter/full values R |
| `deviation.record` | 8 | `DEV-leaf`, pending decision, caller expected/observed text | U, D; A/not-found/full details R |
| `deviation.list` | 8 | `DEV-leaf`, pending decision | U/D/A/S N/A §§1,4; ID-filter/full values R |
| `decision.record-reference` | 9 | `DREF-leaf`, security, approved | U, D; A/not-found/full details R |
| `decision.list` | 9 | `DREF-leaf`, security, approved | U/D/A/S N/A §§1,4; ID-filter/full values R |
| `human-outcome.record` | 10 | `OUT-leaf`, security, approved | U, D; A/not-found/full details R |
| `human-outcome.list` | 10 | `OUT-leaf`, security, approved | U/D/A/S N/A §§1,4; ID-filter/full values R |
| `requirement.record-reference` | 12 | active, null supersession | U, D, A shared; not-found/full details R |
| `requirement.link` | 12 | `REQ-leaf`, implements | U, D; A/not-found/full details R |
| `requirement.get` | 13 | active Requirement, `REQ-leaf` implements link | U/D/A/S N/A §§1,4; not-found/full details R |
| `requirement.list` | 12 | first Requirement active | U/D/A/S N/A §§1,4; ID-filter/full values R |
| `envelope.create` | 2 | caller key and current revision 1 | U, D; A/not-found/full details R |
| `envelope.activate` | 2,9 | caller key | U, D; A/S/not-found/full details R |
| `claim.acquire` | 3 | caller key and mutation mode | U, D; A/not-found/full details R |
| `evidence.link` | 6,10 | caller key and receipt type | U, D; A/not-found/full details R |
| `evidence.evaluate` | 6,13 | first evaluation Obligation `OBL-p22-corpus` | U/D/A/S N/A §§1,4 and contracts §13; optional Obligation not-found/full values R |

</details>

The nineteen `failures.json` rows each name one pre-existing producing `Test...` declaration and its Phase 23 owner task. They are: `review-not-independent` (4), `review-verdict-invalid` (5), `review-changes-required` (6), `review-rework-exhausted` (6), `finding-already-disposed` (5), `finding-not-found` (5), `review-not-found` (13), `concern-not-found` (8), `deviation-not-found` (8), `outcome-not-found` (10), `requirement-not-found` (12), `concern-already-disposed` (8), `claim-mode-invalid` (3), `protected-decision-required` (9), `human-outcome-required` (10), `outcome-not-permitted` (10), `evidence-type-not-accepted` (6), `source-digest-mismatch` (9), and `repository-mismatch` (4). The tagged gate executes their named producer tests; `assertP22ProcessFailure` calls the pinned `p23FailureContracts` fake envelope oracle for those tests.

## Tests run

- Attempt 3 initial focused command: `go test -tags integration -run '^(TestP23FailureCorpusHasEveryAcceptedCode|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestGrantActionSetIsExactlyThirtyNames)$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.979s` (terminated `exit_code: 0`).
- Attempt 3 full-result RED: `go test -tags integration -run '^TestP23SharedLeafConformance$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`), because the old fake held only selected paths. After contract-authored Requirement, Concern, judgment, and Review/Finding models were added, the byte-identical rerun passed: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.373s` (terminated `exit_code: 0`). An intermediate Review/Finding run failed `FAIL` (`exit_code: 1`) because a null input collection projects as the contract's empty `resolved_finding_ids` array.
- Attempt 3 amended-result RED: `go test -tags integration -run '^TestP23SharedAmendedLeafConformance/claim.acquire$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); the full Claim oracle later passed in the complete amended test. Focused `evidence.evaluate`, `envelope.create`, `envelope.activate`, and `evidence.link` selections each failed `FAIL` (`exit_code: 1`) before their contract-authored full models were added. `go test -tags integration -run '^TestP23SharedAmendedLeafConformance$' -count=1 ./cmd/fid/integration/` finally passed: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.659s` (terminated `exit_code: 0`).
- Attempt 3 unknown-field details RED: `go test -tags integration -run '^TestP23SharedLeafConformance/requirement.record-reference$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`) because the fake lacked category and full details. After adding the closed `invalid-argument` details, the combined shared tests passed: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.508s` (terminated `exit_code: 0`).
- Attempt 3 conflicting-payload RED: the same Requirement selection failed `FAIL` (`exit_code: 1`) because `operation_scope` includes a generated Project ID; shared normalization corrected it. The combined shared tests first failed `FAIL` (`exit_code: 1`) for invalid changed inputs to `envelope.activate` and `evidence.link`. A changed Actor exercised the retained Envelope digest rule, and a second valid Receipt exercised the Evidence digest rule. The final amended shared test passed as above. These were fixture corrections, not product changes.
- Attempt 3 final focused command: `go test -tags integration -run '^(TestP23FailureCorpusHasEveryAcceptedCode|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestGrantActionSetIsExactlyThirtyNames)$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.663s` (terminated `exit_code: 0`).
- Attempt 3 touched tagged race: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP23FailureCorpusHasEveryAcceptedCode|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestGrantActionSetIsExactlyThirtyNames)$' ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.876s` (terminated `exit_code: 0`).
- Attempt 3 uncached workflow measurement: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	57.209s` (terminated `exit_code: 0`).

- Attempt 2 initial focused command: `go test -tags integration -run '^(TestP23FailureCorpusHasEveryAcceptedCode|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestGrantActionSetIsExactlyThirtyNames)$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.028s` (terminated `exit_code: 0`).
- Partition RED: `go test -tags integration -run '^TestP23FailureCorpusHasEveryAcceptedCode$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`), first row lacked a marker. Its byte-identical rerun passed: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.647s` (terminated `exit_code: 0`).
- New-leaf semantic RED: `go test -tags integration -run '^TestP23SharedLeafConformance$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`) for the four missing Requirement expectations; a later identical run failed `FAIL` (`exit_code: 1`) for the other fourteen. It passed after all eighteen had literal expected paths: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.335s` (terminated `exit_code: 0`).
- Amended-leaf semantic RED: `go test -tags integration -run '^TestP23SharedAmendedLeafConformance$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`) after the case type gained a field but positional cases had not been extended. After that correction, the same command failed `FAIL` (`exit_code: 1`) for five missing expectations. An intermediate run failed `FAIL` (`exit_code: 1`) on `envelope.revision`; the retained Phase 21 contract names `current_revision`. Its final identical rerun passed: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.535s` (terminated `exit_code: 0`).
- Shared-inventory RED: `go test -tags integration -run '^TestP23SharedLeafConformance$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`) when the 23-row assertion met the old 18-row fixture. The combined shared-test command first failed `FAIL` (`exit_code: 1`) while the Human fixture count still compared against all 23; after selecting its eighteen new leaves, `go test -tags integration -run '^(TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance)$' -count=1 ./cmd/fid/integration/` passed: `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.459s` (terminated `exit_code: 0`).
- Attempt 2 final focused command: `go test -tags integration -run '^(TestP23FailureCorpusHasEveryAcceptedCode|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestGrantActionSetIsExactlyThirtyNames)$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.400s` (terminated `exit_code: 0`).
- Attempt 2 touched tagged race: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestP23FailureCorpusHasEveryAcceptedCode|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestGrantActionSetIsExactlyThirtyNames)$' ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.641s` (terminated `exit_code: 0`).
- Attempt 2 uncached workflow measurement: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	57.091s` (terminated `exit_code: 0`).

- First RED: `go test -tags integration -run '^TestP23ConditionWireGrammarMatchesJournalAndFixtures$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); the retained grammar had six Condition types. After extending the grammar and fixture, the combined focused `TestP23ConditionWireGrammarMatchesJournalAndFixtures|TestP22ClosedValueSchemaCases` command passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.510s` (terminated `exit_code: 0`).
- First grant RED: `go test -tags integration -run '^TestGrantActionSetIsExactlyThirtyNames$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); the pinned corpus file was absent. Its final same-command rerun passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.179s` (terminated `exit_code: 0`). An intermediate run failed `FAIL` (`exit_code: 1`) because the test passed a nil artifact-path array into the existing Activation fixture; `[]string{}` corrected that fixture input.
- Shared leaf RED: `go test -tags integration -run '^TestP23SharedLeafConformance$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); `operations.json` was absent. Its final same-command rerun passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.280s` (terminated `exit_code: 0`). Intermediate authority and stale-version attempts failed `FAIL` (`exit_code: 1`) because the test reused a released Claim and expected an undeclared `run-version-conflict` code; the fresh review-mode fixture now asserts the existing `invalid-argument` stale binding rejection with no Event.
- Retained snapshot RED: `go test -tags integration -run '^TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); it compared the six-type Phase 22 prose directly with the ten-type current fixture. The byte-identical rerun passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.505s` (terminated `exit_code: 0`). The assertion was introduced in `24e3cbdb` (`Refs: phase-22/task-18`), and the Condition fixture in `e04bfbdd` (`Refs: phase-22/task-18`); the accepted Phase 23 four-type extension invalidated only the old exact-six expectation.
- Amended leaf first focused attempt: `go test -tags integration -run '^TestP23SharedAmendedLeafConformance$' -count=1 ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`); the fixture's gate Run had unknown token-budget state. A baseline budget observation before the receipt link made the same command pass with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.489s` (terminated `exit_code: 0`).
- Corpus and surface first focused attempts each ended `FAIL` (`exit_code: 1`) because their pinned `failures.json` and `product-surface.json` files were absent. `go test -tags integration -run '^TestP23FailureCorpusHasEveryAcceptedCode$' -count=1 ./cmd/fid/integration/` later passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	0.502s` (terminated `exit_code: 0`); `go test -tags integration -run '^TestP23CancelledFamiliesAbsentFromProduct$' -count=1 ./cmd/fid/integration/` passed with `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.731s` (terminated `exit_code: 0`).
- Touched tagged race: `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^(TestGrantActionSetIsExactlyThirtyNames|TestP23ConditionWireGrammarMatchesJournalAndFixtures|TestP22ClosedValueSchemaCases|TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestP23FailureCorpusHasEveryAcceptedCode|TestP23CancelledFamiliesAbsentFromProduct|TestP23HumanFixtureStepsRemainRunnable)$' ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	9.054s` (terminated `exit_code: 0`).
- Fixture regeneration: `FID_P23_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-23-task-18-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	27.623s` (terminated `exit_code: 0`); test log named `/home/mark/projects/fid/.cache/phase-23-task-18-human/checkpoints-3135624237/manifest.json`.
- Interactive fixture steps: `go test -tags integration -run '^TestP23HumanFixtureStepsRemainRunnable$' -count=1 ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.853s` (terminated `exit_code: 0`).
- Uncached workflow time: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	55.279s` (terminated `exit_code: 0`).
- Built entrypoint: `go build -o /tmp/fid-p23-audit ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`); `/tmp/fid-p23-audit review --help` - PASS - `Use "fid review [command] --help" for more information about a command.` (terminated `exit_code: 0`).

## Exact quality-gate results

Attempt 1 exact commands ran from `/home/mark/projects/fid` in task-spec order on its then-final test tree. The ordinary workflow package result was cached; the attempt-1 uncached measurement appears under Tests run. Every attempt-1 result below is retained unchanged.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	229.859s` (terminated `exit_code: 0`)

Attempt 2 exact commands, run in task-spec order from `/home/mark/projects/fid` on the partial semantic fixture:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)			60.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	273.836s` (terminated `exit_code: 0`)

Attempt 3 exact commands, run in task-spec order from `/home/mark/projects/fid` on the complete 23-result model with an incomplete authority/stale/not-found matrix:

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`); the same output reports `cmd/fid` at 199.423s.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)			60.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	236.466s` (terminated `exit_code: 0`)

## Deviations from spec

- Attempt 3 remains incomplete: all 23 accepted values and unknown-field/conflicting-payload failures now use independent fake expectations, but per-operation authority, stale version/binding, and not-found classes remain. The first full-result focused attempt failed `FAIL` (`exit_code: 1`); selected-path comparison caused it, full contract-authored models repaired it, and the later identical run passed `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	6.373s` (`exit_code: 0`). The first changed-payload amended run failed `FAIL` (`exit_code: 1`) because its Evidence input named an unbound Receipt; creating a second valid Receipt made the later identical amended command pass `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	1.659s` (`exit_code: 0`).
- Attempt 1 deviation: the full semantic fake/fid equality required by VAL-068 was unfinished. The fake then modeled top-level fields, replay and selected rejection codes. Attempt 2 added partial literal semantic checks for all 23 operations, but complete values and rejection classes remain unfinished. A same-task continuation must finish them without deriving fake expectations from fid output.
- Attempt 1 protected discrepancy: section 5/DEC-230's 18-versus-19 wording prevented an exact Failure-code claim. DEC-265 resolved it without a member or semantic change: 19 declared, 18 Phase-23-added, one retained `requirement-not-found` from `ab4ac8f` (`Refs: phase-20/task-1`). Attempt 2 pins and tests that partition. The original failed evidence remains above.
- The first `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-18-done.md project-docs/docs/phases/phase-23/tasks/task-18.md` attempt reported `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`): the Recovery attempts field included prose after `0`. This report keeps the attempt count as the required integer and retains the explanation in the outcome and concerns sections.

## Concerns, technical debt, or notes for the Tech Lead

- Attempt 1 concern is resolved by DEC-265. The open task-18 concern is the remaining VAL-068 authority, stale version/binding, and not-found matrix, explicit first-response/replay equality for the five amended leaves, stronger generated ID/EventRef/digest syntax checks, and definitive N/A classifications for list ID filters. No product divergence has been established by the accepted, unknown-field, or conflicting-payload comparisons.
- Historical task reports show `cmd/fid` package times at or above 80% of the 600s budget: task 2's later 490.783s, task 4's 496.639s, task 5's 588.268s, tasks 6-9's 535.185/499.695/538.754/513.222s, task 10's 564.337s attempt, and task 11's 531.533s final pass. This task's 214.615s and tagged 229.859s have ample headroom; the earlier measurements remain phase-level warnings.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `2e7fd2f9fa05bce648510789453c8d2ed8633b1d`

## Commit subject

- `test(conformance): compare Phase 23 leaf results`
