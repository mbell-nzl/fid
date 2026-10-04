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

## Spec Sync Verdict

Recorded at wrap-up.

## Triage-Log Audit

Recorded at wrap-up.

## Workflow Notes

Recorded at wrap-up.
