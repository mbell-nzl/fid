# Phase 23 - Task 18: Conformance assembly and wire-surface audit - Recovery report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: in-scope-implementation
- Condition: DEC-265 resolves the Failure-code count and the corpus now pins 19 declared = 18 added + 1 retained. VAL-068 still lacks complete semantic fake/fid result equality and the per-operation rejection-class matrix required by the task handoff. The attempt-2 gates pass, but passing gates do not close that missing assertion.
- Recovery attempts: 2
- Next action: Continue task 18 from this checkpoint. Build independently authored complete expected results and one structural normalization for generated values, use the shared 23-operation case list to compare every applicable rejection class with full failure details and Event invariance, then rerun the focused command and all five exact gates before a COMPLETE report. Do not edit product code or derive fake expectations from fid output.
- Human attention required: no

## Summary of what was implemented

- Added an independently pinned thirty-action corpus. `TestGrantActionSetIsExactlyThirtyNames` compares it separately with the workflow `sortedUniqueCapabilities` map and the journal `const actions`, then compares the two implementations directly. A real activated profile and Envelope grant accept all thirty; five unknown or operation-only names reject without an Event.
- Extended the retained Condition fixture from six to ten types, with 54 accepted and 90 invalid/missing member cases. The per-type shape oracle checks the journal catalog, including the four-field `review-recorded` shape and each exclusive judgment discriminator. The Phase 22 document-snapshot assertion now retains its six historical types and adds the four accepted Phase 23 members.
- Attempt 1 added the shared leaf audit, field-set and replay comparison, and mutation unknown-field probes. Attempt 2 expanded `operations.json` to all 23 operations and made both shared tests consume it. The fake now pins selected contract-authored scalar values for all 23 results. These are partial semantic checks; complete result values, normalization and rejection parity remain open.
- DEC-265's independent literal marker in `failures.json` classifies every declared row. `TestP23FailureCorpusHasEveryAcceptedCode` asserts 19 declared = 18 added + 1 retained, with `requirement-not-found` the sole retained member.
- Pinned the latest product table names and the built command tree. The test traverses real `--help` output, including commands without a short description, and inspects aliases; it compares the complete table and command inventories with the fixture. No cancelled product family was found. The task's `v1-convergence-inventory.md` section 1 already assigns this proof to task 18, so no correction was needed.
- Regenerated task 14's Human fixture manifest with the current binary. Its eighteen Human/JSON leaf pairs and all checkpoint inspections ran after generation; a fresh-state test ran all four recorded interactive steps at their expected exits. Human comprehension and presentation judgments remain for the Human.

## Files created or modified

- `cmd/fid/integration/p22_closed_values_test.go`, `cmd/fid/integration/p22_conformance_test.go`, `cmd/fid/integration/testdata/p22_conformance/wire-inventory.json`.
- `cmd/fid/integration/p23_conformance_test.go`, `p23_grant_action_audit_test.go`, `p23_cancelled_family_audit_test.go`, and `p23_human_fixtures_test.go`.
- `cmd/fid/integration/testdata/p23_conformance/operations.json`, `grant-actions.json`, `failures.json`, and `product-surface.json`.
- Attempt 2 changed `cmd/fid/integration/p23_conformance_test.go` and the `operations.json` and `failures.json` fixtures; no product file changed.
- This report. The Runner-owned `project-docs/docs/phases/phase-23/run-ledger.md` was not edited.

## Documentation updated

- This report preserves attempt 1's evidence and records DEC-265, attempt 2's partial semantic audit, provenance, and exact gates. `v1-convergence-inventory.md` section 1 was already correct (I-10).

## Validation Contract Coverage

