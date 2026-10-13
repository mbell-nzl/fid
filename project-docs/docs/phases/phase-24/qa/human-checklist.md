# Phase 24 Human Verification Checklist

> Phase: 24 - Input authoring and Human take-control
> Date: 2026-10-11
> Tester: _______________

This is the phase-end Human subset. Record pass/fail/partial, quoted output, attempts and your own explanation in each row. Prior Task 11 answers remain interim. Exact output/count/grammar/privacy checks are recorded by QA in final-review.md.

Use the fresh binary `/home/mark/projects/fid/.cache/phase-24-human/fixtures-2666841344/fid` and each workspace below. For live journey/authoring, stay in that original physical root. Source project files belong inside that root; write temporary --input JSON files outside it, for example in a fresh /tmp directory. Read only Human fid output, --help/--example, and Git commands the help names. Acting Actor IDs and project directory are supplied setup; new record IDs, keys, reasoning and files are your choices.

If the prepared Claims have expired, QA reruns the exact generator command in final-review.md and refreshes workspace references. Captured Human-only text remains available below for view judgments. Do not use the internal QA manifest or completed test input files to obtain journey/authoring values.

Fixture refresh (DEC-292 clause 1(c)), 2026-10-13. The original set `fixtures-2145711523` expired about 11:50 NZDT on 2026-10-12. At the Human's instruction, a Tech Lead session reran the exact generator command in final-review.md at HEAD `0006eccc711c0de35ceb03091cc1c1fbfa2acccf` (19:56:47 to 19:57:47 NZDT, exit 0). The product is unchanged since Task 14 (`86e0c26a`): `git diff --name-only 86e0c26a 0006eccc` lists only `project-docs/` paths. The binary embeds `vcs.revision=0006eccc…` and `vcs.modified=true`. The only uncommitted file was the Runner-owned `run-ledger.md`, and the original set was built the same way.
- New set: `/home/mark/projects/fid/.cache/phase-24-human/fixtures-2666841344` (59 manifest rows).
- Manifest SHA-256: `4315ddb90d93e4020cab42608d244d2f3a53d2ec5d80f80ba602a2d86c5cd9f2`.
- Binary SHA-256: `b478ecfe9bc8c3c2521b3fd0c7ee542003c93dba42d8be37633c1d84473e2c1c`.
- Expiry: the live Claims (86400 s TTL) expire about 19:57 NZDT on 2026-10-14.
- Discarded run: an earlier run, with `-v` added (`fixtures-2591929048`), was not the exact command; it was deleted unused.
- What changed: the binary and workspace paths, and the HUMAN-01 start output. That output now shows the new set's record IDs: [fixtures-2666841344/take-control-start.txt](fixtures-2666841344/take-control-start.txt), copied byte for byte from the new set's start refusal. The original is kept in `final-evidence/`.
- What did not change: the captured outputs for HUMAN-02/03/04 still come from the original set.

## Human takes control - VAL-053

Workspace: `/home/mark/projects/fid/.cache/phase-24-human/fixtures-2666841344/take-control`. Acting Actors: `human` and `agent:builder`. Begin with [Human output](fixtures-2666841344/take-control-start.txt) showing the already-refused terminal-success step. Discover the rest from fid --help, the refusal, and subsequent Human reads. Do not copy the setup refusal’s transition arguments as your own solution. Record each source line used.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-01a | VAL-053 | From the initial refusal alone, explain why the step was refused. (human) | You can state the protected Decision reason. | | |
| HUMAN-01b | VAL-053 | From that refusal alone, choose the next legitimate command. (human) | The output gives you an actionable next step. | | |
| HUMAN-01c1 | VAL-053 | Find the acting Issue using Human fid reads. (human) | The Issue identifier is understandable without outside knowledge. | | |
| HUMAN-01c2 | VAL-053 | Find its Envelope using Human fid reads. (human) | The Envelope identifier and scope are understandable without outside knowledge. | | |
| HUMAN-01d | VAL-053 | Author and submit the Decision Reference using discovered help/template and your own contained project file. (human) | You can supply all required values from permitted sources. | | |
| HUMAN-01e | VAL-053 | Author and submit the approving Human Outcome using discovered help/template. (human) | You can record the intended approval without outside assistance. | | |
| HUMAN-01f | VAL-053 | Read the Stops after the approval pair and judge whether resolution is clear. (human) | You can tell the blocker is resolved. | | |
| HUMAN-01g | VAL-053 | Discover and link the qualifying Human Outcome Evidence. (human) | The Human output and help let you choose the required linkage. | | |
| HUMAN-01h | VAL-053 | Use the latest Human resume to choose the terminal-success state. (human) | The available completed state is understandable from the output. | | |
| HUMAN-01i | VAL-053 | Use the latest Human resume to choose the current expected Issue revision. (human) | You can find the precondition without JSON or guessing. | | |
| HUMAN-01j | VAL-053 | Retry the refused terminal transition using your discovered values. (human) | You can move the Issue to terminal success. | | |
| HUMAN-01k | VAL-053 | Discover and submit envelope completion. (human) | You complete the Envelope from fid’s own output/help. | | |
| HUMAN-01l | VAL-053 | Identify any step that required source, specs, JSON internals, .fid contents or other outside knowledge. (human) | No step requires a forbidden source; name any gap. | | |

