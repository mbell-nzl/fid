# Phase 24 - Task 11: Readable Stop, Attention and repository lines - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 2
- Next action: None.
- Human attention required: no

Refs: phase-24/task-11.

The latest TASK-RECOVERY identity was read and confirmed byte-for-byte:

```text
Task: task-11
Trigger signature: task-11:public-interface:plain-copy-run-inspect-repository-view
Attempt count: 2
Category: in-scope-implementation
Repair mode: task-scope
Repair authority: TRG-24-012 / DEC-287
Consecutive no-progress attempts: 0/2
Allowed files: ordinary task-11.md scope
```

This continuation started on branch `phase-24`, HEAD `28dbad27ec666391b7848956830f1835a2cb3cb3`, preserving implementation checkpoint `ab05e496206809860081f941d6c1c171eb845cca` and report checkpoint `9ce8896fd28d4aa0349cdbe7447ac7256e05cfe7`. DEC-287 resolves the protected contradiction below and authorizes only the two plain-copy inspect fixture expectations. It does not change repository behavior or the Human failure renderer. Runner-owned run-ledger.md was never edited, staged, committed or hidden.

**Historical blocked checkpoint, preserved as originally recorded.** Its earlier PASS receipts are historical evidence, UNVERIFIED as fresh results in this session; current observed receipts appear below. The 08:45/08:46 interrupted-session chain remains history and consumes no no-progress budget.

> - Task outcome: BLOCKED
> - Category: specification-contradiction
> - Condition: Task 11/VAL-041 requires every recorded-versus-current Human repository view, including run.inspect, to print different worktree for a plain copy selected with --repository. Phase 24 public-interface.md section 2 simultaneously preserves run.inspect's Phase 21-23 behavior, which rejects that copy before rendering; section 10 and the task prohibit changing JSON. The required success view cannot be produced by a Human-renderer-only change.
> - Recovery attempts: 1
> - Next action: Tech Lead must reconcile the run.inspect plain-copy requirement with the preserved rejection contract by protected decision authority, stating the accepted Human/JSON result and bounded source/test scope. Then resume this same Task 11 from the checkpoint, preserve both failures, finish the real Stop corpus and fixture manifest, update the CLI examples, and run all six exact task gates and independent Runner report validation.
> - Human attention required: yes
>
> One same-task session-continuation and two focused fixture reproductions were performed. The external interruption consumes no no-progress budget.
>
> This is a Developer blocker candidate for Phase Runner triage, not a ledger STOP or task completion. No new task, trigger, TRG or ESC is allocated here.
>
> Recovery identity was checked against the latest ledger entry without paraphrasing its values:
>
> ```text
> Task: task-11
> Category: session-continuation
> Repair mode: task-scope
> Repair authority: task spec
> Trigger signature: task-11:session-continuation:developer-session-interrupted
> Attempt count: 1
> Allowed files: ordinary task-11.md scope
> ```
>
> Starting phase HEAD was `a6b6203d62ec6061dcac84210f6fef6df13d4ef2`. The 2026-10-09 08:45 TASK-RECOVERY and 08:46 REPAIR INTERRUPTED remain the governing recovery. All twelve dispatched partial source/test files were inspected and retained. This continuation only formatted them and created this report; it did not reset, stash, discard or restart the implementation.

## Summary of what was implemented