- VAL-067: PASS for the task-18 audit - `TestGrantActionSetIsExactlyThirtyNames` (`cmd/fid/integration/p23_grant_action_audit_test.go:23`) verifies independent thirty-member profile/journal equality and real allow/reject behavior. The seven Phase 23 names have introducing commits and matching `Refs:` trailers in the provenance table below, all at their owning task. The final tagged gate executed the test.
- VAL-068: PARTIAL - `TestP23SharedLeafConformance` and `TestP23SharedAmendedLeafConformance` in `cmd/fid/integration/p23_conformance_test.go` consume the same 23-row inventory and compare top-level fields, replay, unknown-field code/exit on mutations, and pinned semantic paths. The final focused selection and tagged gate pass. The fake still omits many deterministic nested values and generated-value consistency checks; several rejection classes are missing or only tested in separate producer journeys. The 19/18/1 Failure-code partition is now proved by `TestP23FailureCorpusHasEveryAcceptedCode`, but VAL-068 cannot pass until all result and rejection comparisons are complete.
- VAL-069: PASS - `TestP23CancelledFamiliesAbsentFromProduct` (`cmd/fid/integration/p23_cancelled_family_audit_test.go:22`) compares 59 latest-schema tables and 125 real command paths with a pinned full inventory and checks cancelled names and aliases. The exact tagged gate below executes it.
- VAL-070: FIXTURE READY, HUMAN PENDING - regenerated manifest `checkpoints-3135624237/manifest.json` has the mixed-Finding Review, Finding, and resume inspections; `TestP23HumanFixtures` reran each recorded inspection and leaf pair. Only a Human can judge whether resolved, addressed, and still-blocking are understood.
- VAL-071: FIXTURE READY, HUMAN PENDING - the same manifest records blocked action, Human Outcome, and retry steps. `TestP23HumanFixtureStepsRemainRunnable` (`p23_human_fixtures_test.go:165`) ran all three on a fresh protected-Decision checkpoint at exits 5, 0, 0. Human explanation judgment remains pending.
- VAL-072: FIXTURE READY, HUMAN PENDING - the manifest records the exhaustion Stop and Attention inspections; the fresh-state step test ran the exhausted rework denial at exit 4. Human comprehension judgment remains pending.
- VAL-073: FIXTURE READY, HUMAN PENDING - `TestP23HumanFixtures` reran exact Human text for eighteen leaves against the JSON facts, including purpose and qualifiers. Human acceptance of the presentation remains pending.
- VAL-079: PASS for the measured gate and report-inventory part. Attempt 1 recorded ordinary `cmd/fid` 214.615s, uncached `internal/workflow` 55.279s, and tagged `cmd/fid/integration` 229.859s. Attempt 2 recorded `cmd/fid` 202.293s, uncached `internal/workflow` 57.091s, and tagged `cmd/fid/integration` 273.836s. All are below the 600/600/900s budgets and below 80% headroom. The attempt-1 read-only inventory over task 1-17 reports found all three times in each report.

The attempt-2 executing commands for the rows below were the exact focused command in Tests run and the exact tagged gate in Exact quality-gate results. The owner task for each new leaf comes from `public-interface.md` section 4 and is pinned in `operations.json`. These rows report actual partial checks, not completed VAL-068 equality. The earlier producer tests and owners for the 19 declared Failure codes remain pinned in `failures.json`; its new literal partition marks 18 `added` and only `requirement-not-found` `retained`.

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

The nineteen `failures.json` rows each name one pre-existing producing `Test...` declaration and its Phase 23 owner task. They are: `review-not-independent` (4), `review-verdict-invalid` (5), `review-changes-required` (6), `review-rework-exhausted` (6), `finding-already-disposed` (5), `finding-not-found` (5), `review-not-found` (13), `concern-not-found` (8), `deviation-not-found` (8), `outcome-not-found` (10), `requirement-not-found` (12), `concern-already-disposed` (8), `claim-mode-invalid` (3), `protected-decision-required` (9), `human-outcome-required` (10), `outcome-not-permitted` (10), `evidence-type-not-accepted` (6), `source-digest-mismatch` (9), and `repository-mismatch` (4). The tagged gate executes their named producer tests; `assertP22ProcessFailure` calls the pinned `p23FailureContracts` fake envelope oracle for those tests.

## Tests run

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

## Deviations from spec

- Attempt 1 deviation: the full semantic fake/fid equality required by VAL-068 was unfinished. The fake then modeled top-level fields, replay and selected rejection codes. Attempt 2 added partial literal semantic checks for all 23 operations, but complete values and rejection classes remain unfinished. A same-task continuation must finish them without deriving fake expectations from fid output.
- Attempt 1 protected discrepancy: section 5/DEC-230's 18-versus-19 wording prevented an exact Failure-code claim. DEC-265 resolved it without a member or semantic change: 19 declared, 18 Phase-23-added, one retained `requirement-not-found` from `ab4ac8f` (`Refs: phase-20/task-1`). Attempt 2 pins and tests that partition. The original failed evidence remains above.
- The first `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-18-done.md project-docs/docs/phases/phase-23/tasks/task-18.md` attempt reported `FAIL report-lint found 1 failure(s)` (terminated `exit_code: 1`): the Recovery attempts field included prose after `0`. This report keeps the attempt count as the required integer and retains the explanation in the outcome and concerns sections.

## Concerns, technical debt, or notes for the Tech Lead

- Attempt 1 concern is resolved by DEC-265. The open task-18 concern is the remaining VAL-068 audit work: full independently authored result values, one shared generated-value normalization, cross-case ID and list consistency, and each applicable rejection's complete typed failure and unchanged Event count. No product divergence has been established by the partial checks.
- Historical task reports show `cmd/fid` package times at or above 80% of the 600s budget: task 2's later 490.783s, task 4's 496.639s, task 5's 588.268s, tasks 6-9's 535.185/499.695/538.754/513.222s, task 10's 564.337s attempt, and task 11's 531.533s final pass. This task's 214.615s and tagged 229.859s have ample headroom; the earlier measurements remain phase-level warnings.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
