# Phase 24 - Mid-Phase Review 9 (tasks 10-11)

Reviewer: qa-reviewer
Date: 2026-10-09
Tasks covered: phase-24/task-10, phase-24/task-11
Review type: INITIAL
Prior review: none for Window D2; preceding Window D1 re-review is qa/review-8.md (PASS)

## Verdict

BLOCK

The deterministic Task 10/11 evidence is on track and all fresh D2 checks
passed. The required interim Human walkthrough has not happened. It requires
Human product judgment, outside Developer and QA Reviewer authority; it is
not an implementation repair. Window D2 must not advance until that checkpoint
is recorded. No Human acceptance or rejection is inferred from test results.

## Scope Reviewed

- Task specs and reports: `project-docs/docs/phases/phase-24/tasks/task-10.md`,
  `task-10-done.md`, `task-11.md`, `task-11-done.md`.
- QA plan: `project-docs/docs/phases/phase-24/qa/qa-plan.md`, especially
  sections 4, 5, 7, 8 and 9. Validation contract:
  `project-docs/docs/phases/phase-24/qa/validation-contract.md`.
- Accepted interface and inventory: `project-docs/docs/phases/phase-24/public-interface.md`
  sections 2, 9 and 10; `input-inventory.md` R1-R20 and its recursive walk.
- DEC-285, DEC-286 and DEC-287 read in full, including their mapping tables,
  prohibited changes, preserved assertions and Human checkpoint requirements.
- AGENTS.md, qa-evidence skill, all four required current context documents,
  `08_QUALITY_PLAN.md`, `09_AI_WORKFLOW.md` sections 2.7 and 8, and the
  mid-phase review template. CLI assessment uses `05_DESIGN_PLAN.md` sections
  2, 3, 5, 6 and 7, with delivered examples in `06_DESIGN_PLAN_CLI.md`.
- Prior Window D1 review-7 and independent closure review-8, and Window C
  review-6's recorded PASS. These are preceding-window history, not a D2
  re-review or reopened repair.

Reviewed branch `phase-24` at Task 11 close/current HEAD
`93a1f792d27f001b4052113e9056faf97706c3e5`. Task 10 close is
`a6b6203d62ec6061dcac84210f6fef6df13d4ef2`. The sole initial tracked dirt was
Runner-owned `project-docs/docs/phases/phase-24/run-ledger.md`; it remains
untouched. Current brief/progress snapshots still describe Task 11 recovery;
the inspected close commits and committed COMPLETE reports establish this
review's state. Those snapshots were not edited.

| Checkpoint | Verified identity and assessment |
|---|---|
| Task 10 implementation | `d11bec713c470596131190c1ca9c33274fd6133c`, `fix(cli): complete task 10 root-flag failure mapping and corpus`, `Refs: phase-24/task-10`. Includes earlier same-task checkpoints `d574cf4` and `6607143`. |
| Task 10 report/ending HEAD | `a6b6203d62ec6061dcac84210f6fef6df13d4ef2`, `docs(task): finalize completion report for phase-24 task-10`. Current report bytes equal the blob at this close. |
| Task 11 implementation | `e6b49797cc1dba22012758922778bed805f8f16c`, `test(cli): complete Stop corpus and plain-copy Human fixtures`, `Refs: phase-24/task-11`. Preserves `ab05e49` product changes. |
| Task 11 report/ending HEAD | `93a1f792d27f001b4052113e9056faf97706c3e5`, `docs(task): finalize completion report for phase-24 task-11`. Current report bytes equal this close blob. |

Verified all 47 distinct full commit hashes in the two reports and named close
commits with individual `git rev-parse --verify <hash>^{commit}` calls; each
returned observed exit 0. Verified reported subjects and phase/task references.
Inspected aggregate product, tests and documentation diffs from review-8 to
Task 10 close, and Task 10 close to Task 11 close. The DEC-286 continuation's
product changes are the bounded root guard and two pure helper extractions.
The DEC-287 continuation changes only the new fixture test, CLI documentation
and Task 11 report; no reconciliation, import-binding or failure-renderer
product source changes occur in that continuation. No dependency changes.

## Summary

