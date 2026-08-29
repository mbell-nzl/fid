# Phase 21 Group 1 Independent Checkpoint Review

> Checkpoint attempt: 4, authorized by DEC-168
> Reviewer: Codex GPT-5, fresh independent QA Lead
> Independence: no authorship of tasks 1-12 or QA Review Windows 1-4
> Session date: 2026-08-29 (Pacific/Auckland)
> Repository: `/home/mark/projects/fid`
> Phase HEAD under review: `6825d59e282c3795213fada9f6e9b082f4db3354`
> Contract/spec commit: `ead4258d787bc942bc8394c061afc1e58e903627`
> Review status: REPAIR

This is one unsplit checkpoint session. Prior scratch directories and prior
diagnostic/materialized outputs are excluded. All FEAS evidence below is
produced from a fresh `mktemp -d` allocation created in this session.

## Scope And Authorities

- Governing assertions: `VAL-042` and `VAL-043`.
- Governing checkpoint: `qa/group-1-checkpoint-contract.md`.
- Evidence procedures: QA plan sections 5.1, 7.5, 10, and 11.
- Temporary V1 authority: `v1-convergence-inventory.md` sections 5 and 7.
- Binding recovery authority: DEC-168 attempt 4; DEC-167 governs the tagged
  closure-journey evidence classification.
- Group 2 implementation and capability advertisement are out of scope.

The exact checkpoint/slice-contract commit is `ead4258d787bc942bc8394c061afc1e58e903627`
(`docs(tl): resolve phase 21 task 2 blockers`), which is the current last commit
for both `group-1-checkpoint-contract.md` and `contracts.md`. Related governing
file provenance is: public interface and validation contract
`95b8bc854988e12a5d6a19159b24f5d8ffe8637f`, temporary V1 inventory
`05d3dc6e1aa388d50b7fa842c2de9537b5208f2b`, and Ready QA plan
`e8a636e5f19d8f1d2fcf1e9157d6eb42ac6ebb4c`.

## Session Preconditions

| Check | Result | Evidence |
|---|---|---|
| Fresh independent reviewer | PASS | Reviewer identity and independence declaration above. |
| Authorized attempt 4 wrapper | PASS | `run-ledger.md` contains the DEC-168 `RUN-RESUME` and `TASK-RECOVERY task-13`, trigger `task-13:group-1-checkpoint:qa-lead-feasibility-session-interrupted`, attempt 4, at phase HEAD `6825d59e282c3795213fada9f6e9b082f4db3354`. |
| Prior scratch excluded | PASS | `/tmp/tmp.YC4EUO3gMN`, `/tmp/tmp.zQAE0hBrd2`, and prior diagnostic outputs were not read, copied, reused, or counted. |
| Sole authored repository artifact | PASS SO FAR | Before creation, `git status --short` showed only runner-owned `M project-docs/docs/phases/phase-21/run-ledger.md`; this review is the only QA-authored path. |

## VAL-042 Artifact And Commit Audit

Outcome: PASS.

### Task Reports And Report-Lint

| Task | Final report commit | Reported implementation commit | Fresh two-argument report-lint | Pre-dispatch ledger evidence |
|---:|---|---|---|---|
| 1 | `a169b4114dfbd6f3d19625dd577cf6867bc707cc` | `a094e690815c9940ca926bcad68aedacea663c43` | PASS - `PASS report-lint completed with no failures` (terminated `exit_code: 0`) | `run-ledger.md:523`, lint PASS at `:526`, before task-3 recovery at `:530` |
| 2 | `e6946c891c1ed554dc2ed087f359643cb7ee4f06` | `c0fd2cd8343826fdcd1a597f8f2f3a671e0e15dc` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:329`, lint PASS at `:332`, before task-3 recovery at `:337` |
| 3 | `68cc7ab9c028d9588e2455e18ad56fa5a8474578` | `401f609d4d7e24c0ab5ca45d27a16424000ded9f` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:569`, lint PASS at `:572`, before review-1b at `:576` |
| 4 | `56062540e901553ad5ba1ca7ce208774c05f979f` | `ad9098fac66f2b530a566d87df087fd72fbb584d` plus review repair `cd83fe9f12aee24bfc1fb3e73145bdaffaf731ed` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:732`, lint PASS at `:735`, before review-2b at `:739` |
| 5 | `b794ffe066b0ef16d455948175dd2d117cb4cb3a` | `0c4c04ea0c7d0786e2e6cf88a7573d8a3105aef5` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:631`, lint PASS at `:634`, before task 6 at `:638` |
| 6 | `10fd45f8a58d3c304c56eb0b1349cae5aaeb73b7` | `56dad8fa0dcb9b79f7df2a31b356f1b2eb8a694e` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:638`, lint PASS at `:641`, before review-2 at `:645` |
| 7 | `541869b92471f86df93a574ceeeb75fc878df1d5` | `817ad7c50469c51dc12f88a4bc081ecd7a384d7e` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:748`, lint PASS at `:751`, before task-8 recovery at `:755` |
| 8 | `9b6dccf4ac291e20a50f7f406b285c7ee803aa60` | `9f1c6d03aef12723aad6387b3096c7485655a2ed` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:1715`, lint PASS at `:1718`, before task-10 recovery at `:1722` |
| 9 | `fa2a5ada97181b1c3c04e09441e0ce4dd6b4c396` | `942abad7e9b6a0b0fb9fd009d074aeab79042457` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:1591`, lint PASS at `:1594`, before task-8 recovery at `:1598` |
| 10 | `1c4d222b998374c97134e06420fc1a7dcb0cac79` | `59b079a5a271cc65f7964a471b6383b9c2dd6677` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:2174`, lint PASS at `:2177`, before task-11 recovery at `:2181` |
| 11 | `1f24ea26e0d0ad0b32c33d2dec37ec811b631f27` | `8c763e9cddee1d4ff0df3379e08b286ef74d5300` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:2259`, lint PASS at `:2262`, before task-12 recovery at `:2266` |
| 12 | `20e6760726edc718fe6d0a79df046835202c9b85` | `0ced511b2e6a7e1f252260024176705106e9b793` | PASS - same token (terminated `exit_code: 0`) | `run-ledger.md:2305`, lint PASS at `:2308`, before review-4b at `:2312` |

