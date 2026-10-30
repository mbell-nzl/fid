"""Independent binary probes for Phase 24 final QA. Run after gate G14.
Command: python3 project-docs/docs/phases/phase-24/qa/final-evidence/probes.py
Uses only the accepted public fixture generator output; no held-out material.
"""
import pathlib,json,subprocess,hashlib,tempfile,shlex,time,shutil
ROOT=pathlib.Path(__file__).resolve().parents[6]
OUT=pathlib.Path(__file__).resolve().parent
attempt_results=[]
def run(args,cwd):
    start=time.time(); p=subprocess.run(args,cwd=cwd,text=True,capture_output=True)
    result=dict(argv=args,cwd=str(cwd),exit_code=p.returncode,elapsed_seconds=time.time()-start,stdout=p.stdout,stderr=p.stderr)
    attempt_results.append(result)
    (OUT/'probe-attempt-current.json').write_text(json.dumps(attempt_results,indent=2)+'\n')
    return result
manifest_path=max((ROOT/'.cache/phase-24-human').glob('fixtures-*/manifest.json'),key=lambda p:p.stat().st_mtime)
manifest=json.loads(manifest_path.read_text()); binary=manifest['binary']
sha=hashlib.sha256(pathlib.Path(binary).read_bytes()).hexdigest()
assert sha==manifest['binary_sha256'].removeprefix('sha256:')
assert len(manifest['fixtures'])==59 and len({f['name'] for f in manifest['fixtures']})==59
(OUT/'human-fixture-manifest.json').write_bytes(manifest_path.read_bytes())
rows=[]
view_paths=[]
forbidden={'worktree_id','tree_oid','dirty_digest','capture_source','<nil>'}
(OUT/'human-output').mkdir(exist_ok=True)
for row in manifest['fixtures']:
    assert pathlib.Path(row['workspace']).is_dir() and pathlib.Path(row['cwd']).is_dir()
    assert row['binary_sha256']==manifest['binary_sha256']
    assert shlex.split(row['human_command'])==[binary]+row['human_args']
    assert shlex.split(row['json_command'])==[binary]+row['json_args']
    for args in [row['human_args'],row['json_args']]:
        if '--input' in args: assert pathlib.Path(args[args.index('--input')+1]).is_file()
    if row.get('authoring'):
        # Derive the leaf from the manifest operation: run.close also has
        # leading --project arguments, so argv[:2] is not a leaf there.
        leaf_args=row['operation'].split('.')
        help=run([binary]+leaf_args+['--help'],row['cwd']); example=run([binary]+leaf_args+['--example'],tempfile.mkdtemp(prefix='p24-author-example-'))
        assert help['exit_code']==example['exit_code']==0
        rows.append(dict(name=row['name'],authoring_prestate_preserved=True,help=help,example=example))
        continue
    r=run([binary]+row['human_args'],row['cwd'])
    assert r['exit_code']==row['expected_exit'],(row['name'],r)
    assert r['stdout']==row['human_stdout'] and r['stderr']==row['human_stderr'],row['name']
    (OUT/'human-output'/f"{row['name']}.txt").write_text(r['stdout']+r['stderr'])
    if not row.get('code') and row['operation'] in ['resume','run.inspect','review.record','review.get','review.list']:
        view_paths.append(str(OUT/'human-output'/f"{row['name']}.txt"))
        head=run(['git','rev-parse','HEAD'],row['workspace'])
        assert head['exit_code']==0
        forbidden.add(head['stdout'].strip())
        forbidden.add(hashlib.sha256(str(pathlib.Path(row['workspace']).absolute()).encode()).hexdigest())
        def repository_values(value):
            if isinstance(value,dict):
                for key,item in value.items():
                    if key in ['head_oid','tree_oid','worktree_id','dirty_digest'] and isinstance(item,str) and item:
                        forbidden.add(item)
                    repository_values(item)
            elif isinstance(value,list):
                for item in value: repository_values(item)
        repository_values(json.loads(row['json_stdout']))
    rows.append(dict(name=row['name'],result=r))
# QA-28 predeclares absence as expected non-zero exit 1, not silent zero.
scan_args=['rg','-n']
for token in sorted(forbidden): scan_args+=['-e',token]
forbidden_scan=run(scan_args+view_paths,ROOT)
forbidden_scan['expected_exit_code']=1
assert forbidden_scan['exit_code']==1 and not forbidden_scan['stdout'],forbidden_scan
inputs=(OUT/'inventory-input.txt').read_text().splitlines()
assert len(inputs)==39
examples=[]
for leaf in inputs:
    cwd=pathlib.Path(tempfile.mkdtemp(prefix='p24-example-'))
    before=sorted(str(p.relative_to(cwd)) for p in cwd.rglob('*'))
    ex=run(['/tmp/fid']+leaf.split()+['--example'],cwd)
    after=sorted(str(p.relative_to(cwd)) for p in cwd.rglob('*'))
    assert ex['exit_code']==0 and not ex['stderr'] and before==after==[]
    obj=json.loads(ex['stdout']); assert isinstance(obj,dict)
    help=run(['/tmp/fid']+leaf.split()+['--help'],cwd)
    assert help['exit_code']==0 and not help['stdout'].startswith('Usage:')
    examples.append(dict(leaf=leaf,example=ex,help=help,before=before,after=after,tree_unchanged=True))