- Continued the preserved labelled Stop/Attention renderers, Condition-specific Next guidance, resolution Event rendering and short Human repository lines. Product sources from the implementation checkpoint remain unchanged in attempt 2.
- Corrected only clean-different-run.inspect and dirty-different-run.inspect under DEC-287 clause 4: expected exit 4, recovery-state-conflict, empty stdout in both modes, exact closed failure envelope, Run subject, repository identity changed reason, exact Human stderr and p24ExpectedFailureText parity. Both read-only calls append no Event. Every other original row remains, including plain-copy resume and same-root inspect/resume.
- Added TestP24StopHumanRealCorpus: ten real stop.raise variants and ten Conditions, with qualifying Events from terminated executor observation, failed-Run recovery, finalized activation, budget observation, terminal Run closure, independent Review/Finding resolution, Human Outcome and protected Decision completion. It checks every StopDetails member, subject/actions/condition/resolver/status, open Next plus disposition steps, resolved fact/resolution Events without Next, Attention linkage/status and read-only views.
- Completed and replay-validated the Human fixture manifest: five denial names and all 28 clean/dirty × same/plain-copy view names, with workspace/cwd, both command lines, exit expectations and binary SHA-256. Both unset-root and retained-root generation passed; retained Phase 23 fixtures and import guidance passed unchanged in this continuation.
- Updated Stop/Attention examples and repository guidance in 06_DESIGN_PLAN_CLI.md. Ran the touched-package checks, direct real CLI interactions and all six exact gates in order.

## Files created or modified

Preserved Task 11 implementation checkpoint files (unchanged in attempt 2):

- `cmd/fid/integration/p22_human_guidance_test.go`
- `cmd/fid/integration/p23_human_leaf_pairs_test.go`
- `cmd/fid/p23_human.go`
- `cmd/fid/p23_human_test.go`
- `cmd/fid/resume_human.go`
- `cmd/fid/run_inspect_human.go`
- `cmd/fid/stop_human.go`
- `cmd/fid/stop_test.go`
- `cmd/fid/p24_repository_human_test.go`
- `cmd/fid/p24_stop_corpus_test.go`
- `cmd/fid/p24_stop_human_test.go`

Modified in attempt 2:

- `cmd/fid/integration/p24_human_fixtures_test.go`
- `project-docs/specs/06_DESIGN_PLAN_CLI.md`
- `project-docs/docs/phases/phase-24/tasks/task-11-done.md`

## Documentation updated

- 06_DESIGN_PLAN_CLI.md now shows open Stop and Attention examples, resolved Event lines without Next, all ten Condition commands/wait guidance and the short repository form with the retained plain-copy inspect refusal.
- This report preserves the prior blocked outcome, failed probes, deviations, provenance and interrupted-session history, and adds current repair and gate evidence.

## Validation Contract Coverage

- VAL-040: PASS (deterministic Developer evidence). TestP24StopHumanCorpus and TestP24StopHumanCorpusCompleteness (`cmd/fid/p24_stop_corpus_test.go:35`, `:110`) bind all ten codes and Conditions to the source registries and render all three views, open and resolved. TestP24StopHumanRealCorpus (`cmd/fid/integration/p24_human_fixtures_test.go:73`) raises every code through the real CLI; p24RealStopCase (`:113`) records qualifying resolution and both manual and atomic disposition; p24RealStopFact (`:196`) returns actual journal Event IDs; p24AssertRealStopViews (`:250`) checks each member's label/value, open guidance, resolved Events/no Next and Attention states against the actual CLI views. Real resume preserves its accepted projection: a resolved linked Stop remains while its Attention is open; disposed judgment items and recovered predecessors leave the projection. The new test proves Human/JSON absence after that filtering, while the retained renderer corpus proves all ten variants in both states without changing JSON or lifecycle semantics.
- VAL-041: PASS. TestP24RepositoryHumanViews (`cmd/fid/p24_repository_human_test.go:13`) retains the synthetic clean/dirty × same/different matrix across resume, inspect and Review record/get/list. TestP24HumanFixtures generates both real states, seven views and both worktrees. The DEC-287 branch (`cmd/fid/integration/p24_human_fixtures_test.go:438`) proves both copied inspect refusal pairs, exact Run subject/reasons/metadata/bytes, empty streams and no Event. The original same-root inspect/resume and different-root resume checks still assert their respective comparison labels and forbid full HEAD/identity/digest/nil leakage. Direct dirty-root inspect and plain-copy resume probes below independently demonstrate the real entrypoint. TestP22HumanGuidance preserves the verified imported-destination same-worktree assertion.
- VAL-052: PASS. Full ordinary gate package times: cmd/fid 390.671s; internal/workflow 105.826s; slowest other freshly timed ordinary package internal/export 136.519s (internal/db 53.914s, internal/project 1.976s); tagged cmd/fid/integration 338.886s. Other unchanged ordinary packages reported `(cached)` and supplied no fresh package time; no timing is reconstructed for them. No ordinary package exceeds 600s; tagged integration is below 900s. Headroom notes appear below.
- VAL-054: Developer preparation complete; Human judgment remains for Window D2 and phase end. TestP24HumanFixtures and p24AssertHumanManifest (`cmd/fid/integration/p24_human_fixtures_test.go:550`) prove all five denial names, reproducible Human/JSON pairing and runnable retained commands, including nested input.expected_repository.dirty and protected-decision-required with blocked by. The Human's ability to explain each denial is not attested here.
- VAL-055: Developer preparation complete; Human judgment remains for Window D2 and phase end. The same manifest proves all 28 required view names in clean/dirty and same/plain-copy states, including stop.list, attention.list, resume and review.get. Its two copied inspect rows are the exact DEC-287 refusals. QA must regenerate fixtures at the close commit and obtain the Human walkthrough; this report does not waive either step.

