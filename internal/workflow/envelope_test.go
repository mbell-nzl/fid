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
