# Phase 24 QA Plan - Operating Experience And Input Convergence

> Phase: 24 - Operating experience and input convergence
> Refreshed: 2026-10-04 (first pass)
> QA Lead: Claude Opus 5.5 (Claude Code), QA-plan pass after DEC-277
> Operating mode: **Heavy** (DEC-277). Every result in section 6 was produced
> fresh in this session on `phase-24` at
> `304aa45f12cf3966567231e6986e3070df409c68` (clean tree). Nothing carries
> forward from Phase 23.
>
> **Status: NOT READY. Seventeen Must-fix revisions (R-01..R-17, section 2.3)
> need Tech Lead disposition.** This is not a placeholder `BLOCKED`: every
> required command is concrete and was executed or proven feasible (section
> 6). Three revisions carry blocker candidates (section 2.4). Under the phase
> brief, nothing is dispatched until a readiness refresh of this plan
> confirms the dispositions (section 11).

## 1. Scope And Source Documents

Phase 24 gives a Human control of a stuck or diverging journey through fid
output alone, and converges inputs. fid observes repository identity itself,
so callers send only the expectations a check uses. Request digests cover the
input as written, and per-field digests name the differing field (migration
0019). One input descriptor per `--input` leaf drives help, an `--example`
template and field-path validation. Human failure lines repeat the JSON
failure with reason, blockers and next action. Stop, Attention and repository
lines become readable. Three judgment `get` leaves are added. Judgment views
show recorded reasoning, with one disclaimer per section. The phase has 14
sequential tasks and 57 assertions (52 deterministic, 5 Human).

Sources read for this plan:

- `AGENTS.md`, `CLAUDE.md`
- `project-docs/docs/current/phase_brief.md`, `decisions_active.md`,
  `progress_current.md`, `tl_escalations_current.md` (no open ESC)
- `project-docs/specs/08_QUALITY_PLAN.md` sections 5-8
- `project-docs/specs/09_AI_WORKFLOW.md` sections 2.6, 4.3, 4.4, 6.1
- `project-docs/specs/01_PRD.md` (the 19 PRD IDs the contract cites)
- `project-docs/specs/GLOSSARY.md` (Project-owned prose entry)
- `project-docs/docs/phases/phase-24/qa/validation-contract.md`
- `project-docs/docs/phases/phase-24/public-interface.md`,
  `input-inventory.md`, `decisions.md` (DEC-273..277), `triage-log.md`
- `project-docs/docs/phases/phase-24/tasks/task-1.md` .. `task-14.md`
- `project-docs/docs/phases/phase-23/decisions.md` DEC-227, DEC-259;
  `phase-23/qa/qa-plan.md` sections 1-11 and `qa/human-fixtures.md` (format
  and fixture precedent only)
- `project-docs/docs/phases/phase-21/public-interface.md` lines 80-90,
  402-412, 459 (grammar failures, unknown commands, `operation-unavailable`)
- `project-docs/templates/qa-plan-template.md`
- `project-docs/skills/validation-contract/SKILL.md`,
  `project-docs/skills/qa-evidence/SKILL.md`
- Source files inspected read-only are cited inline where they are evidence.

No held-out scenario, evaluator or excluded-scenario material was read. Mode B
scenario execution is **N/A**, as in Phases 21-23. `AGENTS.md` records
end-to-end and scenario validation as `N/A` until a harness is recorded, and
none is. The `(human)` checks in section 9 remain required. Their "take
control" scenario is the phase's Human run.

## 2. Validation Contract Refinement Check

- All in-scope requirements referenced by at least one `VAL-NNN`: **YES**.
  All 19 PRD IDs in the contract's Coverage table exist in `01_PRD.md`.
  PRD-OPS-011..014 and PRD-PROJ-006/007 were added under DEC-277 (rows at
  `01_PRD.md:150-151, 174-177`). Contract sections 2-11 re-derive to
  assertions, with the gaps listed in 2.1.
- Every `VAL-NNN` has one validator role: **YES**. QA Reviewer 49, QA Lead 3
  (VAL-005, VAL-048, VAL-051), Human 5. No row names two.
- Every `VAL-NNN` has a concrete evidence method: **YES as written; NO as
  executable** for VAL-009 ("no new observation" has no observable, R-03),
  VAL-036 (no producer registry exists, R-06) and VAL-049 (unsatisfiable
  provenance rule, R-13).
- Every `VAL-NNN` maps to a task/window or phase-end check: **YES**, verified
  mechanically. All 57 IDs are cited by at least one task spec, and no task
  cites an ID absent from the contract. The only out-of-range citations are
  the Phase 23 restatement notes (VAL-060, VAL-073).
- Every `VAL-NNN` describes one observable: **NO for six rows.** VAL-001,
  010, 013, 025, 031 and 047 each bundle two independent observables proven
  by different evidence (R-17). Matrix rows (006, 009, 016-019, 022-024, 036,
  040, 051) are acceptable because sections 4 and 7 enumerate their cases.
  The Human rows VAL-053..057 span several observables and are fanned out in
  section 9.
- Conditional or binary assertions require both states: **NO for four rows.**
  VAL-038 omits the in-table code on a non-qualifying leaf (R-08). VAL-040
  omits the Condition types that the code set does not reach (R-10). VAL-041
  omits same-worktree versus different-worktree (R-11). VAL-027/028/029 name
  no refusal code, so a refusal for the wrong reason passes (R-12).
