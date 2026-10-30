"""Phase 24 read-only documentation, identity, inventory and fixture audit.
Run from repo root: python3 project-docs/docs/phases/phase-24/qa/final-evidence/audit.py docs
Additional modes: inventory, supplemental, schema. Heavy modes run serially.
No product, task, contract or existing review is edited.
"""
import pathlib, subprocess, re, json, hashlib, sys, time, sqlite3, tempfile, os
ROOT=pathlib.Path(__file__).resolve().parents[6]
BASE=ROOT/'project-docs/docs/phases/phase-24'
OUT=BASE/'qa/final-evidence'
TASKS=['1','2','3','4','5','6','7','8','9','9a','10','11','12','13','14']
def run(command,cwd=ROOT,expected=0):
    started=time.time()
    p=subprocess.run(command,shell=True,executable='/bin/bash',cwd=cwd,capture_output=True,text=True)
    return dict(command=command,cwd=str(cwd),exit_code=p.returncode,expected_exit_code=expected,elapsed_seconds=time.time()-started,stdout=p.stdout,stderr=p.stderr)
def save(name,value):
    (OUT/name).write_text(json.dumps(value,indent=2)+'\n')
def docs():
    inventory_rows=[]
    for line in (BASE/'input-inventory.md').read_text().splitlines():
        if not re.match(r'\| R\d+ \|',line): continue
        cells=[c.strip() for c in re.split(r'(?<!\\)\|',line)[1:-1]]
        ident,registry,command,count,owner=cells
        spans=[s.replace('\\|','|') for s in re.findall(r'`([^`]+)`',command) if s.startswith(('rg ','ls '))]
        if not spans: continue
        if owner in ['3-9','4-8','6-8']:
            owned={'3-9':['3','4','5','6','7','8','9'],'4-8':['4','5','6','7','8'],'6-8':['6','7','8']}[owner]
        else:
            match=re.match(r'([0-9]+a?)',owner)
            owned=[match.group(1)] if match else []
        for task in owned:
            data=(BASE/f'tasks/task-{task}-done.md').read_text()
            inventory_rows.append(dict(registry=ident,task=task,command=spans[0],receipt_present=spans[0] in data))
    rows=[]
    for t in TASKS:
        path=BASE/f'tasks/task-{t}-done.md'; data=path.read_text()
        current='\n'.join(l for l in data.splitlines() if not l.startswith('>'))
        h=re.search(r'## Commit hash\s+- `([0-9a-f]{40})`',current)
        subject=re.search(r'## Commit subject\s+- `([^`]+)`',current)
        identity=run('git show -s --format=fuller '+h.group(1)) if h else None
        state=run('git status --porcelain -- '+str(path.relative_to(ROOT)))
        blob=run('git show HEAD:'+str(path.relative_to(ROOT)))
        rows.append(dict(task=t,report=str(path.relative_to(ROOT)),sha256=hashlib.sha256(path.read_bytes()).hexdigest(),claimed_hash=h.group(1) if h else None,claimed_subject=subject.group(1) if subject else None,identity=identity,report_clean=not state['stdout'],state=state,equals_HEAD_blob=blob['stdout']==data,commit_receipt=run('git log -1 --format=fuller -- '+str(path.relative_to(ROOT))),inventory_receipt_literals={ 'R4': 'ls migrations' in data and 'ProjectionSchemaVersion' in data, 'R10': 'rg -n \'"error: \' cmd internal' in data},timing_lines=[l for l in current.splitlines() if re.search(r'600|900|480|720|slowest|cached',l)]))
    reviewrows=[]
    for i in range(1,13):
        path=BASE/f'qa/review-{i}.md'
        data=path.read_text()
        verdict=re.search(r'(?m)^.*(?:Verdict|Outcome).*$',data)
        reviewrows.append(dict(review=i,sha256=hashlib.sha256(path.read_bytes()).hexdigest(),receipt=run('git log -1 --format=fuller -- '+str(path.relative_to(ROOT))),clean=run('git status --porcelain -- '+str(path.relative_to(ROOT))),verdict_heading=verdict.group(0) if verdict else None))
    save('document-audit.json',dict(owned_inventory_receipts=inventory_rows,tasks=rows,reviews=reviewrows,dependency_diff=run('git diff 304aa45 -- go.mod go.sum'),spec_diff=run('git diff 304aa45 -- project-docs/specs/03_TRD_SCHEMA.md project-docs/specs/05_DESIGN_PLAN.md project-docs/specs/06_DESIGN_PLAN_CLI.md project-docs/specs/07_DESIGN_PLAN_FLOWS.md project-docs/specs/08_QUALITY_PLAN.md project-docs/specs/GLOSSARY.md'),package_changes=run('git diff --name-status 304aa45 -- cmd internal'),head=run('git rev-parse HEAD'),branch=run('git branch --show-current')))
    for row in rows: print(row['task'],row['claimed_hash'],row['report_clean'],row['equals_HEAD_blob'],row['inventory_receipt_literals'])