## Tests run

Working directory for all test/gate commands: `/home/mark/projects/fid`. Current results below were observed to terminate through the nested execution result. No inherited PASS receipt is presented as a fresh run.

**Current red/green and touched-package checks:**

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

Attempt-2 pre-edit reproduction: 9.099s, clean-different-run.inspect actual exit 4 versus original expected exit 0, recovery-state-conflict with repository identity changed. This failure is preserved.

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	14.424s` (terminated `exit_code: 0`)
- `FID_P24_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/task-11-blocker-evidence go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	27.943s` (terminated `exit_code: 0`)

Both commands are byte-identical to the prior failed probes below, after the authorized DEC-287 correction. The retained new manifest is `.cache/task-11-blocker-evidence/fixtures-102722463/manifest.json`; 33 rows and binary digest `sha256:34ba1e58f07393f58ffa1de037555d2c3b3c78baad8d3fb5da942401fefbb180`. The original incomplete blocker directory `fixtures-1289185125` remains intact.

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24StopHumanRealCorpus$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

Initial corpus-completeness red, 1.695s: real stop.raise fixture contained one code and one Condition, expected ten of each.

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24StopHumanRealCorpus$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

First ten-row construction, 6.993s: invalid short observation digest, incorrect budget blocked_actions, unavailable capability name and newly authored resume assertions that assumed disposed items remain in the canonical projection. No retained assertion changed.

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24StopHumanRealCorpus$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

Second construction, 7.740s: an abandoned Run cannot qualify as a recoverable predecessor. Fixed the fixture to close a failed Run, release its Claim, and record actual recovery.

- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24StopHumanRealCorpus$' -count=1 ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	9.881s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^Test(P24HumanFixtures|P24StopHumanRealCorpus|P23HumanFixtures|P23HumanFixtureStepsRemainRunnable|P23HumanLeafPairs|P22HumanGuidance|P22StopAttentionHumanListsNameReasonAndResolution)$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	23.272s` (terminated `exit_code: 0`)