- Review outcome: BLOCK
- Must-fix findings: 1
- Improvement findings: 0
- Blocker candidates: 1
- Existing escalation references: ESC-24-004, ESC-24-005, ESC-24-006 (resolved history only; none open)
- Covered validation assertions: VAL-027, VAL-028, VAL-029 (message allocation and retained controls), VAL-036, VAL-037, VAL-038, VAL-039, VAL-040, VAL-041, VAL-052, VAL-064, VAL-065, VAL-067; VAL-054 and VAL-055 interim judgment pending

## Validation Contract Coverage

ON TRACK denotes independently corroborated deterministic evidence for the
reviewed allocation, not a new attestation of every historical full gate.
Human-only evidence is separately NOT ON TRACK for the required D2 checkpoint.

| VAL ID | Claimed by | Status | Evidence |
|---|---|---|---|
| VAL-027 | task-10 message half | ON TRACK | I1 actual-binary plain-copy Envelope pair pins `prior registration binding does not match envelope`, metadata, exit and unchanged filesystem; O1 retained original/copy/import controls. |
| VAL-028 | task-10 message half | ON TRACK | I1 Claim pair pins `claim acquire binding is invalid` in both modes with unchanged effects/metadata; O1 original and imported qualifying controls. |
| VAL-029 | task-10 message half | ON TRACK | I1 valid-Claim Run and gate copy pairs pin `run claim binding is invalid` and `run gate: repository binding is invalid`; O1 retained original/copy/import controls. |
| VAL-036 | task-10 | ON TRACK | I1 registry completeness and real-process code pairs: 80 RegisteredCodes members, 74 ordinary pairs, five special pairs and one justified unsupported-kind exclusion. Supplemental reachable idempotency-conflict remains covered. DEC-285 six pairs and DEC-286 exact refusal pair pass; O1 grammar pairing/boundaries. |
| VAL-037 | task-10 | ON TRACK | O1 `TestHumanFailureLines` and I1 blocking denials verify Stop IDs before Attention IDs, and full-text nonblocking pairs verify absence. |
| VAL-038 | task-10 | ON TRACK | I1 every next-table code, including supplemental idempotency conflict; O1 grammar on input/non-input leaves. Full equality proves absence for more than five outside-table codes, operation-unavailable and export.import. |
| VAL-039 | task-10 | ON TRACK | O1 scalar/array/null/object handling, every override, escaped controls, exactly 200 versus 201 Unicode scalars and whole-escape cut. Error payload is not escaped or cut again. |
| VAL-040 | task-11 | ON TRACK | O1 complete ten-code/ten-Condition renderer corpus in open/resolved states across stop/attention/resume; I1 real stop.raise corpus records qualifying resolution Events, correct Next/disposition steps and resolved Event/no-Next states. |
| VAL-041 | task-11 | ON TRACK | F1 manifest's clean/dirty matrices and M1 independent 66-command replay; same-root inspect/resume and plain-copy resume positive controls; clean/dirty copied inspect exact DEC-287 refusal pairs, also directly reproduced B1. I1 retained imported-destination same-worktree control. |
| VAL-052 | tasks 10, 11 | ON TRACK | Inspected both exact-gate timing receipts below. No budget breach; Task 10 cmd/fid warning explicitly retained. Fresh focused timings are separate, not replacements for full-gate timings. |
| VAL-064 | retained task-9a; D2 recheck | ON TRACK | O1 all whole-tree grammar classes/domains, completeness, project/non-project, short/long/512/513/Unicode/escape-cut rows. Human messages and additional lines match accepted JSON-derived expectations; full offending details retained. |
| VAL-065 | retained task-9a; D2 recheck | ON TRACK | O1 unknown root/group words, near-miss isue, sorted alternatives, DEC-279 dotted operation, bounded messages and no-word help controls. |
| VAL-067 | task-10 | ON TRACK | I1 all eighteen actual fallback sites have both-mode domain rows and no exclusions; unchanged details/exit/filesystem. O1 cause-message boundaries and whitespace controls. I1 sensitive activation domain guard reaches the fallback and neither stream exposes the environment sentinel. |
| VAL-054 | Human, D2 interim | NOT ON TRACK: judgment pending | F1/M1 prove five executable denials and captured output, not the Human's ability to state why and what next. B1-H requires the five answers. Phase-end confirmation remains required. |
| VAL-055 | Human, D2 interim | NOT ON TRACK: judgment pending | F1/M1 prove executable Stop, Attention, resume and Review views in the required states. Human accept/reject and reasons for four views are absent. B1-H requires them; phase-end confirmation remains required. |

