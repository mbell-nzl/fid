package main

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/links"
	"github.com/spf13/cobra"
)

func renderLinkError(cmd *cobra.Command, operation string, cause error) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := linkFailure(cmd, operation, cause)
	if err != nil {
		return err
	}
	if !jsonMode {
		return writeHumanFailure(cmd, operation, failure, cause)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
}

func linkFailure(cmd *cobra.Command, operation string, cause error) (cli.Failure, error) {
	if failure, err := issueMutationFailure(cmd, operation, "", cause); err != nil {
		return cli.Failure{}, err
	} else if failure.Code != cli.CodeInvalidArgument {
		return failure, nil
	}

	var invalidType *links.InvalidEndpointTypeError
	if errors.As(cause, &invalidType) {
		return cli.NewFailure(cli.CodeInvalidArgument, "invalid link endpoint type", nil, map[string]any{
			"argument": invalidType.Argument,
			"value":    invalidType.Value,
			"allowed":  []string{"issue"},
		})
	}
	var unknown *links.UnknownRelationError
	if errors.As(cause, &unknown) {
		return cli.NewFailure(cli.CodeUnsupportedRelation, "link relation is unsupported", nil, map[string]any{
			"argument": "relation",
			"value":    unknown.Relation,
			"allowed":  links.SupportedRelations(),
		})
	}
	var endpoint *links.EndpointNotFoundError
	if errors.As(cause, &endpoint) {
		return cli.NewFailure(cli.CodeIssueNotFound, "issue was not found", &cli.Subject{Type: cli.SubjectIssue, ID: endpoint.EndpointID}, map[string]any{"id": endpoint.EndpointID})
	}
	var duplicate *links.DuplicateLinkError
	if errors.As(cause, &duplicate) {
		return cli.NewFailure(cli.CodeDuplicateEntity, "entity already exists", &cli.Subject{Type: cli.SubjectLink, ID: duplicate.FromID}, map[string]any{
			"entity_type": "link",
			"identity":    map[string]any{"fields": []string{"from_id", "relation", "to_id"}, "values": []string{duplicate.FromID, duplicate.Relation, duplicate.ToID}},
		})
	}
	var missing *links.MissingLinkError
	if errors.As(cause, &missing) {
		return cli.NewFailure(cli.CodeLinkNotFound, "link was not found", &cli.Subject{Type: cli.SubjectLink, ID: missing.FromID}, map[string]any{
			"from_type": missing.FromType, "from_id": missing.FromID, "relation": missing.Relation, "to_type": missing.ToType, "to_id": missing.ToID,
		})
	}
	var cycle *links.DependencyCycleError
	if errors.As(cause, &cycle) {
		return cli.NewFailure(cli.CodeDependencyCycle, "link would create a dependency cycle", nil, map[string]any{"path": cycle.Path})
	}
	var limit *links.ResultLimitExceededError
	if errors.As(cause, &limit) {
		return cli.NewFailure(cli.CodeResultLimitExceeded, "read result exceeds its collection limit", nil, map[string]any{"collection": limit.Collection, "limit": limit.Limit, "observed": limit.Observed})
	}
	var self *links.SelfLinkError
	if errors.As(cause, &self) {
		return cli.NewFailure(cli.CodeInvalidArgument, "self link is not permitted", nil, map[string]any{
			"argument": "link",
			"value":    self.FromID,
			"allowed":  []string{"distinct endpoint ids"},
		})
	}
	return genericFailure(cause, cli.CodeInvalidArgument, "invalid link request", nil, map[string]any{
		"argument": "request", "value": nil, "allowed": []string{"accepted " + operation + " input"},
	})
}
