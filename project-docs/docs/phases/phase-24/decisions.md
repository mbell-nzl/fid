# Phase 24 Decisions

> Purpose: phase-local decision history, spec sync verdicts, workflow
> notes, and wrap-up decisions. DEC numbering is global and never resets.
> Phase scope authority is DEC-272 (`../phase-23/decisions.md`). Entries cite
> the contract (`public-interface.md`) rather than restating it.

## Decisions

### DEC-273: Run and gate-run repository checks use the real worktree (D1 = A)

- **Date:** 2026-10-04. ACCEPTED by the Human ("yes, option A") after the
  Tech Lead explained the alternative.
- **Context:** today `run.start` checks the caller's repository block only
  against the Envelope's stored `repository_scope`, and `gate.run` only
  against the Run's `start_repository`. Neither is ever compared with the
  real worktree. Once fid observes the repository (DEC-272 Part 1), each
  comparison must mean something. Option B (real worktree equals the
  snapshot) would refuse every rework Run after an edit, which is impossible
  to escape in Phase 24 because Envelope revision is Phase 25 scope. It would
  also stop a Builder testing their own edits within a Run.
- **Decision:** `run.start` and `gate.run` compare the Envelope or Run binding
  as *same project and same worktree* only. The caller's
  `expected_repository` (`head_oid`, `dirty`) is checked against fid's
  observation, and a mismatch rejects with `repository-mismatch`. The
  observation is recorded as the Run's `start_repository`, and for `gate.run`
  as the Attempt and Receipt repository.
- **Consequences:** a Run or gate run may start after edits or commits made
  since the Envelope or Run began. The existing freshness rule (PRD-EVID-003)
  stales proof whose recorded tree differs from the current one, which
  becomes a real check under task 8. `run.recover` is unaffected because it
  already observes. Contract section 5.
- **Refs:** DEC-272; `public-interface.md` sections 4, 5 and 13;
  PRD-RUN-005; PRD-GATE-003; PRD-EVID-003.

### DEC-274: Add detail read leaves for Concerns, Deviations and Human Outcomes (D2 = A)

- **Date:** 2026-10-04. ACCEPTED by the Human ("A").
- **Context:** DEC-272 A7 shows recorded reasoning cut to one line in lists
  and in full in detail views. Concerns, Deviations and Human Outcomes have
  only `list` leaves, so their full reasoning was readable only through
  `--json` or the recording command's own output. Findings are already
  covered by `review get`.
- **Decision:** add read leaves `concern.get` / `fid concern get <id>`,
  `deviation.get` and `human-outcome.get`, returning `{concern:Concern}`,
  `{deviation:Deviation}` and `{outcome:HumanOutcome}`. Their Human views
  show recorded reasoning in full. Not-found reuses `concern-not-found`,
  `deviation-not-found` and `outcome-not-found`, which gain these producers.
  Option B (`--full` on lists) and option C (JSON only) were rejected.
- **Consequences:** three new public read leaves (task 12). The live leaf
  inventory test (R14), the Phase 23 section 5.1 producer table and the
  Human fixtures gain rows. Capability rows are unchanged.
- **Refs:** DEC-272 A7; `public-interface.md` sections 10 and 13.

### DEC-275: Store per-field request digests to name the differing field (D3 = A)

- **Date:** 2026-10-04. ACCEPTED by the Human ("A").
- **Context:** DEC-272 requires `idempotency-conflict` to name the differing
  field. fid stores one digest per request, so it cannot tell which part
  changed.
- **Decision:** every keyed mutation also stores one digest per top-level
  public input name (`--actor` included). These are fingerprints, never
  input values. Migration `0019` adds
  `idempotency_results.request_field_digests_json`, and the idempotency
  Effect of every keyed Event carries it, so rebuild reproduces it. Conflict
  details gain `differing_fields`. Option B (keep only whole-request
  digests) was rejected.
- **Consequences:** projection schema 19, and manifest `projection_schema` 19
  at phase end. Down migration restores schema 18 exactly (DEC-264).
  Pre-Phase-24 journals fail closed (DEC-270). Task 1 owns storage; task 2
  adopts the rule at all 43 digest call sites (inventory R1-R4).
- **Refs:** DEC-272; `public-interface.md` sections 6 and 7.

