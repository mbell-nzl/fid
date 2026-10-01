#!/usr/bin/env python3
"""Render an index of the full execution results captured by this QA session."""
from pathlib import Path
import json
OUT=Path(__file__).resolve().parent
def esc(value): return str(value).replace('|','\\|').replace('\n','<br>').replace('`','&#96;')
rows=[]
for f in sorted(OUT.glob('*.json')):
 data=json.loads(f.read_text())
 if not isinstance(data,dict) or not {'id','cmd','results'}<=data.keys(): continue
 last=data['results'][-1]
 assert 'exit_code' in last,(f.name,'unterminated')
 output=''.join(r.get('output','') for r in data['results'])
 token=output.rstrip().splitlines()[-1] if output.strip() else f"exit {last['exit_code']} (no output)"
 status='PASS' if last['exit_code']==0 else 'FAIL'
 if data['id'].startswith('absence-'):
  assert last['exit_code']==1 and not output.strip()
  status='PASS (expected exit_code: 1)'
 if data['id'].startswith('setup-'): status='SETUP ONLY'
 if data['id']=='report-timing-excerpts' and not output: status='FILE EVIDENCE'
 rows.append((data,f.name,status,token,last['exit_code']))
text=['# Phase 23 final command receipts','',
'All commands below were executed in the third phase-end QA session. Each linked JSON contains the exact command, working directory (default /home/mark/projects/fid for older entries), every captured result including session IDs, full available output, and observed final exit_code. Prior child sessions provide no passing command evidence.','',
'No full checkptr rerun was used to replace its failed attempt. Focused and planned named selections are separate evidence. Setup and cleanup are not gates. The truncated historical timing grep is retained as a diagnostic read; complete report receipts and timing excerpts supply that audit.','',
'| Receipt | Result | Verbatim final output token | Observed termination |','|---|---|---|---|']
for data,f,status,token,code in rows:
 text.append(f"| [{data['id']}]({f}) | {status} | `{esc(token)}` | `exit_code: {code}` |")
text+=['','## Exact commands','',
'Multiline Python probes retain their exact full source in their receipt and, where applicable, in the adjacent script. Their child-command receipts record actual return codes.','']
for data,f,status,token,code in rows:
 text+=['### '+data['id'],'',f"Working directory: `{data.get('cwd','/home/mark/projects/fid')}`",'']
 cmd=data['cmd']
 if '\n' in cmd: text += [f"Exact multiline invocation: [{f}]({f}).",'']
 else: text += ['\x60\x60\x60text',cmd,'\x60\x60\x60','']
(OUT/'command-index.md').write_text('\n'.join(text).rstrip()+'\n')
print(f'PASS indexed {len(rows)} terminated command receipts with exact commands, final tokens and actual exit codes')