This final tagged pre-gate includes the post-green fixture refactor and retained P23 fixture/step tests, which execute the leaf-pair oracle; P23HumanLeafPairs is only a selector alternative, not a standalone test claim.

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	312.556s` (terminated `exit_code: 0`)

**Direct real-entrypoint probes:** executable `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/fid`. Cwd is the command's `--project` directory (`clean-same` for stop list, `dirty-same` for all other probes). Exit-4 refusals are observed expected contract results, not successful process exits.

- `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/fid --project /home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/clean-same stop list` - PASS - `  Next: fid run close or fid run abandon; after the fact: fid stop resolve, then fid attention resolve where no lifecycle Event disposed it.` (terminated `exit_code: 0`)
- `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/fid --project /home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/dirty-same run inspect RUN-e64bcf --repository /home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/dirty-same` - PASS - `  No observation recorded for this Run; whether its executor is live is unknown.` (terminated `exit_code: 0`)
- `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/fid --project /home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/dirty-same resume P22-9e528d --repository /home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/dirty-different` - PASS - `Budget tokens: remaining unknown input-tokens; consumed unknown; maximum 10; state unknown` (terminated `exit_code: 0`)
- `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/fid --project /home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/dirty-same run inspect RUN-e64bcf --repository /home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/dirty-different` - EXPECTED REFUSAL - `reason: recovery-state-conflict; observation_digest=null, paths=, reasons=repository identity changed, run_id=RUN-e64bcf` (terminated `exit_code: 4`)
- `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/fid --json --project /home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/dirty-same run inspect RUN-e64bcf --repository /home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-102722463/dirty-different` - EXPECTED REFUSAL - `{"contract_version":"fid-cli.v1","error":{"category":"conflict","code":"recovery-state-conflict","details":{"observation_digest":null,"paths":[],"reasons":["repository identity changed"],"run_id":"RUN-e64bcf"},"message":"recovery state conflict","retryable":false,"subject":{"id":"RUN-e64bcf","type":"run"}},"ok":false,"operation":"run.inspect"}` (terminated `exit_code: 4`)

Dirty-root inspect printed `Recorded repository: HEAD 820cdb62a5a7; clean`, `Current repository: HEAD 820cdb62a5a7; dirty` and `Worktree comparison: same worktree`. Plain-copy resume printed short current/recorded/Review lines and `Worktree comparison: different worktree`. Both refusal probes had empty stdout and exactly the DEC-287 stderr. Human stop list showed every operator-stop label and the complete Next/disposition line.

**Historical tests and built-binary probes, verbatim from the blocked report.** These retain their original outcomes; their historical PASS entries are UNVERIFIED as current-session executions and are not used to close a gate:

> All commands in this section were run from `/home/mark/projects/fid` unless their real-entrypoint cwd is stated. Every current receipt is copied from the nested command's own observed termination.
>
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24(StopHuman|RepositoryHuman)' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.404s` (terminated `exit_code: 0`)
>
> This validates incoming partial renderer work, not a new red/green implementation in this continuation.
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
>
> Package time 8.148s; test time 7.19s. Root variable unset, exercising t.TempDir. The failure was clean-different-run.inspect, JSON exit 4 versus required success exit 0, with recovery-state-conflict and reasons:[repository identity changed].
> - `FID_P24_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/task-11-blocker-evidence go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)
>
> Package time 22.536s; test time 21.99s. The same assertion fails with the root variable set, retaining the workspace and binary for real-entrypoint diagnosis. The variable does not evade the test.
> - `git merge-base --is-ancestor ac5e6f9bec80e06bfe5b7b33e3611a894d16038a 304aa45f12cf3966567231e6986e3070df409c68` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
>
> The rejected repository-identity path predates phase approval; this does not excuse the current assertion or assign a different task.
>
> Real entrypoint interactions used the generator's actual executable at `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-1289185125/fid`, cwd and project `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-1289185125/clean-same`, and plain copy `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-1289185125/clean-different`. These interactions were directly observed through exec_command, separately from the Go tests:
>
> - Human `stop list` with the stated binary and `--project` - PASS - `  Next: fid run close or fid run abandon; after the fact: fid stop resolve, then fid attention resolve where no lifecycle Event disposed it.` (terminated `exit_code: 0`)
>
> Output also contains labelled Subject, Code, Blocked actions, Reason, Reference, Resolution condition, Resolver capability and Status for STP-cdf9dc, with no raw JSON.
> - Human `run inspect RUN-d5c15a --repository <clean-same absolute path above>` with the stated binary and `--project` - PASS - `  No observation recorded for this Run; whether its executor is live is unknown.` (terminated `exit_code: 0`)
>
> Both repository lines are `HEAD b4ded0da37b2; clean`; comparison is `same worktree`.
> - Human `resume P22-d9d5d6 --repository <clean-different absolute path above>` with the stated binary and `--project` - PASS - `Budget tokens: remaining unknown input-tokens; consumed unknown; maximum 10; state unknown` (terminated `exit_code: 0`)
>
> Current, recorded and Review repositories use `HEAD b4ded0da37b2; clean`; comparison is `different worktree`.
> - Human `run inspect RUN-d5c15a --repository <clean-different absolute path above>` with the stated binary and `--project` - FAIL against VAL-041 success-view requirement - `reason: recovery-state-conflict; observation_digest=null, paths=, reasons=repository identity changed, run_id=RUN-d5c15a` (terminated `exit_code: 4`)
>
> Output is a failure, not a repository comparison view.
> - JSON `--json run inspect RUN-d5c15a --repository <clean-different absolute path above>` with the stated binary and `--project` - FAIL against VAL-041 success-view requirement - `{"contract_version":"fid-cli.v1","error":{"category":"conflict","code":"recovery-state-conflict","details":{"observation_digest":null,"paths":[],"reasons":["repository identity changed"],"run_id":"RUN-d5c15a"},"message":"recovery state conflict","retryable":false,"subject":{"id":"RUN-d5c15a","type":"run"}},"ok":false,"operation":"run.inspect"}` (terminated `exit_code: 4`)
>
> This is the accepted existing typed rejection; no protected failure mapping was changed.

## Exact quality-gate results

Exactly six task-spec commands were run, in this order, on the final source/test tree. The older ledger prose saying seven supplied no seventh task command. Every receipt uses its own final printed line and observed nested exit_code.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	105.826s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					62.8%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	338.886s` (terminated `exit_code: 0`)
- `rg -n -i "worktree|dirty digest|HeadOID|head_oid" cmd/fid -g '!*_test.go'` - PASS - `cmd/fid/integration/testdata/p21_conformance/manifest.json:144:    "worktree"` (terminated `exit_code: 0`)

