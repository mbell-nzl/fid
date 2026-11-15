package donework

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/github"
)

func localTarget(ctx context.Context, root string) (doneeval.Target, error) {
	git := func(args ...string) (string, error) {
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=1", "GIT_TERMINAL_PROMPT=0")
		out, err := cmd.Output()
		return strings.TrimSpace(string(out)), err
	}
	origin, err := git("remote", "get-url", "origin")
	if err != nil {
		return doneeval.Target{}, refusal(doneeval.PRGate, string(doneeval.RepositoryMismatch), nil)
	}
	repo, err := github.ParseOrigin(origin)
	if err != nil {
		return doneeval.Target{}, refusal(doneeval.PRGate, string(doneeval.RepositoryMismatch), nil)
	}
	ref, err := git("symbolic-ref", "--quiet", "refs/remotes/origin/HEAD")
	const prefix = "refs/remotes/origin/"
	if err != nil || !strings.HasPrefix(ref, prefix) || ref == prefix+"HEAD" {
		return doneeval.Target{}, refusal(doneeval.PRGate, string(doneeval.BranchMismatch), nil)
	}
	return doneeval.Target{Repository: repo, DefaultBranch: strings.TrimPrefix(ref, prefix)}, nil
}

func observe(ctx context.Context, o Options, target doneeval.Target, prURL string) (doneeval.GitHubFacts, error) {
	var facts doneeval.GitHubFacts
	client, err := github.New(target.Repository, target.DefaultBranch, o.Runner)
	if err != nil {
		return facts, err
	}
	pr, err := client.PR(ctx, prURL)
	if err != nil {
		return facts, err
	}
	facts.PR = &pr
	if !pr.Merged {
		return facts, nil
	}
	facts.HeadTree, err = client.CommitTree(ctx, pr.HeadSHA)
	if err != nil {
		return facts, err
	}
	ci, err := client.MergeCommitCI(ctx, pr.MergeCommit)
	if err != nil {
		return facts, err
	}
	facts.CI = &ci
	return facts, nil
}
