# Phase 24 - Task 11: Readable Stop, Attention and repository lines - Checkpoint Report

## Workflow outcome

- Task outcome: BLOCKED
- Category: specification-contradiction
- Condition: Task 11/VAL-041 requires every recorded-versus-current Human repository view, including run.inspect, to print different worktree for a plain copy selected with --repository. Phase 24 public-interface.md section 2 simultaneously preserves run.inspect's Phase 21-23 behavior, which rejects that copy before rendering; section 10 and the task prohibit changing JSON. The required success view cannot be produced by a Human-renderer-only change.
- Recovery attempts: 1
- Next action: Tech Lead must reconcile the run.inspect plain-copy requirement with the preserved rejection contract by protected decision authority, stating the accepted Human/JSON result and bounded source/test scope. Then resume this same Task 11 from the checkpoint, preserve both failures, finish the real Stop corpus and fixture manifest, update the CLI examples, and run all six exact task gates and independent Runner report validation.
- Human attention required: yes

One same-task session-continuation and two focused fixture reproductions were performed. The external interruption consumes no no-progress budget.

This is a Developer blocker candidate for Phase Runner triage, not a ledger STOP or task completion. No new task, trigger, TRG or ESC is allocated here.

Recovery identity was checked against the latest ledger entry without paraphrasing its values:

```text
Task: task-11
Category: session-continuation
Repair mode: task-scope
Repair authority: task spec
Trigger signature: task-11:session-continuation:developer-session-interrupted
Attempt count: 1
Allowed files: ordinary task-11.md scope
```

Starting phase HEAD was `a6b6203d62ec6061dcac84210f6fef6df13d4ef2`. The 2026-10-09 08:45 TASK-RECOVERY and 08:46 REPAIR INTERRUPTED remain the governing recovery. All twelve dispatched partial source/test files were inspected and retained. This continuation only formatted them and created this report; it did not reset, stash, discard or restart the implementation.

## Summary of what was implemented

- Preserved the interrupted implementation of labelled StopDetails for all ten closed codes, labelled subject/actions/condition/resolver/status, Condition-specific open Next guidance with the after-the-fact disposition step, and resolved fact/resolution Events without Next guidance. Attention now renders its linked Stops, including resolved items.
- Preserved short repository summaries for resume, run.inspect and Review record/get/list, and comparisons using workflow.SameWorktree, including its import-binding rule. p23_human.go changes remain restricted to repository lines.
- Preserved the ordinary renderer/registry corpus and the partial tagged TestP24HumanFixtures generator. It supports a supplied FID_P24_HUMAN_FIXTURE_ROOT or t.TempDir, stores a binary digest and both-mode commands, and constructs all five requested denials and the clean/dirty view matrix.
- Confirmed a protected producer mismatch at the clean plain-copy run.inspect fixture. The generator remains failing evidence: no fixture was removed, changed to an expected failure, narrowed, mocked or skipped to obtain green.
- No behavior implementation was added after identifying the contradiction. Current-task acceptance is incomplete.

## Files created or modified

Preserved incoming authorized partial edits, with gofmt/goimports applied:

- `cmd/fid/integration/p22_human_guidance_test.go`
- `cmd/fid/integration/p23_human_leaf_pairs_test.go`
- `cmd/fid/p23_human.go`
- `cmd/fid/p23_human_test.go`
- `cmd/fid/resume_human.go`
- `cmd/fid/run_inspect_human.go`
- `cmd/fid/stop_human.go`
- `cmd/fid/stop_test.go`
- `cmd/fid/integration/p24_human_fixtures_test.go`
- `cmd/fid/p24_repository_human_test.go`
- `cmd/fid/p24_stop_corpus_test.go`
- `cmd/fid/p24_stop_human_test.go`

New checkpoint report:

- `project-docs/docs/phases/phase-24/tasks/task-11-done.md`

Runner-owned `run-ledger.md` remains outside this checkpoint's edits and staging. No task-10, task-12, task-13, domain, decision, task-spec, QA-contract/plan, dependency, control, prompt, template or generated path was changed.

## Documentation updated

- This checkpoint report records the contradiction, immutable failures, provenance and remaining work.
- The required `project-docs/specs/06_DESIGN_PLAN_CLI.md` Stop/Attention examples remain unfinished. They must be updated in the same-task continuation after authority disposition; no documentation completion is claimed.

## Validation Contract Coverage

