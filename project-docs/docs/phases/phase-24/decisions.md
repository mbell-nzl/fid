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

## Spec Sync Verdict

Recorded at wrap-up.

## Triage-Log Audit

Recorded at wrap-up.

## Workflow Notes

Recorded at wrap-up.
