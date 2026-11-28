package epicwork

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/epicrecord"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/project"
)

func git(root string, args ...string) (string, error) {
	c := exec.Command("git", append([]string{"-C", root}, args...)...)
	c.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=1")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func verifiedCommit(q epicplan.Queryer, o Options) (string, []string, error) {
	source, err := policy.DefaultSource(o.WorkspaceRoot)
	if err != nil {
		return "", nil, err
	}
	merges, err := epicplan.Merges(q, o.EpicID)
	if err != nil {
		return "", nil, err
	}
	for _, merge := range merges {
		if _, err := git(o.WorkspaceRoot, "merge-base", "--is-ancestor", merge, source.Commit); err != nil {
			return "", nil, fmt.Errorf("default-branch commit does not contain epic merge %s: %w", merge, err)
		}
	}
	return source.Commit, merges, nil
}

func passEffect(o Options, q *sql.DB, v epicplan.Validation, p epicplan.Pass) (journal.Effect, error) {
	id, err := journal.GenerateID("EVR", func(candidate string) (bool, error) {
		var n int
		err := q.QueryRow(`SELECT count(*) FROM epic_validation_results WHERE id=?`, candidate).Scan(&n)
		return n > 0, err
	})
	if err != nil {
		return journal.Effect{}, err
	}
	p.ID = id
	return journal.Effect{Op: "insert", Table: "epic_validation_results", PK: map[string]any{"id": id}, Values: map[string]any{"epic_id": o.EpicID, "ordinal": v.Ordinal, "generation": v.Generation, "record_json": encode(p)}}, nil
}

func Pass(o Options) (Result, error) {
	o = o.frozen()
	if !policy.Allows(o.Actor, policy.PersonDecisions) {
		return Result{}, &persondecision.DeniedError{Actor: o.Actor}
	}
	if strings.TrimSpace(o.Reason) == "" {
		return Result{}, &persondecision.ReasonRequiredError{}
	}
	return o.write("epic.validation.pass", func(q *sql.DB) ([]journal.Effect, Result, error) {
		p, err := epicplan.Preview(q, o.EpicID)
		if err != nil {
			return nil, Result{}, err
		}
		var v epicplan.Validation
		for _, candidate := range p.Validations {
			if candidate.ID == o.ValidationID {
				v = candidate
			}
		}
		if v.ID == "" {
			return nil, Result{}, fmt.Errorf("validation not found: %s", o.ValidationID)
		}
		if strings.HasPrefix(v.Verifier, "check:") {
			return nil, Result{}, fmt.Errorf("fid must verify check validations")
		}
		commit, merges, err := verifiedCommit(q, o)
		if err != nil {
			return nil, Result{}, err
		}
		if v.Verifier == "human" {
			entry, err := git(o.WorkspaceRoot, "ls-tree", commit, "--", ":(literal)"+v.Recipe)
			if err != nil {
				return nil, Result{}, err
			}
			if !strings.HasPrefix(entry, "100644 blob ") && !strings.HasPrefix(entry, "100755 blob ") {
				return nil, Result{}, fmt.Errorf("recipe %s must exist as a regular file in %s", v.Recipe, commit)
			}
		}
		id, decision, err := persondecision.Effect(q, "pass-validation", o.Actor, "epic", o.EpicID, o.Reason, o.now())
		if err != nil {
			return nil, Result{}, err
		}
		pass := epicplan.Pass{EpicID: o.EpicID, ValidationID: v.ID, Generation: v.Generation, Commit: commit, Merges: merges, Actor: o.Actor, At: o.now().Format(time.RFC3339), Outcome: "passed", DecisionID: id}
		decision.Values["binding_id"] = v.ID
		decision.Values["binding_json"] = encode(pass)
		e, err := passEffect(o, q, v, pass)
		if err != nil {
			return nil, Result{}, err
		}
		return []journal.Effect{decision, e}, Result{EpicID: o.EpicID, DecisionID: id, At: pass.At}, nil
	})
}