## Gate And Evidence Checks

### Exact report-to-spec comparison

Compared both reports' `## Tests run`, `## Exact quality-gate results`, typed
outcome, deviations and timing receipts against their exact task command lists
and applicable QA-plan commands. All seven Task 10 commands and all six Task 11
commands appear with final PASS, an output token and observed termination.
Neither focused selections nor the fixture generator replace a required gate.
No required Developer command is missing, silently skipped or narrowed.

Both specs require these five byte-identical commands:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

Task 10 additionally requires:

```text
rg -n "func (write[A-Za-z]*Failure|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'
rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'
```

Task 11's sixth command is:

```text
rg -n -i "worktree|dirty digest|HeadOID|head_oid" cmd/fid -g '!*_test.go'
```

| Report | Inspected final receipts; these full gates were not independently rerun | Preserved history |
|---|---|---|
| task-10-done.md | Format silent exit 0; lint `0 issues.`/0; ordinary workflow `116.440s`/0; coverage `62.6%`/0; tagged `335.908s`/0; R9/R20 match tokens and exit 0. | Prior exact ordinary/tagged FAIL attempts remain with causal deviations and later identical commands. Independently decoded the complete quoted incoming report at e0495b5: byte-identical 86,395 bytes, SHA-256 `6e3854634443638428063eda4cffa37e2c6ba42b8d5e485cabd617583986cca7`. Earlier activity/root raw failures and protected outcomes remain history. |
| task-11-done.md | Format silent exit 0; lint `0 issues.`/0; ordinary workflow `105.826s`/0; coverage `62.8%`/0; tagged `338.886s`/0; 21 accounted R11 matches/0. | Prior BLOCKED state, focused generator FAILs, interrupted-session uncertainty, report-lint failure/correction and corpus reds remain. DEC-287 correction and later identical focused commands are documented; earlier unexecuted full gates remain explicitly historical BLOCKED, not erased. |

The report checks below are fresh independent mechanical evidence, not
execution of the Developer's historical full suites. Broad race/coverage,
repository lint, benchmarks, archive-scale and full-checkptr were not repeated
for this bounded window: their report receipts were audited and representative
affected behavior reproduced. Later phase-level gates retain their allocations.

| Inspected full-gate timing | cmd/fid | internal/workflow | Slowest other freshly timed ordinary package | Tagged integration |
|---|---:|---:|---|---:|
| Task 10 | 484.734 s | 116.440 s | internal/issues 216.748 s | 335.908 s |
| Task 11 | 390.671 s | 105.826 s | internal/export 136.519 s | 338.886 s |

All recorded numeric times are below 600 s ordinary and 900 s tagged budgets.
Task 10's 484.734 s exceeds the 480 s warning and the report names it. Task 11
explicitly identifies unchanged cached packages as lacking fresh timings;
no cached time is reconstructed or promoted to a new observation.

### Fresh reproduced commands

Working directory `/home/mark/projects/fid`, except manifest commands use each
recorded `cwd`. Full nested tool results were retained. Every yielded session
was polled until its own `exit_code` appeared. No silence or wrapper completion
was treated as command termination. All fresh QA test attempts passed on their
first run; B1 exit-4 refusals are expected process behavior, not test failures.

