package github

import (
	"fmt"
	"strings"
)

type UnavailableReason string

const (
	MissingGH         UnavailableReason = "gh-missing"
	Unauthenticated   UnavailableReason = "not-authenticated"
	APIError          UnavailableReason = "network-or-api"
	PRNotFound        UnavailableReason = "pr-not-found"
	RateLimited       UnavailableReason = "rate-limited"
	MalformedResponse UnavailableReason = "malformed-response"
	OversizedOutput   UnavailableReason = "oversized-output"
	TimedOut          UnavailableReason = "timeout"
)

// UnavailableError never implies that a merge was observed. Durable knowledge
// of an earlier merge belongs to the caller, not to an unsuccessful query.
type UnavailableError struct {
	Reason UnavailableReason
	Detail string
	Cause  error
}

func (err *UnavailableError) Error() string { return "GitHub unavailable: " + err.Detail }
func (err *UnavailableError) Unwrap() error { return err.Cause }

func unavailable(reason UnavailableReason, detail string, cause error) error {
	return &UnavailableError{Reason: reason, Detail: detail, Cause: cause}
}

// gh has no structured error channel. Facts are decoded only from successful
// JSON; its exit code and bounded diagnostic classify unsuccessful requests.
func commandFailure(result Result, pr bool) error {
	diagnostic := strings.ToLower(string(result.Stderr))
	reason, detail := APIError, "gh network or API request failed"
	switch {
	case result.ExitCode == 4 || strings.Contains(diagnostic, "http 401") || strings.Contains(diagnostic, "gh auth login"):
		reason, detail = Unauthenticated, "gh authentication is required"
	case strings.Contains(diagnostic, "rate limit") || strings.Contains(diagnostic, "http 429"):
		reason, detail = RateLimited, "GitHub API rate limit reached"
	case pr && (strings.Contains(diagnostic, "http 404") || strings.Contains(diagnostic, "could not resolve to a pullrequest") || strings.Contains(diagnostic, "no pull requests found")):
		reason, detail = PRNotFound, "PR was not found or is not accessible"
	}
	return unavailable(reason, detail, result.Err)
}

type MismatchKind string

const (
	RepositoryMismatch MismatchKind = "repository"
	BranchMismatch     MismatchKind = "base-branch"
)

type MismatchError struct {
	Kind             MismatchKind
	Expected, Actual string
}

func (err *MismatchError) Error() string {
	return fmt.Sprintf("GitHub %s mismatch: expected %s, got %s", err.Kind, err.Expected, err.Actual)
}
