# Phase 22 Human Verification Checklist

> Phase: 22 — Recovery and portability
> Prepared: 2026-09-12
> Tester: MBEL
> Status: **Attempt 1 completed: FAIL; preserved for task-19 repair and retest.**
> Human responses and outcomes below are the immutable first-attempt judgments.

## Before starting

Use a fresh session without prior chat or an explanation of the intended
answers. Inspect the supplied views and write your own response for each row
before recording an outcome. Record anything unclear in the notes, including
the wording or view that led to your interpretation.

The retained set below was independently inspected after its generator exited.
All twelve actual Human/JSON commands terminated successfully. Binary hash,
clean build revision and four Project identities were independently verified.
This set was generated late on clean final source after deterministic
qualification. The six Human questions retain their original scope; the Human
subsequently recorded the responses below. No agent judgment replaces them.

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
| HUMAN-01a (human) | VAL-076 | What single recovery action would you take next? State what in the view led you to that choice. | It looks like the main issue is that obligation OBL-p22-corpus is not satisfied and qualifying evidence is missing. As such, I'd run the gate next. | PARTIAL | I could identify gate.run as the likely next action only after working through the output with assistance, so I have marked this partial. The default resume view is difficult for a Human to scan: it is a wall of terse text with little vertical structure, many unexplained domain terms, lowercase statuses, and semicolon-separated fields. Lines such as completion incomplete do not clearly distinguish the field name from its value. The Claim line combines effective status, recorded status, and expiry on one line; headings with indented fields or bullet points would be easier to understand. actor-recovery looked like a technical recovery concept, but it is actually an arbitrary Actor ID. The same obligation-evidence-absent message is printed twice on consecutive lines because it applies both to the Obligation and to overall completion, but the text view removes that context and makes it look like accidental duplication. The JSON output formatted through jq was easier to understand than the intended Human view.|
| HUMAN-01b (human) | VAL-076 | What does the view tell you about the proof status of the adopted changes? Explain your interpretation. | I could not determine the proof status of the adopted changes from the view. It said that base.txt was an “adopted path” and “unproven work,” but I did not understand what “adopted” meant or what had happened to the file during recovery. I saw that qualifying evidence was absent and guessed that work may have been done but its quality gate was incomplete or interrupted. I only understood the intended distinction between preserving work and proving it after it was explained to me. | FAIL | I did not initially understand what adopted meant. I first interpreted the missing evidence as possibly meaning that a quality gate had been interrupted, but the view does not establish that. I needed an explanation that adoption means preserving existing uncommitted changes for the successor Run; it does not mean those changes were tested, reviewed, approved, or proved correct. The plain run inspect output added to the confusion because it printed only the Run ID and an unexplained SHA-256 reconciliation digest. A Human cannot infer changed paths, safety, proof status, liveness, or a next action from that digest. The JSON version contains the useful reconciliation details, but the checklist does not explain the purpose of each command or that HUMAN-01a/b primarily concern the resume view. A Human should not need to inspect JSON or receive prior explanation to understand this distinction, so I have marked this partial. |

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
| HUMAN-02a (human) | VAL-077 | In the pre-commit checkpoint, which recovery direction is offered, and what in the view explains it? | The view offers restore-prior. I could identify that label and see that the activation status was prepared, but the view did not explain what “restore prior” meant, what prior state would be restored, or what resources would be affected. I could not explain the recovery direction confidently without additional context. | PARTIAL | The three-column output was much easier to scan than the earlier resume output, but it was still too terse to explain the operation. restore-prior is an unexplained internal term. I needed external context to understand that Activation coordinates the kit, overlay, provider files, project configuration, and backend, and that restore-prior means returning those resources to their pre-activation identities before the point of no return. In this fixture the prior identities are null, which makes the phrase even less intuitive. |
| HUMAN-02b (human) | VAL-077 | In the post-commit checkpoint, is rollback possible? Explain what you understand from the view. |   No. The view says committed, complete-forward, and irreversible commit, which I understood to mean that rollback is no longer possible and recovery must proceed forward. | PASS |   The consequence was clear from the phrase irreversible commit: I understood that rollback was unavailable and that recovery had to continue forward. However, the output did not explain what had been committed, what complete-forward would do, or that this was an activation/install commit rather than a Git commit. Those details required additional context. |

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
| HUMAN-03a (human) | VAL-078 | Which restored Issue, Run and history are you looking at? Record the identities and what helped you recognize the work. |   I could identify Issue P22-2cd612 and Run RUN-b6d6a1 because those IDs were printed in the resume output. However, I could not tell from the Human view that the project had been imported, where it came from, or what history had been restored. I therefore could not confidently recognize the restored work beyond copying the visible Issue and Run IDs. | PARTIAL |   The Human output showed the current Issue, Run, Claim, and repository state, but it did not identify the project as imported or name its source. It also did not explain what “restored history” meant or show how the displayed identities related to the source project. The checklist supplied the source project and archive paths, but the command output itself did not connect them to the current view. The plain run inspect command again returned only a Run ID and an unexplained SHA-256 digest. The JSON contained more identity and repository information, but it was very long and used many fields whose meanings I did not know. I could identify some IDs but could not independently recognize the preserved history with confidence. |
| HUMAN-03b (human) | VAL-078 | What intervention, if any, remains before this imported work can be mutated? State what in the views led you to that conclusion. |   I could not determine the required intervention from the views. I initially answered that the Obligations needed to be completed because the output emphasized missing qualifying Evidence. I did not understand that this was a completion condition rather than the import-specific requirement for mutation. I only learned after receiving an explanation that imported active work requires fresh destination liveness and repository reconciliation before it can be mutated.| FAIL |   The Human view did not explain that importing an active Run does not restore mutation authority in the destination. It displayed the imported Claim as effective active, which suggested to me that the Claim might already be usable. It listed run.inspect and run.observe as ready actions but did not explain that these provide the fresh destination reconciliation and liveness information required after import. The view mixed those actions with missing Obligation evidence, completion status, budgets, and blocked Claim operations, so I incorrectly focused on the Obligation. The source and destination worktree identities, liveness: null, and the absence of a destination checkpoint were visible only in the lengthy JSON and were not meaningful to me without additional domain knowledge. I needed an external explanation to understand what intervention remained.|

## Human results

All six outcomes are recorded: one PASS, two FAIL and three PARTIAL. Agents may prepare and verify the fixtures;
only the Human records the judgments above. Deterministic QA evidence belongs
in the separate final review.

| Section | Total | Pass | Fail | Partial | Blank |
|---|---:|---:|---:|---:|---:|
| Recovery view | 2 | 0 | 1 | 1 | 0 |
| Activation recovery | 2 | 1 | 0 | 1 | 0 |
| Imported work | 2 | 0 | 1 | 1 | 0 |
| **Total** | **6** | **1** | **2** | **3** | **0** |

| Result record | Human entry |
|---|---|
| Overall result | FAIL |
| Date completed | 2026-09-12 |
| Notes | See inline notes above. |

Outcome values are `pass`, `fail` or `partial`; leave an unperformed item
blank. Record the overall result only after completing the six responses.