func Verify(ctx context.Context, o Options) (reply Result, returnErr error) {
	if o.IdempotencyKey == "" {
		return Result{}, fmt.Errorf("verify requires an idempotency key")
	}
	if !policy.Allows(o.Actor, policy.WorkSet) {
		return Result{}, &persondecision.DeniedError{Actor: o.Actor}
	}
	q, err := o.open()
	if err != nil {
		return Result{}, err
	}
	defer func() { _ = q.Close() }()
	req, err := o.request("epic.validation.verify")
	if err != nil {
		return Result{}, err
	}
	prior, found, err := journal.ReplayIdempotent(q, req)
	if err != nil {
		return Result{}, err
	}
	if found {
		var r Result
		err = json.Unmarshal([]byte(prior.ResultJSON), &r)
		r.EventID = prior.EventID
		return r, err
	}
	eligible, err := epicrecord.ReadEligibility(q, o.EpicID)
	if err != nil {
		return Result{}, err
	}
	if !eligible.CanClose() {
		return Result{}, fmt.Errorf("verify requires an open epic with no open tasks")
	}
	p, err := epicplan.Preview(q, o.EpicID)
	if err != nil {
		return Result{}, err
	}
	commit, merges, err := verifiedCommit(q, o)
	if err != nil {
		return Result{}, err
	}
	loaded, err := policy.Load(o.WorkspaceRoot)
	if err != nil {
		return Result{}, err
	}
	if loaded.Source.Commit != commit {
		return Result{}, fmt.Errorf("default-branch policy changed; retry verify")
	}
	base, err := os.MkdirTemp("", "fid-epic-verify-")
	if err != nil {
		return Result{}, err
	}
	defer func() { _ = os.RemoveAll(base) }()
	root := filepath.Join(base, "checkout")
	if _, err := git(o.WorkspaceRoot, "-c", "core.hooksPath=/dev/null", "worktree", "add", "--detach", root, commit); err != nil {
		return Result{}, err
	}
	defer func() {
		_, err := git(o.WorkspaceRoot, "worktree", "remove", "--force", root)
		if err != nil {
			returnErr = errors.Join(returnErr, fmt.Errorf("remove verification worktree: %w", err))
		}
	}()
	results := []epicplan.Pass{}
	artifactEffects := []journal.Effect{}
	for _, v := range p.Validations {
		name, ok := strings.CutPrefix(v.Verifier, "check:")
		if !ok {
			continue
		}
		check, ok := loaded.Effective.Checks[name]
		if !ok {
			return Result{}, fmt.Errorf("unknown check in default-branch policy: %s", name)
		}
		definition := checkrecord.DefinitionFor(check)
		before, err := project.ObserveRepository(root)
		if err != nil {
			return Result{}, err
		}
		process, err := checkrun.ExecuteProcess(ctx, checkrun.ProcessOptions{Root: root, Definition: definition, Version: checkrun.Version()})
		if err != nil {
			return Result{}, err
		}
		after, err := project.ObserveRepository(root)
		if err != nil {
			return Result{}, err
		}
		process.Tree = checkrecord.Tree{Head: after.HeadOID, Tree: after.TreeOID, Dirty: after.Dirty, DirtyDigest: after.DirtyDigest}
		process.FinishedAt = o.now().Format(time.RFC3339)
		process.Reasons = checkrecord.Qualification(checkrecord.Tree{Head: before.HeadOID, Tree: before.TreeOID, Dirty: before.Dirty, DirtyDigest: before.DirtyDigest}, process.Tree, process.Outcome, true, true)
		if err := keepArtifacts(o, q, root, &process.Result, &definition, &artifactEffects); err != nil {
			return Result{}, err
		}
		if len(process.RefusedArtifacts) > 0 {
			process.Reasons = append(process.Reasons, "artifact-not-kept")
		}
		process.Counts = process.Outcome == "passed" && len(process.Reasons) == 0
		outcome := "failed"
		if process.Counts {
			outcome = "passed"
		}
		results = append(results, epicplan.Pass{EpicID: o.EpicID, ValidationID: v.ID, Generation: v.Generation, Commit: commit, Merges: merges, Actor: o.Actor, At: o.now().Format(time.RFC3339), Outcome: outcome, Check: &process.Result, Definition: &definition})
	}
	o = o.frozen()
	for i := range results {
		results[i].At = o.now().Format(time.RFC3339)
	}
	return o.write("epic.validation.verify", func(q *sql.DB) ([]journal.Effect, Result, error) {
		current, err := epicplan.Preview(q, o.EpicID)
		if err != nil {
			return nil, Result{}, err
		}
		eligible, err := epicrecord.ReadEligibility(q, o.EpicID)
		if err != nil {
			return nil, Result{}, err
		}
		nowMerges, err := epicplan.Merges(q, o.EpicID)
		if err != nil {
			return nil, Result{}, err
		}
		if !eligible.CanClose() || current.Binding != p.Binding || !slices.Equal(merges, nowMerges) {
			return nil, Result{}, fmt.Errorf("epic changed during verification; retry")
		}
		effects := append([]journal.Effect{}, artifactEffects...)
		for i, pass := range results {
			var v epicplan.Validation
			for _, candidate := range p.Validations {
				if candidate.ID == pass.ValidationID {
					v = candidate
				}
			}
			e, err := passEffect(o, q, v, pass)
			if err != nil {
				return nil, Result{}, err
			}
			pass.ID = e.PK["id"].(string)
			results[i] = pass
			effects = append(effects, e)
		}
		return effects, Result{EpicID: o.EpicID, Checks: results, At: o.now().Format(time.RFC3339)}, nil
	})
}
