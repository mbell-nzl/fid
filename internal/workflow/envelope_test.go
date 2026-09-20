package workflow

import "testing"

func TestEnvelopeCreateAcceptsTypedObligationSet(t *testing.T) {
	request := testEnvelopeCreateRequest()
	request.Obligations = []ObligationInput{
		testObligation("OBL-gate", "gate-evidence", request.PrimaryIssueID),
		testObligation("OBL-review", "review", request.PrimaryIssueID),
		testObligation("OBL-human", "human-outcome", request.PrimaryIssueID),
	}
	request.Obligations[1].ReviewProfile = stringPointer("review-profile")
	request.Obligations[2].OutcomeClass = stringPointer("outcome-class")

	if err := validateEnvelopeCreate(request); err != nil {
		t.Fatalf("validateEnvelopeCreate() error = %v, want typed obligation set accepted", err)
	}
}

func TestEnvelopeCreateRejectsInvalidTypedObligationAndProtectedDecisionAuthoring(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*EnvelopeCreateRequest)
	}{
		{
			name: "zero obligations",
			mutate: func(request *EnvelopeCreateRequest) {
				request.Obligations = nil
			},
		},
		{
			name: "duplicate review profile",
			mutate: func(request *EnvelopeCreateRequest) {
				request.Obligations = []ObligationInput{testObligation("OBL-review-a", "review", request.PrimaryIssueID), testObligation("OBL-review-b", "review", request.PrimaryIssueID)}
				request.Obligations[0].ReviewProfile = stringPointer("same-profile")
				request.Obligations[1].ReviewProfile = stringPointer("same-profile")
			},
		},
		{
			name: "decision action is not guarded",
			mutate: func(request *EnvelopeCreateRequest) {
				request.Obligations = []ObligationInput{testObligation("OBL-gate", "gate-evidence", request.PrimaryIssueID)}
				request.ProtectedDecisions = []ProtectedDecisionInput{{Ordinal: 0, DecisionClass: "release", RequiredCapabilities: []string{"envelope.complete"}, AffectedActions: []string{"issue.transition-terminal-success"}}}
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := testEnvelopeCreateRequest()
			test.mutate(&request)
			if err := validateEnvelopeCreate(request); err == nil {
				t.Fatal("validateEnvelopeCreate() error = nil, want rejection")
			}
		})
	}
}

func TestEnvelopeCreatePersistsEveryTypedObligationAndProtectedDecision(t *testing.T) {
	request := testEnvelopeCreateRequest()
	request.Obligations = []ObligationInput{
		testObligation("OBL-gate", "gate-evidence", request.PrimaryIssueID),
		testObligation("OBL-review", "review", request.PrimaryIssueID),
		testObligation("OBL-human", "human-outcome", request.PrimaryIssueID),
	}
	request.Obligations[1].ReviewProfile = stringPointer("review-profile")
	request.Obligations[2].OutcomeClass = stringPointer("outcome-class")
	request.ProtectedDecisions = []ProtectedDecisionInput{{Ordinal: 0, DecisionClass: "release", RequiredCapabilities: []string{"envelope.complete"}, AffectedActions: []string{"envelope.complete"}}}

	envelope := envelopeFromCreate("ENV-test", "ACT-test", "2026-09-20T00:00:00Z", request)
	if len(envelope.Obligations) != 3 || len(envelope.ProtectedDecisions) != 1 {
		t.Fatalf("created Envelope = %#v, want all typed authoring values", envelope)
	}
	effects := envelopeCreateEffects(envelope, request)
	obligations, decisions := 0, 0
	for _, effect := range effects {
		switch effect.Table {
		case "envelope_obligations":
			obligations++
		case "envelope_protected_decisions":
			decisions++
		}
	}
	if obligations != 3 || decisions != 1 {
		t.Fatalf("envelope.created effects = %#v, want 3 obligations and 1 protected decision", effects)
	}
}