- Parity or "matches reference" assertions have deterministic proof:
  **YES.** VAL-002, 007, 035, 042 and 050 need byte or semantic comparison,
  and section 4 pins the comparison inputs. VAL-035's leaf set should be
  derived, not listed (I-11).
- Revisions requested from Tech Lead: R-01..R-17 Must-fix and I-01..I-19
  Improvement (section 2.3).

### 2.1 Contract-Coverage Re-Derivation

Contract sections 2-12 were mapped to assertions independently of the
contract's Coverage table. The following clauses have no assertion that
evidences them, or one that cannot:

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

### 2.2 Clauses Verified Sound

- **Accepted 0019 DDL applied, not read** (section 6.4). On a real 0001-0018
  chain, `ALTER TABLE ... ADD COLUMN` with the exact accepted text succeeds
  on an empty table. `DROP COLUMN` restores the schema-18 `sqlite_master`
  dump byte for byte (VAL-001 is feasible as written). The CHECK and NOT NULL
  behave as accepted. A populated table refuses the column (I-14).
- **Inventory registries** R1, R2, R4-R12 and R14 reproduce their recorded
  counts at the planning base (section 6.3). R3, R13 and R15 differ
  (R-01, R-15, I-01).
- **The `--input` leaf set** is exactly the 39 leaves of inventory section 1.
- **First failing tests** for tasks 3, 6, 12 and 13 are red today (section
  6.5). The other first failing tests name new files or behavior that does
  not exist yet.
- **The 100k rebuild** has about 19.6% headroom (worst of three fresh samples
  24128 ms), not the 10% the phase brief records (I-18).
- **Documentation targets** named by tasks 1-14 all exist (section 6.8).

### 2.3 Revisions Requested From Tech Lead

Must-fix (M) items block Ready. Improvement (I) items are applied at the Tech
Lead's discretion. The Tech Lead owns every edit to `validation-contract.md`,
`public-interface.md`, `input-inventory.md` and the task specs. This plan
proposes edits and never applies them.

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

### 2.4 Blocker Candidates

QA emits candidates only. No `ESC` identity is allocated. TL-Lite or the Tech
Lead decides whether each is resolvable within existing authority.

| Candidate | Protected category | Evidence | Authority boundary | Requested action |
|---|---|---|---|---|
| R-04 | public interface (failure-code producer) | `evidence evaluate --help` (no key); contract section 6 producers | Producers of an accepted code are contract text (DEC-277) | Tech Lead states, from accepted text, `evidence.evaluate`'s code for an unobservable worktree; otherwise a contract amendment by DEC |
| R-05 | public interface (closed `reason` domain) | `internal/project/repository.go:61-352` causes; `git init` + `fid init` succeeds | A new member of a closed enum needs contract authority (DEC-195) | Tech Lead maps every cause to the five members, or amends section 6 by DEC |
| R-07 | public interface (Human failure rule; unknown-command JSON) | Silent exit 2 measured; Phase 21 lines 406, 459; `codes.go` lacks `operation-unavailable` | Section 9's "every failure" and the JSON `operation` for an unknown command are contract text | Tech Lead assigns conformance under DEC-175/194 with the `operation` value derived from accepted text, or records a section 9 exception by DEC |

## 3. Task Inventory

| Task | Focus | Covered VAL IDs | Risk | QA Notes |
|---|---|---|---|---|
| 1 | Migration 0019, field-digest storage | 001-005, 052 | high | R-01 file scope; I-07, I-14. DDL measured feasible (6.4). 100k benchmark early sample. |
| 2 | Input-as-written digests, `differing_fields` | 002, 005-009, 052 | high | 43 sites; R-02 per-leaf corpus; I-16. Pins "digest before observe". |
| 3 | Descriptor framework, help, `--example`, field paths | 010-016, 052 | high | New descriptor files inside `internal/cli` (no new package). Fence (DEC-192). R-07 if assigned here; I-01, I-09. |
| 4 | Observation, SourceInput, judgment leaves | 008, 009, 011-023, 035, 052 | high | R-03, R-04, R-05. Human Outcome references become checked. Path-safety fixtures. |
| 5 | `ExpectedRepository`, `review.record` | 009, 011-019, 022-025, 035, 052 | medium | `repository-mismatch` shape change; Phase 23 Review corpora converge. |
| 6 | Registration, Envelope, Claim; same-worktree predicate | 009, 011-019, 022, 023, 026-029, 035, 052 | high | Predicate reused by tasks 7-8 (window after 6, section 8). R-12. |
| 7 | Run leaves (DEC-273) | 009, 011-019, 022-025, 029-031, 035, 052 | high | `base_oid` becomes HEAD; recovery corpora converge. |
| 8 | Gate run, Evidence, real freshness | 009, 011-017, 022-024, 029, 032-035, 052 | high | Highest convergence risk (brief). Artifacts must stay git-ignored. Imported destinations. |
| 9 | Activation descriptors; whole-tree conformance | 011-016, 035, 052 | medium | Fence deleted. Sensitive fields. R-16. |
| 10 | Human failure lines, every family | 036-039, 052 (054) | high | 41 writers (R9), 32 retained assertions (R10). R-06..R-09. |
| 11 | Stop, Attention, repository lines | 040, 041, 052 (054, 055) | medium | R-10, R-11. Human walk-through at its window. |
| 12 | Judgment `get` leaves | 042, 043, 052 | low | Live leaf inventory (R14). |
| 13 | Disclaimers, recorded reasoning, sentinel guard, glossary | 044-048, 052 (056) | medium | R-15; I-15. The glossary already holds a partial "recorded reasoning" note (`GLOSSARY.md:28-31`). |
| 14 | Take-control journey, parity, absence | 049-052 (053-057) | high | R-13, R-14. Absence commands declared exit 1. |