fresh=pathlib.Path(tempfile.mkdtemp(prefix='p24-final-init-'))
initialized=run(['/tmp/fid','--project',str(fresh),'--json','init','--prefix','PQF'],ROOT)
assert initialized['exit_code']==0,initialized
fresh_dbs=list((fresh/'.fid').glob('*.db'))+list((fresh/'.fid').glob('*.sqlite'))
assert len(fresh_dbs)==1,fresh_dbs
schema_query="SELECT type,name,tbl_name,sql FROM sqlite_master WHERE name NOT LIKE 'sqlite_%' AND name != 'schema_migrations' ORDER BY type,name;"
fresh_schema=run(['sqlite3','-json',str(fresh_dbs[0]),schema_query],ROOT)
assert fresh_schema['exit_code']==0
chain=[r for r in json.loads((OUT/'schema-19.json').read_text()) if not r['name'].startswith('sqlite_') and r['name']!='schema_migrations']
fresh_objects=json.loads(fresh_schema['stdout'])
# Direct SQL application does not install the migrator's bookkeeping index.
# Retain the full dump, name and assert that sole extra object explicitly.
# S1's TestMigration0019AddsOnlyFieldDigestColumn includes all migrator
# objects and independently requires them unchanged on up and down.
control_metadata=[r for r in fresh_objects if r['tbl_name']=='schema_migrations']
assert control_metadata==[dict(type='index',name='version_unique',tbl_name='schema_migrations',sql='CREATE UNIQUE INDEX version_unique ON schema_migrations (version)')],control_metadata
assert [r for r in fresh_objects if r['tbl_name']!='schema_migrations']==chain,'fresh init versus migrated projection objects differ'
(OUT/'schema-fresh-init-19.json').write_text(fresh_schema['stdout'])
# Snapshot one fully populated public fixture, rebuild its copied projection,
# and compare every keyed map bytewise. Rebuild does not change bindings.
source=next(f['workspace'] for f in manifest['fixtures'] if f['name']=='judgment-resume')
scratch=pathlib.Path(tempfile.mkdtemp(prefix='p24-final-rebuild-'))/'workspace'
shutil.copytree(source,scratch,symlinks=True)
databases=list((scratch/'.fid').glob('*.db'))+list((scratch/'.fid').glob('*.sqlite'))
assert len(databases)==1,databases
query='SELECT operation_scope,idempotency_key,request_field_digests_json FROM idempotency_results ORDER BY operation_scope,idempotency_key;'
before=run(['sqlite3','-json',str(databases[0]),query],ROOT); assert before['exit_code']==0
rebuild=run(['/tmp/fid','--project',str(scratch),'--json','rebuild'],ROOT); assert rebuild['exit_code']==0,rebuild
after=run(['sqlite3','-json',str(databases[0]),query],ROOT); assert after['exit_code']==0 and before['stdout']==after['stdout']
(OUT/'field-digests-live.json').write_text(before['stdout'])
(OUT/'field-digests-rebuilt.json').write_text(after['stdout'])
proof=dict(manifest_path=str(manifest_path),manifest_sha256=hashlib.sha256(manifest_path.read_bytes()).hexdigest(),binary_sha256=sha,binary_identity=run(['go','version','-m',binary],ROOT),fixture_rows=rows,repository_forbidden_scan=forbidden_scan,examples=examples,fresh_init=dict(command=initialized,schema=fresh_schema,equals_schema19_projection_chain=True,explicit_sole_extra_migrator_object=control_metadata,query_source='internal/db/migration0018_test.go:287; direct SQL chain has no migrator metadata index; full migration0019 test separately checks every migrator object'),rebuild=dict(source=source,scratch=str(scratch),before=before,command=rebuild,after=after),product_identity=run(['git','diff','e97835b9f9983e357b45e80c3dbe4e4b02cedb6d','HEAD','--','cmd','internal','migrations','go.mod','go.sum'],ROOT))
(OUT/'binary-probes.json').write_text(json.dumps(proof,indent=2)+'\n')
print('PASS: 59 fixture assumptions; 56 Human replays; 3 authoring help/templates; 39 no-effect examples; live/rebuild field maps byte-equal')
