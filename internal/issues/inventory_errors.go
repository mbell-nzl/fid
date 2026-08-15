package issues

import "fmt"

type LabelNotFoundError struct{ Label string }

func (err *LabelNotFoundError) Error() string { return fmt.Sprintf("label %q not found", err.Label) }

type DuplicateIssueLabelError struct {
	IssueID string
	Label   string
}

type DuplicateAcceptanceCriterionError struct {
	IssueID string
	Ordinal int
}

func (err *DuplicateAcceptanceCriterionError) Error() string {
	return fmt.Sprintf("issue %q already has acceptance criterion ordinal %d", err.IssueID, err.Ordinal)
}

func (err *DuplicateIssueLabelError) Error() string {
	return fmt.Sprintf("issue %q already has label %q", err.IssueID, err.Label)
}