## 4. Required Evidence Strategy

Rows state what QA requires; they do not prescribe implementation. "Named
tests" means tests the completion report names and whose output it records.
"Independent observation" is the test's own `git rev-parse HEAD`,
`git rev-parse HEAD^{tree}`, `git status --porcelain` and SHA-256 of the
absolute root, never fid output.

| Evidence ID | Required evidence | Covers | Owner | Produced by |
|---|---|---|---|---|
| EVID-01 | `sqlite_master` (type, name, tbl_name, sql) dumps at 18, after up, after down; the up diff is exactly the one column; the down dump is byte-equal to 18 | VAL-001 | Developer | task 1 |
| EVID-02 | Live versus rebuilt `request_field_digests_json`, byte-equal per row, over a journal with keyed mutations (task 1: issue, claim, run, review; task 2: every keyed leaf after R-02) | VAL-002 | Developer | tasks 1, 2 |
| EVID-03 | Journal fixture pair (and malformed map, I-05): replays with the map; without it stops with `journal-replay-stopped`, exit 6, naming the line | VAL-003 | Developer | task 1 |
| EVID-04 | Export, validate, import of a 19 archive; 18-manifest archive rejected (I-06) with the target tree digest unchanged | VAL-004 | Developer | task 1 |
| EVID-05 | `BenchmarkRebuildWorkspace100kEvents` `< 30000 ms/rebuild` with count parity; final review three fresh processes plus instance identity | VAL-005 | Developer, QA Lead | tasks 1, 2; final review |
| EVID-06 | Per-keyed-leaf conflict corpus (R-02): one-field row per leaf; two-field and `--actor` rows per family; exact sorted `differing_fields`, exit 4, journal line count unchanged; canonical retry replays byte-identically except `replayed:true` with the original EventRef | VAL-006, VAL-007 | Developer | task 2; converged by tasks 4-8 |
| EVID-07 | Three presence states: nonempty versus omitted `inapplicable_findings` conflicts naming it; `[]` versus omitted replays; supplied versus omitted Source `digest` conflicts naming the top-level field | VAL-008 | Developer | tasks 2, 4 |
| EVID-08 | Per observing keyed leaf: commit, change the observed value, retry, which replays with no Event; plus the unobservable-worktree and deleted-file retries (R-03) | VAL-009 | Developer | tasks 2, 4-8 |
| EVID-09 | In-process recursive Cobra walk: every leaf `Short` nonempty; `attention resolve --help` describes Attention resolution | VAL-010 | Developer | task 3; kept green by task 12 |
| EVID-10 | Reflection conformance, descriptor versus Go type, for every registered leaf; seeded mismatch fails; registry ∪ fence = live `--input` set until task 9, registry = live set after | VAL-011 | Developer | tasks 3-9 |
| EVID-11 | Help conformance per descriptor field (path, type, requirement, bounds or allowed values, meaning, template line); full text pinned for `stop raise`, `activation component` and one more | VAL-012 | Developer | tasks 3-9 |
| EVID-12 | Built-binary `--example` per leaf from a non-project temp dir: tree snapshot before and after, exit 0, empty stderr; template decoded with unknown fields disallowed and type-checked against the descriptor | VAL-013, VAL-014 | Developer | tasks 3-9 |
| EVID-13 | Exclusivity corpus (`--json`, `--input`, `--actor`, `--idempotency-key`, positional, `--project` per I-09) returns `invalid-argument` exit 2 with no read or write | VAL-015 | Developer | tasks 3-9 |
| EVID-14 | Structural corpus per leaf: 7 failure classes at top level and two levels deep; exact `input.<path>` and `allowed`; `sensitive` value null; no Event | VAL-016 | Developer | tasks 3-9 |
| EVID-15 | Removed-member corpus over the 17 changed leaves, each rejected naming its path | VAL-017 | Developer | tasks 4-8 |
| EVID-16 | Every Source position (I-08) × matching, mismatching, omitted-existing, omitted-missing; five path-safety fixtures; external reference omitted versus supplied | VAL-018-021 | Developer | tasks 4-7 |
| EVID-17 | Stored record repository, each Source repository and the Event `repository` compared with the independent observation (`base_oid` = HEAD, `worktree_id` = SHA-256 of root, `project_id` from config) | VAL-022 | Developer | tasks 4-8 |
| EVID-18 | Unmerged-paths and non-Git fixtures versus a normal worktree per observing leaf; unit rows for every reason and the R-05 mapping | VAL-023 | Developer | tasks 4-8 |
| EVID-19 | Five expectation leaves × match, HEAD differs, dirty differs; mismatch details exactly `{expected:{head_oid,dirty},observed}`; close/abandon unsafe gives `recovery-state-conflict` | VAL-024, VAL-025 | Developer | tasks 5, 7, 8 |
| EVID-20 | Registration fields and `binding_digest` recomputed from the independent observation; same-worktree versus plain-copy refusals with the exact code (R-12); imported destination with binding versus plain copy for four leaves | VAL-026-029 | Developer | tasks 6-8 |
| EVID-21 | `run.start` after commit and edit since Envelope creation records the observation; close/abandon `terminal_repository` observed, with Obligations evaluated against it | VAL-030, VAL-031 | Developer | task 7 |
| EVID-22 | Gate run after an edit since Run start: Attempt, Receipt and both Events carry one identical repository equal to the pre-spawn observation; mismatch leaves no Attempt, no process and no Event | VAL-032, VAL-024 | Developer | task 8 |
| EVID-23 | Freshness: unchanged, edited, reverted × `evidence.evaluate`, terminal-transition denial, `envelope.complete` denial | VAL-033 | Developer | task 8 |
| EVID-24 | `evidence.link` target equals the subject's recorded repository for receipt, review and Human Outcome subjects | VAL-034 | Developer | task 8 |
| EVID-25 | Shared fake/fid corpus over the derived leaf set (I-11), identical acceptance and rejection | VAL-035 | Developer | tasks 4-9 |
| EVID-26 | Built-binary Human/JSON pairing from `RegisteredCodes()` with pinned exclusions (R-06), grammar and unknown-command rows (R-07); `blocked by:` both states; next-table rows plus five outside codes plus the R-08 row; reason-rendering units with the R-09 boundary | VAL-036-039 | Developer | task 10 |
| EVID-27 | `stop list`, `attention list`, `resume` over 10 codes and 10 Condition types, open and resolved (R-10), no raw `{`; repository lines over clean/dirty and same/different worktree (R-11) with a forbidden-token scan | VAL-040, VAL-041 | Developer | task 11 |
| EVID-28 | `get` JSON equals the list item with `event:null`; unknown ID gives `*-not-found` exit 3 `{id}`; missing ID, extra argument, unknown flag give `invalid-argument` exit 2 | VAL-042, VAL-043 | Developer | task 12 |
| EVID-29 | Disclaimer count for 0/1/3 items per view; 80/81/multi-line text per closed field in list and detail views; ESC/CR/BEL/NUL escaped; sentinel guard over every Human fixture; `PRIVATE_` absent | VAL-044-047 | Developer | task 13 |
| EVID-30 | `GLOSSARY.md` entries compared item by item with the contract section 10 closed field list | VAL-048 | QA Lead | task 13; final review |
| EVID-31 | Tagged take-control journey with a provenance record per command and value, classed per R-13; live/rebuilt/imported parity over its workspace, manifest 19 | VAL-049, VAL-050 | Developer | task 14 |
| EVID-32 | Descriptor forbidden-member scan; the two declared absence `rg` lines terminating `exit_code: 1` | VAL-051 | Developer, QA Lead | task 14; final review |
| EVID-33 | Per-task wall times for `cmd/fid/integration`, `cmd/fid`, `internal/workflow` and the slowest ordinary package (I-13) | VAL-052 | Developer | every task |
| EVID-34 | Human fixture generator output (R-14) and Human judgment notes | VAL-053-057 | Developer fixture; Human judgment | task 14; window after 11; phase end |
| EVID-35 | Full phase gate, DEC-259 full-`checkptr` control, archive-scale gate, three 100k samples, reference instance | no VAL | QA Lead | final review |
| EVID-36 | Security review: observation under lock, path containment and symlinks, digest-before-observe, sensitive-field redaction, field-digest privacy (I-17), import binding | no VAL | QA Lead | final review |
| EVID-37 | Inventory registries R1-R15 and section 1 walk regenerated, with deltas explained | no VAL | QA Reviewer, QA Lead | every window; final review |

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
| Registry and absence searches | the `rg` lines in tasks 1, 2, 9, 10, 11, 13, 14 | recorded real status; tasks 13/14 declared `exit_code: 1` | owning task; window; final review |
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
inside Go's default limit. The task-spec form governs. No gate is narrowed, substituted or deferred. Only
tasks 13 and 14 declare expected non-zero commands (R-16 may add task 9).
Cleanup commands are never gate results. Report lint
(`project-docs/tools/report-lint.sh <task-N-done.md> <task-N.md>`) is a Phase
Runner control after the report commit, never a Developer or QA gate. Run
heavy gates serially.