- VAL-040: INCOMPLETE. `TestP24StopHumanBudgetLabel`, `TestP24StopHumanCorpus` and `TestP24StopHumanCorpusCompleteness` in the three ordinary Stop test files pass in the focused command below. The corpus checks ten code/Condition pairs in stop.list, attention.list and resume, with every member label, open Next and disposition steps, resolved Events/no Next, and no raw JSON. These are renderer fixtures; the required complete corpus raised through stop raise with real qualifying resolution Events is still unfinished. The real stop list probe below demonstrates only the open operator-stop case.
- VAL-041: BLOCKED. `TestP24RepositoryHumanViews` passes the synthetic clean/dirty and same/different renderer matrix across resume, run.inspect and Review record/get/list. The tagged fixture at `cmd/fid/integration/p24_human_fixtures_test.go:159` and `:175` requires the actual copied run.inspect success, and both focused executions terminate FAIL because the domain rejects it. Real CLI controls below prove same-root run.inspect renders same worktree and plain-copy resume renders different worktree; the matching plain-copy run.inspect calls in both modes exit 4 before any repository view. The synthetic matrix cannot replace the required real plain-copy evidence.
- VAL-052: UNVERIFIED. The full ordinary and tagged gates were not run after the protected contradiction was established. No current cmd/fid, internal/workflow, slowest ordinary package or full tagged integration time is available. Focused cmd/fid 2.404s and focused integration 8.148s/22.536s are not substitutes. Budgets remain no ordinary package over 600s and tagged integration under 900s; headroom warnings start at 480s ordinary and 720s tagged. No budget PASS is claimed.
- VAL-054: PREPARATION INCOMPLETE; Human judgment pending. The generator executes all five denial constructors before reaching the copied run.inspect failure: idempotency-conflict, repository-mismatch, invalid-argument naming input.expected_repository.dirty, source-digest-mismatch and protected-decision-required with blocked by. Its manifest write/validation is never reached, so the manifest is not acceptance evidence and the Human's ability to explain the denials is not attested.
- VAL-055: PREPARATION INCOMPLETE; Human judgment pending. The generator's clean same/different stop list, attention list, resume and review get calls precede the failure, but the dirty matrix and finalized manifest are not verified. Window D2 and the Human walkthrough remain required after Task 11 completes.

## Tests run

All commands in this section were run from `/home/mark/projects/fid` unless their real-entrypoint cwd is stated. Every current receipt is copied from the nested command's own observed termination.

- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24(StopHuman|RepositoryHuman)' -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	2.404s` (terminated `exit_code: 0`)

This validates incoming partial renderer work, not a new red/green implementation in this continuation.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

Package time 8.148s; test time 7.19s. Root variable unset, exercising t.TempDir. The failure was clean-different-run.inspect, JSON exit 4 versus required success exit 0, with recovery-state-conflict and reasons:[repository identity changed].
- `FID_P24_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/task-11-blocker-evidence go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` - FAIL - `FAIL` (terminated `exit_code: 1`)

Package time 22.536s; test time 21.99s. The same assertion fails with the root variable set, retaining the workspace and binary for real-entrypoint diagnosis. The variable does not evade the test.
- `git merge-base --is-ancestor ac5e6f9bec80e06bfe5b7b33e3611a894d16038a 304aa45f12cf3966567231e6986e3070df409c68` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

The rejected repository-identity path predates phase approval; this does not excuse the current assertion or assign a different task.

Real entrypoint interactions used the generator's actual executable at `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-1289185125/fid`, cwd and project `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-1289185125/clean-same`, and plain copy `/home/mark/projects/fid/.cache/task-11-blocker-evidence/fixtures-1289185125/clean-different`. These interactions were directly observed through exec_command, separately from the Go tests:

- Human `stop list` with the stated binary and `--project` - PASS - `  Next: fid run close or fid run abandon; after the fact: fid stop resolve, then fid attention resolve where no lifecycle Event disposed it.` (terminated `exit_code: 0`)

Output also contains labelled Subject, Code, Blocked actions, Reason, Reference, Resolution condition, Resolver capability and Status for STP-cdf9dc, with no raw JSON.
- Human `run inspect RUN-d5c15a --repository <clean-same absolute path above>` with the stated binary and `--project` - PASS - `  No observation recorded for this Run; whether its executor is live is unknown.` (terminated `exit_code: 0`)

Both repository lines are `HEAD b4ded0da37b2; clean`; comparison is `same worktree`.
- Human `resume P22-d9d5d6 --repository <clean-different absolute path above>` with the stated binary and `--project` - PASS - `Budget tokens: remaining unknown input-tokens; consumed unknown; maximum 10; state unknown` (terminated `exit_code: 0`)

Current, recorded and Review repositories use `HEAD b4ded0da37b2; clean`; comparison is `different worktree`.
- Human `run inspect RUN-d5c15a --repository <clean-different absolute path above>` with the stated binary and `--project` - FAIL against VAL-041 success-view requirement - `reason: recovery-state-conflict; observation_digest=null, paths=, reasons=repository identity changed, run_id=RUN-d5c15a` (terminated `exit_code: 4`)

