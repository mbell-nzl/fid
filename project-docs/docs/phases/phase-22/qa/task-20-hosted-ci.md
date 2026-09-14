# Task 20 — Hosted CI qualification

Verified 2026-09-14 through GitHub CLI, with observed exit_code 0 for the PR
query, run metadata query and full-log retrieval. Raw captures:
.cache/phase22-pr-performance/hosted-after.json and hosted-after.log.

- [Successful run 34805256766](https://github.com/mbell-nzl/fid/actions/runs/34805256766)
- PR24 head: 6b8fe50bcd47148572ba841d561b23348f041dad.
- Run conclusion: success. All six jobs concluded success.
- Run created/started04:13:51Z, updated after completion04:27:44Z:13m53s.
- First job start to final job finish:13m49s. Timing definitions differ by four seconds of orchestration.

| Job | Execution duration |
| --- | ---: |
| checks | 3m42s |
| tests (0) | 13m04s |
| tests (1) | 11m53s |
| tests (2) | 7m52s |
| tests (3) | 10m36s |
| quality | 41s |
| Sum | 47m48s |

Per-job upward rounding estimates50 runner minutes. This is derived from job
metadata, not an account billing receipt; no plan, remaining allowance or charge
has been queried. The old failed quality job used34m35s (approximately35 rounded
minutes), so this run used38.2% more summed execution time while completing
about59.9% sooner. The old run skipped its second full-suite coverage pass and
build; these numbers do not establish cost versus a successful old workflow.

## Hosted evidence

- Every worker: `181 selected targets completed; all 21 packages accounted for`.
- Final reconciliation: `PASS coverage union: 4 shards, 724 targets, 12967 statement blocks`.
- Combined coverage: `65.5%`.
- checks: format, lint, eight controller tests and build all success.
- Artifact upload/download and final required quality check all success.
- Unchanged default package timeout: the longest reported cmd/fid package was393.410s; a job can exceed ten minutes because it executes several packages and discovery.

The four equal-count workers took7m52s–13m04s. Discovery/build before selected
execution accounted for roughly1m27s–1m51s per worker, measured from the test
step's start to its shard announcement. Equal counts do not balance duration;
future optimization could use measured timings. No additional tuning is made.

## Closure and boundary

Task20 is COMPLETE. The original hosted failure is resolved by a successful
hosted run of the tested implementation; failed receipts remain intact. All
original tests and accepted product/Human judgments remain. This addendum and
status updates are documentation only and do not claim another hosted run of
the later documentation commit. The verified PR head above remains the binding
CI evidence. Root does not push these records or merge the PR; Human owns push.
