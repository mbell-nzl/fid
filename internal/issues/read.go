package issues

import "fmt"

const ReadCollectionLimit = 10000

type IssueSummary struct {
	ID             string   `json:"id"`
	Kind           string   `json:"kind"`
	Title          string   `json:"title"`
	Status         string   `json:"status"`
	StatusCategory string   `json:"status_category"`
	WorkflowID     string   `json:"workflow_id"`
	ParentID       *string  `json:"parent_id"`
	Priority       *int     `json:"priority"`
	Sequence       *int     `json:"sequence"`
	SpecMDPath     *string  `json:"spec_md_path"`
	Revision       int      `json:"revision"`
	Labels         []string `json:"labels"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type IssueDetail struct {
	IssueSummary
	AcceptanceCriteria []ReadAcceptanceCriterion `json:"acceptance_criteria"`
}

type ReadAcceptanceCriterion struct {
	ID      string `json:"id"`
	Ordinal int    `json:"ordinal"`
	Text    string `json:"text"`
}

type ResultLimitExceededError struct {
	Collection string
	Limit      int
	Observed   int
}

func (err *ResultLimitExceededError) Error() string {
	return fmt.Sprintf("%s contains %d items, exceeding limit %d", err.Collection, err.Observed, err.Limit)
}

func ReadIssueDetail(options GetOptions) (IssueDetail, error) {
	issue, err := Get(options)
	if err != nil {
		return IssueDetail{}, err
	}
	return issueDetailFromIssue(issue), nil
}

func ListIssueSummaries(options ListOptions) ([]IssueSummary, error) {
	items, err := List(options)
	if err != nil {
		return nil, err
	}
	if err := EnforceReadCollectionLimit("issues", len(items)); err != nil {
		return nil, err
	}
	summaries := make([]IssueSummary, 0, len(items))
	for _, item := range items {
		detail, err := ReadIssueDetail(GetOptions{
			WorkspaceRoot: options.WorkspaceRoot,
			MigrationsDir: options.MigrationsDir,
			ID:            item.ID,
		})
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, detail.IssueSummary)
	}
	return summaries, nil
}

func EnforceReadCollectionLimit(collection string, observed int) error {
	if observed > ReadCollectionLimit {
		return &ResultLimitExceededError{Collection: collection, Limit: ReadCollectionLimit, Observed: observed}
	}
	return nil
}

func issueDetailFromIssue(issue Issue) IssueDetail {
	criteria := make([]ReadAcceptanceCriterion, 0, len(issue.AcceptanceCriteria))
	for _, criterion := range issue.AcceptanceCriteria {
		criteria = append(criteria, ReadAcceptanceCriterion{ID: criterion.ID, Ordinal: criterion.Ordinal, Text: criterion.Text})
	}
	return IssueDetail{
		IssueSummary: IssueSummary{
			ID:             issue.ID,
			Kind:           issue.Kind,
			Title:          issue.Title,
			Status:         issue.Status,
			StatusCategory: issue.StatusCategory,
			WorkflowID:     issue.WorkflowID,
			ParentID:       issue.ParentID,
			Priority:       issue.Priority,
			Sequence:       issue.Sequence,
			SpecMDPath:     issue.SpecMDPath,
			Revision:       issue.Revision,
			Labels:         append([]string{}, issue.Labels...),
			CreatedAt:      requiredReadTimestamp(issue.CreatedAt),
			UpdatedAt:      requiredReadTimestamp(issue.UpdatedAt),
		},
		AcceptanceCriteria: criteria,
	}
}

func requiredReadTimestamp(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