| ID | Exact command or evidence | Result and verbatim token |
|---|---|---|
| L10 | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-10-done.md project-docs/docs/phases/phase-24/tasks/task-10.md` | PASS; `OUTCOME COMPLETE`; last token `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| L11 | `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-11-done.md project-docs/docs/phases/phase-24/tasks/task-11.md` | PASS; `OUTCOME COMPLETE`; last token `PASS report-lint completed with no failures` (terminated `exit_code: 0`). |
| G1 | `gofmt -l cmd internal` | PASS; `exit 0 (no output)` (terminated `exit_code: 0`). Representative exact task command. |
| F1 | Required fixture command below | PASS; `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	29.797s` (terminated `exit_code: 0`). |
| I1 | Focused tagged command below | PASS; `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	46.384s` (terminated `exit_code: 0`). |
| O1 | Focused ordinary command below | PASS; `ok  	github.com/mbell-nzl/fid/internal/cli	1.016s`; final token `ok  	github.com/mbell-nzl/fid/cmd/fid	77.871s` (terminated `exit_code: 0`). |
| O2 | Additional affected retained/root tests below | PASS; `ok  	github.com/mbell-nzl/fid/cmd/fid	8.231s` (terminated `exit_code: 0`). |
| M1 | Independent manifest/hash/66-command replay and historical-report preservation audit | PASS; `PASS manifest: 33 rows, binary digest verified, 66 commands independently replayed with observed statuses and complete streams`; final token `PASS Task 10 incoming report preserved byte-for-byte: 86395 bytes 6e3854634443638428063eda4cffa37e2c6ba42b8d5e485cabd617583986cca7` (terminated `exit_code: 0`). Each process status and both streams retained separately. |
| Inventory | Exact section-1 Bash walk and R1-R20 commands | PASS; walk ends ` 271 total` (terminated `exit_code: 0`). Each search's actual status/token is retained; R5 alone has expected no-match exit 1. |

F1, byte-identical QA-plan regeneration command, executed at the Task 11 close:

```text
FID_P24_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-24-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration
```

I1:

```text
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid/integration -run '^(TestP24(UnavailableImportRootFlagFailurePair|RegisteredCodePairingCompleteness|RegisteredCodePairs|RetainedReadFailurePairs|GenericFallbackDomainPairs|GenericFallbackCorpusCompleteness|GenericFallbackSensitiveProcess|PlainCopyFailurePairs|HumanProtectedDecisionFailure|HumanGateTerminalFailure|StopHumanRealCorpus)|TestP23StopFailureHumanAndJSONParity|TestP23HumanFixtures|TestP23HumanFixtureStepsRemainRunnable|TestP22HumanGuidance)$'
```

O1:

```text
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/cli ./cmd/fid -run '^(TestHumanFailure.*|TestP24(FamilyFailureRenderer|GenericFallbackMessage|HumanFailureCallerAudit|RootFailureAuditNegativeControl|StopHuman.*|RepositoryHumanViews|RunStartOriginalCopyAndImport|GateOriginalCopyAndImport|EnvelopeAndClaimOriginalCopyAndImport)|TestGrammar.*|TestLiveLeafCommandsMatchContractInventories|TestUnavailableImportRootFlag.*|TestRoot.*|TestStatsTokensDefaultsToJSONByActorAttribution)$'
```

O2 covers the newly added root tests and specific retained assertions outside
O1's selector:

```text
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./cmd/fid -run '^(TestP24(UnavailableRootFlagMapping|DiagnosticRootFlagBytesPreserved)|TestStopSuccessHumanNamesBlockingReasonAndQualifyingResolution|TestPhase23HumanLeafTextMatchesRecordedFacts|TestRunStartBindingMatrixRejectsWithoutMutation|TestIssueTransitionRejectsInvalidEdgeAndLeavesStateUnchanged)$'
```

### Preserved fixture manifest and actual-binary observations

Manifest:
`/home/mark/projects/fid/.cache/phase-24-human/fixtures-2268063598/manifest.json`.
Manifest SHA-256:
`ccf9d94d7aab59422bce4be30212a068b56cde3943150e7c33958450f8574beb`.
Binary: the sibling `fid`; independently computed SHA-256:
`61fa7a4c17e05de5e662caf1147349460cb93c2a1ca482cdb6e0c23ba29f55b7`.
This is a new close-commit build, not the Developer's retained binary.

The manifest preserves workspace, cwd, both argv/command forms, expected exit,
binary digest and all four captured streams for every row. Independently
checked all 33 unique names, five denial names and 28 state/view names.
Independently replayed all 66 Human/JSON commands with captured real statuses.
Human output remains byte-identical to the manifest; refusal JSON is exact.
Successful resume JSON includes a live lease projection, so that JSON is not
treated as an immutable byte oracle. Forbidden identity/digest/nil/raw-JSON
tokens are absent from successful captured Human views.

