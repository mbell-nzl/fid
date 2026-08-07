# Phase 19 - TL Escalations (archived)

> Archived at wrap-up (2026-08-08, DEC-106). All eight escalations are
> **resolved**. Identifiers are phase-scoped (`ESC-19-NNN`, DEC-084) and were
> allocated from `triage-log.md`. Full reasoning lives in `decisions.md` and
> `triage-log.md`; these entries are the durable summary.

**Phase total: eight escalations, zero product defects.** Every one landed in
verification apparatus or evidence records. Four (003/004/005/007) were in one
observable's shell ceremony and a fifth (006) was fallout from it.

## Resolved

- **ESC-19-008 — final-review BLOCK on five findings (2026-08-02 → DEC-105,
  2026-08-07; closed by `qa/final-review-2.md` PASS WITH IMPROVEMENTS and
  DEC-106, 2026-08-08).** **None required a product code change.**
  (1) **VAL-026** completed by the Human: **9 PASS**, raising CF-026 (workflow
  inspection + lifecycle), split into **CF-026a** (read surface, schedulable
  now) and **CF-026b** (lifecycle model, needs a decision first).
  (2) **Architecture contradiction is real but the code is right** —
  `internal/workflow` legitimately uses the journal write API under DEC-089;
  DEC-096's `Effect`-only *characterization* is superseded while its operative
  test-packaging rule stands. (3) **Stale canonical records** accepted and
  synced; QA Lead sustained the finding in part, and PRECON-01/02 landed in the
  wrap commit. (4) **Missing `Refs:` on `32614bb`** accepted as a deviation;
  **history is not rewritten** (the hash is cited 30 times across 12 files), and
  enforcement moves to CF-023. QA Lead corrected the supporting evidence: the
  gap is two commits, `32614bb` and `ce4cf5b`, **both from the same irregular
  Task 4 session**, against every other task commit in range carrying the
  trailer. (5) **Structure budgets carried forward as CF-027**: `create.go` was
  **419** lines at phase start and is **412** now — inherited debt the phase
  *improved*, sequenced behind CF-025's parity test.
  **Follow-on findings from the subsequent review, all dispositioned by
  DEC-106:** MAJOR-01 (`02_TRD.md` §3.1's lead sentence was false for ten
  packages — a phase-introduced defect, now fixed), MINOR-01 (`writeMarkdown`
  68 lines → CF-027), MINOR-02 (`ebbcd16` outside the commit convention →
  CF-023).

- **ESC-19-007 — the mandated VAL-059 manifest command never worked, and Task 4's
  report had never passed report-lint (2026-08-02 → DEC-103).** Review 2 re-ran
  all six Task 4 gates **green** (lint `0 issues.`) and confirmed candidate
  parity (six `OK`, clean `diff -r`), then found two evidence defects.
  (1) `sha256sum -c cmd/fid/testdata/v01_parity/SHA256SUMS` from the repo root
  **cannot resolve** the manifest's bare filenames; the manifest was right and
  the command was wrong from the day it was written, and the obvious repair
  (running it inside the artifact directory) is *also* wrong — it checks the
  baselines against their own digests. **DEC-103** made the authoritative form a
  `mktemp -d` candidate capture checked from that directory by absolute manifest
  path, and struck both bad forms.
  (2) `task-4-done.md` failed report-lint 7×. The linter was **correct** — it
  was refusing to close a task on unverified evidence. Fixed by a **Developer
  closure evidence session** (`ba65078`) that re-ran the gates and rewrote the
  report with its own tokens; Review 2's tokens are QA Lead's and were **not**
  transcribed (DEC-102 item 5). **VAL-059 = PASS.** **CF-025 promoted to first
  and mandatory in Task 6**; no sixth manual run authorized.

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
  phase-start; Task 4 → *Complete (gates unverified)* with its lint line
  corrected to FAIL; the one-line repair authorized inside Task 5, whose
  repo-wide green run verified Task 4's gate 2 of 6; **uncaptured gate results
  are `UNVERIFIED`, never `PASS`** — a named future control is not a substitute.

