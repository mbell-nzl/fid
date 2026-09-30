# Phase 23 - Task 18: Conformance assembly and wire-surface audit - Recovery report

## Workflow outcome

- Task outcome: BLOCKED
- Category: specification-contradiction
- Condition: `public-interface.md` section 5 enumerates 19 distinct added Failure codes, but its own prose and DEC-230 say the added set is exactly 18. The existing `p23FailureContracts` table also has 19. Certifying the exact set would choose between conflicting protected authorities.
- Recovery attempts: 0 protected-authority attempts. The in-scope audit and all five exact gates completed; no contract, code, or gate definition was changed to resolve the count.
- Next action: The contract owner must identify the exact accepted Failure-code members/count and record the protected disposition. The Phase Runner can then resume this same task to finish full semantic fake/fid parity and close its report.
- Human attention required: yes

## Summary of what was implemented

- Added an independently pinned thirty-action corpus. `TestGrantActionSetIsExactlyThirtyNames` compares it separately with the workflow `sortedUniqueCapabilities` map and the journal `const actions`, then compares the two implementations directly. A real activated profile and Envelope grant accept all thirty; five unknown or operation-only names reject without an Event.
- Extended the retained Condition fixture from six to ten types, with 54 accepted and 90 invalid/missing member cases. The per-type shape oracle checks the journal catalog, including the four-field `review-recorded` shape and each exclusive judgment discriminator. The Phase 22 document-snapshot assertion now retains its six historical types and adds the four accepted Phase 23 members.
- Added a shared operation inventory for the eighteen new leaves and five amended leaves. A static fake supplies the pinned accepted result-field set and replay expectation; built fid supplies the actual response. All nine new mutations and four amended mutations also reject an unknown field without changing the Event count. The corpus has fresh-state authority and stale-version rejection cases; retained producer tests cover the other named failure cases. The fake currently compares field sets and replay flags, not complete semantic result values, so VAL-068 remains open.
- Pinned the latest product table names and the built command tree. The test traverses real `--help` output, including commands without a short description, and inspects aliases; it compares the complete table and command inventories with the fixture. No cancelled product family was found. The task's `v1-convergence-inventory.md` section 1 already assigns this proof to task 18, so no correction was needed.
- Regenerated task 14's Human fixture manifest with the current binary. Its eighteen Human/JSON leaf pairs and all checkpoint inspections ran after generation; a fresh-state test ran all four recorded interactive steps at their expected exits. Human comprehension and presentation judgments remain for the Human.

## Files created or modified

- `cmd/fid/integration/p22_closed_values_test.go`, `cmd/fid/integration/p22_conformance_test.go`, `cmd/fid/integration/testdata/p22_conformance/wire-inventory.json`.
- `cmd/fid/integration/p23_conformance_test.go`, `p23_grant_action_audit_test.go`, `p23_cancelled_family_audit_test.go`, and `p23_human_fixtures_test.go`.
- `cmd/fid/integration/testdata/p23_conformance/operations.json`, `grant-actions.json`, `failures.json`, and `product-surface.json`.
- This report. The Runner-owned `project-docs/docs/phases/phase-23/run-ledger.md` was not edited.

## Documentation updated

- This report records the audit inventory, exact gates, fixture location, and protected count discrepancy. `v1-convergence-inventory.md` section 1 was already correct (I-10); no documentation change there was required.

## Validation Contract Coverage

- VAL-067: PARTIAL - `TestGrantActionSetIsExactlyThirtyNames` (`cmd/fid/integration/p23_grant_action_audit_test.go:23`) passes the independent thirty-member set equality, real profile/grant acceptance, and unknown-action rejection. The seven Phase 23 owner task numbers are pinned in `grant-actions.json`; historical at-or-after commit provenance still needs an explicit row-by-row audit in the continuation.
- VAL-068: BLOCKED for closure - `TestP23SharedLeafConformance` (`cmd/fid/integration/p23_conformance_test.go:67`) and `TestP23SharedAmendedLeafConformance` (line 153) pass result-field/replay and unknown-field comparisons for all 23 operations; they do not yet compare complete semantic results. The nineteen named failure-code rows in `failures.json` expose the section 5/DEC-230 count conflict. `TestP23FailureCorpusHasEveryAcceptedCode` (line 31) checks every named row against its existing producer test, without deciding which conflicting count governs.
- VAL-069: PASS - `TestP23CancelledFamiliesAbsentFromProduct` (`cmd/fid/integration/p23_cancelled_family_audit_test.go:22`) compares 59 latest-schema tables and 125 real command paths with a pinned full inventory and checks cancelled names and aliases. The exact tagged gate below executes it.
- VAL-070: FIXTURE READY, HUMAN PENDING - regenerated manifest `checkpoints-3135624237/manifest.json` has the mixed-Finding Review, Finding, and resume inspections; `TestP23HumanFixtures` reran each recorded inspection and leaf pair. Only a Human can judge whether resolved, addressed, and still-blocking are understood.
- VAL-071: FIXTURE READY, HUMAN PENDING - the same manifest records blocked action, Human Outcome, and retry steps. `TestP23HumanFixtureStepsRemainRunnable` (`p23_human_fixtures_test.go:165`) ran all three on a fresh protected-Decision checkpoint at exits 5, 0, 0. Human explanation judgment remains pending.
- VAL-072: FIXTURE READY, HUMAN PENDING - the manifest records the exhaustion Stop and Attention inspections; the fresh-state step test ran the exhausted rework denial at exit 4. Human comprehension judgment remains pending.
- VAL-073: FIXTURE READY, HUMAN PENDING - `TestP23HumanFixtures` reran exact Human text for eighteen leaves against the JSON facts, including purpose and qualifiers. Human acceptance of the presentation remains pending.
- VAL-079: PASS for the measured gate and report-inventory part - ordinary `cmd/fid` took 214.615s; an uncached race run reported `internal/workflow` 55.279s because the exact ordinary gate used cached workflow output; tagged `cmd/fid/integration` took 229.859s. These are below 600/600/900s and below 80% headroom. A read-only inventory over task 1-17 reports found all three package-time records in each report; no required gate timed out in this task.

