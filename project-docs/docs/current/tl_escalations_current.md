# TL Escalations - Current

> **Purpose.** Open items that exceed local-decision authority and need full
> Tech Lead or human input.
>
> **Scope.** Open items only. Resolved phase items move to the phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

None.

---

## Recent Archive Notes

- **Phase 20 (in flight)** — **ESC-20-007 resolved by DEC-137, with the
  systemic fix in DEC-138 (2026-08-11).** Task 5's work was finished and
  passing; the dedicated integration gate failed because
  `TestIntegrationClassification` compares against a hardcoded journey list that
  Task 5 legitimately grew by one. The list was stale, so the failure was
  correct — but exact equality against a name list is a change-detector that
  already caused ESC-20-004 and would trip again at Tasks 14-18. DEC-137
  replaces it with required-presence (losing a migrated journey still fails,
  additions are allowed) and leaves DEC-134's repository-wide tagged-test walk
  untouched. **DEC-138** addresses the pattern: ESC-20-005/006/007 all stopped
  finished work for a one-line, self-describing correction to the verification
  apparatus, so a stale bookkeeping inventory is now an in-session repair with a
  mandatory deviation note, and TL-Lite gains a matching fourth RESUME ground.
  Thresholds, product assertions, and evidence-token rules are untouched: you
  may fix the bookkeeping, never the verdict.

- **Phase 20 (in flight)** — **ESC-20-006 resolved by DEC-136 (2026-08-11).**
  Two causes. The `Pending commit` fields were a **Tech Lead defect**: DEC-135
  ordered a single commit, contradicting the Developer prompt's Report And
  Commit Sequence and `report-lint.sh`'s two-step protocol, so the compliant
  Developer left the placeholders. DEC-135 decision 5 is withdrawn, `task-4.md`
  is corrected, and task specs may no longer restate process owned by
  `AGENTS.md` and the role prompts. The truncated tokens `0.009s`,
  `coverage: 80.5% of statements`, and `4.625s` are **sustained**: they drop the
  `ok`/`FAIL`/`PASS` marker, so the DEC-131 control fired correctly and is not
  being relaxed. The repair session re-runs the exact sequence itself, since a
  remembered fragment may not be expanded, then finalizes the report against
  commit `e56b3a0` in a `docs(task): ...` commit. **CF-030** opened.
- **Phase 20 (in flight)** — **ESC-20-005 resolved by DEC-135 (2026-08-11).**
  Task 4's required lint gate correctly failed on dead code: the typed
  `LockUnavailableError` carries only `LockPath` per contracts §2.4/§2.5, so
  `describeLockHolder` lost its caller. No product defect. DEC-135 deletes the
  formatter, refuses `//nolint` suppression and any holder data in the closed
  detail object, and re-runs Task 4 unchanged in one fresh session that inherits
  the uncommitted tree, completes the exact gate sequence with captured tokens,
  commits once, and rewrites `task-4-done.md` with a deviation note quoting the
  original failing token. All seven `BLOCKED` VAL rows re-record from that
  session's own evidence.
- **Phase 20 (in flight)** — **ESC-20-004 resolved by DEC-134 (2026-08-10).**
  Task 3f's direct-parent classifier omitted an integration-tagged in-process
  test in `internal/issues`, leaving it outside both accepted gates. Correction
  `75c856a` makes the inventory repository-wide and removes only the misplaced
  build constraint; the same reviewer returned NON-BLOCK at `cd9e6bc`. No
  product behavior or assertion changed.
- **Phase 20 (in flight)** — **ESC-20-003 resolved by DEC-132 (2026-08-10).**
  The Go command had not terminated: a wrapper returned only empty `output` and
  discarded the live nested `session_id`. Canonical Codex prompts now preserve
  full results, poll yielded sessions, and require `exit_code` or an explicit
  termination signal; `report-lint.sh` ignores wrapped blocker prose while
  keeping its token checks. Broader performance work became the approved
  infrastructure pack.
- **Phase 20 (in flight)** — **ESC-20-002 resolved by DEC-131 (2026-08-09).**
  Both Review 1b findings sustained; **no product defect**. Init's
  derived-prefix rejection is contract-correct (§3.2.1); three e2e journeys were
  stale fixtures, hidden because every Phase 20 full-suite `go test` gate was
  recorded with the impossible token `exit 0 (no output)`. Those lines are now
  **UNVERIFIED**. Recovered by tasks 3c/3d and `qa/review-1c.md`. **CF-029**
  extends the token control to QA review artifacts.
- **Phase 20 (in flight)** — **ESC-20-001 resolved by DEC-130 (2026-08-09).**
  Review 1's two findings were sustained without weakening VAL-002 or VAL-043.
  Recovered by bounded test-led tasks 3a (reason-code detail/subject/redaction
  matrix) and 3b (strict-config and init boundary matrix, symlink regression,
  config split) plus re-review `qa/review-1b.md`. Detail lives in
  `../phases/phase-20/decisions.md` and `../phases/phase-20/tl-escalations.md`.

- **Phase 19** — **eight escalations, ESC-19-001..008, all resolved, zero
  product defects.** Every one landed in verification apparatus or evidence
  records: four (003/004/005/007) in one observable's shell ceremony for
  VAL-059, one (006) in its fallout — a false gate attestation disproven by a
  lint regression the same task introduced — one (001) an evidence-token
  defect, one (002) a test-only import cycle, and one (008) the final review's
  five findings, none in product code. Resolved by DEC-094..DEC-106.
  All eight carried real, unique, phase-scoped IDs (DEC-084); the ID-collision
  defect that affected Phases 17-18 did not recur. Three of eight triage records
  required a Tech Lead backfill (CF-023). Archived to
  `../phases/phase-19/tl-escalations.md`.
- **Phase 18** — ESC-001 (final-review BLOCK, `workflow.Load` deadlock) resolved
  by DEC-072; five other triages adjudicated as DEC-067..071 without an
  `ESC-NNN`. In five of six cases TL-Lite reported an escalation it never minted
  and wrote no triage log — fixed by DEC-085, tracked as **CF-023**. Archived to
  `../phases/phase-18/tl-escalations.md`.
- **Phase 17** — five escalations, all resolved; two retrospective items became
  **CF-023**. Archived to `../phases/phase-17/tl-escalations.md`.

## Standing Escalation Triggers

TL-Lite escalates instead of deciding locally when an issue:

1. Changes a public interface, CLI contract, output shape, journal event schema,
   data model, persistence contract, dependency, tool, quality gate, or workflow.
2. Contradicts `decisions_active.md`, `decisions_principles.md`, `AGENTS.md`, or
   a recorded phase decision.
3. Affects more than one task/subsystem, exposes product-intent ambiguity, or
   changes testing strategy / coverage policy.
4. Resumes context-exhausted or partial work spanning schema/data model,
   dependencies, public interfaces, canonical docs, or multiple Developer
   sessions.

QA findings escalate here if they have architecture implications, ambiguous
acceptance criteria, unclear security/privacy severity, contradict a standing
decision, span tasks/subsystems beyond the review window, plausibly indicate a
material spec deviation, or expose a quality-gate anomaly not anticipated by the
QA plan.