Output is a failure, not a repository comparison view.
- JSON `--json run inspect RUN-d5c15a --repository <clean-different absolute path above>` with the stated binary and `--project` - FAIL against VAL-041 success-view requirement - `{"contract_version":"fid-cli.v1","error":{"category":"conflict","code":"recovery-state-conflict","details":{"observation_digest":null,"paths":[],"reasons":["repository identity changed"],"run_id":"RUN-d5c15a"},"message":"recovery state conflict","retryable":false,"subject":{"id":"RUN-d5c15a","type":"run"}},"ok":false,"operation":"run.inspect"}` (terminated `exit_code: 4`)

This is the accepted existing typed rejection; no protected failure mapping was changed.

## Exact quality-gate results

The task spec contains six exact commands. Ledger prose saying seven does not authorize an invented seventh command. Full task closure is blocked by the protected contradiction; no focused selection replaces the required gates.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

Run after gofmt and goimports on all twelve edited Go files.
- `golangci-lint run ./...` - BLOCKED - Not run: protected contract disposition is required before the final-tree gate; no token or termination result exists for this command in this session.
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - BLOCKED - Not run: protected contract disposition is required before the full task gate; no current ordinary-package wall times, token or termination result exist.
- `go tool cover -func=coverage.out` - BLOCKED - Not run: this session has no completed required coverage gate; an earlier coverage.out would not be fresh Task 11 evidence.
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - BLOCKED - Not run: the required real plain-copy view needs protected authority; focused fixture failures above are retained but do not substitute for this gate.
- `rg -n -i "worktree|dirty digest|HeadOID|head_oid" cmd/fid -g '!*_test.go'` - BLOCKED - Not run as the final ordered inventory gate; preliminary context inventory found the Human renderers and JSON/input sites, but no final-tree inventory PASS is claimed.

QA-plan sections 8 and 9 require independent Window D2 rechecks, fixture regeneration at the Task 11 close commit and Human judgment. Those are pending because Task 11 is not complete; no additional Developer-owned QA gate is attested. The Runner must independently validate this committed report. No Runner report-lint receipt is claimed here.

## Deviations from spec

- Developer report validation first ended `FAIL report-lint found 4 failure(s)` (terminated `exit_code: 4`). Root cause: Recovery attempts contained explanatory prose instead of an integer, and three receipt lines had explanatory prose after the termination field. Repair: retain the exact observed tokens/statuses, make the typed count exactly 1, and move explanations to separate paragraphs. This is report-format repair only; the BLOCKED outcome and every failed fixture/probe remain. The corrected committed report requires a fresh Developer check and independent Runner validation.

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

- Resolution must uniquely specify whether run.inspect on a plain copy retains typed rejection or returns an inspectable comparison, and reconcile VAL-041, section 2 preservation and unchanged JSON. A file-list extension alone cannot settle this semantic conflict. Do not replace the copy with an import, relabel it same worktree, delete the row or change its expected exit under current authority.
- Ordinary work remains: finish the real Stop raise/resolution corpus, both-state fixture manifest, retained Human fixtures, CLI examples, pre-gate touched-package checks and six ordered final gates. No completeness claim is made for the inherited partial implementation and no new task owner is assigned.
- Full-gate headroom is UNVERIFIED. Task 10's committed report records cmd/fid 484.734s, above the 480s warning, internal/workflow 116.440s, internal/issues 216.748s and tagged integration 335.908s. These are historical risk context, not Task 11 PASS receipts. DEC-285/286 built-binary failure-pair rechecks remain mandatory in Window D2.
- The durable diagnostic executable/workspaces are retained under `.cache/task-11-blocker-evidence/fixtures-1289185125/`; generation aborted before manifest.json was written. The original focused TempDir attempt was ephemeral. Regenerate fixtures after resolution rather than presenting this incomplete directory as a valid Human manifest.
- Runner-owned ledger was untouched, unstaged and unhidden; its observed pre-checkpoint SHA-256 was `7df2a146f19e9f5a906caa5c642bdd936f3a9f9a6a2fdf299da555977db5fe94`. The checkpoint preserves all dispatched partial files and the original failing assertion. No held-out scenario or evaluator prompt was read.
- Developer report validation, if run, is separate from independent Runner control and cannot advance the task or replace the Runner receipt.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `ab05e496206809860081f941d6c1c171eb845cca`

## Commit subject

- `feat(cli): checkpoint task 11 Human views and fixture contract blocker`
