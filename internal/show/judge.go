package show

import (
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func (result *Result) readJudges(root string) error {
	if result.Issue.Kind != "task" {
		return nil
	}
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	result.JudgeRounds, err = reviewrecord.JudgeViews(conn, result.Issue.ID)
	return err
}