def inventory():
    data=(BASE/'input-inventory.md').read_text()
    cmds=[]
    for line in data.splitlines():
        m=re.match(r'\| (R\d+) \|',line)
        if not m: continue
        code=re.findall(r'`([^`]+)`',line)
        executable=[span for span in code if span.startswith(('rg ','ls '))]
        if not executable: continue
        command=executable[0].replace('\\|','|')
        expected=1 if m.group(1)=='R5' else 0
        r=run(command,expected=expected); r['registry']=m.group(1); cmds.append(r)
        print(m.group(1),r['exit_code'],len(r['stdout'].splitlines()))
    assert {c['registry'] for c in cmds}=={'R1','R2','R3','R4','R5','R6','R7','R9','R10','R11','R13','R14','R19','R20'}
    assert all(c['exit_code']==c['expected_exit_code'] for c in cmds)
    block=re.search(r'```bash\n(.*?)\n```',data,re.S).group(1)
    r=run(block); r['registry']='section1'; cmds.append(r)
    # wc output retains exact scratch paths; copy every labelled list to durable evidence.
    paths=re.findall(r'^\s*\d+ (\S+/(?:leaves|groups|input|keyed|no-short|typed-flags))$',r['stdout'],re.M)
    for path in paths: (OUT/('inventory-'+pathlib.Path(path).name+'.txt')).write_bytes(pathlib.Path(path).read_bytes())
    cmds.append(run("ls internal/cli/codes.go internal/cli/*codes_test.go cmd/fid/stop_human.go cmd/fid/resume_human.go cmd/fid/resume_human_guidance.go"))
    cmds.append(run("rg -n \"This (Concern|Deviation) is recorded judgment|This Outcome is authorized Human judgment\" cmd/fid -g '!*_test.go'",expected=1))
    cmds.append(run('rg -n "PRIVATE_" cmd/fid/integration',expected=1))
    for c in ['go version','go env GOOS GOARCH','uname -srm','nproc',"grep -m1 '^model name' /proc/cpuinfo","grep -m1 '^MemTotal:' /proc/meminfo"]:
        cmds.append(run(c))
    for c in ['go help testflag','go help build','golangci-lint run --help','go tool cover -help','gofmt -h','rg --help','project-docs/tools/report-lint.sh']:
        cmds.append(run(c,expected=2 if c in ['go tool cover -help','project-docs/tools/report-lint.sh'] else 0))
    save('inventory-and-help.json',cmds)
    coverage=run('go tool cover -func=coverage.out')
    (OUT/'coverage-functions.txt').write_text(coverage['stdout'])
    save('coverage-report-complete.json',dict(result=coverage,profile_sha256=hashlib.sha256((ROOT/'coverage.out').read_bytes()).hexdigest()))
    print('section1',r['exit_code'],r['stdout'])