func TestEvaluateEnvelopeObligationsReturnsIndependentOrderedResults(t *testing.T) {
	request := testEnvelopeCreateRequest()
	request.Obligations = []ObligationInput{
		testObligation("OBL-review", "review", request.PrimaryIssueID),
		testObligation("OBL-gate", "gate-evidence", request.PrimaryIssueID),
		testObligation("OBL-human", "human-outcome", request.PrimaryIssueID),
	}
	request.Obligations[0].ReviewProfile = stringPointer("review-profile")
	request.Obligations[2].OutcomeClass = stringPointer("outcome-class")
	envelope := envelopeFromCreate("ENV-test", "ACT-test", "2026-09-20T00:00:00Z", request)
	states := map[string][]ObligationEvidenceState{
		"OBL-review": {{ObligationID: "OBL-review", EvidenceID: "EVD-review", ReceiptSuccessful: true, Validity: "invalid", Freshness: "stale", ArtifactAvailability: "not-applicable"}},
		"OBL-gate":   {{ObligationID: "OBL-gate"}},
		"OBL-human":  {{ObligationID: "OBL-human"}},
	}

	got := EvaluateEnvelopeObligations(envelope.Obligations, states)
	if len(got) != 3 || got[0].ObligationID != "OBL-gate" || got[1].ObligationID != "OBL-human" || got[2].ObligationID != "OBL-review" {
		t.Fatalf("evaluations = %#v, want ID-ordered results", got)
	}
	if got[0].ArtifactAvailability != "missing" || got[1].ArtifactAvailability != "not-applicable" || got[2].ArtifactAvailability != "not-applicable" {
		t.Fatalf("artifact availability = %#v, want missing only for absent gate Evidence", got)
	}
	if len(got[0].Reasons) != 1 || got[0].Reasons[0].Code != "obligation-evidence-absent" || len(got[1].Reasons) != 1 || got[1].Reasons[0].Code != "human-outcome-absent" || len(got[2].Reasons) != 2 {
		t.Fatalf("evaluation reasons = %#v, want per-type absence and independent review dimensions", got)
	}
}

func TestEnvelopeActivationEvaluationIncludesEveryObligation(t *testing.T) {
	request := testEnvelopeCreateRequest()
	request.Obligations = []ObligationInput{
		testObligation("OBL-review", "review", request.PrimaryIssueID),
		testObligation("OBL-gate", "gate-evidence", request.PrimaryIssueID),
	}
	request.Obligations[0].ReviewProfile = stringPointer("review-profile")
	evaluation := envelopeEvaluation(envelopeFromCreate("ENV-test", "ACT-test", "2026-09-20T00:00:00Z", request))
	if len(evaluation) != 2 || evaluation[0].ObligationID != "OBL-gate" || evaluation[1].ObligationID != "OBL-review" {
		t.Fatalf("activation Evaluation = %#v, want one ordered result per Obligation", evaluation)
	}
}

func TestEnvelopeAuthoringObligationCardinalityAndTypeFieldMatrix(t *testing.T) {
	request := testEnvelopeCreateRequest()
	for _, test := range []struct {
		name        string
		obligations func(string) []ObligationInput
		wantError   bool
	}{
		{name: "zero", obligations: func(string) []ObligationInput { return nil }, wantError: true},
		{name: "one", obligations: func(issueID string) []ObligationInput {
			return []ObligationInput{testObligation("OBL-1", "gate-evidence", issueID)}
		}},
		{name: "sixteen", obligations: func(issueID string) []ObligationInput {
			return envelopeObligationsForCount(issueID)[:16]
		}},
		{name: "seventeen", obligations: envelopeObligationsForCount, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			candidate := request
			count := map[string]int{"zero": 0, "one": 1, "sixteen": 16, "seventeen": 17}[test.name]
			candidate.Obligations = test.obligations(candidate.PrimaryIssueID)
			if len(candidate.Obligations) != count {
				t.Fatalf("obligation count = %d, want %d", len(candidate.Obligations), count)
			}
			err := validateEnvelopeCreate(candidate)
			if (err != nil) != test.wantError {
				t.Fatalf("validateEnvelopeCreate() error = %v, wantError=%t", err, test.wantError)
			}
		})
	}
}

