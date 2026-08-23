package issues

import (
	"fmt"
	"sort"
	"strings"
)

type normalizedUpdateRequest struct {
	title      *string
	priority   *int
	labels     *[]string
	specMDPath *string
}

func normalizeUpdateRequest(options UpdateOptions) (normalizedUpdateRequest, error) {
	request := normalizedUpdateRequest{
		priority: options.Priority,
	}

	if options.Title != nil {
		title := strings.TrimSpace(*options.Title)
		if title == "" {
			return normalizedUpdateRequest{}, fmt.Errorf("update issue: title must not be blank")
		}
		request.title = &title
	}
	if options.Priority != nil && (*options.Priority < 0 || *options.Priority > 5) {
		return normalizedUpdateRequest{}, fmt.Errorf("update issue: priority must be between 0 and 5")
	}
	if options.Labels != nil {
		labels, err := normalizedLabels(*options.Labels)
		if err != nil {
			return normalizedUpdateRequest{}, fmt.Errorf("%s", strings.Replace(err.Error(), "create issue:", "update issue:", 1))
		}
		sort.Strings(labels)
		request.labels = &labels
	}
	if options.SpecMDPath != nil {
		specMDPath := strings.TrimSpace(*options.SpecMDPath)
		if specMDPath == "" {
			return normalizedUpdateRequest{}, fmt.Errorf("update issue: spec_md_path must not be blank")
		}
		request.specMDPath = &specMDPath
	}
	return request, nil
}

func issueUpdateArgv(issueID, actor, idempotencyKey string, request normalizedUpdateRequest, expectedRevision *int) []string {
	argv := []string{"fid", "issue", "update", issueID}
	if request.title != nil {
		argv = append(argv, "--title", *request.title)
	}
	if request.priority != nil {
		argv = append(argv, "--priority", fmt.Sprintf("%d", *request.priority))
	}
	if request.labels != nil {
		if len(*request.labels) == 0 {
			argv = append(argv, "--labels", "")
		}
		for _, label := range *request.labels {
			argv = append(argv, "--labels", label)
		}
	}
	if request.specMDPath != nil {
		argv = append(argv, "--spec-md", *request.specMDPath)
	}
	if expectedRevision != nil {
		argv = append(argv, "--expected-revision", fmt.Sprintf("%d", *expectedRevision))
	}
	argv = append(argv, "--actor", actor, "--idempotency-key", idempotencyKey)
	return argv
}
