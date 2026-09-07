# Phase 22 Final QA Preflight

Reviewer: independent QA Lead `checkpoint12_final`
Date: 2026-09-08
Status: PREPARATION ONLY — tasks 12–18 and final phase gates remain pending.

This is the early QA-84 reconciliation requested while task 12 completes. It
does not declare final QA, re-attest another session's execution, or replace a
product gate. Report tokens and termination fields belong to the originating
Developer/root sessions. This reviewer inspected report bytes and task specs;
no earlier product command was rerun for this preflight.

## Task 1–11 exact-command reconciliation

The exact-command sections contain 70 required occurrences. All are represented
by closing successful receipts with a nonempty token and observed-termination
field. Required repetitions are counted separately. Every task requires:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
```

| Task | Occurrences | Additional obligations inspected | Closing report receipts / result |
|---|---:|---|---|
| 1 | 12 | Environment six commands and 100k Event rebuild benchmark below | Lines 60–114; presence complete |
| 2 | 7 | Workflow package coverage before edits and after implementation | Baseline line 48; recovery baseline line 61; closing lines 120–131; PF-01 encoding finding |
| 3 | 5 | None | Closing QA-rework lines 90–94; PF-01 |
| 4 | 6 | Exact equivalent-Stop regression with count 5 below | Lines 67–73; PF-01 |
| 5 | 5 | None | Lines 59–63; presence complete |
| 6 | 5 | None | Lines 96–102; failed race receipt line 98 preserved; presence complete |
| 7 | 7 | Evidence package coverage before edits and after implementation | Original baseline line 57; recovery baselines 58/83/95; closing lines 105–112; PF-01 |
| 8 | 8 | Exact terminal test count 15 and two full benchmark invocations below | Closing lines 94–101; both distinct benchmark receipts present; PF-01 |
| 9 | 5 | None | Final-source lines 65–69; earlier quotations attributed; presence complete |
| 10 | 5 | None | Lines 53–57; presence complete |
| 11 | 5 | None | Lines 53–57; presence complete |

Line numbers refer to the inspected pre-correction report bytes. Exact extras:

```text
# task 1
go version
go env GOOS GOARCH
uname -srm
nproc
grep -m1 '^model name' /proc/cpuinfo
grep -m1 '^MemTotal:' /proc/meminfo
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal
# task 2: before edits, then after implementation
go test -race -cover ./internal/workflow
# task 4
go test -race ./internal/workflow -run '^TestEquivalentStopsShareAttentionButKeepDistinctFacts$' -count=5
# task 7: before edits, then after implementation
go test -race -cover ./internal/evidence
# task 8
go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=15 ./cmd/fid
go test -bench . -benchmem ./...
go test -bench . -benchmem ./...
```

No task declares a nonzero successful command. Task 8's eight historical
BLOCKED rows explicitly say the commands had not run during the cancelled
Receipt contradiction; they have no invented termination. Its later eight
PASS receipts close the gate after DEC-195.

Task 2 retains the 27.4% original workflow baseline, recovery baselines, failed
full gate, DEC-189 explanation and final 33.7% scoped result. Task 7 retains the
25.2% original evidence baseline and recovery baselines; its closing scoped
result is 24.0%. These are attributed historical measurements, not current-tree
measurements. Tasks 1,2,6,7,8,9,10,11 retain required-gate failure/deviation
history. Task 9 distinguishes originating Developer quotations from final
continuation receipts and root's coverage formatting correction. Tasks 10/11
retain report-control failures and root-assistance attribution.

## PF-01 — encoded tokens in five reports (original finding)

Classification: bounded report repair, no VAL / QA-84. This preparation finding
is not a final-review verdict. Root accepted the display-only repair ownership.

At observed HEAD `0699b660de9890bd769ac2dbceac21934f46d1e4`, these report
lines contained literal backslash-plus-t bytes in captured Go/coverage output.
The displayed tokens were encoded rather than verbatim. No missing execution
was inferred: commands, outcomes, values and termination fields were present.

| Report | Exact affected original lines |
|---|---|
| task-2-done.md | 48,54,57,61,63,69,71,73,75,77,88,93,95,110,114,115,117,118,120,126,128,130 |
| task-3-done.md | 56,60,66,69,72,75,76,79,85,86,87,92,93,94 |
| task-4-done.md | 49,54,59,62,70,71,72,73 |
| task-7-done.md | 57,58,63,66,67,70,71,72,76,79,80,81,82,83,86,91,94,95,100,101,105,110,111,112,127 |
| task-8-done.md | 82,96,97,98,99,100,101 |

The authorized repair decodes only known receipt display escapes, preserving
all original commands, outcomes, times, percentages, termination fields,
ordering, attribution and failure history. It records that this is rendering
correction, not new execution or a status upgrade. No product/test/gate change
is assigned. Resolution must separately cite the correction commit and byte
comparison against these originals.

### PF-01 resolution

RESOLVED by independent byte comparison of root's report-only commit
`be53edadaec38d12c7df6c0eaf43f1bf003f476d`
(`docs(phase-22): decode historical gate token rendering`). The changed-path
set is exactly the five named reports. For each, removing only the appended
`Root report display correction — 2026-09-08` paragraph leaves the original
bytes with literal backslash-t decoded to tab, and no other original-byte
change. Each appended paragraph explicitly preserves originating-session
attribution and disclaims new execution. The original finding remains above.

The Python byte-comparison command used `git diff-tree --no-commit-id
--name-only -r` to assert the exact five-file set, `git show <commit>^:<path>`
and `git show <commit>:<path>` for each version, and asserted:

```python
body, notes = after.split('\n## Root report display correction — 2026-09-08\n', 1)
assert body.rstrip('\n') == before.replace('\\t', '\t').rstrip('\n')
assert 'no new execution or attestation' in notes
```

Observed final token:
`PF-01 RESOLVED: five report-only changes exactly decode tabs and retain attribution; no command, outcome, value or termination changed`
(terminated `exit_code: 0`). The root's applicable report-control receipts are
separate Runner evidence; they have not been executed or re-attested here.

## Command and feasibility validation

Working directory: /home/mark/projects/fid. All referenced task/report files
exist and their exact-command sections contain concrete commands. No historical
binary or deleted test directory is needed to repeat this static inspection.

A Python stdlib inventory extracted fenced commands only from each exact-command
section, retained duplicates with Counter, and emitted each matching receipt's
line, status, token presence and termination-marker presence, plus report SHA-256
and latest report commit/author from git log. It terminated exit_code: 0.
All 70 occurrences had closing successful receipts; prefixed pre-edit baseline
rows were additionally inspected directly. This is file-state evidence, not
product PASS evidence.

The encoding inventory was:

```bash
python3 - <<'PY'
from pathlib import Path
for n in (2,3,4,7,8):
    p=Path(f'project-docs/docs/phases/phase-22/tasks/task-{n}-done.md')
    hits=[str(i) for i,line in enumerate(p.read_text().splitlines(),1)
          if '\\t' in line]
    print(f'{p}: literal-backslash-t lines '+','.join(hits))
PY
```

Observed final token:
`project-docs/docs/phases/phase-22/tasks/task-8-done.md: literal-backslash-t lines 82,96,97,98,99,100,101`
(terminated `exit_code: 0`). Original line lists are retained above; a
later correction does not erase the finding.

## Remaining work

Verified task-12 completion and its independent checkpoint take priority.
Tasks 12–18 still require reconciliation. Final QA must execute its complete
phase commands and deterministic checklist, verify review/repair closures, and
record each validation outcome independently. The separate human checklist has
six blank judgments; task 18 must supply, and QA must verify, its retained
fixture manifest before Human execution. No Human outcome is recorded here.
