# Phase 24 QA Plan - Operating Experience And Input Convergence

> Phase: 24 - Operating experience and input convergence
> Refreshed: 2026-10-04 (DEC-280/281 recheck, third pass)
> QA Lead: Codex. Reviewed HEAD: `c9acdac81972de5f8388cd346f7b78787dbf113c`.
> Operating mode: **Heavy** (DEC-277). Earlier reviews by Claude Opus 5.5
> remain in sections 6 and 12.1-12.5; their command results are historical,
> not freshly attested gate results in this session. Section 12.6 records
> this recheck's own evidence and supersedes earlier readiness verdicts.
>
> **Status: NOT READY (scoped).** DEC-280/281 disposition all nine RR and
> thirteen RI findings. Two new Must-fix contract candidates remain:
> **RC-01**, the pinned plain-copy `run.start` cause (VAL-029; tasks 7, 10),
> and **RC-02**, unbounded grammar messages (VAL-064/065; task 9a).
> A fresh binary returns `run claim binding is invalid`, while DEC-281 pins
> `run envelope binding is invalid`. The Tech Lead must reconcile the
> contract and specify bounded grammar-message precedence before this plan
> can pass. The phase brief's dispatch hold stays
> in force. Review window D is split into D1 after task 9a and D2 after
> task 11 under DEC-280 item 10; the Human walk-through stays after task 11.

## 1. Scope And Source Documents

Phase 24 gives a Human control of a stuck or diverging journey through fid
output alone, and it converges inputs:
- fid observes repository identity itself, so callers send only the
  expectations a check uses;
- request digests cover the input as written, and per-field digests name
  the differing field (migration 0019);
- one input descriptor per `--input` leaf drives help, an `--example`
  template and field-path validation;
- grammar failures name the offending spelling, and unknown commands return
  `operation-unavailable` (Phase 21 allocations, DEC-278);
- Human failure lines repeat the JSON failure with reason, blockers and
  next action;
- Stop, Attention and repository lines become readable;
- three judgment `get` leaves are added;
- judgment views show recorded reasoning, with one disclaimer per section.

The phase has 15 sequential tasks (including 9a) and 67 assertions (62 deterministic, 5
Human).

Sources read for the first two passes (third-pass sources and evidence are
listed in section 12.6):

- `AGENTS.md`, `CLAUDE.md`
- `project-docs/docs/current/phase_brief.md`, `decisions_active.md`,
  `progress_current.md`, `tl_escalations_current.md` (no open ESC)
- `project-docs/specs/08_QUALITY_PLAN.md` sections 5-8
- `project-docs/specs/09_AI_WORKFLOW.md` sections 2.6, 4.3, 4.4, 6.1
- `project-docs/specs/01_PRD.md` (the PRD IDs the contract cites)
- `project-docs/specs/GLOSSARY.md` (Project-owned prose entry)
- `project-docs/docs/phases/phase-24/qa/validation-contract.md` (66 rows)
- `project-docs/docs/phases/phase-24/public-interface.md` (with the DEC-278
  clarifications to sections 2, 5, 6, 9, 10 and D5 in section 13),
  `input-inventory.md` (R1-R16), `decisions.md` (DEC-273..279),
  `triage-log.md`
- `project-docs/docs/phases/phase-24/tasks/task-1.md` .. `task-14.md` (as
  amended by `c80304c` and `21c833f`)
- `project-docs/docs/phases/phase-23/decisions.md` DEC-227, DEC-259;
  `phase-23/qa/qa-plan.md` and `qa/human-fixtures.md` (format and fixture
  precedent only)
- `project-docs/docs/phases/phase-21/public-interface.md` section 1 (lines
  14-45), section 2 (lines 75-95), section 4 (lines 400-415) and the code
  table (line 459)
- `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` (CF-046)
- `project-docs/templates/qa-plan-template.md`
- `project-docs/skills/validation-contract/SKILL.md`,
  `project-docs/skills/qa-evidence/SKILL.md`
- Read-only source evidence, cited inline. This includes the module-cache
  copy of `github.com/spf13/cobra@v1.10.2` (`args.go`, `command.go`).

No held-out scenario, evaluator or excluded-scenario material was read. Mode B
scenario execution is **N/A**, as in Phases 21-23. `AGENTS.md` records
end-to-end and scenario validation as `N/A` until a harness is recorded, and
none is. The `(human)` checks in section 9 remain required. Their "take
control" scenario is the phase's Human run.

## 2. Validation Contract Refinement Check

The following bullets are the historical second-pass check at `21c833f`.
Current refinement results and proposals are in section 12.6. The first-pass
record follows in 2.1-2.4.

- All in-scope requirements referenced by at least one `VAL-NNN`: **YES.**
  Every PRD ID in the Coverage table exists in `01_PRD.md`. The table does
  not list the roadmap-sourced rows VAL-035 and VAL-051..053 (RI-07,
  improvement).
- Every `VAL-NNN` has one validator role: **YES.** QA Reviewer 58, QA Lead 3
  (VAL-005, VAL-048, VAL-051), Human 5.
- Every `VAL-NNN` has a concrete evidence method: **YES**, with one
  infeasible fixture example in task 4 for VAL-023 (RR-09).
- Every `VAL-NNN` maps to a task/window or phase-end check: **YES**, verified
  mechanically. All 66 IDs are cited by at least one task spec. The only
  citation outside the contract is the Phase 23 restatement note (VAL-073).
- Every `VAL-NNN` describes one observable: **YES.** The six bundled rows
  are split (VAL-058..063). VAL-064 and VAL-065 are matrices over enumerated
  cases (failure class × leaf × mode, and command group × mode), and are
  accepted as such.
- Conditional or binary assertions require both states: **YES** for
  VAL-058..063 and VAL-066 (the new-key control is the non-qualifying
  state). **NO** for VAL-038 as task 10 implements it: task 10 and contract
  section 9 give different expected outputs for a grammar failure on an
  `--input` leaf (RR-01).
- Parity or "matches reference" assertions have deterministic proof:
  **YES.** VAL-002, 007, 035, 042 and 050 have pinned comparison inputs.
- Assertions consistent with their task specs: **NO for VAL-036** (stdout
  for post-frame `gate run` failures, RR-07) and **VAL-064** (the
  `journal inspect`/`sql` exclusion, RR-08).
- Revisions requested from Tech Lead: RR-01..RR-09 Must-fix and
  RI-01..RI-13 Improvement (section 12.3).

### 2.1 Contract-Coverage Re-Derivation (first pass, at `304aa45`)

Contract sections 2-12 were mapped to assertions independently of the
contract's Coverage table. At the first pass, the following clauses had no
assertion that evidenced them, or one that could not. Section 12.1 verifies
the dispositions.

| Clause | Mapped to | Gap |
|---|---|---|
| Section 2: "A replay never observes" | VAL-009 | Byte-identical replay after HEAD moves cannot show whether the replay observed (R-03) |
| Section 2: unobservable causes (not a worktree, unmerged, tracked flags, changed during capture) | VAL-023 | Only two of five reasons are exercised. Measured causes such as an unborn HEAD have no reason member (R-05) |
| Section 6 `repository-unavailable` producers "every observing mutation" | VAL-023 | `evidence.evaluate` is in the 17 changed leaves but is an unkeyed, Event-less read (R-04) |
| Section 7: "Every keyed mutation also stores one digest per top-level public input name" | VAL-006, VAL-002 | Evidence samples one leaf per family; a missed call site among 43 is undetectable (R-02) |
| Section 9: "Every Human-mode failure on stderr is exactly these lines" | VAL-036 | Unknown commands print nothing in either mode, and grammar failures lose the offending spelling (R-07) |
| Section 9 next-table row `invalid-argument` "on an `--input` leaf" | VAL-038 | Non-`--input` leaf state not asserted (R-08) |
| Section 9: 200-scalar cut | VAL-039 | Boundary and counting rule ambiguous (R-09) |
| Section 10 A4 `Next:` keyed on Condition type | VAL-040 | One Stop per code cannot be relied on to reach all 10 Condition types (R-10) |
| Section 10: "same worktree / different worktree" | VAL-041 | Only clean/dirty fixtures named (R-11) |
| Roadmap journey "from fid output alone" | VAL-049, VAL-053, VAL-057 | Caller-chosen IDs (DEC-272) and `expected_repository` (contract section 3, DEC-276) never come from fid output (R-13) |

### 2.2 Clauses Verified Sound (first pass)

- **Accepted 0019 DDL applied, not read** (section 6.4). On a real 0001-0018
  chain, `ALTER TABLE ... ADD COLUMN` with the exact accepted text succeeds
  on an empty table. `DROP COLUMN` restores the schema-18 `sqlite_master`
  dump byte for byte (VAL-001 and VAL-058 are feasible). The CHECK and NOT
  NULL constraints behave as accepted. A populated table refuses the column
  (I-14).
- **Inventory registries** R1, R2, R4-R12 and R14 reproduced their recorded
  counts at the planning base (section 6.3). R3, R13 and R15 differed; all
  three are corrected (section 12.1).
- **The `--input` leaf set** is exactly the 39 leaves of inventory section 1.
- **First failing tests** for tasks 3, 6, 12 and 13 are red today (section
  6.5). The other first failing tests name new files or behavior that does
  not exist yet.
- **The 100k rebuild** has about 19.6% headroom (worst of three fresh samples
  24128 ms).
- **Documentation targets** named by tasks 1-14 all exist (section 6.8).

### 2.3 Revisions Requested From Tech Lead (first pass)

This table is the first-pass record. DEC-278 dispositioned every item, and
section 12.1 verifies each disposition against source. New findings are in
section 12.3.

