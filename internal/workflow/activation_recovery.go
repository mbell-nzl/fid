package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

type ActivationResourceObservationRequest struct {
	ActivationID              string  `json:"activation_id"`
	Component                 string  `json:"component"`
	ExpectedActivationVersion int     `json:"expected_activation_version"`
	ExpectedVersion           int     `json:"expected_version"`
	ObservedIdentity          *string `json:"observed_identity"`
	State                     string  `json:"state"`
	EvidenceRef               string  `json:"evidence_ref"`
	ErrorCode                 *string `json:"error_code"`
}

type ActivationRecoveryRequest struct {
	ActivationID    string `json:"activation_id"`
	ExpectedVersion int    `json:"expected_version"`
	ExpectedStatus  string `json:"expected_status"`
	Direction       string `json:"direction"`
	Reason          string `json:"reason"`
}

type ActivationRecoveryCompleteRequest struct {
	ActivationID     string                      `json:"activation_id"`
	ExpectedVersion  int                         `json:"expected_version"`
	SelectedEventID  string                      `json:"selected_event_id"`
	ResourceVersions []ActivationResourceVersion `json:"resource_versions"`
}

type ActivationResourceVersion struct {
	Component string `json:"component"`
	Version   int    `json:"version"`
}

type ActivationResourceResult struct {
	Activation     Activation         `json:"activation"`
	Resource       ActivationResource `json:"resource"`
	IdempotencyKey string             `json:"idempotency_key"`
	Replayed       bool               `json:"replayed"`
	EventID        string
}

func ObserveActivationResource(options ActivationTransitionOptions, request ActivationResourceObservationRequest) (ActivationResourceResult, error) {
	if strings.TrimSpace(request.ActivationID) == "" || request.ExpectedActivationVersion < 1 || request.ExpectedVersion < 1 || strings.TrimSpace(request.EvidenceRef) == "" || (request.State != "old" && request.State != "new" && request.State != "error") {
		return ActivationResourceResult{}, fmt.Errorf("activation resource input is invalid")
	}
	var observedActivation Activation
	var observedResource ActivationResource
	outcome, err := activationMutation(options, "activation.resource", request, func(conn *sql.DB, activation Activation) ([]journal.Effect, error) {
		if err := requireActivationVersion(activation, request.ExpectedActivationVersion); err != nil {
			return nil, err
		}
		if activation.Status == "aborted" || activation.Status == "finalized" || activation.Status == "superseded" {
			return nil, &ActivationInvalidError{Status: activation.Status, ExpectedStatuses: []string{"prepared", "staged", "committed"}}
		}
		resource, ok := activationResourceByComponent(activation.Resources, request.Component)
		if !ok {
			return nil, fmt.Errorf("activation resource input is invalid")
		}
		if resource.Version != request.ExpectedVersion {
			return nil, &ActivationVersionConflictError{ActivationID: activation.ID, ExpectedVersion: request.ExpectedVersion, ActualVersion: resource.Version}
		}
		if !resourceObservationAgrees(resource, request) {
			return nil, fmt.Errorf("activation resource observation does not agree")
		}
		updatedAt := time.Now().UTC().Format(time.RFC3339Nano)
		updated := resource
		updated.Version++
		updated.ObservedIdentity, updated.State, updated.EvidenceRef, updated.ErrorCode = request.ObservedIdentity, request.State, &request.EvidenceRef, request.ErrorCode
		observedActivation = activation
		observedActivation.Version++
		observedActivation.UpdatedAt = updatedAt
		for index := range observedActivation.Resources {
			if observedActivation.Resources[index].Component == updated.Component {
				observedActivation.Resources[index] = updated
			}
		}
		observedResource = updated
		return []journal.Effect{
			activationUpdateEffect(activation, map[string]any{"updated_at": updatedAt}),
			{Op: "update", Table: "activation_resources", PK: map[string]any{"activation_id": activation.ID, "component": resource.Component}, Preconditions: map[string]any{"version": resource.Version}, Values: map[string]any{"version": updated.Version, "observed_identity": updated.ObservedIdentity, "state": updated.State, "evidence_ref": updated.EvidenceRef, "error_code": updated.ErrorCode}},
		}, nil
	}, func(journal.Event) (string, error) {
		return activationResourceResultJSON(observedActivation, observedResource, options.IdempotencyKey)
	})
	if err != nil {
		return ActivationResourceResult{}, err
	}
	return activationResourceResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
}

