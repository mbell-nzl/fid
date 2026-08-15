package links

import (
	"fmt"
	"strings"
)

type UnknownRelationError struct {
	Operation string
	Relation  string
}

func (err *UnknownRelationError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s: unknown relation %q", linkOperationName(err.Operation), err.Relation)
}

type SelfLinkError struct {
	Operation string
	FromID    string
	Relation  string
	ToID      string
}

type DependencyCycleError struct {
	Operation string
	Path      []string
}

type InvalidEndpointTypeError struct {
	Operation string
	Argument  string
	Value     string
}

type ResultLimitExceededError struct {
	Collection string
	Limit      int
	Observed   int
}

func (err *ResultLimitExceededError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s exceeds limit %d with %d items", err.Collection, err.Limit, err.Observed)
}

func (err *InvalidEndpointTypeError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s: %s must be issue or requirement, got %q", linkOperationName(err.Operation), err.Argument, err.Value)
}

func (err *DependencyCycleError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s: dependency cycle: %s", linkOperationName(err.Operation), strings.Join(err.Path, " -> "))
}

func (err *SelfLinkError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s: self link is not allowed: %s %s %s", linkOperationName(err.Operation), err.FromID, err.Relation, err.ToID)
}

type EndpointNotFoundError struct {
	Operation    string
	EndpointType string
	EndpointID   string
	Relation     string
	OtherType    string
	OtherID      string
}

func (err *EndpointNotFoundError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"%s: %s %q does not resolve for relation %q between %s %q and %s %q",
		linkOperationName(err.Operation),
		err.EndpointType,
		err.EndpointID,
		err.Relation,
		err.EndpointType,
		err.EndpointID,
		err.OtherType,
		err.OtherID,
	)
}

type IncompatibleRelationError struct {
	Operation string
	Relation  string
	FromType  string
	FromID    string
	ToType    string
	ToID      string
}

func (err *IncompatibleRelationError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"%s: relation %q is not valid between %s %q and %s %q",
		linkOperationName(err.Operation),
		err.Relation,
		err.FromType,
		err.FromID,
		err.ToType,
		err.ToID,
	)
}

type DuplicateLinkError struct {
	Operation string
	FromID    string
	Relation  string
	ToID      string
}

func (err *DuplicateLinkError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s: link already exists: %s %s %s", linkOperationName(err.Operation), err.FromID, err.Relation, err.ToID)
}

type MissingLinkError struct {
	Operation string
	FromType  string
	FromID    string
	Relation  string
	ToType    string
	ToID      string
}

func (err *MissingLinkError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s: link does not exist: %s %s %s", linkOperationName(err.Operation), err.FromID, err.Relation, err.ToID)
}

func linkOperationName(operation string) string {
	if operation == "" {
		return "link add"
	}
	return operation
}