- **ESC-19-005 — Task 4 Step R cleanup blocked by the filesystem sandbox
  (2026-08-01 → DEC-101).** `git worktree remove` / `git worktree prune` failed
  with `Read-only file system`, leaving the path registered-but-missing and
  cascading into four dependent commands; an authorized rerun then reproduced
  the pinned baseline with all six `SHA256SUMS` entries `OK`. **The stop matched
  the letter of Step R and the reasoning was wrong:** cleanup proves nothing
  about parity, and the Tech Lead had fenced it inside the evidence set with a
  fixed scratch path, so one failed `rm` gated the whole sequence. **DEC-101:**
  cleanup is never evidence; scratch state is unique per attempt (`mktemp -d`);
  task specs declare required environment capabilities and escalate shortfalls
  at step zero. Task 4's gates were then run by the Human **without token
  capture** — recorded honestly as human-attested, with QA Lead's Review 2 rerun
  as the control. **That control fired and the attestation failed (DEC-102).**
  **CF-025** raised.

- **ESC-19-004 — Task 4's first VAL-059 baseline capture failed (2026-08-01 →
  DEC-099/DEC-100).** `capture.sh` resolved a relative output root after
  `cd`-ing into its temp workspace. The task stopped with no commit, six gates
  unrun, and a registered scratch worktree. **The BLOCK rightly stopped the run
  and wrongly reasoned about why:** the failure was a path bug in a *Task 4
  deliverable being authored*, it wrote nothing so could not contaminate
  anything, and DEC-098's guard is **structural** (detached worktree at the
  pinned commit), not temporal. **DEC-099:** preserve the partial work, settle
  baseline provenance by **reproduction against `SHA256SUMS`**, finish normally
  — no restart. **DEC-100:** standing rule — a task-authored harness failing
  while being authored is a **deviation, not a failed gate**.

- **ESC-19-003 — VAL-059 baseline instruction conflicted with the required
  immutable-baseline procedure (2026-07-30 → DEC-098).** Task 4 said
  `git rev-parse HEAD`; QA-plan §6 requires the Task 3 *implementation* commit,
  and at Review 1 the two differed. **QA was right; the task spec had drifted.**
  Baseline pinned to `50533e2b813bf073cd75f75eb70dd09cd15e9040`, moving
  selectors prohibited, Task 4 synchronized to QA-plan §6, normalization fixed
  to the committed no-op. `qa-plan.md` needed no edit.

- **ESC-19-002 — Task 2 lint gate blocked by a test-only import cycle
  (2026-07-29 → DEC-096/DEC-097).** The binding `[]journal.Effect` return type
  makes `workflow → journal → db`, so the in-package
  `internal/db/workflow_tables_test.go` closed a cycle and failed `typecheck`.
  **DEC-096** confirmed the edge and authorized moving the two tests to the
  external `db_test` package; no production or API change. **DEC-097**
  root-caused the missing triage record as a **Codex provider-prompt defect**
  and fixed `project-docs/agents/tl-lite.md`. Backfilled as TRG-19-002.

- **ESC-19-001 — Task 1 VAL-031 evidence token (2026-07-26 → DEC-094/DEC-095).**
  The report recorded a no-match `rg` as `PASS - exit 0 (no output)`; `rg` exits
  1 on no match. The observable itself held. Report corrected; TRG-19-001
  backfilled and the TL-Lite write failure dispositioned.

## Wrap-Up Note

Five of the eight are traceable to Tech Lead artifacts: a task spec that had
drifted from the QA plan, a required command never executed by its author, an
evidence set with cleanup fenced inside it and a fixed scratch path, and — in
the wrap-up sync itself — a false generalisation in `02_TRD.md` §3.1 that the
phase introduced. The single lesson the phase paid the most for is **DEC-104
item 4**: express an observable as an automated test unless the validation
contract records why it cannot be.
