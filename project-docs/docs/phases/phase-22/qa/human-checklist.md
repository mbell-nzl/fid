# Phase 22 Human Verification Checklist

> Phase: 22 — Recovery and portability
> Prepared: 2026-09-12
> Tester: _______________
> Status: **Runnable snapshot independently verified on final reviewed source; Human judgments pending.**
> Human responses and outcomes have not been recorded.

## Before starting

Use a fresh session without prior chat or an explanation of the intended
answers. Inspect the supplied views and write your own response for each row
before recording an outcome. Record anything unclear in the notes, including
the wording or view that led to your interpretation.

The retained set below was independently inspected after its generator exited.
All twelve actual Human/JSON commands terminated successfully. Binary hash,
clean build revision and four Project identities were independently verified.
This set was generated late on clean final source after deterministic
qualification. The six Human questions retain their original scope and blank
responses. No Human judgment is inferred.

| Preparation record | Verified value |
|---|---|
| Fixture guide | [human-fixtures.md](human-fixtures.md) |
| Manifest | [manifest.json](/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/manifest.json) |
| Binary | `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid` |
| SHA-256 | `sha256:f71a1051d0ab5de6dd17e3d7a2cdd222e14a6260973728f3f24764cecf48337b` |
| Clean source revision | `3832df1fbd271d9b868a7feb1d07f48b044efdbb`; source_dirty:false independently matched to go version -m |
| Generated | `2026-09-12T00:19:05.0971514Z` |
| Independent evidence | [final-review-4.md](final-review-4.md), Command And Feasibility Validation |
| Successor lease expires | `2026-09-13T00:19:09.826096583Z` / **13 September 2026, 12:19:09 NZST** |

Keep these directories at their recorded paths. If the lease expires before
this exercise or product change affects the behavior under review, generate a
new set with the accepted command below and have QA verify its manifest and
update the exact paths before recording acceptance. Preserve this original set.

```bash
cd /home/mark/projects/fid
FID_P22_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-22-human go test -tags integration ./cmd/fid/integration/... -run '^TestP22HumanFixtures$' -count=1 -v
```

## Recovery view

**Recovery checkpoint**

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/recovery`.

- `claim_id`: `CLM-2b6d10`
- `issue_id`: `P22-071060`
- `prior_run_id`: `RUN-79a072`
- `project_id`: `project-e329279801dcf616b56fc8a1ddd55e30`
- `run_id`: `RUN-d433a1`

Copy these inspection commands as a block:

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/recovery resume P22-071060 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/recovery --json resume P22-071060 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/recovery run inspect RUN-d433a1 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/recovery
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/recovery --json run inspect RUN-d433a1 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/recovery
```


| Test ID | Covers | Prompt | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-01a (human) | VAL-076 | What single recovery action would you take next? State what in the view led you to that choice. | | | |
| HUMAN-01b (human) | VAL-076 | What does the view tell you about the proof status of the adopted changes? Explain your interpretation. | | | |

## Activation recovery

**Activation checkpoint A**

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/prepared`.

- `activation_id`: `ACT-30a259`
- `project_id`: `project-4e345f3f38d949d1c0d14543e265543e`

Copy these inspection commands as a block:

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/prepared
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/prepared activation status ACT-30a259
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/prepared --json activation status ACT-30a259
```

**Activation checkpoint B**

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/committed`.

- `activation_id`: `ACT-1b512b`
- `project_id`: `project-38b61ea3997ac2d506587789e408bd43`

Copy these inspection commands as a block:

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/committed
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/committed activation status ACT-1b512b
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/committed --json activation status ACT-1b512b
```


| Test ID | Covers | Prompt | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-02a (human) | VAL-077 | In the pre-commit checkpoint, which recovery direction is offered, and what in the view explains it? | | | |
| HUMAN-02b (human) | VAL-077 | In the post-commit checkpoint, is rollback possible? Explain what you understand from the view. | | | |

## Imported work

**Imported project checkpoint**

Project: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/imported`.

- `archive_path`: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/import.zip`
- `claim_id`: `CLM-20648d`
- `issue_id`: `P22-2cd612`
- `project_id`: `project-bcd42e0100d0b70a6f3ece0cb2ab00ee`
- `run_id`: `RUN-b6d6a1`
- `source_project_path`: `/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/import-source`

Copy these inspection commands as a block:

```bash
cd /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/imported resume P22-2cd612 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/imported --json resume P22-2cd612 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/imported run inspect RUN-b6d6a1 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/imported
/home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/fid --project /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/imported --json run inspect RUN-b6d6a1 --repository /home/mark/projects/fid/.cache/phase-22-human/checkpoints-4132086272/imported
```


| Test ID | Covers | Prompt | Your response / evidence | Outcome | Notes |
|---|---|---|---|---|---|
| HUMAN-03a (human) | VAL-078 | Which restored Issue, Run and history are you looking at? Record the identities and what helped you recognize the work. | | | |
| HUMAN-03b (human) | VAL-078 | What intervention, if any, remains before this imported work can be mutated? State what in the views led you to that conclusion. | | | |

## Human results

All six outcomes remain blank. Agents may prepare and verify the fixtures;
only the Human records the judgments above. Deterministic QA evidence belongs
in the separate final review.

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Recovery view | 2 | 0 | 0 | 0 | 2 |
| Activation recovery | 2 | 0 | 0 | 0 | 2 |
| Imported work | 2 | 0 | 0 | 0 | 2 |
| **Total** | **6** | **0** | **0** | **0** | **6** |

| Result record | Human entry |
|---|---|
| Overall result | |
| Date completed | |
| Notes | |

Outcome values are `pass`, `fail` or `partial`; leave an unperformed item
blank. Record the overall result only after completing the six responses.
