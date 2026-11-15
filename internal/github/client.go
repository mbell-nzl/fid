package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Client is immutable after construction and supports concurrent reads. The
// default branch must come from local refs/remotes/origin/HEAD, as in policy
// discovery; this adapter never discovers it over the network or invokes git.
type Client struct {
	repo        Repository
	branch      string
	run         Runner
	environment []string
}

func New(repo Repository, defaultBranch string, run Runner) (*Client, error) {
	parsed, err := ParseOrigin("https://github.com/" + repo.String())
	if err != nil || !parsed.Equal(repo) || defaultBranch == "" {
		return nil, fmt.Errorf("GitHub reader needs origin repository and its locally known default branch")
	}
	if run == nil {
		run = RunGH
	}
	environment := ghEnvironment(os.Environ())
	filtered := environment[:0]
	for _, entry := range environment {
		if !strings.HasPrefix(entry, "GH_REPO=") {
			filtered = append(filtered, entry)
		}
	}
	return &Client{repo: parsed, branch: defaultBranch, run: run, environment: append(filtered, "GH_REPO="+parsed.String())}, nil
}

var oidPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)

func (client *Client) request(ctx context.Context, args []string, pr bool) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	result := client.run(ctx, args, append([]string(nil), client.environment...))
	switch {
	case errors.Is(result.Err, ErrOutputLimit) || len(result.Stdout) > DefaultOutputLimit-len(result.Stderr):
		return nil, unavailable(OversizedOutput, "gh output exceeded size limit", result.Err)
	case errors.Is(result.Err, context.DeadlineExceeded) || errors.Is(result.Err, context.Canceled) || ctx.Err() != nil:
		return nil, unavailable(TimedOut, "gh request timed out or was cancelled", result.Err)
	case errors.Is(result.Err, exec.ErrNotFound):
		return nil, unavailable(MissingGH, "gh executable was not found", result.Err)
	case result.Err != nil || result.ExitCode != 0:
		return nil, commandFailure(result, pr)
	}
	return result.Stdout, nil
}

// requiredObject rejects absent/null required fields as well as malformed or
// trailing JSON. Unknown keys are allowed: GitHub adds response fields.
func requiredObject(data []byte, out any, required ...string) error {
	var fields map[string]json.RawMessage
	if err := validateJSON(data); err != nil {
		return unavailable(MalformedResponse, "invalid or ambiguous JSON response", err)
	}
	if err := json.Unmarshal(data, &fields); err != nil {
		return unavailable(MalformedResponse, "invalid JSON response", err)
	}
	for _, name := range required {
		if value, ok := fields[name]; !ok || string(value) == "null" {
			return unavailable(MalformedResponse, "missing response field "+name, nil)
		}
	}
	if err := json.Unmarshal(data, out); err != nil {
		return unavailable(MalformedResponse, "invalid response fields", err)
	}
	return nil
}