def supplemental():
    cmds=[
    "go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 -v ./cmd/fid ./internal/db ./internal/journal ./internal/project ./internal/export ./internal/workflow ./internal/cli -run 'Test(P24|CurrentIdempotencyField|Schema19|GateReceiptFieldDigest|Migration0019|ComputeRequestDigest|RegistryNextStates|Input|LiveLeaf|RequestDigest|RebuildWorkspacePreserves|Grammar|EveryLiveLeaf|Adopted|Example|Structural|ResolveSource|ObserveRepository|HumanFailure|BoundedMessage|GenericFallback|Phase23HumanLeaf|Task9Prechange|CommitIdempotentReplays)'",
    "go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -count=1 -v ./cmd/fid/integration/... -run 'Test(P24|P23HumanFixtures|P23ConcernListHumanProjection|PortableRoundTripCoversEveryPhase23Table|P22CloseRunUsesObserved|P22DualAdapterRecoveryConformance|P23Shared)'"
    ]
    for i,c in enumerate(cmds,1):
        r=run(c)
        (OUT/f'supplemental-{i}.log').write_text('\n'.join(line.rstrip() for line in (r['stdout']+r['stderr']).splitlines())+'\n')
        save(f'supplemental-{i}.json',r)
        print('S'+str(i),r['exit_code'],r['elapsed_seconds'],r['stdout'].splitlines()[-1] if r['stdout'] else r['stderr'],flush=True)

def schema():
    scratch=pathlib.Path(tempfile.mkdtemp(prefix='p24-final-schema-'))
    attempts=[]
    def sql(text):
        p=subprocess.run(['sqlite3','-json',str(scratch/'chain.db')],input='PRAGMA foreign_keys=ON;\n'+text,capture_output=True,text=True)
        attempts.append(dict(sql=text,exit_code=p.returncode,stdout=p.stdout,stderr=p.stderr))
        assert p.returncode==0,p.stderr
        return p.stdout
    paths=sorted((ROOT/'migrations').glob('*.up.sql'))
    assert len(paths)==19, len(paths)
    for path in paths[:18]: sql('BEGIN;\n'+path.read_text()+'\nCOMMIT;')
    def dump(): return sql('SELECT type,name,tbl_name,sql FROM sqlite_master ORDER BY type,name;')
    before=dump(); (OUT/'schema-18.json').write_text(before)
    sql('BEGIN;\n'+paths[18].read_text()+'\nCOMMIT;')
    after=dump(); (OUT/'schema-19.json').write_text(after)
    old=json.loads(before); new=json.loads(after)
    changed=[(a,b) for a,b in zip(old,new) if a!=b]
    assert len(old)==len(new) and len(changed)==1
    assert changed[0][0]['name']=='idempotency_results'
    column='request_field_digests_json TEXT NOT NULL CHECK (json_valid(request_field_digests_json))'
    assert changed[0][1]['sql']==changed[0][0]['sql'].replace('created_at DATETIME NOT NULL,','created_at DATETIME NOT NULL, '+column+',',1)
    sql('BEGIN;\n'+(ROOT/'migrations/0019_p24_idempotency_fields.down.sql').read_text()+'\nCOMMIT;')
    down=dump(); (OUT/'schema-down-18.json').write_text(down)
    assert down==before
    integrity=sql('PRAGMA integrity_check;')
    foreign=sql('PRAGMA foreign_key_check;')
    assert '"ok"' in integrity and foreign==''
    receipt=run('cmp project-docs/docs/phases/phase-24/qa/final-evidence/schema-18.json project-docs/docs/phases/phase-24/qa/final-evidence/schema-down-18.json')
    save('schema-proof.json',dict(scratch=str(scratch),migrations=[str(p.relative_to(ROOT)) for p in paths],changed=changed,integrity=integrity,foreign_keys=foreign,cmp=receipt,sqlite_attempts=attempts))
    print('schema up: one accepted column; down: byte-equal; integrity: ok; cmp exit',receipt['exit_code'])

if __name__=='__main__':
    {'docs':docs,'inventory':inventory,'supplemental':supplemental,'schema':schema}[sys.argv[1]]()