### DEC-276: `--help` input documentation and `--example` template rules (D4)

- **Date:** 2026-10-04. ACCEPTED by the Human ("yes, accept all").
- **Context:** DEC-272 option C requires every `--input` leaf to document its
  input and print an editable template. Three choices were open.
- **Decision:** templates use static, correctly typed placeholders (never
  pre-filled from the workspace or repository). They contain required fields
  only; optional fields are documented in `--help`. `--json --example` is
  `invalid-argument` because a template is not a `fid-cli.v1` envelope.
  `--example` needs no project, Actor, key or input and writes nothing. Help
  lists every input field with type, requirement, bounds or allowed values
  and a one-line meaning, then names the template command.
- **Consequences:** templates are deterministic and pinned exactly by tests.
  An expectation such as `expected_repository` is always typed by the
  caller, never rubber-stamped from fid's observation.
- **Refs:** DEC-272 Part 2; `public-interface.md` section 8.

### DEC-277: Phase 24 plan approved (contract, Heavy mode, 14 tasks)

- **Date:** 2026-10-04. ACCEPTED by the Human ("Approved").
- **Decision:** the following are approved, with D1-D4 as DEC-273..276:
  - the contract amendment `public-interface.md`, including the Tech Lead
    details the Human was shown: Source digests checked wherever supplied,
    `repository-unavailable`, run close/abandon expectation mismatch as
    `repository-mismatch`, and field-path `invalid-argument`;
  - Heavy mode;
  - the 14-task order in the phase brief.

  The Tech Lead added the *same worktree* rule (contract section 5) when
  drafting assertions. It preserves Phase 22 import semantics and generalizes
  the existing `importedGateDestination` and `importedDestinationBinding`
  rule. It adds no new behaviour.
- **Spec edits under this decision:** `01_PRD.md` gains PRD-OPS-011..014 and
  PRD-PROJ-006/007, which give the validation contract durable sources. The
  roadmap section 10 records the plan.
- **Validation contract:** 57 assertions (52 deterministic, 5 Human), all
  `Draft` until the QA Lead's QA-plan pass. VAL-044/047/056 restate Phase 23
  VAL-060/073.
- **QA cadence (proposed to the QA Lead):**
  - windows after tasks 3, 8, 11 and 14;
  - a Human walk-through of the new failure and Stop output at the
    task-11 window (Phase 23 workflow notes);
  - Mode B held-out runs `N/A` (no harness recorded), as in Phases 21-23.
- **Refs:** DEC-272..276; `qa/validation-contract.md`; `input-inventory.md`.

### DEC-278: QA-plan dispositions and contract clarifications

- **Date:** 2026-10-04. Tech Lead decision under DEC-277. It clarifies
  accepted contract text and adds no behaviour beyond it. It was presented to
  the Human together with D5, which is the one item accepted text does not
  decide.
- **Context:** the QA plan's first pass (`qa/qa-plan.md`, `eb34889`) was
  NOT READY. It raised 17 Must-fix revisions (R-01..R-17), 19 improvements
  (I-01..I-19) and three blocker candidates (R-04, R-05, R-07). The Tech
  Lead checked each against source at `eb34889` before deciding.
