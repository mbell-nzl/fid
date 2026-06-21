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

func Acquire(options AcquireOptions) (Claim, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return Claim{}, fmt.Errorf("acquire claim: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return Claim{}, fmt.Errorf("acquire claim: migrations dir is required")
	}
	if strings.TrimSpace(options.IssueID) == "" {
		return Claim{}, fmt.Errorf("acquire claim: issue id is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, acquireAction(options)); err != nil {
		return Claim{}, err
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return Claim{}, err
	}

	var result Claim
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(issueClaimArgv(options), " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeClaimsConn(conn)

		event, claim, err := buildAcquireEvent(lock, conn, options)
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

func buildAcquireEvent(lock *journal.WriteLock, conn *sql.DB, options AcquireOptions) (journal.Event, Claim, error) {
	ttlSeconds := defaultTTLSeconds
	if options.TTLSeconds != nil {
		ttlSeconds = *options.TTLSeconds
	}
	note := normalizeNote(options.Note)
	if err := assertIssueExists(conn, options.IssueID); err != nil {
		return journal.Event{}, Claim{}, err
	}
	current, err := loadCurrentClaim(conn, options.IssueID)
	if err != nil {
		return journal.Event{}, Claim{}, err
	}
	if current != nil && current.isLiveAt(time.Now().UTC()) && !options.Force {
		return journal.Event{}, Claim{}, &LiveClaimConflictError{
			RefID:     options.IssueID,
			ClaimedBy: current.ClaimedBy,
			ExpiresAt: current.expiresAt().Format(time.RFC3339),
		}
	}

	event, err := journal.PrepareEvent(lock, journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   "issue.claim.acquired",
		Actor: journal.Actor{
			Role: options.ActorRole,
		},
		Command: journal.Command{
			Argv:       issueClaimArgv(options),
			CWD:        ".",
			FIDVersion: fidVersion,
		},
	})
	if err != nil {
		return journal.Event{}, Claim{}, err
	}

	claimedAt := event.OccurredAt.UTC().Format(time.RFC3339)
	values := map[string]any{
		"claimed_by":  options.ActorRole,
		"claimed_at":  claimedAt,
		"ttl_seconds": ttlSeconds,
	}
	if note != nil {
		values["note"] = *note
	}

	claim := Claim{
		RefType:    "issue",
		RefID:      options.IssueID,
		ClaimedBy:  options.ActorRole,
		ClaimedAt:  claimedAt,
		TTLSeconds: ttlSeconds,
		ExpiresAt:  event.OccurredAt.UTC().Add(time.Duration(ttlSeconds) * time.Second).Format(time.RFC3339),
		Note:       note,
	}

	if current == nil {
		event.Effects = []journal.Effect{{
			Op:    "insert",
			Table: "claims",
			PK: map[string]any{
				"ref_type": "issue",
				"ref_id":   options.IssueID,
			},
			Values: values,
		}}
		return event, claim, nil
	}

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
			Values: values,
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
				"event":            "claim_superseded",
				"details": map[string]any{
					"previous_claimed_by": current.ClaimedBy,
					"previous_expires_at": current.expiresAt().Format(time.RFC3339),
					"reason":              takeoverReason(current, options),
				},
				"actor":       options.ActorRole,
				"occurred_at": claimedAt,
			},
		},
	}
	return event, claim, nil
}

func issueClaimArgv(options AcquireOptions) []string {
	argv := []string{"fid", "issue", "claim", options.IssueID, "--by", options.ActorRole}
	if options.TTLSeconds != nil {
		argv = append(argv, "--ttl", fmt.Sprintf("%ds", *options.TTLSeconds))
	}
	if options.Force {
		argv = append(argv, "--force")
	}
	if note := normalizeNote(options.Note); note != nil {
		argv = append(argv, "--note", *note)
	}
	return argv
}

func acquireAction(options AcquireOptions) string {
	if options.Force {
		return actionIssueForceAcquire
	}
	return actionIssueAcquire
}