func TestEnvelopeAuthoringAcceptedAndForbiddenTypedFieldCombinations(t *testing.T) {
	for _, test := range []struct {
		name           string
		obligationType string
		mutate         func(*ObligationInput)
		wantError      bool
	}{
		{name: "gate-evidence accepted", obligationType: "gate-evidence"},
		{name: "review accepted", obligationType: "review", mutate: func(obligation *ObligationInput) {
			obligation.ReviewProfile = stringPointer("profile")
		}},
		{name: "human-outcome accepted", obligationType: "human-outcome", mutate: func(obligation *ObligationInput) {
			obligation.OutcomeClass = stringPointer("outcome")
		}},
		{name: "gate missing registration", obligationType: "gate-evidence", mutate: func(obligation *ObligationInput) {
			obligation.Registration = struct {
				ID      string `json:"id"`
				Version int    `json:"version"`
			}{}
		}, wantError: true},
		{name: "gate has review profile", obligationType: "gate-evidence", mutate: func(obligation *ObligationInput) {
			obligation.ReviewProfile = stringPointer("profile")
		}, wantError: true},
		{name: "gate has outcome class", obligationType: "gate-evidence", mutate: func(obligation *ObligationInput) {
			obligation.OutcomeClass = stringPointer("outcome")
		}, wantError: true},
		{name: "gate evidence type mismatch", obligationType: "gate-evidence", mutate: func(obligation *ObligationInput) {
			obligation.AcceptedEvidenceTypes = []string{"review"}
		}, wantError: true},
		{name: "gate trust source mismatch", obligationType: "gate-evidence", mutate: func(obligation *ObligationInput) {
			obligation.AcceptedTrustSources = []string{"activated-actor"}
		}, wantError: true},
		{name: "gate freshness mismatch", obligationType: "gate-evidence", mutate: func(obligation *ObligationInput) {
			obligation.FreshnessTarget = "none"
		}, wantError: true},
		{name: "gate artifact policy mismatch", obligationType: "gate-evidence", mutate: func(obligation *ObligationInput) {
			obligation.ArtifactPolicy = "metadata-sufficient"
		}, wantError: true},
		{name: "gate independence mismatch", obligationType: "gate-evidence", mutate: func(obligation *ObligationInput) {
			obligation.Independence = "distinct-actor"
		}, wantError: true},
		{name: "review has registration", obligationType: "review", mutate: func(obligation *ObligationInput) {
			obligation.ReviewProfile = stringPointer("profile")
			obligation.Registration.ID = "REG-test"
			obligation.Registration.Version = 1
		}, wantError: true},
		{name: "review missing profile", obligationType: "review", wantError: true},
		{name: "review has outcome class", obligationType: "review", mutate: func(obligation *ObligationInput) {
			obligation.ReviewProfile = stringPointer("profile")
			obligation.OutcomeClass = stringPointer("outcome")
		}, wantError: true},
		{name: "review evidence type mismatch", obligationType: "review", mutate: func(obligation *ObligationInput) {
			obligation.ReviewProfile = stringPointer("profile")
			obligation.AcceptedEvidenceTypes = []string{"receipt"}
		}, wantError: true},
		{name: "review trust source mismatch", obligationType: "review", mutate: func(obligation *ObligationInput) {
			obligation.ReviewProfile = stringPointer("profile")
			obligation.AcceptedTrustSources = []string{"fid-executed"}
		}, wantError: true},
		{name: "review freshness mismatch", obligationType: "review", mutate: func(obligation *ObligationInput) {
			obligation.ReviewProfile = stringPointer("profile")
			obligation.FreshnessTarget = "none"
		}, wantError: true},
		{name: "review artifact policy mismatch", obligationType: "review", mutate: func(obligation *ObligationInput) {
			obligation.ReviewProfile = stringPointer("profile")
			obligation.ArtifactPolicy = "required-available"
		}, wantError: true},
		{name: "review independence mismatch", obligationType: "review", mutate: func(obligation *ObligationInput) {
			obligation.ReviewProfile = stringPointer("profile")
			obligation.Independence = "human"
		}, wantError: true},
		{name: "human has registration", obligationType: "human-outcome", mutate: func(obligation *ObligationInput) {
			obligation.OutcomeClass = stringPointer("outcome")
			obligation.Registration.ID = "REG-test"
			obligation.Registration.Version = 1
		}, wantError: true},
		{name: "human has review profile", obligationType: "human-outcome", mutate: func(obligation *ObligationInput) {
			obligation.OutcomeClass = stringPointer("outcome")
			obligation.ReviewProfile = stringPointer("profile")
		}, wantError: true},
		{name: "human missing outcome class", obligationType: "human-outcome", wantError: true},
		{name: "human evidence type mismatch", obligationType: "human-outcome", mutate: func(obligation *ObligationInput) {
			obligation.OutcomeClass = stringPointer("outcome")
			obligation.AcceptedEvidenceTypes = []string{"review"}
		}, wantError: true},
		{name: "human trust source mismatch", obligationType: "human-outcome", mutate: func(obligation *ObligationInput) {
			obligation.OutcomeClass = stringPointer("outcome")
			obligation.AcceptedTrustSources = []string{"activated-actor"}
		}, wantError: true},
		{name: "human freshness mismatch", obligationType: "human-outcome", mutate: func(obligation *ObligationInput) {
			obligation.OutcomeClass = stringPointer("outcome")
			obligation.FreshnessTarget = "repository-tree"
		}, wantError: true},
		{name: "human artifact policy mismatch", obligationType: "human-outcome", mutate: func(obligation *ObligationInput) {
			obligation.OutcomeClass = stringPointer("outcome")
			obligation.ArtifactPolicy = "required-available"
		}, wantError: true},
		{name: "human independence mismatch", obligationType: "human-outcome", mutate: func(obligation *ObligationInput) {
			obligation.OutcomeClass = stringPointer("outcome")
			obligation.Independence = "distinct-actor"
		}, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := testEnvelopeCreateRequest()
			obligation := testObligation("OBL-test", test.obligationType, request.PrimaryIssueID)
			if test.mutate != nil {
				test.mutate(&obligation)
			}
			request.Obligations = []ObligationInput{obligation}
			if test.obligationType == "review" && obligation.ReviewProfile == nil && !test.wantError {
				t.Fatal("accepted review case must name review_profile")
			}
			err := validateEnvelopeCreate(request)
			if (err != nil) != test.wantError {
				t.Fatalf("validateEnvelopeCreate() error = %v, wantError=%t", err, test.wantError)
			}
		})
	}
}

