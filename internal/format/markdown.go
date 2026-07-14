package format

import (
	"fmt"
	"io"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/labels"
	"github.com/mbell-nzl/fid/internal/links"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/stats"
)

type IssueMarkdownView struct {
	Issue issues.Issue
}

type IssueListMarkdownView struct {
	Items []issues.ListItem
}

func writeMarkdown(writer io.Writer, value any) error {
	switch typed := value.(type) {
	case ReadyView:
		return writeReadyMarkdown(writer, typed)
	case *ReadyView:
		if typed == nil {
			return fmt.Errorf("markdown ready view is nil")
		}
		return writeReadyMarkdown(writer, *typed)
	case IssueMarkdownView:
		return writeIssueMarkdown(writer, typed)
	case *IssueMarkdownView:
		if typed == nil {
			return fmt.Errorf("markdown issue view is nil")
		}
		return writeIssueMarkdown(writer, *typed)
	case issues.Issue:
		return writeIssueMarkdown(writer, IssueMarkdownView{Issue: typed})
	case *issues.Issue:
		if typed == nil {
			return fmt.Errorf("markdown issue is nil")
		}
		return writeIssueMarkdown(writer, IssueMarkdownView{Issue: *typed})
	case IssueListMarkdownView:
		return writeIssueListMarkdown(writer, typed)
	case *IssueListMarkdownView:
		if typed == nil {
			return fmt.Errorf("markdown issue list view is nil")
		}
		return writeIssueListMarkdown(writer, *typed)
	case issues.AcceptanceCriteriaList:
		return writeAcceptanceCriteriaListMarkdown(writer, typed)
	case *issues.AcceptanceCriteriaList:
		if typed == nil {
			return fmt.Errorf("markdown acceptance criteria list is nil")
		}
		return writeAcceptanceCriteriaListMarkdown(writer, *typed)
	case []issues.ListItem:
		return writeIssueListMarkdown(writer, IssueListMarkdownView{Items: typed})
	case []issues.ConcernListItem:
		return writeConcernListMarkdown(writer, typed)
	case []claims.ListItem:
		return writeClaimListMarkdown(writer, typed)
	case []issues.IssueGate:
		return writeIssueGateListMarkdown(writer, typed)
	case []labels.Label:
		return writeLabelListMarkdown(writer, typed)
	case []links.ListItem:
		return writeLinkListMarkdown(writer, typed)
	case []policy.MatrixRow:
		return writePolicyMatrixMarkdown(writer, typed)
	case []stats.TokensRow:
		return writeTokensStatsMarkdown(writer, typed)
	case []stats.CostRow:
		return writeCostStatsMarkdown(writer, typed)
	default:
		return fmt.Errorf("unsupported markdown value type %T", value)
	}
}

func markdownPriority(priority *int) string {
	if priority == nil {
		return "(none)"
	}
	return fmt.Sprintf("P%d", *priority)
}
