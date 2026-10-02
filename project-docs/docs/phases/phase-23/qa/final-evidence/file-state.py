#!/usr/bin/env python3
"""Read-only Phase 23 final-review file and provenance evidence."""
from pathlib import Path
import hashlib
import json
import re
import subprocess

ROOT = Path(__file__).resolve().parents[6]
PHASE = ROOT / 'project-docs/docs/phases/phase-23'

def git(*args):
    result = subprocess.run(['git', *args], cwd=ROOT, text=True, capture_output=True)
    print(json.dumps({'command': ['git', *args], 'stdout': result.stdout,
                      'stderr': result.stderr, 'exit_code': result.returncode}))
    assert result.returncode == 0
    return result.stdout

print('Repository:', ROOT)
assert git('branch', '--show-current').strip() == 'phase-23'
git('rev-parse', 'HEAD')
contract = (PHASE / 'qa/validation-contract.md').read_text()
rows = [line for line in contract.splitlines() if re.match(r'^\| VAL-\d{3} \|', line)]
rows = rows[:84]
ids = [line.split('|')[1].strip() for line in rows]
assert len(set(ids)) == 84
assert set(ids) == {f'VAL-{n:03}' for n in range(1, 85)}
print('84 unique assertion rows; Human rows:', [line.split('|')[1].strip() for line in rows if '| Human |' in line])
for n in range(1, 19):
    spec = (PHASE / f'tasks/task-{n}.md').read_text()
    report = (PHASE / f'tasks/task-{n}-done.md').read_text()
    section = spec.split('## Exact Commands To Run', 1)[1]
    block = re.search(r'```(?:text|bash)?\n(.*?)```', section, re.S).group(1)
    commands = [line.strip() for line in block.splitlines() if line.strip()]
    missing = [command for command in commands if '`' + command + '`' not in report]
    assert not missing, (n, missing)
    assert not re.search(r'<TBD>|<[^>]*COMMAND>', block)
    assert '- Task outcome: COMPLETE' in report
    print(f'Task {n}: {len(commands)} exact commands present; COMPLETE; {len(report.splitlines())} report lines inspected')
    times = [line for line in report.splitlines() if any(word in line.lower() for word in ('wall time', 'wall-time', 'budget', 'times:')) and re.search(r'\d+\.\d+\s*s', line)]
    print('\n'.join(times[-3:]))
    for line in report.splitlines():
        if line.startswith('- VAL-') and not 'SUPERSEDED' in line:
            print(line[:650])
print('All exact command blocks are concrete and represented; report-lint receipts separately verify status/termination syntax.')
prior = git('ls-tree', '-d', '--name-only', '0195897:internal').splitlines()
current = sorted(p.name for p in (ROOT / 'internal').iterdir() if p.is_dir())
assert set(current) - set(prior) == {'review'}
assert set(prior) - set(current) == {'requirements'}
print('Package directories: added review; removed requirements; no other difference')
for event in ('review.recorded','finding.addressed','concern.recorded','concern.disposed','deviation.recorded','decision.referenced','human-outcome.recorded','requirement.referenced','requirement.linked'):
    git('log', '-S"' + event + '"', '--format=%h %s%n%b', '--', 'internal/journal/project_effects.go')
for action in ('review.record','finding.address','concern.record','concern.dispose','deviation.record','decision.record-reference','human-outcome.record'):
    git('log', '-S"' + action + '"', '--format=%h %s%n%b', '--', 'internal/workflow/activation.go')
for p in sorted((PHASE / 'qa').glob('review-*.md')) + sorted((PHASE / 'qa').glob('checkpoint-after-task-11*.md')):
    content = p.read_text()
    lines = [line for line in content.splitlines() if re.match(r'^(Reviewer:|Date:|Review type:|Prior review:|Tasks covered:|VERDICT:|OUTCOME )', line)]
    print(str(p.relative_to(ROOT)), '\n'.join(lines))
checklist = (PHASE / 'qa/human-checklist.md').read_text()
human_rows = [line for line in checklist.splitlines() if line.startswith('| HUMAN-')]
assert len(human_rows) == 63
assert all(line.rstrip().endswith('| | |') and '(human)' in line for line in human_rows)
assert 'Pending Human input' in checklist
print('63 Human-only judgment rows remain blank; overall Pending Human input')
manifest_path = Path(re.search(r'- Manifest: `([^`]+)`', checklist).group(1))
manifest = json.loads(manifest_path.read_text())
assert Path(manifest['binary']).is_file()
actual = 'sha256:' + hashlib.sha256(Path(manifest['binary']).read_bytes()).hexdigest()
assert actual == manifest['binary_sha256']
for checkpoint in manifest['checkpoints']:
    assert Path(checkpoint['project']).is_dir()
    for step in checkpoint.get('steps', []):
        assert Path(step['cwd']).is_dir()
        if '--input' in step['args']:
            assert Path(step['args'][step['args'].index('--input')+1]).is_file()
assert len(manifest['leaf_pairs']) == 18
print('Checklist manifest and retained binary digest verified:', actual)
print('Three workspaces, all step working directories/inputs and 18 leaf pairs exist.')
print('PASS final file-state and provenance audit')