func TestEnvelopeAuthoringProtectedDecisionMatrix(t *testing.T) {
	for _, test := range []struct {
		name      string
		decisions func() []ProtectedDecisionInput
		wantError bool
	}{
		{name: "zero", decisions: func() []ProtectedDecisionInput { return nil }},
		{name: "sixteen", decisions: func() []ProtectedDecisionInput { return protectedDecisionsForCount(16) }},
		{name: "seventeen", decisions: func() []ProtectedDecisionInput { return protectedDecisionsForCount(17) }, wantError: true},
		{name: "duplicate class", decisions: func() []ProtectedDecisionInput {
			return []ProtectedDecisionInput{{Ordinal: 0, DecisionClass: "release", RequiredCapabilities: []string{"envelope.complete"}, AffectedActions: []string{"envelope.complete"}}, {Ordinal: 1, DecisionClass: "release", RequiredCapabilities: []string{"envelope.complete"}, AffectedActions: []string{"envelope.complete"}}}
		}, wantError: true},
		{name: "duplicate ordinal", decisions: func() []ProtectedDecisionInput {
			return []ProtectedDecisionInput{{Ordinal: 0, DecisionClass: "release", RequiredCapabilities: []string{"envelope.complete"}, AffectedActions: []string{"envelope.complete"}}, {Ordinal: 0, DecisionClass: "rollback", RequiredCapabilities: []string{"envelope.complete"}, AffectedActions: []string{"envelope.complete"}}}
		}, wantError: true},
		{name: "unguarded action", decisions: func() []ProtectedDecisionInput {
			return []ProtectedDecisionInput{{Ordinal: 0, DecisionClass: "release", RequiredCapabilities: []string{"envelope.complete"}, AffectedActions: []string{"issue.transition-terminal-success"}}}
		}, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := testEnvelopeCreateRequest()
			request.Obligations = []ObligationInput{testObligation("OBL-gate", "gate-evidence", request.PrimaryIssueID)}
			request.ProtectedDecisions = test.decisions()
			err := validateEnvelopeCreate(request)
			if (err != nil) != test.wantError {
				t.Fatalf("validateEnvelopeCreate() error = %v, wantError=%t", err, test.wantError)
			}
		})
	}
}