| Repository row, both clean and dirty | Observed Human result | Process termination, each mode |
|---|---|---|
| same-root run.inspect | `Worktree comparison: same worktree`; short recorded/current HEAD and state | `exit_code: 0` |
| same-root resume | `Worktree comparison: same worktree` | `exit_code: 0` |
| plain-copy resume | `Worktree comparison: different worktree` | `exit_code: 0` |
| plain-copy run.inspect | Exact DEC-287 recovery-state-conflict, empty stdout, no repository view, no blocked-by/next line | `exit_code: 4` |

B1 directly reproduced all four copied-inspect calls through exec_command,
using the manifest binary, `--project <fixture-root>/<state>-same`,
`run inspect <run-id> --repository <fixture-root>/<state>-different`, with
and without root `--json`. Cwd was each `<state>-same` directory. Clean Run
`RUN-3a5f64`, dirty Run `RUN-136ab8`; all four nested results terminated exit 4.
Human final tokens were respectively:

```text
reason: recovery-state-conflict; observation_digest=null, paths=, reasons=repository identity changed, run_id=RUN-3a5f64
reason: recovery-state-conflict; observation_digest=null, paths=, reasons=repository identity changed, run_id=RUN-136ab8
```

Each JSON final token was the complete DEC-287 clause 2 envelope, with that
Run's subject and ID, `reasons:["repository identity changed"]`, empty paths,
null observation digest, category conflict and operation run.inspect. Full
JSON bytes are retained in the manifest and independent replay receipts.
The tagged test additionally pins no appended Event for the two-mode refusal.

Fresh auxiliary evidence is retained at
`/home/mark/projects/fid/.cache/review-9-d2-zjojhbnw/`: complete inventory
stdout/stderr and nested receipts in `inventory-receipts.json`, full commit
verification in `commit-verification.json`, product/test/doc diffs and all
66 process results in `manifest-replay-receipts.json`. The recursive walk's
binary and six labelled lists remain at `/tmp/tmp.69E6PWrDKH/`.

## Per-Task Findings

### Task 10 - Human failure lines in every family

#### Must-fix

- None in Task 10 implementation/report scope. The shared D2 Human checkpoint
  is B1-H under Task 11 below; it is counted once.

#### Improvement

- None.

Observation: `internal/cli/human_failure.go:13` centralizes presentation from
the canonical Failure without moving authority or repository business rules
into the CLI. `canonical_request.go:191` supplies live leaf/input context.
The caller audit and its root negative controls reproduce; diagnostics retain
their exception. DEC-285 activity/stats/AC callers adopt live JSON mode and the
existing read mapper; DEC-286 root guard preserves precedence and exact
alternatives. I1 checks real process streams, metadata, exit and filesystem
in both modes rather than merely calling the new renderer. O1 proves grammar
and Human pairing after adoption, including control/Unicode cuts. These match
`05_DESIGN_PLAN.md` sections 3, 5 and 6. The RegisteredCodes omission of the
reachable idempotency-conflict is disclosed, with supplemental process
coverage; no new exclusion or registry modification was invented here.

Retained convergence follows the accepted representation/message change:
twenty existing test files change, with no deleted Test functions or new
skips. Inspected grammar, Run binding, protected Decision, link and removed
command assertions: exact Human expectations change while domain subjects,
typed JSON/exit and no-mutation controls remain. Independently verified the
reported introducing hashes, including task-9a `2d5e9d55`, Task 6
`1baaecc7` and Phase 21 Run `8c763e9c`, and their commit references. The
special gate/budget pairs preserve the accepted terminal-frame transport;
they do not invent a new streaming exception to get a green ordinary pair.

### Task 11 - Readable Stop, Attention and repository lines

#### Must-fix

