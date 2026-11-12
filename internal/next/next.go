// Package next assembles fid next: the order fid walks to hand a worktree
// its next task, why it passes over each task ahead, and, with --claim, the
// hold that takes the task. Text and JSON are both rendered from one
// Result.
package next

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
)

type Options struct {
	WorkspaceRoot string
	MigrationsDir string
	// Clock reads the time, in the zone the text shows. fid next reads it
	// once, just before it walks the order: after refreshing the
	// projection, and with --claim under the write lock, so a hold that
	// lapses while fid waits for the lock is read as lapsed.
	Clock func() time.Time
	// Actor takes the task, as agent:<id> or human:<id>; fid next --claim.
	Actor string
}

// Result is what fid next shows. Hold and Event are set only when --claim
// takes or keeps a hold; the rest serves the text.
type Result struct {
	Kind     string
	Worktree string
	Hold     *holds.Hold
	Task     *issues.OverviewTask
	Queue    []Entry
	Event    *Event
	Now      time.Time
	// Served indexes the task handed out in Queue, or is -1.
	Served int
	// Claimed says the result is fid next --claim's.
	Claimed bool
	// Epic is the served task's epic, if it has one.
	Epic *issues.WorkEpic
	// Interrupted is this worktree's lapsed hold that a resume carries on
	// from.
	Interrupted *Working
}

// Event is the Event a hold was taken in.
type Event struct {
	ID         string
	OccurredAt time.Time
}

// DeniedActorError refuses an actor whose type takes no tasks: service
// actors get nothing for now (ADR-0014, trust by actor type).
type DeniedActorError struct{ Actor string }

func (err *DeniedActorError) Error() string {
	return fmt.Sprintf("actor %s can't take tasks", err.Actor)
}

// CheckActor refuses an actor that may not take a task. Agents and humans
// may.
func CheckActor(actor string) error {
	if strings.HasPrefix(actor, "agent:") || strings.HasPrefix(actor, "human:") {
		return nil
	}
	return &DeniedActorError{Actor: actor}
}

// Read is fid next without --claim: what this worktree would be handed. It
// records nothing.
func Read(options Options) (Result, error) {
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return Result{}, fmt.Errorf("read next: %w", err)
	}
	conn, err := openProjection(options.WorkspaceRoot)
	if err != nil {
		return Result{}, err
	}
	defer func() { _ = conn.Close() }()
	result, _, err := assemble(conn, options, options.Clock())
	return result, err
}

// Claim is fid next --claim: it takes the task Read would hand out as a
// hold, ending this worktree's earlier hold and any lapsed hold on the task,
// in one Event under the write lock. A live hold here is returned as it is,
// with no Event.
func Claim(options Options) (Result, error) {
	if err := CheckActor(options.Actor); err != nil {
		return Result{}, err
	}
	var result Result
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid next --claim --actor " + options.Actor}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openProjection(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		var planned order
		if result, planned, err = assemble(conn, options, options.Clock()); err != nil {
			return err
		}
		result.Claimed = true
		return take(lock, conn, options, &result, planned)
	})
	if err != nil {
		return Result{}, err
	}
	return result, nil
}

func openProjection(root string) (*sql.DB, error) {
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	if err != nil {
		return nil, fmt.Errorf("open projection: %w", err)
	}
	return conn, nil
}

// assemble walks the order at now.
func assemble(conn *sql.DB, options Options, now time.Time) (Result, order, error) {
	work, err := issues.ReadWork(conn)
	if err != nil {
		return Result{}, order{}, fmt.Errorf("read next: %w", err)
	}
	open, err := holds.Open(conn)
	if err != nil {
		return Result{}, order{}, fmt.Errorf("read next: %w", err)
	}
	planned := plan(work, open, project.WorktreeID(options.WorkspaceRoot), now)
	result := Result{Kind: planned.kind, Worktree: options.WorkspaceRoot, Queue: planned.queue, Now: now,
		Served: planned.served}
	if planned.served >= 0 {
		task := planned.queue[planned.served].OverviewTask
		result.Task = &task
		if epic, found := work.Epics[epicOf(task)]; found {
			result.Epic = &epic
		}
	}
	if planned.kind == KindResume && planned.own.Lapsed(now) {
		result.Interrupted = working(*planned.own, planned.own.WorktreeID, now)
	}
	return result, planned, nil
}

// take records the hold on the served task at the time the order was
// walked, unless there is none or this worktree's live hold is it.
func take(lock *journal.WriteLock, conn *sql.DB, options Options, result *Result, planned order) error {
	if planned.kind == KindNone {
		return nil
	}
	if planned.kind == KindResume && !planned.own.Lapsed(result.Now) {
		result.Hold = planned.own
		return nil
	}
	ended := map[string]bool{}
	var replaces *string
	if planned.own != nil {
		ended[planned.own.ID] = true
		if planned.kind == KindResume {
			replaces = &planned.own.ID
		}
	}
	if previous, held := planned.held[result.Task.ID]; held {
		ended[previous.ID] = true
	}
	id, err := holds.NewID(conn)
	if err != nil {
		return err
	}
	hold := holds.New(id, result.Task.ID, project.WorktreeID(options.WorkspaceRoot), options.WorkspaceRoot, options.Actor, result.Now, replaces)
	effects := []journal.Effect{}
	for _, previous := range sortedKeys(ended) {
		effects = append(effects, holds.End(previous, holds.EndReplaced, hold.Since))
	}
	actorType, actorID, _ := strings.Cut(options.Actor, ":")
	appended, err := journal.CommitEvent(lock, conn, journal.Event{EventSchema: journal.EventSchemaV2, EventType: "hold.taken", OccurredAt: hold.Since,
		ActorV2:   journal.EventActor{ActorID: actorID, ActorType: actorType, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"},
		Aggregate: journal.Aggregate{ID: hold.ID, Sequence: 1, Type: "hold"}, CommandName: "next.claim", Effects: append(effects, hold.Insert())})
	if err != nil {
		return fmt.Errorf("take hold: %w", err)
	}
	result.Hold, result.Event = &hold, &Event{ID: appended.EventID, OccurredAt: appended.OccurredAt}
	// The queue shows the holds as the Event left them.
	for index, entry := range result.Queue {
		if entry.Working != nil && ended[entry.Working.HoldID] {
			result.Queue[index].Working = nil
		}
	}
	result.Queue[result.Served].Working = working(hold, hold.WorktreeID, result.Now)
	return nil
}

func sortedKeys(set map[string]bool) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Fields is the result as the next.get or next.claim JSON result. hold is
// set only by --claim; epic is the served task's, with its progress as fid
// status gives it.
func (result Result) Fields() map[string]any {
	var hold, epic any
	if result.Hold != nil {
		hold = map[string]any{"id": result.Hold.ID, "task_id": result.Hold.TaskID, "actor": result.Hold.Actor,
			"since": holds.Timestamp(result.Hold.Since), "expires_at": holds.Timestamp(result.Hold.ExpiresAt), "replaces": result.Hold.Replaces}
	}
	if result.Epic != nil {
		epic = map[string]any{"id": result.Epic.ID, "title": result.Epic.Title, "progress": result.Epic.Progress}
	}
	return map[string]any{"kind": result.Kind, "worktree": result.Worktree, "hold": hold, "task": result.Task, "epic": epic, "queue": result.Queue}
}