## Failure clarity - VAL-054

For each captured denial, write the reason and the next action separately. Judge from the Human text alone.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-02-01r | VAL-054 | Explain why idempotency-conflict refused this request. [Human output](final-evidence/human-output/idempotency-conflict.txt) (human) | Reason is clear from the output. | | |
| HUMAN-02-01n | VAL-054 | Choose the next legitimate action after idempotency-conflict. [Human output](final-evidence/human-output/idempotency-conflict.txt) (human) | Next action is clear from the output. | | |
| HUMAN-02-02r | VAL-054 | Explain why repository-mismatch refused this request. [Human output](final-evidence/human-output/repository-mismatch.txt) (human) | Reason is clear from the output. | | |
| HUMAN-02-02n | VAL-054 | Choose the next legitimate action after repository-mismatch. [Human output](final-evidence/human-output/repository-mismatch.txt) (human) | Next action is clear from the output. | | |
| HUMAN-02-03r | VAL-054 | Explain why invalid-argument refused this request. [Human output](final-evidence/human-output/invalid-argument.txt) (human) | Reason is clear from the output. | | |
| HUMAN-02-03n | VAL-054 | Choose the next legitimate action after invalid-argument. [Human output](final-evidence/human-output/invalid-argument.txt) (human) | Next action is clear from the output. | | |
| HUMAN-02-04r | VAL-054 | Explain why source-digest-mismatch refused this request. [Human output](final-evidence/human-output/source-digest-mismatch.txt) (human) | Reason is clear from the output. | | |
| HUMAN-02-04n | VAL-054 | Choose the next legitimate action after source-digest-mismatch. [Human output](final-evidence/human-output/source-digest-mismatch.txt) (human) | Next action is clear from the output. | | |
| HUMAN-02-05r | VAL-054 | Explain why protected-decision-required refused this request. [Human output](final-evidence/human-output/protected-decision-required.txt) (human) | Reason is clear from the output. | | |
| HUMAN-02-05n | VAL-054 | Choose the next legitimate action after protected-decision-required. [Human output](final-evidence/human-output/protected-decision-required.txt) (human) | Next action is clear from the output. | | |

## Stop, Attention and repository views - VAL-055

