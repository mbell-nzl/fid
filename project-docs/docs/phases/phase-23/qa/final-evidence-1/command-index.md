# Phase 23 final re-review — command receipts

Every status below is the nested command’s observed termination. Child processes inside Python observations also retain their individual statuses in linked receipts. These are QA observations, not borrowed Developer gate results. Observation duration includes polling delay; Go package timings are in the raw output. Setup, help, and file inventory commands are identified by their names and are not product gates.

| ID | Observed exit_code | Last output token | Observation seconds |
|---|---:|---|---:|
| [all-action-provenance](all-action-provenance.json) | 0 | `PASS provenance recorded for all 30 grant actions; all seven Phase 23 introductions match owning tasks` | 17.374 |
| [archive-scale](archive-scale.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/export	147.666s` | 172.72 |
| [artifact-audit-final](artifact-audit-final.json) | 0 | `PASS 89 VAL outcomes; 52 deterministic rows; exactly seven blank Human rows; all local evidence links resolve; observed terminations retained; protected hashes/source unchanged` | 0.179 |
| [artifact-audit](artifact-audit.json) | 0 | `PASS 89 VAL outcomes; 52 deterministic rows; exactly seven blank Human rows; all local evidence links resolve; observed terminations retained; protected hashes/source unchanged` | 0.192 |
| [artifact-diff-check](artifact-diff-check.json) | 0 | `exit 0 (no output)` | 0.049 |
| [authority-producer-inventory](authority-producer-inventory.json) | 0 | `cmd/fid/finding.go:126:		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, "active envelope authority is denied", nil, map[string]any{"operation": operation, "actor_id": denied.ActorID, "required_actions": []string{denied.Action}, "required_capabilities": []string{denied.Capability}})` | 0.062 |
| [benchmark-1](benchmark-1.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/journal	26.912s` | 70.593 |
| [benchmark-2](benchmark-2.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/journal	26.312s` | 34.414 |
| [benchmark-3](benchmark-3.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/journal	25.481s` | 25.826 |
| [budget-integration](budget-integration.json) | 0 | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	267.627s` | 291.929 |
| [budget-ordinary](budget-ordinary.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/workflow	101.659s	coverage: 62.7% of statements` | 379.746 |
| [capability-staging](capability-staging.json) | 0 | `PASS task-11-close false and final true capability states on fresh initialized workspaces` | 0.442 |
| [documentation-amendments](documentation-amendments.json) | 0 | `+**Failure code** and **Evaluation reason**)*` | 0.1 |
| [feasibility](feasibility.json) | 0 | `Help observations recorded; help-only nonzero statuses are syntax evidence, not product gate failures` | 0.779 |
| [fixture-identity](fixture-identity.json) | 0 | `PASS fresh manifest and binary digest verified; binary revisions match reviewed HEAD; product diff empty` | 0.233 |
| [full-checkptr](full-checkptr.json) | 1 | `FAIL` | 676.392 |
| [gate-01-format](gate-01-format.json) | 0 | `exit 0 (no output)` | 0.286 |
| [gate-02-lint](gate-02-lint.json) | 0 | `0 issues.` | 21.99 |
| [gate-03-coverage](gate-03-coverage.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 62.7% of statements` | 72.374 |
| [gate-04-coverage-summary](gate-04-coverage-summary.json) | 0 | `total:											(statements)					60.7%` | 0.381 |
| [gate-05-integration](gate-05-integration.json) | 0 | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	250.442s` | 303.607 |
| [gate-06-benchmarks](gate-06-benchmarks.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/workflow	11.876s` | 279.804 |
| [gate-07-build](gate-07-build.json) | 0 | `exit 0 (no output)` | 1.021 |
| [human-fixture-generation](human-fixture-generation.json) | 0 | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	29.199s` | 42.617 |
| [human-guide-verification](human-guide-verification.json) | 0 | `PASS 63 exact guide invocations; Human/JSON facts match; seven Human judgments remain blank; fresh fixtures restored for handoff` | 19.444 |
| [import-failure-diagnosis](import-failure-diagnosis.json) | 0 | `Diagnosis: full-checkptr failed at release of an imported fixture Claim; source/recovered fixture Claims use 60-second leases. Task 16 owns matching import-binding tests. No product change or retry performed.` | 0.543 |
| [import-failure-provenance](import-failure-provenance.json) | 0 | `PASS introducing test and 60-second fixture leases predate Phase 23 planning base; provenance does not waive the failed gate` | 0.158 |
| [live-bootstrap](live-bootstrap.json) | 0 | `ok  	github.com/mbell-nzl/fid/cmd/fid	6.303s` | 39.731 |
| [live-golden](live-golden.json) | 0 | `ok  	github.com/mbell-nzl/fid/cmd/fid	4.538s` | 21.906 |
| [named-audit-retry](named-audit-retry.json) | 0 | `PASS 188 fresh selected tests; no missing/skipped selection or truncated inner result; all corruption receipts verified` | 0.103 |
| [named-audit](named-audit.json) | 1 | `FileNotFoundError: [Errno 2] No such file or directory: '/home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence-1/named-integration.json'` | 0.207 |
| [named-claims](named-claims.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/claims	1.036s` | 1.496 |
| [named-cli](named-cli.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/cli	1.137s` | 1.405 |
| [named-db](named-db.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/db	10.809s` | 11.349 |
| [named-evidence](named-evidence.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/evidence	1.055s` | 61.124 |
| [named-export](named-export.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/export	10.420s` | 10.978 |
| [named-feasibility](named-feasibility.json) | 0 | `PASS 188 named selections have concrete source scope; execution proof follows in named receipts` | 0.18 |
| [named-fid-integration](named-fid-integration.json) | 0 | `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	77.597s` | 294.719 |
| [named-fid](named-fid.json) | 0 | `ok  	github.com/mbell-nzl/fid/cmd/fid	43.372s` | 54.409 |
| [named-issues](named-issues.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/issues	11.079s` | 11.645 |
| [named-journal](named-journal.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/journal	20.810s` | 21.313 |
| [named-links](named-links.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/links	1.018s` | 1.499 |
| [named-maintenance](named-maintenance.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/maintenance	2.891s` | 11.617 |
| [named-review](named-review.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/review	1.153s` | 1.679 |
| [named-workflow](named-workflow.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/workflow	19.621s` | 20.162 |
| [prechange-rebuild](prechange-rebuild.json) | 0 | `PASS DEC-270 expected exit 6 refusal; source journal SHA-256 unchanged` | 0.323 |
| [preservation-rows](preservation-rows.json) | 0 | `PASS linked-seed row preservation; exactly two added columns are null; foreign_key_check=[]; integrity_check=ok` | 41.308 |
| [rebuild-pragmas](rebuild-pragmas.json) | 0 | `PASS full journal rebuild: foreign_key_check=[]; integrity_check=ok; source journal unchanged` | 0.281 |
| [receipt-index](receipt-index.json) | 0 | `PASS 79 receipts indexed; all commands have observed nested exit_code and complete captured output` | 0.119 |
| [reference-binary-identity](reference-binary-identity.json) | 0 | `	build	vcs.modified=false` | 0.091 |
| [reference-comparison](reference-comparison.json) | 0 | `PASS six canonical bytes and Attention key identical at 74c83b7 and 32a55a9; pinned journal identity verified` | 0.108 |
| [reference-current-probe](reference-current-probe.json) | 0 | `ok  	github.com/mbell-nzl/fid/internal/workflow	3.972s` | 13.254 |
| [reference-prechange-probe](reference-prechange-probe.json) | 0 | `{"command": ["go", "test", "-race", "-count=1", "-overlay", "/home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence/reference-overlay.json", "-run", "^TestQAPhase23ReferenceAttentionKeyAndBytes$", "-v", "./internal/workflow"], "cwd": "/tmp/fid-p23-qa-reference-PgIz9Jst/prechange", "environment": {"GIT_DIR": "/home/mark/projects/fid/.git/worktrees/prechange", "GIT_WORK_TREE": "/tmp/fid-p23-qa-reference-PgIz9Jst/prechange"}, "stdout": "=== RUN   TestQAPhase23ReferenceAttentionKeyAndBytes\n=== PAUSE TestQAPhase23ReferenceAttentionKeyAndBytes\n=== CONT  TestQAPhase23ReferenceAttentionKeyAndBytes\n    qa_reference_test.go:19: canonical executor-terminated: {\"type\":\"executor-terminated\",\"subject\":{\"type\":\"run\",\"id\":\"RUN-pinned\"},\"budget_name\":null,\"capability_id\":null}\n    qa_reference_test.go:19: canonical run-recovered: {\"type\":\"run-recovered\",\"subject\":{\"type\":\"run\",\"id\":\"RUN-pinned\"},\"budget_name\":null,\"capability_id\":null}\n    qa_reference_test.go:19: canonical activation-recovered: {\"type\":\"activation-recovered\",\"subject\":{\"type\":\"activation\",\"id\":\"ACT-pinned\"},\"budget_name\":null,\"capability_id\":null}\n    qa_reference_test.go:19: canonical capability-available: {\"type\":\"capability-available\",\"subject\":{\"type\":\"activation\",\"id\":\"ACT-pinned\"},\"budget_name\":null,\"capability_id\":\"attention.v1\"}\n    qa_reference_test.go:19: canonical budget-observed: {\"type\":\"budget-observed\",\"subject\":{\"type\":\"run\",\"id\":\"RUN-pinned\"},\"budget_name\":\"attempts\",\"capability_id\":null}\n    qa_reference_test.go:19: canonical work-ended: {\"type\":\"work-ended\",\"subject\":{\"type\":\"run\",\"id\":\"RUN-pinned\"},\"budget_name\":null,\"capability_id\":null}\n    qa_reference_test.go:25: subject={type:run,id:RUN-pinned} code=budget-exhausted condition={type:budget-observed,subject:{type:run,id:RUN-pinned},budget_name:attempts,capability_id:null} key=sha256:d4b51bb4a6404d4a2579826ad5c08d802ed83c41814cd5fa74d3ee8c7796c3de\n--- PASS: TestQAPhase23ReferenceAttentionKeyAndBytes (0.00s)\nPASS\nok  \tgithub.com/mbell-nzl/fid/internal/workflow\t3.380s\n", "stderr": "", "exit_code": 0}` | 6.444 |
| [reference-task11-build](reference-task11-build.json) | 0 | `{"command": ["go", "build", "-o", "/tmp/fid-p23-final-review-1-task11", "./cmd/fid"], "cwd": "/tmp/fid-p23-qa-reference-PgIz9Jst/task11", "environment": {"GIT_DIR": "/home/mark/projects/fid/.git/worktrees/task111", "GIT_WORK_TREE": "/tmp/fid-p23-qa-reference-PgIz9Jst/task11"}, "stdout": "", "stderr": "", "exit_code": 0}` | 0.763 |
| [report-01](report-01.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-02](report-02.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-03](report-03.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-04](report-04.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-05](report-05.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-06](report-06.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-07](report-07.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-08](report-08.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-09](report-09.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-10](report-10.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-11](report-11.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-12](report-12.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-13](report-13.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-14](report-14.json) | 0 | `PASS report-lint completed with no failures` | not measured |
| [report-15](report-15.json) | 0 | `PASS report-lint completed with no failures` | 22.586 |
| [report-16](report-16.json) | 0 | `PASS report-lint completed with no failures` | 23.17 |
| [report-17](report-17.json) | 0 | `PASS report-lint completed with no failures` | 22.42 |
| [report-18](report-18.json) | 0 | `PASS report-lint completed with no failures` | 23.303 |
| [report-19](report-19.json) | 0 | `PASS report-lint completed with no failures` | 23.427 |
| [report-20](report-20.json) | 0 | `PASS report-lint completed with no failures` | 23.781 |
| [report-21](report-21.json) | 0 | `PASS report-lint completed with no failures` | 23.521 |
| [schema-sql](schema-sql.json) | 0 | `PASS independently applied schema-17 and schema-18 pragma and inventory probes` | 41.036 |
| [security-source-audit](security-source-audit.json) | 0 | `Security source observations captured; behavioral proof is supplied by named tests and real-process journeys` | 0.174 |
| [setup-testcache](setup-testcache.json) | 0 | `exit 0 (no output)` | 0.094 |
| [source-audit](source-audit.json) | 0 | `PASS source, report command, repair scope and provenance audit` | 40.61 |
| [tmp-after-budget-integration](tmp-after-budget-integration.json) | 0 | `[]` | 0.152 |
| [tmp-after-integration](tmp-after-integration.json) | 0 | `[]` | 1.158 |
| [tmp-before-budget-integration](tmp-before-budget-integration.json) | 0 | `[]` | 0.148 |
| [tmp-before-integration](tmp-before-integration.json) | 0 | `[]` | not measured |

## Exact commands and working directories

### all-action-provenance

Working directory: `/home/mark/projects/fid`.

```text
python3 - <<'PY'
import pathlib,json,subprocess
rows=json.loads(pathlib.Path('cmd/fid/integration/testdata/p23_conformance/grant-actions.json').read_text())
assert len(rows)==30
for row in rows:
 args=['git','log','-S"'+row['name']+'"','--format=%H %s%n%b','--','internal/workflow/activation.go']
 r=subprocess.run(args,text=True,capture_output=True)
 print(json.dumps({'owner':row,'command':args,'stdout':r.stdout,'stderr':r.stderr,'exit_code':r.returncode}))
 assert r.returncode==0 and r.stdout.strip()
 if row['phase']==23: assert 'phase-23/task-'+str(row['task']) in r.stdout
print('PASS provenance recorded for all 30 grant actions; all seven Phase 23 introductions match owning tasks')
PY
```

### archive-scale

Working directory: `/home/mark/projects/fid`.

```text
go test -tags archive_scale -race ./internal/export -run '^TestArchive(ExpandedByteBoundaries|CumulativeBoundary|EntryCountBoundaries)$' -count=1
```

### artifact-audit-final

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/audit-final-artifacts.py
```

### artifact-audit

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/audit-final-artifacts.py
```

### artifact-diff-check

Working directory: `/home/mark/projects/fid`.

```text
git diff --check -- project-docs/docs/phases/phase-23/qa
```

### authority-producer-inventory

Working directory: `/home/mark/projects/fid`.

```text
rg -n --glob '!*_test.go' 'CodeAuthority(Missing|UnknownActor|Denied)' cmd internal
```

### benchmark-1

Working directory: `/home/mark/projects/fid`.

```text
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
```

### benchmark-2

Working directory: `/home/mark/projects/fid`.

```text
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
```

### benchmark-3

Working directory: `/home/mark/projects/fid`.

```text
go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/
```

### budget-integration

Working directory: `/home/mark/projects/fid`.

```text
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

### budget-ordinary

Working directory: `/home/mark/projects/fid`.

```text
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
```

### capability-staging

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/capability-staging.py
```

### documentation-amendments

Working directory: `/home/mark/projects/fid`.

```text
git diff 0195897 -- project-docs/specs/02_TRD.md project-docs/specs/03_TRD_SCHEMA.md project-docs/specs/06_DESIGN_PLAN_CLI.md project-docs/specs/GLOSSARY.md
```

### feasibility

Working directory: `/home/mark/projects/fid`.

```text
python3 - <<'PY'
import subprocess, json
from pathlib import Path
commands=[['go','version'],['go','env','GOOS','GOARCH'],['uname','-a'],['nproc'],['go','help','test'],['go','help','testflag'],['go','help','build'],['go','help','clean'],['go','tool','cover','-help'],['gofmt','-h'],['golangci-lint','run','--help'],['sqlite3','--version'],['/tmp/fid','--help'],['/tmp/fid','rebuild','--help'],['/tmp/fid','human-outcome','record','--help'],['/tmp/fid','claim','acquire','--help'],['project-docs/tools/report-lint.sh','--help']]
rows=[]
for args in commands:
 r=subprocess.run(args,text=True,capture_output=True)
 rows.append(dict(command=args,stdout=r.stdout,stderr=r.stderr,exit_code=r.returncode))
 print(args, 'exit_code:',r.returncode, (r.stdout+r.stderr).splitlines()[:2])
Path('project-docs/docs/phases/phase-23/qa/final-evidence-1/feasibility-receipts.json').write_text(json.dumps(rows,indent=2)+'\n')
print(next(x for x in Path('/proc/cpuinfo').read_text().splitlines() if x.startswith('model name')))
print(next(x for x in Path('/proc/meminfo').read_text().splitlines() if x.startswith('MemTotal:')))
print('Help observations recorded; help-only nonzero statuses are syntax evidence, not product gate failures')
PY
```

### fixture-identity

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/fixture-identity.py
```

### full-checkptr

Working directory: `/home/mark/projects/fid`.

```text
go test -race -timeout 30m ./...
```

### gate-01-format

Working directory: `/home/mark/projects/fid`.

```text
gofmt -l cmd internal
```

### gate-02-lint

Working directory: `/home/mark/projects/fid`.

```text
golangci-lint run ./...
```

### gate-03-coverage

Working directory: `/home/mark/projects/fid`.

```text
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
```

### gate-04-coverage-summary

Working directory: `/home/mark/projects/fid`.

```text
go tool cover -func=coverage.out
```

### gate-05-integration

Working directory: `/home/mark/projects/fid`.

```text
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

### gate-06-benchmarks

Working directory: `/home/mark/projects/fid`.

```text
go test -bench . -benchmem ./...
```

### gate-07-build

Working directory: `/home/mark/projects/fid`.

```text
go build -o /tmp/fid ./cmd/fid
```

### human-fixture-generation

Working directory: `/home/mark/projects/fid`.

```text
FID_P23_HUMAN_FIXTURE_ROOT=/home/mark/projects/fid/.cache/phase-23-human go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -run '^TestP23HumanFixtures$' -count=1 -v ./cmd/fid/integration
```

### human-guide-verification

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/refresh-human-guide.py /home/mark/projects/fid/.cache/phase-23-human/checkpoints-935473561/manifest.json
```

### import-failure-diagnosis

Working directory: `/home/mark/projects/fid`.

```text
python3 - <<'PY'
from pathlib import Path
import json,subprocess
receipts=[]
commands=[['git','log','-S','func TestImportRequiresFreshDestinationObservationBeforeRecovery','--format=%H %s%n%b','--','internal/export/import_test.go'],['git','log','-L','81,81:internal/export/archive_fixture_test.go','--format=%H %s%n%b'],['git','log','-L','134,134:internal/export/archive_fixture_test.go','--format=%H %s%n%b'],['git','log','--diff-filter=A','--format=%H %s%n%b','--','internal/export/archive_fixture_test.go']]
for args in commands:
 r=subprocess.run(args,text=True,capture_output=True)
 row=dict(command=args,stdout=r.stdout,stderr=r.stderr,exit_code=r.returncode);receipts.append(row)
 print(json.dumps(row));assert r.returncode==0
for path,lo,hi in [('internal/export/import_test.go',182,290),('internal/export/archive_fixture_test.go',55,140),('project-docs/docs/phases/phase-23/tasks/task-16.md',24,60)]:
 lines=Path(path).read_text().splitlines()
 rows='\n'.join(f'{i+1}: {lines[i]}' for i in range(lo-1,min(hi,len(lines))))
 receipts.append(dict(path=path,lines=rows));print(path+'\n'+rows)
Path('project-docs/docs/phases/phase-23/qa/final-evidence-1/import-failure-diagnosis-receipts.json').write_text(json.dumps(receipts,indent=2)+'\n')
print('Diagnosis: full-checkptr failed at release of an imported fixture Claim; source/recovered fixture Claims use 60-second leases. Task 16 owns matching import-binding tests. No product change or retry performed.')
PY
```

### import-failure-provenance

Working directory: `/home/mark/projects/fid`.

```text
python3 - <<'PY'
import subprocess,json
for commit in ['cc73b092bf5ecd3fdba62a1362062a6d543697ab','1be8e882ef689f15705507c448e3b087e23d5f65','c9b46655d880c8b623be08d06618bffabfa5333a']:
 args=['git','merge-base','--is-ancestor',commit,'0195897a347dfe1872db0f7be3800c5c3ccdefcf']
 r=subprocess.run(args,text=True,capture_output=True);print(json.dumps(dict(command=args,stdout=r.stdout,stderr=r.stderr,exit_code=r.returncode)));assert r.returncode==0
print('PASS introducing test and 60-second fixture leases predate Phase 23 planning base; provenance does not waive the failed gate')
PY
```

### live-bootstrap

Working directory: `/home/mark/projects/fid`.

```text
go test -race -run '^TestLiveGoldenChildSkipsPristineBootstrap$' -count=1 -v ./cmd/fid
```

### live-golden

Working directory: `/home/mark/projects/fid`.

```text
go test -race -run '^TestLiveOutputGoldenFixtures$' -count=1 -v ./cmd/fid
```

### named-audit-retry

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/audit-named.py
```

### named-audit

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/audit-named.py
```

### named-claims

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestClaimAcquireRejectsOpenReviewReworkStop)$' ./internal/claims
```

### named-cli

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestReasonCodeDetailSchemaIsExactForEveryCode|TestPhase23FailureCodesAreRegistered)$' ./internal/cli
```

### named-db

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestMigration0017CreatesExactlyTheAcceptedObjects|TestMigration0017ObligationRegistrationDiscriminator|TestPhase23BaselineMigrationDiffIsExactlySevenAcceptedHunks|TestMigration0017PreservesReceiptEvidenceAndArtifactLinksInPlace|TestMigration0018DeletesOnlyLegacyNonIssueReferences|TestMigration0018CurrentReferenceGuards|TestMigration0018IssuesEnrichedColumnsAndSubjectCount|TestLatestSchemaHasNoV1JudgmentTables|TestMigration0018DownRestoresEmptyV1Tables|TestMigration0018DownMatchesFreshSchema17|TestMigration0018DownPreservesAgentRunAndFullRollback|TestMigration0018RoundTripIntegrity)$' ./internal/db
```

### named-evidence

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestHumanOutcomeAdmissionSatisfiesOnlyMatchingClassWithApprovedResult|TestReviewEvidenceStalesOnTargetTreeChange|TestHumanOutcomeEvaluationKeepsFreshnessAndUsesSupersession|TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant|TestEvidenceAdmissionTask6PureProductionClassifierMatrix|TestReviewAdmissionReturnsBlockingFindingIDs|TestLoadReviewEvidenceBindingEnforcesVerdictTypeAndLinkTimeIndependence|TestLoadHumanOutcomeEvidenceBindingRejectsClassMismatchAsOutcomeNotPermitted|TestReviewCyclesDoNotBlockGateAdmission)$' ./internal/evidence
```

### named-export

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestCreateIsByteIdenticalWithCompleteDeterministicInventory)$' ./internal/export
```

### named-feasibility

Working directory: `/home/mark/projects/fid`.

```text
python3 - <<'PY'
from pathlib import Path
import json,re
root=Path.cwd();out=root/'project-docs/docs/phases/phase-23/qa/final-evidence-1'
inventory=json.loads((out/'named-tests.json').read_text());mapping={}
for pkg,names in inventory.items():
 found={}
 for p in (root/pkg).glob('*_test.go'):
  for num,line in enumerate(p.read_text().splitlines(),1):
   m=re.match(r'func (Test\w+)\(',line)
   if m:found[m[1]]=str(p.relative_to(root))+':'+str(num)
 for name in names:
  assert name in found,(pkg,name)
  mapping[name]=found[name]
 print(pkg,len(names),'named symbols exist')
(out/'named-test-source-map.json').write_text(json.dumps(mapping,indent=2)+'\n')
print('PASS',len(mapping),'named selections have concrete source scope; execution proof follows in named receipts')
PY
```

### named-fid-integration

Working directory: `/home/mark/projects/fid`.

```text
go test -tags integration -race -count=1 -v -run '^(TestP23DecisionReferenceProtectedLifecycleRealProcess|TestP23FailurePrecedenceWithProtectedDecisionAndObligation|TestP23FailurePrecedenceWithReadinessStopAndObligation|TestP23CancelledFamiliesAbsentFromProduct|TestP22FinalActivatedActionSetAcceptsProfilesAndGrants|TestP23HumanOutcomeAuthorityAndCompatibilityMatrixRealProcess|TestP23HumanOutcomeRequiredRaisedByEnvelopeActivation|TestP23HumanOutcomeRequiredRaisedByReviewRecord|TestP23HumanOutcomeRequiredRaisedByArtifactAvailability|TestV01ClosureJourneyAndRebuildParity|TestP23AttentionResolutionThroughReviewRecordedAndFindingResolvedRealProcess|TestP23AttentionResolutionHumanOutcomeReadinessRejectsWrongClassAndWrongResultAndLeavesUnrelatedStopsOpen|TestP23ConcernDeviationRealProcessJourney|TestPhase23LeavesRejectUnknownRequestFields|TestPhase23MutationsRejectMalformedObjectsWithoutEvents|TestGrantActionSetIsExactlyThirtyNames|TestProcessExitContract|TestP23ResumeRejectsJudgmentSectionOverflow|TestP22PortableFakeRejectsRemovedRequirementEvent|TestPortableRoundTripCoversEveryPhase23Table|TestPortableAssertCorruptEventsRejected|TestP23AmendedEnvelopeReplayThroughBuiltBinary|TestP23FindingDispositionRealProcessJourney|TestP23FreshnessAndOutcomeSupersessionRealProcessJourney|TestP23StopResolveHumanOutcomeConditionKeysOnConditionNotRaisingCode|TestP23StopResolveDecisionSatisfiedConditionKeysOnConditionNotRaisingCode|TestP23ConcernDispositionReferenceRejections|TestFinalCapabilityMatrixMatchesAcceptedTable|TestP23ReviewClaimRejectsDirtyPathAdoption|TestP23ReviewRunRetainsPhase22Lifecycle|TestP23HumanFixtures|TestP23HumanFixtureStepsRemainRunnable|TestP23GeneratedValueSyntax|TestP23FailureCorpusHasEveryAcceptedCode|TestP23CallerProvidedReviewIDsRemainDistinct|TestP23SharedLeafConformance|TestP23SharedAmendedLeafConformance|TestP23ReviewObligationGuardingBothActionsBlocksBothRealProcess|TestP23UnsatisfiedObligationBlocksOnlyTheActionItGuardsRealProcess|TestP23ConcernListHumanProjection|TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary|TestP23TenNameActivityInventoryFromProducingLeaves|TestP23RequirementReferenceAndLinkRealProcess|TestP23ConditionWireGrammarMatchesJournalAndFixtures|TestP22ClosedValueSchemaCases|TestP23ReviewRecordRealProcessJourney|TestP23ReviewEvidenceRealProcessJourney|TestP23ReviewCyclesRealProcessJourney|TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild|TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure|TestP23StopFailureHumanAndJSONParity|TestP23StopFailureCodeMatrix|TestP23ReviewReworkManualResolutionWithPriorOutcome|TestP23RequirementNewShapeJournalRoundTrip|TestP22FakeReviewReworkStopCondition|TestP22FakeGenericAdmissionExcludesReviewCycles|TestP22DualAdapterRecoveryConformance)$' ./cmd/fid/integration
```

### named-fid

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestPristineWorkspaceClonePreservesRebuildParityAndIsolation|TestMaintenanceCheckUsesFixed300SecondGraceAndRejectsFormat|TestReviewRecordRejectsNonIndependentActor|TestReviewRecordCreatesImmutableReview|TestReviewRecordIndependentReworkAndReceiptBinding|TestConcernListRetainsIssueScopedV1ReadUntilRemoval|TestConcernListGenericReadExcludesV1Rows|TestPhase23HumanLeafTextMatchesRecordedFacts|TestPhase22AttentionReadBoundary|TestLiveLeafCommandsMatchContractInventories|TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions|TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries|TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange|TestEnvelopeCreateCommandRejectsEveryForbiddenTypedFieldCombinationWithoutPersistenceChange|TestLinkAddRejectsRequirementEndpointsAndInvalidInputBeforeAppend|TestResumeReportsRealFindingsAndJudgmentSections|TestRequirementCommandHasReferenceLeaves|TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent|TestEvidenceEvaluateNullObligationReturnsActiveRevisionInIDOrder|TestClaimAcquireAcceptsReviewMode|TestClaimAcquireKeepsOneActiveClaimAcrossModes|TestDecodeClosedReviewInputRejectsMalformedInapplicability|TestFindingAddressDoesNotResolveOrUnblock|TestReviewRecordMixedTerminalDispositions|TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation|TestGateRegisterUnknownActorReturnsContractAuthorityFailure|TestGateRegisterDeniedActorReturnsContractAuthorityFailure|TestCapabilitiesJSONGoldenAndProjectionRows|TestLiveOutputGoldenFixtures|TestLiveGoldenChildSkipsPristineBootstrap)$' ./cmd/fid
```

### named-issues

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestV1JudgmentProjectionIsAbsentAfterReplacement|TestLinkedToDecisionRequiresApprovedReferenceInIssueScope|TestAcceptanceCriteriaRemainOrderedWithoutRoleVerifications|TestGetReturnsCentralizedDetailShape|TestListFiltersAndOrdering|TestGetReturnsEnrichedConcernFields)$' ./internal/issues
```

### named-journal

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestProjectEventFenceRejectsMalformedPhase23EventsWithoutProjectionChange|TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent|TestTask11ReviewRecordedAndFindingResolvedConditionShapes|TestTask11FindingResolvedConditionDispatchPreservesNullableBinding|TestTask11ConditionShapeRegistryIsClosedAtTenRows|TestConcernDisposedRequiresExactFieldsAndOpenGuard|TestTask8ValueShapesReuseJudgmentSubject|TestTask9ConditionDispatchRetainsPinnedPhase22Bytes|TestTask9ConditionDispatchPreservesNullableBindings|TestCurrentEventColumnDomainsMatchAcceptedDDL|TestCurrentReviewRecordedRegistersItsDistinctJudgmentValues|TestReviewRecordedAloneAdmitsRepeatedFindingDispositionEffects|TestProjectEventPopulatesPhaseTwoRepresentativeTablesViaJournalEffects|TestTask10HumanOutcomeConditionDispatchPreservesNullableBinding|TestTask10HumanOutcomeConditionShapeAndNonStringDiscriminant|TestTask9PrechangeBudgetObservedJournalReplaysWithoutMutation|TestTask9ProtectedDecisionStopDetailsAndConditionAreStrict|TestTask10HumanOutcomeStopDetailsAndConditionAreStrict|TestTask9DecisionConditionShapeAndNonStringDiscriminant|TestCurrentClosedValuesRejectIncompleteAndInvalidFacts|TestCurrentEventRejectsConditionalBindingCorruption|TestExecutionBlockedEventProjectsIdenticallyLiveAndReplay|TestActivationRecoveryEventsProjectIdenticallyLiveAndReplay)$' ./internal/journal
```

### named-links

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestRegistryHelpersCoverEndpointAndCompatibilityRules)$' ./internal/links
```

### named-maintenance

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestCheckErrorsReportsAllSeededStructuralViolationFamilies|TestCheckAcceptsPopulatedGenericReferencesWithoutStructuralErrors)$' ./internal/maintenance
```

### named-review

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestConcernDisposeCarriedForwardBindsExistingIssue|TestConcernSubjectAndDispositionMatrix|TestConcernSelfDisposalKeepsRaisingAndDisposingActors|TestDeviationAcceptsExplicitPredecessorReference|TestChangesRequiredReviewChargesExactlyOneReviewCycle|TestTerminalFindingDispositionAuthorityMatrix|TestValidateTerminalFindingInputBoundaries|TestValidateVerdictFindingAgreement|TestValidateCoverageRejectsSourceInItemAndExclusion|TestCheckIndependenceRejectsClaimAndRunHistoryDirections|TestRequirementRecordReferenceSupersedesPriorActiveRow|TestDecisionReferenceSourceMustExistAndMatchDigest)$' ./internal/review
```

### named-workflow

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -v -run '^(TestReviewClaimRejectsDirtyPathAdoption|TestEnvelopeCreateAcceptsTypedObligationSet|TestEnvelopeCreateRejectsInvalidTypedObligationAndProtectedDecisionAuthoring|TestEnvelopeCreatePersistsEveryTypedObligationAndProtectedDecision|TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults|TestEnvelopeActivationEvaluationIncludesEveryObligation|TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix|TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations|TestEnvelopeAuthoringReviewCycleBudgetMatrix|TestEnvelopeAuthoringReviewCycleAlongsideRetainedBudgets|TestEnvelopeAuthoringProtectedDecisionMatrix|TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState|TestStopResolutionChecksCurrentActivationAgreementAndCapability|TestHumanOutcomeConditionQualifyingFactRules|TestDecisionSatisfiedConditionQualifyingFactRules|TestGuardedActionsEvaluateOnlyObligationsThatGuardThem|TestCompleteEnvelopeDeniedByGuardingObligationCommitsNothing|TestRecoveryFixturesKeepRepositoryAndProjectionIsolation|TestTask9RetainedAttentionKeyMatchesPrechangeReference|TestHumanOutcomeSupersessionFollowsRecordedEventOrderBothDirections|TestHumanOutcomeEvaluationUsesLaterDifferingResult|TestAbsentHumanOutcomeFreshnessRemainsFresh|TestRegistryEqualsV01Definitions|TestRetiredJudgmentRulesAreNotImplemented|TestRecoverRunRejectsApplicableStopWithoutEffects|TestStopAndAttentionListsFilterLimitAndRemainReadOnly|TestProtectedDecisionBlocksOnlyDeclaredAffectedActions|TestTerminalSuccessEvaluationIncludesEnvelopeStop|TestEquivalentStopsShareAttentionButKeepDistinctFacts|TestReviewReworkStopRequiresApprovedOutcome|TestGenericBudgetAdmissionExcludesOnlyReworkAction|TestResumeGateReasonsExcludeReworkBudget)$' ./internal/workflow
```

### prechange-rebuild

Working directory: `/home/mark/projects/fid`.

```text
python3 - <<'PY'
from pathlib import Path
import json,subprocess,hashlib
src=Path('/tmp/fid-task20-prechange.NjCRbu/judgment/.fid/journal/2026-10-01.jsonl')
before=hashlib.sha256(src.read_bytes()).hexdigest()
assert before=='a8fdbdccfc7fd39e27b80f256811d8b70113a0e30285122a166148455d1ba16e'
args=['/tmp/fid','--project',str(src.parents[2]),'--json','rebuild']
r=subprocess.run(args,cwd=src.parents[2],text=True,capture_output=True)
row=dict(command=args,cwd=str(src.parents[2]),stdout=r.stdout,stderr=r.stderr,exit_code=r.returncode,sha256_before=before,sha256_after=hashlib.sha256(src.read_bytes()).hexdigest())
Path('project-docs/docs/phases/phase-23/qa/final-evidence-1/prechange-rebuild-receipt.json').write_text(json.dumps(row,indent=2)+'\n')
print(json.dumps(row))
assert r.returncode==6
assert 'journal-replay-stopped' in r.stdout+r.stderr
assert row['sha256_after']==before
print('PASS DEC-270 expected exit 6 refusal; source journal SHA-256 unchanged')
PY
```

### preservation-rows

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/preservation-rows.py
```

### rebuild-pragmas

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/rebuild-pragmas.py
```

### receipt-index

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/index-receipts.py
```

### reference-binary-identity

Working directory: `/home/mark/projects/fid`.

```text
go version -m /tmp/fid-p23-final-review-1-task11
```

### reference-comparison

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/compare-reference.py
```

### reference-current-probe

Working directory: `/home/mark/projects/fid`.

```text
go test -race -count=1 -overlay /home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence/current-overlay.json -run '^TestQAPhase23ReferenceAttentionKeyAndBytes$' -v ./internal/workflow
```

### reference-prechange-probe

Working directory: `/home/mark/projects/fid`.

```text
python3 /home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence-1/reference-environment.py prechange go test -race -count=1 -overlay /home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence/reference-overlay.json -run '^TestQAPhase23ReferenceAttentionKeyAndBytes$' -v ./internal/workflow
```

### reference-task11-build

Working directory: `/home/mark/projects/fid`.

```text
python3 /home/mark/projects/fid/project-docs/docs/phases/phase-23/qa/final-evidence-1/reference-environment.py task11 go build -o /tmp/fid-p23-final-review-1-task11 ./cmd/fid
```

### report-01

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-1-done.md project-docs/docs/phases/phase-23/tasks/task-1.md
```

### report-02

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-2-done.md project-docs/docs/phases/phase-23/tasks/task-2.md
```

### report-03

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-3-done.md project-docs/docs/phases/phase-23/tasks/task-3.md
```

### report-04

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-4-done.md project-docs/docs/phases/phase-23/tasks/task-4.md
```

### report-05

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-5-done.md project-docs/docs/phases/phase-23/tasks/task-5.md
```

### report-06

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-6-done.md project-docs/docs/phases/phase-23/tasks/task-6.md
```

### report-07

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-7-done.md project-docs/docs/phases/phase-23/tasks/task-7.md
```

### report-08

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-8-done.md project-docs/docs/phases/phase-23/tasks/task-8.md
```

### report-09

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-9-done.md project-docs/docs/phases/phase-23/tasks/task-9.md
```

### report-10

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-10-done.md project-docs/docs/phases/phase-23/tasks/task-10.md
```

### report-11

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-11-done.md project-docs/docs/phases/phase-23/tasks/task-11.md
```

### report-12

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-12-done.md project-docs/docs/phases/phase-23/tasks/task-12.md
```

### report-13

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-13-done.md project-docs/docs/phases/phase-23/tasks/task-13.md
```

### report-14

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-14-done.md project-docs/docs/phases/phase-23/tasks/task-14.md
```

### report-15

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-15-done.md project-docs/docs/phases/phase-23/tasks/task-15.md
```

### report-16

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-16-done.md project-docs/docs/phases/phase-23/tasks/task-16.md
```

### report-17

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-17-done.md project-docs/docs/phases/phase-23/tasks/task-17.md
```

### report-18

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-18-done.md project-docs/docs/phases/phase-23/tasks/task-18.md
```

### report-19

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-19-done.md project-docs/docs/phases/phase-23/tasks/task-19.md
```

### report-20

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-20-done.md project-docs/docs/phases/phase-23/tasks/task-20.md
```

### report-21

Working directory: `/home/mark/projects/fid`.

```text
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-23/tasks/task-21-done.md project-docs/docs/phases/phase-23/tasks/task-21.md
```

### schema-sql

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/schema-sql.py
```

### security-source-audit

Working directory: `/home/mark/projects/fid`.

```text
python3 - <<'PY'
from pathlib import Path
import json
old=json.loads(Path('project-docs/docs/phases/phase-23/qa/final-evidence/security-source-audit.json').read_text())
exec(old['cmd'].split("<<'PY'\n",1)[1].rsplit('\nPY',1)[0])
PY
```

### setup-testcache

Working directory: `/home/mark/projects/fid`.

```text
go clean -testcache
```

### source-audit

Working directory: `/home/mark/projects/fid`.

```text
python3 project-docs/docs/phases/phase-23/qa/final-evidence-1/audit-source.py
```

### tmp-after-budget-integration

Working directory: `/home/mark/projects/fid`.

```text
python3 - <<'PY'
import json, pathlib, shutil, subprocess
r=subprocess.run(['df','-B1','/tmp'],text=True,capture_output=True)
print(json.dumps({'command':['df','-B1','/tmp'],'stdout':r.stdout,'exit_code':r.returncode}))
print(json.dumps(sorted(str(p) for p in pathlib.Path('/tmp').glob('fid-integration-bin*'))))
PY
```

### tmp-after-integration

Working directory: `/home/mark/projects/fid`.

```text
python3 - <<'PY'
import json, pathlib, shutil, subprocess
r=subprocess.run(['df','-B1','/tmp'],text=True,capture_output=True)
print(json.dumps({'command':['df','-B1','/tmp'],'stdout':r.stdout,'exit_code':r.returncode}))
print(json.dumps(sorted(str(p) for p in pathlib.Path('/tmp').glob('fid-integration-bin*'))))
PY
```

### tmp-before-budget-integration

Working directory: `/home/mark/projects/fid`.

```text
python3 - <<'PY'
import json, pathlib, shutil, subprocess
r=subprocess.run(['df','-B1','/tmp'],text=True,capture_output=True)
print(json.dumps({'command':['df','-B1','/tmp'],'stdout':r.stdout,'exit_code':r.returncode}))
print(json.dumps(sorted(str(p) for p in pathlib.Path('/tmp').glob('fid-integration-bin*'))))
PY
```

### tmp-before-integration

Working directory: `/home/mark/projects/fid`.

```text
df -B1 /tmp; python3 -c 'from pathlib import Path; print(sorted(str(p) for p in Path("/tmp").glob("fid-integration-bin*")))'
```
