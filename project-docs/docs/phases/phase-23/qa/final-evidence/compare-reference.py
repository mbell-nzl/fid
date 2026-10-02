#!/usr/bin/env python3
"""Compare independently executed historical/current byte and key probes."""
from pathlib import Path
import base64,gzip,hashlib,json,re
OUT=Path(__file__).resolve().parent
ROOT=OUT.parents[5]
seen=[]
for name in ('reference-prechange-probe','reference-current-probe'):
 data=json.loads((OUT/(name+'.json')).read_text())
 assert data['results'][-1]['exit_code']==0
 if name=='reference-prechange-probe':
  child=json.loads((OUT/'reference-prechange-bound-command.json').read_text())
  assert child['exit_code']==0
  output=child['stdout']
 else:
  output=''.join(row.get('output','') for row in data['results'])
 canonical=dict(re.findall(r'canonical ([^:]+): (\{[^\n]+\})',output))
 keys=re.findall(r' key=(sha256:[0-9a-f]+)',output)
 assert len(canonical)==6 and len(keys)==1,(canonical,keys)
 seen.append((canonical,keys[0]))
assert seen[0]==seen[1]
pinned=(ROOT/'internal/journal/current_event_closed_values_test.go').read_text()
for kind,raw in seen[0][0].items():
 assert raw in pinned,(kind,raw)
 print('MATCH canonical',kind,raw)
print('MATCH attentionSemanticKey',seen[0][1])
fixture=ROOT/'internal/journal/testdata/task9-prechange-budget-observed.jsonl.gz.b64'
raw=gzip.decompress(base64.b64decode(fixture.read_text()))
digest=hashlib.sha256(raw).hexdigest()
assert digest=='8fdf19380c9cf8e1627bbde730e5554e3d0f3b3006ffcff7e21d6162197f11df'
assert len(raw.splitlines())==18
print('Pinned source journal:',str(fixture.relative_to(ROOT)),'events=18 sha256='+digest)
print('PASS six canonical bytes and Attention key identical at 74c83b7 and e0893cd; pinned journal identity verified')
