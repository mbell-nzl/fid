// Package epicwork writes epic validations and person decisions.
package epicwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

type Options struct {
	RecipeSet                                           bool
	WorkspaceRoot                                       string
	MigrationsFS                                        fs.FS
	EpicID, ValidationID, Actor, Reason, IdempotencyKey string
	Validation                                          epicplan.Validation
	Clock                                               func() time.Time
	Confirm                                             func(epicplan.Plan) (bool, error)
}

type Result = epicplan.Result

func encode(v any) string { raw, _ := json.Marshal(v); return string(raw) }
func (o Options) now() time.Time {
	if o.Clock != nil {
		return o.Clock().UTC().Truncate(time.Second)
	}
	return time.Now().UTC().Truncate(time.Second)
}
func (o Options) request(operation string) (journal.IdempotencyRequest, error) {
	cfg, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.ConfigFileName))
	if err != nil {
		return journal.IdempotencyRequest{}, err
	}
	input := map[string]any{"epic": o.EpicID}
	if operation == "epic.validation.add" || operation == "epic.validation.update" {
		input["claim"] = o.Validation.Claim
		input["verifier"] = o.Validation.Verifier
		if o.RecipeSet || o.Validation.Recipe != "" {
			input["recipe"] = o.Validation.Recipe
		}
	}
	if operation == "epic.validation.update" || operation == "epic.validation.remove" || operation == "epic.validation.pass" {
		input["id"] = o.ValidationID
	}
	if operation == "epic.approve" || operation == "epic.validation.pass" {
		input["reason"] = o.Reason
	}
	digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": o.Actor, "operation": operation, "project_id": cfg.ProjectID, "request": input}, o.Actor, map[string]string{"claim": "--claim", "verifier": "--verifier", "recipe": "--recipe", "reason": "--reason"})
	return journal.IdempotencyRequest{OperationScope: operation + ":" + cfg.ProjectID, IdempotencyKey: o.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}, err
}

func (o Options) open() (*sql.DB, error) {
	if _, err := journal.EnsureProjectionFresh(o.WorkspaceRoot, o.MigrationsFS); err != nil {
		return nil, err
	}
	return reviewwork.OpenProjection(o.WorkspaceRoot)
}
func (o Options) write(operation string, fn func(*sql.DB) ([]journal.Effect, Result, error)) (Result, error) {
	var result Result
	if !policy.Allows(o.Actor, policy.WorkSet) {
		return result, &persondecision.DeniedError{Actor: o.Actor}
	}
	err := journal.WithLock(o.WorkspaceRoot, journal.LockOwner{Command: "fid " + strings.ReplaceAll(operation, ".", " ")}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, o.MigrationsFS); err != nil {
			return err
		}
		conn, err := reviewwork.OpenProjection(o.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		req, err := o.request(operation)
		if err != nil {
			return err
		}
		if o.IdempotencyKey != "" {
			prior, found, err := journal.ReplayIdempotent(conn, req)
			if err != nil {
				return err
			}
			if found {
				err = json.Unmarshal([]byte(prior.ResultJSON), &result)
				result.EventID = prior.EventID
				return err
			}
		}
		effects, r, err := fn(conn)
		if err != nil {
			return err
		}
		result = r
		agg := journal.Aggregate{Type: "epic", ID: o.EpicID}
		agg.Sequence, err = journal.NextAggregateSequence(filepath.Join(config.LedgerDir(o.WorkspaceRoot), config.JournalDirName), agg)
		if err != nil {
			return err
		}
		typ, actor, _ := strings.Cut(o.Actor, ":")
		event := journal.Event{EventSchema: 2, EventType: operation, CommandName: operation, OccurredAt: o.now(), Aggregate: agg, ActorV2: journal.EventActor{ActorType: typ, ActorID: actor, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Effects: effects}
		var appended journal.Event
		if o.IdempotencyKey == "" {
			appended, err = journal.CommitEvent(lock, conn, event)
		} else {
			var outcome journal.IdempotencyOutcome
			outcome, err = journal.CommitIdempotent(lock, conn, req, event, func(journal.Event) (journal.IdempotencyMutation, error) {
				return journal.IdempotencyMutation{EntityID: o.EpicID, ResultJSON: encode(result)}, nil
			})
			appended.EventID = outcome.EventID
		}
		result.EventID = appended.EventID
		return err
	})
	return result, err
}

