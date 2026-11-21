package journal

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/issuestate"
)

// Admission applies one rule set to every Event, whichever path brings it
// into the ledger: a live append (checked before the write), replay
// (catch-up and rebuild) and archive import. This file holds the rules that
// need no projection state, and the state-move rule, which reads only the
// moved Issue's kind and state. The family validators in
// validateCurrentEvent follow them, inside the projection transaction.

// journalLineBudget is the journal's line cap: the aggregate and sequence
// scanners, and archive import, read lines shorter than this, newline
// included. A longer line would stop every later write.
const journalLineBudget = 1 << 20

func withinLineBudget(lineLength int) bool { return lineLength+1 < journalLineBudget }

// ValidateEventShape is the path-independent part of admission, shared by
// projection, replay's pre-scan, export and archive import.
func ValidateEventShape(event Event) error {
	if event.EventSchema != EventSchemaV1 && event.EventSchema != EventSchemaV2 {
		return fmt.Errorf("project event %q: unsupported Event schema %d", event.EventID, event.EventSchema)
	}
	encoded, err := MarshalCanonicalJSON(event)
	if err != nil {
		return err
	}
	return validateEventShape(event, len(encoded))
}

func validateEventShape(event Event, encodedLength int) error {
	if err := ValidateEventSupport(event); err != nil {
		return err
	}
	if err := validateDatedEventID(event.EventID); err != nil {
		return err
	}
	if err := validateLowercaseIdentifiers(event); err != nil {
		return err
	}
	if !withinLineBudget(encodedLength) {
		return fmt.Errorf("%s Event exceeds the journal's 1 MiB line budget", eventSizeLabel(event.EventType))
	}
	if !event.IsV2() {
		return validateLegacyEffects(event)
	}
	return validateFamilyShape(event)
}

// maxEventSequence bounds a day's Event sequence far beyond any real day.
// Each dated ID takes its date's next sequence (DatedEventIDs), so only
// that many real Events in one day reach it.
const maxEventSequence = 999999999

func validateDatedEventID(id string) error {
	_, _, _, err := datedEventSequence(id)
	return err
}

// The writers allocate an Event ID by reading every ID of its day
// (nextSequence): each ID of three dash-separated parts whose middle part is
// that date, whatever its prefix. An ID the scanner cannot parse would stop
// every write on that date, so each ID dated that way must parse as the
// scanner parses it. IDs without a date there are never read: dated is
// false for them.
func datedEventSequence(id string) (date string, sequence int, dated bool, err error) {
	parts := strings.Split(id, "-")
	if len(parts) != 3 {
		return "", 0, false, nil
	}
	if _, err := time.Parse("20060102", parts[1]); err != nil {
		return "", 0, false, nil
	}
	if sequence, err = parseEventSequence(id); err != nil || sequence > maxEventSequence {
		return "", 0, true, fmt.Errorf("Event ID %q needs the form EVT-<date>-<sequence>, with a sequence from 1 to %d", id, maxEventSequence)
	}
	return parts[1], sequence, true, nil
}

// DatedEventIDs admits Event IDs in journal order as the writers allocate
// them: a dated ID takes the next sequence of its date, one past the
// highest before it. A date's IDs then run 1, 2, 3, ..., so the sequence
// the writers allocate next is always the one admission takes next, and an
// admitted ID can't skip ahead to leave them one it refuses.
type DatedEventIDs struct{ highest map[string]int }

func NewDatedEventIDs() *DatedEventIDs { return &DatedEventIDs{highest: map[string]int{}} }

// Admit takes the next Event ID in journal order.
func (ids *DatedEventIDs) Admit(id string) error {
	date, sequence, dated, err := datedEventSequence(id)
	if err != nil || !dated {
		return err
	}
	if next := ids.highest[date] + 1; sequence != next {
		return outOfOrderEventIDError(id, next)
	}
	ids.highest[date] = sequence
	return nil
}

func outOfOrderEventIDError(id string, next int) error {
	return fmt.Errorf("Event ID %q is out of order: its date's next sequence is %d", id, next)
}

// The families that checked their own size before it was shared keep
// their wording.
func eventSizeLabel(family string) string {
	labels := map[string]string{"epic.closed": "epic close", "task.finished": "completion", "task.approved": "approval", "finding.closed": "finding close", "review.round-recorded": "review"}
	if label, found := labels[family]; found {
		return label
	}
	return fmt.Sprintf("%q", family)
}

