package checkrun

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"time"

	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/policysnapshot"
	"github.com/mbell-nzl/fid/internal/project"
)

type Options struct {
	WorkspaceRoot, MigrationsDir, HoldID, Actor, IdempotencyKey, Version string
	Names                                                                []string
	Clock                                                                func() time.Time
	Output                                                               io.Writer
	Before                                                               func(string) error
	After                                                                func(checkrecord.Record) error
	Warning                                                              func(string) error
	RenewInterval                                                        time.Duration
}

type execution struct {
	options   Options
	request   journal.IdempotencyRequest
	response  checkrecord.Response
	policy    policy.Loaded
	env       []string
	requestID string
}

func Run(ctx context.Context, options Options) (checkrecord.Response, error) {
	if options.Version == "" {
		options.Version = Version()
	}
	options.Names = append([]string{}, options.Names...)
	slices.Sort(options.Names)
	options.Names = slices.Compact(options.Names)
	if options.IdempotencyKey == "" {
		return checkrecord.Response{}, fmt.Errorf("check requires an idempotency key")
	}
	if !policy.Allows(options.Actor, policy.WorkSet) {
		return checkrecord.Response{}, &checkrecord.Refusal{Code: "authority-denied", Message: "checks require an agent or a person", Details: map[string]any{"actor_id": options.Actor, "operation": "check.execute", "required_actions": []string{"check.execute"}, "required_capabilities": []string{"check.execute"}}}
	}
	ex := execution{options: options, env: environment(), response: checkrecord.Response{Names: options.Names, CheckIDs: []string{}, Checks: []checkrecord.Summary{}, NotStarted: []string{}}}
	replayed := false
	replayRecords := []checkrecord.Record{}
	err := options.write(func(db store) error {
		var err error
		ex.request, err = db.request()
		if err != nil {
			return err
		}
		outcome, found, err := journal.ReplayIdempotent(db.conn, ex.request)
		if err != nil {
			return err
		}
		if found {
			replayed = true
			if err := json.Unmarshal([]byte(outcome.ResultJSON), &ex.response); err != nil {
				return fmt.Errorf("decode check retry: %w", err)
			}
			if !ex.response.Complete {
				return &checkrecord.Refusal{Code: "check-incomplete", TaskID: ex.response.TaskID, Message: "this check request never finished", Details: map[string]any{"check_ids": ex.response.CheckIDs}}
			}
			if options.Before != nil || options.After != nil {
				records, err := checkread.Read(db.conn, ex.response.TaskID)
				if err != nil {
					return err
				}
				for _, record := range records {
					if slices.Contains(ex.response.CheckIDs, record.Start.ID) {
						replayRecords = append(replayRecords, record)
					}
				}
			}
			return nil
		}
		if err := db.current(); err != nil {
			return err
		}
		hold, err := holds.CurrentByID(db.conn, options.HoldID)
		if err != nil {
			return err
		}
		ex.response.TaskID = hold.TaskID
		var comparison policy.Comparison
		if len(options.Names) == 0 {
			comparison, err = policy.Compare(options.WorkspaceRoot, hold.TaskID)
			if err != nil {
				return err
			}
		}
		ex.policy, err = policy.Load(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		if err := ex.validateTask(db); err != nil {
			return err
		}
		if len(options.Names) == 0 {
			ex.options.Names, err = policy.Required(ex.policy.Effective, ex.response.Lane, hold.TaskID, comparison.Paths)
			if err != nil {
				return err
			}
			ex.response.Names = ex.options.Names
		}
		if err := policysnapshot.Ensure(db.lock, db.conn, ex.policy, db.now); err != nil {
			return err
		}
		if len(ex.options.Names) == 0 {
			return ex.empty(db)
		}
		return nil
	})
	if err != nil {
		return ex.response, err
	}
	if replayed {
		for _, record := range replayRecords {
			if options.Before != nil {
				if err := options.Before(record.Start.Name); err != nil {
					return ex.response, err
				}
			}
			if options.After != nil {
				if err := options.After(record); err != nil {
					return ex.response, err
				}
			}
		}
		return ex.response, responseError(ex.response)
	}
	for index, name := range ex.options.Names {
		var record checkrecord.Record
		admissionFailure := false
		err := options.write(func(db store) error {
			if err := db.current(); err != nil {
				admissionFailure = true
				return err
			}
			if err := ex.previous(db); err != nil {
				admissionFailure = true
				return err
			}
			var err error
			record.Start, err = ex.start(db, name)
			return err
		})
		if err != nil {
			if admissionFailure && ex.requestID != "" {
				if stopErr := ex.stop(index, err); stopErr != nil {
					return ex.response, stopErr
				}
				return ex.response, responseError(ex.response)
			}
			return ex.response, err
		}
		if options.Before != nil {
			if err := options.Before(name); err != nil {
				return ex.response, err
			}
		}
		process, err := ExecuteProcess(ctx, ProcessOptions{Root: options.WorkspaceRoot, Definition: record.Start.Definition, Environment: ex.env, Version: options.Version, Output: options.Output, RenewInterval: options.RenewInterval,
			Started: func(process checkrecord.Process) error {
				return options.write(func(db store) error {
					event, err := db.event("check.process-started", record.Start.ID, []journal.Effect{{Op: "insert", Table: "check_processes", PK: map[string]any{"check_id": record.Start.ID}, Values: map[string]any{"process_json": encoded(process)}}})
					if err != nil {
						return err
					}
					_, err = journal.CommitEvent(db.lock, db.conn, event)
					return err
				})
			}, Renew: func() error { return ex.renew(record.Start.ID) }})
		if err != nil {
			return ex.response, err
		}
		result := process.Result
		err = options.write(func(db store) error {
			observation, err := project.ObserveRepository(options.WorkspaceRoot)
			if err != nil {
				return err
			}
			result.Tree = tree(observation)
			result.FinishedAt = db.now.Format(time.RFC3339Nano)
			current := db.current()
			if current != nil {
				var refused *holdwork.Refusal
				if !asHoldRefusal(current, &refused) {
					return current
				}
				process.StopError = current
			}
			if err := ex.keepFiles(db, record.Start, &result); err != nil {
				return err
			}
			result.Reasons = checkrecord.Qualification(record.Start.Tree, result.Tree, result.Outcome, current == nil, len(result.RefusedArtifacts) == 0)
			result.Counts = result.Outcome == "passed" && len(result.Reasons) == 0
			record.Result = &result
			ex.response.Checks = append(ex.response.Checks, summary(record))
			complete := index == len(ex.options.Names)-1 || process.StopError != nil || ctx.Err() != nil
			if complete {
				ex.response.Complete = true
				ex.response.NotStarted = append([]string{}, ex.options.Names[index+1:]...)
				ex.response.Failure = aggregateFailure(ex.response)
				if process.StopError != nil {
					var refused *holdwork.Refusal
					if asHoldRefusal(process.StopError, &refused) {
						ex.response.Failure = &checkrecord.Refusal{Code: refused.Code, TaskID: refused.TaskID, Message: refused.Message, Details: refused.Details}
					} else {
						return process.StopError
					}
				}
			}
			return ex.finish(db, record, complete)
		})
		if err != nil {
			return ex.response, err
		}
		if options.After != nil {
			if err := options.After(record); err != nil {
				return ex.response, err
			}
		}
		if process.StopError != nil || ctx.Err() != nil {
			break
		}
	}
	return ex.response, responseError(ex.response)
}

func (ex *execution) validateTask(db store) error {
	allowed := []string{}
	for name := range ex.policy.Effective.Checks {
		allowed = append(allowed, name)
	}
	slices.Sort(allowed)
	for _, name := range ex.options.Names {
		if _, ok := ex.policy.Effective.Checks[name]; !ok {
			return &checkrecord.Refusal{Code: "invalid-argument", Message: "unknown check: " + name, Details: map[string]any{"argument": "names", "value": name, "allowed": allowed}}
		}
	}
	var state string
	var lane sql.NullString
	if err := db.conn.QueryRow(`SELECT status,lane FROM issues WHERE id=? AND kind='task'`, ex.response.TaskID).Scan(&state, &lane); err != nil {
		return err
	}
	if state != "in_progress" {
		return &holdwork.Refusal{Code: "transition-not-permitted", TaskID: ex.response.TaskID, Message: "start the task before checking it", Details: map[string]any{"from_state": state, "to_state": "in_progress", "unsatisfied": []string{"task-in-progress"}}}
	}
	if !lane.Valid {
		return &holdwork.Refusal{Code: "lane-not-set", TaskID: ex.response.TaskID, Message: "the task has no lane", Details: map[string]any{"task_id": ex.response.TaskID}}
	}
	ex.response.Lane = lane.String
	if !ex.policy.Effective.Lanes[lane.String].ChecksSet {
		return &checkrecord.Refusal{Code: "checks-not-configured", TaskID: ex.response.TaskID, Message: "no checks set for " + lane.String + " in .fid.toml", Details: map[string]any{"lane": lane.String}, Next: "add checks to [lanes." + lane.String + "], then merge .fid.toml and git fetch"}
	}
	return nil
}

func (ex *execution) previous(db store) error {
	records, err := checkread.Read(db.conn, ex.response.TaskID)
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Result != nil {
			continue
		}
		process := record.Start.Owner
		if record.Process != nil {
			process = *record.Process
		}
		alive, err := processAlive(process)
		if err != nil {
			return fmt.Errorf("inspect earlier check %s: %w", record.Start.ID, err)
		}
		if alive {
			return &checkrecord.Refusal{Code: "check-still-running", TaskID: ex.response.TaskID, Message: fmt.Sprintf("%s is still running from an earlier check (%s, pid %d)", record.Start.Name, record.Start.ID, process.PID), Details: map[string]any{"check_id": record.Start.ID, "pid": process.PID}, Next: fmt.Sprintf("wait for it to finish, or stop pid %d, then check again", process.PID)}
		}
		if ex.options.Warning != nil {
			if err := ex.options.Warning("⚠ " + record.Start.ID + " never finished; its process has stopped"); err != nil {
				return err
			}
		}
	}
	return nil
}

