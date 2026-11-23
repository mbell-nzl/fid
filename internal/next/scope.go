package next

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

const ReasonOutsideScope = "outside-scope"

// EpicScopeError identifies an existing issue that cannot scope a worker.
type EpicScopeError struct{ ID, Kind string }

func (err *EpicScopeError) Error() string {
	return fmt.Sprintf("scope %s is a %s, not an epic", err.ID, err.Kind)
}

func validateScope(conn reviewrecord.Queryer, epic string) error {
	if epic == "" {
		return nil
	}
	var kind string
	if err := conn.QueryRow(`SELECT kind FROM issues WHERE id=?`, epic).Scan(&kind); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &issues.NotFoundError{IssueID: epic}
		}
		return fmt.Errorf("read epic scope: %w", err)
	}
	if kind != "epic" {
		return &EpicScopeError{ID: epic, Kind: kind}
	}
	return nil
}
