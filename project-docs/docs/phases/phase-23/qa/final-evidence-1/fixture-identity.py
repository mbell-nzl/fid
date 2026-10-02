#!/usr/bin/env python3
from pathlib import Path
import hashlib,json,subprocess

OUT=Path(__file__).resolve().parent
ROOT=OUT.parents[5]
m=json.loads((OUT/'human-fixture-manifest.json').read_text())
manifest=Path(m['binary']).parent/'manifest.json'
assert manifest.read_bytes()==(OUT/'human-fixture-manifest.json').read_bytes()
receipts=[]
for argv in [['go','version','-m',m['binary']],['go','version','-m','/tmp/fid'],['git','rev-parse','HEAD'],['git','diff','HEAD','--','cmd','internal','migrations','go.mod','go.sum']]:
    r=subprocess.run(argv,cwd=ROOT,text=True,capture_output=True)
    receipts.append(dict(argv=argv,cwd=str(ROOT),stdout=r.stdout,stderr=r.stderr,exit_code=r.returncode))
    assert r.returncode==0
assert receipts[2]['stdout'].strip()==m['source_commit']=='32a55a94efb2a89dbf50b29accb434b9c753418c'
assert not receipts[3]['stdout']
for r in receipts[:2]:
    assert 'vcs.revision='+m['source_commit'] in r['stdout']
    assert 'vcs.modified=true' in r['stdout']
binary_sha='sha256:'+hashlib.sha256(Path(m['binary']).read_bytes()).hexdigest()
assert binary_sha==m['binary_sha256']
phase_sha='sha256:'+hashlib.sha256(Path('/tmp/fid').read_bytes()).hexdigest()
d=dict(manifest=str(manifest),manifest_sha256=hashlib.sha256(manifest.read_bytes()).hexdigest(),binary=m['binary'],binary_sha256=binary_sha,phase_binary_sha256=phase_sha,binaries_equal=binary_sha==phase_sha,source_head=m['source_commit'],source_dirty=m['source_dirty'],product_diff_empty=True,receipts=receipts)
(OUT/'fixture-identity-receipts.json').write_text(json.dumps(d,indent=2)+'\n')
print(json.dumps(d,indent=2))
print('PASS fresh manifest and binary digest verified; binary revisions match reviewed HEAD; product diff empty')
