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

func Renew(options RenewOptions) (Claim, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return Claim{}, fmt.Errorf("renew claim: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return Claim{}, fmt.Errorf("renew claim: migrations dir is required")
	}
	if strings.TrimSpace(options.IssueID) == "" {
		return Claim{}, fmt.Errorf("renew claim: issue id is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionIssueRenew); err != nil {
		return Claim{}, err
	}

	var result Claim
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(issueClaimRenewArgv(options), " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeClaimsConn(conn)

		event, claim, err := buildRenewEvent(lock, conn, options)
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

		result = claim
		return nil
	})
	if err != nil {
		return Claim{}, err
	}

	return result, nil
}

func buildRenewEvent(lock *journal.WriteLock, conn *sql.DB, options RenewOptions) (journal.Event, Claim, error) {
	if err := assertIssueExists(conn, options.IssueID); err != nil {
		return journal.Event{}, Claim{}, err
	}
	current, err := loadCurrentClaim(conn, options.IssueID)
	if err != nil {
		return journal.Event{}, Claim{}, err
	}
	now := operationTime(options.Now)
	if current == nil || !current.isLiveAt(now) || current.ClaimedBy != options.ActorRole {
		return journal.Event{}, Claim{}, claimNotHeldError(options.IssueID, options.ActorRole, current)
	}

	ttlSeconds := current.TTLSeconds
	if options.TTLSeconds != nil {
		ttlSeconds = *options.TTLSeconds
	}

	renewedAt := nextClaimTime(current.ClaimedAt, now)
	event, err := journal.PrepareEvent(lock, journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   "issue.claim.renewed",
		OccurredAt:  renewedAt,
		Actor: journal.Actor{
			Role: options.ActorRole,
		},
		Command: journal.Command{
			Argv:       issueClaimRenewArgv(options),
			CWD:        ".",
			FIDVersion: fidVersion,
		},
	})
	if err != nil {
		return journal.Event{}, Claim{}, err
	}

	claimedAt := event.OccurredAt.UTC().Format(time.RFC3339)
	activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
		return activityIDExists(conn, candidate)
	})
	if err != nil {
		return journal.Event{}, Claim{}, err
	}

	event.Effects = []journal.Effect{
		{
			Op:    "update",
			Table: "claims",
			PK: map[string]any{
				"ref_type": "issue",
				"ref_id":   options.IssueID,
			},
			Values: map[string]any{
				"claimed_at":  claimedAt,
				"ttl_seconds": ttlSeconds,
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
		},
		{
			Op:    "activity",
			Table: "activity_log",
			PK:    map[string]any{"id": activityID},
			Values: map[string]any{
				"journal_event_id": event.EventID,
				"ref_type":         "issue",
				"ref_id":           options.IssueID,
				"event":            "claim_renewed",
				"details": map[string]any{
					"previous_claimed_at":  current.ClaimedAt.Format(time.RFC3339),
					"previous_expires_at":  current.expiresAt().Format(time.RFC3339),
					"previous_ttl_seconds": current.TTLSeconds,
					"ttl_seconds":          ttlSeconds,
				},
				"actor":       options.ActorRole,
				"occurred_at": claimedAt,
			},
		},
	}

	return event, Claim{
		RefType:    "issue",
		RefID:      options.IssueID,
		ClaimedBy:  options.ActorRole,
		ClaimedAt:  claimedAt,
		TTLSeconds: ttlSeconds,
		ExpiresAt:  event.OccurredAt.UTC().Add(time.Duration(ttlSeconds) * time.Second).Format(time.RFC3339),
		Note:       current.Note,
	}, nil
}

func issueClaimRenewArgv(options RenewOptions) []string {
	argv := []string{"fid", "issue", "claim", options.IssueID, "--by", options.ActorRole, "--renew"}
	if options.TTLSeconds != nil {
		argv = append(argv, "--ttl", fmt.Sprintf("%ds", *options.TTLSeconds))
	}
	return argv
}
