// Package github reads GitHub facts through gh. It never changes git or GitHub
// state. Callers supply origin's identity and locally known default branch.
package github

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/handoffrecord"
)

type Repository struct{ Owner, Name string }

func (repo Repository) String() string { return repo.Owner + "/" + repo.Name }
func (repo Repository) Equal(other Repository) bool {
	return strings.EqualFold(repo.Owner, other.Owner) && strings.EqualFold(repo.Name, other.Name)
}

type PRIdentity struct {
	Repository Repository
	Number     int64
}

// ParsePRURL shares hand-off's URL grammar, allowing GitHub's case-insensitive
// hostname. Repository identity is normalized; raw PR body and URL are not.
func ParsePRURL(value string) (PRIdentity, error) {
	parsed, err := url.Parse(value)
	if err != nil {
		return PRIdentity{}, fmt.Errorf("parse PR URL: %w", err)
	}
	parsed.Host = strings.ToLower(parsed.Host)
	if !handoffrecord.ValidPR(parsed.String()) {
		return PRIdentity{}, fmt.Errorf("invalid GitHub PR URL")
	}
	parts := strings.Split(strings.TrimPrefix(parsed.Path, "/"), "/")
	number, err := strconv.ParseInt(handoffrecord.PRNumber(parsed.Path), 10, 64)
	if err != nil {
		return PRIdentity{}, fmt.Errorf("parse PR number: %w", err)
	}
	return PRIdentity{Repository: normalizedRepository(parts[0], parts[1]), Number: number}, nil
}

// ParseOrigin accepts GitHub HTTPS, SSH URL and scp-style SSH remotes.
// Local paths, other hosts and ambiguous URL components are refused.
func ParseOrigin(value string) (Repository, error) {
	if strings.HasPrefix(strings.ToLower(value), "git@github.com:") {
		value = "ssh://git@github.com/" + value[len("git@github.com:"):]
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return Repository{}, fmt.Errorf("parse origin: %w", err)
	}
	if (parsed.Scheme != "https" && parsed.Scheme != "ssh") || !strings.EqualFold(parsed.Host, "github.com") || parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" || parsed.RawPath != "" {
		return Repository{}, fmt.Errorf("origin must name a GitHub HTTPS or SSH repository")
	}
	path := strings.TrimPrefix(parsed.Path, "/")
	if strings.HasSuffix(strings.ToLower(path), ".git") {
		path = path[:len(path)-4]
	}
	identity, err := ParsePRURL("https://github.com/" + path + "/pull/1")
	if err != nil || identity.Repository.Name == "." || identity.Repository.Name == ".." {
		return Repository{}, fmt.Errorf("invalid GitHub origin repository")
	}
	return identity.Repository, nil
}

func normalizedRepository(owner, name string) Repository {
	return Repository{Owner: strings.ToLower(owner), Name: strings.ToLower(name)}
}