func RecoverActivation(options ActivationTransitionOptions, request ActivationRecoveryRequest) (ActivationPrepareResult, error) {
	if strings.TrimSpace(request.ActivationID) == "" || request.ExpectedVersion < 1 || strings.TrimSpace(request.Reason) == "" || (request.Direction != "restore-prior" && request.Direction != "complete-forward") {
		return ActivationPrepareResult{}, fmt.Errorf("activation recover input is invalid")
	}
	var selected Activation
	outcome, err := activationMutation(options, "activation.recover", request, func(conn *sql.DB, activation Activation) ([]journal.Effect, error) {
		if err := requireActivationVersion(activation, request.ExpectedVersion); err != nil {
			return nil, err
		}
		if activation.Status == "aborted" || activation.Status == "finalized" || activation.Status == "superseded" {
			return nil, &ActivationInvalidError{Status: activation.Status, ExpectedStatuses: []string{"prepared", "staged", "committed"}}
		}
		if activation.Status != request.ExpectedStatus {
			return nil, &ActivationInvalidError{Status: activation.Status, ExpectedStatuses: []string{request.ExpectedStatus}}
		}
		required := "restore-prior"
		if activation.PointOfNoReturn {
			required = "complete-forward"
		}
		if request.Direction != required {
			return nil, &ActivationDirectionConflictError{ActivationID: activation.ID, PointOfNoReturn: activation.PointOfNoReturn, RequiredDirection: required}
		}
		updatedAt := time.Now().UTC().Format(time.RFC3339Nano)
		stopEffects, err := activationIncompleteStopEffects(conn, activation, "fid-bootstrap", updatedAt)
		if err != nil {
			return nil, err
		}
		effects := []journal.Effect{
			activationUpdateEffect(activation, map[string]any{"updated_at": updatedAt}),
			{Op: "insert", Table: "activation_recoveries", PK: map[string]any{"activation_id": activation.ID}, Values: map[string]any{"direction": request.Direction, "actor_attribution": options.BootstrapActor, "reason": request.Reason, "selected_event_id": "<event-id>", "completed_event_id": nil}},
		}
		selected = activation
		selected.Version++
		selected.UpdatedAt = updatedAt
		selected.Recovery = &ActivationRecovery{Direction: request.Direction, ActorAttribution: options.BootstrapActor, Reason: request.Reason}
		return append(effects, stopEffects...), nil
	}, func(prepared journal.Event) (string, error) {
		selected.Recovery.SelectedEventID = prepared.EventID
		return activationResultJSON(selected, options.IdempotencyKey)
	})
	if err != nil {
		return ActivationPrepareResult{}, err
	}
	return activationPrepareResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
}

