#!/usr/bin/env python3
"""Verify named selection completion and per-case corruption receipts."""
from pathlib import Path
import json,re
OUT=Path(__file__).resolve().parent
inv=json.loads((OUT/'named-tests.json').read_text())
count=0
for pkg,names in inv.items():
 id='named-fid' if pkg=='cmd/fid' else 'named-fid-integration' if pkg=='cmd/fid/integration' else 'named-'+pkg.rsplit('/',1)[1]
 row=json.loads((OUT/(id+'.json')).read_text())
 output=''.join(r['output'] for r in row['results'])
 assert row['results'][-1]['exit_code']==0
 passed=set(re.findall(r'^--- PASS: (Test\w+) \(',output,re.M))
 assert set(names)<=passed,(pkg,set(names)-passed)
 assert not re.search(r'--- (FAIL|SKIP):',output)
 assert 'Warning: truncated output' not in output
 count+=len(names)
 print(pkg,len(names),'selected top-level tests passed; complete inner capture')
row=json.loads((OUT/'named-fid-integration.json').read_text())
output=''.join(r['output'] for r in row['results'])
records=re.findall(r'workspace=(\S+) validator=(.*?) full projection SHA-256 before=([0-9a-f]{64}) after=([0-9a-f]{64}) tables=(\d+)',output)
primary=[r for r in records if '/TestPortableAssertCorruptEventsRejected' in r[0]]
assert len(primary)==28 and len(set(r[0] for r in primary))==28
assert len(records)==56 and all(a==b and n=='59' for _,_,a,b,n in records)
(OUT/'corruption-probes.json').write_text(json.dumps([dict(zip(('workspace','validator','before','after','tables'),r)) for r in records],indent=2)+'\n')
print('28 independent primary corrupted-event workspaces; 56 total corruption records; every 59-table SHA-256 pair equal')
assert count==188
print('PASS 188 fresh selected tests; no missing/skipped selection or truncated inner result; all corruption receipts verified')