## 6. Command And Feasibility Validation

Heavy mode: nothing carries forward. The working directory for every command
is `/home/mark/projects/fid` unless stated. A serial wrapper captured each
command's own exit status and wall time (`bash -c "<cmd>"; rc=$?`). Logs are
in session scratch storage. `go clean -testcache` ran first as setup, not as
a gate result, so no result is a cached replay. Light read-only probes ran
during gate steps 3-9. Gate wall times therefore include minor contention,
and no step failed.

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
| 100k rebuild (tasks 1, 2; final ×3) | executed three times, fresh processes | PASS | `22653` / `23211` / `24128 ms/rebuild`; `ok ... internal/journal 25.506s` / `25.995s` / `26.885s` (each terminated `exit_code: 0`). Worst-case headroom 19.6% (I-18) |
| Archive-scale gate | executed | PASS | `ok github.com/mbell-nzl/fid/internal/export 143.927s` (terminated `exit_code: 0`, 164.6 s wall) |
| Full-`checkptr` control `go test -race -timeout 30m ./...` | executed | PASS | last line `ok github.com/mbell-nzl/fid/internal/workflow 171.784s`; 0 `FAIL` lines (terminated `exit_code: 0`, 553.2 s wall). `cmd/fid` 539.666 s, `internal/issues` 379.787 s; not held to VAL-052 budgets (DEC-259), bounded by `-timeout 30m` (I-13) |
| R1 `rg -n "RequestDigest:" internal cmd -g '!*_test.go'` (task 2) | executed | PASS | `exit_code: 0`; 43 lines in 32 files, matching the inventory |
| R2 helper `rg` (task 2) | executed | PASS | `exit_code: 0`; 11 helpers in `internal/{links,review(2),evidence(2),workflow(4),labels,issues}`; I-16 |
| R3 `rg -n "idempotency_results" migrations internal -g '!*_test.go'` (task 1) | executed | PASS (syntax); **inventory incomplete** | `exit_code: 0`; 16 files including four closed registries outside task 1's list (R-01) |
| R4 | `ls migrations`; source | PASS | latest `0018_p23_remove_v1_judgment`; `ActiveMigrationVersion = 18`, `ProjectionSchemaVersion = 5` (`internal/db/projection_meta.go:10-11`); manifest uses `ActiveMigrationVersion` (`internal/export/export.go:137`) |
| R5 `rg -n "func decodeClosed[A-Za-z]*Input" cmd/fid` (task 9) | executed | PASS (syntax); status undeclared (R-16) | `exit_code: 0`; 7 helpers (`review`, `run`, `claim`, `activation`, `envelope`, `judgment`, `gate`), 14 calling files |
| R6 | executed | PASS | 4 `same*Repository`, 6 validators, `importedGateDestination` (`internal/evidence/execution.go:338`); `importedDestinationBinding` (`internal/workflow/import_binding.go:18`) |
| R7 | executed | PASS | caller-repository writers in `claims/generic.go`, `evidence/{evidence,execution,registration}.go`, `review/review.go`, `workflow/{envelope,run,run_close}.go` |
| R8 `internal/cli/codes.go` | inspection | PASS | 78 registered codes; `repository-unavailable` and `operation-unavailable` absent; `idempotency-conflict` details 4 strings; `source-digest-mismatch.expected_digest` non-nullable; `repository-mismatch.expected` a full Repository |
| R9 `rg -n "func (write[A-Za-z]*Failure\|render[A-Za-z]*Error)\(" cmd/fid -g '!*_test.go'` (task 10) | executed | PASS | `exit_code: 0`; 41 functions (I-10) |
| R10 | executed | PASS | 32 `"error: ` assertions in 13 test files |
| R11 `rg -n -i "worktree\|dirty digest\|HeadOID\|head_oid" cmd/fid -g '!*_test.go'` (task 11) | executed | PASS | `exit_code: 0`; Human lines in `p23_human.go:139-140`, `resume_human.go:113-117`, `run_inspect_human.go:17-19`; JSON mapping in `review.go:264`; a flag description in `run_recovery.go:68` |
| R12 | `ls` | PASS | `stop_human.go`, `resume_human.go`, `resume_human_guidance.go` present |
| R13 `rg -n "p23CheckJudgmentQualifiers\|PRIVATE_" cmd` | executed | PASS (syntax); **scope gap** | also matches `cmd/fid/integration/testdata/p23_conformance/operations.json` four times (R-15) |
| R14 | executed | PASS | `cmd/fid/contract_surface_test.go:128` |
| R15 Short walk | executed recursively over the live tree | PASS with I-01 | 92 leaves; 20 lack a Short (8 read, 12 `--input`) |
| Section 1 `--input` walk | executed | PASS with I-02 | exactly the 39 inventory leaves; loop status `exit 1` from the last `grep -q` |
| Section 3 fixture families | executed | PASS | 56 test files name caller repository fields (29 integration); 8 pinned literal digests in 7 files; 4 `p2*_fake_test.go` files |
| Task 13/14 per-item disclaimer `rg` | executed | PASS (syntax); matches expected before task 13 | `exit_code: 0`; `cmd/fid/p23_human.go:163, 167, 176` (I-15) |
| Task 13/14 `rg -n "PRIVATE_" cmd/fid/integration` | executed | PASS (syntax); **cannot reach declared exit 1 in scope** | `exit_code: 0`; `p23_human_fixtures_test.go:143`, `p23_human_leaf_pairs_test.go:45-126`, `testdata/p23_conformance/operations.json:67-145` (R-15) |
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
| Keyed leaf set (VAL-006, R-02) | `--help` scan of all 92 leaves | PASS | 50 leaves take `--idempotency-key` across 21 families; `export import` is the only one outside R1; `evidence evaluate` is the only `--input` leaf without a key (R-04) |
| Task 3 first failing test | `fid --json attention resolve --input {"bogus":1}` | PASS (red today) | `invalid-argument`, `argument:"input"`, message `invalid Run request`; exit 2 |
| `--example` absent today | `fid stop raise --example` | PASS (red today) | `error: unknown flag: --example`, exit 2; JSON `operation:"raise"` (not `stop.raise`) |
| `attention resolve --help` text | executed | PASS (red today) | `--input string Run JSON input`; first line `Usage:` |
| Task 12 leaves absent | `fid concern get CON-1` | PASS (red today) | exit 2, **no output in either mode** (R-07) |
| Unknown command and subcommand | `fid bogus`, `fid concern bogus` in both modes | **FAIL vs Phase 21 contract (R-07)** | exit 2, empty stdout and stderr |
| Grammar failure parity | `fid stop list --bogus`, `fid review get` in both modes | **conflicts with section 9 intent (R-07)** | Human `error: unknown flag: --bogus` / `error: accepts 1 arg(s), received 0`; JSON message `invalid command request`, `argument:"request"`, `value:null` |
| Task 6 first failing test | source (`internal/claims/generic.go:483`) | PASS (red today) | empty `worktree_id` fails validation |
| Task 13 first failing test | source (`cmd/fid/p23_human.go:154-163`) | PASS (red today) | Concern item lines print subject, severity, status, disposition, sources and the per-item disclaimer; no `summary` or `reasoning` |
| Task 11 fixture reach | source (`internal/workflow/stop.go:175-189, 440-500`) | PASS | 10 codes, 10 Condition types; `stop raise` accepts any pairing with bindings (R-10) |
| Observation causes | source (`internal/project/repository.go`) | **gap (R-05)** | about 30 `unsafeRepository` causes; operation markers are recorded, not rejected |
| Existing codes used by assertions | `internal/cli/codes.go` | PASS | `journal-replay-stopped`, `recovery-state-conflict`, `evidence-stale`, `concern/deviation/outcome-not-found` (not-found, `{id}`), `archive-invalid`, `idempotency-conflict`, `source-digest-mismatch`, `repository-mismatch` present |
| Plain-copy refusal codes | source | PASS; **unnamed in contract (R-12)** | `cmd/fid/envelope.go:253`, `cmd/fid/claim.go:234` generic `invalid-argument` |
| JSON failure stream | probes | noted | JSON failure envelopes are written to stderr with empty stdout |