// SQLite resolves table and column names without regard to ASCII case, so
// "HOLDS" or "LANE" would reach the same storage as "holds" or "lane" while
// slipping past every exact-name rule. fid's schema is all lowercase.
func validateLowercaseIdentifiers(event Event) error {
	for _, effect := range event.Effects {
		if strings.ToLower(effect.Table) != effect.Table {
			return fmt.Errorf("project event %q: table %q must be lowercase", event.EventID, effect.Table)
		}
		for _, row := range []map[string]any{effect.PK, effect.Values, effect.Prior, effect.Preconditions} {
			for column := range row {
				if strings.ToLower(column) != column {
					return fmt.Errorf("project event %q: column %q must be lowercase", event.EventID, column)
				}
			}
		}
	}
	return nil
}

// Event-v1 is Issue history from before the lean core. No writer produces
// it now and import refuses it; replay still reads it, but only on the
// tables of the schema that history wrote (migrations 0001 and 0002), never
// the lean core's or the projection's own.
var legacyIssueTables = map[string]bool{"issues": true, "roles": true, "epic_metadata": true, "requirements": true, "decisions": true, "escalations": true, "triage_events": true, "task_completions": true, "acceptance_criteria": true, "ac_verifications": true, "agent_runs": true, "gate_runs": true, "labels": true, "issue_labels": true, "links": true, "activity_log": true}

func validateLegacyEffects(event Event) error {
	for _, effect := range event.Effects {
		if !legacyIssueTables[effect.Table] {
			return fmt.Errorf("Event-v1 %q may change only Issue tables, not %q", event.EventID, effect.Table)
		}
	}
	return nil
}

// ValidateEventGrammar checks an Event-v2's family, aggregate type and
// ordered Effect grammar: the tables and operations its family owns.
// Archive import reports it apart from the rules that bind the Event's
// subject, actor and retry result.
func ValidateEventGrammar(event Event) error {
	rule, found := retainedCurrentRules[event.EventType]
	if !found {
		return fmt.Errorf("current Event: unsupported Event family %q", event.EventType)
	}
	if err := validateV2Event(event); err != nil {
		return err
	}
	if event.Aggregate.Type != rule.aggregate {
		return fmt.Errorf("current Event: retained aggregate mismatch")
	}
	for _, effect := range event.Effects {
		switch effect.Op {
		case "insert", "update", "delete", "activity":
		default:
			return fmt.Errorf("project event %q: unsupported effect op %q", event.EventID, effect.Op)
		}
	}
	return validateCurrentStages(event, rule)
}

func validateFamilyShape(event Event) error {
	if err := ValidateEventGrammar(event); err != nil {
		return err
	}
	rule := retainedCurrentRules[event.EventType]
	if !strings.Contains(" "+rule.actors+" ", " "+event.ActorV2.ActorType+" ") {
		return fmt.Errorf("current Event %s: %s actor not admitted", event.EventType, event.ActorV2.ActorType)
	}
	// SQL takes any column as a key, so an issues Effect names its Issue by
	// the primary key alone. An Issue Event changes only its own Issue row.
	for _, effect := range event.Effects {
		if effect.Table != "issues" {
			continue
		}
		id, keyed := effect.PK["id"].(string)
		if rule.aggregate == "issue" && (len(effect.PK) != 1 || id != event.Aggregate.ID) {
			return fmt.Errorf("issues Effect must key the Event aggregate %q by id alone", event.Aggregate.ID)
		}
		if len(effect.PK) != 1 || !keyed {
			return fmt.Errorf("issues Effect must key one Issue by id alone")
		}
	}
	// A retry result belongs to the Event's own request, and only then.
	retried := event.Idempotency != nil
	for _, effect := range event.Effects {
		retried = retried || effect.Table == "idempotency_results"
	}
	if retried {
		_, err := validateIdempotencyBinding(event)
		return err
	}
	return nil
}

// validateStateMoves admits a change of an Issue's state only from the Event
// of the command that owns the move (issuestate.MoveOwner): fid start, the
// hand-off, a review round asking for changes, fid done or a person's
// cancellation for a task or an epic, and fid issue transition for a
// subtask's moves. A task or an epic is created in todo. Each Event changes
// an Issue's state at most once.
//
// It runs before the family validators, against the projection as the Event
// finds it. An earlier Effect of the same Event can create the Issue a later
// one moves, so applyEventEffects checks each Effect again,
// against the row it lands on.
//
// Event-v1 history may still create Issues and move subtasks, but it can't
// move a task or an epic: none of the commands that own those moves ever
// wrote Event-v1.
func validateStateMoves(tx *sql.Tx, event Event) error {
	moved := map[string]bool{}
	for _, effect := range event.Effects {
		if _, sets := effect.Values["status"]; effect.Table == "issues" && effect.Op == "update" && sets {
			id, _ := effect.PK["id"].(string)
			if moved[id] {
				return fmt.Errorf("project event %q: a state change keys one Issue by id alone, once per Event", event.EventID)
			}
			moved[id] = true
		}
		if err := validateStateMove(tx, event, effect); err != nil {
			return err
		}
	}
	return nil
}

