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

### DEC-280: QA readiness-refresh dispositions; task 9 split

- **Date:** 2026-10-04. Tech Lead decision under DEC-277. It clarifies
  accepted text and adds no behaviour beyond it. RR-06 is the one item
  accepted text does not decide, so it goes to the Human as D6 (to be
  DEC-281).
- **Context:** the QA plan's readiness refresh (`qa/qa-plan.md` section 12,
  `87f1071`) was NOT READY. It raised RR-01..RR-09 (RR-06 is a blocker
  candidate) and RI-01..RI-13. The Tech Lead reproduced each finding against
  source and a scratch build of `87f1071`. The product tree there is
  identical to `304aa45`.
- **Decisions:**
  1. **`next:` on grammar failures (RR-01).** Contract section 9 governs:
     every `invalid-argument` on an `--input` leaf prints its row, grammar
     failures included. Task 10 was wrong, and now says the same thing. No
     contract change.
  2. **Unknown-command message (RR-02).** At every level the message is the
     single line `unknown command "<word>" for "<command path>"`, with no
     Cobra suggestions; `allowed` already lists the alternatives. The root's
     `allowed` includes Cobra's `completion` and `help`, because root
     `--help` lists them (RI-12). Those two Cobra commands have no operation
     ID, so they get no VAL-064/065 rows of their own.
  3. **Unknown-command baseline (RR-03; corrects DEC-278 item 4).** Measured
     at `87f1071`: an unknown word under 23 groups prints the group help and
     exits 0 in both modes. The root and five groups (`concern`, `export`,
     `issue`, `issue ac`, `run`) exit 2 silently. Task 9a installs one
     central mechanism on every command with subcommands (R17: the root plus
     28 groups) and replaces the five group-local validators. A group run
     with no word still prints its help. The exit change from 0 to 2 is
     authorized by Phase 21 `public-interface.md` section 5. DEC-278 cited
     that section as 4. The correct section is now cited in VAL-064/065 and
     tasks 9a and 10.
  4. **Grammar classes (RR-04).** These are pinned in task 9a.
     - (a) The unknown-flag `allowed` lists the leaf's local long flags,
       without `--help` or the globals `--json` and `--project`. Phase 21
       section 4 rows list only leaf flags.
     - (b) The malformed-value `allowed` is the flag's declared domain:
       enumerated members, sorted; an integer range written `lo..hi`; or the
       pattern its row states. Otherwise it is the pflag type name.
     - (c) A typed-flag parse failure, or a flag given without a value, is a
       malformed value. Today `issue list --has-concerns=maybe` exits 7
       silently.
     - (d) An extra word on a leaf without subcommands is a wrong positional
       count, never `operation-unavailable`. Its message is
       `<path> accepts no arguments`, after the Phase 21 section 1
       `capabilities` example.
     - An unknown flag on the root or on a group uses the unknown-flag class
       with the DEC-279 `operation` and `allowed:[]`.
     - The `message` is the text of the error that detected the failure:
       Cobra's or pflag's, or the flag-domain check's. It is never a family
       fallback text. This replaces task 9's rule "the text Human mode
       printed". That rule would have pinned `invalid Issue mutation
       request` for `issue update --bogus`, because `issueMutationFlagError`
       prints that text today.
  5. **Flag domains and proof (RR-05).** Each flag with a typed or closed
     value declares its domain once in `cmd/fid`. The check runs before
     workspace resolution, and the corpus is enumerated from the
     declarations. Completeness covers every non-string pflag flag (R18:
     22) and every string domain that a Phase 21-24 row states. Each row
     also runs from a non-project directory, which proves "no read".
     - R19 (`rg '"argument": *"request"'`) must terminate `exit_code: 0`
       with exactly the four family fallbacks left. Each is accounted for in
       the report. The grammar sites (`canonical_request.go`, `root.go`,
       `issueMutationFlagError`) must not match.
     - Flag-domain failures on read leaves also reach `argument:"input"`
       fallbacks (`claim list --limit 0`, `run list --status bogus`). The
       corpus, not the `rg`, is therefore the completeness proof.
     - Flag-combination rules are outside the four classes and keep today's
       failures (CF-047). They are exclusive selectors, exactly-one or
       at-least-one groups, and kind-conditional flags.
  6. **Control characters in messages (TL-01, found during this check).**
     `fid stop raise $'--bo\x1bgus'` writes a raw ESC to stderr today,
     because pflag's `unknown flag: --%s` is unquoted. Section 9's verbatim
     `error:` line would keep doing this, against PRD-OPS-007. Any message
     that embeds caller text escapes control characters when it is built,
     using one shared `internal/cli` function in the existing
     `visibleLiveText` form. The JSON message and the `error:` line stay
     byte-equal. Task 9a escapes the grammar and unknown-command messages.
     Task 10's renderer uses the same function, and its pairing corpus
     asserts that no message holds a control character.
  7. **Post-frame `gate run` (RR-07).** VAL-036 states the stdout exception
     and keeps a `gate-unsuccessful` row. Its JSON side is the failure inside
     the `--json` run's terminal frame.
  8. **Exclusions (RR-08).** VAL-064 excludes `journal inspect` and `sql`
     (Phase 21 section 5), plus Cobra's `help` and `completion`.
  9. **FIFO fixture (RR-09).** Git lists no untracked FIFO. The fixture
     replaces a tracked file with a FIFO. Task 4 also pins
     `invalid unmerged path` as `unmerged-paths`, decided by capture stage
     (RI-02).
  10. **Task 9 split (TL-02).** The Phase 23 sizing lesson applies: split
      any task that crosses several closed registries. Task 9 keeps the
      descriptor and input conformance work. New task 9a owns VAL-064/065
      and runs between tasks 9 and 10. Lettered tasks have a precedent in
      Phase 20 tasks 3a-3i, and the sequence lint accepts them. Window D
      now spans tasks 9, 9a, 10 and 11. Its cadence is the QA Lead's call.
  11. **D6 (RR-06).** Today every generic `invalid-argument` fallback prints
      its domain cause in Human mode. Section 9 would replace that cause
      with the fixed JSON message. Accepted text does not decide this, so
      it is Human decision D6 (contract section 13). Task 10 is not
      dispatched before DEC-281.
  12. **Journey start state (RI-08).** The VAL-049/053 start state has no
      Receipt-backed Obligation, and journey input files are written
      outside the worktree, so real freshness cannot stale the journey.
      `evidence-stale` has no next-table row. That gap is recorded under
      CF-040, because a closed-table change needs the Human.
- **Disposition register:**

| Item | Disposition | Applied in |
|---|---|---|
| RR-01 | Do now (item 1). Task 10 aligned to section 9. New VAL-038 row: a grammar failure on an `--input` leaf | task 10; VAL-038 |
| RR-02 | Do now (item 2). Near-miss rows (`fid isue`) | contract 9; task 9a; VAL-036, VAL-065 |
| RR-03 | Do now (item 3). R17 completeness check | task 9a; VAL-065; inventory R17 |
| RR-04 | Do now (item 4) | task 9a; VAL-064 |
| RR-05 | Do now (item 5). R18, R19 and a flag-domain corpus | task 9a; VAL-064; inventory R18, R19 |
| RR-06 | Escalated to the Human as D6 (item 11). Task 10 waits for DEC-281 | contract 13 |
| RR-07 | Do now (item 7) | VAL-036; task 10 |
| RR-08 | Do now (item 8) | VAL-064; task 9a |
| RR-09 | Do now (item 9) | task 4; VAL-023 |
| RI-01 | Do now. The `artifact` two-field row is `artifact_id` plus `--actor` | VAL-006; task 2 |
| RI-02 | Do now (item 9) | task 4 |
| RI-03 | Record only. An unreadable untracked file stays `changed-during-capture`, and `detail` names it | task 4 note |
| RI-04 | Do now. A control would succeed from the same pre-state in an observable worktree | VAL-066; tasks 2, 4-8 by reference |
| RI-05 | Do now. No shared decoder may match the R5 pattern | tasks 3, 9 |
| RI-06 | Do now. The section 1 block writes one labelled file per list | inventory 1 |
| RI-07 | Do now | validation contract Coverage |
| RI-08 | Do now (item 12) | task 14; CF-040 |
| RI-09 | Do now. The first failing test is the fixture manifest listing the start state | task 14 |
| RI-10 | Do now. Task 12 depends on 4 and 10 | phase brief; task 12 |
| RI-11 | Do now. `issue list --bogus` and `issue get` with no ID are silent in both modes. QA's `issue show` is an unknown word under `issue` (an RR-03 case), not a missing ID | this entry; task 9a |
| RI-12 | Do now (item 2) | task 9a; VAL-065 |
| RI-13 | Do now. `/tmp` had 0.7 GB free at this check | phase brief |
| TL-01 | Do now (item 6) | contract 9; tasks 9a, 10; VAL-036, VAL-064 |
| TL-02 | Do now (item 10) | task 9, new task 9a; tasks 10, 12; phase brief |

- **Consequences:**
  - CF-047 is opened. CF-040 gains the `evidence-stale` next-row note.
  - The validation contract keeps 66 assertions. VAL-064/065 move to task
    9a.
  - Tasks 1-3 change only by RI-01 (task 2) and RI-05 (task 3).
  - Nothing is dispatched until a fresh QA Lead check passes. Task 10 also
    waits for DEC-281.
- **Refs:** DEC-194; DEC-278; DEC-279; Phase 21 `public-interface.md`
  sections 1, 2, 4, 5 and 7; `qa/qa-plan.md` sections 12.3 and 12.4.

### DEC-281: Generic fallbacks carry the domain cause as their message (D6 = A)

- **Date:** 2026-10-04. ACCEPTED by the Human ("A") after the Tech Lead
  showed both options with example output (DEC-280 item 11).
- **Context:** each family's default `invalid-argument` branch catches
  cross-field domain failures with a fixed message (`invalid claim
  request`), `value:null` and a generic `allowed`. Today Human mode prints
  the domain cause instead. Contract section 9 makes the `error:` line the
  JSON message, so the cause would have vanished from both modes.
- **Decision:** a **generic fallback** is a family failure mapper's final
  `invalid-argument` branch for a cause no earlier branch matched. There
  are 18 at `62dddab` (inventory R20). Its JSON `message` is the caught
  error's text (`cause.Error()`, which Human mode prints today). The text is
  control-character-escaped by the DEC-280 function. A text over 512
  scalars keeps its first 511 and gains `…`, with no escape split; this is
  the Phase 21 `message` bound. An empty or all-whitespace cause keeps the
  fixed message. Code, details, subject, category, retryable and exit are
  unchanged. Option B (accept the loss) was rejected.
- **Consequences:**
  - The plain-copy refusals (contract section 5, DEC-278 item 6) keep
    their code and details. Their messages become the domain reasons at the
    refusing sites today, which tasks 6-8 keep:
    - `envelope.create`: `prior registration binding does not match envelope`;
    - `claim.acquire`: `claim acquire binding is invalid`;
    - `run.start`: `run claim binding is invalid` (corrected by DEC-282;
      the original pin named the later Envelope check);
    - `gate.run`: `run gate: repository binding is invalid`.
  - Task 10 delivers the rule, and new VAL-067 proves it. Tasks 6-8 assert
    today's fallback messages. Task 10 converges those assertions, plus
    `cmd/fid/issue_update_test.go` and `cmd/fid/generic_run_test.go`, with
    provenance.
  - Messages are never control flow (Phase 21 section 7). No cause may
    carry a `sensitive` descriptor value, and VAL-067 includes a guard row.
  - The plain-copy `next:` line still points at the input. A dedicated
    worktree-binding failure stays CF-046 (Phase 25).
  - Task 10's dispatch gate is lifted. The validation contract has 67
    assertions.
- **Refs:** DEC-278 item 6; DEC-280; `public-interface.md` sections 5, 9 and
  13; Phase 21 `public-interface.md` sections 1 and 7.

### DEC-282: QA recheck corrections for the Run cause and bounded grammar messages

- **Date:** 2026-10-04. Tech Lead clarification under DEC-277 and the Human's
  request to address QA recheck RC-01/RC-02. Preserves D6=A's actual-cause
  rule, Phase 21's 512-scalar message bound and DEC-280's escaping rule.
- **Context:** QA commit `4819b30c935acdb6827c66316f4f3383343bc76f`,
  `qa/qa-plan.md` sections 12.6.3-12.6.5, reports two bounded contract
  contradictions and retains a NOT READY verdict. The Tech Lead read that
  evidence and inspected `internal/workflow/run.go`'s Claim-before-Envelope
  checks, the existing escape form, and the inherited message contract.
  QA's runtime probes remain attributed to QA; no runtime probe or quality
  gate was run or re-attested in this Tech Lead session.
- **RC-01 — Do now: correct the pin.** The plain-copy `run.start` cause is
  `run claim binding is invalid`. The Claim worktree check rejects before
  the Envelope check. Correct DEC-281's consequence, contract section 5,
  VAL-029 and tasks 7/10 together. Preserve both checks and their order.
  Task 7 proves original-worktree acceptance, import-binding acceptance
  and plain-copy refusal with an otherwise valid active Claim; task 10
  pairs the actual cause in Human/JSON output and converges retained
  message assertions with provenance. A same-worktree snapshot mismatch
  is not a substitute for the plain-copy case. Codes, details, exit and
  no-Event guarantees retain their accepted values.
- **RC-02 — Do now: specify bound precedence.** Grammar failures and
  unknown commands use the same bounded message construction as DEC-281
  generic fallbacks. Construct the specified parser/check/formula text,
  escape controls, then count Unicode scalar values in the resulting text
  before JSON serialization (excluding the Human label and line terminator).
  Up to 512 scalars stays intact. Over 512 keeps the first 511 plus `…`;
  drop an escape sequence that straddles the cut whole, then append `…`
  without padding. Existing `%q` escapes stay escaped and are not split.
  This rule takes precedence over DEC-278/280's exact message-source
  wording. Full offending spelling/value remains in structured details;
  code, category, retryable, subject and exit are unchanged. The Human
  `error:` payload equals the decoded JSON message without a second escape
  or cut. This neither expands the message bound nor changes domain checks.
- **Ownership and evidence:** task 9a owns shared construction and grammar/
  unknown-command producers, with short exact text, 600-character flag/word,
  escaped-text 512/513, non-ASCII and escape-crossing-cut cases (both modes,
  plus shared-construction units). Task 10 reuses that construction for
  generic fallbacks and repeats pairing after adopting its renderer.
  VAL-064/065 own grammar behavior; VAL-036 owns final Human pairing.
  VAL-067 explicitly includes the already required whitespace and combined
  escape/cut cases and a sentinel guard reaching the generic fallback.
  All affected rows keep QA Reviewer as their sole validator. No new
  assertion, task or test-file scope is needed: tasks 7/9a/10 already own
  these source/test locations and require both ordinary race/coverage and
  tagged integration gates, which execute the respective owned suites.
- **Other recheck findings:** adopt the QA Lead's D1 after task 9a and D2
  after task 11 cadence; the Human walkthrough stays after task 11. Correct
  the brief to R1-R20. R20's 18 fallbacks/seven exclusions and the two
  retained literal-message assertions need no allocation change.
- **Artifact/spec sync:** amended `public-interface.md` sections 5/9,
  `qa/validation-contract.md` (029/036/064/065/067 and revision log),
  tasks 7/9a/10, and current brief/progress. Canonical
  `06_DESIGN_PLAN_CLI.md` gains a Phase 24 accepted-amendment entry linking
  these message rules, explicitly pending implementation. The Phase 21
  baseline bound stays intact. No product/data-model/glossary change or new
  carry-forward is required; the full-suite spec audit remains phase-end work.
- **Next:** QA Lead rechecks RC-01/RC-02 and this amendment's changed text,
  recording fresh evidence in a new QA-plan subsection. The prior QA plan
  and its NOT READY verdict remain intact. The phase-wide dispatch hold
  stays until QA recheck passes and the Runner completes readiness checks.
  This decision supplies the requested contract authority, not QA acceptance.
- **Refs:** DEC-277/280/281; `qa/qa-plan.md` section 12.6;
  `public-interface.md` sections 5/9; Phase 21 `public-interface.md`
  sections 1/2/5/7; PRD-OPS-007/014.

### DEC-283: Reconcile the Task 6 STOP-resume recovery tuple; bounded `workflow-sequence-lint.sh` repair; Task 6 retained-fixture scope. Closes ESC-24-001 and ESC-24-002

- **Date:** 2026-10-06. Tech Lead orchestration-integrity decision on the
  Human's instruction to fix the current blocker. This is the separate
  protected authority that `AGENTS.md` and `09_AI_WORKFLOW.md` section 2.3
  require for a `workflow-sequence-lint.sh` change. It follows DEC-148/149.
- **Context:**
  - The 2026-10-05 20:48 RUN-RESUME resolved ESC-24-001 with
    `restored-capability:developer-role-spawn-01a10b08-efe6-7050-8f88-65c524e62519`.
    A Developer session then committed implementation
    `1baaecc7983872b06e2a47e6df0801bfb1370ea0` and report
    `7298964a65beb375a6b50ec45c42ea12dada6c59` (report-lint
    `OUTCOME RECOVERY_REQUIRED`, exit 0). Every path `1baaecc` touches is in
    the 2026-10-05 20:06 TASK-RECOVERY allowlist.
  - That RUN-RESUME contradicts itself. Its `Trigger signature` and
    `Attempt count` are `task-6:session-liveness:developer-role-capacity-after-bounded-retry`
    and `3`, as the verifier requires: TRG-24-003's trigger, attempt plus
    one (`:1049-1054`). Its `Next recovery` instead names
    `task-6:session-continuation:retained-gate-envelope-claim-fixtures-and-nullable-registration`
    at attempt `1`. The bound TASK-RECOVERY must match both (`:723-727`,
    `:742-746`), so no candidate exists. Tech Lead dry runs on a scratch
    copy: the RUN-RESUME-matching successor fails only `TASK-RECOVERY does
    not match the exact schema v2 RUN-RESUME recovery tuple`; the
    tuple-matching one fails only `TASK-RECOVERY does not preserve its
    RUN-RESUME checkpoint`. While `resume_recovery_pending` is set
    (`:473-474`), nothing else may follow, not even TRIAGE, STOP or
    CORRECTION. This is the wedge DEC-149 clause 6 described. The ledger
    still validates `VALID state=REPAIR_PENDING`.
  - Causes: (1) the two spawn failures were routed but never recorded as
    `REPAIR ... INTERRUPTED`, so TRG-24-003 was not a session-liveness
    TRIAGE that keeps the chain trigger (`:1024-1029`). TL-Lite minted a
    new trigger, whereas Phase 23's TRG-23-013 kept the task's. (2) The
    Runner wrote the tuple from ESC-24-001's "preserving the trigger"
    prose, with an attempt that did not advance. (3) The STOP-resume
    branch (`:1034-1100`) never compares `Next recovery` with
    `Continuing task`, `Trigger signature` and `Attempt count`, although
    `phase-runner.md` and section 2.4 say it does. (4) The Runner
    dispatched the Developer before appending the bound TASK-RECOVERY.
  - All 32 phase 20-24 STOP resumes with a tuple were scanned. 30 restate
    their own task, trigger and attempt. Phase 23's 2026-09-26 17:38 is a
    DEC-257 owner-route resume whose first-owner trigger at attempt 1 is
    legitimate. Phase 24's 20:48 is the only defect.
  - Task 6's report records both full gates failing at 21 retained fixture
    definitions outside the 20:06 allowlist, each with introducing-commit
    provenance, and flags `internal/claims/generic.go:380`.
- **Decision:**
  1. **Stop sustained.** Refusing to rewrite the accepted RUN-RESUME or to
     append an invalid STOP was correct. ESC-24-002 is allocated to this
     defect and closed here. ESC-24-001 is closed: the 20:48 RUN-RESUME
     and the later session verify restored capability.
  2. **Accepted semantics.** In a STOP resume without `Repair owners`, the
     verifier-checked `Continuing task`, `Trigger signature` and
     `Attempt count` are the successor's identity. `Next recovery` restates
     them and adds category and repair mode. In the 20:48 tuple the trigger
     and attempt are the defect; its task, category and mode stand. Task 6's
     chain trigger moves to the session-liveness signature at this boundary.
     That is recorded here so it is not later read as a fabricated identity
     change, as DEC-148 did for task-3k. Attempts stay monotonic: 1, TRG 2,
     then 3. No-progress restarts at `0/2`, equal to the pre-STOP count.
  3. **Reconciliation record.** The verifier reads these two lines exactly:

- Reconciles RUN-RESUME: 2026-10-05 20:48 / TRG-24-003
- Reconciled next recovery: task-6 | in-scope-implementation | task-scope | task-6:session-liveness:developer-role-capacity-after-bounded-retry | 3

  4. **Bounded verifier repair (the only authorized change).** In the
     schema-v2 STOP-resume RUN-RESUME branch, without `Repair owners`:
     (a) Split `Next recovery` on ` | `. It is *consistent* when it has
     exactly five fields: task equals `Continuing task`, category is in the
     `:650` set, mode is in the `:653` set, trigger equals
     `Trigger signature`, and attempt equals `Attempt count`.
     (b) An inconsistent tuple is accepted only through a reconciliation.
     Exactly one `### DEC-NNN:` section of the phase decisions file must
     hold a trimmed line `Reconciles RUN-RESUME` (bullet form, as above)
     whose value is `<heading date> <heading time> / <Governing TRG>`. That
     section must hold exactly one `Reconciled next recovery` line. Its tuple
     must be consistent under (a), and its fields 2-3 must equal the recorded
     tuple's. That tuple then binds the successor, and the DEC becomes a
     required citation. Anything else is INVALID `schema v2 RUN-RESUME Next
     recovery does not retain its Continuing task, Trigger signature, and
     Attempt count`.
     (c) A reconciliation keyed to a consistent RUN-RESUME is INVALID
     `RUN-RESUME next-recovery reconciliation targets a consistent Next
     recovery`, so a DEC cannot replace a valid binding.
     (d) The bound TASK-RECOVERY after a reconciled RUN-RESUME must cite the
     reconciling DEC in both `Source` and `Repair authority`; otherwise it is
     INVALID `TASK-RECOVERY does not cite its RUN-RESUME next-recovery
     reconciliation`. The checks at `:714-746` are unchanged; the exact-tuple
     comparison uses the reconciled tuple.
     (e) Owner-route and normal resumes, CORRECTION (evidence only), TRIAGE
     trigger rules, every other message, and every other state transition
     are unchanged.
  5. **Execution** is Tech Lead-authorized workflow maintenance outside the
     run sequence (DEC-148 clause 5, DEC-149 clauses 8-9). It has no
     TRIAGE, no control TASK-RECOVERY or REPAIR, and no ledger event before
     the bound TASK-RECOVERY. The Runner dispatches a Developer into a clean
     temporary worktree at the directive base and integrates the single
     verified commit with `git merge --ff-only`. The Runner-owned ledger stays
     byte-unchanged in the main worktree until then.
  6. **Ledger sequence after integration**, each through the transactional
     append:
     (a) the bound TASK-RECOVERY: task-6, in-scope-implementation,
     task-scope, the session-liveness trigger, attempt 3, `0/2`. `Source` and
     `Repair authority` cite TRG-24-003, ESC-24-001, DEC-283 and the
     restored-capability token. `Allowed files` and `Forbidden changes` are
     byte-identical to the 20:06 entry, since that directive is what the child
     received. Its `Checkpoint` starts at 3ed4824/c845cee on phase HEAD
     70e2106 and cites the maintenance commit. Its `Resolution evidence`
     discloses the dispatch-before-append deviation and both rejected
     candidates with their messages and exit status.
     (b) `REPAIR task-6`, attempt 3, `MATERIAL-PROGRESS`, ending at
     1baaecc/7298964 with report-lint `OUTCOME RECOVERY_REQUIRED`.
     (c) TASK-RECOVERY task-6, same trigger, attempt 4, `0/2`, sourcing the
     report and DEC-283, using clause 7's allowlist. Then dispatch.
     Tech Lead dry runs of (a)-(c) on a scratch ledger copy, with the binding
     reconciled, each returned `APPEND_CONFIRMED state=REPAIR_PENDING`.
  7. **Task 6 scope (attempt 4).** The 21 report paths are retained fixtures
     whose inputs Task 6's accepted removals displaced (gate.register
     identities, envelope.create Source repository and real bytes,
     envelope.complete `repository`, claim.acquire `worktree_id`), or which
     copy observed history to another root. Sampled introducing commits
     resolve to the recorded phase/task. None comes from a later Phase 24
     task. The seven untagged files run under the ordinary race/coverage
     gate and the fourteen `integration`-tagged files under the tagged gate,
     so both are Task 6 required commands. Convergence is Task 6's, with
     per-row provenance in the report. The exact `Allowed files` value is
     the 20:06 list followed by the 21 paths in report order:

     `internal/evidence/registration.go,internal/workflow/envelope.go,internal/workflow/import_binding.go,internal/claims/generic.go,cmd/fid/gate.go,cmd/fid/envelope.go,cmd/fid/claim.go,cmd/fid/input_descriptor_command.go,cmd/fid/input_descriptors_binding.go,cmd/fid/gate_test.go,cmd/fid/envelope_test.go,cmd/fid/generic_claim_test.go,cmd/fid/input_descriptor_test.go,cmd/fid/integration/p21_conformance_test.go,cmd/fid/integration/p22_conformance_test.go,cmd/fid/integration/p22_fake_test.go,cmd/fid/integration/p22_portable_conformance_test.go,cmd/fid/integration/p22_portable_fake_test.go,cmd/fid/integration/p23_conformance_test.go,cmd/fid/integration/p23_amended_envelope_oracle_test.go,cmd/fid/integration/p23_amended_claim_oracle_test.go,cmd/fid/integration/p23_claim_mode_test.go,cmd/fid/integration/p23_review_record_test.go,cmd/fid/integration/p23_review_leaf_oracle_test.go,cmd/fid/integration/p23_judgment_leaf_oracle_test.go,cmd/fid/integration/p23_portable_recovery_test.go,cmd/fid/integration/p23_leaf_surface_test.go,project-docs/specs/06_DESIGN_PLAN_CLI.md,project-docs/specs/07_DESIGN_PLAN_FLOWS.md,project-docs/docs/phases/phase-24/tasks/task-6-done.md,cmd/fid/generic_run_test.go,cmd/fid/evidence_test.go,cmd/fid/review_test.go,cmd/fid/task18_completion_resume_test.go,cmd/fid/task21_maintenance_rebuild_test.go,cmd/fid/stop_evaluation_matrix_test.go,internal/export/archive_fixture_test.go,cmd/fid/integration/p22_resume_actions_test.go,cmd/fid/integration/p22_recovery_journey_test.go,cmd/fid/integration/p23_review_evidence_budget_test.go,cmd/fid/integration/p23_obligation_guard_test.go,cmd/fid/integration/p23_human_outcome_test.go,cmd/fid/integration/p23_human_fixtures_test.go,cmd/fid/integration/p23_human_failure_test.go,cmd/fid/integration/p23_gate_grants_test.go,cmd/fid/integration/p23_decision_reference_test.go,cmd/fid/integration/p23_task13_precedence_test.go,cmd/fid/integration/p23_concern_deviation_test.go,cmd/fid/integration/p24_gate_evidence_corpus_test.go,cmd/fid/integration/p24_single_journal_test.go,cmd/fid/integration/p24_idempotency_test.go`

     `Forbidden changes` is the 20:06 text plus: build path-bound history in
     the root that uses it, never by copying observed history to another
     root; never create a fictitious import binding; never relax the shared
     predicate.
  8. **Stored Claim authority (contract section 5 detail).**
     `validateStoredClaimAuthority` compares the stored Claim worktree to the
     Envelope scope. Task 6 now stores the observed worktree, so a Claim
     acquired under VAL-029's imported-destination equivalence would be
     refused by its own renew, release or expire. Section 5 already says the
     destination's worktree stands for the source worktree recorded in
     imported history. That comparison therefore uses the shared
     `workflow.SameWorktree` predicate, never a copy. Inputs, the
     `claim worktree binding is invalid` reason, the code, the exit and
     no-Event behavior are unchanged.
     Red first, in an imported destination: acquire, then renew and release
     (expire too if a public path reaches it). Green after convergence, with
     the original-worktree control unchanged. If the red does not reproduce,
     record that and leave the line unchanged. Evidence is reported under
     VAL-029's Task 6 slice and reviewed at QA window B. No VAL row changes.
  9. No other product behavior, public interface, data model, dependency,
     VAL row, task boundary, gate command, QA cadence or Task 6 command
     changes.