- **B1-H — Human checkpoint pending.**
  [qa-plan.md:531](qa-plan.md#8-mid-phase-review-cadence) and
  [validation-contract.md:122](validation-contract.md) require the interim
  VAL-054/VAL-055 Human walkthrough after Task 11. The required
  `project-docs/docs/phases/phase-24/qa/human-walkthrough-after-task-11.md`
  is absent, and no Human answers/acceptance were supplied. Owning checkpoint:
  Window D2, tasks 10/11; owning authority: Human, with Runner/QA recording
  coordination. **No Developer continuation is assigned.** Repair boundary:
  obtain and record only the actual Human judgments using F1's preserved
  fixtures; no product, report, task, spec, QA-contract or gate changes.
  Required evidence: five reason/next-action answers for VAL-054 and four
  accept/reject judgments with reasons for VAL-055, bound to the regenerated
  manifest/binary/close commit. Independent follow-up scope: validate those
  nine answers and provenance, preserve this review and all deterministic
  evidence, and record the checkpoint's resulting review outcome. Do not
  mark the finding resolved from deterministic output alone.

#### Improvement

- None.

Observation: `stop_human.go:52` renders labelled closed-variant details and
switches open Next guidance to resolved fact/resolution Events; Attention
loads linked Stops without changing lifecycle selection. O1 covers all ten
variants in both states and I1 proves real qualifying Events. Real resume
keeps its canonical projection: disposed items can leave it, while an open
Attention can retain its resolved Stop. That is explicitly documented/tested,
not a forced change to JSON to make an unreachable history row render.
`run_inspect_human.go:86` and `p23_human.go:139` shorten repository presentation;
SameWorktree remains a domain predicate. DEC-287's two exact refused rows retain
their names, clean/dirty states and plain-copy identity. The successful plain
copy is proven on resume, separately from the retained imported-destination
same-worktree control. This matches `05_DESIGN_PLAN.md` sections 2, 3 and 6
and the delivered CLI examples.

Four existing Task 11 test files converge. Their reported provenance resolves:
Phase 23 review facts `59d2653d`/`b71ed1ba` (phase-23/task-14), imported
inspection `eff28fb1` (phase-22/task-19), and Attention `a0a93c7c`
(phase-22/task-9). Inspected changed assertions retain exact Human/JSON fact
and read-only checks. Synthetic different-worktree inspect renderer coverage
is not used as proof of CLI plain-copy success.

## Prior Finding Resolution

This is the initial D2 review, not a fresh re-review of D1. Review-7's F1
Activation identity scalar finding remains closed by independent review-8
PASS, and Window C remains review-6 PASS. No earlier artifact was amended or
superseded by a self-correction. Task 10/11 recovery histories and their
resolved DEC-285/286/287 authority boundaries remain intact.

## Repair Plan

No bounded Developer repair is required. B1-H requires the external Human
checkpoint below. It must not be routed to Task 10 or Task 11 by the filename
of the missing walkthrough, or converted into product changes without a
separate acceptance finding.

## Cross-cutting observations

R1-R20 and the exact recursive Bash walk were reproduced at current HEAD.
Markdown table alternations were interpreted as actual regex alternation.
Each nested search/list status and real last output token is saved, including
the unchanged no-match R5 expected/observed exit 1. Counts versus review-8:

| Registry | Current evidence and accounted delta |
|---|---|
| R1/R2 | 45 request-digest sites in 33 files; one shared helper. Unchanged. |
| R3/R4 | 126 storage matches in 18 files; 38 migration files through 0019; 17 projection/schema references. Unchanged. |
| R5 | No closed decoder functions; real `exit 1 (no output)` (expected and terminated `exit_code: 1`). |
| R6/R7 | Four repository helper matches; 16 resolved-fact references in two files. No caller-owned identity revival. |
| R8 | 85 definition/declaration matches. Registered corpus proves the actual 80-member code list plus declared exclusion/supplement. |
| R9/R10 | 43 named failure functions; 61 Human assertion matches in 26 files, increased from D1's 30 by Task 10's renderer/corpora/root tests and Task 11 fixture assertions. Added stable callers and root hook are separately audited. |
| R11/R12 | 21 repository references versus 19 at D1, accounting for comparison/helper adoption. All three Stop/Attention renderer files exist; no remaining raw Human repository identity line. |
| R13/R14 | 15 existing judgment/PRIVATE references, unchanged; their removal remains Task 13. Live leaf completeness test exists and passes O1. |
| R15/R16 | 92 leaves, zero lacking Short; 50 idempotency-key leaves, 49 canonical keyed leaves after named export.import exclusion. |
| R17/R18 | 29 root/groups; 61 typed flags, 17 integers and 44 booleans including 39 example flags. Unchanged from D1. |
| R19 | Exactly four permitted request-level family fallbacks, all exercised by I1. |
| R20 | 23 textual matches: eighteen actual generic fallbacks and five typed/combination sites. AST completeness and both-mode domain rows pass I1; no fallback is hidden behind an unreachability claim. |

Representative exact task inventory receipts: R9's last token
`cmd/fid/issue_get_list.go:98:func renderIssueReadError(cmd *cobra.Command, jsonMode bool, operation string, readErr error) error {`;
R20's last token is the `cmd/fid/issue_mutation_response.go:218` genericFailure
line; R11's last token is `cmd/fid/gate.go:300` repository-mismatch mapper
line. Each nested command terminated `exit_code: 0`; complete verbatim lines
are saved in the inventory receipt file. Search traversal order need not match
the Developer's last line; the full sets/counts and observed status match the
accepted inventory.

## Blocker Candidates

- **B1-H — protected Human product judgment / external action required.**
  Evidence: QA-plan section 8 says, “The window-D2 walk-through records Human
  judgment only.” VAL-054 and VAL-055 assign the validator to Human and the
  interim evidence to this checkpoint. The walkthrough file and actual
  Human judgments are absent. F1's 33 rows and M1's 66 passing process
  replays cannot answer whether a Human understands or accepts the output.
  This exceeds Developer repair and QA Reviewer judgment authority. The
  invocation permits a separate deterministic walkthrough artifact only
  when QA is authorized to record it; the plan permits Human judgment
  only, and this role writes exactly one review. Accordingly no separate
  walkthrough file was written, and all judgment stays pending. Required
  action: Human completes the nine rows below from the retained Human
  streams/commands; Runner/authorized QA records their actual answers at
  the named walkthrough path and obtains independent checkpoint closure.
  No ESC is allocated and no existing open ESC covers this pending action.

### Exact pending Human action

Use the manifest at
`.cache/phase-24-human/fixtures-2268063598/manifest.json`, generated at
`93a1f792d27f001b4052113e9056faf97706c3e5`. Its `human_stderr` and
`human_stdout` fields preserve the actual output, and `human_command`/`cwd`
allow execution. Record the Human's answers verbatim or explicitly attributed
without using JSON, source or specs to supply the answers:

| Interim row | Manifest fixture(s) | Required Human evidence | Current status |
|---|---|---|---|
| HUMAN-02a / VAL-054 | idempotency-conflict | State why refused and next action from Human output alone. | PENDING |
| HUMAN-02b / VAL-054 | repository-mismatch | State why refused and next action. | PENDING |
| HUMAN-02c / VAL-054 | invalid-argument, nested input.expected_repository.dirty | State why refused and next action. | PENDING |
| HUMAN-02d / VAL-054 | protected-decision-required, including blocked by | State why refused and next action. | PENDING |
| HUMAN-02e / VAL-054 | source-digest-mismatch | State why refused and next action. | PENDING |
| HUMAN-03a / VAL-055 | clean/dirty × same/different stop.list | Accept/reject the Stop view, with reason. | PENDING |
| HUMAN-03b / VAL-055 | clean/dirty × same/different attention.list | Accept/reject the Attention view, with reason. | PENDING |
| HUMAN-03c / VAL-055 | clean/dirty × same/different resume | Accept/reject resume, including short HEAD/state and worktree comparison, with reason. | PENDING |
| HUMAN-03d / VAL-055 | clean/dirty × same/different review.get | Accept/reject the Review view and repository lines, with reason. | PENDING |

The additional clean/dirty plain-copy run.inspect rows are preserved refusals,
not successful different-worktree views. These interim answers do not replace
the phase-end HUMAN-02/HUMAN-03 confirmation. If the Human rejects a row,
classify the concrete acceptance issue under its existing owner/authority;
do not infer a design change from the mere absence of today's judgment.

## Existing Escalation References

- ESC-24-004: resolved by DEC-285; activity/read callers independently verified.
- ESC-24-005: resolved by DEC-286; root-flag refusal independently verified.
- ESC-24-006: resolved by DEC-287; exact clean/dirty copy refusal independently verified.
- No open escalation; no escalation identity created or register edited.

## Final Lines

Review written: project-docs/docs/phases/phase-24/qa/review-9.md
Must-fix: 1  Improvement: 0  Blocker candidates: 1
VERDICT: BLOCK