R11 match accounting: run_inspect_human.go:17,18 and resume_human.go:124,126 call repositoryHumanSummary; p23_human.go:139 reads the JSON HEAD before the same short formatter. run_inspect_human.go:23,109,110,112 and resume_human.go:25 are the accepted SameWorktree predicate and same/different worktree labels. review.go:144,184, run.go:269 and gate.go:300 map JSON expectations/observations/failure details. input_descriptors_review.go:24, input_descriptors_run.go:12,53 and input_descriptors_binding.go:75,88 are accepted input descriptions/fields; run_recovery.go:68 is the repository-path flag description. integration/testdata/p21_conformance/manifest.json:144 is the retained Phase 21 command-inventory parameter, not Human output. All 21 matches are accounted for; the Human matrix proves no long-identity/digest/nil leakage.

QA-plan sections 8 and 9 leave independent Window D2 rechecks, fixture regeneration at the Task 11 close commit and Human judgment to QA/Runner. There is no additional Developer-owned task-level QA command. Developer report-lint below is separate from the independent Runner receipt and does not advance the ledger.

**Preserved historical gate disposition, before DEC-287:**

> The task spec contains six exact commands. Ledger prose saying seven does not authorize an invented seventh command. Full task closure is blocked by the protected contradiction; no focused selection replaces the required gates.
>
> - `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
>
> Run after gofmt and goimports on all twelve edited Go files.
> - `golangci-lint run ./...` - BLOCKED - Not run: protected contract disposition is required before the final-tree gate; no token or termination result exists for this command in this session.
> - `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - Not run: protected contract disposition is required before the full task gate; no current ordinary-package wall times, token or termination result exist.
> - `go tool cover -func=coverage.out` - BLOCKED - Not run: this session has no completed required coverage gate; an earlier coverage.out would not be fresh Task 11 evidence.
> - `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - Not run: the required real plain-copy view needs protected authority; focused fixture failures above are retained but do not substitute for this gate.
> - `rg -n -i "worktree|dirty digest|HeadOID|head_oid" cmd/fid -g '!*_test.go'` - BLOCKED - Not run as the final ordered inventory gate; preliminary context inventory found the Human renderers and JSON/input sites, but no final-tree inventory PASS is claimed.
>
> QA-plan sections 8 and 9 require independent Window D2 rechecks, fixture regeneration at the Task 11 close commit and Human judgment. Those are pending because Task 11 is not complete; no additional Developer-owned QA gate is attested. The Runner must independently validate this committed report. No Runner report-lint receipt is claimed here.

## Deviations from spec

- Causal DEC-287 repair: the first focused failure remained `FAIL` (terminated `exit_code: 1`, 9.099s), matching the preserved 8.148s/22.536s failures. Root cause was the planning contradiction documented in DEC-287: the original plain-copy inspect success expectation conflicted with the retained producer refusal. Clause 4 authorizes correcting only the clean/dirty different-run.inspect rows. Both byte-identical commands now pass (14.424s/27.943s); all other rows and the five/28 manifest name inventory remain. No InspectRun, reconciliation, import binding, SameWorktree, failure mapping or section 9 renderer was edited.
- New real-corpus red/green construction retained all failures (`FAIL`, exit_code 1, 1.695s/6.993s/7.740s). Fixture repair used the accepted digest shape, exact budget action set, available capability registry ID and actual failed-Run recovery sequence. Real resume assertions follow the preserved canonical projection rather than inventing a history view. No retained assertion was deleted, weakened, replaced, narrowed or re-pointed. The all-variant renderer oracle remains unchanged.

**All prior deviations preserved verbatim as historical evidence.** Statements that authority or later green was absent describe the blocked checkpoint, not this DEC-287 continuation:

- Developer report validation first ended `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`). Root cause: Recovery attempts contained explanatory prose instead of an integer, and three receipt lines had explanatory prose after the termination field. Repair: retain the exact observed tokens/statuses, make the typed count exactly 1, and move explanations to separate paragraphs. This is report-format repair only; the BLOCKED outcome and every failed fixture/probe remain. The byte-identical Developer report command subsequently ended `PASS report-lint completed with no failures` (terminated `exit_code: 0`) and printed `OUTCOME BLOCKED`. This is Developer report validation, not an independent Runner receipt; Runner validation remains pending.

- Protected contradiction: task-11.md:91-95 and VAL-041 require a successful different-worktree repository view for every recorded-versus-current R11 view. Phase 24 public-interface.md section 2 retains run.inspect's previous behavior; section 10 and the task keep JSON unchanged. InspectRun calls reconcileRunRepository and returns its error before composing RunInspection (`internal/workflow/run_inspect.go:45`); reconcileRepository rejects an unbound changed WorktreeID at `internal/workflow/reconciliation.go:53-54`. A plain copy has no import binding. Human renderer edits cannot observe or render a result that the producer never returns. Altering the domain rejection, JSON result/status, or reducing the accepted view corpus requires protected decision authority. No such repair was made and there is no later green to report.
- Failed attempts remain immutable: the unset-root focused command ended `FAIL` (terminated `exit_code: 1`, 8.148s), and the set-root command ended `FAIL` (terminated `exit_code: 1`, 22.536s). Both have the same repository-identity root cause. The second run retains a concrete fixture for diagnosis; it is not a mined green. Built-binary probes confirm Human and JSON exit 4, and successful same-root run.inspect/plain-copy resume controls isolate the producer mismatch. These are focused failures; no required long-gate failure is invented.
- External interruption evidence is preserved from the governing ledger: Developer session `01a11d00-ae11-7240-8714-4245a84373dd` was lost after server restart; wait_agent returned `not_found`. It returned no report, commit or typed task outcome. Elapsed session duration and command in flight are unknown. There is no observed prior-command termination available to this replacement, so prior TDD or gate results are UNVERIFIED, never PASS. The specified first Budget-label test is present and currently green; its original red transcript cannot be reconstructed from partial code.
- Historical Runner rejection evidence, not fresh Developer command attestation: the 08:45 ledger retains `INVALID REPAIR Trigger signature does not match TASK-RECOVERY`; `INVALID REPAIR Attempt count does not match TASK-RECOVERY`; `INVALID REPAIR has no preceding TASK-RECOVERY` with its recorded exit_code: 1. The Runner left the ledger unchanged and recorded VALID state=ADVANCEABLE, then appended the matching 08:46 INTERRUPTED event. No copy of that history is counted as a fresh PASS or a new recovery attempt here.
- Retained-test convergence was already present on entry and is preserved with introducing-commit provenance, established using git log -S, file-addition history and commit bodies:

| Retained assertion | Introducing commit and provenance | Accepted Task 11 change and preserved strength |
|---|---|---|
| Review repository facts in p23_human_test.go | `59d2653da7cfe14699dc4d1cd85ef032bec26a5c`, Refs: phase-23/task-14; expanded exact leaf assertions in `b71ed1ba18a277da22a71c5fdcf18eadb181be26`, same phase/task | Long identity/capture lines become exact 12-hex HEAD/clean lines; all other golden facts remain. New matrix forbids raw identity/digest/nil leaks and tests dirty output. |
| Review repository lines in p23_human_leaf_pairs_test.go, including TestP23HumanFixtures oracle | `b71ed1ba18a277da22a71c5fdcf18eadb181be26`, Refs: phase-23/task-14 | The JSON-derived oracle renders short HEAD and actual clean/dirty state; full Human byte equality and JSON fact/replay equality remain. Its complete retained generator still requires a fresh green after disposition. |
| Imported run.inspect repository labels in p22_human_guidance_test.go | `eff28fb1b7ea0a543017bef00ccb183cc79fc396`, Refs: phase-22/task-19 | Recorded/Current repository labels replace exposed worktree identity; same worktree is additionally asserted for a verified import. Historical/fresh liveness and read-only checks remain. |
| Resolved Attention result in stop_test.go | `a0a93c7ce4dd0c7dfda4eb39a6ba40760068e460`, Refs: phase-22/task-9 | Supply the existing Stop as known fixture data because resolved Attention must now render linked Stops too. The Attention ID/code/resolution-Event assertions remain intact. |

- The repository rejection itself was introduced by `ac5e6f9bec80e06bfe5b7b33e3611a894d16038a`, `feat(recovery): add repository reconciliation observation`, Refs: phase-22/task-4. Its ancestry before Phase 24 approval `304aa45f12cf3966567231e6986e3070df409c68` was freshly observed with exit_code: 0. VAL-041 was introduced by that Phase 24 approval commit, whose body explicitly approves the contract and generalizes imported-destination same-worktree meaning. This provenance is evidence of the conflict, not a routing of repair to a historical task.

## Concerns, technical debt, or notes for the Tech Lead

- No headroom warning on current recorded times: cmd/fid 390.671s is below 480s, and tagged integration 338.886s is below 720s. The earlier Task 10 cmd/fid 484.734s warning stays historical context, not current timing evidence.
- Window D2, DEC-285/286 failure-pair rechecks, close-commit fixture regeneration and the Human walkthrough remain required. Human VAL-054/055 judgment is pending; deterministic preparation is complete.
- Resume keeps its canonical selection of open Stops/open Attention and current Run. Real-corpus tests prove resolved linked-Stop output while Attention remains open, plus absence after lifecycle disposition and predecessor replacement. The unchanged ordinary renderer corpus proves all ten resolved variants for all three renderers; no JSON or selection semantics were changed to force unreachable real rows.
- Runner ledger SHA-256 before/after this continuation is `b81acd682ca84e38f24ce876f745d38f91ad2e574d061adcb0d61acde69086b0`. No hidden index flags were used, and no held-out scenario or evaluator prompt was read.
- Independent Runner report verification remains separate from Developer report validation and task completion dispatch.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `e6b49797cc1dba22012758922778bed805f8f16c`

## Commit subject

- `test(cli): complete Stop corpus and plain-copy Human fixtures`