### 6.6 Scope And Gate Reach

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Gate reach covers owned paths | task file lists against registry output | **PASS with R-01, R-15, R-07** | both race gates compile every `cmd/`/`internal/` package; three tasks would need files they do not list |
| Same-worktree predicate home | source | PASS | `internal/workflow/import_binding.go` exists (task 6 names it) |
| Fake adapters present | `ls` | PASS | `p22_fake_test.go`, `p22_semantic_fake_test.go`, `p22_portable_fake_test.go`, `p22_activation_fake_test.go` (I-11) |
| Glossary starting state | `GLOSSARY.md:21-31` | PASS | "recorded reasoning" exists only as a note inside Project-owned prose, with a partial list; VAL-048's non-qualifying state exists today |

### 6.7 Human Checklist Feasibility

| Referenced step | Validation method | Result | Evidence |
|---|---|---|---|
| Disposable project in a non-Git directory | executed | PASS | `fid init --prefix PNG` gives `initialized fid project in .fid/` (also the VAL-023 non-Git fixture) |
| Disposable project with unborn HEAD | executed | PASS (and R-05) | `git init` then `fid init --prefix PUB` succeeds |
| Phase 23 fixture generator | source | PASS | `TestP23HumanFixtures` with `FID_P23_HUMAN_FIXTURE_ROOT` (`cmd/fid/integration/p23_human_fixtures_test.go:57-63`), including a protected-Decision checkpoint (line 264) |
| Phase 24 stuck-journey fixture | task specs | **FAIL - drives R-14** | no task owns a generator; task 14 assigns preparation to the QA Lead |
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
| `/tmp` filesystem | `df -h /tmp` | CONSTRAINT (I-19) | tmpfs 3.9 GB, 1.4 GB free; no stale `fid-integration-binary-*` directories |
| Disk-backed alternative | `df -h /` | PASS | 896 GB available |
| Per-binary test time | measured | PASS | tagged 293.1 s (33% of 900 s); `cmd/fid` 330.8 s (55% of 600 s); `internal/workflow` 110.3 s; all well inside the `-timeout 30m` bound |