Accept or reject each individual view/state, explaining what made it clear or unclear.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-03-01 | VAL-055 | Judge stop.list in fixture clean-same-stop.list: [Human output](final-evidence/human-output/clean-same-stop.list.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-02 | VAL-055 | Judge attention.list in fixture clean-same-attention.list: [Human output](final-evidence/human-output/clean-same-attention.list.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-03 | VAL-055 | Judge resume in fixture clean-same-resume: [Human output](final-evidence/human-output/clean-same-resume.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-04 | VAL-055 | Judge review.get in fixture clean-same-review.get: [Human output](final-evidence/human-output/clean-same-review.get.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-05 | VAL-055 | Judge stop.list in fixture clean-different-stop.list: [Human output](final-evidence/human-output/clean-different-stop.list.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-06 | VAL-055 | Judge attention.list in fixture clean-different-attention.list: [Human output](final-evidence/human-output/clean-different-attention.list.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-07 | VAL-055 | Judge resume in fixture clean-different-resume: [Human output](final-evidence/human-output/clean-different-resume.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-08 | VAL-055 | Judge review.get in fixture clean-different-review.get: [Human output](final-evidence/human-output/clean-different-review.get.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-09 | VAL-055 | Judge stop.list in fixture dirty-same-stop.list: [Human output](final-evidence/human-output/dirty-same-stop.list.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-10 | VAL-055 | Judge attention.list in fixture dirty-same-attention.list: [Human output](final-evidence/human-output/dirty-same-attention.list.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-11 | VAL-055 | Judge resume in fixture dirty-same-resume: [Human output](final-evidence/human-output/dirty-same-resume.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-12 | VAL-055 | Judge review.get in fixture dirty-same-review.get: [Human output](final-evidence/human-output/dirty-same-review.get.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-13 | VAL-055 | Judge stop.list in fixture dirty-different-stop.list: [Human output](final-evidence/human-output/dirty-different-stop.list.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-14 | VAL-055 | Judge attention.list in fixture dirty-different-attention.list: [Human output](final-evidence/human-output/dirty-different-attention.list.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-15 | VAL-055 | Judge resume in fixture dirty-different-resume: [Human output](final-evidence/human-output/dirty-different-resume.txt) (human) | You accept this view’s readability and repository context. | | |
| HUMAN-03-16 | VAL-055 | Judge review.get in fixture dirty-different-review.get: [Human output](final-evidence/human-output/dirty-different-review.get.txt) (human) | You accept this view’s readability and repository context. | | |

## Judgment presentation - VAL-056

Review each individual judgment view below. Structural disclaimer counts and escaping already have deterministic evidence; here judge whether the section labels, purpose, qualifier and reasoning communicate their intended meaning.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-04-01a | VAL-056 | requirement.record-reference Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-requirement.record-reference.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-01b | VAL-056 | requirement.record-reference Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-requirement.record-reference.txt) (human) | Purpose is clear. | | |
| HUMAN-04-01d | VAL-056 | requirement.record-reference Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-requirement.record-reference.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-01e | VAL-056 | requirement.record-reference Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-requirement.record-reference.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-02a | VAL-056 | requirement.link Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-requirement.link.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-02b | VAL-056 | requirement.link Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-requirement.link.txt) (human) | Purpose is clear. | | |
| HUMAN-04-02d | VAL-056 | requirement.link Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-requirement.link.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-02e | VAL-056 | requirement.link Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-requirement.link.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-03a | VAL-056 | concern.record Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-concern.record.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-03b | VAL-056 | concern.record Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-concern.record.txt) (human) | Purpose is clear. | | |
| HUMAN-04-03d | VAL-056 | concern.record Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-concern.record.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-03e | VAL-056 | concern.record Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-concern.record.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-03c1 | VAL-056 | Judge how the Concern section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-concern.record.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-04a | VAL-056 | concern.dispose Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-concern.dispose.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-04b | VAL-056 | concern.dispose Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-concern.dispose.txt) (human) | Purpose is clear. | | |
| HUMAN-04-04d | VAL-056 | concern.dispose Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-concern.dispose.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-04e | VAL-056 | concern.dispose Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-concern.dispose.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-04c1 | VAL-056 | Judge how the Concern section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-concern.dispose.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-05a | VAL-056 | deviation.record Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-deviation.record.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-05b | VAL-056 | deviation.record Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-deviation.record.txt) (human) | Purpose is clear. | | |
| HUMAN-04-05d | VAL-056 | deviation.record Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-deviation.record.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-05e | VAL-056 | deviation.record Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-deviation.record.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-05c1 | VAL-056 | Judge how the Deviation section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-deviation.record.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-06a | VAL-056 | decision.record-reference Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-decision.record-reference.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-06b | VAL-056 | decision.record-reference Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-decision.record-reference.txt) (human) | Purpose is clear. | | |
| HUMAN-04-06d | VAL-056 | decision.record-reference Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-decision.record-reference.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-06e | VAL-056 | decision.record-reference Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-decision.record-reference.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-07a | VAL-056 | human-outcome.record Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-human-outcome.record.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-07b | VAL-056 | human-outcome.record Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-human-outcome.record.txt) (human) | Purpose is clear. | | |
| HUMAN-04-07d | VAL-056 | human-outcome.record Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-human-outcome.record.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-07e | VAL-056 | human-outcome.record Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-human-outcome.record.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-07c1 | VAL-056 | Judge how the Human Outcome section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-human-outcome.record.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-08a | VAL-056 | review.record Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-review.record.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-08b | VAL-056 | review.record Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-review.record.txt) (human) | Purpose is clear. | | |
| HUMAN-04-08d | VAL-056 | review.record Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-review.record.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-08e | VAL-056 | review.record Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-review.record.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-09a | VAL-056 | finding.address Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-finding.address.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-09b | VAL-056 | finding.address Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-finding.address.txt) (human) | Purpose is clear. | | |
| HUMAN-04-09d | VAL-056 | finding.address Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-finding.address.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-09e | VAL-056 | finding.address Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-finding.address.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-10a | VAL-056 | review.get Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-review.get.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-10b | VAL-056 | review.get Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-review.get.txt) (human) | Purpose is clear. | | |
| HUMAN-04-10d | VAL-056 | review.get Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-review.get.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-10e | VAL-056 | review.get Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-review.get.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-11a | VAL-056 | review.list Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-review.list.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-11b | VAL-056 | review.list Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-review.list.txt) (human) | Purpose is clear. | | |
| HUMAN-04-11d | VAL-056 | review.list Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-review.list.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-11e | VAL-056 | review.list Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-review.list.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-12a | VAL-056 | finding.list Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-finding.list.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-12b | VAL-056 | finding.list Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-finding.list.txt) (human) | Purpose is clear. | | |
| HUMAN-04-12d | VAL-056 | finding.list Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-finding.list.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-12e | VAL-056 | finding.list Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-finding.list.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-13a | VAL-056 | concern.list Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-concern.list.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-13b | VAL-056 | concern.list Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-concern.list.txt) (human) | Purpose is clear. | | |
| HUMAN-04-13d | VAL-056 | concern.list Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-concern.list.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-13e | VAL-056 | concern.list Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-concern.list.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-13c1 | VAL-056 | Judge how the Concern section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-concern.list.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-14a | VAL-056 | deviation.list Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-deviation.list.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-14b | VAL-056 | deviation.list Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-deviation.list.txt) (human) | Purpose is clear. | | |
| HUMAN-04-14d | VAL-056 | deviation.list Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-deviation.list.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-14e | VAL-056 | deviation.list Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-deviation.list.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-14c1 | VAL-056 | Judge how the Deviation section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-deviation.list.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-15a | VAL-056 | decision.list Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-decision.list.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-15b | VAL-056 | decision.list Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-decision.list.txt) (human) | Purpose is clear. | | |
| HUMAN-04-15d | VAL-056 | decision.list Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-decision.list.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-15e | VAL-056 | decision.list Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-decision.list.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-16a | VAL-056 | human-outcome.list Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-human-outcome.list.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-16b | VAL-056 | human-outcome.list Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-human-outcome.list.txt) (human) | Purpose is clear. | | |
| HUMAN-04-16d | VAL-056 | human-outcome.list Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-human-outcome.list.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-16e | VAL-056 | human-outcome.list Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-human-outcome.list.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-16c1 | VAL-056 | Judge how the Human Outcome section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-human-outcome.list.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-17a | VAL-056 | requirement.get Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-requirement.get.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-17b | VAL-056 | requirement.get Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-requirement.get.txt) (human) | Purpose is clear. | | |
| HUMAN-04-17d | VAL-056 | requirement.get Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-requirement.get.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-17e | VAL-056 | requirement.get Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-requirement.get.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-18a | VAL-056 | requirement.list Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-requirement.list.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-18b | VAL-056 | requirement.list Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-requirement.list.txt) (human) | Purpose is clear. | | |
| HUMAN-04-18d | VAL-056 | requirement.list Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-requirement.list.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-18e | VAL-056 | requirement.list Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-requirement.list.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-19a | VAL-056 | deviation.get Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-deviation.get.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-19b | VAL-056 | deviation.get Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-deviation.get.txt) (human) | Purpose is clear. | | |
| HUMAN-04-19d | VAL-056 | deviation.get Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-deviation.get.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-19e | VAL-056 | deviation.get Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-deviation.get.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-19c1 | VAL-056 | Judge how the Deviation section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-deviation.get.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-20a | VAL-056 | human-outcome.get Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-human-outcome.get.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-20b | VAL-056 | human-outcome.get Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-human-outcome.get.txt) (human) | Purpose is clear. | | |
| HUMAN-04-20d | VAL-056 | human-outcome.get Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-human-outcome.get.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-20e | VAL-056 | human-outcome.get Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-human-outcome.get.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-20c1 | VAL-056 | Judge how the Human Outcome section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-human-outcome.get.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-21a | VAL-056 | concern.get Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-concern.get.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-21b | VAL-056 | concern.get Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-concern.get.txt) (human) | Purpose is clear. | | |
| HUMAN-04-21d | VAL-056 | concern.get Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-concern.get.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-21e | VAL-056 | concern.get Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-concern.get.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-21c1 | VAL-056 | Judge how the Concern section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-concern.get.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-22a | VAL-056 | resume Do the labels make the judgment section boundaries clear?: [Human output](final-evidence/human-output/judgment-resume.txt) (human) | Sections are easy to identify. | | |
| HUMAN-04-22b | VAL-056 | resume Is each displayed section’s purpose understandable?: [Human output](final-evidence/human-output/judgment-resume.txt) (human) | Purpose is clear. | | |
| HUMAN-04-22d | VAL-056 | resume Does the view keep routine progress narration separate from the judgment you are reading?: [Human output](final-evidence/human-output/judgment-resume.txt) (human) | Judgment content is easy to find without routine narration. | | |
| HUMAN-04-22e | VAL-056 | resume Where reasoning is displayed, is it readable; where none is displayed, is that presentation understandable?: [Human output](final-evidence/human-output/judgment-resume.txt) (human) | You can understand the recorded reasoning or its absent presentation. | | |
| HUMAN-04-22c1 | VAL-056 | Judge how the Concern section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-resume.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-22c2 | VAL-056 | Judge how the Deviation section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-resume.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |
| HUMAN-04-22c3 | VAL-056 | Judge how the Human Outcome section communicates its judgment/proof limits (or its empty state): [Human output](final-evidence/human-output/judgment-resume.txt) (human) | The qualifier is clear without distracting repetition; an empty section is understandable. | | |

