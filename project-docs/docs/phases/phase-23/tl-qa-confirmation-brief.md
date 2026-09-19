# Phase 23 - Tech Lead Brief For The QA Lead Confirmation Pass

> Author: Tech Lead, 2026-09-19. Audience: a fresh QA Lead session.
> Subject: confirm the edits applied under **DEC-233** in response to R-19 and
> I-32/I-33/I-34 (`qa/qa-plan.md` section 16).
> Base: branch `phase-23`; the working tree was clean at commit.
> This supersedes the previous brief, which covered DEC-232. That pass is
> closed: section 16 confirmed DEC-232 in full, including that the
> `JudgmentSubject` split is better than the R-18 proposal it replaced and that
> R-17's owner derivation re-derives to task 4 alone.

## 1. What This Pass Is

The fifth confirmation cycle. Sections 12 through 16 of the QA plan all stand in
full; **nothing in any of them is reopened**. R-19 is confirmed correct in every
particular and accepted as proposed, with the follow-on you suggested. I-32,
I-33, and I-34 are all accepted.

**No product file changed.** Every edited file is under `project-docs/docs/`.
Do not re-measure the DEC-227 budgets, the DEC-223 migration chain, the DEC-224
raising sequences, the failure-code sets, or the contract row counts. Those
results are fresh and unaffected. The contract is still **79** rows.

One thing is different about this pass: **it does not end in `READY` on its
own.** Taking I-33 seriously surfaced a blocking question I cannot disposition,
raised as **ESC-23-001** and awaiting the Human. Your verdict should be `READY
EXCEPT ESC-23-001` or `NOT READY`, and should say explicitly whether you agree
that tasks 1-8 are unaffected by it. Record the outcome as a new section 17 in
`qa/qa-plan.md`.

## 2. The Dispositions, In One Line Each

- **R-19 accepted in full, owners are your five.** Confirmed independently in
  source, including the fixture path: `p21_conformance_test.go:96`,
  `p22_fake_test.go:308`, and `p22_resume_actions_test.go:111` all build
  capability lists and stage them through Kit Activation, so the rejection is on
  the ordinary setup path and not a synthetic one. Tasks 4, 5, 8, 9, and 10 each
  gain `internal/journal/current_event_value_relations.go`.
- **Your follow-on accepted.** Task 18's VAL-067 audit now asserts set equality
  between the two copies directly, rather than each against a literal.
- **I-32 accepted**, twelve not eleven, corrected in DEC-232's R-17c paragraph
  and in section 3 of this brief. Assignment unchanged.
- **I-33 accepted, and it did the work you expected it to.** The frame now reads
  "every closed journal registry a Phase 23 value or Effect reaches, whichever
  family carries it", with the two entry points named. Walking the second entry
  point found two more items, below.
- **I-34 accepted**, implementation authoritative: 19 implemented names plus six
  is **25**. The prose set is recorded as drifted and its reconciliation is a
  wrap-up spec-sync item.

## 3. My Own Findings On The Widened Frame

Three, from asking what else a Phase 23 **value** reaches. Verify each; two are
dispositioned and one is not.

1. **`semanticPhase22Events` has no Phase 23 analogue.** Owner: task 1, same
   file and change as the DEC-229 mirror. Once the nine families are in
   `semanticEventRules`, `found` is true at every profile, so
   `archive_semantics.go:328` would accept `review.recorded` inside a legacy
   projection-14 archive. VAL-065 is unaffected; the manifest check at
   `archive_validation.go:130` already rejects cross-version archives. Check my
   reading of that boolean and my claim that task 1 can prove it with no Event.
2. **The `Condition` nullable-binding relation fails open.** Inventory row 8,
   owners tasks 9, 10, 11. Same file as R-19, different case. Reproduced under
   overlay with the shape widened: a `human-outcome-recorded` Condition with all
   three new fields non-null, with none, and a `decision-satisfied` Condition
   carrying `outcome_class` are all accepted, while the Phase 22 control still
   returns `condition nullable binding mismatch`. Not a blocker - it fails open.
3. **ESC-23-001, and this one is a blocker.** Checking (2) showed the widening
   itself is not monotonic. `Condition` is the only field-adding widening of a
   live journal value shape in the phase - every other section 2 amendment
   widens a closed enum or lands on a retained family, which returns at
   `current_event_validation.go:64` before value validation. Two measured
   consequences: `validateCurrentShape` matches field count exactly and
   `ReplayJournal` (`replay.go:153`) revalidates every historical Event, so no
   workspace holding a pre-widening Stop can rebuild and `UpgradeToCurrent` has
   no hook because Phase 23 does not bump `EventSchema`; and
   `attentionSemanticKey` (`stop.go:422`) hashes the whole struct, measured to
   shift from `sha256:8aa57e2a…` to `sha256:3407ff3b…` on the same
   `budget-observed` condition, so every stored `semantic_key` stops matching and
   an identical Stop raises a duplicate Attention. Full record and three options
   in `tl-escalations.md`.

## 4. What I Want From You

1. **R-19's application.** Five task specs, one inventory row, one audit change.
   Is the 23 -> 30 split right in both copies, and does every one of the five
   tasks now name the file?
2. **Finding 3 above, hard.** This is the one that matters. Is the replay claim
   right - does `ReplayJournal` really revalidate historical Events against
   current shapes, and is there genuinely no upgrade hook? Is the semantic-key
   shift right? And is my Option A - non-applicable fields **absent** rather
   than null, with per-type Condition shapes keyed on `type` exactly as
   `currentStopDetailShapes` is keyed on `code` - actually cheaper than Option
   B, or am I under-counting its public-surface cost? Argue with the
   recommendation; the Human decides, and a wrong recommendation from me is
   worse than no recommendation.
3. **Scope of the block.** I claim only tasks 9, 10, and 11 are held, because no
   earlier task touches the `Condition` shape or struct. Verify that against all
   eighteen specs. If task 6 or any earlier task touches it, the whole phase is
   held and I need to know now.
4. **Inventory completeness, on the widened frame this time.** Eight registries,
   plus eight recorded positives (two from DEC-232 and six new). This is the
   sixth pass to find something inside the previous fix, and the fifth found its
   item only after the frame was widened, so the honest prior is still that there
   may be a ninth. The new question to ask is the one that found R-19: for every
   **value** in `public-interface.md` section 2, which live Phase 21/22 family
   carries it, and what does that family's path validate?
5. **Mechanical.** 79 contract rows unchanged; no VAL row renumbered or
   reworded; no product file touched; no protected definition, command,
   threshold, or budget changed. The one `public-interface.md` edit (I-34) adds
   and removes no name.

## 5. What Would Make This NOT READY

A ninth unowned registry; my ESC-23-001 evidence failing source verification;
the block being wider than tasks 9-11; or the `ErrorSubjectType` reading being
wrong. Anything smaller is an improvement, and improvements do not block.

Note that `READY EXCEPT ESC-23-001` is a real and expected outcome here: the
escalation is a Human contract decision, not a defect in the plan's ownership,
and tasks 1-8 are ready to dispatch without it.
