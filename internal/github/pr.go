package github

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
)

const prJSONFields = "number,url,state,mergedAt,mergeCommit,headRefOid,headRefName,baseRefName,baseRefOid,isCrossRepository,body"

type PullRequest struct {
	Number                             int64
	URL, State                         string
	Merged                             bool
	MergedAt                           *time.Time
	MergeCommit                        string
	HeadSHA, HeadRef, BaseRef, BaseSHA string
	BaseRepository                     Repository
	CrossRepository                    bool
	Body                               string
}

type prResponse struct {
	Number      int64      `json:"number"`
	URL         string     `json:"url"`
	State       string     `json:"state"`
	MergedAt    *time.Time `json:"mergedAt"`
	MergeCommit *struct {
		OID string `json:"oid"`
	} `json:"mergeCommit"`
	HeadSHA         string `json:"headRefOid"`
	HeadRef         string `json:"headRefName"`
	BaseRef         string `json:"baseRefName"`
	BaseSHA         string `json:"baseRefOid"`
	CrossRepository bool   `json:"isCrossRepository"`
	Body            string `json:"body"`
}

func (client *Client) PR(ctx context.Context, value string) (PullRequest, error) {
	identity, err := ParsePRURL(value)
	if err != nil {
		return PullRequest{}, err
	}
	if !identity.Repository.Equal(client.repo) {
		return PullRequest{}, client.repoMismatch(identity.Repository)
	}
	data, err := client.request(ctx, []string{"pr", "view", strconv.FormatInt(identity.Number, 10), "--repo", client.repo.String(), "--json", prJSONFields}, true)
	if err != nil {
		return PullRequest{}, err
	}
	return client.decodePR(data, identity)
}

func (client *Client) decodePR(data []byte, identity PRIdentity) (PullRequest, error) {
	var response prResponse
	err := requiredObject(data, &response, "number", "url", "state", "headRefOid", "headRefName", "baseRefName", "baseRefOid", "isCrossRepository", "body")
	if err != nil {
		return PullRequest{}, err
	}
	// gh's canonical PR URL names the base repository, including for fork heads.
	// headRepository and headRepositoryOwner are deliberately never consulted.
	base, err := ParsePRURL(response.URL)
	if err != nil {
		return PullRequest{}, unavailable(MalformedResponse, "invalid PR identity", err)
	}
	if !base.Repository.Equal(client.repo) {
		return PullRequest{}, client.repoMismatch(base.Repository)
	}
	if response.Number != identity.Number || base.Number != identity.Number || !oidPattern.MatchString(response.HeadSHA) || !oidPattern.MatchString(response.BaseSHA) || response.HeadRef == "" || response.BaseRef == "" {
		return PullRequest{}, unavailable(MalformedResponse, "inconsistent PR identity or commit", nil)
	}
	if response.BaseRef != client.branch {
		return PullRequest{}, &MismatchError{Kind: BranchMismatch, Expected: client.branch, Actual: response.BaseRef}
	}
	if err := validateMerge(response, data); err != nil {
		return PullRequest{}, err
	}
	facts := PullRequest{Number: response.Number, URL: response.URL, State: response.State, Merged: response.State == "MERGED", MergedAt: response.MergedAt,
		HeadSHA: response.HeadSHA, HeadRef: response.HeadRef, BaseRef: response.BaseRef, BaseSHA: response.BaseSHA, BaseRepository: base.Repository, CrossRepository: response.CrossRepository, Body: response.Body}
	if response.MergeCommit != nil {
		facts.MergeCommit = response.MergeCommit.OID
	}
	return facts, nil
}

func validateMerge(response prResponse, data []byte) error {
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(data, &fields) // Already decoded by requiredObject.
	if _, ok := fields["mergedAt"]; !ok {
		return unavailable(MalformedResponse, "missing merge time field", nil)
	}
	if _, ok := fields["mergeCommit"]; !ok {
		return unavailable(MalformedResponse, "missing merge commit field", nil)
	}
	switch response.State {
	case "MERGED":
		if response.MergedAt == nil || response.MergedAt.IsZero() || response.MergeCommit == nil || !oidPattern.MatchString(response.MergeCommit.OID) {
			return unavailable(MalformedResponse, "merged PR lacks merge time or commit", nil)
		}
	case "OPEN", "CLOSED":
		if response.MergedAt != nil || response.MergeCommit != nil {
			return unavailable(MalformedResponse, "unmerged PR has merge facts", nil)
		}
	default:
		return unavailable(MalformedResponse, "unknown PR state", nil)
	}
	return nil
}

func (client *Client) repoMismatch(actual Repository) error {
	return &MismatchError{Kind: RepositoryMismatch, Expected: client.repo.String(), Actual: actual.String()}
}