// validateStateMove checks one Effect's change of an Issue's state against
// the row it changes, as the projection stands: an insert's row merged from
// its key and values, as the projector inserts it, or the row an update's
// id names.
func validateStateMove(tx *sql.Tx, event Event, effect Effect) error {
	to, sets := effect.Values["status"]
	if effect.Table != "issues" || (!sets && effect.Op != "insert") {
		return nil
	}
	move := StateMoveError{EventID: event.EventID, Family: event.EventType}
	switch effect.Op {
	case "insert":
		// Every insert is a creation, with the first state its row holds.
		row, err := mergeEffectValues(effect.PK, effect.Values)
		if err != nil {
			return fmt.Errorf("project event %q: insert issues: %w", event.EventID, err)
		}
		to = row["status"]
		move.IssueID, _ = row["id"].(string)
		move.Kind, _ = row["kind"].(string)
	case "update":
		id, keyed := effect.PK["id"].(string)
		if len(effect.PK) != 1 || !keyed {
			return fmt.Errorf("project event %q: a state change keys one Issue by id alone, once per Event", event.EventID)
		}
		move.IssueID = id
		err := tx.QueryRow(`SELECT kind, status FROM issues WHERE id=?`, id).Scan(&move.Kind, &move.From)
		if errors.Is(err, sql.ErrNoRows) {
			// An earlier Effect of this Event may create the Issue: the move
			// is judged when it lands. An Issue that is never there changes
			// no row, and the write refuses it (applyUpdate).
			return nil
		}
		if err != nil {
			return fmt.Errorf("project event %q: read the Issue whose state changes: %w", event.EventID, err)
		}
	default:
		return fmt.Errorf("project event %q: an Issue's state changes only by insert or update", event.EventID)
	}
	move.To, _ = to.(string)
	owner, owned := issuestate.MoveOwner(move.Kind, move.From, move.To)
	move.Owner = owner
	switch {
	case !event.IsV2() && issuestate.HasFixedStates(move.Kind) && move.From != "":
		move.Legacy = true
	case !owned:
	case move.From == "":
		// A creation needs only its first state here: Event-v2's grammar
		// lets only issue.created insert an Issue, and Event-v1 history
		// named its families freely.
		return nil
	case owner != event.EventType:
	case owner == issuestate.FamilyTransition && issuestate.HasFixedStates(move.Kind) && cancelDecisions(event) != 1:
		move.Undecided = true
	default:
		return nil
	}
	return &move
}

func cancelDecisions(event Event) int {
	n := 0
	for _, effect := range event.Effects {
		if effect.Table == "person_decisions" && effect.Op == "insert" && effect.Values["kind"] == "cancel" {
			n++
		}
	}
	return n
}

// StateMoveError refuses a change of an Issue's state that its Event's
// family doesn't own. Owner is the family that owns the move, if any.
type StateMoveError struct {
	EventID, Family, IssueID, Kind, From, To, Owner string
	// Legacy: Event-v1 moved a task or an epic. Undecided: a task's or an
	// epic's cancellation came without exactly one person's cancel decision.
	Legacy, Undecided bool
}

func (err *StateMoveError) Error() string {
	move := fmt.Sprintf("move %s %s from %s to %s", err.Kind, err.IssueID, err.From, err.To)
	if err.From == "" {
		move = fmt.Sprintf("create %s %s in %s", err.Kind, err.IssueID, err.To)
	}
	why := "no command makes that move"
	switch {
	case issuestate.FinishedCancellation(err.Kind, err.From, err.To):
		why = issuestate.FinishedWorkMessage
	case err.Legacy:
		why = "Event-v1 can't move a task or an epic"
	case err.Undecided:
		why = "it needs exactly one person's cancel decision"
	case err.Owner != "":
		why = "only " + issuestate.OwnerCommand(err.Owner, err.Kind) + " makes that move"
	}
	return fmt.Sprintf("Event %q (%s) can't %s: %s", err.EventID, err.Family, move, why)
}