Each fresh report-lint result also emitted `OUTCOME COMPLETE`, verified every
exact task-spec command as recorded PASS, resolved and subject-matched the
reported implementation commit, and verified the report committed and clean.

### Review Windows And Finding Closure

| Window | Initial artifact | Initial outcome/findings | Fresh re-review | Final outcome | Resolution evidence |
|---:|---|---|---|---|---|
| 1 | `review-1.md` at `17d4e22d15b6fdc0024a8fd8ff539ae6f7ceeb5b` | REPAIR / 2 Must-fix (tasks 1, 3) | `review-1b.md` at `44c0b21f990dd31281d0069886e6dc61c17f8343` | PASS | Both findings RESOLVED with fresh DB/package gates and report-lint; ledger ordered task-1 then task-3 repair before review-1b. |
| 2 | `review-2.md` at `8ca2f8aacfbcd75809a0e661f037063fd8a2fd32` | REPAIR / 1 Must-fix (task 4) | `review-2b.md` at `df1b46a5e5e7fbf1656c2824d886b84d2b8e0e82` | PASS | Finding RESOLVED with lifecycle/abort/augmentation evidence, exact package gate, and report-lint. |
| 3 | `review-3.md` at `303c75e8de363b5f30b614b7eea878a137fb7cac` | REPAIR / 1 Must-fix (task 9) | `review-3b.md` at `cba74981c16c6855139708beccc8dff9f159de7d` | PASS | Finding RESOLVED with actor/epic/model positives, role/invalid/`--by` negatives, journal digest stability, exact package gate, and report-lint. |
| 4 | `review-4.md` at `11396b45a79244a00c7379800d7a647fa32c0cb0` | REPAIR / 3 Must-fix (tasks 10, 11, 12) | `review-4b.md` at `6caaa6f7dc7cca0c5e9e4306f81ae53228f25d25` | PASS | All three findings RESOLVED with fresh Claim, Run, transition, tagged-integration, and report-lint evidence. |

## Independent Group 1 Contract Review

Result: REPAIR because two bounded Group 1 defects remain; see Findings And
Verdict. The accepted behavior not implicated by those defects passed as
follows.

| Contract item | Result | Fresh evidence |
|---|---|---|
| Event v2, migration, rebuild, and idempotency authority | PARTIAL | Race-enabled `TestMigration0013OperationalKernelInventoryAndRoundTrip` and the three `TestCommitIdempotent*` replay/conflict/concurrency tests passed (`ok .../internal/db 1.735s`, `ok .../internal/journal 2.620s`, terminated `exit_code: 0`). The clean real-process Activation replay defect in MF-1 prevents an overall PASS. |
| Capability matrix / no Group 2 advertisement | PASS | `TestCapabilitiesJSONGoldenAndProjectionRows` and `TestCapabilitiesRejectsOverridesAndKeepsDeferredRowsUnavailable` passed; all seven deferred rows stayed false on their named surfaces (terminated `exit_code: 0`). The hostile fixture's pre/post capability JSON was byte-identical. |
| Finalized Activation trust | PASS | `TestFinalizedActivationIsTheOnlyTrustedActorSource` passed and the hostile journey reached Registration only after `ACT-72da41` finalized (terminated `exit_code: 0`). |
| Registration ordering and no policy override | PASS | `TestGateRegisterRequiresFinalizedActivatedDefinition` and `TestGateRegisterRejectsRuntimePolicyOverride` passed; hostile Registration `REG-a612f3` resolved the activated definition and immutable binding (terminated `exit_code: 0`). |
| Envelope, Obligation, budget, Grant, Claim, Run, transition | PASS | Focused race tests for the exact Envelope, Claim, Run, Grant, and execution-transition matrices passed, including every Grant mismatch and ordered transition reason. The hostile real-process journey below independently committed each qualifying operation. |
| Seven-name bridge and explicit-role precedence | PASS | `TestKnownRolesUsesV1CompatibilityVocabularyForRoleFreeV2Bootstrap` and `TestKnownRolesKeepsExplicitRoleFactsAuthoritative` passed (terminated `exit_code: 0`); the generic source and dependency checks below show no new-path reachability. |
| Public `fid policy` absent | PASS | `TestPublicPolicyCommandIsAbsent` passed (terminated `exit_code: 0`). |
| Judgment mutation fail-closed | PARTIAL | AC verify, Concern raise/address, and maintenance vacuum all returned exact `capability-unavailable` envelopes and appended no Event, but real-process exit 1 violates the accepted exit 5 contract (MF-2). |
| No Group 2 pull-forward | PASS | Deferred capability tests keep all seven Group 2 rows false; `fid gate --help` exposes only `register` and `resolve`; no Group 2 behavior was used by this checkpoint. |

## VAL-043 Isolation And Temporary Inventory

Outcome: PASS. The exact temporary inventory remains present only in its
contracted checkpoint disposition, every fresh inventory match was classified,
and the seeded hostile journey proves that no new generic path consults it.

### Section 5 Temporary Inventory

