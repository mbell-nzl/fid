Phase 23 delivers Review and Human judgment. An independent Reviewer records coverage and Findings against an exact Envelope revision and repository identity. Builder rework keeps review independence; only a later qualifying Review resolves a Finding. Protected Decisions and subjective Obligations receive provenanced Human Outcomes, and completion is re-evaluated from those facts.

Scope:

- Review Records, Findings and their dispositions (addressed is nonterminal; resolved and inapplicable commit with the qualifying Review).
- Review-cycle budgets. A `review-rework-exhausted` Stop is resolved by an approved `review-rework-extension` Human Outcome.
- Concerns, Deviations, and Decision References, with protected-Decision blocking.
- Human Outcomes and Attention resolution keyed on the recorded Stop Condition.
- Requirement references with a stable `requirement_id`. These replace V1 requirement authoring.
- Migrations 0017 (link-preserving Evidence rebuild) and 0018 (V1 judgment removal, exact schema-17 down).
- Portable recovery parity for every Phase 23 table.
- `reviews.coverage.v1` and `human.outcomes.v1` qualified.
- Package boundary: `internal/review` added, `internal/requirements` removed, `internal/links` restricted to Issues.

Validation: 18 planned tasks plus supplemental tasks 19-21; all 89 validation assertions pass (85 deterministic, 4 Human).

- The first Human checklist failed (55 pass, 2 fail, 6 partial). It found an unresolvable rework Stop, silent Human-mode Stop failures, and a missing stable requirement ID. Tasks 19-21 repaired them, and the Human re-check passed 7/7.
- Independent final review 2 passes with no Must-fix.
- VAL-071 and VAL-073 pass by Tech Lead disposition. Their residual Human notes are carried forward as CF-045.

Recorded checks (final review 2):

- Format, lint, and the race/coverage suite (60.7%).
- Tagged real-process integration and the phase-end full-`checkptr` control.
- Benchmarks, build, and a 100k rebuild at 23705 ms (limit 30000 ms).
- Uncached package times, inside the 600/900 s budgets: `cmd/fid` 359.6 s, `internal/workflow` 113.6 s, tagged integration 289.6 s.

Hosted CI does not run the tagged integration suite (CF-041), so a green PR is not phase-level evidence.

Acceptance and review trail:

- `project-docs/docs/phases/phase-23/qa/final-acceptance.md`
- `project-docs/docs/phases/phase-23/qa/final-review-2.md`
- `project-docs/docs/phases/phase-23/qa/human-checklist.md` and `human-checklist-2.md`
- `project-docs/docs/phases/phase-23/wrap-up-audit.md`
- `project-docs/docs/phases/phase-23/decisions.md` (DEC-271 sign-off)

Remaining boundaries:

- Human operating experience (raw-JSON Stop details, denial reasons, hint verbosity, `--input` authoring) is CF-045, to be scheduled before Phase 24 planning.
- Pre-change Phase 23 development journals fail closed (DEC-270, ADR-0011).
- Parallel Claims and Envelope revision remain Phase 24; adapter and all-capability qualification remain Phase 25.
- This is a pre-adoption baseline, not a supported release.

🤖 Generated with [Claude Code](https://claude.com/claude-code)
