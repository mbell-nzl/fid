#!/usr/bin/env python3
"""Reproduce retained fixture text without consuming Human judgment steps."""
from pathlib import Path
import hashlib,json,re,subprocess
ROOT=Path(__file__).resolve().parents[6]
OUT=Path(__file__).resolve().parent
qa=OUT.parent
checklist=(qa/'human-checklist.md').read_text()
guide=(qa/'human-fixtures.md').read_text()
manifest_path=Path(re.search(r'- Manifest: `([^`]+)`',checklist).group(1))
manifest=json.loads(manifest_path.read_text())
binary='/tmp/fid'
assert hashlib.sha256(Path(binary).read_bytes()).hexdigest()==manifest['binary_sha256'].removeprefix('sha256:')
records=[]
def call(cwd,args,expected=0):
 r=subprocess.run([binary,*args],cwd=cwd,text=True,capture_output=True)
 record={'cwd':cwd,'command':[binary,*args],'stdout':r.stdout,'stderr':r.stderr,'exit_code':r.returncode}
 records.append(record)
 (OUT/'human-live-receipts.json').write_text(json.dumps(records,indent=2)+'\n')
 assert r.returncode==expected,(args,r.returncode,r.stderr)
 assert r.stdout and not r.stderr,(args,r.stderr)
 return r.stdout
for cp in manifest['checkpoints']:
 for inspection in cp['inspections']:
  text=call(inspection['cwd'],inspection['args'])
  if '--json' in inspection['args']: json.loads(text)
  else: assert text==cp['expected_human']['\0'.join(inspection['args'])],(cp['name'],inspection['args'])
 print('PASS checklist checkpoint inspections:',cp['name'],len(cp['inspections']))
for pair in manifest['leaf_pairs']:
 human=call(pair['cwd'],pair['human_args'])
 facts=json.loads(call(pair['cwd'],pair['json_args']))
 assert human==pair['expected_human']
 assert facts['operation']==pair['operation'] and facts['result']==pair['json_facts']
 assert 'PRIVATE_' not in human
 print('PASS checklist pair:',pair['operation'],'exact text and JSON facts')
# The guide's 18-leaf appendix is its own retained identity; reproduce its literal pairs.
sections=re.findall(r'### `([^`]+)`\n(.*?)(?=\n### `|\Z)',guide,re.S)
assert len(sections)==18,len(sections)
appendix=guide.split('## 2026-09-30 QA repair',1)[1]
guide_manifest_path=Path(re.search(r'- Manifest: `([^`]+)`',appendix).group(1))
guide_manifest=json.loads(guide_manifest_path.read_text())
guide_pairs={pair['operation']:pair for pair in guide_manifest['leaf_pairs']}
assert len(guide_pairs)==18
def published_facts(expected,actual,path=''):
 if isinstance(expected,dict):
  assert isinstance(actual,dict) and expected.keys()<=actual.keys(),path
  for key,value in expected.items(): published_facts(value,actual[key],path+'/'+key)
 elif isinstance(expected,list):
  assert isinstance(actual,list) and len(expected)==len(actual),path
  for i,(a,b) in enumerate(zip(expected,actual)): published_facts(a,b,path+'/'+str(i))
 else: assert expected==actual,(path,expected,actual)
for operation,body in sections:
 cwd=re.search(r'- Working directory: `([^`]+)`',body).group(1)
 human_args=json.loads(re.search(r'- Human arguments: `([^\n]+)`',body).group(1))
 json_args=json.loads(re.search(r'- JSON arguments: `([^\n]+)`',body).group(1))
 expected_human=re.search(r'```text\n(.*?)\n```',body,re.S).group(1)+'\n'
 expected_json=json.loads(re.search(r'```json\n(.*?)\n```',body,re.S).group(1))
 human=call(cwd,human_args)
 actual=json.loads(call(cwd,json_args))
 assert human==expected_human,operation
 assert actual['operation']==operation and actual['result']==guide_pairs[operation]['json_facts'],operation
 published_facts(expected_json,actual['result'])
 assert 'PRIVATE_' not in human
 if operation.startswith(('human-outcome.','concern.','deviation.')):
  assert 'not deterministic proof' in human,operation
 print('PASS guide literal pair:',operation,'exact Human and complete manifest JSON equality; every published prose-omitted fact matches')
assert len([s for s in checklist.splitlines() if s.startswith('| HUMAN-') and s.endswith('| | |')])==63
print('PASS 36 retained Human/JSON leaf pairs and all checkpoint reads; 63 judgments remain blank')
