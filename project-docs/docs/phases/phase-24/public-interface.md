# Phase 24 Public Interface Amendment

> **Status: DRAFT for Human review (planning, 2026-10-03).** Authority:
> DEC-272. Items marked **[D2]**..**[D4]** await the Human decision named in
> section 13. Everything else is Tech Lead contract detail under DEC-272 and
> becomes binding when the Human approves this file with the task list.
> D1 is decided (DEC-273).
> Phase 21, 22 and 23 `public-interface.md` remain the baseline for every
> leaf, value, code and exit this file does not replace. Inventory and task
> ownership: [`input-inventory.md`](input-inventory.md).

## 1. Principle

A Human must be able to run any fid command unaided, using only fid's own
output (DEC-272 A8). `--input` JSON stays the single write contract. Callers
send only values a check actually uses. fid observes repository identity
itself and records what it observed. Outputs and the journal keep the full
`Repository` shape.

## 2. Repository Observation

- **What fid observes.** fid observes the Git worktree that contains the
  selected project (`--project`, or the current directory). It uses the
  existing `ObserveRepository` capture: a double read that must agree, `.fid/`
  excluded, ignored files excluded, and unsafe states rejected.
- **The recorded value.** fid records the full `Repository`: `project_id` from
  project config, `worktree_id` = SHA-256 of the absolute worktree root,
  `capture_source:"git"`, `head_oid`, `tree_oid`, `base_oid` = observed HEAD,
  `dirty`, and `dirty_digest`. This is the shape that `resume` and the
  judgment leaves already record.
- **No input names a worktree by ID.** A read that must name a different
  worktree uses a path flag (`--repository <path>`), as `resume`,
  `run inspect` and `export create` already do. Parallel worktrees are
  Phase 25.
- **When.** An observing mutation observes at most once, under the write
  lock, after closed-input validation and after the idempotency replay check.
  A replay never observes (section 7).
- **Unobservable repository.** If observation fails (not a Git worktree,
  unmerged paths, assume-unchanged/skip-worktree entries, or a change during
  capture), the request returns `repository-unavailable` (section 6).
  `run.close` and `run.abandon` keep their Phase 22 `recovery-state-conflict`.

## 3. Shared Input Values

```text
ExpectedRepository = {head_oid:oid, dirty:boolean}
SourceInput = {source_type:"project-file"|"external-reference",
  display_label:string[1..128], path:text|null, section:text|null,
  external_id:ID|null, optional(digest:digest)}
```

- **ExpectedRepository.** fid compares `head_oid` and `dirty` with its
  observation. Any difference rejects with `repository-mismatch`, with no
  Event. A Human can supply both values from `git rev-parse HEAD` and
  `git status --porcelain`.
- **SourceInput** is the Phase 21 `Source` without `repository`. The
  path/external-ID cross-field rule is unchanged.
  - `project-file`: `digest` is optional. If present, fid reads `path` and
    rejects on a mismatch. If omitted, fid computes and records the digest.
    `path` resolves under the project root with the Phase 21 relative-path
    rule. It must be a regular file that stays inside the root after symlink
    resolution. A missing, unreadable or escaping file rejects with
    `source-digest-mismatch`, `observed_digest:null`.
  - `external-reference`: `digest` is required, because fid cannot observe
    external content. It is stored as supplied.
- **Stored and output `Source`** keep the Phase 21 shape. `digest` is the
  supplied or observed value. `repository` is fid's observation for that
  command.

## 4. Per-Leaf Input Amendments (handoff item a)

17 of the 39 `--input` leaves change. The other 22 keep their input shape but
gain help, an `--example` template and field-named validation (section 8).