## 7. Deterministic Checklist

The QA Lead executes or verifies these at final review. "Rerun named tests"
means `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run
'^(<names>)$' <packages>` (adding `-tags integration` for tagged tests), using
the exact names in the completion report, with output recorded. Real-process
probes use the `/tmp/fid` built by the phase gate against a fresh `git init` +
`fid init --prefix <P>` project. PASS without evidence is invalid.

| ID | Check | Covers | Evidence Required |
|---|---|---|---|
| QA-01 | Full seven-command phase gate in the section 5 form, each with token and observed termination `(qa-lead)` | no VAL | command transcript |
| QA-02 | DEC-259 full-`checkptr` control and DEC-217 archive-scale gate `(qa-lead)` | no VAL | transcripts |
| QA-03 | Three consecutive fresh 100k samples, each `< 30000 ms/rebuild`, with the reference instance captured in the same session `(qa-lead)` | VAL-005 | transcripts, instance values |
| QA-04 | Fresh-init `sqlite_master` dump versus a schema-18 chain (one-column diff); apply the 0019 down migration and `cmp` with 18; rerun task 1 named tests `(qa-lead)` | VAL-001 | dumps, `cmp` status, test output |
| QA-05 | On one workspace, `fid rebuild` then compare `request_field_digests_json` live versus rebuilt, plus the task 1/2 parity tests `(qa-lead)` | VAL-002 | SQL dumps, test output |
| QA-06 | Rerun the replay fixture pair (and malformed map); confirm the stop line and exit 6 `(qa-lead)` | VAL-003 | test output |
| QA-07 | Export, validate, import a 19 archive; edit a copy's manifest to 18 and confirm the named rejection with the target unchanged `(qa-lead)` | VAL-004 | outputs, tree digests |
| QA-08 | Rerun the per-keyed-leaf conflict and canonical-retry corpus; reproduce one conflict per family through the binary `(qa-lead)` | VAL-006, VAL-007 | test output, probe JSON |
| QA-09 | Reproduce the three VAL-008 presence states through the binary `(qa-lead)` | VAL-008 | probe JSON |
| QA-10 | Per observing keyed leaf, rerun the advanced-HEAD retry and the unobservable-worktree retry (R-03); journal line count unchanged `(qa-lead)` | VAL-009 | test output, line counts |
| QA-11 | Recursive live-tree walk: every leaf has a Short; `attention resolve --help` text `(qa-lead)` | VAL-010 | walk output |
| QA-12 | Registry equals the live `--input` set (39), no fence, reflection conformance and seeded-mismatch tests `(qa-lead)` | VAL-011 | test output, walk diff |
| QA-13 | For every `--input` leaf from a non-project temp dir: `--example` (tree snapshot before and after), template decode, and `--help` lists every descriptor field and the template line `(qa-lead)` | VAL-012, VAL-013, VAL-014 | per-leaf outputs, snapshot digests |
| QA-14 | Exclusivity corpus rerun, plus binary probes of `--json --example` and `--project <dir> --example` `(qa-lead)` | VAL-015 | outputs |
| QA-15 | Structural corpus rerun, plus binary probes of one nested unknown key, one enum violation and one sensitive field `(qa-lead)` | VAL-016 | outputs |
| QA-16 | Removed-member probe for each of the 17 changed leaves `(qa-lead)` | VAL-017 | outputs |
| QA-17 | Source corpus rerun; binary probes of match, mismatch, omitted-existing, omitted-missing, absolute path, `..`, escaping symlink, directory, in-root symlink, external reference `(qa-lead)` | VAL-018-021 | outputs, stored rows |
| QA-18 | For one leaf per task 4-8, stored and Event repositories against an independent observation `(qa-lead)` | VAL-022 | SQL rows, git output |
| QA-19 | Unmerged-paths and non-Git probes per observing leaf versus normal; unit rows for every reason `(qa-lead)` | VAL-023 | outputs, test output |
| QA-20 | Expectation corpus (5 leaves × 3 states); mismatch details shape; close/abandon unsafe `(qa-lead)` | VAL-024, VAL-025 | outputs |
| QA-21 | Registration binding digest recomputed; same-worktree versus plain copy for `envelope.create`, `claim.acquire`; imported destination versus plain copy for four leaves, each refusal with its exact code (R-12) `(qa-lead)` | VAL-026-029 | outputs, recomputation |
| QA-22 | `run.start` after commit and edit; close/abandon terminal repository and Obligation evaluation `(qa-lead)` | VAL-030, VAL-031 | outputs, SQL rows |
| QA-23 | Gate run after an edit: one identical repository across Attempt, Receipt and both Events; mismatch leaves no Attempt row, no Event, and no marker file from the gate's command (proving no spawn) `(qa-lead)` | VAL-032 | SQL rows, journal lines |
| QA-24 | Freshness three states × three surfaces `(qa-lead)` | VAL-033 | outputs |
| QA-25 | `evidence.link` target per subject type `(qa-lead)` | VAL-034 | SQL rows |
| QA-26 | Rerun the fake/fid corpus; confirm the derived leaf set has no one-sided case `(qa-lead)` | VAL-035 | test output, case inventory |
| QA-27 | Rerun the pairing corpus; confirm rows plus exclusions equal `RegisteredCodes()`; binary probes of one Stop-blocked denial, one next-table code per row, one non-`--input` `invalid-argument`, one unknown command (per R-07 disposition), and the 200/201-scalar boundary (R-09) `(qa-lead)` | VAL-036-039 | test output, stderr captures |
| QA-28 | Stop and Attention views over every code and Condition type, open and resolved; repository lines clean/dirty and same/different worktree; forbidden-token `rg` over captured Human output `(qa-lead)` | VAL-040, VAL-041 | captured output, `rg` status |
| QA-29 | `get` leaves versus list items; not-found and grammar probes `(qa-lead)` | VAL-042, VAL-043 | outputs |
| QA-30 | Disclaimer counts, reasoning cut and full text, control escaping, sentinel guard rerun; `PRIVATE_` search status `(qa-lead)` | VAL-044-047 | captured output, test output |
| QA-31 | `GLOSSARY.md` entries compared item by item with contract section 10 `(qa-lead)` | VAL-048 | quoted entries and comparison |
| QA-32 | Rerun the take-control journey; check each provenance record's class against the R-13 rule; parity test `(qa-lead)` | VAL-049, VAL-050 | test output, provenance log |
| QA-33 | Descriptor forbidden-member scan; both absence `rg` lines terminate `exit_code: 1`; no raw-JSON Stop renderer `(qa-lead)` | VAL-051 | test output, `rg` status |
| QA-34 | Every completion report records the three package times (plus the slowest ordinary package); none breaches 600/900 s `(qa-lead)` | VAL-052 | report excerpts |
| QA-35 | Security review: observation under the write lock and after replay; path containment and symlink escape; digest-before-observe; sensitive fields never echoed; field-digest privacy (I-17); import-binding equivalence limited to the bound root `(qa-lead)` | no VAL | findings with file references |
| QA-36 | Every completion report lists every exact task-spec command with token and observed termination; declared non-zero commands recorded as declared; converged rows carry introducing commit, invalidating change and executing command (DEC-162/165/166) `(qa-lead)` | no VAL | report-to-spec comparison |
| QA-37 | `git diff 304aa45 -- go.mod go.sum` empty, or each dependency recorded per the budget; no new package outside `internal/*` boundaries in `02_TRD.md` `(qa-lead)` | no VAL | diff output |
| QA-38 | Registries R1-R15 and the section 1 walk regenerated; R2 lists only the shared helper; R5 matches the R-16 disposition; deltas explained `(qa-lead)` | no VAL | command output |
| QA-39 | Spec amendments named by tasks 1-14 applied (03, 05, 06, 07, 08, GLOSSARY); spec sync inputs for wrap-up `(qa-lead)` | no VAL | spec diffs |
| QA-40 | Every window review is `PASS`; each earlier `REPAIR` closed by a fresh independent re-review; the task-11 Human walk-through recorded `(qa-lead)` | no VAL | review files |

