# Phase 19 Task 4 — Executive Briefing

**Date:** 2026-08-02 · **Author:** Tech Lead · **Status:** Awaiting your decision

---

## The one-line version

The code is written and passing. We are stuck on the *paperwork* that proves it
works — and specifically on a temp-directory cleanup command that failed on a
file permission.

## Where things actually stand

| | Status |
|---|---|
| Task 4 product code | **Written, tests passing** |
| Proof it doesn't break existing behaviour | **Obtained and verified** |
| Final quality checks (6 commands) | Not yet run |
| Committed to git | **No — 3 sessions of work sitting uncommitted** |

Nothing is wrong with the software. We have stopped three times, and all three
stops were in the verification scaffolding, never in the product.

## What happened

Task 4 changes issue creation — the most heavily used code path in the product.
To prove we didn't break it, the plan says: build the *old* version of the app,
record its output, then confirm the *new* version produces identical output.

Reasonable. But that proof needs a stack of temporary machinery: a second
checkout of the repo, a separately built binary, a custom capture script, and
checksum files. Six moving parts, none of which are the product, all of which
must succeed in one clean run before the real checks may start.

Three stops, three different pieces of that machinery breaking:

1. **Stop 1** — two documents disagreed about which commit to compare against.
2. **Stop 2** — the capture script had a path bug (it was being written by this
   very task).
3. **Stop 3 (current)** — the cleanup command that deletes a leftover temp
   folder hit a read-only filesystem and failed. Everything after it fell over.

**Critically, on stop 3 the work was then redone successfully.** The comparison
ran, and all six checks came back `OK`. We have the proof. The Developer stopped
anyway, because the instructions I wrote said to stop if any step failed — and
that cleanup line was inside the "must not fail" fence, where it did not belong.

## Root cause

Two things, both mine:

- **The proof mechanism is heavier than the thing it's proving.** Each stop was
  a different brick falling out of scaffolding that shouldn't have been this
  tall. We kept fixing the individual brick instead of asking whether the
  scaffolding was the right size.
- **Uncommitted work auto-escalates.** Because Task 4 has unsaved work spanning
  multiple sessions on a critical file, the rules require *every* interruption
  — however trivial — to come to me for a decision. The task can no longer fail
  small. That's an accident of process design, not a real risk signal.

The Developer followed my instructions correctly each time. This is a spec
problem, not an execution problem.

## Proposed path forward

1. **Commit the finished work now.** Stops the auto-escalation loop immediately.
2. **Accept the proof we already have.** It was verified six ways; redoing it
   buys nothing.
3. **Turn the one-off ceremony into a permanent automated test.** Same guarantee,
   stronger over time, and the fragile temp machinery goes away for good.
4. **Stop treating cleanup as evidence.** Use a fresh temp folder each time so no
   run ever depends on the last one being tidied up.

Then Task 4 finishes normally: run the six checks, commit, done. My estimate is
one clean Developer session.

## The ask

**Approve, and I'll write it up and unblock Task 4 today.**

If you'd rather keep it minimal, I can do items 1, 2 and 4 now and defer the
test conversion (item 3) to Task 6. Slightly less durable, marginally faster.

No other phase work is blocked behind this — Tasks 5 and 6 depend on Task 4, so
the phase is idle until it clears.

---

*Full detail: `tl_escalations_current.md` (ESC-19-005), `triage-log.md`
(TRG-19-005), `run-ledger.md`.*