- **Decisions:**
  1. **Observing leaves (R-04).** The contract section 2 rule "the request
     returns `repository-unavailable`" applies to every leaf that observes,
     `evidence.evaluate` included. Section 6's producer wording ("every
     observing mutation") omitted the one observing read, so it now reads
     "every observing leaf". The list of 18 observing leaves is in section 2.
     `evidence.link` does not observe: it copies the subject's recorded
     repository. Leaves that observe under earlier contracts and are not
     listed keep their Phase 21-23 behaviour. CF-046 records their failure
     codes.
  2. **Reason mapping (R-05).** Every measured unsafe cause maps by capture
     stage to one of the five accepted members (contract section 6). No new
     member is needed. A repository with no commit is `not-a-repository`,
     because no Repository value can be formed without HEAD.
  3. **Grammar failures (R-07, first half).** Phase 21 `public-interface.md`
     section 2 accepted that grammar failures name the offending spelling or
     value. That was never delivered: today they say
     `argument:"request"`, and `issue list --bogus` prints nothing in JSON
     mode. Delivery is bounded repair of an accepted allocation (DEC-194).
     Task 9 owns it over the whole live tree, with the details in its spec.
     The JSON `message` becomes the parser's text, so the section 9 `error:`
     line keeps the spelling that Human mode prints today.
  4. **Unknown commands (R-07, second half).** Phase 21 section 4 accepted
     `operation-unavailable` for unknown commands, but `internal/cli/codes.go`
     never registered it, and both modes are silent today. Task 9 delivers it
     (DEC-194). The JSON envelope `operation` for a command that does not
     exist is not in accepted text. It is Human decision D5 (contract
     section 13), to be recorded as DEC-279. Task 9 is not dispatched before
     DEC-279.
  5. **Reason-line cut (R-09).** The `reason:` line is counted label
     included, after escaping. A line over 200 Unicode scalars keeps its
     first 200 and gains `…`. A line of exactly 200 is unchanged. An escape
     sequence is never split: one that would straddle the cut is dropped
     whole. This matches the A7 rule (81 scalars become 80 plus `…`).
  6. **Plain-copy refusal (R-12).** "Refused as today" (contract section 5)
     means each leaf's existing fallback failure, measured at `eb34889`:
     `invalid-argument`, exit 2, `argument:"input"`, `value:null`, with
     messages `invalid envelope input`, `invalid claim request`,
     `invalid Run request` and `invalid gate request`. Task 3's field-path
     rule covers structural failures only and does not change these. A
     dedicated worktree-binding failure belongs with parallel worktrees
     (CF-046, Phase 25).
  7. **Gate-receipt update (R-01).** The field-digest map is an insert value
     on every keyed Event. On the `gate.receipt-recorded` update it is a
     guard beside `request_digest`, never a value. That keeps the column
     immutable after insert, like `request_digest`. Archive
     `semanticUpdateColumns` and the update values stay unchanged.
  8. **Stop `Next:` states (R-10).** An open Stop prints a `Next:` line with
     its Condition command, then the after-the-fact step (`fid stop resolve`,
     then `fid attention resolve`). A resolved Stop prints its resolution
     Event and no `Next:` line.
  9. **Journey provenance (R-13).** Every value in the take-control journey
     has one of five classes:
     - `fid-output`: earlier Human-mode output. JSON output is not a source.
     - `fid-help`: `--help` or `--example` text.
     - `caller-authored`: new record IDs, idempotency keys, reasoning,
       display labels, and project files the caller writes.
     - `git`: a git command that the leaf's help names for that field.
     - `setup`: the acting Actor IDs and the project directory only.

     A value with no class is a product gap and a blocker candidate. It is
     never a reason to widen the classes.
  10. **Human fixtures (R-14).** Task 11 adds a tagged fixture generator,
      `TestP24HumanFixtures`, for the window-D walk-through. Task 14 extends
      it for the phase-end checklist. QA regenerates fixtures by running the
      test with an environment variable and writes no test code.
  11. **Field-digest privacy (I-17).** Field digests cover top-level public
      names only. The sensitive members (`environment[].value`,
      `redact_literals`) are nested inside `content`, which the
      whole-request digest already covers today, so exposure is unchanged.
      Task 9 asserts that no top-level descriptor field is `sensitive`.
  12. **Assertions.** Six bundled rows are split into new rows VAL-058..063
      (R-17). VAL-064 and VAL-065 cover grammar failures and unknown
      commands (R-07). VAL-066 covers "a replay never observes" (R-03).
      Existing IDs are not renumbered.
  13. **QA cadence.** The QA Lead's five windows are accepted: after tasks
      3, 6, 8, 11 and 14. The window after task 11 includes the Human
      walk-through.
- **Disposition register:**

