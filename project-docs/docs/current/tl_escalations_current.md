# TL Escalations - Current

> **Purpose.** Open items that exceed local-decision authority and need full
> Tech Lead or human input.
>
> **Scope.** Open items only. Resolved phase items move to the phase archive.
>
> **Size target:** under 80 lines.

## Open Escalations

None. **All five final-review blocking findings are dispositioned (DEC-105);
QA Lead's subsequent final-review disposition is the remaining gate.**

## Resolved This Phase (archive at wrap)

- **ESC-19-008 — final-review BLOCK on five findings (2026-08-02 → DEC-105,
  2026-08-07).** **None required a product code change.**
  (1) **VAL-026** completed by the Human: **9 PASS**, raising CF-026 (workflow
  inspection + lifecycle), which Tech Lead split into **CF-026a** (read surface,
  schedulable now) and **CF-026b** (lifecycle model, needs a decision first).
  (2) **Architecture contradiction is real but the code is right** —
  `internal/workflow` legitimately uses the journal write API under DEC-089;
  DEC-096's `Effect`-only *characterization* is superseded while its operative
  test-packaging rule stands. (3) **Stale canonical records** accepted; Tech Lead
  syncs at wrap. (4) **Missing `Refs:` on `32614bb`** accepted as a deviation —
  10 sibling commits carry it, so the convention holds and only enforcement is
  missing (→ CF-023); **history is not rewritten**, since the hash is cited in
  eight places. (5) **Structure budgets carried forward as CF-027**:
  `create.go` was **419** lines at `50533e2` and is **412** now — inherited debt
  that Phase 19 *improved*, to be decomposed after CF-025's parity test makes it
  safe.

Full reasoning lives in `../phases/phase-19/decisions.md` and
`../phases/phase-19/triage-log.md`; these are pointers.

- **ESC-19-007 — the mandated VAL-059 manifest command never worked, and Task 4's
  report has never passed report-lint (2026-08-02 → DEC-103).** Review 2 re-ran
  all six Task 4 gates **green** (lint `0 issues.`) and confirmed candidate
  parity (six `OK`, clean `diff -r`), then found two evidence defects.
  (1) `sha256sum -c cmd/fid/testdata/v01_parity/SHA256SUMS` from the repo root
  **cannot resolve** the manifest's bare filenames — `WARNING: 6 listed files
  could not be read`. The manifest is right; my command was wrong from the day I
  wrote it, and the obvious repair (running it inside the artifact directory) is
  *also* wrong — it checks the baselines against their own digests. **DEC-103**
  makes the authoritative form a `mktemp -d` candidate capture checked from that
  directory by absolute manifest path; both bad forms are struck.
  (2) `task-4-done.md` fails report-lint 7× (missing `## Tests run`; six
  commands `UNVERIFIED`/`FAIL`). The linter is **correct** — it is refusing to
  close a task on unverified evidence. Fixed by a **Developer closure evidence
  session** that re-runs the gates and rewrites the report with its own tokens;
  Review 2's tokens are QA Lead's and are **not** transcribed (DEC-102 item 5).
  **VAL-059 = PASS** on Review 2's executed evidence. **CF-025 is promoted to
  first and mandatory in Task 6**; no sixth manual run is authorized.

- **ESC-19-006 — Task 5's lint block was misclassified as pre-existing; it
  disproved Task 4's attested lint PASS (2026-08-02 → DEC-102).**
  `golangci-lint run ./...` failed `* errcheck: 1` at
  `internal/workflow/resolve.go:61`. The Developer called it pre-existing on the
  strength of `git show HEAD` and declined the out-of-scope repair; TL-Lite
  adopted that framing. **Stopping was right and declining was right, but the
  provenance was never checked:** `git log --diff-filter=A` returns **`32614bb`**,
  Task 4's own commit, and `task-4-done.md` lists the file as `(new)`. So the
  gate Task 4 recorded as `PASS (human-attested)` **could not have passed**.
  **DEC-102:** "pre-existing" requires an introducing commit older than
  phase-start; DEC-101's void condition has triggered (early, from a Developer
  rather than QA Lead — honoured in substance); Task 4 → **Complete (gates
  unverified)** with its lint line corrected to FAIL and four gates relabelled
  `UNVERIFIED`; the one-line repair is authorized inside Task 5, whose repo-wide
  green run verifies Task 4's gate 2 of 6; **uncaptured gate results are
  `UNVERIFIED`, never `PASS`** — a named future control is not a substitute.

