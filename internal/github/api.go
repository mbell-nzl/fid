package github

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	runsProjection       = "{total_count, runs: [.workflow_runs[] | {id, run_attempt, workflow_id, name, path, event, head_sha, status, conclusion, html_url}]}"
	commitTreeProjection = "{sha, tree: .tree.sha}"
	pageSize             = 100
	maxCIPages           = 100
)

// api names both host and repository explicitly. gh api has no --repo flag;
// it uses a literal repos/owner/name path rather than inferred placeholders.
func (client *Client) api(ctx context.Context, path, projection string) ([]byte, error) {
	return client.request(ctx, []string{"api", "--hostname", "github.com", "--method", "GET", "repos/" + client.repo.String() + "/" + path, "--jq", projection}, false)
}

// CommitTree validates the returned commit identity as well as its tree. It
// neither requires a local object nor attempts to fetch one.
func (client *Client) CommitTree(ctx context.Context, sha string) (string, error) {
	if !oidPattern.MatchString(sha) {
		return "", fmt.Errorf("invalid commit SHA")
	}
	data, err := client.api(ctx, "git/commits/"+sha, commitTreeProjection)
	if err != nil {
		return "", err
	}
	var response struct {
		SHA  string `json:"sha"`
		Tree string `json:"tree"`
	}
	if err := requiredObject(data, &response, "sha", "tree"); err != nil {
		return "", err
	}
	if response.SHA != sha || !oidPattern.MatchString(response.Tree) {
		return "", unavailable(MalformedResponse, "commit identity or tree does not match", nil)
	}
	return response.Tree, nil
}

// MergeCommitCI reads every page before evaluating. A short, duplicate, changed
// or oversized listing cannot produce even a partially green observation.
// The timeout and output budget cover the whole observation, not just a page.
func (client *Client) MergeCommitCI(ctx context.Context, sha string) (CI, error) {
	if !oidPattern.MatchString(sha) {
		return CI{}, fmt.Errorf("invalid commit SHA")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	var observed []WorkflowRun
	total, bytes := -1, 0
	seen := make(map[[2]int64]bool)
	identities := make(map[int64]WorkflowRun)
	for page := 1; page <= maxCIPages; page++ {
		path := "actions/runs?head_sha=" + sha + "&per_page=" + strconv.Itoa(pageSize) + "&page=" + strconv.Itoa(page)
		data, err := client.api(ctx, path, runsProjection)
		if err != nil {
			return CI{}, err
		}
		bytes += len(data)
		if bytes > DefaultOutputLimit {
			return CI{}, unavailable(OversizedOutput, "CI response exceeded size limit", nil)
		}
		response, err := client.decodeCIPage(data)
		if err != nil {
			return CI{}, err
		}
		if total < 0 {
			total = response.Total
		}
		if response.Total != total || len(response.Runs) != min(pageSize, total-len(observed)) {
			return CI{}, unavailable(MalformedResponse, "incomplete or changing CI pagination", nil)
		}
		for _, run := range response.Runs {
			key := [2]int64{run.ID, run.Attempt}
			previous, exists := identities[run.ID]
			if seen[key] || exists && (previous.WorkflowID != run.WorkflowID || previous.HeadSHA != run.HeadSHA || previous.Event != run.Event) {
				return CI{}, unavailable(MalformedResponse, "duplicate or inconsistent CI identity", nil)
			}
			seen[key] = true
			identities[run.ID] = run
			observed = append(observed, run)
		}
		if len(observed) == total {
			return EvaluateCI(sha, observed), nil
		}
	}
	return CI{}, unavailable(OversizedOutput, "CI pagination exceeded page limit", nil)
}

type ciPageResponse struct {
	Total int           `json:"total_count"`
	Runs  []WorkflowRun `json:"runs"`
}

func (client *Client) decodeCIPage(data []byte) (ciPageResponse, error) {
	var response ciPageResponse
	if err := requiredObject(data, &response, "total_count", "runs"); err != nil {
		return response, err
	}
	if response.Total < 0 {
		return response, unavailable(MalformedResponse, "invalid CI total", nil)
	}
	var raw struct {
		Runs []json.RawMessage `json:"runs"`
	}
	_ = json.Unmarshal(data, &raw)
	for i, run := range response.Runs {
		var checked WorkflowRun
		if err := requiredObject(raw.Runs[i], &checked, "id", "run_attempt", "workflow_id", "name", "path", "event", "head_sha", "status", "html_url"); err != nil {
			return response, err
		}
		if !validWorkflow(run) || !strings.EqualFold(run.URL, "https://github.com/"+client.repo.String()+"/actions/runs/"+strconv.FormatInt(run.ID, 10)) {
			return response, unavailable(MalformedResponse, "incomplete or inconsistent CI workflow", nil)
		}
	}
	return response, nil
}

func validWorkflow(run WorkflowRun) bool {
	if run.ID < 1 || run.Attempt < 1 || run.WorkflowID < 1 || run.Name == "" || run.Path == "" || run.Event == "" || run.URL == "" || !oidPattern.MatchString(run.HeadSHA) {
		return false
	}
	if run.Status == "completed" {
		return run.Conclusion != ""
	}
	return run.Conclusion == "" && slices.Contains([]string{"queued", "in_progress", "waiting", "pending", "requested"}, run.Status)
}