func CompleteActivationRecovery(options ActivationTransitionOptions, request ActivationRecoveryCompleteRequest) (ActivationPrepareResult, error) {
	if strings.TrimSpace(request.ActivationID) == "" || request.ExpectedVersion < 1 || strings.TrimSpace(request.SelectedEventID) == "" || !validResourceVersions(request.ResourceVersions) {
		return ActivationPrepareResult{}, fmt.Errorf("activation recover-complete input is invalid")
	}
	var completed Activation
	outcome, err := activationMutation(options, "activation.recover-complete", request, func(conn *sql.DB, activation Activation) ([]journal.Effect, error) {
		if err := requireActivationVersion(activation, request.ExpectedVersion); err != nil {
			return nil, err
		}
		if activation.Recovery == nil || activation.Recovery.SelectedEventID != request.SelectedEventID || activation.Recovery.CompletedEventID != nil {
			return nil, fmt.Errorf("activation recovery selection is invalid")
		}
		for _, resource := range activation.Resources {
			expected := versionForResource(request.ResourceVersions, resource.Component)
			if expected != resource.Version {
				return nil, &ActivationVersionConflictError{ActivationID: activation.ID, ExpectedVersion: expected, ActualVersion: resource.Version}
			}
		}
		if activation.Recovery.Direction == "complete-forward" {
			if err := requireNewResourceAgreement(conn, activation.ID); err != nil {
				return nil, err
			}
		} else if err := requirePriorResourceAgreement(activation); err != nil {
			return nil, err
		}
		updatedAt := time.Now().UTC().Format(time.RFC3339Nano)
		status := "aborted"
		componentStatus := "aborted"
		if activation.Recovery.Direction == "complete-forward" {
			status, componentStatus = "finalized", "finalized"
		}
		effects := []journal.Effect{}
		if status == "finalized" && activation.PriorActivationID != nil {
			prior, err := loadActivation(conn, *activation.PriorActivationID)
			if err != nil {
				return nil, err
			}
			effects = append(effects, activationUpdateEffect(prior, map[string]any{"status": "superseded", "updated_at": updatedAt}))
		}
		effects = append(effects, activationUpdateEffect(activation, map[string]any{"status": status, "updated_at": updatedAt}), journal.Effect{Op: "update", Table: "activation_recoveries", PK: map[string]any{"activation_id": activation.ID}, Values: map[string]any{"completed_event_id": "<event-id>"}})
		components, err := activationComponentIdentities(conn, activation.ID)
		if err != nil {
			return nil, err
		}
		effects = append(effects, activationComponentStatusEffects(activation.ID, components, componentStatus)...)
		resolutionEffects, err := activationIncompleteStopResolutionEffects(conn, activation.ID, "fid-bootstrap", time.Now().UTC().Format(time.RFC3339Nano))
		if err != nil {
			return nil, err
		}
		effects = append(effects, resolutionEffects...)
		completed = activation
		completed.Status = status
		completed.Version++
		completed.UpdatedAt = updatedAt
		completed.Recovery = &ActivationRecovery{Direction: activation.Recovery.Direction, ActorAttribution: activation.Recovery.ActorAttribution, Reason: activation.Recovery.Reason, SelectedEventID: activation.Recovery.SelectedEventID}
		return effects, nil
	}, func(prepared journal.Event) (string, error) {
		completionID := prepared.EventID
		completed.Recovery.CompletedEventID = &completionID
		return activationResultJSON(completed, options.IdempotencyKey)
	})
	if err != nil {
		return ActivationPrepareResult{}, err
	}
	return activationPrepareResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
}

func requirePriorResourceAgreement(activation Activation) error {
	resources := make(map[string]ActivationResource, len(activation.Resources))
	for _, resource := range activation.Resources {
		resources[resource.Component] = resource
	}
	components := make([]string, 0)
	reasons := make([]EvaluationReason, 0)
	for _, component := range activationResourceComponents {
		resource, found := resources[component]
		if !found {
			components = append(components, component)
			reasons = append(reasons, evaluationReason("activation-incomplete", "resource "+component+" is missing", activation.ID, "activation"))
			continue
		}
		if resource.State == "old" && sameOptionalIdentity(resource.ObservedIdentity, resource.OldIdentity) {
			continue
		}
		components = append(components, component)
		reasons = append(reasons, evaluationReason("activation-incomplete", "resource "+component+" does not agree with the prior identity", activation.ID, "activation"))
	}
	if len(components) == 0 {
		return nil
	}
	return &ActivationIncompleteError{ActivationID: activation.ID, Components: components, Reasons: reasons}
}