func TestEvaluateEnvelopeObligationsIncludesQualifyingValidFreshAvailableState(t *testing.T) {
	request := testEnvelopeCreateRequest()
	request.Obligations = []ObligationInput{
		testObligation("OBL-gate", "gate-evidence", request.PrimaryIssueID),
		testObligation("OBL-review", "review", request.PrimaryIssueID),
		testObligation("OBL-human", "human-outcome", request.PrimaryIssueID),
	}
	request.Obligations[1].ReviewProfile = stringPointer("profile")
	request.Obligations[2].OutcomeClass = stringPointer("outcome")
	envelope := envelopeFromCreate("ENV-test", "ACT-test", "2026-09-20T00:00:00Z", request)
	states := map[string][]ObligationEvidenceState{}
	for _, obligation := range envelope.Obligations {
		states[obligation.ID] = []ObligationEvidenceState{{
			ObligationID: obligation.ID, EvidenceID: "EVD-" + obligation.ID, ReceiptID: "RCT-" + obligation.ID,
			ReceiptSuccessful: true, Validity: "valid", Freshness: "fresh", ArtifactAvailability: "available",
		}}
	}

	evaluations := EvaluateEnvelopeObligations(envelope.Obligations, states)
	if len(evaluations) != 3 {
		t.Fatalf("evaluations = %#v, want one qualifying result per Obligation", evaluations)
	}
	for _, evaluation := range evaluations {
		if evaluation.Status != "satisfied" || evaluation.Validity != "valid" || evaluation.Freshness != "fresh" || len(evaluation.Reasons) != 0 {
			t.Fatalf("evaluation = %#v, want satisfied valid/fresh with no reasons", evaluation)
		}
		if evaluation.ObligationID == "OBL-review" || evaluation.ObligationID == "OBL-human" {
			if evaluation.ArtifactAvailability != "not-applicable" {
				t.Fatalf("metadata evaluation = %#v, want not-applicable artifact availability", evaluation)
			}
		} else if evaluation.ArtifactAvailability != "available" {
			t.Fatalf("gate evaluation = %#v, want available artifact availability", evaluation)
		}
	}
}

