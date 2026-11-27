Review the task's whole change in this detached review worktree, at the
hand-off head supplied below. Read AGENTS.md, the approved brief, and any
repository reviewer instructions. Use handoff.base...handoff.head for the
diff. Review every required_profiles entry in one fresh session, for the
declared stream. A separate stream is an independent whole-change review.
Actors are distinct; the default same-model review uses Mark's standing
fallback when the other model is unavailable. Name the actual model.

This hand-off head's required checks have passed on its tree through fid check;
fid refuses hand-off otherwise. Do not rerun full suites (including make
fmt-check, lint, test, build or the integration suite), even if repository
instructions tell builders to run them. Read the change, run focused tests
and probes with overlay mutations, and record your commands and results.

Assess correctness, scope and test adequacy. Refuse only on missing evidence;
size and coverage numbers inform judgement. For each open earlier finding,
assess whether the change resolves it, with evidence. Do not edit or commit
files, post a comment, push, merge, approve as a person or record the review.
fid records your result and posts its canonical PR comment.

fid posts these fields to the PR as GitHub Markdown for a person to read: in
every text field, put commands, file paths, test names, flags, identifiers and
literal output in `backticks`. Write verification as a short Markdown list,
one command per item with its result. Use concise evidence.

Write only the review input JSON to LOOP_REVIEW, with this schema:
{"reviewer":{"actor":"<LOOP_ACTOR>","model":"<LOOP_MODEL>",
"session_id":"<LOOP_SESSION>"},"stream":"<requested stream>",
"profiles":["<every required profile>"],"verdict":"APPROVE or CHANGES",
"purpose":"<what the task does>","verification":"<commands, results, limits>",
"findings":[{"profile":"code or test-adequacy","severity":"major, minor or nit",
"summary":"<problem>","location":{"file":"<relative path>","line":1},
"evidence":"<observable evidence>","fix":"<requested action>"}],
"dispositions":[{"finding_id":"<earlier ID>","disposition":"resolved, not-applicable or open",
"reason":"<evidence>"}]}
Use findings: [] when there are none; location may be null. CHANGES means a
major remains open; minor means fix or answer with a reason; nit is optional.
Close an earlier major through dispositions after verifying its resolution.
The context JSON, including earlier findings, is also at LOOP_CONTEXT.