| Leaf | Removed from input | Kept or added expectation | fid observes and records |
|---|---|---|---|
| `gate.register` | `worktree_id`, `repository` | none | Registration `project_id`, `worktree_id` and binding digest from the observation |
| `envelope.create` | `repository_scope`; `repository` in `risk_result.source`, `sources[]`, `obligations[].source` | Source `digest` optional per section 3 | `repository_scope` = observation; Sources per section 3; existing registration/worktree checks compare observed values |
| `envelope.complete` | `repository` (already discarded and replaced by observation today) | none | unchanged: completion already evaluates against the observation |
| `claim.acquire` | `worktree_id` | none | Claim `worktree_id` = observed worktree; it must equal the Envelope scope's worktree (today's rule, now checked against reality) |
| `run.start` | `repository` | `expected_repository` required | `start_repository` = observation; Envelope comparison per section 5 (DEC-273) |
| `run.close` | `repository` | `expected_repository` required | `terminal_repository` = observation; obligations evaluated against it |
| `run.abandon` | `repository` | `expected_repository` required | `terminal_repository` = observation |
| `gate.run` | `repository` | `expected_repository` required | One observation before spawn is recorded as Attempt `start_repository`, Receipt `terminal_repository` and Event `repository` (today all three are the caller's block). Run comparison per section 5 (DEC-273) |
| `evidence.link` | `target_repository` | none | `target_repository` = the subject's recorded repository (Receipt terminal, Review repository, or the Human Outcome Event's observation) |
| `evidence.evaluate` | `target_repository` | none | freshness is evaluated against the current observation, as terminal transition and completion already are |
| `review.record` | `repository`; `repository` in coverage items, exclusions and `evidence_refs` | `expected_repository` required; Source digests optional | `repository` = observation (today the caller's block must equal it exactly) |
| `concern.record` | `repository` in `references[]` | Source digests optional | references checked or observed per section 3 |
| `deviation.record` | `repository` in `references[]` | same | same |
| `decision.record-reference` | `repository` in `source` | `source.digest` optional | already checked when supplied; now observed when omitted |
| `human-outcome.record` | `repository` in `references[]` | Source digests optional | references become checked (today stored with nothing checked) |
| `requirement.record-reference` | `repository` in `source` | `source.digest` optional | as `decision.record-reference` |
| `run.observe-budget` | `repository` in `measurement_ref` | Source digest optional | `measurement_ref` per section 3 |

The Event `repository` of every Event above is a fid observation or a stored
fid-observed value, never caller input. Claim and Finding Events keep the
Envelope scope, which is now itself observed.

## 5. Comparisons That Change Meaning (D1, DEC-273)

Today `run.start` requires the caller's block to equal the Envelope's stored
`repository_scope` byte for byte. `gate.run` requires it to equal the Run's
`start_repository`. Neither is ever compared with the real worktree.

**Decided (DEC-273, option A).** Each comparison becomes *same project and
same worktree*. The caller's `expected_repository` is checked against the
real worktree instead. A Run or gate run may therefore start after edits or
commits made since the Envelope or Run began. Freshness (PRD-EVID-003)
stales proof whose recorded tree differs from the current one. The rejected
option B (the real worktree must equal the snapshot) would have refused
every rework Run until an Envelope revision (Phase 25), and blocked a
Builder from testing their own edits within a Run.

## 6. Failure Codes

No code is removed. Category, exit, retryable and subject are unchanged
unless stated.

| Code | Change |
|---|---|
| `repository-mismatch` | Details become `{expected:ExpectedRepository, observed:Repository}`. Producers become `run.start`, `run.close`, `run.abandon`, `gate.run` (before Attempt append) and `review.record`. `run.close`/`run.abandon` expectation mismatch moves here from `recovery-state-conflict`. |
| `source-digest-mismatch` | Details become `{display_label:text, path:?text, expected_digest:?digest, observed_digest:?digest}`. `expected_digest` is null when the caller omitted the digest. Producers become every leaf that accepts a `SourceInput`. |
| `repository-unavailable` | **Implemented for the first time** (Phase 21 prose listed it; `internal/cli/codes.go` never did). precondition / retryable true / exit 4 / subject null; details `{reason:"not-a-repository"\|"unmerged-paths"\|"tracked-flags"\|"changed-during-capture"\|"unsafe-path", detail:text}`. Producers: every observing mutation except `run.close`/`run.abandon`. |
| `idempotency-conflict` | Details gain `differing_fields:string[]` (sorted unique, nonempty). Each name is a public input name: a top-level `--input` field, a flag spelling (`--title`), a positional argument name, or `--actor`. See section 7. |
| `invalid-argument` | Detail **shape unchanged**. For a structural `--input` failure, `argument` names the field path (`input.coverage.items[0].source.path`). `value` is the offending key, JSON type or value (null when missing). `allowed` lists the accepted keys, type or members. Today Human mode prints Go decoder text that JSON omits; after this change JSON carries it (required by section 9). |

## 7. Request Digest Rule (handoff item b) [D3]

- The request digest covers the input **as the caller wrote it**: the
  canonical form of the closed input object, plus operation, project ID and
  resolved Actor (as today). An optional member is present exactly when the
  caller wrote it. The DEC-242 rule (`inapplicable_findings: []` canonicalizes
  to omission) is retained as the only exception.
- fid-observed values (repository observation, observed Source digests,
  current clock) never enter the digest. Implementations compute the digest
  from the decoded input before observing and never mutate the request after
  that.
- **Consequence.** A retry with the same key and input after HEAD moved, or
  after a referenced file changed, replays the original result. It does not
  re-check, re-observe or fail.
- **Naming the differing field [D3].** Every keyed mutation also stores one
  digest per top-level public input name, `--actor` included. On a conflict,
  `differing_fields` lists the names whose digest differs, or that are
  present on only one side. Storage: migration `0019` adds
  `idempotency_results.request_field_digests_json` (JSON object, name →
  digest). The idempotency Effect of every keyed Event carries it, so rebuild
  reproduces it. Projection schema becomes 19 and manifest
  `projection_schema` 19 at phase end. Down migration restores schema 18
  exactly (DEC-264).
- Journals written before this change are unsupported and fail closed
  (ADR-0011, DEC-270). No task criterion requires them to replay.

## 8. Input Help And Templates (Part 2) [D4]

- **Single source.** One input descriptor per `--input` leaf (39 leaves)
  drives three things: the leaf's `--help` input section, its `--example`
  template, and closed structural validation. A conformance test proves each
  descriptor's field tree equals the Go request type's JSON field tree.
  Cross-field and domain rules stay in their domain packages.
- **Structural validation** covers unknown, duplicate, missing and null
  members; JSON type; closed enum membership; and declared formats (ID,
  digest, OID, relative path, integer range, string bound). Failures name the
  field path (section 6).
- **`--help`.** Every leaf in the live tree, reads included, has a one-line
  description. Every `--input` leaf's help states its purpose and lists each
  input field: path, type, required/optional/nullable, allowed values or
  bounds, and a one-line meaning. It then names the template command.
  `attention resolve --help` no longer says "Run JSON input".
- **`--example` [D4].** It prints one editable JSON object to stdout and exits
  0. It reads no workspace (it works outside a project), needs no `--actor`,
  `--idempotency-key` or `--input`, and writes nothing.
  - Keys appear in documented order with 2-space indentation.
  - Required fields only. Optional fields are listed in `--help`.
  - Values are correctly typed placeholders: IDs like `"<run-id>"`, minimal
    valid integers, `false`, the first allowed enum member, `null` for
    nullable fields, and one element in each required nonempty array.
  - Any other flag or argument with `--example` is `invalid-argument`, and so
    is `--json --example`, because the template is not a `fid-cli.v1`
    envelope.
- **Record IDs stay caller-chosen** (DEC-272).

## 9. Human Failure Lines (Part 3 A5; handoff item c)

This **replaces** Phase 21 `public-interface.md` section 1's "`error:
<message>` plus at most two non-semantic guidance lines". Every Human-mode
failure on stderr is exactly these lines, in this order, each at most once:

```text
error: <error.message from the JSON failure, verbatim>
reason: <error.code>[; <key detail>]
blocked by: <Stop and Attention IDs>     # only when details carry nonempty stop_ids/attention_ids
next: <resolving action>                 # only for the codes in the next-table below
```

- Every line derives from the JSON failure (code, message, subject and
  details). Exit status and the JSON failure are unchanged.
- The `error:` line repairs A5. Today `issue transition` prints a different
  sentence from its JSON `message`.
- **Key detail.** Scalar detail members render as `name=value` in key order;
  scalar arrays are comma-joined. A closed override table renders nested
  details:
  - `repository-mismatch`: `expected HEAD <short> clean|dirty, observed HEAD <short> clean|dirty`
  - `idempotency-conflict`: `differing fields: <names>`
  - `transition-not-permitted`, `obligation-unsatisfied`, `stop-unresolved`, `resolution-invalid`: reason codes joined
  - `invalid-argument`: `argument=<path>, allowed=<a|b>`

  Values are control-character-escaped, and the line is cut to 200 Unicode
  scalars with `…`.
- **next-table (closed):**

| Code | `next:` text (IDs and classes from details) |
|---|---|
| `invalid-argument` on an `--input` leaf | `see fid <leaf> --help; template: fid <leaf> --example` |
| `repository-mismatch` | `check HEAD and uncommitted changes, then correct expected_repository` |
| `source-digest-mismatch` | `correct or omit digest for <display_label>` |
| `idempotency-conflict` | `resend the original input, or use a new --idempotency-key` |
| `stop-unresolved` | `read each Stop's qualifying resolution: fid stop list --status open` |
| `review-rework-exhausted` | `an activated Human records an approved Human Outcome of class review-rework-extension: fid human-outcome record (template: fid human-outcome record --example)` |
| `protected-decision-required` | `record the <decision_class> Decision Reference and an approving Human Outcome: fid decision record-reference, fid human-outcome record (templates: --example)` |
| `human-outcome-required` | `an activated Human records a <outcome_class> Human Outcome, then link it: fid human-outcome record, fid evidence link (templates: --example)` |
| `obligation-unsatisfied` | `see per-Obligation reasons: fid evidence evaluate --input <file> (template: fid evidence evaluate --example)` |

## 10. Human Views (Part 3 A4, A6, A7) [D2]

- **A4 Stops and Attention.** `stop list`, `attention list` and `resume`
  render each Stop as labelled lines and never print raw JSON: subject, code,
  blocked actions, one line per StopDetails member for its closed variant,
  resolution condition, resolver capability, status, and a `Next:` line keyed
  on the Condition type:
  - `executor-terminated` → `fid run observe`
  - `run-recovered` → `fid run recover`
  - `activation-recovered` → `fid activation recover`
  - `capability-available` → no command; wait for the capability
  - `budget-observed` → `fid run observe-budget`
  - `work-ended` → `fid run close` or `fid run abandon`
  - `review-recorded`, `finding-resolved` → `fid review record`
  - `human-outcome-recorded` → `fid human-outcome record`
  - `decision-satisfied` → `fid decision record-reference` then `fid human-outcome record`

  After the fact, `fid stop resolve`, then `fid attention resolve`, where no
  lifecycle Event disposed it.
- **Repository lines.** Every Human repository line shows short HEAD (12 hex
  digits) and `clean`/`dirty` only. Recorded-versus-current comparisons say
  `same worktree` or `different worktree`. No Human view prints `worktree_id`,
  `tree_oid`, `dirty_digest`, `capture_source` or a nil marker. JSON is
  unchanged.
- **A6.** Keep each view's `Purpose:` line. Each section that lists Human
  Outcomes, Concerns or Deviations prints its judgment disclaimer exactly
  once, before the items. Items carry none. An empty section prints `none`
  and no disclaimer. The on/off switch is Phase 26.
- **A7 recorded reasoning.** Human views show recorded reasoning. The closed
  list:
  - Concern `summary` and `reasoning`;
  - Concern disposition `reason`;
  - Finding `summary` and every FindingDisposition `reason`;
  - Deviation `expected`, `observed` and `reasoning`;
  - Human Outcome `reasoning`.

  Lists show one line: whitespace collapsed, cut to 80 Unicode scalars with
  `…`. Detail views show full text. Control characters are always escaped
  (PRD-OPS-007).
  - Phase 23 section 7's "never reproduces project prose" now means never
    reproducing the **content of a referenced project file**.
  - The `PRIVATE_` leak checks are replaced by a guard that seeds referenced
    files with sentinel content and proves no Human view prints it.
  - GLOSSARY distinguishes project-owned prose from recorded reasoning.
- **[D2] Detail views.** Concerns, Deviations and Human Outcomes have no
  detail view today; only list leaves exist.
  - **Proposed (option A):** add read leaves `concern.get` /
    `fid concern get <id>`, `deviation.get` and `human-outcome.get`, with
    results `{concern:Concern}`, `{deviation:Deviation}` and
    `{outcome:HumanOutcome}`. Not-found uses the existing
    `concern-not-found`, `deviation-not-found` and `outcome-not-found`, which
    gain these producers.
  - Mutation outputs and `review get` are also detail views.

## 11. Restated Assertions (handoff item d)

Phase 23 VAL-060 and VAL-073 are restated as new Phase 24 assertions. They
are not weakened:

- **VAL-060 →** Human text never presents a Human Outcome, Concern or
  Deviation as deterministic proof. Every section listing any of them
  carries the judgment disclaimer exactly once, and no item line claims
  proof. Human text never reproduces the content of a referenced project
  file. It shows recorded reasoning: one cut line in lists, full text in
  detail views.
- **VAL-073 →** Human views follow the accepted guidance: labelled sections,
  an operation purpose line, judgment disclaimers once per section, and no
  narration of routine state.

## 12. Non-Goals

- Flag-based input forms.
- The A6 on/off switch (Phase 26).
- The CF-040 full Human-output design revisit.
- Naming worktrees by path in inputs, or parallel worktrees (Phase 25).
- Envelope revision.
- Gate terminal re-observation after process exit.
- Compatibility with pre-Phase-24 journals.
- Capability manifest changes: no row is added, removed or re-advertised.
- Changes to cross-field domain validation messages beyond section 6.

## 13. Decisions For The Human

| ID | Question | Tech Lead recommendation |
|---|---|---|
| D1 | Does `run.start` / `gate.run` require the real worktree to equal the Envelope/Run snapshot, or only the same project and worktree? | **Decided A (DEC-273)**: same project and worktree, plus the caller's expectation |
| D2 | Where do Humans read full recorded reasoning for Concerns, Deviations and Human Outcomes? | A: add `concern get`, `deviation get`, `human-outcome get` |
| D3 | Naming the differing field needs per-field digests stored in the journal and migration 0019. Accept the cost? | Accept |
| D4 | `--example` rules in section 8 (static typed placeholders, required fields only, refuses `--json`) | Accept |