- **Repair directive (verifier maintenance, clause 4):**
  - **Originating task:** task-6
  - **Control trigger signature:**
    `task-6:session-liveness:developer-role-capacity-after-bounded-retry`
  - **Control contract:** a STOP resume's `Next recovery` must be
    satisfiable by its own bound successor; a recorded defect is reconciled
    only by one DEC line pair, never by rewriting history.
  - **Control mismatch:** the STOP-resume branch accepts a tuple that
    contradicts its verified `Trigger signature` and `Attempt count`. The
    bound successor is then unsatisfiable, and the ledger admits no event.
  - **Authorized change:** clause 4 (a)-(e) only.
  - **Control repair base:** the phase-24 HEAD that records this decision.
    The Runner records it from `git rev-parse --verify HEAD^{commit}`,
    with only its own `run-ledger.md` dirty, and keeps it immutable.
  - **Allowed files:**
    `project-docs/tools/workflow-sequence-lint.sh,project-docs/tools/drills/run-drills.sh`
  - **Forbidden changes:** any run ledger, decisions file, triage log, task
    spec or report, product and test files, `report-lint.sh`,
    `workflow-control-scope-lint.sh`, `workflow-route.sh`, agent prompts,
    templates, specs, and any rule outside clause 4.
  - **Exact verification** (record verbatim output and observed `exit_code`):
    1. Red, gap: fixture F1 is a `phase-N` directory with a triage log and a
       decisions file. Its sequence is RUN-START; TASK-RECOVERY task-1
       (trigger A, attempt 1); TRIAGE BLOCK with trigger B, attempt 2, and
       an ESC; STOP; a restored-capability RUN-RESUME with trigger B,
       attempt 3, and `Next recovery` naming trigger A, attempt 1. It has no
       reconciliation. Before the change it is `VALID state=REPAIR_PENDING`.
       After, it is INVALID with the clause 4(b) message.
    2. Red, wedge: F2 is F1 plus a reconciling DEC and the bound successor
       (trigger B, attempt 3, the recorded category and mode, `0/2`), citing
       the TRG, ESC, DEC and restored-capability token. Before the change
       its only diagnostic is `TASK-RECOVERY does not match the exact schema
       v2 RUN-RESUME recovery tuple`. After, it is `VALID
       state=REPAIR_PENDING`.
    3. Guards after the change, each INVALID with the named message:
       F2's successor without the DEC citation (4d); a successor at trigger
       A, attempt 1 (`TASK-RECOVERY does not preserve its RUN-RESUME
       checkpoint`). Each of these RUN-RESUME reconciliation faults gives
       the 4(b) message: a wrong key time or TRG, a wrong reconciled attempt,
       a changed category or mode, two DECs with the key, and a missing or
       duplicated tuple line. A reconciliation of a consistent resume gives
       the 4(c) message.
    4. The existing consistent STOP-resume, DEC-257 owner-route,
       normal-resume and DEC-148/149 drills are unchanged, and
       `sh project-docs/tools/drills/run-drills.sh` ends `exit_code: 0`.
    5. Regression: the base verifier, materialized from the base into fresh
       scratch storage, and the repaired verifier give byte-identical
       combined output and exit status on the real phase 16-23 ledgers.
       The 2026-10-06 Tech Lead baseline: 16-19 `VALID state=STOPPED`; 20
       `INVALID REVIEW PASS must use Owning task: none: review-7b`; 21-23
       `VALID state=ADVANCEABLE`. Both verifiers return
       `VALID state=REPAIR_PENDING` on
       `/home/mark/projects/fid/project-docs/docs/phases/phase-24/run-ledger.md`
       (read-only, by that absolute path, because the worktree copy lacks the
       uncommitted 20:48 entry).
    6. `sh project-docs/tools/workflow-lint.sh` exits 0. It does at this
       base.
    7. Scope, proved directly as in DEC-149 item 7:
       `git diff --name-only --no-renames <base> <head>` prints exactly the
       two allowed paths; `git status --short` is empty; one commit on
       `<base>`; regular files with no mode change.
- **Consequences:** Once (a)-(c) are confirmed, Task 6 continues. QA window B
  follows Task 6 completion as planned. If item 1 or 2 cannot be produced
  as stated, that is new information: return to the Tech Lead, not to a
  fixture substitution (DEC-149).