func envelopeObligationsForCount(issueID string) []ObligationInput {
	obligations := []ObligationInput{testObligation("OBL-gate", "gate-evidence", issueID)}
	for index := 1; index < 17; index++ {
		if index%2 == 1 {
			obligation := testObligation("OBL-review-"+string(rune('a'+index)), "review", issueID)
			obligation.ReviewProfile = stringPointer("profile-" + string(rune('a'+index)))
			obligations = append(obligations, obligation)
			continue
		}
		obligation := testObligation("OBL-human-"+string(rune('a'+index)), "human-outcome", issueID)
		obligation.OutcomeClass = stringPointer("outcome-" + string(rune('a'+index)))
		obligations = append(obligations, obligation)
	}
	return obligations
}

func protectedDecisionsForCount(count int) []ProtectedDecisionInput {
	decisions := make([]ProtectedDecisionInput, 0, count)
	for ordinal := 0; ordinal < count; ordinal++ {
		decisions = append(decisions, ProtectedDecisionInput{
			Ordinal: ordinal, DecisionClass: "decision-" + string(rune('a'+ordinal)),
			RequiredCapabilities: []string{"envelope.complete"}, AffectedActions: []string{"envelope.complete"},
		})
	}
	return decisions
}

func testEnvelopeCreateRequest() EnvelopeCreateRequest {
	repository := Repository{
		CaptureSource: "git",
		HeadOID:       testEnvelopeOID,
		ProjectID:     "PRJ-test",
		TreeOID:       testEnvelopeOID,
		WorktreeID:    "WT-test",
	}
	source := Source{
		Digest:       "sha256:" + testEnvelopeDigest,
		DisplayLabel: "contract",
		Path:         stringPointer("project-docs/contracts.md"),
		Repository:   repository,
		SourceType:   "project-file",
	}
	request := EnvelopeCreateRequest{
		Outcome:         "Deliver the typed Envelope contract.",
		RepositoryScope: repository,
		Sources:         []Source{source},
		PrimaryIssueID:  "ISS-test",
		AuthorityGrants: []AuthorityGrant{{
			Actions: []string{"envelope.complete"},
		}},
	}
	request.RiskResult.Source = source
	request.RiskResult.ResultDigest = "sha256:" + testEnvelopeDigest
	request.AuthorityGrants[0].ActorSelector.ActorID = "ACT-test"
	request.AuthorityGrants[0].ActorSelector.RequiredCapabilities = []string{"envelope.complete"}
	request.AuthorityGrants[0].Scope.IssueID = request.PrimaryIssueID
	return request
}

func testObligation(id, obligationType, issueID string) ObligationInput {
	request := testEnvelopeCreateRequest()
	obligation := ObligationInput{
		ID:             id,
		Type:           obligationType,
		Applicability:  "required",
		GuardedActions: []string{"envelope.complete"},
		Source:         request.Sources[0],
	}
	obligation.Subject.Type = "issue"
	obligation.Subject.ID = issueID
	switch obligationType {
	case "gate-evidence":
		obligation.Registration.ID = "REG-test"
		obligation.Registration.Version = 1
		obligation.AcceptedEvidenceTypes = []string{"receipt"}
		obligation.AcceptedTrustSources = []string{"fid-executed"}
		obligation.FreshnessTarget = "repository-tree"
		obligation.ArtifactPolicy = "required-available"
		obligation.Independence = "fid-executed"
	case "review":
		obligation.AcceptedEvidenceTypes = []string{"review"}
		obligation.AcceptedTrustSources = []string{"activated-actor"}
		obligation.FreshnessTarget = "repository-tree"
		obligation.ArtifactPolicy = "metadata-sufficient"
		obligation.Independence = "distinct-actor"
	case "human-outcome":
		obligation.AcceptedEvidenceTypes = []string{"human-outcome"}
		obligation.AcceptedTrustSources = []string{"activated-human"}
		obligation.FreshnessTarget = "none"
		obligation.ArtifactPolicy = "metadata-sufficient"
		obligation.Independence = "human"
	}
	return obligation
}

func stringPointer(value string) *string { return &value }

const testEnvelopeOID = "0123456789abcdef0123456789abcdef01234567"
const testEnvelopeDigest = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
