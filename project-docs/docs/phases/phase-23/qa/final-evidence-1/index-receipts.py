#!/usr/bin/env python3
"""Index independently captured command results without inventing termination."""
from pathlib import Path
import json
from datetime import datetime

OUT=Path(__file__).resolve().parent
entries=[]
for p in sorted(OUT.glob('*.json')):
    data=json.loads(p.read_text())
    if not isinstance(data,dict) or not {'id','cmd','results'}<=data.keys():continue
    last=data['results'][-1]
    assert 'exit_code' in last,(p,'missing observed termination')
    output=''.join(x.get('output','') for x in data['results'])
    assert 'Warning: truncated output' not in output,(p,'inner output truncated')
    lines=output.strip().splitlines()
    token=lines[-1] if lines else f'exit {last["exit_code"]} (no output)'
    # Start/end delimit observation, including any polling delay.
    elapsed=None
    if data.get('started') and data.get('ended'):
        elapsed=(datetime.fromisoformat(data['ended'])-datetime.fromisoformat(data['started'])).total_seconds()
    entries.append(dict(id=data['id'],command=data['cmd'],cwd=data.get('cwd'),
                        token=token,exit_code=last['exit_code'],
                        observation_seconds=elapsed,receipt=p.name))
(OUT/'command-index.json').write_text(json.dumps(entries,indent=2)+'\n')
text='# Phase 23 final re-review — command receipts\n\n'
text+='Every status below is the nested command’s observed termination. Child processes inside Python observations also retain their individual statuses in linked receipts. These are QA observations, not borrowed Developer gate results. Observation duration includes polling delay; Go package timings are in the raw output. Setup, help, and file inventory commands are identified by their names and are not product gates.\n\n'
text+='| ID | Observed exit_code | Last output token | Observation seconds |\n|---|---:|---|---:|\n'
for r in entries:
    token=r['token'].replace('|','\\|').replace('\x60','')
    text+=f"| [{r['id']}]({r['receipt']}) | {r['exit_code']} | \x60{token}\x60 | {r['observation_seconds'] if r['observation_seconds'] is not None else 'not measured'} |\n"
text+='\n## Exact commands and working directories\n'
for r in entries:
    text+=f"\n### {r['id']}\n\nWorking directory: \x60{r['cwd'] or '/home/mark/projects/fid'}\x60.\n\n\x60\x60\x60text\n{r['command']}\n\x60\x60\x60\n"
(OUT/'command-index.md').write_text(text)
print(f'PASS {len(entries)} receipts indexed; all commands have observed nested exit_code and complete captured output')