| Inventory class | Contracted checkpoint state | Fresh classification |
|---|---|---|
| Tables and role/judgment requirement rows | May remain only as exactly enumerated | PASS - all fresh SQL/table and requirement-row references map to inventory section 4; no unlisted table family or rule was found. |
| Replacement-bound commands/writers and Events | May remain only as exactly enumerated | PASS - `dev`, legacy completion/claim/run compatibility, and their V1 Events map to sections 2 and 3; accepted Event v2 writers remain distinct. |
| Replacement-bound source/readers | May remain only as exactly enumerated | PASS - frozen readers, maintenance diagnostics, historical statistics, and structural workflow readers map to sections 1.1 and 4. |
| Internal V1 policy implementation | Reachable only by enumerated dev/completion compatibility callers | PASS - production caller inventory is limited to the enumerated compatibility paths; the generic dependency and source checks below are clean. |
| Seven-name nonpersistent bridge | Exact role-free bootstrap only; explicit role Effects win; generic paths unreachable | PASS - bootstrap metadata and the retained standalone compatibility bridge are classified; focused policy and generic-authority tests below exercise precedence and non-reachability. |
| Final V1 absence | Reserved for task 20; not required here | NOT REQUIRED |
| Final-schema parity | Reserved for task 21; not required here | NOT REQUIRED |

### Section 7 Fresh Classified Inventories

| Recipe | Expected checkpoint status | Observed termination | Classification |
|---|---:|---|---|
| Role-policy references | `exit_code: 0` | `cmd/fid/dev.go:...Role: role` through `internal/workflow/transition.go:...Role:` (terminated `exit_code: 0`) | Every match classified to section 1 or 1.1: temporary dev/completion/claim/run policy, role-free bootstrap/lock metadata, or retained standalone attribution. |
| `--by` references | `exit_code: 0` | `cmd/fid/claim.go:...GetString("by")` through `internal/runs/record.go:...RecordedBy` (terminated `exit_code: 0`) | Every match classified to temporary replacement-bound command input, legacy issue-mutation rejection, or internal compatibility diagnostics. |
| V1 table references | `exit_code: 0` | `cmd/fid/dev.go:...task_completions` through `internal/workflow/store.go:...workflow_transition_requires` (terminated `exit_code: 0`) | Every match classified to section 4 as a temporary writer, frozen reader, diagnostic/statistical reader, or structural workflow reader. |
| Event families | `exit_code: 0` | `cmd/fid/dev.go:...task_completion_recorded` through `internal/workflow/transition.go:...issue_transitioned` (terminated `exit_code: 0`) | Every match classified to section 3 as temporary/frozen V1 or accepted v2/current event output. |

UNCLASSIFIED: 0.

### Generic-Path Dependency And Adversarial Evidence

| Evidence | Result | Fresh command/test evidence |
|---|---|---|
| Generic dependency graph | PASS | `go list -deps ./internal/workflow ./internal/policy ./internal/claims ./internal/evidence` ended with `internal/evidence` (terminated `exit_code: 0`); no `internal/dev` or `internal/runs` package appeared. |
| Generic source inventory absence | PASS | Dedicated generic implementation file set searched for `roles`, all six temporary evidence/judgment tables, and `internal/dev`: `exit 1 (no output)` (terminated `exit_code: 1`, expected absence). A separate broad diagnostic including the hybrid `internal/issues/transition.go` found only its retained standalone error text `missing roles`; source inspection and focused execution tests prove `execution.Bound` returns before the legacy rule evaluator. |
| Activation isolation | PASS | In fresh `/tmp/tmp.AtLMawnhdn`, `activation.prepare`, both component stages, commit, and finalize succeeded after all eight hostile rows were seeded; final token included `status":"finalized"` (each terminated `exit_code: 0`). |
| Registration isolation | PASS | `gate.register` produced `REG-a612f3`, version 1, bound to `ACT-72da41` and `gate-unit`, after the hostile seed (terminated `exit_code: 0`). |
| Envelope isolation | PASS | `envelope.create` produced draft `ENV-caf49d`, and `envelope.activate` produced active revision 1 with the exact Obligation and null budget counters (each terminated `exit_code: 0`). |
| Grant isolation | PASS | The hostile Claim, Run, and transition were authorized only by the finalized profile and exact Envelope Grant. `TestEnvelopeAuthorityRequiresEveryBindingAndRejectsRoleAttribution` passed every single-dimension mismatch and all 13 seeded-role protected-operation negatives (terminated `exit_code: 0`). |
| Claim isolation | PASS | V2 Claim `CLM-82fe8f` committed despite the hostile V1 `claims` row for the same Issue (terminated `exit_code: 0`). |
| Run isolation | PASS | Run `RUN-a50561` committed despite hostile `agent_runs` and `gate_runs`; `TestGenericRunUsesNoHistoricalAgentRunsTable` also passed (terminated `exit_code: 0`). |
| Nonterminal transition isolation | PASS | `QA-2cf0e4 todo -> in_progress` committed at revision 2 despite the hostile V1 Claim, open Concern, hostile AC verification, and an added unsatisfied `all_acs_verified` human requirement on that exact edge (terminated `exit_code: 0`). |
| Capability isolation | PASS | Pre-seed and post-seed `fid capabilities --json` files were byte-identical: `PASS capabilities byte-identical after hostile seed` (terminated `exit_code: 0`). |
| Idempotency isolation | PASS | Repeating the exact seeded `claim.acquire` returned original Event `EVT-20260828-000012`, original Claim `CLM-82fe8f`, and only changed `replayed` to `true` (terminated `exit_code: 0`; no new Event). |
| AC verification fail-closed | PASS for VAL-043 isolation | Seeded `ac_verifications` did not enable mutation: response code `capability-unavailable`, capability `reviews.coverage.v1`, no Event. The real process terminated exit 1 rather than accepted exit 5; this separate Task 9 defect is MF-2. |
| Concern mutation fail-closed | PASS for VAL-043 isolation | Seeded open `concerns` and completion rows did not enable raise or address: both returned `capability-unavailable`, capability `human.outcomes.v1`, and no Event. Both real processes terminated exit 1 rather than accepted exit 5; see MF-2. |