| ID | Sev | Where | Finding and proposal |
|---|---|---|---|
| R-01 | M | task 1 `Files To Create Or Modify`; inventory R3; VAL-001-004 | **The new column reaches closed registries that task 1 does not list (DEC-167/173).** The R3 command (`rg -n "idempotency_results" migrations internal -g '!*_test.go'`) returns 16 files, not the three the inventory summarizes. Four are closed registries the column must enter: (a) `internal/journal/current_event_columns.go:28`, the per-table column and type list that "the contract test checks against DDL"; (b) `internal/journal/current_event_json.go:96`, where any `_json` column that is not listed fails `unregistered current JSON column`; (c) `internal/journal/current_event_updates.go:45-46`, the closed values/guards of the gate-receipt `update` (values `result_entity_id result_event_id aggregate_sequence result_json`, guard `request_digest`); (d) `internal/export/archive_updates.go:24` `semanticUpdateColumns`, where archive validation rejects any other update value as "immutable in this Event family". Task 1 says the map "travels ... on the gate-receipt `update` path". If it travels as an update value, (c) and (d) must change, or VAL-004 and VAL-050 fail on any archive that contains a Receipt. Retained tests displaced by the NOT NULL column are also outside the "pinning schema 18" phrase: `internal/maintenance/check_test.go:347`, `internal/db/operational_kernel_migration_test.go:583-595` and the closed portable column list `cmd/fid/integration/p22_portable_schema_test.go:34`. Proposal: add (a)-(d) and these tests to task 1. State whether the update effect carries the map as a value or as a guard. Correct the R3 count. |
| R-02 | M | VAL-006, VAL-007, VAL-002; task 2 | **Conflict corpus samples one leaf per family, but adoption is per call site.** There are 43 `RequestDigest:` sites in 32 files (R1), serving 49 keyed leaves: the live tree has 50 leaves with `--idempotency-key`, and `export import` uses an install marker and no request digest (`internal/export/import.go:104, 175`). A site missed by task 2 keeps task 1's `{}` map. A conflict there would carry an empty `differing_fields`, which the contract forbids, and no family-sampled test would catch it. The contract's own Shared Terms and the brief's Known Risks require corpora "enumerated from registries, not hand-picked". Proposal: VAL-006/007 enumerate every keyed leaf from the live tree (the `--idempotency-key` leaf set minus `export import`), with one one-field-change row and one canonical-retry row per leaf. Keep the two-field and `--actor` rows per family. A completeness test fails when a keyed leaf has no row. VAL-002's rebuild journal contains one keyed mutation per keyed leaf. |
| R-03 | M | VAL-009; tasks 2, 4-8 | **"No new observation" has no observable.** A replay that re-observes an advanced HEAD and then discards the result still returns byte-identical output. The listed method cannot tell the two apart, so the contract section 2 rule "A replay never observes" is unproven. Proposal: add one retry per observing keyed leaf, made after the worktree has become unobservable (unmerged paths, or for `run.close`/`run.abandon` the same state that yields `recovery-state-conflict`). Also add one retry per SourceInput leaf after the referenced file has been deleted. Each retry must replay the original result with no Event. |
| R-04 | M | Shared Terms "Observing mutation"; VAL-009, VAL-022, VAL-023; task 8 | **`evidence.evaluate` is counted as an observing mutation but is neither keyed nor Event-producing.** Measured: `fid evidence evaluate --help` has `--input` and no `--idempotency-key` (section 6.5), and it is the only `--input` leaf without a key. VAL-009 (replay) and VAL-022's Event-repository clause cannot apply to it. Contract section 6 names "every observing mutation" as the `repository-unavailable` producers, so it is unclear whether evaluate on an unobservable worktree returns `repository-unavailable` or something else. Proposal: split the term into "observing leaf" and "observing mutation". State per assertion whether `evidence.evaluate` is in or out. State its code for an unobservable worktree from accepted text. Blocker candidate if the code is not determinable (2.4). |
| R-05 | M | contract section 6 `repository-unavailable.reason`; task 4; VAL-023 | **The closed reason domain does not cover measured causes.** `internal/project/repository.go` raises `ErrUnsafeRepository` for about 30 causes (section 6.5). Task 4 requires classification by type into exactly five reasons. Causes with no member include: unborn HEAD (`missing base` at line 77, `missing tree` at line 81), path-count and output-size limits (lines 161, 318, 345), and Git command failures (lines 141-149, 256, 295). An unborn HEAD is reachable on a new project: `git init` followed by `fid init` succeeds (measured). VAL-023 also exercises only `not-a-repository` and `unmerged-paths`. Proposal: state the mapping from each cause class to a reason. Add unit rows for `tracked-flags` (`git update-index --assume-unchanged`), `changed-during-capture` (injected) and `unsafe-path`. Blocker candidate if a new reason member is needed (DEC-195; 2.4). |
| R-06 | M | VAL-036; task 10 | **No registry of "failure codes with a CLI producer" exists.** `internal/cli/codes.go` registers 78 codes (`registeredCodes`) and no producers. A corpus "enumerated from every failure code with a CLI producer" is therefore filtered by the test author, and can silently drop hard codes. Proposal: enumerate rows from `cli.RegisteredCodes()`. Pin an explicit exclusion list in the test, each entry with its reason (for example, codes with no reachable CLI producer). A new code without a row or exclusion fails. Count grammar-path producers of `invalid-argument` (`cmd/fid/canonical_request.go:41`) as separate rows from the `--input` structural path. |
| R-07 | M | contract section 9; Phase 21 `public-interface.md:84-87, 406, 459`; task 10; VAL-036, VAL-054 | **Unknown commands are silent, and grammar failures would lose their specific text.** Measured (section 6.5): `fid bogus`, `fid concern bogus` and `fid concern get X` exit 2 with **no output in either mode** (root `SilenceErrors`, `cmd/fid/root.go:67-68`; `main.go` exits through `cli.ExitCode`). Phase 21 accepted "Cobra unknown-command failure, exit 2, Human stderr; JSON callers receive `operation-unavailable`". `operation-unavailable` is in the Phase 21 code table but absent from `internal/cli/codes.go`; it exists only as a p21 testdata inventory file. Grammar failures (unknown flag, missing positional) emit JSON `{argument:"request", value:null, allowed:["accepted <op> input"]}` with message `invalid command request`, while Human prints Cobra text (`error: unknown flag: --bogus`). Applying section 9 verbatim makes the Human line `error: invalid command request` with `argument=request`, which loses the offending spelling. That is a Human regression, and it contradicts Phase 21's rule that grammar failures carry the offending spelling/value and accepted alternatives (the file was introduced in `62a3e46`, 2026-08-15). Proposal: under DEC-175/194, assign the conformance repair to task 10 (or task 3), adding `cmd/fid/root.go`, `cmd/fid/main.go`, `cmd/fid/canonical_request.go` and `internal/cli/codes.go` to its scope. Add grammar and unknown-command rows to VAL-036. Alternatively, record an explicit section 9 exception by DEC. Blocker candidate for the JSON `operation` value of an unknown command (2.4). |
| R-08 | M | VAL-038; task 10 | **Both states of the `invalid-argument` next-row.** The row applies "on an `--input` leaf". VAL-038's non-qualifying state uses only codes outside the table. Proposal: add `invalid-argument` from a non-`--input` leaf (for example, `fid issue create` with a missing flag) and require no `next:` line. |
| R-09 | M | contract section 9; VAL-039; task 10 | **The 200-scalar cut is ambiguous at the boundary.** "The line is cut to 200 Unicode scalars with `…`" does not say whether the count includes the `reason: ` label, or whether `…` is inside the 200. VAL-045 states its rule exactly ("81 cut to 80 plus `…`"). VAL-039 says "over 200", while task 10's AC says "a 201-scalar detail". Proposal: state the counting rule. Require evidence of 200 uncut and 201 cut. |
| R-10 | M | VAL-040; task 11 | **`Next:` is keyed on the Condition type, but the corpus is keyed on the Stop code.** There are 10 Stop codes and 10 Condition types (`internal/workflow/stop.go:175-189`). One Stop per code need not reach `work-ended`, `review-recorded` or `finding-resolved`. `stop raise` accepts any code with any Condition plus its binding (`validateStopRaise`, lines 440-500), so all ten are reachable. Proposal: the corpus covers all 10 Condition types, including the no-command `capability-available` case. It covers both statuses (open, and resolved with its Event). It also covers an open Stop without lifecycle disposition (names `fid stop resolve` then `fid attention resolve`) versus one with disposition. |
| R-11 | M | VAL-041; task 11 | **Both states of the recorded-versus-current comparison.** The pass criterion names `same worktree` and `different worktree`, but the method lists only clean/dirty fixtures. Proposal: add a different-worktree fixture (for example, `fid resume --repository <copy>` or `fid run inspect --repository <copy>`) and a same-worktree fixture. |
| R-12 | M | VAL-027, VAL-028, VAL-029; tasks 6-8 | **The refusal state names no code.** Today a plain-copy refusal falls through to generic `invalid-argument` (`argument:"input"`; message `invalid envelope input` at `cmd/fid/envelope.go:253`, `invalid claim request` at `cmd/fid/claim.go:234`). A refusal for another reason (authority, `repository-unavailable`) would also satisfy "refused". Proposal: name the exact code, message and details per leaf for the plain-copy state, and say whether task 3's field-path rule changes `argument` for these domain failures. Prove the copy itself is observable, so the refusal comes from the binding. |
| R-13 | M | VAL-049, VAL-053, VAL-057; task 14 | **The provenance rule cannot be satisfied as written.** Record IDs are caller-chosen (DEC-272). Reasoning text and project-file `path`/`display_label` are authored by the caller. `expected_repository` comes from `git rev-parse HEAD` and `git status --porcelain` (contract section 3), "typed by the caller, never rubber-stamped" (DEC-276). After task 11, Human repository lines show only a 12-hex HEAD, while `head_oid` needs the full OID. None of these values is "present in prior fid output". Proposal: enumerate the permitted provenance classes: prior fid output (state whether JSON counts for the deterministic test), `--help`/`--example`, caller-authored IDs, text and files, and the git commands that fid help names for `ExpectedRepository`. Require the test to record the class of each value. Require help for `expected_repository` to say how to obtain each value, so VAL-057's "guided only by fid output" can pass for `run close`. |
| R-14 | M | VAL-053..057; task 14 note "The QA Lead prepares those fixtures"; task-11 window | **Human fixtures have no reproducible owner.** The QA Lead has no test-code write authority to build a generator. A hand-scripted, activated protected-Decision workspace is not reproducible evidence. The Phase 23 precedent (R-10, DEC-228) is a tagged generator emitting workspaces plus a manifest (`TestP23HumanFixtures`, `FID_P23_HUMAN_FIXTURE_ROOT`, `cmd/fid/integration/p23_human_fixtures_test.go:57-63`, still present). Proposal: task 14 adds a fixture-emitting mode (env-var root, manifest, binary SHA-256) covering the VAL-049 start state, at least five denials (VAL-054), the four views (VAL-055), the judgment views (VAL-056) and authoring workspaces (VAL-057). For the task-11 walk-through, name who regenerates the converged `TestP23HumanFixtures` at the task-11 close commit (proposed: the window QA Reviewer) and the artifact path (section 8). |
| R-15 | M | task 13 (and task 14) `Files To Create Or Modify`; inventory R13; VAL-047, VAL-051 | **A declared expected-exit-1 command cannot reach exit 1 inside task scope.** `rg -n "PRIVATE_" cmd/fid/integration` also matches `cmd/fid/integration/testdata/p23_conformance/operations.json:67, 114, 144, 145`. That file is read by `p23_conformance_test.go:142, 334` and was added in `3335149` (Phase 23). Neither file is in task 13's list, and the inventory's R13 row names only two lines. Proposal: add both files to task 13 (DEC-167/173), and correct R13. |
| R-16 | M | task 9 `Exact Commands To Run` | **Undeclared exit status.** If task 9 removes the seven family-local decoders, as contract section 8's single source implies, `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` exits 1. That is an undeclared non-zero, which `AGENTS.md` treats as a failed gate. Proposal: declare the intended end state, either exit 0 with named retained helpers or exit 1 under `## Expected Non-Zero Commands`. |
| R-17 | M | VAL-001, 010, 013, 025, 031, 047 | **Bundled assertions.** Each row asserts two observables that are proven by different evidence and can fail independently. Proposal: split each into two rows (the Tech Lead assigns the new IDs), or reword it to one observable. Splits: VAL-001 into (a) up adds exactly the one column and (b) down restores the schema-18 dump byte-equal. VAL-010 into (a) every live leaf has a one-line `Short` and (b) `attention resolve --help` describes Attention resolution. VAL-013 into (a) template form (one object, required fields only, documented order, 2-space indent) and (b) no side effects (works outside a project, exit 0, empty stderr, tree unchanged). VAL-025 into (a) the closed details shape and (b) the close/abandon producer mapping (mismatch gives `repository-mismatch`, unsafe gives `recovery-state-conflict`). VAL-031 into (a) the recorded terminal repository and (b) close evaluating Obligations against it. VAL-047 into (a) the sentinel guard and (b) removal of the `PRIVATE_` checks (the latter may join VAL-051's absence inventory). |
| I-01 | I | inventory R15 | 8 read leaves lack a Short, not 9: `concern list` prints `List concerns.`. The 12 `--input` leaves are confirmed. The live tree has 92 leaves (excluding `help`/`completion`), including depth 1 (`capabilities`, `init`, `ready`, `rebuild`, `resume`, `sql`) and depth 3 (`issue ac add/list`, `issue label add/remove`). VAL-010's walk must recurse. The section 1 two-level walk is sufficient only for `--input` leaves. |
| I-02 | I | inventory section 1 walk | The loop's status is the last leaf's `grep -q` (measured `exit 1`, with the correct 39 lines). Say that the line set is the pass condition and the status is recorded as observed. |
| I-03 | I | VAL-006 | Exclude `export import` by name (install-marker idempotency; outside R1). |
| I-04 | I | VAL-007 | For flag and positional leaves, define a canonically equal retry as flag-order permutation. Reordered keys apply to `--input` leaves only. |
| I-05 | I | VAL-003 | Add a malformed-map row (non-object, or a non-digest value) that stops replay. |
| I-06 | I | VAL-004 | Name the failure exactly: `archive-invalid` with check `schema-unsupported` on `manifest.json` (`internal/export/archive_validation.go:130-131`). |
| I-07 | I | VAL-005; task 1 | The 100k journal has no keyed Events (`internal/journal/rebuild_scale_test.go`: `CommandName: "benchmark.replay"`, no `Idempotency`). Task 1's "if the generator builds keyed Events" condition is therefore false, and VAL-005 measures unchanged replay paths. Changing the journal would change a protected benchmark definition (DEC-158); record only. |
| I-08 | I | VAL-018, VAL-019 | Enumerate Source positions from the descriptor registry (every SourceInput-typed field), not from the "for example" list, with a completeness check. |
| I-09 | I | VAL-015 | Add `--project <dir>` (a global flag) to the exclusivity corpus. Contract section 8 says any other flag is invalid. |
| I-10 | I | task 10 | Declare the expected status of the R9 `rg` line (exit 0 expected while family writers remain). |
| I-11 | I | VAL-035; task 9 | Derive the fake/fid leaf set from the fake's operation dispatch intersected with the 17 changed leaves. Measured candidates: `gate.register`, `envelope.create`, `envelope.complete`, `claim.acquire`, `run.start`, `run.close`, `run.abandon`, `gate.run`, `run.observe-budget`. The fake references no judgment leaf as an operation. |
| I-12 | I | VAL-044 | Enumerate the judgment views (lists, the three `get` leaves, mutation outputs, `resume` sections) from the renderers, not by hand. |
| I-13 | I | VAL-052 | Have completion reports also record the slowest ordinary package. Baseline: `cmd/fid` 330.8 s (55% of 600 s), `internal/issues` 204.2 s, tagged 293.1 s (33% of 900 s). Tasks 3-10 add built-binary corpora to `cmd/fid`; watch the 80% warning (480 s). Under the unbudgeted DEC-259 control, `cmd/fid` already takes 539.7 s, so its `-timeout 30m` bound matters. |
| I-14 | I | task 1; DEC-270 | ADD COLUMN with the accepted text fails on a populated table (`Cannot add a NOT NULL column with default value NULL`, measured). Init, rebuild and import migrate empty tables. Catch-up, however, runs `MigrateUp` on the live projection (`internal/journal/projection.go:184`). A Phase 23 workspace with keyed rows therefore fails at migration. No compatibility is owed, but task 1 should record the observed fail-closed behavior and confirm that a failed (dirty) migration does not wedge `fid rebuild`. |
| I-15 | I | task 13 | The section disclaimer must not reuse the three per-item sentences (`cmd/fid/p23_human.go:163, 167, 176`), or the absence `rg` keeps matching. |
| I-16 | I | task 2 | The shared helper's name must match `func [a-zA-Z]*(RequestDigest\|requestDigest\|judgmentDigest)\(`, or "shows only the shared helper" becomes an undeclared exit 1. |
| I-17 | I | security (no VAL) | Per-field digests make each sensitive member (`environment[].value`, `redact_literals`) individually digest-addressable in the journal. Record a disposition. The final-review security check (QA-35) covers it. |
| I-18 | I | phase brief Known Risks | Rebuild headroom is about 19.6% (worst fresh sample 24128 ms), not about 10%. |
| I-19 | I | environment | `/tmp` is a 3.9 GB tmpfs with 1.4 GB free at assessment. Others' scratch (`pytest-of-mark` 773 MB, `kauri-*` ~1 GB) is not fid's to remove. Integration binary directories are about 17-21 MB each. If space falls, a gate fails for environmental reasons (CF-042). |

### 2.4 Blocker Candidates (first pass)

All three were resolved from accepted text or by Human decision: R-04 and
R-05 by DEC-278 items 1-2, and R-07's `operation` value by DEC-279. Section
12.1 verifies them. The open candidate from this refresh is in section 12.4.

| Candidate | Protected category | Evidence | Authority boundary | Requested action |
|---|---|---|---|---|
| R-04 | public interface (failure-code producer) | `evidence evaluate --help` (no key); contract section 6 producers | Producers of an accepted code are contract text (DEC-277) | Tech Lead states, from accepted text, `evidence.evaluate`'s code for an unobservable worktree; otherwise a contract amendment by DEC |
| R-05 | public interface (closed `reason` domain) | `internal/project/repository.go:61-352` causes; `git init` + `fid init` succeeds | A new member of a closed enum needs contract authority (DEC-195) | Tech Lead maps every cause to the five members, or amends section 6 by DEC |
| R-07 | public interface (Human failure rule; unknown-command JSON) | Silent exit 2 measured; Phase 21 lines 406, 459; `codes.go` lacks `operation-unavailable` | Section 9's "every failure" and the JSON `operation` for an unknown command are contract text | Tech Lead assigns conformance under DEC-175/194 with the `operation` value derived from accepted text, or records a section 9 exception by DEC |

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| 1 | Migration 0019, field-digest storage | 001-005, 058, 052 | high | Map is an insert value and a receipt-update guard (DEC-278 item 7); registries and retained tests now in scope (12.1). I-14 note. 100k early sample. |
| 2 | Input-as-written digests, `differing_fields` | 002, 005-009, 066, 052 | high | 49 keyed leaves in 20 families (R16); completeness test; RI-01 (artifact two-field row); RI-04 (control construction). |
| 3 | Descriptor framework, help, `--example`, field paths | 010-016, 059, 060, 052 | high | Fence (DEC-192). VAL-059 non-qualifying leaves are exactly `attention resolve`, `stop raise`, `stop resolve` (12.2). |
| 4 | Observation, SourceInput, judgment leaves | 008, 009, 011-023, 035, 060, 066, 052 | high | Five-reason mapping; tracked-FIFO fixture and capture-stage classification verified under DEC-280. Human Outcome references become checked. |
| 5 | `ExpectedRepository`, `review.record` | 009, 011-019, 022-025, 035, 060, 066, 052 | medium | Help names the git commands (VAL-012); `dirty` meaning matches fid's capture (12.1 R-13). |
| 6 | Registration, Envelope, Claim; same-worktree predicate | 009, 011-019, 022, 023, 026-029, 035, 060, 066, 052 | high | Predicate reused by tasks 7-8. Plain copy is a working workspace (12.2). Exact fallback refusals (DEC-278 item 6). |
| 7 | Run leaves (DEC-273) | 009, 011-019, 022-025, 029-031, 035, 060-062, 066, 052 | high | `base_oid` becomes HEAD; `run.close` records `failed` on stale Evidence (`run_close.go:81-87`). |
| 8 | Gate run, Evidence, real freshness | 009, 011-017, 022-025, 029, 032-035, 060, 066, 052 | high | Highest convergence risk. `evidence.evaluate` is the one observing read. Artifacts stay git-ignored. |
| 9 | Activation descriptors; whole-tree input and help conformance | 011-016, 035, 060, 052 | **high** | Completes the descriptor registry and fake/fid assembly; R5 naming constraint. |
| 9a | Grammar failures and unknown commands | 064, 065, 052 (036 advanced) | **high** | R17 root/groups, R18 typed flags plus cited string domains, R19 exactly four residual fallbacks; control escaping; non-project repetitions. RC-02 requires a bounded-message rule. Review before task 10. |
| 10 | Human failure lines and generic-fallback messages | 036-039, 067, 027-029 messages, 052 (054) | **high** | 41 writers (R9); R20 has 18 fallbacks and seven excluded matches. DEC-281 adopted; RC-01 still needs the TL's plain-copy message correction. |
| 11 | Stop, Attention, repository lines; Human fixture generator | 040, 041, 052 (054, 055) | medium | `TestP24HumanFixtures` (DEC-278 item 10). Every Condition type has existing resolve fixtures (12.1 R-10). |
| 12 | Judgment `get` leaves | 042, 043, 052 | low | Depends on tasks 4 and 10; extends the task 9a/10 corpora (DEC-280). |
| 13 | Disclaimers, recorded reasoning, sentinel guard, glossary | 044-048, 063, 052 (056) | medium | `operations.json` and reader in scope (R-15); I-15. |
| 14 | Take-control journey, parity, absence | 049-052, 063 (053-057) | high | Closed provenance classes (DEC-278 item 9). RI-08 (real freshness in the journey), RI-09 (red-first test). |

## 4. Required Evidence Strategy

Rows state what QA requires; they do not prescribe implementation. "Named
tests" means tests the completion report names and whose output it records.
"Independent observation" is the test's own `git rev-parse HEAD`,
`git rev-parse HEAD^{tree}`, `git status --porcelain` and SHA-256 of the
absolute root, never fid output. EVID-38..44 were added in this refresh;
existing IDs are not renumbered.

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | `sqlite_master` (type, name, tbl_name, sql) dumps at 18, after up, after down; the up diff is exactly the one column; the down dump is byte-equal to 18. Each half is reported against its own assertion | VAL-001, VAL-058 | Developer | task 1 |
| EVID-02 | Live versus rebuilt `request_field_digests_json`, byte-equal per row (task 1: issue, claim, run, review; task 2: one keyed mutation per keyed leaf, R16) | VAL-002 | Developer | tasks 1, 2 |
| EVID-03 | Journal fixtures: map present replays; removed, non-object and non-digest value each stop with `journal-replay-stopped`, exit 6, naming the line; receipt-update guard mismatch fails closed | VAL-003 | Developer | task 1 |
| EVID-04 | Export, validate, import of a 19 archive; 18-manifest archive rejected `archive-invalid`/`schema-unsupported` on `manifest.json`, target tree digest unchanged | VAL-004 | Developer | task 1 |
| EVID-05 | `BenchmarkRebuildWorkspace100kEvents` `< 30000 ms/rebuild` with count parity; final review three fresh processes plus instance identity | VAL-005 | Developer, QA Lead | tasks 1, 2; final review |
| EVID-06 | Per-keyed-leaf conflict corpus (49 leaves, R16) with a completeness test: one-field row per leaf; two-field and `--actor` rows per family; exact sorted `differing_fields`, exit 4, journal line count unchanged; canonical retry replays byte-identically except `replayed:true` with the original EventRef | VAL-006, VAL-007 | Developer | task 2; converged by tasks 4-9 |
| EVID-07 | Three presence states: nonempty versus omitted `inapplicable_findings` conflicts naming it; `[]` versus omitted replays; supplied versus omitted Source `digest` conflicts naming the top-level field | VAL-008 | Developer | tasks 2, 4 |
| EVID-08 | Per observing mutation: commit, change the observed value (HEAD, dirty, file with digest omitted), retry, which replays with no Event and no failure | VAL-009 | Developer | tasks 2, 4-8 |
| EVID-09 | In-process recursive Cobra walk: every leaf `Short` nonempty, no help starts `Usage:` | VAL-010 | Developer | task 3; kept green by task 12 |
| EVID-10 | Reflection conformance, descriptor versus Go type, for every registered leaf; seeded mismatch fails; registry ∪ fence = live `--input` set until task 9, registry = live set after | VAL-011 | Developer | tasks 3-9 |
| EVID-11 | Help conformance per descriptor field (path, type, requirement, bounds or allowed values, meaning, template line); `expected_repository` meanings name `git rev-parse HEAD` and `git status --porcelain`; full text pinned for `stop raise`, `activation component` and one more | VAL-012 | Developer | tasks 3-9 |
| EVID-12 | Built-binary `--example` per leaf: one object, required fields only, documented order, 2-space indent; template decoded with unknown fields disallowed and type-checked against the descriptor | VAL-013, VAL-014 | Developer | tasks 3-9 |
| EVID-13 | Exclusivity corpus (`--json`, `--project <dir>`, `--input`, `--actor`, `--idempotency-key`, positional) returns `invalid-argument` exit 2 with task 3's details and no read or write | VAL-015 | Developer | tasks 3-9 |
| EVID-14 | Structural corpus per leaf: 7 failure classes at top level and two levels deep; exact `input.<path>` and `allowed`; `sensitive` value null; no Event | VAL-016 | Developer | tasks 3-9 |
| EVID-15 | Removed-member corpus over the 17 changed leaves, each rejected naming its path | VAL-017 | Developer | tasks 4-8 |
| EVID-16 | Every Source position (from the descriptor registry, with completeness check) × matching, mismatching, omitted-existing, omitted-missing; five path-safety fixtures; external reference omitted versus supplied | VAL-018-021 | Developer | tasks 4-7 |
| EVID-17 | Stored record repository, each Source repository and the Event `repository` compared with the independent observation (`base_oid` = HEAD, `worktree_id` = SHA-256 of root, `project_id` from config) | VAL-022 | Developer | tasks 4-8 |
| EVID-18 | Unmerged-paths and non-Git fixtures versus a normal worktree per observing leaf (except close/abandon); cause-class rows for no commit, tracked flag, unsupported file type (tracked file replaced by a FIFO, RR-09), injected capture change, injected size limit | VAL-023 | Developer | tasks 4-8 |
| EVID-19 | Five expectation leaves × match, HEAD differs, dirty differs; mismatch details exactly `{expected:{head_oid,dirty},observed}` | VAL-024, VAL-025 | Developer | tasks 5, 7, 8 |
| EVID-20 | Registration fields and `binding_digest` recomputed from the independent observation; same-worktree versus plain-copy and imported destination with binding for four leaves. Tasks 6-8 retain fixed JSON messages; task 10 converges to DEC-281 domain causes, after RC-01 is dispositioned | VAL-026-029 | Developer | tasks 6-8; message convergence task 10 |
| EVID-21 | `run.start` after commit and edit since Envelope creation records the observation; close/abandon `terminal_repository` observed with `base_oid` = HEAD | VAL-030, VAL-031 | Developer | task 7 |
| EVID-22 | Gate run after an edit since Run start: Attempt, Receipt and both Events carry one identical repository equal to the pre-spawn observation; mismatch leaves no Attempt, no process and no Event | VAL-032, VAL-024 | Developer | task 8 |
| EVID-23 | Freshness: unchanged, edited, reverted × `evidence.evaluate`, terminal-transition denial, `envelope.complete` denial | VAL-033 | Developer | task 8 |
| EVID-24 | `evidence.link` target and Event `repository` equal the subject's recorded repository for receipt, review and Human Outcome subjects | VAL-034 | Developer | task 8 |
| EVID-25 | Shared fake/fid corpus over the derived leaf set, identical acceptance and rejection, no one-sided case | VAL-035 | Developer | tasks 4-9 |
| EVID-26 | Built-binary Human/JSON pairing from `RegisteredCodes()` with pinned exclusions; grammar and unknown-command rows (including a near-miss word, RR-02); the post-frame `gate run` row (RR-07); `blocked by:` both states; next-table rows plus five outside codes plus the non-`--input` row and the `--input`-leaf grammar row (RR-01); reason-rendering units with the 200/201 boundary | VAL-036-039 | Developer | task 10 |
| EVID-27 | `stop list`, `attention list`, `resume` over 10 codes and 10 Condition types, open and resolved, no raw `{`; repository lines over clean/dirty and same/different worktree (plain copy via `--repository`) with a forbidden-token scan | VAL-040, VAL-041 | Developer | task 11 |
| EVID-28 | `get` JSON equals the list item with `event:null`; unknown ID gives `*-not-found` exit 3 `{id}`; missing ID, extra argument, unknown flag give `invalid-argument` exit 2 | VAL-042, VAL-043 | Developer | task 12 |
| EVID-29 | Disclaimer count for 0/1/3 items per enumerated view; 80/81/multi-line text per closed field in list and detail views; ESC/CR/BEL/NUL escaped; sentinel guard over every Human fixture | VAL-044-047 | Developer | task 13 |
| EVID-30 | `GLOSSARY.md` entries compared item by item with the contract section 10 closed field list | VAL-048 | QA Lead | task 13; final review |
| EVID-31 | Tagged take-control journey with a provenance record (DEC-278 class and source line) per command and value; live/rebuilt/imported parity over its workspace, manifest 19 | VAL-049, VAL-050 | Developer | task 14 |
| EVID-32 | Descriptor forbidden-member scan; per-item disclaimer absence `rg` terminating `exit_code: 1`; no raw-JSON Stop renderer | VAL-051 | Developer, QA Lead | task 14; final review |
| EVID-33 | Per-task wall times for `cmd/fid/integration`, `cmd/fid`, `internal/workflow` and the slowest ordinary package | VAL-052 | Developer | every task |
| EVID-34 | `TestP24HumanFixtures` output (manifest with workspace, cwd, command lines, expected exit, binary SHA-256) and Human judgment notes | VAL-053-057 | Developer fixture (tasks 11, 14); Human judgment | window D2; phase end |
| EVID-35 | Full phase gate, DEC-259 full-`checkptr` control, archive-scale gate, three 100k samples, reference instance | no VAL | QA Lead | final review |
| EVID-36 | Security review: observation under lock and never on replay; path containment and symlinks; digest-before-observe; sensitive-field redaction; no top-level sensitive descriptor field (I-17); import binding | no VAL | QA Lead | final review |
| EVID-37 | Inventory registries R1-R20 and the section 1 block regenerated, with deltas explained | no VAL | QA Reviewer, QA Lead | every window; final review |
| EVID-38 | Per observing mutation: retry in an unobservable worktree (unmerged paths) and, for project-file Sources, after deleting the file, each replaying with no Event; a new-key control that would succeed in an observable worktree fails with no Event (from task 4: `repository-unavailable`, `recovery-state-conflict` for close/abandon, `source-digest-mismatch`) | VAL-066 | Developer | tasks 2, 4-8 |
| EVID-39 | `attention resolve --help` text; every leaf that used `runMutationFlags` outside the Run family (`attention resolve`, `stop raise`, `stop resolve`) has no Run wording | VAL-059 | Developer | task 3 |
| EVID-40 | Built-binary `--example` per leaf from a non-project temp dir: tree snapshot digest before and after, exit 0, empty stderr | VAL-060 | Developer | tasks 3-9 |
| EVID-41 | Close and abandon: mismatching expectation gives `repository-mismatch`; unmerged paths give `recovery-state-conflict`; neither appends | VAL-061 | Developer | task 7 |
| EVID-42 | Receipt-backed Obligation: close in the proven tree gives `succeeded`; separate fixture after a tracked edit gives `failed` | VAL-062 | Developer | task 7 |
| EVID-43 | `rg -n "PRIVATE_" cmd/fid/integration` terminating `exit_code: 1` | VAL-063 | Developer, QA Lead | task 13; task 14; final review |
| EVID-44 | Whole-tree grammar corpus per task 9a: leaf/group completeness, every typed flag and cited string domain, ESC spelling, missing value, non-project repetitions; unknown word under each group, near miss, and no-word controls. Exclude Human-only diagnostics and generated Cobra commands as declared; include help/completion in root alternatives | VAL-064, VAL-065 | Developer | task 9a |
| EVID-45 | R20 site-to-row map: both-mode reachable fallback causes or individually justified unreachability; unchanged metadata, details and exit; matching escaped JSON/Human text; 512/513 scalars, empty and whitespace-only causes, escape crossing the cut; sensitive activation-component guard reaches the generic fallback, not an earlier typed failure | VAL-067 | Developer | task 10 |

## 5. Deterministic Gates

| Gate | Command | Expected Result | Required When |
|---|---|---|---|
| Format | `gofmt -l cmd internal` | no output, `exit_code: 0` | every task, final review |
| Lint | `golangci-lint run ./...` | `0 issues.`, `exit_code: 0` | every task, final review |
| Race + coverage | `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` | all packages `ok`, `exit_code: 0`; no ordinary package over 600 s | every task, final review |
| Coverage report | `go tool cover -func=coverage.out` | `total:` line recorded, `exit_code: 0` | every task, final review |
| Integration | `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` | `ok ... cmd/fid/integration`, under 900 s, `exit_code: 0` | every task, final review (hosted CI does not run it, CF-041) |
| Reference rebuild | `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` | `ok`, count parity, `< 30000 ms/rebuild` | tasks 1, 2; final review three consecutive fresh processes (DEC-158) |
| Reference instance | `go version`; `go env GOOS GOARCH`; `uname -srm`; `nproc`; `grep -m1 '^model name' /proc/cpuinfo`; `grep -m1 '^MemTotal:' /proc/meminfo` | recorded with the samples | final review |
| Registry and absence searches | the `rg` lines in tasks 1, 2, 9, 9a, 10, 11, 13, 14 | recorded real status; task 9 R5 and the tasks 13/14 absence lines declared `exit_code: 1`; task 9a R19 and task 10 R20 expect `exit_code: 0` with accounted matches | owning task; window; final review |
| Human fixtures | `FID_P24_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-24-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP24HumanFixtures$' -count=1 ./cmd/fid/integration` | `ok`, manifest written | window D2 (QA Reviewer, task-11 close); phase end (QA Lead, task-14 close) |
| Benchmarks | `go test -bench . -benchmem ./...` | all packages `ok`, `exit_code: 0` | final review |
| Archive scale | `go test -tags archive_scale -race ./internal/export -run '^TestArchive(ExpandedByteBoundaries\|CumulativeBoundary\|EntryCountBoundaries)$' -count=1` | `ok`, `exit_code: 0` | final review (DEC-217 phase-end gate) |
| Full-`checkptr` control | `go test -race -timeout 30m ./...` | `ok`, `exit_code: 0`; not held to VAL-052 budgets | final review (DEC-259); any task that changes a `modernc.org` version |
| Build | `go build -o /tmp/fid ./cmd/fid` | `exit 0 (no output)` | final review |
| Full phase gate | format, lint, race + coverage, coverage report, integration, benchmarks, build, as above | every command `PASS` with observed termination | final review |
| End-to-end | `N/A` | no held-out harness recorded | - |

The two race gates carry `-timeout 30m`, as every Phase 24 task spec does
(DEC-227 added per phase; DEC-259 flag). `AGENTS.md` still lists them
without the bound. Only the bounded form was executed (section 6.2). At
today's package times, all under 600 s, the unbounded form would also stay
inside Go's default limit. The task-spec form governs. No gate is narrowed,
substituted or deferred. Tasks 9, 13 and 14 declare expected non-zero
commands; task 9a R19 instead expects four matches and exit 0. Cleanup commands are never gate
results. Report lint
(`project-docs/tools/report-lint.sh <task-N-done.md> <task-N.md>`) is a Phase
Runner control after the report commit, never a Developer or QA gate. Run
heavy gates serially.

## 6. Command And Feasibility Validation

Sections 6.1-6.9 are the first-pass results at `304aa45`. The product trees at
`21c833f` are byte-identical (section 12.2), and no gate command's text
changed, so these results still hold. Every command, fixture, scope and probe
that DEC-278 or DEC-279 changed is validated fresh in section 12.2.

Heavy mode: nothing carries forward from Phase 23. The working directory for
every command is `/home/mark/projects/fid` unless stated. A serial wrapper
captured each command's own exit status and wall time
(`bash -c "<cmd>"; rc=$?`). Logs are in session scratch storage.
`go clean -testcache` ran first as setup, not as a gate result, so no result
is a cached replay. Light read-only probes ran during gate steps 3-9. Gate
wall times therefore include minor contention, and no step failed.

### 6.1 Toolchain And Reference Instance

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Go toolchain | executed | PASS | `go version go1.26.2 linux/amd64`; `go env GOOS GOARCH` gives `linux` `amd64` |
| Kernel / CPU / memory | executed | PASS | `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64`; `nproc` 6; `AMD Ryzen 7 5700G with Radeon Graphics`; `MemTotal: 8133236 kB` |
| Linter | executed | PASS | `golangci-lint has version 2.12.2 built with go1.26.2` |
| Git, ripgrep, SQLite CLI | executed | PASS | `git version 2.53.0`; `ripgrep 15.1.0`; `sqlite3 3.50.6` |
| Runtime SQLite and migrator | `go.mod` inspection | PASS | `modernc.org/sqlite v1.51.0`; `github.com/golang-migrate/migrate/v4 v4.19.1` |
| `goimports` | `which` | PASS | `/home/mark/go/bin/goimports` |

### 6.2 Phase Gate Commands On The Planning Base

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `gofmt -l cmd internal` | executed | PASS | `exit 0 (no output)` (terminated `exit_code: 0`, 0.2 s) |
| `golangci-lint run ./...` | executed | PASS | `0 issues.` (terminated `exit_code: 0`, 14.9 s) |
| `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` | executed | PASS | last line `ok github.com/mbell-nzl/fid/internal/workflow 110.296s coverage: 62.7% of statements`; 0 `FAIL` lines (terminated `exit_code: 0`, 347.0 s wall). Slowest: `cmd/fid` 330.801 s, `internal/issues` 204.209 s, `internal/export` 189.427 s, `internal/journal` 183.282 s, `internal/workflow` 110.296 s |
| `go tool cover -func=coverage.out` | executed | PASS | `total: (statements) 60.7%` (terminated `exit_code: 0`); `coverage.out` is gitignored |
| `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` | executed | PASS | `ok github.com/mbell-nzl/fid/cmd/fid/integration 293.064s` (terminated `exit_code: 0`, 303.5 s wall) |
| `go test -bench . -benchmem ./...` | executed | PASS | last line `ok github.com/mbell-nzl/fid/internal/workflow 8.370s`; 0 `FAIL` lines; in-run `22597 ms/rebuild` (100k), `2206 ms/rebuild` (10k) (terminated `exit_code: 0`, 121.8 s) |
| `go build -o /tmp/fid ./cmd/fid` | executed | PASS | `exit 0 (no output)` (terminated `exit_code: 0`); 16899656-byte binary |

### 6.3 Task-Specific And Registry Commands

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| 100k rebuild (tasks 1, 2; final ×3) | executed three times, fresh processes | PASS | `22653` / `23211` / `24128 ms/rebuild`; `ok ... internal/journal 25.506s` / `25.995s` / `26.885s` (each terminated `exit_code: 0`). Worst-case headroom 19.6% |
| Archive-scale gate | executed | PASS | `ok github.com/mbell-nzl/fid/internal/export 143.927s` (terminated `exit_code: 0`, 164.6 s wall) |
| Full-`checkptr` control `go test -race -timeout 30m ./...` | executed | PASS | last line `ok github.com/mbell-nzl/fid/internal/workflow 171.784s`; 0 `FAIL` lines (terminated `exit_code: 0`, 553.2 s wall). `cmd/fid` 539.666 s, `internal/issues` 379.787 s; not held to VAL-052 budgets (DEC-259), bounded by `-timeout 30m` |
| R1 `rg -n "RequestDigest:" internal cmd -g '!*_test.go'` (task 2) | executed | PASS | `exit_code: 0`; 43 lines in 32 files, matching the inventory |
| R2 helper `rg` (task 2) | executed | PASS | `exit_code: 0`; 11 helpers in `internal/{links,review(2),evidence(2),workflow(4),labels,issues}` |
| R3 `rg -n "idempotency_results" migrations internal -g '!*_test.go'` (task 1) | executed | PASS (syntax); inventory incomplete at first pass (R-01, since corrected) | `exit_code: 0`; 16 files |
| R4 | `ls migrations`; source | PASS | latest `0018_p23_remove_v1_judgment`; `ActiveMigrationVersion = 18`, `ProjectionSchemaVersion = 5` (`internal/db/projection_meta.go:10-11`); manifest uses `ActiveMigrationVersion` (`internal/export/export.go:137`) |
| R5 `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` (task 9) | executed | PASS (syntax) | `exit_code: 0`; 7 helpers, 14 calling files |
| R6 | executed | PASS | 4 `same*Repository`, 6 validators, `importedGateDestination` (`internal/evidence/execution.go:338`); `importedDestinationBinding` (`internal/workflow/import_binding.go:18`) |
| R7 | executed | PASS | caller-repository writers in `claims/generic.go`, `evidence/{evidence,execution,registration}.go`, `review/review.go`, `workflow/{envelope,run,run_close}.go` |
| R8 `internal/cli/codes.go` | inspection | PASS | 78 registered codes; `repository-unavailable` and `operation-unavailable` absent; `idempotency-conflict` details 4 strings; `source-digest-mismatch.expected_digest` non-nullable; `repository-mismatch.expected` a full Repository |
| R9 `rg -n "func (write[A-Za-z]*Failure\|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'` (task 10) | executed | PASS | `exit_code: 0`; 41 functions |
| R10 | executed | PASS | 32 `"error: ` assertions in 13 test files |
| R11 `rg -n -i "worktree\|dirty digest\|HeadOID\|head_oid" cmd/fid -g '!*_test.go'` (task 11) | executed | PASS | `exit_code: 0`; Human lines in `p23_human.go:139-140`, `resume_human.go:113-117`, `run_inspect_human.go:17-19`; JSON mapping in `review.go:264`; a flag description in `run_recovery.go:68` |
| R12 | `ls` | PASS | `stop_human.go`, `resume_human.go`, `resume_human_guidance.go` present |
| R13 `rg -n "p23CheckJudgmentQualifiers\|PRIVATE_" cmd` | executed | PASS (syntax) | also matches `cmd/fid/integration/testdata/p23_conformance/operations.json` four times (R-15, since corrected) |
| R14 | executed | PASS | `cmd/fid/contract_surface_test.go:128` |
| R15 Short walk | executed recursively over the live tree | PASS | 92 leaves; 20 lack a Short (8 read, 12 `--input`) |
| Section 1 `--input` walk | executed | PASS | exactly the 39 inventory leaves |
| Section 3 fixture families | executed | PASS | 56 test files name caller repository fields (29 integration); 8 pinned literal digests in 7 files; 4 `p2*_fake_test.go` files |
| Task 13/14 per-item disclaimer `rg` | executed | PASS (syntax); matches expected before task 13 | `exit_code: 0`; `cmd/fid/p23_human.go:163, 167, 176` |
| Task 13/14 `rg -n "PRIVATE_" cmd/fid/integration` | executed | PASS (syntax) | `exit_code: 0`; `p23_human_fixtures_test.go:143`, `p23_human_leaf_pairs_test.go:45-126`, `testdata/p23_conformance/operations.json:67-145` |
| `project-docs/tools/report-lint.sh` | usage output | PASS | `Usage: project-docs/tools/report-lint.sh <completion-report.md> <task-spec.md> [--repo <git-dir>]` |
| New file paths | `ls` | PASS | `migrations/0019_*`, `internal/db/migration0019_test.go`, `cmd/fid/input_descriptor_test.go`, `cmd/fid/integration/p24_take_control_test.go` absent; `internal/db/migration0018_test.go` present as pattern |

### 6.4 Accepted DDL - Applied, Not Read

Every `migrations/*.up.sql` from 0001 to 0018 was applied in order to a
scratch database. Each file ran under `PRAGMA foreign_keys=ON` inside one
`BEGIN ... COMMIT`, mirroring the runner. The database was not the repository.

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Fresh chain 0001-0018 | applied | PASS | `integrity_check` ok, `foreign_key_check` empty; `idempotency_results` has 8 columns plus `PRIMARY KEY (operation_scope, idempotency_key)` and `sqlite_autoindex_idempotency_results_1` |
| Accepted column, empty table | `ALTER TABLE idempotency_results ADD COLUMN request_field_digests_json TEXT NOT NULL CHECK (json_valid(request_field_digests_json))` | PASS | rc 0; the 18→19 `sqlite_master` diff is exactly one line (`created_at DATETIME NOT NULL, request_field_digests_json TEXT NOT NULL CHECK (json_valid(request_field_digests_json)),`) |
| Constraint behavior at 19 | INSERT probes | PASS | non-JSON rejected `CHECK constraint failed` (19); omitted column rejected `NOT NULL constraint failed` (19); `'{}'` accepted |
| Down to 18 | `ALTER TABLE idempotency_results DROP COLUMN request_field_digests_json` | PASS | `sqlite_master` (type, name, tbl_name, sql) byte-equal to the schema-18 dump (`cmp` silent) |
| Accepted column, populated table | same ADD COLUMN after one row | **FAIL (informational, I-14)** | `Cannot add a NOT NULL column with default value NULL`, rc 1; init, rebuild and import migrate empty tables; catch-up migrates in place |

### 6.5 CLI Surface Assumptions

Probes used a scratch build of the planning base against fresh projects in
session scratch storage.

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Keyed leaf set (VAL-006) | `--help` scan of all 92 leaves | PASS | 50 leaves take `--idempotency-key` across 21 families; `export import` is the only one outside R1; `evidence evaluate` is the only `--input` leaf without a key |
| Task 3 first failing test | `fid --json attention resolve --input {"bogus":1}` | PASS (red today) | `invalid-argument`, `argument:"input"`, message `invalid Run request`; exit 2 |
| `--example` absent today | `fid stop raise --example` | PASS (red today) | `error: unknown flag: --example`, exit 2; JSON `operation:"raise"` (not `stop.raise`) |
| `attention resolve --help` text | executed | PASS (red today) | `--input string Run JSON input`; first line `Usage:` |
| Task 12 leaves absent | `fid concern get CON-1` | PASS (red today) | exit 2, no output in either mode |
| Unknown command and subcommand | `fid bogus`, `fid concern bogus` in both modes | FAIL vs Phase 21 contract (R-07) | exit 2, empty stdout and stderr |
| Grammar failure parity | `fid stop list --bogus`, `fid review get` in both modes | conflicts with section 9 intent (R-07) | Human `error: unknown flag: --bogus` / `error: accepts 1 arg(s), received 0`; JSON message `invalid command request`, `argument:"request"`, `value:null` |
| Task 6 first failing test | source (`internal/claims/generic.go:483`) | PASS (red today) | empty `worktree_id` fails validation |
| Task 13 first failing test | source (`cmd/fid/p23_human.go:154-163`) | PASS (red today) | Concern item lines print subject, severity, status, disposition, sources and the per-item disclaimer; no `summary` or `reasoning` |
| Task 11 fixture reach | source (`internal/workflow/stop.go:175-189, 440-500`) | PASS | 10 codes, 10 Condition types; `stop raise` accepts any pairing with bindings |
| Observation causes | source (`internal/project/repository.go`) | gap at first pass (R-05, since resolved) | about 30 `unsafeRepository` causes; operation markers are recorded, not rejected |
| Existing codes used by assertions | `internal/cli/codes.go` | PASS | `journal-replay-stopped`, `recovery-state-conflict`, `evidence-stale`, `concern/deviation/outcome-not-found` (not-found, `{id}`), `archive-invalid`, `idempotency-conflict`, `source-digest-mismatch`, `repository-mismatch` present |
| Plain-copy refusal codes | source | PASS | `cmd/fid/envelope.go:253`, `cmd/fid/claim.go:234` generic `invalid-argument` |
| JSON failure stream | probes | noted | JSON failure envelopes are written to stderr with empty stdout |

### 6.6 Scope And Gate Reach

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Gate reach covers owned paths | task file lists against registry output | PASS at first pass with R-01, R-15, R-07 (since corrected); see RR-03 | both race gates compile every `cmd/`/`internal/` package |
| Same-worktree predicate home | source | PASS | `internal/workflow/import_binding.go` exists (task 6 names it) |
| Fake adapters present | `ls` | PASS | `p22_fake_test.go`, `p22_semantic_fake_test.go`, `p22_portable_fake_test.go`, `p22_activation_fake_test.go` |
| Glossary starting state | `GLOSSARY.md:21-31` | PASS | "recorded reasoning" exists only as a note inside Project-owned prose, with a partial list; VAL-048's non-qualifying state exists today |

### 6.7 Human Checklist Feasibility

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Disposable project in a non-Git directory | executed | PASS | `fid init --prefix PNG` gives `initialized fid project in .fid/` (also the VAL-023 non-Git fixture) |
| Disposable project with unborn HEAD | executed | PASS | `git init` then `fid init --prefix PUB` succeeds (VAL-023 `not-a-repository` row) |
| Phase 23 fixture generator | source | PASS | `TestP23HumanFixtures` with `FID_P23_HUMAN_FIXTURE_ROOT` (`cmd/fid/integration/p23_human_fixtures_test.go:57-63`), including a protected-Decision checkpoint (line 264) |
| Phase 24 stuck-journey fixture | task specs | PASS after DEC-278 item 10 (first pass FAIL, R-14) | `TestP24HumanFixtures`, task 11, extended by task 14 (12.1) |
| Judgment not delegable | role rules | PASS | VAL-053..057 are `Human`; no `(qa-lead)` item substitutes |

### 6.8 Documentation Targets

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `03_TRD_SCHEMA.md` (task 1) | inspection | PASS | `idempotency_results` row at line 71; DDL at line 310 |
| `05_DESIGN_PLAN.md` section 6 (task 10) | header | PASS | `## 6. Errors And Stops` at line 111 |
| `06_DESIGN_PLAN_CLI.md` (tasks 2-13) | header | PASS | section 5 failure codes at line 259 (`repository-mismatch` line 344, `source-digest-mismatch` line 343); section 7 Human rendering at line 379 |
| `07_DESIGN_PLAN_FLOWS.md` (tasks 6-8) | header | PASS | sections 3, 5, 6, 8 at lines 26, 58, 73, 118 |
| `08_QUALITY_PLAN.md` (task 14) | header | PASS | section 4 critical scenarios at line 147 |
| `GLOSSARY.md` (tasks 4, 13) | inspection | PASS | target entries absent today |

### 6.9 Storage And Time Budgets

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| `/tmp` filesystem | `df -h /tmp` | CONSTRAINT | first pass: tmpfs 3.9 GB, 1.4 GB free; refresh: 1.2 GB free (RI-13); no stale `fid-integration-binary-*` directories |
| Disk-backed alternative | `df -h /` | PASS | 896 GB available |
| Per-binary test time | measured | PASS | tagged 293.1 s (33% of 900 s); `cmd/fid` 330.8 s (55% of 600 s); `internal/workflow` 110.3 s; all well inside the `-timeout 30m` bound |

## 7. Deterministic Checklist

The QA Lead executes or verifies these at final review. "Rerun named tests"
means `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run
'^(<names>)$' <packages>` (adding `-tags integration` for tagged tests), using
the exact names in the completion report, with output recorded. Real-process
probes use the `/tmp/fid` built by the phase gate against a fresh `git init` +
`fid init --prefix <P>` project. PASS without evidence is invalid. Each
assertion keeps its own outcome slot even where one item covers two.
QA-41..48 were added in this refresh.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Full seven-command phase gate in the section 5 form, each with token and observed termination `(qa-lead)` | no VAL | command transcript |
| QA-02 | DEC-259 full-`checkptr` control and DEC-217 archive-scale gate `(qa-lead)` | no VAL | transcripts |
| QA-03 | Three consecutive fresh 100k samples, each `< 30000 ms/rebuild`, with the reference instance captured in the same session `(qa-lead)` | VAL-005 | transcripts, instance values |
| QA-04 | Fresh-init `sqlite_master` dump versus a schema-18 chain (one-column diff, VAL-001); apply the 0019 down migration and `cmp` with 18 (VAL-058); rerun task 1 named tests `(qa-lead)` | VAL-001, VAL-058 | dumps, `cmp` status, test output |
| QA-05 | On one workspace, `fid rebuild` then compare `request_field_digests_json` live versus rebuilt, plus the task 1/2 parity tests `(qa-lead)` | VAL-002 | SQL dumps, test output |
| QA-06 | Rerun the replay fixtures (present, removed, non-object, non-digest value); confirm the stop line and exit 6 `(qa-lead)` | VAL-003 | test output |
| QA-07 | Export, validate, import a 19 archive; edit a copy's manifest to 18 and confirm the named rejection with the target unchanged `(qa-lead)` | VAL-004 | outputs, tree digests |
| QA-08 | Rerun the per-keyed-leaf conflict and canonical-retry corpus and its completeness test; confirm 49 rows against R16; reproduce one conflict per family through the binary `(qa-lead)` | VAL-006, VAL-007 | test output, probe JSON |
| QA-09 | Reproduce the three VAL-008 presence states through the binary `(qa-lead)` | VAL-008 | probe JSON |
| QA-10 | Per observing mutation, rerun the advanced-HEAD (and changed-file) retry; journal line count unchanged `(qa-lead)` | VAL-009 | test output, line counts |
| QA-11 | Recursive live-tree walk: every leaf has a Short and no help starts `Usage:` `(qa-lead)` | VAL-010 | walk output |
| QA-12 | Registry equals the live `--input` set (39), no fence, reflection conformance and seeded-mismatch tests `(qa-lead)` | VAL-011 | test output, walk diff |
| QA-13 | For every `--input` leaf from a non-project temp dir: `--example` form, template decode, and `--help` lists every descriptor field and the template line; `expected_repository` names the git commands `(qa-lead)` | VAL-012, VAL-013, VAL-014 | per-leaf outputs |
| QA-14 | Exclusivity corpus rerun, plus binary probes of `--json --example` and `--project <dir> --example` `(qa-lead)` | VAL-015 | outputs |
| QA-15 | Structural corpus rerun, plus binary probes of one nested unknown key, one enum violation and one sensitive field `(qa-lead)` | VAL-016 | outputs |
| QA-16 | Removed-member probe for each of the 17 changed leaves `(qa-lead)` | VAL-017 | outputs |
| QA-17 | Source corpus rerun; binary probes of match, mismatch, omitted-existing, omitted-missing, absolute path, `..`, escaping symlink, directory, in-root symlink, external reference `(qa-lead)` | VAL-018-021 | outputs, stored rows |
| QA-18 | For one leaf per task 4-8, stored and Event repositories against an independent observation `(qa-lead)` | VAL-022 | SQL rows, git output |
| QA-19 | Unmerged-paths and non-Git probes per observing leaf (except close/abandon) versus normal; cause-class rows: no commit, `--assume-unchanged`, tracked file replaced by a FIFO, injected capture change, injected size limit `(qa-lead)` | VAL-023 | outputs, test output |
| QA-20 | Expectation corpus (5 leaves × 3 states); mismatch details shape `(qa-lead)` | VAL-024, VAL-025 | outputs |
| QA-21 | Registration binding digest recomputed; same-worktree versus plain copy for `envelope.create`, `claim.acquire`; imported destination versus plain copy for four leaves, each refusal with the exact DEC-278 message and details `(qa-lead)` | VAL-026-029 | outputs, recomputation |
| QA-22 | `run.start` after commit and edit; close/abandon terminal repository `(qa-lead)` | VAL-030, VAL-031 | outputs, SQL rows |
| QA-23 | Gate run after an edit: one identical repository across Attempt, Receipt and both Events; mismatch leaves no Attempt row, no Event, and no marker file from the gate's command (proving no spawn) `(qa-lead)` | VAL-032 | SQL rows, journal lines |
| QA-24 | Freshness three states × three surfaces `(qa-lead)` | VAL-033 | outputs |
| QA-25 | `evidence.link` target and Event repository per subject type `(qa-lead)` | VAL-034 | SQL rows, journal lines |
| QA-26 | Rerun the fake/fid corpus; confirm the derived leaf set has no one-sided case `(qa-lead)` | VAL-035 | test output, case inventory |
| QA-27 | Rerun the pairing corpus; confirm rows plus exclusions equal `RegisteredCodes()`; binary probes of one Stop-blocked denial, one next-table code per row, one non-`--input` `invalid-argument`, one grammar failure on an `--input` leaf (per the RR-01 disposition), one unknown command and one near-miss word such as `fid isue` (RR-02), one post-frame `gate run` failure (RR-07), and the 200/201-scalar boundary `(qa-lead)` | VAL-036-039 | test output, stderr captures |
| QA-28 | Stop and Attention views over every code and Condition type, open (with `Next:` and the after-the-fact step) and resolved (resolution Event, no `Next:`); repository lines clean/dirty and same/different worktree; forbidden-token `rg` over captured Human output `(qa-lead)` | VAL-040, VAL-041 | captured output, `rg` status |
| QA-29 | `get` leaves versus list items; not-found and grammar probes `(qa-lead)` | VAL-042, VAL-043 | outputs |
| QA-30 | Disclaimer counts over the enumerated view set, reasoning cut and full text, control escaping, sentinel guard rerun `(qa-lead)` | VAL-044-047 | captured output, test output |
| QA-31 | `GLOSSARY.md` entries compared item by item with contract section 10 `(qa-lead)` | VAL-048 | quoted entries and comparison |
| QA-32 | Rerun the take-control journey; check each provenance record's class against DEC-278 item 9 and that no value comes from JSON; parity test `(qa-lead)` | VAL-049, VAL-050 | test output, provenance log |
| QA-33 | Descriptor forbidden-member scan; per-item disclaimer `rg` terminates `exit_code: 1`; no raw-JSON Stop renderer `(qa-lead)` | VAL-051 | test output, `rg` status |
| QA-34 | Every completion report records the three package times plus the slowest ordinary package; none breaches 600/900 s `(qa-lead)` | VAL-052 | report excerpts |
| QA-35 | Security review: observation under the write lock and never on replay; path containment and symlink escape; digest-before-observe; sensitive fields never echoed; no top-level sensitive descriptor field (I-17); import-binding equivalence limited to the bound root `(qa-lead)` | no VAL | findings with file references |
| QA-36 | Every completion report lists every exact task-spec command with token and observed termination; declared non-zero commands (tasks 9, 13, 14) recorded as declared; converged rows carry introducing commit, invalidating change and executing command (DEC-162/165/166) `(qa-lead)` | no VAL | report-to-spec comparison |
| QA-37 | `git diff 304aa45 -- go.mod go.sum` empty, or each dependency recorded per the budget; no new package outside `internal/*` boundaries in `02_TRD.md` `(qa-lead)` | no VAL | diff output |
| QA-38 | Registries R1-R20 and the section 1 block regenerated; R2 lists only the shared helper; R5 terminates `exit_code: 1`; R19 matches exactly four family fallbacks; R20 matches accounted sites with no missing fallback; deltas explained `(qa-lead)` | no VAL | command output |
| QA-39 | Spec amendments named by tasks 1-14 applied (03, 05, 06, 07, 08, GLOSSARY); spec sync inputs for wrap-up `(qa-lead)` | no VAL | spec diffs |
| QA-40 | Every window review is `PASS`; each earlier `REPAIR` closed by a fresh independent re-review; the task-11 Human walk-through recorded `(qa-lead)` | no VAL | review files |
| QA-41 | Per observing mutation, rerun the unobservable-worktree and deleted-file retries and their new-key controls; replays append nothing; controls fail with the VAL-066 codes `(qa-lead)` | VAL-066 | test output, line counts, control outputs |
| QA-42 | `attention resolve --help`, `stop raise --help`, `stop resolve --help`: no Run wording; Attention wording present `(qa-lead)` | VAL-059 | help captures |
| QA-43 | For every `--input` leaf from a non-project temp dir: `--example` exit 0, empty stderr, tree snapshot digest unchanged `(qa-lead)` | VAL-060 | per-leaf statuses, snapshot digests |
| QA-44 | Close and abandon with a mismatching expectation (`repository-mismatch`) and with unmerged paths (`recovery-state-conflict`); journal line count unchanged `(qa-lead)` | VAL-061 | outputs, line counts |
| QA-45 | Rerun the two VAL-062 fixtures; closed Run status `succeeded` and `failed` `(qa-lead)` | VAL-062 | test output, SQL rows |
| QA-46 | Rerun the whole-tree grammar corpus and its completeness check; binary probes of one unknown flag, one missing required flag, one malformed typed value (e.g. `issue list --has-concerns=maybe`), one malformed enumerated value (e.g. `issue list --order bogus`) and one extra positional, both modes `(qa-lead)` | VAL-064 | test output, captures |
| QA-47 | `fid bogus`, `fid isue`, and an unknown word under every command group, both modes: `operation-unavailable`, DEC-279 `operation`, sorted `allowed`, exit 2, nonempty Human stderr `(qa-lead)` | VAL-065 | captures per group |
| QA-48 | `rg -n "PRIVATE_" cmd/fid/integration` terminates `exit_code: 1` `(qa-lead)` | VAL-063 | `rg` status |
| QA-49 | Rerun the fallback corpus and its R20 completeness check, the 512/513, empty/whitespace and escape-cut boundaries, and an activation-component sentinel guard that reaches the fallback; compare both streams and unchanged failure metadata; converge plain-copy causes to the TL-dispositioned RC-01 rule `(qa-lead)` | VAL-067, VAL-027-029 messages | per-site captures, unit results, exclusion reasons, retained-test provenance |

## 8. Mid-Phase Review Cadence

DEC-278 item 13 accepted five windows. DEC-280 item 10 delegates the expanded
window D cadence to the QA Lead. This recheck splits it into D1 and D2: two
tasks each, with grammar reviewed before task 10 changes the renderers. Six
windows now cover all 15 tasks; the task-11 Human checkpoint is preserved.

| Window | After task | Reviews | In-window focus |
|---|---|---|---|
| A | 3 | tasks 1-3 | VAL-001-009, 058 and task 2's rows of 066 (migration, replay fail-closed, schema 19, benchmark, field digests); VAL-010-016, 059, 060 framework |
| B | 6 | tasks 4-6 | Observation, SourceInput, expectations and the same-worktree predicate before tasks 7-8 reuse them (VAL-008 Source case, 009, 011-029 for these leaves, 035, 060, 066). Fresh feasibility for the unmerged-path, plain-copy and import-binding fixtures that tasks 7-8 prescribe |
| C | 8 | tasks 7-8 | DEC-273 Runs and gate runs, real freshness (VAL-009, 011-019, 022-035, 060-062, 066 for these leaves). Fresh feasibility for every non-default command and corpus that tasks 9-14 prescribe |
| D1 | 9a | tasks 9, 9a | Whole-tree descriptors (VAL-011-016, 060 complete, 035), grammar and unknown commands (064, 065), R17-R19 completeness and control escaping, before task 10 adopts their messages. |
| D2 | 11 | tasks 10, 11 | Failure lines (036-039), generic fallbacks (067 and 027-029 messages), Stop and repository lines (040, 041). Recheck grammar/Human pairing after task 10. **Human walk-through** of new failure and Stop output (VAL-054, VAL-055 interim), recorded in `qa/human-walkthrough-after-task-11.md` from `TestP24HumanFixtures` output regenerated by the window QA Reviewer at the task-11 close commit (section 5) |
| E | 14 | tasks 12-14 | `get` leaves, judgment views, journey, parity, absence (VAL-042-052, 063) |

- Each window uses `project-docs/templates/mid-phase-review-template.md` and
  the next unused `qa/review-M.md`. It reruns the in-window named tests and
  verifies in-window `VAL-NNN` in both states. It compares each report's
  commands with its spec. It reruns inventory registries R1-R20 (inventory
  header). It checks VAL-052 wall times for every in-window task, and records
  gate evidence or why a gate was not run.
- The window-D2 walk-through records Human judgment only. Its interim results
  do not replace the phase-end rows HUMAN-02 and HUMAN-03.
- Each review records exactly one outcome: `PASS`, `REPAIR` or `BLOCK`.
  `REPAIR` names the owning task, the affected assertion, the exact repair
  boundary, the required evidence and the fresh re-review scope. The window
  does not advance. After repair, a fresh independent QA Reviewer writes the
  next unused review path, cites the prior review and verifies each original
  finding.
- `BLOCK` is reserved for protected authority, work outside the reviewed
  task/window scope, or an exhausted no-progress budget. QA records a blocker
  candidate and never allocates an `ESC` ID.

## 9. Human Verification Scope

Fixtures come from `TestP24HumanFixtures` (task 11, extended by task 14;
DEC-278 item 10). The QA Lead regenerates them at the task-14 close commit
with the section 5 command and writes no test code. The Human records
judgment only. The HUMAN-01 step list is re-derived at final review from the
task-14 journey's recorded steps. Any added or removed step keeps the parent
VAL-053. If RR-06 is dispositioned as an accepted consequence, HUMAN-02 gains
a sixth row for a generic-fallback denial (the plain-copy `claim acquire`
refusal).

| ID | Check | Covers | Human Evidence Expected |
|---|---|---|---|
| HUMAN-01a | From the refused completion step's output alone, can you say why it was refused? `(human)` | VAL-053 | answer and quoted output |
| HUMAN-01b | From that output alone, can you name the next command to run? `(human)` | VAL-053 | answer and quoted output |
| HUMAN-01c | Using only `--help` and `--example`, can you record the Decision Reference? `(human)` | VAL-053 | yes/no, attempts |
| HUMAN-01d | Using only `--help` and `--example`, can you record the approving Human Outcome? `(human)` | VAL-053 | yes/no, attempts |
| HUMAN-01e | From fid output, can you tell that the blocking Stop is resolved? `(human)` | VAL-053 | answer and quoted output |
| HUMAN-01f | Can you retry the refused step and complete the Envelope? `(human)` | VAL-053 | observation |
| HUMAN-01g | Did any step need source, specs, JSON internals or `.fid/`? `(human)` | VAL-053 | yes/no with the step named |
| HUMAN-02a | `idempotency-conflict`: state the reason and the next action from the output alone `(human)` | VAL-054 | answer |
| HUMAN-02b | `repository-mismatch`: state the reason and the next action `(human)` | VAL-054 | answer |
| HUMAN-02c | `invalid-argument` on an `--input` leaf (nested field): state the reason and the next action `(human)` | VAL-054 | answer |
| HUMAN-02d | `protected-decision-required` (with `blocked by:`): state the reason and the next action `(human)` | VAL-054 | answer |
| HUMAN-02e | `source-digest-mismatch`: state the reason and the next action `(human)` | VAL-054 | answer |
| HUMAN-03a | Do you accept the `stop list` view? `(human)` | VAL-055 | accept/reject with reason |
| HUMAN-03b | Do you accept the `attention list` view? `(human)` | VAL-055 | accept/reject with reason |
| HUMAN-03c | Do you accept the `resume` view? `(human)` | VAL-055 | accept/reject with reason |
| HUMAN-03d | Do you accept the `review get` view? `(human)` | VAL-055 | accept/reject with reason |
| HUMAN-04a | Does every judgment view present labelled sections? `(human)` | VAL-056 | accept/reject per view |
| HUMAN-04b | Does every judgment view state its purpose? `(human)` | VAL-056 | accept/reject per view |
| HUMAN-04c | Does every judgment section carry its disclaimer once? `(human)` | VAL-056 | accept/reject per view |
| HUMAN-04d | Is every judgment view free of routine-state narration? `(human)` | VAL-056 | accept/reject per view |
| HUMAN-04e | Is recorded reasoning readable in every judgment view? `(human)` | VAL-056 | accept/reject per view |
| HUMAN-05a | Author a `run close` input with `expected_repository`, using only `--help`, `--example` and the git commands the help names; accepted on the first or second try? `(human)` | VAL-057 | attempts, final input |
| HUMAN-05b | Author a `decision record-reference` input with the digest omitted; accepted on the first or second try? `(human)` | VAL-057 | attempts, final input |
| HUMAN-05c | Author an input for one more leaf of your choice; accepted on the first or second try? `(human)` | VAL-057 | leaf, attempts, final input |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Every required deterministic gate in section 5 passes with observed
  termination evidence, including the three fresh 100k samples, the
  archive-scale gate and the full-`checkptr` control.
- Every `VAL-NNN` (all 66, plus any row added under section 12.3) has planned
  evidence and a final-review outcome slot. Every outcome is `PASS`, `FAIL`,
  `BLOCKED` or `CARRY-FORWARD` with cited evidence.
- Every `(qa-lead)` item has reproducible evidence captured in the
  final-review session.
- Human checks are judgment-only and mapped to `human-checklist.md`.
- Every window review (A-E) is `PASS`; each earlier `REPAIR` has a fresh
  independent `PASS` re-review. The window-D walk-through is recorded.
- No unresolved Must-fix finding or blocker candidate remains.
- The spec amendments named by tasks 1-14 are applied and named for wrap-up
  sync, across all seven specs plus the glossary (Heavy).

An in-scope failed gate or incomplete evidence triggers repair. Recording
`BLOCKED` does not make it acceptable. A protected missing or contradictory
contract prevents final review and is emitted as a blocker candidate for
TL-Lite or the Tech Lead, which alone may allocate an `ESC` identity.

## 11. Readiness Conditions

First-pass conditions 1 and 2 are met: DEC-278 dispositioned R-01..R-17 and
I-01..I-19, DEC-279 decided D5, and the edits are applied (verified in 12.1).
Condition 3 is this refresh. The plan becomes Ready only after:

1. The Tech Lead dispositions RR-01..RR-09 (section 12.3). For RR-06, the
   Tech Lead records whether accepted text decides it or a DEC with Human
   approval is needed. RR-01 is decision-free if task 10 is aligned to
   contract section 9; restricting section 9 instead needs a DEC.
   Improvement items RI-01..RI-13 are applied or declined.
2. The Tech Lead applies accepted edits to the contract, the validation
   contract and the affected task specs (at least tasks 4, 9 and 10).
3. A fresh QA Lead check revalidates every item those edits change, in a new
   subsection of section 12. It reuses nothing that changed.
4. The Phase Runner's readiness check passes.

## 12. Readiness Refresh (Second Pass)

Performed 2026-10-04 at `21c833f` (clean tree, branch `phase-24`). It
covers the changes in `c80304c` (DEC-278) and `21c833f` (DEC-279). Each
first-pass disposition was checked against source or a fresh probe, not
against the disposition register.

### 12.1 First-Pass Disposition Verification

| Item | DEC-278 disposition | Applied in | Fresh evidence (this session) | Result |
|---|---|---|---|---|
| R-01 | Map is an insert value on every keyed insert and a **guard** on the `gate.receipt-recorded` update (item 7) | task 1 Binding Constraints and Scope; inventory R3 | R3 `exit_code: 0`, 16 files. The task-1 constructor `rg` finds exactly `internal/journal/idempotency.go:218` and `internal/evidence/execution.go:456` (insert) and `:757` (update, guard `request_digest`). `validateSemanticUpdate` (`internal/export/archive_updates.go:28-41`) iterates `effect.Values` only, so a new guard leaves archive validation unchanged. The R3 "read, unchanged" sites are all SELECTs: `workflow/activation.go:845`, `maintenance/operational.go:128`, `journal/current_event_stop_references.go:37`, `current_event_result_projection.go:27`. The only tests that insert rows directly are the two task 1 names (`maintenance/check_test.go:344` `seedTask21CorruptProjection`; `db/operational_kernel_migration_test.go`) | Resolved |
| R-02 | Per-keyed-leaf corpus over R16 with a completeness test | VAL-002/006/007; task 2; inventory R16 | 50 keyed leaves; 49 after `export import` across 20 families. All idempotency scopes are `<operation>:<project_id>` with no actor component, so an `--actor` row reaches a conflict | Resolved (RI-01) |
| R-03 | New VAL-066 with a new-key control | VAL-066; tasks 2, 4-8 | The seven `internal/review` observation calls run unconditionally after validation and authority (`decision.go:129`, `human_outcome.go:114`, `requirement.go:94, 145`, `deviation.go:93`, `concern.go:96, 163`), so a valid control reaches observation | Resolved (RI-04) |
| R-04 | 18 observing leaves; `evidence.evaluate` returns `repository-unavailable` (item 1) | contract 2, 6; Shared Terms; VAL-009/017/022/023 | Live observation sites: the eight review leaves; `envelope.go:448` (complete); `run_close.go:291`; and exactly the six excluded leaves: `recovery.go:193`, `run_checkpoint.go:106`, `evidence/artifact.go:452`, `evaluation.go:463` (terminal-success transition), `resume.go:90`, `run_inspect.go:41`. `evidence.evaluate` does not observe today (task 8 adds it) | Resolved; not a blocker |
| R-05 | Five-member mapping by capture stage (item 2) | contract 6; task 4; VAL-023 | Every `unsafeRepository` cause in `internal/project/repository.go:61-364` falls under one stated rule: toplevel/HEAD/tree → `not-a-repository`; `ls-files -v` flags → `tracked-flags`; `ls-files -u` and stage≠0 images → `unmerged-paths`; `validPath`, non-UTF-8, unsupported type/mode/object, path count, output size → `unsafe-path`; DeepEqual and read failures → `changed-during-capture`. One cause fits two rules (RI-02). The example fixture for unsupported file types is infeasible (RR-09) | Resolved, with RR-09 |
| R-06 | Rows from `cli.RegisteredCodes()` with pinned exclusions | VAL-036; task 10 | `internal/cli/codes.go:399` exports `RegisteredCodes()` | Resolved |
| R-07 | Task 9 delivers grammar failures and `operation-unavailable` (items 3-4); D5 decided A (DEC-279) | contract 9, 13; VAL-064/065; task 9 | DEC-279 recorded and VAL-065 pins it. The measured baseline differs from the record (RR-03, RI-11), and the message and `allowed` rules are not fully determined (RR-02, RR-04, RR-05). Section 9 still drops cause text in the generic fallbacks (RR-06) | Disposition applied; new findings |
| R-08 | VAL-038 non-`--input` row | VAL-038; task 10 | Row present. Task 10 contradicts section 9 for grammar failures on `--input` leaves (RR-01) | Applied; RR-01 |
| R-09 | Counting rule (item 5) | contract 9; VAL-039; task 10 | Identical rule in all three: whole line, label included, after escaping; 200 uncut; 201 → first 200 + `…`; a straddling escape is dropped whole | Resolved |
| R-10 | 10 Condition types, open and resolved `Next:` states (item 8) | contract 10; VAL-040; task 11 | 10 codes, 10 Condition types (`stop.go:175-189`). Resolution needs a real qualifying fact Event per type (`qualifyingFact`, `stop.go:614-641`). Every type already has resolve fixtures in 1-8 test files each | Resolved (window D confirms) |
| R-11 | Same-worktree and plain-copy fixtures | VAL-041; task 11 | `resume --repository` and `run inspect --repository` exist. A `cp -a` plain copy is a working independent workspace (`issue create` in the copy: `created issue PC-b43c67 (task) "two"`, rc 0) | Resolved |
| R-12 | Exact fallback refusals (item 6) | contract 5; VAL-027/028/029; tasks 6-8 | `envelope.go:255` `invalid envelope input`; `claim.go:234` `invalid claim request`; `run.go:321` `invalid Run request`, each `allowed:["accepted <op> input"]`; `gate.go:329` `invalid gate request`, `allowed:["accepted gate registration input"]` (default branch for every gate operation). Only `review.record` maps `repository-mismatch` today (`review.go:225-227`), so no earlier branch intercepts the refusal | Resolved |
| R-13 | Five provenance classes (item 9); help names the git commands | VAL-012/049; tasks 5, 14 | fid's `dirty` is staged ∪ unstaged ∪ untracked-not-ignored, excluding `.fid` (`repository.go:138-169`). That is exactly "`git status --porcelain` lists a path outside `.fid/`" (task 5) | Resolved (RI-08) |
| R-14 | `TestP24HumanFixtures` in task 11, extended by task 14 (item 10) | tasks 11, 14 | `.cache/phase-24-human` is ignored (`.gitignore:51:/.cache/`). The command form is valid: the same flags with `-list '^TestP23HumanFixtures$'` print `TestP23HumanFixtures`, `ok ... 1.686s` (terminated `exit_code: 0`). The env-root-else-`t.TempDir()` pattern exists (`p23_human_fixtures_test.go:63-65`) | Resolved |
| R-15 | Task 13 adds `operations.json` and its reader | task 13; inventory R13 | `rg -n "PRIVATE_" cmd/fid/integration` `exit_code: 0` in exactly three files: `p23_human_leaf_pairs_test.go` (7), `p23_human_fixtures_test.go` (1), `testdata/p23_conformance/operations.json` (4). All three are in task 13's scope | Resolved |
| R-16 | Declared `exit_code: 1` | task 9 | 7 decoders; callers are all `--input` families (14 files); test references (`envelope_test.go:341`, `finding_test.go:21-27`) converge in tasks 6/5 | Resolved (RI-05) |
| R-17 | Split into VAL-058..063 (item 12) | contract | Each new row has one observable and both states | Resolved |
| I-01, I-02 | Recursive walk that exits 0; corrected counts | inventory 1, R15 | Block terminates `exit_code: 0`; 92 leaves, 39 `--input`, 50 keyed, 20 without a Short (8 read + 12 `--input`) | Resolved (RI-06) |
| I-03..I-06 | Applied | Shared Terms; VAL-007, 003, 004 | Row text present | Resolved |
| I-07 | Record only | VAL-005; task 1 | Note present | Recorded |
| I-08, I-09, I-11, I-12 | Applied | Shared Terms; VAL-015, 018, 035, 044; tasks 3, 4, 9, 13 | Row text present | Resolved |
| I-10, I-16 | Statuses and names declared | tasks 10, 2 | Present | Resolved |
| I-13, I-14, I-15 | Applied | VAL-052 and all tasks; task 1 note; task 13 | Present | Resolved |
| I-17 | No top-level sensitive field (item 11) | task 9 | `environment` and `redact_literals` exist only in `GateDefinition` (`internal/evidence/policy.go:20, 24`). That type reaches input only through activation `content` | Resolved |
| I-18, I-19 | Phase brief | brief Known Risks | Present; `/tmp` free space has since fallen to 1.2 GB (RI-13) | Resolved |

### 12.2 Fresh Command And Feasibility Validation

Working directory `/home/mark/projects/fid`. Probes used a scratch build of
`21c833f` against fresh projects in session scratch storage.

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Product tree identity `304aa45` → `21c833f` | `git rev-parse <rev>:<path>`; `git diff --stat 304aa45 HEAD -- . ':!project-docs'` | PASS | `cmd` `072971c9…`, `internal` `fb9d687b…`, `migrations` `9ef9c27e…`, `go.mod` `9825eba4…`, `go.sum` `b9915dbe…` identical; diff empty. The section 6 gate results stand |
| Inventory section 1 block, verbatim under `bash` | executed | PASS | terminated `exit_code: 0`; 109 stdout lines = 39 + 50 + 20, concatenated (RI-06); empty stderr |
| `--input`, R15, R16 lists | the same loops, outputs split | PASS | 39 (identical to the section 1 rows); 92 leaves (6 at depth 1, 82 at depth 2, 4 at depth 3); 20 without a Short; 49 keyed leaves in 20 families after `export import` |
| R1, R2 (task 2) | executed | PASS | `exit_code: 0`; 43 sites in 32 files; 11 helpers |
| R3 (task 1) | executed | PASS | `exit_code: 0`; 16 files |
| `rg -n "EventIdempotency\{\|Table: \"idempotency_results\"" internal -g '!*_test.go'` (task 1 scope) | executed | PASS | `exit_code: 0`; `idempotency.go` (2), `execution.go` (4) |
| R4 | executed | PASS | latest `0018_p23_remove_v1_judgment`; `ActiveMigrationVersion = 18`; `ProjectionSchemaVersion = 5` |
| R5 (task 9, declared `exit_code: 1` at task close) | executed | PASS (syntax; exit 0 expected before task 9) | `exit_code: 0`; 7 decoders |
| R6, R7, R9, R10, R14 | executed | PASS | 11 functions; 8 files; 41 writers; 32 assertions in 13 files; `contract_surface_test.go:128` |
| R11 (task 11) | executed | PASS | `exit_code: 0`; `review.go`, `resume_human.go`, `run_recovery.go`, `p23_human.go`, `run_inspect_human.go`, `testdata/p21_conformance/manifest.json` |
| Tasks 13/14 absence lines | executed | PASS (syntax; exit 0 expected before task 13) | disclaimer `rg` `exit_code: 0` at `p23_human.go:163, 167, 176`; `PRIVATE_` `rg` `exit_code: 0` in three in-scope files |
| `rg -n '"argument": *"request"' cmd/fid -g '!*_test.go'` (proposed for task 9, RR-05) | executed | PASS (syntax) | `exit_code: 0`; 7 sites: `canonical_request.go:47`, `root.go:33`, `link_response.go:89`, `issue_mutation_response.go:123, 236`, `issue_get_list.go:129`, `receipt.go:87` |
| Human fixture regeneration (task 11 note) | `-list` with the same flags against the P23 analogue; `git check-ignore -v` | PASS | see R-14 above |
| VAL-023 unsupported-file-type fixture | scratch repository | **untracked FIFO infeasible (RR-09)**; tracked-file variant PASS | untracked FIFO: `git ls-files --others --exclude-standard -z` and `git status --porcelain` both empty. Tracked `a.txt` replaced by a FIFO: `git diff --name-only` lists `a.txt`; porcelain ` M a.txt` |
| VAL-023 tracked flag and no-commit fixtures | source | PASS | `ls-files -v` lowercase/`S` → flagged (`repository.go:292-303`); unborn HEAD → `missing base` (`:75-78`) |
| Plain copy (VAL-027..029, 041) | `cp -a` of a fid project | PASS | see R-11 above |
| VAL-062 close status | source | PASS | `run_close.go:81-87`: `failed` unless Obligation evaluation returns no reasons |
| VAL-059 baseline | `--help` scan of the 39 `--input` leaves | PASS | "Run JSON input" on the seven Run leaves plus `attention resolve`, `stop raise`, `stop resolve` (all task 3) |
| Unknown word under every command group, both modes (VAL-065 baseline) | built-binary probe of 28 groups and the root | **baseline differs from DEC-278 (RR-03)** | 23 groups (`activation`, `activity`, `artifact`, `attention`, `claim`, `decision`, `deviation`, `envelope`, `evidence`, `finding`, `gate`, `human-outcome`, `issue label`, `journal`, `label`, `link`, `maintenance`, `receipt`, `requirement`, `review`, `stats`, `stop`, `workflow`) print help on stdout and exit 0 in both modes. `concern`, `export`, `issue`, `issue ac`, `run` and the root exit 2 with no output |
| Cobra unknown-command text (RR-02) | module source `cobra@v1.10.2` | noted | `args.go:36` appends `findSuggestions`; `command.go:781-795` emits `\n\nDid you mean this?\n\t<cmd>\n` unless `DisableSuggestions`; fid's root sets none (`root.go:61-73`) |
| Grammar baseline (VAL-064) | built-binary probes, both modes | noted (RR-04, RR-05, RI-11) | `issue list --bogus` and `issue show` (no ID): rc 2, silent in **both** modes. `issue list --has-concerns=maybe`: **rc 7**, silent in both modes (pflag `invalid argument` is not matched by `isCobraInputError`, `root.go:45-52`). `ready extra`, `concern list extra`: rc 2, silent (`cobra.NoArgs` text `unknown command …`). `stop raise --bogus`: JSON `operation:"raise"`. `capabilities extra` already conforms |
| Human versus JSON message (RR-05, RR-06) | probes; source | noted | `issue update --priority 99`: Human `error: update issue: priority must be between 0 and 5`, JSON `invalid Issue mutation request`/`argument:"request"`. `issue list --order bogus`: Human `unsupported issue list order "bogus"; supported orders: …`, JSON `invalid issue read argument`. 25 `CodeInvalidArgument, "invalid …"` fallback sites in `cmd/fid`; the Human branch prints the cause (`renderIssueError`, `issue_helpers.go:43`) |
| Post-frame `gate run` failure (RR-07) | source | noted | Human mode renders live output to stdout (`gate.go:195`); `gate-unsuccessful` failure built at `gate.go:138` |
| Root `--help` subcommands (RI-12) | executed | noted | lists `completion` and `help` among 34 names |
| Validation-contract cross-check | script over the contract and 14 task specs | PASS | 66 rows, all `Ready`; validators QA Reviewer 58, QA Lead 3, Human 5; every ID cited by at least one task; only out-of-contract citation is Phase 23 VAL-073; Coverage table omits 035, 051, 052, 053 (RI-07) |
| `/tmp` | `df -h /tmp /` | CONSTRAINT | 3.9 GB tmpfs, 1.2 GB free; `/` 896 GB free; no stale `fid-integration-binary-*` |

### 12.3 New Revisions Requested From Tech Lead

Must-fix (M) items block Ready. Improvement (I) items are applied at the
Tech Lead's discretion. The Tech Lead owns every edit to the contract, the
validation contract, the inventory and the task specs. This plan proposes
edits and never applies them.

| ID | Sev | Where | Finding and proposal |
|---|---|---|---|
| RR-01 | M | task 10 Binding Constraints (`next:`); contract 9 next-table; VAL-038 | **Task 10 contradicts section 9 for grammar failures on `--input` leaves.** Task 10 says "Grammar failures and failures on leaves without `--input` print no `next:`". Contract section 9 applies the row to every "`invalid-argument` on an `--input` leaf", with no grammar exception, and VAL-038 follows the contract. So `fid stop raise --bogus`, and every task 3 `--example` exclusivity failure, has two different expected outputs. Proposal: align task 10 with section 9, so that only leaves without `--input` print no `next:`. Alternatively, amend section 9 by DEC to limit the row to structural `input.<path>` failures. Either way, add a VAL-038 row for a grammar failure on an `--input` leaf in the chosen state. |
| RR-02 | M | task 9 (unknown commands); contract 9; VAL-036, VAL-065 | **The unknown-command `message` is multi-line for a near-miss word.** Task 9 makes the message "Cobra's unknown-command text". At the root, Cobra builds that text with suggestions (`cobra@v1.10.2` `args.go:36`, `command.go:781-795`), and fid sets no `DisableSuggestions`. For `fid isue` the message would be four lines (`unknown command "isue" for "fid"`, blank, `Did you mean this?`, `\tissue`). The Human `error:` line then breaks section 9's "exactly these lines" rule and VAL-036's label check. The group validators already use the single-line form (`concern.go:25`, `issue.go:17`, `issue_ac.go:17`, `run.go:31`, `export.go:24`). Proposal: pin the message at every level as the single line `unknown command "<word>" for "<command path>"`; `allowed` already lists the alternatives. Add a near-miss row (`fid isue`) to VAL-065 and to VAL-036's unknown-command rows. |
| RR-03 | M | task 9 Objective and Scope; DEC-278 item 4; VAL-065 | **The unknown-command baseline and file scope are wrong for 23 of 28 groups.** DEC-278 item 4 and task 9 say both modes are "silent today". Measured: an unknown word under 23 groups prints the group help on stdout and **exits 0** in both modes, so `fid --json activation bogus` returns Human text as success. Only `concern`, `export`, `issue`, `issue ac`, `run` and the root exit 2 silently (12.2). Task 9's file list names `root.go`, `main.go`, `canonical_request.go`, `codes.go` and leaf files, but not the 23 group constructors. Proposal: correct the stated baseline. Add to task 9's scope either every group constructor in `cmd/fid`, or one central mechanism in `root.go` applied to every command with subcommands. Phase 21 section 4 (exit 2) already authorizes the exit change. |
| RR-04 | M | task 9 grammar table; VAL-064 | **The grammar `allowed` values and class boundaries are not fully determined.** (a) Unknown flag: "the leaf's accepted long flag spellings" does not say whether the globals `--json` and `--project` are included, and they are enabled per leaf (`rootFlagCapabilities`, `root.go:21, 134`). Phase 21 rows list only leaf-local flags; the globals sit outside the row (Phase 21 line 18). Proposal: the row's local flags only. (b) Malformed value: `[<the flag's type or format>]` differs from Phase 21's "row's exact accepted alternatives" for enumerated flags (`issue list --order <created\|priority\|sequence\|updated>`, Phase 21 line 333). Proposal: an enumerated flag lists its members; a typed flag lists its pflag type name (`int`, `bool`, …). (c) Typed-flag parse failures are a path the table does not name. `issue list --has-concerns=maybe` exits **7** and is silent in both modes today. Proposal: name them in the malformed-value class. (d) A `cobra.NoArgs` leaf reports an extra word as `unknown command "<word>" for "<leaf>"` (`ready`, `concern list`, `activity list`). Proposal: state that a leaf's extra word is a wrong positional count (`invalid-argument`), never `operation-unavailable`, and pin its `message`. |
| RR-05 | M | task 9 Exact Commands; VAL-064 | **The grammar rule is universal, but its proof is sampled (the first-pass R-02 pattern).** VAL-064 takes one row per class per leaf, while seven request-level fallbacks remain (`rg -n '"argument": *"request"' cmd/fid -g '!*_test.go'`, 12.2). Flag checks done by fid itself, not the parser, reach these fallbacks with Human text richer than the JSON message: `issue update --priority 99`, `issue create --actor bogus`, `issue list --order bogus` (12.2). A site that is missed keeps `argument:"request"`. After task 10, it prints the generic text in Human mode, and no assertion fails. Proposal: add that `rg` to task 9 with a declared status: `exit_code: 1` if every request-level fallback is converted, or exit 0 with each remaining match accounted for as non-grammar in the report. Also require one malformed-value row per flag with an enumerated, bounded or typed value, enumerated from the Phase 21-23 rows. |
| RR-06 | M | contract 9; task 10; VAL-036, VAL-054 | **Blocker candidate (12.4): generic domain fallbacks lose the Human cause.** Every family's default failure branch maps unmatched errors to a fixed `invalid-argument` message (25 sites). Human mode prints the cause instead (`writeClaimFailure`, `claim.go:262-268`, via `renderIssueError`, `issue_helpers.go:43`). Section 9 makes the Human `error:` line the JSON message verbatim and keeps JSON unchanged. Concrete case, the VAL-028 refusal: today a worktree that does not match the Envelope's prints `error: claim acquire binding is invalid` (`internal/claims/generic.go:361-362`). After task 10 it prints `error: invalid claim request`, then `reason: invalid-argument; argument=input, allowed=accepted claim.acquire input`, then `next: see fid claim acquire --help; template: fid claim acquire --example`. That points the Human at the input when the worktree is the cause. No assertion detects this: VAL-036 checks form, and VAL-054's five denials include no fallback. Proposal: see 12.4. |
| RR-07 | M | VAL-036 pass criterion; task 10 (`gate run`) | **"stdout is empty" cannot hold for post-frame `gate run` failures.** Task 10 prints the failure lines after the streamed output for a terminal non-success in Human mode. The live output goes to stdout (`gate.go:195`), and Phase 21 section 1 keeps the streaming exception. A `gate-unsuccessful` row therefore has nonempty stdout. Excluding the code would drop a reachable producer. Proposal: VAL-036 states the exception (post-frame `gate run`: stdout carries the streamed Human output, and stderr has the section 9 lines) and keeps a row for it. |
| RR-08 | M | VAL-064 evidence method; task 9 | **VAL-064 says "every live leaf", but task 9 excludes `journal inspect` and `sql`.** Phase 21 section 4 makes those two Human-only diagnostics outside `fid-cli.v1` that reject `--json`. A reviewer applying VAL-064 as written would fail it on both. Proposal: state the exclusion in VAL-064, citing Phase 21 section 4. |
| RR-09 | M | task 4 Acceptance Criteria (VAL-023 cause classes) | **The named unsupported-file-type fixture cannot reach `unsafe-path`.** Task 4 suggests "an untracked FIFO". Git does not list FIFOs as untracked (measured: `ls-files --others --exclude-standard` and `status --porcelain` both empty), so observation succeeds. A tracked file replaced by a FIFO is listed by `git diff --name-only` and reaches `unsupported worktree file type` (`repository.go:247-248`). Proposal: name the tracked-file-replaced-by-a-FIFO fixture. |
| RI-01 | I | VAL-006; task 2 | The `artifact` family has one keyed leaf, and `artifact verify` has a single top-level field (`artifact_id`). Its two-field row can only be `artifact_id` + `--actor`. Say so. |
| RI-02 | I | task 4 reason mapping | `invalid unmerged path` (`repository.go:261-262`) fits both `unmerged-paths` (capture stage) and `unsafe-path` (invalid path). Pin one; capture stage suggests `unmerged-paths`. |
| RI-03 | I | task 4 | An unreadable untracked file (`os.ReadFile` fails, `repository.go:243-245`) maps to `changed-during-capture`, which is retryable, although a retry cannot succeed. `detail` names it. Record only. |
| RI-04 | I | VAL-066; tasks 2, 4-8 | Define the control as "a request, under a new key, that would succeed in an observable worktree in the same pre-state". A same-input control fails on domain checks first for create-style leaves (`concern.go:96` checks the ID before observing) and for already-completed subjects (`envelope.complete`, `run.close`). |
| RI-05 | I | task 9 | The R5 pattern `func decodeClosed[A-Za-z]*Input` also matches a shared `decodeClosedInput`. Name any shared `cmd/fid` decoder outside the pattern, or the declared exit 1 fails (cf. I-16). |
| RI-06 | I | inventory section 1 | The block prints three lists back to back (109 lines). Label or split them so a window reviewer can diff each list. |
| RI-07 | I | validation contract Coverage | Add the roadmap-sourced rows: VAL-035 (roadmap section 7), VAL-051 (removal column), VAL-052 (Phase 23 VAL-079, DEC-227), VAL-053 (journey). |
| RI-08 | I | task 14 journey | Under real freshness, files the caller writes inside the worktree change the dirty digest and can stale Receipt-backed Evidence. That covers the Decision source and input files built from `--example`. `evidence-stale` has no next-table row, so the journey could reach a step with no `next:` at the last task. Design the start state so this cannot happen. For example, write input files outside the worktree, write Sources under a path the start state's `.gitignore` already ignores, or guard completion with no Receipt-backed Evidence. |
| RI-09 | I | task 14 First Failing Test | The named first step (a `next:` line naming `fid decision record-reference`) is delivered by task 10, so it is likely green when task 14 starts. Use an assertion that only task 14 makes true (for example the fixture manifest listing the start state). |
| RI-10 | I | phase brief Task Index | Task 12 lists dependency 4, but its acceptance criteria use the task 9/10 corpora. Sequential order makes this harmless; align the column. |
| RI-11 | I | DEC-278 item 3; task 9 | `issue list --bogus` is silent in **both** modes, not only JSON, and so is `issue show` with no ID (12.2). Task 9's "where Human mode printed nothing, use the parser's text" rule covers this. Correct the record. |
| RI-12 | I | DEC-279; task 9 | Root `--help` lists Cobra's `completion` and `help`. By DEC-279's literal rule they belong in `allowed` for `fid bogus`. Say so, so a test does not filter them out. |
| RI-13 | I | environment | `/tmp` free space is now 1.2 GB (1.4 GB at first pass). CF-042 risk persists; the brief's figure is stale. |

### 12.4 Blocker Candidates

QA emits candidates only. No `ESC` identity is allocated. TL-Lite or the Tech
Lead decides whether each is resolvable within existing authority.

| Candidate | Protected category | Evidence | Authority boundary | Requested action |
|---|---|---|---|---|
| RR-06 | public interface (Human failure content under contract section 9) | 25 generic `invalid-argument` fallbacks whose Human branch prints the cause today (`claim.go:262-268`; `issue_helpers.go:43`); refusal at `claims/generic.go:361-362`. Section 9 sets `error:` to the JSON message verbatim, and "the JSON failure [is] unchanged" | Changing a JSON message, or the section 9 rule, is contract text (DEC-277). Accepting a Human-visible loss of the reason is product judgment | The Tech Lead states whether accepted text intends this consequence. If it does, record that by DEC and add a generic-fallback denial (the plain-copy refusal) to VAL-054 / HUMAN-02, so the Human judges it. Otherwise, amend section 9 by DEC with Human approval (for example, the fallback `message` carries the domain cause) |

RR-01 becomes a contract question only if the Tech Lead chooses to restrict
section 9 instead of aligning task 10.

### 12.5 Verdict And Dispatch Scope

- **NOT READY.** Nine Must-fix items are open (RR-01..RR-09), one of them a
  blocker candidate (RR-06). Every first-pass item is closed (12.1).
- **Scope of the open items:**
  - none touches tasks 1-3 or their assertions (VAL-001..016, 058..060,
    task 2's rows of 066, 052);
  - RR-09 gates task 4;
  - RR-02, RR-03, RR-04, RR-05 and RR-08 gate task 9 (VAL-064/065);
  - RR-01, RR-06 and RR-07 gate task 10 (VAL-036/038).
- The phase brief holds all dispatch until this refresh passes. DEC-278
  gated a single task (task 9, on D5), so this plan supports the Tech Lead
  lifting the gate for tasks 1-3 by decision. The plan itself stays NOT
  READY until the open items are dispositioned and the changed text passes a
  fresh QA Lead check (section 11).

### 12.6 DEC-280/281 Recheck (Codex, 2026-10-04)

Reviewed `c9acdac81972de5f8388cd346f7b78787dbf113c`
(`docs(phase-24): record D6 generic-fallback messages (DEC-281)`) and
`62dddab02584c4c2e4c122dc1122b95e1944f35a`
(`docs(phase-24): apply readiness-refresh dispositions; split task 9 (DEC-280)`),
against prior review commit `87f107132543b0609c2cf1860cd2f883f5e5724e`.
Working directory for repository commands: `/home/mark/projects/fid`.

Loaded the four current-context documents, AGENTS, quality plan, both invoked
skills, QA-plan and validation-contract templates, workflow sections 2.6,
6.1 and 8, Phase 24 contract/inventory/decisions and changed task specs.
Read Phase 21 transport/grammar/failure clauses and the relevant glossary and
PRD entries. There are 15 task specs and **zero Phase 24 completion reports**;
this is a planning recheck, not phase-end acceptance. No product file changed.

#### 12.6.1 Disposition Verification

These results verify planning changes; they do not mark unimplemented VALs
PASS. Earlier evidence in sections 6 and 12.1-12.5 remains attributable to
its original session. No prior gate PASS is re-attested here.

| Finding | Rechecked change and evidence | Disposition |
|---|---|---|
| RR-01 | Task 10 Binding Constraints and AC, section 9 next-table and VAL-038 all include grammar failures on `--input` leaves; the non-`--input` state excludes `next:` | Closed |
| RR-02 | Task 9a and VAL-065 pin one-line unknown-command text without suggestions; task 10/VAL-036 include `fid isue`. The long-message case is separately RC-02 | Original finding closed |
| RR-03 | Task 9a owns the central root mechanism and five old validators; live R17 probe reproduces 23 groups giving help/exit 0 and root plus five groups silent/exit 2, in both modes | Closed |
| RR-04 | Task 9a's class table specifies local flags, declared domains, typed parse/missing-value cases, root/group failures and leaf positional errors; Phase 21 sections 1, 2, 5 and 7 references are correct | Closed |
| RR-05 | Task 9a/VAL-064 require every R18 flag, cited string domain and non-project repetition. R19 has seven current sites and explicitly expects the four residual family fallbacks after implementation; CF-047 names combination/flag-set deferrals | Closed |
| RR-06 | D6=A is recorded in DEC-281 and contract sections 5, 9, 13; task 10 owns message convergence and VAL-067. R20 classification and two retained literal-message assertions verified. One newly pinned reason is incorrect: RC-01 below | Decision resolved; RC-01 remains |
| RR-07 | Task 10 and VAL-036 expressly preserve post-frame Human stdout and compare the JSON terminal-frame failure; `gate.go` uses stdout for the live sink | Closed |
| RR-08 | Task 9a and VAL-064 exclude `journal inspect`, `sql`, and generated Cobra commands; the latter remain in root `allowed` as specified | Closed |
| RR-09 | Task 4/VAL-023 now use a tracked file replaced by a FIFO. Fresh Git fixture lists that file and omits an untracked FIFO | Closed |
| RI-01 | Task 2/VAL-006 name `artifact_id` plus `--actor` for the two-field artifact row | Closed |
| RI-02, RI-03 | Task 4 pins `invalid unmerged path` to capture-stage `unmerged-paths`; unreadable-file behavior is explicitly recorded, not silently changed | Closed / recorded |
| RI-04 | VAL-066 and tasks 2/4 require controls that would succeed from the same pre-state when observable, with fresh record/subject identities as needed; tasks 5-8 cite the shared assertion | Closed |
| RI-05 | Tasks 3 and 9 prohibit shared decoder names matching R5, preserving the declared absence proof | Closed |
| RI-06 | Executed the inventory section 1 block verbatim: six separate files, labelled by pathname, and final `252 total` | Closed |
| RI-07 | Coverage now includes roadmap sources for 035, 049, 051, 053 and Phase 23/DEC-227 for 052 | Closed |
| RI-08, RI-09 | Task 14 excludes Receipt-backed Obligations from the start state, writes input files outside the tree, retains contained Source files, and makes the new manifest entry its first red assertion; CF-040 records stale-Evidence guidance | Closed |
| RI-10 | Task 12 and phase brief both depend on 4 and 10 and use task 9a grammar | Closed |
| RI-11, RI-12 | DEC-280 corrects `issue get` versus the prior `issue show` example; fresh probes confirm silence in both modes. Root help includes `help`/`completion`, which task 9a/VAL-065 retain in alternatives | Closed |
| RI-13 | Brief records the 0.7 GB measurement at DEC-280. Fresh measurement is 1.4 GB free of 3.9 GB; this is changing capacity, not a gate failure | Recorded |
| TL-01 | Contract section 9, tasks 9a/10 and VAL-036/064 require shared escaping. Fresh `stop raise` ESC probe reproduces the current raw Human control byte | Applied, implementation pending |
| TL-02 | Task 9 retains descriptors; task 9a owns grammar/unknown commands with both full task test gates; task 10 depends on it. QA exercises the delegated cadence choice: D1 after 9a, D2 after 11 (section 8) | Applied |

#### 12.6.2 Command And Feasibility Validation

Exploratory failure probes below deliberately invoke refused requests and
record their actual exit 2 or 7. They are observations, not task-gate PASS
claims with undeclared non-zero expectations. Every subprocess was waited for;
its actual return code, stdout and stderr were captured. The build's nested
session was polled to `exit_code: 0`.

| Command or step | Validation method | Result and observed evidence |
|---|---|---|
| `go build -o /tmp/fid ./cmd/fid` | Execute exact AGENTS build command | PASS: `exit 0 (no output)` (terminated `exit_code: 0`) |
| Inventory section 1 bash block | Execute verbatim from repository root | Feasible: `252 total` (terminated `exit_code: 0`); `/tmp/tmp.DulWIudkuG/{leaves,groups,input,keyed,no-short,typed-flags}` contain 92, 29, 39, 50, 20, 22 lines respectively. Keyed minus `export import` is 49. R18 is 17 int + 5 bool |
| R17 unknown words | Fresh binary under root/all 28 groups, both modes, cwd `/tmp` | 58 observed terminations: 46 help responses (`exit_code: 0`), 12 empty responses (`exit_code: 2`), matching DEC-280 |
| R19 exact task 9a `rg` | Execute `rg -n '"argument": *"request"' cmd/fid -g '!*_test.go'` | Seven sites (terminated `exit_code: 0`), exactly inventory R19; the final four-match condition is future task evidence |
| R20 exact task 10 `rg` | Execute `rg -n 'CodeInvalidArgument, "invalid [^"]*"' cmd/fid -g '!*_test.go'`, inspect each enclosing mapper and also the broader `CodeInvalidArgument` search | 25 matches (terminated `exit_code: 0`); 18 final family branches, three grammar branches, four typed argument branches. No extra family fallback in the broader search |
| Retained literal messages | Extract all 18 R20 fallback strings and scan `cmd`/`internal` `*_test.go` | Exactly `generic_run_test.go:404` and `issue_update_test.go:331`; counting script terminated `exit_code: 0`. Task 10 also owns later task 6-8 message convergence |
| Task 9a format/lint/test/coverage/integration command syntax | Fresh `go help testflag`, `go help build`, `golangci-lint run --help`; `go tool cover -help`; compare exact task commands with section 5 | Help documents race, gcflags, tags, timeout, coverprofile, func, run/bench options; first three help calls terminated `exit_code: 0`; cover help terminated `exit_code: 2` and printed usage. The latter is a usage observation, not a failed coverage run. Scopes `cmd/fid`, `internal/cli`, and tagged integration exist; no command placeholder in task 9a |
| Grammar examples and TL-01 | Built-binary probes from `/tmp`, each in both modes | `issue list --bogus` and `issue get`: empty streams, `exit_code: 2`. `issue list --has-concerns=maybe`: empty streams, `exit_code: 7`. ESC in `stop raise --bo<ESC>gus`: Human raw ESC, `exit_code: 2`; JSON has fixed `invalid command request`, operation `raise`, `exit_code: 2` |
| Grammar message bound, RC-02 | Run `/tmp/fid stop raise` with one unknown long flag consisting of `--` plus 600 ASCII `x` characters; count the parser message without Human label/LF | Parser text has 616 scalars, observed process `exit_code: 2`, empty stdout. Task 9a's pinned root unknown-command formula for a 600-character word has 628 scalars. Both exceed inherited `message:string[1..512]` |
| Tracked FIFO, VAL-023 | Fresh committed `a.txt`, replace with FIFO; also create untracked FIFO; run Git inventory | `git diff --name-only`: `a.txt`; `git status --porcelain`: ` M a.txt`; `git ls-files --others --exclude-standard`: empty. Each terminated `exit_code: 0`. Scratch `/tmp/fid-p24-fifo-cmts92mh` |
| Plain-copy causes, VAL-027..029 | Fresh CLI-built activation/Registration/Envelope/Claim/Run; whole-directory copies including `.git`; submit destination-root repository identity through today's caller-input format | Three expected reasons confirmed, one contradicted; details in 12.6.3. Each refusal in both modes terminated `exit_code: 2`; stdout empty; journal bytes unchanged |
| Sensitive guard, VAL-067 | Stage valid gate component, then same component identity/version with new key and a literal sentinel in `environment[].value`; both modes | Reaches activation's generic `invalid-argument`, exit 2; Human cause contains `UNIQUE constraint failed: activation_components.activation_id, activation_components.component_type, activation_components.component_id, activation_components.component_version (1555)`. Neither stream contains `QA_P24_SECRET_SENTINEL`. This is feasible fallback evidence, not a structural-decoder-only guard |
| Contract/task/artifact counts | Mechanical file-state scan | `contract: 67 Ready: 67 roles: {'QA Reviewer': 59, 'QA Lead': 3, 'Human': 5}`; 15 specs, zero completion reports (terminated `exit_code: 0`) |
| Updated QA-plan coverage and edit integrity | Expand VAL ranges in all task specs and QA/Human checklist rows; check task 9a/10 command blocks for placeholders; `git diff --check` | `VAL IDs: 67 task citations missing: [] QA/Human mappings missing: []`; task 9a has six exact gate/search lines and task 10 seven; scanner terminated `exit_code: 0`. Diff check: `exit 0 (no output)` (terminated `exit_code: 0`) |
| Product tree change since previous QA | `git diff --stat 87f1071 HEAD -- cmd internal migrations go.mod go.sum` | Empty, terminated `exit_code: 0`; identity evidence only, not a re-attestation of prior test results |
| Storage | `df -h /tmp /` | `/tmp` 3.9G total, 1.4G available; `/` 896G available (terminated `exit_code: 0`) |

The earlier full race, integration, coverage, benchmark, archive-scale and
full-checkptr results are **UNVERIFIED in this session**. No phase-end gate is
claimed from them. This recheck validates changed planning commands through
current help/file-state evidence and executes the new behavioral probes; the
full gates remain required at their task/window/phase-end boundaries.

#### 12.6.3 Plain-Copy Reproduction And RC-01

The source-only D6 claim needed a running check. The scratch probe was
`python3 /tmp/fid-p24-recheck-ipmy0o/probe.py` (final attempt terminated
`exit_code: 0`, last token `original-worktree controls: envelope.create,
claim.acquire, run.start accepted`). Full command argv, cwd, input files,
stdout/stderr and subprocess exits are retained in
`/tmp/fid-p24-binding-dm_48t18/receipts.json` and `input-*.json`.
These are session scratch, so the durable reproduction procedure follows.

1. Build the reviewed HEAD with the exact build command above. In a fresh
   Git worktree commit `base.txt`, run `fid init --prefix PQA` and create a
   task Issue. Stage/finalize an activation as in
   `cmd/fid/integration/p22_fake_test.go:314`: one Actor with
   `claim.acquire`, `envelope.activate`, `envelope.create`, `gate.register`,
   `gate.run`, `run.start`; a `/bin/true` Gate Definition; five observed
   activation resources. Write CLI input files outside the worktree.
2. Use the actual project ID, HEAD, tree and SHA-256 of the absolute root
   for today's caller Repository. Register the Gate scoped to the Issue.
   Build a gate-evidence Envelope request from the same helper, with matching
   Source repository blocks and Registration. Copy the whole workspace
   before the next successful original-worktree action at each stage.
3. For `envelope.create`, update all caller Repository worktree IDs in the
   copy's request to the copy-root hash. For `claim.acquire`, copy an active
   Envelope before Claim acquisition; use mode `mutation`, a nonempty path
   scope, and `ttl_seconds:3600`, updating caller `worktree_id` in the copy.
4. Acquire the Claim in the original workspace. Copy it with that still
   active Claim before `run.start`. Keep Claim ID/version, Issue, Envelope,
   Actor and executor valid; change only the request Repository worktree ID
   to the copy-root hash. After successful original `run.start`, similarly
   copy before `gate.run`, retaining its Run/Registration/Obligation IDs.
5. Run each refusal once with `--json` and once without, under fresh keys.
   Hash every `.fid/journal/*.jsonl` before and after. Original-worktree
   create/acquire/start succeed; all eight copy calls below fail with empty
   stdout and unchanged journal bytes. This proves current binding-check
   reachability using current inputs, not the unimplemented Phase 24 input
   shape or future import-equivalence behavior.

| Leaf | Observed Human stderr (LF after text) | Observed JSON message | Exit in both modes |
|---|---|---|---|
| `envelope.create` | `error: prior registration binding does not match envelope` | `invalid envelope input` | 2 |
| `claim.acquire` | `error: claim acquire binding is invalid` | `invalid claim request` | 2 |
| `run.start` | **`error: run claim binding is invalid`** | `invalid Run request` | 2 |
| `gate.run` | `error: run gate: repository binding is invalid` | `invalid gate request` | 2 |

Every JSON failure has `code:invalid-argument`, `category:invalid-input`,
`retryable:false`, `subject:null`, `argument:input`, `value:null`. Allowed
arrays are respectively `accepted envelope.create input`,
`accepted claim.acquire input`, `accepted run.start input`, and
`accepted gate registration input`.

**RC-01 — Must-fix; specification/public-interface blocker candidate.**
VAL-029, contract section 5, DEC-281 consequences, task 7's plain-copy
constraint and task 10's message constraint pin `run envelope binding is
invalid`. At `internal/workflow/run.go:228-229`, the earlier Claim check
compares the Claim worktree with the request repository and returns
`run claim binding is invalid`. The Envelope check is later, at 235-236.
Task 7 also requires the Claim check to use the same-worktree predicate;
in a plain copy with no import binding that check will still fail first.

The pinned text does exist, but a separate baseline probe reaches it only
after preserving the Claim's worktree identity and changing the supplied
`tree_oid` to a different valid 40-hex value. That probe returned
`error: run envelope binding is invalid` (terminated `exit_code: 2`). It
does not establish the plain-copy reason. See
`/tmp/fid-p24-binding-dt_3ikw9/same-worktree-different-tree.json`.

| Category | Authority boundary | Requested action | Required recheck |
|---|---|---|---|
| Specification/public interface | QA cannot edit an accepted message in VAL-029 or choose to reorder/rename domain denials. This is a contradiction in newly accepted text, not a product failure needing speculative implementation | TL reconciles the five named locations, preferably correcting the plain-copy cause to `run claim binding is invalid` in keeping with D6=A's actual-cause rule. If preserving the Envelope text is intentional, record explicit authority for changing check ordering/denial behavior and its affected cases | Verify the TL disposition against the live refusal and task 7/10 scopes; require task 7's original/copy/import cases and task 10's Human/JSON message pairing. Do not weaken Claim binding or substitute a same-worktree snapshot mismatch for a plain-copy row |

The probe's first scaffold attempt used invalid Claim mode `mutating`; its
original-worktree control failed with `invalid claim request` (process exit
2; probe terminated `exit_code: 1`). It was not evidence for the Claim cause.
Correcting it to accepted `mutation` and using fresh scratch produced the
successful controls and binding refusals above. A first sentinel scaffold
included an extra `credential_name` member and reached the closed-shape
check; its output is not the claimed domain guard. The final guard uses
exactly `{name,source,value}` and reaches the duplicate-component fallback.
Earlier scratch receipts remain; no failed attempt is erased.

#### 12.6.4 RC-02: Grammar Message Bound

**Must-fix; specification/public-interface blocker candidate.** Task 9a's
class table and VAL-064 require the detecting error's text; its unknown-command
rule and VAL-065 require exactly `unknown command "<word>" for "<command path>"`.
Phase 21 contract section 1 still limits JSON `message` to 512 Unicode
scalars. Neither command words nor unknown flag spellings have a stated bound
that keeps those messages within 512. Escaping can also expand the length.
DEC-281's first-511-plus-ellipsis rule expressly governs generic fallbacks,
not these grammar/unknown-command branches.

Reproduction (cwd `/tmp`, Python standard library, reviewed `/tmp/fid`):

```python
import subprocess
word = "x" * 600
p = subprocess.run(["/tmp/fid", "stop", "raise", "--" + word],
                   capture_output=True, text=True)
message = p.stderr.removeprefix("error: ").removesuffix("\n")
print(p.returncode, len(message))  # observed: 2 616
print(len(f'unknown command "{word}" for "fid"'))  # observed: 628
```

This proves a planning contradiction; the baseline JSON still uses its fixed
message, so it is not a claim that task 9a is already implemented incorrectly.
`internal/cli/codes.go:412-441` also provides no message-length enforcement to
resolve it automatically.

- **Authority boundary:** QA cannot choose how much parser text to omit while
  the assertion says exact text, or silently expand the accepted 512 bound.
- **Requested action:** TL specifies precedence and a bounded representation
  for grammar/unknown-command messages in contract section 9, task 9a and
  VAL-064/065. Reusing the accepted fallback escape-then-truncate rule would
  be consistent, but requires that scope to be stated. Details must retain
  the offending spelling/value.
- **Required evidence and recheck:** both modes, short exact text plus 512/513
  output boundaries, a long flag, a long unknown word and an escape crossing
  the cut; no raw control character, unchanged code/details/exit, and Human
  `error:` byte-equal to JSON. Task 9a owns construction; task 10's pairing
  corpus must use the same rule. Recheck the amended planning clauses before
  dispatch.

#### 12.6.5 Refinements, Cadence And Verdict

- Coverage rollup: all in-scope sources have assertions; every assertion has
  a validator and task/window. All 67 are labelled `Ready`, but VAL-029's
  exact reason is contradicted by RC-01 and VAL-064/065 need RC-02's bounded
  message rule. Those labels cannot authorize dispatch.
- VAL-067 is one fallback-message transformation checked over sites and
  boundary cases. QA-49/EVID-45 now cover it. Include whitespace-only causes
  (DEC-281/task 10), non-ASCII scalar counting and an escape crossing the cut;
  separate plain ASCII boundary and ESC-only rows do not prove the combined
  boundary. Check the sentinel guard actually reaches the generic branch.
  These evidence refinements implement existing rules; no assertion is edited.
- R20's split is **correct**: three excluded grammar sites move under task
  9a; four excluded typed sites are Finding-target twice, link endpoint type,
  and project-init argument. `review.list` is included even though its
  `argument` is `filter`, and `issueReadFailure` keeps its empty-string
  `value`. D6 does not normalize those details to `input`/null.
- Nonblocking documentation sync: the TL should update the phase brief's
  R1-R19 label to R1-R20 and its five-window summary to the six windows in
  section 8. The QA plan owns the cadence under DEC-280 item 10. Window D1
  covers 9/9a; D2 covers 10/11; Human VAL-054/055 stays after task 11.
- **Verdict: NOT READY (scoped), two Must-fix candidates RC-01/RC-02.** The earlier
  nine RR findings are dispositioned; the newly pinned `run.start` cause
  requires TL reconciliation, as does the grammar-message bound. The findings
  concern tasks 7/9a/10 and VAL-029/064/065 and do not invalidate tasks 1-3.
  The phase brief nevertheless holds all
  dispatch until recheck passes; QA does not lift that hold or allocate an
  ESC. After the TL corrections, recheck RC-01/RC-02 and any changed text. No final
  review or Human acceptance outcome is produced at this planning stage.