The executing command for every operation and Failure-code row below is the final exact tagged gate, `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` (PASS, `cmd/fid/integration 229.859s`, exit 0). The focused shared audit command is recorded under Tests run. Owner provenance for new leaves is the accepted `public-interface.md` section 4 task column; `operations.json` pins it per row. The earlier producer test and owning task for each named Failure-code row are pinned in `failures.json` and checked by `TestP23FailureCorpusHasEveryAcceptedCode`.

| Operation | Owner task | Fake/fid comparison executed |
|---|---:|---|
| `review.record` | 4 | accepted keys/replay; unknown-field, fresh stale-version denial |
| `review.get` | 13 | accepted keys |
| `review.list` | 13 | accepted keys |
| `finding.address` | 5 | accepted keys/replay; unknown-field |
| `finding.list` | 5 | accepted keys |
| `concern.record` | 8 | accepted keys/replay; unknown-field |
| `concern.dispose` | 8 | accepted keys/replay; unknown-field |
| `concern.list` | 8 | accepted keys |
| `deviation.record` | 8 | accepted keys/replay; unknown-field |
| `deviation.list` | 8 | accepted keys |
| `decision.record-reference` | 9 | accepted keys/replay; unknown-field |
| `decision.list` | 9 | accepted keys |
| `human-outcome.record` | 10 | accepted keys/replay; unknown-field |
| `human-outcome.list` | 10 | accepted keys |
| `requirement.record-reference` | 12 | accepted keys/replay; unknown-field and fresh unauthorized Actor denial |
| `requirement.link` | 12 | accepted keys/replay; unknown-field |
| `requirement.get` | 13 | accepted keys |
| `requirement.list` | 12 | accepted keys |
| `envelope.create` | 2 | accepted keys/replay; unknown-field |
| `envelope.activate` | 2, 9 | accepted keys/replay; unknown-field |
| `claim.acquire` | 3 | accepted keys/replay; unknown-field |
| `evidence.link` | 6, 10 | accepted keys/replay; unknown-field |
| `evidence.evaluate` | 6, 13 | accepted `evaluations` result keys |

The nineteen `failures.json` rows each name one pre-existing producing `Test...` declaration and its Phase 23 owner task. They are: `review-not-independent` (4), `review-verdict-invalid` (5), `review-changes-required` (6), `review-rework-exhausted` (6), `finding-already-disposed` (5), `finding-not-found` (5), `review-not-found` (13), `concern-not-found` (8), `deviation-not-found` (8), `outcome-not-found` (10), `requirement-not-found` (12), `concern-already-disposed` (8), `claim-mode-invalid` (3), `protected-decision-required` (9), `human-outcome-required` (10), `outcome-not-permitted` (10), `evidence-type-not-accepted` (6), `source-digest-mismatch` (9), and `repository-mismatch` (4). The tagged gate executes their named producer tests; `assertP22ProcessFailure` calls the pinned `p23FailureContracts` fake envelope oracle for those tests.

## Tests run

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

Working directory `/home/mark/projects/fid`; commands executed in the task-spec order on the final test tree. The ordinary gate's workflow package output was cached, so the separate uncached timing command above supplies its current measured time.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)			60.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	229.859s` (terminated `exit_code: 0`)

## Deviations from spec

- The full semantic fake/fid equality required by VAL-068 is unfinished. The current fake models exact top-level fields, replay, and selected rejection codes; earlier producer tests separately exercise domain facts. A same-task continuation must extend the shared cases to compare accepted semantic values and the required rejection classes across every applicable operation without deriving fake expectations from fid output.
- The section 5/DEC-230 18-versus-19 conflict prevents a final exact Failure-code corpus claim. The first failed audit attempts and their later green results remain recorded above; no failure code or protected threshold was changed.

## Concerns, technical debt, or notes for the Tech Lead

- Protected decision needed: section 5 enumerates nineteen named codes after deleting `finding-unresolved` and `decision-not-found` from the draft, while the same section and DEC-230 say eighteen. The current implementation and `p23FailureContracts` have nineteen. Identify the exact intended member set before task closure; this report does not assign a different task to repair it.
- Historical task reports show `cmd/fid` package times at or above 80% of the 600s budget: task 2's later 490.783s, task 4's 496.639s, task 5's 588.268s, tasks 6-9's 535.185/499.695/538.754/513.222s, task 10's 564.337s attempt, and task 11's 531.533s final pass. This task's 214.615s and tagged 229.859s have ample headroom; the earlier measurements remain phase-level warnings.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `333514946aeab8328fa618f0f619fdce84694f9a`

## Commit subject

- `test(conformance): audit Phase 23 wire surface`