## Author inputs from help - VAL-057

Use your own input files and new keys; do not read the generator’s completed input. Each workspace is restored to its pre-action state. The acting Actor is supplied per row below. Discover domain identifiers from Human reads. For the third leaf, a concern.record workspace is prepared; choose it or another valid --input leaf available in that workspace and record your choice.

| Test ID | Covers | Test | Expected Result | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-05a | VAL-057 | Author run close with expected_repository. Workspace `/home/mark/projects/fid/.cache/phase-24-human/fixtures-2666841344/authoring-run-close`. Acting Actor `actor-recovery`. Start from `/home/mark/projects/fid/.cache/phase-24-human/fixtures-2666841344/fid run close --help` and --example. (human) | Accepted on your first or second try using only permitted guidance. | | Record chosen leaf, attempts, sources and final input. |
| HUMAN-05b | VAL-057 | Author decision record-reference with project-file digest omitted. Workspace `/home/mark/projects/fid/.cache/phase-24-human/fixtures-2666841344/authoring-decision-reference`. Acting Actor `actor-recovery`. Start from `/home/mark/projects/fid/.cache/phase-24-human/fixtures-2666841344/fid decision record-reference --help` and --example. (human) | Accepted on your first or second try using only permitted guidance. | | Record chosen leaf, attempts, sources and final input. |
| HUMAN-05c | VAL-057 | Choose and author a third --input leaf. Workspace `/home/mark/projects/fid/.cache/phase-24-human/fixtures-2666841344/authoring-third-leaf`. Acting Actor `actor-recovery`. Start from `/home/mark/projects/fid/.cache/phase-24-human/fixtures-2666841344/fid concern record --help` and --example. (human) | Accepted on your first or second try using only permitted guidance. | | Record chosen leaf, attempts, sources and final input. |

## Summary

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Take control | 13 | 0 | 0 | 0 | 13 |
| Failure clarity | 10 | 0 | 0 | 0 | 10 |
| Repository views | 16 | 0 | 0 | 0 | 16 |
| Judgment presentation | 101 | 0 | 0 | 0 | 101 |
| Input authoring | 3 | 0 | 0 | 0 | 3 |
| **Total** | 143 | 0 | 0 | 0 | 143 |

| | |
|---|---|
| **Overall result** | Pending Human |
| **Date completed** | |
| **Notes** | All cells intentionally blank; QA supplies feasibility, the Human supplies judgment. |
