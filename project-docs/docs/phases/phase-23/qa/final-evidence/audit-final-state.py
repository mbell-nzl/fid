#!/usr/bin/env python3
"""Verify the QA evidence's terminal states and protected workspace boundaries."""
from pathlib import Path
import hashlib, json, re, subprocess

OUT = Path(__file__).resolve().parent
ROOT = OUT.parents[5]
QA = OUT.parent

def git(*args, cwd=ROOT):
    r = subprocess.run(['git', *args], cwd=cwd, text=True, capture_output=True)
    print(json.dumps({'command':['git', *args], 'cwd':str(cwd), 'stdout':r.stdout, 'stderr':r.stderr, 'exit_code':r.returncode}))
    assert r.returncode == 0
    return r.stdout.strip()

assert git('branch', '--show-current') == 'phase-23'
assert git('rev-parse', 'HEAD') == 'e0893cdee915013371455198906852081633e051'
baseline = json.loads((OUT/'protected-file-hashes.json').read_text())
original = ''.join(r.get('output','') for r in baseline['results'])
for line in original.splitlines():
    match = re.fullmatch(r'([a-f0-9]{64})  (.+)',line)
    if not match: continue
    digest, path = match.groups()
    if path.endswith('/human-checklist.md'): continue
    actual = hashlib.sha256((ROOT/path).read_bytes()).hexdigest()
    assert actual == digest, path
    print('UNCHANGED', actual, path)

def human_rows(path):
    return [s for s in path.read_text().splitlines() if s.startswith('| HUMAN-')]
rows = human_rows(QA/'human-checklist.md')
assert rows == human_rows(OUT/'human-checklist-before-refresh.md')
assert len(rows) == 63 and all(s.endswith('| | |') and '(human)' in s for s in rows)
print('PASS 63 original Human rows remain byte-identical and blank; sha256='+hashlib.sha256('\n'.join(rows).encode()).hexdigest())
tracked = git('diff','--name-only').splitlines()
assert set(tracked) == {'project-docs/docs/current/tl_escalations_current.md','project-docs/docs/phases/phase-23/run-ledger.md','project-docs/docs/phases/phase-23/triage-log.md'}
assert git('diff','--cached','--name-only') == ''
for path in git('ls-files','--others','--exclude-standard').splitlines():
    assert path.startswith('project-docs/docs/phases/phase-23/qa/'),path
print('PASS tracked dirt equals the three pre-existing Runner/TL files; all new files are QA-owned')
for suffix,head in [('task11','3dff4574a71438900100449e36d7044969ba9162'),('prechange','74c83b7033729785eb99403ecec9dcde30ac4280')]:
    wd=Path('/tmp/fid-p23-qa-reference-PgIz9Jst')/suffix
    assert git('rev-parse','HEAD',cwd=wd) == head
    assert git('status','--porcelain',cwd=wd) == ''
print('PASS both reference worktrees retain the exact clean detached commits')
failed=[]
count=0
for path in sorted(OUT.glob('*.json')):
    data=json.loads(path.read_text())
    if not isinstance(data,dict) or not {'id','cmd','results'} <= data.keys(): continue
    count+=1
    result=data['results'][-1]
    assert 'exit_code' in result,(path.name,'unterminated')
    if result['exit_code'] and not data['id'].startswith('absence-'):
        failed.append(data['id'])
    if data['id'].startswith('absence-'):
        assert result['exit_code']==1 and not ''.join(x.get('output','') for x in data['results']).strip()
expected={'full-checkptr','reference-worktree-setup','feasibility-fid','human-live','human-live-rerun','reference-task11-build','reference-task11-build-diagnostic'}
assert set(failed)==expected,failed
print('Failed attempts retained:',json.dumps(failed))
checks=json.loads((OUT/'checklist-evidence.json').read_text())
assert len(checks)==52 and {r['id'] for r in checks}=={f'QA-{n:02}' for n in range(1,53)}
vals={v for r in checks for v in r['vals']}
assert vals=={f'VAL-{n:03}' for n in range(1,85)}-{'VAL-070','VAL-071','VAL-072','VAL-073'}
for row in checks:
    assert all((OUT/(receipt+'.json')).is_file() for receipt in row['receipts'])
print(f'PASS {count} terminated command receipts; seven failed attempts retained; 52 QA checks cover 80 deterministic assertions')