func Change(o Options, action string) (Result, error) {
	o = o.frozen()
	if o.IdempotencyKey == "" {
		return Result{}, fmt.Errorf("validation changes require an idempotency key")
	}
	return o.write("epic.validation."+action, func(q *sql.DB) ([]journal.Effect, Result, error) {
		p, err := epicplan.Preview(q, o.EpicID)
		if err != nil {
			return nil, Result{}, err
		}
		list, generation, ordinal, err := epicplan.Read(q, o.EpicID)
		if err != nil {
			return nil, Result{}, err
		}
		v := o.Validation
		if action == "add" {
			ordinal++
			v.Ordinal = ordinal
			v.ID = fmt.Sprintf("V%d", ordinal)
		} else {
			found := false
			for _, existing := range list {
				if existing.ID == o.ValidationID {
					v.Ordinal, v.ID = existing.Ordinal, existing.ID
					found = true
					if action == "remove" {
						v = existing
					}
				}
			}
			if !found {
				return nil, Result{}, fmt.Errorf("validation %s not found", o.ValidationID)
			}
		}
		if p.State == "done" || p.State == "cancelled" {
			return nil, Result{}, fmt.Errorf("validation changes require an open epic")
		}
		if action != "remove" {
			v, err = epicplan.Normalize(v)
			if err != nil {
				return nil, Result{}, err
			}
			if err := validateVerifier(o.WorkspaceRoot, q, v.Verifier); err != nil {
				return nil, Result{}, err
			}
		}
		v.Generation = generation + 1
		r := Result{EpicID: o.EpicID, Validation: &v, At: o.now().Format(time.RFC3339)}
		if action == "add" && len(list) >= 10 {
			r.Warnings = []string{"More than 10 validations; consider splitting the epic"}
		}
		e := journal.Effect{Op: "insert", Table: "epic_validation_changes", PK: map[string]any{"epic_id": o.EpicID, "generation": v.Generation}, Values: map[string]any{"ordinal": v.Ordinal, "claim": v.Claim, "verifier": v.Verifier, "recipe": v.Recipe, "removed": action == "remove"}}
		return []journal.Effect{e}, r, nil
	})
}

func validateVerifier(root string, q *sql.DB, verifier string) error {
	if name, ok := strings.CutPrefix(verifier, "check:"); ok {
		p, err := policy.Load(root)
		if err != nil {
			return err
		}
		if _, ok := p.Effective.Checks[name]; !ok {
			return fmt.Errorf("unknown check in default-branch policy: %s", name)
		}
	}
	if name, ok := strings.CutPrefix(verifier, "review:"); ok {
		if name == "code" || name == "test-adequacy" || name == "merge" {
			return nil
		}
		rows, err := q.Query(`SELECT record_json FROM review_rounds`)
		if err != nil {
			return err
		}
		defer func() { _ = rows.Close() }()
		for rows.Next() {
			var raw string
			if err := rows.Scan(&raw); err != nil {
				return err
			}
			var r struct {
				Profiles []string `json:"profiles"`
			}
			if err := json.Unmarshal([]byte(raw), &r); err != nil {
				return err
			}
			for _, p := range r.Profiles {
				if p == name {
					return nil
				}
			}
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return fmt.Errorf("unknown review profile: %s", name)
	}
	return nil
}

func Approve(o Options) (Result, error) {
	if !policy.Allows(o.Actor, policy.PersonDecisions) {
		return Result{}, &persondecision.DeniedError{Actor: o.Actor}
	}
	if strings.TrimSpace(o.Reason) == "" {
		return Result{}, &persondecision.ReasonRequiredError{}
	}
	q, err := o.open()
	if err != nil {
		return Result{}, err
	}
	if o.IdempotencyKey != "" {
		req, err := o.request("epic.approve")
		if err != nil {
			_ = q.Close()
			return Result{}, err
		}
		prior, found, err := journal.ReplayIdempotent(q, req)
		if err != nil {
			_ = q.Close()
			return Result{}, err
		}
		if found {
			var r Result
			err = json.Unmarshal([]byte(prior.ResultJSON), &r)
			r.EventID = prior.EventID
			_ = q.Close()
			return r, err
		}
	}
	preview, err := epicplan.Preview(q, o.EpicID)
	_ = q.Close()
	if err != nil {
		return Result{}, err
	}
	if o.Confirm == nil {
		return Result{}, fmt.Errorf("epic approval requires interactive confirmation")
	}
	yes, err := o.Confirm(preview)
	if err != nil {
		return Result{}, err
	}
	if !yes {
		return Result{}, fmt.Errorf("epic approval cancelled")
	}
	o = o.frozen()
	return o.write("epic.approve", func(q *sql.DB) ([]journal.Effect, Result, error) {
		p, err := epicplan.Preview(q, o.EpicID)
		if err != nil {
			return nil, Result{}, err
		}
		if encode(p) != encode(preview) {
			return nil, Result{}, fmt.Errorf("epic plan changed after preview; retry approval")
		}
		if p.State == "done" || p.State == "cancelled" {
			return nil, Result{}, fmt.Errorf("approval requires an open epic")
		}
		id, e, err := persondecision.Effect(q, "approve-epic", o.Actor, "epic", o.EpicID, o.Reason, o.now())
		if err != nil {
			return nil, Result{}, err
		}
		e.Values["binding_id"] = p.Binding.Digest
		e.Values["binding_json"] = encode(p.Binding)
		effects := []journal.Effect{e}
		r := Result{EpicID: o.EpicID, DecisionID: id, At: o.now().Format(time.RFC3339)}
		for _, t := range p.Tasks {
			if t.Accept {
				effects = append(effects, journal.Effect{Op: "insert", Table: "epic_plan_acceptances", PK: map[string]any{"issue_id": t.ID}, Values: map[string]any{"decision_id": id}})
				r.Accepted = append(r.Accepted, t.ID)
			}
		}
		return effects, r, nil
	})
}

func (o Options) frozen() Options {
	now := o.now()
	o.Clock = func() time.Time { return now }
	return o
}
