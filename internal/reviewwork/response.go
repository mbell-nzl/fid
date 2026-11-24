package reviewwork

import (
	"database/sql"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

type ResponseOptions struct {
	WorkspaceRoot                               string
	MigrationsFS                                fs.FS
	FindingID, Actor, Text, Title, EpicID, Lane string
	Defer                                       bool
	Clock                                       func() time.Time
}

type ResponseResult struct {
	Response reviewrecord.Response `json:"response"`
	EventID  string                `json:"-"`
}

func Respond(o ResponseOptions) (ResponseResult, error) {
	var result ResponseResult
	err := journal.WithLock(o.WorkspaceRoot, journal.LockOwner{Command: "fid finding response"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, o.MigrationsFS); err != nil {
			return err
		}
		conn, err := OpenProjection(o.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		task, err := reviewrecord.FindingTask(conn, o.FindingID)
		if err != nil {
			return err
		}
		id, err := responseID(conn, "FRS", "finding_responses")
		if err != nil {
			return err
		}
		now := time.Now()
		if o.Clock != nil {
			now = o.Clock()
		}
		now = now.UTC().Truncate(time.Second)
		h, err := reviewrecord.LatestHandoff(conn, task)
		if err != nil {
			return err
		}
		r := reviewrecord.Response{HandoffID: h.ID, ID: id, FindingID: o.FindingID, TaskID: task, Actor: o.Actor, Text: o.Text, At: now.Format(time.RFC3339)}
		verb, family := "answer", "answered"
		if o.Defer {
			verb, family = "defer", "deferred"
			cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.ConfigFileName))
			if err != nil {
				return err
			}
			newTask, err := responseID(conn, cfg.IssuePrefix, "issues")
			if err != nil {
				return err
			}
			var epic sql.NullString
			if err := conn.QueryRow(`SELECT parent_id FROM issues WHERE id=?`, task).Scan(&epic); err != nil {
				return err
			}
			if o.EpicID != "" {
				epic = sql.NullString{String: o.EpicID, Valid: true}
			}
			var parent *string
			if epic.Valid {
				parent = &epic.String
			}
			lane := o.Lane
			if lane == "" {
				lane = "standard"
			}
			r.DeferredTask = &reviewrecord.DeferredTask{ID: newTask, Title: o.Title, EpicID: parent, Lane: lane}
			f, err := reviewrecord.ReadFinding(conn, o.FindingID)
			if err != nil {
				return err
			}
			floor, err := reviewrecord.DeferralLane(conn, f, r.HandoffID)
			if err != nil {
				return err
			}
			if floor == "assured" && lane == "standard" {
				r.DeferredTask.Lane = "assured"
			}
			r.Text = reviewrecord.DeferralCriterion(f)
		}
		effects, err := journal.FindingResponseEffects(conn, r)
		if err != nil {
			return err
		}
		kind, actor, _ := strings.Cut(o.Actor, ":")
		e, err := journal.CommitEvent(lock, conn, journal.Event{EventSchema: journal.EventSchemaV2, EventType: "finding." + family, CommandName: "finding." + verb, OccurredAt: now, ActorV2: journal.EventActor{ActorType: kind, ActorID: actor, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: id, Type: "finding-response", Sequence: 1}, Effects: effects})
		result = ResponseResult{r, e.EventID}
		return err
	})
	return result, err
}

func responseID(conn *sql.DB, prefix, table string) (string, error) {
	return journal.GenerateID(prefix, func(id string) (bool, error) {
		var n int
		err := conn.QueryRow(fmt.Sprintf(`SELECT count(*) FROM %s WHERE id=?`, table), id).Scan(&n)
		return n != 0, err
	})
}
