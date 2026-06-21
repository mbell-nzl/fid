package claims

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

func Release(options ReleaseOptions) (ReleaseResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return ReleaseResult{}, fmt.Errorf("release claim: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return ReleaseResult{}, fmt.Errorf("release claim: migrations dir is required")
	}
	if strings.TrimSpace(options.IssueID) == "" {
		return ReleaseResult{}, fmt.Errorf("release claim: issue id is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionIssueRelease); err != nil {
		return ReleaseResult{}, err
	}

	var result ReleaseResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(issueClaimReleaseArgv(options), " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeClaimsConn(conn)

		event, releaseResult, err := buildReleaseEvent(lock, conn, options)
		if err != nil {
			return err
		}
		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}
		if err := projectClaimsEvent(options.WorkspaceRoot, conn, appended); err != nil {
			return err
		}

		result = releaseResult
		return nil
	})
	if err != nil {
		return ReleaseResult{}, err
	}

	return result, nil
}

func buildReleaseEvent(lock *journal.WriteLock, conn *sql.DB, options ReleaseOptions) (journal.Event, ReleaseResult, error) {
	if err := assertIssueExists(conn, options.IssueID); err != nil {
		return journal.Event{}, ReleaseResult{}, err
	}
	releaseEffect, err := BuildIssueReleaseEffect(conn, options.IssueID, options.ActorRole)
	if err != nil {
		return journal.Event{}, ReleaseResult{}, err
	}

	event, err := journal.PrepareEvent(lock, journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   "issue.claim.released",
		Actor: journal.Actor{
			Role: options.ActorRole,
		},
		Command: journal.Command{
			Argv:       issueClaimReleaseArgv(options),
			CWD:        ".",
			FIDVersion: fidVersion,
		},
	})
	if err != nil {
		return journal.Event{}, ReleaseResult{}, err
	}

	event.Effects = []journal.Effect{releaseEffect}

	return event, ReleaseResult{
		RefType:    "issue",
		RefID:      options.IssueID,
		ReleasedBy: options.ActorRole,
	}, nil
}

func BuildIssueReleaseEffect(conn *sql.DB, issueID, actorRole string) (journal.Effect, error) {
	current, err := requireLiveIssueClaimHeld(conn, issueID, actorRole)
	if err != nil {
		return journal.Effect{}, err
	}
	return buildIssueReleaseEffect(issueID, current), nil
}

func MaybeBuildIssueReleaseEffect(conn *sql.DB, issueID, actorRole string) (*journal.Effect, error) {
	current, err := loadCurrentClaim(conn, issueID)
	if err != nil {
		return nil, err
	}
	if current == nil || !current.isLiveAt(time.Now().UTC()) || current.ClaimedBy != actorRole {
		return nil, nil
	}
	effect := buildIssueReleaseEffect(issueID, current)
	return &effect, nil
}

func buildIssueReleaseEffect(issueID string, current *currentClaim) journal.Effect {
	return journal.Effect{
		Op:    "delete",
		Table: "claims",
		PK: map[string]any{
			"ref_type": "issue",
			"ref_id":   issueID,
		},
		Prior: map[string]any{
			"claimed_by":  current.ClaimedBy,
			"claimed_at":  current.ClaimedAt.Format(time.RFC3339),
			"ttl_seconds": current.TTLSeconds,
			"note":        dereferenceString(current.Note),
		},
		Preconditions: map[string]any{
			"claimed_by":  current.ClaimedBy,
			"claimed_at":  current.ClaimedAt.Format(time.RFC3339),
			"ttl_seconds": current.TTLSeconds,
			"note":        dereferenceString(current.Note),
		},
	}
}

func issueClaimReleaseArgv(options ReleaseOptions) []string {
	return []string{"fid", "issue", "unclaim", options.IssueID, "--by", options.ActorRole}
}