func tree(observation project.Observation) checkrecord.Tree {
	return checkrecord.Tree{Head: observation.HeadOID, Tree: observation.TreeOID, Dirty: observation.Dirty, DirtyDigest: observation.DirtyDigest}
}

func (ex *execution) start(db store, name string) (checkrecord.Start, error) {
	id, err := db.newCheckID(journal.GenerateID)
	if err != nil {
		return checkrecord.Start{}, err
	}
	observation, err := project.ObserveRepository(ex.options.WorkspaceRoot)
	if err != nil {
		return checkrecord.Start{}, err
	}
	owner, err := processIdentity(os.Getpid())
	if err != nil {
		return checkrecord.Start{}, err
	}
	check := ex.policy.Effective.Checks[name]
	definition := checkrecord.DefinitionFor(check)
	start := checkrecord.Start{ID: id, TaskID: ex.response.TaskID, HoldID: ex.options.HoldID, Actor: ex.options.Actor, Name: name, Definition: definition, DefinitionDigest: checkrecord.Digest(definition), EnvironmentDigest: checkrecord.Digest(ex.env), StartedAt: db.now.Format(time.RFC3339Nano), Tree: tree(observation), Owner: owner}
	if err := db.conn.QueryRow(`SELECT sequence FROM policy_snapshots WHERE digest=? ORDER BY sequence DESC LIMIT 1`, ex.policy.Effective.Digest()).Scan(&start.PolicySnapshot); err != nil {
		return start, err
	}
	effects := []journal.Effect{}
	first := ex.requestID == ""
	if first {
		ex.requestID = id
		effects = append(effects, journal.Effect{Op: "insert", Table: "check_requests", PK: map[string]any{"id": id}, Values: map[string]any{"task_id": start.TaskID, "hold_id": start.HoldID, "actor": start.Actor, "names_json": encoded(ex.options.Names), "response_json": nil}})
	}
	effects = append(effects, journal.Effect{Op: "insert", Table: "check_starts", PK: map[string]any{"id": id}, Values: map[string]any{"request_id": ex.requestID, "task_id": start.TaskID, "name": name, "start_json": encoded(start)}})
	ex.response.CheckIDs = append(ex.response.CheckIDs, id)
	event, err := db.event("check.started", id, effects)
	if err != nil {
		return start, err
	}
	if first {
		_, err = journal.CommitIdempotent(db.lock, db.conn, ex.request, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: id, ResultJSON: encoded(ex.response)}, nil
		})
	} else {
		err = ex.commitProgress(db, event)
	}
	return start, err
}

func responseError(response checkrecord.Response) error {
	if response.Failure != nil {
		return response.Failure
	}
	return nil
}