The seed transaction printed exactly one row for each of `roles`, V1 `claims`,
`task_completions`, `ac_verifications`, `agent_runs`, `gate_runs`, `concerns`,
and the hostile `workflow_transition_requires` row, then terminated
`exit_code: 0`. The hostile journal remained at 14 Events after all four
fail-closed judgment probes, confirming their zero-Event effect.

## Fresh FEAS-01..FEAS-19 Evidence

Fresh scratch path: `/tmp/tmp.mbG09jOf8J`. This path was allocated by
`mktemp -d` in this checkpoint session and is distinct from every excluded
prior path. A new `go.mod`, six-mode `child/main.go`, and `child-bin` were
materialized there. `go vet ./...` and
`go build -buildvcs=false -o /tmp/tmp.mbG09jOf8J/child-bin ./child` were both
silent and terminated `exit_code: 0`.

| Probe | Result | Output token / observable | Observed termination | Budget |
|---|---|---|---|---|
| FEAS-01 | PASS | `argv=["arg-1","arg-2"]`, `out-1`, `err-1`, `out-2` | terminated `exit_code: 0` | <0.01s |
| FEAS-02 | PASS | stdout file: `argv=[]`, `out-1`, `out-2`; stderr file: only `err-1` | child and both inspections terminated `exit_code: 0` | <0.01s |
| FEAS-03 | PASS | `00000000: 41ff fe1b 5b33 316d 0d0a` | writer and `xxd` terminated `exit_code: 0` | <0.01s |
| FEAS-04 | PASS | `exit 124 (no output)` after TERM timeout | terminated `exit_code: 124` | 4.898s |
| FEAS-05 | PASS | `exit 124 (no output)` after KILL timeout; `pgrep -x child-bin` then found no survivor | timeout terminated `exit_code: 124`; process check terminated `exit_code: 1` | 4.896s |
| FEAS-06 | PASS | `does-not-exist: No such file or directory` | terminated `exit_code: 127` | <0.01s |
| FEAS-07 | PASS | `8388608` bytes stdout, `8388608` bytes stderr, `16777216 total` | child and `wc -c` terminated `exit_code: 0` | <0.01s |
| FEAS-08 | PASS | exactly `HOME=/tmp`, `PATH=/usr/bin`, `PROBE_MODE=env` | terminated `exit_code: 0` | <0.01s |
| FEAS-09 | PASS | `1:token=CANARY-g1-a4-9f27` | writer and positive `rg` terminated `exit_code: 0` | <0.01s |
| FEAS-10 | PASS | base `/tmp/tmp.mbG09jOf8J/cwd`; resolved `/tmp/tmp.mbG09jOf8J/outside/target.txt`; `ESCAPES` | terminated `exit_code: 0` | <0.01s |
| FEAS-11 | PASS | before/copy digest `1e6537...a30f68`; mutated digest `40eab3...b06b0` | hash/copy/mutate comparison terminated `exit_code: 0` | <0.01s |
| FEAS-12 | PASS | `ATOMIC_STORE digest=309150...cc3bb tmp-empty object-present` | terminated `exit_code: 0` | <0.01s |
| FEAS-13 | PASS | actual nested exits retained for FEAS-04..07: `124`, `124`, `127`, `0` | copied from each command result, not wrapper completion | derived |
| FEAS-14 | PASS | tree OIDs `893364...147fb` then `c0615a...f3735` | every Git procedure command terminated `exit_code: 0` | <0.01s |
| FEAS-15 | PASS | `MISSING_OBSERVABLE: PASS`; `RESTORED_OBSERVABLE: PASS` | both observations terminated `exit_code: 0` | <0.01s |
| FEAS-16 | PASS | both archives `a68d8f...a5b55`; `ZIP_DETERMINISM: PASS` after 1.1s separation | writer and comparison terminated `exit_code: 0` | 1.132s |
| FEAS-17 | PASS after preserved retry | valid archive: `No errors detected...`, unchanged digest `a68d8f...a5b55`; byte 7 corruption: `mismatch between local and central GPF bit 11` | valid `unzip -t` terminated `exit_code: 0`; corrupt run terminated `exit_code: 1` | <0.01s per attempt |
| FEAS-18 | PASS | all 16 distinct section 5.1 nonzero shapes emitted no output: 11 `rg`/`grep` shapes at exit 1 and five `ls` shapes at exit 2 | each observed status exactly matched its declaration | <0.01s each |
| FEAS-19 | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 5.483s` | terminated `exit_code: 0` | 5.827s, independent |

Aggregate FEAS-01..FEAS-18 observed command wall time was approximately
10.94s, far below 900s; no individual probe approached 120s. FEAS-19 took
5.827s, independently below 300s.

FEAS-17 deviation record: the first corruption selection overwrote byte 50
with zero. It produced `bad CRC 716e0fc0 (should be 9f606eec)` and terminated
`exit_code: 2`, so it did not satisfy the contract and is not the closing
evidence. The same one-byte-corruption procedure was repeated from the fresh
uncorrupted archive at byte 7; GNU `unzip` produced its warning classification
and the required observed `exit_code: 1`. No command definition, accepted
status, or archive-construction rule changed.

## DEC-167 Tagged Integration Note

`TestV01ClosureJourneyAndRebuildParity` had never been green in Phase 21 before
the first green task-12/review evidence. Fresh FEAS-19 is green in this
checkpoint (`ok github.com/mbell-nzl/fid/cmd/fid/integration 5.483s`, terminated
`exit_code: 0`) and is recorded as new checkpoint evidence. It is not
continuity evidence and is not classified as regression-free carry-forward.

## Required Artifact Checks

| Check | Result | Evidence |
|---|---|---|
| Review completeness | PASS | Exact identity/HEAD/contract provenance present; `rg -c '^\| FEAS-[0-9][0-9] '` returned `19`; separate VAL-042/043 outcomes, DEC-167 note, findings, and literal verdict are present (terminated `exit_code: 0`). |
| `git diff --check` for this artifact | PASS | `git diff --cached --check` produced no output (terminated `exit_code: 0`). |
| Worktree scope / only review staged | PASS | `git diff --cached --name-only` listed only this review (terminated `exit_code: 0`); `git status --short` separately showed the Phase Runner's `run-ledger.md` modification unstaged and untouched. |

## Findings And Verdict

### MF-1 - Activation prepare cannot replay its own committed result

- Classification: Must-fix; bounded `REPAIR`.
- Owning task: task 4, Activation lifecycle implementation.
- Affected assertion: VAL-005 as exercised by the accepted Activation mutation,
  plus the checkpoint Independent Contract Review; VAL-042 and VAL-043 are not
  failed by this finding.
- Reproduction: in fresh clean `/tmp/tmp.eL7w0FYFsM`, initialize a project and
  run byte-identical
  `fid --json activation prepare --input prepare.json --idempotency-key clean-activation-replay`
  twice. The first returned `EVT-20260828-000002`, `status:"prepared"`, and
  `replayed:false` (terminated `exit_code: 0`). The second emitted no stdout or
  stderr and terminated `exit_code: 1`; the journal remained exactly 2 lines.
- Contract mismatch: `public-interface.md` section 6 requires every keyed
  mutation's same-key/same-digest replay to return the original Event and
  domain fields, change only `replayed:true`, and append no Event. The generic
  `CommitIdempotent` replay/rebuild/conflict/concurrency tests are green, so the
  defect is in the Activation caller path.
- Exact repair boundary: within task 4's `internal/workflow/` and `cmd/fid/`
  scope, make `activation.prepare` recognize and return an existing exact
  idempotency result before a fresh-preparation lifecycle conflict can reject
  it. Preserve new-key lifecycle conflicts, digest conflicts, Event schema,
  result shape, and every protected contract/command definition. Add a
  real-process or equivalent command-level regression that proves the two-call
  result and zero Event delta.
- Required evidence: red/green reproduction with both observed exits and Event
  count; task 4's exact five commands with captured termination; committed task
  report and fresh two-argument report-lint. A fresh independent QA Lead then
  repeats the entire Group 1 checkpoint, including all FEAS probes, because the
  checkpoint contract requires a complete unsplit re-review.

### MF-2 - Fail-closed leaves use the wrong process exit class

- Classification: Must-fix; bounded `REPAIR`.
- Owning task: task 9, retained inventory attribution and judgment freeze.
- Affected assertion: VAL-030 and the checkpoint Independent Contract Review;
  VAL-043's V1-isolation observable still passes because every operation
  remained unavailable and mutation-free.
- Reproduction: against the fully seeded hostile workspace, `issue ac verify`,
  `concern raise`, `concern address`, and `maintenance vacuum` each emitted the
  exact `capability-unavailable` JSON for respectively
  `reviews.coverage.v1`, `human.outcomes.v1`, `human.outcomes.v1`, and
  `telemetry.trajectory.v1`, but each terminated `exit_code: 1`. The journal
  stayed at 14 lines. `cmd/fid/main.go` currently selects `os.Exit(1)` for every
  returned error.
- Contract mismatch: `public-interface.md` section 5 and the section 7 failure
  table require `capability-unavailable` to terminate exit 5; QA plan QA-29
  requires the same exit-5 observable.
- Exact repair boundary: within task 9's accepted `cmd/fid/` and
  `internal/cli/` scope, select the accepted process exit for
  `capability-unavailable` without changing any JSON error field, command
  grammar, capability availability, writer absence, or other protected exit
  mapping. Cover all four fail-closed leaves with real built-binary tests.
- Required evidence: each of the four commands returning the exact current
  envelope and observed `exit_code: 5`, unchanged journal/projection state,
  task 9's exact five commands, committed report, and fresh two-argument
  report-lint. The fresh independent reviewer reruns the complete checkpoint.

Blocker candidates: none. Both findings are decision-free and bounded inside
existing Group 1 tasks; no protected contract change or external capability is
required.

- VAL-042: PASS.
- VAL-043: PASS.
- Fresh feasibility FEAS-01..FEAS-19: PASS.
- Independent Group 1 Contract Review: REPAIR (MF-1, MF-2).

CHECKPOINT VERDICT: REPAIR

---

# Phase 21 Group 1 Checkpoint - Second Independent Review

Date: 2026-08-29

Reviewer: fresh independent QA Lead

Group 1 HEAD reviewed: `5aeb19b1e51a78386e0a51aac00fb20f9ec44172`

Fresh feasibility scratch: `/tmp/tmp.tauoSIKOsG`

Review scope: complete checkpoint rerun after the DEC-169 task 4 and task 9
repairs and before task 13 dispatch

The first review above, including its `CHECKPOINT VERDICT: REPAIR`, is retained
verbatim. Before this append it occupied 276 lines and had SHA-256
`2c476e24b31185c730f4cd58720c4bf2b7ae5501c07c5b594702166f55f499c1`.
No prior checkpoint scratch or diagnostic path was used as evidence in this
review. Task 13 was not dispatched or inspected as implementation work.

## Authority And Provenance

The checkpoint was rerun against these exact governing revisions:

- Group 1 checkpoint contract:
  `ead4258d787bc942bc8394c061afc1e58e903627`.
- QA plan, including sections 7.5 and 11:
  `e8a636e5f19d8f1d2fcf1e9157d6eb42ac6ebb4c`.
- Validation contract:
  `95b8bc854988e12a5d6a19159b24f5d8ffe8637f`.
- Public-interface contract:
  `95b8bc854988e12a5d6a19159b24f5d8ffe8637f`.
- V1 convergence inventory:
  `05d3dc6e1aa388d50b7fa842c2de9537b5208f2b`.
- Phase decisions, including DEC-167 item 8, DEC-168, and DEC-169:
  `8daf8dd88d852e7ac81ca3e1d20a244d2ad56cb1`.

The task-spec, implementation, and committed completion-report provenance is:

| Task | Spec revision | Implementation | Report revision |
|---|---|---|---|
| 1 | `95b8bc854988e12a5d6a19159b24f5d8ffe8637f` | `a094e690815c9940ca926bcad68aedacea663c43` | `a169b4114dfbd6f3d19625dd577cf6867bc707cc` |
| 2 | `05d3dc6e1aa388d50b7fa842c2de9537b5208f2b` | `c0fd2cd8343826fdcd1a597f8f2f3a671e0e15dc` | `e6946c891c1ed554dc2ed087f359643cb7ee4f06` |
| 3 | `b367a3be2a34abd91bb094069afae8cd4c1e3bb8` | `401f609d4d7e24c0ab5ca45d27a16424000ded9f` | `68cc7ab9c028d9588e2455e18ad56fa5a8474578` |
| 4 | `8daf8dd88d852e7ac81ca3e1d20a244d2ad56cb1` | `0444ccc458f3ac5a720724825f7f81589cb41e09` | `ae8ddb2b8e4faf03fb95d28586686f3f81a4149f` |
| 5 | `b367a3be2a34abd91bb094069afae8cd4c1e3bb8` | `0c4c04ea0c7d0786e2e6cf88a7573d8a3105aef5` | `b794ffe066b0ef16d455948175dd2d117cb4cb3a` |
| 6 | `b367a3be2a34abd91bb094069afae8cd4c1e3bb8` | `56dad8fa0dcb9b79f7df2a31b356f1b2eb8a694e` | `10fd45f8a58d3c304c56eb0b1349cae5aaeb73b7` |
| 7 | `b367a3be2a34abd91bb094069afae8cd4c1e3bb8` | `817ad7c50469c51dc12f88a4bc081ecd7a384d7e` | `541869b92471f86df93a574ceeeb75fc878df1d5` |
| 8 | `b367a3be2a34abd91bb094069afae8cd4c1e3bb8` | `9f1c6d03aef12723aad6387b3096c7485655a2ed` | `9b6dccf4ac291e20a50f7f406b285c7ee803aa60` |
| 9 | `8daf8dd88d852e7ac81ca3e1d20a244d2ad56cb1` | `1d89aa409bc71443c712948c54c2c10faba1c809` | `5aeb19b1e51a78386e0a51aac00fb20f9ec44172` |
| 10 | `b367a3be2a34abd91bb094069afae8cd4c1e3bb8` | `59b079a5a271cc65f7964a471b6383b9c2dd6677` | `1c4d222b998374c97134e06420fc1a7dcb0cac79` |
| 11 | `b367a3be2a34abd91bb094069afae8cd4c1e3bb8` | `8c763e9cddee1d4ff0df3379e08b286ef74d5300` | `1f24ea26e0d0ad0b32c33d2dec37ec811b631f27` |
| 12 | `aa2c263190801659383d2b1e900d0cbcb08d757b` | `0ced511b2e6a7e1f252260024176705106e9b793` | `20e6760726edc718fe6d0a79df046835202c9b85` |

The worktree started at the required HEAD with only Phase Runner-owned changes
to `run-ledger.md` and `triage-log.md`. This review did not stage or modify
either file.

## Checkpoint Preconditions And VAL-042

`VAL-042: PASS`.

- All task 1 through task 12 specs and committed completion reports exist.
  Each fresh two-argument `report-lint.sh <report> <spec>` run terminated
  `exit_code: 0` with `PASS report-lint completed with no failures` and each
  report states `OUTCOME COMPLETE`.
- The run-ledger audit found each task's completion only after its committed
  implementation, committed report, and Phase Runner two-argument report-lint
  evidence. The task 4 and task 9 repair completions are represented by their
  later report commits shown above.
- Window 1-3 review `17d4e22d15b6fdc0024a8fd8ff539ae6f7ceeb5b`
  is closed by independent re-review
  `44c0b21f990dd31281d0069886e6dc61c17f8343` with `PASS`.
- Window 4-6 review `8ca2f8aacfbcd75809a0e661f037063fd8a2fd32`
  is closed by independent re-review
  `df1b46a5e5e7fbf1656c2824d886b84d2b8e0e82` with `PASS`.
- Window 7-9 review `303c75e8de363b5f30b614b7eea878a137fb7cac`
  is closed by independent re-review
  `cba74981c16c6855139708beccc8dff9f159de7d` with `PASS`.
- Window 10-12 review `11396b45a79244a00c7379800d7a647fa32c0cb0`
  is closed by independent re-review
  `6caaa6f7dc7cca0c5e9e4306f81ae53228f25d25` with `PASS`.

No task-window finding remains open.

## Independent Contract Review

Result: `PASS`.

- Fresh, cache-bypassed package execution
  `timeout 120s go test -race -count=1 ./internal/journal ./internal/policy
  ./internal/workflow ./internal/claims ./internal/evidence ./internal/issues`
  terminated `exit_code: 0`. The last package token was
  `ok github.com/mbell-nzl/fid/internal/issues 76.728s`.
- Fresh focused command tests for public `fid policy` absence, exact deferred
  capability rows, complete Envelope authority binding, finalized Activation
  requirements for Gate registration, generic Run V1 independence, and
  execution-bound transition V1 independence terminated `exit_code: 0` with
  `ok github.com/mbell-nzl/fid/cmd/fid 10.279s`.
- Event-v2 validation, journal-first idempotency, capability publication,
  Activation, Gate registration, Envelope values, authority evaluation,
  generic Claim, generic Run, and execution-bound transition remain within
  the accepted Group 1 contract. No Group 2 executor, live-frame, receipt,
  Artifact, evidence-link, or export implementation was pulled forward.
- The production absence inventory found no `internal/executor`,
  `internal/artifacts`, or `internal/export` package (`ls` terminated
  `exit_code: 2`, no output). The only Group 2 term found in Group 1 source was
  the accepted `gate.run` capability vocabulary in the closed Activation
  authority profile, not an executor or streaming implementation.

### DEC-169 Task 4 Repair

The former MF-1 is resolved within task 4's accepted boundary.

- Fresh
  `TestActivationReplaysAcrossLifecycleWithoutPreemptingLifecycleChecks`
  execution terminated `exit_code: 0` with
  `ok github.com/mbell-nzl/fid/cmd/fid 3.242s`.
- Exact same-key/same-digest replay is selected before lifecycle rejection for
  prepare, authority-profile component, gate-definition component, commit,
  abort, and finalize. Each replay retained the original Event, Entity, and
  Domain result, emitted no new Event, and changed only `replayed` to `true`.
- New-key lifecycle conflicts and same-key/different-digest conflicts retained
  byte-identical typed conflict envelopes and emitted no Event.
- The task 4 report preserves both red streams and observed failure: the first
  same-key replay returned empty stdout and stderr; the intermediate repair
  exposed the optional-ID rendering defect and remained red; the final green
  followed the renderer repair. This rerun is the first fresh green checkpoint
  evidence after the repair, not a transcription of that report.

### DEC-169 Task 9 Repair

The former MF-2 is resolved within task 9's accepted boundary.

- `internal/cli.ExitCode` implements the complete public-interface section 7
  map: invalid input 2, not found 3, conflict and precondition 4, authority 5,
  integrity 6, internal 7, and gate-unsuccessful 8. `FailureError` carries both
  code and category through the renderer boundary; `main.go` exits with that
  mapping. Cobra unknown-command, unknown-flag, and usage failures map to 2.
- Fresh category/typed-error unit tests terminated `exit_code: 0`; all eight
  category cases passed. Fresh real-process `TestProcessExitContract`
  terminated `exit_code: 0` with
  `ok github.com/mbell-nzl/fid/cmd/fid/integration 1.838s`.
- The four fail-closed leaves, AC verify, concern raise, concern address, and
  maintenance vacuum, now retain their exact capability-unavailable envelope
  and terminate 5. A fresh Human-mode vacuum process also terminated 5 with
  `error: capability is unavailable`.
- All 28 pre-repair `cli.WriteFailure` sites were audited. Current production
  source has no direct `cli.WriteFailure(` call (`rg` terminated
  `exit_code: 1`, no output); each error-returning leaf propagates the typed
  result, while Human renderers retain the same category.
- The repair left capability values, domain writers, envelope fields, and
  output writers unchanged. The gate command diff only propagates existing
  registration failures. No live gate-run terminal-frame implementation
  exists before tasks 14-15, so accepted future streaming semantics are
  untouched.

## Temporary Inventory And VAL-043

`VAL-043: PASS`.

The four section 7 classification commands were executed fresh at Group 1
HEAD. Every `role_policy`/role attribution match, every `--by` match, every
temporary-table reference, and every `EventType` reference maps to the exact
retained, replaced, frozen, or current category in
`v1-convergence-inventory.md`. `UNCLASSIFIED: 0`. Final V1 absence remains
correctly deferred to tasks 20-21 and was not imposed at this checkpoint.

Fresh structural isolation evidence:

- `go list -deps ./internal/workflow ./internal/policy ./internal/claims
  ./internal/evidence` terminated `exit_code: 0` and listed neither
  `internal/dev` nor `internal/runs`.
- The direct generic-package import scan for those two packages terminated
  `exit_code: 1`, no output.
- The exact generic source and extracted execution-bound transition-path scan
  for `roles`, `gate_runs`, `task_completions`, `agent_runs`,
  `ac_verifications`, `concerns`, and `internal/dev` terminated
  `exit_code: 1`, no output. The only broader transition match is the retained
  standalone V1 fail-closed reason; source control flow proves that evaluator
  runs only when the transition is not execution-bound.

A scratch-only compiler overlay then seeded hostile rows in all temporary
tables: `roles`, `claims`, `task_completions`, `ac_verifications`,
`agent_runs`, `gate_runs`, `concerns`, and `workflow_transition_requires`.
Both attempts used the exact command:

```text
timeout 120s go test -overlay=/tmp/tmp.tauoSIKOsG/overlay.json -race ./cmd/fid -run '^TestGroup1CheckpointAdversarialTemporaryV1Isolation$' -count=1 -v
```

The first overlay attempt terminated `exit_code: 1` because the QA fixture
omitted the required concern ID and therefore observed canonical invalid-input
instead of the intended fail-closed leaf; its failing token identified the
concern-address probe and the `invalid-argument` category. No product file
changed. After
repairing only the fresh scratch fixture, the byte-identical test command
terminated `exit_code: 0` with these tokens:

```text
PASS hostile rows: roles=1 claims=1 task_completions=1 ac_verifications=1 agent_runs=1 gate_runs=1 concerns=1 workflow_transition_requires=1
PASS generic paths: activation=ACT-2d4094 registration=REG-b760ba envelope=ENV-dfce46 claim=CLM-091ec6 run=active transition=in_progress capabilities=byte-identical replay=true fail-closed=4
ok github.com/mbell-nzl/fid/cmd/fid 4.781s
```

This fresh adversarial run proves unchanged Activation prepare/component/
commit/finalize, Registration, Envelope create/activate, direct Grant, Claim,
Run, nonterminal execution-bound transition, capabilities, and same-key Claim
idempotency despite hostile V1 state. AC verification, concern raise/address,
and maintenance vacuum remained fail-closed and produced no Event. Separate
focused fail-closed and no-writer tests also terminated `exit_code: 0`.

## Command And Feasibility Validation

Fresh scratch was allocated in this checkpoint session by `mktemp -d` as
`/tmp/tmp.tauoSIKOsG`. Setup `go vet ./...` and
`go build -buildvcs=false -o /tmp/tmp.tauoSIKOsG/child-bin ./child` each
terminated `exit_code: 0`. Every FEAS-01..FEAS-18 child command was separately
bounded to 120 seconds. Their observed nested wall times total less than 30
seconds, within the separate 900-second aggregate budget. FEAS-19 was run
separately under its 300-second bound.

| Probe | Result | Fresh token and observed termination |
|---|---|---|
| FEAS-01 | PASS | `argv=["arg-1","arg-2"]`, `out-1`, `out-2`, `err-1`; `exit_code: 0` |
| FEAS-02 | PASS | stdout held only argv/out lines; stderr only `err-1`; each inspection and child `exit_code: 0` |
| FEAS-03 | PASS | `00000000: 41ff fe1b 5b33 316d 0d0a`; writer and `xxd` `exit_code: 0` |
| FEAS-04 | PASS | TERM timeout produced no output; `exit_code: 124`, 4.913s |
| FEAS-05 | PASS | KILL timeout `exit_code: 124`, 4.918s; survivor scan `exit_code: 1`, no output |
| FEAS-06 | PASS | fresh nonexistent executable reported `No such file or directory`; `exit_code: 127` |
| FEAS-07 | PASS | `8388608` bytes on each stream, `16777216 total`; producer and `wc` `exit_code: 0` |
| FEAS-08 | PASS | exact sorted environment `HOME=/tmp`, `PATH=/usr/bin`, `PROBE_MODE=env`; `exit_code: 0` |
| FEAS-09 | PASS | fresh canary `CANARY-g1b-tauoSIKOsG-829` found as `1:token=...`; writer and `rg` `exit_code: 0` |
| FEAS-10 | PASS | `ESCAPES /tmp/tmp.tauoSIKOsG/cwd -> /tmp/tmp.tauoSIKOsG/outside/target.txt`; `exit_code: 0` |
| FEAS-11 | PASS | source changed from `465599f9...cf70` to `d3e9d084...30a2`; copy retained first digest; `exit_code: 0` |
| FEAS-12 | PASS | digest `55b62c00...b6a3`; `ATOMIC_STORE ... tmp-empty object-present`; `exit_code: 0` |
| FEAS-13 | PASS | nested statuses copied directly: FEAS-04 `124`, FEAS-05 `124`, FEAS-06 `127`, FEAS-07 `0` |
| FEAS-14 | PASS | scratch tree OIDs `15e8e4f9...cef99` then `4708f942...4709c`; all git commands `exit_code: 0` |
| FEAS-15 | PASS | `MISSING_OBSERVABLE: PASS`, then `RESTORED_OBSERVABLE: PASS`; each `exit_code: 0` |
| FEAS-16 | PASS | two archives, at least 1.1s apart, both `a68d8fed...5b55`; `cmp` `exit_code: 0` |
| FEAS-17 | PASS after preserved retry | clean `unzip -t` `exit_code: 0`, digest unchanged; first corruption attempt honestly failed the contract with `exit_code: 2`; retry from clean changed local flag byte 7 and produced required `exit_code: 1` |
| FEAS-18 | PASS | all 16 fresh absence shapes produced no output: shapes 1-11 observed expected `exit_code: 1`; shapes 12-16 observed expected `exit_code: 2` |
| FEAS-19 | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 7.536s`; `exit_code: 0`, observed wall 8.206s under 300s |

FEAS-17's first corruption choice zeroed offset 50 and made `unzip -t`
terminate 2 with `bad CRC` and `At least one error was detected`. That failed
attempt is immutable evidence. The successful retry started from the clean
archive, wrote `0x08` at byte 7, and made `unzip -t` report
`mismatch between local and central GPF bit 11` and terminate 1. It closed the
exact probe without changing its accepted status.

Per DEC-167 item 8,
`TestV01ClosureJourneyAndRebuildParity` had never been green in Phase 21 before
its first green run. The fresh FEAS-19 green result above is new evidence, not
a continuity claim.

## Findings And Outcomes

Must-fix findings: none.

Improvement findings: none.

Blocker candidates: none.

- VAL-042: PASS.
- VAL-043: PASS.
- Fresh feasibility FEAS-01..FEAS-19: PASS.
- Independent Group 1 Contract Review: PASS.

CHECKPOINT VERDICT: PASS
