# Real gh responses captured by the orchestrator, 2026-11-15 (gh 2.92.0, repo mbell-nzl/fid)

- pr74.json: `gh pr view 74 --json url,number,state,isDraft,mergedAt,mergeCommit,headRefOid,headRefName,baseRefName,baseRefOid,isCrossRepository,headRepository,headRepositoryOwner,closed`
  (a merged PR; mergeCommit.oid is main's merge commit 4517e22d; headRefOid 14b530c9 is the branch tip at merge).
- runs-merge74.json: `gh run list --commit <merge sha> --json databaseId,attempt,headSha,headBranch,event,status,conclusion,workflowName,workflowDatabaseId,url,createdAt,updatedAt`
  (one push run of workflow `ci`, success, 18 s: tree reuse skipped the test jobs and the run still succeeded).
- rest-runs-merge74.json: `gh api 'repos/<o>/<r>/actions/runs?head_sha=<merge sha>&per_page=100'` reduced with --jq to id,
  run_attempt, head_sha, head_branch, event, status, conclusion, name, path, workflow_id, html_url, plus total_count.
- head-commit74.json: `gh api repos/<o>/<r>/git/commits/<head sha>` reduced to sha, tree, parents; the tree matches
  local `git rev-parse <head>^{tree}` (3ceee6a1).
- runs-failed-sample.json: three failed pull_request runs (conclusion "failure").

Fingerprint experiment (for the approval slice, D10 `git patch-id`):
- `git diff $(git merge-base <baseRefOid> <headRefOid>) <headRefOid> | git patch-id --stable` equals
  `gh pr diff <n> | git patch-id --stable` for #67 (5 files), #66 (194 files) and #63 (2 files); `git diff --binary`
  gave the same patch-id too.
- `gh pr diff` FAILS for large PRs: #71 and #72 -> "HTTP 406: Sorry, the diff exceeded the maximum number of files
  (300)". So done can't depend on gh for the diff; compute locally from objects already present (no fetch, lazy
  fetch disabled), and refuse with a clear repair when they're missing.
