# TL Escalations - Current

> **Purpose.** Open items that exceed local-decision authority and need full
> Tech Lead or human input.
>
> **Scope.** Open items only. Resolved phase items move to the phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

_None._

---

## Recent Archive Notes

- **Phase 20 (in flight)** — **ESC-20-003 resolved by DEC-132 (2026-08-10).**
  The Go command had not terminated: a wrapper returned only empty `output`
  and discarded the live nested `session_id`. Canonical Codex prompts now
  preserve full results, poll yielded sessions, and require `exit_code` or an
  explicit termination signal. `report-lint.sh` now ignores wrapped blocker
  prose while retaining its token checks. Task 3d and all exact gates pass;
  the integration gate proved the protocol by staying silent through multiple
  polls before `cmd/fid` completed in 160.765s. Broader performance work is
  deliberately deferred to the Human-requested pre-Task-4 follow-up pack.
- **Phase 20 (in flight)** — **ESC-20-002 resolved by DEC-131 (2026-08-09).**
  Both Review 1b Blocking findings sustained; **no product defect**. Init's
  derived-prefix rejection is contract-correct (§3.2.1); three e2e journeys are
  stale fixtures that never converged, hidden because every Phase 20 full-suite
  `go test` gate was recorded with the impossible token `exit 0 (no output)`.
  Those lines are now **UNVERIFIED**. Recovery: **task-3c** (fixture
  convergence, `max_excerpt_bytes` rows, real-token full gate) and **task-3d**
  (report-lint token control, report corrections), then `qa/review-1c.md`, then
  Task 4. A third BLOCK returns to the Human. **CF-029** extends the control to
  QA review artifacts.
- **Phase 20 (in flight)** — **ESC-20-001 resolved by DEC-130 (2026-08-09).**
  Review 1's two Blocking findings were sustained without weakening VAL-002 or
  VAL-043. Recovery is two bounded test-led tasks, **task-3a** (exhaustive
  reason-code detail/subject/redaction matrix) and **task-3b** (strict-config
  and init boundary matrix, dangling-symlink collision regression, config file
  split), then a bounded re-review `qa/review-1b.md`, then Task 4. Detail lives
  in `../phases/phase-20/decisions.md` DEC-130 and
  `../phases/phase-20/tl-escalations.md`.

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
