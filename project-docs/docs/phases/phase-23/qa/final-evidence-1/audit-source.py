#!/usr/bin/env python3
"""Independent QA file-state, report-command, and repair-boundary audit."""
from pathlib import Path
import hashlib, json, re, subprocess

ROOT = Path(__file__).resolve().parents[6]
OUT = Path(__file__).resolve().parent
PHASE = ROOT / 'project-docs/docs/phases/phase-23'
records = []
def command(args, expected=0):
    r = subprocess.run(args, cwd=ROOT, text=True, capture_output=True)
    records.append(dict(command=args, cwd=str(ROOT), stdout=r.stdout,
                        stderr=r.stderr, exit_code=r.returncode))
    (OUT/'source-audit-receipts.json').write_text(json.dumps(records, indent=2)+'\n')
    assert r.returncode == expected, records[-1]
    return r.stdout
def git(*args):
    return command(['git', *args])

assert git('rev-parse','HEAD').strip() == '32a55a94efb2a89dbf50b29accb434b9c753418c'
contract=(PHASE/'qa/validation-contract.md').read_text()
rows=[x for x in contract.splitlines() if re.match(r'^\| VAL-\d{3} \|',x)]
assert len(rows)==89
assert {x.split('|')[1].strip() for x in rows}=={f'VAL-{n:03}' for n in range(1,90)}
print('89 unique assertions; Human-owned VAL-070..073')
report_audit=[]
for n in range(1,22):
    spec=(PHASE/f'tasks/task-{n}.md').read_text()
    report=(PHASE/f'tasks/task-{n}-done.md').read_text()
    section=spec.split('## Exact Commands To Run',1)[1]
    block=re.search(r'\x60\x60\x60(?:text|bash)?\n(.*?)\x60\x60\x60',section,re.S).group(1)
    commands=[x.strip() for x in block.splitlines() if x.strip()]
    assert all('\x60'+x+'\x60' in report for x in commands),(n,commands)
    assert not re.search(r'<TBD>|<[^>]*COMMAND>',block)
    assert '- Task outcome: COMPLETE' in report,n
    command_rows=[x for x in report.splitlines() if re.match(r'^\s*- \x60.+\x60 - (?:PASS|FAIL|BLOCKED)',x)]
    spec_sections=re.findall(r'(?ms)^## (?:Acceptance Criteria|Validation Contract Coverage|Expected Non-Zero Commands|Concerns.*|Scope|Deviations.*|Dependency.*|Exact Commands To Run).*?(?=^## |\Z)',spec)
    report_sections=re.findall(r'(?ms)^## (?:Validation Contract Coverage|Concerns.*|Deviations.*|Dependency.*).*?(?=^## |\Z)',report)
    timing=[x for x in report.splitlines() if any(w in x.lower() for w in ('wall time','wall-time','budget','times:','seconds')) and re.search(r'\d+\.\d+',x)]
    report_audit.append(dict(task=n,commands=commands,command_rows=command_rows,
        spec_sections=spec_sections,report_sections=report_sections,timing=timing,
        spec_sha256=hashlib.sha256(spec.encode()).hexdigest(),
        report_sha256=hashlib.sha256(report.encode()).hexdigest()))
    print(f'Task {n}: exact command block ({len(commands)}) present; COMPLETE; report-lint evidence recorded separately')
(OUT/'task-audit.json').write_text(json.dumps(report_audit,indent=2)+'\n')

repair_commits={
 'b84b17c7b7218f981f368f908225bd139fab926e':{
 'cmd/fid/workspace_fixture_test.go','cmd/fid/live_golden_child_bootstrap_test.go',
 'project-docs/docs/phases/phase-23/tasks/task-18-done.md'},
 'b37c16986d3dbba8bc786089c2869e3016009b56':{
 'project-docs/docs/phases/phase-23/tasks/task-18-done.md'},
 '9ea1f5533000d7cb12a48ba0ce92d0729d3d268b':{
 'project-docs/docs/phases/phase-23/tasks/task-19-done.md'},
 '025aa410c9efb5c344398d4f45cce6c5798c7529':{
 'project-docs/docs/phases/phase-23/tasks/task-21-done.md'}}
for commit,allowed in repair_commits.items():
    git('merge-base','--is-ancestor',commit,'HEAD')
    names=set(git('diff-tree','--no-commit-id','--name-only','-r',commit).splitlines())
    assert names==allowed,(commit,names,allowed)
    git('show','--format=fuller',commit,'--')
    print('Verified ancestry and exact repair boundary:',commit,sorted(names))
assert not git('diff','b84b17c^','b37c169','--','cmd/fid/live_output_test.go','cmd/fid/testdata/val-036')
assert '5000' in (ROOT/'cmd/fid/live_output_test.go').read_text()
print('DEC-267: golden bytes, live-output assertions and 5000ms limit unchanged by repair')
for n in (19,21):
    report=(PHASE/f'tasks/task-{n}-done.md').read_text()
    reds=re.findall(r'^- \x60([^\x60]+)\x60 - FAIL ',report,re.M)
    greens=re.findall(r'^- \x60([^\x60]+)\x60 - PASS ',report,re.M)
    missing=[c for c in reds if c not in greens]
    assert not missing,(n,missing)
    print(f'Task {n}: {len(reds)} failed command strings each have byte-identical PASS follow-up')
assert not git('diff','0195897','--','go.mod','go.sum')
prior=set(git('ls-tree','-d','--name-only','0195897:internal').splitlines())
current={p.name for p in (ROOT/'internal').iterdir() if p.is_dir()}
assert current-prior=={'review'} and prior-current=={'requirements'}
print('Dependency files unchanged; only package addition review and removal requirements')
git('diff','0195897','--','migrations/0013_p21_operational_kernel.up.sql','migrations/0015_p22_recovery.up.sql')
events=('review.recorded','finding.addressed','concern.recorded','concern.disposed',
        'deviation.recorded','decision.referenced','human-outcome.recorded',
        'requirement.referenced','requirement.linked')
for name in events:
    git('log','-S"'+name+'"','--format=%h %s%n%b','--','internal/journal/project_effects.go')
for name in ('review.record','finding.address','concern.record','concern.dispose',
             'deviation.record','decision.record-reference','human-outcome.record'):
    git('log','-S"'+name+'"','--format=%h %s%n%b','--','internal/workflow/activation.go')
for pattern in ('task_completions|ac_verifications',r'requirement\.synced|last_synced_hash',
                'has_completion_report|concerns_addressed'):
    command(['rg','-n','--glob','!*_test.go',pattern,'cmd','internal'],expected=1)
command(['rg','-n','unavailableJudgmentEvents','cmd','internal'],expected=1)
print('Four expected-nonzero absence commands: exit 1 (no output); observed exit_code: 1')
for p in sorted((PHASE/'qa').glob('review-*.md'))+sorted((PHASE/'qa').glob('checkpoint-after-task-11*.md')):
    lines=[x for x in p.read_text().splitlines() if re.match(r'^(Reviewer:|Review type:|Prior review:|VERDICT:|OUTCOME )',x)]
    print(p.name, '; '.join(lines))
old=(PHASE/'qa/human-checklist.md').read_text()
humans=[x for x in old.splitlines() if x.startswith('| HUMAN-')]
outcomes={}
for row in humans:
    status=row.split('|')[5].strip().upper()
    outcomes[status]=outcomes.get(status,0)+1
assert len(humans)==63 and outcomes=={'PASS':55,'PARTIAL':6,'FAIL':2},outcomes
print('Preserved Human results:',outcomes)
print('PASS source, report command, repair scope and provenance audit')