func sameOptionalIdentity(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func activationMutation(options ActivationTransitionOptions, operation string, request any, build func(*sql.DB, Activation) ([]journal.Effect, error), resultJSON func(journal.Event) (string, error)) (journal.IdempotencyOutcome, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" {
		return journal.IdempotencyOutcome{}, fmt.Errorf("%s: workspace root and migrations directory are required", operation)
	}
	project, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, "config.toml"))
	if err != nil {
		return journal.IdempotencyOutcome{}, err
	}
	var outcome journal.IdempotencyOutcome
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "fid " + strings.ReplaceAll(operation, ".", " ")}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}
		conn, err := openWorkflowDB(options.WorkspaceRoot)
		if err != nil {
			return err
		}
		defer closeWorkflowDB(conn)
		digest, fields, err := journal.ComputeRequestDigest(map[string]any{"bootstrap_actor": options.BootstrapActor, "operation": operation, "project_id": project.ProjectID, "request": request}, options.BootstrapActor, nil)
		if err != nil {
			return err
		}
		idempotency := journal.IdempotencyRequest{OperationScope: operation + ":" + project.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
		if replayed, found, err := journal.ReplayIdempotent(conn, idempotency); err != nil {
			return err
		} else if found {
			outcome = replayed
			return nil
		}
		activationID := activationIDFromRequest(request)
		activation, err := loadActivation(conn, activationID)
		if err != nil {
			return err
		}
		effects, err := build(conn, activation)
		if err != nil {
			return err
		}
		eventType := map[string]string{"activation.resource": "activation.resource-observed", "activation.recover": "activation.recovery-selected", "activation.recover-complete": "activation.recovery-completed"}[operation]
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: eventType, ActorV2: journal.EventActor{ActorID: "fid-bootstrap", ActorType: "service", Producer: "fid", ProvenanceRef: "activation:bootstrap"}, Aggregate: journal.Aggregate{ID: activation.ID, Sequence: activation.Version + 1, Type: "activation"}, CommandName: operation, Effects: effects}
		outcome, err = journal.CommitIdempotent(lock, conn, idempotency, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			for index := range prepared.Effects {
				for field, value := range prepared.Effects[index].Values {
					if value == "<event-id>" {
						prepared.Effects[index].Values[field] = prepared.EventID
					}
				}
			}
			result := `{"replayed":false}`
			if resultJSON != nil {
				var resultErr error
				result, resultErr = resultJSON(prepared)
				if resultErr != nil {
					return journal.IdempotencyMutation{}, resultErr
				}
			}
			return journal.IdempotencyMutation{EntityID: activation.ID, ResultJSON: result}, nil
		})
		return err
	})
	return outcome, err
}

func activationIncompleteStopEffects(conn *sql.DB, activation Activation, actorID, raisedAt string) ([]journal.Effect, error) {
	components := make([]string, 0, len(activation.Resources))
	for _, resource := range activation.Resources {
		components = append(components, resource.Component)
	}
	details, err := json.Marshal(map[string]any{"activation_id": activation.ID, "point_of_no_return": activation.PointOfNoReturn, "components": components})
	if err != nil {
		return nil, err
	}
	subject := Subject{Type: "activation", ID: activation.ID}
	request := StopRaiseRequest{Subject: subject, Code: "activation-incomplete", BlockedActions: []string{"claim.acquire", "envelope.complete", "gate.run", "issue.transition-terminal-success", "run.start"}, Details: details, ResolutionCondition: Condition{Type: "activation-recovered", Subject: subject}, ResolverCapability: "stop.resolve"}
	if err := validateStopRaise(request); err != nil {
		return nil, err
	}
	stopID, err := newStopID(conn)
	if err != nil {
		return nil, err
	}
	key, err := attentionSemanticKey(subject, request.Code, request.ResolutionCondition)
	if err != nil {
		return nil, err
	}
	attention, found, err := loadOpenAttentionByKey(conn, key)
	if err != nil {
		return nil, err
	}
	if !found {
		attention.ID, err = newAttentionID(conn)
		if err != nil {
			return nil, err
		}
		attention = Attention{ID: attention.ID, SemanticKey: key, Subject: subject, Code: request.Code, Status: "open", Version: 1}
	}
	attention.RaisedAt = raisedAt
	stop := Stop{ID: stopID, Subject: subject, Code: request.Code, BlockedActions: request.BlockedActions, Details: details, ResolutionCondition: request.ResolutionCondition, ResolverCapability: request.ResolverCapability, ActorID: actorID, RaisedAt: raisedAt}
	return stopRaisedEventForDomain(stop, attention, !found, "<event-id>").Effects, nil
}