| Item | Disposition | Applied in |
|---|---|---|
| R-01 | Do now (item 7). Task 1 adds the column, JSON-kind and update-guard registries, the direct gate-run Effects and three retained tests. R3 is corrected. | task 1; inventory R3 |
| R-02 | Do now. Per-keyed-leaf corpus over 49 leaves, enumerated from the new R16 with a completeness test. | VAL-002/006/007; task 2; inventory R16 |
| R-03 | Do now. New VAL-066 with a new-key control row. | VAL-066; tasks 2, 4-8 |
| R-04 | Do now (item 1). Resolved from accepted text, so it is not a blocker. | contract 2, 6; Shared Terms; VAL-009/017/022/023 |
| R-05 | Do now (item 2). Resolved within the closed domain, so it is not a blocker. | contract 6; task 4; VAL-023 |
| R-06 | Do now. Rows come from `cli.RegisteredCodes()`, with a pinned exclusion list. | VAL-036; task 10 |
| R-07 | Do now (items 3-4). The operation value awaits D5 (DEC-279). | VAL-064/065; tasks 9, 10; contract 9, 13 |
| R-08 | Do now. | VAL-038; task 10 |
| R-09 | Do now (item 5). | contract 9; VAL-039; task 10 |
| R-10 | Do now (item 8). Covers all 10 Condition types, open and resolved. | VAL-040; task 11 |
| R-11 | Do now. Same-worktree and plain-copy fixtures. | VAL-041; task 11 |
| R-12 | Do now (item 6). | VAL-027/028/029; tasks 6-8 |
| R-13 | Do now (item 9). Help for `expected_repository` names its git commands. | VAL-012/049; tasks 5, 14 |
| R-14 | Do now (item 10). | tasks 11, 14 |
| R-15 | Do now. Adds `operations.json` and its reader to task 13. R13 is corrected. | task 13; inventory R13 |
| R-16 | Do now. Declared `exit_code: 1`. | task 9 |
| R-17 | Do now (item 12). | VAL-058..063 |
| I-01, I-02 | Do now. Recursive walk that exits 0; corrected counts. | inventory 1, R15; VAL-010 |
| I-03, I-04, I-05, I-06 | Do now. | VAL-006/007/003/004 |
| I-07 | Record only. The 100k journal has no keyed Events. Its definition is protected (DEC-158). | task 1; VAL-005 |
| I-08, I-09, I-11, I-12 | Do now. | VAL-018/019/015/035/044 |
| I-10, I-16 | Do now. Expected statuses and names declared. | tasks 10, 2 |
| I-13 | Do now. Reports also record the slowest ordinary package. | VAL-052; all tasks |
| I-14 | Do now as a report note. No compatibility path is added (DEC-270). | task 1 |
| I-15 | Do now. | task 13 |
| I-17 | Do now (item 11). | task 9 |
| I-18, I-19 | Do now. Recorded in the phase brief. | phase brief |

- **Consequences:**
  - The contract gains DEC-278 clarifications in sections 2, 5, 6, 9 and 10,
    and D5 in section 13.
  - The validation contract has 66 rows (61 deterministic, 5 Human). All are
    `Ready` except VAL-065, which waits for DEC-279.
  - Tasks 1-14 are amended.
  - A QA Lead readiness refresh must revalidate every changed item before
    task 1 is dispatched.
- **Refs:** DEC-194; DEC-195; DEC-264; DEC-277; `qa/qa-plan.md` sections 2.3
  and 2.4; Phase 21 `public-interface.md` sections 2 and 4.

### DEC-279: Unknown-command `operation` names the receiving command (D5 = A)

- **Date:** 2026-10-04. ACCEPTED by the Human ("A") after the Tech Lead
  showed both options with example output.
- **Context:** Phase 21 accepted `operation-unavailable` for unknown
  commands, but no accepted text gave the JSON envelope's `operation` value
  for a command that does not exist (DEC-278 item 4).
- **Decision:** `operation` is the dotted path of the known command that
  received the unknown word. The root is `fid`. For example, `fid bogus`
  gives `fid`, `fid concern bogus` gives `concern`, and `fid issue ac bogus`
  gives `issue.ac`. The unknown word goes in `details.value`, with
  `argument:"command"` and `allowed` = the subcommand names that command's
  `--help` lists, sorted. Option B (the typed path, `concern.bogus`) was
  rejected because it puts caller text in a field that otherwise holds only
  known names.
- **Consequences:** VAL-065 is `Ready`, and task 9's dispatch gate is
  lifted. No other assertion or task changes.
- **Refs:** DEC-278; `public-interface.md` sections 9 and 13; task 9.

## Spec Sync Verdict

Recorded at wrap-up.

## Triage-Log Audit

Recorded at wrap-up.

## Workflow Notes

Recorded at wrap-up.