- **Workflow notes for wrap-up:** (i) during an active TASK-RECOVERY, record a
  spawn failure as `REPAIR ... INTERRUPTED` so its TRIAGE keeps the chain
  trigger (DEC-148 clause 7's namespace asymmetry recurred); (ii) append the
  bound TASK-RECOVERY before dispatch; (iii) an owner-route tuple
  satisfiability check, and a legal record for a second STOP while
  resume-pending (DEC-149 clause 6), remain open.
- **Artifact/spec sync:** `09_AI_WORKFLOW.md` section 2.4 and
  `templates/run-ledger-template.md` gain the satisfiable-tuple rule and the
  reconciliation line pair. `public-interface.md` section 5 gains the stored
  Claim comparison. `tasks/task-6.md` gains a DEC-283 amendment. The ESC
  register, archive, brief and progress are refreshed. No PRD, TRD, data
  model, glossary or validation-contract change.
- **Refs:** TRG-24-003; ESC-24-001; ESC-24-002; DEC-148/149/150/257/277/278;
  `run-ledger.md` 2026-10-05 20:06 TASK-RECOVERY, 20:13 TRIAGE, 20:14
  STOP, 20:48 RUN-RESUME; `tasks/task-6-done.md` at `7298964`;
  `workflow-sequence-lint.sh:473-474`, `:650-654`, `:714-746`,
  `:1024-1029`, `:1034-1100`, `:1312-1345`; `internal/claims/generic.go:380-388`;
  `09_AI_WORKFLOW.md` sections 2.3-2.4, 9.

### DEC-284: Correct unconsumed TRG-24-006 through a DEC-190 correcting TRIAGE; no verifier change; late QA window C. Closes ESC-24-003

- **Date:** 2026-10-07. Tech Lead ESC resolution on the Human's instruction
  "Resolve ESC-24-003". Precedents: DEC-190 (Phase 22 correcting TRIAGE and
  its contract), DEC-283 clause 2 (chain trigger is identity).
- **Context:**
  - Ledger tail: 02:17 TASK-RECOVERY task-9 (`report-repair`, trigger
    `task-9:report-lint:noncanonical-pass-evidence-termination`, attempt 1);
    02:20 REPAIR `MATERIAL-PROGRESS`; 02:24 TRIAGE TRG-24-006 `REPAIR` with
    trigger `task-9:scope-gap:shared-descriptor-variant-adoption`, attempt 2,
    `Repair category` `scope-gap`. The ledger validates
    `VALID state=REPAIR_PENDING` at `f680b20`.
  - TRG-24-006 is right on substance (the four-file shared extension) and
    wrong on identity, twice. (1) `scope-gap` is not a recovery category
    (`run-ledger-template.md:155`). (2) It minted a new trigger while the
    report-repair continuation was active. Section 2.3 says REPAIR keeps the
    task's trigger budget, and the DEC-190 rule (section 9) says "an active
    continuation still retains its trigger and next attempt". The verifier
    checks neither at TRIAGE time, but enforces both at TASK-RECOVERY time
    (`workflow-sequence-lint.sh:710`, `:751-757`, `:817-819`). This is the
    DEC-283 defect class again: an accepted record with no satisfiable
    successor. DEC-283 workflow note (i) predicted the trigger half.
  - The in-section "Correction to TRG-24-006" is not a form the accepted
    contract parses. `section_field_matches` reads only the original labels,
    and the H3 sits inside the TRG-24-006 section. The accepted correction
    form is DEC-190's later correcting REPAIR TRIAGE for the same task and
    attempt (`pending_triage_reconciliation`, `:488-511`).
  - TRG-24-007's BLOCK was never recorded in the ledger, and could not be,
    because a pending TRIAGE REPAIR admits no TRIAGE BLOCK. No STOP exists.
    That is exactly the DEC-190 situation ("neither rejected BLOCK nor STOP
    was persisted").
  - Tech Lead dry runs were made on scratch copies of the phase files with
    the unchanged verifier at `f680b20` (results in the clause 5 table).
    Both TL-Lite rejections reproduce. No TASK-RECOVERY variant satisfies
    the current state, and TRIAGE BLOCK is rejected. The DEC-190 route
    appends at every step, and its guards reject.
  - Separately, QA plan section 8 window C (after task 8, tasks 7-8, plus
    feasibility of the commands tasks 9-14 prescribe) never ran. The ledger
    goes from 01:51 `TASK task-8 COMPLETE` to 02:17 TASK-RECOVERY task-9, and
    `qa/review-5.md` does not exist. No decision waived it.
- **Decision:**
  1. **ESC-24-003 closed; BLOCK premise rejected.** TL-Lite was right
     not to self-repair a protected verifier or rewrite history. But no
     verifier change is needed, so none is authorized. Parsing in-section
     corrections would add a second correction form beside DEC-190's,
     which is a protected definition change, so it is rejected.
     TRG-24-007 stays immutable as a triage-log-only BLOCK. It gets no
     ledger event, STOP, control amendment or STOP-resume RUN-RESUME.
  2. **Binding semantics.**
     (a) The effective external authority is **TRG-24-008**: a Tech Lead
     correcting TRG under this decision, written like TRG-22-005. It
     copies TRG-24-006's operative fields byte-for-byte: accepted contract
     and mismatch, `Allowed files`, recovery scope, `Forbidden changes`,
     required evidence, `Exact verification`, `Originating report` and
     independent re-review. Product scope is therefore unchanged: the
     `task-9.md` ordinary scope plus exactly
     `internal/cli/input_descriptor.go,internal/cli/input_render.go,internal/cli/input_validate.go,internal/cli/input_validate_test.go`.
     (b) TRG-24-006, its in-section correction, and the 02:24 ledger TRIAGE
     stay unedited. They are superseded, not rewritten. The in-section
     correction is evidence only and binds nothing.
     (c) Identity: `task-9`; trigger
     `task-9:report-lint:noncanonical-pass-evidence-termination`; attempt
     `2`; category `in-scope-implementation`; authority
     `task spec; DEC-284`. The trigger is chain identity, not a
     description (DEC-283 clause 2). The product defect keeps its
     descriptive name `task-9:scope-gap:shared-descriptor-variant-adoption`
     in TRG-24-006/008 prose only.
     (d) Budget: attempts are 1 (report repair), then 2. The 02:20
     `MATERIAL-PROGRESS` gives `0/2`, so nothing resets and nothing is
     consumed by rejected candidates.
  3. **Route.** This is ordinary task-scope continuation in the phase
     worktree. No verifier, drill, scope-verifier, route or report-lint
     change, no out-of-sequence maintenance, and no temporary worktree. The
     "ledger-only dirt" precondition that TL-Lite cited governs isolated
     control-conformance dispatch, which this route does not use. This
     decision's commit still includes the triage log, so the Runner-owned
     ledger is the only dirty path.
  4. **Ledger sequence** (Runner, transactional append; each step requires
     `APPEND_CONFIRMED state=REPAIR_PENDING` with observed `exit_code: 0`;
     append before dispatch):
     (a) A new Runner session first appends a normal schema-v2 RUN-RESUME:
     STOP-only fields `N/A`, `Continuing task: task-9` (DEC-190 runner
     restart).
     (b) `TRIAGE TRG-24-008` with exactly these five values:
     `Originating task: task-9`,
     `Trigger signature: task-9:report-lint:noncanonical-pass-evidence-termination`,
     `Attempt count: 2`, `Repair category: in-scope-implementation`,
     `Repair authority: task spec; DEC-284`, and `Verdict: REPAIR`.
     (c) `TASK-RECOVERY task-9` with the same task, trigger, attempt,
     `Category`, and `Repair authority`. `Repair mode: task-scope`,
     `Triage disposition: N/A`, `0/2`, `Source` citing TRG-24-008 and
     DEC-284, control fields and base `N/A`. `Allowed files` is
     `ordinary task-9.md scope plus exact shared extension ` followed by
     the four paths. `Forbidden changes`, `Exact verification` and
     `Originating report` are copied byte-for-byte from TRG-24-008.
     `Resolution evidence` discloses TRG-24-007/ESC-24-003, DEC-284, and
     both rejected candidates with their messages and `exit_code: 1`.
     Then dispatch the same-task Developer continuation.
     Its later REPAIR uses the same trigger and attempt.
  5. **Verification contract.** Positive: (a)-(c) each give the receipt
     above, and the ledger then validates `VALID state=REPAIR_PENDING`. The
     Tech Lead dry run gave exactly that both with and without (a).
     Negative guards from the same dry run, each `exit_code: 1` with the
     ledger byte-unchanged:

     | Candidate | Diagnostic |
     |---|---|
     | Rejected corrected candidate (product trigger, `in-scope-implementation`) | `does not preserve its governing schema v2 TRIAGE REPAIR authority`; `continuing recovery must preserve task and Trigger signature` |
     | Rejected original candidate (`scope-gap`) | `has missing or invalid Category`; `continuing recovery must preserve task and Trigger signature` |
     | Chain trigger with either category, before TRG-24-008 | `does not preserve its governing schema v2 TRIAGE REPAIR authority` |
     | TRIAGE BLOCK TRG-24-007 | `TRIAGE REPAIR must be followed by its bound TASK-RECOVERY before TRIAGE` |
     | TRG-24-008 with the product trigger | the line above plus `schema v2 TRIAGE identity does not match its unique external record` |
     | TRG-24-008 whose category differs from its record | `schema v2 TRIAGE identity does not match its unique external record` |
     | After TRG-24-008: product trigger; `Source` without TRG-24-008; authority `task spec`; attempt 3 | `does not preserve its governing schema v2 TRIAGE REPAIR authority` (product trigger and attempt 3 also give their continuation messages) |
     | After TRG-24-008: `1/2` | `continuing recovery has invalid Consecutive no-progress attempts; expected 0/2` |

     Immutable history: every existing ledger event, TRG-24-001..007, the
     in-section correction and both rejected candidates stay as recorded.
     If (a)-(c) as written is rejected, that is new information: return
     to the Tech Lead. Never substitute field values.
  6. **QA window C, late.** The Runner's omission is a deviation, not a
     cadence change. The window's scope is unchanged. A window C `REPAIR`
     naming tasks 7 or 8 cannot be routed while Task 9's recovery is open
     (overlapping-recovery rule). So the Runner runs window C immediately
     after `TASK task-9 COMPLETE` and before dispatching task 9a, as the
     next unused `qa/review-M.md`. It reviews tasks 7-8 at their recorded
     close commits plus current HEAD, and feasibility for tasks 9a-14.
     Task 9's own gates evidence its own feasibility. D1 (tasks 9, 9a)
     is unchanged. Dry run: `REVIEW` PASS after Task 9 completion appends
     `APPEND_CONFIRMED state=ADVANCEABLE`.
  7. **Recurrence guard (spec sync, clarifying only).** Section 2.3 and
     `tl-lite.md` now state three accepted rules explicitly: a continuation
     TRG keeps the chain trigger and the next attempt; `Repair category` is a
     ledger recovery category; and an unconsumed ordinary REPAIR TRG is
     corrected only by a later correcting TRG. No rule changes.
  8. No other product behavior, public interface, VAL row, task boundary,
     gate command or budget changes. Task 9 gains a DEC-284 amendment
     that records the TRG-24-006/008 extension. The required ordinary
     race/coverage gate and `golangci-lint run ./...` already execute
     those four files.
- **Workflow notes for wrap-up:** (i) This is the second identity defect at
  a continuation TRIAGE in this phase, after DEC-283 note (i). (ii) A
  candidate hardening is to fail fast at TRIAGE: reject a REPAIR whose
  category is outside the enum, or whose trigger differs from an active
  continuation's, rather than wedge at TASK-RECOVERY. This is protected
  verifier maintenance and is not authorized here. (iii) The ledger has no
  legal record for a BLOCK while a TRIAGE REPAIR is pending (compare DEC-149
  clause 6 and DEC-283 note iii). (iv) The Runner skipped a QA window
  without a recorded decision. The route or Runner should check section 8
  cadence before dispatching the next task.
- **Artifact/spec sync:** `09_AI_WORKFLOW.md` section 2.3 and
  `agents/tl-lite.md` (with generated copies) get clause 7. `tasks/task-9.md`
  gets the amendment. The triage log gets TRG-24-008. The ESC register,
  archive, brief and progress are refreshed. No PRD, TRD, data model,
  glossary, template, validation-contract or QA-plan change.
- **Refs:** ESC-24-003; TRG-24-006 (and correction), TRG-24-007, TRG-24-008;
  DEC-190, DEC-283; `run-ledger.md` 02:17, 02:20, 02:24;
  `workflow-sequence-lint.sh:488-511`, `:535-538`, `:705-760`, `:817-819`, `:1343-1377`;
  `qa/qa-plan.md` section 8; `09_AI_WORKFLOW.md` sections 2.3, 9.

### DEC-285: `activity.list` failures use the accepted retained-read mapper; stats and AC-list Human failures adopt the renderer. Closes ESC-24-004

- **Date:** 2026-10-07. Tech Lead ESC resolution on the Human's instruction
  "Resolve ESC-24-004". This decision clarifies accepted text. It adds no
  failure code, detail shape, message rule or exit class. Precedents:
  DEC-280 item 3 (an exit change that an accepted Phase 21 section
  authorizes), DEC-194 (an accepted allocation that was never delivered),
  DEC-175 (a latent defect first exposed by a task's new requirement belongs
  to that task) and DEC-283 clause 2 (STOP-resume identity).
- **Context:**
  - TRG-24-010 (BLOCK) and the 20:50 STOP. The Task 10 report at `9c3e003`
    records built-binary `activity list --since invalid` and
    `--on invalid`. Both print raw prose in both modes and exit 7. From a
    directory with no project the command prints nothing and exits 7.
    Neither mode emits a `fid-cli.v1` failure. Report-lint returns
    `OUTCOME BLOCKED` (exit 11).
  - Source at `9c3e003`: `cmd/fid/activity.go:33-35` returns the
    workspace-resolution error wrapped and untyped. Lines 45-49 call the raw
    `renderIssueError` and return the untyped error. JSON mode is read only
    after success. `internal/activity` `validateListOptions` returns untyped
    `--on` and `--since` causes before any read.
  - **What is accepted.**
    (a) `activity.list` is a Phase 21 section 4.1 `fid-cli.v1` leaf. Phase 20
    `contracts.md` section 6 kept it as unenveloped legacy JSON only until
    Phase 22. Phase 21 section 1 requires one failure envelope on stderr with
    the section 7 exit class. Phase 20 section 2.4 says that
    `internal-error`'s "presence in a test is a defect, not a contract". So
    the raw exit-7 output is a latent Phase 21 nonconformance, not an
    accepted JSON code, details or exit. On these paths there is nothing for
    TRG-24-009's "preserve accepted JSON code/details/exit" to preserve. The
    raw caller dates from `977974c` (phase-16/task-1). `ab09b47`
    (phase-21/task-9) enveloped only the success path.
    (b) The retained standalone reads in the same section 4.1 table already
    share one failure mapper: `issueReadFailure`, reached through
    `renderIssueReadError` (`cmd/fid/issue_get_list.go`). `issue.get`,
    `issue.list`, `issue.ac.list`, `stats.tokens` and `stats.cost` use it.
    `stats.*` resolve the workspace with the same `issueMutationWorkspaceRoot`
    that `activity.list` uses. The mapper's branches are: `project-not-found`,
    `issue-not-found`, `result-limit-exceeded`, `journal-replay-stopped`, then
    the R20 generic fallback (DEC-281).
    (c) `--since` and `--on` have no declared flag domain, and none was
    required. DEC-280 item 5 limits completeness to non-string flags (R18)
    and string domains that a Phase 21-24 row states. The `activity.list`
    row states no `--since` domain and has no `--on` (CF-047(b)). Their
    checks are therefore domain-package checks. Task 9a left those
    unchanged, and they reach the family mapper. DEC-280 item 5 cites
    `receipt list --outcome bogus` reaching `argument:"request"` as a check
    of the same kind. VAL-064 is not engaged.
    (d) Task 10's One renderer constraint says: "A family that formats Human
    failures without a JSON failure constructs the same failure it would
    emit in JSON mode." For `activity.list`, that failure is (a) produced
    through (b).
  - **Caller audit** (report Concerns). `stats.go:38,48,73,83` and
    `issue_ac.go:109,118` call `renderIssueReadError(cmd, true, …)`, so
    Human mode prints the JSON envelope on stderr. Contract section 9 says
    "Every Human-mode failure on stderr is exactly these lines". Their JSON
    mode is already correct. The forced arguments came from `e51afe3`
    (phase-21/task-9) and `16f4f8d` (phase-20/task-10), both before the
    Phase 24 task pack. Neither TRG-24-009 nor TRG-24-010 grants an edit
    there, so Task 10 cannot meet its R9 no-bypass criterion without this
    decision.
  - `stats.tokens` and `stats.cost` print the JSON success envelope in Human
    mode too; the retained `TestStatsTokensDefaultsToJSONByActorAttribution`
    asserts this. Success output is outside Task 10.
- **Decision:**
  1. **ESC-24-004 closed.** The BLOCK was correct: TRG-24-009 forbade
     inventing a mapping, and this decision invents none. The accepted
     mapping is Phase 21 sections 1, 4.1 and 7, applied through the
     retained-read family mapper. No new contract rule, code, detail shape,
     message rule or Human decision is needed; section 9 gains only a
     clarifying note. The exit change
     from 7 to 2, 3 or 6 restores Phase 21 conformance, as DEC-280 item 3
     did for unknown words.
  2. **`activity.list` mapping.** `activity.go` reads JSON mode first, as
     `issue get` does. It passes both the unwrapped `issueMutationWorkspaceRoot`
     error (as `stats.go` does) and every `activity.List` error to
     `renderIssueReadError(cmd, jsonMode, "activity.list", err)`. It adds no
     activity-specific branch. Expected failures, all with
     `retryable:false` and `operation:"activity.list"`:

     | Cause | Code / category / exit | `message` | `subject` | `details` |
     |---|---|---|---|---|
     | `--since` refused by `internal/activity` (e.g. `invalid`) | `invalid-argument` / invalid-input / 2 | `activity list: --since must be YYYY-MM-DD or RFC3339` (the DEC-281 cause) | null | `{argument:"request", value:"", allowed:["accepted issue read input"]}` |
     | `--on` refused by `internal/activity` (e.g. `invalid`) | same | `activity list: --on must be <type>/<id>` | null | same |
     | No `.fid/` in an existing directory (cwd or `--project`) | `project-not-found` / not-found / 3 | `fid project was not found` | `{type:"project", id:""}` | `{path:""}` |
     | Journal replay stopped | `journal-replay-stopped` / integrity / 6 | `journal replay stopped` | null | the existing `journalReplayStoppedFailure` details |
     | Any other cause | the mapper's final branch, unchanged (DEC-281) | | | |

     Human mode renders through `writeHumanFailure` with leaf
     `fid activity list`. The leaf has no `--input`, so there is no `next:`
     line. `--event` keeps its task 9a grammar failure. stdout and success
     output are unchanged. The no-project failure is byte-identical to
     `stats tokens` except for `operation`.
  3. **Stats and AC-list callers.** Task 10 gains `cmd/fid/stats.go` and
     `cmd/fid/issue_ac.go`, for one change only: the six forced `true`
     arguments become the live JSON mode. JSON failures stay byte-identical.
     Human failures become the section 9 lines. Exits are unchanged. Human
     success output (stats JSON) is out of scope (CF-047(c)).
  4. **No other product change.** Do not edit `internal/activity`,
     `issueReadFailure`, `genericFailure`, the flag-domain declarations or
     their texts. Declaring `--since` and `--on` domains belongs with the
     `activity.list` row convergence (CF-047(b), Phase 26).
  5. **Retained convergence.** If these four callers displace a retained
     assertion, converge it under DEC-162/165/166 with introducing-commit
     provenance, in the `cmd/fid` or `cmd/fid/integration` test file that
     holds it. The ordinary race gate and the tagged integration gate
     already run both directories (DEC-167/173). Keep the assertion: never
     delete or weaken it.
  6. **Evidence (VAL-036, VAL-067).** The built-binary pairing corpus gains
     a row for each clause 2 cause the CLI can reach: `--since`, `--on`, and
     no project. It also gains a no-project row each for `stats tokens`,
     `stats cost` and `issue ac list`. Each row is run in both modes. It
     checks the Human `error:` and `reason:` lines against JSON, the exit,
     empty stdout, and that nothing in the filesystem changes. The `--since`
     row may drive the R20 row for the `issueReadFailure` fallback. The
     report's R9 table lists all three callers.
  7. **Recovery identity** (STOP resume; DEC-283 clause 2). The chain trigger
     moves to TRG-24-010's trigger at attempt 2: `task-10`,
     `task-10:public-interface:activity-list-failure-envelope-mapping`,
     attempt `2`, `in-scope-implementation`, `task-scope`, `0/2`. Attempts
     stay monotonic: 1 (the TRG-24-009 continuation, which returned BLOCKED),
     then 2. The TRG-24-009 trigger stays in history.
  8. **Ledger sequence** (Runner, transactional append). Each step must
     return `APPEND_CONFIRMED state=REPAIR_PENDING` with observed
     `exit_code: 0`, and the step must be appended before dispatch.
     (a) **RUN-RESUME** (schema v2 STOP resume). Set
     `Prior run/session or STOP: 2026-10-07 20:50 STOP / TRG-24-010 / ESC-24-004`,
     `Governing TRG: TRG-24-010`, `Governing ESC: ESC-24-004`, the trigger
     and attempt from clause 7, and `Resolution authority: DEC-285`. Set
     `Continuing task: task-10` and
     `Next recovery: task-10 | in-scope-implementation | task-scope | task-10:public-interface:activity-list-failure-envelope-mapping | 2`.
     `Checkpoint` names the post-DEC-285 phase HEAD and the preserved
     `d574cf4`/`9c3e003` checkpoints.
     (b) **TASK-RECOVERY task-10**, with the same tuple and `0/2`.
     `Source: TRG-24-010 / ESC-24-004 / DEC-285; project-docs/docs/phases/phase-24/tasks/task-10-done.md`.
     `Repair authority: TRG-24-010 / DEC-285`.
     `Triage disposition: N/A`. Control fields and base are `N/A`.
     `Originating report: N/A`. `Exact verification: N/A`. Set
     `Allowed files: ordinary task-10.md scope plus cmd/fid/activity.go (TRG-24-009) and cmd/fid/stats.go,cmd/fid/issue_ac.go (DEC-285 clause 3)`.
     `Forbidden changes` is TRG-24-009's list, plus: no edit to
     `internal/activity`, `issueReadFailure`, `genericFailure` or flag-domain
     declarations, and no stats success-output change. Dispatch the
     same-task Developer continuation only after (b) is appended. Its later
     REPAIR uses the same trigger and attempt.
  9. **Verification.** The Tech Lead dry-ran the sequence on a scratch copy
     of the phase files with the unchanged verifier. The current ledger
     validates `VALID state=STOPPED`, and (a) and (b) each returned
     `APPEND_CONFIRMED state=REPAIR_PENDING` (exit 0). The ledger then
     validates `VALID state=REPAIR_PENDING`. The negative guards below each
     returned exit 1 with the scratch ledger byte-unchanged:

     | Candidate | Diagnostic |
     |---|---|
     | RUN-RESUME at attempt 1 | `RUN-RESUME Attempt count must increment its governing TRIAGE` |
     | RUN-RESUME with the TRG-24-009 trigger | `RUN-RESUME Trigger signature does not match its governing TRIAGE` |
     | TASK-RECOVERY whose `Source` omits DEC-285 | `TASK-RECOVERY does not retain its RUN-RESUME decision authority`; `… typed RUN-RESUME resolution authority` |
     | TASK-RECOVERY with the TRG-24-009 trigger | `TASK-RECOVERY does not preserve its RUN-RESUME checkpoint`; `… does not match the exact schema v2 RUN-RESUME recovery tuple` |

     If (a) or (b) is rejected as written, that is new information. Return to
     the Tech Lead, and never substitute field values.
  10. **Unchanged.** TRG-24-009, TRG-24-010, the STOP, the report, every
      failed probe and the recovery history stay as recorded. Task 10's
      seven commands, VAL identity, thresholds, budgets and the QA cadence
      are unchanged. Window D2 after Task 11 and the Human walkthrough are
      not waived.
- **Artifact/spec sync:** `tasks/task-10.md` gains a DEC-285 amendment.
  `public-interface.md` section 9 gains one clarifying note. The VAL-036
  method and the change log record the clause 6 rows. Inventory R9 notes
  the three callers. CF-047 gains (b) and (c) wording in
  `decisions_active.md`. The ESC register, archive, brief and progress are
  refreshed. There is no PRD, TRD, data model, glossary, template, QA-plan,
  tool or prompt change.
- **Workflow note for wrap-up:** TRG-24-009 bounded the file scope correctly.
  But no one checked whether the target leaf had a JSON failure at all. A
  caller-scope disposition should first confirm that the "existing mapping"
  it relies on exists.
- **Refs:** ESC-24-004; TRG-24-009, TRG-24-010; `run-ledger.md` 20:15,
  20:17, 20:49, 20:50; `tasks/task-10-done.md` at `9c3e003`; Phase 20
  `contracts.md` sections 2.4, 6; Phase 21 `public-interface.md` sections
  1, 2, 4.1, 7; DEC-175, DEC-194, DEC-280, DEC-281, DEC-283.

### DEC-286: A root flag a stable leaf does not accept is a Phase 21 section 2 grammar failure; `export import --project` adopts it. Closes ESC-24-005

- **Date:** 2026-10-09. Tech Lead ESC resolution on the Human's instruction
  "Resolve ESC-24-005". This decision clarifies accepted text. It adds no
  failure code, detail shape or exit class, and it changes no success
  output or existing JSON failure. One message is chosen (clause 3). The
  Human may override that wording by DEC before Task 10 closes.
  Precedents: DEC-285 (a leaf path with no `fid-cli.v1` failure has no
  accepted exit to preserve), DEC-280 item 4 (task 9a's class table) and
  DEC-175 (a latent defect first exposed by a task's new requirement
  belongs to that task).
- **Context:**
  - TRG-24-011 (BLOCK) and the 2026-10-08 05:24 STOP. The Task 10 report
    at `e0495b5` records built-binary
    `fid [--json] --project <dir> export import missing.zip --repository <dir> --idempotency-key probe`.
    In both modes it prints
    `error: --project is not available for export import until its fid-cli.v1 transport is migrated`,
    exits 7, leaves stdout empty and changes no file. Report-lint returns
    `OUTCOME BLOCKED` (exit 23). The Tech Lead reproduced both modes on a
    scratch build of `e0495b5`.
  - **Source.** `rejectUnavailableRootFlags` (`cmd/fid/root.go:106-128`)
    is the root `PersistentPreRunE`. For a leaf whose `rootFlagCapabilities`
    entry disables a root flag that was set, it writes raw `error: …` and
    returns an untyped error, which `cli.ExitCode` maps to 7. Cobra v1.10.2
    runs it after flag parsing and the positional check, and before the
    handler (`command.go:968,986`), so before any workspace read. Two kinds
    of leaf reach it. (i) `sql` and `journal inspect` have no registry
    entry, so both flags are refused. (ii) `export import` is registered
    `{JSON: true}` by `5ec69e6` (phase-22/task-15). Every other stable
    leaf is registered for both flags. So `export import --project` is the
    only stable-leaf instance. The guard dates from `cd53c43`
    (phase-20/task-2) as the temporary "not yet migrated" rule.
  - **What is accepted.**
    (a) Phase 21 section 2: "Missing required arguments/flags, forbidden or
    unknown flags, extra arguments, mutually exclusive flag violations, and
    malformed flag values return `invalid-argument` with the offending
    spelling/value and the row's exact accepted alternatives in `allowed`;
    they perform no read or write." Section 7 makes `invalid-argument`
    invalid-input, not retryable, exit 2, subject null, with details
    `{argument, value, allowed}`. The refusal itself is delivered Phase 22
    behaviour. The Phase 22 `export.import` row names its target as
    `--repository` and lists no `--project`, and the registry refuses it.
    So `--project` is a forbidden flag on that leaf, and accepted text
    fixes its code, category, retryable, exit, subject and detail shape.
    (b) `export.import` is a `fid-cli.v1` leaf (Phase 22 section 1 and its
    row), so the raw exit 7 is a Phase 21 section 1/7 nonconformance, as
    `activity.list` was under DEC-285. No accepted exit or message exists
    to preserve. "Until its transport is migrated" describes Phase 20
    scaffolding; the transport was migrated in Phase 22.
    (c) Task 9a's envelope rule gives a grammar failure the leaf's
    operation ID (`export.import`). Its unknown-flag class sets `argument`
    and `value` to the spelling, and `allowed` to the leaf's local long
    flags, sorted, without `--help`, `--json` and `--project`. Live
    `fid export import a.zip --bogus` returns
    `allowed:["--idempotency-key","--repository"]`, the Phase 22 row's
    two flags.
    (d) Contract section 9 derives the Human lines from the JSON failure.
    `export import` has no `--input`, so there is no `next:` line.
  - **What was not determined.** Section 2's "spelling/value" does not say
    which belongs in `value` for a recognized flag that carries a value.
    Task 9a's class table has no row for a recognized but forbidden flag.
    No accepted message exists. TRG-24-011 was right not to choose them.
- **Decision:**
  1. **ESC-24-005 closed.** The BLOCK was correct. The code, category,
     retryable, exit, subject, detail shape and operation come from
     accepted text (Context a-d). Clause 2 fixes the remaining values, and
     section 9 gains one clarifying note. No code, detail shape, exit
     class, renderer rule or success output changes.
  2. **Mapping: a root flag that a stable leaf's row does not accept.**
     This applies to every leaf except `sql` and `journal inspect`. Today
     the only instance is `export import` with `--project`. It is a
     grammar failure, a fifth class beside task 9a's table:

     | Field | Value |
     |---|---|
     | `code` / `category` / `retryable` / exit | `invalid-argument` / `invalid-input` / `false` / 2 |
     | `subject` | `null` |
     | `operation` | the leaf's operation ID (task 9a): `export.import` |
     | `details.argument` | the flag spelling: `--project` |
     | `details.value` | the flag spelling: `--project`, as for task 9a's unknown flag. The supplied path is not echoed: the flag is what is refused, and Phase 21 section 7 bars unrelated absolute paths from details. |
     | `details.allowed` | computed exactly as task 9a's unknown-flag class: the leaf's local long flags, sorted bytewise, without `--help`, `--json` and `--project`. For `export import`: `["--idempotency-key","--repository"]` |
     | `message`, before the DEC-282 bounded rule | `--<flag> is not available for <command path without "fid ">`: `--project is not available for export import` |

     Exact JSON failure (stderr; stdout empty; exit 2):

     ```json
     {"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["--idempotency-key","--repository"],"argument":"--project","value":"--project"},"message":"--project is not available for export import","retryable":false,"subject":null},"ok":false,"operation":"export.import"}
     ```

     Exact Human failure (stderr; stdout empty; exit 2), from the shared
     renderer with leaf `fid export import`:

     ```text
     error: --project is not available for export import
     reason: invalid-argument; argument=--project, allowed=--idempotency-key|--repository
     ```

     JSON mode is the live `--json` value. The check keeps its position:
     after flag parsing and the positional check, before the handler. No
     workspace is read, so the failure is identical from a non-project
     directory and with a nonexistent `--project` path. An unknown flag or
     a wrong positional count still fails first, with its task 9a class. A
     missing `--repository` or `--idempotency-key` is reported only when
     `--project` is absent. Nothing is written.
  3. **Message.** The Phase 20 clause "until its fid-cli.v1 transport is
     migrated" is dropped because it is false for a `fid-cli.v1` leaf. The
     kept stem is the detecting check's own text, as in task 9a's "missing
     required flag" class. Messages are never control flow (Phase 21
     section 7), so later rewording needs only a DEC.
  4. **Unchanged.** `sql` and `journal inspect` keep their byte-identical
     raw refusal and exit. They are outside `fid-cli.v1` (Phase 21 section
     5), and Task 10 excludes them. Their retained assertions in
     `root_flags_test.go` and `contract_surface_test.go`
     (`assertLegacyLeafRejectsRootFlag`) stay unchanged. The
     `rootFlagCapabilities` registry is unchanged, so `export import` still
     refuses `--project`; this decision does not consider accepting it as a
     target. `writeGrammarFailure`, `writeHumanFailure`, `genericFailure`,
     `internal/cli` and the flag-domain declarations are not edited.
  5. **Task 10 scope.** Task 10's One renderer requirement is the first to
     expose this latent Phase 20 guard (DEC-175). Product changes:
     - `cmd/fid/root.go` `rejectUnavailableRootFlags`: for any leaf other
       than `sql` and `journal inspect`, return
       `writeGrammarFailure(command, leafOperation(command), cli.CodeInvalidArgument, <clause 2 message>, <clause 2 details>)`.
       The raw branch for the two diagnostics stays byte-identical.
     - `cmd/fid/canonical_request.go`: pure extraction with no behaviour
       change. The unknown-flag `allowed` computation in
       `canonicalFlagError`, and the `sql`/`journal inspect` test in
       `installLeafGrammar`, each move to one unexported helper. The
       original caller and the root guard both call it. The task 9a corpus
       in the ordinary race gate proves nothing else changed.

     Tests. Each is new or Task 10's own, and no retained assertion
     changes:
     - `cmd/fid/root_flags_test.go`: add in-process clause 2 rows in both
       modes. Include a non-project working directory, a nonexistent
       `--project` path, and a check that the two diagnostics still print
       their raw text. Existing tests in this file are unchanged.
     - `cmd/fid/human_failure_caller_audit_test.go` (added by `6607143`,
       phase-24/task-10): the `root.go` check may permit an `error: `
       literal only in the branch that `sql` and `journal inspect` reach.
       It must still fail on any other raw `error: ` text in `root.go`. A
       negative control proves it, for example a parsed source fixture
       whose stable-leaf branch writes raw text.
     - `cmd/fid/integration/p24_retained_read_failure_test.go`
       `TestP24UnavailableImportRootFlagFailurePair` (added by `6607143`):
       assert the exact clause 2 JSON bytes and the exact two Human lines,
       exit 2 in both modes, empty stdout and an unchanged filesystem. Keep
       it separate from the valid `export.import --repository` fixture.

     `gofmt -l cmd internal`, `golangci-lint run ./...`, the ordinary
     race/coverage gate and the tagged integration gate already execute
     every file above (DEC-167/173). The seven commands are unchanged.
  6. **Evidence (VAL-036).** The built-binary pairing corpus gains the
     clause 2 row in both modes. The report's R9 table lists the root guard
     as a renderer caller beside the DEC-285 callers.
  7. **Recovery identity** (STOP resume; DEC-283 clause 2, DEC-285 clause
     7). The chain trigger moves to TRG-24-011's trigger at attempt 2:
     `task-10`, `task-10:public-interface:unavailable-root-flag-failure-mapping`,
     attempt `2`, `in-scope-implementation`, `task-scope`, `0/2`. Attempt 1
     of this trigger is the DEC-285 continuation, which returned BLOCKED
     at `e0495b5` after material progress (the six DEC-285 pairs went
     green). The TRG-24-009 and TRG-24-010 triggers and the DEC-285
     recovery stay in history.
  8. **Ledger sequence** (Runner, transactional append). Each step must
     return `APPEND_CONFIRMED state=REPAIR_PENDING` with observed
     `exit_code: 0`, and each must be appended before dispatch.
     (a) **RUN-RESUME** (schema v2 STOP resume). Set
     `Prior run/session or STOP: 2026-10-08 05:24 STOP / TRG-24-011 / ESC-24-005`,
     `Governing TRG: TRG-24-011`, `Governing ESC: ESC-24-005`, the trigger
     and attempt from clause 7, and `Resolution authority: DEC-286`. Set
     `Continuing task: task-10` and
     `Next recovery: task-10 | in-scope-implementation | task-scope | task-10:public-interface:unavailable-root-flag-failure-mapping | 2`.
     `Checkpoint` names the post-DEC-286 phase HEAD and the preserved
     `6607143` (implementation) and `e0495b5` (report) checkpoints.
     (b) **TASK-RECOVERY task-10**, with the same tuple and `0/2`.
     `Source: TRG-24-011 / ESC-24-005 / DEC-286; project-docs/docs/phases/phase-24/tasks/task-10-done.md`.
     `Repair authority: TRG-24-011 / DEC-286`.
     `Triage disposition: N/A`. Control fields and base are `N/A`.
     `Originating report: N/A`. `Exact verification: N/A`. Set
     `Allowed files: ordinary task-10.md scope plus cmd/fid/activity.go (TRG-24-009), cmd/fid/stats.go,cmd/fid/issue_ac.go (DEC-285 clause 3), and cmd/fid/root.go,cmd/fid/canonical_request.go,cmd/fid/root_flags_test.go,cmd/fid/human_failure_caller_audit_test.go,cmd/fid/integration/p24_retained_read_failure_test.go (DEC-286 clause 5)`.
     `Forbidden changes` is the 04:27 TASK-RECOVERY list, plus: no change
     to the `rootFlagCapabilities` registry, to the `sql`/`journal inspect`
     refusal bytes or their retained assertions, or to `writeGrammarFailure`
     or `writeHumanFailure`, and no `canonical_request.go` change beyond the
     clause 5 extraction. Dispatch the same-task Developer continuation only
     after (b) is appended. Its later REPAIR uses the same trigger and
     attempt.
  9. **Verification.** The Tech Lead dry-ran the sequence on a scratch copy
     of the phase files (with this decision) and the unchanged verifier.
     The current ledger validates `VALID state=STOPPED`. (a) and (b) each
     returned `APPEND_CONFIRMED state=REPAIR_PENDING` (exit 0), and the
     ledger then validates `VALID state=REPAIR_PENDING`. Each negative
     guard below returned exit 1 with the scratch ledger byte-unchanged:

     | Candidate | Diagnostic |
     |---|---|
     | RUN-RESUME at attempt 1 | `RUN-RESUME Attempt count must increment its governing TRIAGE` |
     | RUN-RESUME with the TRG-24-010 trigger at attempt 3 | `RUN-RESUME Trigger signature does not match its governing TRIAGE`; `… Attempt count must increment …` |
     | TASK-RECOVERY whose `Source` omits DEC-286 | `TASK-RECOVERY does not retain its RUN-RESUME decision authority`; `… typed RUN-RESUME resolution authority` |
     | TASK-RECOVERY with the TRG-24-010 trigger | `TASK-RECOVERY does not preserve its RUN-RESUME checkpoint`; `… does not match the exact schema v2 RUN-RESUME recovery tuple` |

     The verifier also accepted a RUN-RESUME citing
     `Resolution authority: DEC-285`. It checks only that the cited DEC
     exists, not that it resolves the governing ESC. This is observed
     control friction for CF-036, not a blocker: (a) must cite DEC-286 as
     written. If (a) or (b) is rejected as written, that is new
     information. Return to the Tech Lead, and never substitute field
     values.
  10. **Unchanged.** TRG-24-009, TRG-24-010, TRG-24-011, both STOPs, every
      report, failed probe and recovery entry stay as recorded, and so do
      DEC-285 and its six pairs. The seventeen missing `RegisteredCodes`
      rows and the activation, projection and budget fixtures remain
      ordinary Task 10 work under its existing rules: no invented exclusion
      or stdout exception. Task 10's seven commands, VAL identity,
      thresholds, budgets and the QA cadence are unchanged. Window D2 after
      Task 11 and the Human walkthrough are not waived.
- **Artifact/spec sync:** `tasks/task-10.md` gains a DEC-286 amendment.
  `public-interface.md` section 9 gains one clarifying note. The VAL-036
  method and the change log record the clause 6 row. Inventory R9 notes the
  root guard. CF-047 records the `sql`/`journal inspect` refusal wording in
  `decisions_active.md`. The ESC register, archive, brief and progress are
  refreshed. There is no PRD, TRD, data model, glossary, template, QA-plan,
  tool or prompt change.
- **Workflow note for wrap-up:** TRG-24-009 and DEC-285 audited family
  writers and named callers, but not the root `PersistentPreRunE`, which
  emits failures before any family writer runs. A renderer-adoption audit
  should walk pre-handler hooks as well as handlers.
- **Refs:** ESC-24-005; TRG-24-011; `run-ledger.md` 2026-10-08 04:26,
  04:27, 05:24; `tasks/task-10-done.md` at `e0495b5`; Phase 21
  `public-interface.md` sections 1, 2, 5, 7; Phase 22 `public-interface.md`
  section 1 and the `export.import` row; task 9a class table; DEC-175,
  DEC-280, DEC-282, DEC-283, DEC-285.

### DEC-287: A plain copy keeps `run.inspect`'s retained refusal; VAL-041 proves `different worktree` on `resume`. Closes ESC-24-006

- **Date:** 2026-10-09. Tech Lead ESC resolution on the Human's instruction
  to resolve ESC-24-006. This decision clarifies accepted text. It changes
  no behaviour, output, code, exit or JSON in either mode. It corrects the
  VAL-041 method and the Task 11 binding constraint and acceptance
  criterion, which required a successful plain-copy view from every
  recorded-versus-current view, so that they agree with contract section 2.
  The Human may override by DEC before Task 11 closes (rejected option B,
  clause 1).
- **Context:**
  - TRG-24-012 (BLOCK) and the 2026-10-09 08:55 STOP. Task 11 checkpoints:
    `ab05e49` (implementation) and `9ce8896` (report, report-lint
    `OUTCOME BLOCKED`, exit 0). The tagged generator row
    `clean-different-run.inspect` expects exit 0 and `different worktree`.
    Both focused `TestP24HumanFixtures` runs fail (exit 1) at that row.
  - **Tech Lead reproduction** with the retained binary
    `.cache/task-11-blocker-evidence/fixtures-1289185125/fid`, project
    `clean-same`, Run `RUN-d5c15a`. With `--repository` set to the same
    root, both modes exit 0 and Human prints
    `Worktree comparison: same worktree`. With `--repository` set to the
    plain copy `clean-different`, both modes exit 4, stdout is empty, and
    stderr carries the clause 2 bytes.
  - **Source.** `InspectRun` reconciles before composing its result
    (`internal/workflow/run_inspect.go:45`). `reconcileRunRepository`
    (`import_binding.go:11`) takes the imported-destination path only with
    a matching import binding. Otherwise `reconcileRepository` refuses a
    changed `WorktreeID` with `repository identity changed`
    (`reconciliation.go:53`; `ac5e6f9`, phase-22/task-4). So a successful
    `run.inspect` always compares the recorded worktree or a bound
    destination, and `SameWorktree` accepts both. Its view can only say
    `same worktree`. Under accepted contracts, `different worktree` is
    unreachable in `run.inspect`.
  - **What is accepted.**
    (a) Contract section 2 names `run.inspect` among the leaves that keep
    their Phase 21-23 behaviour.
    (b) Phase 22 `contracts.md` section 3: `run.inspect` is read-only and
    compares the observation with Run start or the latest checkpoint.
    "Changed repository identity" is among the states that "do not
    qualify". The delivered refusal, `recovery-state-conflict` (category
    `conflict`, exit 4) with `repository identity changed`, dates from
    `ac5e6f9` (phase-22/task-4). Phase 22 closed with it in place.
    (c) Contract section 10 and Task 11: JSON is unchanged. Section 10 says
    what a comparison line says when a view renders. It does not require a
    view whose producer refuses a plain copy to render one.
    (d) DEC-278 R-11 disposed "Same-worktree and plain-copy fixtures". The
    QA plan's R-11 proposal offered
    "`fid resume --repository <copy>` or `fid run inspect --repository <copy>`".
    `resume` renders a plain copy: the Task 11 report records Human
    `different worktree`, exit 0. So the plain-copy state is provable on
    `resume`.
    (e) CF-046(b) and CF-026a already leave `run.inspect`'s repository
    behaviour, and inspection beyond accepted recovery, to Phase 25 and
    future design.
  - **Cause.** The VAL-041 method and Task 11 (DEC-277/278) generalized
    R-11's "either leaf" proposal to every recorded-versus-current view
    without checking that `run.inspect` refuses a plain copy. This is a
    Tech Lead planning defect. The Developer and TL-Lite were right not
    to choose.
- **Decision:**
  1. **ESC-24-006 closed; option A.** Keep the accepted behaviour. No
     Human public-interface choice is needed, because nothing changes in
     either mode. Rejected option B: make `run.inspect` render a plain copy,
     for example by skipping reconciliation or by returning a
     non-qualifying comparison. B turns an accepted exit-4 failure into a
     success, so it changes JSON (section 10) and reverses section 2. It
     also needs a result shape for a non-qualifying reconciliation that
     Phase 22 does not define. That is CF-026a and CF-046 work for Phase 25
     parallel worktrees, and it needs Human approval.
  2. **Plain-copy `run.inspect` outcome (unchanged Phase 22 refusal).** It
     applies when the observed worktree (`--repository`, else the current
     worktree) differs from the Run's recorded base worktree and no import
     binding matches, clean or dirty. Both modes exit 4 with empty stdout and append
     no Event. Exact JSON (stderr):

     ```json
     {"contract_version":"fid-cli.v1","error":{"category":"conflict","code":"recovery-state-conflict","details":{"observation_digest":null,"paths":[],"reasons":["repository identity changed"],"run_id":"<run-id>"},"message":"recovery state conflict","retryable":false,"subject":{"id":"<run-id>","type":"run"}},"ok":false,"operation":"run.inspect"}
     ```

     Exact Human (stderr), the live Task 10 section 9 rendering of that
     failure:

     ```text
     error: recovery state conflict
     reason: recovery-state-conflict; observation_digest=null, paths=, reasons=repository identity changed, run_id=<run-id>
     ```

     There is no `blocked by:` or `next:` line, and no repository view line
     (`Recorded repository:`, `Current repository:` or
     `Worktree comparison:`). The identity check runs before the base,
     ancestry, unmerged and operation checks, so a dirty plain copy gives
     the same bytes with its own Run ID.
  3. **VAL-041 amended.**
     - Method: every Human view that shows a repository (inventory R11),
       with clean and dirty fixtures. `resume` gets a same-worktree fixture
       and a plain-copy fixture (`--repository <copy>`). `run.inspect` gets a
       same-root fixture and a plain-copy refusal row.
     - Pass: shows `HEAD <12 hex>` with `clean` or `dirty`. `resume` says
       `same worktree` for the same root and `different worktree` for the
       plain copy. `run.inspect` says `same worktree` for the same root. For
       the plain copy it returns the clause 2 refusal in both modes. No
       forbidden text appears.
     - Fail: any forbidden text, a wrong state, or a plain-copy
       `run.inspect` that succeeds or differs from clause 2.
     - Notes: JSON is unchanged. A successful `run.inspect` only ever says
       `same worktree` (DEC-287); different-worktree inspection is
       CF-046(c).
  4. **Task 11 amended; preserved assertions.** No row is deleted or
     renamed, and both states (clean, dirty) are kept. In
     `cmd/fid/integration/p24_human_fixtures_test.go`, only
     `<state>-different-run.inspect` changes its expectation, from exit 0
     and `different worktree` to the clause 2 refusal. The row has code
     `recovery-state-conflict` and expected exit 4. Its Human stderr equals
     the shared-renderer derivation of its JSON failure
     (`p24ExpectedFailureText`) and the clause 2 bytes with that fixture's
     Run ID. Its details carry `reasons:["repository identity changed"]`
     and its subject is that Run. Human stdout stays empty. These stay
     unchanged:
     - `<state>-different-resume`: `different worktree`, the short form and
       the forbidden-token scan.
     - `<state>-same-run.inspect` and `<state>-same-resume`:
       `same worktree`.
     - Every other row.
     - The manifest check's five denial names and 28 view names.
     - The imported-destination `same worktree` assertion in
       `p22_human_guidance_test.go`.
     `TestP24RepositoryHumanViews` may keep its synthetic `run.inspect`
     different-worktree row as renderer coverage. It is not VAL-041
     plain-copy evidence.
  5. **Scope.** Ordinary `task-11.md` scope, with no new path. This
     decision changes only `p24_human_fixtures_test.go`, Task 11's new
     file, which the tagged gate executes (DEC-167/173). No product source
     changes for this row. Forbidden:
     - edits to `InspectRun`, `reconcileRunRepository`,
       `reconcileRepository`, `importedDestinationBinding` or
       `SameWorktree`;
     - edits to the `recovery-state-conflict` mapping or the section 9
       renderer;
     - `run_inspect_human.go` changes beyond Task 11's repository lines;
     - substituting an import for the copy, relabelling it, or removing the
       row.
     Task 11's six commands are unchanged. The 08:45 ledger prose says
     "seven"; Task 11 has six.
  6. **Recovery identity** (STOP resume; DEC-283 clause 2, DEC-286
     clause 7). `task-11`,
     `task-11:public-interface:plain-copy-run-inspect-repository-view`,
     attempt `2`, `in-scope-implementation`, `task-scope`, `0/2`. Attempt 1
     is the replacement continuation dispatched after the 08:46
     `INTERRUPTED` REPAIR. It returned BLOCKED at `9ce8896` after material
     progress (renderers, Stop corpus, five denials).
     `task-11:session-continuation:developer-session-interrupted` attempt 1
     and its 08:45 and 08:46 entries stay in history unchanged.
  7. **Ledger sequence** (Runner, transactional append). Each step must
     return `APPEND_CONFIRMED state=REPAIR_PENDING` with observed
     `exit_code: 0`, and each must be appended before dispatch.
     (a) **RUN-RESUME** (schema v2 STOP resume). Set
     `Prior run/session or STOP: 2026-10-09 08:55 STOP / TRG-24-012 / ESC-24-006`,
     `Governing TRG: TRG-24-012`, `Governing ESC: ESC-24-006`, the trigger
     and attempt from clause 6, and `Resolution authority: DEC-287`. Set
     `Continuing task: task-11` and
     `Next recovery: task-11 | in-scope-implementation | task-scope | task-11:public-interface:plain-copy-run-inspect-repository-view | 2`.
     `Checkpoint` names the post-DEC-287 phase HEAD and the preserved
     `ab05e49` (implementation) and `9ce8896` (report) checkpoints.
     (b) **TASK-RECOVERY task-11**, with the same tuple and `0/2`.
     `Source: TRG-24-012 / ESC-24-006 / DEC-287; project-docs/docs/phases/phase-24/tasks/task-11-done.md`.
     `Repair authority: TRG-24-012 / DEC-287`.
     `Triage disposition: N/A`. Control fields and base are `N/A`.
     `Originating report: N/A`. `Exact verification: N/A`.
     `Allowed files: ordinary task-11.md scope`. `Forbidden changes` is
     the 08:45 TASK-RECOVERY list, plus clause 5's list. The one
     validation-assertion change permitted is the clause 4 row. Dispatch
     the same-task Developer continuation only after (b) is appended. Its
     later REPAIR uses the same trigger and attempt.
  8. **Verification.** The Tech Lead dry-ran the sequence on a scratch copy
     of the phase files (with this decision) and the unchanged verifier.
     The live ledger validates `VALID state=STOPPED`. (a) and (b) each
     returned `APPEND_CONFIRMED state=REPAIR_PENDING` (exit 0), and the
     ledger then validates `VALID state=REPAIR_PENDING`. Each negative
     guard below returned exit 1 with the scratch ledger byte-unchanged:

     | Candidate | Diagnostic |
     |---|---|
     | RUN-RESUME at attempt 1 | `RUN-RESUME Attempt count must increment its governing TRIAGE` |
     | RUN-RESUME with the session-continuation trigger | `RUN-RESUME Trigger signature does not match its governing TRIAGE` |
     | RUN-RESUME citing an unrecorded DEC-288 | `RUN-RESUME cites an unrecorded or malformed phase decision` |
     | TASK-RECOVERY whose `Source` omits DEC-287 | `TASK-RECOVERY does not retain its RUN-RESUME decision authority`; `… typed RUN-RESUME resolution authority` |
     | TASK-RECOVERY with the session-continuation trigger, attempt 1 | `TASK-RECOVERY does not preserve its RUN-RESUME checkpoint`; `… exact schema v2 RUN-RESUME recovery tuple` |
     | TASK-RECOVERY with no RUN-RESUME | `continuing recovery must preserve task and Trigger signature` |

     As under DEC-286, the verifier also accepted a RUN-RESUME citing
     `Resolution authority: DEC-286`, because it checks only that the DEC
     exists (CF-036). (a) must cite DEC-287 as written. If (a) or (b) is
     rejected as written, that is new information. Return to the Tech
     Lead, and never substitute field values.
  9. **Unchanged.** TRG-24-012, the STOP, the Task 11 report, both focused
     failures, the built-binary probes and the session-continuation
     recovery stay as recorded. The rest of Task 11 is ordinary work under
     its existing rules: the real `stop raise` corpus with qualifying
     resolution Events, the dirty matrix, the manifest, retained Human
     fixtures, the `06_DESIGN_PLAN_CLI.md` examples and the six commands in
     order. Report-lint `OUTCOME COMPLETE` is still required. Window D2
     (tasks 10-11, with the DEC-285/286 rows and this clause 2 pair), the
     fixture regeneration at Task 11's close commit and the Human
     walkthrough are not waived.
- **Artifact/spec sync:** `public-interface.md` section 10 gains one
  clarifying note. VAL-041 and the change log record clause 3.
  `tasks/task-11.md` gains a DEC-287 amendment. `decisions_active.md`
  CF-046 gains item (c). The ESC register, archive, brief and progress are
  refreshed. There is no PRD, TRD, data model, glossary, inventory,
  QA-plan, template, tool or prompt change.
- **Workflow note for wrap-up:** when a VAL row requires a state for every
  view of a class, check that each view's producer can reach that state
  under the retained contracts. R-11 offered either leaf; planning
  generalized it to both.
- **Refs:** ESC-24-006; TRG-24-012; `run-ledger.md` 2026-10-09 08:45,
  08:46, 08:55; `tasks/task-11-done.md` at `9ce8896`; contract sections 2
  and 10; Phase 22 `contracts.md` section 3 and `public-interface.md`
  `run.inspect` row; `qa/qa-plan.md` R-11; DEC-278 R-11; CF-026a; CF-046;
  DEC-283, DEC-286.

### DEC-288: Window D2 Human walkthrough recorded; review-10 closes B1-H through an evidence-only Task 11 recovery. Closes ESC-24-007

- **Date:** 2026-10-10. Tech Lead ESC resolution on the Human's instruction
  to resolve ESC-24-007. This decision records no product judgment of its
  own and changes no behaviour, output, code, JSON, assertion, gate or QA
  artifact. It binds the ledger route from the 2026-10-09 16:28 STOP to a
  fresh independent QA checkpoint review.
- **Context:**
  - Review-9 (`b9083f3b`) returned BLOCK with one finding, B1-H. Its
    deterministic D2 checks and the fixture regeneration at Task 11 close
    `93a1f792` passed, but the interim VAL-054/VAL-055 Human walkthrough
    (QA plan section 8, window D2) was missing. TL-Lite recorded
    TRG-24-013, then the identity corrections TRG-24-014 and TRG-24-015
    (governing), and raised ESC-24-007. The Runner appended TRIAGE
    TRG-24-015 and the STOP.
  - **Human evidence.** The Human (Mark Bell) recorded his own answers in
    `qa/human-walkthrough-after-task-11.md`, commit
    `a23b9f33` (`docs(qa): record phase-24 Window D2 Human walkthrough`).
    The Tech Lead created the blank template on the Human's request. At the
    Human's instruction it then committed `fcc60def`
    (`docs(qa): clarify phase-24 Window D2 Human walkthrough answers`). That
    commit fixes HUMAN-02a wording and HUMAN-02c's boolean spelling, and
    adds to HUMAN-03a..03d what each accepted view shows. The artifact's
    edit note attributes those edits; no judgment changed. The artifact
    names the manifest
    `.cache/phase-24-human/fixtures-2268063598/manifest.json`, binary
    `sha256:61fa7a4c…f55b7` and Task 11 close `93a1f792`. The Tech Lead
    confirmed that the retained binary still has that digest. All nine
    rows are answered: HUMAN-02a..02e each state a reason and a next action
    and are marked clear from the output alone; HUMAN-03a..03d are each
    `accept` with a reason. No row is rejected or unclear, so no acceptance
    finding needs an owner or authority disposition.
  - **Ledger constraint.** The unchanged sequence verifier lets a STOP
    resume only through `RUN-RESUME` with a continuing task, followed by
    that task's `TASK-RECOVERY`. It has no direct route from a STOP to a
    REVIEW. `REVIEW BLOCK` with `Owning task: none` opens no QA repair, so
    the Repair owners route does not apply. The governing trigger is
    task-prefixed (`task-11:…`, TRG-24-014/015). The only legal route is
    therefore a Task 11 recovery whose single obligation, the Human
    evidence, is already met.
- **Decision:**
  1. **ESC-24-007 closed.** The Human authority it waited for has been
     supplied. Closing B1-H stays with QA (clause 4): this decision
     neither validates the answers nor claims VAL-054/055 PASS.
  2. **Recovery identity** (STOP resume; DEC-283 clause 2). `task-11`,
     `task-11:qa-window-d2:human-walkthrough:val-054-055-pending`, attempt
     `2`, `qa-rework`, `task-scope`, `0/2`. Attempt 1 is the TRG-24-015
     BLOCK. No Developer is dispatched: the continuation's work was
     Human-authored and is already committed. Task 11's product and report
     state is unchanged since `93a1f792`.
  3. **Ledger sequence** (Runner, transactional append, each with observed
     `exit_code: 0`):
     (a) **RUN-RESUME** (schema v2 STOP resume):
     `Prior run/session or STOP: 2026-10-09 16:28 STOP / TRG-24-015 / ESC-24-007`,
     `Governing TRG: TRG-24-015`, `Governing ESC: ESC-24-007`, the clause 2
     trigger and attempt, `Resolution authority: DEC-288`,
     `Continuing task: task-11` and
     `Next recovery: task-11 | qa-rework | task-scope | task-11:qa-window-d2:human-walkthrough:val-054-055-pending | 2`.
     `Checkpoint` names the post-DEC-288 phase HEAD, walkthrough commits
     `a23b9f33` (answers) and `fcc60def` (attributed edits), review-9
     `b9083f3b` and Task 11 close `93a1f792`. Expect
     `APPEND_CONFIRMED state=REPAIR_PENDING`.
     (b) **TASK-RECOVERY task-11**, with the same tuple and `0/2`.
     `Source: TRG-24-015 / ESC-24-007 / DEC-288; project-docs/docs/phases/phase-24/qa/review-9.md`.
     `Repair authority: TRG-24-015 / DEC-288`. `Triage disposition: N/A`.
     Control fields, `Originating report`, `Exact verification` and
     `Control repair base` are `N/A`.
     `Allowed files: none; no Developer dispatch`. `Forbidden changes`:
     no edits to product code, tests, task specs, the Task 11 report, QA
     artifacts (including review-9 and the walkthrough), the validation
     contract, QA plan, decisions, workflow controls or Human fixtures; no
     fixture regeneration; no invented or edited Human answer. `Next
     action`: verify `a23b9f33` and `fcc60def` (`git rev-parse --verify`,
     exact subjects, walkthrough unchanged since `fcc60def`), then append
     (c). Expect `APPEND_CONFIRMED state=REPAIR_PENDING`.
     (c) **REPAIR task-11**, same trigger and attempt, `Outcome: RESOLVED`.
     `Ending checkpoint` is the unchanged phase HEAD. It has no new
     Developer commit, and the evidence is `a23b9f33` and `fcc60def`.
     Control and integration fields are `N/A`. Expect
     `APPEND_CONFIRMED state=COMPLETION_PENDING`.
     (d) **TASK task-11 COMPLETE.** Rerun
     `project-docs/tools/report-lint.sh project-docs/docs/phases/phase-24/tasks/task-11-done.md project-docs/docs/phases/phase-24/tasks/task-11.md`
     on the unchanged report and record its `OUTCOME COMPLETE` and
     termination. `Commit` is `93a1f792…` with its subject. Expect
     `APPEND_CONFIRMED state=ADVANCEABLE`.
  4. **Fresh QA checkpoint review.** Spawn a fresh independent QA Reviewer
     for the next unused path, `qa/review-10.md`, as the window D2
     re-review of review-9 (tasks 10-11). Scope: B1-H only. Check the nine
     walkthrough rows against VAL-054/055 and QA plan section 9, and the
     provenance (manifest, binary digest, close commit, Human attribution).
     Keep every review-9 deterministic result as evidence and don't repeat
     the full gates. Check the `fcc60def` edit note against its diff
     from `a23b9f33`. Append `REVIEW review-10` with `Re-review of: review-9`
     and `Owning task: none` for PASS or BLOCK. A PASS closes B1-H and window D2, and
     the Runner then dispatches Task 12. A REPAIR or BLOCK goes through
     ordinary TL-Lite triage. A rejected or unclear row found by QA is an
     acceptance finding for its existing owner. It is never re-answered on
     the Human's behalf.
  5. **Unchanged.** TRG-24-013/014/015, the STOP, review-9, the Task 10/11
     COMPLETE records, all failed and recovery history, VAL-054/055, the
     QA plan and every gate. These interim answers don't replace the
     phase-end HUMAN-02/HUMAN-03 rows. The QA Lead still regenerates the
     fixtures at Task 14 close.
  6. **Verification.** The Tech Lead dry-ran (a)-(d) and a review-10 PASS
     on a scratch copy of the phase files (with this decision) and the
     unchanged verifier. The live ledger validates `VALID state=STOPPED`.
     (a) through (d) returned the expected receipts in order, all exit 0:
     `REPAIR_PENDING`, `REPAIR_PENDING`, `COMPLETION_PENDING`,
     `ADVANCEABLE`. A review-10 PASS then returned
     `APPEND_CONFIRMED state=ADVANCEABLE`. Report-lint on the unchanged
     Task 11 report returned `OUTCOME COMPLETE` (exit 0). Each negative
     guard below returned exit 1 with the scratch ledger byte-unchanged:

     | Candidate | Diagnostic |
     |---|---|
     | RUN-RESUME at attempt 1 | `RUN-RESUME Attempt count must increment its governing TRIAGE` |
     | RUN-RESUME citing an unrecorded DEC-289 | `RUN-RESUME cites an unrecorded or malformed phase decision` |
     | REVIEW right after (a), no TASK-RECOVERY | `RUN-RESUME must be followed by its bound TASK-RECOVERY before REVIEW` |
     | TASK-RECOVERY with `in-scope-implementation` | `TASK-RECOVERY does not match the exact schema v2 RUN-RESUME recovery tuple` |
     | TASK-RECOVERY whose `Source` omits DEC-288 | `TASK-RECOVERY does not retain its RUN-RESUME decision authority` |
     | TASK COMPLETE before (c) | `task completed before a resolved REPAIR event: task-11` |

     A REVIEW appended straight after the STOP, with no RUN-RESUME, is
     accepted but leaves `VALID state=STOPPED`. So the Runner must not
     dispatch review-10 before (d). If any append is rejected as written,
     that is new information: return to the Tech Lead and never substitute
     field values.
- **Artifact/spec sync:** none. No PRD, TRD, design, data model, glossary,
  public interface, inventory, validation contract, QA plan, task spec,
  template, tool or prompt change. The ESC register, archive, brief and
  progress are refreshed. CF-036 records the control friction.
- **Workflow note for wrap-up:** a planned mid-phase Human checkpoint
  reaches the ledger as a QA BLOCK and a STOP. It can only resume through
  an evidence-only task recovery. Two ways to avoid this next phase: ask
  the Human for the walkthrough before the window review is dispatched, or
  give the verifier a Human-evidence resume route (CF-036).
- **Refs:** ESC-24-007; TRG-24-013, TRG-24-014, TRG-24-015; `qa/review-9.md`
  B1-H; `qa/human-walkthrough-after-task-11.md` at `a23b9f33` and
  `fcc60def`;
  `qa/qa-plan.md` sections 8 and 9; VAL-054, VAL-055; `run-ledger.md`
  2026-10-09 16:08, 16:19, 16:28; DEC-283, DEC-287.

### DEC-289: Task 13 owns convergence of the retained Concern-list hidden-summary predicate; scope amended; empty marker is `None.`

- **Date:** 2026-10-11. Tech Lead resolution of a Phase Runner return. No
  TRG or ESC was allocated: the Runner returned before any triage, and the
  ledger was never stopped. This decision changes no product behaviour,
  gate command, threshold or validation-assertion semantics.
- **Context:**
  - Task 13 checkpointed implementation `595738e7` and report `50e75fc8`
    (`RECOVERY_REQUIRED`, `gate-repair`; report-lint `OUTCOME
    RECOVERY_REQUIRED`, exit 0). Required command 5, the full tagged
    gate, failed. The in-scope cause (the sentinel helper made earlier
    Evidence stale) was repaired in scope, and a focused selection then
    passed. One cause remains:
    `cmd/fid/integration/p23_concern_list_human_test.go:29`, in
    `TestP23ConcernListHumanProjection`, rejects Human output containing
    the recorded summary `private concern prose`. Task 13's accepted A7
    change (contract section 10) requires that summary to be shown. The
    file is outside Task 13's named scope.
  - **Provenance (Tech Lead verified, `git log` and `git blame`).** Lines
    20-31, including the predicate, were introduced by `59d2653d`
    (`feat(review): checkpoint Phase 23 Human projections`,
    `Refs: phase-23/task-14`). Lines 32+ came from `cdcc4f01` (same
    task). No other test under `cmd/fid` matches the hidden-prose values.
    The file has `//go:build integration` and is in
    `./cmd/fid/integration`, so Task 13 required command 5 compiles and
    runs it.
  - **Runner rejection.** The Runner's `TASK-RECOVERY task-13` cited
    DEC-162/165/166 as authority. Those are standing Phase 21 decisions
    (`decisions_active.md`), not Phase 24 decisions, so the verifier
    correctly returned
    `INVALID TASK-RECOVERY cites an unrecorded or malformed phase decision`.
    The live ledger is unchanged, with `VALID state=ADVANCEABLE` and
    SHA-256 `82b416d1…69d07`.
  - **Planning miss.** Task 13's scope named the other Phase 23 Human
    tests but missed this one. That broke DEC-167/173 (scope must cover
    every path the gates walk). It is a Tech Lead planning defect, not a
    Developer error.
  - **Empty-section marker.** Contract section 10 A6, VAL-044 and the
    Task 13 acceptance criterion say an empty section prints `none`. The
    retained Phase 23 renderer (`59d2653d`, `p23_human.go`) has always
    printed `  None.` with no disclaimer. Task 13 kept that marker, and
    `06_DESIGN_PLAN_CLI.md` now documents it. A6 never meant to change the
    marker: it removes the disclaimer from empty sections. A literal
    reading would be a false VAL-044 failure at window E.
- **Decision:**
  1. **Ownership.** Subject-preserving convergence of the line 29
     predicate belongs to Task 13. Task 13's accepted A7 change displaced
     it, and the required tagged gate executes it. No gate, command or
     other task changes.
  2. **Scope amendment.** Task 13 gains
     `cmd/fid/integration/p23_concern_list_human_test.go`. The only
     permitted edit is to replace the
     `strings.Contains(human, "private concern prose")` rejection on line
     29 with positive checks for both of these list-view lines in `human`:
     `Summary: private concern prose` and `Reasoning: private reasoning`.
     Both values are under 80 scalars, so they appear uncut. Keep every
     other assertion and value in the file byte-unchanged. That includes
     the `err`/`stderr`, `Purpose: Inspect a raised concern` and
     `CON-human-list` checks; the JSON list; both recorded Concerns; the
     four filter subtests; and every later assertion. Record the
     provenance row (`59d2653d`, phase-23/task-14) in the report.
  3. **Empty-marker clarification.** In contract section 10 A6, VAL-044
     and the Task 13 acceptance criterion, "prints `none`" means the
     retained `None.` line. The text now says `None.`. The observable and
     its fail condition are unchanged, so no code change follows.
  4. **Ledger route** (from `ADVANCEABLE`; each append transactional,
     expected exit 0):
     (a) **TASK-RECOVERY task-13** with these exact scalars:
     `Source: project-docs/docs/phases/phase-24/tasks/task-13-done.md typed RECOVERY_REQUIRED and report-lint OUTCOME RECOVERY_REQUIRED; DEC-289`;
     `Category: gate-repair`; `Repair mode: task-scope`;
     `Repair authority: task spec; DEC-289`; `Triage disposition: N/A`;
     `Trigger signature: task-13:gate-repair:retained-concern-list-hidden-summary-predicate`;
     `Attempt count: 1`; `Consecutive no-progress attempts: 0/2`.
     `Control contract`, `Control mismatch`, `Originating report`,
     `Exact verification`, `Control repair base` and
     `Temporary repair worktree`: `N/A`.
     `Allowed files: cmd/fid/p23_human.go,cmd/fid/resume_human.go,cmd/fid/p23_human_test.go,cmd/fid/integration/p23_human_leaf_pairs_test.go,cmd/fid/integration/p23_human_fixtures_test.go,cmd/fid/integration/testdata/p23_conformance/operations.json,cmd/fid/integration/p23_conformance_test.go,cmd/fid/integration/p23_concern_list_human_test.go,project-docs/specs/GLOSSARY.md,project-docs/specs/06_DESIGN_PLAN_CLI.md,project-docs/docs/phases/phase-24/tasks/task-13-done.md`.
     `Forbidden changes`: no change to accepted A6/A7 rendering, the
     closed field list, the 80-scalar cut, escaping, the sentinel guard,
     JSON output, Task 13 gate commands, validation assertions, coverage
     policy, thresholds or budgets. In
     `p23_concern_list_human_test.go`, no edit other than clause 2. Do not
     delete, weaken, narrow, re-point or rename any other retained
     assertion. Do not edit Task 11/12/14 surfaces, task specs,
     decisions, QA artifacts, workflow controls, prompts, templates,
     dependencies, `run-ledger.md` or `triage-log.md`. Preserve the
     attempt-1 failures and provenance table in the report.
     `Checkpoint` names implementation `595738e7`, report `50e75fc8` and
     the DEC-289 commit. `Next action`: one same-task Developer
     continuation applies clause 2. It then runs all seven Task 13
     commands byte-identically, in order, with observed termination; the
     two `rg` commands must exit 1. It finalizes `task-13-done.md`
     (`Refs: phase-24/task-13`), and independent report-lint must return
     `OUTCOME COMPLETE`, exit 0. Expect
     `APPEND_CONFIRMED state=REPAIR_PENDING`.
     (b) **REPAIR task-13**, same trigger and attempt, `Outcome: RESOLVED`
     after the continuation's commits. Expect
     `APPEND_CONFIRMED state=COMPLETION_PENDING`.
     (c) **TASK task-13 COMPLETE** after report-lint `OUTCOME COMPLETE`.
     Expect `APPEND_CONFIRMED state=ADVANCEABLE`. Then dispatch Task 14.
     QA window E stays after Task 14 (QA plan section 8).
  5. **Unchanged.** Task 13's other scope, its commands and expected
     non-zero exits, VAL-044..048, VAL-052, VAL-063, the QA plan, all
     Task 1-12 records and the attempt-1 history. DEC-162/165/166 remain
     standing principles for how convergence is done. Phase 24 ledger
     authority fields cite Phase 24 decisions only.
  6. **Verification.** The Tech Lead dry-ran (a)-(c) on a scratch copy of
     the phase files (with this decision) using the unchanged verifier.
     (a) returned `APPEND_CONFIRMED state=REPAIR_PENDING`, (b)
     `APPEND_CONFIRMED state=COMPLETION_PENDING` and (c)
     `APPEND_CONFIRMED state=ADVANCEABLE`, all exit 0. A final validate
     returned `VALID state=ADVANCEABLE`. Each negative guard below
     returned exit 1 with the scratch ledger byte-unchanged:

     | Candidate | Diagnostic |
     |---|---|
     | (a) with `Repair authority: task spec; DEC-162/165/166` (the Runner's rejected append) | `TASK-RECOVERY cites an unrecorded or malformed phase decision` |
     | (a) against the decisions file before DEC-289 | `TASK-RECOVERY cites an unrecorded or malformed phase decision` |
     | TASK task-13 COMPLETE straight after (a) | `task completed before a resolved REPAIR event: task-13` |
     | (b) with a different trigger signature | `REPAIR Trigger signature does not match TASK-RECOVERY` |

     If any append is rejected as written, that is new information. Return
     to the Tech Lead and never substitute field values.
- **Artifact/spec sync:** `public-interface.md` section 10 A6 and
  `qa/validation-contract.md` VAL-044 now say `None.` (clause 3).
  `tasks/task-13.md` has the clause 2 scope and the clause 3 acceptance
  wording. `06_DESIGN_PLAN_CLI.md` already says `None.`. No change to the
  PRD, TRD, data model, glossary, inventory, QA plan, tools or prompts.
- **Workflow note for wrap-up:** when a task's accepted change reverses a
  prior phase's suppression rule, grep for that rule's retained
  predicates at planning, not only for the named markers. Here,
  `PRIVATE_` was found but the lowercase `private concern prose` was
  missed. Task specs should name DEC-162/165/166 as method, not ledger
  authority, so the Runner doesn't cite them in authority fields.
- **Refs:** `tasks/task-13.md`, `tasks/task-13-done.md` at `50e75fc8`;
  `595738e7`; `59d2653d`, `cdcc4f01`; contract sections 10 and 11;
  VAL-044, VAL-045; DEC-162/165/166, DEC-167/173 (Phase 21 standing);
  DEC-283 (scope-amendment precedent), DEC-288.

### DEC-290: `resume` shows the Issue's revision and next workflow states; Task 14 owns the change. Closes ESC-24-008

- **Date:** 2026-10-11. Tech Lead ESC resolution on the Human's
  instruction "Resolve the ESC-24-008 issue". This is a bounded
  public-interface addition: one JSON member on `Resume.issue`, two Human
  lines in the `resume` Issue block, and two flag usage strings. No failure
  code, exit, input, Event or other output changes. The Human may
  override the choice by DEC before Task 14 closes (rejected options,
  clause 2).
- **Context:**
  - TRG-24-016 (BLOCK) and the 2026-10-11 08:18 STOP. Task 14 checkpoints
    are implementation `4b6a3837` and report `98d45617` (report-lint
    `OUTCOME BLOCKED`, exit 0). `TestP24TakeControlTerminalStateProvenance`
    walks the stuck start's Human reads (`issue list`, `resume`,
    `workflow list`, `stop list`, `attention list`) and the whole advertised
    help tree. It finds no source for the `--to` value `qa_pass`, and the
    required tagged gate fails on it (exit 1).
  - **Tech Lead reproduction.** A focused `-v` run of the diagnostic at
    `98d45617` (investigation, not gate evidence) shows the stuck start in
    Human mode. `resume` prints `Status: qa_review` and the completion
    reason `issue-not-terminal-success for issue:<id>: completed workflow
    category`, but no state in that category. `workflow list` prints only
    workflow names. `issue transition --help` says `--to string Target
    state`.
  - **No output in either mode names workflow states.** `workflow list`
    JSON carries `states` and `transitions` as counts (`workflow/list.go:30`).
    `Resume.issue` is `{id,status,status_category,revision}` (Phase 21
    `public-interface.md` Resume row). The default workflows are seeded by
    migration `0007`, so a caller never typed them. This is a real product
    gap, not a provenance-class gap. DEC-278 item 9 correctly refuses
    `setup` or `caller-authored` for it.
  - **Second gap on the same step (the report's follow-up concern).** An
    Envelope-bound transition needs `--expected-revision` equal to the
    Issue's current revision, or it is refused with
    `issue-revision-conflict` (`workflow/evaluation.go:456`). No Human
    output prints that revision: `issue get` and `issue list` print
    `id | status | title`, and Human `resume` omits the `issue.revision` its
    JSON already carries. Deciding it here avoids a second stop on the same
    step.
  - **Rest of the journey.** The Tech Lead walked the remaining steps
    against current output and `--example` templates: `decision
    record-reference`, `human-outcome record`, `evidence link` and
    `envelope complete`. Their values come from `resume`, `stop list` and
    help (IDs, classes, capabilities, `OBL-security`, Envelope revision 1),
    from caller-authored IDs, keys, reasoning and files, or from `setup`
    actor IDs. Recording the approving Decision Reference and Human Outcome
    pair resolves both Stops atomically (Phase 23 DEC-250/269).
    Envelope completion does not require a closed Run. No other value
    without a class was found. That is a planning walk, not proof: Task 14
    still records every value's class.
- **Decision:**
  1. **ESC-24-008 closed.** The BLOCK was correct. `resume`, the
     take-control entry point, gains the Issue's revision and the workflow
     states it can move to, in both modes (clause 3). The journey then
     takes `--to` and `--expected-revision` as `fid-output`. DEC-278's
     five classes and VAL-049 are unchanged.
  2. **Rejected options.** (a) A Human-only line: Human output would carry
     data that JSON lacks, and a JSON caller taking control would stay
     stuck. (b) States in `workflow list`: the caller would still have to
     infer which workflow the Issue uses, and no output names it. (c) A
     new `workflow get` leaf: a new operation across every whole-tree
     corpus (VAL-010/064/065), and it still needs the workflow ID.
     (d) Lines on `issue get`: its JSON and Human output are pinned by the
     v0.1 parity goldens (`cmd/fid/testdata/v01_parity/issue-get-*.json`),
     and its object is shared with `issue list`. (e) Widening DEC-278:
     forbidden by DEC-278 item 9.
  3. **Contract.**
     - **JSON.** `Resume.issue` gains `next_states: NextState[]`, with
       `NextState = {state:string, category:string}`. The Phase 21 shape
       becomes `issue:{id,status,status_category,revision,next_states}`.
       It holds every state `S` for which
       `Registry.RulesForEdge(<issue workflow>, <status>, S)` reports an
       edge. That is the same edge predicate `issue.transition` uses
       (`issues/transition.go:185`), so edges to `cancelled` are included.
       The `Definition.Transitions` view omits them for v0.1 parity
       (`registry.go:165`) and must not be used. `category` is `S`'s
       workflow category. Items are sorted by `state`, bytewise ascending.
       It is `[]`, never `null`, when the status has no outgoing edge.
       These are edges only: rules, authority, Claim, Run, Stop and
       Obligation checks still run when the request is made. A workflow
       has at most 64 states, so no result limit is added. No other
       `resume` member, order, code or exit changes.
     - **Human.** The `resume` Issue block becomes:

       ```text
       Issue <issue-id>
         Status: <status>
         Revision: <revision>
         Next states: <state> (<category>), <state> (<category>)
         Completion: <completion status>
         Envelope: <envelope-id> (revision <n>, <status>)
       ```

       Items appear in JSON order, joined by `, `. An empty or nil list
       prints `  Next states: none.`. Values need no extra escaping: state
       names match `^[a-z][a-z0-9_]{0,63}$` (`workflow/authoring.go:11`), and
       categories are a closed set. No other `resume` line, section or
       order changes. For the take-control start (`qa_review`, revision 4)
       the lines read `Revision: 4` and
       `Next states: cancelled (cancelled), qa_block (started), qa_pass (completed)`.
     - **Help.** `issue transition` flag usages become
       `--to`: `Target workflow state; fid resume <issue-id> lists the Issue's next states`
       and `--expected-revision`:
       `Required current Issue revision; fid resume <issue-id> shows it`.
       No other help text changes, including `issue update`'s.
  4. **Journey provenance (VAL-049, unchanged classes).** The Builder's
     retried terminal-success transition records `--to qa_pass` as
     `fid-output`, sourced from the `resume` `Next states:` line. It
     records `--expected-revision` as `fid-output`, from the `Revision:`
     line of a `resume` read after the last mutation before the transition.
     The command `issue transition` is `fid-help` from `fid issue --help`.
     The acting `--actor` stays `setup`. The diagnostic
     `TestP24TakeControlTerminalStateProvenance` (Task 14's own, `4b6a3837`)
     stays and must pass. It may be tightened to require these two lines.
     Its forbidden-source guard (`--json`, `sql`, `journal`) must not be
     loosened.
  5. **Task 14 scope.** Clause 3 is the first product change this journey
     needs, and Task 14's required gates are the first to expose the gap.
     Task 14 owns it (DEC-175 precedent). Product files:
     `internal/workflow/registry.go` (a registry method enumerating
     `edges` from a state), `internal/workflow/resume.go`
     (`ResumeIssue.NextStates`, filled in `resumeIssue`),
     `cmd/fid/resume_human.go` (`resumeContextHumanLines`) and
     `cmd/fid/issue_transition.go` (the two usages). New test files:
     `internal/workflow/registry_next_states_test.go` (several edges
     including `cancelled`, sorted with categories; a state with none) and
     `cmd/fid/p24_resume_issue_test.go` (in-process Human lines for a full,
     empty and nil list; the help usages). The JSON/Human pairing row for
     the start state goes in `cmd/fid/integration/p24_take_control_test.go`,
     which is already Task 14 scope; JSON is allowed there because that row
     is not part of the journey. The ordinary race/coverage gate executes
     `internal/workflow` and `cmd/fid`, and the tagged gate executes
     `cmd/fid/integration` (DEC-167/173). The seven commands are
     unchanged. Red first: a clause 3 row fails before product code. The
     retained diagnostic is already red.
  6. **Retained convergence.** The Tech Lead found no retained test pinning
     the exact `resume` `issue` object or Issue block lines. The search
     covered `status_category`, `Completion: `, `  Status: `,
     `DisallowUnknownFields` and resume goldens across `cmd` and
     `internal`. `TestHumanRecoveryGuidance` builds a nil list, and the
     `p24_stop_corpus_test.go` `Next:` check does not match `Next states:`.
     If clause 3 still displaces a retained assertion, Task 14 owns its
     subject-preserving convergence, in any test file under `cmd/fid/`,
     `cmd/fid/integration/` or `internal/workflow/`. The only permitted
     edit is adding the two Human lines or the `next_states` member to an
     exact expected value. Record a provenance row for each edit:
     introducing commit, phase and task. Any other kind of displacement is
     new information: return it to the Tech Lead.
  7. **Recovery identity** (STOP resume; DEC-286 clause 7 precedent). The
     chain trigger stays TRG-24-016's, at attempt 2: `task-14`,
     `task-14:public-interface-decision:terminal-state-provenance-qa-pass`,
     `in-scope-implementation`, `task-scope`, `0/2`. The signature stays
     exact even though clause 3 also covers the revision. Attempt 1 is the
     BLOCKED checkpoint `4b6a3837`/`98d45617`. Its failed diagnostic,
     failed tagged gate, timing table and fixture-setup attempts stay in
     history.
  8. **Ledger sequence** (Runner, transactional append, each expected
     exit 0, each appended before any dispatch).
     (a) **RUN-RESUME** (schema v2 STOP resume):
     `Sequence contract: v2`; `Recovery directive schema: v2`;
     `Prior run/session or STOP: 2026-10-11 08:18 STOP / TRG-24-016 / ESC-24-008`;
     `Governing TRG: TRG-24-016`; `Governing ESC: ESC-24-008`;
     `Trigger signature: task-14:public-interface-decision:terminal-state-provenance-qa-pass`;
     `Attempt count: 2`; `Resolution authority: DEC-290`;
     `Continuing task: task-14`;
     `Next recovery: task-14 | in-scope-implementation | task-scope | task-14:public-interface-decision:terminal-state-provenance-qa-pass | 2`.
     `Checkpoint` names the post-DEC-290 phase HEAD and the preserved
     `4b6a3837` (implementation) and `98d45617` (report) checkpoints.
     `Resolution evidence` cites DEC-290 at the verified HEAD and the
     unchanged ledger validating `VALID state=STOPPED`. Expect
     `APPEND_CONFIRMED state=REPAIR_PENDING`.
     (b) **TASK-RECOVERY task-14**, same tuple:
     `Source: TRG-24-016 / ESC-24-008 / DEC-290; project-docs/docs/phases/phase-24/tasks/task-14-done.md`;
     `Category: in-scope-implementation`; `Repair mode: task-scope`;
     `Repair authority: TRG-24-016 / DEC-290`; `Triage disposition: N/A`;
     `Attempt count: 2`; `Consecutive no-progress attempts: 0/2`.
     `Control contract`, `Control mismatch`, `Originating report`,
     `Exact verification`, `Control repair base` and
     `Temporary repair worktree`: `N/A`.
     `Allowed files: ordinary task-14.md scope plus internal/workflow/registry.go,internal/workflow/resume.go,cmd/fid/resume_human.go,cmd/fid/issue_transition.go,internal/workflow/registry_next_states_test.go,cmd/fid/p24_resume_issue_test.go (DEC-290 clause 5), and clause 6 convergence edits only`.
     `Forbidden changes`: no `resume` JSON change other than
     `issue.next_states`; no other `resume` Human line, section or order;
     no change to `issue.transition` behaviour, validation, codes or
     exits, or to help text beyond the two clause 3 usages; no change to
     `issue get`/`list`, `workflow list`, `issue update` or the registry's
     `Definition.Transitions` view. No change to DEC-278's provenance
     classes, VAL-049 or any validation assertion, to Task 14 gate
     commands, coverage policy, thresholds or budgets. No JSON, SQL,
     journal or setup-state substitution for any journey value. Keep
     `TestP24TakeControlTerminalStateProvenance` and its forbidden-source
     guard: tighten, never loosen. Do not delete, weaken, narrow,
     re-point or rename retained assertions beyond clause 6. Do not edit
     task specs, decisions, the contract, the validation contract, QA
     artifacts, the QA plan, workflow controls, prompts, templates,
     dependencies, `run-ledger.md` or `triage-log.md`. Preserve the
     attempt-1 BLOCKED evidence.
     `Next action`: one same-task Developer continuation. It implements
     clause 3 red-first, then completes ordinary Task 14: the VAL-049
     journey with clause 4 provenance, VAL-050 parity, the VAL-051 scan,
     the VAL-053..057 fixture extensions, VAL-068 rows and the
     `08_QUALITY_PLAN.md` update. It runs all seven Task 14 commands
     byte-identically, in order, with observed termination; both `rg`
     commands must exit 1. It finalizes `task-14-done.md`
     (`Refs: phase-24/task-14`). Independent report-lint must then
     return `OUTCOME COMPLETE`, exit 0. Expect
     `APPEND_CONFIRMED state=REPAIR_PENDING`.
     (c) **REPAIR task-14**, same trigger and attempt, `Outcome: RESOLVED`
     after the continuation's commits. Expect
     `APPEND_CONFIRMED state=COMPLETION_PENDING`.
     (d) **TASK task-14 COMPLETE** after report-lint `OUTCOME COMPLETE`.
     Expect `APPEND_CONFIRMED state=ADVANCEABLE`. Then QA window E (QA
     plan section 8: tasks 12-14) also validates VAL-068.
  9. **Verification.** The Tech Lead dry-ran (a)-(d) on a scratch copy of
     the phase ledger, decisions (with this decision) and triage log, using
     the unchanged verifier. The current ledger validates
     `VALID state=STOPPED`. (a) and (b) each returned
     `APPEND_CONFIRMED state=REPAIR_PENDING`, then (c)
     `APPEND_CONFIRMED state=COMPLETION_PENDING` and (d)
     `APPEND_CONFIRMED state=ADVANCEABLE`, all exit 0. A final validate
     returned `VALID state=ADVANCEABLE`. Each negative guard below
     returned exit 1 with the scratch ledger byte-unchanged:

     | Candidate | Diagnostic |
     |---|---|
     | (a) at attempt 1 | `RUN-RESUME Attempt count must increment its governing TRIAGE` |
     | (a) against the decisions file before DEC-290 | `RUN-RESUME cites an unrecorded or malformed phase decision`; `… invalid or unverifiable Resolution authority` |
     | (b) whose `Source` omits DEC-290 | `TASK-RECOVERY does not retain its RUN-RESUME decision authority`; `… typed RUN-RESUME resolution authority` |
     | (b) with `Category: gate-repair` | `TASK-RECOVERY does not match the exact schema v2 RUN-RESUME recovery tuple` |
     | TASK task-14 COMPLETE straight after (b) | `task completed before a resolved REPAIR event: task-14` |

     The verifier also accepted (b) appended straight after the STOP with
     no RUN-RESUME (`APPEND_CONFIRMED state=STOPPED`). This is control
     friction recorded under CF-036, not a route: (a) must come first, as
     written. If any append is rejected as written, that is new
     information. Return to the Tech Lead and never substitute field
     values.
  10. **Unchanged.** TRG-24-016, the STOP, the attempt-1 report and every
      failed result; DEC-278 items 9-10; VAL-049..057 and VAL-063; Task
      14's seven commands, expected non-zero exits and budgets; the QA
      cadence (window E after Task 14, then the final review and Human
      checklist). CF-040's Human-output design revisit is not pulled
      forward.
- **Artifact/spec sync:** `public-interface.md` section 10 and its header
  gain the clause 3 contract. `qa/validation-contract.md` gains VAL-068
  (68 rows: 63 deterministic, 5 Human), with its coverage and revision log.
  `tasks/task-14.md` gains the DEC-290 amendment.
  `06_DESIGN_PLAN_CLI.md` describes the `resume` Issue block. The ESC
  register, archive, brief and progress are refreshed. No PRD, TRD, schema,
  data model, glossary, inventory, QA plan, tool or prompt change:
  PRD-PROJ-006 and the roadmap journey already require taking control
  from fid output alone, and nothing is persisted.
- **Workflow note for wrap-up:** DEC-278 R-13 fixed the provenance
  classes, but planning never walked each journey value against current
  Human output. A per-step value walk at planning would have found both
  gaps before Task 14. Journey tasks should carry that walk in their spec.
- **Refs:** ESC-24-008; TRG-24-016; `run-ledger.md` 2026-10-11 08:17 and
  08:18; `tasks/task-14.md`; `tasks/task-14-done.md` at `98d45617`;
  `4b6a3837`; Phase 21 `public-interface.md` Resume row; contract sections
  1 and 10; DEC-175, DEC-278 items 9-10, DEC-286 (route precedent).

### DEC-291: R13 runs as its rendered Markdown command; the Task 13 report-only recovery resumes. Closes ESC-24-009

- **Date:** 2026-10-11. Tech Lead ESC resolution on the Human's instruction
  "Resolve ESC-24-009". This decision changes no product code, test, gate
  command, registry, assertion, threshold or budget. It says how a command
  written in a Markdown table is read, adds that rule to the inventory
  header, and binds the ledger route from the 2026-10-11 10:43 STOP to a
  fresh window E re-review.
- **Context:**
  - Review-11 (`562583d9`, REPAIR) F1: Task 13's report does not record
    its owned R13 inventory command. The 10:26 TASK-RECOVERY task-13
    (`task-13:qa-rework:review-11-r13-inventory-receipt`, attempt 1)
    told the Developer to run
    `rg -n "p23CheckJudgmentQualifiers\|PRIVATE_" cmd` and to observe
    exit 0 with three qualifier references. The Developer ran exactly
    those bytes. It got no output and exit 1, and recorded a BLOCKED
    receipt (`f446a3be`, report-lint `OUTCOME COMPLETE`, exit 0).
    TL-Lite returned BLOCK in TRG-24-017 and raised ESC-24-009. The
    Runner appended TRIAGE and STOP (10:42, 10:43).
  - **Cause: a Markdown table escape was copied as command bytes.** The
    inventory R13 row and review-11's Repair Plan and regeneration rows
    are GitHub-Flavored Markdown table cells. A `|` inside a cell must be
    written `\|`, even inside a code span, or it splits the cell. The
    renderer removes the backslash. The Tech Lead rendered inventory
    line 94 with `markdown-it` (`gfm-like`). The cell renders as
    `<code>rg -n "p23CheckJudgmentQualifiers|PRIVATE_" cmd</code>`.
    The escaped bytes appear outside a table only where someone copied
    them from a cell: the 10:26 TASK-RECOVERY `Next action` and its
    dispatch. Outside a table the backslash is literal. Inside a
    double-quoted shell string it reaches `rg`, which reads `\|` as a
    literal pipe, so nothing matches.
  - **Every other record reads the cell this way.** (a) The R13 planning
    count (15 matches, DEC-278) and the R13 count at each earlier window
    (reviews 1 and 3-9) match only with alternation. (b) Review-5 wrote
    R13 outside a table as `rg -n "p23CheckJudgmentQualifiers|PRIVATE_" cmd`.
    (c) R2, R4, R6, R7 and R11 use the same escape. Tasks 2, 4-8 and 11
    recorded their owned commands with a plain `|`, and every QA window
    accepted those receipts. (d) Review-11's F1 prose writes the plain
    `|`. Its regeneration section says "Markdown-escaped table
    alternations were executed as regex alternation." Rendered, review-11
    agrees with itself.
  - **Tech Lead reproduction** (investigation, not gate evidence), at
    `f446a3be` from `/home/mark/projects/fid`: the escaped bytes printed
    nothing and exited 1. The rendered command printed three lines and
    exited 0: `cmd/fid/integration/p23_human_leaf_pairs_test.go:132` and
    `:152`, and `cmd/fid/integration/p23_human_fixtures_test.go:144`.
    That matches both attempts the Runner attached.
  - **The BLOCK was correct.** The directive fixed the bytes and forbade
    command changes. Deciding how canonical text is read needs Tech Lead
    authority, which TL-Lite does not have.
- **Decision:**
  1. **ESC-24-009 closed.** R13 is unchanged. Its command is the
     rendered cell, `rg -n "p23CheckJudgmentQualifiers|PRIVATE_" cmd`. Its
     required result is unchanged: exit 0, three
     `p23CheckJudgmentQualifiers` references in two test files, and no
     `PRIVATE_` reference. This is not a command change: the escaped
     spelling was never the accepted command.
  2. **Reading rule** (Phase 24, every role). A command in a Markdown
     table cell is its rendered text, and `\|` is the cell escape for
     `|`. A command copied out of a table into a ledger field, a dispatch,
     a report line or a code block is written with the plain `|`.
     "Byte-identical" compares rendered commands. The rule covers
     inventory R2, R4, R6, R7, R11 and R13, and every other command in a
     Phase 24 table (QA plan, reviews, phase-end evidence).
  3. **Reconciliation.** (a) The `input-inventory.md` header now states
     the clause 2 rule and gives R13's rendered command. The R13 row stays
     as it is; it is correct Markdown. (b) Review-11 is a committed QA
     record and is not edited. Rendered, its Repair Plan cell gives the
     same command as F1. Review-12 cites this decision. (c) The 10:26
     TASK-RECOVERY, TRG-24-017, the STOP and the `f446a3be` report
     subsection stay as history. The attempt-1 receipt stays in the report
     as a failed attempt of a command that was not R13. It is never
     relabelled as R13 evidence.
  4. **Budget.** The attempt-1 directive named bytes that are not R13 and
     forbade changing them, so it could not make progress. It does not
     count as a no-progress attempt (TRG-24-017 agrees). The count is
     `0/2`.
  5. **Recovery identity** (STOP resume; DEC-290 clause 7 precedent).
     `task-13`, `qa-rework`, `task-scope`, `0/2`, trigger
     `task-13:qa-rework:review-11-r13-command-contradiction`. That is
     TRG-24-017's trigger; the verifier binds the RUN-RESUME to it and to
     attempt 3 (TRG-24-017 is attempt 2). The earlier
     `task-13:qa-rework:review-11-r13-inventory-receipt` attempt-1 chain
     stays in history and is not renamed.
  6. **Ledger sequence** (Runner, transactional append, each expected
     exit 0, each appended before the dispatch it allows).
     (a) **RUN-RESUME** (schema v2 STOP resume):
     `Sequence contract: v2`; `Recovery directive schema: v2`;
     `Prior run/session or STOP: 2026-10-11 10:43 STOP / TRG-24-017 / ESC-24-009`;
     `Governing TRG: TRG-24-017`; `Governing ESC: ESC-24-009`;
     `Trigger signature: task-13:qa-rework:review-11-r13-command-contradiction`;
     `Attempt count: 3`; `Resolution authority: DEC-291`;
     `Continuing task: task-13`;
     `Next recovery: task-13 | qa-rework | task-scope | task-13:qa-rework:review-11-r13-command-contradiction | 3`.
     `Checkpoint` names the post-DEC-291 phase HEAD, the Task 13 report
     commits `fd6783f9` (close) and `f446a3be` (blocked receipt), and
     review-11 `562583d9`. `Resolution evidence` cites DEC-291 at the
     verified HEAD and the unchanged ledger validating
     `VALID state=STOPPED`. Expect `APPEND_CONFIRMED state=REPAIR_PENDING`.
     (b) **TASK-RECOVERY task-13**, same tuple:
     `Source: TRG-24-017 / ESC-24-009 / DEC-291; project-docs/docs/phases/phase-24/qa/review-11.md`;
     `Category: qa-rework`; `Repair mode: task-scope`;
     `Repair authority: TRG-24-017 / DEC-291; QA review; review-11`;
     `Triage disposition: N/A`; `Attempt count: 3`;
     `Consecutive no-progress attempts: 0/2`. `Control contract`,
     `Control mismatch`, `Originating report`, `Exact verification`,
     `Control repair base` and `Temporary repair worktree`: `N/A`.
     `Allowed files: project-docs/docs/phases/phase-24/tasks/task-13-done.md`.
     `Forbidden changes`: no change to product code, tests, specs, the
     inventory, task specs, decisions, the validation contract, QA
     artifacts (including review-11), the QA plan, workflow controls,
     prompts, templates, dependencies, `run-ledger.md` or
     `triage-log.md`. Do not edit, delete or relabel the attempt-1
     subsection or its exit-1 receipt. Do not rerun any product gate or
     record a new result for one. Do not copy TL, TL-Lite or QA
     reproduction into a Developer result. Do not weaken the
     three-reference, no-`PRIVATE_` result.
     `Next action`: one same-task Developer continuation edits only
     `task-13-done.md`. From `/home/mark/projects/fid` it runs the
     rendered R13 command `rg -n "p23CheckJudgmentQualifiers|PRIVATE_" cmd`
     (DEC-291 clause 1). In a new subsection it records the command, an
     output token, the working directory and the observed exit
     `exit_code: 0`, and it accounts for three
     `p23CheckJudgmentQualifiers` references in two test files and no
     `PRIVATE_` reference. It cites DEC-291 for the spelling. It
     finalizes the report (`Refs: phase-24/task-13`) and commits only the
     report. Independent two-argument report-lint must then return
     `OUTCOME COMPLETE`, exit 0. Expect
     `APPEND_CONFIRMED state=REPAIR_PENDING`.
     (c) **REPAIR task-13**, same trigger and attempt, `Outcome: RESOLVED`
     after the report commit. Expect
     `APPEND_CONFIRMED state=COMPLETION_PENDING`.
     (d) **TASK task-13 COMPLETE**, with the new report commit, after
     report-lint `OUTCOME COMPLETE`. Expect
     `APPEND_CONFIRMED state=REREVIEW_PENDING`: review-11 stays open until
     its re-review (clause 7). Do not dispatch any other task or review
     first.
  7. **Fresh window E re-review.** Spawn a fresh independent QA Reviewer
     for `qa/review-12.md`, the window E re-review of review-11. Scope:
     F1 only. It checks the new receipt and the report delta and history.
     It regenerates R13 with the rendered command, and reruns
     representative judgment and sentinel evidence, as the review-11
     Repair Plan says. It keeps every other review-11 result as evidence
     and does not repeat the full gates. Append `REVIEW review-12` with
     `Re-review of: review-11` and `Owning task: none` for PASS. A PASS
     closes F1 and window E, with expected
     `APPEND_CONFIRMED state=ADVANCEABLE`. The Runner then spawns the QA
     Lead for the final review and Human checklist. A REPAIR or BLOCK goes
     through ordinary TL-Lite triage.
  8. **Unchanged.** Task 13's product state, its seven commands and their
     receipts, VAL-044..048 and VAL-063, the inventory rows and counts,
     review-11, the QA plan, all earlier history and the window cadence.
  9. **Verification.** The Tech Lead dry-ran (a)-(d) and a review-12
     PASS on a scratch copy of the phase ledger, decisions (with this
     decision) and triage log, using the unchanged verifier. The live
     ledger validates `VALID state=STOPPED` (SHA-256 `d08356ae…cecd9f8`).
     (a) and (b) each returned `APPEND_CONFIRMED state=REPAIR_PENDING`,
     (c) `APPEND_CONFIRMED state=COMPLETION_PENDING`, (d)
     `APPEND_CONFIRMED state=REREVIEW_PENDING` and the review-12 PASS
     `APPEND_CONFIRMED state=ADVANCEABLE`, all exit 0. A final validate
     returned `VALID state=ADVANCEABLE`. Each negative guard below
     returned exit 1 with the scratch ledger byte-unchanged:

     | Candidate | Diagnostic |
     |---|---|
     | (a) at attempt 2 | `RUN-RESUME Attempt count must increment its governing TRIAGE` |
     | (a) against the decisions file before DEC-291 | `RUN-RESUME cites an unrecorded or malformed phase decision`; `… invalid or unverifiable Resolution authority` |
     | (b) whose `Source` and `Repair authority` omit DEC-291 | `TASK-RECOVERY does not retain its RUN-RESUME decision authority`; `… typed RUN-RESUME resolution authority` |
     | (b) with `Category: in-scope-implementation` | `TASK-RECOVERY does not match the exact schema v2 RUN-RESUME recovery tuple` |
     | (b) with the attempt-1 trigger `…-r13-inventory-receipt` | `TASK-RECOVERY does not preserve its RUN-RESUME checkpoint`; `… exact schema v2 RUN-RESUME recovery tuple` |
     | (b) straight after the STOP, with no (a) | `overlapping TASK-RECOVERY precedes required REPAIR` |
     | TASK task-13 COMPLETE straight after (b) | `task completed before a resolved REPAIR event: task-13` |
     | review-12 PASS straight after (c) | `QA re-review passed before owning task completed: task-13` |
     | review-12 PASS with `Re-review of: none` | `QA PASS does not reference required re-review review-11: review-12` |

     Unlike DEC-290, (b) without (a) is rejected here, because the
     attempt-1 recovery is still open. The (b) `Checkpoint` must name the
     same post-DEC-291 phase HEAD as (a). The `Next action` carries the
     plain `|` in the command, and the verifier accepted it.
     If any append is rejected as written, that is new information.
     Return to the Tech Lead and never substitute field values.
- **Artifact/spec sync:** the `input-inventory.md` header (clause 3a).
  No PRD, TRD, design, data model, glossary, public interface, validation
  contract, QA plan, task spec, template, tool or prompt change. The ESC
  register, archive, brief and progress are refreshed. CF-036 records the
  control friction.
- **Workflow note for wrap-up:** tables are a poor place for commands
  that contain `|`. In later inventories, put each registry command in a
  fenced code block (or one block per section), and let the table cite
  the registry ID. The workflow's "byte-identical" rule could also say
  that it compares rendered commands.
- **Refs:** ESC-24-009; TRG-24-017; `run-ledger.md` 2026-10-11 10:25,
  10:26, 10:42 and 10:43; `qa/review-11.md` F1, Repair Plan and R13
  regeneration; `input-inventory.md` R13; `tasks/task-13-done.md` at
  `f446a3be`; `qa/review-5.md` R13; DEC-278 (R13 correction), DEC-288 and
  DEC-290 (route precedent).

## Spec Sync Verdict

Recorded at wrap-up.

## Triage-Log Audit

Recorded at wrap-up.

## Workflow Notes

Recorded at wrap-up.