func activationIncompleteStopResolutionEffects(conn *sql.DB, activationID, actorID, resolvedAt string) ([]journal.Effect, error) {
	rows, err := conn.Query(`SELECT s.id FROM stops s LEFT JOIN stop_resolutions r ON r.stop_id = s.id WHERE s.subject_type = 'activation' AND s.subject_id = ? AND s.code = 'activation-incomplete' AND r.stop_id IS NULL`, activationID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	effects := []journal.Effect{}
	for rows.Next() {
		var stopID string
		if err := rows.Scan(&stopID); err != nil {
			return nil, err
		}
		effects = append(effects, journal.Effect{Op: "insert", Table: "stop_resolutions", PK: map[string]any{"stop_id": stopID}, Values: map[string]any{"actor_id": actorID, "fact_event_id": "<event-id>", "resolved_at": resolvedAt, "created_event_id": "<event-id>"}})
	}
	return effects, rows.Err()
}

func activationIDFromRequest(request any) string {
	switch value := request.(type) {
	case ActivationResourceObservationRequest:
		return value.ActivationID
	case ActivationRecoveryRequest:
		return value.ActivationID
	case ActivationRecoveryCompleteRequest:
		return value.ActivationID
	}
	return ""
}

func activationResourceResultJSON(activation Activation, resource ActivationResource, key string) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"activation": activation, "resource": resource, "idempotency_key": key, "replayed": false})
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func activationResourceResultFromJSON(raw, eventID string, replayed bool) (ActivationResourceResult, error) {
	var value struct {
		Activation     Activation         `json:"activation"`
		Resource       ActivationResource `json:"resource"`
		IdempotencyKey string             `json:"idempotency_key"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return ActivationResourceResult{}, fmt.Errorf("decode stored activation resource result: %w", err)
	}
	return ActivationResourceResult{Activation: value.Activation, Resource: value.Resource, IdempotencyKey: value.IdempotencyKey, Replayed: replayed, EventID: eventID}, nil
}
func resourceObservationAgrees(resource ActivationResource, request ActivationResourceObservationRequest) bool {
	if request.State == "new" {
		return request.ObservedIdentity != nil && *request.ObservedIdentity == resource.NewIdentity
	}
	if request.State == "old" {
		return (resource.OldIdentity == nil && request.ObservedIdentity == nil) || (resource.OldIdentity != nil && request.ObservedIdentity != nil && *resource.OldIdentity == *request.ObservedIdentity)
	}
	return request.ObservedIdentity == nil && request.ErrorCode != nil
}
func validResourceVersions(values []ActivationResourceVersion) bool {
	if len(values) != 5 {
		return false
	}
	seen := map[string]bool{}
	for _, value := range values {
		if value.Version < 1 || seen[value.Component] {
			return false
		}
		seen[value.Component] = true
	}
	for _, component := range activationResourceComponents {
		if !seen[component] {
			return false
		}
	}
	return true
}
func versionForResource(values []ActivationResourceVersion, component string) int {
	for _, value := range values {
		if value.Component == component {
			return value.Version
		}
	}
	return 0
}