## 8. Mid-Phase Review Cadence

DEC-277 proposed windows after tasks 3, 8, 11 and 14. This plan adopts them
and adds a window after task 6, so that five dependent tasks do not share one
window:

| Window | After task | Reviews | In-window focus |
|---|---|---|---|
| A | 3 | tasks 1-3 | VAL-001-009 (migration, replay fail-closed, schema 19, benchmark, field digests; task 2 rows of 008/009), VAL-010-016 framework |
| B | 6 | tasks 4-6 | Observation, SourceInput, expectations and the same-worktree predicate before tasks 7-8 reuse them (VAL-008 Source case, 009, 011-029 for these leaves, 035). Fresh feasibility for the unmerged-path, plain-copy and import-binding fixtures that tasks 7-8 prescribe |
| C | 8 | tasks 7-8 | DEC-273 Runs and gate runs, real freshness (VAL-009, 011-019, 022-025, 029-035 for these leaves). Fresh feasibility for every non-default command and corpus that tasks 9-14 prescribe |
| D | 11 | tasks 9-11 | Whole-tree descriptors (VAL-011-016 complete), failure lines (036-039), Stop and repository lines (040, 041). **Human walk-through** of new failure and Stop output (VAL-054, VAL-055 interim), recorded in `qa/human-walkthrough-after-task-11.md` from regenerated fixtures (R-14) |
| E | 14 | tasks 12-14 | `get` leaves, judgment views, journey, parity, absence (VAL-042-052) |

