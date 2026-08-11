package issues

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mbell-nzl/fid/internal/policy"
)

type normalizedUpdateRequest struct {
	title      *string
	priority   *int
	labels     *[]string
	specMDPath *string
	ownerRole  *string
	ownerActor *string
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
	if options.OwnerRole != nil {
		ownerRole := strings.TrimSpace(*options.OwnerRole)
		if ownerRole == "" {
			return normalizedUpdateRequest{}, fmt.Errorf("update issue: owner_role must not be blank")
		}
		if err := policy.RequireKnownRole(options.WorkspaceRoot, ownerRole); err != nil {
			return normalizedUpdateRequest{}, err
		}
		request.ownerRole = &ownerRole
	}
	if options.OwnerActor != nil {
		ownerActor := strings.TrimSpace(*options.OwnerActor)
		if ownerActor == "" {
			return normalizedUpdateRequest{}, fmt.Errorf("update issue: owner_actor must not be blank")
		}
		request.ownerActor = &ownerActor
	}

	return request, nil
}

func issueUpdateArgv(issueID, actorRole string, request normalizedUpdateRequest) []string {
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
	if request.ownerRole != nil {
		argv = append(argv, "--owner-role", *request.ownerRole)
	}
	if request.ownerActor != nil {
		argv = append(argv, "--owner-actor", *request.ownerActor)
	}
	if request.specMDPath != nil {
		argv = append(argv, "--spec-md", *request.specMDPath)
	}
	argv = append(argv, "--by", actorRole)
	return argv
}
