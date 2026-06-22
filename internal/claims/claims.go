package claims

import (
	"fmt"
	"time"
)

const (
	defaultTTLSeconds       = 1800
	actionIssueAcquire      = "issue.claim.acquire"
	actionIssueRenew        = "issue.claim.renew"
	actionIssueForceAcquire = "issue.claim.force_acquire"
	actionIssueRelease      = "issue.claim.release"
	fidVersion              = "0.1.0"
)

type AcquireOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	ActorRole     string
	TTLSeconds    *int
	Note          *string
	Force         bool
}

type Claim struct {
	RefType    string  `json:"ref_type"`
	RefID      string  `json:"ref_id"`
	ClaimedBy  string  `json:"claimed_by"`
	ClaimedAt  string  `json:"claimed_at"`
	TTLSeconds int     `json:"ttl_seconds"`
	ExpiresAt  string  `json:"expires_at"`
	Note       *string `json:"note"`
	Previous   *ClaimPrevious
}

type ClaimPrevious struct {
	ClaimedBy string
	ExpiresAt string
}

type RenewOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	ActorRole     string
	TTLSeconds    *int
}

type ReleaseOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	ActorRole     string
}

type ReleaseResult struct {
	RefType    string
	RefID      string
	ReleasedBy string
}

type LiveClaimConflictError struct {
	RefID     string
	ClaimedBy string
	ExpiresAt string
}

func (err *LiveClaimConflictError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("issue %q already claimed by %q until %s", err.RefID, err.ClaimedBy, err.ExpiresAt)
}

type ClaimNotHeldError struct {
	RefID       string
	RequestedBy string
	ClaimedBy   string
	ExpiresAt   string
}

func (err *ClaimNotHeldError) Error() string {
	if err == nil {
		return "<nil>"
	}
	if err.ClaimedBy == "" {
		return fmt.Sprintf("issue %q has no live claim held by %q", err.RefID, err.RequestedBy)
	}
	return fmt.Sprintf(
		"issue %q live claim not held by %q; held by %q until %s",
		err.RefID,
		err.RequestedBy,
		err.ClaimedBy,
		err.ExpiresAt,
	)
}

type currentClaim struct {
	ClaimedBy  string
	ClaimedAt  time.Time
	TTLSeconds int
	Note       *string
}

func (claim *currentClaim) expiresAt() time.Time {
	return claim.ClaimedAt.Add(time.Duration(claim.TTLSeconds) * time.Second)
}

func (claim *currentClaim) isLiveAt(now time.Time) bool {
	return now.UTC().Before(claim.expiresAt())
}