- Each window uses `project-docs/templates/mid-phase-review-template.md` and
  the next unused `qa/review-M.md`. It reruns the in-window named tests and
  verifies in-window `VAL-NNN` in both states. It compares each report's
  commands with its spec. It reruns inventory registries R1-R15 (inventory
  header). It checks VAL-052 wall times for every in-window task, and records
  gate evidence or why a gate was not run.
- The window-D walk-through records Human judgment only. Its interim results
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

Fixtures come from the owner assigned under R-14. The Human records judgment
only. The HUMAN-01 step list is re-derived at final review from the task-14
journey's recorded steps. Any added or removed step keeps the parent
VAL-053.

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
| HUMAN-05a | Author a `run close` input with `expected_repository` from `--help`/`--example`; accepted on the first or second try? `(human)` | VAL-057 | attempts, final input |
| HUMAN-05b | Author a `decision record-reference` input with the digest omitted; accepted on the first or second try? `(human)` | VAL-057 | attempts, final input |
| HUMAN-05c | Author an input for one more leaf of your choice; accepted on the first or second try? `(human)` | VAL-057 | leaf, attempts, final input |

## 10. Phase-End Decision Rule

The phase can proceed to final review only when:

- Every required deterministic gate in section 5 passes with observed
  termination evidence, including the three fresh 100k samples, the
  archive-scale gate and the full-`checkptr` control.
- Every `VAL-NNN` (all 57, plus any row added under section 2.3) has planned
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

This plan becomes Ready only after:

1. The Tech Lead dispositions R-01..R-17. For R-04, R-05 and R-07, the Tech
   Lead records whether accepted text decides them or a DEC (with Human
   approval where protected) is needed. Improvement items are applied or
   declined.
2. The Tech Lead applies accepted edits to the contract, the inventory and the
   affected task specs (at least tasks 1, 2, 3 or 10, 4, 6-11, 13, 14), and
   marks rows `Ready`.
3. A fresh QA Lead readiness refresh revalidates every changed command,
   fixture owner, file scope and assertion in a new section of this file. It
   reuses no result above for anything that changed.
4. The Phase Runner's readiness check passes.