- **ESC-19-005 — Task 4 Step R cleanup blocked by the filesystem sandbox
  (2026-08-01 → DEC-101).** `git worktree remove` / `git worktree prune` failed
  with `Read-only file system`, leaving the path registered-but-missing and
  cascading into four dependent commands; an authorized rerun then reproduced
  the pinned baseline with all six `SHA256SUMS` entries `OK`. **The stop matched
  the letter of Step R and the reasoning was wrong:** cleanup proves nothing
  about parity, and I had fenced it inside the evidence set with a fixed scratch
  path, so one failed `rm` gated the whole sequence. **DEC-101:** cleanup is
  never evidence; scratch state is unique per attempt (`mktemp -d`); task specs
  declare required environment capabilities and escalate shortfalls at step
  zero. Baseline provenance is settled by the six `OK` and closed. Task 4's
  gates were then run by the Human **without token capture** — recorded
  honestly as human-attested, with **QA Lead's Review 2 rerun as the control**;
  any failure there voids the disposition and reopens Task 4. **CF-025** carries
  the VAL-059 golden-file test conversion to Task 6.
- **ESC-19-004 — Task 4's first VAL-059 baseline capture failed (2026-08-01 →
  DEC-099/DEC-100).** `capture.sh` resolved a relative output root after
  `cd`-ing into its temp workspace. The task stopped with no commit, six gates
  unrun, and a registered scratch worktree. **The BLOCK rightly stopped the run
  and wrongly reasoned about why:** the failure was a path bug in a *Task 4
  deliverable being authored*, it wrote nothing so could not contaminate
  anything, and DEC-098's guard is **structural** (detached worktree at the
  pinned commit), not temporal. **DEC-099:** preserve the partial work, settle
  baseline provenance by **reproduction against `SHA256SUMS`**, remove the
  worktree via `git worktree remove`, finish normally — no restart. **DEC-100:**
  standing rule — a task-authored harness failing while being authored is a
  **deviation, not a failed gate**. Task 4 cleared to resume.
- **ESC-19-003 — VAL-059 baseline instruction conflicted with the required
  immutable-baseline procedure (2026-07-30 → DEC-098).** Task 4 said `git
  rev-parse HEAD`; QA-plan §6 requires the Task 3 *implementation* commit, and
  at Review 1 the two differed. **QA was right; the task spec had drifted.** The
  baseline is now pinned to `50533e2b813bf073cd75f75eb70dd09cd15e9040`, moving
  selectors are prohibited, Task 4 is synchronized to QA-plan §6, and
  normalization is the committed no-op. `qa-plan.md` needed no edit.
- **ESC-19-002 — Task 2 lint gate blocked by a test-only import cycle
  (2026-07-29 → DEC-096/DEC-097).** The binding `[]journal.Effect` return type
  makes `workflow → journal → db`, so the in-package
  `internal/db/workflow_tables_test.go` closed a cycle and failed `typecheck`.
  **DEC-096** confirms the `Effect`-only edge and authorizes moving the two
  tests to the external `db_test` package; no production or API change.
  **DEC-097** root-causes the missing triage record as a **Codex
  provider-prompt defect** and fixes `project-docs/agents/tl-lite.md`;
  `sync-agents.sh` has been run and verified. Backfilled as TRG-19-002.
- **ESC-19-001 — Task 1 VAL-031 evidence token (2026-07-26 → DEC-094/DEC-095).**
  The report recorded a no-match `rg` as `PASS - exit 0 (no output)`; `rg` exits
  1 on no match. The observable itself held. Report corrected; TRG-19-001
  backfilled and the TL-Lite write failure dispositioned.

## Recent Archive Notes

- **Phase 18** — ESC-001 (final-review BLOCK, `workflow.Load` deadlock) resolved
  by DEC-072; five other triages adjudicated as DEC-067..071 without an
  `ESC-NNN`. In five of six cases TL-Lite reported an escalation it never minted
  and wrote no triage log — fixed by DEC-085, tracked as **CF-023**. Archived to
  `../phases/phase-18/tl-escalations.md`.
- **Phase 17** — five escalations, all resolved; two retrospective items became
  **CF-023**. Archived to `../phases/phase-17/tl-escalations.md`.
- **Phase 16** — no `ESC-NNN`; DEC-043..048 and TRG-046 archived in
  `../phases/phase-16/`.

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
