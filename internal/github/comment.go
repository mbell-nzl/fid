package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var commentFragment = regexp.MustCompile(`^issuecomment-[1-9][0-9]*$`)

func ValidCommentURL(pr, comment string) bool {
	p, err := ParsePRURL(pr)
	if err != nil {
		return false
	}
	u, err := url.Parse(comment)
	if err != nil || u.Scheme != "https" || u.Host != "github.com" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.RawPath != "" || !commentFragment.MatchString(u.Fragment) {
		return false
	}
	base := "/" + p.Repository.String() + "/"
	return strings.EqualFold(u.Path, base+"issues/"+strconv.FormatInt(p.Number, 10)) || strings.EqualFold(u.Path, base+"pull/"+strconv.FormatInt(p.Number, 10))
}

// PostComment reconciles a stable marker before writing, including recovery
// when GitHub accepted a post but the local receipt could not be appended.
// The caller serializes attempts under the journal lock.
func (client *Client) PostComment(ctx context.Context, pr, marker, body string) (string, error) {
	p, err := ParsePRURL(pr)
	if err != nil || !p.Repository.Equal(client.repo) {
		return "", fmt.Errorf("comment PR differs from repository")
	}
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	full := body + "\n" + marker + "\n"
	total := 0
	for page := 1; page <= maxCIPages; page++ {
		data, err := client.api(ctx, fmt.Sprintf("issues/%d/comments?per_page=100&page=%d", p.Number, page), "map({body,html_url})")
		if err != nil {
			return "", err
		}
		total += len(data)
		if total > DefaultOutputLimit {
			return "", fmt.Errorf("comment listing exceeded size limit")
		}
		var comments []struct {
			Body string `json:"body"`
			URL  string `json:"html_url"`
		}
		if err := validateJSON(data); err != nil {
			return "", err
		}
		if err := json.Unmarshal(data, &comments); err != nil || comments == nil {
			return "", fmt.Errorf("invalid comment listing")
		}
		for _, c := range comments {
			if strings.Contains(c.Body, marker) {
				if c.Body != full || !ValidCommentURL(pr, c.URL) {
					return "", fmt.Errorf("existing comment differs from canonical record")
				}
				return c.URL, nil
			}
		}
		if len(comments) < pageSize {
			return client.writeComment(ctx, pr, full)
		}
	}
	return "", fmt.Errorf("comment listing exceeded page limit")
}

func (client *Client) writeComment(ctx context.Context, pr, body string) (string, error) {
	f, err := os.CreateTemp("", "fid-review-comment-*")
	if err != nil {
		return "", err
	}
	defer func() { _ = os.Remove(f.Name()) }()
	_, writeErr := f.WriteString(body)
	closeErr := f.Close()
	if writeErr != nil {
		return "", writeErr
	}
	if closeErr != nil {
		return "", closeErr
	}
	data, err := client.request(ctx, []string{"pr", "comment", pr, "--repo", client.repo.String(), "--body-file", f.Name()}, true)
	if err != nil {
		return "", err
	}
	result := strings.TrimSpace(string(data))
	if !ValidCommentURL(pr, result) {
		return "", fmt.Errorf("invalid posted comment URL")
	}
	return result, nil
}
